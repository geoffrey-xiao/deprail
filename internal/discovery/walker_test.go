package discovery

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestWalkReturnsSortedRepositoryRelativePathsAndAppliesIgnores(t *testing.T) {
	root := t.TempDir()
	mustWriteFile(t, filepath.Join(root, "z", "last.txt"), "last")
	mustWriteFile(t, filepath.Join(root, "a", "first.txt"), "first")
	mustWriteFile(t, filepath.Join(root, ".git", "ignored"), "ignored")
	mustWriteFile(t, filepath.Join(root, "node_modules", "ignored"), "ignored")

	result, err := Walk(context.Background(), root, WalkOptions{})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"a/first.txt", "z/last.txt"}
	if len(result.Paths) != len(want) {
		t.Fatalf("paths = %#v, want %#v", result.Paths, want)
	}
	for i := range want {
		if result.Paths[i] != want[i] {
			t.Errorf("paths[%d] = %q, want %q", i, result.Paths[i], want[i])
		}
	}
}

func TestWalkReportsExternalSymlinkWithoutFollowingIt(t *testing.T) {
	root := t.TempDir()
	external := t.TempDir()
	mustWriteFile(t, filepath.Join(external, "secret.txt"), "secret")
	if err := os.Symlink(external, filepath.Join(root, "escape")); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}

	result, err := Walk(context.Background(), root, WalkOptions{})
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range result.Paths {
		if path == "escape/secret.txt" {
			t.Fatal("walker followed external symlink")
		}
	}
	if !hasDiagnostic(result.Diagnostics, "PATH_OUTSIDE_ROOT", "escape") {
		t.Fatalf("diagnostics = %#v, want PATH_OUTSIDE_ROOT for escape", result.Diagnostics)
	}
}

func TestWalkRejectsFileRoot(t *testing.T) {
	root := filepath.Join(t.TempDir(), "file")
	mustWriteFile(t, root, "not a directory")
	if _, err := Walk(context.Background(), root, WalkOptions{}); err == nil {
		t.Fatal("expected file root to be rejected")
	}
}

func TestWalkCancellationStopsBeforeTraversalCompletes(t *testing.T) {
	root := t.TempDir()
	mustWriteFile(t, filepath.Join(root, "file.txt"), "content")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := Walk(ctx, root, WalkOptions{}); err == nil {
		t.Fatal("expected cancellation error")
	}
}

func hasDiagnostic(diagnostics []Diagnostic, code, scope string) bool {
	for _, diagnostic := range diagnostics {
		if diagnostic.Code == code && diagnostic.Scope == scope {
			return true
		}
	}
	return false
}

func mustWriteFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
