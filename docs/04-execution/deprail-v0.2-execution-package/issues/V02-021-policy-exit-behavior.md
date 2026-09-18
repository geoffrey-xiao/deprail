# V02-021 Policy Exit Behavior

## Planning Metadata

- Type: feature
- Area: cli
- Priority: P0
- Risk: R2
- Target version: 0.2.0
- Sprint: S6
- Dependencies: V02-019, V02-020
- Parent epic: EPIC-007

## Definition of Ready

- [ ] Contract inputs, outputs, and failure behavior are confirmed.
- [ ] Dependencies and acceptance evidence are explicit.
- [ ] Compatibility and security boundaries are reviewed.

## Goal

Expose stable CI exit behavior for policy decisions without conflating configuration, scanner, and policy failures.

## Scope

`deprail policy check`, pass/warn/block results, incomplete scans, expired exceptions, invalid policy, and machine-output separation.

## Out of Scope

GitHub Action packaging, remote policy, or automatic remediation.

## Inputs, Outputs, and Failure Behavior

Policy pass, warning, and block outcomes map to documented exit behavior. Invalid configuration returns code 2; unusable scan results remain code 3; policy failure uses the approved stable code.

## Required Tests

Command contract, exit matrix, JSON output, stderr separation, invalid policy, incomplete scan, and expired exception cases.

## Acceptance Criteria

- Exit behavior is stable and documented.
- Policy failure is distinguishable from scanner failure and configuration error.
- No policy check emits a safe pass from unusable input.
- JSON output is schema-valid.

## Human Review

Review CI compatibility and failure precedence.

## Evidence Required

Exit matrix, command transcripts, schema examples, and compatibility statement.
## Final Acceptance

- [ ] Acceptance evidence is linked.
- [ ] CI and human review are complete.
- [ ] Remaining risk is recorded.
