package discovery

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestPNPMDetectorMapsDeclaredWorkspacesAndRootLockfile(t *testing.T) {
	root := t.TempDir()
	writeJSFile(t, root, "package.json", `{"name":"root","packageManager":"pnpm@9.0.0"}`)
	writeJSFile(t, root, "pnpm-workspace.yaml", "packages:\n  - 'apps/*'\n")
	writeJSFile(t, root, "pnpm-lock.yaml", "lockfileVersion: '9.0'\n")
	writeJSFile(t, root, "apps/web/package.json", `{"name":"web"}`)
	result, err := (PNPMDetector{}).Detect(context.Background(), RepositoryView{Root: root, Paths: []string{
		"apps/web/package.json", "pnpm-lock.yaml", "package.json", "pnpm-workspace.yaml",
	}})
	if err != nil {
		t.Fatal(err)
	}
	reordered, err := (PNPMDetector{}).Detect(context.Background(), RepositoryView{Root: root, Paths: []string{
		"package.json", "pnpm-workspace.yaml", "apps/web/package.json", "pnpm-lock.yaml",
	}})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(result, reordered) {
		t.Fatalf("input order changed result:\nfirst=%#v\nsecond=%#v", result, reordered)
	}
	if len(result.Workspaces) != 2 {
		t.Fatalf("workspaces = %#v", result.Workspaces)
	}
	if result.Workspaces[0].RelativePath != "." || result.Workspaces[1].RelativePath != "apps/web" {
		t.Fatalf("workspace paths = %#v", result.Workspaces)
	}
	for _, workspace := range result.Workspaces {
		if workspace.Completeness != Complete || len(workspace.Lockfiles) != 1 || workspace.Lockfiles[0] != "pnpm-lock.yaml" {
			t.Fatalf("workspace = %#v", workspace)
		}
	}
}

func TestYarnDetectorSupportsObjectWorkspaceDeclaration(t *testing.T) {
	root := t.TempDir()
	writeJSFile(t, root, "package.json", `{"name":"root","workspaces":{"packages":["packages/*"]}}`)
	writeJSFile(t, root, "yarn.lock", "# yarn lockfile v1\n")
	writeJSFile(t, root, "packages/ui/package.json", `{"name":"ui"}`)
	result, err := (YarnDetector{}).Detect(context.Background(), RepositoryView{Root: root, Paths: []string{
		"packages/ui/package.json", "package.json", "yarn.lock",
	}})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Workspaces) != 2 || result.Workspaces[1].Scope != "ui" {
		t.Fatalf("workspaces = %#v", result.Workspaces)
	}
}

func TestPNPMDetectorMarksMissingLockfileIncomplete(t *testing.T) {
	root := t.TempDir()
	writeJSFile(t, root, "package.json", `{"name":"root","packageManager":"pnpm@9.0.0"}`)
	writeJSFile(t, root, "pnpm-workspace.yaml", "packages:\n  - 'packages/*'\n")
	result, err := (PNPMDetector{}).Detect(context.Background(), RepositoryView{Root: root, Paths: []string{"package.json", "pnpm-workspace.yaml"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Workspaces) != 1 || result.Workspaces[0].Completeness != Partial {
		t.Fatalf("workspaces = %#v", result.Workspaces)
	}
	if !hasDiagnostic(result.Workspaces[0].Diagnostics, "DISCOVERY_INCOMPLETE", "package.json") {
		t.Fatalf("diagnostics = %#v", result.Workspaces[0].Diagnostics)
	}
}

func TestYarnDetectorRejectsMalformedWorkspaceManifest(t *testing.T) {
	root := t.TempDir()
	writeJSFile(t, root, "package.json", `{not-json`)
	result, err := (YarnDetector{}).Detect(context.Background(), RepositoryView{Root: root, Paths: []string{"package.json", "yarn.lock"}})
	if err != nil {
		t.Fatal(err)
	}
	if !hasDiagnostic(result.Diagnostics, "MANIFEST_INVALID", "package.json") {
		t.Fatalf("diagnostics = %#v", result.Diagnostics)
	}
}

func writeJSFile(t *testing.T, root, relative, content string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
