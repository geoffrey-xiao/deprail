# V05-000: Prepare the v0.5.0 Development Baseline

- GitHub Issue: [#387](https://github.com/geoffrey-xiao/deprail/issues/387)
- GitHub parent epic: [EPIC-007 / #390](https://github.com/geoffrey-xiao/deprail/issues/390)
- Epic: [EPIC-007](../epics/EPIC-007-local-history-readiness.md)
- Target: `v0.5.0`
- Status: Open; planning only
- Priority: P0
- Risk: R2
- Area: docs
- Owner: `@geoffrey-xiao`
- Reviewer: Independent architecture/security reviewer remains to be named; the owner's planning-ticket exception does not change #387's separate review gate.
- Dependencies: #356; product, architecture, roadmap, release-plan, and execution-package reconciliation.

## Value

Prepare a bounded, evidence-backed v0.5 plan for local scan history and the embedded console without starting runtime implementation before the required design and review gates.

## Scope

Reconcile the v0.5 plan/package with the product design, architecture, roadmap, and #356 predecessor evidence; maintain the explicit v0.5 boundary; record the owner's disposition that all five #356 follow-ups remain in #356; maintain the Definition of Ready; and link the dependency-ordered planning work under EPIC-007.

## Exclusions

No runtime code, schema, migration, UI/API implementation, implementation child issue, #356 closure, or claim that an open predecessor gap is complete.

## Inputs, outputs, and failure behavior

Inputs are the versioned planning hierarchy, v0.4.0-preview.2 evidence/retrospective, #356, and the v0.5 design package. Outputs are reconciled scope, explicit predecessor dispositions, linked planning contracts, and a current readiness state. Conflicting contracts or missing approvals keep implementation blocked and are recorded rather than silently resolved.

## Required verification

Check source links, owner disposition crosswalk, issue hierarchy, milestone, labels, Project membership/status, and the v0.5 view filter. No runtime tests are in scope.

## Acceptance criteria

- The included and excluded v0.5 scope remains traceable to the product, architecture, and roadmap.
- All five #356 follow-ups are explicitly kept in #356; this does not close them or imply their completion.
- Planning issues are linked under EPIC-007; implementation issues remain absent until the v0.5 DoR passes.
- The independent reviewer and owner approval remain separate, visible gates for implementation readiness.
