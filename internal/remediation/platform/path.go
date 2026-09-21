package platform

import (
	"errors"
	"path/filepath"
	"strings"
)

func RelativePath(path string) (string, error) {
	if path == "" || filepath.IsAbs(path) {
		return "", errors.New("path must be relative")
	}
	portable := strings.ReplaceAll(path, "\\", "/")
	if len(portable) >= 2 && portable[1] == ':' {
		return "", errors.New("path must not contain a volume")
	}
	clean := filepath.Clean(filepath.FromSlash(portable))
	if clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", errors.New("path escapes workspace")
	}
	return filepath.ToSlash(clean), nil
}
