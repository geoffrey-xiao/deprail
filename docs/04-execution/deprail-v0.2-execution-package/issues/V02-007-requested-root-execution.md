# V02-007 Verify Scanner Requested-Root Execution

- Type: bug
- Area: adapter
- Priority: P0
- Risk: R2
- Target version: 0.2.0
- Milestone: v0.2.0
- Sprint: Sprint 0
- Owner: TBD
- Reviewer: TBD
- GitHub Issue: [#117](https://github.com/geoffrey-xiao/deprail/issues/117)
- Parent epic: [EPIC-002 / #113](https://github.com/geoffrey-xiao/deprail/issues/113)
- Dependencies: None

## Definition of Ready

- [ ] Canonical root and caller working-directory behavior are defined.
- [ ] Artifact placement contract is identified.
- [ ] Outside-root and containment scenarios are named.
- [ ] Owner and security reviewer are assigned.

## Goal

Ensure scan execution and raw artifacts use the requested repository root, not the caller process directory.

## Scope

Trace requested-root resolution through application, scan plan, process working directory, adapter invocation, and artifact storage. Add containment assertions.

## Out of Scope

Changing repository discovery boundaries or adding remote storage.

## Inputs, Outputs, and Failure Behavior

Running from outside the target scans only the canonical target. Traversal or symlink escape fails with a security diagnostic. Artifact write failure is explicit and cannot produce a complete persisted result.

## Required Tests

- Scan from outside target directory.
- Two repositories with distinguishable manifests.
- Relative and absolute target paths.
- External symlink and traversal escape.
- Artifact path containment.

## Acceptance Criteria

- [ ] Scanner working directory equals requested canonical root.
- [ ] Caller CWD cannot change scan scope.
- [ ] Raw artifacts remain within the approved root.
- [ ] Hostile paths fail safely.

## Human Review

Security reviewer inspects path resolution and process invocation.

## Evidence Required

End-to-end transcript, artifact paths, hostile-path results, and platform test output.

## Final Acceptance

- [ ] Owner reviewed every criterion during PR review.
- [ ] Security review and CI evidence were completed.
- [ ] Remaining risk is recorded.
- [ ] PR review and merge evidence are linked.
