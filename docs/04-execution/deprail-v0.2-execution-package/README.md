# DepRail v0.2.0 Execution Package

**Status:** Planning baseline; not implementation-ready until Definition of Ready is reviewed
**Version:** v0.2.0
**Predecessor:** [v0.1 execution package](../deprail-v0.1-execution-package/README.md)
**Product:** DepRail dependency security guardrail

## Purpose

This package is the durable planning and execution context for the v0.2.0 product line. It preserves v0.1 contracts unless a v0.2 document explicitly classifies a change as clarified, extended, breaking, deferred, or accepted limitation.

## Source-of-truth order

1. [PRD-v0.2](PRD-v0.2.md)
2. [v0.1-to-v0.2 transition](V0.1-TO-V0.2-TRANSITION.md)
3. [requirements](requirements/)
4. [architecture baseline](ARCHITECTURE-v0.2.md)
5. [epics](epics/)
6. [issues](issues/)
7. [sprint plans](tracking/)
8. [Master Checklist](tracking/MASTER-CHECKLIST.md)
9. [release procedure](../../RELEASE-v0.2.md)

Repository-wide workflow rules remain in [`AGENTS.md`](../../../AGENTS.md). Shared Project conventions remain in the [v0.1 Project setup guide](../deprail-v0.1-execution-package/tracking/GITHUB-PROJECT-SETUP.md) until a later approved update.

## v0.2 planning position

v0.2 begins as a hardening and product-readiness line. It must first close the highest-value v0.1 gaps:

- stable empty JSON collections;
- strict command argument validation;
- real OSV-Scanner v2 fixtures and exit-code coverage;
- tag/commit version identity;
- scan-from-outside-root regression coverage;
- release-mode evidence and workflow discipline.

New product capabilities are candidates, not commitments, until they have an approved issue contract and acceptance evidence.

## Package status

This package is local planning context only. It does not create GitHub Issues, Project views, milestones, labels, or remote release records. Those may be established later during stage kickoff after this package is reviewed.

## Required v0.2 gate

Before implementation begins, the project owner must review the PRD, transition record, target version decision, and Master Checklist. Before v0.2 release, all checklist items require linked evidence and owner acceptance during PR review.
