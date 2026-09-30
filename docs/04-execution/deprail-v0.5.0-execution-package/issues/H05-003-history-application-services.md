# H05-003 History Application Services

## Planning metadata

- Type: `feature`
- Area: `foundation`
- Priority: `P0`
- Risk: `R3`
- Target version: `v0.5.0`
- Milestone: [`v0.5.0`](https://github.com/geoffrey-xiao/deprail/milestone/11)
- Sprint: Unassigned
- Owner: `@geoffrey-xiao`
- Reviewer: `@geoffrey-xiao`
- Dependencies: [H05-001](H05-001-safe-history-projection.md), [H05-002](H05-002-private-history-store.md); accepted planning/design contracts and parent-owned timeout follow-up disposition.
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

Provide one shared application boundary for same-operation capture and safe, bounded read queries, without rescanning or changing existing scan behavior.

## Scope

Own `internal/app/history.go` and an additive optional `ScanOptions.Capture` sink. `CaptureSink.Capture(ctx, CaptureInput) (SaveResult,error)` is a single terminal callback, not a new begin/complete/finish lifecycle. CaptureInput is app-owned same-operation graph/report/error/typed diagnostics/provenance; it maps to neutral normalize.HistoryInput. Preserve existing Scan signature, serialized report and all default callers. Allocate one UUID/time for the selected occurrence; invoke the sink once at every terminal boundary with trustworthy available context, before typed errors are lost to strings. Never rescan/rediscover. Selected terminal persistence has a fresh 3s context, 2s store busy cap and second-termination cancellation; save errors remain separate from original scan error/report.

The application query boundary has `List`, `Detail`, `Workspaces`, `Findings` and `ArtifactIntegrity` operations with explicit parent IDs/continuations/limits where applicable. ArtifactIntegrity returns only verified/missing/digest_mismatch/unavailable metadata, never artifact bytes. Domain types import no scanner/SQL/Cobra/HTTP types; HTTP owns envelope serialization and cursor wire encoding while shared app/domain services own ordering, continuation and no-truncation semantics.

Queries return validated safe domain views and ordered, count-bounded candidates with continuation values: history uses `(recorded_at_us DESC, history_entry_id DESC)`, workspaces/finding children use immutable contract ordering, default25/max50. HTTP transport serializes the final envelope/cursor under 1MiB, returning whole records and advancing from the last emitted—not merely last loaded—record. No truncation, skipped rows or empty nonterminal page; detail is not split. Trusted digest resolution streams verification and returns explicit integrity metadata; absent resolver is unavailable. Missing storage is empty without initialization.

Source contracts: [FR-501–506](../requirements/FUNCTIONAL-REQUIREMENTS.md), [STORE-01](../requirements/FAILURE-AND-DATA-CONTRACT.md), [API/OpenAPI](../API-DESIGN.md), [artifact integrity](../requirements/FAILURE-AND-DATA-CONTRACT.md), [ADR-0004](../../../adr/ADR-0004-local-scan-history.md#detailed-engineering-selections-for-owner-review), [release plan §14](../../../03-planning/deprail-development-plan-v0.5.0.md), product/architecture baseline [product design](../../../01-product/deprail-product-design-v1-ai.md) and [architecture](../../../02-architecture/deprail-architecture-and-tech-stack-v1.md).

## Out of Scope

No HTTP listener/browser, automatic capture, CLI flag, scan API signature change, scan/rescan, raw report/error parsing, unsafe path resolver, entry deletion/export, or persistence-success claim after failed write.

## Inputs, Outputs, and Failure Behavior

Capture consumes the actual same-operation context and yields saved occurrence or typed `HISTORY_WRITE_FAILED`; save failure is separate from scan outcome, never mutates report/stdout, and does not retry. If both scan and persistence fail, preserve original scan failure and expose safe separate persistence diagnostic. Queries return validated DTO/page or typed missing/incompatible/corrupt/too-large/artifact-integrity errors. Missing store returns empty list. Every DTO's report completeness, operation outcome, and artifact integrity remain independent.

## Required Tests

- Real operation integration captures truthful context once and compares returned scan report/output; no wiring/mock echo/source-signature assertion. Actual scan/scanner invocation evidence confirms capture does not initiate a second scan or discovery.
- Completed, partial, failed and cancelled capture; missing graph/report null; unsafe/unknown diagnostic or >128 refuses capture without false empty success.
- Same-operation retry stable ID/no duplicate; new operation same source ID distinct; simultaneous operations distinct.
- Terminal persistence deadline, cancellation and store lock timeout; primary scan nonzero outcome remains primary, no false saved indication.
- Read absent store without creating; deterministic pagination including equal timestamps, Unicode byte-boundary, one record too large; no skip/duplicate across pages.
- Trusted artifact present/missing/digest mismatch, no resolver (`unavailable`), unreferenced digest and hostile path; return integrity metadata only, never bytes or absolute paths.
- Smoke actual scan operation capture, close/reopen, query pages and artifact integrity states; ensure no rescan via operation-level counter/trace evidence.

## Acceptance Criteria

- [ ] Optional capture consumes same-operation typed graph/report/error provenance exactly once; no rescan or rediscovery.
- [ ] Existing `Scan` signature, report bytes and default behavior are unchanged; absence of opt-in leaves history untouched.
- [ ] Terminal save stays within 3s budget; persistence failure never claims saved or overwrites primary scan result.
- [ ] Query DTOs are validated safe projections; absent store is empty and read-only.
- [ ] Paging is stable and whole-record byte bounded, with explicit too-large failure; no omission/truncation.
- [ ] Artifact integrity queries verify digests and distinguish verified/missing/digest_mismatch/unavailable without returning bytes, unsafe resolution or deletion.
- [ ] Runtime consumer-edge, fault, and actual-operation smoke evidence is linked after implementation.

## Owner Review

`@geoffrey-xiao` reviews application contract and distinct security/architecture acceptance; external review optional. Runtime proof, CI and actual API-consumer evidence are future gates.

## Evidence Required

- Verification commands or scenarios: application boundary and store integration tests; real scan-operation capture/query/artifact smoke; lock/cancel/deadline and byte-page tests.
- Expected artifacts, logs, screenshots, or links: operation-level trace proving one capture/no rescan, before/after report byte comparison, persistence outcome evidence, pages and digest-verification results; owner decisions.

## Final Acceptance

- [ ] Owner reviewed every acceptance criterion during PR review.
- [ ] Required verification and CI results were reviewed.
- [ ] Owner review and remaining risk are recorded.
- [ ] Evidence links are attached.
- [ ] Owner review and merge evidence are linked.

## Rollback

Remove optional capture/query integration while retaining the original scan contract and committed history/artifacts. Owner `@geoffrey-xiao` records any recovery action.