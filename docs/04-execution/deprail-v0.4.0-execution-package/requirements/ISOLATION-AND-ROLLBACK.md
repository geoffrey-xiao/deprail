# v0.4 Isolation and Rollback Requirements

**Status:** Proposed; requires ADR-0002 and architecture/security review.
**Related:** #299, #302, #303

## Requirements

### ISO-001 Canonical source root

The executor MUST resolve the requested repository root canonically and reject traversal, non-directory roots, and symlink escape before creating an isolated workspace.

### ISO-002 Worktree isolation

The default workspace MUST be a unique temporary Git worktree created from the reviewed source commit. The caller's worktree MUST remain outside the mutation path.

### ISO-003 Permissions

The workspace and evidence directory MUST use restrictive permissions. The implementation MUST not broaden permissions inherited from an unsafe parent without an explicit reviewed decision.

### ISO-004 Authorization binding

Execution MUST bind the plan digest, source commit, canonical root, authorized paths, and approval expiry. Any mismatch MUST fail before mutation.

### ISO-005 Initial snapshot

Before mutation, capture deterministic relative paths, file metadata required for restoration, and content digests. Evidence MUST identify the source commit and workspace.

### ISO-006 Authorized writes

Writes MUST remain below the isolated workspace and within plan-authorized paths. Symlinks and path re-resolution MUST be checked at the write boundary.

### ISO-007 Atomic writes

File replacement MUST write to a temporary file in the same filesystem and atomically rename it. Interrupted writes MUST leave the target either at its old or complete new content.

### ISO-008 Cleanup

Cleanup MUST be idempotent and report success, partial cleanup, or failure. Cleanup failure MUST remain visible and MUST NOT be reported as rollback success.

### ISO-009 Failure and cancellation

Timeout, cancellation, process failure, verification failure, rescan regression, and malformed tool output MUST trigger isolated-workspace cleanup and an explicit terminal state.

### ISO-010 Source immutability

The caller's repository tree, index, refs, manifests, and lockfiles MUST remain unchanged across success, failure, timeout, cancellation, and interrupted cleanup scenarios.

### ISO-011 Evidence

Evidence MUST include workspace identity, source commit, plan digest, changed paths, before/after digests, cleanup result, error code, and redacted diagnostics. It MUST exclude credentials and unrestricted environment values.

## Acceptance matrix

| Case | Expected result |
| --- | --- |
| Valid root and plan | Workspace created and identified |
| Traversal or external symlink | Rejected before workspace mutation |
| Approval mismatch or expiry | Rejected before mutation |
| Atomic write interruption | Old or complete new file, never a partial target |
| Process timeout | Process tree terminated; workspace cleanup recorded |
| User cancellation | Cancellation terminal state; source unchanged |
| Verification failure | No success result; cleanup/rollback evidence retained |
| Repeated cleanup | Deterministic, safe, and idempotent |
| Cleanup failure | Partial/failed state with explicit reason |
| All supported platforms | Equivalent serialized meaning and source immutability |
