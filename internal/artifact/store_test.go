package artifact

import (
	"context"
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

func TestStoreVerifyStreamsTrustedArtifactAndReportsMissing(t *testing.T) {
	store := Store{Root: t.TempDir(), MaxBytes: 1024}
	artifact, err := store.Put([]byte("content-addressed output"))
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Verify(context.Background(), artifact.Digest); err != nil {
		t.Fatalf("verify: %v", err)
	}
	if err := store.Verify(context.Background(), "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"); !IsCode(err, ErrNotFound) {
		t.Fatalf("missing artifact error = %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := store.Verify(ctx, artifact.Digest); !IsCode(err, ErrWriteFailed) {
		t.Fatalf("cancelled verification error = %v", err)
	}
}

func TestStoreVerifyRejectsSymlinkedArtifactPathAndDigestMismatch(t *testing.T) {
	store := Store{Root: t.TempDir(), MaxBytes: 1024}
	artifact, err := store.Put([]byte("original"))
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(store.Root, artifact.Digest[:2], artifact.Digest)
	if err := os.WriteFile(path, []byte("changed"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := store.Verify(context.Background(), artifact.Digest); !IsCode(err, ErrIntegrity) {
		t.Fatalf("digest mismatch error = %v", err)
	}
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(t.TempDir(), "target")
	if err := os.WriteFile(target, []byte("original"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, path); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	if err := store.Verify(context.Background(), artifact.Digest); !IsCode(err, ErrIntegrity) {
		t.Fatalf("symlink verification error = %v", err)
	}
}
