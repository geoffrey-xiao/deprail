# EPIC-001: v0.5 Local Web and History Readiness
**GitHub epic:** [#390](https://github.com/geoffrey-xiao/deprail/issues/390)

- Target: `v0.5.0`
- Status: Planning; design artifacts remain drafts
- Owner: `@geoffrey-xiao`
- Reviewer: `@geoffrey-xiao` for planning-ticket oversight only, per explicit owner direction; this is not independent architecture/security review
- Milestone: `v0.5.0`

## Outcome

Produce a reviewable v0.5 design package for local scan history and an embedded console, with explicit UX, API/security, persistence, failure, compatibility, and test contracts. Derive the implementation issue map only after the package Definition of Ready is satisfied.

## Boundaries

Includes completion of UX evidence, OpenAPI/API and local-listener decisions, SQLite schema/lifecycle/recovery ADR, cross-cutting security/failure/compatibility/test evidence, and the final execution package/implementation issue map. The five v0.4.0-preview.2 follow-ups remain tracked in #356 and are not duplicated in v0.5.

Excludes runtime Go, SQL, HTTP, React, or build-tool implementation; implementation issue creation before the v0.5 Definition of Ready; browser-triggered scans or mutation; repository/package changes; hosted/team services; identity/RBAC; remote publishing; and any v0.6 capability.

## Dependencies

- Preparation issue #387, which remains open and is represented locally as V05-000.
- Product design, architecture, roadmap, v0.5 release plan, and the draft execution package.
- #356 remains open; the owner's disposition is to keep its five follow-ups in that issue, without implying completion.

## Child issues and order

1. [V05-000 / #387 — Prepare the v0.5.0 development baseline](https://github.com/geoffrey-xiao/deprail/issues/387).
2. [V05-001 / #394 — Finalize history and scan-detail UX acceptance](https://github.com/geoffrey-xiao/deprail/issues/394).
3. [V05-002 / #396 — Freeze the local API and OpenAPI security contract](https://github.com/geoffrey-xiao/deprail/issues/396); depends on V05-001.
4. [V05-003 / #391 — Decide SQLite history schema and lifecycle](https://github.com/geoffrey-xiao/deprail/issues/391); may proceed in parallel after shared history vocabulary is recorded.
5. [V05-004 / #393 — Complete failure, security, compatibility, and test contracts](https://github.com/geoffrey-xiao/deprail/issues/393); depends on V05-001 through V05-003.
6. [V05-005 / #392 — Assemble the execution package and implementation readiness map](https://github.com/geoffrey-xiao/deprail/issues/392); depends on V05-000 through V05-004.

Every child is a planning/design deliverable, not an implementation authorization. The owner directed creation of these planning tickets without an independent reviewer; no independent review or design approval is claimed. The independent reviewer and separate human approval remain required by the v0.5 Definition of Ready before runtime implementation.

## Acceptance

- Each planning child has its outcome, scope, dependencies, failure behavior, evidence, owner, and observable acceptance recorded locally and on GitHub.
- GitHub child relationships, v0.5.0 milestone, required labels, and Project membership/status are verified.
- UX, API, SQLite, security, compatibility, failure, and test artifacts are complete enough for the required human review; draft status is preserved until actual approvals are recorded.
- The v0.5 Definition of Ready is not marked complete by this epic. No runtime implementation issue is created or started by this epic.
