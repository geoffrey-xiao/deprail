# v0.5 Failure and Data Contract

**Status:** Detailed proposed v0.5 contract. This document does not constitute owner acceptance or authorize implementation.

## Current detailed proposal (takes precedence over historical draft below)

The sections below preserve the historical draft; where they disagree, this current detailed proposal and its linked schema control. These choices remain proposals and do not constitute owner acceptance or authorize implementation. The projection schema is [`schemas/history-v1/history.schema.json`](../../../../schemas/history-v1/history.schema.json); examples are in its `examples/` directory.

### Projection, nullability, and identity

The projection is an explicit allowlist with `additionalProperties:false` at every object. Unknown source fields are never copied or silently discarded to claim a complete faithful capture: fields outside the projection are ignored only when they are not required to derive a projected value; an unclassified field that affects safe interpretation, or an unsupported source report schema version, refuses capture as `HISTORY_WRITE_FAILED`. On read, unknown persisted fields or unsupported `schema_version` are `HISTORY_SCHEMA_UNSUPPORTED` for a future projection version and `HISTORY_CORRUPT` for malformed data claiming `history-v1`. Source `source_schema_version` is metadata, distinct from history projection version.

The application allocates one canonical lowercase UUIDv4 `history_entry_id` and one nonnegative UTC Unix-microsecond `recorded_at_us` per selected operation occurrence. Retries within that same capture call reuse the same ID; a later scan operation always gets a new ID, including when `source_scan_id` repeats. A persisted retry after uncertain commit first checks that exact ID and verifies identical content; identical committed content is idempotent success, different content is `HISTORY_WRITE_FAILED`, and no retry creates a second occurrence.

`report:null` means no trustworthy `ScanReport` exists; it is never converted to an empty report. A non-null report with `findings:[]` is a trustworthy report containing zero findings. `workspaces:null` means trustworthy same-operation discovery context is unavailable; `workspaces:[]` means discovery context exists and found zero workspaces. Null report means the SQL report identity/status/schema-version/finding-count columns are all NULL. Non-null report means those values are present and equal projection values; its finding count is the array length. Workspace count is NULL exactly when `workspaces` is null, otherwise its length. Report status is independent of operation outcome.

Source projection uses the graph from the same scan invocation only. No rediscovery is permitted after scanning; absent graph context produces null workspace data/count. Graph completeness is discovery-only, never scanner completeness. Preserve the exported `Scan` signature, `ScanReport` JSON keys/status, and every existing caller: `ScanOptions` has an optional internal capture sink; `Scan` invokes it once with the already computed `ProjectGraph` (when available), final report, operation error and typed boundary diagnostics before returning. No serialization fields are added and scanning is never repeated.

Finding component identity is the nonempty PURL when present, otherwise nonempty Component, matching `baseline.ConvertScan`; however component display name is always required and bounded to 256 characters to satisfy API representation. Missing display name, identity, version, workspace ID/path, ecosystem, or TargetID refuses capture. Stable identity is the existing `normalize.StableFindingKey` over selected component identity, workspace ID, version, TargetID and aliases; it is never TargetID. Duplicate stable keys reject capture. Findings sort by stable key; aliases and artifact digests are sorted and unique. No dependency paths or unsupported severity provenance are invented. Preserve ecosystem values at API boundary: validated graph ecosystem is copied as-is; finding `javascript`/`java` map through established app conversion to `npm`/`maven`; other values must satisfy the OpenAPI ecosystem enum.

All workspace and finding paths are repository-relative slash paths, with no absolute paths, backslashes, NUL, empty/dot traversal segments, or `..`; containment and symlink checks are additionally runtime responsibilities. Repository root is never stored. Artifact references in history are digest identifiers only. Exact artifact integrity reads use the existing content-addressed artifact store's configured root and validated digest, verify bytes against the digest on access, and return missing/mismatch/unavailable integrity without exposing filesystem paths. Since history stores no absolute artifact path, the implementation must resolve by digest through an explicitly supplied trusted artifact-store root; if no safe root is available, report unavailable rather than guessing from repository paths.

Projection size is at most 16 MiB serialized UTF-8; additionally cap at 10,000 findings, 1,000 workspaces, 128 diagnostics, and 4,096 digest references. Store quota is 1,000 entries and 256 MiB logical serialized projection bytes. No eviction or delete: admission exceeding either quota rejects the selected save atomically as `HISTORY_WRITE_FAILED`; no cap accounting includes SQLite/WAL overhead or external artifact bytes. Lock wait is 2,000 ms and there are no unbounded retries.

### SQL and source mapping

The SQL shape in §4 is retained. Before exposing any row, parse and validate the JSON Schema, then check canonical UUIDv4, UTC microsecond integer, enum values, source metadata, row/projection equality, finding key uniqueness and ordering, sorted unique digest set, and exact equality of `history_artifact_refs`. Any inconsistency makes that record `HISTORY_CORRUPT`; it is not skipped into an apparently successful empty page. Use WAL, `synchronous=FULL`, `foreign_keys=ON` on every connection, one serialized writer, read-only snapshot readers, and a single `BEGIN IMMEDIATE` transaction for row and all refs. Commit only after full projection validation and quota admission; any failure rolls back all row/ref changes. POSIX root/database/backup modes are 0700/0600; Windows owner-only ACL. Canonical root must be outside scanned tree.

| Typed source | Stored projection | API view |
| --- | --- | --- |
| Allocated canonical operation UUID | `history_entry_id` | `historyEntryID` |
| UTC Unix microseconds (epoch through year 9999) | `recorded_at_us` | RFC3339 UTC `recordedAt`, six fractional digits |
| Terminal operation error classification | `operation_outcome` | `operationOutcome` (independent of report status) |
| Validated `RepositoryIdentity.Repository`, otherwise null | `repository_label` | `repositoryLabel`, never the root |
| Same-operation graph absent / present | `workspaces:null` / array | null / exact length `workspaceCount`; child unavailable / available state |
| Graph workspace ID, relative path, ecosystem, manager, completeness | `workspace_id`, `path`, `ecosystem`, `package_manager`, `discovery_completeness` | `workspaceID`, `relativePath`, `ecosystem`, `packageManager`, `discoveryCompleteness` |
| No trustworthy report / supported report | `report:null` / object | null report metadata/counts / available report and finding collection |
| Report `SchemaVersion`, `ScanID`, `Status`, `RepositoryState` | `source_schema_version`, `source_scan_id`, `status`, `repository_state` | `sourceReportSchemaVersion`/`sourceSchemaVersion`, `sourceScanID`, `reportStatus`, `repositoryState` |
| Existing normalized identity computation (§Projection) | `stable_finding_key` | `stableFindingKey`, not TargetID |
| Finding Component, PURL, Version, WorkspaceID, WorkspacePath, TargetID | `component_name`, `component_purl` (null if absent), `version`, `workspace_id`, `relative_path`, `vulnerability_id` | `componentName`, `componentPURL`, `version`, `workspaceID`, `relativePath`, `vulnerabilityID` |
| Finding ecosystem (§Projection), aliases, severity, fixed version | `ecosystem`, sorted unique `aliases`, supported `severity`, nullable `fixed_version` | `ecosystem`, `aliases`, `severity`, `fixedVersion`; no unsupported detail invented |
| Sorted unique Report `ArtifactDigests` | `artifact_digests` and matching SQL refs | `artifactReferences[{digest,integrity}]`; unavailable without resolver |
| Typed same-operation provenance, absent values null | `report.provenance.{deprail_version,scanner_name,scanner_version,scanner_database_version}` | required `provenance.{deprailVersion,scannerName,scannerVersion,scannerDatabaseVersion}` |
| Typed diagnostic registry | `diagnostics[{code,scope,message,workspace_id}]` | `{code,scope,message,workspaceID}`; static paired messages only |
| `ScanReport.Errors`, root, raw output, descriptions and unknown source members | never copied | never exposed |


Each provenance string is either null or a nonempty bounded reviewed value. Populate only from typed application/scanner metadata obtained at an actual source boundary; the current report does not itself carry reliable scanner/database version fields. Do not infer scanner name (including `OSV`) or versions from findings, command output strings, or errors. Unavailable metadata is null; a present but unsafe value containing a credential, URL, absolute path, control bytes, or otherwise unreviewed source text refuses capture rather than being copied or replaced with null.
The projection permits maxima larger than an API page. API serialization remains governed by existing OpenAPI bounds: whole-record byte-aware pages, default count 25, maximum 50, 1 MiB response cap. API projection is a bounded view of stored history, not raw projection pass-through; each returned component must satisfy OpenAPI maxima. A stored value that cannot be represented safely is a scoped integrity/response failure, never truncation or altered identity.

### Typed diagnostics and operation/capture failures

Only emit diagnostics at observed typed boundaries; never parse `Error()` strings. Schema `oneOf` enforces exact code/message pairs. Diagnostic scope is `repository` or `workspace`, with nullable `workspace_id` nonnull exactly for workspace scope. Graph diagnostic source codes map statically: `MANIFEST_INVALID` → invalid manifest; `DISCOVERY_INCOMPLETE` → incomplete discovery; `WALK_ENTRY_FAILED` → failed walk entry; `PATH_OUTSIDE_ROOT` → outside root; `SYMLINK_SKIPPED` → skipped symlink. Their schema messages are fixed; source message/path is never copied. Workspace scope only when containing graph workspace ID is known; otherwise repository scope. `Discover` error → `DISCOVERY_FAILED`; plan validation/adapter `ErrInvalidPlan` → `CONFIG_INVALID`; adapter `ErrUnsupportedTarget` → `SCANNER_VERSION_UNSUPPORTED`; Execute adapter errors use `adapter.IsCode` to retain `SCANNER_NOT_FOUND`, `SCANNER_EXIT_NONZERO`, `SCANNER_TIMEOUT`, `SCANNER_OUTPUT_LIMIT`, or `SCANNER_OUTPUT_INVALID`; unknown Execute errors refuse capture. Artifact `Put` → `ARTIFACT_STORE_FAILED`; Parse `ErrInvalidOutput` → `SCANNER_OUTPUT_INVALID`, unknown parse errors refuse capture; `Normalize` → `FINDING_NORMALIZATION_FAILED`; `NormalizeComponent` → `COMPONENT_IDENTITY_INVALID`. Typed cancellation at a scan boundary may yield `CANCELLED`; never infer from text. Event-sink/other unclassified errors refuse capture as `HISTORY_WRITE_FAILED`. Diagnostics sort by `(scope, workspace_id-or-empty, code)`.

Only source report `v1alpha` is supported. API bounds apply before commit: scan ID <=128; state/key lowercase 64-hex; component display name <=256; PURL `pkg:` <=2048; version/fixed <=128; workspace ID <=256; aliases <=32 of <=256; API severity/ecosystem enums; paths <=4096. Workspace root path `.` is explicitly allowed; no other dot segment, traversal, absolute/drive-root, backslash, or control byte is allowed.

An operation returning nil error is `completed`, including a partial/failed report. Context cancellation is `cancelled`; any other fatal operation error is `failed`. A report accompanying cancellation remains unchanged and is retained only if trustworthy; otherwise `report:null`. Capture uses a fresh persistence context bounded to 3 seconds overall and a 2-second SQLite busy timeout; timeout/cancellation interrupts and rolls back capture without changing operation outcome/report. Selected-save failure preserves scan result and is separately reported. CLI stdout JSON remains unchanged; `--save-history` failure uses stderr and exit 6 only if scan otherwise succeeded. Existing scan failure/cancellation exit takes precedence. No save request means no capture.

### Fresh store and retention

The first database created by v0.5 starts directly at `PRAGMA user_version=1` in one transaction; there is no prior shipped history schema or import migration. A database with `user_version=0` is initialized only if it has no DepRail history tables/data; any unexpected existing objects/data are refused as corrupt, never overwritten. Future migrations are forward-only transactions, preceded by validated SQLite Online Backup API backup with unique non-overwriting name; backup failure blocks migration. Unsupported future DB version is read-only unavailable; no downgrade, auto-repair, reset, restore, automatic eviction, or deletion. Preserve current DB/WAL/SHM and backup for manual recovery.

### Validation scenarios

Schema examples must validate. Reject: 10001 findings, 1001 workspaces, >128 diagnostics, >4096 digests, >16 MiB serialization, uppercase/non-v4 UUID, unsorted/duplicate finding keys or digests, absolute/traversing/backslash paths, missing both PURL and component, duplicate stable keys, null report with non-null report SQL columns, row/projection count/version/identity disagreement, artifact-ref mismatch, future projection marker, unclassified source error, missing/unsafe identity, and capture exceeding entry/store quota. Confirm null report differs from complete report with zero findings; null workspaces differs from empty graph; repeated source scan IDs yield distinct UUIDs; canceled scan capture uses a fresh persistence context and remains canceled if capture fails; transaction fault leaves no row or refs.

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

**Owner-selected direction; technical decisions pending:** The earlier `project_json`/`report_json`/`operation_diagnostics_json` payload proposal is superseded as a *candidate* because unchanged CLI scan JSON contains an absolute host path and non-v1alpha finding/error shapes. [ADR-0004](../../../adr/ADR-0004-local-scan-history.md#owner-selected-direction-independent-history-projection-not-accepted) preserves the original proposal and records this choice. The revised candidate below persists only a separately versioned history projection; existing CLI/report serialization and remediation inputs remain unchanged. Exact fields, bounds, error codes and SQL require owner acceptance and evidence; external review is optional under ADR-0005.

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
| `repository_label`, `workspaces` | Nullable redacted project display label and nullable workspace collection from the validated `ProjectGraph` of the **same scan operation**, never after-the-fact rediscovery. `app.Scan` currently returns only `ScanReport`, so access to that graph requires an owner-accepted application contract; absent trustworthy context gives null workspaces/count, not an empty confirmed list. When available, each workspace contains reviewed ID, repository-relative `/` path, ecosystem, package manager and explicitly **discovery-only** completeness, not inferred scanner success. No raw `ProjectGraph` JSON or unreviewed diagnostics/details are copied. |
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
