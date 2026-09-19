package baseline

import (
	"crypto/sha256"
	"encoding/hex"
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

func TestStoreRejectsUnsafeIDsUppercaseDigestsUnknownFieldsAndOverwrite(t *testing.T) {
	unsafe := validBaseline()
	unsafe.BaselineID = "../escape"
	if _, err := (Store{Root: t.TempDir()}).Save(unsafe); err == nil {
		t.Fatal("unsafe baseline ID was accepted")
	}
	uppercase := validBaseline()
	uppercase.ArtifactDigests[0] = "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
	if _, err := (Store{Root: t.TempDir()}).Save(uppercase); err == nil {
		t.Fatal("uppercase artifact digest was accepted")
	}
	store := Store{Root: t.TempDir()}
	path, err := store.Save(validBaseline())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.Save(validBaseline()); err == nil {
		t.Fatal("existing baseline was overwritten")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	data = append(data[:len(data)-2], []byte(`,"unexpected":true}`+"\n")...)
	digest := sha256.Sum256(data)
	unknownPath := filepath.Join(store.Root, "base-001-"+hex.EncodeToString(digest[:])+".json")
	if err := os.WriteFile(unknownPath, data, 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Load(unknownPath); err == nil {
		t.Fatal("unknown baseline field was accepted")
	}
}
func TestValidateRejectsInvalidAliases(t *testing.T) {
	empty := validBaseline()
	empty.Findings[0].VulnerabilityAliases = []string{""}
	if err := Validate(empty); err == nil {
		t.Fatal("empty alias was accepted")
	}
	duplicate := validBaseline()
	duplicate.Findings[0].VulnerabilityAliases = []string{"CVE-1", "CVE-1"}
	if err := Validate(duplicate); err == nil {
		t.Fatal("duplicate alias was accepted")
	}
}
