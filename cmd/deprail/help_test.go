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

func TestDiffHelpDocumentsComparisonFlags(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := run([]string{"diff", "--help"}, &stdout, &stderr); code != 0 {
		t.Fatalf("exit code = %d, stderr=%q", code, stderr.String())
	}
	for _, text := range []string{"--base", "--head", "--format"} {
		if !strings.Contains(stdout.String(), text) {
			t.Fatalf("help missing %q: %q", text, stdout.String())
		}
	}
}

func TestHelpRejectsUnknownCommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := run([]string{"scna", "--help"}, &stdout, &stderr); code != 2 {
		t.Fatalf("exit code = %d, want 2; stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	if stdout.Len() != 0 || !strings.Contains(stderr.String(), "CONFIG_INVALID") {
		t.Fatalf("stdout=%q stderr=%q", stdout.String(), stderr.String())
	}
}

func TestFixPlanHelpDocumentsOutput(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := run([]string{"fix", "plan", "--help"}, &stdout, &stderr); code != 0 {
		t.Fatalf("exit code = %d, stderr=%q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "--output") {
		t.Fatalf("help missing --output: %q", stdout.String())
	}
}

func TestFixApplyHelpDocumentsApprovalAndDryRun(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := run([]string{"fix", "apply", "--help"}, &stdout, &stderr); code != 0 {
		t.Fatalf("exit code = %d, stderr=%q", code, stderr.String())
	}
	for _, text := range []string{"--plan", "--approval", "--dry-run"} {
		if !strings.Contains(stdout.String(), text) {
			t.Fatalf("help missing %q: %q", text, stdout.String())
		}
	}
}
