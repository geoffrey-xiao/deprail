# V031-002 Add Terminal Capability Selection

## Planning metadata

- Type: feature
- Area: cli
- Priority: P1
- Risk: R2
- Epic: EPIC-001
- Dependencies: V031-001

## Goal
Select human, JSON, TTY, non-TTY, color, width, CI, and animation behavior explicitly without changing domain results.

## Acceptance criteria

- [ ] TTY detection is evaluated per approved output stream.
- [ ] JSON mode disables banners, progress, and ANSI on stdout.
- [ ] Non-TTY and CI modes are stable and line-oriented.
- [ ] Color, width, and animation fallbacks are explicit and testable.
- [ ] Cancellation and terminal cleanup capabilities are represented.
- [ ] Linux, macOS, and Windows capability tests cover supported differences.

## Evidence
Capability unit tests, stream contract tests, and cross-platform smoke captures.

## Exclusions
No full-screen TUI, scanner behavior, or machine-schema change.

## GitHub tracking

- Issue: Pending creation
- Parent epic: Pending creation
