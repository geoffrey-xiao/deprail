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
| Independent reviewer | [`@geoffreyxiaoai`](https://github.com/geoffreyxiaoai); approved the then-current package in [PR #412](https://github.com/geoffrey-xiao/deprail/pull/412); review of the current OpenAPI/security candidate remains pending. |

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

The owner recorded separate DoR approval for the then-current package in [#392](https://github.com/geoffrey-xiao/deprail/issues/392#issuecomment-5807906998), and the named independent reviewer approved that full-package review in [PR #412](https://github.com/geoffrey-xiao/deprail/pull/412). PR #412 predates the current [`OpenAPI 3.1 candidate`](../../../schemas/openapi/v1/openapi.yaml) and its concrete local security profile; those later API details remain unreviewed. Keep each gate unchecked until its current-version evidence and required acceptance are linked.

| Gate | Current evidence, not acceptance | Owner's next action / blocking decision |
| --- | --- | --- |
| Scope and predecessor, plan §12 items 1–2 | Draft plan [PR #388](https://github.com/geoffrey-xiao/deprail/pull/388), package [PR #389](https://github.com/geoffrey-xiao/deprail/pull/389), open [#387](https://github.com/geoffrey-xiao/deprail/issues/387). The owner [assigned/targeted all five #356 follow-ups](https://github.com/geoffrey-xiao/deprail/issues/387#issuecomment-5792798436) in the [v0.4 retrospective](../../retrospectives/RETROSPECTIVE-v0.4.0-preview.2.md#follow-up-actions); [#356](https://github.com/geoffrey-xiao/deprail/issues/356) stays open outside v0.5. | Reconcile and accept the plan and each predecessor disposition in #387 without claiming the v0.4 gaps passed or silently importing them into v0.5. |
| Owner scope, compatibility, and preview gate | Proposed included/excluded work and release gate are in the [plan](../../03-planning/deprail-development-plan-v0.5.0.md#3-included-scope-proposed), [PRD](PRD-v0.5.md), and [crosswalk](CONTRACT-CROSSWALK.md#remaining-decisions-and-approval-gates). | Explicitly accept or revise the release boundary, history capture/CLI exit behavior, measured size limits or retention risk, and preview acceptance criteria; record any material change through the planning hierarchy. |
| UX and accessibility | [#394](https://github.com/geoffrey-xiao/deprail/issues/394) is closed; the [UX proposal](UX-DESIGN.md) and mockups were part of the full-package review approved in [PR #412](https://github.com/geoffrey-xiao/deprail/pull/412). | Component/dependency/tokens, supported browser/OS/assistive-technology matrix and actual-surface keyboard, screen-reader, contrast and responsive evidence remain to resolve; package review does not claim browser verification. |
| API/OpenAPI and local listener | [PR #413](https://github.com/geoffrey-xiao/deprail/pull/413) merged the candidate; the owner merged five schema/security corrections in [PR #414](https://github.com/geoffrey-xiao/deprail/pull/414). Its body accidentally triggered auto-close of [#396](https://github.com/geoffrey-xiao/deprail/issues/396); the issue is reopened and its Project status is Review. | PR #414 has no independent human APPROVED review recorded; its only review is Codex COMMENTED. Owner merge accepts the correction patch, not the independent-review gate for the current contract. Keep implementation issue creation blocked until the named reviewer decides on this exact candidate and the current DoR is reconciled. |
| SQLite and data lifecycle | [#391](https://github.com/geoffrey-xiao/deprail/issues/391) closed after owner-approved [PR #407](https://github.com/geoffrey-xiao/deprail/pull/407); the owner-selected `history-v1` direction is in [merged PR #411](https://github.com/geoffrey-xiao/deprail/pull/411), and the then-current package was reviewed in [PR #412](https://github.com/geoffrey-xiao/deprail/pull/412). The ADR remains Proposed; no SQL/runtime schema is approved. | Exact projection fields, row validation, migration/lifecycle, retention-growth risk, opt-in capture and persistence-failure precedence remain explicit design decisions; review permissions, backups, recovery and artifact contracts before implementation. |
| Security, failure, test, compatibility and rollback | Draft [PR #408](https://github.com/geoffrey-xiao/deprail/pull/408) merged with three-OS CI; [#393](https://github.com/geoffrey-xiao/deprail/issues/393) remains in Review. The then-current threat/test contracts were included in PR #412's package review; FR/SEC/STORE scenarios remain planned evidence, not executed runtime checks. | Review the newer API local-origin/auth/error profile and failure/rollback changes; pin/select toolchains, SQLite driver, frontend/browser matrix, response limits and recovery owner/procedure. |
| Independent review and DoR sign-off | The owner comment on [#392](https://github.com/geoffrey-xiao/deprail/issues/392#issuecomment-5807906998) authorized task creation, and [`@geoffreyxiaoai`](https://github.com/geoffreyxiaoai) approved the then-current package in [PR #412](https://github.com/geoffrey-xiao/deprail/pull/412). Both predate [PR #414](https://github.com/geoffrey-xiao/deprail/pull/414); that correction was merged by the owner but has no independent APPROVED review. | The earlier authorization/review does not establish approval of the later OpenAPI contract. Record a reviewer decision on the exact current candidate and reconcile the current-version DoR before implementation issue creation. |
| Ordered issue contracts and tracking | [#392](https://github.com/geoffrey-xiao/deprail/issues/392) is open in Project Review with [PR #409](https://github.com/geoffrey-xiao/deprail/pull/409); its [conditional delivery map](issues/V05-005-readiness-map.md#conditional-post-dor-delivery-map) is not an issue backlog. [EPIC-001 / #390](https://github.com/geoffrey-xiao/deprail/issues/390) has all six planning children; the [v0.5.0 milestone](https://github.com/geoffrey-xiao/deprail/milestone/11) and [Project view](https://github.com/orgs/geoffrey-xiao/projects/1/views/6) exist, with view filter `milestone:v0.5.0`. | After contracts and separate DoR approvals, finalize one-outcome issue contracts with reviewer, risk, acceptance/evidence/rollback and verified labels/milestone/Project links. Do not create or activate implementation issues beforehand. |

**Current API design blocker:** The actual scan report is neither privacy-safe unchanged JSON nor v1alpha finding/error-schema-valid. [ADR-0004 records the validation finding and owner-selected separate projection](../../adr/ADR-0004-local-scan-history.md#owner-selected-direction-independent-history-projection-not-accepted). PR #414 merged the five schema corrections; its merge closed #396 and the issue was reopened because no independent review is recorded. Exact `history-v1` fields, safe-diagnostic handling, SQL validation, row/API mapping, source-version compatibility, and independent review remain DoR gates; no runtime work is authorized.

## No implementation authorization

The owner authorized implementation task creation in [#392](https://github.com/geoffrey-xiao/deprail/issues/392#issuecomment-5807906998), and the reviewer approved the then-current package in PR #412; both decisions predate the corrected OpenAPI candidate. PR #414 was later merged by the owner without an independent APPROVED review. Keep the v0.5 plan in Draft and do not create implementation issues or begin runtime work until the current contract and complete DoR have separate linked decisions.
