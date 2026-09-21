# EPIC-003: Verification and Rescan

**Status:** Proposed

## Outcome

A remediation result is verified with supported checks and compared with a deterministic before/after scan.

## Boundaries

Includes test/build/type-check discovery, command ranking, bounded execution, pre/post scan identity, and resolved/residual/introduced/unknown classification. Excludes policy gates and unsupported ecosystem behavior.

## Acceptance

Successful, failed, incomplete, and unavailable verification are distinct; scan transitions are stable and zero findings are safe only for complete scans.
