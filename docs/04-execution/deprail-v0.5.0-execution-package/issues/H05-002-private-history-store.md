# H05-002 Private History Store
- GitHub Issue: [#421](https://github.com/geoffrey-xiao/deprail/issues/421).
- GitHub parent: [#418](https://github.com/geoffrey-xiao/deprail/issues/418). Published dependencies: #420; global pre-start gate #419.

## Planning metadata

- Type: `feature`
- Area: `foundation`
- Priority: `P0`
- Risk: `R3`
- Target version: `v0.5.0`
- Milestone: [`v0.5.0`](https://github.com/geoffrey-xiao/deprail/milestone/11)
- Sprint: Sprint 4
- Owner: `@geoffrey-xiao`
- Reviewer: `@geoffrey-xiao`
- Dependencies: [H05-001](H05-001-safe-history-projection.md); accepted planning/design contracts and parent-owned timeout follow-up disposition.
- Blocked reason: None unless Project Status is Blocked; runtime remains gated on acceptance below.

## Definition of Ready

- [x] Value and user impact are stated.
- [x] Scope and explicit exclusions are stated.
- [x] Inputs, outputs, and failure behavior are defined.
- [x] Required tests or smoke scenarios are named.
- [x] Acceptance criteria are observable.
- [x] Owner reviewer is assigned; external reviewer is optional.
- [x] Dependencies and target version are recorded.

## Goal

Persist explicitly selected history occurrences privately and transactionally, while absent-store reads stay read-only and all refusal paths preserve prior committed data.

## Scope

Own `internal/store/history` with `OpenReadOnly(ctx)`, `OpenWriter(ctx)`, `Append(ctx, Entry)`, `Get(ctx,id)`, and `List(ctx,Page)` boundaries. Store only a validated H05-001 projection plus exact digest reference rows. Use `modernc.org/sqlite` pinned candidate `v1.59.0` after dependency/license/SBOM/vulnerability review; CGo-free `database/sql`. Store path is canonical `os.UserConfigDir()/.deprail/history.sqlite3`, never repository fallback; refuse if inside scanned repository or unsafe symlink/reparse boundary. POSIX directory 0700/database and backup 0600; Windows owner-only ACL; failure to establish permissions refuses without alternate write location.

SQLite `user_version=1`, projection `history-v1` and source version are distinct. Verify WAL, foreign keys on every connection, FULL synchronous on writers, one serialized writer/BEGIN IMMEDIATE and 2,000ms busy timeout with cancellation precedence; no retry loop. The application supplies one occurrence UUID/time; storage never regenerates identity. An explicit identical-ID retry is idempotent only for identical validated bytes/refs; conflicting content fails. Different occurrences may share source scan ID. Commit row plus all refs atomically and validate every denormalized value against projection on write/read.

Before future migrations, SQLite online backup to unique sibling, restrictive permissions, validate before finalizing; never overwrite backup or copy only main DB during WAL. Apply ordered forward-only migration and `user_version` in one transaction. Preserve old DB and backup on failure. No downgrade/auto-repair/restore. Current schema v1 opens without migration; future unsupported schema, corrupt schema or failed validation is read-only refusal, never modified. Missing DB reads as empty without creating directories/files; only explicit writer initialization may create it.

Enforce per entry <=16 MiB serialized projection, <=10,000 findings, <=1,000 workspaces, <=128 diagnostics, <=4,096 digests; total <=1,000 entries and <=256 MiB summed projection JSON evaluated inside serialized transaction. Refuse new write at quota; do not evict. Logical quotas do not bound SQLite/WAL/index/backup/external artifact physical growth. Raw artifact files remain in existing store; never delete them.

Trace to [ADR-0004](../../../adr/ADR-0004-local-scan-history.md#detailed-engineering-selections-for-owner-review), [failure/data contract](../requirements/FAILURE-AND-DATA-CONTRACT.md), [FR-508](../requirements/FUNCTIONAL-REQUIREMENTS.md), [STORE-01](../requirements/FAILURE-AND-DATA-CONTRACT.md), [SEC-02/08](../requirements/SECURITY-REQUIREMENTS.md), [release plan §14](../../../03-planning/deprail-development-plan-v0.5.0.md).

## Out of Scope

No UI/HTTP, automatic capture, retention/deletion/eviction, artifact garbage collection, downgrade, automatic repair or restore, source report rewrite, raw artifacts in DB, repository-local fallback, or claim of hard disk quota.

## Inputs, Outputs, and Failure Behavior

Input: schema-valid projection, same-operation occurrence ID, digest references and context for canonical repository boundary. Output: committed entry/ref set, read-only absence as empty, or explicit typed failure. Busy deadline (2s), cancellation, disk full, invalid projection, quota, permission failure, corruption, unsafe path or unsupported future schema never report saved and leave prior committed rows intact. `HISTORY_WRITE_FAILED` covers write/projection/quota/lock/durable failure; compatibility/corruption/permission states remain explicit typed failures per error model. No automatic retries or destructive recovery.

## Required Tests

- Fresh/current-v1, missing read-only (prove no filesystem creation), first explicit initialization and unsupported future/corrupt schema nonmutation.
- Atomic insert+refs under process interruption, cancellation, lock held beyond 2s, disk-full injection and quota exact/one-over; compare prior DB and refs.
- Idempotent same-ID identical retry; conflicting retry refusal; distinct entries with repeated source ID.
- POSIX mode and Windows ACL evidence; unsafe in-repo root/symlink/reparse refusal; WAL/foreign-key/FULL verification.
- Online backup during WAL writes, backup validation, migration failure preservation, no overwrite; manual restore rehearsal future gate.
- Missing and digest-mismatch artifact refs remain integrity states; no deletion.
- Smoke the real store API with a validated occurrence, close/reopen and list/read projection plus refs; absent-store reads create nothing. Full scan/CLI opt-in smoke belongs to H05-003/004, preventing a dependency on those later slices.

## Acceptance Criteria

- [ ] Canonical per-user store is owner-private and unsafe roots/boundaries fail closed on declared OSes.
- [ ] Absent-store reads are empty and do not create anything; explicit save alone initializes.
- [ ] Current v1 persists unique occurrences and validated refs atomically; same-ID retry semantics are deterministic.
- [ ] Busy timeout is 2s; all cancellation/lock/disk-full/quota failures leave previous committed data intact and never claim save.
- [ ] Entry/store bounds refuse rather than truncate/evict; disclose logical versus physical quota.
- [ ] Future schema/corruption is never auto-repaired/downgraded; migration backup behavior is validated before future migrations.
- [ ] Runtime permission, transaction, failure and real-store smoke evidence are attached after implementation, not asserted now.

## Owner Review

`@geoffrey-xiao` reviews technical behavior and separately records security/architecture acceptance, including physical-growth risk and platform permission evidence. External review optional. Runtime/platform evidence and CI are future acceptance.

## Evidence Required

- Verification commands or scenarios: store unit/integration cases, SQLite fault injection, actual fresh/current DB smoke; Linux/macOS/Windows permission and WAL/backup scenarios before release.
- Expected artifacts, logs, screenshots, or links: DB schema/version evidence, modes/ACL, before/after DB integrity, backup validation and failure-injection results; owner decisions.

## Implementation evidence (acceptance pending)

The living implementation and verification record is [`H05-002-STORE-EVIDENCE.md`](../tracking/H05-002-STORE-EVIDENCE.md), with the reviewed dependency SBOM at [`H05-002-dependencies.cdx.json`](../tracking/H05-002-dependencies.cdx.json). [PR #436](https://github.com/geoffrey-xiao/deprail/pull/436) is ready for owner review and Project Status is `Review`. The latest exact-head CI run [36805833391](https://github.com/geoffrey-xiao/deprail/actions/runs/36805833391) passed Ubuntu/macOS/Windows after the first four runs exposed Windows path ownership and DACL-shape assumptions. Owner technical/security review remains pending; acceptance boxes remain unchecked.

## Final Acceptance

- [ ] Owner reviewed every acceptance criterion during PR review.
- [ ] Required verification and CI results were reviewed.
- [ ] Owner review and remaining risk are recorded.
- [ ] Evidence links are attached.
- [ ] Owner review and merge evidence are linked.

## Rollback

Disable history-store integration; preserve DB/WAL/SHM, validated backups and artifacts. Never reset, downgrade, restore automatically or delete data. Owner `@geoffrey-xiao` owns recovery.
