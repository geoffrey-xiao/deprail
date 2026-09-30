# H05-001 Safe History Projection
- GitHub Issue: [#420](https://github.com/geoffrey-xiao/deprail/issues/420).
- GitHub parent: [#418](https://github.com/geoffrey-xiao/deprail/issues/418). Completed global QA prerequisites: #419 investigation, #430 capture correction and #432 fixture correction.

## Planning metadata

- Type: `feature`
- Area: `normalization`
- Priority: `P0`
- Risk: `R3`
- Target version: `v0.5.0`
- Milestone: [`v0.5.0`](https://github.com/geoffrey-xiao/deprail/milestone/11)
- Sprint: Sprint 4.
- Owner: `@geoffrey-xiao`
- Reviewer: `@geoffrey-xiao`
- Dependencies: Accepted PR #417 design baseline; accepted #419 investigation, owner-merged capture fix #431 / #430 and owner-accepted fixture fix #433 / #432; reviewed main `6c74cf9861981928c29f7a8cb09799e63f3fd6dd`.
- Blocked reason: None for this bounded issue after the linked owner kickoff decision and issue-specific DoR; other H05 issues remain Todo.
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

The accepted design baseline and [V05-QA-001](V05-QA-001-verification-timeout-disposition.md) resolution/disposition gate runtime starts. [#432 completion](V05-QA-003-darwin-helper-fixture.md#owner-acceptance-and-completion) supplies the accepted correction and exact local/CI evidence; original failed runs remain historical. Owner “可以 继续吧” separately authorized this issue after the displayed implementation plan. [Kickoff record](https://github.com/geoffrey-xiao/deprail/issues/420#issuecomment-5913450360) links checked DoR, Sprint 4, Project In Progress and branch `feat/420-safe-history-projection` from synchronized reviewed-main `6c74cf9`. Dependents are [H05-002](H05-002-private-history-store.md) and [H05-003](H05-003-history-application-services.md); neither starts through this authorization. Publication or Sprint assignment alone is not runtime permission.

## Implementation evidence (PR #434 merged)

- `ProjectHistory` and `ValidateProjection` implement a pure typed projection; the source `ScanReport` and its JSON remain unchanged. Workspaces/report distinguish unavailable (`null`) from known empty arrays; operation outcome remains independent of report completeness.
- Projection validates IDs, source schema, repository digest, PURLs, repository-relative paths, provenance, diagnostic registry/scope, finding identities, exact allowlisted history fields, sorted/deduplicated aliases/digests/diagnostics, and encoded-byte/count limits. Unsupported future schema markers are distinguished from corrupt stored values. No truncation or parsing of raw report errors.
- Scanner severity enums are canonicalized case-insensitively to the history labels. CVSS vectors are not converted to a guessed label; severity projects as `null`. Unknown values reject projection.
- Targeted projection/schema tests passed. `make verify` passed (`go generate ./...`, `go vet ./...`, `go test ./...`, `go build ./...`).
- Actual local OSV-Scanner 2.6.0 scan of `testdata/fixtures/npm-basic` produced 5 findings and 1 retained artifact. The exact `history-v1` schema validated both complete and missing-scanner failed projections. Round-trip equality, source report immutability, unchanged fixture inputs, unavailable workspace preservation, repeated source scan ID with distinct operation IDs, diagnostic `SCANNER_NOT_FOUND`, and exclusion of root/raw error/credential sentinel all passed. Output SHA-256: `56e0eddd076b34b3f951fce1f440439063ac6f7058788502a1c56116fa260ccc` (temporary artifact `/tmp/deprail-h05-001-d181a2eced014d4c822797fb62aaeb30.json`; harness removed).
- Owner merged PR [#434](https://github.com/geoffrey-xiao/deprail/pull/434) as merge commit `ea12fbc6928560a9a70a8d466748f67a90e0bf55`. Exact-head [CI run 36740438976](https://github.com/geoffrey-xiao/deprail/actions/runs/36740438976) passed on Linux, macOS and Windows; the owner-merged PR was labeled and its Project item is Done. The PR uses `Refs #420`, so the issue remains open. Acceptance checkboxes above remain pending an explicit owner closeout record; no custom Project URL was added.

## Rollback

Withdraw the projection/capture support without changing CLI report bytes; retain any already committed history and artifacts. Owner `@geoffrey-xiao` records rollback and evidence.