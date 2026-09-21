package evidence

import (
	"os"
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
