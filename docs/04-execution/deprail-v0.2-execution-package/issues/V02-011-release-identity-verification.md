# V02-011 Verify Release Identity Across Artifacts

- Type: test
- Area: cli
- Priority: P1
- Risk: R1
- Target version: 0.2.0
- Milestone: v0.2.0
- Sprint: Sprint 1
- Owner: TBD
- Reviewer: TBD
- Dependencies: V02-003
- Status: Local planning; GitHub issue not created

## Definition of Ready

- [ ] Version identity source and output fields are defined.
- [ ] Artifact matrix and release tag scenarios are listed.
- [ ] Development-build behavior is defined.
- [ ] Owner and reviewer are assigned.

## Goal

Verify that release tag, CLI output, artifact metadata, and source commit identify the same build.

## Scope

Build the documented release matrix and compare tag/commit identity across artifacts and `doctor --format json` output.

## Out of Scope

Implementing version injection itself, signing, SBOM generation, or changing artifact names without a compatibility decision.

## Inputs, Outputs, and Failure Behavior

A mismatch fails release verification and blocks publication. Missing development metadata is reported as development/unknown, never as the stable release version.

## Required Tests

- Release artifact identity smoke.
- Development build identity smoke.
- Cross-platform artifact comparison.
- Manifest/tag/CLI consistency check.

## Acceptance Criteria

- [ ] Release artifacts report the intended tag and commit.
- [ ] CLI output and artifacts agree.
- [ ] Development builds are clearly non-release.
- [ ] Mismatch blocks release evidence acceptance.

## Human Review

Release reviewer checks identity evidence and rollback implications.

## Evidence Required

Build commands, output transcripts, artifact inventory, and comparison result.

## Final Acceptance

- [ ] Owner reviewed every criterion during PR review.
- [ ] CI and release evidence were reviewed.
- [ ] Remaining risk is recorded.
- [ ] PR review and merge evidence are linked.
