# V03-002 Add Versioned Remediation-Plan Schema

## Planning metadata

- Type: feature
- Area: normalization
- Priority: P1
- Risk: R2
- Epic: EPIC-001
- Dependencies: V03-001, v0.3 data contract

## Goal
Publish the machine-readable remediation-plan schema, examples, and compatibility tests.

## Acceptance criteria

- [ ] Schema includes plan identity, source report, finding, candidates, recommendation, files, commands, risks, assumptions, verification, rollback, and provenance.
- [ ] Structured commands require a repository-relative working directory.
- [ ] Examples cover recommended, multiple-candidate, no-recommendation, unknown, rejected, stale, and incomplete plans.
- [ ] Schema validation rejects missing provenance and unsafe command paths.
- [ ] Golden serialization is deterministic.
- [ ] Schema version and compatibility notes are documented.

## Evidence
Schema validation output, examples, golden tests, and compatibility review.

## Exclusions
No breaking changes to v0.2 schemas and no executor implementation.
## GitHub tracking

- Issue: [#207](https://github.com/geoffrey-xiao/deprail/issues/207)
- Parent epic: [#202](https://github.com/geoffrey-xiao/deprail/issues/202)
