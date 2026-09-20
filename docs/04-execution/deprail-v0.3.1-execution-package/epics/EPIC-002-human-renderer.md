# EPIC-002 Human Renderer and Outcomes

## Outcome
DepRail commands use shared, readable human presentation primitives and truthfully distinguish complete, partial, failed, and cancelled outcomes.

## Scope

- Shared headers, sections, statuses, tables, summaries, and guidance.
- `doctor`, `discover`, `scan`, `diff`, and `fix plan` human surfaces.
- Complete zero-finding, partial, failed, and cancelled summaries.
- Narrow-terminal and monochrome fallbacks.

## Issues

- V031-003: Implement shared human presentation primitives.
- V031-004: Implement truthful outcome summaries and command integration.

## Acceptance

Human output is consistent across the listed commands, useful without color, readable at narrow widths, and never describes an incomplete or failed scan as a successful empty result.

## Exclusions

No JSON schema changes, scanner semantics, mutation, network, or full-screen interactive TUI.

## GitHub tracking

- Parent issue: [#242](https://github.com/geoffrey-xiao/deprail/issues/242)
- Child issues: [#245](https://github.com/geoffrey-xiao/deprail/issues/245), [#256](https://github.com/geoffrey-xiao/deprail/issues/256)
