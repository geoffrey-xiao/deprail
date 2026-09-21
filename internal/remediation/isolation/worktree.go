package isolation

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Workspace struct {
	Root      string
	Path      string
	SourceRef string
	Created   bool
}

func Create(ctx context.Context, root, sourceRef string) (Workspace, error) {
	canonical, err := canonicalDirectory(root)
	if err != nil {
		return Workspace{}, err
	}
	if sourceRef == "" || strings.ContainsAny(sourceRef, "\r\n") {
		return Workspace{}, errors.New("source ref is required")
	}
	parent, err := os.MkdirTemp("", "deprail-worktree-")
	if err != nil {
		return Workspace{}, fmt.Errorf("create workspace parent: %w", err)
	}
	path := filepath.Join(parent, "worktree")
	cmd := exec.CommandContext(ctx, "git", "-C", canonical, "worktree", "add", "--detach", path, sourceRef)
	if output, err := cmd.CombinedOutput(); err != nil {
		_ = os.RemoveAll(parent)
		return Workspace{}, fmt.Errorf("create isolated worktree: %w: %s", err, strings.TrimSpace(string(output)))
	}
	return Workspace{Root: canonical, Path: path, SourceRef: sourceRef, Created: true}, nil
}

func (w Workspace) Remove(ctx context.Context) error {
	if !w.Created || w.Path == "" {
		return nil
	}
	cmd := exec.CommandContext(ctx, "git", "-C", w.Root, "worktree", "remove", "--force", w.Path)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("remove isolated worktree: %w: %s", err, strings.TrimSpace(string(output)))
	}
	return os.RemoveAll(filepath.Dir(w.Path))
}

func canonicalDirectory(root string) (string, error) {
	if root == "" {
		return "", errors.New("repository root is required")
	}
	absolute, err := filepath.Abs(root)
	if err != nil {
		return "", fmt.Errorf("resolve repository root: %w", err)
	}
	canonical, err := filepath.EvalSymlinks(absolute)
	if err != nil {
		return "", fmt.Errorf("resolve repository root: %w", err)
	}
	info, err := os.Stat(canonical)
	if err != nil {
		return "", fmt.Errorf("stat repository root: %w", err)
	}
	if !info.IsDir() {
		return "", errors.New("repository root must be a directory")
	}
	return canonical, nil
}
