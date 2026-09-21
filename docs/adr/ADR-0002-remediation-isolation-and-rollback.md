# ADR-0002 Remediation Isolation and Rollback

- Status: Accepted
- Date: 2026-09-21
- Owners: `@geoffrey-xiao`
- Related issues: [#299](https://github.com/geoffrey-xiao/deprail/issues/299), [#302](https://github.com/geoffrey-xiao/deprail/issues/302), [#303](https://github.com/geoffrey-xiao/deprail/issues/303)
- Release: `v0.4.0-preview.1`

## Context

v0.4 introduces the first repository mutation boundary. Product design requires approved remediation plans to apply in an isolated branch or worktree, followed by bounded verification and rescan. Architecture requires canonical path containment, safe process invocation, atomic writes, cancellation, rollback, and explicit human approval. The caller's repository must remain unchanged when an operation fails or is cancelled.

## Decision

Use a temporary Git worktree created from the reviewed source commit as the default isolation boundary.

The executor MUST:

1. Resolve and canonicalize the source repository root before creating the worktree.
2. Reject traversal, external symlink escape, non-directory roots, and ambiguous repository identity.
3. Create a unique temporary worktree outside the caller's working tree with restrictive permissions.
4. Bind the execution to the source commit, plan digest, requested operations, and approval record.
5. Capture initial repository file paths and content digests before mutation.
6. Permit writes only below the isolated worktree and only for plan-authorized paths.
7. Use atomic temporary-file writes and rename within the worktree.
8. Discard the worktree on success after evidence retention, or on failure, timeout, and cancellation after recording cleanup status.
9. Never mutate, reset, clean, delete, or restore the caller's original worktree.
10. Report cleanup and rollback as explicit evidence; never infer success from process exit alone.

A future platform without Git worktree support may define an equivalent adapter only through a new ADR. The first implementation does not support arbitrary workspace providers.

Approval is single-use and expires. It is invalid if the plan digest, source commit, repository root, or authorized path set changes.

## Alternatives Considered

### Mutate the caller's worktree

Rejected. It violates the product safety boundary and makes rollback and concurrent use unsafe.

### Temporary directory copy

Deferred. Copy semantics, ignored files, symlinks, permissions, and repository identity are harder to prove equivalent than a Git worktree.

### Container or VM isolation

Deferred. Stronger isolation is not required for the first local core slice and would add platform and distribution complexity.

## Consequences

### Positive

- Original source remains protected by construction.
- Worktree identity and source commit are inspectable.
- Failure cleanup can be deterministic and evidence-backed.
- The boundary supports later verification and rescan without widening mutation scope.

### Negative

- Requires Git availability and worktree support.
- Temporary disk usage and cleanup failures must be handled.
- Windows permissions and open-handle behavior require platform-specific tests.

## Compatibility and Migration

No existing v0.3 or v0.3.1 command mutates repositories. v0.4 adds the boundary behind the new apply workflow; scan, plan, diff, policy, and machine-output contracts remain unchanged.

## Validation

Acceptance requires contract tests for canonical paths, symlink escape, restrictive permissions, source-tree immutability, repeated cleanup, interrupted cleanup, timeout, cancellation, and Linux/macOS/Windows behavior. These remain implementation acceptance criteria for #302.
