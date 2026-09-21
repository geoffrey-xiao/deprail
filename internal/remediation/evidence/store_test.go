package evidence

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestStoreSaveIsDeterministicAndNonOverwriting(t *testing.T) {
	store := Store{Root: t.TempDir()}
	record := testRecord()
	firstPath, firstDigest, err := store.Save(record)
	if err != nil {
		t.Fatal(err)
	}
	secondPath, secondDigest, err := store.Save(record)
	if err != nil {
		t.Fatal(err)
	}
	if firstPath != secondPath || firstDigest != secondDigest {
		t.Fatalf("save paths/digests differ: %q/%q vs %q/%q", firstPath, firstDigest, secondPath, secondDigest)
	}
	data, err := os.ReadFile(firstPath)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 || string(data) == "" {
		t.Fatal("evidence file is empty")
	}
}

func TestStoreRejectsMissingRoot(t *testing.T) {
	if _, _, err := (Store{}).Save(testRecord()); err == nil {
		t.Fatal("expected missing root rejection")
	}
}
func TestStoreRejectsSymlinkedShard(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	record := testRecord()
	digest, err := Digest(record)
	if err != nil {
		t.Fatal(err)
	}
	shard := filepath.Join(root, digest[:2])
	if err := os.Symlink(outside, shard); err != nil {
		t.Skip("symlinks unavailable")
	}
	if _, _, err := (Store{Root: root}).Save(record); err == nil {
		t.Fatal("expected symlinked shard rejection")
	}
}

func TestStoreRestrictsPermissiveRoot(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Windows does not expose POSIX directory permission semantics")
	}
	root := t.TempDir()
	if err := os.Chmod(root, 0o777); err != nil {
		t.Fatal(err)
	}
	path, _, err := (Store{Root: root}).Save(testRecord())
	if err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(root)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm()&0o077 != 0 {
		t.Fatalf("root permissions remain permissive: %o", info.Mode().Perm())
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatal(err)
	}
}
