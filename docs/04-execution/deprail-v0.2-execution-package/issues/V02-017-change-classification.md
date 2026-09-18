# V02-017 New Resolved Unchanged Classification

## Planning Metadata

- Type: feature
- Area: normalization
- Priority: P0
- Risk: R1
- Target version: 0.2.0
- Sprint: S5
- Dependencies: V02-016
- Parent epic: EPIC-006

## Definition of Ready

- [ ] Contract inputs, outputs, and failure behavior are confirmed.
- [ ] Dependencies and acceptance evidence are explicit.
- [ ] Compatibility and security boundaries are reviewed.

## Goal

Classify dependency and vulnerability changes for pull-request risk review.

## Scope

New, resolved, unchanged, added dependency, upgraded dependency, removed dependency, and affected-workspace classifications.

## Out of Scope

Severity policy, exception approval, exploitability scoring, or remediation planning.

## Inputs, Outputs, and Failure Behavior

Classification is deterministic for compatible inputs. Unknown or ambiguous identity remains explicit and cannot be classified as safe unchanged work.

## Required Tests

Table-driven classifications, alias and version changes, dependency additions/removals, workspace changes, and unknown identity cases.

## Acceptance Criteria

- Every comparable finding receives an explicit classification.
- New findings cannot be hidden by ordering or alias changes.
- Unknown identity is preserved as uncertainty.
- Classification is stable across equivalent input permutations.

## Human Review

Review consumer-visible meaning and stable identity boundaries.

## Evidence Required

Classification matrix, fixtures, output examples, and property/determinism results.
## Final Acceptance

- [ ] Acceptance evidence is linked.
- [ ] CI and human review are complete.
- [ ] Remaining risk is recorded.
