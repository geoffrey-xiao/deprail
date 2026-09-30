# V05-QA-002 Enforce Bounded Subprocess Output Capture
- GitHub Issue: [#430](https://github.com/geoffrey-xiao/deprail/issues/430).
- GitHub parent: [#418](https://github.com/geoffrey-xiao/deprail/issues/418).

## Planning metadata

- Type: `bug`
- Area: `foundation`
- Priority: `P0`
- Risk: `R3`
- Target version: `0.5.0`
- Milestone: `v0.5.0`
- Sprint: Sprint 4 (current live delivery tracking stage).
- Owner: `@geoffrey-xiao`
- Reviewer: `@geoffrey-xiao`; external reviewer optional under ADR-0005.
- Dependencies: owner-merged investigation PR #429 / closed #419; parent delivery epic #418. Distinct outcome from diagnostic-only #419.
- Blocked reason: None for this bounded correction; H05 runtime readiness remains separately gated.

## Definition of Ready

- [x] Value and user impact are stated.
- [x] Scope and explicit exclusions are stated.
- [x] Inputs, outputs, and failure behavior are defined.
- [x] Required tests and real smoke scenarios are named.
- [x] Acceptance criteria are observable.
- [x] Owner reviewer is assigned; external reviewer optional.
- [x] Dependencies and target version are recorded.

The owner's “请继续我们的任务” authorizes the announced bounded output-capture repair, not merging or releasing. Product Security and Privacy requires bounded hostile tool output; architecture process/testing boundaries, roadmap cross-release invariants, release plan §14/§15 acceptance gates, PROC-001 and the accepted #419 investigation map this bug to the existing security contract. No new product capability, deadline, error code or API is introduced.

## Goal

Ensure stdout and stderr retained by the real process runner never exceed Request.OutputCap per stream, including os/exec pipe-copy fast paths. Preserve explicit SCANNER_OUTPUT_LIMIT failure instead of returning oversized or empty successful capture.

## Scope

Own internal/process/runner.go and runner_test.go, with documentation/evidence in the v0.5 execution package. Confirm the promoted bytes.Buffer.ReadFrom bypass through a controlled pipe-copy probe; remove the bypass by keeping the backing buffer as a named field, routing pipe writes through the existing bounded Write. No new allocation or copy layer. Add real-child regression cases for stdout/stderr, below/at/over cap, single large write and repeated chunks, plus unaffected small output and existing timeout/cancellation/nonzero-exit behavior.

## Out of Scope

No darwin symlink-fixture correction, arbitrary longer deadlines, retries, discarded assertions, process-tree redesign, scanner/parser/schema/CLI contract changes, dependencies, history/UI runtime, release or automatic merge. Do not relabel the original make verify exit 2 as passed. Investigation records remain historical.

## Inputs, Outputs, and Failure Behavior

Input: unchanged process.Request with positive output cap, explicit argument array, cwd/environment and bounded context. Output: exact retained prefix of each stream, at most OutputCap bytes, and existing error classification/precedence. Preserve current conservative at-cap SCANNER_OUTPUT_LIMIT classification; do not silently redefine the limit threshold. Below-cap successful output remains byte-identical. Timeout/cancellation remain effective. Buffer capacity may include bounded allocator slack and fixed pipe-copy overhead; no input-size-dependent unbounded capture.

## Required Tests

- Real helper child emits deterministic stdout/stderr without a testing-framework trailer; below/at/over cap and chunked cases assert exact retained prefixes and failure behavior.
- Controlled io.Copy pipe probe establishes the fast-path mechanism; throwaway, not a source/implementation-detail permanent test.
- Existing process timeout/cancellation/nonzero/argument-array tests; targeted process suite, race check and make verify on this corrected baseline. Record any unchanged known failure separately rather than retry it away.
- Actual deprail CLI invokes a controlled scanner; noisy scanner must produce incomplete/failed result rather than a successful empty scan, while ordinary controlled scan remains compatible.
- Three-OS CI; no cross-platform runtime evidence invented from local compilation.

## Acceptance Criteria

- [x] Real stdout/stderr captures retain exact deterministic prefixes of at most OutputCap bytes for all oversized/chunked cases.
- [x] Promoted fast-path root cause is supported by actual pre-fix observation; same consumer regression fails before and passes after the correction.
- [x] Below-cap output and current conservative at-cap/error precedence are preserved; no deadline, retry, API or error-code change.
- [x] Actual runner/CLI smoke, targeted/full verification and three-OS CI outcomes are separately linked accurately, including any original unresolved local timeout.
- [x] Owner performs technical and separate security/architecture review of the correction and records the remaining runtime-start gate; no H05 readiness inferred from this issue's creation.

## Owner Review

Review actual bounds, both streams, allocation/copy behavior, timeout/cancellation and compatibility evidence. The owner may provide both technical and security/architecture decisions; they remain distinct. No runtime readiness disposition or merge authorization is inferred.

## Evidence Required

Baseline/commit/toolchain/platform; exact before/after bytes and commands; causal pipe-copy observation; real-child regression and CLI smoke outputs/exit codes; local verification and CI separately; immutable source/evidence links; owner review and residual risk. No sensitive environment or credentials.

## Rollback

Revert the correction only through a reviewed change, retaining the vulnerability/failing-regression evidence and runtime gate. No data migration or user repository modification occurs.

## Final Acceptance

- [x] Owner reviewed each acceptance criterion.
- [x] Required local/CI results and remaining risk were reviewed.
- [x] Separate owner security/architecture decision is recorded.
- [x] Evidence and owner-reviewed merge are linked.
- [x] Runtime-start disposition is explicit; unrelated #419 diagnostic acceptance is not substituted.

## Contract links

[Product security/privacy](../../../01-product/deprail-product-design-v1-ai.md#security-and-privacy), [architecture testing/process boundary](../../../02-architecture/deprail-architecture-and-tech-stack-v1.md#10-test-architecture), [roadmap invariants](../../../03-planning/deprail-roadmap-v1.md#cross-release-invariants), [release plan §14](../../../03-planning/deprail-development-plan-v0.5.0.md#14-detailed-contract-reconciliation-for-current-owner-review), [PROC-001](../../deprail-v0.1-execution-package/issues/PROC-001-safe-process-runner.md), [#419 investigation](../tracking/V05-QA-001-INVESTIGATION.md#separate-output-capture-finding--unresolved).

## Correction evidence

[Actual before/after and runner/CLI evidence](../tracking/V05-QA-002-EVIDENCE.md) confirms the copy fast-path cause, six failing-before cases and ten passing-after boundaries, seven real-runner scenarios and ordinary/noisy/race-instrumented CLI behavior. Local `make verify` still exits 2 on the unchanged copied-helper timeout; no retry or helper fix is included. Owner review, separate security/architecture decision and merge remain pending. This correction does not authorize H05 runtime kickoff.

## Additive owner security/architecture decision — 2026-10-01

The owner separately accepted the security/architecture assessment of owner-merged [PR #431](https://github.com/geoffrey-xiao/deprail/pull/431) in the current conversation. The fix removes anonymous `*bytes.Buffer` embedding so `io.Copy` cannot select a promoted `ReadFrom` fast path around the bounded `Write`; it preserves the existing byte cap, error classification, deadlines, cancellation and APIs. Real-child tests exercise stdout and stderr at below/at/over cap and large single/chunked writes; actual runner, normal/noisy CLI and race-instrumented CLI smoke are linked above. Exact-head three-OS CI passed.

Accepted residual risk: this bounds retained stream bytes, not total child CPU/output or allocator capacity; the pre-existing helper startup uncertainty is separate and remains under #419. The separate owner decision does not itself authorize H05 runtime work. The global start gate remains tracked in #419.
