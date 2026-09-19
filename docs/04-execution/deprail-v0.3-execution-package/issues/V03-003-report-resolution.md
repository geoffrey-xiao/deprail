# V03-003 Implement Explicit Report and Finding Resolution

## Planning metadata

- Type: feature
- Area: cli
- Priority: P1
- Risk: R2
- Epic: EPIC-002
- Dependencies: V03-001, v0.3 CLI/data/error contracts

## Goal
Resolve one finding only from an explicit normalized scan report while preserving repository, scan, workspace, and artifact provenance.

## Acceptance criteria

- [ ] `--report` and `--finding` are required.
- [ ] Missing, malformed, incompatible, stale, and ambiguous inputs fail explicitly.
- [ ] No latest-report guessing or implicit rescan occurs.
- [ ] Finding resolution preserves source scan ID, repository state, workspace, and artifact digests.
- [ ] Stable error codes and stdout/stderr behavior are preserved.
- [ ] Tests cover multiple reports and duplicate finding identifiers.

## Exclusions
No report index, implicit discovery, network fetch, or repository mutation.
## GitHub tracking

- Issue: [#208](https://github.com/geoffrey-xiao/deprail/issues/208)
- Parent epic: [#203](https://github.com/geoffrey-xiao/deprail/issues/203)
