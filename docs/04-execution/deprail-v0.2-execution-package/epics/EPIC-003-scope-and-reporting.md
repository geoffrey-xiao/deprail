# EPIC-003 Scope and Reporting

## Outcome

Ensure requested-root scope, completeness, deterministic normalization, and machine output remain trustworthy under repeated use.

## Scope

- Outside-root scan regression coverage.
- Empty array serialization and schema examples.
- Deterministic report and artifact ordering.
- Path and artifact containment evidence.

## Acceptance

Equivalent inputs produce equivalent semantic JSON; incomplete or failed scope remains visible and cannot be described as safe.

## Issues

- `V02-008` Add scan-from-outside-root end-to-end coverage.
- `V02-009` Update report schema examples and empty collection behavior.
- `V02-010` Add deterministic ordering and artifact containment regression coverage.
