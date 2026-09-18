# V02-020 Expiring Exceptions

## Planning Metadata

- Type: feature
- Area: foundation
- Priority: P1
- Risk: R2
- Target version: 0.2.0
- Sprint: S6
- Dependencies: V02-019
- Parent epic: EPIC-007

## Definition of Ready

- [ ] Contract inputs, outputs, and failure behavior are confirmed.
- [ ] Dependencies and acceptance evidence are explicit.
- [ ] Compatibility and security boundaries are reviewed.

## Goal

Represent temporary risk acceptance without turning an ignored finding into a permanent safe result.

## Scope

Exception schema, finding and scope identity, rationale, approver, creation time, expiry, review condition, and expired behavior.

## Out of Scope

Remote approval workflows, hosted storage, identity/RBAC, and automatic exception creation.

## Inputs, Outputs, and Failure Behavior

Valid unexpired exceptions affect only their declared scope. Missing, malformed, expired, or mismatched exceptions remain visible and cannot suppress policy failure.

## Required Tests

Schema, scope matching, expiry boundaries, malformed records, deterministic evaluation, and unknown finding identity.

## Acceptance Criteria

- Required exception fields are versioned and validated.
- Expiry is deterministic and enforced.
- Exceptions cannot hide unrelated findings or incomplete scans.
- Expired exceptions reappear in decisions and output.

## Human Review

Review security, auditability, and time-bound semantics.

## Evidence Required

Schema, examples, boundary tests, and decision transcripts.
## Final Acceptance

- [ ] Acceptance evidence is linked.
- [ ] CI and human review are complete.
- [ ] Remaining risk is recorded.
