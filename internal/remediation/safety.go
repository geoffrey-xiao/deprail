package remediation

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var (
	ErrUnsafePath       = errors.New("unsafe repository path")
	ErrOutputExists     = errors.New("external output already exists")
	ErrRepositoryChange = errors.New("repository changed during planning")
)

// CanonicalRepositoryRoot resolves root without following paths supplied by callers later.
func CanonicalRepositoryRoot(root string) (string, error) {
	if root == "" {
		return "", fmt.Errorf("repository root is required: %w", ErrUnsafePath)
	}
	absolute, err := filepath.Abs(root)
	if err != nil {
		return "", fmt.Errorf("resolve repository root: %w", err)
	}
	resolved, err := filepath.EvalSymlinks(absolute)
	if err != nil {
		return "", fmt.Errorf("resolve repository root: %w", err)
	}
	info, err := os.Stat(resolved)
	if err != nil || !info.IsDir() {
		return "", fmt.Errorf("repository root is not a directory: %w", ErrUnsafePath)
	}
	return filepath.Clean(resolved), nil
}

// ValidateExternalOutput rejects output paths that resolve inside the repository.
// The returned path is absolute and its existing parent is symlink-resolved.
func ValidateExternalOutput(root, output string) (string, error) {
	canonicalRoot, err := CanonicalRepositoryRoot(root)
	if err != nil {
		return "", err
	}
	absoluteRoot, err := filepath.Abs(root)
	if err != nil {
		return "", fmt.Errorf("resolve repository root: %w", err)
	}
	if output == "" {
		return "", fmt.Errorf("output path is required: %w", ErrUnsafePath)
	}
	if hasTraversalComponent(output) {
		return "", fmt.Errorf("output path contains traversal: %w", ErrUnsafePath)
	}
	absolute, err := filepath.Abs(output)
	if err != nil {
		return "", fmt.Errorf("resolve output path: %w", err)
	}
	clean := filepath.Clean(absolute)
	if pathWithin(absoluteRoot, clean) || pathWithin(canonicalRoot, clean) {
		return "", fmt.Errorf("output path is inside repository root: %w", ErrUnsafePath)
	}
	parent, err := filepath.EvalSymlinks(filepath.Dir(clean))
	if err != nil {
		return "", fmt.Errorf("resolve output parent: %w", err)
	}
	resolved := filepath.Join(parent, filepath.Base(clean))
	if pathWithin(canonicalRoot, resolved) || pathWithin(absoluteRoot, resolved) {
		return "", fmt.Errorf("output path is inside repository root: %w", ErrUnsafePath)
	}
	return resolved, nil
}

func pathWithin(root, path string) bool {
	rel, err := filepath.Rel(root, path)
	return err == nil && rel != ".." && !hasParentPrefix(rel)
}
func hasTraversalComponent(path string) bool {
	for _, component := range strings.FieldsFunc(path, func(r rune) bool {
		return r == '/' || r == '\\'
	}) {
		if component == ".." {
			return true
		}
	}
	return false
}

func hasParentPrefix(path string) bool {
	return path == ".." || len(path) > 3 && path[:3] == ".."+string(filepath.Separator)
}

// WriteExternalOutput atomically publishes restrictive output and never overwrites.
func WriteExternalOutput(root, output string, data []byte) (string, error) {
	path, err := ValidateExternalOutput(root, output)
	if err != nil {
		return "", err
	}
	if err := publishExternalOutput(path, data); err != nil {
		return "", err
	}
	return path, nil
}

type RepositorySnapshot struct {
	Entries map[string]string
}

// SnapshotRepository records regular files, directories, and symlink targets without following symlinks.
func SnapshotRepository(root string) (RepositorySnapshot, error) {
	canonicalRoot, err := CanonicalRepositoryRoot(root)
	if err != nil {
		return RepositorySnapshot{}, err
	}
	entries := make(map[string]string)
	err = filepath.WalkDir(canonicalRoot, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, err := filepath.Rel(canonicalRoot, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		key := filepath.ToSlash(rel)
		switch info.Mode() & os.ModeType {
		case os.ModeSymlink:
			target, err := os.Readlink(path)
			if err != nil {
				return err
			}
			entries[key] = "symlink:" + target
		case os.ModeDir:
			entries[key] = "directory:" + info.Mode().String()
		case 0:
			digest, err := fileDigest(path)
			if err != nil {
				return err
			}
			entries[key] = "file:" + info.Mode().String() + ":" + digest
		default:
			entries[key] = "special:" + info.Mode().String()
		}
		return nil
	})
	if err != nil {
		return RepositorySnapshot{}, fmt.Errorf("snapshot repository: %w", err)
	}
	return RepositorySnapshot{Entries: entries}, nil
}

func (s RepositorySnapshot) Equal(other RepositorySnapshot) bool {
	if len(s.Entries) != len(other.Entries) {
		return false
	}
	for key, value := range s.Entries {
		if other.Entries[key] != value {
			return false
		}
	}
	return true
}
func EnsureRepositoryUnchanged(before, after RepositorySnapshot) error {
	if !before.Equal(after) {
		return ErrRepositoryChange
	}
	return nil
}

func fileDigest(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func (s RepositorySnapshot) Paths() []string {
	paths := make([]string, 0, len(s.Entries))
	for path := range s.Entries {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	return paths
}
