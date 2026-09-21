# v0.4 Issue Contracts

These issue contracts are local source contracts for proposed GitHub issues. They remain subject to PR #312, owner/reviewer approval, and the crosswalk in [`CONTRACT-CROSSWALK.md`](CONTRACT-CROSSWALK.md).

| Local issue | GitHub issue | Epic | Dependencies | Outcome |
|---|---|---|---|---|
| V04-001 | #302 | EPIC-001 | None | Isolated workspace and rollback boundary |
| V04-002 | #307 candidate | EPIC-002 | V04-001 | Approval bound to plan and source identity |
| V04-003 | #308 candidate | EPIC-002 | V04-002 | Bounded supported mutation |
| V04-004 | #309 candidate | EPIC-003 | V04-003 | Verification and before/after scan classification |
| V04-005 | #310 candidate | EPIC-004 | V04-001 through V04-004 | Patch evidence and failure outcomes |
| V04-006 | #311 candidate | EPIC-005 | V04-001 through V04-005 | Cross-platform acceptance matrix |

The `candidate` designation means the existing issue must be reconciled with this local contract before it becomes implementation-ready. No issue authorizes caller-worktree mutation, publication, autonomous approval, or unsupported ecosystems.
