package baseline

import (
	"reflect"
	"testing"
)

func TestCompareIsOrderIndependentAndClassifiesChanges(t *testing.T) {
	base := validBaseline()
	head := validBaseline()
	base.Findings = []Finding{{StableKey: "resolved", Component: "b", Version: "1"}, {StableKey: "same", Component: "a", Version: "1"}}
	head.Findings = []Finding{{StableKey: "same", Component: "a", Version: "1"}, {StableKey: "added", Component: "c", Version: "2"}}
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
	if got := []ChangeKind{first.Changes[0].Kind, first.Changes[1].Kind, first.Changes[2].Kind}; !reflect.DeepEqual(got, []ChangeKind{Added, Resolved, Unchanged}) {
		t.Fatalf("kinds = %v", got)
	}
}

func TestCompareRejectsIncompleteOrIncompatibleInputs(t *testing.T) {
	base := validBaseline()
	head := validBaseline()
	head.Status = "partial"
	if _, err := Compare(base, head); err == nil {
		t.Fatal("incomplete head was accepted")
	}
	head = validBaseline()
	head.SchemaVersion = "v2"
	if _, err := Compare(base, head); err == nil {
		t.Fatal("incompatible head was accepted")
	}
}
