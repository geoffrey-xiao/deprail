# V031-005 Wire Real Lifecycle Progress

## Planning metadata

- Type: feature
- Area: cli
- Priority: P1
- Dependencies: V031-001, V031-002, V031-003, V031-007
- Epic: EPIC-003

## Goal
Show concise TTY stderr progress from real application lifecycle events without invented work or completion claims.

## Acceptance criteria

- [ ] Progress covers validation, discovery, planning, scanning, artifacts, normalization, and final report boundaries where events exist.
- [ ] Progress is TTY-only and never writes machine stdout.
- [ ] No artificial delays, percentages, durations, or fabricated work are emitted.
- [ ] Concurrent workspace output is serialized safely.
- [ ] Cancellation and failure terminate or finalize progress cleanly.

## Evidence
Controlled scanner integration tests, interactive TTY captures, redirected-stream tests, and cancellation evidence.

## Exclusions
No scanner invocation changes, full-screen TUI, telemetry, or network behavior.

## GitHub tracking
- Issue: [#252](https://github.com/geoffrey-xiao/deprail/issues/252)
- Parent epic: [#243](https://github.com/geoffrey-xiao/deprail/issues/243)
## Definition of Ready

- [ ] Owner and named reviewer are assigned.
- [ ] v0.3.1 target, sprint, dependencies, and contract mapping are confirmed.
- [ ] Inputs, outputs, failure behavior, and required evidence are reviewed.
- [ ] Implementation remains blocked until release-level Definition of Ready approval.

## Contract mapping

- Requirements: `FUNCTIONAL-REQUIREMENTS.md` UX-003 and UX-004.
- CLI contract: `CLI-CONTRACT.md` stream behavior.
- Evidence: controlled-process, TTY, redirected-output, and cancellation captures.

## Inputs, outputs, and failure behavior

- Inputs: validated lifecycle events and approved terminal capabilities.
- Outputs: concise TTY stderr progress.
- Failure: missing events, scanner failure, timeout, or cancellation closes progress truthfully and never fabricates completion.

## Final acceptance

- [ ] Owner reviewed every criterion.
- [ ] Architecture/security review completed.
- [ ] CI and required evidence links are recorded.
- [ ] Remaining risk and follow-up are documented.
