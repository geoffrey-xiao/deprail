package baseline

import (
	"reflect"
	"testing"
)

func TestCompareIsOrderIndependentAndClassifiesChanges(t *testing.T) {
	base := validBaseline()
	head := validBaseline()
	base.Findings = []Finding{{StableKey: "same", Component: "pkg:npm/a@1", Version: "1"}, {StableKey: "resolved", Component: "pkg:npm/b@1", Version: "1"}}
	head.Findings = []Finding{{StableKey: "same", Component: "pkg:npm/a@2", Version: "2"}, {StableKey: "added", Component: "pkg:npm/c@2", Version: "2"}}
	base.ArtifactDigests = []string{"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}
	head.ArtifactDigests = []string{"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"}
	first, err := Compare(base, head)
	if err != nil {
		t.Fatal(err)
	}
	head.Findings[0], head.Findings[1] = head.Findings[1], head.Findings[0]
	second, err := Compare(base, head)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(first, second) {
		t.Fatalf("comparison changed with input order: %#v != %#v", first, second)
	}
	if len(first.Changes) != 3 {
		t.Fatalf("changes = %#v", first.Changes)
	}
	if first.Changes[0].Kind != Added || first.Changes[1].Kind != Resolved || first.Changes[2].Kind != Changed {
		t.Fatalf("changes = %#v", first.Changes)
	}
	if !reflect.DeepEqual(first.BaseArtifactDigests, base.ArtifactDigests) || !reflect.DeepEqual(first.HeadArtifactDigests, head.ArtifactDigests) {
		t.Fatalf("provenance was not retained: %#v", first)
	}
}

func TestCompareRejectsIncompleteOrIncompatibleInputs(t *testing.T) {
	base := validBaseline()
	head := validBaseline()
	head.Status = "partial"
	if _, err := Compare(base, head); err == nil {
		t.Fatal("incomplete head was accepted")
	}

	head.SchemaVersion = "v2"
	if _, err := Compare(base, head); err == nil {
		t.Fatal("incompatible head was accepted")
	}
}
