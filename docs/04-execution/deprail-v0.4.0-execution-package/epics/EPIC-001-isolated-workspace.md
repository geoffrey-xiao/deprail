# EPIC-001: Isolated Workspace and Rollback

**Status:** Proposed

## Outcome

Approved remediation work executes in a temporary workspace, preserves the caller repository, and records cleanup or rollback explicitly.

## Boundaries

Includes canonical path validation, worktree lifecycle, authorized paths, atomic writes, snapshot/digest capture, cancellation cleanup, and rollback evidence. Excludes package-manager behavior and publication.

## Acceptance

Hostile paths fail closed; source digests remain unchanged; cleanup is idempotent; timeout/cancellation leaves evidence; Linux/macOS/Windows semantics are equivalent.
