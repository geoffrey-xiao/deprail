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
## Definition of Ready

- [ ] Owner and named reviewer are assigned.
- [ ] v0.3.1 target, sprint, dependencies, and contract mapping are confirmed.
- [ ] Inputs, outputs, failure behavior, and required evidence are reviewed.
- [ ] Implementation remains blocked until release-level Definition of Ready approval.

## Contract mapping

- Requirements: `FUNCTIONAL-REQUIREMENTS.md` UX-005.
- Error model: `ERROR-MODEL.md` complete, partial, failed, and cancelled states.
- Evidence: outcome fixtures and manual captures.

## Inputs, outputs, and failure behavior

- Inputs: completeness, findings, diagnostics, cancellation, and report state.
- Outputs: truthful human summaries and unchanged machine results.
- Failure: incomplete or failed analysis never renders as successful empty output.

## Final acceptance

- [ ] Owner reviewed every criterion.
- [ ] Architecture/security review completed.
- [ ] CI and required evidence links are recorded.
- [ ] Remaining risk and follow-up are documented.
