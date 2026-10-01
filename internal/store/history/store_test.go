package history

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	domain "github.com/geoffrey-xiao/deprail/internal/domain/history"
	"github.com/geoffrey-xiao/deprail/internal/normalize"
)

const testUUID1 = "2d931510-d99f-494a-8c67-87feb05e1594"
const testUUID2 = "2d931510-d99f-494a-8c67-87feb05e1595"
const testRepositoryState = "a1b2c3d4e5f67890a1b2c3d4e5f67890a1b2c3d4e5f67890a1b2c3d4e5f67890"
const testArtifactDigest = "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"

func init() {
	if os.Getenv("DEPRAIL_H05_CRASH_HELPER") != "1" {
		return
	}
	entryBytes, err := base64.StdEncoding.DecodeString(os.Getenv("DEPRAIL_H05_CRASH_ENTRY"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid crash-test entry")
		os.Exit(23)
	}
	var entry Entry
	if err := json.Unmarshal(entryBytes, &entry); err != nil {
		fmt.Fprintln(os.Stderr, "invalid crash-test projection")
		os.Exit(23)
	}
	ctx := context.Background()
	store, err := openAt(ctx, os.Getenv("DEPRAIL_H05_CRASH_CONFIG"), os.Getenv("DEPRAIL_H05_CRASH_REPO"), true, productionLimits)
	if err != nil {
		fmt.Fprintln(os.Stderr, "crash-test writer unavailable")
		os.Exit(23)
	}
	projectionBytes, err := encodeAndValidate(entry)
	if err != nil {
		fmt.Fprintln(os.Stderr, "crash-test projection invalid")
		os.Exit(23)
	}
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "crash-test transaction unavailable")
		os.Exit(23)
	}
	if err := admit(ctx, tx, int64(len(projectionBytes)), store.limits); err != nil {
		fmt.Fprintln(os.Stderr, "crash-test admission failed")
		os.Exit(23)
	}
	if err := insertEntry(ctx, tx, entry, string(projectionBytes)); err != nil {
		fmt.Fprintln(os.Stderr, "crash-test row insertion failed")
		os.Exit(23)
	}
	os.Exit(17)
}

func testRoots(t *testing.T) (string, string) {
	t.Helper()
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	base, err := os.MkdirTemp(home, "deprail-history-test-")
	if err != nil {
		t.Fatal(err)
	}
	canonical, err := filepath.EvalSymlinks(base)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(canonical) })
	config := filepath.Join(canonical, "config")
	repository := filepath.Join(canonical, "repository")
	if err := os.Mkdir(config, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(repository, 0o700); err != nil {
		t.Fatal(err)
	}
	return config, repository
}

func testEntry(t *testing.T, id string, recordedAt int64, digest string) Entry {
	t.Helper()
	digests := []string{}
	if digest != "" {
		digests = []string{digest}
	}
	projection, err := normalize.ProjectHistory(normalize.HistoryInput{
		HistoryEntryID:     id,
		RecordedAtUS:       recordedAt,
		OperationOutcome:   "completed",
		RepositoryLabel:    stringPointer("repo"),
		Workspaces:         []domain.Workspace{},
		Report:             &normalize.HistoryReportInput{SchemaVersion: "v1alpha", ScanID: "same-scan", Status: "complete", RepositoryState: testRepositoryState, Findings: []normalize.HistoryFindingInput{}, ArtifactDigests: digests},
		Diagnostics:        []normalize.HistoryDiagnosticInput{},
		UnclassifiedSource: false,
	})
	if err != nil {
		t.Fatalf("build history projection: %v", err)
	}
	return Entry{Projection: projection, ArtifactRefs: append([]string{}, digests...)}
}

func stringPointer(value string) *string { return &value }

func openTestWriter(t *testing.T, config, repository string, quota limits) *Store {
	t.Helper()
	store, err := openAt(context.Background(), config, repository, true, quota)
	if err != nil {
		t.Fatalf("open history writer: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })
	return store
}

func TestReadOnlyMissingStoreDoesNotCreatePaths(t *testing.T) {
	config, repository := testRoots(t)
	config = filepath.Join(config, "missing-config")
	store, err := openAt(context.Background(), config, repository, false, productionLimits)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	page, err := store.List(context.Background(), Page{})
	if err != nil || len(page.Entries) != 0 || page.HasMore {
		t.Fatalf("missing read-only store result = %#v, %v", page, err)
	}
	if _, err := os.Lstat(config); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("read-only open created config path: %v", err)
	}
}

func TestAppendReopenIdempotencyAndKeysetOrder(t *testing.T) {
	config, repository := testRoots(t)
	store := openTestWriter(t, config, repository, productionLimits)
	first := testEntry(t, testUUID1, 1780000000123456, "")
	second := testEntry(t, testUUID2, 1780000000123456, "")
	if err := store.Append(context.Background(), first); err != nil {
		t.Fatal(err)
	}
	if err := store.Append(context.Background(), first); err != nil {
		t.Fatalf("identical retry failed: %v", err)
	}
	if err := store.Append(context.Background(), second); err != nil {
		t.Fatal(err)
	}
	conflict := first
	conflict.Projection.RepositoryLabel = stringPointer("different")
	if err := store.Append(context.Background(), conflict); !IsCode(err, ErrWriteFailed) {
		t.Fatalf("conflicting retry error = %v", err)
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	reader, err := openAt(context.Background(), config, repository, false, productionLimits)
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	got, err := reader.Get(context.Background(), testUUID1)
	if err != nil || !reflect.DeepEqual(got.Projection, first.Projection) || len(got.ArtifactRefs) != 0 {
		t.Fatalf("reopened entry = %#v, %v", got, err)
	}
	page, err := reader.List(context.Background(), Page{Limit: 1})
	if err != nil || len(page.Entries) != 1 || !page.HasMore || page.Entries[0].HistoryEntryID != testUUID2 {
		t.Fatalf("first keyset page = %#v, %v", page, err)
	}
	cursor := Cursor{RecordedAtUS: page.Entries[0].RecordedAtUS, HistoryEntryID: page.Entries[0].HistoryEntryID}
	page, err = reader.List(context.Background(), Page{Limit: 1, Before: &cursor})
	if err != nil || len(page.Entries) != 1 || page.HasMore || page.Entries[0].HistoryEntryID != testUUID1 {
		t.Fatalf("second keyset page = %#v, %v", page, err)
	}
}

func TestAppendAtomicRefsAndPreservesEarlierRows(t *testing.T) {
	config, repository := testRoots(t)
	store := openTestWriter(t, config, repository, productionLimits)
	prior := testEntry(t, testUUID1, 1780000000123456, "")
	if err := store.Append(context.Background(), prior); err != nil {
		t.Fatal(err)
	}
	if _, err := store.db.Exec(`CREATE TRIGGER fail_refs BEFORE INSERT ON history_artifact_refs BEGIN SELECT RAISE(ABORT, 'injected write failure'); END`); err != nil {
		t.Fatal(err)
	}
	failed := testEntry(t, testUUID2, 1780000000123457, testArtifactDigest)
	if err := store.Append(context.Background(), failed); !IsCode(err, ErrWriteFailed) {
		t.Fatalf("injected ref failure error = %v", err)
	}
	if _, err := store.Get(context.Background(), testUUID2); !IsCode(err, ErrEntryNotFound) {
		t.Fatalf("partially committed entry lookup = %v", err)
	}
	got, err := store.Get(context.Background(), testUUID1)
	if err != nil || !reflect.DeepEqual(got.Projection, prior.Projection) {
		t.Fatalf("earlier committed entry changed: %#v, %v", got, err)
	}
	var refs int
	if err := store.db.QueryRow(`SELECT count(*) FROM history_artifact_refs`).Scan(&refs); err != nil || refs != 0 {
		t.Fatalf("reference count = %d, %v", refs, err)
	}
}

func TestAdmissionExactAndOneOverPreservesCommittedData(t *testing.T) {
	config, repository := testRoots(t)
	first := testEntry(t, testUUID1, 1780000000123456, "")
	second := testEntry(t, testUUID2, 1780000000123456, "")
	encoded, err := encodeAndValidate(first)
	if err != nil {
		t.Fatal(err)
	}
	store := openTestWriter(t, config, repository, limits{entries: 2, bytes: int64(len(encoded))})
	if err := store.Append(context.Background(), first); err != nil {
		t.Fatalf("exact quota append failed: %v", err)
	}
	if err := store.Append(context.Background(), second); !IsCode(err, ErrWriteFailed) {
		t.Fatalf("one-over quota error = %v", err)
	}
	got, err := store.Get(context.Background(), testUUID1)
	if err != nil || !reflect.DeepEqual(got.Projection, first.Projection) {
		t.Fatalf("quota failure changed prior row: %#v, %v", got, err)
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	config2, repository2 := testRoots(t)
	countQuota := openTestWriter(t, config2, repository2, limits{entries: 1, bytes: maxHistoryBytes})
	if err := countQuota.Append(context.Background(), first); err != nil {
		t.Fatalf("exact entry-count quota append failed: %v", err)
	}
	if err := countQuota.Append(context.Background(), second); !IsCode(err, ErrWriteFailed) {
		t.Fatalf("entry-count quota one-over error = %v", err)
	}
	if _, err := countQuota.Get(context.Background(), testUUID1); err != nil {
		t.Fatalf("entry-count quota failure changed prior row: %v", err)
	}
	if err := countQuota.Close(); err != nil {
		t.Fatal(err)
	}

	config3, repository3 := testRoots(t)
	tooSmall := openTestWriter(t, config3, repository3, limits{entries: 2, bytes: int64(len(encoded) - 1)})
	if err := tooSmall.Append(context.Background(), first); !IsCode(err, ErrWriteFailed) {
		t.Fatalf("byte quota one-under error = %v", err)
	}
	page, err := tooSmall.List(context.Background(), Page{})
	if err != nil || len(page.Entries) != 0 {
		t.Fatalf("byte quota failure left rows: %#v, %v", page, err)
	}
}
func TestCancelledAppendPreservesCommittedData(t *testing.T) {
	config, repository := testRoots(t)
	store := openTestWriter(t, config, repository, productionLimits)
	prior := testEntry(t, testUUID1, 1780000000123456, "")
	if err := store.Append(context.Background(), prior); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := store.Append(ctx, testEntry(t, testUUID2, 1780000000123457, "")); !IsCode(err, ErrWriteFailed) {
		t.Fatalf("cancelled append error = %v", err)
	}
	got, err := store.Get(context.Background(), prior.Projection.HistoryEntryID)
	if err != nil || !reflect.DeepEqual(got, prior) {
		t.Fatalf("prior row changed after cancelled append: %#v, %v", got, err)
	}
	if _, err := store.Get(context.Background(), testUUID2); !IsCode(err, ErrEntryNotFound) {
		t.Fatalf("cancelled row result = %v", err)
	}
}

func TestBusyWriterDeadlinePreservesCommittedData(t *testing.T) {
	config, repository := testRoots(t)
	store := openTestWriter(t, config, repository, productionLimits)
	prior := testEntry(t, testUUID1, 1780000000123456, "")
	if err := store.Append(context.Background(), prior); err != nil {
		t.Fatal(err)
	}
	locker, err := sql.Open("sqlite", dataSourceName(store.path, true))
	if err != nil {
		t.Fatal(err)
	}
	defer locker.Close()
	locker.SetMaxOpenConns(1)
	lockTx, err := locker.BeginTx(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	started := time.Now()
	err = store.Append(context.Background(), testEntry(t, testUUID2, 1780000000123457, ""))
	elapsed := time.Since(started)
	if err == nil || !IsCode(err, ErrWriteFailed) {
		t.Fatalf("locked append error = %v", err)
	}
	if elapsed < time.Second || elapsed > 5*time.Second {
		t.Fatalf("writer lock returned after %s, expected bounded ~2s wait", elapsed)
	}
	if err := lockTx.Rollback(); err != nil {
		t.Fatal(err)
	}
	got, err := store.Get(context.Background(), testUUID1)
	if err != nil || !reflect.DeepEqual(got.Projection, prior.Projection) {
		t.Fatalf("busy failure changed prior row: %#v, %v", got, err)
	}
	if _, err := store.Get(context.Background(), testUUID2); !IsCode(err, ErrEntryNotFound) {
		t.Fatalf("busy entry lookup = %v", err)
	}
}

func TestConnectionSafetySettingsAndWAL(t *testing.T) {
	config, repository := testRoots(t)
	store := openTestWriter(t, config, repository, productionLimits)
	var mode string
	var foreignKeys, synchronous int
	if err := store.db.QueryRow(`PRAGMA journal_mode`).Scan(&mode); err != nil || mode != "wal" {
		t.Fatalf("journal mode = %q, %v", mode, err)
	}
	if err := store.db.QueryRow(`PRAGMA foreign_keys`).Scan(&foreignKeys); err != nil || foreignKeys != 1 {
		t.Fatalf("foreign_keys = %d, %v", foreignKeys, err)
	}
	if err := store.db.QueryRow(`PRAGMA synchronous`).Scan(&synchronous); err != nil || synchronous != 2 {
		t.Fatalf("synchronous = %d, %v", synchronous, err)
	}
	if err := store.Append(context.Background(), testEntry(t, testUUID1, 1780000000123456, "")); err != nil {
		t.Fatal(err)
	}
}

func TestFutureSchemaRefusalPreservesDatabase(t *testing.T) {
	config, repository := testRoots(t)
	store := openTestWriter(t, config, repository, productionLimits)
	entry := testEntry(t, testUUID1, 1780000000123456, "")
	if err := store.Append(context.Background(), entry); err != nil {
		t.Fatal(err)
	}
	path := store.path
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	migrationDB, err := sql.Open("sqlite", dataSourceName(path, true))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := migrationDB.Exec(`PRAGMA user_version=2`); err != nil {
		migrationDB.Close()
		t.Fatal(err)
	}
	if err := migrationDB.Close(); err != nil {
		t.Fatal(err)
	}
	before, err := fileSHA256(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := openAt(context.Background(), config, repository, false, productionLimits); !IsCode(err, ErrSchemaUnsupported) {
		t.Fatalf("future schema open error = %v", err)
	}
	after, err := fileSHA256(path)
	if err != nil {
		t.Fatal(err)
	}
	if before != after {
		t.Fatal("future schema refusal modified the database file")
	}
}

func TestReadRejectsRowProjectionMismatch(t *testing.T) {
	config, repository := testRoots(t)
	store := openTestWriter(t, config, repository, productionLimits)
	entry := testEntry(t, testUUID1, 1780000000123456, "")
	if err := store.Append(context.Background(), entry); err != nil {
		t.Fatal(err)
	}
	entry.Projection.RepositoryLabel = stringPointer("tampered")
	raw, err := json.Marshal(entry.Projection)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.db.Exec(`UPDATE history_entries SET projection_json=? WHERE history_entry_id=?`, string(raw), testUUID1); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Get(context.Background(), testUUID1); !IsCode(err, ErrCorrupt) {
		t.Fatalf("mismatched row read error = %v", err)
	}
}

func TestOnlineBackupIncludesCommittedWALState(t *testing.T) {
	config, repository := testRoots(t)
	store := openTestWriter(t, config, repository, productionLimits)
	entry := testEntry(t, testUUID1, 1780000000123456, testArtifactDigest)
	if err := store.Append(context.Background(), entry); err != nil {
		t.Fatal(err)
	}
	backup, err := store.backupForMigration(context.Background())
	if err != nil {
		t.Fatalf("online backup failed: %v", err)
	}
	if err := validateBackup(context.Background(), backup); err != nil {
		t.Fatalf("completed backup failed validation: %v", err)
	}
	firstDigest, err := fileSHA256(backup)
	if err != nil {
		t.Fatal(err)
	}
	secondBackup, err := store.backupForMigration(context.Background())
	if err != nil {
		t.Fatalf("second unique online backup failed: %v", err)
	}
	if secondBackup == backup {
		t.Fatal("online backup reused an existing backup path")
	}
	secondDigest, err := fileSHA256(backup)
	if err != nil {
		t.Fatal(err)
	}
	if firstDigest != secondDigest {
		t.Fatal("creating a second backup overwrote the first backup")
	}
	if err := validateBackup(context.Background(), secondBackup); err != nil {
		t.Fatalf("second online backup failed validation: %v", err)
	}
	copyDB, err := sql.Open("sqlite", dataSourceName(backup, false))
	if err != nil {
		t.Fatal(err)
	}
	copyDB.SetMaxOpenConns(1)
	copyTx, err := copyDB.BeginTx(context.Background(), &sql.TxOptions{ReadOnly: true})
	if err != nil {
		t.Fatal(err)
	}
	snapshot, _, err := readEntryByID(context.Background(), copyTx, testUUID1)
	if err == nil {
		err = copyTx.Commit()
	} else {
		_ = copyTx.Rollback()
	}
	if err != nil || !reflect.DeepEqual(snapshot.Projection, entry.Projection) || !reflect.DeepEqual(snapshot.ArtifactRefs, entry.ArtifactRefs) {
		t.Fatalf("backup omitted committed WAL entry: %#v, %v", snapshot, err)
	}
	if err := copyDB.Close(); err != nil {
		t.Fatal(err)
	}
	if err := secureFile(backup, false); err != nil {
		t.Fatalf("backup permissions are not private: %v", err)
	}
}

func fileSHA256(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	digest := sha256.Sum256(data)
	return hex.EncodeToString(digest[:]), nil
}
func TestStoreRefusesRepositoryLocalHistoryPath(t *testing.T) {
	_, repository := testRoots(t)
	config := filepath.Join(repository, ".config")
	if _, err := openAt(context.Background(), config, repository, true, productionLimits); !IsCode(err, ErrWriteFailed) {
		t.Fatalf("writer repository containment error = %v", err)
	}
	if _, err := os.Lstat(config); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("writer created repository-local config path: %v", err)
	}
	if _, err := openAt(context.Background(), config, repository, false, productionLimits); !IsCode(err, ErrUnavailable) {
		t.Fatalf("reader repository containment error = %v", err)
	}
}
func TestUnversionedDatabaseWithObjectsIsRefusedWithoutModification(t *testing.T) {
	config, repository := testRoots(t)
	directory := filepath.Join(config, ".deprail")
	if err := ensurePrivateDirectory(directory); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(directory, "history.sqlite3")
	db, err := sql.Open("sqlite", dataSourceName(path, true))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`CREATE TABLE unexpected (value TEXT)`); err != nil {
		db.Close()
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	if err := secureFile(path, true); err != nil {
		t.Fatal(err)
	}
	before, err := fileSHA256(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := openAt(context.Background(), config, repository, true, productionLimits); !IsCode(err, ErrCorrupt) {
		t.Fatalf("unversioned object database error = %v", err)
	}
	after, err := fileSHA256(path)
	if err != nil {
		t.Fatal(err)
	}
	if before != after {
		t.Fatal("unversioned corruption refusal modified the database")
	}
}

func TestStoredFutureProjectionVersionIsUnsupported(t *testing.T) {
	config, repository := testRoots(t)
	store := openTestWriter(t, config, repository, productionLimits)
	entry := testEntry(t, testUUID1, 1780000000123456, "")
	if err := store.Append(context.Background(), entry); err != nil {
		t.Fatal(err)
	}
	entry.Projection.SchemaVersion = "history-v2"
	raw, err := json.Marshal(entry.Projection)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.db.Exec(`UPDATE history_entries SET projection_json=? WHERE history_entry_id=?`, string(raw), testUUID1); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Get(context.Background(), testUUID1); !IsCode(err, ErrSchemaUnsupported) {
		t.Fatalf("future projection version error = %v", err)
	}
}

func TestCancelledBackupPreservesCommittedHistory(t *testing.T) {
	config, repository := testRoots(t)
	store := openTestWriter(t, config, repository, productionLimits)
	entry := testEntry(t, testUUID1, 1780000000123456, "")
	if err := store.Append(context.Background(), entry); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := store.backupForMigration(ctx); !IsCode(err, ErrMigrationFailed) {
		t.Fatalf("cancelled backup error = %v", err)
	}
	got, err := store.Get(context.Background(), testUUID1)
	if err != nil || !reflect.DeepEqual(got.Projection, entry.Projection) {
		t.Fatalf("cancelled backup changed committed history: %#v, %v", got, err)
	}
}
func TestSQLitePageLimitFailureRollsBackAppend(t *testing.T) {
	config, repository := testRoots(t)
	store := openTestWriter(t, config, repository, productionLimits)
	prior := testEntry(t, testUUID1, 1780000000123456, "")
	if err := store.Append(context.Background(), prior); err != nil {
		t.Fatal(err)
	}
	var pageCount int
	if err := store.db.QueryRow(`PRAGMA page_count`).Scan(&pageCount); err != nil {
		t.Fatal(err)
	}
	if _, err := store.db.Exec(fmt.Sprintf(`PRAGMA max_page_count=%d`, pageCount)); err != nil {
		t.Fatal(err)
	}
	workspaces := make([]domain.Workspace, 1000)
	for i := range workspaces {
		workspaces[i] = domain.Workspace{
			WorkspaceID:           fmt.Sprintf("w%04d", i),
			Path:                  fmt.Sprintf("packages/p%04d", i),
			Ecosystem:             "npm",
			PackageManager:        "npm",
			DiscoveryCompleteness: "complete",
		}
	}
	largeProjection, err := normalize.ProjectHistory(normalize.HistoryInput{
		HistoryEntryID:   testUUID2,
		RecordedAtUS:     1780000000123457,
		OperationOutcome: "completed",
		Workspaces:       workspaces,
		Diagnostics:      []normalize.HistoryDiagnosticInput{},
	})
	if err != nil {
		t.Fatal(err)
	}
	err = store.Append(context.Background(), Entry{Projection: largeProjection, ArtifactRefs: []string{}})
	if !IsCode(err, ErrWriteFailed) {
		t.Fatalf("SQLite page-limit append error = %v", err)
	}
	got, err := store.Get(context.Background(), testUUID1)
	if err != nil || !reflect.DeepEqual(got.Projection, prior.Projection) {
		t.Fatalf("SQLite full simulation changed prior row: %#v, %v", got, err)
	}
	if _, err := store.Get(context.Background(), testUUID2); !IsCode(err, ErrEntryNotFound) {
		t.Fatalf("SQLite full simulation left a new row: %v", err)
	}
}
func TestProcessExitBeforeCommitRollsBackRowAndRefs(t *testing.T) {
	config, repository := testRoots(t)
	writer := openTestWriter(t, config, repository, productionLimits)
	prior := testEntry(t, testUUID1, 1780000000123456, "")
	if err := writer.Append(context.Background(), prior); err != nil {
		t.Fatal(err)
	}
	path := writer.path
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	uncommitted := testEntry(t, testUUID2, 1780000000123457, testArtifactDigest)
	entryBytes, err := json.Marshal(uncommitted)
	if err != nil {
		t.Fatal(err)
	}
	command := exec.Command(os.Args[0], "-test.run=^TestReadOnlyMissingStoreDoesNotCreatePaths$")
	command.Env = append(os.Environ(),
		"DEPRAIL_H05_CRASH_HELPER=1",
		"DEPRAIL_H05_CRASH_CONFIG="+config,
		"DEPRAIL_H05_CRASH_REPO="+repository,
		"DEPRAIL_H05_CRASH_ENTRY="+base64.StdEncoding.EncodeToString(entryBytes),
	)
	output, err := command.CombinedOutput()
	exit, ok := err.(*exec.ExitError)
	if !ok || exit.ExitCode() != 17 {
		t.Fatalf("crash helper exit = %v, output %q", err, output)
	}
	reader, err := openAt(context.Background(), config, repository, false, productionLimits)
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	if got, err := reader.Get(context.Background(), testUUID1); err != nil || !reflect.DeepEqual(got.Projection, prior.Projection) {
		t.Fatalf("process interruption changed prior row: %#v, %v", got, err)
	}
	if _, err := reader.Get(context.Background(), testUUID2); !IsCode(err, ErrEntryNotFound) {
		t.Fatalf("uncommitted process row survived: %v", err)
	}
	var refs int
	verify, err := sql.Open("sqlite", dataSourceName(path, false))
	if err != nil {
		t.Fatal(err)
	}
	if err := verify.QueryRow(`SELECT count(*) FROM history_artifact_refs`).Scan(&refs); err != nil {
		verify.Close()
		t.Fatal(err)
	}
	if err := verify.Close(); err != nil {
		t.Fatal(err)
	}
	if refs != 0 {
		t.Fatalf("uncommitted process refs survived: %d", refs)
	}
}
