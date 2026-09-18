# V02-001 Normalize Empty Report Collections

- Type: bug
- Area: normalization
- Priority: P0
- Risk: R1
- Target version: 0.2.0
- Milestone: v0.2.0
- Sprint: Sprint 0
- Owner: TBD
- Reviewer: TBD
- GitHub Issue: [#112](https://github.com/geoffrey-xiao/deprail/issues/112)
- Parent epic: [EPIC-001 / #106](https://github.com/geoffrey-xiao/deprail/issues/106)
- Dependencies: None
- Status: GitHub issue #112; child of EPIC-001

## Definition of Ready

- [ ] Consumer-visible JSON behavior is confirmed.
- [ ] Current schema and serializers are identified.
- [ ] Empty, non-empty, partial, and failed cases are named.
- [ ] Owner and reviewer are assigned before implementation.

## Goal

Serialize empty report collections as JSON arrays instead of `null`.

## Scope

Initialize findings, errors, workspace, diagnostic, and artifact collections wherever the stable contract defines arrays. Update schema examples and focused contract tests.

## Out of Scope

Changing finding meaning, completeness rules, stable keys, or unrelated schema fields.

## Inputs, Outputs, and Failure Behavior

Given a complete scan with no findings, output contains `[]` for empty collections. Non-empty and partial/failed reports retain their existing values. Serialization or schema failure remains an error and never becomes a successful scan.

## Required Tests

- Complete empty scan JSON contract test.
- Non-empty findings and errors regression test.
- Partial and failed report serialization test.
- Schema validation and deterministic comparison.

## Acceptance Criteria

- [ ] Empty collections serialize as `[]`.
- [ ] Existing non-empty and error collections remain semantically unchanged.
- [ ] Current schema examples validate.
- [ ] No stable key depends on collection initialization.

## Human Review

Review consumer compatibility and every changed report constructor/serializer.

## Evidence Required

Focused test output, schema validation, before/after JSON examples, and compatibility statement.

## Final Acceptance

- [ ] Owner reviewed every criterion during PR review.
- [ ] CI and evidence were reviewed.
- [ ] Remaining risk is recorded.
- [ ] PR review and merge evidence are linked.
