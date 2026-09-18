# V02-004 Establish Release-Mode Evidence Checklist

- Type: docs
- Area: docs
- Priority: P0
- Risk: R0
- Target version: 0.2.0
- Milestone: v0.2.0
- Sprint: Sprint 0
- Owner: TBD
- Reviewer: TBD
- Dependencies: V02-000 context baseline
- GitHub Issue: [#108](https://github.com/geoffrey-xiao/deprail/issues/108)
- Parent epic: [EPIC-001 / #106](https://github.com/geoffrey-xiao/deprail/issues/106)
- Status: GitHub issue #108; child of EPIC-001

## Definition of Ready

- [ ] Preview, RC, and stable modes are defined.
- [ ] Required artifact, checksum, platform, and manual-scan evidence is listed.
- [ ] SBOM/signing treatment is explicit.
- [ ] Owner and reviewer are assigned.

## Goal

Provide one executable checklist for preview, release-candidate, and stable release decisions.

## Scope

Document version baseline checks, tag rules, CI, artifacts, checksums, smoke tests, manual repository scans, supply-chain gaps, approval, rollback, and final evidence.

## Out of Scope

Implementing SBOM generation, signing, release automation, or artifact workflow behavior.

## Inputs, Outputs, and Failure Behavior

The checklist consumes release metadata and verification outputs. Missing evidence blocks the relevant release mode; it cannot be replaced by an unchecked assertion.

## Required Tests

- Link and command verification.
- Dry review against the v0.2 Master Checklist.
- Preview/RC/stable scenario review.

## Acceptance Criteria

- [ ] Each release mode has explicit entry and exit criteria.
- [ ] Required evidence and owners are named.
- [ ] SBOM/signing gaps cannot be hidden.
- [ ] Rollback and immutable-tag rules are documented.

## Human Review

Release owner reviews operational clarity and remaining manual steps.

## Evidence Required

Reviewed release procedure and checklist examples.

## Final Acceptance

- [ ] Owner reviewed every criterion during PR review.
- [ ] Documentation evidence was reviewed.
- [ ] Remaining risk is recorded.
- [ ] PR review and merge evidence are linked.
