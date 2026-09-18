# V02-008 Add Scan-from-Outside-Root Regression

- Type: test
- Area: test
- Priority: P1
- Risk: R2
- Target version: 0.2.0
- Milestone: v0.2.0
- Sprint: Sprint 1
- Owner: TBD
- Reviewer: TBD
- GitHub Issue: [#128](https://github.com/geoffrey-xiao/deprail/issues/128)
- Parent epic: [EPIC-003 / #114](https://github.com/geoffrey-xiao/deprail/issues/114)
- Dependencies: V02-007

## Definition of Ready

- [ ] Test fixture repositories are identified.
- [ ] Caller and target directory layout is defined.
- [ ] Expected scanner and artifact paths are specified.
- [ ] Owner and reviewer are assigned.

## Goal

Permanently prevent regressions where scanning from outside the target inspects the wrong repository.

## Scope

Add an end-to-end test that launches the CLI from a separate working directory and scans a fixed target fixture.

## Out of Scope

Changing scanner adapter behavior beyond the requested-root contract or adding new fixture ecosystems.

## Inputs, Outputs, and Failure Behavior

The test uses distinguishable target and caller repositories. It fails if findings, workspaces, scanner invocation, or artifacts reflect the caller directory. Containment failures are explicit.

## Required Tests

- CLI scan from outside target.
- Relative and absolute target paths where supported.
- Artifact location assertion.
- Cross-platform path semantics.

## Acceptance Criteria

- [ ] Test fails for caller-CWD regression.
- [ ] Test passes with requested-root behavior.
- [ ] Assertions cover both scan result and artifact placement.
- [ ] Test is deterministic and fixture-local.

## Human Review

Review whether the test proves consumer-visible root behavior rather than implementation details.

## Evidence Required

Test output, fixture layout, and platform result.

## Final Acceptance

- [ ] Owner reviewed every criterion during PR review.
- [ ] CI and regression evidence were reviewed.
- [ ] Remaining risk is recorded.
- [ ] PR review and merge evidence are linked.
