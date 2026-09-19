# EPIC-002 Report Input and Safety Boundaries

## Outcome
Planning consumes an explicit normalized scan report and enforces repository/output containment without guessing or mutating.

## Scope

- Explicit `--report` and `--finding` resolution.
- Report schema/provenance validation.
- Stale and ambiguous finding rejection.
- Canonical root and symlink containment.
- External output boundary and atomic writes.
- Mutation-attempt and hostile-input tests.

## Issues

- V03-003: Implement explicit scan-report and finding resolution.
- V03-004: Enforce plan output and repository safety boundaries.

## Acceptance

A finding resolves only from the supplied report; missing, invalid, stale, or ambiguous reports fail with stable diagnostics. In-repository output and traversal are rejected. Planning leaves the target tree unchanged.
## GitHub tracking

- Parent issue: [#203](https://github.com/geoffrey-xiao/deprail/issues/203)
- Child issues: [#208](https://github.com/geoffrey-xiao/deprail/issues/208), [#209](https://github.com/geoffrey-xiao/deprail/issues/209)
