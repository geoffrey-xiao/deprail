package discovery

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestMavenDetectorFindsProjectIdentity(t *testing.T) {
	root := t.TempDir()
	writeJavaFile(t, root, "pom.xml", `<project><modelVersion>4.0.0</modelVersion><groupId>com.example</groupId><artifactId>app</artifactId><version>1.0.0</version></project>`)
	result, err := (MavenDetector{}).Detect(context.Background(), RepositoryView{Root: root, Paths: []string{"pom.xml"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Workspaces) != 1 {
		t.Fatalf("workspaces = %#v", result.Workspaces)
	}
	workspace := result.Workspaces[0]
	if workspace.Scope != "com.example:app" || workspace.PackageManager != "maven" || workspace.Completeness != Complete {
		t.Fatalf("workspace = %#v", workspace)
	}
}

func TestMavenDetectorKeepsNestedModulesIndependent(t *testing.T) {
	root := t.TempDir()
	writeJavaFile(t, root, "pom.xml", `<project><artifactId>root</artifactId></project>`)
	writeJavaFile(t, root, "service/pom.xml", `<project><parent><groupId>com.example</groupId></parent><artifactId>service</artifactId></project>`)
	result, err := (MavenDetector{}).Detect(context.Background(), RepositoryView{Root: root, Paths: []string{"service/pom.xml", "pom.xml"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Workspaces) != 2 || result.Workspaces[0].RelativePath != "." || result.Workspaces[1].RelativePath != "service" {
		t.Fatalf("workspaces = %#v", result.Workspaces)
	}
	if result.Workspaces[1].Scope != "com.example:service" {
		t.Fatalf("nested workspace = %#v", result.Workspaces[1])
	}
}

func TestGradleDetectorMapsSettingsAndBuildFiles(t *testing.T) {
	root := t.TempDir()
	writeJavaFile(t, root, "settings.gradle", "rootProject.name = 'sample'\ninclude ':app'")
	writeJavaFile(t, root, "build.gradle", "plugins { id 'java' }\n")
	writeJavaFile(t, root, "app/build.gradle", "plugins { id 'java' }\n")
	writeJavaFile(t, root, "gradle.lockfile", "empty = false\n")
	result, err := (GradleDetector{}).Detect(context.Background(), RepositoryView{Root: root, Paths: []string{
		"app/build.gradle", "gradle.lockfile", "build.gradle", "settings.gradle",
	}})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Workspaces) != 2 {
		t.Fatalf("workspaces = %#v", result.Workspaces)
	}
	if result.Workspaces[0].RelativePath != "." || result.Workspaces[0].Scope != "sample" || len(result.Workspaces[0].Manifests) != 2 {
		t.Fatalf("root workspace = %#v", result.Workspaces[0])
	}
	if result.Workspaces[1].RelativePath != "app" {
		t.Fatalf("nested workspace = %#v", result.Workspaces[1])
	}
}

func TestGradleDetectorRejectsHostileBuildInput(t *testing.T) {
	root := t.TempDir()
	writeJavaFile(t, root, "settings.gradle", "\x00")
	result, err := (GradleDetector{}).Detect(context.Background(), RepositoryView{Root: root, Paths: []string{"settings.gradle"}})
	if err != nil {
		t.Fatal(err)
	}
	if !hasDiagnostic(result.Diagnostics, "MANIFEST_INVALID", "settings.gradle") {
		t.Fatalf("diagnostics = %#v", result.Diagnostics)
	}
}

func TestMavenDetectorIsInputOrderIndependent(t *testing.T) {
	root := t.TempDir()
	writeJavaFile(t, root, "a/pom.xml", `<project><artifactId>a</artifactId></project>`)
	writeJavaFile(t, root, "b/pom.xml", `<project><artifactId>b</artifactId></project>`)
	first, err := (MavenDetector{}).Detect(context.Background(), RepositoryView{Root: root, Paths: []string{"b/pom.xml", "a/pom.xml"}})
	if err != nil {
		t.Fatal(err)
	}
	second, err := (MavenDetector{}).Detect(context.Background(), RepositoryView{Root: root, Paths: []string{"a/pom.xml", "b/pom.xml"}})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(first, second) {
		t.Fatalf("input order changed result: %#v != %#v", first, second)
	}
}

func writeJavaFile(t *testing.T, root, relative, content string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
