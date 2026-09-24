# ADR-0004: Local Scan-History Storage and Lifecycle

- Status: Proposed; owner-approved, independent architecture/security review pending; runtime authorization not granted
- Date: 2026-09-23
- Owners: `@geoffrey-xiao`
- Related issues: [#391 V05-003](https://github.com/geoffrey-xiao/deprail/issues/391); [#396 V05-002](https://github.com/geoffrey-xiao/deprail/issues/396)
- Related pull request: [#407](https://github.com/geoffrey-xiao/deprail/pull/407) (merged)
- Independent reviewer: Not assigned by owner direction; independent architecture/security review remains required before acceptance or runtime authorization.

## Context

v0.5 proposes retained local scan-operation history backed by SQLite. History identity must distinguish repeated operations that share the same `ScanReport.ScanID`; operation outcome must remain distinct from report completeness; raw scanner artifacts remain in the existing content-addressed store. The persisted report must retain its versioned meaning, provenance, and safe unknown fields. Existing CLI behavior must not depend on history unless a separately reviewed opt-in is used.

The existing v1alpha project and scan documents are the source snapshots. The project schema includes a repository-root field, while the local-data privacy contract prohibits persisting absolute host paths by default. The history projection therefore must store the repository root as `.` and use repository-relative workspace paths. Exact API response fields, numeric limits, SQLite driver/build strategy, and platform recovery evidence remain unresolved.

This proposal follows the v0.5 product, architecture, roadmap, release-plan, UX, API, failure/data, security, and compatibility contracts. It is not implementation authorization.

## Decision

Propose schema version 1 for review, with the following boundaries:

1. Place the per-user store at `os.UserConfigDir()/.deprail/history.sqlite3`; use no repository-relative fallback. Canonicalize the path and refuse history access if it resolves inside the scanned repository or crosses an unsafe symlink/reparse-point boundary. A missing database is an empty store for read-only history queries; reads do not create it. The first explicitly requested save may initialize it.
2. On POSIX, create the data directory with mode `0700` and database/backup files with mode `0600`. On Windows, enforce an owner-only ACL. If the required permission boundary cannot be established, fail the history operation without weakening permissions or writing elsewhere. Exact platform APIs and evidence remain implementation gates.
3. Record schema version in `PRAGMA user_version`, initially `1`. Keep it distinct from each stored project/report document's existing `schema_version` (`v1alpha`). Use sequential, forward-only migrations; refuse a newer or unsupported schema without modifying it.
4. Use one `history_entries` row per selected scan-operation occurrence. Generate a canonical UUIDv4 `history_entry_id` from a cryptographically secure random source once for that occurrence. `source_scan_id` is nullable when no report exists and non-unique when present. Retries of the same in-memory save reuse the same history ID; a new scan operation gets a new history ID even if its source scan ID repeats. Never use `INSERT OR REPLACE` or deduplicate by `source_scan_id`.
5. Store the recorded UTC time as integer Unix microseconds for deterministic ordering. Order list queries by `recorded_at_us DESC, history_entry_id DESC`. Keep operation outcome (`completed`, `failed`, `cancelled`) separate from optional report completeness (`complete`, `partial`, `failed`). A missing report is null/unavailable, never an empty successful report.
6. Store an optional validated v1alpha project snapshot and optional unchanged v1alpha scan-report JSON text. Normalize the project's `repository_root` to `.`; retain only safe repository-relative workspace paths and the non-path display label. Preserve safe unknown fields. Store no raw scanner bytes, source files, absolute host paths, tokens, environment values, or credential-bearing URLs. Keep `report_schema_version`, `source_scan_id`, operation diagnostics, and bounded list summaries as explicit columns; validate denormalized values against their snapshot before commit and on read.
7. Keep raw artifacts outside SQLite. Insert a `history_artifact_refs(history_entry_id, digest)` reference row for every validated report digest in the same transaction as its history entry. Verify content digests on evidence reads. Never cascade history deletion into artifact-file deletion; v0.5 proposes no per-entry delete/export operation or automatic artifact garbage collection.
8. Capture only an explicitly selected set of operations. Candidate CLI trigger: an additive `deprail scan --save-history` option; default CLI scans do not initialize or write history, and the browser remains read-only without scan-triggering or mutation routes. When explicitly selected, attempt to save terminal outcomes, including failed/cancelled operations when safe diagnostics exist. The exact flag and its CLI contract require separate compatibility approval.
9. Use SQLite WAL with one serialized writer and concurrent readers; set and verify WAL mode during initialization, enable foreign keys on every connection, and set/verify `synchronous=FULL` on each writer connection. Use a finite busy timeout. Each history row, project/report snapshot, and artifact-reference set commits in one transaction. A lock timeout or disk-full error rolls back the new entry and is reported as `HISTORY_WRITE_FAILED`; an existing scan result and its exit semantics are not rewritten. Candidate behavior for a successful scan whose requested history write fails is an unchanged report on stdout, a safe diagnostic on stderr, and additive CLI exit code `6`. When both scanning and persistence fail, preserve the scan failure and report the persistence failure separately. The exit-code precedence and timeout value require compatibility review and evidence; code `6` is not approved.
10. Before a future migration, make a consistent pre-migration snapshot with SQLite's online backup API into a unique sibling file with restrictive permissions; validate it before finalizing, never overwrite an existing backup, and do not copy only the main database file while WAL is active. Apply one ordered migration in a transaction and update `user_version` in that transaction. On failure, roll back and preserve both the original database and snapshot. No downgrade, automatic repair, or automatic restore is proposed. Restore remains a documented manual recovery action after preserving the damaged/current files.
11. Retain explicitly selected history without automatic eviction in v0.5. Do not expose entry deletion/export controls in the browser/API and do not remove raw artifacts. This has an unbounded disk-growth risk; the owner must explicitly accept that risk or select a measured cap before DoR. Numeric per-entry/report bounds and total-store policy remain unresolved until representative payload and platform evidence exists.

The proposed v1 relational shape and validation invariants are detailed in [`FAILURE-AND-DATA-CONTRACT.md`](../04-execution/deprail-v0.5.0-execution-package/requirements/FAILURE-AND-DATA-CONTRACT.md). The table layout is not a public API schema; API projection remains governed by the separately reviewed OpenAPI contract.

## Alternatives Considered

- Store one database per scanned repository: rejected as the default because it couples history to untrusted repository locations, risks storing data inside the scanned tree, and fragments one user's history.
- Persist every CLI scan automatically: rejected because it silently changes existing CLI behavior and conflicts with the v0.5 requirement that history be explicitly selected.
- Put raw artifact bytes in SQLite: rejected because the existing content-addressed artifact store owns raw evidence and SQLite should retain only verified references.
- Use `source_scan_id` as the history primary key: rejected because repeated operations can legitimately share a source scan ID.
- Normalize every report field into API-shaped relational tables now: deferred because the exact OpenAPI resource and pagination contract is not accepted; duplicating the report would add consistency and migration obligations before evidence justifies it.
- Delete unreferenced artifacts during history cleanup: rejected for v0.5 because artifacts may be shared and no approved reference-counted deletion workflow exists.

## Consequences

### Positive

- Repeated operations remain individually addressable without changing existing report identity.
- Scan outcome, report completeness, persistence outcome, and artifact integrity remain distinct.
- Existing report/project schemas and raw artifact ownership are preserved.
- Default CLI behavior remains independent of the optional local store.
- Migration failure has a defined rollback boundary and a separately validated recovery copy.

### Negative

- Explicit opt-in history is less discoverable and requires a reviewed CLI affordance.
- JSON snapshots plus API paging require payload limits and measurement; parsing large snapshots can consume memory.
- A per-user store and WAL require platform-specific permission, locking, and backup evidence.
- No automatic eviction or artifact collection can grow disk usage without bound.
- `--save-history` and candidate exit code `6` are additive CLI decisions that may require a compatibility revision.

## Compatibility and Migration

Schema version 1 is a new local format; no existing v0.1-v0.4 scan report is migrated into history automatically. Existing reports remain unchanged and retain their own schema version. The initial schema has no downgrade path. Future migrations must be sequential, deterministic, transactional, non-destructive, and preceded by a validated backup. Unsupported future versions remain untouched and produce an explicit history compatibility failure. Any change to report meaning, CLI exit codes, API fields, or artifact-store deletion semantics requires a separate reviewed compatibility decision.

## Validation

Before acceptance, provide:

- Schema validation for fresh v1 stores, required columns/constraints, report/project round trips, unknown safe fields, and repeated `source_scan_id` values.
- An idempotency scenario proving a same-operation retry cannot duplicate an entry and a later run with the same source ID remains distinct.
- Migration matrix coverage for absent/fresh store, current v1, supported prior version, interrupted migration, corrupt store, and unsupported future version; each failure must preserve the prior committed database and backup.
- Transaction scenarios for lock timeout, cancellation, disk full, malformed/oversized snapshots, and write interruption; no partial entry or false `saved` result.
- Artifact-reference round trips plus missing/digest-mismatch behavior without deleting or trusting invalid bytes.
- POSIX and Windows permission, canonical data-root, symlink/reparse-point, WAL/backup, and restore evidence.
- Representative payload-size and memory measurements before choosing numeric bounds or accepting unbounded retention risk.
- Owner and independent architecture/security review recorded separately. Until those reviews and the v0.5 DoR pass, this ADR remains Proposed and no SQLite driver, migration, or runtime store may be implemented.
