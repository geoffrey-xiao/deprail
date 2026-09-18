# V02-025 Pull-Request Validation

- Type: feature
- Area: test
- Priority: P0
- Risk: R2
- Target version: 0.2.0
- Sprint: S7
- Dependencies: V02-023, V02-024
- Parent epic: EPIC-008

## Goal

Prove the guardrail on real base/head pull-request scenarios.

## Scope

Representative repositories, baseline setup, dependency additions and upgrades, new/resolved findings, policy pass/block, incomplete scan, SARIF artifact, and pull-request summary.

## Out of Scope

Automatic dependency changes, merge, hosted dashboards, or non-DepRail security scanning.

## Inputs, Outputs, and Failure Behavior

A pull request receives a truthful summary and evidence artifact. New policy violations block according to contract; historical findings do not block unless policy says so; incomplete scans never pass silently.

## Required Tests

Controlled real pull-request fixtures or repository runs, fork behavior, policy matrix, artifact retrieval, and failure-path validation.

## Acceptance Criteria

- New-risk classification works on a real base/head pair.
- Policy pass and block behavior is observable.
- SARIF and full evidence are retrievable.
- Incomplete and scanner-failure cases remain unsafe.

## Human Review

Manual pull-request validation and security review.

## Evidence Required

Pull-request URLs or controlled transcripts, artifacts, exit codes, and remaining-risk record.
