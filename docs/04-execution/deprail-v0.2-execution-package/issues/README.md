# v0.2 Issue Contract Index

These are local planning contracts. Each item maps to the reconciled v0.2 release plan and must pass Definition of Ready before implementation.

| ID | Title | Priority | Risk | Area | Epic | Sprint | Status |
| --- | --- | --- | --- | --- | --- | --- | --- |
| [V02-000](V02-000-context-baseline.md) | Establish v0.2.0 context baseline | P0 | R1 | foundation | Package gate | Sprint 0 | GitHub #98; review required |
| [V02-001](V02-001-empty-report-collections.md) | Normalize empty report collections | P0 | R1 | normalization | EPIC-001 | Sprint 0 | GitHub [#112](https://github.com/geoffrey-xiao/deprail/issues/112), child of #106 |
| [V02-002](V02-002-strict-command-arguments.md) | Reject unexpected command arguments | P0 | R1 | cli | EPIC-001 | Sprint 0 | GitHub [#110](https://github.com/geoffrey-xiao/deprail/issues/110), child of #106 |
| [V02-003](V02-003-build-version-identity.md) | Inject truthful build version identity | P1 | R1 | cli | EPIC-001 | Sprint 1 | GitHub [#107](https://github.com/geoffrey-xiao/deprail/issues/107), child of #106 |
| [V02-004](V02-004-release-mode-checklist.md) | Establish release-mode evidence checklist | P0 | R0 | docs | EPIC-001 | Sprint 0 | GitHub [#108](https://github.com/geoffrey-xiao/deprail/issues/108), child of #106 |
| [V02-005](V02-005-osv-v2-fixtures.md) | Add real OSV-Scanner v2 fixtures | P0 | R1 | adapter | EPIC-002 | Sprint 0 | Ready for review |
| [V02-006](V02-006-scanner-exit-matrix.md) | Add scanner exit-code matrix | P0 | R1 | adapter | EPIC-002 | Sprint 0 | Ready for review |
| [V02-007](V02-007-requested-root-execution.md) | Verify scanner requested-root execution | P0 | R2 | adapter | EPIC-002 | Sprint 0 | Ready for review |
| [V02-008](V02-008-outside-root-regression.md) | Add outside-root scan regression | P1 | R2 | test | EPIC-003 | Sprint 1 | Draft |
| [V02-009](V02-009-report-arrays-and-examples.md) | Update report examples and arrays | P0 | R1 | normalization | EPIC-003 | Sprint 0 | Ready for review |
| [V02-010](V02-010-deterministic-artifacts-ordering.md) | Verify deterministic artifacts and ordering | P1 | R2 | test | EPIC-003 | Sprint 1 | Draft |
| [V02-011](V02-011-release-identity-verification.md) | Verify release identity across artifacts | P1 | R1 | cli | EPIC-004 | Sprint 1 | Draft |
| [V02-012](V02-012-release-mode-smoke.md) | Add release-mode smoke procedure | P1 | R1 | test | EPIC-004 | Sprint 1 | Draft |
| [V02-013](V02-013-release-evidence.md) | Record release artifact and supply-chain evidence | P0 | R2 | docs | EPIC-004 | Sprint 0 | Ready for review |
| [V02-014](V02-014-capability-boundary-decision.md) | Decide candidate capability boundary | P1 | R1 | foundation | EPIC-005 | Sprint 0 | Decision required |

| [V02-015](V02-015-baseline-representation-storage.md) | Baseline representation and storage | P0 | R1 | normalization | EPIC-006 | S5 | Planned |
| [V02-016](V02-016-base-head-comparison.md) | Base/head scan comparison | P0 | R2 | normalization | EPIC-006 | S5 | Planned |
| [V02-017](V02-017-change-classification.md) | New/resolved/unchanged classification | P0 | R1 | normalization | EPIC-006 | S5 | Planned |
| [V02-018](V02-018-diff-command-output.md) | Deterministic diff output and `deprail diff` | P0 | R1 | cli | EPIC-006 | S5 | Planned |
| [V02-019](V02-019-policy-evaluator.md) | Typed policy evaluation | P0 | R2 | foundation | EPIC-007 | S6 | Planned |
| [V02-020](V02-020-expiring-exceptions.md) | Expiring exceptions | P1 | R2 | foundation | EPIC-007 | S6 | Planned |
| [V02-021](V02-021-policy-exit-behavior.md) | Policy exit behavior | P0 | R2 | cli | EPIC-007 | S6 | Planned |
| [V02-022](V02-022-sarif-output.md) | SARIF output | P0 | R1 | cli | EPIC-007 | S6 | Planned |
| [V02-023](V02-023-github-action-packaging.md) | GitHub Action packaging | P0 | R2 | docs | EPIC-008 | S7 | Planned |
| [V02-024](V02-024-action-permissions-caching.md) | Action permissions and caching | P0 | R2 | docs | EPIC-008 | S7 | Planned |
| [V02-025](V02-025-pull-request-validation.md) | Pull-request validation | P0 | R2 | test | EPIC-008 | S7 | Planned |
| [V02-026](V02-026-public-preview-evidence.md) | Public-preview release evidence | P0 | R2 | docs | EPIC-008 | S7 | Planned |

## Common Definition of Ready

Every item needs a confirmed owner, reviewer, target version, dependencies, observable acceptance criteria, failure behavior, required tests/evidence, and explicit exclusions before branch creation.

## Common Definition of Done

The reviewed PR is merged, required CI passes, acceptance evidence is linked, remaining risk is recorded, and the Project item can move directly from `Review` to `Done` under the reviewed-merge workflow.
