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

- Issue: [#246](https://github.com/geoffrey-xiao/deprail/issues/246)
## Definition of Ready

- [ ] Owner and named reviewer are assigned.
- [ ] v0.3.1 target, sprint, dependencies, and contract mapping are confirmed.
- [ ] Inputs, outputs, failure behavior, and required evidence are reviewed.
- [ ] Implementation remains blocked until release-level Definition of Ready approval.

## Contract mapping

- Requirements: `FUNCTIONAL-REQUIREMENTS.md` UX-003, UX-004, UX-007, UX-009.
- CLI contract: `CLI-CONTRACT.md` stream behavior.
- Evidence: capability unit tests and cross-platform stream captures.

## Inputs, outputs, and failure behavior

- Inputs: output mode, stream handles, TTY/CI state, color, width, and cancellation capabilities.
- Outputs: explicit renderer capability selection.
- Failure: ambiguous or unavailable capabilities select a safe non-interactive fallback; machine output is never contaminated.

## Final acceptance

- [ ] Owner reviewed every criterion.
- [ ] Architecture/security review completed.
- [ ] CI and required evidence links are recorded.
- [ ] Remaining risk and follow-up are documented.
- Parent epic: [#244](https://github.com/geoffrey-xiao/deprail/issues/244)

