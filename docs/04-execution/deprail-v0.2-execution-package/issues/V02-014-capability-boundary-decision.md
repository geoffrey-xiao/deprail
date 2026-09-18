# V02-014 Decide v0.2 Candidate Capability Boundary

- Type: decision
- Area: foundation
- Priority: P1
- Risk: R1
- Target version: 0.2.0
- Milestone: v0.2.0
- Sprint: Sprint 0
- Owner: Project owner
- Reviewer: Maintainer
- GitHub Issue: [#126](https://github.com/geoffrey-xiao/deprail/issues/126)
- Parent epic: [EPIC-005 / #116](https://github.com/geoffrey-xiao/deprail/issues/116)
- Dependencies: V02-000 context baseline

## Definition of Ready

- [ ] Candidate capabilities are listed with user value.
- [ ] Compatibility, security, operational, and schedule trade-offs are identified.
- [ ] Non-selection and deferral consequences are understood.
- [ ] Owner and reviewer are assigned.

## Goal

Decide which candidate capabilities, if any, may enter v0.2 implementation without expanding scope implicitly.

## Scope

Evaluate baseline comparison, policy gates, SARIF, additional scanner ecosystems, SBOM/signing, and remote publishing/history. Select, defer, or reject each candidate.

## Out of Scope

Implementing any candidate capability before the decision and contract updates are approved.

## Inputs, Outputs, and Failure Behavior

The output is a written scope decision with compatibility and risk treatment. If evidence is insufficient, the capability remains deferred and no implementation issue becomes Ready.

## Required Tests

- Decision record review.
- Compatibility and security impact review.
- PRD and requirements consistency check after selection.

## Acceptance Criteria

- [ ] Every candidate is selected, deferred, or rejected.
- [ ] Selected capabilities have owners, dependencies, and acceptance paths.
- [ ] Deferred capabilities remain explicitly out of scope.
- [ ] PRD and requirements are updated if scope changes.

## Human Review

Project owner and maintainer review the trade-offs and final boundary.

## Evidence Required

Decision record, updated PRD/requirements if needed, and remaining-risk statement.

## Final Acceptance

- [ ] Owner reviewed every criterion during PR review.
- [ ] Required human reviewers approved.
- [ ] Evidence and remaining risk are recorded.
- [ ] PR review and merge evidence are linked.
