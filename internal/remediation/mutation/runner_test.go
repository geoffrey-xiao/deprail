package mutation

import (
	"context"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestRunUsesApprovedDirectExecutable(t *testing.T) {
	goPath, err := exec.LookPath("go")
	if err != nil {
		t.Skip("go is unavailable")
	}
	result, err := Run(context.Background(), Request{Path: goPath, Args: []string{"version"}, WorkspaceRoot: filepath.Dir(goPath), Approved: true, Timeout: 10 * time.Second, OutputCap: 1 << 20})
	if err != nil {
		t.Fatal(err)
	}
	if result.ExitCode != 0 {
		t.Fatalf("exit code = %d", result.ExitCode)
	}
}

func TestRunRejectsUnapprovedOrShellExecution(t *testing.T) {
	if _, err := Run(context.Background(), Request{Path: "/usr/bin/env", WorkspaceRoot: "/tmp", Timeout: time.Second, OutputCap: 1024}); err == nil {
		t.Fatal("expected approval rejection")
	}
	if _, err := Run(context.Background(), Request{Path: "/bin/sh", WorkspaceRoot: "/tmp", Approved: true, Timeout: time.Second, OutputCap: 1024}); err == nil {
		t.Fatal("expected shell rejection")
	}
}

func TestRunRejectsScriptsAndNetworkPolicy(t *testing.T) {
	base := Request{Path: "/usr/bin/env", WorkspaceRoot: "/tmp", Approved: true, Timeout: time.Second, OutputCap: 1024}
	base.AllowScripts = true
	if _, err := Run(context.Background(), base); err == nil {
		t.Fatal("expected script policy rejection")
	}
	base.AllowScripts = false
	base.AllowNetwork = true
	if _, err := Run(context.Background(), base); err == nil {
		t.Fatal("expected network policy rejection")
	}
}
