# V05-QA-004: Synchronize Timed Apply E2E With Mutation Startup

GitHub issue: [#454](https://github.com/geoffrey-xiao/deprail/issues/454).

## Planning metadata

- Release/milestone: v0.5.0; Sprint4; owner/reviewer @geoffrey-xiao; external review optional.
- Type:test; area:test; priority:P0; risk:R1 (test-only correction).
- Parent verification #427; existing cancellation/evidence contract from v0.4 apply #368 and v0.5 TEST-STRATEGY/FR-504 cross-platform compatibility. No new release capability or CLI contract.
- Dependencies: synchronized main d6362f3, existing controlled package-manager helper and real isolated-worktree E2E; unrelated default-console #452/#453 remains separate.
- Blocked reason: none after own preparation PR #445 returned to Draft, preserving the two-active-review-PR limit.

## Definition of Ready

- [x] Observed Windows failures: cancellation in runs37551712270/37564839801; timeout in37569906460 at bc4d732; unexpected end of JSON input before later Windows build/smoke.
- [x] Source cancellation unconditionally fires after500ms; timeout starts a1s budget before worktree/mutation startup. Neither asserts the intended phase was reached.
- [x] Controlled Darwin preflight-cancellation smoke yielded exit3, stdoutBytes0, PATH_OUTSIDE_ROOT and exact same JSON decode error. This proves the mechanism; original Windows stderr was not logged, so exact failed preflight phase remains unknown.
- [x] Scope, observable assertions, failure bounds, owner and evidence explicit below. No rerun/classification/weakening substituted for investigation.

## Goal and scope

Make cancellation and timeout E2E exercise an actually running controlled mutation child, independent of preflight speed. Child publishes a private ready marker before sleeping. Parent waits for readiness or bounded startup failure/early operation completion, then triggers cancellation or an actual standard-library deadline. Preserve distinct cancelled versus partial outcome, exit3, nonempty durable evidence and successful cleanup. Test-context deadline propagation uses the existing Context port; no product test hook or context contract change.

Improve JSON decode failure diagnostics with exit/stderr so preflight failures are no longer reported as unexplained empty JSON. Do not change product JSON/error behavior or suppress decode failure. No retries, long fixed sleeps, broader helper refactoring, timeout-only increase, default artifact fix, production changes or release acceptance.

## Inputs, outputs and failure behavior

Real isolated fixture repository/plan/approval/worktree and controlled native child. Output: genuine apply result/evidence and cleanup. Missing marker, child exit before readiness or30s readiness expiry fails with phase/exit diagnostics and cancels/drains the operation; it is never counted as successful cancellation. Helper startup signal is test-only and cannot mutate user repositories.

## Required tests and acceptance

- [ ] Both timed cases confirm child readiness before triggering termination; cancellation returns cancelled/exit3, expired deadline returns partial/exit3.
- [ ] Existing evidence and cleanup assertions preserved; additionally validate persisted evidence outcome, unchanged caller state and removed isolated worktree.
- [ ] Early preflight failure remains an error with exit/stderr diagnostics, not accepted malformed/empty JSON.
- [ ] Focused actual-process E2E, race check and full make verify/make test-integration pass; exact-head Linux/macOS/Windows CI and native package smoke recorded.
- [ ] Historical failures remain linked; explicit explanation distinguishes observed mechanism from unobserved original Windows phase.

## Human review and evidence

Owner reviews readiness ordering, genuine deadline propagation, startup bounds and unchanged assertions. Evidence records source/head, commands, timed runtime outcome/evidence/cleanup and CI. No permanent tests of source text or incidental wording. H05-008 prerequisite/security/release decisions remain outstanding.

[Diagnosis and actual-process verification](../tracking/V05-QA-004-TIMED-APPLY-EVIDENCE.md) preserves failure history, the controlled preflight mechanism, phase/deadline correction, race/full-suite results and remaining release limits.

## Final acceptance

- [ ] Owner criteria/residual risk reviewed.
- [ ] Required verification/CI accepted.
- [ ] Owner-reviewed merge recorded before Project Done/closure.

## Rollback

Revert test/helper synchronization only; no product binary, schema, history, source repository or release-tag changes.
