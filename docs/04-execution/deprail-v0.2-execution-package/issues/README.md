# v0.2 Issue Contract Index

These are local planning contracts only. They are not GitHub Issues yet. Each item must be converted into the repository issue template and pass Definition of Ready before implementation.

| ID | Title | Priority | Risk | Area | Epic | Sprint | Status |
| --- | --- | --- | --- | --- | --- | --- | --- |
| [V02-000](V02-000-context-baseline.md) | Establish v0.2.0 context baseline | P0 | R1 | foundation | Package gate | Sprint 0 | GitHub #98; review required |
| V02-002 | Reject unexpected command arguments | P0 | R1 | cli | EPIC-001 | Sprint 0 | Ready for review |
| V02-003 | Inject truthful build version identity | P1 | R1 | cli | EPIC-001 | Sprint 1 | Draft |
| V02-004 | Establish release-mode evidence checklist | P0 | R0 | docs | EPIC-001 | Sprint 0 | Ready for review |
| V02-005 | Add real OSV-Scanner v2 fixtures | P0 | R1 | adapter | EPIC-002 | Sprint 0 | Ready for review |
| V02-006 | Add scanner exit-code matrix | P0 | R1 | adapter | EPIC-002 | Sprint 0 | Ready for review |
| V02-007 | Verify scanner requested-root execution | P0 | R2 | adapter | EPIC-002 | Sprint 0 | Ready for review |
| V02-008 | Add outside-root scan regression | P1 | R2 | test | EPIC-003 | Sprint 1 | Draft |
| V02-009 | Update report examples and arrays | P0 | R1 | normalization | EPIC-003 | Sprint 0 | Ready for review |
| V02-010 | Verify deterministic artifacts and ordering | P1 | R2 | test | EPIC-003 | Sprint 1 | Draft |
| V02-011 | Verify release identity across artifacts | P1 | R1 | cli | EPIC-004 | Sprint 1 | Draft |
| V02-012 | Add release-mode smoke procedure | P1 | R1 | test | EPIC-004 | Sprint 1 | Draft |
| V02-013 | Record release artifact and supply-chain evidence | P0 | R2 | docs | EPIC-004 | Sprint 0 | Ready for review |
| V02-014 | Decide candidate capability boundary | P1 | R1 | foundation | EPIC-005 | Sprint 0 | Decision required |

## Common Definition of Ready

Every item needs a confirmed owner, reviewer, target version, dependencies, observable acceptance criteria, failure behavior, required tests/evidence, and explicit exclusions before branch creation.

## Common Definition of Done

The reviewed PR is merged, required CI passes, acceptance evidence is linked, remaining risk is recorded, and the Project item can move directly from `Review` to `Done` under the reviewed-merge workflow.
