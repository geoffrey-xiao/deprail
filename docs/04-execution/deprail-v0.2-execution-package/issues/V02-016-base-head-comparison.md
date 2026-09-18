# V02-016 Base and Head Scan Comparison

- Type: feature
- Area: normalization
- Priority: P0
- Risk: R2
- Target version: 0.2.0
- Sprint: S5
- Dependencies: V02-015
- Parent epic: EPIC-006

## Goal

Compare trusted base and head results without depending on input ordering or volatile fields.

## Scope

Load two compatible scan results, compare stable finding keys and dependency identities, and preserve completeness and provenance differences.

## Out of Scope

Policy decisions, exception handling, SARIF, or source mutation.

## Inputs, Outputs, and Failure Behavior

Compatible inputs produce a deterministic comparison. Missing, incompatible, incomplete, or malformed inputs produce explicit errors and no safe diff.

## Required Tests

Equivalent input permutations, added/resolved findings, changed dependency versions, incomplete inputs, incompatible schemas, and malformed data.

## Acceptance Criteria

- Comparison is order-independent.
- New and resolved findings are retained with evidence.
- Incomplete inputs cannot produce a safe complete comparison.
- Volatile fields do not change comparison identity.

## Human Review

Review stable-key and compatibility semantics.

## Evidence Required

Fixtures, comparison JSON, failure transcripts, and determinism results.
