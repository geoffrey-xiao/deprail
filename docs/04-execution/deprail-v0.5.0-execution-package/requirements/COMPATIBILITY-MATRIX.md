# v0.5 Compatibility Matrix

**Status:** Draft; supported exact versions are not yet frozen.

| Surface | Compatibility commitment | Required evidence / open decision |
| --- | --- | --- |
| CLI commands and exit codes | Preserve existing command, stdout/stderr, JSON purity, error-code, and exit-code contracts unless a reviewed additive/compatibility decision says otherwise. | Regression suite with history absent, unavailable, and enabled; approval for any behavior change. |
| Scan report schema | Preserve current versioned report meaning, stable keys, ordering, provenance, and completeness enum (`complete`, `partial`, `failed`). Cancellation is not added to `ScanReport.Status` in v0.5. | Validate the unchanged report schema; verify cancellation is represented outside the report and cannot be shown as completed based on retained completeness. |
| Persisted history schema | New local format with a unique `historyEntryID` distinct from source `ScanReport.ScanID`; exact version and supported migrations TBD. | ADR defines ID generation, ingestion retry/idempotency, source-version fields, forward path, backup/recovery, and downgrade policy. |
| Local API | New versioned contract; draft URL version `/api/v1` is a proposal. No compatibility promise until OpenAPI is accepted. | OpenAPI 3.1 validation, additive/breaking policy, client/version mismatch behavior. |
| Embedded web client | Shipped with matching binary; exact supported browsers TBD. | Package smoke on selected browsers/OS, UI/API version skew test, no external CDN/font/analytics. |
| Web frontend toolchain | React/TypeScript/Vite per architecture; exact pinned versions and shadcn/ui component/primitives strategy TBD. | Lockfile, license and dependency review, clean build, bundle/startup budget. |
| SQLite engine/driver | SQLite local store per roadmap/architecture; engine/driver/version and build strategy TBD. | Cross-platform availability, CGO/static binary implications, migration and locking tests. |
| Local listener | Proposed loopback-only; address/port/lifecycle/auth/origin behavior TBD. | Linux/macOS/Windows bound-address and hostile-origin test evidence. |
| Filesystem paths/permissions | Canonical containment, `/` serialized relative paths, platform-equivalent meaning. | Unicode, long path, symlink, permission, interrupted write, data-root tests per OS. |
| Data retention/deletion | Not defined by roadmap; must be decided, documented, and migration-safe before implementation. | Capacity, disk-full, export/deletion, artifact-reference behavior review. |
| v0.4 remediation/scan use | No remediation contract changes. Existing CLI scanning should not depend on web or history absent an approved decision. | Regression/manual comparison; explicit disposition of #356 follow-ups. |

## Compatibility decision gates

- Pin Go/Node/package manager, SQLite driver/runtime, frontend packages, supported browser versions, and build embedding tool only after a reproducibility/security review.
- Do not silently convert existing scan JSON or discard fields; unknown safe fields follow the current contract's preservation rules.
- Define API additive changes, deprecation/version transition, persisted schema migrations, downgrade refusal, and recovery before implementation.
- New history failure behavior must not silently alter CLI status/outputs. If the accepted design changes exit status or scan persistence semantics, record a compatibility decision and update CLI contracts first.
- Confirm binaries and data handling on Linux, macOS, Windows. Platform-specific storage limitations must be explicit, not hidden by fallback behavior.
