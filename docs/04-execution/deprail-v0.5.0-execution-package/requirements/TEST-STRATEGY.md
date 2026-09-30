# v0.5 Test Strategy

**Status:** Current proposed acceptance matrix for owner review. Planning/schema validation precedes implementation; actual runtime and packaged-surface evidence follows implementation. Commands/toolchain candidates and support policy are specified in the compatibility matrix; no runtime pass is claimed.

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

History capture is explicitly selected through `deprail scan --save-history`; ordinary scans do not initialize/write history. Use the failure/data schema examples for offline contract scenarios and real representative JS/Python/Java repositories for packaged release smoke. Browser/API calls never trigger scans or writes.

## 8. Release manual evidence

Record reviewed binary/version/commit, OS/architecture/browser/version, exact commands, exit codes, listener address, history DB location/state, API results, screenshot/accessibility observations, tree comparison, artifact/checksum/SBOM/signature/provenance status, optional external reviewer (if invited), owner decision, and remaining gaps. Follow `docs/RELEASE-CHECKLIST.md`; tests alone do not approve a release.

## 9. Requirement and threat evidence matrix

The IDs below are planned verification scenarios, not implemented tests or completed evidence. The v0.5.0 release-evidence record must attach exact command/result, commit and binary identity, fixture/database digests, exit/status, OS/architecture, and relevant tool versions. Browser scenarios also record browser and assistive-technology versions. Use the concrete dependency candidates and supported platform/browser policy in the compatibility matrix; selection is not a platform pass.

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
| SEC-05 | Prohibited request bodies, request-target/header limits, whole-record byte-limited pages, oversized indivisible records/details, concurrent and slow requests hit explicit bounds without truncation or resource exhaustion. | Fully serialize JSON as UTF-8 and enforce the 1 MiB cap including escaping/envelope/cursor. Exercise maximum Unicode, exact cap/one-byte-over, short pages with every record returned once, empty terminal page, eight-request saturation and the 10-second deadline. A record/detail that cannot fit returns bounded `500 API_RESPONSE_TOO_LARGE` before any success body; never emit a zero-item nonterminal page. | Linux, macOS, Windows; browser cases on every supported browser. |
| SEC-06 | Hostile repository labels/findings render as text; scripts, unsafe URLs, and markup do not execute. | Actual-browser hostile-content capture and console/error record. | Every approved browser/OS combination. |
| SEC-07 | Use actual nonempty and failed `ScanReport` values plus hostile project diagnostics/path/credential fixtures; persisted `history-v1`, API responses, and logs must not leak raw root/errors, tokens, credential URLs, full environments, raw source, SQL/stack traces, or sensitive absolute host paths. Unknown fields must be classified or rejected, never silently promoted to safe detail. | Source/projection/response field comparison, redaction corpus and sanitized response/log evidence. | Linux, macOS, Windows; every approved browser for rendered output. |
| SEC-08 | Corrupt pages/rows, unsupported future DB/projection schema, and row/projection mismatch produce explicit errors; bytes are preserved and no automatic repair/replacement occurs. | Corruption/future-version/mismatch matrix and pre/post DB digests. | Linux, macOS, Windows. |
| SEC-09 | Missing or digest-mismatched artifact is disclosed, not trusted, substituted, or deleted. | Artifact fixture digests and integrity-result record. | Linux, macOS, Windows. |
| SEC-10 | UI/API contract mismatch is visible; no stale or empty-success state is substituted. | Packaged version-skew result and browser capture. | Linux, macOS, Windows; every approved browser/OS combination. |
| SEC-11 | Inspect the real listener address; default listener is loopback-only and never silently falls back to wildcard/LAN/public binding. | Bound-address and startup-failure evidence. | Linux, macOS, Windows. |
| SEC-12 | Encoded/static asset routes serve only the embedded allowlist and disclose no host files. | Route corpus, response inventory, filesystem canary result. | Linux, macOS, Windows; every approved browser. |
| SEC-13 | Pinned toolchain/dependencies build reproducibly; review licenses, lockfile, and software bill of materials before selecting packages. | Lockfile/build identity, license review, SBOM, and CI results. | Linux, macOS, Windows build jobs. |

All matrix rows are delivery/release gates, not current pass claims. Apply the supported browser/assistive-technology version policy in the compatibility matrix and record the exact versions used per evidence run; never infer runtime support from CI runner availability. Public release assets exclude raw reports/logs and credentials. Retain sanitized evidence in the version-specific release-evidence record with links to controlled artifacts.

## 10. Detailed boundary and recovery scenarios

These acceptance cases make the engineering admission limits observable; they do not require adding runtime code during planning.

- Entry budget: exact 16 MiB serialized projection versus one byte over; exact and one-over 10,000 findings/1,000 workspaces/128 diagnostics/4,096 digests. All-or-nothing refusal, no silently dropped aliases, diagnostics or digests.
- Store admission: exact/one-over 1,000 entries and 256 MiB summed projection bytes; race two captures near the limit. The writer transaction admits only permitted totals. Indexes/WAL/backups/artifacts are not falsely included in this logical quota; forced disk-full rolls back independently.
- Identity and repeated save: retry a known same-operation ID with identical content versus conflicting content; verify no duplicate or replacement. A later scan with the same source scan ID gets another history occurrence.
- Safe partial/failed/cancelled capture: exercise typed errors before they become CLI strings, early failure without graph/report, graph-only context, partial source report and cancellation with unchanged pre-cancellation report status. Unknown diagnostic code, unsupported source version, ambiguous finding identity and an unsafe required value refuse capture rather than fabricate safe detail.
- Cancellation: stop scanner work when canceled; selected capture uses only the independently bounded terminal-save context specified by failure/data, preserves the original cancellation exit and never resumes scanning. A second termination or save deadline stops capture without partial rows.
- Browser bootstrap: initial top-level navigation remains possible under the static-route Fetch Metadata rules; API requires bearer plus strict origin/host checks. Clear fragment before rendering/request; reload is an authentication-recovery state, not empty history. No token in printed URL, diagnostics, logs, requests targets, storage or referrer. Browser/OS launch exposure remains a reviewed residual risk, not a guarantee of same-user isolation.
- Artifact provenance: without a trusted current artifact resolver, return `unavailable`, not `verified` or `missing`. Exercise trusted present/missing/mismatched bytes separately; no absolute repository root or artifact path is persisted/reconstructed from labels.
- Recovery rehearsal: stop console/writers, preserve the original DB/WAL/SHM plus external artifacts and candidate backups, validate a separate backup copy, open a separately restored copy with the matching binary, and perform only explicitly owner-controlled replacement. Unsupported schema, failed backup validation or permissions leaves the only copy untouched.
- Delivery identity: package known UI routes and embedded assets with the matching API/schema contract. Record clean-lockfile build, license/SBOM review and candidate vulnerability results; budget failures are explicit release blockers or owner-dispositioned preview gaps, never unsupported success claims.

## 11. Planning contract smoke evidence (PR #417)

Executed on macOS/darwin arm64 with Go `1.27.1`, host Node `26.7.0`/npm `11.19.0`; proposed frontend build pins were not installed. Python validation used the isolated environment `/tmp/deprail-392-contract-review` with PyYAML `6.0.3`, jsonschema `4.26.0` and openapi-spec-validator `0.9.0`.

- Unique-key YAML parsing and OpenAPI 3.1 validation passed; all **62 operation-response examples** validated against their schemas. There remain exactly five authenticated GET API operations; static routes are separately declared.
- Draft 2020-12 schema validation and **three history-v1 examples** passed. Throwaway projection smoke validated **eight API resources**, including null report/context versus trustworthy empty report and nonempty findings. Provenance is explicit nullable typed metadata, not fabricated scanner/database versions.
- **29 hostile/boundary mutations** plus three extra trailing-control/dot-segment path cases were rejected; **17 diagnostic code/message pairs** accepted and mismatched messages rejected. Root `.`, dotfiles and Unicode remain valid. A regex-alternative anchoring defect exposed by `../x` and trailing newline cases was corrected before publication.
- The nonempty fixture key was corrected after exercising actual `normalize.StableFindingKey` with `go run ./.deprail-contract-review`; it returned `629c2d635fcf2c1b477bf527cde2fb15d587aa86abab81146c6aded0e2d37feb`. The throwaway program was removed. This verifies the existing function, not implemented history capture.
- Maximum legal Unicode finding fields in a 25-record API envelope measured **1,540,420 UTF-8 bytes** or **4,598,820 ASCII-escaped bytes**, including a conservative 512-byte cursor. Whole-record simulation returned 17/8 records at 1,047,676/492,829 bytes for UTF-8, or five pages of five records at 920,240 bytes (last 919,730) for escaped JSON. All simulated envelopes validated, remained ≤1 MiB and retained every record in order. This is serialization evidence, not HTTP runtime/paging performance.
- Proposed cursor packing roundtrips produced history/workspace/finding decoded lengths **26/54/50 bytes**, unpadded base64url lengths **35/72/67 ASCII bytes**. No workspace label/path is encoded. Production cursor rejection/authentication remains a future acceptance scenario.
- WCAG luminance calculations for the six text/accent tokens against white and card surfaces yielded minimum ratios **6.18–15.55:1**; control border `#64748B` gives **4.55:1**. Decorative `#CBD5E1` is only **1.48:1** against white and is explicitly forbidden as the sole meaningful control boundary. Actual rendered states/browser accessibility remain untested.
- Script-free temporary dependency review: `npm install --package-lock-only --ignore-scripts --no-fund --no-audit && npm audit --json`. Candidate Vite `7.1.4` failed (one vulnerable direct dependency, high severity); amended `7.3.6` set returned audit exit **0**, **zero reported vulnerabilities**, **121 resolved dependency entries**. Lockfile SHA-256: `71a3114f95631e20a348c00216546797559f9f84fc8ca97389dcfac535a2304c`. No package scripts ran or project dependencies were installed. License metadata includes MIT/Apache-2.0/ISC/BSD-3-Clause and CC-BY-4.0; complete license/provenance/SBOM and pinned-toolchain build review remains required.
- **186 local document file/heading targets** passed. Reviewed artifact SHA-256: history schema `52bca367f8e91f9947bf05bcd567b977fc4c47fc16dce85e493d9c186f6fdaae`; OpenAPI `3cf30ba33ba382e7ec62bbc0a8fcbd0f09b54cb09e01ff189d5912479ac4df9f`.
- `make verify` exited **2**: generation and vet passed, but existing `internal/remediation/verification.TestRunExecutesSelectedCommandsInWorkspace` failed with `SCANNER_TIMEOUT: process exceeded configured deadline`. No assertion/deadline was changed and no passing rerun is substituted for that result. Separate `make build` exited **0**. CI for the final reviewed commit must be linked independently; no local/full verification pass is claimed.

No SQLite/history/HTTP/React runtime was implemented. Transactions, permission enforcement, request security, actual frontend build, browser/AT, platform smoke and recovery remain planned future checks. Owner technical acceptance and owner security/architecture disposition of this exact proposal are still pending.
