package history

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type resolvedPaths struct {
	database string
}

func preparePaths(configDir, repositoryRoot string, create bool) (resolvedPaths, bool, error) {
	if configDir == "" {
		return resolvedPaths{}, false, errors.New("user config directory is unavailable")
	}
	absoluteConfig, err := filepath.Abs(configDir)
	if err != nil {
		return resolvedPaths{}, false, err
	}
	absoluteConfig = filepath.Clean(absoluteConfig)
	dataDir := filepath.Join(absoluteConfig, ".deprail")
	canonicalRepository := ""
	if repositoryRoot != "" {
		canonicalRepository, err = canonicalDirectory(repositoryRoot)
		if err != nil {
			return resolvedPaths{}, false, errors.New("repository boundary is not canonical")
		}
		if pathInside(canonicalRepository, dataDir) {
			return resolvedPaths{}, false, errors.New("history directory is inside the scanned repository")
		}
	}
	missing, err := checkDirectoryChain(absoluteConfig, create)
	if err != nil {
		return resolvedPaths{}, false, err
	}
	if missing {
		return resolvedPaths{database: filepath.Join(dataDir, "history.sqlite3")}, true, nil
	}
	if create {
		if err := ensurePrivateDirectory(dataDir); err != nil {
			return resolvedPaths{}, false, err
		}
	} else {
		missing, err = checkDirectoryChain(dataDir, false)
		if err != nil {
			return resolvedPaths{}, false, err
		}
		if missing {
			return resolvedPaths{database: filepath.Join(dataDir, "history.sqlite3")}, true, nil
		}
		if err := secureDirectory(dataDir, false, false); err != nil {
			return resolvedPaths{}, false, err
		}
	}
	if canonicalRepository != "" && pathInside(canonicalRepository, dataDir) {
		return resolvedPaths{}, false, errors.New("history directory is inside the scanned repository")
	}
	resolvedData, err := canonicalDirectory(dataDir)
	if err != nil {
		return resolvedPaths{}, false, err
	}
	return resolvedPaths{database: filepath.Join(resolvedData, "history.sqlite3")}, false, nil
}

func canonicalDirectory(path string) (string, error) {
	absolute, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	absolute = filepath.Clean(absolute)
	missing, err := checkDirectoryChain(absolute, false)
	if err != nil || missing {
		if err == nil {
			err = os.ErrNotExist
		}
		return "", err
	}
	resolved, err := filepath.EvalSymlinks(absolute)
	if err != nil {
		return "", err
	}
	info, err := os.Stat(resolved)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", errors.New("path is not a directory")
	}
	return filepath.Clean(resolved), nil
}

func checkDirectoryChain(path string, create bool) (bool, error) {
	volume := filepath.VolumeName(path)
	rest := strings.TrimPrefix(path, volume)
	root := volume + string(os.PathSeparator)
	if volume == "" {
		root = string(os.PathSeparator)
	}
	current := filepath.Clean(root)
	info, err := os.Lstat(current)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return true, nil
		}
		return false, err
	}
	if !info.IsDir() {
		return false, errors.New("path root is not a directory")
	}
	if unsafe, err := isUnsafePathLink(current, info); err != nil || unsafe {
		if err == nil {
			err = errors.New("path contains a symlink or reparse point")
		}
		return false, err
	}
	separator := string(os.PathSeparator)
	for _, part := range strings.Split(strings.Trim(rest, separator), separator) {
		if part == "" || part == "." {
			continue
		}
		current = filepath.Join(current, part)
		info, err = os.Lstat(current)
		if errors.Is(err, os.ErrNotExist) {
			if !create {
				return true, nil
			}
			if err := os.Mkdir(current, 0o700); err != nil && !errors.Is(err, os.ErrExist) {
				return false, err
			}
			info, err = os.Lstat(current)
		}
		if err != nil {
			return false, err
		}
		if !info.IsDir() {
			return false, errors.New("history path component is not a directory")
		}
		if unsafe, err := isUnsafePathLink(current, info); err != nil || unsafe {
			if err == nil {
				err = errors.New("path contains a symlink or reparse point")
			}
			return false, err
		}
	}
	return false, nil
}

func ensurePrivateDirectory(path string) error {
	_, statErr := os.Lstat(path)
	created := errors.Is(statErr, os.ErrNotExist)
	if statErr != nil && !created {
		return statErr
	}
	missing, err := checkDirectoryChain(path, true)
	if err != nil || missing {
		if err == nil {
			err = os.ErrNotExist
		}
		return err
	}
	return secureDirectory(path, true, created)
}

func pathInside(root, candidate string) bool {
	relative, err := filepath.Rel(root, candidate)
	if err != nil {
		return false
	}
	if runtime.GOOS == "windows" {
		relative = strings.ToLower(relative)
	}
	return relative == "." || relative != ".." && !strings.HasPrefix(relative, ".."+string(os.PathSeparator)) && !filepath.IsAbs(relative)
}
func rejectOrphanSidecars(database string) error {
	for _, suffix := range []string{"-wal", "-shm"} {
		if _, err := os.Lstat(database + suffix); err == nil {
			return errors.New("database sidecar exists without its main database")
		} else if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}
	return nil
}
