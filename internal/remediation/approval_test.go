package remediation

import (
	"testing"
	"time"
)

func TestApprovalBindsPlanSourceAndExpiry(t *testing.T) {
	plan := testPlan()
	plan.RepositoryIdentity.Revision = "commit-1"
	plan.PlanID = StablePlanID(plan)
	approval, err := NewApproval(plan, "/repo", time.Now().UTC().Add(time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if err := approval.Validate(plan, "/repo", "commit-1", time.Now().UTC()); err != nil {
		t.Fatal(err)
	}
	changed := plan
	changed.Risks[0].Details = "changed"
	if err := approval.Validate(changed, "/repo", "commit-1", time.Now().UTC()); err == nil {
		t.Fatal("expected changed plan rejection")
	}
	if err := approval.Consume(); err != nil {
		t.Fatal(err)
	}
	if err := approval.Validate(plan, "/repo", "commit-1", time.Now().UTC()); err == nil {
		t.Fatal("expected single-use rejection")
	}
}

func TestApprovalRejectsExpiredAndMismatchedSource(t *testing.T) {
	plan := testPlan()
	plan.RepositoryIdentity.Revision = "commit-1"
	plan.PlanID = StablePlanID(plan)
	approval, err := NewApproval(plan, "/repo", time.Now().UTC().Add(time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if err := approval.Validate(plan, "/other", "commit-1", time.Now().UTC()); err == nil {
		t.Fatal("expected source root rejection")
	}
	if err := approval.Validate(plan, "/repo", "commit-1", time.Now().UTC().Add(2*time.Hour)); err == nil {
		t.Fatal("expected expiry rejection")
	}
}
