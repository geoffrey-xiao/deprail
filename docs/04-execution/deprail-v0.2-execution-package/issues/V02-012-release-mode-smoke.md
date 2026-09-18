# V02-012 Add Release-Mode Smoke Procedure

- Type: test
- Area: test
- Priority: P1
- Risk: R1
- Target version: 0.2.0
- Milestone: v0.2.0
- Sprint: Sprint 1
- Owner: TBD
- Reviewer: TBD
- GitHub Issue: [#127](https://github.com/geoffrey-xiao/deprail/issues/127)
- Parent epic: [EPIC-004 / #115](https://github.com/geoffrey-xiao/deprail/issues/115)
- Dependencies: V02-004

## Definition of Ready

- [ ] Preview, RC, and stable entry criteria are available.
- [ ] Artifact and platform smoke commands are defined.
- [ ] Required evidence locations are defined.
- [ ] Owner and reviewer are assigned.

## Goal

Make release-mode verification executable rather than an informal checklist.

## Scope

Exercise the release procedure against preview, RC, and stable scenarios, including clean checkout, pinned toolchain, artifact smoke, checksums, and evidence capture.

## Out of Scope

Publishing a release, changing protected environments, or implementing signing/SBOM.

## Inputs, Outputs, and Failure Behavior

Missing evidence blocks the release mode under review. A preview result cannot be presented as stable acceptance.

## Required Tests

- Clean-checkout `make verify`.
- Platform artifact smoke.
- Checksum review.
- Version/tag/commit identity.
- Failure scenario review.

## Acceptance Criteria

- [ ] Procedure can be followed without undocumented steps.
- [ ] Preview, RC, and stable differences are clear.
- [ ] Evidence record is reproducible.
- [ ] Failure and rollback instructions are usable.

## Human Review

Release owner reviews operational completeness.

## Evidence Required

Executed smoke record, checklist review, and identified gaps.

## Final Acceptance

- [ ] Owner reviewed every criterion during PR review.
- [ ] CI and smoke evidence were reviewed.
- [ ] Remaining risk is recorded.
- [ ] PR review and merge evidence are linked.
