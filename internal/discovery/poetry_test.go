package discovery

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestPoetryDetectorAssociatesNearestLockfile(t *testing.T) {
	root := t.TempDir()
	writePoetryFile(t, root, "pyproject.toml", "[tool.poetry]\nname = \"service\"\nversion = \"1.0.0\"\n")
	writePoetryFile(t, root, "poetry.lock", "content-hash = \"abc\"\n")
	result, err := (PoetryDetector{}).Detect(context.Background(), RepositoryView{Root: root, Paths: []string{"poetry.lock", "pyproject.toml"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Workspaces) != 1 {
		t.Fatalf("workspaces = %#v", result.Workspaces)
	}
	workspace := result.Workspaces[0]
	if workspace.PackageManager != "poetry" || workspace.Scope != "service" || workspace.Completeness != Complete {
		t.Fatalf("workspace = %#v", workspace)
	}
	if len(workspace.Lockfiles) != 1 || workspace.Lockfiles[0] != "poetry.lock" {
		t.Fatalf("lockfiles = %#v", workspace.Lockfiles)
	}
}

func TestPoetryDetectorDoesNotCrossAssociateNestedProjects(t *testing.T) {
	root := t.TempDir()
	writePoetryFile(t, root, "pyproject.toml", "[tool.poetry]\nname = \"root\"\n")
	writePoetryFile(t, root, "poetry.lock", "content-hash = \"root\"\n")
	writePoetryFile(t, root, "services/api/pyproject.toml", "[tool.poetry]\nname = \"api\"\n")
	result, err := (PoetryDetector{}).Detect(context.Background(), RepositoryView{Root: root, Paths: []string{
		"services/api/pyproject.toml", "poetry.lock", "pyproject.toml",
	}})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Workspaces) != 2 {
		t.Fatalf("workspaces = %#v", result.Workspaces)
	}
	if result.Workspaces[0].RelativePath != "." || result.Workspaces[0].Completeness != Complete {
		t.Fatalf("root workspace = %#v", result.Workspaces[0])
	}
	if result.Workspaces[1].RelativePath != "services/api" || result.Workspaces[1].Completeness != Partial {
		t.Fatalf("nested workspace = %#v", result.Workspaces[1])
	}
}

func TestPoetryDetectorMarksMissingAndMalformedInputs(t *testing.T) {
	root := t.TempDir()
	writePoetryFile(t, root, "missing/pyproject.toml", "[tool.poetry]\nname = \"missing\"\n")
	writePoetryFile(t, root, "broken/pyproject.toml", "[tool.poetry]\nname =\n")
	result, err := (PoetryDetector{}).Detect(context.Background(), RepositoryView{Root: root, Paths: []string{
		"broken/pyproject.toml", "missing/pyproject.toml",
	}})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Workspaces) != 1 || result.Workspaces[0].Completeness != Partial {
		t.Fatalf("workspaces = %#v", result.Workspaces)
	}
	if !hasDiagnostic(result.Diagnostics, "MANIFEST_INVALID", "broken/pyproject.toml") || !hasDiagnostic(result.Diagnostics, "DISCOVERY_INCOMPLETE", "missing/pyproject.toml") {
		t.Fatalf("diagnostics = %#v", result.Diagnostics)
	}
}

func TestPoetryDetectorIsInputOrderIndependent(t *testing.T) {
	root := t.TempDir()
	writePoetryFile(t, root, "a/pyproject.toml", "[tool.poetry]\nname = \"a\"\n")
	writePoetryFile(t, root, "a/poetry.lock", "content-hash = \"a\"\n")
	writePoetryFile(t, root, "b/pyproject.toml", "[tool.poetry]\nname = \"b\"\n")
	first, err := (PoetryDetector{}).Detect(context.Background(), RepositoryView{Root: root, Paths: []string{"b/pyproject.toml", "a/poetry.lock", "a/pyproject.toml"}})
	if err != nil {
		t.Fatal(err)
	}
	second, err := (PoetryDetector{}).Detect(context.Background(), RepositoryView{Root: root, Paths: []string{"a/pyproject.toml", "b/pyproject.toml", "a/poetry.lock"}})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(first, second) {
		t.Fatalf("input order changed result: %#v != %#v", first, second)
	}
}

func writePoetryFile(t *testing.T, root, relative, content string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
