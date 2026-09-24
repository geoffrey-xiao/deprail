# v0.5 Compatibility Matrix

**Status:** Draft; supported exact versions are not yet frozen.

| Surface | Compatibility commitment | Required evidence / open decision |
| --- | --- | --- |
| CLI commands and exit codes | Preserve existing command, stdout/stderr, JSON purity, error-code, and exit-code contracts. Candidate `deprail scan --save-history` and additive exit code `6` remain unapproved; they require a separate compatibility decision and must not alter existing scans. | Regression scenarios for history absent, unavailable, and explicitly selected; review precedence when scan and persistence both fail. |
| Scan report schema | Preserve existing CLI `ScanReport` JSON bytes/meaning, stable keys where defined, ordering, provenance, and completeness enum (`complete`, `partial`, `failed`). Do not claim current CLI JSON validates as the proposed v1alpha scan document or rewrite it for history. Cancellation is not added to `ScanReport.Status` in v0.5. | Compare actual nonempty/failed CLI JSON before/after; explicitly reconcile the known app-report/v1alpha-schema discrepancy independently before any CLI/schema change. Verify cancellation remains separate from report status. |
| Persisted history schema | Owner-selected candidate: DB `PRAGMA user_version=1`, distinct allowlisted `history-v1` projection and source-report version metadata; `/api/v1` is another separately versioned contract. UUIDv4 `historyEntryID` differs from repeatable source `ScanReport.ScanID`; no raw report/project JSON or automatic legacy import. None of these versions is interchangeable or yet accepted for implementation. | Approve projection fields, stable-key mapping, safe diagnostics, privacy/unknown handling and row/JSON validation; then review uniqueness/retry semantics, refs, migrations, backup/recovery, retention and downgrade refusal. Numeric limits need payload evidence. |
| Local API | New versioned contract; draft URL version `/api/v1` is a proposal. No compatibility promise until OpenAPI is accepted. | OpenAPI 3.1 validation, additive/breaking policy, client/version mismatch behavior. |
| Embedded web client | Shipped with matching binary; exact supported browsers TBD. | Package smoke on selected browsers/OS, UI/API version skew test, no external CDN/font/analytics. |
| Web frontend toolchain | React/TypeScript/Vite per architecture; exact pinned versions and shadcn/ui component/primitives strategy TBD. | Lockfile, license and dependency review, clean build, bundle/startup budget. |
| SQLite engine/driver | SQLite local store per roadmap/architecture; engine/driver/version and build strategy TBD. | Cross-platform availability, CGO/static binary implications, migration and locking tests. |
| Local listener | Proposed loopback-only; address/port/lifecycle/auth/origin behavior TBD. | Linux/macOS/Windows bound-address and hostile-origin test evidence. |
| Filesystem paths/permissions | Canonical containment, `/` serialized relative paths, platform-equivalent meaning. | Unicode, long path, symlink, permission, interrupted write, data-root tests per OS. |
| Data retention/deletion | Draft proposal retains explicitly selected history without automatic eviction, per-entry deletion/export, or raw-artifact cleanup. Unbounded growth risk is unresolved pending a measured cap or explicit owner acceptance. | Capacity, disk-full, artifact-reference behavior, user guidance, and cross-platform recovery review before DoR. |
| v0.4 remediation/scan use | No remediation contract changes. Existing CLI scanning should not depend on web or history absent an approved decision. | Regression/manual comparison; explicit disposition of #356 follow-ups. |

## Compatibility decision gates

- Pin Go/Node/package manager, SQLite driver/runtime, frontend packages, supported browser versions, and build embedding tool only after a reproducibility/security review.
- Do not silently convert existing scan JSON or discard fields from its established CLI contract. The candidate history projection is a different, allowlisted representation: classify every source field needed for its semantics, reject unsafe/unclassified required data before save, and never pretend omitted source data is a complete report.
- Define API additive changes, deprecation/version transition, persisted schema migrations, downgrade refusal, and recovery before implementation.
- New history failure behavior must not silently alter CLI status/outputs. Candidate `--save-history`/exit code `6` behavior is unapproved; any accepted change requires an explicit compatibility decision and updated CLI contracts.
- Confirm binaries and data handling on Linux, macOS, Windows. Platform-specific storage limitations must be explicit, not hidden by fallback behavior.

## Component and platform decision inputs (unselected)

These are review criteria, not selected versions or dependencies. Record selected versions, owners, and evidence in the approved execution package before implementation.

| Component | Decision inputs | Required evidence |
| --- | --- | --- |
| Go toolchain | Preserve the currently approved repository toolchain unless a reviewed compatibility change is made; identify supported OS/architectures. | Pinned toolchain declaration, reproducible build, Linux/macOS/Windows build results. |
| SQLite engine/driver | Confirm maintenance/security support, license, WAL/foreign-key/online-backup/`synchronous=FULL` behavior, cancellation and busy handling, CGO/static-build consequences, and supported SQLite versions. Do not select a driver in this contract. | Driver/engine version matrix, license/SBOM review, feature tests, cross-platform package evidence. |
| Node/package manager/frontend | Pin Node and package-manager versions; establish lockfile integrity, supported React/Vite/build tool versions, transitive-license review, reproducible embedded assets, and binary-size/startup budgets. | Clean build from lockfile, license report, SBOM, asset digest and size/startup record. |
| Browsers/assistive technology | Owner and reviewer select browser products/versions and OS combinations from the accepted UX needs; select screen-reader/assistive-technology pairings. Do not infer support from CI or local availability. | Exact browser/OS/AT matrix, versioned actual-surface accessibility and API/UI-skew results. |
| Packaging/platform | Freeze OS/architecture targets and embedded-asset strategy; no filesystem fallback or unsupported listener exposure. | Per-target binary identity/checksum and smoke evidence on Linux, macOS, Windows. |

## Compatibility, rollback, and release-evidence matrix

| Change/failure | Preserved contract | Rollback/recovery rule | Required evidence |
| --- | --- | --- | --- |
| CLI run without history/web or with local service unavailable | Existing v0.1/v0.2 commands, JSON/stdout/stderr, report meaning, and exit meanings remain unchanged. Candidate `--save-history` and exit code `6` are unapproved. | Revert the application change; do not alter established scan codes or reserved exit-code meanings. | Linux/macOS/Windows CLI regression and tree comparison; `FR-505`. |
| History write fails after a selected scan | Do not claim saved; rollback the entire new row/reference set and preserve scan outcome/report. | Preserve previous DB and external raw artifacts; any exit-status change needs a separate compatibility decision. | Forced rollback, exact stdout/stderr/exit record, DB and artifact digests; `STORE-01`. |
| Source cannot form a safe `history-v1` projection or stored row/projection disagrees | Keep original CLI bytes/outcome; do not claim saved on selected capture, and do not serve an unsafe or falsely empty entry. | Refuse pre-commit save or quarantine read as typed failure without rewriting prior entries; any exit-status change needs separate approval. | Nonempty/failed source report, absolute root, raw diagnostic, missing/duplicate stable key, unknown field, row mismatch; `STORE-01`, `SEC-07`, `SEC-08`. |
| Forward migration fails or the database is corrupt/future-version | Never downgrade/repair/replace automatically; retain committed history and validated pre-migration copy. Older binaries must refuse unsupported newer schemas. | Stop the application; preserve current DB and WAL/SHM; restore a validated backup to a new file only through an explicit operator action. | Fresh/current/prior/interrupted/corrupt/future migration matrix, backup digests and manual recovery record; `FR-508`, `SEC-08`. |
| Artifact reference missing or digest mismatched | Keep only trustworthy metadata; never trust, fabricate, or delete shared raw bytes. | Restore/recover artifact separately from DB; history transaction rollback does not garbage-collect artifacts. | Missing/mismatch fixtures and content digests; `FR-503`, `SEC-09`. |
| API/UI contract or embedded asset mismatch | Fail visibly; never return/truncate into empty success; CLI remains independent. | Roll back to a matching binary/API/asset set, not a mixed-version package. | Packaged version-skew capture on every approved OS/browser; `FR-509`, `SEC-10`. |
| Toolchain, SQLite driver, frontend, or browser support changes | No dependency/version is selected by this draft; existing CLI meaning stays stable and separately versioned `history-v1` changes require their own migration review. | Restore the previous reviewed lockfile/toolchain/package set; any persisted-schema change follows the migration/backup rule above. | Reproducible build, license/SBOM review, package checksums, cross-platform/browser results; `SEC-13`. |

The v0.5.0 release-evidence record must link these artifacts to the reviewed commit/binary, exact command, result, OS/architecture, database/fixture digest, and reviewer decision. These rows define future evidence; no runtime or platform pass is claimed here.
