package python

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/geoffrey-xiao/deprail/internal/remediation"
)

func request(t *testing.T, fixture, name, version string) remediation.PlanningRequest {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", "..", "..", "testdata", "fixtures", fixture))
	if err != nil {
		t.Fatal(err)
	}
	return remediation.PlanningRequest{Repository: remediation.RepositoryIdentity{Root: root, Repository: fixture}, Workspace: remediation.WorkspaceIdentity{ID: "root", Path: "."}, Finding: remediation.FindingIdentity{StableKey: "python-finding", CurrentVersion: version, FixedVersions: []string{"2.31.1", "3.0.0"}}, Component: remediation.Component{PURL: "pkg:pypi/" + name + "@" + version, Name: name, Version: version}, RepositoryState: "tree-1"}
}

func TestPythonAdapterPlansUVFixture(t *testing.T) {
	r := request(t, "python-uv", "requests", "2.31.0")
	evidence, err := (Adapter{}).Plan(context.Background(), r)
	if err != nil {
		t.Fatal(err)
	}
	if err := remediation.ValidatePlanningEvidence(evidence); err != nil {
		t.Fatal(err)
	}
	if evidence.State != remediation.AdapterUnknown || len(evidence.Candidates) != 2 || len(evidence.Commands) != 0 || evidence.Candidates[0].State != remediation.CandidateRejected {
		t.Fatalf("evidence = %#v", evidence)
	}
}

func TestPythonAdapterRejectsAmbiguousOwner(t *testing.T) {
	r := request(t, "python-requirements", "transitive", "1.0.0")
	r.Finding.DependencyPath = []string{"root", "transitive"}
	evidence, err := (Adapter{}).Plan(context.Background(), r)
	if err != nil {
		t.Fatal(err)
	}
	if evidence.State != remediation.AdapterUnknown || len(evidence.Commands) != 0 {
		t.Fatalf("evidence = %#v", evidence)
	}
}
