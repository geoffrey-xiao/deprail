# V02-023 GitHub Action Packaging

## Planning Metadata

- Type: feature
- Area: docs
- Priority: P0
- Risk: R2
- Target version: 0.2.0
- Sprint: S7
- Dependencies: V02-018, V02-021, V02-022
- Parent epic: EPIC-008

## Definition of Ready

- [ ] Contract inputs, outputs, and failure behavior are confirmed.
- [ ] Dependencies and acceptance evidence are explicit.
- [ ] Compatibility and security boundaries are reviewed.

## Goal

Package the v0.2 guardrail as a reviewable GitHub Action without duplicating security logic.

## Scope

Action metadata, pinned or controlled execution, input/output contracts, artifact handling, version selection, and failure propagation.

## Out of Scope

Hosted DepRail services, automatic merge, source upload by default, and unreviewed third-party execution.

## Inputs, Outputs, and Failure Behavior

The Action invokes the CLI and preserves its stable outputs and exit behavior. Missing tools, invalid inputs, policy blocks, and incomplete scans remain explicit failures.

## Required Tests

Action metadata, local runner smoke, input validation, output artifacts, failure propagation, and version reproducibility.

## Acceptance Criteria

- Action inputs and outputs are documented and versioned.
- CLI behavior is reused rather than reimplemented.
- Failure and policy outcomes reach the pull request correctly.
- Action packaging is reproducible and reviewable.

## Human Review

Review supply-chain, permissions, and release behavior.

## Evidence Required

Action metadata, local smoke output, artifact inventory, and rollback procedure.
## Final Acceptance

- [ ] Acceptance evidence is linked.
- [ ] CI and human review are complete.
- [ ] Remaining risk is recorded.
