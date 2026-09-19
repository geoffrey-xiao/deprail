package javascript

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/geoffrey-xiao/deprail/internal/remediation"
)

func requestForFixture(t *testing.T, name string) remediation.PlanningRequest {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", "..", "..", "testdata", "fixtures", name))
	if err != nil {
		t.Fatal(err)
	}
	return remediation.PlanningRequest{
		Repository:      remediation.RepositoryIdentity{Root: root, Repository: name, Revision: "rev-1"},
		Workspace:       remediation.WorkspaceIdentity{ID: "root", Path: "."},
		Finding:         remediation.FindingIdentity{StableKey: "finding-1", CurrentVersion: "4.17.20", FixedVersions: []string{"4.17.21", "5.0.0"}},
		Component:       remediation.Component{PURL: "pkg:npm/lodash@4.17.20", Name: "lodash", Version: "4.17.20"},
		RepositoryState: "tree-1",
	}
}

func TestJavaScriptAdapterPlansDeterministicallyFromOfflineEvidence(t *testing.T) {
	adapter := Adapter{}
	request := requestForFixture(t, "npm-basic")
	assessment, err := adapter.Assess(context.Background(), request)
	if err != nil || assessment.State != remediation.AdapterSupported {
		t.Fatalf("assessment = %#v, err = %v", assessment, err)
	}
	first, err := adapter.Plan(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	if err := remediation.ValidatePlanningEvidence(first); err != nil {
		t.Fatal(err)
	}
	if len(first.Candidates) != 2 || first.Candidates[0].Version != "4.17.21" || first.Candidates[0].State != remediation.CandidateViable || first.Candidates[1].State != remediation.CandidateRejected || !first.Candidates[0].Direct {
		t.Fatalf("candidates = %#v", first.Candidates)
	}
	if first.Commands[0].Executable != "npm" || first.Commands[0].WorkingDirectory != "." || first.Commands[0].Arguments[0] != "install" {
		t.Fatalf("commands = %#v", first.Commands)
	}
	second, err := adapter.Plan(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	if first.Candidates[0].Version != second.Candidates[0].Version || first.Commands[0].Arguments[1] != second.Commands[0].Arguments[1] {
		t.Fatalf("nondeterministic plans: %#v %#v", first, second)
	}
}

func TestJavaScriptAdapterRejectsConflictingOrMalformedInputs(t *testing.T) {
	adapter := Adapter{}
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "package.json"), []byte(`{"dependencies":{"demo":"^1.0.0"}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "package-lock.json"), []byte("not json"), 0o600); err != nil {
		t.Fatal(err)
	}
	request := requestForFixture(t, "npm-basic")
	request.Repository.Root = root
	request.Component.Name = "demo"
	request.Component.Version = "1.0.0"
	request.Finding.CurrentVersion = "1.0.0"
	request.Finding.FixedVersions = []string{"1.0.1"}
	evidence, err := adapter.Plan(context.Background(), request)
	if err != nil || evidence.State != remediation.AdapterRejected {
		t.Fatalf("malformed evidence = %#v, err = %v", evidence, err)
	}
	if err := os.WriteFile(filepath.Join(root, "yarn.lock"), []byte("demo@^1.0.0:\n  version \"1.0.1\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	evidence, err = adapter.Plan(context.Background(), request)
	if err != nil || evidence.State != remediation.AdapterRejected {
		t.Fatalf("conflicting evidence = %#v, err = %v", evidence, err)
	}
}

func TestJavaScriptAdapterMarksMissingLockfileUnavailable(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "package.json"), []byte(`{"dependencies":{"demo":"^1.0.0"}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	request := requestForFixture(t, "npm-basic")
	request.Repository.Root = root
	assessment, err := (Adapter{}).Assess(context.Background(), request)
	if err != nil || assessment.State != remediation.AdapterUnavailable {
		t.Fatalf("assessment = %#v, err = %v", assessment, err)
	}
}
func TestJavaScriptAdapterRejectsSymlinkWorkspaceAndUsesRetainedRootLockfile(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink fixture requires platform support")
	}
	root := t.TempDir()
	outside := t.TempDir()
	if err := os.WriteFile(filepath.Join(outside, "package.json"), []byte(`{"dependencies":{"demo":"^1.0.0"}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(root, "linked")); err != nil {
		t.Fatal(err)
	}
	request := requestForFixture(t, "npm-basic")
	request.Repository.Root = root
	request.Workspace.Path = "linked"
	if assessment, err := (Adapter{}).Assess(context.Background(), request); err != nil || assessment.State != remediation.AdapterRejected {
		t.Fatalf("symlink assessment = %#v, err = %v", assessment, err)
	}

	if err := os.WriteFile(filepath.Join(root, "package.json"), []byte(`{"dependencies":{"demo":"^1.0.0"}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "package-lock.json"), []byte(`{"lockfileVersion":3,"packages":{}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	request.Workspace.Path = "apps/web"
	if err := os.MkdirAll(filepath.Join(root, "apps", "web"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "apps", "web", "package.json"), []byte(`{"dependencies":{"demo":"^1.0.0"}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	request.CurrentState.Lockfile = "package-lock.json"
	request.Component.Name = "demo"
	request.Component.Version = "1.0.0"
	request.Finding.CurrentVersion = "1.0.0"
	request.Finding.FixedVersions = []string{"1.0.1"}
	evidence, err := (Adapter{}).Plan(context.Background(), request)
	if err != nil || evidence.State != remediation.AdapterSupported || evidence.Commands[0].WorkingDirectory != "apps/web" {
		t.Fatalf("root lockfile evidence = %#v, err = %v", evidence, err)
	}
}

func TestJavaScriptAdapterHandlesTransitiveOwnersAndMalformedLockfiles(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "package.json"), []byte(`{"dependencies":{"parent":"^2.0.0"}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "package-lock.json"), []byte(`{"lockfileVersion":3,"packages":{}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	request := requestForFixture(t, "npm-basic")
	request.Repository.Root = root
	request.Component = remediation.Component{PURL: "pkg:npm/child@1.0.0", Name: "child", Version: "1.0.0"}
	request.Finding = remediation.FindingIdentity{StableKey: "transitive", CurrentVersion: "1.0.0", FixedVersions: []string{"1.0.1"}, DependencyPath: []string{"root", "parent", "child"}}
	evidence, err := (Adapter{}).Plan(context.Background(), request)
	if err != nil || evidence.State != remediation.AdapterSupported || evidence.Commands[0].Arguments[1] != "parent@1.0.1" {
		t.Fatalf("transitive evidence = %#v, err = %v", evidence, err)
	}

	if err := os.Remove(filepath.Join(root, "package-lock.json")); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "pnpm-lock.yaml"), []byte("not a lockfile"), 0o600); err != nil {
		t.Fatal(err)
	}
	request.CurrentState.Lockfile = "pnpm-lock.yaml"
	evidence, err = (Adapter{}).Plan(context.Background(), request)
	if err != nil || evidence.State != remediation.AdapterRejected {
		t.Fatalf("malformed pnpm evidence = %#v, err = %v", evidence, err)
	}
}
