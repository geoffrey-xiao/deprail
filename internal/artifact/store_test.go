package artifact

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStorePutGetIsContentAddressedAndIdempotent(t *testing.T) {
	store := Store{Root: t.TempDir(), MaxBytes: 1024}
	data := []byte("scanner output")
	first, err := store.Put(data)
	if err != nil {
		t.Fatalf("put: %v", err)
	}
	second, err := store.Put(data)
	if err != nil || first.Digest != second.Digest {
		t.Fatalf("repeat put = %#v, %#v; err = %v", first, second, err)
	}
	got, err := store.Get(first.Digest)
	if err != nil || string(got) != string(data) {
		t.Fatalf("get = %q, err = %v", got, err)
	}
	if _, err := os.Stat(filepath.Join(store.Root, first.Digest[:2], first.Digest)); err != nil {
		t.Fatalf("artifact path: %v", err)
	}
}

func TestStoreRejectsOversizeAndTamperedArtifacts(t *testing.T) {
	store := Store{Root: t.TempDir(), MaxBytes: 3}
	if _, err := store.Put([]byte("four")); !IsCode(err, ErrOutputLimit) {
		t.Fatalf("oversize error = %v", err)
	}
	store.MaxBytes = 100
	artifact, err := store.Put([]byte("safe"))
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(store.Root, artifact.Digest[:2], artifact.Digest)
	if err := os.WriteFile(path, []byte("tampered"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Get(artifact.Digest); !IsCode(err, ErrIntegrity) {
		t.Fatalf("tamper error = %v", err)
	}
}
