# V02-024 Action Permissions and Caching

- Type: docs
- Area: docs
- Priority: P0
- Risk: R2
- Target version: 0.2.0
- Sprint: S7
- Dependencies: V02-023
- Parent epic: EPIC-008

## Goal

Document and verify least-privilege permissions, cache boundaries, and data handling for the GitHub Action.

## Scope

Permissions, token use, cache keys and invalidation, artifact retention, fork behavior, and source/privacy boundaries.

## Out of Scope

Hosted storage, telemetry, broad organization administration, and automatic write-back.

## Inputs, Outputs, and Failure Behavior

The Action works with documented minimum permissions. Cache misses remain safe; stale or untrusted cache data cannot create a safe result. Fork and permission failures are explicit.

## Required Tests

Permission matrix, fork pull request behavior, cache hit/miss, invalidation, artifact handling, and secret-redaction cases.

## Acceptance Criteria

- Minimum permissions are stated and tested.
- Cache identity includes relevant repository, configuration, tool, and database inputs.
- No source or secret leakage occurs through artifacts or logs.
- Unsupported permission contexts fail clearly.

## Human Review

Security review is required.

## Evidence Required

Permission matrix, cache design, fork smoke output, and redaction evidence.
