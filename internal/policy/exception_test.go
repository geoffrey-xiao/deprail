package policy

import (
	"testing"
	"time"
)

func TestExceptionExpiryAndScope(t *testing.T) {
	created := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	e := Exception{SchemaVersion: "v1alpha", ID: "ex-1", FindingKey: "finding-1", Scope: "workspace-a", Rationale: "temporary acceptance", Approver: "reviewer", CreatedAt: created, ExpiresAt: created.Add(time.Hour), ReviewCondition: "review weekly"}
	if !e.Matches("finding-1", "workspace-a", created.Add(30*time.Minute)) {
		t.Fatal("active matching exception did not match")
	}
	if e.Matches("finding-1", "workspace-a", created.Add(time.Hour)) {
		t.Fatal("exception matched at expiry")
	}
	if e.Matches("finding-1", "workspace-b", created.Add(30*time.Minute)) {
		t.Fatal("exception crossed scope")
	}
}

func TestExceptionRejectsMalformedRecords(t *testing.T) {
	if err := (Exception{}).Validate(); err == nil {
		t.Fatal("malformed exception accepted")
	}
}
