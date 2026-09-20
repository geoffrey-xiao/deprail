package app

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/geoffrey-xiao/deprail/internal/adapter"
	"github.com/geoffrey-xiao/deprail/internal/remediation"
)

func TestPlanBuildsReadOnlyJavaScriptPlanFromReport(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", "..", "testdata", "fixtures", "npm-basic"))
	if err != nil {
		t.Fatal(err)
	}
	state, err := CurrentRepositoryState(root)
	if err != nil {
		t.Fatal(err)
	}
	report := remediation.Report{
		SchemaVersion: remediation.ReportSchemaVersion, DocumentType: remediation.ReportDocumentType, ReportID: "scan-1",
		RepositoryIdentity: remediation.RepositoryIdentity{Root: root, Repository: "npm-basic"}, ArtifactDigests: []string{"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"}, RepositoryState: state, Status: remediation.ReportComplete,
		Findings: []remediation.ReportFinding{{StableKey: "finding-1", Workspace: remediation.WorkspaceIdentity{ID: "root", Path: "."}, Component: remediation.Component{PURL: "pkg:npm/lodash@4.17.20", Name: "lodash", Version: "4.17.20"}, CurrentVersion: "4.17.20", FixedVersions: []string{"4.17.21"}}},
	}
	path := filepath.Join(t.TempDir(), "report.json")
	data, err := json.Marshal(report)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	plan, err := Plan(context.Background(), path, "finding-1", root, PlanOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if err := plan.Validate(); err != nil {
		t.Fatal(err)
	}
	if plan.Component.Name != "lodash" || len(plan.Commands) != 1 || plan.Commands[0].WorkingDirectory != "." {
		t.Fatalf("plan = %#v", plan)
	}
}

func TestPlanRejectsStaleExplicitState(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", "..", "testdata", "fixtures", "npm-basic"))
	if err != nil {
		t.Fatal(err)
	}
	report := remediation.Report{SchemaVersion: remediation.ReportSchemaVersion, DocumentType: remediation.ReportDocumentType, ReportID: "scan-1", RepositoryIdentity: remediation.RepositoryIdentity{Root: root, Repository: "npm-basic"}, ArtifactDigests: []string{"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"}, RepositoryState: "tree-1", Status: remediation.ReportComplete, Findings: []remediation.ReportFinding{{StableKey: "finding-1", Workspace: remediation.WorkspaceIdentity{ID: "root", Path: "."}, Component: remediation.Component{PURL: "pkg:npm/lodash@4.17.20", Name: "lodash", Version: "4.17.20"}, CurrentVersion: "4.17.20"}}}
	path := filepath.Join(t.TempDir(), "report.json")
	data, err := json.Marshal(report)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	_, err = Plan(context.Background(), path, "finding-1", root, PlanOptions{CurrentRepositoryState: "tree-2"})
	if err == nil {
		t.Fatal("expected stale report error")
	}
	if reportErr, ok := err.(*remediation.ReportError); !ok || reportErr.Code != remediation.ReportInputStale {
		t.Fatalf("err = %v", err)
	}
}
func TestPlanAcceptsScanJSONOutput(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", "..", "testdata", "fixtures", "npm-basic"))
	if err != nil {
		t.Fatal(err)
	}
	state, err := CurrentRepositoryState(root)
	if err != nil {
		t.Fatal(err)
	}
	scan := ScanReport{
		SchemaVersion: remediation.ReportSchemaVersion, DocumentType: remediation.ReportDocumentType, ScanID: "scan-1",
		RepositoryIdentity: remediation.RepositoryIdentity{Root: root, Repository: "npm-basic"}, RepositoryState: state,
		Status: "complete", Findings: []adapter.Finding{{Component: "pkg:npm/lodash", Version: "4.17.20", Fixed: "4.17.21", TargetID: "finding-1"}},
		ArtifactDigests: []string{"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"},
	}
	path := filepath.Join(t.TempDir(), "scan.json")
	data, err := json.Marshal(scan)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	plan, err := Plan(context.Background(), path, "finding-1", root, PlanOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if plan.FindingIdentity.StableKey != "finding-1" {
		t.Fatalf("plan = %#v", plan)
	}
}
func TestLoadPlanningReportPreservesScanComponentName(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", "..", "testdata", "fixtures", "npm-basic"))
	if err != nil {
		t.Fatal(err)
	}
	state, err := CurrentRepositoryState(root)
	if err != nil {
		t.Fatal(err)
	}
	scan := ScanReport{
		SchemaVersion: remediation.ReportSchemaVersion, DocumentType: remediation.ReportDocumentType, ScanID: "scan-scoped",
		RepositoryIdentity: remediation.RepositoryIdentity{Root: root, Repository: "npm-basic"}, RepositoryState: state,
		Status: "complete", Findings: []adapter.Finding{{
			Component: "@scope/pkg", PURL: "pkg:npm/%40scope%2Fpkg@1.2.3", Version: "1.2.3", TargetID: "finding-scoped", Aliases: []string{},
		}},
		Errors: []string{}, ArtifactDigests: []string{"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"},
	}
	path := filepath.Join(t.TempDir(), "scan.json")
	data, err := json.Marshal(scan)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	report, err := loadPlanningReport(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := report.Findings[0].Component.Name; got != "@scope/pkg" {
		t.Fatalf("component name = %q, want %q", got, "@scope/pkg")
	}
	if got := report.Findings[0].Component.PURL; got != "pkg:npm/%40scope%2Fpkg@1.2.3" {
		t.Fatalf("component PURL = %q", got)
	}
}
