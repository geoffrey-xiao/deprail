package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestDiscoverJSONWritesDataToStdoutOnly(t *testing.T) {
	fixture := fixturePath(t, "npm-basic")
	var stdout, stderr bytes.Buffer
	code := run([]string{"discover", fixture, "--format", "json"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit code = %d, stderr = %q", code, stderr.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
	var graph struct {
		DocumentType string `json:"document_type"`
		Completeness string `json:"completeness"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &graph); err != nil {
		t.Fatalf("stdout is not JSON: %v", err)
	}
	if graph.DocumentType != "project" || graph.Completeness != "complete" {
		t.Fatalf("graph = %#v", graph)
	}
}

func TestDiscoverIncompleteProjectReturnsCodeThreeAndJSON(t *testing.T) {
	fixture := fixturePath(t, "npm-basic")
	var stdout, stderr bytes.Buffer
	code := run([]string{"discover", fixture, "--format", "json"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("fixture should remain complete, exit code = %d, stderr = %q", code, stderr.String())
	}

	root := t.TempDir()
	writeTestFile(t, filepath.Join(root, "package.json"), `{"name":"incomplete"}`)
	stdout.Reset()
	stderr.Reset()
	code = run([]string{"discover", root, "--format", "json"}, &stdout, &stderr)
	if code != 3 {
		t.Fatalf("exit code = %d, want 3; stderr = %q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"completeness": "partial"`) {
		t.Fatalf("stdout = %q, want partial graph", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty for graph result", stderr.String())
	}
}

func TestDiscoverRejectsMultiplePaths(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"discover", ".", "other"}, &stdout, &stderr)
	if code != 2 || stdout.Len() != 0 || !strings.Contains(stderr.String(), "CONFIG_INVALID") {
		t.Fatalf("code=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
}

func fixturePath(t *testing.T, fixture string) string {
	t.Helper()
	_, source, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("locate command test source")
	}
	return filepath.Join(filepath.Dir(source), "..", "..", "testdata", "fixtures", fixture)
}

func writeTestFile(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}
