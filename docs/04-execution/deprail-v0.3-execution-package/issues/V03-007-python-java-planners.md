# V03-007 Implement Python and Java Planning Adapters

## Planning metadata

- Type: feature
- Area: adapter
- Priority: P1
- Risk: R2
- Epic: EPIC-003
- Dependencies: V03-005

## Goal
Plan requirements/uv/Poetry and Maven/Gradle upgrades through the shared read-only adapter contract.

## Acceptance criteria

- [ ] Python constraints, lockfiles, direct/transitive ownership, and runtime risks are represented.
- [ ] Maven/Gradle constraints, dependency paths, lockfiles, and major-version risks are represented.
- [ ] Future commands include canonical workspace working directories.
- [ ] Unsupported syntax, malformed files, missing metadata, and ambiguity fail safely.
- [ ] Python, Maven, and Gradle fixtures produce deterministic plans.
- [ ] No package-manager command or script executes.

## Exclusions
No new scanner family, install, update, build, or rescan.
## GitHub tracking

- Issue: [#212](https://github.com/geoffrey-xiao/deprail/issues/212)
- Parent epic: [#204](https://github.com/geoffrey-xiao/deprail/issues/204)
