package app

import (
	"context"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/geoffrey-xiao/deprail/internal/discovery"
)

func TestDiscoverMixedRepositoryFindsAllSupportedWorkspaceTypes(t *testing.T) {
	root := fixturePath(t, "mixed-repository")
	graph, err := Discover(context.Background(), root, DiscoverOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if graph.Completeness != discovery.Complete {
		t.Fatalf("completeness = %q, want complete because all workspaces have authoritative inputs", graph.Completeness)
	}
	if len(graph.Workspaces) != 3 {
		t.Fatalf("workspaces = %#v", graph.Workspaces)
	}
	want := map[string]struct {
		ecosystem    discovery.Ecosystem
		completeness discovery.Completeness
	}{
		"frontend":        {discovery.EcosystemNPM, discovery.Complete},
		"services/api":    {discovery.EcosystemPython, discovery.Complete},
		"services/worker": {discovery.EcosystemMaven, discovery.Complete},
	}
	for _, workspace := range graph.Workspaces {
		expected, ok := want[workspace.RelativePath]
		if !ok {
			t.Fatalf("unexpected workspace: %#v", workspace)
		}
		if workspace.Ecosystem != expected.ecosystem || workspace.Completeness != expected.completeness {
			t.Fatalf("workspace = %#v, want ecosystem=%q completeness=%q", workspace, expected.ecosystem, expected.completeness)
		}
		for _, manifest := range workspace.Manifests {
			if strings.Contains(manifest, "\\") || filepath.IsAbs(filepath.FromSlash(manifest)) {
				t.Fatalf("manifest path is not portable: %q", manifest)
			}
		}
	}
}

func TestDiscoverMixedRepositoryOutputIsDeterministic(t *testing.T) {
	root := fixturePath(t, "mixed-repository")
	first, err := Discover(context.Background(), root, DiscoverOptions{})
	if err != nil {
		t.Fatal(err)
	}
	second, err := Discover(context.Background(), root, DiscoverOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if first.SchemaVersion != second.SchemaVersion || first.Completeness != second.Completeness {
		t.Fatalf("discovery metadata changed: %#v != %#v", first, second)
	}
	if len(first.Workspaces) != len(second.Workspaces) || len(first.Diagnostics) != len(second.Diagnostics) {
		t.Fatalf("discovery shape changed: %#v != %#v", first, second)
	}
	for i := range first.Workspaces {
		if first.Workspaces[i].RelativePath != second.Workspaces[i].RelativePath || first.Workspaces[i].Ecosystem != second.Workspaces[i].Ecosystem {
			t.Fatalf("workspace ordering changed: %#v != %#v", first.Workspaces, second.Workspaces)
		}
	}
}

func fixturePath(t *testing.T, fixture string) string {
	t.Helper()
	_, source, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("locate app test source")
	}
	return filepath.Join(filepath.Dir(source), "..", "..", "testdata", "fixtures", fixture)
}
