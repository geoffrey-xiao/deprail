package main

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/geoffrey-xiao/deprail/internal/app"
)

func TestBaselineCreateRejectsStaleScan(t *testing.T) {
	root := t.TempDir()
	scanPath := filepath.Join(root, "scan.json")
	outputPath := filepath.Join(root, "baseline.json")
	state, err := app.CurrentRepositoryState(root)
	if err != nil {
		t.Fatal(err)
	}
	writeCompleteScan(t, scanPath, root, state)
	if err := os.WriteFile(filepath.Join(root, "changed.txt"), []byte("changed"), 0o600); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	if code := run([]string{"baseline", "create", "--scan", scanPath, "--output", outputPath}, &stdout, &stderr); code != 3 || !strings.Contains(stderr.String(), "BASELINE_INPUT_STALE") {
		t.Fatalf("stale code=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	if _, err := os.Stat(outputPath); !os.IsNotExist(err) {
		t.Fatalf("stale scan created output: err=%v", err)
	}
}

func TestBaselineCreateRejectsSymlinkedOutputDirectory(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Windows CI does not guarantee symlink creation privileges")
	}
	root := t.TempDir()
	outside := t.TempDir()
	scanPath := filepath.Join(root, "scan.json")
	state, err := app.CurrentRepositoryState(root)
	if err != nil {
		t.Fatal(err)
	}
	writeCompleteScan(t, scanPath, root, state)
	link := filepath.Join(root, "escape")
	if err := os.Symlink(outside, link); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	outputPath := filepath.Join(link, "baseline.json")
	if code := run([]string{"baseline", "create", "--scan", scanPath, "--output", outputPath}, &stdout, &stderr); code != 3 || !strings.Contains(stderr.String(), "PATH_OUTSIDE_ROOT") {
		t.Fatalf("symlink output code=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	if _, err := os.Stat(filepath.Join(outside, "baseline.json")); !os.IsNotExist(err) {
		t.Fatalf("symlink output escaped root: err=%v", err)
	}
}
