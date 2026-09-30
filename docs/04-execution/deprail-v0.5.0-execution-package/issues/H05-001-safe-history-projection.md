# H05-001 Safe History Projection

## Planning metadata

- Type: `feature`
- Area: `normalization`
- Priority: `P0`
- Risk: `R3`
- Target version: `v0.5.0`
- Milestone: [`v0.5.0`](https://github.com/geoffrey-xiao/deprail/milestone/11)
- Sprint: Unassigned
- Owner: `@geoffrey-xiao`
- Reviewer: `@geoffrey-xiao`
- Dependencies: Accepted PR #417 design baseline and [V05-QA-001](V05-QA-001-verification-timeout-disposition.md) resolution or explicit owner pre-start disposition.
- Blocked reason: None unless Project Status is Blocked; runtime remains gated on acceptance below.
- Epic: [V05-EPIC-002](../epics/EPIC-002-local-history-delivery.md).
- Contract mapping: release plan §14 safe storage/diagnostics, [product baseline](../../../01-product/deprail-product-design-v1-ai.md), [architecture baseline](../../../02-architecture/deprail-architecture-and-tech-stack-v1.md), FR-501/503/504 and SEC-07.

## Definition of Ready

- [x] Value and user impact are stated.
- [x] Scope and explicit exclusions are stated.
- [x] Inputs, outputs, and failure behavior are defined.
- [x] Required tests or smoke scenarios are named.
- [x] Acceptance criteria are observable.
- [x] Owner reviewer is assigned; external reviewer is optional.
- [x] Dependencies and target version are recorded.

## Goal

Produce a deterministic, privacy-safe, separately versioned `history-v1` projection from the actual same-operation scan inputs, without changing `ScanReport`, CLI JSON, or remediation inputs.

## Scope

Own plain projection entities in `internal/domain/history` and `ProjectHistory(input HistoryInput) (history.Projection, error)` / `ValidateProjection(history.Projection) error` in `internal/normalize/history.go`. `normalize.HistoryInput` owns neutral typed source fields; H05-003 maps its actual graph/report/diagnostic context to those fields without an inward import of `internal/app` or scanner-specific types into domain. The input carries canonical operation UUID/time/outcome, optional graph/report values, typed diagnostics and approved provenance. [`history.schema.json`](../../../../schemas/history-v1/history.schema.json) is the exact allowlist; examples are fixtures, not actual runtime evidence.

Map fields explicitly: repository label only (nullable); workspaces nullable when same-operation graph unavailable, otherwise actual safe repository-relative summaries; report null when unavailable, otherwise schema version `v1alpha`, nonempty source scan ID, completeness, repository-state digest, approved nullable provenance, normalized findings and unique validated artifact digests. Null means unavailable; empty means known empty. Operation outcome (`completed|failed|cancelled`) is independent from report status (`complete|partial|failed`). Never infer successful empty results. Preserve source ordering where required; sort/deduplicate diagnostics by scope, workspace ID, code and exact record; stable finding identity uses existing `normalize.StableFindingKey` with workspace, component/PURL, version, vulnerability ID and sorted aliases. `TargetID` is vulnerability ID, not a finding key.

Map recognized typed diagnostic codes and fixed schema messages; require validated scope/workspace association. Never parse `ScanReport.Errors`. Unknown required sources, ambiguous identity, unsafe paths/digests or malformed data reject capture. Missing provenance values are null, not a required known tool version; present unsafe values refuse capture. Exclude absolute roots, raw errors/secrets/credential URLs/environment/source contents and unreviewed unknown fields. Reject absolute/drive/UNC/traversal/control/backslash paths and dot components except exact `.`.

Admission before encoding: complete UTF-8 `projection_json` <=16 MiB; <=10,000 findings, 1,000 workspaces, 128 diagnostics, 4,096 artifact digests. No truncation. Encoded output must validate against schema and round-trip without semantic change. Concrete source mapping and typed diagnostics are defined by [failure/data contract](../requirements/FAILURE-AND-DATA-CONTRACT.md), [error model](../requirements/ERROR-MODEL.md), and [ADR-0004 selections](../../../adr/ADR-0004-local-scan-history.md#detailed-engineering-selections-for-owner-review).

## Out of Scope

No source-report/CLI format changes, rescanning or rediscovery, storage/API/CLI wiring, raw artifact reads, migration, schema loosening, new diagnostic codes, or claiming values unavailable from trusted same-operation inputs.

## Inputs, Outputs, and Failure Behavior

Input is typed validated operation context, not arbitrary JSON. Output is one schema-valid projection with immutable operation ID/timestamp or a typed projection failure. Missing graph/report remains null only where schema permits; unknown diagnostic, excessive bounds, invalid identity/path/digest or serialization/schema mismatch fails closed. Failure must not alter scan outcome, stdout bytes, or existing database state.

## Required Tests

- Schema-valid complete/nonempty, failed, cancelled-without-context, partial and known-empty examples; null versus empty distinction.
- Real `ScanReport` finding mapping including duplicate source scan IDs, stable key consistency with `StableFindingKey`, aliases ordering, vulnerability ID mapping and repeated-run distinct occurrence UUIDs.
- Reject absolute/UNC/drive/traversal/control/backslash paths, unsafe repository labels, missing/ambiguous identity, invalid artifact digest, unknown diagnostic/code, raw string errors, and invalid provenance; assert no source values leak into serialized projection.
- Exact limits and one-over for bytes/counts; UTF-8 and escaping byte accounting; reject rather than truncate.
- Determinism under reordered equivalent inputs; no mutation of source report/CLI serialization.
- Smoke a real representative scan input through projection and schema validation; inspect serialized bytes for root, raw errors and secret sentinels.

## Acceptance Criteria

- [ ] Projection matches strict `history-v1` allowlist and conveys source/report/operation states without conflation.
- [ ] Repeated scan IDs remain distinct operation occurrences; stable finding keys match existing algorithm.
- [ ] Missing context is null, known-empty context remains empty, and no rediscovery or fabricated provenance occurs.
- [ ] Unsafe/unknown/incomplete input and every exceeded bound fail closed without truncation or source mutation.
- [ ] Privacy test proves no absolute root, raw error, secret, credential URL, environment value, or unknown field persists.
- [ ] Required behavior, failure/security tests and real-input smoke are attached as future evidence; none is claimed passed in planning.
- [ ] Documentation/crosswalk and evidence requirements identify schema and source contracts.

## Owner Review

`@geoffrey-xiao` reviews every criterion and separately records technical contract acceptance. Security/architecture acceptance is a distinct decision/evidence item; external review is optional. Runtime behavior, actual smoke and CI remain future evidence.

## Evidence Required

- Verification commands or scenarios: targeted projection/schema tests, byte/count boundary and privacy leakage cases, representative actual scan-input smoke; future CI command recorded by implementation PR.
- Expected artifacts, logs, screenshots, or links: test results, schema-valid emitted sample, safe serialized-byte inspection, source-to-field mapping and owner technical/security decisions.

## Final Acceptance

- [ ] Owner reviewed every acceptance criterion during PR review.
- [ ] Required verification and CI results were reviewed.
- [ ] Owner review and remaining risk are recorded.
- [ ] Evidence links are attached.
- [ ] Owner review and merge evidence are linked.

## Dependencies and pre-start gate

The accepted design baseline and [V05-QA-001](V05-QA-001-verification-timeout-disposition.md) resolution/disposition gate runtime starts. Dependents are [H05-002](H05-002-private-history-store.md) and [H05-003](H05-003-history-application-services.md). No runtime implementation is authorized by publication alone.

## Rollback

Withdraw the projection/capture support without changing CLI report bytes; retain any already committed history and artifacts. Owner `@geoffrey-xiao` records rollback and evidence.