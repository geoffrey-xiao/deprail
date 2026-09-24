# V05-003: Decide SQLite History Schema and Lifecycle
- GitHub Issue: [#391](https://github.com/geoffrey-xiao/deprail/issues/391)

- Epic: [EPIC-001 / #390](https://github.com/geoffrey-xiao/deprail/issues/390)
- Target: `v0.5.0`
- Status: Closed in GitHub; original proposal was owner-approved in merged PR #407; the owner-selected history-projection revision remains proposed after PR #411 and still needs owner technical decisions. Independent review is optional under ADR-0005; v0.5 technical DoR and runtime authorization remain open.
- Type: decision
- Area: docs
- Priority: P0
- Risk: R3
- Owner: `@geoffrey-xiao`
- Reviewer: `@geoffrey-xiao` (owner); external review is optional under [ADR-0005](../../../adr/ADR-0005-solo-owner-review-policy.md).
- Dependencies: V05-000 / #387; shared history vocabulary from V05-001

## Value

A durable local history store must preserve existing report meaning and provenance and must not lose the only prior database copy during a failed migration or recovery.

## Scope

Create the SQLite data-lifecycle ADR and versioned schema contract. Decide canonical data root and restrictive permissions; DB/schema version and ID generation; stored fields versus artifact references; ingestion trigger and retry/idempotency; transaction scope and writer/readers/locking/journal/busy behavior; size bounds; retention/deletion/export and artifact reference accounting; migration/backup/restore/downgrade policy; disk-full, corruption, interrupted migration/write, permission, and unsupported-future-schema recovery; and behavior when persistence fails while a scan succeeds. State data classification, privacy, and platform implications.

## Out of scope

No SQL schema migration code, SQLite driver selection/install, store adapter, destructive repair, automatic eviction, backup command, artifact deletion, remote sync, or runtime implementation issue.

## Inputs, outputs, and failure behavior

Inputs: `ARCHITECTURE-v0.5.md`, `FAILURE-AND-DATA-CONTRACT.md`, compatibility/security requirements, existing report/artifact contracts, and the accepted history workflow. Outputs: the proposed [ADR-0004](../../../adr/ADR-0004-local-scan-history.md) and versioned schema/lifecycle specification. These are owner-review drafts; numeric bounds and owner acceptance of unbounded-growth risk remain DoR gates. Storage open/lock/full/corrupt, migration interruption, artifact loss/digest mismatch, and unsupported schema must remain explicit failures; failed history persistence must not rewrite scan outcome or report completeness.

## Required verification and evidence

Review schema/identity round trips for repeated `sourceScanID`, migration/rollback/recovery sequences, restrictive permissions and artifact-reference integrity. Provide a deterministic migration matrix covering fresh, supported prior, latest, interrupted, corrupt, and future-version stores. No runtime DB implementation or platform pass is asserted.

## Acceptance criteria

- The schema has unique `historyEntryID` distinct from repeatable `sourceScanID`; operation outcome and report completeness remain separate.
- The ADR defines transaction/concurrency, migration/backup/recovery, retention/deletion/export, permissions, size limits, and artifact-reference behavior.
- Failure outcomes preserve old committed data and never present failed writes as saved or empty-success.
- CLI compatibility, data-root privacy, cross-platform constraints, and unresolved risk are explicit.
- Owner technical acceptance, security evidence, and runtime authorization remain outstanding until the v0.5 Definition of Ready passes; independent review is optional.
