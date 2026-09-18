package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/geoffrey-xiao/deprail/internal/discovery"
)

type DiscoverOptions struct {
	NoIgnore bool
}

func Discover(ctx context.Context, root string, options DiscoverOptions) (discovery.ProjectGraph, error) {
	absoluteRoot, err := filepath.Abs(root)
	if err != nil {
		return discovery.ProjectGraph{}, fmt.Errorf("resolve discovery path: %w", err)
	}
	result, err := discovery.Walk(ctx, absoluteRoot, discovery.WalkOptions{NoIgnore: options.NoIgnore})
	if err != nil {
		return discovery.ProjectGraph{}, err
	}
	detected, err := (discovery.NPMDetector{}).Detect(ctx, discovery.RepositoryView{Root: absoluteRoot, Paths: result.Paths})
	if err != nil {
		return discovery.ProjectGraph{}, err
	}
	detected = discovery.FinalizeDetection(detected)
	graph := discovery.ProjectGraph{
		SchemaVersion:  discovery.SchemaVersion,
		DocumentType:   "project",
		RepositoryRoot: ".",
		DisplayName:    filepath.Base(absoluteRoot),
		Workspaces:     detected.Workspaces,
		Completeness:   discovery.CompletenessForWorkspaces(detected.Workspaces),
		Diagnostics:    append([]discovery.Diagnostic{}, result.Diagnostics...),
	}
	graph.Diagnostics = append(graph.Diagnostics, detected.Diagnostics...)
	graph.Canonicalize()
	return graph, nil
}

func ValidateDiscoverRoot(root string) error {
	info, err := os.Stat(root)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("path is not a directory")
	}
	return nil
}
