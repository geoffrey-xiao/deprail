# V03-005 Define Read-Only Planning Adapter Contract

## Planning metadata

- Type: feature
- Area: adapter
- Priority: P1
- Risk: R2
- Epic: EPIC-003
- Dependencies: V03-001, V03-002

## Goal
Define one adapter port for dependency ownership, constraints, candidates, risks, commands, expected files, and verification discovery.

## Acceptance criteria

- [ ] Contract is independent of package-manager and scanner types.
- [ ] Adapter methods are read-only and return normalized domain evidence.
- [ ] Commands include executable, arguments, and working directory.
- [ ] Unsupported, unavailable, unknown, and rejected states are explicit.
- [ ] Contract tests run with offline fixtures.
- [ ] No adapter can execute package-manager commands in v0.3.

## Exclusions
No ecosystem implementation, network service, or package installation.
## GitHub tracking

- Issue: [#210](https://github.com/geoffrey-xiao/deprail/issues/210)
- Parent epic: [#204](https://github.com/geoffrey-xiao/deprail/issues/204)
