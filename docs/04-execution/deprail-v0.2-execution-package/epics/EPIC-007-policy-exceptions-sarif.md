# EPIC-007 Policy, Exceptions, and SARIF

## Outcome

Evaluate new dependency risk deterministically, represent expiring exceptions safely, and emit interoperable SARIF evidence.

## Scope

Typed policy evaluation, completeness and severity gates, new-finding decisions, expiring exceptions, stable policy exit behavior, `deprail policy check`, and SARIF output.

## Out of Scope

A general policy language, remote exception approval, hosted services, automatic remediation, or repository mutation.

## Child issues

- V02-019 Typed policy evaluation.
- V02-020 Expiring exceptions.
- V02-021 Policy exit behavior.
- V02-022 SARIF output.

## Dependencies

EPIC-006 diff classifications and the prerequisite normalized scan contract.

## Acceptance

- New/high/critical findings and incomplete scans follow explicit policy rules.
- Exceptions include scope, rationale, approver, creation time, expiry, and review condition.
- Expired exceptions reappear and cannot silently make a scan safe.
- Policy exits are stable and machine-output remains separate from diagnostics.
- SARIF output validates against supported consumers and preserves provenance.
