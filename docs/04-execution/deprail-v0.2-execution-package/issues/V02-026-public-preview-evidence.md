# V02-026 Public Preview Release Evidence

- Type: docs
- Area: docs
- Priority: P0
- Risk: R2
- Target version: 0.2.0
- Sprint: S7
- Dependencies: V02-023, V02-024, V02-025
- Parent epic: EPIC-008

## Goal

Record whether the CI guardrail is ready for the v0.2 public preview.

## Scope

Preview artifact identity, checksums, platform smoke, Action validation, representative repositories, documentation, known limitations, rollback, SBOM/signing status, and owner decision.

## Out of Scope

Production hosted service, automatic remediation, signing implementation, or new feature scope after freeze.

## Inputs, Outputs, and Failure Behavior

The evidence record states pass, gap, or blocked for every gate. Missing evidence prevents a stable preview claim and records the remaining risk.

## Required Tests

Cross-platform CI, installation smoke, real pull-request validation, representative scans, checksum verification, and documentation walkthrough.

## Acceptance Criteria

- Every v0.2 release gate has linked evidence.
- Artifact, CLI, tag, and commit identity agree.
- Known SBOM/signing gaps are explicit.
- Rollback and support limitations are documented.
- Owner records the public-preview decision.

## Human Review

Release and security review are required.

## Evidence Required

Artifact inventory, checksums, CI links, pull-request evidence, smoke output, known-risk record, and decision.
