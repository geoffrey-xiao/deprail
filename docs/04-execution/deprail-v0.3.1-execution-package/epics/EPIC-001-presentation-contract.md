# EPIC-001 Presentation Contract

## Outcome
A renderer-independent presentation contract describes real application lifecycle events and terminal capabilities without importing presentation concerns into domain or scanner code.

## Scope

- Structured lifecycle event vocabulary and ordering.
- TTY, CI, color, width, and animation capability detection.
- Output-mode selection and cancellation lifecycle.
- Dependency boundary for terminal libraries.

## Issues

- V031-001: Define structured presentation events and lifecycle semantics.
- V031-002: Add terminal capability and output-mode selection.

## Acceptance

Application services emit deterministic structured events; renderers consume them without domain coupling; TTY/non-TTY and cancellation decisions are explicit and testable.

## Exclusions

No visual styling, scanner changes, repository mutation, full-screen TUI, or machine-schema changes.

## GitHub tracking

- Parent issue: Pending creation
- Child issues: Pending creation
