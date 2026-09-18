# V02-005 Add Real OSV-Scanner v2 Fixtures

- Type: test
- Area: adapter
- Priority: P0
- Risk: R1
- Target version: 0.2.0
- Milestone: v0.2.0
- Sprint: Sprint 0
- Owner: TBD
- Reviewer: TBD
- Dependencies: None
- Status: Local planning; GitHub issue not created

## Definition of Ready

- [ ] Supported OSV-Scanner v2 version range is recorded.
- [ ] Real output samples are obtained from supported commands.
- [ ] Fixture provenance and license safety are documented.
- [ ] Owner and reviewer are assigned.

## Goal

Make the OSV-Scanner v2 adapter contract reproducible from checked-in raw output.

## Scope

Add fixtures for complete empty output, vulnerability findings, severity arrays/database severity, affected ranges, fixed events, and relevant metadata. Document fixture origin and normalization expectations.

## Out of Scope

Supporting OSV-Scanner v1, installing scanners automatically, or changing unrelated ecosystems.

## Inputs, Outputs, and Failure Behavior

Fixtures represent real supported v2 output. Fixture parse failure fails the contract test and cannot be interpreted as an empty scan.

## Required Tests

- Raw fixture parsing.
- Findings and fixed-version normalization.
- Empty result parsing.
- Malformed and incompatible fixture rejection.

## Acceptance Criteria

- [ ] At least one real v2 empty fixture exists.
- [ ] At least one real v2 findings fixture exists.
- [ ] Fixed versions and severity provenance are asserted.
- [ ] Fixture provenance is recorded.

## Human Review

Adapter maintainer reviews fidelity to the supported scanner version.

## Evidence Required

Raw fixture files, provenance note, parser test output, and expected normalized result.

## Final Acceptance

- [ ] Owner reviewed every criterion during PR review.
- [ ] Adapter review and CI evidence were completed.
- [ ] Remaining risk is recorded.
- [ ] PR review and merge evidence are linked.
