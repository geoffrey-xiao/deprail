package discovery

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestRequirementsDetectorFindsIndependentProjects(t *testing.T) {
	root := t.TempDir()
	writePythonFile(t, root, "services/api/requirements.txt", "requests==2.31.0\n")
	result, err := (RequirementsDetector{}).Detect(context.Background(), RepositoryView{Root: root, Paths: []string{"services/api/requirements.txt"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Workspaces) != 1 {
		t.Fatalf("workspaces = %#v", result.Workspaces)
	}
	workspace := result.Workspaces[0]
	if workspace.RelativePath != "services/api" || workspace.PackageManager != "pip" || workspace.Completeness != Complete {
		t.Fatalf("workspace = %#v", workspace)
	}
}

func TestUVDetectorAssociatesPyprojectAndLock(t *testing.T) {
	root := t.TempDir()
	writePythonFile(t, root, "pyproject.toml", "[project]\nname = \"sample\"\nversion = \"1.0.0\"\n")
	writePythonFile(t, root, "uv.lock", "version = 1\n")
	result, err := (UVDetector{}).Detect(context.Background(), RepositoryView{Root: root, Paths: []string{"uv.lock", "pyproject.toml"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Workspaces) != 1 {
		t.Fatalf("workspaces = %#v", result.Workspaces)
	}
	workspace := result.Workspaces[0]
	if workspace.Scope != "sample" || workspace.Lockfiles[0] != "uv.lock" || workspace.Completeness != Complete {
		t.Fatalf("workspace = %#v", workspace)
	}
}

func TestUVDetectorMarksMissingLockIncomplete(t *testing.T) {
	root := t.TempDir()
	writePythonFile(t, root, "pyproject.toml", "[project]\nname = \"sample\"\n")
	result, err := (UVDetector{}).Detect(context.Background(), RepositoryView{Root: root, Paths: []string{"pyproject.toml"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Workspaces) != 1 || result.Workspaces[0].Completeness != Partial {
		t.Fatalf("workspaces = %#v", result.Workspaces)
	}
	if !hasDiagnostic(result.Workspaces[0].Diagnostics, "DISCOVERY_INCOMPLETE", "pyproject.toml") {
		t.Fatalf("diagnostics = %#v", result.Workspaces[0].Diagnostics)
	}
}

func TestUVDetectorRejectsMalformedPyprojectAndIgnoresPoetryMetadata(t *testing.T) {
	root := t.TempDir()
	writePythonFile(t, root, "broken/pyproject.toml", "\x00")
	writePythonFile(t, root, "poetry/pyproject.toml", "[tool.poetry]\nname = \"poetry-project\"\n")
	result, err := (UVDetector{}).Detect(context.Background(), RepositoryView{Root: root, Paths: []string{"poetry/pyproject.toml", "broken/pyproject.toml"}})
	if err != nil {
		t.Fatal(err)
	}
	if !hasDiagnostic(result.Diagnostics, "MANIFEST_INVALID", "broken/pyproject.toml") || len(result.Workspaces) != 0 {
		t.Fatalf("result = %#v", result)
	}
}

func TestRequirementsDetectorIsInputOrderIndependent(t *testing.T) {
	root := t.TempDir()
	writePythonFile(t, root, "a/requirements.txt", "a==1\n")
	writePythonFile(t, root, "b/requirements-dev.txt", "b==1\n")
	first, err := (RequirementsDetector{}).Detect(context.Background(), RepositoryView{Root: root, Paths: []string{"b/requirements-dev.txt", "a/requirements.txt"}})
	if err != nil {
		t.Fatal(err)
	}
	second, err := (RequirementsDetector{}).Detect(context.Background(), RepositoryView{Root: root, Paths: []string{"a/requirements.txt", "b/requirements-dev.txt"}})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(first, second) {
		t.Fatalf("input order changed result: %#v != %#v", first, second)
	}
}

func writePythonFile(t *testing.T, root, relative, content string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
