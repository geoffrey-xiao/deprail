package presenter

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/geoffrey-xiao/deprail/internal/remediation"
)

func presenterPlan(t *testing.T) remediation.Plan {
	t.Helper()
	plan := remediation.Plan{SchemaVersion: remediation.SchemaVersion, CreatedFrom: remediation.CreatedFrom{ReportDigest: strings.Repeat("a", 64), SourceScanID: "scan-1", Scanner: "osv-scanner", RepositoryState: "tree-1"}, RepositoryIdentity: remediation.RepositoryIdentity{Root: t.TempDir(), Repository: "repo"}, WorkspaceIdentity: remediation.WorkspaceIdentity{ID: "root", Path: "."}, FindingIdentity: remediation.FindingIdentity{StableKey: "finding-1", CurrentVersion: "1.0.0"}, Component: remediation.Component{PURL: "pkg:npm/example@1.0.0", Name: "example", Version: "1.0.0"}, Candidates: []remediation.Candidate{{ID: "1.1.0", Version: "1.1.0", State: remediation.CandidateViable}}, AffectedFiles: []remediation.AffectedFile{{Path: "package.json", Kind: "manifest", Effect: "dependency constraint"}}, Commands: []remediation.Command{{Executable: "npm", Arguments: []string{"install", "example@1.1.0"}, WorkingDirectory: "."}}, Risks: []remediation.Risk{{Code: "UNKNOWN_RUNTIME", Severity: "medium", Detected: true}}, Assumptions: []remediation.Assumption{}, Verification: []remediation.Verification{}, Provenance: remediation.Provenance{ArtifactDigests: []string{}, Sources: []string{"test"}}, Rollback: remediation.Rollback{Steps: []string{"restore package.json"}}}
	plan.PlanID = remediation.StablePlanID(plan)
	return plan
}

func TestWritePlanTerminalIncludesSafetySections(t *testing.T) {
	plan := presenterPlan(t)
	var out bytes.Buffer
	if err := WritePlanTerminal(&out, plan); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"Read-only: yes", "Recommendation:", "Risks:", "Affected files:", "Future commands:", "Verification:", "Provenance:"} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("output missing %q: %s", want, out.String())
		}
	}
}

func TestWritePlanJSONIsCanonicalAndValid(t *testing.T) {
	left := presenterPlan(t)
	right := left
	var first, second bytes.Buffer
	if err := WritePlanJSON(&first, left); err != nil {
		t.Fatal(err)
	}
	if err := WritePlanJSON(&second, right); err != nil {
		t.Fatal(err)
	}
	if first.String() != second.String() {
		t.Fatal("equivalent plans serialized differently")
	}
	var decoded remediation.Plan
	if err := json.Unmarshal(first.Bytes(), &decoded); err != nil {
		t.Fatal(err)
	}
	if err := decoded.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestWritePlanAtomicRejectsRepositoryAndDoesNotOverwrite(t *testing.T) {
	plan := presenterPlan(t)
	inside := filepath.Join(plan.RepositoryIdentity.Root, "plan.json")
	if err := WritePlanAtomic(inside, plan); err == nil {
		t.Fatal("expected repository-root output rejection")
	}
	outside := filepath.Join(t.TempDir(), "plan.json")
	if err := WritePlanAtomic(outside, plan); err != nil {
		t.Fatal(err)
	}
	if err := WritePlanAtomic(outside, plan); err == nil {
		t.Fatal("expected overwrite rejection")
	}
	if _, err := os.Stat(outside); err != nil {
		t.Fatal(err)
	}
}
func TestWritePlanAtomicRejectsTraversalAndUnresolvedSymlinkParent(t *testing.T) {
	plan := presenterPlan(t)
	traversal := filepath.Join(t.TempDir(), "child") + string(os.PathSeparator) + ".." + string(os.PathSeparator) + "plan.json"
	if err := WritePlanAtomic(traversal, plan); err == nil {
		t.Fatal("expected traversal rejection")
	}
	external := t.TempDir()
	link := filepath.Join(external, "link")
	if err := os.Symlink(plan.RepositoryIdentity.Root, link); err != nil {
		t.Fatal(err)
	}
	output := filepath.Join(external, "link", "new", "plan.json")
	if err := WritePlanAtomic(output, plan); err == nil {
		t.Fatal("expected unresolved symlink-parent rejection")
	}
}
