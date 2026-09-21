# V04-002: Bind Approval to Plan and Source Identity

**Epic:** EPIC-002
**Status:** Proposed; implementation blocked

## Scope

Define the versioned approval record, plan digest, source commit/root identity, authorized paths, expiry, single-use semantics, dry-run behavior, and stable approval errors.

## Acceptance

Changed plan, source, root, target, paths, expired approval, and reused approval are rejected before mutation. Approval records are redaction-safe and auditable.
