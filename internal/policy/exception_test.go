package policy

import (
	"testing"
	"time"
)

func testException() Exception {
	created := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	return Exception{SchemaVersion: "v1alpha", DocumentType: "exception", ID: "ex-1", FindingKey: "finding-1", Scope: "workspace-a", Rationale: "temporary acceptance", Approver: "reviewer", CreatedAt: created, ExpiresAt: created.Add(time.Hour), ReviewCondition: "review weekly"}
}
func TestExceptionExpiryAndScope(t *testing.T) {
	e := testException()
	if !e.Matches("finding-1", "workspace-a", e.CreatedAt.Add(30*time.Minute)) {
		t.Fatal("active matching exception did not match")
	}
	if e.Matches("finding-1", "workspace-a", e.ExpiresAt) || e.Matches("finding-1", "workspace-b", e.CreatedAt.Add(30*time.Minute)) {
		t.Fatal("exception boundary or scope failed")
	}
}
func TestExceptionCannotActivateBeforeCreation(t *testing.T) {
	e := testException()
	if e.ActiveAt(e.CreatedAt.Add(-time.Minute)) {
		t.Fatal("future exception activated early")
	}
}
func TestExceptionRejectsMalformedRecords(t *testing.T) {
	if err := (Exception{}).Validate(); err == nil {
		t.Fatal("malformed exception accepted")
	}
}
