package discovery

import (
	"context"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type MavenDetector struct{}
type GradleDetector struct{}

func (MavenDetector) Name() string  { return "maven" }
func (GradleDetector) Name() string { return "gradle" }

type mavenProject struct {
	XMLName    xml.Name `xml:"project"`
	Model      string   `xml:"modelVersion"`
	GroupID    string   `xml:"groupId"`
	ArtifactID string   `xml:"artifactId"`
	Version    string   `xml:"version"`
	Parent     struct {
		GroupID string `xml:"groupId"`
	} `xml:"parent"`
}

func (MavenDetector) Detect(ctx context.Context, repository RepositoryView) (DetectionResult, error) {
	paths := sortedRepositoryPaths(repository.Paths)
	result := DetectionResult{}
	for _, manifest := range paths {
		select {
		case <-ctx.Done():
			return DetectionResult{}, ctx.Err()
		default:
		}
		if filepath.Base(filepath.FromSlash(manifest)) != "pom.xml" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(repository.Root, filepath.FromSlash(manifest)))
		if err != nil {
			result.Diagnostics = append(result.Diagnostics, Diagnostic{Code: "MANIFEST_INVALID", Message: err.Error(), Scope: manifest})
			continue
		}
		var project mavenProject
		if err := xml.Unmarshal(data, &project); err != nil || project.XMLName.Local != "project" || project.ArtifactID == "" {
			if err == nil {
				err = fmt.Errorf("Maven project is missing artifactId")
			}
			result.Diagnostics = append(result.Diagnostics, Diagnostic{Code: "MANIFEST_INVALID", Message: err.Error(), Scope: manifest})
			continue
		}
		groupID := project.GroupID
		if groupID == "" {
			groupID = project.Parent.GroupID
		}
		scope := project.ArtifactID
		if groupID != "" {
			scope = groupID + ":" + project.ArtifactID
		}
		directory := filepath.ToSlash(filepath.Dir(filepath.FromSlash(manifest)))
		workspace := Workspace{
			SchemaVersion: SchemaVersion, DocumentType: "workspace", WorkspaceID: workspaceID(directory),
			RelativePath: directory, Ecosystem: EcosystemMaven, PackageManager: "maven",
			Manifests: []string{manifest}, Scope: scope, Completeness: Complete,
		}
		result.Workspaces = append(result.Workspaces, workspace)
	}
	canonicalizeDetectionResult(&result)
	return result, nil
}

func (GradleDetector) Detect(ctx context.Context, repository RepositoryView) (DetectionResult, error) {
	paths := sortedRepositoryPaths(repository.Paths)
	pathSet := make(map[string]struct{}, len(paths))
	for _, path := range paths {
		pathSet[path] = struct{}{}
	}
	result := DetectionResult{}
	for _, manifest := range paths {
		base := filepath.Base(filepath.FromSlash(manifest))
		if base != "settings.gradle" && base != "settings.gradle.kts" && base != "build.gradle" && base != "build.gradle.kts" {
			continue
		}
		select {
		case <-ctx.Done():
			return DetectionResult{}, ctx.Err()
		default:
		}
		data, err := os.ReadFile(filepath.Join(repository.Root, filepath.FromSlash(manifest)))
		if err != nil {
			result.Diagnostics = append(result.Diagnostics, Diagnostic{Code: "MANIFEST_INVALID", Message: err.Error(), Scope: manifest})
			continue
		}
		if strings.ContainsRune(string(data), '\x00') {
			result.Diagnostics = append(result.Diagnostics, Diagnostic{Code: "MANIFEST_INVALID", Message: "Gradle build file contains a NUL byte", Scope: manifest})
			continue
		}
		directory := filepath.ToSlash(filepath.Dir(filepath.FromSlash(manifest)))
		if base == "settings.gradle" || base == "settings.gradle.kts" {
			workspace := gradleWorkspace(directory, manifest, pathSet, string(data), true)
			result.Workspaces = append(result.Workspaces, workspace)
			continue
		}
		if _, hasSettings := pathSet[settingsPath(directory, pathSet)]; hasSettings {
			continue
		}
		workspace := gradleWorkspace(directory, manifest, pathSet, string(data), false)
		result.Workspaces = append(result.Workspaces, workspace)
	}
	for _, workspace := range result.Workspaces {
		result.Diagnostics = append(result.Diagnostics, workspace.Diagnostics...)
	}
	canonicalizeDetectionResult(&result)
	return result, nil
}

func gradleWorkspace(directory, manifest string, paths map[string]struct{}, content string, settings bool) Workspace {
	scope := extractQuotedValue(content, "rootProject.name")
	if scope == "" {
		scope = filepath.Base(directory)
		if directory == "." {
			scope = "."
		}
	}
	workspace := Workspace{
		SchemaVersion: SchemaVersion, DocumentType: "workspace", WorkspaceID: workspaceID(directory),
		RelativePath: directory, Ecosystem: EcosystemGradle, PackageManager: "gradle",
		Manifests: []string{manifest}, Scope: scope, Completeness: Complete,
	}
	if settings {
		for _, candidate := range []string{"build.gradle", "build.gradle.kts"} {
			path := candidate
			if directory != "." {
				path = directory + "/" + candidate
			}
			if _, ok := paths[path]; ok {
				workspace.Manifests = append(workspace.Manifests, path)
			}
		}
	}
	lockfile := "gradle.lockfile"
	if directory != "." {
		lockfile = directory + "/gradle.lockfile"
	}
	if _, ok := paths[lockfile]; ok {
		workspace.Lockfiles = []string{lockfile}
	}
	return workspace
}

func settingsPath(directory string, paths map[string]struct{}) string {
	for _, name := range []string{"settings.gradle", "settings.gradle.kts"} {
		candidate := name
		if directory != "." {
			candidate = directory + "/" + name
		}
		if _, ok := paths[candidate]; ok {
			return candidate
		}
	}
	return "__missing_settings__"
}

func extractQuotedValue(content, key string) string {
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, key) {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		return strings.Trim(strings.TrimSpace(parts[1]), "\"'")
	}
	return ""
}

func _sortJavaWorkspaces(workspaces []Workspace) {
	sort.SliceStable(workspaces, func(i, j int) bool { return workspaceKey(workspaces[i]) < workspaceKey(workspaces[j]) })
}
