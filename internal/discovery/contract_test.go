package discovery

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestCanonicalizeMakesWorkspaceOrderingStable(t *testing.T) {
	left := ProjectGraph{
		SchemaVersion:  SchemaVersion,
		DocumentType:   "project",
		RepositoryRoot: ".",
		DisplayName:    "fixture",
		Completeness:   Complete,
		Workspaces: []Workspace{
			{SchemaVersion: SchemaVersion, DocumentType: "workspace", WorkspaceID: "b", RelativePath: "services/z", Ecosystem: EcosystemPython, PackageManager: "pip", Manifests: []string{"requirements.txt", "constraints.txt"}},
			{SchemaVersion: SchemaVersion, DocumentType: "workspace", WorkspaceID: "a", RelativePath: "frontend", Ecosystem: EcosystemNPM, PackageManager: "npm", Manifests: []string{"package.json"}},
		},
		Diagnostics: []Diagnostic{
			{Code: "Z_CODE", Message: "later", Scope: "repository"},
			{Code: "A_CODE", Message: "first", Scope: "repository"},
		},
	}
	right := ProjectGraph{
		SchemaVersion:  SchemaVersion,
		DocumentType:   "project",
		RepositoryRoot: ".",
		DisplayName:    "fixture",
		Completeness:   Complete,
		Workspaces: []Workspace{
			{SchemaVersion: SchemaVersion, DocumentType: "workspace", WorkspaceID: "a", RelativePath: "frontend", Ecosystem: EcosystemNPM, PackageManager: "npm", Manifests: []string{"package.json"}},
			{SchemaVersion: SchemaVersion, DocumentType: "workspace", WorkspaceID: "b", RelativePath: "services/z", Ecosystem: EcosystemPython, PackageManager: "pip", Manifests: []string{"constraints.txt", "requirements.txt"}},
		},
		Diagnostics: []Diagnostic{
			{Code: "A_CODE", Message: "first", Scope: "repository"},
			{Code: "Z_CODE", Message: "later", Scope: "repository"},
		},
	}
	left.Canonicalize()
	right.Canonicalize()
	leftJSON, err := json.Marshal(left)
	if err != nil {
		t.Fatal(err)
	}
	rightJSON, err := json.Marshal(right)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(leftJSON, rightJSON) {
		t.Fatalf("canonical JSON differs:\n%s\n%s", leftJSON, rightJSON)
	}
}

func TestValidateProjectGraphRejectsMissingContractIdentity(t *testing.T) {
	graph := ProjectGraph{SchemaVersion: SchemaVersion, DocumentType: "project", RepositoryRoot: ".", DisplayName: "fixture"}
	graph.Workspaces = []Workspace{{SchemaVersion: SchemaVersion, DocumentType: "workspace", RelativePath: "frontend"}}
	if err := ValidateProjectGraph(graph); err == nil {
		t.Fatal("expected invalid workspace identity")
	}
}

func TestValidateProjectGraphAcceptsMinimalProject(t *testing.T) {
	graph := ProjectGraph{
		SchemaVersion:  SchemaVersion,
		DocumentType:   "project",
		RepositoryRoot: ".",
		DisplayName:    "fixture",
		Completeness:   Complete,
		Workspaces: []Workspace{{
			SchemaVersion: SchemaVersion, DocumentType: "workspace", WorkspaceID: "frontend", RelativePath: "frontend", Ecosystem: EcosystemNPM, PackageManager: "npm",
		}},
	}
	if err := ValidateProjectGraph(graph); err != nil {
		t.Fatalf("expected valid project: %v", err)
	}
}
