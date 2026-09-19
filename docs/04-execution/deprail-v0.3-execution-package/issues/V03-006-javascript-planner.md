# V03-006 Implement JavaScript Planning Adapter

## Planning metadata

- Type: feature
- Area: adapter
- Priority: P1
- Risk: R2
- Epic: EPIC-003
- Dependencies: V03-005

## Goal
Plan npm, pnpm, and Yarn dependency upgrades from manifests, lockfiles, and retained scan evidence without mutation.

## Acceptance criteria

- [ ] Direct and transitive ownership are distinguished.
- [ ] Semver constraints, lockfile locations, peer/engine risks, and major jumps are represented.
- [ ] Minimum viable non-vulnerable candidates are ordered deterministically.
- [ ] Future package-manager commands include workspace working directory.
- [ ] Malformed manifests, conflicting lockfiles, scripts, traversal, and missing metadata fail safely.
- [ ] npm/pnpm/Yarn fixtures and golden plans pass.

## Exclusions
No install, update, script, lockfile rewrite, or network fallback.
