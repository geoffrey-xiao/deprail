package historyv1_test

import (
	"embed"
	"encoding/json"
	"regexp"
	"strings"
	"testing"

	"github.com/geoffrey-xiao/deprail/internal/domain/history"
	"github.com/geoffrey-xiao/deprail/internal/normalize"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

//go:embed history.schema.json examples/*.json
var schemaFiles embed.FS

func TestHistoryExamplesValidateAgainstSchema(t *testing.T) {
	schemaBytes, err := schemaFiles.ReadFile("history.schema.json")
	if err != nil {
		t.Fatal(err)
	}
	var schemaDocument any
	if err := json.Unmarshal(schemaBytes, &schemaDocument); err != nil {
		t.Fatal(err)
	}
	const schemaURL = "https://deprail.dev/schemas/history-v1/history.schema.json"
	compiler := newHistorySchemaCompiler()
	if err := compiler.AddResource(schemaURL, schemaDocument); err != nil {
		t.Fatal(err)
	}
	compiled, err := compiler.Compile(schemaURL)
	if err != nil {
		t.Fatalf("compile schema: %v", err)
	}
	for _, name := range []string{"examples/cancelled-no-context.json", "examples/complete-empty.json", "examples/nonempty.json"} {
		t.Run(name, func(t *testing.T) {
			data, err := schemaFiles.ReadFile(name)
			if err != nil {
				t.Fatal(err)
			}
			var document any
			if err := json.Unmarshal(data, &document); err != nil {
				t.Fatal(err)
			}
			if err := compiled.Validate(document); err != nil {
				t.Fatalf("example violates history-v1 schema: %v", err)
			}
		})
	}
}

func TestHistorySchemaRejectsUnknownFields(t *testing.T) {
	schemaBytes, err := schemaFiles.ReadFile("history.schema.json")
	if err != nil {
		t.Fatal(err)
	}
	var schemaDocument any
	if err := json.Unmarshal(schemaBytes, &schemaDocument); err != nil {
		t.Fatal(err)
	}
	const schemaURL = "https://deprail.dev/schemas/history-v1/history.schema.json"
	compiler := newHistorySchemaCompiler()
	if err := compiler.AddResource(schemaURL, schemaDocument); err != nil {
		t.Fatal(err)
	}
	compiled, err := compiler.Compile(schemaURL)
	if err != nil {
		t.Fatalf("compile schema: %v", err)
	}
	data, err := schemaFiles.ReadFile("examples/complete-empty.json")
	if err != nil {
		t.Fatal(err)
	}
	var document map[string]any
	if err := json.Unmarshal(data, &document); err != nil {
		t.Fatal(err)
	}
	document["absolute_root"] = "/private/worktree"
	if err := compiled.Validate(document); err == nil {
		t.Fatal("history schema accepted a property outside the exact allowlist")
	}
}
func TestHistorySchemaRejectsFutureProjectionMarker(t *testing.T) {
	schemaBytes, err := schemaFiles.ReadFile("history.schema.json")
	if err != nil {
		t.Fatal(err)
	}
	var schemaDocument any
	if err := json.Unmarshal(schemaBytes, &schemaDocument); err != nil {
		t.Fatal(err)
	}
	const schemaURL = "https://deprail.dev/schemas/history-v1/history.schema.json"
	compiler := newHistorySchemaCompiler()
	if err := compiler.AddResource(schemaURL, schemaDocument); err != nil {
		t.Fatal(err)
	}
	compiled, err := compiler.Compile(schemaURL)
	if err != nil {
		t.Fatalf("compile schema: %v", err)
	}
	data, err := schemaFiles.ReadFile("examples/nonempty.json")
	if err != nil {
		t.Fatal(err)
	}
	var document map[string]any
	if err := json.Unmarshal(data, &document); err != nil {
		t.Fatal(err)
	}
	document["schema_version"] = "history-v2"
	if err := compiled.Validate(document); err == nil {
		t.Fatal("history schema accepted a future projection marker")
	}
}
func TestRuntimeProjectionStatesValidateAgainstHistorySchema(t *testing.T) {
	schemaBytes, err := schemaFiles.ReadFile("history.schema.json")
	if err != nil {
		t.Fatal(err)
	}
	var schemaDocument any
	if err := json.Unmarshal(schemaBytes, &schemaDocument); err != nil {
		t.Fatal(err)
	}
	const schemaURL = "https://deprail.dev/schemas/history-v1/history.schema.json"
	compiler := newHistorySchemaCompiler()
	if err := compiler.AddResource(schemaURL, schemaDocument); err != nil {
		t.Fatal(err)
	}
	compiled, err := compiler.Compile(schemaURL)
	if err != nil {
		t.Fatalf("compile schema: %v", err)
	}

	state := "a1b2c3d4e5f67890a1b2c3d4e5f67890a1b2c3d4e5f67890a1b2c3d4e5f67890"
	input := normalize.HistoryInput{
		HistoryEntryID:   "2d931510-d99f-494a-8c67-87feb05e1594",
		RecordedAtUS:     1780000000123456,
		OperationOutcome: "completed",
		Workspaces:       []history.Workspace{},
		Report:           &normalize.HistoryReportInput{SchemaVersion: "v1alpha", ScanID: "scan", Status: "complete", RepositoryState: state, Findings: []normalize.HistoryFindingInput{}, ArtifactDigests: []string{}},
	}
	cases := []struct {
		name   string
		mutate func(*normalize.HistoryInput)
	}{
		{"complete-empty", func(*normalize.HistoryInput) {}},
		{"nonempty", func(in *normalize.HistoryInput) {
			in.Workspaces = nil
			in.Report.Findings = []normalize.HistoryFindingInput{{Component: `lib <parser> "quoted" 東京`, PURL: "pkg:npm/lodash@4.17.20", Version: "4.17.20", WorkspaceID: "workspace", WorkspacePath: ".", Ecosystem: "npm", TargetID: "GHSA-example", Aliases: []string{}, Severity: "high", FixedVersion: "4.17.21"}}
		}},
		{"partial", func(in *normalize.HistoryInput) {
			in.Report.Status = "partial"
			in.Diagnostics = []normalize.HistoryDiagnosticInput{{Code: "SCANNER_OUTPUT_LIMIT", Scope: "repository"}}
		}},
		{"failed", func(in *normalize.HistoryInput) {
			in.OperationOutcome = "failed"
			in.Report.Status = "failed"
			in.Diagnostics = []normalize.HistoryDiagnosticInput{{Code: "SCANNER_EXIT_NONZERO", Scope: "repository"}}
		}},
		{"cancelled-without-report", func(in *normalize.HistoryInput) {
			in.OperationOutcome = "cancelled"
			in.Report = nil
			in.Workspaces = nil
			in.Diagnostics = []normalize.HistoryDiagnosticInput{{Code: "CANCELLED", Scope: "repository"}}
		}},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			candidate := input
			candidate.Report = new(normalize.HistoryReportInput)
			if input.Report != nil {
				*candidate.Report = *input.Report
			}
			candidate.Diagnostics = nil
			test.mutate(&candidate)
			projection, err := normalize.ProjectHistory(candidate)
			if err != nil {
				t.Fatalf("project runtime case: %v", err)
			}
			data, err := json.Marshal(projection)
			if err != nil {
				t.Fatalf("encode projection: %v", err)
			}
			var document any
			if err := json.Unmarshal(data, &document); err != nil {
				t.Fatalf("decode projection: %v", err)
			}
			if err := compiled.Validate(document); err != nil {
				t.Fatalf("runtime projection violates history schema: %v", err)
			}
		})
	}
}

// The history schema's reviewed path pattern uses negative lookaheads unsupported
// by Go regexp. This test-only engine preserves those exact path constraints.
const historyPathPattern = `^(?![\s\S]*[\u0000-\u001f])(?:\.|(?!/)(?![A-Za-z]:)(?!.*(?:^|/)\.\.(?:/|$))(?!.*(?:^|/)\.(?:/|$))(?!.*\\)[^/]+(?:/[^/]+)*)$`

type historySchemaRegexp struct {
	source string
	re     *regexp.Regexp
	path   bool
}

func (r historySchemaRegexp) String() string { return r.source }
func (r historySchemaRegexp) MatchString(value string) bool {
	if !r.re.MatchString(value) {
		return false
	}
	if !r.path {
		return true
	}
	if value == "." {
		return true
	}
	if value == "" || strings.HasPrefix(value, "/") || strings.Contains(value, `\`) {
		return false
	}
	if len(value) >= 2 && ((value[0] >= 'A' && value[0] <= 'Z') || (value[0] >= 'a' && value[0] <= 'z')) && value[1] == ':' {
		return false
	}
	for _, char := range value {
		if char < 0x20 {
			return false
		}
	}
	for _, part := range strings.Split(value, "/") {
		if part == "" || part == "." || part == ".." {
			return false
		}
	}
	return true
}

func historyRegexpEngine(pattern string) (jsonschema.Regexp, error) {
	if pattern == historyPathPattern {
		return historySchemaRegexp{source: pattern, re: regexp.MustCompile(`^(?:\.|[^/]+(?:/[^/]+)*)$`), path: true}, nil
	}
	compiled, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return historySchemaRegexp{source: pattern, re: compiled}, nil
}

func newHistorySchemaCompiler() *jsonschema.Compiler {
	compiler := jsonschema.NewCompiler()
	compiler.UseRegexpEngine(historyRegexpEngine)
	return compiler
}

func TestHistorySchemaPathRegexpAdapter(t *testing.T) {
	matcher, err := historyRegexpEngine(historyPathPattern)
	if err != nil {
		t.Fatal(err)
	}
	for _, value := range []string{".", "src/package.json", "dir/東京 & package!"} {
		if !matcher.MatchString(value) {
			t.Errorf("schema path matcher rejected allowed value %q", value)
		}
	}
	for _, value := range []string{"", "/absolute", `\\server\share`, "C:/repo", "../escape", "src/../escape", "src/./file", "src//file", `src\file`, "bad\x00path"} {
		if matcher.MatchString(value) {
			t.Errorf("schema path matcher accepted hostile value %q", value)
		}
	}
	if _, err := historyRegexpEngine(`^(?!unsupported)$`); err == nil {
		t.Fatal("unsupported non-path lookaround pattern silently received an approximation")
	}
}
