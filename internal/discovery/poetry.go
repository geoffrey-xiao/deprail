package discovery

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type PoetryDetector struct{}

func (PoetryDetector) Name() string { return "poetry" }

func (PoetryDetector) Detect(ctx context.Context, repository RepositoryView) (DetectionResult, error) {
	paths := sortedRepositoryPaths(repository.Paths)
	pathSet := make(map[string]struct{}, len(paths))
	for _, path := range paths {
		pathSet[path] = struct{}{}
	}
	result := DetectionResult{}
	for _, manifest := range paths {
		select {
		case <-ctx.Done():
			return DetectionResult{}, ctx.Err()
		default:
		}
		if filepath.Base(filepath.FromSlash(manifest)) != "pyproject.toml" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(repository.Root, filepath.FromSlash(manifest)))
		if err != nil {
			result.Diagnostics = append(result.Diagnostics, Diagnostic{Code: "MANIFEST_INVALID", Message: err.Error(), Scope: manifest})
			continue
		}
		name, poetryProject, parseErr := parsePoetryProject(data)
		if parseErr != nil {
			result.Diagnostics = append(result.Diagnostics, Diagnostic{Code: "MANIFEST_INVALID", Message: parseErr.Error(), Scope: manifest})
			continue
		}
		if !poetryProject {
			continue
		}
		directory := filepath.ToSlash(filepath.Dir(filepath.FromSlash(manifest)))
		lockfile := "poetry.lock"
		if directory != "." {
			lockfile = directory + "/poetry.lock"
		}
		workspace := Workspace{
			SchemaVersion: SchemaVersion, DocumentType: "workspace", WorkspaceID: workspaceID(directory),
			RelativePath: directory, Ecosystem: EcosystemPython, PackageManager: "poetry",
			Manifests: []string{manifest}, Scope: name, Completeness: Complete,
		}
		if _, ok := pathSet[lockfile]; !ok {
			workspace.Completeness = Partial
			workspace.Diagnostics = append(workspace.Diagnostics, Diagnostic{Code: "DISCOVERY_INCOMPLETE", Message: "Poetry project has no authoritative lockfile", Scope: manifest})
		} else {
			workspace.Lockfiles = []string{lockfile}
			lockData, readErr := os.ReadFile(filepath.Join(repository.Root, filepath.FromSlash(lockfile)))
			if readErr != nil {
				workspace = poetryDiagnostic(workspace, readErr, lockfile)
			} else if !strings.Contains(string(lockData), "content-hash") && !strings.Contains(string(lockData), "[[package]]") {
				workspace = poetryDiagnostic(workspace, fmt.Errorf("Poetry lockfile has no package or content-hash declaration"), lockfile)
			}
		}
		result.Workspaces = append(result.Workspaces, workspace)
		result.Diagnostics = append(result.Diagnostics, workspace.Diagnostics...)
	}
	sort.SliceStable(result.Workspaces, func(i, j int) bool {
		return workspaceKey(result.Workspaces[i]) < workspaceKey(result.Workspaces[j])
	})
	sortDiagnostics(result.Diagnostics)
	return result, nil
}

func parsePoetryProject(data []byte) (string, bool, error) {
	content := string(data)
	if strings.ContainsRune(content, '\x00') {
		return "", false, fmt.Errorf("pyproject.toml contains a NUL byte")
	}
	inPoetry := false
	foundPoetry := false
	for _, rawLine := range strings.Split(content, "\n") {
		line := strings.TrimSpace(strings.SplitN(rawLine, "#", 2)[0])
		if line == "[tool.poetry]" {
			inPoetry = true
			foundPoetry = true
			continue
		}
		if strings.HasPrefix(line, "[") {
			inPoetry = false
		}
		if inPoetry && strings.HasPrefix(line, "name") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 || strings.TrimSpace(parts[1]) == "" {
				return "", true, fmt.Errorf("invalid Poetry name declaration")
			}
			return strings.Trim(strings.TrimSpace(parts[1]), "\"'"), true, nil
		}
	}
	if foundPoetry {
		return "", true, fmt.Errorf("Poetry metadata has no name declaration")
	}
	return "", false, nil
}

func poetryDiagnostic(workspace Workspace, err error, scope string) Workspace {
	workspace.Completeness = Failed
	workspace.Diagnostics = append(workspace.Diagnostics, Diagnostic{Code: "MANIFEST_INVALID", Message: err.Error(), Scope: scope})
	return workspace
}
