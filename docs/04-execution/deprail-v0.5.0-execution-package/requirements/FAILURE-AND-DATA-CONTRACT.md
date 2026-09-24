# v0.5 Failure and Data Contract

**Status:** Draft proposal; SQLite schema, storage errors, and retention are not approved.

## 1. Data integrity invariants

- A stored history entry has a unique, stable `historyEntryID` for that scan-operation occurrence. It is the history primary key and API resource identifier; it is not derived solely from repository state.
- Preserve the source report's `ScanReport.ScanID` separately as `sourceScanID`. Current `ScanReport.ScanID` is derived from repository state and repeated runs can legitimately share it; never use it as a unique history key or silently change its existing meaning.
- A history entry separately records an operation outcome (`completed`, `failed`, or `cancelled`) and, when a report exists, its unchanged report completeness (`complete`, `partial`, or `failed`). A cancelled operation is not a new `ScanReport.Status` value.
- A cancelled operation may carry a partial report returned by the current scan API. The history/UI must retain the operation outcome as Cancelled and must not present the report's pre-cancellation completeness as proof that the run completed. If no trustworthy report exists, do not synthesize an empty report.
- The report preserves its report/schema version, tool/scanner/database provenance, and artifact digest references required by existing contracts.
- Persistence must not convert a failed/partial report to complete, or a missing/invalid report to empty-success.
- The history index and referenced raw artifacts are separate resources. Dangling/missing artifacts and digest mismatches are explicit, scoped errors; do not silently delete metadata or invent evidence.
- Writes are atomic at the documented transaction boundary. Assign a history ID once per operation occurrence; define retry/idempotency behavior separately and never deduplicate entries solely by `sourceScanID`.
- No source files, credentials, complete process environments, or unrelated repository contents are stored by default.
- Exact stored fields, ID generation/ingestion idempotency, payload-size thresholds, artifact retention/reference counting, and history capture trigger require an approved schema/ADR.

## 2. State vocabulary

Keep three concepts distinct:

- Report completeness: `complete`, `partial`, or `failed`, matching the current v1alpha contract. Do not add `cancelled` to the existing report schema without a separately approved schema change.
- History operation outcome: `completed` when the scan operation returns without an operation error, `cancelled` when its error is cancellation, or `failed` for another fatal operation error. A completed operation may still have a `partial` or `failed` report; neither axis is inferred from the other.
- History resource identity: unique `historyEntryID` for this operation occurrence; `sourceScanID` preserves the existing report identity and may repeat across distinct runs.
- API request outcome: transport/application success or a typed request failure/cancellation. Cancelling a read request does not rewrite persisted history or scan-report state.

These are semantic states, not finalized wire codes. Stable error codes and HTTP status mapping must be cross-walked with [`ERROR-MODEL.md`](ERROR-MODEL.md) and the exact prior versioned contracts before implementation.

## 3. Failure matrix (proposed)

| Condition | Required observable behavior | Data-safety requirement | Verification ID |
| --- | --- | --- | --- |
| Empty history | Successful empty collection and empty-state UI | Do not synthesize a scan record | `FR-502` |
| Missing data directory/database | Read-only history query returns an empty collection without creating files; only an explicitly approved save initializes storage | Do not make a normal scan depend on history initialization | `FR-502`, `FR-505` |
| Permission denied | Explicit storage-unavailable error | Do not weaken permissions or write elsewhere silently | `SEC-02`, `FR-505` |
| Database locked/concurrent writer | Wait only for the bounded busy interval; then return a typed read/write failure | No partial transaction visible; retain the independent scan outcome | `STORE-01` |
| Disk full/quota | Abort the selected save with `HISTORY_WRITE_FAILED` | Roll back the full transaction; preserve the scan outcome and prior entries | `STORE-01`, `FR-508` |
| Interrupted write/process crash | SQLite transaction rollback leaves the prior committed state readable | No half-record or ambiguous “saved” result | `FR-508`, `SEC-08` |
| Unsupported future schema | Refuse unsafe reads/writes with upgrade guidance | Never downgrade, delete, or overwrite automatically | `FR-508`, `SEC-08` |
| Migration error/interruption | Roll back the versioned transaction and preserve the pre-migration backup | Preserve the old database and provide manual recovery guidance | `FR-508`, `SEC-08` |
| Corrupt DB/page | Explicit integrity/storage error | No automatic destructive reset; preserve evidence for recovery | `FR-508`, `SEC-08` |
| Missing referenced artifact | Scan metadata may be shown if valid; integrity warning/error for evidence | Never substitute empty artifact or claim provenance verified | `FR-503`, `SEC-09` |
| Digest mismatch | Integrity failure at scoped read | Do not trust corrupted bytes or mutate the source artifact | `FR-503`, `SEC-09` |
| Malformed stored `history-v1` projection or row/JSON mismatch | Explicit invalid-record result scoped to entry | Do not expose malformed content as a valid report or empty page | `FR-503`, `SEC-08` |
| Unsafe source projection (absolute path, unclassified raw error/field, incomplete or duplicate finding identity) | Refuse selected history save with candidate `HISTORY_WRITE_FAILED` before commit | No partial row/reference, no saved claim, original scan result and previous DB unchanged | `STORE-01`, `SEC-07` |
| Requested history save fails after scan success | Report `HISTORY_WRITE_FAILED` separately; candidate exit code `6` is subject to CLI review | Preserve the scan report and outcome; do not claim it was saved | `STORE-01` |
| Selected history entry exceeds its approved size bound | Fail the save explicitly before commit | No partial entry; preserve the scan outcome | `SEC-05`, `FR-508` |
| Request oversized/malformed | Client error; bounded work | Reject before expensive parsing/DB query | `FR-507`, `SEC-05` |
| Cancelled scan operation | History entry records `operationOutcome=cancelled`; an available projected report retains the source report's unchanged status | Never treat retained pre-cancellation completeness as proof of completion | `FR-504` |
| API timeout/request cancellation | Explicit request outcome; no fabricated complete response | Do not rewrite stored history or `ScanReport.Status` | `FR-502`, `FR-503` |
| UI assets missing/version mismatch | Explicit console unavailable/error | CLI remains usable | `FR-509`, `SEC-10` |
| Shutdown during query/write | Cancel/drain according to the accepted transaction model | Durable commit or rollback; no ambiguous “saved” result | `FR-508`, `STORE-01` |

## 4. Proposed persisted schema v1 (unapproved)

Database schema version is the monotonic integer in `PRAGMA user_version`; the proposed initial value is `1`. The owner-selected, still-unapproved stored projection has its own candidate `schema_version: history-v1`, separate from a source report's existing `v1alpha` version and from `/api/v1`. The database schema is not an API response or a claim that current `app.ScanReport` validates against the v1alpha scan JSON Schema.

| Object | Proposed v1 representation |
| --- | --- |
| `history_entries.history_entry_id` | Canonical UUIDv4 text primary key, generated once per operation occurrence; never derived from `source_scan_id`. |
| `history_entries.recorded_at_us` | UTC Unix microseconds when the history entry is recorded; list ordering is `recorded_at_us DESC, history_entry_id DESC`. |
| `history_entries.operation_outcome` | Required enum `completed`, `failed`, or `cancelled`. |
| `history_entries.source_scan_id` | Nullable source `ScanReport.ScanID`; non-unique because distinct operations can share it. |
| `history_entries.report_status` | Nullable report completeness `complete`, `partial`, or `failed`; independent of operation outcome. |
| `history_entries.repository_label` | Nullable reviewed/redacted `ProjectGraph.DisplayName` display value, never a copied project snapshot or absolute host path. |
| `history_entries.workspace_count`, `finding_count` | Nullable nonnegative list-summary counts validated against the versioned projection; absent workspace/report data is null, never a fabricated zero. |
| `history_entries.projection_schema_version`, `projection_json` | Required `history-v1` marker and schema-validated, allowlisted projection of safe project/report fields and diagnostics; never an unchanged project/scan document or raw error string. The version is repeated inside JSON and must agree with the column. |
| `history_entries.source_report_schema_version` | Nullable original `ScanReport.SchemaVersion` metadata, present only with a trustworthy projected report; not the projection's schema version. |
| Projected diagnostics | Bounded, safe structured diagnostics inside `projection_json`; raw `ScanReport.Errors` strings must not be copied or inferred to be safe. |
| `history_artifact_refs` | `(history_entry_id, digest)` rows for validated report artifact digests; foreign-keyed to the history row, with a non-unique digest index for reference queries. |

**Owner-selected direction, independent review pending:** The earlier `project_json`/`report_json`/`operation_diagnostics_json` payload proposal is superseded as a *candidate* because unchanged CLI scan JSON contains an absolute host path and non-v1alpha finding/error shapes. [ADR-0004](../../../adr/ADR-0004-local-scan-history.md#owner-selected-direction-independent-history-projection-not-accepted) preserves the original proposal and records this choice. The revised candidate below persists only a separately versioned history projection; existing CLI/report serialization and remediation inputs remain unchanged. No selection of exact fields, numeric bounds, error codes or SQL has independent approval.

The revised candidate relational shape is:

```sql
CREATE TABLE history_entries (
  history_entry_id TEXT PRIMARY KEY NOT NULL,
  recorded_at_us INTEGER NOT NULL,
  operation_outcome TEXT NOT NULL
    CHECK (operation_outcome IN ('completed', 'failed', 'cancelled')),
  source_scan_id TEXT,
  report_status TEXT
    CHECK (report_status IS NULL OR report_status IN ('complete', 'partial', 'failed')),
  repository_label TEXT,
  workspace_count INTEGER CHECK (workspace_count IS NULL OR workspace_count >= 0),
  finding_count INTEGER CHECK (finding_count IS NULL OR finding_count >= 0),
  projection_schema_version TEXT NOT NULL
    CHECK (projection_schema_version = 'history-v1'),
  projection_json TEXT NOT NULL,
  source_report_schema_version TEXT,
  CHECK (
    (source_scan_id IS NULL AND source_report_schema_version IS NULL AND
     report_status IS NULL AND finding_count IS NULL)
    OR
    (source_scan_id IS NOT NULL AND source_report_schema_version IS NOT NULL AND
     report_status IS NOT NULL AND finding_count IS NOT NULL)
  )
);

CREATE INDEX history_entries_order
  ON history_entries (recorded_at_us DESC, history_entry_id DESC);

CREATE TABLE history_artifact_refs (
  history_entry_id TEXT NOT NULL
    REFERENCES history_entries(history_entry_id) ON DELETE CASCADE,
  digest TEXT NOT NULL
    CHECK (length(digest) = 64 AND digest NOT GLOB '*[^0-9a-f]*'),
  PRIMARY KEY (history_entry_id, digest)
);

CREATE INDEX history_artifact_refs_digest
  ON history_artifact_refs (digest);
```

Enable and verify `PRAGMA foreign_keys=ON` for every connection. Application validation must check UUIDv4 syntax, the independent projection schema version, source-report metadata, row/projection identity and summary-count consistency, unique deterministic finding keys, safe relative paths/diagnostics, and that artifact-reference rows exactly match the validated projection digest set. `source_scan_id` is intentionally not unique. No raw project/report JSON or artifact bytes are stored in SQLite, and a foreign key never implies deletion of an artifact file. No per-entry deletion operation is proposed in v0.5.

### Candidate `history-v1` projection object

This is a proposed internal storage shape, **not** the existing `v1alpha` scan schema or the public OpenAPI schema. Every object uses an explicit allowlist and is validated on write and read. Row columns repeated in the projection must match exactly; source CLI/report bytes remain unchanged.

| Projected field | Source and required behavior |
| --- | --- |
| `schema_version`, `history_entry_id`, `recorded_at_us`, `operation_outcome` | `history-v1` and the corresponding row values; recorded time is one UTC microsecond value per occurrence. |
| `repository_label`, `workspaces` | Nullable reviewed/redacted project display label and nullable workspace collection from the validated `ProjectGraph` of the **same scan operation**, never after-the-fact rediscovery. `app.Scan` currently returns only `ScanReport`, so access to that graph requires a separately reviewed application contract; absent trustworthy context gives null workspaces/count, not an empty confirmed list. When available, each workspace contains reviewed ID, repository-relative `/` path, ecosystem, package manager and explicitly **discovery-only** completeness, not inferred scanner success. No raw `ProjectGraph` JSON or unreviewed diagnostics/details are copied. |
| `report` | Null if there is no trustworthy source report. Otherwise retain source schema version, source scan ID, unchanged completeness, opaque repository-state identity, sorted finding summaries and validated/sorted raw-artifact digests. **Never** include `repository_identity.root`, raw report/errors, source files, host paths, or unverified tool/database provenance. No report means nullable report columns and `finding_count`, not zero findings. |
| `report.findings[]` | Reviewed component identity (`PURL` when present, otherwise `Component`, matching [`baseline.ConvertScan`](../../../../internal/baseline/convert.go)'s fallback), version, workspace ID/relative path, ecosystem, vulnerability ID (`adapter.Finding.TargetID`), sorted aliases, available severity/fixed version, and stable finding key. A missing PURL remains absent, never a fabricated PURL. Reuse [`normalize.StableFindingKey`](../../../../internal/normalize/key.go) with the selected component identity, workspace ID, version, vulnerability ID, and aliases; reject missing/duplicate identity, do not label `TargetID` as a finding key, and sort by stable key. Do not invent dependency paths, severity sources or raw evidence absent from the input. |
| `diagnostics[]` | Only bounded, redacted, typed application diagnostics with approved code/scope/message. The current `ScanReport.Errors []string` is not copied. If safe diagnostics cannot truthfully represent a selected partial/failed/cancelled operation, fail the history save before commit and retain the original scan result. |

Reject an unsupported projection version, unreviewed required source field, invalid path/identifier, malformed digest, or row/projection mismatch with a typed history failure; do not silently drop security-relevant unknown data or coerce it into an empty/complete result. If graph context is unavailable, finding workspace paths are validated independently as repository-relative source-report values; neither those findings nor the discovery graph alone establish all scanned workspaces or per-workspace scan completeness. Unknown fields are preserved only after explicit safe-field classification and versioned projection review. The exact schema artifact, same-operation graph/result contract, limits and error-code mapping remain approval gates; a successful scan by itself is not proof that its history projection was saved.

The initial query index supports deterministic history-list ordering. Additional filters/indexes, child-row normalization, and exact API projections remain open until OpenAPI and payload evidence are accepted. Enforce per-entry and response bounds before expensive work; numeric thresholds are not selected from small fixtures. Large projected finding lists may make child-page extraction expensive, so memory/performance evidence is an acceptance gate rather than an assumed property.

## 5. Migration and recovery (draft recommendation)

Use `PRAGMA user_version` as the single database-schema version. Migrations are ordered, forward-only, and transactional; set the new version in the same transaction as the schema/data changes. Before a future migration, create a consistent SQLite Online Backup API snapshot in the same per-user data root with a unique, non-overwriting name and restrictive permissions. Validate the backup before migration; do not copy only the main database file while WAL is active. If backup creation/validation fails, refuse migration and leave the original untouched.

On migration failure or interruption, roll back and retain the prior database and backup. A newer/unsupported version is read-only unavailable for this binary and is never downgraded. Do not auto-repair, replace, or restore over the only current copy. Recovery is a documented manual action: stop DepRail, preserve the current database and WAL/SHM files, validate a backup, restore to a new file, and only then perform an explicit operator-controlled replacement. No backup/restore command is included in v0.5.

Required migration evidence covers: absent/fresh store; current v1; every explicitly supported prior version; interruption before/during/after migration; corrupt database; and unsupported future `user_version`. Every failed case must retain the last committed history and a recoverable copy.

## 6. Retention and deletion (draft recommendation)

Only operations explicitly selected for history capture are retained. Proposed v0.5 behavior has no automatic eviction, per-entry delete, export, browser/API mutation, or automatic raw-artifact cleanup. History references remain separate from content-addressed artifact bytes; missing or mismatched artifacts are reported explicitly. No artifact is removed based only on a history-row change.

This retention choice permits unbounded database, backup, and artifact growth. A numeric entry/store cap or an explicitly accepted disk-growth risk is required before DoR; no cap is inferred from synthetic fixtures. If a cap is selected, refusal to save a new entry must be explicit and must not modify the scan outcome. User-facing location/retention guidance is required; destructive reset/delete controls remain out of scope.

## 7. Transaction and concurrency contract (draft recommendation)

Use SQLite WAL with `synchronous=FULL`, one serialized writer and read-only readers. Set and verify WAL mode when initializing the database; enable foreign keys on every connection and set/verify `synchronous=FULL` on each writer connection. A save uses one write transaction for the history row, validated projection, and all artifact-reference rows; use `BEGIN IMMEDIATE` and a finite busy timeout. The numeric timeout is an evidence-backed open choice. Readers use a stable snapshot for each bounded query. Do not perform unbounded application retries.

If a lock timeout, disk-full condition, permission failure, or transaction error aborts the save, roll back the entire entry and return candidate `HISTORY_WRITE_FAILED`. A successful scan remains successful as a scan; an explicitly requested history save failure is a separate persistence outcome. Candidate CLI behavior is unchanged report JSON on stdout, a safe persistence diagnostic on stderr, and exit code `6` only when the scan otherwise succeeded. If scan and history persistence both fail, preserve the original scan error and report the persistence failure separately. The additive flag/exit mapping needs a separate CLI compatibility decision before adoption.

## 8. Redaction and privacy

Use the proposed per-user `.deprail` root outside the scanned repository, with canonical-path and symlink/reparse-point checks. If the configured root lies inside the scanned tree or permissions cannot be enforced, fail history access/save without a repository-relative fallback. POSIX directory/database/backup modes are `0700`/`0600`; Windows requires an owner-only ACL. Do not persist either source `repository_root` or `ScanReport.repository_identity.root`; project workspace paths in the `history-v1` projection are repository-relative `/` paths. Do not persist raw source files/reports/errors, secrets, full environments, credential-bearing URLs, or sensitive absolute host paths.
