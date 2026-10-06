package app

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	domain "github.com/geoffrey-xiao/deprail/internal/domain/history"
	"github.com/geoffrey-xiao/deprail/internal/normalize"
)

type fixedHistoryReader struct{ entry domain.Entry }

func (r fixedHistoryReader) Get(context.Context, string) (domain.Entry, error) {
	return r.entry, nil
}

func (fixedHistoryReader) List(context.Context, domain.Page) (domain.PageResult, error) {
	return domain.PageResult{Entries: []domain.Summary{}}, nil
}

func TestWorkspacePageUsesUTF8BytesForResponseBudget(t *testing.T) {
	const entryID = "00000000-0000-4000-8000-000000000101"
	path := strings.Repeat("界", 10)
	projection, err := normalize.ProjectHistory(normalize.HistoryInput{
		HistoryEntryID: entryID, RecordedAtUS: 1, OperationOutcome: "completed",
		Workspaces: []domain.Workspace{{
			WorkspaceID: "workspace-1", Path: path, Ecosystem: "npm",
			PackageManager: "npm", DiscoveryCompleteness: "complete",
		}},
		Diagnostics: []normalize.HistoryDiagnosticInput{},
	})
	if err != nil {
		t.Fatalf("project query fixture: %v", err)
	}
	service := &HistoryService{Reader: fixedHistoryReader{entry: domain.Entry{Projection: projection, ArtifactRefs: []string{}}}}
	itemBytes, err := json.Marshal(WorkspaceView{
		WorkspaceID: "workspace-1", RelativePath: path, Ecosystem: "npm",
		PackageManager: "npm", DiscoveryCompleteness: "complete",
	})
	if err != nil {
		t.Fatal(err)
	}
	emptyEnvelope := `{"collectionState":"available","items":[],"nextCursor":null}`
	reserved := len(emptyEnvelope) - 2 - 4 + maxEncodedCursorSize + 2
	// A rune-counting implementation would fit this candidate; the public query
	// must count the actual UTF-8 response bytes and reject the whole record.
	maxBytes := reserved + len(itemBytes) - len([]rune(path))
	if maxBytes < 600 {
		t.Fatalf("test budget %d is below the minimum request limit", maxBytes)
	}
	_, err = service.Workspaces(context.Background(), WorkspacePageRequest{
		HistoryEntryID: entryID, Limit: 25, MaxBytes: maxBytes,
	})
	var historyErr *HistoryError
	if !errors.As(err, &historyErr) || historyErr.Code != HistoryResponseTooLarge {
		t.Fatalf("workspace page did not enforce UTF-8 byte limit: err=%v budget=%d candidate=%d", err, maxBytes, len(itemBytes))
	}
}

func TestHistoryQueryRejectsMalformedContinuations(t *testing.T) {
	service := &HistoryService{}
	_, err := service.List(context.Background(), HistoryPageRequest{
		Before: &domain.Cursor{RecordedAtUS: 0, HistoryEntryID: "not-a-uuid"},
	})
	requireHistoryErrorCode(t, err, HistoryRequestInvalid)
	_, err = service.Workspaces(context.Background(), WorkspacePageRequest{
		HistoryEntryID: "00000000-0000-4000-8000-000000000101",
		Before:         &WorkspaceCursor{WorkspaceID: "../outside"},
	})
	requireHistoryErrorCode(t, err, HistoryRequestInvalid)
	_, err = service.Findings(context.Background(), FindingPageRequest{
		HistoryEntryID: "00000000-0000-4000-8000-000000000101",
		Before:         &FindingCursor{StableFindingKey: "not-a-digest"},
	})
	requireHistoryErrorCode(t, err, HistoryRequestInvalid)
}

func requireHistoryErrorCode(t *testing.T, err error, code HistoryErrorCode) {
	t.Helper()
	var historyErr *HistoryError
	if !errors.As(err, &historyErr) || historyErr.Code != code {
		t.Fatalf("history error = %v, want %s", err, code)
	}
}

func TestWorkspaceHashedOrdinalCursorContinuesAndRejectsTampering(t *testing.T) {
	const entryID = "00000000-0000-4000-8000-000000000101"
	projection, err := normalize.ProjectHistory(normalize.HistoryInput{
		HistoryEntryID: entryID, RecordedAtUS: 1, OperationOutcome: "completed",
		Workspaces: []domain.Workspace{
			{WorkspaceID: "workspace-b", Path: "b", Ecosystem: "npm", PackageManager: "npm", DiscoveryCompleteness: "complete"},
			{WorkspaceID: "workspace-a", Path: "a", Ecosystem: "npm", PackageManager: "npm", DiscoveryCompleteness: "complete"},
		},
		Diagnostics: []normalize.HistoryDiagnosticInput{},
	})
	if err != nil {
		t.Fatalf("project query fixture: %v", err)
	}
	service := &HistoryService{Reader: fixedHistoryReader{entry: domain.Entry{Projection: projection, ArtifactRefs: []string{}}}}
	first, err := service.Workspaces(context.Background(), WorkspacePageRequest{HistoryEntryID: entryID, Limit: 1, MaxBytes: 1 << 20})
	if err != nil || !first.Available || len(first.Items) != 1 || first.Items[0].WorkspaceID != "workspace-a" || first.Next == nil {
		t.Fatalf("first workspace page = %#v, err=%v", first, err)
	}
	cursor := first.Next
	if !cursor.Encoded || cursor.Ordinal != 0 || cursor.WorkspaceID != "" {
		t.Fatalf("application cursor payload = %#v", cursor)
	}
	second, err := service.Workspaces(context.Background(), WorkspacePageRequest{HistoryEntryID: entryID, Limit: 1, Before: cursor, MaxBytes: 1 << 20})
	if err != nil || len(second.Items) != 1 || second.Items[0].WorkspaceID != "workspace-b" || second.HasMore {
		t.Fatalf("continued workspace page = %#v, err=%v", second, err)
	}

	withIdentity := *cursor
	withIdentity.WorkspaceID = "workspace-a"
	_, err = service.Workspaces(context.Background(), WorkspacePageRequest{HistoryEntryID: entryID, Limit: 1, Before: &withIdentity, MaxBytes: 1 << 20})
	requireHistoryErrorCode(t, err, HistoryRequestInvalid)
	tampered := *cursor
	tampered.WorkspaceHash[0] ^= 1
	_, err = service.Workspaces(context.Background(), WorkspacePageRequest{HistoryEntryID: entryID, Limit: 1, Before: &tampered, MaxBytes: 1 << 20})
	requireHistoryErrorCode(t, err, HistoryRequestInvalid)
}
