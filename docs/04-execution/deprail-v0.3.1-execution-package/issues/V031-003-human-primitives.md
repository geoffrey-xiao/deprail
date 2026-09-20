# V031-003 Implement Shared Human Presentation Primitives

## Planning metadata

- Type: feature
- Area: cli
- Priority: P1
- Risk: R1
- Epic: EPIC-002
- Dependencies: V031-001, V031-002

## Goal
Implement consistent headers, sections, statuses, tables, summaries, guidance, and error presentation for human CLI output.

## Acceptance criteria

- [ ] `doctor`, `discover`, `scan`, `diff`, and `fix plan` use shared primitives.
- [ ] Meaning is readable without color.
- [ ] Narrow-terminal fallback remains readable.
- [ ] Human rendering consumes structured results without changing them.
- [ ] Golden/contract cases avoid pinning incidental ANSI or library internals.

## Evidence
Presenter tests and reviewed golden captures for each command surface, monochrome, and narrow width.

## Exclusions
No JSON schema changes, scanner changes, mutation, or full-screen TUI.

## GitHub tracking

- Issue: [#245](https://github.com/geoffrey-xiao/deprail/issues/245)
- Parent epic: [#242](https://github.com/geoffrey-xiao/deprail/issues/242)
## Definition of Ready

- [ ] Owner and named reviewer are assigned.
- [ ] v0.3.1 target, sprint, dependencies, and contract mapping are confirmed.
- [ ] Inputs, outputs, failure behavior, and required evidence are reviewed.
- [ ] Implementation remains blocked until release-level Definition of Ready approval.

## Contract mapping

- Requirements: `FUNCTIONAL-REQUIREMENTS.md` UX-001 and UX-009.
- Architecture: `ARCHITECTURE-v0.3.1.md` presenter responsibility.
- Evidence: command presenter captures and semantic golden tests.

## Inputs, outputs, and failure behavior

- Inputs: structured domain results and presentation events.
- Outputs: human terminal summaries for `doctor`, `discover`, `scan`, `diff`, and `fix plan`.
- Failure: renderer errors remain explicit and cannot alter domain results or machine output.

## Final acceptance

- [ ] Owner reviewed every criterion.
- [ ] Architecture/security review completed.
- [ ] CI and required evidence links are recorded.
- [ ] Remaining risk and follow-up are documented.
