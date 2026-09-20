# v0.3.1 Manual Test Guide

Use fixed fixtures and controlled scanner behavior. Record platform, terminal, width, color mode, command, stdout, stderr, exit code, and repository-tree hash.

## Scenarios

1. Run a complete scan with findings in an interactive TTY.
2. Run a complete scan with zero findings and confirm the explicit successful summary.
3. Produce a partial scan and confirm it cannot appear empty or successful.
4. Produce a scanner failure and confirm stable error guidance.
5. Redirect stdout and stderr and confirm no animation or ANSI control sequences.
6. Run `--format json` and parse stdout independently from stderr.
7. Run quiet and verbose modes and compare their documented differences.
8. Use narrow terminal width and monochrome output.
9. Inject newlines, tabs, Unicode, ANSI, and long values into workspace and package labels.
10. Cancel a running operation and confirm terminal state restoration.
11. Repeat representative commands on Linux, macOS, and Windows.

## Evidence rule

A scenario passes only when the observed output, stream, exit behavior, repository state, and serialized meaning match the v0.3.1 contracts. Do not treat a visually attractive output as sufficient evidence if it violates machine-output or security invariants.
