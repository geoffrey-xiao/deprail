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
