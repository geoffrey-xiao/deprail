package remediation

import (
	"os/exec"
	"testing"
	"time"
)

func approvalTestPlan(t *testing.T) (Plan, string) {
	t.Helper()
	root := t.TempDir()
	if out, err := exec.Command("git", "-C", root, "init", "-q").CombinedOutput(); err != nil {
		t.Fatalf("git init: %v: %s", err, out)
	}
	plan := testPlan()
	plan.RepositoryIdentity.Root = root
	plan.RepositoryIdentity.Revision = "commit-1"
	plan.PlanID = StablePlanID(plan)
	return plan, root
}

func TestApprovalBindsPlanSourceAndExpiry(t *testing.T) {
	plan, root := approvalTestPlan(t)
	approval, err := NewApproval(plan, root, time.Now().UTC().Add(time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if err := approval.Validate(plan, root, "commit-1", time.Now().UTC()); err != nil {
		t.Fatal(err)
	}
	changed := plan
	changed.Risks[0].Details = "changed"
	if err := approval.Validate(changed, root, "commit-1", time.Now().UTC()); err == nil {
		t.Fatal("expected changed plan rejection")
	}
	copy := approval
	if err := approval.Consume(); err != nil {
		t.Fatal(err)
	}
	if err := copy.Validate(plan, root, "commit-1", time.Now().UTC()); err == nil {
		t.Fatal("expected copied approval rejection")
	}
}

func TestApprovalRejectsExpiredAndMismatchedSource(t *testing.T) {
	plan, root := approvalTestPlan(t)
	approval, err := NewApproval(plan, root, time.Now().UTC().Add(time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if err := approval.Validate(plan, t.TempDir(), "commit-1", time.Now().UTC()); err == nil {
		t.Fatal("expected source root rejection")
	}
	if err := approval.Validate(plan, root, "commit-1", time.Now().UTC().Add(2*time.Hour)); err == nil {
		t.Fatal("expected expiry rejection")
	}

}
