package discovery

import (
	"context"
	"errors"
	"sort"
)

const SchemaVersion = "v1alpha"

type Completeness string

const (
	Complete Completeness = "complete"
	Partial  Completeness = "partial"
	Failed   Completeness = "failed"
)

type Ecosystem string

const (
	EcosystemNPM     Ecosystem = "npm"
	EcosystemPNPM    Ecosystem = "pnpm"
	EcosystemYarn    Ecosystem = "yarn"
	EcosystemPython  Ecosystem = "python"
	EcosystemMaven   Ecosystem = "maven"
	EcosystemGradle  Ecosystem = "gradle"
	EcosystemUnknown Ecosystem = "unknown"
)

type Diagnostic struct {
	Code      string         `json:"code"`
	Message   string         `json:"message"`
	Scope     string         `json:"scope"`
	Retryable bool           `json:"retryable"`
	Details   map[string]any `json:"details,omitempty"`
	Help      string         `json:"help,omitempty"`
}

type Workspace struct {
	SchemaVersion  string       `json:"schema_version"`
	DocumentType   string       `json:"document_type"`
	WorkspaceID    string       `json:"workspace_id"`
	RelativePath   string       `json:"relative_path"`
	Ecosystem      Ecosystem    `json:"ecosystem"`
	PackageManager string       `json:"package_manager"`
	Manifests      []string     `json:"manifests"`
	Lockfiles      []string     `json:"lockfiles"`
	Scope          string       `json:"scope,omitempty"`
	Completeness   Completeness `json:"completeness"`
	Warnings       []string     `json:"warnings,omitempty"`
	Diagnostics    []Diagnostic `json:"diagnostics,omitempty"`
}

type ProjectGraph struct {
	SchemaVersion       string       `json:"schema_version"`
	DocumentType        string       `json:"document_type"`
	RepositoryRoot      string       `json:"repository_root"`
	DisplayName         string       `json:"display_name"`
	ConfigurationDigest string       `json:"configuration_digest,omitempty"`
	Workspaces          []Workspace  `json:"workspaces"`
	Completeness        Completeness `json:"completeness"`
	Diagnostics         []Diagnostic `json:"diagnostics"`
}

type RepositoryView struct {
	Root  string
	Paths []string
}

type Detector interface {
	Name() string
	Detect(context.Context, RepositoryView) (DetectionResult, error)
}

type DetectionResult struct {
	Workspaces  []Workspace
	Diagnostics []Diagnostic
}

func (g *ProjectGraph) Canonicalize() {
	for i := range g.Workspaces {
		workspace := &g.Workspaces[i]
		sort.Strings(workspace.Manifests)
		sort.Strings(workspace.Lockfiles)
		sort.Strings(workspace.Warnings)
		sortDiagnostics(workspace.Diagnostics)
	}
	sort.SliceStable(g.Workspaces, func(i, j int) bool {
		left, right := g.Workspaces[i], g.Workspaces[j]
		return workspaceKey(left) < workspaceKey(right)
	})
	sortDiagnostics(g.Diagnostics)
}

func workspaceKey(workspace Workspace) string {
	return workspace.RelativePath + "\x00" + string(workspace.Ecosystem) + "\x00" + workspace.PackageManager + "\x00" + workspace.WorkspaceID
}

func sortDiagnostics(diagnostics []Diagnostic) {
	sort.SliceStable(diagnostics, func(i, j int) bool {
		left, right := diagnostics[i], diagnostics[j]
		if left.Scope != right.Scope {
			return left.Scope < right.Scope
		}
		if left.Code != right.Code {
			return left.Code < right.Code
		}
		return left.Message < right.Message
	})
}

func ValidateProjectGraph(graph ProjectGraph) error {
	if graph.SchemaVersion != SchemaVersion {
		return errors.New("unsupported schema version")
	}
	if graph.DocumentType != "project" {
		return errors.New("invalid project document type")
	}
	if graph.RepositoryRoot == "" || graph.DisplayName == "" {
		return errors.New("project identity is required")
	}
	for _, workspace := range graph.Workspaces {
		if workspace.SchemaVersion != SchemaVersion || workspace.DocumentType != "workspace" {
			return errors.New("workspace schema identity is required")
		}
		if workspace.WorkspaceID == "" || workspace.RelativePath == "" || workspace.PackageManager == "" {
			return errors.New("workspace identity is required")
		}
	}
	return nil
}
