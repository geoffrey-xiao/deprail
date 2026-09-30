# V05-QA-001 Resolve the Local Verification Timeout Readiness Gate
- GitHub Issue: [#419](https://github.com/geoffrey-xiao/deprail/issues/419).
- GitHub parent: [#418](https://github.com/geoffrey-xiao/deprail/issues/418). Dependencies: accepted PR #417 / owner decision.

## Planning metadata

- Type: `bug`
- Area: `test`
- Priority: `P0`
- Risk: `R2`
- Target version: `0.5.0`
- Milestone: `v0.5.0`
- Sprint: Unassigned; no runtime Sprint kickoff in this planning handoff.
- Owner: `@geoffrey-xiao`
- Reviewer: `@geoffrey-xiao`; external review optional under ADR-0005.
- Dependencies: [owner acceptance](https://github.com/geoffrey-xiao/deprail/pull/417#issuecomment-5907670878), [merged PR #417](https://github.com/geoffrey-xiao/deprail/pull/417).
- Parent: [V05-EPIC-002](../epics/EPIC-002-local-history-delivery.md).
- Project status at publication: Todo.
- Blocked reason: None for investigation. Resolution or explicit owner disposition gates every runtime implementation start.

## Definition of Ready

- [x] Value and user impact are stated.
- [x] Scope and explicit exclusions are stated.
- [x] Inputs, outputs, and failure behavior are defined.
- [x] Required tests or smoke scenarios are named.
- [x] Acceptance criteria are observable.
- [x] Owner reviewer is assigned; external reviewer optional.
- [x] Dependencies and target version are recorded.

Specification readiness is not a passing verification result or implementation authorization. Publication is authorized by the linked owner decision; this issue's evidence and disposition are not complete.

## Goal

Establish a defensible resolution/disposition of the observed local full-suite timeout before v0.5 runtime work starts. Do not turn the known failed check into a pass by replacing it with a successful rerun.

## Scope

- Investigate the recorded `make verify` exit 2 on darwin/arm64, Go 1.27.1: `internal/remediation/verification.TestRunExecutesSelectedCommandsInWorkspace` failed with `SCANNER_TIMEOUT: process exceeded configured deadline`.
- The existing test calls `Run` with a one-second deadline (`run_test.go:52`) and a copied helper executable; root cause is not established. Record startup/execution/cancellation timing and relevant scheduling/process evidence in a controlled diagnostic scenario, not a repeat merely to confirm the already observed failure.
- Own investigation/evidence in `internal/remediation/verification/run_test.go`, `run.go` and the process boundary only as necessary to locate the cause. A behavior-changing fix requires its own bounded plan and review before edits; no predetermined deadline increase or retry is authorized.
- Record the final owner decision in this contract and release tracking: verified correction, or an explicit narrowly bounded residual-risk disposition with owner and next gate.
- Traceability: release plan §14 acceptance/tracking, [test strategy §11](../requirements/TEST-STRATEGY.md#11-planning-contract-smoke-evidence-pr-417), architecture external-process/testing boundary and [general release checklist](../../../RELEASE-CHECKLIST.md).

## Out of Scope

No history/UI runtime, scanner semantics, package installs, production retries, output-limit removal, blanket test skipping, weakened assertions, or arbitrary longer deadlines. No full-verification pass is claimed by this issue's creation. No release approval.

## Inputs, Outputs, and Failure Behavior

Input: the recorded failing command/result and unchanged test/runner deadline semantics. Output: controlled diagnostic evidence, causal conclusion or explicitly stated uncertainty, proposed bounded correction if needed, and a distinct owner readiness disposition. Unknown cause remains unknown; it is not asserted to be a harmless machine flake because CI passed. Until disposition, runtime issues stay unstarted; the existing failure remains in the evidence record.

## Required Tests

- Exercise a controlled real helper process and observe startup, working-directory output, execution/deadline, cancellation and exit status; use deterministic isolated inputs.
- If a correction is approved, verify the changed consumer-visible failure path and preserve genuine timeout/cancellation/output-cap behavior. Do not add tests for source text or incidental timing defaults.
- Required eventual commands: relevant targeted verification/process tests and `make verify` from the reviewed clean baseline; record exact versions, platform, commit and exits. A new result complements rather than erases the original failure.

## Acceptance Criteria

- [x] The original command, failure, platform and one-second helper boundary remain linked accurately.
- [x] Diagnostic evidence identifies the failing lifecycle phase; root cause is supported, or remaining uncertainty is explicit with a bounded risk assessment.
- [x] Any correction preserves effective timeout/cancellation and workspace execution boundaries and has actual smoke/regression evidence; no blind timeout increase, retry or skip masks the failure.
- [x] The owner records explicit resolution or risk disposition and the exact remaining verification gate before any H05 runtime issue starts.
- [x] Local/CI outcomes are recorded separately; no known failed run is relabeled passed.

## Owner Review

Owner reviews timing/cancellation safety, causal evidence, any proposed fix, uncertainty and the pre-start decision. Security/architecture effects of any external-process change require a separate owner record. This issue does not grant that change in advance.

## Evidence Required

Sanitized command/transcript, commit/toolchain/platform identity, lifecycle timing, proposed/implemented fix scope if applicable, targeted/full-check results, three-OS CI links as applicable, owner disposition and links from the delivery tracking record. No raw environment, credentials or repository-sensitive paths in public evidence.

## Rollback

Withdraw an unproven correction; retain the original bounded process semantics and recorded failure. Keep runtime starts blocked until owner disposition; do not hide the issue by deleting evidence.

## Final Acceptance

- [x] Owner reviewed every acceptance criterion during review.
- [x] Required verification and CI results were reviewed.
- [x] Owner decision and remaining risk are recorded.
- [x] Evidence links are attached and runtime gate updated.
- [x] Any implementation PR and owner-reviewed merge are linked; investigation-only disposition states that no runtime code changed.

## Investigation evidence — 2026-09-30

The owner's `go` instruction authorizes this bounded investigation after owner-merged PR #428. [Actual lifecycle evidence](../tracking/V05-QA-001-INVESTIGATION.md) records cold helper launches over three seconds, warm/symlink observations, unchanged production-runner workspace/timeout/cancellation smoke and an unresolved output-capture boundary finding (155 bytes captured with a 32-byte cap). The original local full-verification failure remains failed. Exact OS cause is not established; no production or test correction was applied. A darwin-only symlink-fixture correction is proposed for owner authorization, not declared verified. Technical acceptance, separate security/architecture review and runtime pre-start disposition remain pending; no H05 issue starts or acceptance checkbox is completed.

## Additive evidence reconciliation — 2026-10-01

The original Darwin/arm64 `make verify` failure on baseline `325aa23627f8b08ad6b0e6535b2a79a9e968a941` remains a failure: exit 2 at `TestRunExecutesSelectedCommandsInWorkspace`, with the one-second helper deadline. The lifecycle probe localizes the observed delay to after process start and before first output, but does not prove the earlier run's exact phase or an OS/security-subsystem cause.

- The separate 155-byte capture under a 32-byte cap was reproduced to the promoted `bytes.Buffer.ReadFrom` fast path and corrected by owner-merged [PR #431](https://github.com/geoffrey-xiao/deprail/pull/431). Real-child stdout/stderr boundary regressions pass; exact-head [CI run 36720770813](https://github.com/geoffrey-xiao/deprail/actions/runs/36720770813) passed on Linux, macOS and Windows. The earlier #430 local `make verify` failure and its evidence remain preserved in [V05-QA-002-EVIDENCE](../tracking/V05-QA-002-EVIDENCE.md).
- The Darwin-only helper copy was replaced by a symlink to the existing test executable in owner-merged [PR #433](https://github.com/geoffrey-xiao/deprail/pull/433), preserving production runner, deadlines and assertions. Actual fixture probes and corrected-baseline `make verify` passed; exact-head [CI run 36725991235](https://github.com/geoffrey-xiao/deprail/actions/runs/36725991235) passed on all three OSes.
- The owner accepted the residual cold-start root-cause uncertainty and the separate fixture security/architecture assessment in [#432](https://github.com/geoffrey-xiao/deprail/issues/432#issuecomment-5913275072). The integrated projection branch was subsequently owner-merged as [PR #434](https://github.com/geoffrey-xiao/deprail/pull/434); local `make verify` and exact-head [CI run 36740438976](https://github.com/geoffrey-xiao/deprail/actions/runs/36740438976) passed.

These results complement; they do not erase or relabel historical failures. The owner accepted the residual cold-start uncertainty in #432 and separately accepted the #430 capture correction's security/architecture assessment in this conversation; the decision is documented in [V05-QA-002](V05-QA-002-bounded-output-capture.md#additive-owner-securityarchitecture-decision--2026-10-01). The exact owner-approved readiness disposition and verification gate are specified below; this closeout PR records and links them. H05-002/#421 remains Blocked until the owner merges this closeout PR. Every later H05 issue still requires its own DoR, dependency check, targeted behavioral/security evidence, `make verify`, exact-head three-OS CI and owner review. No release approval is inferred.

**Proposed readiness disposition:** accept the bounded residual uncertainty already explicitly accepted in #432; do not claim an exact OS cause. The separate output-bound defect is corrected and its security/architecture assessment is now recorded. Owner merge of this closeout PR is the acceptance action for the combined #419 gate. On merge, H05 runtime work may begin only issue-by-issue after each issue's DoR and dependency check. The original failure and pre-correction failures remain historical.
