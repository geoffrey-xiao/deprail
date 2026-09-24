# ADR-0004: Local Scan-History Storage and Lifecycle

- Status: Proposed; owner accepted the high-level v0.5 direction on 2026-09-24; exact projection/API schema, token bootstrap, numeric bounds, driver/lifecycle evidence, and runtime authorization remain pending
- Date: 2026-09-23
- Owners: `@geoffrey-xiao`
- Related issues: [#391 V05-003](https://github.com/geoffrey-xiao/deprail/issues/391); [#396 V05-002](https://github.com/geoffrey-xiao/deprail/issues/396)
- Related pull request: [#407](https://github.com/geoffrey-xiao/deprail/pull/407) (merged)
- Owner reviewer: `@geoffrey-xiao`; external architecture/security review is optional under [ADR-0005](ADR-0005-solo-owner-review-policy.md).

Governance note: [ADR-0005](ADR-0005-solo-owner-review-policy.md) supersedes all independent-review and dual-approval requirements in this proposed ADR. Statements below about unavailable or pending external review preserve historical context; current technical acceptance requires the owner's decision and the listed evidence only.

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
- Owner security/architecture review and technical evidence are recorded. External review is optional under ADR-0005. Until owner decisions and the v0.5 technical DoR pass, this ADR remains Proposed and no SQLite driver, migration, or runtime store may be implemented.

## Post-proposal validation finding: source report is not the proposed stored document

The owner-approved proposal in [PR #407](https://github.com/geoffrey-xiao/deprail/pull/407) remains Proposed in its revised form; the original unchanged, schema-validated `report_json` and no-absolute-host-path requirements cannot both be implemented against the current CLI report:

- [`app.ScanReport`](../../internal/app/scan.go) contains `repository_identity.root`, `[]adapter.Finding`, and `[]string` errors. `Scan` constructs `repository_identity.root` from the canonical **absolute** `scanRoot`; normalizing only the project snapshot's `repository_root` to `.` does not remove that value from the report.
- [`adapter.Finding`](../../internal/adapter/contract.go) does not have the required `schema_version`, `document_type`, `stable_key`, `vulnerability`, `evidence`, or `state` fields of the [`v1alpha` finding schema](../../schemas/v1alpha/deprail.schema.json). Current nonempty `ScanReport.Findings` therefore cannot validate as an unchanged v1alpha scan document. `ScanReport.Errors` is a string list, while the v1alpha scan schema requires diagnostic objects.
- The existing CLI/report format and remediation inputs must not be silently rewritten to make storage pass validation. Storing the unchanged bytes would persist an absolute host path and violate the proposed data-minimization boundary; calling a redacted copy “unchanged v1alpha report JSON” would be false.

**Pending owner decision:** The conservative candidate is a separately versioned, allowlisted **history projection** derived from the report, retaining source scan identity, report completeness, safe provenance and validated artifact digests without persisting absolute paths. It would not be an unchanged `ScanReport` or a v1alpha scan document; the proposed `report_json` column, snapshot-validation and API detail contracts must be revised explicitly, with compatibility, unknown-safe-field and missing-artifact behavior evidenced. The alternative of changing CLI/report serialization is a separate compatibility/schema change; storing the raw absolute-root report requires an explicit privacy decision.

**Evidence needed before acceptance:** Validate a real `app.ScanReport` with nonempty findings, a report with nonempty errors, and an absolute source root against the proposed persistence schema and the eventual OpenAPI projection. Assert no absolute host path or unreviewed field is persisted/exposed, no original CLI/report or remediation input changes silently, and failed projection/serialization leaves the scan result and prior database unchanged. Small zero-finding schema examples alone do not exercise this boundary.

## Owner-selected direction: independent history projection (not accepted)

The owner selected a separate, versioned history projection rather than changing existing `ScanReport`/CLI JSON or reducing v0.5 scan-detail scope. This is a design direction, **not** owner acceptance of the exact storage schema, API fields, error codes, retention, numeric limits, or implementation. Original decision item 6's unchanged v1alpha `report_json` assumption is superseded as a *candidate payload design* by this addendum; all other proposed boundaries still require owner acceptance and technical evidence. External review is optional under ADR-0005. No database has been shipped for v0.5, so this is a revision of a proposed initial schema, not an automatic migration of user data.

The preceding “pending owner decision” finding records the state before this selection. Snapshot-specific wording elsewhere in the original proposal (item 3's stored-document versions, item 9's two snapshots, the repository-root-as-`.` note, and Validation's project/report round trips) is historical, not an instruction to persist original documents alongside `history-v1`. The revised candidate retains SQLite versioning, atomicity and path safety while validating one projected payload and its artifact references. No root field is stored; `.` denotes only an actual repository-relative workspace path where applicable.

- A projection with its own schema version (candidate `history-v1`, separate from SQLite `PRAGMA user_version=1`, source `ScanReport.SchemaVersion`, and `/api/v1`) is built from trusted, validated `app.ScanReport` and `ProjectGraph` inputs. Preserve the original CLI/report bytes and remediation contracts. Persist no unchanged project/report JSON, absolute repository root, raw scanner bytes, raw error strings, credential URLs, unreviewed unknown fields, or original source-file contents. Accept only reviewed, bounded fields; a required value that cannot be safely classified rejects the selected save without modifying the scan result or prior database.
- Retain the unique history occurrence ID, recorded UTC time, operation outcome, nullable source scan ID and report completeness as separate values. The projected report may include its source schema version, opaque repository-state identity, safe workspace summaries, normalized finding summaries, and validated artifact digests. Tool/database provenance is present only if supplied by an approved source; never fabricate it or infer completeness from an empty list. API responses are reviewed projections of these values, not serialized `ScanReport` or the entire stored JSON blob.
- Map actual `adapter.Finding` fields deliberately: `TargetID` is a vulnerability ID, **not** the finding key. Reuse [`normalize.StableFindingKey`](../../internal/normalize/key.go) with workspace ID, PURL/component identity, version, vulnerability ID, and sorted aliases as in [`baseline.ConvertScan`](../../internal/baseline/convert.go); reject missing/ambiguous identity rather than invent a key. Preserve report/finding order-independence, status, source identity, and safe provenance. A missing/mismatched artifact remains an integrity state; do not synthesize evidence bytes or a clean result.
Project only bounded, redacted, typed application diagnostics from trusted structured context. Current `ScanReport.Errors` strings are not a trusted API/storage diagnostic schema and may contain local paths or secrets. If the selected capture cannot represent a partial/failed/cancelled operation truthfully and safely, fail that history save with a typed persistence diagnostic; do not store an empty successful record or rewrite the scan outcome. The exact error-code distinction remains an owner decision in the [error model](../04-execution/deprail-v0.5.0-execution-package/requirements/ERROR-MODEL.md).

`app.Scan` currently consumes `ProjectGraph` internally but returns only `ScanReport`; an external history adapter cannot recover the original discovered workspace set from that report. A future capture path must obtain validated graph context from the **same operation** through an owner-accepted application contract, not re-discover later and silently combine different repository states. If that context is unavailable, projected workspaces/count are null rather than a guessed list. `Workspace.Completeness` describes discovery, not per-workspace scanner success; existing completion events do not provide a reliable replacement for a typed workspace scan result. Finding paths are separately validated relative to the repository, and the API must not label discovery completeness as scanner completeness.

The revised proposed column/JSON shape and its read/write consistency checks are in the [failure/data contract](../04-execution/deprail-v0.5.0-execution-package/requirements/FAILURE-AND-DATA-CONTRACT.md#4-proposed-persisted-schema-v1-unapproved). The owner must accept that shape, the privacy/unknown-field policy, real partial/failure diagnostics, payload bounds, and compatibility evidence before any implementation issue or OpenAPI schema is frozen; external review is optional under ADR-0005. The existing small zero-finding examples and local reports with up to 14 findings are not sufficient to establish numeric capacity limits.

## Owner-accepted v0.5 design direction (2026-09-24; detailed contract pending)

The owner accepted the high-level v0.5 directions in the [execution-package decision table](../04-execution/deprail-v0.5.0-execution-package/README.md#owner-accepted-technical-directions-detailed-dor-pending). This acceptance establishes a planning baseline only. It does not accept the exact API/storage schemas, error mappings, dependencies, numeric bounds, token transfer mechanism, runtime behavior, or implementation authorization.

- Keep v0.5 local, read-only, and limited to scan history/console. History capture is explicit opt-in through the proposed `deprail scan --save-history`; ordinary scans remain unchanged by default. A selected save failure must not rewrite the scan outcome; the additive flag and exit-code behavior remain to be recorded in the current contract.
- Preserve existing `ScanReport` and CLI JSON. Persist a separately versioned allowlisted `history-v1` projection, use same-operation graph context or null rather than rediscovery, and reject unsafe projections without modifying the scan or prior database.
- Keep the local API loopback-only, read-only, same-origin, without CORS, and bounded by a 1 MiB serialized UTF-8 response ceiling. Whole-record byte-aware paging is the direction; the fragment-token bootstrap and its transient browser/OS exposure remain unaccepted pending a safer, evidenced transfer choice.
- Keep an owner-only per-user SQLite root, serialized transactional writes, WAL/`synchronous=FULL`, validated online-backup snapshots, and no automatic repair, deletion, eviction, or artifact cleanup. Explicit capture and documented manual cleanup are preferred for the preview.
- Retain the React/TypeScript/Vite baseline without adding UI component/icon/font/CDN dependencies; use project-owned semantic UI, target WCAG 2.2 AA, and use Chrome/NVDA on Windows, Safari/VoiceOver on macOS, and Firefox/Orca on Linux as the proposed support matrix.
- For #356, assign rollback/recovery to `@geoffrey-xiao`; defer full Python/Java remediation beyond v0.5 while preserving representative history coverage; require Linux/macOS/Windows representative smoke before stable release; record SBOM, signature, and provenance as supplied, unavailable, or deferred without inference.

The remaining DoR is not complete. Exact projection fields and diagnostics, CLI compatibility/error mapping, API authentication and resource bounds, SQLite driver/timeout/store limits, dependency/build evidence, release-plan carry-forward details, and design-stage acceptance evidence must be recorded before implementation issue creation. Runtime, hostile-input, browser/accessibility, recovery, and cross-platform results remain implementation/release evidence, not passes inferred from this direction. ADR-0004 remains Proposed and authorizes no runtime work.
