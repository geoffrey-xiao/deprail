# V03-004 Enforce Plan Output and Repository Safety

## Planning metadata

- Type: security
- Area: foundation
- Priority: P1
- Risk: R3
- Epic: EPIC-002
- Dependencies: V03-001, V03-003

## Goal
Prove planning cannot mutate the target repository and all paths/outputs remain within approved boundaries.

## Acceptance criteria

- [ ] Output paths inside the target root are rejected.
- [ ] Traversal and symlink escapes are rejected.
- [ ] External output is atomic, restrictive, and non-overwriting by default.
- [ ] Planning performs no package-manager execution, script execution, or writes to target files.
- [ ] Hostile path, manifest, command, and metadata tests pass.
- [ ] Before/after repository snapshots are identical after successful and failed planning.

## Evidence
Security tests, process audit, filesystem snapshots, and platform smoke results.
