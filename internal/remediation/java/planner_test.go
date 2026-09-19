package java

import (
	"context"
	"github.com/geoffrey-xiao/deprail/internal/remediation"
	"os"
	"path/filepath"
	"testing"
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
func TestJavaPlannerRejectsEscapingManifest(t *testing.T) {
	root := t.TempDir()
	outside := filepath.Join(t.TempDir(), "pom.xml")
	if err := os.WriteFile(outside, []byte(`<project><dependencies/></project>`), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(root, "pom.xml")); err != nil {
		t.Fatal(err)
	}
	r := javaRequest(t)
	r.Repository.Root = root
	evidence, err := (Adapter{}).Plan(context.Background(), r)
	if err != nil {
		t.Fatal(err)
	}
	if evidence.State != remediation.AdapterRejected {
		t.Fatalf("evidence = %#v", evidence)
	}
}

func TestGradlePlannerWithholdsUnverifiedCommand(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "build.gradle"), []byte(`implementation "com.google.guava:guava:32.1.2"`), 0o600); err != nil {
		t.Fatal(err)
	}
	r := javaRequest(t)
	r.Repository.Root = root
	evidence, err := (Adapter{}).Plan(context.Background(), r)
	if err != nil {
		t.Fatal(err)
	}
	if evidence.State != remediation.AdapterUnknown || len(evidence.Commands) != 0 {
		t.Fatalf("evidence = %#v", evidence)
	}
}
