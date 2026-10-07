# H05-008: Verify Local History Release Evidence
- GitHub Issue: [#427](https://github.com/geoffrey-xiao/deprail/issues/427).
- GitHub parent: [#418](https://github.com/geoffrey-xiao/deprail/issues/418). Published dependencies: #420–#426; global pre-start gate #419.

## Planning metadata

- Epic: [EPIC-002 — Local History Delivery](../epics/EPIC-002-local-history-delivery.md)
- Type: `test`
- Area: `test`
- Priority: `P0`
- Risk: `R3`
- Target version: `v0.5.0`
- Milestone: [`v0.5.0`](https://github.com/geoffrey-xiao/deprail/milestone/11)
- Sprint: Sprint4
- Owner: `@geoffrey-xiao`
- Reviewer: `@geoffrey-xiao` (owner); independent review is not a gate under ADR-0005.
- Dependencies: H05-001–007, [release checklist](../../../RELEASE-CHECKLIST.md), [release plan §14](../../../03-planning/deprail-development-plan-v0.5.0.md#14-detailed-contract-reconciliation-for-current-owner-review), and #452 artifact-fix main cutover through #456. #454/#455 timed-apply is owner-merged to main 9181e76; #453 merged only to its temporary base. Earlier #419/#432 QA gate remains resolved. No release acceptance is inferred from component completion.
- Blocked reason: formal execution awaits prerequisite technical/security acceptance or explicit dispositions, reviewed blocker fixes on one immutable main candidate, and actual candidate/platform/browser-AT/recovery/supply-chain evidence. Preparation #445 is Draft; Project Blocked.

## Definition of Ready

- [x] User/release value grounded in accepted proposal [PR #417](https://github.com/geoffrey-xiao/deprail/pull/417) and owner decision [#417](https://github.com/geoffrey-xiao/deprail/pull/417#issuecomment-5907670878).
- [x] Scope and exclusions set by release checklist, [test strategy](../requirements/TEST-STRATEGY.md) and [compatibility matrix](../requirements/COMPATIBILITY-MATRIX.md).
- [x] Required workflows, outputs and failure evidence specified below.
- [x] Verification scenarios and artifacts named.
- [x] Observable acceptance specified; runtime/release acceptance remains unchecked.
- [x] Owner assigned; technical/security assessment and release decision are distinct owner decisions; no independent reviewer gate.
- [x] H05-001–007 dependencies and target version recorded; prerequisite dispositions remain a pre-start gate.

## Goal

Assemble and review integrated behavioral, browser/accessibility, security, recovery, packaging and supply-chain evidence for the v0.5 local history delivery so the owner can make explicit technical, security and release decisions. This issue adds verification/evidence, not product features or automatic release authority.

## Scope

Own version-specific release-evidence procedures/bundle using [general release checklist](../../../RELEASE-CHECKLIST.md). Exercise actual JS/Python/Java scan/capture/store/query/API/UI workflows with truthful success/partial/failed/cancelled/empty/unavailable states; actual packaged binaries and declared browser/AT combinations on Linux/macOS/Windows. Verify hostile inputs/auth/recovery/privacy, interrupted-write/rollback and manual backup recovery without touching the only copy. Compare repository trees around read-only console use. Record exact checksums and supplied/verified, unavailable or explicitly deferred SBOM/signature/provenance; no implied attestation. Trace FR-501–511, SEC-01–13, STORE-01, [test strategy](../requirements/TEST-STRATEGY.md), [product](../../../01-product/deprail-product-design-v1-ai.md), [architecture](../../../02-architecture/deprail-architecture-and-tech-stack-v1.md), [release plan §14](../../../03-planning/deprail-development-plan-v0.5.0.md#14-detailed-contract-reconciliation-for-current-owner-review).

## Out of Scope

No implementation/features, mock-only or unit-only release claims, automatic publication/tagging, weakening acceptance, silent SBOM/signature/provenance gaps, or substituting an independent reviewer for owner security/release decision. Do not mutate repositories for fixtures; use offline/license-safe source fixtures and isolated temporary repositories.

## Inputs, Outputs, and Failure Behavior

Inputs: completed H05-001–007 candidate, release targets, representative offline-safe repositories, threat/recovery scenarios and release checklist. Outputs: traceable evidence index, binary checksums, browser/AT results, repository-tree comparison, recovery record and explicit supply-chain/attestation dispositions. Any missing or failed critical scenario, unexplained tree change, security boundary failure, unreviewed vulnerability, incompatible version, or unverifiable binary is a no-go; report exact gap, do not substitute a mock or claim success. Owner makes separate technical/security acceptance and release/no-go decision; this issue never releases automatically.

## Required Tests

- Run end-to-end representative JS, Python and Java histories through actual candidate binaries; verify distinct operation/report states and truthful artifact integrity.
- Exercise actual browser UI and supported assistive technology on target OS/browser matrix: keyboard, focus, zoom/reflow, status announcements, hostile text, auth loss and recovery.
- Threat/recovery scenarios: Host/Origin/auth/cursor/path attacks; read-only proof; missing/corrupt/unsupported store; backup restore, interrupted write, quota/lock failure; confirm existing CLI behavior.
- Compare repository tree around console use; verify no scan or mutation occurred and no private path/token is exposed.
- Hash release binaries and record SBOM/signature/provenance presence/verification or explicit gap; establish artifact-to-source/build identity evidence.

## Acceptance Criteria

- [ ] Representative JavaScript, Python and Java workflows pass through integrated real runtime; failures/partial/cancelled and unavailable states remain truthful.
- [ ] Actual browser/AT and four-target packaged evidence records keyboard/accessibility, security, launch, recovery and compatibility outcomes.
- [ ] Threat, storage recovery and read-only tests pass; repository tree remains unchanged during console-only use; no sensitive token/path leaks.
- [ ] Each candidate binary has checksum and explicit SBOM/signature/provenance status with evidence or stated absence; supply-chain gaps cannot be implicit.
- [ ] Owner records distinct technical/security acceptance and explicit release/no-go decision; no automatic release or independent-review gate is imposed.

## Owner Review

Owner reviews every criterion and evidence artifact, records remaining risks, performs a distinct security assessment and makes an explicit release/no-go decision. A passing test bundle alone is not owner acceptance or release authorization.

## Evidence Required

- Verification commands or scenarios: future end-to-end matrix for JS/Python/Java; packaged OS/browser/AT matrix; hostile-input and recovery rehearsal; tree comparison; binary checksum and SBOM/signature/provenance checks.
- Expected artifacts, logs, screenshots, or links: evidence index with exact binary/source identity, sanitized scenario results, browser/AT records, recovery and before/after tree manifests, checksums, SBOM and explicit attestation status, owner security and release decisions.

### Readiness preparation — 2026-10-07

[Manual verification guide](../MANUAL-TEST-GUIDE.md) and [readiness assessment](../tracking/H05-008-READINESS.md) specify integrated procedure and outstanding decisions. Earlier #432 QA and owner-merged #455 are resolved; #453's temporary-base merge did not deliver the artifact fix to main, now tracked by #456 under #452. Prerequisite acceptance for #421 and #423–#426 remains unrecorded. Preparation #445 remains Draft; formal #427 remains Blocked. No H05-008 runtime or final acceptance box is completed.

## Final Acceptance

- [ ] Owner reviewed every acceptance criterion during PR review.
- [ ] Required verification and CI results were reviewed.
- [ ] Owner review and remaining risk are recorded.
- [ ] Evidence links are attached.
- [ ] Owner review and merge evidence are linked.

## Rollback

Record no-go, preserve immutable tags and existing user data, and publish a corrected preview under a new tag if separately authorized; never move an immutable tag or automatically downgrade/reset storage.
