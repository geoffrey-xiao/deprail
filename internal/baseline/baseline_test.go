package baseline

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/geoffrey-xiao/deprail/internal/discovery"
)

func validBaseline() Document {
	return Document{
		SchemaVersion: SchemaVersion,
		DocumentType:  DocumentType,
		BaselineID:    "base-001",
		SourceScanID:  "scan-001",
		Status:        discovery.Complete,
		Findings: []Finding{
			{StableKey: "b", Component: "pkg:npm/b@1", Version: "1"},
			{StableKey: "a", Component: "pkg:npm/a@1", Version: "1"},
		},
		ArtifactDigests: []string{"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
	}
}

func TestStoreRoundTripCanonicalizesAndPreservesIdentity(t *testing.T) {
	store := Store{Root: t.TempDir()}
	path, err := store.Save(validBaseline())
	if err != nil {
		t.Fatal(err)
	}
	got, err := store.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	want := validBaseline()
	want.Canonicalize()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("baseline = %#v, want %#v", got, want)
	}
	if filepath.Dir(path) != store.Root {
		t.Fatalf("path = %q, want store root %q", path, store.Root)
	}
}

func TestStoreRejectsIncompleteAndTamperedBaselines(t *testing.T) {
	invalid := validBaseline()
	invalid.Status = discovery.Partial
	if _, err := (Store{Root: t.TempDir()}).Save(invalid); err == nil {
		t.Fatal("partial baseline was accepted")
	}
	store := Store{Root: t.TempDir()}
	path, err := store.Save(validBaseline())
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("tampered"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Load(path); err == nil {
		t.Fatal("tampered baseline was accepted")
	}
}
