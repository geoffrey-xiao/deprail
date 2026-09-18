package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

type PNPMDetector struct{}
type YarnDetector struct{}

func (PNPMDetector) Name() string { return "pnpm" }
func (YarnDetector) Name() string { return "yarn" }

func (PNPMDetector) Detect(ctx context.Context, repository RepositoryView) (DetectionResult, error) {
	return detectJavaScriptWorkspaces(ctx, repository, jsDetectorConfig{
		ecosystem:      EcosystemPNPM,
		packageManager: "pnpm",
		lockfile:       "pnpm-lock.yaml",
		workspaceFile:  "pnpm-workspace.yaml",
	})
}

func (YarnDetector) Detect(ctx context.Context, repository RepositoryView) (DetectionResult, error) {
	return detectJavaScriptWorkspaces(ctx, repository, jsDetectorConfig{
		ecosystem:      EcosystemYarn,
		packageManager: "yarn",
		lockfile:       "yarn.lock",
	})
}

type jsDetectorConfig struct {
	ecosystem      Ecosystem
	packageManager string
	lockfile       string
	workspaceFile  string
}

type packageManifest struct {
	Name           string          `json:"name"`
	PackageManager string          `json:"packageManager"`
	Workspaces     json.RawMessage `json:"workspaces"`
}

func detectJavaScriptWorkspaces(ctx context.Context, repository RepositoryView, config jsDetectorConfig) (DetectionResult, error) {
	paths := append([]string(nil), repository.Paths...)
	sort.Strings(paths)
	pathSet := make(map[string]struct{}, len(paths))
	for _, candidate := range paths {
		pathSet[normalizeRelativePath(candidate)] = struct{}{}
	}
	if _, ok := pathSet["package.json"]; !ok {
		return DetectionResult{}, nil
	}

	rootData, err := os.ReadFile(filepath.Join(repository.Root, "package.json"))
	if err != nil {
		return DetectionResult{}, err
	}
	var rootPackage packageManifest
	if err := json.Unmarshal(rootData, &rootPackage); err != nil {
		return DetectionResult{Diagnostics: []Diagnostic{{Code: "MANIFEST_INVALID", Message: err.Error(), Scope: "package.json"}}}, nil
	}
	patterns, err := workspacePatterns(repository.Root, pathSet, rootPackage, config)
	if err != nil {
		return DetectionResult{Diagnostics: []Diagnostic{{Code: "MANIFEST_INVALID", Message: err.Error(), Scope: config.workspaceFile}}}, nil
	}
	if !isManagedByDetector(rootPackage.PackageManager, config.packageManager, pathSet, config) {
		return DetectionResult{}, nil
	}

	result := DetectionResult{}
	for _, manifest := range paths {
		select {
		case <-ctx.Done():
			return DetectionResult{}, ctx.Err()
		default:
		}
		manifest = normalizeRelativePath(manifest)
		if filepath.Base(filepath.FromSlash(manifest)) != "package.json" || !isWorkspaceManifest(manifest, patterns) {
			continue
		}
		workspace := detectJSWorkspace(repository.Root, manifest, pathSet, config)
		result.Workspaces = append(result.Workspaces, workspace)
		result.Diagnostics = append(result.Diagnostics, workspace.Diagnostics...)
	}
	sort.SliceStable(result.Workspaces, func(i, j int) bool {
		return workspaceKey(result.Workspaces[i]) < workspaceKey(result.Workspaces[j])
	})
	sortDiagnostics(result.Diagnostics)
	return result, nil
}

func workspacePatterns(root string, paths map[string]struct{}, manifest packageManifest, config jsDetectorConfig) ([]string, error) {
	if config.workspaceFile != "" {
		if _, ok := paths[config.workspaceFile]; ok {
			data, err := os.ReadFile(filepath.Join(root, config.workspaceFile))
			if err != nil {
				return nil, err
			}
			return parsePNPMWorkspacePatterns(string(data)), nil
		}
	}
	if len(manifest.Workspaces) == 0 {
		return nil, nil
	}
	var patterns []string
	if err := json.Unmarshal(manifest.Workspaces, &patterns); err == nil {
		return patterns, nil
	}
	var object struct {
		Packages []string `json:"packages"`
	}
	if err := json.Unmarshal(manifest.Workspaces, &object); err != nil {
		return nil, err
	}
	return object.Packages, nil
}

func parsePNPMWorkspacePatterns(content string) []string {
	var patterns []string
	inPackages := false
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "packages:" {
			inPackages = true
			continue
		}
		if !inPackages || !strings.HasPrefix(trimmed, "-") {
			continue
		}
		value := strings.TrimSpace(strings.TrimPrefix(trimmed, "-"))
		value = strings.Trim(value, "\"'")
		if value != "" {
			patterns = append(patterns, value)
		}
	}
	return patterns
}

func isManagedByDetector(packageManager, expected string, paths map[string]struct{}, config jsDetectorConfig) bool {
	if strings.HasPrefix(packageManager, expected+"@") {
		return true
	}
	if config.workspaceFile != "" {
		_, ok := paths[config.workspaceFile]
		return ok
	}
	_, ok := paths[config.lockfile]
	return ok
}

func isWorkspaceManifest(manifest string, patterns []string) bool {
	if manifest == "package.json" {
		return true
	}
	directory := path.Dir(manifest)
	for _, pattern := range patterns {
		pattern = strings.TrimSuffix(strings.TrimSpace(pattern), "/")
		if pattern == "." || pattern == "" {
			continue
		}
		if ok, _ := path.Match(pattern+"/package.json", manifest); ok {
			return true
		}
		if strings.HasSuffix(pattern, "/*") && strings.HasPrefix(directory, strings.TrimSuffix(pattern, "/*")+"/") {
			return true
		}
	}
	return false
}

func detectJSWorkspace(root, manifest string, paths map[string]struct{}, config jsDetectorConfig) Workspace {
	directory := filepath.ToSlash(filepath.Dir(filepath.FromSlash(manifest)))
	workspace := Workspace{
		SchemaVersion: SchemaVersion, DocumentType: "workspace", WorkspaceID: workspaceID(directory),
		RelativePath: directory, Ecosystem: config.ecosystem, PackageManager: config.packageManager,
		Manifests: []string{manifest}, Completeness: Complete,
	}
	var packageData packageManifest
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(manifest)))
	if err != nil {
		return jsDiagnostic(workspace, err, manifest)
	}
	if err := json.Unmarshal(data, &packageData); err != nil {
		return jsDiagnostic(workspace, err, manifest)
	}
	workspace.Scope = packageData.Name
	if _, ok := paths[config.lockfile]; ok {
		workspace.Lockfiles = []string{config.lockfile}
		return workspace
	}
	workspace.Completeness = Partial
	workspace.Diagnostics = append(workspace.Diagnostics, Diagnostic{Code: "DISCOVERY_INCOMPLETE", Message: fmt.Sprintf("%s project has no authoritative lockfile", config.packageManager), Scope: manifest})
	return workspace
}

func jsDiagnostic(workspace Workspace, err error, scope string) Workspace {
	workspace.Completeness = Failed
	workspace.Diagnostics = append(workspace.Diagnostics, Diagnostic{Code: "MANIFEST_INVALID", Message: err.Error(), Scope: scope})
	return workspace
}
