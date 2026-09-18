# V02-003 Inject Truthful Build Version Identity

- Type: feature
- Area: cli
- Priority: P1
- Risk: R1
- Target version: 0.2.0
- Milestone: v0.2.0
- Sprint: Sprint 1
- Owner: TBD
- Reviewer: TBD
- Dependencies: None
- GitHub Issue: [#107](https://github.com/geoffrey-xiao/deprail/issues/107)
- Parent epic: [EPIC-001 / #106](https://github.com/geoffrey-xiao/deprail/issues/106)
- Status: GitHub issue #107; child of EPIC-001

## Definition of Ready

- [ ] Release, preview, and development identity rules are agreed.
- [ ] Existing version output locations are identified.
- [ ] Linker/build metadata strategy is selected.
- [ ] Owner and reviewer are assigned.

## Goal

Make CLI version output identify the actual release tag and source commit without falsely identifying development builds as stable releases.

## Scope

Add build metadata injection and expose it through the documented doctor/version output. Define behavior when Git metadata is unavailable.

## Out of Scope

Changing finding keys, report semantics, release automation, or adding a remote version service.

## Inputs, Outputs, and Failure Behavior

Release artifacts report the release version, tag, and commit when available. Development builds report a development marker. Missing metadata uses a truthful unknown/development value and does not fail otherwise valid scans.

## Required Tests

- Release metadata smoke test.
- Development/unknown metadata test.
- Doctor JSON and terminal output tests.
- Stable-key regression test proving build metadata is excluded.

## Acceptance Criteria

- [ ] Release tag and commit are visible in documented output.
- [ ] Development builds do not claim a stable release.
- [ ] Missing metadata is explicit and safe.
- [ ] Stable finding keys remain unchanged.

## Human Review

Review release identity, reproducibility, and compatibility with scripts consuming version output.

## Evidence Required

Build command, output examples, tests, and release compatibility record.

## Final Acceptance

- [ ] Owner reviewed every criterion during PR review.
- [ ] CI and evidence were reviewed.
- [ ] Remaining risk is recorded.
- [ ] PR review and merge evidence are linked.
