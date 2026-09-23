# v0.5 Issue Contracts

These local planning contracts map to GitHub issues. All new issues belong beneath EPIC-001 as GitHub subissues; existing #387 is added as the V05-000 child. None is a runtime implementation issue.

| Local issue | GitHub issue | Parent epic | Dependencies | Outcome |
|---|---|---|---|---|
| V05-000 | [#387](https://github.com/geoffrey-xiao/deprail/issues/387) | [EPIC-001 / #390](https://github.com/geoffrey-xiao/deprail/issues/390) | #356 | Prepare the v0.5.0 development baseline |
| V05-001 | [#394](https://github.com/geoffrey-xiao/deprail/issues/394) | EPIC-001 / #390 | V05-000 | Approve-ready UX evidence for history and scan detail |
| V05-002 | [#396](https://github.com/geoffrey-xiao/deprail/issues/396) | EPIC-001 / #390 | V05-001 | Complete OpenAPI and local API security contract |
| V05-003 | [#391](https://github.com/geoffrey-xiao/deprail/issues/391) | EPIC-001 / #390 | V05-000; shared history vocabulary from V05-001 | Decide SQLite schema, lifecycle, migration, and recovery |
| V05-004 | [#393](https://github.com/geoffrey-xiao/deprail/issues/393) | EPIC-001 / #390 | V05-001 through V05-003 | Complete cross-cutting security, failure, compatibility, and test contracts |
| V05-005 | [#392](https://github.com/geoffrey-xiao/deprail/issues/392) | EPIC-001 / #390 | V05-000 through V05-004 | Assemble the execution package and readiness map; no implementation issue before DoR |

The existing #356 predecessor follow-ups remain separately tracked there by owner direction. The planning-ticket reviewer exception is limited to creating/tracking these design tasks; it does not satisfy the independent-review gate or authorize implementation.
