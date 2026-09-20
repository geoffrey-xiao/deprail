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
- Issue: [#255](https://github.com/geoffrey-xiao/deprail/issues/255)
- Parent epic: [#244](https://github.com/geoffrey-xiao/deprail/issues/244)

## Definition of Ready

- [ ] Owner and named reviewer are assigned.
- [ ] v0.3.1 target, sprint, dependencies, and contract mapping are confirmed.
- [ ] Inputs, outputs, failure behavior, and required evidence are reviewed.
- [ ] Implementation remains blocked until release-level Definition of Ready approval.

## Contract mapping

- Requirements: `FUNCTIONAL-REQUIREMENTS.md` UX-002.
- Architecture: `ARCHITECTURE-v0.3.1.md` event contract and dependency direction.
- Evidence: event ordering, cancellation, and dependency-boundary tests.

## Inputs, outputs, and failure behavior

- Inputs: application lifecycle transitions and structured domain results.
- Outputs: renderer-independent presentation events.
- Failure: invalid lifecycle ordering or cancellation state remains explicit and never renders false completion.

## Final acceptance

- [ ] Owner reviewed every criterion.
- [ ] Architecture/security review completed.
- [ ] CI and required evidence links are recorded.
- [ ] Remaining risk and follow-up are documented.
