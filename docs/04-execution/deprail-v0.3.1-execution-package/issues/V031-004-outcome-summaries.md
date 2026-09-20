# V031-004 Implement Truthful Outcome Summaries

## Planning metadata

- Type: feature
- Area: cli
- Priority: P1
- Risk: R2
- Epic: EPIC-002
- Dependencies: V031-001, V031-003

## Goal
Make complete, partial, failed, and cancelled results visibly distinct across CLI commands.

## Acceptance criteria

- [ ] Only complete zero-finding scans render “No known vulnerabilities found.”
- [ ] Partial output states incomplete status and affected scope safely.
- [ ] Failed output states failure and stable error code where applicable.
- [ ] Cancellation states cancellation and restores terminal state.
- [ ] `diff` and `fix plan` preserve explicit errors and incomplete evidence.
- [ ] Equivalent domain results produce deterministic summaries.

## Evidence
Complete/partial/failed/cancelled fixtures, presenter contract tests, and manual captures.

## Exclusions
No changes to scan semantics, exit codes, or remediation execution.

## GitHub tracking

- Issue: [#256](https://github.com/geoffrey-xiao/deprail/issues/256)
- Parent epic: [#242](https://github.com/geoffrey-xiao/deprail/issues/242)
