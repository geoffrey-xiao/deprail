# V031-005 Wire Real Lifecycle Progress

## Planning metadata

- Type: feature
- Area: cli
- Priority: P1
- Risk: R2
- Epic: EPIC-003
- Dependencies: V031-001, V031-002, V031-003

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

- Issue: Pending creation
- Parent epic: Pending creation
