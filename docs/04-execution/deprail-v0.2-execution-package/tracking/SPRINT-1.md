# v0.2 Sprint 1 — Product Reliability

## Goal

Complete the highest-value reliability and developer-experience improvements after the contract baseline is stable.

## Planned items

- V02-003 truthful build version identity.
- V02-008 scan-from-outside-root regression.
- V02-010 deterministic artifacts and ordering.
- V02-011 release identity verification across artifacts.
- V02-012 release-mode smoke procedure.

## Dependencies

Sprint 0 must establish the output, scanner, and release contracts. V02-011 depends on the version identity decision in V02-003; V02-008 depends on V02-007.

## Exit criteria

- Release and development builds report truthful identity.
- Outside-root execution has a permanent regression test.
- Deterministic output and artifact behavior are evidenced.
- Preview, RC, and stable smoke procedures are executable and reviewed.
