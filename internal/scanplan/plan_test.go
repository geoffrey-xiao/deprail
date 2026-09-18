package scanplan

import (
	"testing"

	"github.com/geoffrey-xiao/deprail/internal/discovery"
)

func TestBuildOrdersCompleteTargetsAndSkipsIncomplete(t *testing.T) {
	graph := discovery.ProjectGraph{Workspaces: []discovery.Workspace{
		{WorkspaceID: "z", RelativePath: "services/z", Ecosystem: discovery.EcosystemMaven, Completeness: discovery.Complete, Manifests: []string{"services/z/pom.xml"}},
		{WorkspaceID: "broken", RelativePath: "frontend", Ecosystem: discovery.EcosystemNPM, Completeness: discovery.Partial, Manifests: []string{"frontend/package.json"}},
		{WorkspaceID: "a", RelativePath: "services/a", Ecosystem: discovery.EcosystemPython, Completeness: discovery.Complete, Manifests: []string{"services/a/pyproject.toml"}, Lockfiles: []string{"services/a/uv.lock"}},
	}}
	units, err := Build(graph)
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	if len(units) != 2 || units[0].Target.RelativePath != "services/a" || units[1].Target.RelativePath != "services/z" {
		t.Fatalf("units = %#v", units)
	}
	if units[0].Target.PackageFiles[0] != "services/a/pyproject.toml" || units[0].Target.PackageFiles[1] != "services/a/uv.lock" {
		t.Fatalf("files = %#v", units[0].Target.PackageFiles)
	}
}

func TestBuildRejectsGraphWithoutCompleteTargets(t *testing.T) {
	_, err := Build(discovery.ProjectGraph{Workspaces: []discovery.Workspace{{WorkspaceID: "broken", RelativePath: "broken", Completeness: discovery.Failed}}})
	if err == nil {
		t.Fatal("expected no-target error")
	}
}
