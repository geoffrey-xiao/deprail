package scanplan

import (
	"fmt"
	"sort"

	"github.com/geoffrey-xiao/deprail/internal/adapter"
	"github.com/geoffrey-xiao/deprail/internal/discovery"
)

// Unit is one deterministic scanner execution unit for a discovered workspace.
type Unit struct {
	Target adapter.Target
	Args   []string
}

// Build creates units only for complete workspaces. Incomplete workspaces are
// retained by the caller for diagnostics and must not become successful scans.
func Build(graph discovery.ProjectGraph) ([]Unit, error) {
	units := make([]Unit, 0, len(graph.Workspaces))
	for _, workspace := range graph.Workspaces {
		if workspace.Completeness != discovery.Complete {
			continue
		}
		files := append([]string(nil), workspace.Manifests...)
		files = append(files, workspace.Lockfiles...)
		sort.Strings(files)
		if len(files) == 0 {
			return nil, fmt.Errorf("workspace %q has no authoritative files", workspace.WorkspaceID)
		}
		units = append(units, Unit{Target: adapter.Target{
			WorkspaceID: workspace.WorkspaceID, RelativePath: workspace.RelativePath,
			Ecosystem: string(workspace.Ecosystem), PackageFiles: files,
		}, Args: []string{"scan", "--format", "json", workspace.RelativePath}})
	}
	if len(units) == 0 {
		return nil, fmt.Errorf("project has no complete scan targets")
	}
	sort.Slice(units, func(i, j int) bool {
		return units[i].Target.RelativePath+"\x00"+units[i].Target.WorkspaceID < units[j].Target.RelativePath+"\x00"+units[j].Target.WorkspaceID
	})
	return units, nil
}
