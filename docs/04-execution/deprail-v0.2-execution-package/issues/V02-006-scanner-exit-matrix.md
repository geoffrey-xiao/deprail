# V02-006 Add Scanner Exit-Code Matrix

- Type: test
- Area: adapter
- Priority: P0
- Risk: R1
- Target version: 0.2.0
- Milestone: v0.2.0
- Sprint: Sprint 0
- Owner: TBD
- Reviewer: TBD
- Dependencies: V02-005
- Status: Local planning; GitHub issue not created

## Definition of Ready

- [ ] Supported command and version are fixed.
- [ ] Exit meanings are sourced from the real scanner contract.
- [ ] Process runner and adapter boundaries are identified.
- [ ] Owner and reviewer are assigned.

## Goal

Ensure OSV-Scanner exit codes distinguish findings from execution failure.

## Scope

Define and test the matrix for zero findings, vulnerabilities found, malformed output, unsupported version, missing tool, timeout, cancellation, and unexpected non-zero exit.

## Out of Scope

Changing the process runner API or adding retries.

## Inputs, Outputs, and Failure Behavior

A vulnerability-found exit with valid JSON produces a trustworthy report with findings. Tool failure, timeout, malformed output, and unsupported output produce explicit failure/partial outcomes and never safe success.

## Required Tests

- Controlled child-process exit fixtures.
- Valid findings with non-zero exit.
- Invalid output with zero/non-zero exit.
- Timeout and cancellation.
- Missing tool and unsupported version.

## Acceptance Criteria

- [ ] Every supported exit case has a documented result.
- [ ] Vulnerability-found exit is accepted only with valid output.
- [ ] Unexpected exits remain scanner failures.
- [ ] Exit behavior is stable across supported platforms.

## Human Review

Adapter and process maintainers review safety of each classification.

## Evidence Required

Exit matrix, controlled process fixtures, platform test output, and error-code mapping.

## Final Acceptance

- [ ] Owner reviewed every criterion during PR review.
- [ ] CI and adapter evidence were reviewed.
- [ ] Remaining risk is recorded.
- [ ] PR review and merge evidence are linked.
