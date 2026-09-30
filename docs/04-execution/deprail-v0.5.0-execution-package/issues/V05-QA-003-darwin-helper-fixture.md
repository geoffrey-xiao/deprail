# V05-QA-003 Stabilize the Darwin Verification Helper Fixture
- GitHub Issue: [#432](https://github.com/geoffrey-xiao/deprail/issues/432).
- GitHub parent: [#418](https://github.com/geoffrey-xiao/deprail/issues/418).

## Planning metadata

- Type: `bug`
- Area: `test`
- Priority: `P0`
- Risk: `R2`
- Target version: `0.5.0`
- Milestone: `v0.5.0`
- Sprint: Sprint 4.
- Owner: `@geoffrey-xiao`
- Reviewer: `@geoffrey-xiao`; independent review optional under ADR-0005.
- Dependencies: accepted #419 investigation / merged PR #429, corrected capture #430 / owner-merged PR #431; parent delivery epic #418.
- Blocked reason: None for this authorized test-only correction; H05 runtime remains gated on verified correction and owner review.

## Definition of Ready

- [x] Value and user impact are stated.
- [x] Scope and explicit exclusions are stated.
- [x] Inputs, outputs, and failure behavior are defined.
- [x] Required tests and actual smoke scenarios are named.
- [x] Acceptance criteria are observable.
- [x] Owner reviewer and dependencies are recorded.
- [x] Release/version/Sprint and evidence are explicit.

The owner's “请继续” authorizes the announced darwin-only fixture correction following owner merge of #431. This restores reliable verification of the existing workspace/deadline contract: product verification/security, architecture testing/external-process boundary, roadmap cross-platform invariants and release plan §14/§15 readiness, not a new feature or relaxed runtime deadline.

## Goal

Remove avoidable cold-image startup from the darwin test fixture without weakening command execution, one-second success deadline, timeout/cancellation checks or production process behavior. Preserve the original failed runs as evidence.

## Scope

Only internal/remediation/verification/run_test.go testTool: on darwin, create the allowlisted npm symlink to os.Executable(), the currently running test image; otherwise retain the existing byte-copy fixture, including npm.exe on Windows. Fail fixture creation on error; no fallback, retries, prewarming or skipped assertions. Record exact corrected-baseline targeted/full verification, real helper/cwd/deadline smoke and three-OS CI. Reuse existing consumer behavior tests; no permanent test of symlink implementation or source text.

## Out of Scope

No production process/verification/scanner/API/CLI/storage changes, new error code, arbitrary longer deadline, timeout removal, assertion weakening, dependency, H05 history/UI, release or automatic merge. Do not claim a specific OS security subsystem caused the historical failure; that remains unknown. This issue does not automatically grant runtime kickoff.

## Inputs, Outputs, and Failure Behavior

Fixture input: current test executable and isolated temporary directory. Output: allowlisted executable path which launches the same controlled helper modes. Darwin uses a symlink to the active image; Linux/Windows keep byte copying and existing names. Symlink/path resolution failure is a test failure, not a skip or fallback. Production Run deadlines, working-directory containment, output cap and error precedence remain unchanged.

## Required Tests

Existing TestRunExecutesSelectedCommandsInWorkspace, stop-after-failure, timeout, cancellation, untrusted-command and workspace-escape tests; targeted verification/process suites; one full make verify on the corrected baseline. Observe actual helper success/failure/sleep and workspace output in controlled real subprocess smoke. Three-OS CI proves respective test execution, not unrun native local platforms. Original recorded failures are not rerun merely to confirm them or relabeled passed.

## Acceptance Criteria

- [x] Existing actual command execution and partial-result tests pass without modifying their deadlines or assertions.
- [x] Intentional 20ms timeout, cancellation and workspace containment behavior remain effective; output capture regressions from #430 pass.
- [x] Corrected-baseline make verify passes, or any new failure is reported accurately with no success-by-retry substitution.
- [x] Actual helper/cwd/deadline smoke and Linux/macOS/Windows CI are linked; darwin-only change does not alter Linux/Windows fixture behavior.
- [x] Owner reviews the correction and explicitly records the remaining runtime-start decision; historical cold-start/OS uncertainty and original failures remain preserved.

## Owner Review

Inspect the test-only scope and no-copy active-image fixture, genuine timeout/cancellation/containment evidence, unchanged production semantics and remaining uncertainty. Record technical acceptance and security/architecture assessment separately; external reviewer optional. Owner merge or disposition of this correction is required before declaring the timeout readiness gate resolved.

## Evidence Required

Baseline/commit/platform/toolchain, exact source change and smoke/targeted/full commands/results, historical failure links, three-OS CI, owner review and remaining readiness decision. No raw environment or credential/path disclosure.

## Rollback

Revert the test-only correction through a reviewed change; production/user data remains unchanged. Retain the known cold-copy failure and keep H05 startup gated until corrected or explicitly dispositioned.

## Final Acceptance

- [x] Owner reviewed all acceptance criteria and local/CI evidence.
- [x] Technical acceptance and separate security/architecture assessment are recorded.
- [x] Owner-reviewed merge and remaining runtime-start decision are linked.
- [x] Historical failed runs remain accurate; no release approval or H05 implementation is inferred.

## Contract links

[Product verification](../../../01-product/deprail-product-design-v1-ai.md#remediation-engine), [architecture testing](../../../02-architecture/deprail-architecture-and-tech-stack-v1.md#10-test-architecture), [roadmap invariants](../../../03-planning/deprail-roadmap-v1.md#cross-release-invariants), [release plan readiness](../../../03-planning/deprail-development-plan-v0.5.0.md#15-owner-acceptance-and-implementation-backlog-publication), [accepted helper proposal](../tracking/V05-QA-001-INVESTIGATION.md#proposed-bounded-correction-and-next-decision), [capture correction](V05-QA-002-bounded-output-capture.md).

## Actual correction evidence

[Helper/targeted/full evidence](../tracking/V05-QA-003-EVIDENCE.md) records the real corrected fixture's success, exit7, 20ms timeout, cancellation and canonical workspace behavior; uncached verification/process suites and single-run full `make verify` passed. The throwaway workspace probe's initial argument-separator error is disclosed and corrected; no permanent assertion/deadline changed. At review submission, three-OS CI, owner decisions and merge were pending; the additive completion record below resolves those gates without changing historical failed runs or OS-cause uncertainty.

## Owner acceptance and completion

Owner-merged [PR #433](https://github.com/geoffrey-xiao/deprail/pull/433) at `6c74cf9861981928c29f7a8cb09799e63f3fd6dd` contains the reviewed correction. [CI run 36725991235](https://github.com/geoffrey-xiao/deprail/actions/runs/36725991235) passed on tested head `b6a96ee34b1053b1c860fa9bf22dad17d19c741b`: Ubuntu 26s, macOS 1m5s, Windows 1m40s.

The owner selected “确认验收并完成关联”, confirming every criterion, technical acceptance and the separate test-fixture security/architecture assessment; residual macOS cold-start root-cause uncertainty was accepted. [Durable owner record](https://github.com/geoffrey-xiao/deprail/issues/432#issuecomment-5913275072) records that conversation decision without fabricating a native APPROVED review. Issue #432 is closed, Project Done, with native Linked pull requests referencing #433.

This completed QA gate does not itself start another issue or approve a release. The subsequent “可以 继续吧” separately authorizes H05-001 / #420 under its own DoR and synchronized-main gate; see [kickoff decision](https://github.com/geoffrey-xiao/deprail/issues/420#issuecomment-5913450360). Original failed runs remain failed.
