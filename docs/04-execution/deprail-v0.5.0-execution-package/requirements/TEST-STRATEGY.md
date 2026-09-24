# v0.5 Test Strategy

**Status:** Draft; exact commands/tool versions are assigned in the implementation execution package after design approval.

## 1. Principles

Test consumer-observable behavior across service, storage, API, browser, packaging, and platforms. Keep fixtures offline, deterministic, minimal, and license-safe. A passing unit suite alone does not satisfy release acceptance. Preserve CLI tests and add no UI/API framework-specific assertions where externally observable contracts suffice.

## 2. Domain and history contract tests

- Distinct runs with the same source `ScanReport.ScanID` have unique history-entry IDs and remain independently addressable; ordering is deterministic.
- Report completeness (`complete|partial|failed`) and operation outcome (`completed|failed|cancelled`) survive serialization/storage roundtrip as separate fields; cancelled operations never acquire a fabricated report status.
- Empty history differs from missing/corrupt/unavailable store.
- Finding keys, provenance, artifact digests, workspace identity, and diagnostics survive permitted data transforms.
- Duplicate ingestion/idempotency and conflict behavior once specified.
- Missing artifacts/digest mismatch are explicit and do not become empty evidence.

## 3. SQLite integration tests

- Fresh DB initialization, supported migration from each prior schema, latest-schema reopen.
- Migration interruption/failure, old DB preservation, backup/restore according to approved plan.
- Transaction rollback after mid-write failure and process-like interruption.
- Concurrent readers/writer, lock/busy timeout, cancellation, disk-full/quota simulation, permission denial, corrupt file/page, unsupported future schema.
- Retention/deletion and artifact reference accounting once defined.
- Atomic and restrictive filesystem behavior on Linux/macOS/Windows.

Use temporary directories/DBs; never require network, real user data, or host-global paths.

## 4. API contract and adversarial tests

- Validate operations and responses against reviewed OpenAPI 3.1 schema and examples.
- Verify bounded pages, stable cursors/order, invalid/unknown filters, unknown IDs, unsupported versions, malformed/oversized bodies, wrong content type/method.
- Assert HTTP status/error envelopes for store unavailable, corrupt data, artifact missing/mismatch, timeout, and read-request cancellation; separately verify persisted scan-operation cancellation does not mutate `ScanReport.Status`.
- Security cases: traversal and encoded separators, unusual Unicode, symlink/path boundary if any filesystem resolution exists, SQL metacharacters, hostile Host/Origin, DNS-rebinding patterns, forbidden CORS, credential-bearing headers, and oversized result/query workloads.
- Ensure APIs do not return raw credentials, full environments, unsafe absolute paths, or source file content.

## 5. Browser UX and accessibility verification

Exercise the actual packaged application/browser surface for history, scan detail, report completeness, completed/failed/cancelled operation outcomes, no-findings, empty history, loading, storage/API errors, incompatible schema, missing artifact, not found, repeated source scan IDs with distinct history entries, and narrow viewport. Verify shadcn-inspired tokens/components match approved design without asserting CSS implementation details. Test keyboard-only navigation, focus movement/restoration, accessible names/landmarks/table semantics/status announcements, contrast, reduced motion, and hostile repository labels. Capture browser/OS/version and evidence; no simulated component test substitutes for the actual-surface smoke.

## 6. Cross-platform and compatibility

Run unit/contract/integration and packaging checks on Linux, macOS, and Windows. Verify loopback-only listener, origin behavior, paths/permissions, SQLite locking/shutdown, embedded assets, and schema meaning on each supported OS. Compare semantic records and error categories; do not pin platform-specific path strings or incidental messages.

## 7. CLI and end-to-end checks

- Existing CLI commands and JSON/exit contracts pass regression coverage with console disabled and history unavailable according to accepted compatibility decisions.
- End-to-end representative local workflow: produce/use saved history entries from a fixed valid report fixture, including repeated source scan IDs and a cancelled operation with a returned report, open console, browse page/detail/findings, and confirm operation outcome, report completeness, and provenance remain distinct.
- Compare repository tree before/after browser/API workflows; browser must not mutate source repository or package files.
- Exercise no-data startup and storage failure without installing scanners or fetching network resources.

The exact method for creating history fixtures and whether CLI scan automatically records history is a design decision; do not invent a behavior in tests before it is approved.

## 8. Release manual evidence

Record reviewed binary/version/commit, OS/architecture/browser/version, exact commands, exit codes, listener address, history DB location/state, API results, screenshot/accessibility observations, tree comparison, artifact/checksum/SBOM/signature/provenance status, independent reviewer, owner decision, and remaining gaps. Follow `docs/RELEASE-CHECKLIST.md`; tests alone do not approve a release.

## 9. Requirement and threat evidence matrix

The IDs below are planned verification scenarios, not implemented tests or completed evidence. The v0.5.0 release-evidence record must attach exact command/result, commit and binary identity, fixture/database digests, exit/status, OS/architecture, and relevant tool versions. Browser scenarios also record browser and assistive-technology versions. No supported browser list or SQLite driver is selected by this matrix.

| ID | Consumer-observable scenario | Required future evidence | Required matrix |
| --- | --- | --- | --- |
| FR-501 | Save/read two operation entries with the same source `ScanReport.ScanID`; require distinct history IDs, unchanged source report identity, correct `history-v1`/source schema metadata, and no invented provenance. Retry the same operation using its ID; it must not create a duplicate. | Storage integration output, fixture and DB digests, projected round-trip values. | SQLite behavior on Linux, macOS, Windows. |
| FR-502 | Page history with a stable tie-breaker/cursor while another entry is inserted; distinguish confirmed empty page from loading, invalid cursor, store unavailable, timeout/cancel, and API failure. | API contract results with sanitized status/body/cursor; actual-browser empty/error-state capture. | Linux, macOS, Windows; every approved browser/OS combination. |
| FR-503 | Load projected detail with a real nonempty report, without a trustworthy report, missing artifact, and digest mismatch. A null report is unavailable, not zero findings; a valid zero-finding complete report remains distinct. Finding IDs follow `normalize.StableFindingKey`, not adapter `TargetID`, and API response allowlists preserve only trustworthy metadata. A missing same-operation `ProjectGraph` yields unavailable workspaces/count (not inferred from findings or a later rediscovery). | Storage/API result and fixture digests, finding-key comparison, browser detail-state capture. | Linux, macOS, Windows; every approved browser/OS combination. |
| FR-504 | Exercise completed+complete, completed+failed, cancelled+complete, partial with zero findings, and failed-without-report. Verify both state axes survive and no incomplete result appears clean; discovery-only workspace completeness must not be presented as per-workspace scan success when scan events/status are insufficient. | Serialized round-trip comparison, API responses, and actual-browser state/announcement record. | Linux, macOS, Windows; every approved browser/OS combination. |
| FR-505 | Run existing CLI commands with web disabled, history absent, and history unavailable; default stdout/stderr, JSON purity, and exit codes stay unchanged, including a real nonempty/failed scan with source `repository_identity.root` and errors in CLI JSON but never history. A listener startup failure reports safely as `API_LISTENER_UNAVAILABLE` without breaking CLI use. Test explicit history-save behavior only after its CLI decision is approved. | CLI regression output, startup diagnostic, and repository-tree comparison. | Linux, macOS, Windows. |
| STORE-01 | After explicit history capture is approved and selected, force projection validation failure (absolute root, raw/unclassified error or field, missing/duplicate stable finding identity) and transaction failure separately; require typed persistence failure, no “saved” claim/partial row, and unchanged scan report/outcome. Additive exit code `6` is tested only if separately approved. | Projection/transaction rollback, safe stderr diagnostic, stdout and exit-status record, before/after DB digest. | Linux, macOS, Windows. |
| FR-506 | Run equivalent domain scenarios through CLI and HTTP adapters over the same application services; observable report/status/provenance meaning matches. | Service/adapter contract results and sanitized request/response pair. | Linux, macOS, Windows. |
| FR-507 | Send invalid IDs/filters, unknown routes/versions, encoded separators, SQL metacharacters, unsupported methods/content types, and hostile Host/Origin values; reject safely with no filesystem escape, query injection, mutation, or false-empty result. | Sanitized adversarial request/result report and DB/tree before-after digests. | Linux, macOS, Windows; browser-origin cases on every approved browser. |
| FR-508 | Exercise fresh/current/future/corrupt stores and interrupted migration/write; test each explicitly supported prior schema (if any) and record none if no prior version is supported. Verify old committed data and backup survive and retention behavior matches the accepted decision. | Migration/rollback matrix, backup and DB digests, recovery procedure record. | SQLite/filesystem behavior on Linux, macOS, Windows. |
| FR-509 | Package UI assets with matching API version; missing assets or version mismatch produces visible failure and never serves arbitrary filesystem paths. | Packaged binary smoke log, asset/API version identity, browser capture. | Linux, macOS, Windows; every approved browser/OS combination. |
| FR-510 | Exercise history/detail using keyboard and assistive technology; verify focus return, accessible names/status announcements, contrast, narrow layout, and reduced motion. | Actual-surface accessibility checklist, screenshots, browser and assistive-technology versions. | Every approved browser/OS/assistive-technology combination. |
| FR-511 | Inspect routes, methods, navigation, and browser controls; scan, delete, export, remediation, remote publishing, team, and agent-write actions are absent and cannot be invoked. | Route/method inventory, browser workflow capture, repository-tree comparison. | Linux, macOS, Windows; every approved browser/OS combination. |
| SEC-01 | Hostile web origin, DNS rebinding pattern, and unexpected Host/Origin cannot read history; no permissive CORS. | Sanitized request/response trace and listener/origin test result. | Linux, macOS, Windows; every approved browser. |
| SEC-02 | Verify per-user data-root ownership/modes or ACL; unauthorized local user cannot read history; no unsafe fallback. | POSIX mode/Windows ACL evidence and data-root path record. | Linux, macOS, Windows. |
| SEC-03 | Traversal, encoded separators, symlink/reparse escape, and arbitrary asset path requests fail without outside-root access. | Hostile path corpus/result and filesystem before-after record. | Linux, macOS, Windows. |
| SEC-04 | SQL metacharacters and unknown sort/filter input cannot change query meaning or trigger unbounded work. | Adversarial query results and unchanged DB digest. | Linux, macOS, Windows. |
| SEC-05 | Oversized bodies (only if a body route is approved), oversized responses, concurrent requests, and slow requests hit explicit bounds without truncation or resource exhaustion. | Limit-boundary results and bounded memory/runtime profile. | Linux, macOS, Windows; browser cases on every approved browser. |
| SEC-06 | Hostile repository labels/findings render as text; scripts, unsafe URLs, and markup do not execute. | Actual-browser hostile-content capture and console/error record. | Every approved browser/OS combination. |
| SEC-07 | Use actual nonempty and failed `ScanReport` values plus hostile project diagnostics/path/credential fixtures; persisted `history-v1`, API responses, and logs must not leak raw root/errors, tokens, credential URLs, full environments, raw source, SQL/stack traces, or sensitive absolute host paths. Unknown fields must be classified or rejected, never silently promoted to safe detail. | Source/projection/response field comparison, redaction corpus and sanitized response/log evidence. | Linux, macOS, Windows; every approved browser for rendered output. |
| SEC-08 | Corrupt pages/rows, unsupported future DB/projection schema, and row/projection mismatch produce explicit errors; bytes are preserved and no automatic repair/replacement occurs. | Corruption/future-version/mismatch matrix and pre/post DB digests. | Linux, macOS, Windows. |
| SEC-09 | Missing or digest-mismatched artifact is disclosed, not trusted, substituted, or deleted. | Artifact fixture digests and integrity-result record. | Linux, macOS, Windows. |
| SEC-10 | UI/API contract mismatch is visible; no stale or empty-success state is substituted. | Packaged version-skew result and browser capture. | Linux, macOS, Windows; every approved browser/OS combination. |
| SEC-11 | Inspect the real listener address; default listener is loopback-only and never silently falls back to wildcard/LAN/public binding. | Bound-address and startup-failure evidence. | Linux, macOS, Windows. |
| SEC-12 | Encoded/static asset routes serve only the embedded allowlist and disclose no host files. | Route corpus, response inventory, filesystem canary result. | Linux, macOS, Windows; every approved browser. |
| SEC-13 | Pinned toolchain/dependencies build reproducibly; review licenses, lockfile, and software bill of materials before selecting packages. | Lockfile/build identity, license review, SBOM, and CI results. | Linux, macOS, Windows build jobs. |

All matrix rows are release gates, not current pass claims. Freeze exact supported browser/assistive-technology versions and evidence storage location in the approved execution package before implementation; do not infer them from CI runners or local browser availability.
