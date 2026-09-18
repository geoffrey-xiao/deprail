package discovery

import "sort"

// FinalizeDetection applies cross-workspace conflict diagnostics and stable ordering.
// Individual detectors remain responsible for parsing their own inputs.
func FinalizeDetection(result DetectionResult) DetectionResult {
	for i := range result.Workspaces {
		workspace := &result.Workspaces[i]
		sort.Strings(workspace.Manifests)
		sort.Strings(workspace.Lockfiles)
		sort.Strings(workspace.Warnings)
		if len(workspace.Lockfiles) > 1 && !containsDiagnostic(workspace.Diagnostics, "LOCKFILE_CONFLICT", workspace.RelativePath) {
			workspace.Completeness = mergeCompleteness(workspace.Completeness, Partial)
			workspace.Diagnostics = append(workspace.Diagnostics, Diagnostic{
				Code:    "LOCKFILE_CONFLICT",
				Message: "multiple lockfiles are associated with one workspace",
				Scope:   workspace.RelativePath,
				Help:    "Keep one authoritative lockfile for the workspace",
			})
		}
	}
	result.Workspaces = append([]Workspace(nil), result.Workspaces...)
	sort.SliceStable(result.Workspaces, func(i, j int) bool {
		return workspaceKey(result.Workspaces[i]) < workspaceKey(result.Workspaces[j])
	})
	result.Diagnostics = append([]Diagnostic(nil), result.Diagnostics...)
	sortDiagnostics(result.Diagnostics)
	for i := range result.Workspaces {
		sortDiagnostics(result.Workspaces[i].Diagnostics)
		for _, diagnostic := range result.Workspaces[i].Diagnostics {
			if !containsDiagnostic(result.Diagnostics, diagnostic.Code, diagnostic.Scope) {
				result.Diagnostics = append(result.Diagnostics, diagnostic)
			}
		}
	}
	sortDiagnostics(result.Diagnostics)
	return result
}

func CompletenessForWorkspaces(workspaces []Workspace) Completeness {
	status := Complete
	for _, workspace := range workspaces {
		if workspace.Completeness == Failed {
			return Failed
		}
		if workspace.Completeness == Partial {
			status = Partial
		}
	}
	return status
}

func mergeCompleteness(current, candidate Completeness) Completeness {
	if current == Failed || candidate == Failed {
		return Failed
	}
	if current == Partial || candidate == Partial {
		return Partial
	}
	return Complete
}

func containsDiagnostic(diagnostics []Diagnostic, code, scope string) bool {
	for _, diagnostic := range diagnostics {
		if diagnostic.Code == code && diagnostic.Scope == scope {
			return true
		}
	}
	return false
}
