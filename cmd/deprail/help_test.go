package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRootHelpIsSuccessfulAndDeterministic(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := run([]string{"--help"}, &stdout, &stderr); code != 0 {
		t.Fatalf("exit code = %d, stderr = %q", code, stderr.String())
	}
	if stderr.Len() != 0 || !strings.Contains(stdout.String(), "Usage: deprail <command>") {
		t.Fatalf("stdout=%q stderr=%q", stdout.String(), stderr.String())
	}
}

func TestScanHelpDocumentsQuietAndVerbose(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := run([]string{"scan", "--help"}, &stdout, &stderr); code != 0 {
		t.Fatalf("exit code = %d, stderr = %q", code, stderr.String())
	}
	for _, text := range []string{"--quiet", "--verbose", "--format", "--output"} {
		if !strings.Contains(stdout.String(), text) {
			t.Fatalf("help missing %q: %q", text, stdout.String())
		}
	}
}
