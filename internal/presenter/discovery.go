package presenter

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/geoffrey-xiao/deprail/internal/discovery"
)

func WriteDiscoveryJSON(w io.Writer, graph discovery.ProjectGraph) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(graph)
}

func WriteDiscoveryTerminal(w io.Writer, graph discovery.ProjectGraph, verbose bool) error {
	if _, err := fmt.Fprintf(w, "Repository: %s\n", graph.DisplayName); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "Completeness: %s\n", graph.Completeness); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "Workspaces: %d\n", len(graph.Workspaces)); err != nil {
		return err
	}
	for _, workspace := range graph.Workspaces {
		if _, err := fmt.Fprintf(w, "- %s [%s] %s\n", workspace.RelativePath, workspace.Ecosystem, workspace.Completeness); err != nil {
			return err
		}
	}
	if verbose {
		for _, diagnostic := range graph.Diagnostics {
			if _, err := fmt.Fprintf(w, "Diagnostic: %s (%s): %s\n", diagnostic.Code, diagnostic.Scope, diagnostic.Message); err != nil {
				return err
			}
		}
	}
	return nil
}
