package policy

import (
	"testing"

	"github.com/geoffrey-xiao/deprail/internal/baseline"
)

func TestEvaluateBlocksNewFindingsDeterministically(t *testing.T) {
	comparison := baseline.Comparison{Changes: []baseline.Change{{Kind: baseline.Added, Key: "new-key"}}}
	got, err := Evaluate(Policy{BlockOnNew: true}, comparison)
	if err != nil {
		t.Fatal(err)
	}
	if got.Outcome != Block || len(got.Reasons) != 1 {
		t.Fatalf("decision = %#v", got)
	}
}

func TestEvaluateRejectsUnknownSeverity(t *testing.T) {
	if _, err := Evaluate(Policy{MinimumSeverity: "urgent"}, baseline.Comparison{}); err == nil {
		t.Fatal("invalid severity accepted")
	}
}
