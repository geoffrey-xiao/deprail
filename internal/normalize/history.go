package normalize

import (
	"encoding/json"
	"regexp"
	"slices"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/geoffrey-xiao/deprail/internal/domain/history"
)

const (
	HistoryMaxBytes           = 16 << 20
	HistoryMaxFindings        = 10000
	HistoryMaxWorkspaces      = 1000
	HistoryMaxDiagnostics     = 128
	HistoryMaxArtifactDigests = 4096
)

type ProjectionError struct {
	Code    string
	Message string
}

func (e *ProjectionError) Error() string { return e.Code + ": " + e.Message }

func historyFailure(corrupt bool) error {
	if corrupt {
		return &ProjectionError{Code: "HISTORY_CORRUPT", Message: "Stored history projection is invalid."}
	}
	return &ProjectionError{Code: "HISTORY_WRITE_FAILED", Message: "History projection could not be safely validated."}
}

type HistoryInput struct {
	HistoryEntryID     string
	RecordedAtUS       int64
	OperationOutcome   string
	RepositoryLabel    *string
	Workspaces         []history.Workspace
	Report             *HistoryReportInput
	Diagnostics        []HistoryDiagnosticInput
	UnclassifiedSource bool
}
type HistoryReportInput struct {
	SchemaVersion   string
	ScanID          string
	Status          string
	RepositoryState string
	Provenance      history.Provenance
	Findings        []HistoryFindingInput
	ArtifactDigests []string
}
type HistoryFindingInput struct {
	Component     string
	PURL          string
	Version       string
	WorkspaceID   string
	WorkspacePath string
	Ecosystem     string
	TargetID      string
	Aliases       []string
	Severity      string
	FixedVersion  string
}
type HistoryDiagnosticInput struct {
	Code, Scope string
	WorkspaceID *string
}

type historyWorkspaceScope struct {
	path, ecosystem string
}

var uuidV4 = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
var hex64 = regexp.MustCompile(`^[0-9a-f]{64}$`)
var projectionVersion = regexp.MustCompile(`^history-v[1-9][0-9]*$`)

var credentialPattern = regexp.MustCompile(`(?i)(password|passwd|token|secret|authorization|api[_-]?key|credential|auth|user|key)\s*[:=]`)

var cvssVectorPattern = regexp.MustCompile(`(?i)^CVSS:[34]\.[0-9]+/[A-Z0-9:/._-]+$`)
var drivePathPattern = regexp.MustCompile(`(^|[[:space:]])[A-Za-z]:`)

func ProjectHistory(in HistoryInput) (history.Projection, error) {
	var zero history.Projection
	if in.UnclassifiedSource || len(in.Workspaces) > HistoryMaxWorkspaces {
		return zero, historyFailure(false)
	}
	if in.Report != nil && len(in.Report.Findings) > HistoryMaxFindings {
		return zero, historyFailure(false)
	}
	diagnosticCapacity := len(in.Diagnostics)
	if diagnosticCapacity > HistoryMaxDiagnostics+1 {
		diagnosticCapacity = HistoryMaxDiagnostics + 1
	}
	p := history.Projection{
		SchemaVersion: history.SchemaVersion, HistoryEntryID: in.HistoryEntryID,
		RecordedAtUS: in.RecordedAtUS, OperationOutcome: in.OperationOutcome,
		RepositoryLabel: cloneString(in.RepositoryLabel),
		Diagnostics:     make([]history.Diagnostic, 0, diagnosticCapacity),
	}
	if in.Workspaces != nil {
		p.Workspaces = append(make([]history.Workspace, 0, len(in.Workspaces)), in.Workspaces...)
	}
	if in.Report != nil {
		src := in.Report
		if src.SchemaVersion != "v1alpha" || !safeIdentifier(src.ScanID, 128) ||
			!oneOf(src.Status, "complete", "partial", "failed") || !hex64.MatchString(src.RepositoryState) {
			return zero, historyFailure(false)
		}
		for _, digest := range src.ArtifactDigests {
			if !hex64.MatchString(digest) {
				return zero, historyFailure(false)
			}
		}
		for _, value := range []*string{src.Provenance.DeprailVersion, src.Provenance.ScannerName, src.Provenance.ScannerVersion, src.Provenance.ScannerDatabaseVersion} {
			if !safeProvenance(value, 128) {
				return zero, historyFailure(false)
			}
		}
		for _, finding := range src.Findings {
			if !validFindingInput(finding) {
				return zero, historyFailure(false)
			}
			for _, alias := range finding.Aliases {
				if !safeIdentifier(alias, 256) {
					return zero, historyFailure(false)
				}
			}
		}
		digests, ok := canonicalStrings(src.ArtifactDigests, HistoryMaxArtifactDigests)
		if !ok {
			return zero, historyFailure(false)
		}
		r := &history.Report{
			SourceSchemaVersion: src.SchemaVersion, SourceScanID: src.ScanID,
			Status: src.Status, RepositoryState: src.RepositoryState,
			Provenance:      cloneProvenance(src.Provenance),
			Findings:        make([]history.Finding, 0, len(src.Findings)),
			ArtifactDigests: digests,
		}
		for _, f := range src.Findings {
			aliases, ok := canonicalStrings(f.Aliases, 32)
			if !ok {
				return zero, historyFailure(false)
			}
			identity := f.PURL
			if identity == "" {
				identity = f.Component
			}
			var purl, fixed *string
			if f.PURL != "" {
				purl = new(f.PURL)
			}
			if f.FixedVersion != "" {
				fixed = new(f.FixedVersion)
			}
			severity, severityOK := projectedSeverity(f.Severity)
			if !severityOK {
				return zero, historyFailure(false)
			}
			key := StableFindingKey(FindingInput{WorkspaceID: f.WorkspaceID, ComponentPURL: identity, ComponentVersion: f.Version, VulnerabilityID: f.TargetID, VulnerabilityAliases: aliases})
			r.Findings = append(r.Findings, history.Finding{StableFindingKey: key, ComponentName: f.Component, ComponentPURL: purl, Version: f.Version, WorkspaceID: f.WorkspaceID, RelativePath: f.WorkspacePath, Ecosystem: f.Ecosystem, VulnerabilityID: f.TargetID, Aliases: aliases, Severity: severity, FixedVersion: fixed})
		}
		sort.Slice(r.Findings, func(i, j int) bool { return r.Findings[i].StableFindingKey < r.Findings[j].StableFindingKey })
		p.Report = r
	}
	type diagnosticKey struct {
		code, scope, workspaceID string
		hasWorkspace             bool
	}
	seenDiagnostics := make(map[diagnosticKey]struct{}, diagnosticCapacity)
	for _, diagnostic := range in.Diagnostics {
		message, ok := diagnosticMessages[diagnostic.Code]
		if !ok || !oneOf(diagnostic.Scope, "repository", "workspace") ||
			diagnostic.Scope == "repository" && diagnostic.WorkspaceID != nil ||
			diagnostic.Scope == "workspace" && (diagnostic.WorkspaceID == nil || !safeIdentifier(*diagnostic.WorkspaceID, 256)) {
			return zero, historyFailure(false)
		}
		key := diagnosticKey{code: diagnostic.Code, scope: diagnostic.Scope}
		if diagnostic.WorkspaceID != nil {
			key.workspaceID, key.hasWorkspace = *diagnostic.WorkspaceID, true
		}
		if _, exists := seenDiagnostics[key]; exists {
			continue
		}
		seenDiagnostics[key] = struct{}{}
		p.Diagnostics = append(p.Diagnostics, history.Diagnostic{
			Code: diagnostic.Code, Scope: diagnostic.Scope, Message: message,
			WorkspaceID: cloneString(diagnostic.WorkspaceID),
		})
		if len(p.Diagnostics) > HistoryMaxDiagnostics {
			return zero, historyFailure(false)
		}
	}
	sort.Slice(p.Diagnostics, func(i, j int) bool { return diagnosticLess(p.Diagnostics[i], p.Diagnostics[j]) })
	if err := validateProjection(p, false); err != nil {
		return zero, err
	}
	return p, nil
}

func ValidateProjection(p history.Projection) error { return validateProjection(p, true) }

func validateProjection(p history.Projection, corrupt bool) error {
	fail := func() error { return historyFailure(corrupt) }
	if p.SchemaVersion != history.SchemaVersion {
		if projectionVersion.MatchString(p.SchemaVersion) {
			return &ProjectionError{Code: "HISTORY_SCHEMA_UNSUPPORTED", Message: "History projection schema is unsupported."}
		}
		return fail()
	}
	if !uuidV4.MatchString(p.HistoryEntryID) || p.RecordedAtUS < 0 || p.RecordedAtUS > 253402300799999999 || !oneOf(p.OperationOutcome, "completed", "failed", "cancelled") || !safeLabel(p.RepositoryLabel, 256) {
		return fail()
	}
	if p.Diagnostics == nil || len(p.Diagnostics) > HistoryMaxDiagnostics || len(p.Workspaces) > HistoryMaxWorkspaces {
		return fail()
	}
	workspaces := make(map[string]historyWorkspaceScope, len(p.Workspaces))
	for _, w := range p.Workspaces {
		if !safeIdentifier(w.WorkspaceID, 256) || !safePath(w.Path) || !validEco(w.Ecosystem) || !safeScalar(w.PackageManager, 128) || !oneOf(w.DiscoveryCompleteness, "complete", "partial", "failed") {
			return fail()
		}
		if _, exists := workspaces[w.WorkspaceID]; exists {
			return fail()
		}
		workspaces[w.WorkspaceID] = historyWorkspaceScope{path: w.Path, ecosystem: w.Ecosystem}
	}
	findingWorkspaces := make(map[string]historyWorkspaceScope)
	if p.Report != nil && !validateReport(*p.Report, workspaces, p.Workspaces != nil, corrupt) {
		return fail()
	}
	if p.Report != nil && p.Workspaces == nil {
		for _, finding := range p.Report.Findings {
			if prior, ok := findingWorkspaces[finding.WorkspaceID]; ok &&
				(prior.path != finding.RelativePath || prior.ecosystem != finding.Ecosystem) {
				return fail()
			}
			findingWorkspaces[finding.WorkspaceID] = historyWorkspaceScope{path: finding.RelativePath, ecosystem: finding.Ecosystem}
		}
	}
	needsEvidence := p.OperationOutcome != "completed" || (p.Report != nil && p.Report.Status != "complete")
	if needsEvidence && len(p.Diagnostics) == 0 {
		return fail()
	}
	cancelledEvidence := false
	for i, d := range p.Diagnostics {
		message, ok := diagnosticMessages[d.Code]
		if !ok || d.Message != message || !oneOf(d.Scope, "repository", "workspace") {
			return fail()
		}
		if d.Code == "CANCELLED" {
			cancelledEvidence = true
		}
		if d.Scope == "repository" {
			if d.WorkspaceID != nil {
				return fail()
			}
		} else if d.WorkspaceID == nil || !safeIdentifier(*d.WorkspaceID, 256) ||
			p.Workspaces != nil && !workspaceExists(workspaces, *d.WorkspaceID) ||
			p.Workspaces == nil && !findingWorkspaceExists(findingWorkspaces, *d.WorkspaceID) {
			return fail()
		}
		if i > 0 && !diagnosticLess(p.Diagnostics[i-1], d) {
			return fail()
		}
	}
	if p.OperationOutcome == "cancelled" && !cancelledEvidence {
		return fail()
	}
	encoded, err := json.Marshal(p)
	if err != nil || !utf8.Valid(encoded) || len(encoded) > HistoryMaxBytes {
		return fail()
	}
	return nil
}

func validFindingInput(finding HistoryFindingInput) bool {
	var purl, fixed *string
	if finding.PURL != "" {
		purl = new(finding.PURL)
	}
	if finding.FixedVersion != "" {
		fixed = new(finding.FixedVersion)
	}
	severity, ok := projectedSeverity(finding.Severity)
	if !ok {
		return false
	}
	identity := finding.PURL
	if identity == "" {
		identity = finding.Component
	}
	return safeComponentName(finding.Component, 256) && safePURL(purl) &&
		safeScalar(finding.Version, 128) && safeIdentifier(finding.WorkspaceID, 256) &&
		safePath(finding.WorkspacePath) && validEco(finding.Ecosystem) &&
		safeIdentifier(finding.TargetID, 256) && validSeverity(severity) &&
		safeNullableScalar(fixed, 128) && safeScalar(identity, 2048)
}

func validateReport(r history.Report, workspaces map[string]historyWorkspaceScope, graph, verifyKeys bool) bool {
	if r.Findings == nil || r.ArtifactDigests == nil || len(r.Findings) > HistoryMaxFindings || len(r.ArtifactDigests) > HistoryMaxArtifactDigests ||
		r.SourceSchemaVersion != "v1alpha" || !safeIdentifier(r.SourceScanID, 128) ||
		!oneOf(r.Status, "complete", "partial", "failed") || !hex64.MatchString(r.RepositoryState) {
		return false
	}
	for _, value := range []*string{r.Provenance.DeprailVersion, r.Provenance.ScannerName, r.Provenance.ScannerVersion, r.Provenance.ScannerDatabaseVersion} {
		if !safeProvenance(value, 128) {
			return false
		}
	}
	for i, digest := range r.ArtifactDigests {
		if !hex64.MatchString(digest) || i > 0 && r.ArtifactDigests[i-1] >= digest {
			return false
		}
	}
	for i, finding := range r.Findings {
		if finding.Aliases == nil || !hex64.MatchString(finding.StableFindingKey) ||
			!safeComponentName(finding.ComponentName, 256) || !safePURL(finding.ComponentPURL) ||
			!safeScalar(finding.Version, 128) || !safeIdentifier(finding.WorkspaceID, 256) ||
			!safePath(finding.RelativePath) || !validEco(finding.Ecosystem) ||
			!safeIdentifier(finding.VulnerabilityID, 256) || len(finding.Aliases) > 32 ||
			!validSeverity(finding.Severity) || !safeNullableScalar(finding.FixedVersion, 128) {
			return false
		}
		for j, alias := range finding.Aliases {
			if !safeIdentifier(alias, 256) || j > 0 && finding.Aliases[j-1] >= alias {
				return false
			}
		}
		identity := finding.ComponentName
		if finding.ComponentPURL != nil {
			identity = *finding.ComponentPURL
		}
		if verifyKeys && StableFindingKey(FindingInput{WorkspaceID: finding.WorkspaceID, ComponentPURL: identity, ComponentVersion: finding.Version, VulnerabilityID: finding.VulnerabilityID, VulnerabilityAliases: finding.Aliases}) != finding.StableFindingKey {
			return false
		}
		if i > 0 && r.Findings[i-1].StableFindingKey >= finding.StableFindingKey {
			return false
		}
		if graph {
			workspace, ok := workspaces[finding.WorkspaceID]
			if !ok || workspace.path != finding.RelativePath || workspace.ecosystem != finding.Ecosystem {
				return false
			}
		}
	}
	return true
}

func projectedSeverity(value string) (*string, bool) {
	if value == "" {
		return nil, true
	}
	normalized := strings.ToLower(value)
	if oneOf(normalized, "unknown", "low", "moderate", "high", "critical") {
		return new(normalized), true
	}
	if cvssVectorPattern.MatchString(value) {
		// The history contract accepts normalized severity labels, not CVSS
		// vectors. Do not infer a label or persist an unmodeled raw value.
		return nil, true
	}
	return nil, false
}

func canonicalStrings(source []string, max int) ([]string, bool) {
	if len(source) <= max {
		values := append(make([]string, 0, len(source)), source...)
		sort.Strings(values)
		return slices.Compact(values), true
	}
	capacity := len(source)
	if capacity > max+1 {
		capacity = max + 1
	}
	values := make([]string, 0, capacity)
	seen := make(map[string]struct{}, capacity)
	for _, value := range source {
		if _, exists := seen[value]; exists {
			continue
		}
		if len(values) == max {
			return nil, false
		}
		seen[value] = struct{}{}
		values = append(values, value)
	}
	sort.Strings(values)
	return values, true
}

func safeBase(s string, max int) bool {
	if s == "" || len(s) > max*utf8.UTFMax || !utf8.ValidString(s) || utf8.RuneCountInString(s) > max {
		return false
	}
	return !hasUnsafeRunes(s) && !credentialPattern.MatchString(s) &&
		!drivePathPattern.MatchString(s) && !strings.Contains(s, "://") &&
		!strings.Contains(s, "\\")
}
func safeScalar(s string, max int) bool { return safeBase(s, max) && !looksAbsolute(s) }
func safeIdentifier(s string, max int) bool {
	return safeScalar(s, max) && !strings.ContainsAny(s, "\\/") && !strings.Contains(s, ":")
}
func safeLabel(value *string, max int) bool {
	return value == nil || safeScalar(*value, max) && !strings.ContainsAny(*value, "/\\")
}
func safeProvenance(value *string, max int) bool {
	return value == nil || safeScalar(*value, max) && !strings.ContainsAny(*value, "\\/")
}
func safeNullableScalar(value *string, max int) bool { return value == nil || safeScalar(*value, max) }
func safeComponentName(value string, max int) bool   { return safeScalar(value, max) }
func safePURL(value *string) bool {
	if value == nil {
		return true
	}
	purl := *value
	if !safeBase(purl, 2048) || !strings.HasPrefix(purl, "pkg:") {
		return false
	}
	lower := strings.ToLower(purl)
	for _, marker := range []string{"?token=", "&token=", "?password=", "&password=", "?secret=", "&secret=", "?username=", "&username="} {
		if strings.Contains(lower, marker) {
			return false
		}
	}
	// The '@' delimiter is valid in npm PURLs both for scoped names and versions.
	return true
}
func validSeverity(value *string) bool {
	return value == nil || oneOf(*value, "unknown", "low", "moderate", "high", "critical")
}
func safePath(path string) bool {
	if !safeBase(path, 4096) || strings.ContainsAny(path, "\\:") || strings.HasPrefix(path, "/") {
		return false
	}
	if path == "." {
		return true
	}
	for part := range strings.SplitSeq(path, "/") {
		if part == "" || part == "." || part == ".." {
			return false
		}
	}
	return true
}
func looksAbsolute(value string) bool {
	if strings.HasPrefix(value, "~") {
		return true
	}
	for i, r := range value {
		if r != '/' {
			continue
		}
		if i == 0 {
			return true
		}
		previous, _ := utf8.DecodeLastRuneInString(value[:i])
		if unicode.IsSpace(previous) || strings.ContainsRune("`\"'=:([]{}", previous) {
			return true
		}
	}
	return false
}
func hasUnsafeRunes(value string) bool {
	for _, r := range value {
		if unicode.IsControl(r) || unicode.Is(unicode.Cf, r) {
			return true
		}
	}
	return false
}
func validEco(value string) bool {
	return oneOf(value, "npm", "pnpm", "yarn", "python", "maven", "gradle", "unknown")
}
func oneOf(value string, options ...string) bool {
	for _, option := range options {
		if value == option {
			return true
		}
	}
	return false
}
func cloneString(value *string) *string {
	if value == nil {
		return nil
	}
	return new(*value)
}
func cloneProvenance(value history.Provenance) history.Provenance {
	return history.Provenance{
		DeprailVersion: cloneString(value.DeprailVersion), ScannerName: cloneString(value.ScannerName),
		ScannerVersion: cloneString(value.ScannerVersion), ScannerDatabaseVersion: cloneString(value.ScannerDatabaseVersion),
	}
}
func workspaceExists(values map[string]historyWorkspaceScope, id string) bool {
	_, ok := values[id]
	return ok
}
func findingWorkspaceExists(values map[string]historyWorkspaceScope, id string) bool {
	_, ok := values[id]
	return ok
}
func diagnosticLess(a, b history.Diagnostic) bool {
	if a.Scope != b.Scope {
		return a.Scope < b.Scope
	}
	aid, bid := "", ""
	if a.WorkspaceID != nil {
		aid = *a.WorkspaceID
	}
	if b.WorkspaceID != nil {
		bid = *b.WorkspaceID
	}
	if aid != bid {
		return aid < bid
	}
	if a.Code != b.Code {
		return a.Code < b.Code
	}
	return a.Message < b.Message
}

var diagnosticMessages = map[string]string{
	"DISCOVERY_FAILED": "Workspace discovery failed.", "DISCOVERY_INCOMPLETE": "Workspace discovery is incomplete.", "MANIFEST_INVALID": "A workspace manifest is invalid.", "WALK_ENTRY_FAILED": "A repository entry could not be inspected.", "PATH_OUTSIDE_ROOT": "A path resolves outside the repository.", "SYMLINK_SKIPPED": "A symlink was not traversed.", "CONFIG_INVALID": "The scan configuration is invalid.", "SCANNER_VERSION_UNSUPPORTED": "The scanner version is unsupported.", "SCANNER_NOT_FOUND": "The scanner executable is unavailable.", "SCANNER_EXIT_NONZERO": "The scanner exited unsuccessfully.", "SCANNER_TIMEOUT": "The scanner timed out.", "SCANNER_OUTPUT_LIMIT": "Scanner output exceeded its limit.", "SCANNER_OUTPUT_INVALID": "Scanner output is invalid.", "ARTIFACT_STORE_FAILED": "Scanner evidence could not be stored.", "FINDING_NORMALIZATION_FAILED": "Findings could not be normalized.", "COMPONENT_IDENTITY_INVALID": "A finding lacked safe component identity.", "CANCELLED": "The scan operation was cancelled.",
}
