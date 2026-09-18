package discovery

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestNPMDetectorFindsAuthoritativePackageLock(t *testing.T) {
	root := t.TempDir()
	writeNPMFile(t, root, "package.json", `{"name":"fixture-root"}`)
	writeNPMFile(t, root, "package-lock.json", `{ "lockfileVersion": 3 }`)

	result, err := (NPMDetector{}).Detect(context.Background(), RepositoryView{Root: root, Paths: []string{"package-lock.json", "package.json"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Workspaces) != 1 {
		t.Fatalf("workspaces = %d, want 1", len(result.Workspaces))
	}
	workspace := result.Workspaces[0]
	if workspace.Completeness != Complete || workspace.RelativePath != "." || workspace.WorkspaceID != "." {
		t.Fatalf("workspace = %#v", workspace)
	}
	if len(workspace.Lockfiles) != 1 || workspace.Lockfiles[0] != "package-lock.json" {
		t.Fatalf("lockfiles = %#v", workspace.Lockfiles)
	}
	if workspace.Scope != "fixture-root" {
		t.Fatalf("scope = %q, want fixture-root", workspace.Scope)
	}
}

func TestNPMDetectorMarksMissingLockfileIncomplete(t *testing.T) {
	root := t.TempDir()
	writeNPMFile(t, root, "frontend/package.json", `{"name":"frontend"}`)
	result, err := (NPMDetector{}).Detect(context.Background(), RepositoryView{Root: root, Paths: []string{"frontend/package.json"}})
	if err != nil {
		t.Fatal(err)
	}
	workspace := result.Workspaces[0]
	if workspace.Completeness != Partial {
		t.Fatalf("completeness = %q, want partial", workspace.Completeness)
	}
	if !hasDiagnostic(workspace.Diagnostics, "DISCOVERY_INCOMPLETE", "frontend/package.json") {
		t.Fatalf("diagnostics = %#v", workspace.Diagnostics)
	}
}

func TestNPMDetectorRejectsMalformedManifest(t *testing.T) {
	root := t.TempDir()
	writeNPMFile(t, root, "package.json", `{not-json`)
	result, err := (NPMDetector{}).Detect(context.Background(), RepositoryView{Root: root, Paths: []string{"package.json"}})
	if err != nil {
		t.Fatal(err)
	}
	workspace := result.Workspaces[0]
	if workspace.Completeness != Failed || !hasDiagnostic(workspace.Diagnostics, "MANIFEST_INVALID", "package.json") {
		t.Fatalf("workspace = %#v", workspace)
	}
}

func TestNPMDetectorUsesInputOrderIndependently(t *testing.T) {
	root := t.TempDir()
	writeNPMFile(t, root, "b/package.json", `{"name":"b"}`)
	writeNPMFile(t, root, "b/package-lock.json", `{}`)
	writeNPMFile(t, root, "a/package.json", `{"name":"a"}`)
	writeNPMFile(t, root, "a/npm-shrinkwrap.json", `{}`)

	first, err := (NPMDetector{}).Detect(context.Background(), RepositoryView{Root: root, Paths: []string{"b/package.json", "a/npm-shrinkwrap.json", "a/package.json", "b/package-lock.json"}})
	if err != nil {
		t.Fatal(err)
	}
	second, err := (NPMDetector{}).Detect(context.Background(), RepositoryView{Root: root, Paths: []string{"a/package.json", "b/package-lock.json", "b/package.json", "a/npm-shrinkwrap.json"}})
	if err != nil {
		t.Fatal(err)
	}
	firstGraph := ProjectGraph{SchemaVersion: SchemaVersion, DocumentType: "project", RepositoryRoot: ".", DisplayName: "fixture", Workspaces: first.Workspaces, Completeness: Complete}
	secondGraph := ProjectGraph{SchemaVersion: SchemaVersion, DocumentType: "project", RepositoryRoot: ".", DisplayName: "fixture", Workspaces: second.Workspaces, Completeness: Complete}
	firstGraph.Canonicalize()
	secondGraph.Canonicalize()
	firstJSON, _ := json.Marshal(firstGraph)
	secondJSON, _ := json.Marshal(secondGraph)
	if string(firstJSON) != string(secondJSON) {
		t.Fatalf("detector output depends on input order: %s != %s", firstJSON, secondJSON)
	}
}

func writeNPMFile(t *testing.T, root, relative, contents string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}
