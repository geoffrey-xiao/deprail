package isolation

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func TestCreateAndRemoveWorktreePreservesSource(t *testing.T) {
	root := t.TempDir()
	run := func(args ...string) {
		t.Helper()
		if out, err := exec.Command(args[0], args[1:]...).CombinedOutput(); err != nil {
			t.Fatalf("%v: %v: %s", args, err, out)
		}
	}
	run("git", "-C", root, "init", "-q")
	path := filepath.Join(root, "manifest.txt")
	if err := os.WriteFile(path, []byte("before\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	run("git", "-C", root, "add", "manifest.txt")
	run("git", "-C", root, "-c", "user.name=Test", "-c", "user.email=test@example.invalid", "commit", "-qm", "initial")
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	ws, err := Create(context.Background(), root, "HEAD")
	if err != nil {
		t.Fatal(err)
	}
	canonical, err := filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatal(err)
	}
	if ws.Root != canonical {
		t.Fatalf("root = %q, want %q", ws.Root, canonical)
	}
	if err := os.WriteFile(filepath.Join(ws.Path, "manifest.txt"), []byte("changed\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := ws.Remove(context.Background()); err != nil {
		t.Fatal(err)
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != string(before) {
		t.Fatalf("source changed: %q", after)
	}
	if err := ws.Remove(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestSubdirectoryAllowsNestedWorkspaceWithoutTouchingSource(t *testing.T) {
	root := t.TempDir()
	run := func(args ...string) {
		t.Helper()
		if out, err := exec.Command(args[0], args[1:]...).CombinedOutput(); err != nil {
			t.Fatalf("%v: %v: %s", args, err, out)
		}
	}
	run("git", "-C", root, "init", "-q")
	nested := filepath.Join(root, "services", "api")
	if err := os.MkdirAll(nested, 0o700); err != nil {
		t.Fatal(err)
	}
	source := filepath.Join(nested, "package.json")
	if err := os.WriteFile(source, []byte("{}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	run("git", "-C", root, "add", ".")
	run("git", "-C", root, "-c", "user.name=Test", "-c", "user.email=test@example.invalid", "commit", "-qm", "initial")

	ws, err := Create(context.Background(), root, "HEAD")
	if err != nil {
		t.Fatal(err)
	}
	nestedWS, err := ws.Subdirectory("services/api")
	if err != nil {
		t.Fatal(err)
	}
	if nestedWS.Path != filepath.Join(ws.Path, "services", "api") || !nestedWS.Prepared() {
		t.Fatalf("nested workspace = %#v", nestedWS)
	}
	if err := os.WriteFile(filepath.Join(nestedWS.Path, "changed.txt"), []byte("isolated\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := ws.Remove(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, "services", "api", "changed.txt")); !os.IsNotExist(err) {
		t.Fatalf("isolated mutation reached source repository: %v", err)
	}
}

func TestSubdirectoryRejectsTraversalAndMissingPath(t *testing.T) {
	ws := Workspace{Root: t.TempDir(), Path: t.TempDir(), Created: true, verified: true}
	for _, relative := range []string{"../outside", filepath.Join("..", "outside"), "missing"} {
		if _, err := ws.Subdirectory(relative); err == nil {
			t.Fatalf("expected %q to be rejected", relative)
		}
	}
}

func TestSubdirectoryRejectsSymlinkedPath(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink permissions vary on Windows")
	}
	root := t.TempDir()
	target := t.TempDir()
	link := filepath.Join(root, "linked")
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}
	ws := Workspace{Root: root, Path: root, Created: true, verified: true}
	if _, err := ws.Subdirectory("linked"); err == nil {
		t.Fatal("expected symlinked workspace path to be rejected")
	}
}

func TestRemovePreservesStateWhenGitCleanupFails(t *testing.T) {
	path := filepath.Join(t.TempDir(), "worktree")
	ws := Workspace{Root: t.TempDir(), Path: path, Created: true, verified: true}
	result := ws.RemoveWithEvidence(context.Background())
	if result.Err == nil || result.GitRemoved || !result.FilesystemRemoved || !result.Retryable {
		t.Fatalf("cleanup result = %#v", result)
	}
	if !ws.Created || !ws.verified || ws.Path != path {
		t.Fatalf("workspace state was lost after cleanup failure: %#v", ws)
	}
}

func TestRemoveAttemptsFilesystemCleanupAfterGitCancellation(t *testing.T) {
	root := t.TempDir()
	run := func(args ...string) {
		t.Helper()
		if out, err := exec.Command(args[0], args[1:]...).CombinedOutput(); err != nil {
			t.Fatalf("%v: %v: %s", args, err, out)
		}
	}
	run("git", "-C", root, "init", "-q")
	file := filepath.Join(root, "manifest.txt")
	if err := os.WriteFile(file, []byte("before\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	run("git", "-C", root, "add", "manifest.txt")
	run("git", "-C", root, "-c", "user.name=Test", "-c", "user.email=test@example.invalid", "commit", "-qm", "initial")
	ws, err := Create(context.Background(), root, "HEAD")
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := ws.Remove(ctx); err == nil {
		t.Fatal("expected cancelled Git cleanup error")
	}
	if _, err := os.Stat(filepath.Dir(ws.Path)); !os.IsNotExist(err) {
		t.Fatalf("workspace parent still exists: %v", err)
	}
	if !ws.Created || ws.Path == "" {
		t.Fatalf("cleanup state was not retained for Git retry: %#v", ws)
	}
}

func TestCreateRejectsRepositorySubdirectory(t *testing.T) {
	root := t.TempDir()
	run := func(args ...string) {
		t.Helper()
		if out, err := exec.Command(args[0], args[1:]...).CombinedOutput(); err != nil {
			t.Fatalf("%v: %v: %s", args, err, out)
		}
	}
	run("git", "-C", root, "init", "-q")
	subdir := filepath.Join(root, "nested")
	if err := os.Mkdir(subdir, 0o700); err != nil {
		t.Fatal(err)
	}
	if _, err := Create(context.Background(), subdir, "HEAD"); err == nil {
		t.Fatal("expected repository subdirectory rejection")
	}
}

func TestCreateRejectsInvalidInput(t *testing.T) {
	if _, err := Create(context.Background(), "", "HEAD"); err == nil {
		t.Fatal("expected empty root rejection")
	}
	if _, err := Create(context.Background(), t.TempDir(), "HEAD\n--force"); err == nil {
		t.Fatal("expected source ref rejection")
	}
}
