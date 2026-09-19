//go:build !windows

package remediation

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"
)

func publishExternalOutput(path string, data []byte) error {
	parent := filepath.Dir(path)
	base := filepath.Base(path)
	dirfd, err := openDirectoryNoFollow(parent)
	if err != nil {
		return fmt.Errorf("open output directory securely: %w", err)
	}
	defer unix.Close(dirfd)

	tempName, err := secureTempName()
	if err != nil {
		return fmt.Errorf("create output temporary name: %w", err)
	}
	fd, err := unix.Openat(dirfd, tempName, unix.O_WRONLY|unix.O_CREAT|unix.O_EXCL|unix.O_NOFOLLOW|unix.O_CLOEXEC, 0o600)
	if err != nil {
		return fmt.Errorf("create output temporary file: %w", err)
	}
	temp := os.NewFile(uintptr(fd), tempName)
	removeTemp := true
	defer func() {
		if removeTemp {
			_ = unix.Unlinkat(dirfd, tempName, 0)
		}
	}()
	if _, err := temp.Write(data); err != nil {
		temp.Close()
		return fmt.Errorf("write external output: %w", err)
	}
	if err := temp.Sync(); err != nil {
		temp.Close()
		return fmt.Errorf("sync external output: %w", err)
	}
	if err := temp.Close(); err != nil {
		return fmt.Errorf("close external output: %w", err)
	}
	if err := unix.Linkat(dirfd, tempName, dirfd, base, 0); err != nil {
		if err == unix.EEXIST {
			return ErrOutputExists
		}
		return fmt.Errorf("publish external output without overwrite: %w", err)
	}
	removeTemp = false
	if err := unix.Unlinkat(dirfd, tempName, 0); err != nil {
		return fmt.Errorf("remove output temporary file: %w", err)
	}
	return nil
}

func openDirectoryNoFollow(path string) (int, error) {
	volume := filepath.VolumeName(path)
	clean := filepath.Clean(path)
	if volume != "" {
		return 0, fmt.Errorf("unsupported output volume %q", volume)
	}
	fd, err := unix.Open(string(filepath.Separator), unix.O_RDONLY|unix.O_DIRECTORY|unix.O_NOFOLLOW|unix.O_CLOEXEC, 0)
	if err != nil {
		return 0, err
	}
	for _, component := range splitAbsolutePath(clean) {
		next, err := unix.Openat(fd, component, unix.O_RDONLY|unix.O_DIRECTORY|unix.O_NOFOLLOW|unix.O_CLOEXEC, 0)
		if err != nil {
			unix.Close(fd)
			return 0, err
		}
		unix.Close(fd)
		fd = next
	}
	return fd, nil
}

func splitAbsolutePath(path string) []string {
	parts := make([]string, 0)
	start := 0
	for start < len(path) {
		for start < len(path) && path[start] == filepath.Separator {
			start++
		}
		end := start
		for end < len(path) && path[end] != filepath.Separator {
			end++
		}
		if start < end {
			parts = append(parts, path[start:end])
		}
		start = end
	}
	return parts
}

func secureTempName() (string, error) {
	var random [16]byte
	if _, err := rand.Read(random[:]); err != nil {
		return "", err
	}
	return ".deprail-output-" + hex.EncodeToString(random[:]), nil
}
