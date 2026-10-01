package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/geoffrey-xiao/deprail/internal/adapter"
	"github.com/geoffrey-xiao/deprail/internal/artifact"
	"github.com/geoffrey-xiao/deprail/internal/discovery"
	domain "github.com/geoffrey-xiao/deprail/internal/domain/history"
	"github.com/geoffrey-xiao/deprail/internal/normalize"
	"github.com/geoffrey-xiao/deprail/internal/store/history"
)

type countingScanner struct {
	fakeScanner
	plans    int
	executes int
}

func (s *countingScanner) Plan(ctx context.Context, targets []adapter.Target) (adapter.Plan, error) {
	s.plans++
	return s.fakeScanner.Plan(ctx, targets)
}

func (s *countingScanner) Execute(ctx context.Context, plan adapter.Plan) (adapter.RawResult, error) {
	s.executes++
	return s.fakeScanner.Execute(ctx, plan)
}

type scanEventCounter struct {
	discoveryStarts int
}

func (c *scanEventCounter) Emit(event Event) error {
	if event.Type == EventWorkspaceDiscoveryStarted {
		c.discoveryStarts++
	}
	return nil
}

func TestScanCapturePersistsSameOperationAndQueriesAfterReopen(t *testing.T) {
	ctx := context.Background()
	root := fixturePath(t, "npm-basic")
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	configHome, err := os.MkdirTemp(home, "deprail-history-app-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(configHome) })
	t.Setenv("HOME", configHome)
	t.Setenv("XDG_CONFIG_HOME", configHome)
	t.Setenv("APPDATA", configHome)
	artifacts := artifact.Store{Root: filepath.Join(t.TempDir(), "artifacts"), MaxBytes: 1024}
	writer, err := history.OpenWriter(ctx, history.Options{RepositoryRoot: root})
	if err != nil {
		t.Fatalf("open history writer: %v", err)
	}
	service := &HistoryService{Writer: writer, Reader: writer, Artifacts: artifacts}
	scanner := &countingScanner{}
	events := &scanEventCounter{}
	capturedReport, err := Scan(ctx, root, ScanOptions{Scanner: scanner, Artifacts: artifacts, Events: events, Capture: service})
	if err != nil {
		t.Fatalf("scan with capture: %v", err)
	}
	if scanner.plans != 1 || scanner.executes != 1 {
		t.Fatalf("capture started extra scanner work: plans=%d executes=%d", scanner.plans, scanner.executes)
	}
	if events.discoveryStarts != 1 {
		t.Fatalf("capture repeated workspace discovery: starts=%d", events.discoveryStarts)
	}
	if len(capturedReport.Findings) != 1 || len(capturedReport.ArtifactDigests) != 1 {
		t.Fatalf("captured scan report = %#v", capturedReport)
	}
	capturedBytes, err := json.Marshal(capturedReport)
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{
		"00000000-0000-4000-8000-000000000001",
		"00000000-0000-4000-8000-000000000002",
		"00000000-0000-4000-8000-000000000003",
	} {
		if _, err := service.Capture(ctx, CaptureInput{
			HistoryEntryID: id, RecordedAtUS: 253402300000000000,
			Diagnostics: []normalize.HistoryDiagnosticInput{},
		}); err != nil {
			t.Fatalf("capture pagination fixture: %v", err)
		}
	}
	_, err = service.Capture(ctx, CaptureInput{
		HistoryEntryID: "00000000-0000-4000-8000-000000000004",
		RecordedAtUS:   253402300000000000,
		Diagnostics:    []normalize.HistoryDiagnosticInput{{Code: "UNKNOWN_PRIVATE_CODE", Scope: "repository"}},
	})
	if err == nil {
		t.Fatal("unknown diagnostic was persisted")
	}
	tooManyDiagnostics := make([]normalize.HistoryDiagnosticInput, 129)
	historyWorkspaces := make([]discovery.Workspace, len(tooManyDiagnostics))
	for i := range tooManyDiagnostics {
		id := fmt.Sprintf("ws-%03d", i)
		historyWorkspaces[i] = discovery.Workspace{
			WorkspaceID: id, RelativePath: "workspace-" + id, Ecosystem: discovery.EcosystemNPM,
			PackageManager: "npm", Completeness: discovery.Complete,
		}
		tooManyDiagnostics[i] = normalize.HistoryDiagnosticInput{
			Code: "CANCELLED", Scope: "workspace", WorkspaceID: &historyWorkspaces[i].WorkspaceID,
		}
	}
	_, err = service.Capture(ctx, CaptureInput{
		HistoryEntryID: "00000000-0000-4000-8000-000000000005",
		RecordedAtUS:   253402300000000000,
		Graph:          &discovery.ProjectGraph{Workspaces: historyWorkspaces},
		Diagnostics:    tooManyDiagnostics,
	})
	if err == nil {
		t.Fatal("oversized diagnostic projection was persisted")
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	reader, err := history.OpenReadOnly(ctx, history.Options{RepositoryRoot: root})
	if err != nil {
		t.Fatalf("reopen history read-only: %v", err)
	}
	defer reader.Close()
	service = &HistoryService{Reader: reader, Artifacts: artifacts}
	var cursor *domain.Cursor
	var orderedIDs []string
	var capturedEntryID string
	for {
		page, err := service.List(ctx, HistoryPageRequest{Limit: 1, Before: cursor, MaxBytes: 1 << 20})
		if err != nil {
			t.Fatalf("list history: %v", err)
		}
		for _, item := range page.Items {
			orderedIDs = append(orderedIDs, item.HistoryEntryID)
			if item.SourceScanID != nil && *item.SourceScanID == capturedReport.ScanID {
				capturedEntryID = item.HistoryEntryID
			}
		}
		if !page.HasMore {
			break
		}
		if page.Next == nil {
			t.Fatal("nonterminal history page omitted cursor")
		}
		cursor = page.Next
	}
	if len(orderedIDs) != 4 || orderedIDs[0] != "00000000-0000-4000-8000-000000000003" || orderedIDs[1] != "00000000-0000-4000-8000-000000000002" || orderedIDs[2] != "00000000-0000-4000-8000-000000000001" {
		t.Fatalf("equal-timestamp history order or pagination = %v", orderedIDs)
	}
	if capturedEntryID == "" {
		t.Fatal("same-operation scan was absent after reopen")
	}
	if _, err := service.List(ctx, HistoryPageRequest{Limit: 1, MaxBytes: 600}); err == nil {
		t.Fatal("expected a single oversized summary to fail")
	} else {
		var historyErr *HistoryError
		if !errors.As(err, &historyErr) || historyErr.Code != HistoryResponseTooLarge {
			t.Fatalf("oversized page error = %v", err)
		}
	}
	detail, err := service.Detail(ctx, capturedEntryID)
	if err != nil {
		t.Fatalf("history detail: %v", err)
	}
	if detail.Report == nil || detail.Report.RepositoryState != capturedReport.RepositoryState || len(detail.ArtifactReferences) != 1 || detail.ArtifactReferences[0].Integrity != "verified" {
		t.Fatalf("same-operation detail = %#v", detail)
	}
	detailBytes, err := json.Marshal(detail)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(detailBytes), root) || strings.Contains(string(detailBytes), `{"results":[]`) {
		t.Fatalf("history detail exposed repository path or raw artifact: %s", detailBytes)
	}
	unavailableService := &HistoryService{Reader: reader}
	unavailable, err := unavailableService.ArtifactIntegrity(ctx, capturedEntryID)
	if err != nil || len(unavailable) != 1 || unavailable[0].Integrity != "unavailable" {
		t.Fatalf("unresolved artifact integrity = %#v, err=%v", unavailable, err)
	}
	artifactPath := filepath.Join(artifacts.Root, capturedReport.ArtifactDigests[0][:2], capturedReport.ArtifactDigests[0])
	if err := os.Remove(artifactPath); err != nil {
		t.Fatal(err)
	}
	missing, err := service.ArtifactIntegrity(ctx, capturedEntryID)
	if err != nil || len(missing) != 1 || missing[0].Integrity != "missing" {
		t.Fatalf("missing artifact integrity = %#v, err=%v", missing, err)
	}
	if err := os.WriteFile(artifactPath, []byte("tampered"), 0o600); err != nil {
		t.Fatal(err)
	}
	mismatch, err := service.ArtifactIntegrity(ctx, capturedEntryID)
	if err != nil || len(mismatch) != 1 || mismatch[0].Integrity != "digest_mismatch" {
		t.Fatalf("mismatched artifact integrity = %#v, err=%v", mismatch, err)
	}
	workspaces, err := service.Workspaces(ctx, WorkspacePageRequest{HistoryEntryID: capturedEntryID, Limit: 25, MaxBytes: 1 << 20})
	if err != nil || !workspaces.Available || len(workspaces.Items) == 0 {
		t.Fatalf("workspaces = %#v, err=%v", workspaces, err)
	}
	findings, err := service.Findings(ctx, FindingPageRequest{HistoryEntryID: capturedEntryID, Limit: 25, MaxBytes: 1 << 20})
	if err != nil || !findings.Available || len(findings.Items) != 1 {
		t.Fatalf("findings = %#v, err=%v", findings, err)
	}
	plainArtifacts := artifact.Store{Root: filepath.Join(t.TempDir(), "plain-artifacts"), MaxBytes: 1024}
	plainReport, err := Scan(ctx, root, ScanOptions{Scanner: fakeScanner{}, Artifacts: plainArtifacts})
	if err != nil {
		t.Fatalf("scan without capture: %v", err)
	}
	plainBytes, err := json.Marshal(plainReport)
	if err != nil {
		t.Fatal(err)
	}
	if string(capturedBytes) != string(plainBytes) {
		t.Fatalf("opt-in capture changed scan report bytes:\nwith capture: %s\nwithout capture: %s", capturedBytes, plainBytes)
	}
	unchanged, err := service.List(ctx, HistoryPageRequest{Limit: 25, MaxBytes: 1 << 20})
	if err != nil || len(unchanged.Items) != 4 {
		t.Fatalf("non-opt-in scan changed history: entries=%d, err=%v", len(unchanged.Items), err)
	}
}

type captureRecorder struct {
	calls      int
	input      CaptureInput
	contextErr error
	err        error
}

func (r *captureRecorder) Capture(ctx context.Context, input CaptureInput) (SaveResult, error) {
	r.calls++
	r.input = input
	r.contextErr = ctx.Err()
	if r.err != nil {
		return SaveResult{}, r.err
	}
	return SaveResult{HistoryEntryID: input.HistoryEntryID, Saved: true}, nil
}

func TestScanCaptureRunsOnceAtFailedTerminalAndKeepsPersistenceFailureSeparate(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	recorder := &captureRecorder{err: errors.New("private store failure")}
	report, err := Scan(ctx, filepath.Join(t.TempDir(), "missing"), ScanOptions{Capture: recorder})
	if err == nil {
		t.Fatal("expected scan failure")
	}
	var captureErr *ScanCaptureError
	if !errors.As(err, &captureErr) || captureErr.OperationErr != recorder.input.OperationError || captureErr.PersistenceErr == nil {
		t.Fatalf("capture failure did not preserve both errors: %T %v", err, err)
	}
	if recorder.calls != 1 || recorder.contextErr != nil {
		t.Fatalf("capture calls/context = %d/%v", recorder.calls, recorder.contextErr)
	}
	if recorder.input.Graph != nil || recorder.input.Report != nil || recorder.input.OperationError == nil || report.Status != "failed" {
		t.Fatalf("failed terminal context = %#v, report=%#v", recorder.input, report)
	}
	if len(recorder.input.Diagnostics) != 1 || recorder.input.Diagnostics[0].Code != "DISCOVERY_FAILED" {
		t.Fatalf("failed terminal diagnostics = %#v", recorder.input.Diagnostics)
	}
	if !errors.Is(err, recorder.input.OperationError) || !errors.Is(err, recorder.err) {
		t.Fatalf("combined failure lost a cause: %v", err)
	}
	if got := captureErr.PersistenceErr.(*HistoryError).Code; got != HistoryWriteFailed {
		t.Fatalf("persistence error code = %s", got)
	}
}
