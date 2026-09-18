package discovery

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type RequirementsDetector struct{}
type UVDetector struct{}

func (RequirementsDetector) Name() string { return "requirements" }
func (UVDetector) Name() string           { return "uv" }

func (RequirementsDetector) Detect(ctx context.Context, repository RepositoryView) (DetectionResult, error) {
	paths := sortedRepositoryPaths(repository.Paths)
	result := DetectionResult{}
	for _, candidate := range paths {
		select {
		case <-ctx.Done():
			return DetectionResult{}, ctx.Err()
		default:
		}
		base := filepath.Base(filepath.FromSlash(candidate))
		if !strings.HasPrefix(base, "requirements") || !strings.HasSuffix(base, ".txt") {
			continue
		}
		workspace := pythonWorkspace(candidate, EcosystemPython, "pip", []string{candidate})
		result.Workspaces = append(result.Workspaces, workspace)
		result.Diagnostics = append(result.Diagnostics, workspace.Diagnostics...)
	}
	canonicalizeDetectionResult(&result)
	return result, nil
}

func (UVDetector) Detect(ctx context.Context, repository RepositoryView) (DetectionResult, error) {
	paths := sortedRepositoryPaths(repository.Paths)
	pathSet := make(map[string]struct{}, len(paths))
	for _, candidate := range paths {
		pathSet[candidate] = struct{}{}
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
		name, parseErr := parsePythonProject(data)
		if parseErr != nil {
			result.Diagnostics = append(result.Diagnostics, Diagnostic{Code: "MANIFEST_INVALID", Message: parseErr.Error(), Scope: manifest})
			continue
		}
		if name == "" {
			continue
		}
		directory := filepath.ToSlash(filepath.Dir(filepath.FromSlash(manifest)))
		if directory == "." {
			directory = "."
		}
		lockfile := filepath.ToSlash(filepath.Join(directory, "uv.lock"))
		if directory == "." {
			lockfile = "uv.lock"
		}
		workspace := Workspace{
			SchemaVersion: SchemaVersion, DocumentType: "workspace", WorkspaceID: workspaceID(directory),
			RelativePath: directory, Ecosystem: EcosystemPython, PackageManager: "uv",
			Manifests: []string{manifest}, Scope: name, Completeness: Complete,
		}
		if _, ok := pathSet[lockfile]; !ok {
			workspace.Completeness = Partial
			workspace.Diagnostics = append(workspace.Diagnostics, Diagnostic{Code: "DISCOVERY_INCOMPLETE", Message: "uv project has no authoritative lockfile", Scope: manifest})
		} else {
			workspace.Lockfiles = []string{lockfile}
			lockData, readErr := os.ReadFile(filepath.Join(repository.Root, filepath.FromSlash(lockfile)))
			if readErr != nil {
				workspace = pythonDiagnostic(workspace, readErr, lockfile)
			} else if !strings.Contains(string(lockData), "version") {
				workspace = pythonDiagnostic(workspace, fmt.Errorf("uv lockfile has no version declaration"), lockfile)
			}
		}
		result.Workspaces = append(result.Workspaces, workspace)
		result.Diagnostics = append(result.Diagnostics, workspace.Diagnostics...)
	}
	canonicalizeDetectionResult(&result)
	return result, nil
}

func pythonWorkspace(manifest string, ecosystem Ecosystem, packageManager string, manifests []string) Workspace {
	directory := filepath.ToSlash(filepath.Dir(filepath.FromSlash(manifest)))
	return Workspace{
		SchemaVersion: SchemaVersion, DocumentType: "workspace", WorkspaceID: workspaceID(directory),
		RelativePath: directory, Ecosystem: ecosystem, PackageManager: packageManager,
		Manifests: manifests, Completeness: Complete,
	}
}

func parsePythonProject(data []byte) (string, error) {
	content := string(data)
	if strings.ContainsRune(content, '\x00') {
		return "", fmt.Errorf("pyproject.toml contains a NUL byte")
	}
	inProject := false
	for _, rawLine := range strings.Split(content, "\n") {
		line := strings.TrimSpace(strings.SplitN(rawLine, "#", 2)[0])
		if line == "[project]" {
			inProject = true
			continue
		}
		if strings.HasPrefix(line, "[") {
			inProject = false
		}
		if inProject && strings.HasPrefix(line, "name") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 || strings.TrimSpace(parts[1]) == "" {
				return "", fmt.Errorf("invalid project name declaration")
			}
			return strings.Trim(strings.TrimSpace(parts[1]), "\"'"), nil
		}
	}
	return "", nil
}

func pythonDiagnostic(workspace Workspace, err error, scope string) Workspace {
	workspace.Completeness = Failed
	workspace.Diagnostics = append(workspace.Diagnostics, Diagnostic{Code: "MANIFEST_INVALID", Message: err.Error(), Scope: scope})
	return workspace
}

func sortedRepositoryPaths(paths []string) []string {
	normalized := make([]string, 0, len(paths))
	for _, path := range paths {
		normalized = append(normalized, normalizeRelativePath(path))
	}
	sort.Strings(normalized)
	return normalized
}

func canonicalizeDetectionResult(result *DetectionResult) {
	sort.SliceStable(result.Workspaces, func(i, j int) bool {
		return workspaceKey(result.Workspaces[i]) < workspaceKey(result.Workspaces[j])
	})
	sortDiagnostics(result.Diagnostics)
}
