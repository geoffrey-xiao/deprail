# V031-001 Define Structured Presentation Events

## Planning metadata

- Type: feature
- Area: cli
- Priority: P1
- Risk: R2
- Epic: EPIC-001
- Dependencies: v0.3.1 execution package and Definition of Ready

## Goal
Define deterministic renderer-independent lifecycle events for CLI application progress and final reporting.

## Acceptance criteria

- [ ] Event vocabulary covers validation, discovery, planning, workspace scans, artifacts, normalization, report readiness, and cancellation.
- [ ] Event fields have explicit lifecycle, identity, ordering, and terminal-state semantics.
- [ ] Events contain no ANSI, spinner, terminal-control, or renderer-specific data.
- [ ] Application/domain packages remain independent of terminal libraries.
- [ ] Tests prove valid ordering, cancellation, repeated workspace behavior, and terminal-state invariants.

## Evidence
Event contract documentation, unit/contract tests, and architecture review of dependency direction.

## Exclusions
No visual styling, scanner changes, JSON schema changes, or repository mutation.

## GitHub tracking

- Issue: Pending creation
- Parent epic: Pending creation
