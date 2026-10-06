package localhttp

import (
	"encoding/base64"
	"strings"
	"testing"

	"github.com/geoffrey-xiao/deprail/internal/app"
	domain "github.com/geoffrey-xiao/deprail/internal/domain/history"
)

func TestCursorCodecsPreserveOrderAndBindResourceParent(t *testing.T) {
	const parent = "00000000-0000-4000-8000-000000000101"
	const otherParent = "00000000-0000-4000-8000-000000000102"

	history := domain.Cursor{RecordedAtUS: 1780000000123456, HistoryEntryID: parent}
	historyEncoded, err := encodeHistoryCursor(history)
	if err != nil || len(historyEncoded) != 35 {
		t.Fatalf("encode history cursor = %q, err=%v", historyEncoded, err)
	}
	historyDecoded, err := decodeHistoryCursor(historyEncoded)
	if err != nil || *historyDecoded != history {
		t.Fatalf("decode history cursor = %#v, err=%v", historyDecoded, err)
	}
	if _, err := decodeFindingCursor(historyEncoded, parent); err == nil {
		t.Fatal("history cursor was accepted for findings")
	}
	if _, err := decodeHistoryCursor(historyEncoded + "="); err == nil {
		t.Fatal("padded cursor was accepted")
	}

	workspace := &app.WorkspaceCursor{Ordinal: 17, Encoded: true}
	copy(workspace.WorkspaceHash[:], strings.Repeat("w", len(workspace.WorkspaceHash)))
	workspaceEncoded, err := encodeWorkspaceCursor(parent, workspace)
	if err != nil || len(workspaceEncoded) != 72 {
		t.Fatalf("encode workspace cursor = %q, err=%v", workspaceEncoded, err)
	}
	workspaceDecoded, err := decodeWorkspaceCursor(workspaceEncoded, parent)
	if err != nil || workspaceDecoded.Ordinal != workspace.Ordinal || workspaceDecoded.WorkspaceHash != workspace.WorkspaceHash || !workspaceDecoded.Encoded {
		t.Fatalf("decode workspace cursor = %#v, err=%v", workspaceDecoded, err)
	}
	if _, err := decodeWorkspaceCursor(workspaceEncoded, otherParent); err == nil {
		t.Fatal("workspace cursor was reused for a different parent")
	}

	finding := &app.FindingCursor{StableFindingKey: strings.Repeat("a", 64)}
	findingEncoded, err := encodeFindingCursor(parent, finding)
	if err != nil || len(findingEncoded) != 67 {
		t.Fatalf("encode finding cursor = %q, err=%v", findingEncoded, err)
	}
	findingDecoded, err := decodeFindingCursor(findingEncoded, parent)
	if err != nil || findingDecoded.StableFindingKey != finding.StableFindingKey {
		t.Fatalf("decode finding cursor = %#v, err=%v", findingDecoded, err)
	}
	if _, err := decodeFindingCursor(findingEncoded, otherParent); err == nil {
		t.Fatal("finding cursor was reused for a different parent")
	}
}

func TestCursorCodecsRejectNoncanonicalAndAlteredValues(t *testing.T) {
	const parent = "0000000a-0000-4000-8000-000000000101"
	encoded, err := encodeHistoryCursor(domain.Cursor{RecordedAtUS: 1, HistoryEntryID: parent})
	if err != nil {
		t.Fatal(err)
	}
	raw, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatal(err)
	}
	for _, mutation := range []struct {
		index int
		value byte
	}{
		{index: 0, value: 2},
		{index: 1, value: cursorFindings},
	} {
		altered := append([]byte(nil), raw...)
		altered[mutation.index] = mutation.value
		if _, err := decodeHistoryCursor(base64.RawURLEncoding.EncodeToString(altered)); err == nil {
			t.Fatalf("accepted cursor mutation at byte %d", mutation.index)
		}
	}
	if _, err := encodeHistoryCursor(domain.Cursor{RecordedAtUS: -1, HistoryEntryID: parent}); err == nil {
		t.Fatal("encoded a negative timestamp")
	}
	if _, err := encodeHistoryCursor(domain.Cursor{RecordedAtUS: 1, HistoryEntryID: strings.ToUpper(parent)}); err == nil {
		t.Fatal("encoded a noncanonical UUID")
	}
}
