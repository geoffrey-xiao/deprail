package discovery

import (
	"reflect"
	"testing"
)

func TestFinalizeDetectionDiagnosesConflictingLockfiles(t *testing.T) {
	result := FinalizeDetection(DetectionResult{Workspaces: []Workspace{{
		SchemaVersion: SchemaVersion, DocumentType: "workspace", WorkspaceID: ".", RelativePath: ".",
		Ecosystem: EcosystemNPM, PackageManager: "npm", Manifests: []string{"package.json"},
		Lockfiles: []string{"npm-shrinkwrap.json", "package-lock.json"}, Completeness: Complete,
	}}})
	workspace := result.Workspaces[0]
	if workspace.Completeness != Partial {
		t.Fatalf("completeness = %q, want partial", workspace.Completeness)
	}
	if !containsDiagnostic(workspace.Diagnostics, "LOCKFILE_CONFLICT", ".") {
		t.Fatalf("workspace diagnostics = %#v", workspace.Diagnostics)
	}
	if !containsDiagnostic(result.Diagnostics, "LOCKFILE_CONFLICT", ".") {
		t.Fatalf("result diagnostics = %#v", result.Diagnostics)
	}
}

func TestCompletenessForWorkspacesPreservesPartialSuccess(t *testing.T) {
	workspaces := []Workspace{{Completeness: Complete}, {Completeness: Partial}}
	if got := CompletenessForWorkspaces(workspaces); got != Partial {
		t.Fatalf("completeness = %q, want partial", got)
	}
	workspaces = append(workspaces, Workspace{Completeness: Failed})
	if got := CompletenessForWorkspaces(workspaces); got != Failed {
		t.Fatalf("completeness = %q, want failed", got)
	}
}

func TestFinalizeDetectionKeepsSuccessfulWorkspaceWhenAnotherFails(t *testing.T) {
	input := DetectionResult{Workspaces: []Workspace{
		{SchemaVersion: SchemaVersion, DocumentType: "workspace", WorkspaceID: "broken", RelativePath: "broken", Ecosystem: EcosystemPython, PackageManager: "uv", Completeness: Failed, Diagnostics: []Diagnostic{{Code: "MANIFEST_INVALID", Scope: "broken/pyproject.toml"}}},
		{SchemaVersion: SchemaVersion, DocumentType: "workspace", WorkspaceID: "ok", RelativePath: "ok", Ecosystem: EcosystemPython, PackageManager: "pip", Completeness: Complete},
	}}
	got := FinalizeDetection(input)
	if len(got.Workspaces) != 2 || got.Workspaces[0].RelativePath != "broken" || got.Workspaces[1].RelativePath != "ok" {
		t.Fatalf("workspaces = %#v", got.Workspaces)
	}
	if len(got.Diagnostics) != 1 || got.Diagnostics[0].Code != "MANIFEST_INVALID" {
		t.Fatalf("diagnostics = %#v", got.Diagnostics)
	}
}

func TestFinalizeDetectionIsOrderIndependent(t *testing.T) {
	left := FinalizeDetection(DetectionResult{Workspaces: []Workspace{
		{WorkspaceID: "b", RelativePath: "b", Ecosystem: EcosystemNPM, PackageManager: "npm", Lockfiles: []string{"z", "a"}, Completeness: Complete},
		{WorkspaceID: "a", RelativePath: "a", Ecosystem: EcosystemNPM, PackageManager: "npm", Completeness: Complete},
	}})
	right := FinalizeDetection(DetectionResult{Workspaces: []Workspace{
		{WorkspaceID: "a", RelativePath: "a", Ecosystem: EcosystemNPM, PackageManager: "npm", Completeness: Complete},
		{WorkspaceID: "b", RelativePath: "b", Ecosystem: EcosystemNPM, PackageManager: "npm", Lockfiles: []string{"a", "z"}, Completeness: Complete},
	}})
	if !reflect.DeepEqual(left, right) {
		t.Fatalf("finalization depends on input order: %#v != %#v", left, right)
	}
}
