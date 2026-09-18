# V02-009 Update Report Examples and Empty Arrays

- Type: bug
- Area: normalization
- Priority: P0
- Risk: R1
- Target version: 0.2.0
- Milestone: v0.2.0
- Sprint: Sprint 0
- Owner: TBD
- Reviewer: TBD
- GitHub Issue: [#120](https://github.com/geoffrey-xiao/deprail/issues/120)
- Parent epic: [EPIC-003 / #114](https://github.com/geoffrey-xiao/deprail/issues/114)
- Dependencies: V02-001

## Definition of Ready

- [ ] Current schema and examples are enumerated.
- [ ] Empty complete, partial, failed, and findings examples are selected.
- [ ] Consumer compatibility impact is documented.
- [ ] Owner and reviewer are assigned.

## Goal

Align versioned report examples and schema fixtures with the required empty-array serialization contract.

## Scope

Update stable examples, golden files, schema fixtures, and documentation that currently imply `null` for empty collections.

## Out of Scope

Changing report fields, completeness semantics, or finding normalization beyond empty collection representation.

## Inputs, Outputs, and Failure Behavior

Examples validate against the current schema. Any schema/example mismatch fails verification and is corrected before merge; examples never claim a failed scan is safe.

## Required Tests

- Validate every changed JSON example.
- Compare complete, partial, and failed examples.
- Verify deterministic ordering.

## Acceptance Criteria

- [ ] All current examples use arrays for empty collections.
- [ ] Examples cover complete and incomplete outcomes.
- [ ] Documentation matches runtime behavior.
- [ ] No broad unrelated golden refresh is included.

## Human Review

Schema and compatibility review is required for generated or golden-file changes.

## Evidence Required

Schema validator output, changed-example inventory, and compatibility statement.

## Final Acceptance

- [ ] Owner reviewed every criterion during PR review.
- [ ] CI and schema evidence were reviewed.
- [ ] Remaining risk is recorded.
- [ ] PR review and merge evidence are linked.
