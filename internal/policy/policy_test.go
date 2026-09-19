package policy

import (
	"testing"
	"time"

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
	if err != nil || got.Outcome != Unusable {
		t.Fatalf("incomplete decision = %#v, err = %v", got, err)
	}
}

func TestEvaluateRejectsUnknownSeverity(t *testing.T) {
	if _, err := Evaluate(Policy{MinimumSeverity: "urgent"}, Input{Complete: true}); err == nil {
		t.Fatal("invalid severity accepted")
	}
}
func TestEvaluateHonorsActiveScopedException(t *testing.T) {
	now := time.Date(2026, 1, 1, 1, 0, 0, 0, time.UTC)
	comparison := baseline.Comparison{Changes: []baseline.Change{{Kind: baseline.Added, Key: "new-key", Head: &baseline.Finding{Severity: "critical"}}}}
	exception := Exception{SchemaVersion: "v1alpha", DocumentType: "exception", ID: "ex-1", FindingKey: "new-key", Scope: "workspace-a", Rationale: "temporary", Approver: "reviewer", CreatedAt: now.Add(-time.Hour), ExpiresAt: now.Add(time.Hour), ReviewCondition: "review"}
	got, err := Evaluate(Policy{BlockOnNew: true}, Input{Comparison: comparison, Complete: true, Scope: "workspace-a", Now: now, Exceptions: []Exception{exception}})
	if err != nil || got.Outcome != Pass {
		t.Fatalf("decision = %#v, err = %v", got, err)
	}
}
