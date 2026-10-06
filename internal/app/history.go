package app

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/geoffrey-xiao/deprail/internal/adapter"
	"github.com/geoffrey-xiao/deprail/internal/artifact"
	"github.com/geoffrey-xiao/deprail/internal/discovery"
	domain "github.com/geoffrey-xiao/deprail/internal/domain/history"
	"github.com/geoffrey-xiao/deprail/internal/normalize"
)

const (
	historyDefaultPageSize = 25
	historyMaxPageSize     = 50
	historyMaxResponseSize = 1 << 20
	maxEncodedCursorSize   = 512
	captureTimeout         = 3 * time.Second
)

// HistoryWriter is the least-privilege capture port.
type HistoryWriter interface {
	Append(context.Context, domain.Entry) error
}

// HistoryReader is the read-only query port.
type HistoryReader interface {
	Get(context.Context, string) (domain.Entry, error)
	List(context.Context, domain.Page) (domain.PageResult, error)
}

// ArtifactVerifier confirms a content-addressed artifact without returning bytes.
type ArtifactVerifier interface {
	Verify(context.Context, string) error
}

// CaptureSink receives one terminal capture for a selected scan operation.
type CaptureSink interface {
	Capture(context.Context, CaptureInput) (SaveResult, error)
}

// SaveResult identifies the exact occurrence durably saved by a capture sink.
type SaveResult struct {
	HistoryEntryID string
	Saved          bool
}

// CaptureInput contains same-operation context. Report.Errors and raw paths are
// intentionally never copied into the persisted projection.
type CaptureInput struct {
	HistoryEntryID     string
	RecordedAtUS       int64
	Graph              *discovery.ProjectGraph
	Report             *ScanReport
	OperationError     error
	Diagnostics        []normalize.HistoryDiagnosticInput
	Provenance         domain.Provenance
	UnclassifiedSource bool
}

// HistoryErrorCode is a stable application-level history or query failure.
type HistoryErrorCode string

const (
	HistoryWriteFailed       HistoryErrorCode = "HISTORY_WRITE_FAILED"
	HistoryUnavailable       HistoryErrorCode = "HISTORY_UNAVAILABLE"
	HistorySchemaUnsupported HistoryErrorCode = "HISTORY_SCHEMA_UNSUPPORTED"
	HistoryMigrationFailed   HistoryErrorCode = "HISTORY_MIGRATION_FAILED"
	HistoryCorrupt           HistoryErrorCode = "HISTORY_CORRUPT"
	HistoryEntryNotFound     HistoryErrorCode = "HISTORY_ENTRY_NOT_FOUND"
	HistoryRequestInvalid    HistoryErrorCode = "API_REQUEST_INVALID"
	HistoryResponseTooLarge  HistoryErrorCode = "API_RESPONSE_TOO_LARGE"
)

// HistoryError exposes a stable code and never formats storage or path details.
type HistoryError struct {
	Code  HistoryErrorCode
	cause error
}

func (e *HistoryError) Error() string { return string(e.Code) }
func (e *HistoryError) Unwrap() error { return e.cause }

// ScanCaptureError keeps persistence failure separate from the scan's primary
// error while preserving both for typed callers.
type ScanCaptureError struct {
	OperationErr   error
	PersistenceErr error
}

func (e *ScanCaptureError) Error() string {
	if e.OperationErr != nil {
		return e.OperationErr.Error()
	}
	return e.PersistenceErr.Error()
}

func (e *ScanCaptureError) Unwrap() []error {
	result := make([]error, 0, 2)
	if e.OperationErr != nil {
		result = append(result, e.OperationErr)
	}
	if e.PersistenceErr != nil {
		result = append(result, e.PersistenceErr)
	}
	return result
}

// HistoryService implements capture and read-only history use cases.
type HistoryService struct {
	Writer    HistoryWriter
	Reader    HistoryReader
	Artifacts ArtifactVerifier
}

// Capture validates and commits one projection and its exact digest references.
func (s *HistoryService) Capture(ctx context.Context, input CaptureInput) (SaveResult, error) {
	if s == nil || s.Writer == nil {
		return SaveResult{}, &HistoryError{Code: HistoryWriteFailed}
	}
	projection, err := normalize.ProjectHistory(historyInput(input))
	if err != nil {
		return SaveResult{}, &HistoryError{Code: HistoryWriteFailed, cause: err}
	}
	refs := []string{}
	if projection.Report != nil {
		refs = append(refs, projection.Report.ArtifactDigests...)
	}
	if err := s.Writer.Append(ctx, domain.Entry{Projection: projection, ArtifactRefs: refs}); err != nil {
		return SaveResult{}, &HistoryError{Code: HistoryWriteFailed, cause: err}
	}
	return SaveResult{HistoryEntryID: projection.HistoryEntryID, Saved: true}, nil
}

func historyInput(input CaptureInput) normalize.HistoryInput {
	outcome := "completed"
	if input.OperationError != nil {
		outcome = "failed"
		if errors.Is(input.OperationError, context.Canceled) {
			outcome = "cancelled"
		}
	}
	var workspaces []domain.Workspace
	if input.Graph != nil {
		workspaces = make([]domain.Workspace, 0, len(input.Graph.Workspaces))
		for _, workspace := range input.Graph.Workspaces {
			workspaces = append(workspaces, domain.Workspace{
				WorkspaceID: workspace.WorkspaceID, Path: workspace.RelativePath,
				Ecosystem: string(workspace.Ecosystem), PackageManager: workspace.PackageManager,
				DiscoveryCompleteness: string(workspace.Completeness),
			})
		}
	}
	var report *normalize.HistoryReportInput
	var repositoryLabel *string
	if input.Report != nil {
		source := input.Report
		repositoryLabel = normalize.HistoryRepositoryLabel(source.RepositoryIdentity.Repository)
		findings := make([]normalize.HistoryFindingInput, 0, len(source.Findings))
		for _, finding := range source.Findings {
			findings = append(findings, normalize.HistoryFindingInput{
				Component: finding.Component, PURL: finding.PURL, Version: finding.Version,
				WorkspaceID: finding.WorkspaceID, WorkspacePath: finding.WorkspacePath,
				Ecosystem: finding.Ecosystem, TargetID: finding.TargetID,
				Aliases: append([]string(nil), finding.Aliases...), Severity: finding.Severity,
				FixedVersion: finding.Fixed,
			})
		}
		report = &normalize.HistoryReportInput{
			SchemaVersion: source.SchemaVersion, ScanID: source.ScanID, Status: string(source.Status),
			RepositoryState: source.RepositoryState, Provenance: input.Provenance,
			Findings: findings, ArtifactDigests: append([]string{}, source.ArtifactDigests...),
		}
	}
	return normalize.HistoryInput{
		HistoryEntryID: input.HistoryEntryID, RecordedAtUS: input.RecordedAtUS,
		OperationOutcome: outcome, RepositoryLabel: repositoryLabel, Workspaces: workspaces,
		Report: report, Diagnostics: append([]normalize.HistoryDiagnosticInput{}, input.Diagnostics...),
		UnclassifiedSource: input.UnclassifiedSource,
	}
}

// HistorySummary is an allowlisted list projection; no filesystem root is returned.
type HistorySummary struct {
	HistoryEntryID            string  `json:"historyEntryID"`
	RecordedAt                string  `json:"recordedAt"`
	SourceScanID              *string `json:"sourceScanID"`
	OperationOutcome          string  `json:"operationOutcome"`
	ReportStatus              *string `json:"reportStatus"`
	RepositoryLabel           *string `json:"repositoryLabel"`
	WorkspaceCount            *int    `json:"workspaceCount"`
	FindingCount              *int    `json:"findingCount"`
	SourceReportSchemaVersion *string `json:"sourceReportSchemaVersion"`
	cursor                    domain.Cursor
}

type HistoryPageRequest struct {
	Limit    int
	Before   *domain.Cursor
	MaxBytes int
}

type HistoryPage struct {
	Items   []HistorySummary
	Next    *domain.Cursor
	HasMore bool
}

// HistoryReport contains only source metadata approved for history detail.
type HistoryReport struct {
	SourceSchemaVersion string            `json:"sourceSchemaVersion"`
	RepositoryState     string            `json:"repositoryState"`
	Provenance          HistoryProvenance `json:"provenance"`
}

type HistoryProvenance struct {
	DeprailVersion         *string `json:"deprailVersion"`
	ScannerName            *string `json:"scannerName"`
	ScannerVersion         *string `json:"scannerVersion"`
	ScannerDatabaseVersion *string `json:"scannerDatabaseVersion"`
}

type HistoryDiagnostic struct {
	Code        string  `json:"code"`
	Message     string  `json:"message"`
	Scope       string  `json:"scope"`
	WorkspaceID *string `json:"workspaceID,omitempty"`
}

type ArtifactReference struct {
	Digest    string `json:"digest"`
	Integrity string `json:"integrity"`
}

type HistoryDetail struct {
	Summary            HistorySummary      `json:"summary"`
	Report             *HistoryReport      `json:"report"`
	Diagnostics        []HistoryDiagnostic `json:"diagnostics"`
	ArtifactReferences []ArtifactReference `json:"artifactReferences"`
}

type WorkspaceView struct {
	WorkspaceID           string `json:"workspaceID"`
	RelativePath          string `json:"relativePath"`
	Ecosystem             string `json:"ecosystem"`
	PackageManager        string `json:"packageManager"`
	DiscoveryCompleteness string `json:"discoveryCompleteness"`
}

type WorkspaceCursor struct {
	WorkspaceID   string
	Ordinal       uint32
	WorkspaceHash [sha256.Size]byte
	Encoded       bool
}

type WorkspacePageRequest struct {
	HistoryEntryID string
	Limit          int
	Before         *WorkspaceCursor
	MaxBytes       int
}

type WorkspacePage struct {
	Available bool
	Reason    string
	Items     []WorkspaceView
	Next      *WorkspaceCursor
	HasMore   bool
}

type FindingView struct {
	StableFindingKey string   `json:"stableFindingKey"`
	ComponentName    string   `json:"componentName"`
	ComponentPURL    *string  `json:"componentPURL"`
	Version          string   `json:"version"`
	WorkspaceID      string   `json:"workspaceID"`
	RelativePath     string   `json:"relativePath"`
	Ecosystem        string   `json:"ecosystem"`
	VulnerabilityID  string   `json:"vulnerabilityID"`
	Aliases          []string `json:"aliases"`
	Severity         *string  `json:"severity"`
	FixedVersion     *string  `json:"fixedVersion"`
}

type FindingCursor struct{ StableFindingKey string }

type FindingPageRequest struct {
	HistoryEntryID string
	Limit          int
	Before         *FindingCursor
	MaxBytes       int
}

type FindingPage struct {
	Available    bool
	ReportStatus *string
	Reason       string
	Items        []FindingView
	Next         *FindingCursor
	HasMore      bool
}

// List returns an ordered whole-record page and an application cursor. The byte
// budget reserves the maximum wire cursor; HTTP still validates its final envelope.
func (s *HistoryService) List(ctx context.Context, request HistoryPageRequest) (HistoryPage, error) {
	limit, maxBytes, err := pageBounds(request.Limit, request.MaxBytes)
	if err != nil || request.Before != nil && !validHistoryCursor(*request.Before) {
		return HistoryPage{}, requestError()
	}
	if s == nil || s.Reader == nil {
		return HistoryPage{}, &HistoryError{Code: HistoryUnavailable}
	}
	stored, err := s.Reader.List(ctx, domain.Page{Limit: limit, Before: request.Before})
	if err != nil {
		return HistoryPage{}, readError(err)
	}
	items := make([]HistorySummary, 0, len(stored.Entries))
	for _, summary := range stored.Entries {
		items = append(items, summarize(summary))
	}
	count, trimmed, err := fitPage(items, maxBytes, `{"items":[],"nextCursor":null}`)
	if err != nil {
		return HistoryPage{}, err
	}
	items = items[:count]
	hasMore := trimmed || stored.HasMore
	page := HistoryPage{Items: items, HasMore: hasMore}
	if hasMore {
		if len(items) == 0 {
			return HistoryPage{}, corruptError(nil)
		}
		cursor := items[len(items)-1].cursor
		page.Next = &cursor
	}
	return page, nil
}

// Detail returns only a validated allowlist and artifact integrity metadata.
func (s *HistoryService) Detail(ctx context.Context, id string) (HistoryDetail, error) {
	entry, err := s.getEntry(ctx, id)
	if err != nil {
		return HistoryDetail{}, err
	}
	p := entry.Projection
	detail := HistoryDetail{
		Summary: summarizeProjection(p), Diagnostics: make([]HistoryDiagnostic, 0, len(p.Diagnostics)),
		ArtifactReferences: []ArtifactReference{},
	}
	if p.Report != nil {
		detail.Report = &HistoryReport{
			SourceSchemaVersion: p.Report.SourceSchemaVersion, RepositoryState: p.Report.RepositoryState,
			Provenance: HistoryProvenance{
				DeprailVersion:         cloneString(p.Report.Provenance.DeprailVersion),
				ScannerName:            cloneString(p.Report.Provenance.ScannerName),
				ScannerVersion:         cloneString(p.Report.Provenance.ScannerVersion),
				ScannerDatabaseVersion: cloneString(p.Report.Provenance.ScannerDatabaseVersion),
			},
		}
	}
	for _, diagnostic := range p.Diagnostics {
		detail.Diagnostics = append(detail.Diagnostics, HistoryDiagnostic{
			Code: diagnostic.Code, Message: diagnostic.Message, Scope: diagnostic.Scope,
			WorkspaceID: cloneString(diagnostic.WorkspaceID),
		})
	}
	refs, err := s.artifactReferences(ctx, entry)
	if err != nil {
		return HistoryDetail{}, err
	}
	detail.ArtifactReferences = refs
	return detail, nil
}

// Workspaces returns a deterministic page or an explicit unavailable state.
func (s *HistoryService) Workspaces(ctx context.Context, request WorkspacePageRequest) (WorkspacePage, error) {
	limit, maxBytes, err := pageBounds(request.Limit, request.MaxBytes)
	if err != nil || request.Before != nil && !validWorkspaceCursor(*request.Before) {
		return WorkspacePage{}, requestError()
	}
	entry, err := s.getEntry(ctx, request.HistoryEntryID)
	if err != nil {
		return WorkspacePage{}, err
	}
	if entry.Projection.Workspaces == nil {
		return WorkspacePage{Available: false, Reason: "source_graph_unavailable"}, nil
	}
	workspaces := append([]domain.Workspace(nil), entry.Projection.Workspaces...)
	sort.Slice(workspaces, func(i, j int) bool { return workspaces[i].WorkspaceID < workspaces[j].WorkspaceID })
	start := 0
	if request.Before != nil {
		if request.Before.Encoded {
			if uint64(request.Before.Ordinal) >= uint64(len(workspaces)) || sha256.Sum256([]byte(workspaces[request.Before.Ordinal].WorkspaceID)) != request.Before.WorkspaceHash {
				return WorkspacePage{}, requestError()
			}
			start = int(request.Before.Ordinal) + 1
		} else {
			start = sort.Search(len(workspaces), func(i int) bool { return workspaces[i].WorkspaceID > request.Before.WorkspaceID })
		}
	}
	end := min(start+limit, len(workspaces))
	items := make([]WorkspaceView, 0, end-start)
	for _, workspace := range workspaces[start:end] {
		items = append(items, WorkspaceView{
			WorkspaceID: workspace.WorkspaceID, RelativePath: workspace.Path,
			Ecosystem: workspace.Ecosystem, PackageManager: workspace.PackageManager,
			DiscoveryCompleteness: workspace.DiscoveryCompleteness,
		})
	}
	trimmedCount, trimmed, err := fitPage(items, maxBytes, `{"collectionState":"available","items":[],"nextCursor":null}`)
	if err != nil {
		return WorkspacePage{}, err
	}
	items = items[:trimmedCount]
	hasMore := trimmed || end < len(workspaces)
	page := WorkspacePage{Available: true, Items: items, HasMore: hasMore}
	if hasMore {
		lastIndex := start + len(items) - 1
		lastWorkspaceID := items[len(items)-1].WorkspaceID
		page.Next = &WorkspaceCursor{
			Ordinal: uint32(lastIndex), WorkspaceHash: sha256.Sum256([]byte(lastWorkspaceID)), Encoded: true,
		}
	}
	return page, nil
}

// Findings returns a deterministic page or an explicit unavailable state.
func (s *HistoryService) Findings(ctx context.Context, request FindingPageRequest) (FindingPage, error) {
	limit, maxBytes, err := pageBounds(request.Limit, request.MaxBytes)
	if err != nil || request.Before != nil && !validFindingCursor(*request.Before) {
		return FindingPage{}, requestError()
	}
	entry, err := s.getEntry(ctx, request.HistoryEntryID)
	if err != nil {
		return FindingPage{}, err
	}
	if entry.Projection.Report == nil {
		return FindingPage{Available: false, Reason: "source_report_unavailable"}, nil
	}
	findings := append([]domain.Finding(nil), entry.Projection.Report.Findings...)
	sort.Slice(findings, func(i, j int) bool { return findings[i].StableFindingKey < findings[j].StableFindingKey })
	start := 0
	if request.Before != nil {
		start = sort.Search(len(findings), func(i int) bool { return findings[i].StableFindingKey > request.Before.StableFindingKey })
	}
	end := min(start+limit, len(findings))
	items := make([]FindingView, 0, end-start)
	for _, finding := range findings[start:end] {
		items = append(items, FindingView{
			StableFindingKey: finding.StableFindingKey, ComponentName: finding.ComponentName,
			ComponentPURL: cloneString(finding.ComponentPURL), Version: finding.Version,
			WorkspaceID: finding.WorkspaceID, RelativePath: finding.RelativePath,
			Ecosystem: finding.Ecosystem, VulnerabilityID: finding.VulnerabilityID,
			Aliases: append([]string{}, finding.Aliases...), Severity: cloneString(finding.Severity),
			FixedVersion: cloneString(finding.FixedVersion),
		})
	}
	emptyEnvelope := fmt.Sprintf(`{"collectionState":"available","reportStatus":%q,"items":[],"nextCursor":null}`, entry.Projection.Report.Status)
	trimmedCount, trimmed, err := fitPage(items, maxBytes, emptyEnvelope)
	if err != nil {
		return FindingPage{}, err
	}
	items = items[:trimmedCount]
	hasMore := trimmed || end < len(findings)
	status := entry.Projection.Report.Status
	page := FindingPage{Available: true, ReportStatus: &status, Items: items, HasMore: hasMore}
	if hasMore {
		if len(items) == 0 {
			return FindingPage{}, corruptError(nil)
		}
		page.Next = &FindingCursor{StableFindingKey: items[len(items)-1].StableFindingKey}
	}
	return page, nil
}

// ArtifactIntegrity verifies only digests referenced by the selected entry.
func (s *HistoryService) ArtifactIntegrity(ctx context.Context, id string) ([]ArtifactReference, error) {
	entry, err := s.getEntry(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.artifactReferences(ctx, entry)
}

func (s *HistoryService) getEntry(ctx context.Context, id string) (domain.Entry, error) {
	if !validHistoryID(id) {
		return domain.Entry{}, requestError()
	}
	if s == nil || s.Reader == nil {
		return domain.Entry{}, &HistoryError{Code: HistoryUnavailable}
	}
	entry, err := s.Reader.Get(ctx, id)
	if err != nil {
		return domain.Entry{}, readError(err)
	}
	if err := normalize.ValidateProjection(entry.Projection); err != nil {
		return domain.Entry{}, projectionReadError(err)
	}
	var want []string
	if entry.Projection.Report != nil {
		want = entry.Projection.Report.ArtifactDigests
	}
	if len(want) != len(entry.ArtifactRefs) {
		return domain.Entry{}, corruptError(nil)
	}
	for i := range want {
		if want[i] != entry.ArtifactRefs[i] {
			return domain.Entry{}, corruptError(nil)
		}
	}
	return entry, nil
}

func (s *HistoryService) artifactReferences(ctx context.Context, entry domain.Entry) ([]ArtifactReference, error) {
	refs := make([]ArtifactReference, 0, len(entry.ArtifactRefs))
	for _, digest := range entry.ArtifactRefs {
		if !validDigest(digest) {
			return nil, corruptError(nil)
		}
		integrity := "unavailable"
		if s.Artifacts != nil {
			err := s.Artifacts.Verify(ctx, digest)
			if err == nil {
				integrity = "verified"
			} else if ctx.Err() != nil {
				return nil, ctx.Err()
			} else if artifact.IsCode(err, artifact.ErrNotFound) {
				integrity = "missing"
			} else if artifact.IsCode(err, artifact.ErrIntegrity) {
				integrity = "digest_mismatch"
			}
		}
		refs = append(refs, ArtifactReference{Digest: digest, Integrity: integrity})
	}
	return refs, nil
}

func summarize(summary domain.Summary) HistorySummary {
	return HistorySummary{
		HistoryEntryID: summary.HistoryEntryID,
		RecordedAt:     time.UnixMicro(summary.RecordedAtUS).UTC().Format("2006-01-02T15:04:05.000000Z"),
		SourceScanID:   cloneString(summary.SourceScanID), OperationOutcome: summary.OperationOutcome,
		ReportStatus: cloneString(summary.ReportStatus), RepositoryLabel: cloneString(summary.RepositoryLabel),
		WorkspaceCount: cloneInt(summary.WorkspaceCount), FindingCount: cloneInt(summary.FindingCount),
		SourceReportSchemaVersion: cloneString(summary.SourceReportSchemaVersion),
		cursor:                    domain.Cursor{RecordedAtUS: summary.RecordedAtUS, HistoryEntryID: summary.HistoryEntryID},
	}
}

func summarizeProjection(projection domain.Projection) HistorySummary {
	summary := domain.Summary{
		HistoryEntryID: projection.HistoryEntryID, RecordedAtUS: projection.RecordedAtUS,
		OperationOutcome: projection.OperationOutcome, RepositoryLabel: projection.RepositoryLabel,
	}
	if projection.Workspaces != nil {
		count := len(projection.Workspaces)
		summary.WorkspaceCount = &count
	}
	if projection.Report != nil {
		status, findings := projection.Report.Status, len(projection.Report.Findings)
		sourceID, sourceSchema := projection.Report.SourceScanID, projection.Report.SourceSchemaVersion
		summary.ReportStatus, summary.FindingCount = &status, &findings
		summary.SourceScanID, summary.SourceReportSchemaVersion = &sourceID, &sourceSchema
	}
	return summarize(summary)
}

func pageBounds(limit, maxBytes int) (int, int, error) {
	if limit == 0 {
		limit = historyDefaultPageSize
	}
	if limit < 1 || limit > historyMaxPageSize {
		return 0, 0, requestError()
	}
	if maxBytes == 0 {
		maxBytes = historyMaxResponseSize
	}
	if maxBytes < 600 || maxBytes > historyMaxResponseSize {
		return 0, 0, requestError()
	}
	return limit, maxBytes, nil
}

func fitPage[T any](items []T, maxBytes int, emptyEnvelope string) (int, bool, error) {
	const emptyArrayBytes = 2
	const nullCursorBytes = 4
	const quotedCursorBytes = maxEncodedCursorSize + 2
	reserved := len(emptyEnvelope) - emptyArrayBytes - nullCursorBytes + quotedCursorBytes
	budget := maxBytes - reserved
	if budget < 0 {
		return 0, false, responseTooLargeError()
	}
	used := 0
	for i := range items {
		encoded, err := json.Marshal(items[i])
		if err != nil {
			return 0, false, corruptError(err)
		}
		cost := len(encoded)
		if i > 0 {
			cost++
		}
		if cost > budget-used {
			if i == 0 {
				return 0, false, responseTooLargeError()
			}
			return i, true, nil
		}
		used += cost
	}
	return len(items), false, nil
}

func readError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return err
	}
	var projectionErr *normalize.ProjectionError
	if errors.As(err, &projectionErr) {
		return projectionReadError(err)
	}
	var coded interface{ HistoryErrorCode() string }
	if errors.As(err, &coded) {
		switch HistoryErrorCode(coded.HistoryErrorCode()) {
		case HistoryEntryNotFound, HistorySchemaUnsupported, HistoryMigrationFailed, HistoryCorrupt, HistoryUnavailable:
			return &HistoryError{Code: HistoryErrorCode(coded.HistoryErrorCode()), cause: err}
		}
	}
	return &HistoryError{Code: HistoryUnavailable, cause: err}
}

func projectionReadError(err error) error {
	var projectionErr *normalize.ProjectionError
	if errors.As(err, &projectionErr) {
		switch HistoryErrorCode(projectionErr.Code) {
		case HistorySchemaUnsupported:
			return &HistoryError{Code: HistorySchemaUnsupported, cause: err}
		case HistoryCorrupt:
			return &HistoryError{Code: HistoryCorrupt, cause: err}
		}
	}
	return corruptError(err)
}

func requestError() error            { return &HistoryError{Code: HistoryRequestInvalid} }
func corruptError(cause error) error { return &HistoryError{Code: HistoryCorrupt, cause: cause} }
func responseTooLargeError() error   { return &HistoryError{Code: HistoryResponseTooLarge} }

func validHistoryCursor(cursor domain.Cursor) bool {
	return cursor.RecordedAtUS >= 0 && cursor.RecordedAtUS <= 253402300799999999 && validHistoryID(cursor.HistoryEntryID)
}

func validHistoryID(id string) bool {
	if len(id) != 36 {
		return false
	}
	for i, value := range id {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			if value != '-' {
				return false
			}
			continue
		}
		if !(value >= '0' && value <= '9' || value >= 'a' && value <= 'f') {
			return false
		}
	}
	return id[14] == '4' && strings.ContainsRune("89ab", rune(id[19]))
}

func validWorkspaceCursor(cursor WorkspaceCursor) bool {
	if cursor.Encoded {
		return cursor.WorkspaceID == ""
	}
	return safeCursorText(cursor.WorkspaceID, 256) && !strings.ContainsAny(cursor.WorkspaceID, "/\\:")
}

func validFindingCursor(cursor FindingCursor) bool { return validDigest(cursor.StableFindingKey) }

func validDigest(value string) bool {
	if len(value) != 64 {
		return false
	}
	for _, char := range value {
		if !(char >= '0' && char <= '9' || char >= 'a' && char <= 'f') {
			return false
		}
	}
	return true
}

func safeCursorText(value string, maxRunes int) bool {
	if value == "" || len(value) > maxRunes*utf8.UTFMax || !utf8.ValidString(value) || utf8.RuneCountInString(value) > maxRunes {
		return false
	}
	for _, char := range value {
		if unicode.IsControl(char) || unicode.Is(unicode.Cf, char) {
			return false
		}
	}
	return true
}

func cloneString(value *string) *string {
	if value == nil {
		return nil
	}
	return new(*value)
}

func cloneInt(value *int) *int {
	if value == nil {
		return nil
	}
	return new(*value)
}

func randomHistoryEntryID() (string, error) {
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", err
	}
	raw[6] = raw[6]&0x0f | 0x40
	raw[8] = raw[8]&0x3f | 0x80
	var encoded [36]byte
	hex.Encode(encoded[0:8], raw[0:4])
	encoded[8] = '-'
	hex.Encode(encoded[9:13], raw[4:6])
	encoded[13] = '-'
	hex.Encode(encoded[14:18], raw[6:8])
	encoded[18] = '-'
	hex.Encode(encoded[19:23], raw[8:10])
	encoded[23] = '-'
	hex.Encode(encoded[24:36], raw[10:16])
	return string(encoded[:]), nil
}

func safeDiagnosticCode(err error) (string, bool) {
	if errors.Is(err, context.Canceled) {
		return "CANCELLED", true
	}
	var scannerErr *adapter.Error
	if errors.As(err, &scannerErr) {
		switch scannerErr.Code {
		case adapter.ErrInvalidPlan:
			return "CONFIG_INVALID", true
		case adapter.ErrUnsupportedTarget:
			return "SCANNER_VERSION_UNSUPPORTED", true
		case adapter.ErrScannerNotFound:
			return "SCANNER_NOT_FOUND", true
		case adapter.ErrExecution:
			return "SCANNER_EXIT_NONZERO", true
		case adapter.ErrTimeout:
			return "SCANNER_TIMEOUT", true
		case adapter.ErrOutputLimit:
			return "SCANNER_OUTPUT_LIMIT", true
		case adapter.ErrInvalidOutput:
			return "SCANNER_OUTPUT_INVALID", true
		}
	}
	return "", false
}

type scanCaptureState struct {
	recordedAtUS    int64
	graph           *discovery.ProjectGraph
	reportAvailable bool
	diagnostics     []normalize.HistoryDiagnosticInput
	provenance      domain.Provenance
	unclassified    bool
}

func emitScanEvent(state *scanCaptureState, sink EventSink, event Event) error {
	err := sink.Emit(event)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			state.addDiagnostic("CANCELLED", "repository", nil)
		} else {
			state.markUnclassified()
		}
	}
	return err
}

func (s *scanCaptureState) markUnclassified() {
	if s != nil {
		s.unclassified = true
	}
}

func knownGraphDiagnostic(code string) bool {
	switch code {
	case "DISCOVERY_INCOMPLETE", "MANIFEST_INVALID", "WALK_ENTRY_FAILED", "PATH_OUTSIDE_ROOT", "SYMLINK_SKIPPED":
		return true
	default:
		return false
	}
}

func (s *scanCaptureState) addDiagnostic(code, scope string, workspaceID *string) {
	if s == nil {
		return
	}
	if !knownGraphDiagnostic(code) && !validDiagnosticCode(code) {
		s.unclassified = true
		return
	}
	s.diagnostics = append(s.diagnostics, normalize.HistoryDiagnosticInput{Code: code, Scope: scope, WorkspaceID: cloneString(workspaceID)})
}

func validDiagnosticCode(code string) bool {
	switch code {
	case "DISCOVERY_FAILED", "CONFIG_INVALID", "SCANNER_VERSION_UNSUPPORTED", "SCANNER_NOT_FOUND", "SCANNER_EXIT_NONZERO", "SCANNER_TIMEOUT", "SCANNER_OUTPUT_LIMIT", "SCANNER_OUTPUT_INVALID", "ARTIFACT_STORE_FAILED", "FINDING_NORMALIZATION_FAILED", "COMPONENT_IDENTITY_INVALID", "CANCELLED":
		return true
	default:
		return false
	}
}

func (s *scanCaptureState) addGraph(graph discovery.ProjectGraph) {
	if s == nil {
		return
	}
	for _, workspace := range graph.Workspaces {
		workspaceID := workspace.WorkspaceID
		for _, diagnostic := range workspace.Diagnostics {
			if !knownGraphDiagnostic(diagnostic.Code) {
				s.unclassified = true
				continue
			}
			s.addDiagnostic(diagnostic.Code, "workspace", &workspaceID)
		}
	}
	for _, diagnostic := range graph.Diagnostics {
		if !knownGraphDiagnostic(diagnostic.Code) {
			s.unclassified = true
			continue
		}
		workspaceID := graphDiagnosticWorkspace(graph.Workspaces, diagnostic.Scope)
		if workspaceID == nil {
			s.addDiagnostic(diagnostic.Code, "repository", nil)
		} else {
			s.addDiagnostic(diagnostic.Code, "workspace", workspaceID)
		}
	}
}

func graphDiagnosticWorkspace(workspaces []discovery.Workspace, path string) *string {
	var best *discovery.Workspace
	for i := range workspaces {
		workspace := &workspaces[i]
		root := workspace.RelativePath
		matches := path == root || root == "." || strings.HasPrefix(path, strings.TrimSuffix(root, "/")+"/")
		if matches && (best == nil || len(root) > len(best.RelativePath)) {
			best = workspace
		}
	}
	if best == nil || !safeCursorText(best.WorkspaceID, 256) {
		return nil
	}
	return new(best.WorkspaceID)
}

func (s *scanCaptureState) addBoundaryError(err error, kind string) {
	if s == nil || err == nil {
		return
	}
	if code, ok := safeDiagnosticCode(err); ok {
		s.addDiagnostic(code, "repository", nil)
		return
	}
	switch kind {
	case "discover":
		s.addDiagnostic("DISCOVERY_FAILED", "repository", nil)
	case "normalize":
		s.addDiagnostic("FINDING_NORMALIZATION_FAILED", "repository", nil)
	case "component":
		s.addDiagnostic("COMPONENT_IDENTITY_INVALID", "repository", nil)
	case "artifact":
		s.addDiagnostic("ARTIFACT_STORE_FAILED", "repository", nil)
	default:
		s.unclassified = true
	}
}

func (s *scanCaptureState) addWorkspaceError(err error, kind, workspaceID string) {
	if s == nil || err == nil {
		return
	}
	if code, ok := safeDiagnosticCode(err); ok {
		s.addDiagnostic(code, "workspace", &workspaceID)
		return
	}
	s.addBoundaryError(err, kind)
}

func stableHistoryError(err error) error {
	if err == nil {
		return nil
	}
	return &HistoryError{Code: HistoryWriteFailed, cause: err}
}
