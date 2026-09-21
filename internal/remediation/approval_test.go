package remediation

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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

func TestApprovalRejectsDuplicateAndEscapingAffectedPaths(t *testing.T) {
	plan, root := approvalTestPlan(t)
	plan.AffectedFiles = append(plan.AffectedFiles, plan.AffectedFiles[0])
	plan.PlanID = StablePlanID(plan)
	if _, err := NewApproval(plan, root, time.Now().UTC().Add(time.Hour)); err == nil {
		t.Fatal("expected duplicate affected path rejection")
	}

	plan, root = approvalTestPlan(t)
	plan.AffectedFiles[0].Path = "../outside.json"
	plan.PlanID = StablePlanID(plan)
	if _, err := NewApproval(plan, root, time.Now().UTC().Add(time.Hour)); err == nil {
		t.Fatal("expected escaping affected path rejection")
	}
}

func TestApprovalRejectsSymlinkedAffectedPathOutsideRoot(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink permissions vary on Windows")
	}
	plan, root := approvalTestPlan(t)
	outside := t.TempDir()
	target := filepath.Join(outside, "package.json")
	if err := os.WriteFile(target, []byte("{}"), 0o600); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(root, "package.json")
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}
	if _, err := NewApproval(plan, root, time.Now().UTC().Add(time.Hour)); err == nil {
		t.Fatal("expected symlink escape rejection")
	}
}

func TestApprovalRejectsMissingPathUnderSymlinkedDirectory(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink permissions vary on Windows")
	}
	plan, root := approvalTestPlan(t)
	outside := t.TempDir()
	link := filepath.Join(root, "nested")
	if err := os.Symlink(outside, link); err != nil {
		t.Fatal(err)
	}
	plan.AffectedFiles[0].Path = "nested/new.json"
	plan.PlanID = StablePlanID(plan)
	if _, err := NewApproval(plan, root, time.Now().UTC().Add(time.Hour)); err == nil {
		t.Fatal("expected missing path under symlink rejection")
	}
}
