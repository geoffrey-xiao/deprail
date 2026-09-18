# V02-010 Verify Deterministic Artifacts and Ordering

- Type: test
- Area: test
- Priority: P1
- Risk: R2
- Target version: 0.2.0
- Milestone: v0.2.0
- Sprint: Sprint 1
- Owner: TBD
- Reviewer: TBD
- GitHub Issue: [#124](https://github.com/geoffrey-xiao/deprail/issues/124)
- Parent epic: [EPIC-003 / #114](https://github.com/geoffrey-xiao/deprail/issues/114)
- Dependencies: V02-009

## Definition of Ready

- [ ] Stable ordering keys and allowed run metadata are identified.
- [ ] Artifact naming/digest invariants are listed.
- [ ] Permutation and repeated-run fixtures are selected.
- [ ] Owner and reviewer are assigned.

## Goal

Prove equivalent inputs produce equivalent semantic reports and stable artifact identity.

## Scope

Add permutation/property/golden coverage for findings, evidence, workspaces, diagnostics, and artifact references.

## Out of Scope

Removing explicitly allowed timestamps or run metadata, changing stable-key semantics, or introducing a new persistence layer.

## Inputs, Outputs, and Failure Behavior

Permuting equivalent input order produces identical semantic JSON except declared metadata. Any nondeterministic stable field fails the test and is investigated rather than masked.

## Required Tests

- Input-order permutation.
- Repeated normalization.
- Artifact digest/name stability.
- Cross-platform path serialization.

## Acceptance Criteria

- [ ] Equivalent input permutations compare equal semantically.
- [ ] Stable keys exclude prohibited fields.
- [ ] Artifact references remain reproducible.
- [ ] Tests identify allowed metadata differences explicitly.

## Human Review

Normalization and security reviewers inspect test boundaries and equality rules.

## Evidence Required

Golden/property output, allowed-difference definition, and platform results.

## Final Acceptance

- [ ] Owner reviewed every criterion during PR review.
- [ ] CI and determinism evidence were reviewed.
- [ ] Remaining risk is recorded.
- [ ] PR review and merge evidence are linked.
