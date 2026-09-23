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
- [`UX-DESIGN.md`](UX-DESIGN.md) — workflows, information architecture, shadcn/ui-inspired direction, interaction and accessibility states.
- [`API-DESIGN.md`](API-DESIGN.md) — proposed local REST resource and operation contract; OpenAPI artifact remains a design-gate deliverable.
- [`requirements/FUNCTIONAL-REQUIREMENTS.md`](requirements/FUNCTIONAL-REQUIREMENTS.md) — observable behaviors and release gates.
- [`requirements/FAILURE-AND-DATA-CONTRACT.md`](requirements/FAILURE-AND-DATA-CONTRACT.md) — proposed persistence and failure semantics.
- [`requirements/SECURITY-REQUIREMENTS.md`](requirements/SECURITY-REQUIREMENTS.md) — local threat boundaries and security acceptance.
- [`requirements/COMPATIBILITY-MATRIX.md`](requirements/COMPATIBILITY-MATRIX.md) — compatibility commitments and unresolved decisions.
- [`requirements/TEST-STRATEGY.md`](requirements/TEST-STRATEGY.md) — contract, migration, UI, integration, and platform evidence.

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

## No implementation authorization

All documents in this package are review drafts. No SQL schema, public API, UI dependency, endpoint, serving behavior, or migration is final. No issue may authorize runtime implementation until the v0.5 plan and this package pass the Definition of Ready and the owner and independent reviewer record their separate decisions.
