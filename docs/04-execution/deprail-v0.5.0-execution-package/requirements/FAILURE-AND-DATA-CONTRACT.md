# v0.5 Failure and Data Contract

**Status:** Draft proposal; SQLite schema, storage errors, and retention are not approved.

## 1. Data integrity invariants

- A stored scan refers to one stable scan identity and preserves the source report's explicit outcome, report/schema version, tool/scanner/database provenance, and artifact digest references required by existing contracts.
- Persistence must not convert a failed/partial/cancelled scan to complete, or a missing/invalid report to empty-success.
- The history index and referenced raw artifacts are separate resources. Dangling/missing artifacts and digest mismatches are explicit, scoped errors; do not silently delete metadata or invent evidence.
- Writes are atomic at the documented transaction boundary. Duplicate ingestion/idempotency and conflict behavior must be chosen before implementation.
- No source files, credentials, complete process environments, or unrelated repository contents are stored by default.
- Exact stored fields, payload-size thresholds, artifact retention/reference counting, and history capture trigger require an approved schema/ADR.

## 2. State vocabulary

Keep scan outcome distinct from history/API operation outcome:

- Scan: `complete`, `partial`, `failed`, or `cancelled` as current product contracts require.
- History operation: success, unavailable, invalid/incompatible, corrupt, integrity-failed, full/limit-exceeded, conflict, or cancelled/timeout as defined by final error contract.
- API request: success, client validation/not-found, unsupported version, rate/size bound, unavailable, internal failure, or cancellation/timeout.

These labels are semantic states, not finalized wire codes. Stable error codes and HTTP status mapping must be cross-walked with `requirements/ERROR-MODEL.md` before implementation.

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
| API timeout/cancel | Explicit operation outcome | No incomplete response described as complete |
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
