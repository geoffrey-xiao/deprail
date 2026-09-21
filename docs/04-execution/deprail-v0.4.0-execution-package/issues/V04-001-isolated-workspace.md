# V04-001: Establish Isolated Workspace and Rollback Boundary

**Epic:** EPIC-001
**Status:** Proposed; implementation blocked

## Scope

Canonicalize the source root, create a temporary detached worktree from the reviewed commit, capture source/workspace identity and snapshots, restrict writes to authorized paths, and provide idempotent cleanup/rollback evidence.

## Excludes

Package-manager execution, dependency mutation, verification discovery, rescan, commit, push, PR, merge, and publication.

## Acceptance

Source tree remains unchanged after success, failure, timeout, and cancellation. Traversal and external symlinks fail closed. Cleanup is idempotent and records failure explicitly. Contract tests cover hostile paths and platform semantics.
