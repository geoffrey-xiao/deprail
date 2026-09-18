# V02-013 Record Release Artifact and Supply-Chain Evidence

- Type: docs
- Area: docs
- Priority: P0
- Risk: R2
- Target version: 0.2.0
- Milestone: v0.2.0
- Sprint: Sprint 0
- Owner: TBD
- Reviewer: TBD
- Dependencies: V02-004
- Status: Local planning; GitHub issue not created

## Definition of Ready

- [ ] Required artifact platforms are confirmed.
- [ ] Checksum, manual-scan, SBOM, and signing evidence fields are defined.
- [ ] Preview exception policy is understood.
- [ ] Owner and reviewer are assigned.

## Goal

Create a durable release evidence record that accurately states artifact and supply-chain readiness.

## Scope

Define the evidence record format and capture artifact names, commit/tag, checksums, platform smoke results, scanner version, representative repository scans, SBOM status, signing status, reviewer, and remaining risks.

## Out of Scope

Implementing SBOM or signing, changing release infrastructure, or claiming unavailable controls.

## Inputs, Outputs, and Failure Behavior

Missing or inconsistent evidence blocks the relevant release mode. An unavailable SBOM or signature is recorded as a gap rather than implied to exist.

## Required Tests

- Evidence-record completeness review.
- Checksum/artifact cross-check.
- Platform and manual-scan record review.
- Supply-chain gap wording review.

## Acceptance Criteria

- [ ] Every required artifact has a checksum and smoke result.
- [ ] Manual repository scan evidence is linked.
- [ ] SBOM and signing status are explicit.
- [ ] Remaining risks and owner decision are recorded.

## Human Review

Release and security reviewers inspect evidence accuracy.

## Evidence Required

Completed release evidence record and linked artifact/checksum/smoke outputs.

## Final Acceptance

- [ ] Owner reviewed every criterion during PR review.
- [ ] Release and security evidence were reviewed.
- [ ] Remaining risk is recorded.
- [ ] PR review and merge evidence are linked.
