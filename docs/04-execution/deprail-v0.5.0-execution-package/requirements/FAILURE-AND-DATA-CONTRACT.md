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

| Condition | Required observable behavior | Data-safety requirement |
| --- | --- | --- |
| Empty history | Successful empty collection and empty-state UI | Do not synthesize a scan record |
| Missing data directory/database | Follow reviewed first-run initialization policy | No silent scan failure or unintended creation beyond approved policy |
| Permission denied | Explicit storage unavailable error | Do not weaken permissions or write elsewhere silently |
| Database locked/concurrent writer | Bounded wait or explicit conflict/unavailable result, per decision | No partial transaction visible |
| Disk full/quota | Explicit persistence failure; separate from scan outcome | Roll back incomplete DB transaction; do not report record persisted |
| Interrupted write/process crash | Recovery follows SQLite transaction/atomicity decision | Prior committed history remains readable; no half-record |
| Unsupported future schema | Refuse unsafe writes with upgrade guidance | Never downgrade/delete/overwrite automatically |
| Migration error/interruption | Explicit migration failure and recovery instruction | Preserve old database/backup; migration is atomic or recoverable |
| Corrupt DB/page | Explicit integrity/storage error | No automatic destructive reset; preserve evidence for recovery |
| Missing referenced artifact | Scan metadata may be shown if valid; integrity warning/error for evidence | Never substitute empty artifact or claim provenance verified |
| Digest mismatch | Integrity failure at scoped read | No trust in corrupted bytes; no mutation of source artifact |
| Malformed stored report | Explicit invalid-record result scoped to scan | Do not expose malformed content as valid report |
| Request oversized/malformed | Client error; bounded work | Reject before expensive parsing/DB query |
| Cancelled scan operation | History entry records `operationOutcome=cancelled`; any returned report keeps its existing status | Never present the report's retained pre-cancellation completeness as proof of a completed run |
| API timeout/request cancellation | Explicit request outcome; no fabricated complete response | Must not rewrite a stored history entry or `ScanReport.Status` |
| UI assets missing/version mismatch | Explicit console unavailable/error | CLI remains usable |
| Shutdown during query/write | Cancel/drain according to accepted transaction model | Durable commit or rollback; no ambiguous “saved” result |

## 4. Migration and recovery (decision required)

The persistence ADR must define schema version storage, supported upgrade path, migration transaction, pre-migration backup/snapshot, free-space check, interrupted-migration recovery, backup naming/permissions, restore instructions, and any downgrade prohibition. Migration must never replace the only prior copy before successful validation. How backup retention and restore are exposed is a human decision; this draft does not assume a backup command.

## 5. Retention and deletion (decision required)

Specify default retention behavior, whether user-controlled deletion/export is included, how report metadata and content-addressed artifacts remain consistent, and how size limits are surfaced. Do not invent automatic eviction or delete shared artifacts without proving they are unreferenced and confirming the product decision. If no retention limit is selected, the disk-growth risk must be explicitly accepted and documented.

## 6. Transaction and concurrency contract (decision required)

Select writer model, reader concurrency, SQLite journal mode, busy timeout, transaction granularity, and read snapshot semantics. A scan operation that succeeds while history persistence fails must preserve the scan's outcome and communicate that persistence failed. If capture is automatic, the CLI contract must identify whether persistence failure changes command exit status; this decision cannot be left to incidental adapter behavior.

## 7. Redaction and privacy

API and UI expose only approved fields. Do not persist environment values, tokens, credential-bearing URLs, absolute host user paths, raw source files, or package-manager secrets. Repository labels/paths are untrusted and escaped in HTML; diagnostics are bounded and redacted. Local does not mean trusted: another local process and a hostile browser origin may reach the service.
