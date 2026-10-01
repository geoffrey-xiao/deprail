//go:build !windows

package history

import (
	"errors"
	"os"
	"syscall"
)

func isUnsafePathLink(_ string, info os.FileInfo) (bool, error) {
	return info.Mode()&os.ModeSymlink != 0, nil
}

func secureDirectory(path string, writable, _ bool) error {
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.IsDir() {
		return errors.New("history directory is not a real directory")
	}
	if !ownedByCurrentUser(info) {
		return errors.New("history directory is not owned by the current user")
	}
	if writable {
		if err := os.Chmod(path, 0o700); err != nil {
			return err
		}
		info, err = os.Lstat(path)
		if err != nil {
			return err
		}
	}
	if info.Mode().Perm() != 0o700 {
		return errors.New("history directory permissions are not private")
	}
	return nil
}

func secureFile(path string, writable, _ bool) error {
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.Mode().IsRegular() {
		return errors.New("history database is not a regular file")
	}
	if !ownedByCurrentUser(info) {
		return errors.New("history database is not owned by the current user")
	}
	if writable {
		if err := os.Chmod(path, 0o600); err != nil {
			return err
		}
		info, err = os.Lstat(path)
		if err != nil {
			return err
		}
	}
	if info.Mode().Perm() != 0o600 {
		return errors.New("history database permissions are not private")
	}
	return nil
}

func secureSQLiteSidecars(database string, _ bool) error {
	for _, suffix := range []string{"-wal", "-shm"} {
		path := database + suffix
		if err := secureFile(path, true, false); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}
	return nil
}

func ownedByCurrentUser(info os.FileInfo) bool {
	stat, ok := info.Sys().(*syscall.Stat_t)
	return ok && stat.Uid == uint32(os.Geteuid())
}
