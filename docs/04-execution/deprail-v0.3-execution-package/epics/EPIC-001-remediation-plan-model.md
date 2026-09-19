# EPIC-001 Remediation Plan Model

## Outcome
A deterministic, schema-valid domain model represents reviewable remediation plans without repository mutation.

## Scope

- Plan and candidate entities.
- Stable plan identity and canonicalization.
- Structured commands with repository-relative working directories.
- Risk, assumptions, verification, rollback, and provenance fields.
- JSON Schema, examples, converters, and compatibility tests.

## Issues

- V03-001: Define remediation-plan domain model and stable identity.
- V03-002: Add versioned remediation-plan schema and golden examples.

## Acceptance

Plans are deterministic, schema-valid, provenance-preserving, explicit about unknown evidence, and contain no executable shell fragments.

## Exclusions

No CLI, package-manager adapter, file mutation, command execution, or network metadata.
