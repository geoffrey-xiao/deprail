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

## Decision

The required v0.2 CI-guardrail capabilities remain in scope and proceed through their existing implementation issues:

| Capability | Decision | Implementation path |
| --- | --- | --- |
| Baseline representation, storage, and compatibility | Retained for v0.2 | V02-015 |
| Base/head comparison and deterministic diff | Retained for v0.2 | V02-016, V02-018 |
| New/resolved/unchanged classification | Retained for v0.2 | V02-017 |
| Typed policy gates, completeness, and exit behavior | Retained for v0.2 | V02-019, V02-021 |
| Expiring exceptions | Retained for v0.2 | V02-020 |
| SARIF output | Retained for v0.2 | V02-022 |

The following capabilities remain deferred because they require separate contracts, compatibility analysis, security review, and explicit inclusion approval:

| Capability | Decision | Rationale |
| --- | --- | --- |
| Additional scanner ecosystems | Deferred | The v0.2 scope remains limited to the currently supported JavaScript, Python, and Java flows. |
| SBOM/signing implementation | Deferred | No reviewed implementation, platform process, or supply-chain evidence exists yet. |
| Remote publishing/history | Deferred | Requires network, storage, authentication, permissions, and retention decisions outside the current release boundary. |

V02-015 through V02-022 remain implementation-ready according to their individual contracts; this decision does not remove or defer those issues. No PRD or requirements update is required because this decision reconciles the issue package with the existing v0.2 plan.

### Consequences and remaining risk

V0.2 retains its CI-guardrail outcome: deterministic baselines and diffs, policy evaluation, and SARIF output remain required before the v0.2 boundary is complete. Additional scanner ecosystems, SBOM/signing implementation, and remote history remain explicit future work. Each deferred capability requires a new or updated issue contract before implementation.

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
