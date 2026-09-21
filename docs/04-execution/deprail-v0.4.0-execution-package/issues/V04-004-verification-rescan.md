# V04-004: Verify and Classify Before/After Scans

**Epic:** EPIC-003
**Status:** Proposed; implementation blocked

## Scope

Discover approved test/build/type-check commands, execute them with process limits, run the existing scan contract before and after mutation, and classify resolved, residual, introduced, and unknown findings.

## Acceptance

Verification success, failure, incomplete, unavailable, and cancellation remain distinct. Scan identity and ordering are stable. A zero-finding result is safe only when complete.
