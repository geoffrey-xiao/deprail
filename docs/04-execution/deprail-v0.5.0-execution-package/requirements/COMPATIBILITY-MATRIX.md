# v0.5 Compatibility Matrix

**Status:** Draft; supported exact versions are not yet frozen.

| Surface | Compatibility commitment | Required evidence / open decision |
| --- | --- | --- |
| CLI commands and exit codes | Preserve existing command, stdout/stderr, JSON purity, error-code, and exit-code contracts. Owner-accepted v0.5 addition: explicit `deprail scan --save-history` and exit code `6` only when a requested save fails after an otherwise-successful scan; ordinary scans remain unchanged. | Regression scenarios for history absent, unavailable, and explicitly selected; review precedence when scan and persistence both fail. |
| Scan report schema | Preserve existing CLI `ScanReport` JSON bytes/meaning, stable keys where defined, ordering, provenance, and completeness enum (`complete`, `partial`, `failed`). Do not claim current CLI JSON validates as the proposed v1alpha scan document or rewrite it for history. Cancellation is not added to `ScanReport.Status` in v0.5. | Compare actual nonempty/failed CLI JSON before/after; explicitly reconcile the known app-report/v1alpha-schema discrepancy before any CLI/schema change. Verify cancellation remains separate from report status. |
| Persisted history schema | Owner-selected candidate: DB `PRAGMA user_version=1`, distinct allowlisted `history-v1` projection and source-report version metadata; `/api/v1` is another separately versioned contract. UUIDv4 `historyEntryID` differs from repeatable source `ScanReport.ScanID`; no raw report/project JSON or automatic legacy import. None of these versions is yet accepted for implementation. | Owner accepts projection fields, stable-key mapping, safe diagnostics, privacy/unknown handling and row/JSON validation; then verifies uniqueness/retry semantics, refs, migrations, backup/recovery, retention and downgrade refusal. Numeric limits need payload evidence. |
| Local API | Candidate contract: OpenAPI 3.1 `info.version: 1.0.0` at `/api/v1`; route, response, error, pagination, and security profile remain proposals. No client compatibility promise until the exact API/projection contract is owner-accepted. | Candidate OpenAPI references/examples and status-specific error constraints were validated offline; the 1 MiB serialized-byte cap remains a runtime gate and is not inferred from schema `maxLength`. Owner accepts the exact projection, version behavior, and UI/API mismatch handling. |
| Embedded web client | Matching binary/assets and current-stable browser/AT support policy in §Selected engineering contract. | Actual package smoke, UI/API skew and browser/AT evidence after implementation; no CDN/font/analytics. |
| Web frontend toolchain | Exact proposed pins, project-owned semantic primitives, lockfile and embed ownership below. | Owner contract review now; locked build, complete license/SBOM review and measured budgets after implementation. |
| SQLite engine/driver | CGo-free `modernc.org/sqlite v1.59.0` review candidate; linked SQLite engine version recorded from the resolved module, not inferred. | Driver/engine licenses and provenance; actual four-target build, migrations, locks and recovery before release. |
| Local listener | Candidate `127.0.0.1` ephemeral port with process-scoped bearer, strict Host/Origin, no CORS, bounded requests, and explicit foreground lifecycle; none is runtime-approved. | Owner security decision plus Linux/macOS/Windows bind, hostile-origin, lifecycle, and shutdown evidence. |
| Filesystem paths/permissions | Canonical containment, `/` serialized relative paths, platform-equivalent meaning. | Unicode, long path, symlink, permission, interrupted write, data-root tests per OS. |
| Data retention/deletion | No eviction/delete/export; atomic admission refuses >1,000 entries or >256 MiB logical projection bytes. Per-entry max16 MiB; physical SQLite/WAL/index/backup and external artifacts are not capped. | Owner accepts remaining physical-growth risk; actual quota/disk-full/ref-integrity/recovery evidence after implementation. |
| v0.4 remediation/scan use | No remediation contract changes. Existing CLI scanning should not depend on web or history absent an approved decision. | Regression/manual comparison; explicit disposition of #356 follow-ups. |

## Compatibility decision gates

### Accepted v0.5 CLI decision

Earlier draft wording: “Candidate `deprail scan --save-history` and additive exit code `6` remain unapproved; they require a separate compatibility decision and must not alter existing scans.” That wording records prior status and is superseded by the exact CLI contract in [v0.5 development plan §§14–15](../../../03-planning/deprail-development-plan-v0.5.0.md). The owner accepted the `--save-history` flag and additive exit code `6`; implementation must preserve scan output bytes, avoid history initialization without the flag, keep scan failure primary, and return `6` only when the scan would otherwise exit `0` and requested persistence fails. Existing meanings for exit codes `1`–`5` do not change.


- Pin Go/Node/package manager, SQLite driver/runtime, frontend packages, supported browser versions, and build embedding tool only after the owner reviews reproducibility and security evidence.
- Do not silently convert existing scan JSON or discard fields from its established CLI contract. The candidate history projection is a different, allowlisted representation: classify every source field needed for its semantics, reject unsafe/unclassified required data before save, and never pretend omitted source data is a complete report.
- Define API additive changes, deprecation/version transition, persisted schema migrations, downgrade refusal, and recovery before implementation.
- New history failure behavior must preserve scan status and outputs. The owner-accepted v0.5 `--save-history`/exit code `6` behavior is defined above; do not alter existing scans or codes `1`–`5`.
- Confirm binaries and data handling on Linux, macOS, Windows. Platform-specific storage limitations must be explicit, not hidden by fallback behavior.

## Component and platform decision rationale

The inputs below explain the selected engineering proposal later in this document; they are not unresolved alternatives. Owner acceptance and implementation evidence remain distinct.

| Component | Decision inputs | Required evidence |
| --- | --- | --- |
| Go toolchain | Preserve the currently approved repository toolchain unless a reviewed compatibility change is made; identify supported OS/architectures. | Pinned toolchain declaration, reproducible build, Linux/macOS/Windows build results. |
| SQLite engine/driver | Confirm maintenance/security support, license, WAL/foreign-key/online-backup/`synchronous=FULL` behavior, cancellation and busy handling, CGO/static-build consequences, and supported SQLite versions. Do not select a driver in this contract. | Driver/engine version matrix, license/SBOM review, feature tests, cross-platform package evidence. |
| Node/package manager/frontend | Pin Node and package-manager versions; establish lockfile integrity, supported React/Vite/build tool versions, transitive-license review, reproducible embedded assets, and binary-size/startup budgets. | Clean build from lockfile, license report, SBOM, asset digest and size/startup record. |
| Browsers/assistive technology | Owner selects browser products/versions and OS combinations from accepted UX needs; select screen-reader/assistive-technology pairings. Do not infer support from CI or local availability. | Exact browser/OS/AT matrix, versioned actual-surface accessibility and API/UI-skew results. |
| Packaging/platform | Freeze OS/architecture targets and embedded-asset strategy; no filesystem fallback or unsupported listener exposure. | Per-target binary identity/checksum and smoke evidence on Linux, macOS, Windows. |

## Compatibility, rollback, and release-evidence matrix

| Change/failure | Preserved contract | Rollback/recovery rule | Required evidence |
| --- | --- | --- | --- |
| CLI scan without `--save-history` | Existing v0.1/v0.2 commands, JSON/stdout/stderr, report meaning, and exit meanings remain unchanged; history is not initialized. | Revert the opt-in wiring without changing established scan codes or reserved exit-code meanings. | Linux/macOS/Windows CLI regression and tree comparison; `FR-505`. |
| History write fails after a selected scan | Do not claim saved; rollback the entire new row/reference set and preserve scan output. Preserve a nonzero/incomplete scan result; return `6` only when the requested save is the sole failure. | Preserve previous DB and external raw artifacts. | Forced rollback, exact stdout/stderr/exit record, DB and artifact digests; `STORE-01`. |
| History source cannot form a safe `history-v1` projection or stored row/projection disagrees | Keep original CLI bytes/outcome; do not claim saved on selected capture, and do not serve an unsafe or falsely empty entry. A failure of the requested save returns `6` only if the scan otherwise succeeds. | Refuse pre-commit save or quarantine read as typed failure without rewriting prior entries. | Nonempty/failed source report, absolute root, raw diagnostic, missing/duplicate stable key, unknown field, row mismatch; `STORE-01`, `SEC-07`, `SEC-08`.
| Forward migration fails or the database is corrupt/future-version | Never downgrade/repair/replace automatically; retain committed history and validated pre-migration copy. Older binaries must refuse unsupported newer schemas. | Stop the application; preserve current DB and WAL/SHM; restore a validated backup to a new file only through an explicit operator action. | Fresh/current/prior/interrupted/corrupt/future migration matrix, backup digests and manual recovery record; `FR-508`, `SEC-08`. |
| Artifact reference missing or digest mismatched | Keep only trustworthy metadata; never trust, fabricate, or delete shared raw bytes. | Restore/recover artifact separately from DB; history transaction rollback does not garbage-collect artifacts. | Missing/mismatch fixtures and content digests; `FR-503`, `SEC-09`. |
| API/UI contract or embedded asset mismatch | Fail visibly; never return/truncate into empty success; CLI remains independent. | Roll back to a matching binary/API/asset set, not a mixed-version package. | Packaged version-skew capture on every approved OS/browser; `FR-509`, `SEC-10`. |
| Toolchain, SQLite driver, frontend, or browser support changes | No dependency/version is selected by this draft; existing CLI meaning stays stable and separately versioned `history-v1` changes require their own migration review. | Restore the previous reviewed lockfile/toolchain/package set; any persisted-schema change follows the migration/backup rule above. | Reproducible build, license/SBOM review, package checksums, cross-platform/browser results; `SEC-13`. |

The v0.5.0 release-evidence record must link these artifacts to the owner-reviewed commit/binary, exact command, result, OS/architecture, database/fixture digest, and owner decision. These rows define future evidence; no runtime or platform pass is claimed here.

## Selected engineering frontend/platform contract (owner review; evidence not executed)

Exact selected review candidates: Node `22.23.3`, npm `10.9.9`, React/React DOM `19.1.1`, TypeScript `5.9.2`, Vite `7.3.6`, `@vitejs/plugin-react` `5.0.2`, `@types/react` `19.1.10`, `@types/react-dom` `19.1.7`, `@types/node` `22.18.1`. No frontend pins exist in the repository. Sources: [Node metadata](https://registry.npmjs.org/node/22.23.3), [npm metadata](https://registry.npmjs.org/npm/10.9.9), [Vite metadata](https://registry.npmjs.org/vite/7.3.6); other package/version metadata uses the same registry path. Isolated script-free lock resolution/audit rejected Vite `7.1.4` with known high-severity advisories and reported zero vulnerabilities for the amended set. Audit used host Node `26.7.0`/npm `11.19.0`, not proposed build pins. Complete license/SPDX/SBOM, provenance, tarball-integrity and reproducible pinned build remains a later gate; audit is not security approval.

Project-owned semantic React UI; no third-party component/font/icon runtime dependencies. Source/lockfile `web/`, generated `web/dist`, embedded at compile time by Go-owned `embed.FS`; no filesystem fallback. Static allowlist: `/console/`, `/console/history`, `/console/scans/{historyEntryID}` (canonical UUIDv4), `/console/about`, declared hashed assets; unknown paths 404; API dispatch separate. HTTP UI is query-only; capture remains explicit CLI/application use case.

Policy budgets (not benchmarks): gzip JS ≤250 KiB, CSS ≤50 KiB, embedded UI ≤1 MiB, build ≤60s, asset-attributable binary growth ≤1 MiB against same Go/toolchain build without frontend assets, startup listener-ready ≤2s on recorded reference machine. Report total binary growth including SQLite separately.

| Surface | Supported policy | Evidence recorded per run |
| --- | --- | --- |
| Windows browser/AT | Current stable Chrome with NVDA | Exact OS/build, Chrome/NVDA versions, hardware, commit, scenario/result |
| macOS browser/AT | Current stable Safari with VoiceOver | Exact OS/build, Safari/VoiceOver versions, hardware, commit, scenario/result |
| Linux browser/AT | Current stable Firefox with Orca | Exact distribution/kernel, Firefox/Orca versions, hardware, commit, scenario/result |
| Binary targets | Existing release workflow: linux/amd64, darwin/amd64, darwin/arm64, windows/amd64 | Exact OS version, target, checksum and smoke result |

No installed versions or passing platform/browser evidence are asserted. Record actual versions exercised, never guessed.


Offline recovery: stop binary; preserve DB, matching WAL/SHM, artifact directory together. Copy byte-for-byte while stopped to separate recovery location; record hashes and retain originals. Never reset, delete, repair-in-place or downgrade. On a copy only, use SQLite read-only mode to inspect `PRAGMA user_version` and `PRAGMA quick_check`; retain output. Future schema is refused without writes. Restore only validated backup to separate new location by explicit operator action; test quick_check/schema/artifact digests before deliberate switch. No validated backup: retain all files and escalate. Artifact unavailable/unverified remains explicit; never serialize absolute root, raw errors or unknown fields.

Release/test recommendations for parent: locked build + license/SBOM review; budgets and embed integrity; static route/API distinction, skew and asset-integrity checks; four-target package smoke; schema future-version refusal and offline WAL/SHM recovery rehearsal; artifact absence without absolute-root leakage; actual browser/AT keyboard, focus, contrast, 320px/200% zoom and hostile-text checks. All future gates, not executed evidence.
