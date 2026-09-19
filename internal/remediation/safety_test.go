package remediation

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestValidateExternalOutputRejectsRootTraversalAndSymlinkEscapes(t *testing.T) {
	root := t.TempDir()
	if _, err := ValidateExternalOutput(root, filepath.Join(root, "plan.json")); !errors.Is(err, ErrUnsafePath) {
		t.Fatalf("inside output error = %v", err)
	}
	if _, err := ValidateExternalOutput(root, filepath.Join(root, "nested", "..", "plan.json")); !errors.Is(err, ErrUnsafePath) {
		t.Fatalf("traversal output error = %v", err)
	}
	external := t.TempDir()
	link := filepath.Join(root, "external")
	if err := os.Symlink(external, link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	if _, err := ValidateExternalOutput(root, filepath.Join(link, "plan.json")); !errors.Is(err, ErrUnsafePath) {
		t.Fatalf("symlink output error = %v", err)
	}
}

func TestWriteExternalOutputIsRestrictiveAtomicAndNonOverwriting(t *testing.T) {
	root := t.TempDir()
	external := t.TempDir()
	path := filepath.Join(external, "plan.json")
	written, err := WriteExternalOutput(root, path, []byte("first"))
	if err != nil {
		t.Fatal(err)
	}
	expectedRoot, err := filepath.EvalSymlinks(external)
	if err != nil {
		t.Fatal(err)
	}
	expected := filepath.Join(expectedRoot, "plan.json")
	if written != expected {
		t.Fatalf("written path = %q, want %q", written, expected)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if runtime.GOOS != "windows" && info.Mode().Perm() != 0o600 {
		t.Fatalf("mode = %o, want 600", info.Mode().Perm())
	}
	data, err := os.ReadFile(path)
	if err != nil || string(data) != "first" {
		t.Fatalf("data = %q, err = %v", data, err)
	}
	if _, err := WriteExternalOutput(root, path, []byte("second")); !errors.Is(err, ErrOutputExists) {
		t.Fatalf("overwrite error = %v", err)
	}
	data, err = os.ReadFile(path)
	if err != nil || string(data) != "first" {
		t.Fatalf("data after overwrite = %q, err = %v", data, err)
	}
}

func TestSnapshotRepositoryDetectsSuccessfulAndFailedPlanningMutation(t *testing.T) {
	root := t.TempDir()
	manifest := filepath.Join(root, "package.json")
	if err := os.WriteFile(manifest, []byte(`{"name":"fixture"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	before, err := SnapshotRepository(root)
	if err != nil {
		t.Fatal(err)
	}
	if !before.Equal(before) || len(before.Paths()) != 1 {
		t.Fatalf("baseline snapshot = %#v", before)
	}
	if err := os.WriteFile(filepath.Join(root, "generated.tmp"), []byte("mutation"), 0o600); err != nil {
		t.Fatal(err)
	}
	after, err := SnapshotRepository(root)
	if err != nil {
		t.Fatal(err)
	}
	if err := EnsureRepositoryUnchanged(before, after); !errors.Is(err, ErrRepositoryChange) {
		t.Fatalf("mutation error = %v", err)
	}
}

func TestSnapshotRepositoryRecordsSymlinkMetadataWithoutFollowingIt(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink creation may require elevated privileges on Windows")
	}
	root := t.TempDir()
	target := filepath.Join(t.TempDir(), "outside.txt")
	if err := os.WriteFile(target, []byte("outside"), 0o600); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(root, "link.txt")
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}
	snapshot, err := SnapshotRepository(root)
	if err != nil {
		t.Fatal(err)
	}
	if got := snapshot.Entries["link.txt"]; got != "symlink:"+target {
		t.Fatalf("symlink entry = %q", got)
	}
}
