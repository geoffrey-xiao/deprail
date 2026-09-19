# EPIC-004 CLI and Evidence Delivery

## Outcome
Users can generate, inspect, validate, and safely store remediation plans through the CLI with deterministic terminal/JSON output.

## Scope

- `deprail fix plan --report --finding` command.
- Terminal and JSON presenters.
- External plan output and optional local plan storage.
- Stable errors, stdout/stderr separation, and exit behavior.
- Cross-platform smoke, representative fixtures, and v0.3 release evidence.

## Issues

- V03-008: Wire the read-only fix-plan application service and CLI.
- V03-009: Add terminal/JSON presenters and safe plan persistence.
- V03-010: Complete cross-platform evidence and v0.3 release gate.

## Acceptance

The command generates a valid plan from an explicit report, never mutates the repository, remains deterministic across supported platforms, and publishes complete evidence for the v0.3 decision.
## GitHub tracking

- Parent issue: [#205](https://github.com/geoffrey-xiao/deprail/issues/205)
- Child issues: [#213](https://github.com/geoffrey-xiao/deprail/issues/213), [#214](https://github.com/geoffrey-xiao/deprail/issues/214), [#215](https://github.com/geoffrey-xiao/deprail/issues/215)
