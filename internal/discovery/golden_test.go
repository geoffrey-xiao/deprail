package discovery

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestNPMSmallFixtureProducesStableProjectGraphGolden(t *testing.T) {
	fixtureRoot := fixturePath(t, "npm-basic")
	walked, err := Walk(context.Background(), fixtureRoot, WalkOptions{})
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range walked.Paths {
		if strings.Contains(path, "\\") || filepath.IsAbs(filepath.FromSlash(path)) {
			t.Fatalf("path %q is not repository-relative slash-separated", path)
		}
	}
	detected, err := (NPMDetector{}).Detect(context.Background(), RepositoryView{Root: fixtureRoot, Paths: walked.Paths})
	if err != nil {
		t.Fatal(err)
	}
	graph := ProjectGraph{
		SchemaVersion:  SchemaVersion,
		DocumentType:   "project",
		RepositoryRoot: ".",
		DisplayName:    "npm-basic",
		Workspaces:     detected.Workspaces,
		Completeness:   Complete,
		Diagnostics:    append([]Diagnostic{}, detected.Diagnostics...),
	}
	graph.Canonicalize()
	if err := ValidateProjectGraph(graph); err != nil {
		t.Fatalf("generated graph is invalid: %v", err)
	}
	actual, err := json.MarshalIndent(graph, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	actual = append(actual, '\n')
	goldenPath := filepath.Join(fixtureRoot, "..", "..", "golden", "discovery-npm-basic.json")
	golden, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatal(err)
	}
	goldenText := normalizeGoldenLineEndings(string(golden))
	if string(actual) != goldenText {
		t.Fatalf("discovery golden changed:\nactual:\n%s\nwant:\n%s", actual, goldenText)
	}
}
func TestNormalizeGoldenLineEndings(t *testing.T) {
	if got := normalizeGoldenLineEndings("first\r\nsecond\r\n"); got != "first\nsecond\n" {
		t.Fatalf("normalized line endings = %q", got)
	}
}

func normalizeGoldenLineEndings(content string) string {
	return strings.ReplaceAll(content, "\r\n", "\n")
}

func fixturePath(t *testing.T, fixture string) string {
	t.Helper()
	_, source, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("locate discovery test source")
	}
	return filepath.Join(filepath.Dir(source), "..", "..", "testdata", "fixtures", fixture)
}
