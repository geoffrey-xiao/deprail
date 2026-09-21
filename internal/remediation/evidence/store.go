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
	dir := filepath.Join(s.Root, digest[:2])
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", "", fmt.Errorf("create evidence directory: %w", err)
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
