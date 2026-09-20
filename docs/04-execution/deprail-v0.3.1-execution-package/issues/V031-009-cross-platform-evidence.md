# V031-009 Produce Cross-Platform Terminal Evidence

## Planning metadata

- Type: release
- Area: docs
- Priority: P1
- Risk: R2
- Epic: EPIC-004
- Dependencies: V031-004, V031-005, V031-006, V031-007, V031-008

## Goal
Verify equivalent CLI meaning and safe terminal behavior on Linux, macOS, and Windows.

## Acceptance criteria

- [ ] Interactive TTY behavior is captured on each supported platform.
- [ ] Non-TTY, JSON, quiet, verbose, narrow, monochrome, partial, failed, hostile-label, and cancellation cases are covered.
- [ ] Serialized meaning and exit semantics are equivalent across platforms.
- [ ] Platform-specific styling differences are documented and non-semantic.
- [ ] Repository trees, scanner arguments, and artifacts remain unchanged by presentation.

## Evidence
Platform smoke captures, compatibility matrix updates, repository-state checks, and reviewer sign-off.

## Exclusions
No stable-release claim, v0.4 mutation, verification, or publication.

## GitHub tracking
- Issue: [#254](https://github.com/geoffrey-xiao/deprail/issues/254)
- Parent epic: [#241](https://github.com/geoffrey-xiao/deprail/issues/241)
