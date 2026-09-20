# EPIC-003 Progress, Compatibility, and Safety

## Outcome
Interactive progress reflects real work, CLI flags behave truthfully, and hostile repository-derived labels cannot control terminals or leak sensitive data.

## Scope

- TTY-only progress on stderr from real lifecycle events.
- Non-TTY and JSON stream purity.
- Explicit `--quiet` implementation and `--verbose` compatibility.
- Terminal-safe escaping, redaction, serialization, and cancellation cleanup.

## Issues

- V031-005: Wire real lifecycle progress to the renderer.
- V031-006: Implement quiet mode and preserve verbose behavior.
- V031-007: Harden hostile-label rendering and redaction.
- V031-008: Verify non-TTY, JSON, and cancellation stream behavior.

## Acceptance

Progress never invents work, contaminates machine streams, or survives cancellation incorrectly; quiet and verbose behavior is covered; hostile labels remain data.

## Exclusions

No artificial delays, scanner changes, package-manager execution, repository mutation, or telemetry.

## GitHub tracking

- Parent issue: Pending creation
- Child issues: Pending creation
