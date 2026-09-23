# v0.5 Compatibility Matrix

**Status:** Draft; supported exact versions are not yet frozen.

| Surface | Compatibility commitment | Required evidence / open decision |
| --- | --- | --- |
| CLI commands and exit codes | Preserve existing command, stdout/stderr, JSON purity, error-code, and exit-code contracts. Candidate `deprail scan --save-history` and additive exit code `6` remain unapproved; they require a separate compatibility decision and must not alter existing scans. | Regression scenarios for history absent, unavailable, and explicitly selected; review precedence when scan and persistence both fail. |
| Scan report schema | Preserve current versioned report meaning, stable keys, ordering, provenance, and completeness enum (`complete`, `partial`, `failed`). Cancellation is not added to `ScanReport.Status` in v0.5. | Validate the unchanged report schema; verify cancellation is represented outside the report and cannot be shown as completed based on retained completeness. |
| Persisted history schema | Proposed schema v1 uses a UUIDv4 `historyEntryID` distinct from repeatable source `ScanReport.ScanID`; `PRAGMA user_version` versions the database, while stored report/project documents retain their own schema versions. No legacy scan is imported automatically. | ADR defines uniqueness/retry semantics, stored fields and artifact references, forward-only migrations, validated backup/recovery, retention, and downgrade refusal; numeric limits remain open pending payload evidence. |
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
- Do not silently convert existing scan JSON or discard fields; unknown safe fields follow the current contract's preservation rules.
- Define API additive changes, deprecation/version transition, persisted schema migrations, downgrade refusal, and recovery before implementation.
- New history failure behavior must not silently alter CLI status/outputs. Candidate `--save-history`/exit code `6` behavior is unapproved; any accepted change requires an explicit compatibility decision and updated CLI contracts.
- Confirm binaries and data handling on Linux, macOS, Windows. Platform-specific storage limitations must be explicit, not hidden by fallback behavior.
