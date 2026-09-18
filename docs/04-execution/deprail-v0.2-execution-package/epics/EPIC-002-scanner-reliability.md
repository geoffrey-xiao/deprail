# EPIC-002 Scanner Reliability

## Outcome

Make OSV-Scanner v2 behavior reproducible from checked-in evidence and safe failure handling.

## Scope

- Real v2 raw-output fixtures.
- Exit-code compatibility matrix.
- Malformed output, timeout, missing tool, and incompatible version evidence.
- Scanner working-directory and requested-root assertions.

## Acceptance

The adapter contract tests distinguish findings, successful empty output, scanner failure, malformed output, timeout, and unsupported versions without false-safe results.

## Issues

- `V02-005` Add real OSV-Scanner v2 fixtures.
- `V02-006` Add explicit scanner exit-code matrix.
- `V02-007` Verify scanner execution from requested root.
