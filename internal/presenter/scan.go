package presenter

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/geoffrey-xiao/deprail/internal/app"
)

func WriteScanJSON(w io.Writer, report app.ScanReport) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(report)
}

func WriteScanTerminal(w io.Writer, report app.ScanReport, verbose bool) error {
	if err := WriteStatus(w, string(report.Status)); err != nil {
		return err
	}
	if err := WriteSummary(w, "Findings", len(report.Findings)); err != nil {
		return err
	}
	if verbose || len(report.Errors) > 0 {
		if err := WriteSection(w, "Diagnostics"); err != nil {
			return err
		}
		for _, diagnostic := range report.Errors {
			if err := WriteError(w, "", diagnostic); err != nil {
				return err
			}
		}
	}
	return nil
}

func WriteAtomic(path string, data []byte) error {
	if path == "" {
		return fmt.Errorf("output path is required")
	}
	if _, err := os.Lstat(path); err == nil {
		return fmt.Errorf("output file already exists")
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("inspect output file: %w", err)
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}
	tmp, err := os.CreateTemp(dir, ".deprail-output-*")
	if err == nil {
		_ = tmp.Chmod(0o600)
	}
	if err != nil {
		return fmt.Errorf("create output file: %w", err)
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)
	if _, err = tmp.Write(data); err == nil {
		err = tmp.Close()
	} else {
		_ = tmp.Close()
	}
	if err != nil {
		return fmt.Errorf("write output file: %w", err)
	}
	if err := os.Rename(tmpName, path); err != nil {
		return fmt.Errorf("publish output file: %w", err)
	}
	return nil
}
