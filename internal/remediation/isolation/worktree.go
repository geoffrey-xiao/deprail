package isolation

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/geoffrey-xiao/deprail/internal/process"
)

const (
	gitTimeout   = 30 * time.Second
	gitOutputCap = 1 << 20
)

type Workspace struct {
	Root      string
	Path      string
	SourceRef string
	Created   bool
	verified  bool
}

func (w Workspace) Prepared() bool { return w.Created && w.verified && w.Path != "" }

// Subdirectory returns a verified workspace rooted at a repository-relative
// subdirectory. It preserves the isolation guarantees of the parent workspace.
func (w Workspace) Subdirectory(relative string) (Workspace, error) {
	if !w.Prepared() {
		return Workspace{}, errors.New("verified isolated workspace is required")
	}
	if relative == "" {
		relative = "."
	}
	if filepath.IsAbs(relative) {
		return Workspace{}, errors.New("workspace subdirectory must be relative")
	}
	path := filepath.Clean(filepath.Join(w.Path, filepath.FromSlash(relative)))
	rel, err := filepath.Rel(w.Path, path)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return Workspace{}, errors.New("workspace subdirectory escapes isolated workspace")
	}
	return Workspace{Root: w.Root, Path: path, SourceRef: w.SourceRef, Created: true, verified: true}, nil
}

func Create(ctx context.Context, root, sourceRef string) (Workspace, error) {
	canonical, err := canonicalDirectory(ctx, root)
	if err != nil {
		return Workspace{}, err
	}
	if sourceRef == "" || strings.ContainsAny(sourceRef, "\r\n") || strings.HasPrefix(sourceRef, "-") {
		return Workspace{}, errors.New("source ref is invalid")
	}
	resolved, err := runGit(ctx, canonical, "rev-parse", "--verify", sourceRef+"^{commit}")
	if err != nil {
		return Workspace{}, fmt.Errorf("resolve source ref: %w", err)
	}
	commit := strings.TrimSpace(string(resolved.Stdout))
	if commit == "" || strings.ContainsAny(commit, "\r\n") {
		return Workspace{}, errors.New("source ref did not resolve to a commit")
	}
	parent, err := os.MkdirTemp("", "deprail-worktree-")
	if err != nil {
		return Workspace{}, fmt.Errorf("create workspace parent: %w", err)
	}
	if err := os.Chmod(parent, 0o700); err != nil {
		_ = os.RemoveAll(parent)
		return Workspace{}, fmt.Errorf("restrict workspace parent: %w", err)
	}
	path := filepath.Join(parent, "worktree")
	if _, err := runGit(ctx, canonical, "worktree", "add", "--detach", path, commit); err != nil {
		_ = os.RemoveAll(parent)
		return Workspace{}, fmt.Errorf("create isolated worktree: %w", err)
	}
	return Workspace{Root: canonical, Path: path, SourceRef: commit, Created: true, verified: true}, nil
}

type CleanupResult struct {
	GitRemoved        bool
	FilesystemRemoved bool
	Retryable         bool
	Err               error
}

func (r CleanupResult) Succeeded() bool { return r.Err == nil }

func (w *Workspace) Remove(ctx context.Context) error {
	return w.RemoveWithEvidence(ctx).Err
}

func (w *Workspace) RemoveWithEvidence(ctx context.Context) CleanupResult {
	if w == nil || (!w.Created && w.Path == "") {
		return CleanupResult{}
	}
	result := CleanupResult{Retryable: true}
	var gitErr error
	if w.Created {
		if _, err := runGit(ctx, w.Root, "worktree", "remove", "--force", w.Path); err != nil {
			gitErr = fmt.Errorf("remove isolated worktree: %w", err)
		} else {
			result.GitRemoved = true
			w.Created = false
			w.verified = false
		}
	}
	if err := os.RemoveAll(filepath.Dir(w.Path)); err != nil {
		result.Err = fmt.Errorf("remove workspace parent: %w", err)
	} else {
		result.FilesystemRemoved = true
	}
	if gitErr != nil && result.Err != nil {
		result.Err = fmt.Errorf("%v; %w", gitErr, result.Err)
	} else if gitErr != nil {
		result.Err = gitErr
	}
	if result.Err == nil {
		result.Retryable = false
		w.Path = ""
	}
	return result
}

func canonicalDirectory(ctx context.Context, root string) (string, error) {
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
	result, err := runGit(ctx, canonical, "rev-parse", "--show-toplevel")
	if err != nil {
		return "", fmt.Errorf("resolve Git repository root: %w", err)
	}
	gitRoot, err := filepath.EvalSymlinks(strings.TrimSpace(string(result.Stdout)))
	if err != nil || gitRoot != canonical {
		return "", errors.New("repository root must be the canonical Git top-level")
	}
	return canonical, nil
}

func CanonicalRepositoryRoot(ctx context.Context, root string) (string, error) {
	return canonicalDirectory(ctx, root)
}

func runGit(ctx context.Context, dir string, args ...string) (process.Result, error) {
	git, err := exec.LookPath("git")
	if err != nil {
		return process.Result{}, err
	}
	return process.Run(ctx, process.Request{Path: git, Args: append([]string{"-C", dir}, args...), Dir: dir, Timeout: gitTimeout, OutputCap: gitOutputCap})
}
