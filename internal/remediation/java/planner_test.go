package java

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/geoffrey-xiao/deprail/internal/remediation"
)

func javaRequest(t *testing.T) remediation.PlanningRequest {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", "..", "..", "testdata", "fixtures", "java-maven"))
	if err != nil {
		t.Fatal(err)
	}
	return remediation.PlanningRequest{Repository: remediation.RepositoryIdentity{Root: root, Repository: "java-maven"}, Workspace: remediation.WorkspaceIdentity{ID: "root", Path: "."}, Finding: remediation.FindingIdentity{StableKey: "java-finding", CurrentVersion: "32.1.2-jre", FixedVersions: []string{"32.1.3-jre", "33.0.0-jre"}}, Component: remediation.Component{PURL: "pkg:maven/com.google.guava/guava@32.1.2-jre", Name: "guava", Version: "32.1.2-jre"}, RepositoryState: "tree-1"}
}

func TestMavenAdapterPlansDeterministically(t *testing.T) {
	evidence, err := (Adapter{}).Plan(context.Background(), javaRequest(t))
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

func TestJavaAdapterRejectsMissingManifest(t *testing.T) {
	r := javaRequest(t)
	r.Repository.Root = t.TempDir()
	assessment, err := (Adapter{}).Assess(context.Background(), r)
	if err != nil {
		t.Fatal(err)
	}
	if assessment.State != remediation.AdapterUnavailable {
		t.Fatalf("assessment = %#v", assessment)
	}
}
