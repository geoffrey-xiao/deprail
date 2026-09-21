# DepRail v0.4.0 Execution Package

**Version:** `v0.4.0-preview.1`
**Status:** Context package draft; implementation remains blocked pending owner review
**Parent epic:** Not finalized; candidate epic #299 exists for review
**Plan:** [`docs/03-planning/deprail-development-plan-v0.4.0.md`](../../03-planning/deprail-development-plan-v0.4.0.md)

## Source order

1. Product design and architecture.
2. Whole-project roadmap.
3. v0.4 development plan.
4. This execution package and its requirements.
5. Approved issue contracts.
6. Implementation and evidence.

Lower-level issue contracts must not contradict the plan or architecture. Any change to mutation, approval, process, path, compatibility, schema, or publication boundaries requires explicit change control.

## Release boundary

```text
v0.3: finding -> candidate analysis -> reviewable plan
v0.4: approved plan -> isolated apply -> verify -> rescan -> patch evidence
```

## Required contracts

The following documents are the normative v0.4 context drafts:

- [`PRD-v0.4.md`](PRD-v0.4.md)
- [`ARCHITECTURE-v0.4.md`](ARCHITECTURE-v0.4.md)
- [`requirements/FUNCTIONAL-REQUIREMENTS.md`](requirements/FUNCTIONAL-REQUIREMENTS.md)
- [`requirements/FAILURE-AND-DATA-CONTRACT.md`](requirements/FAILURE-AND-DATA-CONTRACT.md)
- [`requirements/TEST-STRATEGY.md`](requirements/TEST-STRATEGY.md)
- [`requirements/ISOLATION-AND-ROLLBACK.md`](requirements/ISOLATION-AND-ROLLBACK.md)

The contract set covers:

| Contract | Required decision | Evidence |
| --- | --- | --- |
| Mutation boundary | Only approved plan operations in isolated workspace | Path and tree comparison |
| Isolation | Temporary worktree/equivalent with restrictive permissions and cleanup | Workspace lifecycle record |
| Approval | Approval bound to plan digest, source commit, target, expiry | Approval contract test |
| Atomic writes | Temp file plus atomic rename; no source-tree writes | Interrupted-write test |
| Rollback | Discard/restore on failure, timeout, cancellation, or regression | Rollback evidence |
| Package manager | Direct argv, bounded process, scripts/network denied by default | Process contract tests |
| Limits | Deadline, cancellation, output cap, process-tree termination | Failure matrix |
| Verification | Supported test/build/type-check discovery only | Command-selection evidence |
| Rescan | Same scan contract; stable finding transition classification | Before/after scan evidence |
| Patch evidence | Diff, digests, commands, tools, outcomes, redacted diagnostics | Schema/golden evidence |
| Failure | Explicit complete/partial/failed/cancelled states | Partial-failure fixtures |

## Required scenario matrix

- Valid approved plan.
- Invalid or stale plan.
- Approval mismatch, expiry, and cancellation.
- External symlink and traversal escape.
- Shell metacharacters and hostile package names.
- Package-manager missing, incompatible, non-zero, timeout, malformed output, and output limit.
- Install-script and network denial.
- Interrupted atomic write.
- Partial mutation followed by rollback.
- Test/build/type-check success and failure.
- Resolved, residual, introduced, and unknown findings after rescan.
- Original repository tree unchanged.
- Linux, macOS, and Windows semantic equivalence.

## Definition of Ready checklist

The v0.4 Definition of Ready and ADR-0002 approval are already recorded in
[`tracking/MASTER-CHECKLIST.md`](tracking/MASTER-CHECKLIST.md). This package
does not reset those completed records.

The supplemental context review required before implementation is:

- [ ] PRD and architecture context accepted.
- [ ] Functional requirements accepted.
- [ ] Failure/data and redaction contract accepted.
- [ ] Test strategy and cross-platform evidence plan accepted.
- [ ] Child issue contracts derived from the accepted documents.
- [ ] Owner and named architecture/security reviewer confirm the complete package.

## No implementation authorization

This package does not authorize code changes. Runtime implementation begins only after every Definition of Ready item is checked with linked evidence and owner/reviewer approval. No issue may add repository mutation, publication, autonomous approval, web/team services, or agent write scope without a new decision record.
