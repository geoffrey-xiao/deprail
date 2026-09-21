package app

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

func CurrentRepositoryRevision(ctx context.Context, root string) (string, error) {
	git, err := exec.LookPath("git")
	if err != nil {
		return "", err
	}
	command := exec.CommandContext(ctx, git, "-C", root, "rev-parse", "HEAD")
	output, err := command.Output()
	if err != nil {
		return "", err
	}
	revision := strings.TrimSpace(string(output))
	if revision == "" {
		return "", fmt.Errorf("git revision is empty")
	}
	return revision, nil
}

func CurrentRepositoryState(root string) (string, error) {
	return CurrentRepositoryStateExcluding(root)
}

// CurrentRepositoryStateExcluding computes the repository digest while ignoring
// generated files that were not present when a scan was taken.
func CurrentRepositoryStateExcluding(root string, excluded ...string) (string, error) {
	canonical, err := filepath.EvalSymlinks(root)
	if err != nil {
		return "", fmt.Errorf("resolve repository root: %w", err)
	}
	excludedPaths := make(map[string]struct{}, len(excluded))
	for _, path := range excluded {
		if path == "" {
			continue
		}
		absolute, err := filepath.Abs(path)
		if err != nil {
			return "", fmt.Errorf("resolve excluded path: %w", err)
		}
		resolved := absolute
		if evaluated, evalErr := filepath.EvalSymlinks(absolute); evalErr == nil {
			resolved = evaluated
		} else if parent, parentErr := filepath.EvalSymlinks(filepath.Dir(absolute)); parentErr == nil {
			resolved = filepath.Join(parent, filepath.Base(absolute))
		}
		relative, err := filepath.Rel(canonical, resolved)
		if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(os.PathSeparator)) {
			continue
		}
		excludedPaths[filepath.Clean(relative)] = struct{}{}
	}
	type entry struct{ path, digest string }
	entries := []entry{}
	err = filepath.WalkDir(canonical, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, err := filepath.Rel(canonical, path)
		if err != nil {
			return err
		}
		if rel == ".deprail" || strings.HasPrefix(rel, ".deprail"+string(os.PathSeparator)) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if d.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("repository contains symlink: %s", filepath.ToSlash(rel))
		}
		if _, skip := excludedPaths[filepath.Clean(rel)]; skip {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		if !d.Type().IsRegular() {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		digest := sha256.Sum256(data)
		entries = append(entries, entry{filepath.ToSlash(rel), hex.EncodeToString(digest[:])})
		return nil
	})
	if err != nil {
		return "", err
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].path < entries[j].path })
	h := sha256.New()
	for _, item := range entries {
		fmt.Fprintf(h, "%s\x00%s\x00", item.path, item.digest)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
