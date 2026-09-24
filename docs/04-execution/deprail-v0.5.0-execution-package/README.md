# DepRail v0.5.0 Execution Package

| Attribute | Value |
| --- | --- |
| Release | v0.5.0; preview first |
| Status | Draft for owner and independent architecture/security review; no runtime implementation authorized |
| Preparation issue | [#387](https://github.com/geoffrey-xiao/deprail/issues/387) |
| Release plan | [`../../03-planning/deprail-development-plan-v0.5.0.md`](../../03-planning/deprail-development-plan-v0.5.0.md) |
| Product baseline | [`../../01-product/deprail-product-design-v1-ai.md`](../../01-product/deprail-product-design-v1-ai.md) |
| Architecture baseline | [`../../02-architecture/deprail-architecture-and-tech-stack-v1.md`](../../02-architecture/deprail-architecture-and-tech-stack-v1.md) |
| Roadmap | [`../../03-planning/deprail-roadmap-v1.md`](../../03-planning/deprail-roadmap-v1.md) |
| Predecessor closeout | [#356 v0.4.0-preview.2](https://github.com/geoffrey-xiao/deprail/issues/356), still open |
| Owner | `@geoffrey-xiao` |
| Independent reviewer | To be named |

This package translates the proposed v0.5 release outcome into reviewable UX, API, architecture, data, failure, security, compatibility, and test contracts. It is a draft subordinate to the product design, architecture, roadmap, and v0.5 release plan. Conflicts stop implementation and require a versioned decision; lower-level requirements cannot silently override those sources.

## Source order

1. Product design and architecture.
2. Whole-project roadmap.
3. [`v0.5.0 development plan`](../../03-planning/deprail-development-plan-v0.5.0.md).
4. This package and its requirements.
5. Owner- and reviewer-approved issue contracts.
6. Implementation, verification, and release evidence.

## Package documents

- [`PRD-v0.5.md`](PRD-v0.5.md) — outcome, scope, users, exclusions, proposed acceptance criteria.
- [`ARCHITECTURE-v0.5.md`](ARCHITECTURE-v0.5.md) — component boundaries and local system context.
- [ADR-0004](../../adr/ADR-0004-local-scan-history.md) — owner-approved SQLite proposal; independent review and DoR remain outstanding.
- [`UX-DESIGN.md`](UX-DESIGN.md) — workflows, information architecture, shadcn/ui-inspired direction, interaction and accessibility states.
- [`API-DESIGN.md`](API-DESIGN.md) — proposed local REST resource and operation contract; OpenAPI artifact remains a design-gate deliverable.
- [`requirements/FUNCTIONAL-REQUIREMENTS.md`](requirements/FUNCTIONAL-REQUIREMENTS.md) — observable behaviors and release gates.
- [`requirements/FAILURE-AND-DATA-CONTRACT.md`](requirements/FAILURE-AND-DATA-CONTRACT.md) — proposed persistence and failure semantics.
- [`requirements/ERROR-MODEL.md`](requirements/ERROR-MODEL.md) — candidate history/API error categories; codes/status mapping remain draft.
- [`requirements/SECURITY-REQUIREMENTS.md`](requirements/SECURITY-REQUIREMENTS.md) — local threat boundaries and security acceptance.
- [`requirements/COMPATIBILITY-MATRIX.md`](requirements/COMPATIBILITY-MATRIX.md) — compatibility commitments and unresolved decisions.
- [`requirements/TEST-STRATEGY.md`](requirements/TEST-STRATEGY.md) — contract, migration, UI, integration, and platform evidence.
- [`CONTRACT-CROSSWALK.md`](CONTRACT-CROSSWALK.md) — source-contract and FR-501–FR-511 traceability; open gates remain explicit.

- [`epics/README.md`](epics/README.md) and [`issues/README.md`](issues/README.md) — EPIC-001 / #390 and its planning-only subissue contracts.

## Release boundary

v0.4: approved remediation plan → isolated apply → verification → rescan → patch evidence.
v0.5 proposal: local scan history → local API → embedded browser console.
v0.6: team collaboration, PostgreSQL, identity, RBAC, and exceptions (not included here).

The v0.5 proposal is local-first and read-only from the browser. It does not add scan triggering, remediation, repository mutation, remote publishing, or team services. Exact endpoints, persistence schema, retention policy, UI routes, component dependencies, and serving command behavior are unapproved until the design gates below pass.

## Mandatory design-first gates

1. Reconcile user workflows and domain/resource vocabulary; disposition each unresolved #356 predecessor item.
2. Review and approve [`UX-DESIGN.md`](UX-DESIGN.md), including the requested shadcn/ui style, accessibility states, and frontend dependency choice.
3. Review and approve [`API-DESIGN.md`](API-DESIGN.md), including OpenAPI 3.1, listener/origin/security, errors, bounds, and compatibility behavior.
4. Review and approve the SQLite schema, lifecycle, migration, and recovery decision. Record the approved decision as an ADR before storage implementation.
5. Complete security, failure/data, test, and compatibility contracts; name an independent architecture/security reviewer.
6. Derive ordered epics and implementation issues from these accepted contracts and update the v0.5 plan and Definition of Ready.
7. Owner and independent reviewer separately approve the complete Definition of Ready.

UX and API design may iterate in parallel only after their shared resource concepts are reconciled. No runtime implementation issue may be created or started before these gates pass. This package does not create implementation issues or mark any gate complete.

## Owner-directed planning-ticket exception

The owner directed creation of a v0.5 planning epic and design/readiness subissues without waiting for an independent reviewer. This exception applies only to those planning tickets: no independent review is claimed, the design contracts remain drafts, and runtime implementation remains blocked until the Definition of Ready and separate human-review requirements are satisfied. The owner also confirmed that the five open v0.4 follow-ups remain tracked in [#356](https://github.com/geoffrey-xiao/deprail/issues/356), not carried into v0.5; none is thereby completed or closed.

## Definition of Ready

- [ ] Owner approves the v0.5 outcome, inclusions, exclusions, compatibility boundary, failure model, and preview acceptance gate.
- [ ] Every open #356 follow-up is explicitly assigned, carried forward, or deferred with owner and target; no predecessor gate is inferred complete.
- [ ] UX design is approved and browser interaction/accessibility criteria are observable.
- [ ] OpenAPI/API security design is reviewed; endpoint and data schemas are versioned and validated.
- [ ] SQLite schema, transaction, migration, retention/deletion, corruption, backup/recovery, and permissions decision is accepted.
- [ ] Security and failure contracts include local-origin, path, process, privacy, and sensitive-data boundaries.
- [ ] Compatibility matrix, cross-platform test strategy, packaging constraints, and rollback/recovery evidence are accepted.
- [ ] Named independent architecture/security reviewer records approval separately from the owner.
- [ ] Ordered issue contracts each have one primary outcome, dependencies, inputs/outputs, failures, acceptance evidence, reviewer, and rollback notes.
- [ ] Project/milestone/view metadata and evidence links are verified.

### Readiness evidence and blockers

The rows below map this checklist and [release-plan §12](../../03-planning/deprail-development-plan-v0.5.0.md#12-definition-of-ready) to current evidence. All checklist boxes above remain unchecked: a merged proposal or a closed planning issue is not approval of the complete gate. `@geoffrey-xiao` owns disposition and evidence links; the independent reviewer must be named by the owner and record a separate decision.

| Gate | Current evidence, not acceptance | Owner's next action / blocking decision |
| --- | --- | --- |
| Scope and predecessor, plan §12 items 1–2 | Draft plan [PR #388](https://github.com/geoffrey-xiao/deprail/pull/388), package [PR #389](https://github.com/geoffrey-xiao/deprail/pull/389), open [#387](https://github.com/geoffrey-xiao/deprail/issues/387). The owner [assigned/targeted all five #356 follow-ups](https://github.com/geoffrey-xiao/deprail/issues/387#issuecomment-5792798436) in the [v0.4 retrospective](../../retrospectives/RETROSPECTIVE-v0.4.0-preview.2.md#follow-up-actions); [#356](https://github.com/geoffrey-xiao/deprail/issues/356) stays open outside v0.5. | Reconcile and accept the plan and each predecessor disposition in #387 without claiming the v0.4 gaps passed or silently importing them into v0.5. |
| Owner scope, compatibility, and preview gate | Proposed included/excluded work and release gate are in the [plan](../../03-planning/deprail-development-plan-v0.5.0.md#3-included-scope-proposed), [PRD](PRD-v0.5.md), and [crosswalk](CONTRACT-CROSSWALK.md#remaining-decisions-and-approval-gates). | Explicitly accept or revise the release boundary, history capture/CLI exit behavior, measured size limits or retention risk, and preview acceptance criteria; record any material change through the planning hierarchy. |
| UX and accessibility | [#394](https://github.com/geoffrey-xiao/deprail/issues/394) is closed; the [UX proposal](UX-DESIGN.md) and mockups exist, but a closed design issue does not supply independent review. | Approve component/dependency/tokens, browser/OS/assistive-technology matrix and response states after independent UX/security review; link decision and actual-surface acceptance criteria. |
| API/OpenAPI and local listener | [PR #404](https://github.com/geoffrey-xiao/deprail/pull/404) merged a candidate profile; [#396](https://github.com/geoffrey-xiao/deprail/issues/396) remains in Review. The owner selected the separately versioned history-projection direction, now proposed in [PR #410](https://github.com/geoffrey-xiao/deprail/pull/410); [API-DESIGN.md](API-DESIGN.md) has no accepted OpenAPI artifact. | Independently review exact projected fields/diagnostics and freeze/validate OpenAPI 3.1 resources/examples, pagination/cursors, version/error mapping, numeric bounds and listener/Host/Origin/CORS/CSRF/auth/shutdown rules. |
| SQLite and data lifecycle | [#391](https://github.com/geoffrey-xiao/deprail/issues/391) closed after owner-approved [PR #407](https://github.com/geoffrey-xiao/deprail/pull/407); [ADR-0004](../../adr/ADR-0004-local-scan-history.md) remains Proposed and its original unchanged-report payload is superseded by a candidate `history-v1` projection in PR #410. | Approve exact projection/SQL separately from the CLI report; resolve payload/store bounds, retention-growth risk, opt-in capture and persistence-failure precedence; record independent review of permissions, migrations, backup/recovery and artifacts. |
| Security, failure, test, compatibility and rollback | Draft [PR #408](https://github.com/geoffrey-xiao/deprail/pull/408) merged with three-OS CI; [#393](https://github.com/geoffrey-xiao/deprail/issues/393) remains in Review. FR/SEC/STORE scenarios and rollback expectations in the [test strategy](requirements/TEST-STRATEGY.md#9-requirement-and-threat-evidence-matrix) are future evidence, not executed checks. | Review local-origin/privacy/error controls; pin/select supported toolchains, SQLite driver, frontend/browser matrix, response limits and recovery owner/procedure. Accept the evidence plan without claiming runtime passes. |
| Independent review and DoR sign-off | No independent architecture/security reviewer is named. Planning-ticket owner oversight and merged PRs do not replace that review. | Owner names an independent reviewer; reviewer records architecture/security (including UX/API/SQLite) decision, then owner separately records the complete DoR decision and remaining risks. |
| Ordered issue contracts and tracking | [#392](https://github.com/geoffrey-xiao/deprail/issues/392) is open in Project Review with [PR #409](https://github.com/geoffrey-xiao/deprail/pull/409); its [conditional delivery map](issues/V05-005-readiness-map.md#conditional-post-dor-delivery-map) is not an issue backlog. [EPIC-001 / #390](https://github.com/geoffrey-xiao/deprail/issues/390) has all six planning children; the [v0.5.0 milestone](https://github.com/geoffrey-xiao/deprail/milestone/11) and [Project view](https://github.com/orgs/geoffrey-xiao/projects/1/views/6) exist, with view filter `milestone:v0.5.0`. | After contracts and separate DoR approvals, finalize one-outcome issue contracts with reviewer, risk, acceptance/evidence/rollback and verified labels/milestone/Project links. Do not create or activate implementation issues beforehand. |

**Additional DoR blocker:** The actual scan report is neither privacy-safe unchanged JSON nor v1alpha finding/error-schema-valid. [ADR-0004 records the validation finding and owner-selected separate projection](../../adr/ADR-0004-local-scan-history.md#owner-selected-direction-independent-history-projection-not-accepted). Exact `history-v1` allowlisted fields, unknown/safe-diagnostic handling, SQL validation, API response mapping and source-version compatibility still require owner and independent reviewer acceptance; no raw-report persistence, public API freeze or runtime work is authorized.

## No implementation authorization

All documents in this package are review drafts. No SQL schema, public API, UI dependency, endpoint, serving behavior, or migration is final. No issue may authorize runtime implementation until the v0.5 plan and this package pass the Definition of Ready and the owner and independent reviewer record their separate decisions.
