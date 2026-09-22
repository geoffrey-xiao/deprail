package verification

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func testTool(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	name := "npm"
	if runtime.GOOS == "windows" {
		name = "npm.exe"
	}
	path := filepath.Join(dir, name)
	data, err := os.ReadFile(os.Args[0])
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o700); err != nil {
		t.Fatal(err)
	}
	return path
}
func TestValidateCommandSpecRestrictsNodeAndAllowlist(t *testing.T) {
	if err := ValidateCommandSpec("node", []string{"--check", "verify.js"}); err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		path string
		args []string
	}{
		{"node", []string{"-e", "process.exit(0)"}},
		{"node", []string{"--check", "../outside.js"}},
		{"/tmp/node", []string{"--check", "verify.js"}},
		{"make", []string{"test"}},
	} {
		if err := ValidateCommandSpec(test.path, test.args); err == nil {
			t.Fatalf("accepted unsafe command %q %#v", test.path, test.args)
		}
	}
}

func TestValidateCommandInWorkspaceRejectsSymlinkEscape(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink fixture requires platform support")
	}
	workspace := t.TempDir()
	outside := t.TempDir()
	if err := os.WriteFile(filepath.Join(outside, "verify.js"), []byte("const outside = true;\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(outside, "verify.js"), filepath.Join(workspace, "verify.js")); err != nil {
		t.Fatal(err)
	}
	err := ValidateCommandInWorkspace(Command{Path: "/usr/bin/node", Args: []string{"--check", "verify.js"}, WorkingDirectory: "."}, workspace)
	if err == nil {
		t.Fatal("expected symlink escape to be rejected")
	}
}
func TestRunExecutesSelectedCommandsInWorkspace(t *testing.T) {
	workspace := t.TempDir()
	tool := testTool(t)
	results, err := Run(context.Background(), workspace, []Command{{ID: "test", Kind: Test, Path: tool, Args: []string{"-test.run=TestVerificationHelper", "--", "ok"}, WorkingDirectory: ".", Enabled: true}}, time.Second, 1024)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0].Process.ExitCode != 0 || !strings.HasPrefix(string(results[0].Process.Stdout), "ok\n") {
		t.Fatalf("results = %#v", results)
	}
}

func TestRunStopsAfterFailureAndRetainsPartialResults(t *testing.T) {
	workspace := t.TempDir()
	tool := testTool(t)
	commands := []Command{
		{ID: "fail", Kind: Test, Path: tool, Args: []string{"-test.run=TestVerificationHelper", "--", "fail"}, WorkingDirectory: ".", Enabled: true},
		{ID: "later", Kind: Test, Path: tool, Args: []string{"-test.run=TestVerificationHelper", "--", "ok"}, WorkingDirectory: ".", Enabled: true},
	}
	results, err := Run(context.Background(), workspace, commands, time.Second, 1024)
	if err == nil || len(results) != 1 || results[0].Command.ID != "fail" {
		t.Fatalf("results=%#v err=%v", results, err)
	}
}

func TestRunClassifiesTimeout(t *testing.T) {
	tool := testTool(t)
	results, err := Run(context.Background(), t.TempDir(), []Command{{ID: "slow", Kind: Test, Path: tool, Args: []string{"-test.run=TestVerificationHelper", "--", "sleep"}, WorkingDirectory: ".", Enabled: true}}, 20*time.Millisecond, 1024)
	if err == nil || len(results) != 1 || !strings.Contains(err.Error(), string("SCANNER_TIMEOUT")) {
		t.Fatalf("results=%#v err=%v", results, err)
	}
}

func TestRunRejectsEmptyAndUntrustedCommands(t *testing.T) {
	if _, err := Run(context.Background(), t.TempDir(), nil, time.Second, 1024); !errors.Is(err, ErrNoCommands) {
		t.Fatalf("empty commands error = %v", err)
	}
	_, err := Run(context.Background(), t.TempDir(), []Command{{ID: "unsafe", Kind: Test, Path: "/usr/bin/rm", WorkingDirectory: ".", Enabled: true}}, time.Second, 1024)
	if err == nil {
		t.Fatal("expected untrusted executable rejection")
	}
}

func TestRunPreservesCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	tool := testTool(t)
	_, err := Run(ctx, t.TempDir(), []Command{{ID: "cancelled", Kind: Test, Path: tool, Args: []string{"-test.run=TestVerificationHelper", "--", "sleep"}, WorkingDirectory: ".", Enabled: true}}, time.Second, 1024)
	if err == nil || !strings.Contains(err.Error(), "SCANNER_CANCELLED") {
		t.Fatalf("cancellation error = %v", err)
	}
}

func TestVerificationHelper(t *testing.T) {
	for i, arg := range os.Args {
		if arg != "--" || i+1 >= len(os.Args) {
			continue
		}
		switch os.Args[i+1] {
		case "ok":
			_, _ = os.Stdout.WriteString("ok\n")
		case "fail":
			_, _ = os.Stderr.WriteString("failed\n")
			os.Exit(7)
		case "sleep":
			time.Sleep(time.Second)
		}
	}
}

func TestRunRejectsWorkspaceEscape(t *testing.T) {
	workspace := t.TempDir()
	outside := t.TempDir()
	if err := os.Symlink(outside, filepath.Join(workspace, "outside")); err != nil {
		t.Skip("symlinks unavailable")
	}
	tool := testTool(t)
	_, err := Run(context.Background(), workspace, []Command{{ID: "escape", Kind: Test, Path: tool, WorkingDirectory: "outside", Enabled: true}}, time.Second, 1024)
	if err == nil {
		t.Fatal("expected workspace escape rejection")
	}
}
