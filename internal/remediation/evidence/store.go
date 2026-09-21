package evidence

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Store struct {
	Root string
}

func (s Store) Save(record Record) (string, string, error) {
	data, err := record.MarshalJSON()
	if err != nil {
		return "", "", err
	}
	digest, err := Digest(record)
	if err != nil {
		return "", "", err
	}
	if s.Root == "" {
		return "", "", errors.New("evidence root is required")
	}
	root, err := secureRoot(s.Root)
	if err != nil {
		return "", "", err
	}
	dir := filepath.Join(root, digest[:2])
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", "", fmt.Errorf("create evidence directory: %w", err)
	}
	if err := os.Chmod(dir, 0o700); err != nil {
		return "", "", fmt.Errorf("restrict evidence directory: %w", err)
	}
	if err := validatePrivateDirectory(dir); err != nil {
		return "", "", err
	}
	path := filepath.Join(dir, digest+".json")
	if existing, err := os.ReadFile(path); err == nil {
		if string(existing) != string(data) {
			return "", "", errors.New("evidence digest collision")
		}
		return path, digest, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return "", "", fmt.Errorf("inspect evidence: %w", err)
	}
	tmp, err := os.CreateTemp(dir, ".evidence-*")
	if err != nil {
		return "", "", fmt.Errorf("create evidence temporary file: %w", err)
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)
	if err := tmp.Chmod(0o600); err != nil {
		_ = tmp.Close()
		return "", "", fmt.Errorf("restrict evidence temporary file: %w", err)
	}
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return "", "", fmt.Errorf("write evidence: %w", err)
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return "", "", fmt.Errorf("sync evidence: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return "", "", fmt.Errorf("close evidence: %w", err)
	}
	if err := os.Link(tmpName, path); err != nil {
		if errors.Is(err, os.ErrExist) {
			return s.Save(record)
		}
		return "", "", fmt.Errorf("publish evidence: %w", err)
	}
	return path, digest, nil
}

func secureRoot(root string) (string, error) {
	absolute, err := filepath.Abs(root)
	if err != nil {
		return "", fmt.Errorf("resolve evidence root: %w", err)
	}
	if err := os.MkdirAll(absolute, 0o700); err != nil {
		return "", fmt.Errorf("create evidence root: %w", err)
	}
	if err := os.Chmod(absolute, 0o700); err != nil {
		return "", fmt.Errorf("restrict evidence root: %w", err)
	}
	info, err := os.Lstat(absolute)
	if err != nil {
		return "", fmt.Errorf("inspect evidence root: %w", err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return "", errors.New("evidence root must not be a symlink")
	}
	resolved, err := filepath.EvalSymlinks(absolute)
	if err != nil {
		return "", fmt.Errorf("resolve evidence root: %w", err)
	}
	if err := validatePrivateDirectory(resolved); err != nil {
		return "", err
	}
	return resolved, nil
}

func validatePrivateDirectory(path string) error {
	info, err := os.Lstat(path)
	if err != nil {
		return fmt.Errorf("inspect evidence directory: %w", err)
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.IsDir() {
		return errors.New("evidence write boundary must be a private directory")
	}
	if info.Mode().Perm()&0o077 != 0 {
		return errors.New("evidence directory permissions are too permissive")
	}
	return nil
}
