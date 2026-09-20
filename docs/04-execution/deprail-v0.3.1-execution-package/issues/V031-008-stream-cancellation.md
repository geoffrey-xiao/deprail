# V031-008 Verify Non-TTY JSON and Cancellation Streams

## Planning metadata

- Type: test
- Area: cli
- Priority: P1
- Risk: R2
- Epic: EPIC-003
- Dependencies: V031-002, V031-004, V031-005, V031-006, V031-007

## Goal
Prove stream separation, JSON purity, stable redirected output, and terminal cleanup under cancellation and failure.

## Acceptance criteria

- [ ] JSON stdout parses independently of stderr for all covered commands.
- [ ] Progress, diagnostics, and guidance never contaminate JSON stdout.
- [ ] Non-TTY output has no animation, ANSI, or carriage-return updates.
- [ ] Cancellation restores terminal state and preserves truthful status.
- [ ] Timeouts, missing tools, malformed output, and partial failures remain errors.

## Evidence
Contract tests, controlled process tests, redirected captures, JSON parsing, and cancellation smoke evidence.

## Exclusions
No new error codes, scanner behavior, or repository mutation.

## GitHub tracking
- Issue: [#247](https://github.com/geoffrey-xiao/deprail/issues/247)
- Parent epic: [#243](https://github.com/geoffrey-xiao/deprail/issues/243)
## Definition of Ready

- [ ] Owner and named reviewer are assigned.
- [ ] v0.3.1 target, sprint, dependencies, and contract mapping are confirmed.
- [ ] Inputs, outputs, failure behavior, and required evidence are reviewed.
- [ ] Implementation remains blocked until release-level Definition of Ready approval.

## Contract mapping

- Requirements: `FUNCTIONAL-REQUIREMENTS.md` UX-004, UX-007, UX-012.
- Error model: `ERROR-MODEL.md` failure and cancellation semantics.
- Evidence: stream contract, JSON parsing, timeout, and cancellation tests.

## Inputs, outputs, and failure behavior

- Inputs: controlled CLI processes, stream handles, JSON reports, and cancellation signals.
- Outputs: separated stdout/stderr with truthful exit and cancellation behavior.
- Failure: timeout, malformed output, missing tools, and partial failures remain explicit errors.

## Final acceptance

- [ ] Owner reviewed every criterion.
- [ ] Architecture/security review completed.
- [ ] CI and required evidence links are recorded.
- [ ] Remaining risk and follow-up are documented.
