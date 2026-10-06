# H05-007: Package Reproducible Embedded Console Binaries
- GitHub Issue: [#426](https://github.com/geoffrey-xiao/deprail/issues/426).
- GitHub parent: [#418](https://github.com/geoffrey-xiao/deprail/issues/418). Published dependencies: #424, #425; global pre-start gate #419.

## Planning metadata

- Epic: [EPIC-002 — Local History Delivery](../epics/EPIC-002-local-history-delivery.md)
- Type: `feature`
- Area: `foundation`
- Priority: `P0`
- Risk: `R2`
- Target version: `v0.5.0`
- Milestone: [`v0.5.0`](https://github.com/geoffrey-xiao/deprail/milestone/11)
- Sprint: Unassigned
- Owner: `@geoffrey-xiao`
- Reviewer: `@geoffrey-xiao` (owner); independent review is not a gate under ADR-0005.
- Dependencies: H05-005 transport/static-asset contract and H05-006 complete UI source/lockfile; accepted compatibility and security contracts. H05-008 consumes packaged outputs.
- Blocked reason: None for source implementation; [PR #444](https://github.com/geoffrey-xiao/deprail/pull/444) open in Project Review. #419 timeout/OutputCap follow-up and #425 owner acceptance remain open; package release gates still apply.

## Definition of Ready

- [x] User value supported by accepted proposal [PR #417](https://github.com/geoffrey-xiao/deprail/pull/417) and owner decision [#417](https://github.com/geoffrey-xiao/deprail/pull/417#issuecomment-5907670878).
- [x] Scope and exclusions specified by [API](../API-DESIGN.md), [compatibility](../requirements/COMPATIBILITY-MATRIX.md), and [security](../requirements/SECURITY-REQUIREMENTS.md) contracts.
- [x] Build inputs, outputs, mismatch and missing-asset failures specified.
- [x] Reproducibility/package scenarios are named below.
- [x] Observable acceptance specified; actual binary evidence remains unchecked.
- [x] Owner reviewer assigned; independent review is not required.
- [x] H05-005/006 dependencies and target version recorded; runtime remains gated by prerequisite disposition.

## Goal

Embed the reviewed complete local console into reproducibly built, version-matched DepRail binaries for the four declared release targets, without weakening existing CLI independence or security boundaries.

## Scope

Own `web/assets.go` (`//go:embed dist`), locked frontend build and Makefile/CI/release packaging integration. Reuse pre-introduction dependency/license/vulnerability/provenance decisions and verify the complete packaged supply chain, generated asset identity and reproducibility. Only embedded allowlisted assets; missing/skew is a visible failure, never disk/network/stale fallback. Produce exactly `linux/amd64`, `darwin/amd64`, `darwin/arm64`, `windows/amd64` per [compatibility matrix](../requirements/COMPATIBILITY-MATRIX.md). Measure gzip JS<=250KiB, CSS<=50KiB, embedded UI<=1MiB, asset-attributable binary growth<=1MiB, frontend build<=60s and listener-ready<=2s on a recorded runner; report total binary/driver delta separately. Preserve independent CLI behavior. Trace FR-509, SEC-10/12/13, [product](../../../01-product/deprail-product-design-v1-ai.md), [architecture](../../../02-architecture/deprail-architecture-and-tech-stack-v1.md) and [release plan §14](../../../03-planning/deprail-development-plan-v0.5.0.md#14-detailed-contract-reconciliation-for-current-owner-review).

## Out of Scope

No SQLite driver implementation (owned by H05-002), API response-cap increase, filesystem/network/CDN asset fallback, unreviewed dependencies/lock drift, release publication, store migration or unrelated CLI changes. Do not impose an asset-only 1 MiB budget on total binary growth including SQLite.

## Inputs, Outputs, and Failure Behavior

Inputs: H05-006 source plus immutable lockfile, H05-005 route/security profile, selected Go/platform toolchains, license/vulnerability review and build config. Outputs: reproducible embedded assets and four target binaries with matching UI/API identity, attributable asset-size evidence and supply-chain disposition. Missing assets, version mismatch, failed generation, lock drift or unaccepted dependency risk fails visibly and does not silently serve another version. Existing CLI commands remain available independently of browser assets.

## Required Tests

- From clean checkout/cache conditions, perform locked UI build twice and compare generated asset identity; verify no undeclared network fetch and record toolchain/lock versions.
- Build and launch each of four target binaries (native or declared cross-target smoke); verify matching console asset/API version and existing CLI commands.
- Remove/mismatch embedded asset in a throwaway build harness and verify explicit failure; verify no disk/CDN fallback and exact static allowlist.
- Measure compressed and uncompressed asset-attributable binary contribution and confirm API responses still enforce 1,048,576 UTF-8-byte maximum independently.
- Review all direct/transitive licenses, SBOM and vulnerability findings with explicit accepted/mitigated/rejected status.

## Acceptance Criteria

- [ ] Locked clean builds reproducibly yield the declared embedded UI; dependency, license, vulnerability and supply-chain risks have explicit disposition.
- [ ] Four target binaries contain matching assets/API; absent/skewed assets fail visibly and no filesystem/network fallback exists.
- [ ] Asset-attributable budgets and startup behavior are measured; the separate hard 1 MiB API response limit is unchanged.
- [ ] Existing CLI operates independently of console asset availability and its output/exit behavior is unchanged.
- [ ] Actual platform/package scenarios and dependency review are evidenced; no result is asserted from source-only tests.

## Owner Review

Owner reviews reproducibility, supply-chain dispositions, target outputs and compatibility evidence. Planning DoR is not actual packaged acceptance or release approval; independent review is not a gate.

## Evidence Required

- Verification commands or scenarios: future clean locked build/rebuild comparison, four-target binary launch, missing/skew injection using throwaway harness, existing CLI regression, asset/API size measurement, dependency/license/SBOM review.
- Expected artifacts, logs, screenshots, or links: binary hashes, toolchain and lock identity, asset-attributable size table, SBOM and explicit vulnerability/license disposition, sanitized startup/error results.

Implementation evidence: [H05-007 packaging evidence](../tracking/H05-007-PACKAGING-EVIDENCE.md) and [PR #444](https://github.com/geoffrey-xiao/deprail/pull/444) record the pinned Darwin build, offline repeat, actual console browser/static smoke, throwaway failure binaries, size/hash table, four target cross-builds, and unresolved platform/supply-chain owner gates. Acceptance boxes remain unchecked until CI/native platform evidence and owner review are recorded.

## Final Acceptance

- [ ] Owner reviewed every acceptance criterion during PR review.
- [ ] Required verification and CI results were reviewed.
- [ ] Owner review and remaining risk are recorded.
- [ ] Evidence links are attached.
- [ ] Owner review and merge evidence are linked.

## Rollback

Use a matching previously reviewed binary/assets; preserve history database/WAL/artifacts and never automatically downgrade the store. Existing CLI remains usable.
