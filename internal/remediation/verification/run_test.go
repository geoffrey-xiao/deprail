package verification

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRunExecutesSelectedCommandsInWorkspace(t *testing.T) {
	workspace := t.TempDir()
	results, err := Run(context.Background(), workspace, []Command{{ID: "test", Kind: Test, Path: os.Args[0], Args: []string{"-test.run=TestVerificationHelper", "--", "ok"}, WorkingDirectory: ".", Enabled: true}}, time.Second, 1024)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0].Process.ExitCode != 0 || !strings.HasPrefix(string(results[0].Process.Stdout), "ok\n") {
		t.Fatalf("results = %#v", results)
	}
}

func TestRunStopsAfterFailureAndRetainsPartialResults(t *testing.T) {
	workspace := t.TempDir()
	commands := []Command{
		{ID: "fail", Kind: Test, Path: os.Args[0], Args: []string{"-test.run=TestVerificationHelper", "--", "fail"}, WorkingDirectory: ".", Enabled: true},
		{ID: "later", Kind: Test, Path: os.Args[0], Args: []string{"-test.run=TestVerificationHelper", "--", "ok"}, WorkingDirectory: ".", Enabled: true},
	}
	results, err := Run(context.Background(), workspace, commands, time.Second, 1024)
	if err == nil || len(results) != 1 || results[0].Command.ID != "fail" {
		t.Fatalf("results=%#v err=%v", results, err)
	}
}

func TestRunClassifiesTimeout(t *testing.T) {
	results, err := Run(context.Background(), t.TempDir(), []Command{{ID: "slow", Kind: Test, Path: os.Args[0], Args: []string{"-test.run=TestVerificationHelper", "--", "sleep"}, WorkingDirectory: ".", Enabled: true}}, 20*time.Millisecond, 1024)
	if err == nil || len(results) != 1 || !strings.Contains(err.Error(), string("SCANNER_TIMEOUT")) {
		t.Fatalf("results=%#v err=%v", results, err)
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
	_, err := Run(context.Background(), workspace, []Command{{ID: "escape", Kind: Test, Path: os.Args[0], WorkingDirectory: "outside", Enabled: true}}, time.Second, 1024)
	if err == nil {
		t.Fatal("expected workspace escape rejection")
	}
}
