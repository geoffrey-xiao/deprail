# v0.4 Issue Contracts

These issue contracts are local source contracts for the GitHub sub-issues. They remain subject to PR #312, the ADR-0003 scope-change decision, owner/reviewer approval, and the crosswalk in [`CONTRACT-CROSSWALK.md`](CONTRACT-CROSSWALK.md).

| Local issue | GitHub issue | Parent epic | Dependencies | Outcome |
|---|---|---|---|---|
| V04-001 | [#302](https://github.com/geoffrey-xiao/deprail/issues/302) | [EPIC-001 / #313](https://github.com/geoffrey-xiao/deprail/issues/313) | None | Isolated workspace and rollback boundary |
| V04-002 | [#307](https://github.com/geoffrey-xiao/deprail/issues/307) | [EPIC-002 / #314](https://github.com/geoffrey-xiao/deprail/issues/314) | V04-001 | Approval bound to plan and source identity |
| V04-003 | [#308](https://github.com/geoffrey-xiao/deprail/issues/308) | [EPIC-002 / #314](https://github.com/geoffrey-xiao/deprail/issues/314) | V04-002 | Bounded supported mutation |
| V04-004 | [#309](https://github.com/geoffrey-xiao/deprail/issues/309) | [EPIC-003 / #315](https://github.com/geoffrey-xiao/deprail/issues/315) | V04-003 | Verification and before/after scan classification |
| V04-005 | [#310](https://github.com/geoffrey-xiao/deprail/issues/310) | [EPIC-004 / #316](https://github.com/geoffrey-xiao/deprail/issues/316) | V04-001 through V04-004 | Patch evidence and failure outcomes |
| V04-007 | [#329](https://github.com/geoffrey-xiao/deprail/issues/329) | Release parent [#299](https://github.com/geoffrey-xiao/deprail/issues/299); extends EPIC-003 | V04-004 and existing baseline contract | Trusted baseline generated from complete scan output |
| V04-006 | [#311](https://github.com/geoffrey-xiao/deprail/issues/311) | [EPIC-005 / #317](https://github.com/geoffrey-xiao/deprail/issues/317) | V04-001 through V04-005 | Cross-platform acceptance matrix |

The release parent is [#299](https://github.com/geoffrey-xiao/deprail/issues/299). The GitHub hierarchy is #299 → #313–#317 → #302, #307–#311, with V04-007 tracked directly under #299 as an additive scope change. No issue authorizes caller-worktree mutation, publication, autonomous approval, or unsupported ecosystems.
