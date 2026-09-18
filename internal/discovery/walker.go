package discovery

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var defaultIgnoredDirectories = map[string]struct{}{
	".deprail":     {},
	".git":         {},
	".hg":          {},
	".svn":         {},
	"node_modules": {},
	"target":       {},
	".venv":        {},
}

type WalkOptions struct {
	IgnoreDirectories []string
	NoIgnore          bool
}

type WalkResult struct {
	Paths       []string
	Diagnostics []Diagnostic
}

func Walk(ctx context.Context, root string, options WalkOptions) (WalkResult, error) {
	canonicalRoot, err := canonicalRoot(root)
	if err != nil {
		return WalkResult{}, err
	}
	ignored := ignoredDirectories(options)
	result := WalkResult{}
	if err := walkDirectory(ctx, canonicalRoot, canonicalRoot, ignored, &result); err != nil {
		return WalkResult{}, err
	}
	sort.Strings(result.Paths)
	sortDiagnostics(result.Diagnostics)
	return result, nil
}

func canonicalRoot(root string) (string, error) {
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
		return "", errors.New("repository root is not a directory")
	}
	return canonical, nil
}

func walkDirectory(ctx context.Context, root, directory string, ignored map[string]struct{}, result *WalkResult) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	entries, err := os.ReadDir(directory)
	if err != nil {
		return fmt.Errorf("read directory %s: %w", directory, err)
	}
	for _, entry := range entries {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		fullPath := filepath.Join(directory, entry.Name())
		info, err := os.Lstat(fullPath)
		if err != nil {
			result.Diagnostics = append(result.Diagnostics, Diagnostic{Code: "WALK_ENTRY_FAILED", Message: err.Error(), Scope: relativePath(root, fullPath)})
			continue
		}
		if info.Mode()&os.ModeSymlink != 0 {
			handleSymlink(root, fullPath, result)
			continue
		}
		if info.IsDir() {
			if _, skip := ignored[entry.Name()]; skip {
				continue
			}
			if err := walkDirectory(ctx, root, fullPath, ignored, result); err != nil {
				return err
			}
			continue
		}
		if info.Mode().IsRegular() {
			result.Paths = append(result.Paths, relativePath(root, fullPath))
		}
	}
	return nil
}

func ignoredDirectories(options WalkOptions) map[string]struct{} {
	if options.NoIgnore {
		return directorySet(options.IgnoreDirectories)
	}
	ignored := directorySet(options.IgnoreDirectories)
	for name := range defaultIgnoredDirectories {
		ignored[name] = struct{}{}
	}
	return ignored
}

func directorySet(names []string) map[string]struct{} {
	ignored := make(map[string]struct{}, len(names))
	for _, name := range names {
		if name != "" {
			ignored[name] = struct{}{}
		}
	}
	return ignored
}

func handleSymlink(root, path string, result *WalkResult) {
	target, err := filepath.EvalSymlinks(path)
	if err != nil {
		result.Diagnostics = append(result.Diagnostics, Diagnostic{Code: "PATH_OUTSIDE_ROOT", Message: err.Error(), Scope: relativePath(root, path)})
		return
	}
	if !isWithinRoot(root, target) {
		result.Diagnostics = append(result.Diagnostics, Diagnostic{Code: "PATH_OUTSIDE_ROOT", Message: "symlink target escapes repository root", Scope: relativePath(root, path)})
		return
	}
	result.Diagnostics = append(result.Diagnostics, Diagnostic{Code: "SYMLINK_SKIPPED", Message: "symlink is not traversed", Scope: relativePath(root, path)})
}

func isWithinRoot(root, path string) bool {
	relative, err := filepath.Rel(root, path)
	return err == nil && relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator)) && !filepath.IsAbs(relative)
}

func relativePath(root, path string) string {
	relative, err := filepath.Rel(root, path)
	if err != nil {
		return filepath.ToSlash(path)
	}
	return filepath.ToSlash(relative)
}
