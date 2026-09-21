# V04-008: Implement Isolated Fix Apply Orchestration

- GitHub Issue: [#332](https://github.com/geoffrey-xiao/deprail/issues/332)
- Epic: [EPIC-006](../epics/EPIC-006-fix-apply.md)
- Target: `v0.4.0-preview.2`
- Status: Todo; child issues created and tracked in the epic
- Priority: P0
- Risk: R3
- Area: CLI
- Reviewer: Project owner plus independent architecture/security reviewer

## Value

Users need a safe, reviewable path from an approved `deprail fix plan` to a verified remediation result without mutating the caller worktree or hiding failure.

## Scope

Implement `deprail fix apply` with plan validation, approval/source binding, dry-run, isolated detached worktree execution, allowlisted adapter mutations, bounded verification, rescan, transition classification, rollback/discard, cleanup evidence, and deterministic terminal/JSON output.

## Exclusions

No caller-worktree mutation, automatic commit/push/PR/merge, autonomous approval, arbitrary shell commands, unbounded scripts/network access, new ecosystems, remote publishing, policy-language expansion, or release publication.

## Dependencies

- `deprail fix plan` and versioned plan schema.
- V04-001/V04-002/V04-003: isolation, approval, and mutation boundaries.
- V04-004/V04-005/V04-006: verification, evidence, and cross-platform semantics.
- V04-007: trusted baseline generation.

## Child issues

- [#337 — Validate apply plans and approval binding](https://github.com/geoffrey-xiao/deprail/issues/337)
- [#338 — Implement isolated detached worktree lifecycle](https://github.com/geoffrey-xiao/deprail/issues/338)
- [#336 — Implement allowlisted adapter mutations](https://github.com/geoffrey-xiao/deprail/issues/336)
- [#340 — Add bounded verification and rescan orchestration](https://github.com/geoffrey-xiao/deprail/issues/340)
- [#339 — Classify remediation finding transitions](https://github.com/geoffrey-xiao/deprail/issues/339)
- [#345 — Implement rollback, discard, and cleanup evidence](https://github.com/geoffrey-xiao/deprail/issues/345)
- [#346 — Implement apply evidence and deterministic output](https://github.com/geoffrey-xiao/deprail/issues/346)
- [#347 — Add security and cross-platform apply coverage](https://github.com/geoffrey-xiao/deprail/issues/347)

## Acceptance criteria

1. Valid approved plans apply only inside a temporary detached worktree; caller tree status and digest remain unchanged.
2. Dry-run validates and reports without mutation.
3. Stale, unapproved, malformed, unsupported, path-escaping, timeout, cancellation, process, verification, rescan, and cleanup failures produce explicit stable results.
4. Verification and rescan classify resolved, unchanged, introduced, and residual findings deterministically.
5. Rollback/discard is idempotent and evidenced across success and failure paths.
6. Security coverage includes hostile paths, symlinks, shell metacharacters, output limits, credential-bearing values, process-tree termination, and interrupted cleanup.
7. Linux/macOS/Windows contract and smoke evidence is attached.

## Required evidence

Record exact reviewed source commit, binary identity, commands, exit codes, artifact paths, repository-tree comparison, checksums, cleanup/rollback outcomes, and remaining risk. Obtain independent architecture/security review before closure.
