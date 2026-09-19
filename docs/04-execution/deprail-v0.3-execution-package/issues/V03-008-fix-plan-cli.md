# V03-008 Wire Read-Only Fix-Plan Service and CLI

## Planning metadata

- Type: feature
- Area: cli
- Priority: P1
- Risk: R2
- Epic: EPIC-004
- Dependencies: V03-001 through V03-005

## Goal
Expose deterministic remediation planning through `deprail fix plan` without mutation.

## Acceptance criteria

- [ ] Command requires `--report` and `--finding`.
- [ ] Missing, invalid, stale, unsupported, incomplete, and ambiguous inputs use stable errors.
- [ ] Terminal and JSON modes preserve stdout/stderr contracts.
- [ ] Planning never edits the target tree or executes external package-manager operations.
- [ ] Exit behavior is reconciled with existing v0.2 codes.
- [ ] End-to-end fixtures cover all supported ecosystems.

## Exclusions
No `fix apply`, worktree, verification execution, PR creation, or publication.
## GitHub tracking

- Issue: [#213](https://github.com/geoffrey-xiao/deprail/issues/213)
- Parent epic: [#205](https://github.com/geoffrey-xiao/deprail/issues/205)
