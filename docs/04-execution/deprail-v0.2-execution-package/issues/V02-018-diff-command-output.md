# V02-018 Deterministic Diff Output and deprail diff

- Type: feature
- Area: cli
- Priority: P0
- Risk: R1
- Target version: 0.2.0
- Sprint: S5
- Dependencies: V02-017
- Parent epic: EPIC-006

## Goal

Expose baseline comparison through a stable CLI command and machine-readable output.

## Scope

`deprail diff`, base/head inputs, terminal output, JSON output, exit behavior, stdout/stderr separation, and path normalization.

## Out of Scope

Policy enforcement, SARIF, GitHub Action packaging, and remote baseline retrieval.

## Inputs, Outputs, and Failure Behavior

Valid compatible inputs produce deterministic diff output. Invalid arguments, missing baselines, incompatible inputs, or incomplete scans return documented errors and never claim no change.

## Required Tests

CLI contract, JSON schema, terminal output, invalid input, incomplete input, exit codes, and stdout/stderr separation.

## Acceptance Criteria

- `deprail diff` supports documented base/head forms.
- JSON is schema-valid and deterministic.
- Invalid or unsafe comparisons fail before a safe result is emitted.
- Exit behavior is documented and tested.

## Human Review

Review CLI compatibility and diff interpretation.

## Evidence Required

Command transcripts, JSON examples, exit-code matrix, and compatibility note.
