# v0.3.1 Test Strategy

## Unit

Test event ordering, capability selection, quiet/verbose precedence, width and monochrome decisions, outcome summaries, escaping, redaction, and terminal cleanup.

## Contract

Capture stdout and stderr independently. Assert JSON purity, deterministic machine output, TTY-only progress, non-TTY line behavior, stable error codes, and truthful complete/partial/failed states.

## Golden

Cover `doctor`, `discover`, `scan`, and `fix plan` with findings, zero findings, partial results, failures, cancellation, quiet, verbose, monochrome, narrow width, and hostile labels. Avoid pinning incidental ANSI or library internals.

## Integration and smoke

Exercise real CLI paths with controlled scanner fixtures and cancellation. Verify scanner arguments, repository tree, artifact behavior, and exit codes are unchanged. Capture interactive and redirected output.

## Cross-platform

Run Linux, macOS, and Windows smoke scenarios. Platform-specific styling may differ, but serialized meaning, status, safety behavior, and stream contracts may not.
