package history

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	domain "github.com/geoffrey-xiao/deprail/internal/domain/history"
	"github.com/geoffrey-xiao/deprail/internal/normalize"
	historyv1 "github.com/geoffrey-xiao/deprail/schemas/history-v1"
	"io"
	"strings"
)

type limits struct {
	entries int64
	bytes   int64
}

var productionLimits = limits{entries: maxHistoryEntries, bytes: maxHistoryBytes}

const (
	createEntries = `CREATE TABLE history_entries (
  history_entry_id TEXT PRIMARY KEY NOT NULL,
  recorded_at_us INTEGER NOT NULL,
  operation_outcome TEXT NOT NULL CHECK (operation_outcome IN ('completed', 'failed', 'cancelled')),
  source_scan_id TEXT,
  report_status TEXT CHECK (report_status IS NULL OR report_status IN ('complete', 'partial', 'failed')),
  repository_label TEXT,
  workspace_count INTEGER CHECK (workspace_count IS NULL OR workspace_count >= 0),
  finding_count INTEGER CHECK (finding_count IS NULL OR finding_count >= 0),
  projection_schema_version TEXT NOT NULL CHECK (projection_schema_version = 'history-v1'),
  projection_json TEXT NOT NULL,
  source_report_schema_version TEXT,
  CHECK (
    (source_scan_id IS NULL AND source_report_schema_version IS NULL AND report_status IS NULL AND finding_count IS NULL)
    OR
    (source_scan_id IS NOT NULL AND source_report_schema_version IS NOT NULL AND report_status IS NOT NULL AND finding_count IS NOT NULL)
  )
)`
	createEntriesOrder = `CREATE INDEX history_entries_order ON history_entries (recorded_at_us DESC, history_entry_id DESC)`
	createArtifactRefs = `CREATE TABLE history_artifact_refs (
  history_entry_id TEXT NOT NULL REFERENCES history_entries(history_entry_id) ON DELETE CASCADE,
  digest TEXT NOT NULL CHECK (length(digest) = 64 AND digest NOT GLOB '*[^0-9a-f]*'),
  PRIMARY KEY (history_entry_id, digest)
)`
	createArtifactDigestIndex = `CREATE INDEX history_artifact_refs_digest ON history_artifact_refs (digest)`
)

var expectedSchemaObjects = map[string]string{
	"table:history_entries":              createEntries,
	"index:history_entries_order":        createEntriesOrder,
	"table:history_artifact_refs":        createArtifactRefs,
	"index:history_artifact_refs_digest": createArtifactDigestIndex,
}

func initializeOrValidate(ctx context.Context, db *sql.DB, writable bool) error {
	if err := verifyConnectionPragmas(ctx, db, writable); err != nil {
		return storeError(ErrUnavailable, err)
	}
	var version int
	if err := db.QueryRowContext(ctx, "PRAGMA user_version").Scan(&version); err != nil {
		return storeError(ErrCorrupt, err)
	}
	switch {
	case version > 1:
		return storeError(ErrSchemaUnsupported, nil)
	case version < 0:
		return storeError(ErrCorrupt, nil)
	case version == 1:
		if err := verifyJournalMode(ctx, db); err != nil {
			return storeError(ErrCorrupt, err)
		}
		if err := verifySchema(ctx, db); err != nil {
			return storeError(ErrCorrupt, err)
		}
		return nil
	default:
		if err := verifyEmptyDatabase(ctx, db); err != nil {
			return storeError(ErrCorrupt, err)
		}
		if !writable {
			return storeError(ErrCorrupt, errors.New("uninitialized database is not readable history"))
		}
		if err := enableWAL(ctx, db); err != nil {
			return storeError(ErrWriteFailed, err)
		}
		if err := initializeSchema(ctx, db); err != nil {
			return storeError(ErrWriteFailed, err)
		}
		if err := verifyJournalMode(ctx, db); err != nil {
			return storeError(ErrCorrupt, err)
		}
		if err := verifySchema(ctx, db); err != nil {
			return storeError(ErrCorrupt, err)
		}
		return nil
	}
}

func verifyConnectionPragmas(ctx context.Context, db *sql.DB, writable bool) error {
	var foreignKeys int
	if err := db.QueryRowContext(ctx, "PRAGMA foreign_keys").Scan(&foreignKeys); err != nil {
		return err
	}
	if foreignKeys != 1 {
		return errors.New("foreign keys are disabled")
	}
	if writable {
		var synchronous int
		if err := db.QueryRowContext(ctx, "PRAGMA synchronous").Scan(&synchronous); err != nil {
			return err
		}
		if synchronous != 2 {
			return errors.New("full synchronous mode is disabled")
		}
	}
	return nil
}

func verifyJournalMode(ctx context.Context, db *sql.DB) error {
	var mode string
	if err := db.QueryRowContext(ctx, "PRAGMA journal_mode").Scan(&mode); err != nil {
		return err
	}
	if !strings.EqualFold(mode, "wal") {
		return errors.New("database is not in WAL mode")
	}
	return nil
}

func enableWAL(ctx context.Context, db *sql.DB) error {
	var mode string
	if err := db.QueryRowContext(ctx, "PRAGMA journal_mode=WAL").Scan(&mode); err != nil {
		return err
	}
	if !strings.EqualFold(mode, "wal") {
		return errors.New("SQLite refused WAL mode")
	}
	return nil
}

func verifyEmptyDatabase(ctx context.Context, db *sql.DB) error {
	var objects, pages, freePages int
	if err := db.QueryRowContext(ctx, `SELECT count(*) FROM sqlite_master WHERE name NOT LIKE 'sqlite_%'`).Scan(&objects); err != nil {
		return err
	}
	if err := db.QueryRowContext(ctx, "PRAGMA page_count").Scan(&pages); err != nil {
		return err
	}
	if err := db.QueryRowContext(ctx, "PRAGMA freelist_count").Scan(&freePages); err != nil {
		return err
	}
	if objects != 0 || freePages != 0 || pages > 1 {
		return errors.New("unversioned database contains objects or data")
	}
	return nil
}

func initializeSchema(ctx context.Context, db *sql.DB) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, statement := range []string{createEntries, createEntriesOrder, createArtifactRefs, createArtifactDigestIndex} {
		if _, err := tx.ExecContext(ctx, statement); err != nil {
			return err
		}
	}
	if _, err := tx.ExecContext(ctx, "PRAGMA user_version=1"); err != nil {
		return err
	}
	return tx.Commit()
}

func verifySchema(ctx context.Context, db *sql.DB) error {
	rows, err := db.QueryContext(ctx, `SELECT type, name, sql FROM sqlite_master WHERE name NOT LIKE 'sqlite_%' ORDER BY type, name`)
	if err != nil {
		return err
	}
	defer rows.Close()
	seen := make(map[string]bool, len(expectedSchemaObjects))
	for rows.Next() {
		var kind, name, statement string
		if err := rows.Scan(&kind, &name, &statement); err != nil {
			return err
		}
		key := kind + ":" + name
		expected, ok := expectedSchemaObjects[key]
		if !ok || normalizeSQL(statement) != normalizeSQL(expected) {
			return errors.New("database schema does not match history-v1")
		}
		seen[key] = true
	}
	if err := rows.Err(); err != nil {
		return err
	}
	if len(seen) != len(expectedSchemaObjects) {
		return errors.New("database schema is incomplete")
	}
	return nil
}

func normalizeSQL(statement string) string {
	return strings.Join(strings.Fields(strings.ToLower(statement)), " ")
}

func admit(ctx context.Context, tx *sql.Tx, bytes int64, quota limits) error {
	var entries, storedBytes int64
	if err := tx.QueryRowContext(ctx, `SELECT count(*), COALESCE(sum(length(CAST(projection_json AS BLOB))), 0) FROM history_entries`).Scan(&entries, &storedBytes); err != nil {
		return storeError(ErrCorrupt, err)
	}
	if entries >= quota.entries || bytes > quota.bytes-storedBytes {
		return storeError(ErrWriteFailed, errors.New("history admission limit exceeded"))
	}
	return nil
}

func insertEntry(ctx context.Context, tx *sql.Tx, entry Entry, projectionJSON string) error {
	p := entry.Projection
	var sourceScanID, reportStatus, sourceSchema any
	var findingCount any
	if p.Report != nil {
		sourceScanID = p.Report.SourceScanID
		reportStatus = p.Report.Status
		sourceSchema = p.Report.SourceSchemaVersion
		findingCount = len(p.Report.Findings)
	}
	var repositoryLabel any
	if p.RepositoryLabel != nil {
		repositoryLabel = *p.RepositoryLabel
	}
	var workspaceCount any
	if p.Workspaces != nil {
		workspaceCount = len(p.Workspaces)
	}
	_, err := tx.ExecContext(ctx, `INSERT INTO history_entries (
 history_entry_id, recorded_at_us, operation_outcome, source_scan_id, report_status,
 repository_label, workspace_count, finding_count, projection_schema_version, projection_json,
 source_report_schema_version
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		p.HistoryEntryID, p.RecordedAtUS, p.OperationOutcome, sourceScanID, reportStatus,
		repositoryLabel, workspaceCount, findingCount, p.SchemaVersion, projectionJSON, sourceSchema)
	if err != nil {
		return err
	}
	for _, digest := range entry.ArtifactRefs {
		if _, err := tx.ExecContext(ctx, `INSERT INTO history_artifact_refs (history_entry_id, digest) VALUES (?, ?)`, p.HistoryEntryID, digest); err != nil {
			return err
		}
	}
	return nil
}

type rowQuerier interface {
	QueryRowContext(context.Context, string, ...any) *sql.Row
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
}

func readEntryByID(ctx context.Context, query rowQuerier, id string) (Entry, []string, error) {
	var entry Entry
	var storedID, outcome, projectionVersion string
	var recordedAt int64
	var sourceScanID, reportStatus, repositoryLabel, sourceSchema sql.NullString
	var workspaceCount, findingCount sql.NullInt64
	var raw []byte
	err := query.QueryRowContext(ctx, `SELECT history_entry_id, recorded_at_us, operation_outcome,
 source_scan_id, report_status, repository_label, workspace_count, finding_count,
 projection_schema_version, projection_json, source_report_schema_version
 FROM history_entries WHERE history_entry_id=?`, id).Scan(
		&storedID, &recordedAt, &outcome, &sourceScanID, &reportStatus, &repositoryLabel,
		&workspaceCount, &findingCount, &projectionVersion, &raw, &sourceSchema)
	if err != nil {
		return Entry{}, nil, err
	}
	if len(raw) == 0 || projectionVersion != domain.SchemaVersion {
		return Entry{}, nil, storeError(ErrCorrupt, errors.New("stored row has an invalid projection version"))
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&entry.Projection); err != nil {
		return Entry{}, nil, storeError(ErrCorrupt, err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return Entry{}, nil, storeError(ErrCorrupt, errors.New("stored projection has trailing JSON data"))
	}
	if err := normalize.ValidateProjection(entry.Projection); err != nil {
		var projectionError *normalize.ProjectionError
		if errors.As(err, &projectionError) && projectionError.Code == "HISTORY_SCHEMA_UNSUPPORTED" {
			return Entry{}, nil, storeError(ErrSchemaUnsupported, err)
		}
		return Entry{}, nil, storeError(ErrCorrupt, err)
	}
	if err := historyv1.ValidateJSON(raw); err != nil {
		return Entry{}, nil, storeError(ErrCorrupt, err)
	}
	canonical, err := json.Marshal(entry.Projection)
	if err != nil || !bytes.Equal(raw, canonical) {
		return Entry{}, nil, storeError(ErrCorrupt, errors.New("stored projection is not canonical JSON"))
	}
	p := entry.Projection
	if storedID != p.HistoryEntryID || recordedAt != p.RecordedAtUS || outcome != p.OperationOutcome || projectionVersion != p.SchemaVersion ||
		!nullableStringEqual(repositoryLabel, p.RepositoryLabel) || !nullableCountEqual(workspaceCount, p.Workspaces != nil, len(p.Workspaces)) {
		return Entry{}, nil, storeError(ErrCorrupt, errors.New("row identity or summary does not match projection"))
	}
	if p.Report == nil {
		if sourceScanID.Valid || reportStatus.Valid || sourceSchema.Valid || findingCount.Valid {
			return Entry{}, nil, storeError(ErrCorrupt, errors.New("null report has non-null row metadata"))
		}
	} else if !sourceScanID.Valid || !reportStatus.Valid || !sourceSchema.Valid || !findingCount.Valid ||
		sourceScanID.String != p.Report.SourceScanID || reportStatus.String != p.Report.Status ||
		sourceSchema.String != p.Report.SourceSchemaVersion || findingCount.Int64 != int64(len(p.Report.Findings)) {
		return Entry{}, nil, storeError(ErrCorrupt, errors.New("report row metadata does not match projection"))
	}
	refs, err := readRefs(ctx, query, storedID)
	if err != nil {
		return Entry{}, nil, err
	}
	wantRefs := []string{}
	if p.Report != nil {
		wantRefs = p.Report.ArtifactDigests
	}
	if !equalStrings(refs, wantRefs) {
		return Entry{}, nil, storeError(ErrCorrupt, errors.New("artifact reference rows do not match projection"))
	}
	entry.ArtifactRefs = refs
	return entry, refs, nil
}

func readRefs(ctx context.Context, query rowQuerier, id string) ([]string, error) {
	rows, err := query.QueryContext(ctx, `SELECT digest FROM history_artifact_refs WHERE history_entry_id=? ORDER BY digest`, id)
	if err != nil {
		return nil, storeError(ErrCorrupt, err)
	}
	defer rows.Close()
	refs := []string{}
	for rows.Next() {
		var digest string
		if err := rows.Scan(&digest); err != nil {
			return nil, storeError(ErrCorrupt, err)
		}
		refs = append(refs, digest)
	}
	if err := rows.Err(); err != nil {
		return nil, storeError(ErrCorrupt, err)
	}
	return refs, nil
}

func nullableStringEqual(stored sql.NullString, value *string) bool {
	return stored.Valid == (value != nil) && (!stored.Valid || stored.String == *value)
}

func nullableCountEqual(stored sql.NullInt64, present bool, count int) bool {
	return stored.Valid == present && (!present || stored.Int64 == int64(count))
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
