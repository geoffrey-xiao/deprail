package localhttp

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/geoffrey-xiao/deprail/internal/app"
	domain "github.com/geoffrey-xiao/deprail/internal/domain/history"
)

const (
	cursorVersion    = 1
	cursorHistory    = 1
	cursorWorkspaces = 2
	cursorFindings   = 3
	cursorMaxBytes   = 512
	maxRecordedAtUS  = 253402300799999999
)

var errCursorInvalid = errors.New("invalid continuation")

func encodeHistoryCursor(cursor domain.Cursor) (string, error) {
	id, ok := parseHistoryID(cursor.HistoryEntryID)
	if !ok || cursor.RecordedAtUS < 0 || cursor.RecordedAtUS > maxRecordedAtUS {
		return "", errCursorInvalid
	}
	var raw [26]byte
	raw[0], raw[1] = cursorVersion, cursorHistory
	binary.BigEndian.PutUint64(raw[2:10], uint64(cursor.RecordedAtUS))
	copy(raw[10:], id[:])
	return base64.RawURLEncoding.EncodeToString(raw[:]), nil
}

func decodeHistoryCursor(value string) (*domain.Cursor, error) {
	raw, err := decodeCursor(value, 26)
	if err != nil || raw[0] != cursorVersion || raw[1] != cursorHistory {
		return nil, errCursorInvalid
	}
	timestamp := int64(binary.BigEndian.Uint64(raw[2:10]))
	id := formatHistoryID(raw[10:26])
	if timestamp < 0 || timestamp > maxRecordedAtUS {
		return nil, errCursorInvalid
	}
	return &domain.Cursor{RecordedAtUS: timestamp, HistoryEntryID: id}, nil
}

func encodeWorkspaceCursor(parentID string, cursor *app.WorkspaceCursor) (string, error) {
	parent, ok := parseHistoryID(parentID)
	if !ok || cursor == nil || !cursor.Encoded {
		return "", errCursorInvalid
	}
	var raw [54]byte
	raw[0], raw[1] = cursorVersion, cursorWorkspaces
	copy(raw[2:18], parent[:])
	binary.BigEndian.PutUint32(raw[18:22], cursor.Ordinal)
	copy(raw[22:], cursor.WorkspaceHash[:])
	return base64.RawURLEncoding.EncodeToString(raw[:]), nil
}

func decodeWorkspaceCursor(value, parentID string) (*app.WorkspaceCursor, error) {
	parent, ok := parseHistoryID(parentID)
	if !ok {
		return nil, errCursorInvalid
	}
	raw, err := decodeCursor(value, 54)
	if err != nil || raw[0] != cursorVersion || raw[1] != cursorWorkspaces || !equalBytes(raw[2:18], parent[:]) {
		return nil, errCursorInvalid
	}
	cursor := &app.WorkspaceCursor{Ordinal: binary.BigEndian.Uint32(raw[18:22]), Encoded: true}
	copy(cursor.WorkspaceHash[:], raw[22:54])
	return cursor, nil
}

func encodeFindingCursor(parentID string, cursor *app.FindingCursor) (string, error) {
	parent, ok := parseHistoryID(parentID)
	if !ok || cursor == nil {
		return "", errCursorInvalid
	}
	var key [32]byte
	if !decodeLowerHex(key[:], cursor.StableFindingKey) {
		return "", errCursorInvalid
	}
	var raw [50]byte
	raw[0], raw[1] = cursorVersion, cursorFindings
	copy(raw[2:18], parent[:])
	copy(raw[18:], key[:])
	return base64.RawURLEncoding.EncodeToString(raw[:]), nil
}

func decodeFindingCursor(value, parentID string) (*app.FindingCursor, error) {
	parent, ok := parseHistoryID(parentID)
	if !ok {
		return nil, errCursorInvalid
	}
	raw, err := decodeCursor(value, 50)
	if err != nil || raw[0] != cursorVersion || raw[1] != cursorFindings || !equalBytes(raw[2:18], parent[:]) {
		return nil, errCursorInvalid
	}
	return &app.FindingCursor{StableFindingKey: hex.EncodeToString(raw[18:50])}, nil
}

func decodeCursor(value string, decodedLength int) ([]byte, error) {
	if value == "" || len(value) > cursorMaxBytes || strings.Contains(value, "=") {
		return nil, errCursorInvalid
	}
	raw, err := base64.RawURLEncoding.Strict().DecodeString(value)
	if err != nil || len(raw) != decodedLength {
		return nil, errCursorInvalid
	}
	return raw, nil
}

func parseHistoryID(value string) ([16]byte, bool) {
	var id [16]byte
	if len(value) != 36 || value[8] != '-' || value[13] != '-' || value[18] != '-' || value[23] != '-' {
		return id, false
	}
	byteIndex := 0
	for i := 0; i < len(value); {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			i++
			continue
		}
		high, highOK := lowerHexNibble(value[i])
		low, lowOK := lowerHexNibble(value[i+1])
		if !highOK || !lowOK || byteIndex >= len(id) {
			return [16]byte{}, false
		}
		id[byteIndex] = high<<4 | low
		byteIndex++
		i += 2
	}
	if byteIndex != len(id) || id[6]>>4 != 4 || id[8]&0xc0 != 0x80 {
		return [16]byte{}, false
	}
	return id, true
}

func decodeLowerHex(dst []byte, value string) bool {
	if len(value) != len(dst)*2 {
		return false
	}
	for i := range dst {
		high, highOK := lowerHexNibble(value[i*2])
		low, lowOK := lowerHexNibble(value[i*2+1])
		if !highOK || !lowOK {
			return false
		}
		dst[i] = high<<4 | low
	}
	return true
}

func lowerHexNibble(value byte) (byte, bool) {
	switch {
	case value >= '0' && value <= '9':
		return value - '0', true
	case value >= 'a' && value <= 'f':
		return value - 'a' + 10, true
	default:
		return 0, false
	}
}

func formatHistoryID(raw []byte) string {
	var text [36]byte
	hex.Encode(text[0:8], raw[0:4])
	text[8] = '-'
	hex.Encode(text[9:13], raw[4:6])
	text[13] = '-'
	hex.Encode(text[14:18], raw[6:8])
	text[18] = '-'
	hex.Encode(text[19:23], raw[8:10])
	text[23] = '-'
	hex.Encode(text[24:36], raw[10:16])
	return string(text[:])
}

func equalBytes(left, right []byte) bool {
	if len(left) != len(right) {
		return false
	}
	var diff byte
	for i := range left {
		diff |= left[i] ^ right[i]
	}
	return diff == 0
}
