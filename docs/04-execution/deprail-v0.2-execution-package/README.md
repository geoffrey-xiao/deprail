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

v0.2 has two execution layers:

1. **Prerequisite hardening:** EPIC-001 through EPIC-004 close v0.1 contract, scanner, scope, and release-evidence gaps.
2. **CI guardrail capabilities:** EPIC-006 through EPIC-008 deliver baselines, diff, policy, exceptions, SARIF, GitHub Action validation, and public-preview evidence.

The prerequisite layer is necessary but is not the complete v0.2 product outcome. EPIC-005 remains a decision boundary for capabilities not already selected by the product and architecture baselines.

## Package status

This package is the local planning context for v0.2. GitHub milestone, Project, epic, issue, and parent-child records are created only from reconciled contracts during stage kickoff.

## Required v0.2 gate

Before implementation begins for a new capability, the project owner must review the PRD, transition record, release plan, target version decision, and Master Checklist. Before v0.2 release, prerequisite hardening, S5-S7 guardrail work, all checklist items, linked evidence, cross-platform verification, and owner acceptance during PR review are required.

