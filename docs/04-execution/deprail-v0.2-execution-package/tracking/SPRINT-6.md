# v0.2 Sprint 6 — Policy and SARIF

## Goal

Turn deterministic diff results into safe policy decisions and interoperable evidence.

## Planned items

- V02-019 typed policy evaluation.
- V02-020 expiring exceptions.
- V02-021 policy exit behavior.
- V02-022 SARIF output.

## Dependencies

V02-017 classifications are required. V02-020 depends on V02-019; V02-021 depends on V02-019 and V02-020; V02-022 depends on V02-018 and V02-021.

## Exit criteria

- New-risk, severity, completeness, scanner failure, and expiry behavior are evidenced.
- Policy exit behavior is stable.
- SARIF validates and preserves provenance.
