# DepRail v0.4.0 Development Plan

**Release target:** `v0.4.0-preview.1`
**Mode:** Preview first; stable follow-up only after complete evidence
**Status:** Planning baseline; implementation is blocked pending Definition of Ready approval
**Parent epic:** [#299](https://github.com/geoffrey-xiao/deprail/issues/299)
**Milestone:** `v0.4.0`
**Owner:** `@geoffrey-xiao`
**Reviewer:** Architecture/security reviewer to be named

## 1. Source reconciliation

This plan follows the product design, architecture, and roadmap hierarchy:

- Product: `docs/01-product/deprail-product-design-v1-ai.md`
- Architecture: `docs/02-architecture/deprail-architecture-and-tech-stack-v1.md`
- Roadmap: `docs/03-planning/deprail-roadmap-v1.md`
- Predecessor evidence: `docs/release-evidence/RELEASE-v0.3.1-preview.1-EVIDENCE.md`
- Parent epic: issue #299

The roadmap defines v0.4 as remediation and verification. v0.3 remains plan-only and v0.3.1 remains presentation-only. This plan introduces no web, team, MCP, broad scanner, or autonomous publication scope.

## 2. User outcome

A user can apply an explicitly approved remediation plan inside an isolated workspace, run bounded verification, rescan the result, and receive durable patch evidence without risking the source repository or hiding partial failure.

## 3. Included scope

- Read and validate a versioned remediation plan.
- Require explicit approval before any mutation.
- Create an isolated temporary worktree or equivalent workspace.
- Apply only allowlisted, plan-described dependency changes.
- Execute package-manager operations without a shell and without implicit install scripts unless explicitly approved.
- Discover and run bounded tests, builds, and type checks.
- Capture before/after trees, diffs, commands, exit codes, and artifacts.
- Rescan after verification and classify resolved, residual, introduced, and unverifiable findings.
- Roll back or discard the isolated workspace on failure.
- Emit deterministic machine-readable patch and verification evidence.

## 4. Explicit exclusions

- Direct mutation of the caller's working tree.
- Automatic commit, branch push, pull request, merge, or publication.
- Autonomous approval or policy bypass.
- Unbounded package installation or arbitrary script execution.
- New scanner families, web services, team storage, MCP write tools, or agent mutation.
- Stable-release claim without a separate release gate.

## 5. Safety contracts

### Mutation boundary

Mutation is permitted only inside an isolated workspace created after plan validation and approval. The original repository root is read-only. Every changed path must be declared by the plan or generated package-manager result and remain inside the isolated workspace.

### Isolation

The preferred implementation is a temporary Git worktree from the reviewed source commit. The worktree must have a unique path, restrictive permissions, canonical-root validation, cleanup ownership, and a recoverable reference to the source commit. Failure cleanup must not remove or modify the caller's repository.

### Approval and authorization

The application must receive an explicit approval token/state tied to the plan digest, source commit, repository root, and requested operations. Approval is single-use, expires, and is invalidated when the plan, source, or target files change. Dry-run and plan inspection require no mutation approval; apply requires it.

### Atomicity and rollback

Writes use temporary files and atomic rename within the isolated workspace. Before mutation, capture the original tree/digests. On any failed operation, timeout, cancellation, verification failure, or rescan regression, discard the isolated workspace or restore from the recorded snapshot. Never claim rollback succeeded without evidence.

### Package-manager execution

Commands use direct argument arrays, explicit working directories, approved environment variables, bounded output, deadlines, and process-tree cancellation. Package-manager scripts and network access are denied by default. Any permitted network or script behavior must be represented in the plan and approval record.

### Process limits

Every child process has a deadline, cancellation path, stdout/stderr byte limits, and process-tree termination. Timeouts, cancellation, non-zero exit, malformed output, and output-limit failures are typed errors, never empty success.

### Verification discovery

Tests, builds, and type checks are discovered from repository conventions and supported ecosystem contracts. Commands are ranked and presented before execution. Only supported, non-interactive commands inside the isolated workspace may run. Unknown commands require explicit approval and remain unexecuted by default.

### Rescan and comparison

The post-change scan uses the same repository root semantics and scanner contract as the pre-change scan. Evidence compares before/after findings by stable keys and classifies resolved, residual, introduced, and unknown states. A scan failure yields incomplete verification, not success.

### Patch evidence

Evidence includes plan digest, source commit, isolated workspace identity, changed files, before/after digests, diff, commands and arguments, tool versions, exit codes, durations, scan identities, findings transition, rollback state, and redaction-safe diagnostics. Credentials, full environments, and source content outside the diff are excluded.

### Failure and partial failure

The result has explicit `complete`, `partial`, `failed`, and `cancelled` states. Partial application never becomes success. Successful intermediate evidence is retained, failed operation scope is named, rollback status is explicit, and the original repository remains unchanged.

## 6. Acceptance gate

The preview is acceptable only when:

- Approved plans apply only inside an isolated workspace.
- Original repository trees remain unchanged.
- Process, path, symlink, credential, timeout, cancellation, and output-limit boundaries are tested.
- Verification and post-change rescan results are deterministic and complete/partial/failed states are explicit.
- Rollback and interrupted execution have evidence.
- Linux, macOS, and Windows semantics are equivalent where supported.
- Security/architecture review and owner approval are recorded separately.

## 7. Definition of Ready

Implementation remains blocked until the owner and named reviewer confirm:

- This plan is reconciled with product, architecture, roadmap, and v0.3 evidence.
- The execution package and requirements are complete.
- Each implementation issue maps to one frozen contract and has failure behavior, acceptance evidence, reviewer, and dependencies.
- Isolation, approval, rollback, process, verification, rescan, and patch-evidence decisions are accepted.
- Required fixtures, hostile-input scenarios, and cross-platform matrix are named.
- No issue silently expands into v0.5, web, team, MCP, or autonomous mutation scope.

## 8. Rollback and displaced work

If the preview fails, do not mutate the source repository or move the release tag. Preserve evidence, discard isolated workspaces, and publish a corrected immutable preview only after the failed contract is fixed. v0.5 web/history and v0.7 agent write capabilities remain unchanged and deferred.
