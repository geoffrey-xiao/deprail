package app

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/geoffrey-xiao/deprail/internal/remediation"
)

func TestPlanBuildsReadOnlyJavaScriptPlanFromReport(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", "..", "testdata", "fixtures", "npm-basic"))
	if err != nil {
		t.Fatal(err)
	}
	report := remediation.Report{
		SchemaVersion: remediation.ReportSchemaVersion, DocumentType: remediation.ReportDocumentType, ReportID: "scan-1",
		RepositoryIdentity: remediation.RepositoryIdentity{Root: root, Repository: "npm-basic"}, ArtifactDigests: []string{"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"}, RepositoryState: "tree-1", Status: remediation.ReportComplete,
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
