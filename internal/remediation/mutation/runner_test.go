package mutation

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/geoffrey-xiao/deprail/internal/remediation/isolation"
)

func testWorkspace(t *testing.T) isolation.Workspace {
	t.Helper()
	root := t.TempDir()
	run := func(args ...string) {
		t.Helper()
		if out, err := exec.Command(args[0], args[1:]...).CombinedOutput(); err != nil {
			t.Fatalf("%v: %v: %s", args, err, out)
		}
	}
	run("git", "-C", root, "init", "-q")
	if err := os.WriteFile(filepath.Join(root, "package.json"), []byte("{}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	run("git", "-C", root, "add", ".")
	run("git", "-C", root, "-c", "user.name=Test", "-c", "user.email=test@example.invalid", "commit", "-qm", "initial")
	workspace, err := isolation.Create(context.Background(), root, "HEAD")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = workspace.Remove(context.Background()) })
	return workspace
}

func TestRunRejectsUnapprovedUnisolatedOrUnsupportedExecution(t *testing.T) {
	if _, err := Run(context.Background(), Request{Path: "/usr/bin/npm", Approved: true, DenyScripts: true, DenyNetwork: true, Timeout: time.Second, OutputCap: 1024}); err == nil {
		t.Fatal("expected isolated workspace rejection")
	}
	workspace := testWorkspace(t)
	if _, err := Run(context.Background(), Request{Path: "/bin/sh", Workspace: workspace, Approved: true, DenyScripts: true, DenyNetwork: true, Timeout: time.Second, OutputCap: 1024}); err == nil {
		t.Fatal("expected shell rejection")
	}
	if _, err := Run(context.Background(), Request{Path: "/usr/bin/env", Workspace: workspace, Approved: true, DenyScripts: true, DenyNetwork: true, Timeout: time.Second, OutputCap: 1024}); err == nil {
		t.Fatal("expected unsupported executable rejection")
	}
}

func TestSecuredArgsAddOfflineAndScriptFlags(t *testing.T) {
	args, err := securedArgs("/usr/bin/npm", []string{"install"})
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"install", "--ignore-scripts", "--offline"} {
		found := false
		for _, arg := range args {
			if arg == want {
				found = true
			}
		}
		if !found {
			t.Fatalf("missing secured flag %q in %v", want, args)
		}
	}
}
