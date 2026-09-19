//go:build windows

package remediation

import (
	"fmt"
	"os"
	"path/filepath"
)

func publishExternalOutput(path string, data []byte) error {
	parent := filepath.Dir(path)
	if _, err := os.Stat(parent); err != nil {
		return fmt.Errorf("output parent is unavailable: %w", err)
	}
	temp, err := os.CreateTemp(parent, ".deprail-output-*")
	if err != nil {
		return fmt.Errorf("create output temporary file: %w", err)
	}
	tempName := temp.Name()
	defer os.Remove(tempName)
	if err := temp.Chmod(0o600); err != nil {
		temp.Close()
		return fmt.Errorf("restrict output temporary file: %w", err)
	}
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
	if err := os.Link(tempName, path); err != nil {
		if os.IsExist(err) {
			return ErrOutputExists
		}
		return fmt.Errorf("publish external output without overwrite: %w", err)
	}
	return nil
}
