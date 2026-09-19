package policy

import (
	"testing"

	"github.com/geoffrey-xiao/deprail/internal/baseline"
)

func TestEvaluateBlocksNewFindingsDeterministically(t *testing.T) {
	comparison := baseline.Comparison{Changes: []baseline.Change{{Kind: baseline.Added, Key: "new-key"}}}
	got, err := Evaluate(Policy{BlockOnNew: true}, Input{Comparison: comparison, Complete: true})
	if err != nil {
		t.Fatal(err)
	}
	if got.Outcome != Block || len(got.Reasons) != 1 {
		t.Fatalf("decision = %#v", got)
	}
}

func TestEvaluateEnforcesSeverityAndCompleteness(t *testing.T) {
	comparison := baseline.Comparison{Changes: []baseline.Change{{Kind: baseline.Added, Key: "critical", Head: &baseline.Finding{Severity: "critical"}}}}
	got, err := Evaluate(Policy{MinimumSeverity: "high"}, Input{Comparison: comparison, Complete: true})
	if err != nil || got.Outcome != Block {
		t.Fatalf("decision = %#v, err = %v", got, err)
	}
	got, err = Evaluate(Policy{}, Input{Complete: false})
	if err != nil || got.Outcome != Block {
		t.Fatalf("incomplete decision = %#v, err = %v", got, err)
	}
}

func TestEvaluateRejectsUnknownSeverity(t *testing.T) {
	if _, err := Evaluate(Policy{MinimumSeverity: "urgent"}, Input{Complete: true}); err == nil {
		t.Fatal("invalid severity accepted")
	}
}
