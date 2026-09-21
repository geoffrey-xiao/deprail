package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/geoffrey-xiao/deprail/internal/remediation"
	"github.com/geoffrey-xiao/deprail/internal/remediation/verification"
)

func TestApplyVerificationCommandsResolvesSupportedTools(t *testing.T) {
	dir := t.TempDir()
	tool := filepath.Join(dir, "npm")
	data, err := os.ReadFile(os.Args[0])
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(tool, data, 0o700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir)

	commands, err := applyVerificationCommands(remediation.Plan{Verification: []remediation.Verification{{
		ID:      "tests",
		Command: remediation.Command{Executable: "npm", Arguments: []string{"test"}, WorkingDirectory: "services/api"},
		Reason:  "run repository tests",
	}}})
	if err != nil {
		t.Fatal(err)
	}
	if len(commands) != 1 {
		t.Fatalf("commands = %#v", commands)
	}
	command := commands[0]
	if command.ID != "tests" || command.Kind != verification.Test || command.Path != tool || !command.Enabled || command.WorkingDirectory != "services/api" {
		t.Fatalf("command = %#v", command)
	}
	if len(command.Args) != 1 || command.Args[0] != "test" {
		t.Fatalf("args = %#v", command.Args)
	}
}

func TestApplyVerificationCommandsRejectsMissingTools(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	_, err := applyVerificationCommands(remediation.Plan{Verification: []remediation.Verification{{
		ID:      "tests",
		Command: remediation.Command{Executable: "npm", Arguments: []string{"test"}, WorkingDirectory: "."},
	}}})
	if err == nil {
		t.Fatal("expected missing verification tool error")
	}
}
