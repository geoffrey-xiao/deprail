# V02-015 Baseline Representation and Storage

## Planning Metadata

- Type: feature
- Area: normalization
- Priority: P0
- Risk: R1
- Target version: 0.2.0
- Sprint: S5
- Dependencies: V02-001, V02-010
- Parent epic: EPIC-006

## Definition of Ready

- [ ] Contract inputs, outputs, and failure behavior are confirmed.
- [ ] Dependencies and acceptance evidence are explicit.
- [ ] Compatibility and security boundaries are reviewed.

## Goal

Define a versioned, deterministic baseline representation for trusted scan results.

## Scope

Baseline identity, source scan references, stable findings, completeness, provenance, compatibility, and local storage boundaries.

## Out of Scope

Remote publishing, hosted history, policy decisions, and automatic baseline acceptance.

## Inputs, Outputs, and Failure Behavior

A valid complete baseline is stored and reloadable. Invalid, incomplete, incompatible, or tampered baselines fail explicitly and never become trusted comparison input.

## Required Tests

Schema, round-trip, deterministic serialization, incomplete baseline rejection, provenance retention, and tamper detection tests.

## Acceptance Criteria

- Baselines have versioned schemas and stable identity.
- Incomplete or incompatible baselines are rejected explicitly.
- Stable findings and raw provenance remain available.
- Equivalent baselines serialize identically.

## Human Review

Review compatibility, persistence, and trust-boundary behavior.

## Evidence Required

Schema, examples, round-trip output, failure cases, and compatibility note.
## Final Acceptance

- [ ] Acceptance evidence is linked.
- [ ] CI and human review are complete.
- [ ] Remaining risk is recorded.
