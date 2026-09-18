# V02-022 SARIF Output

- Type: feature
- Area: cli
- Priority: P0
- Risk: R1
- Target version: 0.2.0
- Sprint: S6
- Dependencies: V02-018, V02-021
- Parent epic: EPIC-007

## Goal

Export normalized findings and policy decisions in interoperable SARIF without losing DepRail provenance.

## Scope

SARIF mapping, rule/result identity, locations, severity, remediation metadata, policy status, schema validation, and stdout/output-file behavior.

## Out of Scope

SARIF upload services, hosted dashboards, or consumer-specific undocumented extensions.

## Inputs, Outputs, and Failure Behavior

Valid reports produce SARIF. Unknown or incomplete evidence remains represented as such. Mapping or output failures are explicit and never silently produce an empty successful SARIF file.

## Required Tests

SARIF schema validation, finding mapping, policy mapping, empty output, incomplete output, deterministic serialization, and atomic file output.

## Acceptance Criteria

- SARIF validates against the supported version.
- Findings retain stable identity and source provenance.
- Empty and incomplete results remain distinguishable.
- Output is deterministic and diagnostics stay off machine output.

## Human Review

Review standards compatibility and evidence preservation.

## Evidence Required

SARIF fixtures, validator output, deterministic comparison, and consumer smoke result.
