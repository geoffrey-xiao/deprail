package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/geoffrey-xiao/deprail/internal/baseline"
)

func writeCompleteScan(t *testing.T, path string) {
	t.Helper()
	data := []byte(`{
  "schema_version": "v1alpha",
  "document_type": "scan",
  "scan_id": "scan-001",
  "repository_state": "tree-001",
  "status": "complete",
  "findings": [{
    "component": "lodash",
    "purl": "pkg:npm/lodash@4.17.20",
    "version": "4.17.20",
    "aliases": ["CVE-2021-23337"],
    "severity": "high",
    "target_id": "OSV-1",
    "workspace_id": "root"
  }],
  "errors": [],
  "artifact_digests": []
}
`)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestBaselineCreateWritesDeterministicBaselineConsumedByDiff(t *testing.T) {
	root := t.TempDir()
	scanPath := filepath.Join(root, "scan.json")
	basePath := filepath.Join(root, "baseline.json")
	headPath := filepath.Join(root, "baseline-head.json")
	writeCompleteScan(t, scanPath)

	var stdout, stderr bytes.Buffer
	if code := run([]string{"baseline", "create", "--scan", scanPath, "--output", basePath, "--format", "json"}, &stdout, &stderr); code != 0 {
		t.Fatalf("baseline create code=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr=%q", stderr.String())
	}
	var output baseline.Document
	if err := json.Unmarshal(stdout.Bytes(), &output); err != nil {
		t.Fatalf("stdout is not baseline JSON: %v", err)
	}
	if err := baseline.Validate(output); err != nil {
		t.Fatal(err)
	}
	if len(output.Findings) != 1 || output.Findings[0].TargetID != "OSV-1" {
		t.Fatalf("finding target identity = %#v", output.Findings)
	}
	if _, err := baseline.LoadFile(basePath); err != nil {
		t.Fatalf("generated baseline cannot be loaded: %v", err)
	}
	data, err := os.ReadFile(basePath)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(headPath, data, 0o600); err != nil {
		t.Fatal(err)
	}

	stdout.Reset()
	stderr.Reset()
	if code := run([]string{"diff", "--base", basePath, "--head", headPath, "--format", "json"}, &stdout, &stderr); code != 0 {
		t.Fatalf("diff code=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), `"document_type":"diff"`) {
		t.Fatalf("diff output=%q", stdout.String())
	}
}

func TestBaselineCreateRejectsIncompleteAndExistingOutput(t *testing.T) {
	root := t.TempDir()
	scanPath := filepath.Join(root, "scan.json")
	outputPath := filepath.Join(root, "baseline.json")
	writeCompleteScan(t, scanPath)
	data, err := os.ReadFile(scanPath)
	if err != nil {
		t.Fatal(err)
	}
	data = bytes.Replace(data, []byte(`"complete"`), []byte(`"partial"`), 1)
	if err := os.WriteFile(scanPath, data, 0o600); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	if code := run([]string{"baseline", "create", "--scan", scanPath, "--output", outputPath}, &stdout, &stderr); code != 3 {
		t.Fatalf("incomplete code=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	if _, err := os.Stat(outputPath); !os.IsNotExist(err) {
		t.Fatalf("incomplete scan created output: err=%v", err)
	}

	writeCompleteScan(t, scanPath)
	stdout.Reset()
	stderr.Reset()
	if code := run([]string{"baseline", "create", "--scan", scanPath, "--output", outputPath}, &stdout, &stderr); code != 0 {
		t.Fatalf("valid create code=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	info, err := os.Stat(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("output permissions = %o, want 600", info.Mode().Perm())
	}
	before, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	stdout.Reset()
	stderr.Reset()
	if code := run([]string{"baseline", "create", "--scan", scanPath, "--output", outputPath}, &stdout, &stderr); code != 3 || !strings.Contains(stderr.String(), "BASELINE_OUTPUT_EXISTS") {
		t.Fatalf("existing output code=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	after, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(before, after) {
		t.Fatal("existing baseline was modified")
	}
}

func TestBaselineCreateHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := run([]string{"baseline", "create", "--help"}, &stdout, &stderr); code != 0 {
		t.Fatalf("help code=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "--scan path") || !strings.Contains(stdout.String(), "--output path") {
		t.Fatalf("help=%q", stdout.String())
	}
}
