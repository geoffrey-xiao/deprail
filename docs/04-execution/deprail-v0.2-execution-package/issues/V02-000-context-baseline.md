# V02-000 Establish v0.2.0 Context Baseline

- Type: decision
- Area: foundation
- Priority: P0
- Risk: R1
- Target version: 0.2.0
- Milestone: v0.2.0
- Sprint: Sprint 0
- GitHub Issue: [#98](https://github.com/geoffrey-xiao/deprail/issues/98)
- Dependencies: v0.1 release decision and current `origin/main` baseline
- Owner: Project owner
- Reviewer: Maintainer

## Definition of Ready

- [x] Value and user impact are stated.
- [x] Scope and explicit exclusions are stated.
- [x] Inputs, outputs, and failure behavior are defined.
- [x] Required tests or smoke scenarios are named.
- [x] Acceptance criteria are observable.
- [x] Owner and reviewer are assigned.
- [x] Dependencies and target version are recorded.

## Goal

Review and establish the v0.2.0 execution context before implementation issues are created.

## Scope

- Review the v0.2.0 PRD, transition contract, architecture baseline, requirements, epics, backlog, sprint plans, Master Checklist, evidence guide, and release procedure.
- Confirm the latest stable, preview, and RC tag baseline, current `origin/main`, and release manifest.
- Confirm the v0.2.0 target version and release-mode assumptions.
- Confirm Sprint 0 scope, dependencies, owners, and reviewers.
- Record unresolved v0.1 release-gate risks and their treatment in v0.2 planning.

## Out of Scope

- Runtime implementation.
- Creating all v0.2 implementation issues.
- Creating new product capabilities without a separate decision.
- Publishing a release or changing schemas, CLI behavior, or adapters.

## Inputs, Outputs, and Failure Behavior

Inputs are the local v0.2 execution package, v0.1 contracts and retrospective, repository tags, `origin/main`, and the release manifest. The output is a reviewed context decision with linked evidence and an explicit list of blocked or deferred work.

If the v0.1 baseline, target version, or package scope is not accepted, this issue remains open or moves to `Blocked`; no implementation issue is treated as ready.

## Required Tests

- Review all linked v0.2 package documents for internal consistency.
- Verify the current tag, `origin/main`, and release manifest baseline.
- Verify the v0.2 Master Checklist and release procedure contain the required gates.
- Record review notes and remaining risks in the GitHub issue.

## Human Review

Project owner reviews the package and confirms whether implementation issues may be created.

## Acceptance Criteria

- [ ] The v0.2 PRD and transition contract are reviewed.
- [ ] The v0.2 target version and release mode are confirmed.
- [ ] The current repository, tag, and manifest baseline are recorded.
- [ ] Sprint 0 scope, dependencies, owners, and reviewers are confirmed.
- [ ] v0.1 release-gate risks are explicitly accepted, deferred, or assigned.
- [ ] The v0.2 package is approved as the planning source of truth.

## Evidence Required

- Links to the reviewed v0.2 package documents.
- Tag, `origin/main`, and manifest baseline values.
- Review decision and remaining-risk record.
- Updated v0.2 Master Checklist or linked decision notes.

## Final Acceptance

- [ ] Owner reviewed every acceptance criterion during review.
- [ ] Required verification and evidence were reviewed.
- [ ] Remaining risk and deferred work are recorded.
- [ ] The reviewed context decision authorizes creation of implementation issues.
