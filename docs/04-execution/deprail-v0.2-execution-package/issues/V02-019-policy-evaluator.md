# V02-019 Typed Policy Evaluation

- Type: feature
- Area: foundation
- Priority: P0
- Risk: R2
- Target version: 0.2.0
- Sprint: S6
- Dependencies: V02-017
- Parent epic: EPIC-007

## Goal

Evaluate scan and diff results using deterministic typed policy rules.

## Scope

New-risk gates, severity thresholds, required-complete-scan behavior, scanner failure handling, and policy decision output.

## Out of Scope

General policy languages, remote policy services, remediation, or hidden risk scoring.

## Inputs, Outputs, and Failure Behavior

Policy decisions are explicit pass, warn, or block outcomes. Incomplete scans, scanner failures, unknown data, or invalid policy must not become a pass.

## Required Tests

Rule boundaries, new/high/critical findings, incomplete scans, scanner failures, unknown values, and deterministic repeated evaluation.

## Acceptance Criteria

- Policy behavior is typed and deterministic.
- New-risk and severity gates are explicit.
- Incomplete or failed scans cannot pass silently.
- Decision evidence identifies the applied rules and inputs.

## Human Review

Security and compatibility review of default policy behavior.

## Evidence Required

Policy fixtures, decision JSON, exit mapping, and failure transcripts.
