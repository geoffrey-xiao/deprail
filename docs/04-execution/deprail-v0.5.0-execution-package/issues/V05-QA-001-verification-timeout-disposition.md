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

- [ ] The original command, failure, platform and one-second helper boundary remain linked accurately.
- [ ] Diagnostic evidence identifies the failing lifecycle phase; root cause is supported, or remaining uncertainty is explicit with a bounded risk assessment.
- [ ] Any correction preserves effective timeout/cancellation and workspace execution boundaries and has actual smoke/regression evidence; no blind timeout increase, retry or skip masks the failure.
- [ ] The owner records explicit resolution or risk disposition and the exact remaining verification gate before any H05 runtime issue starts.
- [ ] Local/CI outcomes are recorded separately; no known failed run is relabeled passed.

## Owner Review

Owner reviews timing/cancellation safety, causal evidence, any proposed fix, uncertainty and the pre-start decision. Security/architecture effects of any external-process change require a separate owner record. This issue does not grant that change in advance.

## Evidence Required

Sanitized command/transcript, commit/toolchain/platform identity, lifecycle timing, proposed/implemented fix scope if applicable, targeted/full-check results, three-OS CI links as applicable, owner disposition and links from the delivery tracking record. No raw environment, credentials or repository-sensitive paths in public evidence.

## Rollback

Withdraw an unproven correction; retain the original bounded process semantics and recorded failure. Keep runtime starts blocked until owner disposition; do not hide the issue by deleting evidence.

## Final Acceptance

- [ ] Owner reviewed every acceptance criterion during review.
- [ ] Required verification and CI results were reviewed.
- [ ] Owner decision and remaining risk are recorded.
- [ ] Evidence links are attached and runtime gate updated.
- [ ] Any implementation PR and owner-reviewed merge are linked; investigation-only disposition states that no runtime code changed.
