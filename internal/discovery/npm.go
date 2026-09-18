package discovery

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type NPMDetector struct{}

func (NPMDetector) Name() string { return "npm" }

func (NPMDetector) Detect(ctx context.Context, repository RepositoryView) (DetectionResult, error) {
	paths := append([]string(nil), repository.Paths...)
	sort.Strings(paths)
	pathSet := make(map[string]struct{}, len(paths))
	for _, path := range paths {
		pathSet[normalizeRelativePath(path)] = struct{}{}
	}

	result := DetectionResult{}
	for _, path := range paths {
		select {
		case <-ctx.Done():
			return DetectionResult{}, ctx.Err()
		default:
		}
		path = normalizeRelativePath(path)
		if filepath.Base(filepath.FromSlash(path)) != "package.json" {
			continue
		}
		workspace := detectNPMWorkspace(repository.Root, path, pathSet)
		result.Workspaces = append(result.Workspaces, workspace)
		result.Diagnostics = append(result.Diagnostics, workspace.Diagnostics...)
	}
	return result, nil
}

func detectNPMWorkspace(root, manifest string, paths map[string]struct{}) Workspace {
	directory := filepath.ToSlash(filepath.Dir(filepath.FromSlash(manifest)))
	workspace := Workspace{
		SchemaVersion:  SchemaVersion,
		DocumentType:   "workspace",
		WorkspaceID:    workspaceID(directory),
		RelativePath:   directory,
		Ecosystem:      EcosystemNPM,
		PackageManager: "npm",
		Manifests:      []string{manifest},
		Completeness:   Complete,
	}

	manifestPath := filepath.Join(root, filepath.FromSlash(manifest))
	var packageManifest struct {
		Name string `json:"name"`
	}
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return npmDiagnostic(workspace, "MANIFEST_INVALID", err.Error(), manifest)
	}
	if err := json.Unmarshal(data, &packageManifest); err != nil {
		return npmDiagnostic(workspace, "MANIFEST_INVALID", err.Error(), manifest)
	}
	workspace.Scope = packageManifest.Name

	lockfile := chooseNPMLockfile(directory, paths)
	if lockfile == "" {
		workspace.Completeness = Partial
		workspace.Diagnostics = append(workspace.Diagnostics, Diagnostic{
			Code: "DISCOVERY_INCOMPLETE", Message: "npm project has no authoritative lockfile", Scope: manifest,
			Help: "Add package-lock.json or npm-shrinkwrap.json",
		})
		return workspace
	}
	workspace.Lockfiles = []string{lockfile}
	lockData, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(lockfile)))
	if err != nil {
		return npmDiagnostic(workspace, "MANIFEST_INVALID", err.Error(), lockfile)
	}
	var lockDocument map[string]any
	if err := json.Unmarshal(lockData, &lockDocument); err != nil {
		return npmDiagnostic(workspace, "MANIFEST_INVALID", err.Error(), lockfile)
	}
	return workspace
}

func chooseNPMLockfile(directory string, paths map[string]struct{}) string {
	candidates := []string{"package-lock.json", "npm-shrinkwrap.json"}
	for _, name := range candidates {
		path := name
		if directory != "" && directory != "." {
			path = directory + "/" + name
		}
		if _, ok := paths[path]; ok {
			return path
		}
	}
	return ""
}

func npmDiagnostic(workspace Workspace, code, message, scope string) Workspace {
	workspace.Completeness = Failed
	workspace.Diagnostics = append(workspace.Diagnostics, Diagnostic{Code: code, Message: message, Scope: scope})
	return workspace
}

func workspaceID(relativePath string) string {
	if relativePath == "" || relativePath == "." {
		return "."
	}
	return relativePath
}

func normalizeRelativePath(path string) string {
	path = filepath.ToSlash(path)
	path = strings.TrimPrefix(path, "./")
	if path == "." {
		return ""
	}
	return path
}
