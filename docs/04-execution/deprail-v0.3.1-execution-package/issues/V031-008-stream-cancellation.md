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

- Issue: Pending creation
- Parent epic: Pending creation
