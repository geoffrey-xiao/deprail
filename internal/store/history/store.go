package history

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"

	domain "github.com/geoffrey-xiao/deprail/internal/domain/history"
	"github.com/geoffrey-xiao/deprail/internal/normalize"
	historyv1 "github.com/geoffrey-xiao/deprail/schemas/history-v1"
	_ "modernc.org/sqlite"
)

const (
	busyTimeout       = 2 * time.Second
	defaultPageSize   = 25
	maxPageSize       = 50
	maxHistoryEntries = 1000
	maxHistoryBytes   = 256 << 20
)

var historyIDPattern = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

// ErrorCode identifies a stable local-history storage failure.
type ErrorCode string

const (
	ErrUnavailable       ErrorCode = "HISTORY_UNAVAILABLE"
	ErrWriteFailed       ErrorCode = "HISTORY_WRITE_FAILED"
	ErrSchemaUnsupported ErrorCode = "HISTORY_SCHEMA_UNSUPPORTED"
	ErrMigrationFailed   ErrorCode = "HISTORY_MIGRATION_FAILED"
	ErrCorrupt           ErrorCode = "HISTORY_CORRUPT"
	ErrEntryNotFound     ErrorCode = "HISTORY_ENTRY_NOT_FOUND"
)

// Error carries a stable code without exposing database, path, or OS details.
type Error struct {
	Code  ErrorCode
	cause error
}

func (e *Error) HistoryErrorCode() string { return string(e.Code) }
func (e *Error) Error() string            { return string(e.Code) }
func (e *Error) Unwrap() error            { return e.cause }

func IsCode(err error, code ErrorCode) bool {
	var target *Error
	return errors.As(err, &target) && target.Code == code
}

func storeError(code ErrorCode, cause error) error { return &Error{Code: code, cause: cause} }

// Options supplies the scanned repository boundary when opening a writer or
// when a read must be checked against a currently selected repository.
type Options struct {
	RepositoryRoot string
}

// Entry, Cursor, Page, Summary, and PageResult remain storage-facing names for
// the dependency-neutral history models.
type Entry = domain.Entry
type Cursor = domain.Cursor
type Page = domain.Page
type Summary = domain.Summary
type PageResult = domain.PageResult

// Store owns one SQLite handle. Read-only stores never create missing paths.
type Store struct {
	db       *sql.DB
	path     string
	readOnly bool
	writers  chan struct{}
	limits   limits
}

// OpenReadOnly opens the per-user store without creating any missing directory
// or file. RepositoryRoot, when supplied, is checked as a containment boundary.
func OpenReadOnly(ctx context.Context, options Options) (*Store, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, storeError(ErrUnavailable, err)
	}
	return openAt(ctx, configDir, options.RepositoryRoot, false, productionLimits)
}

// OpenWriter is the explicit write boundary. It requires the active scanned
// repository root so the data directory can be rejected if it is inside it.
func OpenWriter(ctx context.Context, options Options) (*Store, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, storeError(ErrWriteFailed, err)
	}
	return openAt(ctx, configDir, options.RepositoryRoot, true, productionLimits)
}

func openAt(ctx context.Context, configDir, repositoryRoot string, writable bool, quota limits) (*Store, error) {
	code := ErrUnavailable
	if writable {
		code = ErrWriteFailed
	}
	if err := ctx.Err(); err != nil {
		return nil, storeError(code, err)
	}
	if writable && repositoryRoot == "" {
		return nil, storeError(ErrWriteFailed, errors.New("repository boundary is required"))
	}
	paths, missing, err := preparePaths(configDir, repositoryRoot, writable)
	if err != nil {
		return nil, storeError(code, err)
	}
	store := &Store{path: paths.database, readOnly: !writable, writers: make(chan struct{}, 1), limits: quota}
	if missing && !writable {
		return store, nil
	}
	if writable {
		if err := secureFile(paths.database, true, false); err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return nil, storeError(ErrWriteFailed, err)
			}
			if err := rejectOrphanSidecars(paths.database); err != nil {
				return nil, storeError(ErrCorrupt, err)
			}
			file, createErr := os.OpenFile(paths.database, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0o600)
			created := createErr == nil
			if createErr != nil && !errors.Is(createErr, os.ErrExist) {
				return nil, storeError(ErrWriteFailed, createErr)
			}
			if created {
				if closeErr := file.Close(); closeErr != nil {
					return nil, storeError(ErrWriteFailed, closeErr)
				}
			}
			if err := secureFile(paths.database, true, created); err != nil {
				return nil, storeError(ErrWriteFailed, err)
			}
		}
	} else if err := secureFile(paths.database, false, false); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return store, nil
		}
		return nil, storeError(ErrUnavailable, err)
	}

	if writable {
		if err := secureSQLiteSidecars(paths.database, false); err != nil {
			return nil, storeError(ErrWriteFailed, err)
		}
	}
	db, err := sql.Open("sqlite", dataSourceName(paths.database, writable))
	if err != nil {
		return nil, storeError(code, err)
	}
	if writable {
		db.SetMaxOpenConns(1)
		db.SetMaxIdleConns(1)
	} else {
		db.SetMaxOpenConns(8)
		db.SetMaxIdleConns(8)
	}
	store.db = db
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, storeError(code, err)
	}
	if err := initializeOrValidate(ctx, db, writable); err != nil {
		db.Close()
		return nil, err
	}
	if writable {
		if err := secureSQLiteSidecars(paths.database, true); err != nil {
			db.Close()
			return nil, storeError(ErrWriteFailed, err)
		}
		if err := secureFile(paths.database, true, false); err != nil {
			db.Close()
			return nil, storeError(ErrWriteFailed, err)
		}
	}
	return store, nil
}

// Close releases the connection pool. Closing an absent read-only store is safe.
func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	if err := s.db.Close(); err != nil {
		return storeError(ErrUnavailable, err)
	}
	return nil
}

// Append validates the exact projection and atomically inserts its row and refs.
func (s *Store) Append(ctx context.Context, entry Entry) error {
	if s == nil || s.readOnly || s.db == nil {
		return storeError(ErrWriteFailed, errors.New("writer is not open"))
	}
	if err := ctx.Err(); err != nil {
		return storeError(ErrWriteFailed, err)
	}
	if err := s.acquireWriter(ctx, ErrWriteFailed); err != nil {
		return err
	}
	defer func() { <-s.writers }()

	projectionBytes, err := encodeAndValidate(entry)
	if err != nil {
		return storeError(ErrWriteFailed, err)
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return storeError(ErrWriteFailed, err)
	}
	defer tx.Rollback()

	existing, refs, err := readEntryByID(ctx, tx, entry.Projection.HistoryEntryID)
	if err == nil {
		existingBytes, marshalErr := json.Marshal(existing.Projection)
		if marshalErr != nil {
			return storeError(ErrCorrupt, marshalErr)
		}
		if bytes.Equal(existingBytes, projectionBytes) && slices.Equal(refs, entry.ArtifactRefs) {
			if err := tx.Commit(); err != nil {
				return storeError(ErrWriteFailed, err)
			}
			if err := secureSQLiteSidecars(s.path, false); err != nil {
				return storeError(ErrWriteFailed, err)
			}
			return nil
		}
		return storeError(ErrWriteFailed, errors.New("history ID conflicts with committed content"))
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if err := admit(ctx, tx, int64(len(projectionBytes)), s.limits); err != nil {
		return err
	}
	if err := insertEntry(ctx, tx, entry, string(projectionBytes)); err != nil {
		return storeError(ErrWriteFailed, err)
	}
	if err := tx.Commit(); err != nil {
		return storeError(ErrWriteFailed, err)
	}
	if err := secureSQLiteSidecars(s.path, false); err != nil {
		return storeError(ErrWriteFailed, err)
	}
	return nil
}

// Get returns one fully validated entry or HISTORY_ENTRY_NOT_FOUND.
func (s *Store) Get(ctx context.Context, id string) (Entry, error) {
	if err := ctx.Err(); err != nil {
		return Entry{}, storeError(ErrUnavailable, err)
	}
	if !historyIDPattern.MatchString(id) {
		return Entry{}, errors.New("invalid history entry ID")
	}
	if s == nil || s.db == nil {
		return Entry{}, storeError(ErrEntryNotFound, sql.ErrNoRows)
	}
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return Entry{}, storeError(ErrUnavailable, err)
	}
	defer tx.Rollback()
	entry, _, err := readEntryByID(ctx, tx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return Entry{}, storeError(ErrEntryNotFound, err)
	}
	if err != nil {
		return Entry{}, err
	}
	if err := tx.Commit(); err != nil {
		return Entry{}, storeError(ErrUnavailable, err)
	}
	return entry, nil
}

// List returns a bounded, deterministic page of validated index summaries.
func (s *Store) List(ctx context.Context, page Page) (PageResult, error) {
	if err := ctx.Err(); err != nil {
		return PageResult{}, storeError(ErrUnavailable, err)
	}
	limit := page.Limit
	if limit == 0 {
		limit = defaultPageSize
	}
	if limit < 1 || limit > maxPageSize || page.Before != nil && (page.Before.RecordedAtUS < 0 || !historyIDPattern.MatchString(page.Before.HistoryEntryID)) {
		return PageResult{}, errors.New("invalid history page boundary")
	}
	if s == nil || s.db == nil {
		return PageResult{Entries: []Summary{}}, nil
	}
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return PageResult{}, storeError(ErrUnavailable, err)
	}
	defer tx.Rollback()
	query := `SELECT history_entry_id FROM history_entries ORDER BY recorded_at_us DESC, history_entry_id DESC LIMIT ?`
	args := []any{limit + 1}
	if page.Before != nil {
		query = `SELECT history_entry_id FROM history_entries WHERE recorded_at_us < ? OR (recorded_at_us = ? AND history_entry_id < ?) ORDER BY recorded_at_us DESC, history_entry_id DESC LIMIT ?`
		args = []any{page.Before.RecordedAtUS, page.Before.RecordedAtUS, page.Before.HistoryEntryID, limit + 1}
	}
	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return PageResult{}, storeError(ErrCorrupt, err)
	}
	ids := make([]string, 0, limit+1)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return PageResult{}, storeError(ErrCorrupt, err)
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return PageResult{}, storeError(ErrCorrupt, err)
	}
	if err := rows.Close(); err != nil {
		return PageResult{}, storeError(ErrCorrupt, err)
	}
	result := PageResult{Entries: make([]Summary, 0, min(limit, len(ids)))}
	result.HasMore = len(ids) > limit
	if result.HasMore {
		ids = ids[:limit]
	}
	for _, id := range ids {
		entry, _, err := readEntryByID(ctx, tx, id)
		if err != nil {
			return PageResult{}, err
		}
		result.Entries = append(result.Entries, summary(entry.Projection))
	}
	if err := tx.Commit(); err != nil {
		return PageResult{}, storeError(ErrUnavailable, err)
	}
	return result, nil
}

func encodeAndValidate(entry Entry) ([]byte, error) {
	if err := normalize.ValidateProjection(entry.Projection); err != nil {
		return nil, err
	}
	wantRefs := []string{}
	if entry.Projection.Report != nil {
		wantRefs = entry.Projection.Report.ArtifactDigests
	}
	if !slices.Equal(entry.ArtifactRefs, wantRefs) {
		return nil, errors.New("artifact references do not match the projection")
	}
	data, err := json.Marshal(entry.Projection)
	if err != nil {
		return nil, err
	}
	if err := historyv1.ValidateJSON(data); err != nil {
		return nil, err
	}
	return data, nil
}

func summary(p domain.Projection) Summary {
	result := Summary{HistoryEntryID: p.HistoryEntryID, RecordedAtUS: p.RecordedAtUS, OperationOutcome: p.OperationOutcome, RepositoryLabel: cloneString(p.RepositoryLabel)}
	if p.Workspaces != nil {
		count := len(p.Workspaces)
		result.WorkspaceCount = &count
	}
	if p.Report != nil {
		status, count := p.Report.Status, len(p.Report.Findings)
		sourceID, schema := p.Report.SourceScanID, p.Report.SourceSchemaVersion
		result.ReportStatus, result.FindingCount = &status, &count
		result.SourceScanID, result.SourceReportSchemaVersion = &sourceID, &schema
	}
	return result
}

func cloneString(value *string) *string {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

func dataSourceName(path string, writable bool) string {
	path = filepath.ToSlash(path)
	if volume := filepath.VolumeName(filepath.FromSlash(path)); volume != "" && !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	uri := url.URL{Scheme: "file", Path: path}
	query := uri.Query()
	query.Set("mode", "ro")
	query.Set("_busy_timeout", fmt.Sprint(busyTimeout.Milliseconds()))
	query.Add("_pragma", "foreign_keys(ON)")
	if writable {
		query.Set("mode", "rwc")
		query.Set("_txlock", "immediate")
		query.Add("_pragma", "synchronous(FULL)")
	} else {
		query.Add("_pragma", "query_only(ON)")
	}
	uri.RawQuery = query.Encode()
	return uri.String()
}

func (s *Store) acquireWriter(ctx context.Context, code ErrorCode) error {
	if err := ctx.Err(); err != nil {
		return storeError(code, err)
	}
	timer := time.NewTimer(busyTimeout)
	defer timer.Stop()
	select {
	case s.writers <- struct{}{}:
		if err := ctx.Err(); err != nil {
			<-s.writers
			return storeError(code, err)
		}
		return nil
	case <-ctx.Done():
		return storeError(code, ctx.Err())
	case <-timer.C:
		if err := ctx.Err(); err != nil {
			return storeError(code, err)
		}
		return storeError(code, errors.New("writer busy deadline exceeded"))
	}
}
