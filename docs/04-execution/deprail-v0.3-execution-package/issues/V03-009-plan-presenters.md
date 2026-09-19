# V03-009 Add Plan Presenters and Safe Persistence

## Planning metadata

- Type: feature
- Area: cli
- Priority: P1
- Risk: R2
- Epic: EPIC-004
- Dependencies: V03-002, V03-008

## Goal
Render concise human plans and stable machine JSON, with safe external plan persistence.

## Acceptance criteria

- [ ] Terminal output shows finding, recommendation/no recommendation, risks, files, future commands, verification, provenance, and read-only notice.
- [ ] JSON validates against the v0.3 schema.
- [ ] Equivalent plans serialize identically.
- [ ] Output inside target root is rejected.
- [ ] External output is atomic and non-overwriting by default.
- [ ] Diagnostics never contaminate machine stdout.

## Exclusions
No web presenter, database history, source upload, or mutation.
