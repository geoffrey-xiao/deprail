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
	if err := WriteHeader(w, "Repository: "+graph.DisplayName); err != nil {
		return err
	}
	if err := WriteStatus(w, string(graph.Completeness)); err != nil {
		return err
	}
	if err := WriteSummary(w, "Workspaces", len(graph.Workspaces)); err != nil {
		return err
	}
	if err := WriteSection(w, "Workspace results"); err != nil {
		return err
	}
	for _, workspace := range graph.Workspaces {
		if _, err := fmt.Fprintf(w, "- %s [%s] %s\n", workspace.RelativePath, workspace.Ecosystem, workspace.Completeness); err != nil {
			return err
		}
	}
	if verbose {
		if err := WriteSection(w, "Diagnostics"); err != nil {
			return err
		}
		for _, diagnostic := range graph.Diagnostics {
			if err := WriteError(w, diagnostic.Code, diagnostic.Scope+": "+diagnostic.Message); err != nil {
				return err
			}
		}
	}
	return nil
}
