package normalize_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/geoffrey-xiao/deprail/internal/domain/history"
	"github.com/geoffrey-xiao/deprail/internal/normalize"
)

const historyEntryID = "2d931510-d99f-494a-8c67-87feb05e1594"
const repositoryState = "a1b2c3d4e5f67890a1b2c3d4e5f67890a1b2c3d4e5f67890a1b2c3d4e5f67890"

func historyInput() normalize.HistoryInput {
	return normalize.HistoryInput{
		HistoryEntryID:   historyEntryID,
		RecordedAtUS:     1780000000123456,
		OperationOutcome: "completed",
		RepositoryLabel:  new("service"),
		Workspaces:       []history.Workspace{{WorkspaceID: "ws-npm-root", Path: ".", Ecosystem: "npm", PackageManager: "npm", DiscoveryCompleteness: "complete"}},
		Report: &normalize.HistoryReportInput{
			SchemaVersion: "v1alpha", ScanID: "scan-a1b2c3d4", Status: "complete", RepositoryState: repositoryState,
			Findings: []normalize.HistoryFindingInput{{Component: "lodash", PURL: "pkg:npm/lodash@4.17.20", Version: "4.17.20", WorkspaceID: "ws-npm-root", WorkspacePath: ".", Ecosystem: "npm", TargetID: "GHSA-example", Aliases: []string{"CVE-2024-0001"}, Severity: "high", FixedVersion: "4.17.21"}},
		},
	}
}

func TestProjectHistoryPreservesReportAndFindingIdentity(t *testing.T) {
	input := historyInput()
	projection, err := normalize.ProjectHistory(input)
	if err != nil {
		t.Fatalf("project history: %v", err)
	}
	if projection.SchemaVersion != history.SchemaVersion || projection.HistoryEntryID != input.HistoryEntryID || projection.RecordedAtUS != input.RecordedAtUS {
		t.Fatalf("operation identity was not preserved: %#v", projection)
	}
	if projection.Report == nil || projection.Report.SourceSchemaVersion != "v1alpha" || projection.Report.SourceScanID != "scan-a1b2c3d4" {
		t.Fatalf("source report identity was not preserved: %#v", projection.Report)
	}
	if len(projection.Report.Findings) != 1 {
		t.Fatalf("findings = %d, want one", len(projection.Report.Findings))
	}
	finding := projection.Report.Findings[0]
	wantKey := normalize.StableFindingKey(normalize.FindingInput{WorkspaceID: input.Report.Findings[0].WorkspaceID, ComponentPURL: input.Report.Findings[0].PURL, ComponentVersion: input.Report.Findings[0].Version, VulnerabilityID: input.Report.Findings[0].TargetID, VulnerabilityAliases: input.Report.Findings[0].Aliases})
	if finding.StableFindingKey != wantKey || finding.VulnerabilityID != "GHSA-example" || finding.ComponentPURL == nil || *finding.ComponentPURL != input.Report.Findings[0].PURL {
		t.Fatalf("finding identity was not preserved: %#v", finding)
	}
	if err := normalize.ValidateProjection(projection); err != nil {
		t.Fatalf("projected value does not validate: %v", err)
	}
}

func TestProjectHistoryDistinguishesUnavailableFromKnownEmpty(t *testing.T) {
	unavailable := historyInput()
	unavailable.Workspaces = nil
	unavailable.Report = nil
	gotUnavailable, err := normalize.ProjectHistory(unavailable)
	if err != nil {
		t.Fatalf("project unavailable context: %v", err)
	}
	if gotUnavailable.Workspaces != nil || gotUnavailable.Report != nil {
		t.Fatalf("unavailable context became known empty: %#v", gotUnavailable)
	}

	empty := historyInput()
	empty.Workspaces = []history.Workspace{}
	empty.Report.Findings = []normalize.HistoryFindingInput{}
	gotEmpty, err := normalize.ProjectHistory(empty)
	if err != nil {
		t.Fatalf("project known-empty context: %v", err)
	}
	if gotEmpty.Workspaces == nil || len(gotEmpty.Workspaces) != 0 || gotEmpty.Report == nil || gotEmpty.Report.Findings == nil || len(gotEmpty.Report.Findings) != 0 {
		t.Fatalf("known-empty context lost its availability: %#v", gotEmpty)
	}
}

func TestProjectHistoryKeepsOutcomeIndependentFromReportStatus(t *testing.T) {
	input := historyInput()
	input.OperationOutcome = "failed"
	input.Report.Status = "partial"
	input.Diagnostics = []normalize.HistoryDiagnosticInput{{Code: "SCANNER_EXIT_NONZERO", Scope: "repository"}}
	projection, err := normalize.ProjectHistory(input)
	if err != nil {
		t.Fatalf("project independent states: %v", err)
	}
	if projection.OperationOutcome != "failed" || projection.Report.Status != "partial" {
		t.Fatalf("states conflated: outcome=%q status=%q", projection.OperationOutcome, projection.Report.Status)
	}
}

func TestProjectHistoryCopiesInputsAndReturnedProjection(t *testing.T) {
	input := historyInput()
	projection, err := normalize.ProjectHistory(input)
	if err != nil {
		t.Fatalf("project history: %v", err)
	}
	*input.RepositoryLabel = "mutated source pointer"
	input.Workspaces[0].Path = "mutated source workspace"
	input.Report.Findings[0].Aliases[0] = "mutated source alias"
	input.Report.Findings[0].Component = "mutated source component"
	if *projection.RepositoryLabel != "service" || projection.Workspaces[0].Path != "." || projection.Report.Findings[0].Aliases[0] != "CVE-2024-0001" || projection.Report.Findings[0].ComponentName != "lodash" {
		t.Fatalf("source mutation changed projected value: %#v", projection)
	}
	*projection.RepositoryLabel = "mutated projected pointer"
	projection.Workspaces[0].Path = "mutated projected workspace"
	projection.Report.Findings[0].Aliases[0] = "mutated projected alias"
	if *input.RepositoryLabel != "mutated source pointer" || input.Workspaces[0].Path != "mutated source workspace" || input.Report.Findings[0].Aliases[0] != "mutated source alias" {
		t.Fatal("projection mutation changed source-owned data")
	}
}
func TestProjectHistoryCanonicalizesEquivalentInputWithoutReorderingWorkspaces(t *testing.T) {
	left := historyInput()
	left.Workspaces = append(left.Workspaces, history.Workspace{WorkspaceID: "ws-python", Path: "tools", Ecosystem: "python", PackageManager: "pip", DiscoveryCompleteness: "partial"})
	left.Report.Findings = append(left.Report.Findings, normalize.HistoryFindingInput{Component: "zlib", Version: "1.2.13", WorkspaceID: "ws-python", WorkspacePath: "tools", Ecosystem: "python", TargetID: "CVE-2024-1234", Aliases: []string{"Z-ALIAS", "A-ALIAS", "A-ALIAS"}})
	left.Report.ArtifactDigests = []string{"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"}
	right := historyInput()
	right.Workspaces = append([]history.Workspace(nil), left.Workspaces...)
	right.Workspaces[0], right.Workspaces[1] = right.Workspaces[1], right.Workspaces[0]
	right.Report.Findings = append([]normalize.HistoryFindingInput(nil), left.Report.Findings...)
	right.Report.Findings[0], right.Report.Findings[1] = right.Report.Findings[1], right.Report.Findings[0]
	right.Report.Findings[0].Aliases = []string{"A-ALIAS", "Z-ALIAS"}
	right.Report.ArtifactDigests = []string{left.Report.ArtifactDigests[1], left.Report.ArtifactDigests[0]}
	a, err := normalize.ProjectHistory(left)
	if err != nil {
		t.Fatalf("project left input: %v", err)
	}
	b, err := normalize.ProjectHistory(right)
	if err != nil {
		t.Fatalf("project reordered input: %v", err)
	}
	if !reflect.DeepEqual(a.Report, b.Report) {
		t.Fatalf("report projection depends on input ordering:\nleft %#v\nright %#v", a.Report, b.Report)
	}
	if !reflect.DeepEqual(a.Workspaces, left.Workspaces) || !reflect.DeepEqual(b.Workspaces, right.Workspaces) {
		t.Fatalf("workspace source order was not preserved: left=%#v right=%#v", a.Workspaces, b.Workspaces)
	}
}

func TestProjectHistoryUsesOccurrenceIdentityAndAllowsRepeatedSourceScanID(t *testing.T) {
	first, err := normalize.ProjectHistory(historyInput())
	if err != nil {
		t.Fatalf("project first occurrence: %v", err)
	}
	secondInput := historyInput()
	secondInput.HistoryEntryID = "2d931510-d99f-494a-8c67-87feb05e1595"
	second, err := normalize.ProjectHistory(secondInput)
	if err != nil {
		t.Fatalf("project second occurrence: %v", err)
	}
	if first.Report.SourceScanID != second.Report.SourceScanID || first.HistoryEntryID == second.HistoryEntryID {
		t.Fatalf("source scan identity replaced occurrence identity: first=%#v second=%#v", first, second)
	}
}

func TestProjectHistoryDeduplicatesAndScopesTypedDiagnostics(t *testing.T) {
	input := historyInput()
	input.Diagnostics = []normalize.HistoryDiagnosticInput{
		{Code: "MANIFEST_INVALID", Scope: "workspace", WorkspaceID: new("ws-npm-root")},
		{Code: "MANIFEST_INVALID", Scope: "workspace", WorkspaceID: new("ws-npm-root")},
		{Code: "DISCOVERY_INCOMPLETE", Scope: "repository"},
	}
	projection, err := normalize.ProjectHistory(input)
	if err != nil {
		t.Fatalf("project typed diagnostics: %v", err)
	}
	if len(projection.Diagnostics) != 2 || projection.Diagnostics[0].Code != "DISCOVERY_INCOMPLETE" || projection.Diagnostics[0].Message != "Workspace discovery is incomplete." || projection.Diagnostics[0].WorkspaceID != nil || projection.Diagnostics[1].Code != "MANIFEST_INVALID" || projection.Diagnostics[1].Message != "A workspace manifest is invalid." || projection.Diagnostics[1].WorkspaceID == nil || *projection.Diagnostics[1].WorkspaceID != "ws-npm-root" {
		t.Fatalf("diagnostics did not use the fixed registry, scope and exact-record deduplication: %#v", projection.Diagnostics)
	}
}

func TestValidateProjectionDistinguishesFutureAndCorruptVersions(t *testing.T) {
	input := historyInput()
	input.Report.Findings = append(input.Report.Findings, normalize.HistoryFindingInput{Component: "other", PURL: "pkg:npm/other@1", Version: "1", WorkspaceID: "ws-npm-root", WorkspacePath: ".", Ecosystem: "npm", TargetID: "GHSA-second"})
	projection, err := normalize.ProjectHistory(input)
	if err != nil {
		t.Fatalf("project history: %v", err)
	}
	future := projection
	future.SchemaVersion = "history-v2"
	var futureErr *normalize.ProjectionError
	if err := normalize.ValidateProjection(future); !errors.As(err, &futureErr) || futureErr.Code != "HISTORY_SCHEMA_UNSUPPORTED" {
		t.Fatalf("future projection error = %#v, want HISTORY_SCHEMA_UNSUPPORTED", err)
	}
	unknownMarker := projection
	unknownMarker.SchemaVersion = "future"
	var unknownErr *normalize.ProjectionError
	if err := normalize.ValidateProjection(unknownMarker); !errors.As(err, &unknownErr) || unknownErr.Code != "HISTORY_CORRUPT" {
		t.Fatalf("unknown projection marker error = %#v, want HISTORY_CORRUPT", err)
	}
	assertCorrupt := func(name string, value history.Projection) {
		t.Helper()
		var corruptErr *normalize.ProjectionError
		if err := normalize.ValidateProjection(value); !errors.As(err, &corruptErr) || corruptErr.Code != "HISTORY_CORRUPT" {
			t.Errorf("%s error = %#v, want HISTORY_CORRUPT", name, err)
		}
	}
	badKey := projection
	badKey.Report = new(history.Report)
	*badKey.Report = *projection.Report
	badKey.Report.Findings = append([]history.Finding(nil), projection.Report.Findings...)
	badKey.Report.Findings[0].StableFindingKey = strings.Repeat("0", 64)
	assertCorrupt("tampered stable key", badKey)
	badOrder := projection
	badOrder.Report = new(history.Report)
	*badOrder.Report = *projection.Report
	badOrder.Report.Findings = append([]history.Finding(nil), projection.Report.Findings...)
	badOrder.Report.Findings[0], badOrder.Report.Findings[1] = badOrder.Report.Findings[1], badOrder.Report.Findings[0]
	assertCorrupt("tampered finding order", badOrder)
	badMessage := projection
	badMessage.Diagnostics = []history.Diagnostic{{Code: "DISCOVERY_INCOMPLETE", Scope: "repository", Message: "tampered"}}
	assertCorrupt("tampered diagnostic message", badMessage)
}

func TestProjectHistoryRejectsUnsafeAndUnclassifiedSourceWithoutLeakingIt(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*normalize.HistoryInput)
	}{
		{"absolute workspace path", func(in *normalize.HistoryInput) { in.Workspaces[0].Path = "/private/secret-root" }},
		{"UNC workspace path", func(in *normalize.HistoryInput) { in.Workspaces[0].Path = `\\server\share` }},
		{"drive workspace path", func(in *normalize.HistoryInput) { in.Workspaces[0].Path = "C:/private" }},
		{"backslash workspace path", func(in *normalize.HistoryInput) { in.Workspaces[0].Path = "src\\package.json" }},
		{"dot component workspace path", func(in *normalize.HistoryInput) { in.Workspaces[0].Path = "src/./package.json" }},
		{"traversal workspace path", func(in *normalize.HistoryInput) { in.Workspaces[0].Path = "src/../secret-root" }},
		{"control workspace path", func(in *normalize.HistoryInput) { in.Workspaces[0].Path = "src/\x00file" }},
		{"traversal finding path", func(in *normalize.HistoryInput) { in.Report.Findings[0].WorkspacePath = "../secret-root" }},
		{"unsafe repository label", func(in *normalize.HistoryInput) { in.RepositoryLabel = new("https://user:secret@example.test") }},
		{"credential-bearing provenance", func(in *normalize.HistoryInput) {
			in.Report.Provenance.ScannerName = new("https://user:secret@example.test")
		}},
		{"unclassified source", func(in *normalize.HistoryInput) { in.UnclassifiedSource = true }},
		{"unknown diagnostic", func(in *normalize.HistoryInput) {
			in.Diagnostics = []normalize.HistoryDiagnosticInput{{Code: "RAW_ERROR-secret", Scope: "repository"}}
		}},
		{"invalid operation ID", func(in *normalize.HistoryInput) { in.HistoryEntryID = "not-a-canonical-uuid" }},
		{"missing component display name", func(in *normalize.HistoryInput) { in.Report.Findings[0].Component = "" }},
		{"missing component identity", func(in *normalize.HistoryInput) {
			in.Report.Findings[0].Component = ""
			in.Report.Findings[0].PURL = ""
		}},
		{"missing version", func(in *normalize.HistoryInput) { in.Report.Findings[0].Version = "" }},
		{"missing workspace identity", func(in *normalize.HistoryInput) { in.Report.Findings[0].WorkspaceID = "" }},
		{"missing vulnerability identity", func(in *normalize.HistoryInput) { in.Report.Findings[0].TargetID = "" }},
		{"invalid component PURL", func(in *normalize.HistoryInput) { in.Report.Findings[0].PURL = "https://example.test/package" }},
		{"unsafe provenance path", func(in *normalize.HistoryInput) { in.Report.Provenance.ScannerVersion = new("/private/secret-root") }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := historyInput()
			test.mutate(&input)
			if _, err := normalize.ProjectHistory(input); err == nil {
				t.Fatal("unsafe source unexpectedly projected")
			} else {
				for _, sentinel := range []string{"secret-root", "user:secret", "RAW_ERROR-secret", "example.test"} {
					if strings.Contains(err.Error(), sentinel) {
						t.Fatalf("projection error leaked input sentinel %q: %v", sentinel, err)
					}
				}
			}
		})
	}
}
func TestProjectHistoryAppliesCountLimitsAfterCanonicalization(t *testing.T) {
	digest := func(n int) string { return fmt.Sprintf("%064x", n) }
	input := historyInput()
	input.Report.ArtifactDigests = make([]string, normalize.HistoryMaxArtifactDigests+1)
	for i := range input.Report.ArtifactDigests {
		input.Report.ArtifactDigests[i] = digest(i)
	}
	input.Report.ArtifactDigests[normalize.HistoryMaxArtifactDigests] = input.Report.ArtifactDigests[0]
	if _, err := normalize.ProjectHistory(input); err != nil {
		t.Fatalf("duplicate digest should be deduplicated before its bound: %v", err)
	}
	input.Report.ArtifactDigests = append(input.Report.ArtifactDigests, digest(normalize.HistoryMaxArtifactDigests+1))
	if _, err := normalize.ProjectHistory(input); err == nil {
		t.Fatal("projection exceeding the canonical digest limit was accepted")
	}

	input = historyInput()
	input.Report.Findings[0].Aliases = make([]string, 40)
	for i := range input.Report.Findings[0].Aliases {
		input.Report.Findings[0].Aliases[i] = "CVE-2024-0001"
	}
	if _, err := normalize.ProjectHistory(input); err != nil {
		t.Fatalf("duplicate aliases should canonicalize before the per-finding bound: %v", err)
	}
	input.Report.Findings[0].Aliases = make([]string, 33)
	for i := range input.Report.Findings[0].Aliases {
		input.Report.Findings[0].Aliases[i] = fmt.Sprintf("CVE-%08d", i)
	}
	if _, err := normalize.ProjectHistory(input); err == nil {
		t.Fatal("finding exceeding the canonical alias limit was accepted")
	}
	input = historyInput()
	input.Report = nil
	input.Workspaces = make([]history.Workspace, normalize.HistoryMaxWorkspaces+1)
	for i := range input.Workspaces {
		input.Workspaces[i] = history.Workspace{WorkspaceID: fmt.Sprintf("workspace-%04d", i), Path: ".", Ecosystem: "unknown", PackageManager: "tool", DiscoveryCompleteness: "complete"}
	}
	if _, err := normalize.ProjectHistory(input); err == nil {
		t.Fatal("projection exceeding the workspace limit was accepted")
	}
	input.Workspaces = input.Workspaces[:normalize.HistoryMaxWorkspaces]
	if _, err := normalize.ProjectHistory(input); err != nil {
		t.Fatalf("projection at workspace limit rejected: %v", err)
	}
}

func TestProjectHistoryEnforcesFindingCountLimit(t *testing.T) {
	input := historyInput()
	input.Workspaces = nil
	input.Report.Findings = make([]normalize.HistoryFindingInput, normalize.HistoryMaxFindings)
	for i := range input.Report.Findings {
		input.Report.Findings[i] = normalize.HistoryFindingInput{Component: "pkg", Version: "1", WorkspaceID: "workspace", WorkspacePath: ".", Ecosystem: "npm", TargetID: fmt.Sprintf("ISSUE-%05d", i)}
	}
	if _, err := normalize.ProjectHistory(input); err != nil {
		t.Fatalf("projection at finding limit rejected: %v", err)
	}
	input.Report.Findings = append(input.Report.Findings, normalize.HistoryFindingInput{Component: "pkg", Version: "1", WorkspaceID: "workspace", WorkspacePath: ".", Ecosystem: "npm", TargetID: "ISSUE-over"})
	if _, err := normalize.ProjectHistory(input); err == nil {
		t.Fatal("projection exceeding the finding limit was accepted")
	}
}
func TestProjectHistoryPreservesNullabilityAndComponentFallback(t *testing.T) {
	input := historyInput()
	input.RepositoryLabel = nil
	input.Workspaces = nil
	input.Report.Findings[0].PURL = ""
	input.Report.Findings[0].Component = "lodash"
	input.Report.Findings[0].Severity = ""
	input.Report.Findings[0].FixedVersion = ""
	projection, err := normalize.ProjectHistory(input)
	if err != nil {
		t.Fatalf("project fallback identity: %v", err)
	}
	finding := projection.Report.Findings[0]
	wantKey := normalize.StableFindingKey(normalize.FindingInput{WorkspaceID: "ws-npm-root", ComponentPURL: "lodash", ComponentVersion: "4.17.20", VulnerabilityID: "GHSA-example", VulnerabilityAliases: []string{"CVE-2024-0001"}})
	if projection.RepositoryLabel != nil || projection.Workspaces != nil || finding.ComponentPURL != nil || finding.Severity != nil || finding.FixedVersion != nil || finding.StableFindingKey != wantKey {
		t.Fatalf("nullability or component fallback identity changed: projection=%#v finding=%#v", projection, finding)
	}
}

func TestProjectHistoryNormalizesScannerSeverityCase(t *testing.T) {
	input := historyInput()
	input.Report.Findings[0].Severity = "MODERATE"
	projection, err := normalize.ProjectHistory(input)
	if err != nil {
		t.Fatalf("project scanner severity: %v", err)
	}
	if got := projection.Report.Findings[0].Severity; got == nil || *got != "moderate" {
		t.Fatalf("scanner enum severity = %v, want canonical moderate", got)
	}
}
func TestProjectHistoryDoesNotInventSeverityFromCVSSVector(t *testing.T) {
	input := historyInput()
	input.Report.Findings[0].Severity = "CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:N/I:N/A:L"
	projection, err := normalize.ProjectHistory(input)
	if err != nil {
		t.Fatalf("project supported CVSS source value: %v", err)
	}
	if projection.Report.Findings[0].Severity != nil {
		t.Fatalf("CVSS vector was misrepresented as normalized severity: %q", *projection.Report.Findings[0].Severity)
	}
	input.Report.Findings[0].Severity = "unrecognized-secret"
	if _, err := normalize.ProjectHistory(input); err == nil {
		t.Fatal("unclassified severity source was silently projected")
	}
}

func TestProjectHistoryEnforcesDiagnosticLimitAfterExactDeduplication(t *testing.T) {
	input := historyInput()
	input.Report = nil
	input.Workspaces = make([]history.Workspace, normalize.HistoryMaxDiagnostics)
	input.Diagnostics = make([]normalize.HistoryDiagnosticInput, normalize.HistoryMaxDiagnostics+1)
	for i := range input.Workspaces {
		id := fmt.Sprintf("workspace-%03d", i)
		input.Workspaces[i] = history.Workspace{WorkspaceID: id, Path: ".", Ecosystem: "unknown", PackageManager: "tool", DiscoveryCompleteness: "complete"}
		input.Diagnostics[i] = normalize.HistoryDiagnosticInput{Code: "MANIFEST_INVALID", Scope: "workspace", WorkspaceID: new(id)}
	}
	input.Diagnostics[normalize.HistoryMaxDiagnostics] = input.Diagnostics[0]
	projection, err := normalize.ProjectHistory(input)
	if err != nil || len(projection.Diagnostics) != normalize.HistoryMaxDiagnostics {
		t.Fatalf("exact-record duplicate should canonicalize to diagnostic limit: count=%d err=%v", len(projection.Diagnostics), err)
	}
	input.Diagnostics[normalize.HistoryMaxDiagnostics] = normalize.HistoryDiagnosticInput{Code: "DISCOVERY_INCOMPLETE", Scope: "repository"}
	if _, err := normalize.ProjectHistory(input); err == nil {
		t.Fatal("projection exceeding canonical diagnostic limit was accepted")
	}
}
func TestValidateProjectionEnforcesExactUTF8AndEscapedJSONByteLimit(t *testing.T) {
	const findings = 2000
	const target = normalize.HistoryMaxBytes
	makeProjection := func(aliasLength int) history.Projection {
		p := history.Projection{
			SchemaVersion: history.SchemaVersion, HistoryEntryID: historyEntryID, RecordedAtUS: 1780000000123456,
			OperationOutcome: "completed", RepositoryLabel: new("service"),
			Report: &history.Report{
				SourceSchemaVersion: "v1alpha", SourceScanID: "scan", Status: "complete", RepositoryState: repositoryState,
				Findings: make([]history.Finding, findings), ArtifactDigests: []string{},
			},
			Workspaces: nil, Diagnostics: []history.Diagnostic{},
		}
		aliases := make([]string, 32)
		for i := range aliases {
			aliases[i] = fmt.Sprintf("%02d", i) + strings.Repeat("a", aliasLength-2)
		}
		for i := range p.Report.Findings {
			f := history.Finding{ComponentName: "x", ComponentPURL: new("pkg:npm/x@1"), Version: "1", WorkspaceID: "workspace", RelativePath: ".", Ecosystem: "npm", VulnerabilityID: fmt.Sprintf("ISSUE-%04d", i), Aliases: append([]string(nil), aliases...)}
			f.StableFindingKey = normalize.StableFindingKey(normalize.FindingInput{WorkspaceID: f.WorkspaceID, ComponentPURL: *f.ComponentPURL, ComponentVersion: f.Version, VulnerabilityID: f.VulnerabilityID, VulnerabilityAliases: f.Aliases})
			p.Report.Findings[i] = f
		}
		sort.Slice(p.Report.Findings, func(i, j int) bool {
			return p.Report.Findings[i].StableFindingKey < p.Report.Findings[j].StableFindingKey
		})
		return p
	}
	var exact history.Projection
	exactSize := -1
	for aliasLength := 256; aliasLength >= 2; aliasLength-- {
		candidate := makeProjection(aliasLength)
		encoded, err := json.Marshal(candidate)
		if err != nil {
			t.Fatalf("encode boundary fixture: %v", err)
		}
		deficit := target - len(encoded)
		if deficit >= 0 && deficit <= findings*255 {
			for i := 0; deficit > 0; i++ {
				pad := deficit
				if pad > 255 {
					pad = 255
				}
				candidate.Report.Findings[i].ComponentName += strings.Repeat("a", pad)
				deficit -= pad
			}
			exact = candidate
			encoded, err = json.Marshal(exact)
			if err != nil {
				t.Fatalf("encode exact-size fixture: %v", err)
			}
			exactSize = len(encoded)
			break
		}
	}
	if exactSize != target {
		t.Fatalf("could not construct exact byte boundary fixture: got %d, want %d", exactSize, target)
	}
	if err := normalize.ValidateProjection(exact); err != nil {
		t.Fatalf("projection exactly at byte limit rejected: %v", err)
	}
	oneOver := exact
	oneOver.Report = new(history.Report)
	*oneOver.Report = *exact.Report
	oneOver.Report.Findings = append([]history.Finding(nil), exact.Report.Findings...)
	oneOver.Report.Findings[0].ComponentName = "é" + oneOver.Report.Findings[0].ComponentName[1:]
	encoded, err := json.Marshal(oneOver)
	if err != nil {
		t.Fatalf("encode UTF-8 one-byte-over fixture: %v", err)
	}
	if len(encoded) != target+1 {
		t.Fatalf("UTF-8 boundary change encoded to %d bytes, want %d", len(encoded), target+1)
	}
	if err := normalize.ValidateProjection(oneOver); err == nil {
		t.Fatal("projection one UTF-8 byte above the limit was accepted")
	}
	escapedOver := exact
	escapedOver.Report = new(history.Report)
	*escapedOver.Report = *exact.Report
	escapedOver.Report.Findings = append([]history.Finding(nil), exact.Report.Findings...)
	escapedOver.Report.Findings[0].ComponentName = "\"" + escapedOver.Report.Findings[0].ComponentName[1:]
	escapedBytes, err := json.Marshal(escapedOver)
	if err != nil {
		t.Fatalf("encode escaped one-byte-over fixture: %v", err)
	}
	if len(escapedBytes) != target+1 {
		t.Fatalf("JSON escape boundary change encoded to %d bytes, want %d", len(escapedBytes), target+1)
	}
	if err := normalize.ValidateProjection(escapedOver); err == nil {
		t.Fatal("projection one escaped JSON byte above the limit was accepted")
	}
}
func TestProjectHistoryRefusesNonSuccessfulContextWithoutDiagnosticEvidence(t *testing.T) {
	tests := []struct {
		name    string
		outcome string
		status  string
	}{
		{"failed operation", "failed", "complete"},
		{"partial report", "completed", "partial"},
		{"failed report", "completed", "failed"},
		{"cancelled operation", "cancelled", "complete"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := historyInput()
			input.OperationOutcome = test.outcome
			input.Report.Status = test.status
			if _, err := normalize.ProjectHistory(input); err == nil {
				t.Fatal("incomplete context without typed diagnostic evidence was accepted")
			}
		})
	}
}
func TestValidateProjectionRejectsNullWhereSchemaRequiresArrays(t *testing.T) {
	base, err := normalize.ProjectHistory(historyInput())
	if err != nil {
		t.Fatalf("project valid baseline: %v", err)
	}
	tests := []struct {
		name   string
		mutate func(*history.Projection)
	}{
		{"diagnostics", func(p *history.Projection) { p.Diagnostics = nil }},
		{"findings", func(p *history.Projection) { p.Report.Findings = nil }},
		{"artifact digests", func(p *history.Projection) { p.Report.ArtifactDigests = nil }},
		{"aliases", func(p *history.Projection) { p.Report.Findings[0].Aliases = nil }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := base
			p.Report = new(history.Report)
			*p.Report = *base.Report
			p.Report.Findings = append([]history.Finding(nil), base.Report.Findings...)
			p.Report.ArtifactDigests = append([]string(nil), base.Report.ArtifactDigests...)
			p.Diagnostics = append([]history.Diagnostic(nil), base.Diagnostics...)
			p.Report.Findings[0].Aliases = append([]string(nil), base.Report.Findings[0].Aliases...)
			test.mutate(&p)
			if err := normalize.ValidateProjection(p); err == nil {
				t.Fatal("schema-invalid null array was accepted")
			}
		})
	}
}
func TestProjectHistoryRequiresCancellationDiagnosticForCancelledOutcome(t *testing.T) {
	input := historyInput()
	input.OperationOutcome = "cancelled"
	input.Diagnostics = []normalize.HistoryDiagnosticInput{{Code: "SCANNER_TIMEOUT", Scope: "repository"}}
	if _, err := normalize.ProjectHistory(input); err == nil {
		t.Fatal("cancelled operation with a non-CANCELLED diagnostic was accepted")
	}
	input.Diagnostics[0].Code = "CANCELLED"
	projection, err := normalize.ProjectHistory(input)
	if err != nil {
		t.Fatalf("cancelled operation with CANCELLED evidence was rejected: %v", err)
	}
	projection.Diagnostics[0].Code = "SCANNER_TIMEOUT"
	projection.Diagnostics[0].Message = "The scanner timed out."
	if err := normalize.ValidateProjection(projection); err == nil {
		t.Fatal("persisted cancelled projection with non-CANCELLED evidence was accepted")
	}
}
func TestProjectHistoryRejectsConflictingFindingWorkspaceAssociationsWithoutGraph(t *testing.T) {
	input := historyInput()
	input.Workspaces = nil
	input.Report.Findings = append(input.Report.Findings, normalize.HistoryFindingInput{Component: "other", PURL: "pkg:npm/other@1", Version: "1", WorkspaceID: "ws-npm-root", WorkspacePath: "subdir", Ecosystem: "npm", TargetID: "GHSA-second"})
	if _, err := normalize.ProjectHistory(input); err == nil {
		t.Fatal("findings with one workspace ID mapped to conflicting paths were accepted without graph context")
	}
}
