# EPIC-004: Patch Evidence and Failure Reporting

**Status:** Proposed

## Outcome

Every apply attempt produces redaction-safe, content-addressed evidence that distinguishes success, partial completion, failure, cancellation, and cleanup failure.

## Boundaries

Includes versioned JSON, stable diagnostics, command/tool identity, before/after digests, diff metadata, artifact retention, redaction, and terminal outcomes. Excludes remote publishing and user credential storage.

## Acceptance

Evidence validates against its schema, retains successful partial evidence, excludes secrets, and never represents an incomplete result as a successful empty scan.
