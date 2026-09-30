# V05-EPIC-002 Deliver Local Scan History and the Read-Only Console
- GitHub Issue: [#418](https://github.com/geoffrey-xiao/deprail/issues/418).
- Published children: quality investigation #419; bounded-capture correction #430; darwin fixture correction #432; H05-001–H05-008 #420–#427. See [durable backlog](../tracking/IMPLEMENTATION-BACKLOG.md).

## Planning metadata

- Type: `feature`
- Area: `foundation`
- Priority: `P0`
- Risk: `R3`
- Target version: `0.5.0`; preview first, no release authorized here.
- Milestone: `v0.5.0`
- Sprint: Sprint 4 in the live Project following the owner's tracking assignment; this does not authorize H05 runtime work.
- Owner: `@geoffrey-xiao`
- Reviewer: `@geoffrey-xiao`; independent review optional under ADR-0005.
- Dependencies: [EPIC-001 / #390 design baseline](https://github.com/geoffrey-xiao/deprail/issues/390), [owner exact-contract acceptance](https://github.com/geoffrey-xiao/deprail/pull/417#issuecomment-5907670878), merged [PR #417](https://github.com/geoffrey-xiao/deprail/pull/417).
- Project status at publication: Todo; publication does not start implementation.
- Blocked reason: None (Todo, not started). Runtime starts require quality preflight disposition, issue-specific DoR/owner review and synchronized reviewed main; no runtime code authorized in this handoff.

## Definition of Ready

- [x] Value and user impact are stated.
- [x] Scope and explicit exclusions are stated.
- [x] Inputs, outputs, and failure behavior are defined.
- [x] Required tests or smoke scenarios are named.
- [x] Acceptance criteria are observable.
- [x] Owner reviewer is assigned; external review optional.
- [x] Dependencies and target version are recorded.

This is specification readiness for publishing the approved backlog, not completion of children or permission to skip their pre-start gates.

## Goal

A developer explicitly saves trustworthy local scan history and browses it through an embedded, read-only local console, with operation outcome, report completeness, provenance and missing evidence visibly distinct. The optional console/store never become prerequisites for ordinary scans.

## Scope

The accepted source baseline is reviewed head `b0649547d1e0f033f10f3f28920b74eab9b16e1b`, merged into main as `40e43a83eaac1b84dd038ee1ea4dccbd3e7234f6`. Implementation capabilities map one-to-one to the eight approved H05 slices; the ninth child is a readiness preflight, not an added product capability.

| Order | Child contract | Primary outcome | Dependencies |
| --- | --- | --- | --- |
| 0 | [V05-QA-001](../issues/V05-QA-001-verification-timeout-disposition.md) | Resolve/disposition the known local verification timeout before runtime starts | Accepted owner decision; original failure evidence |
| 1 | [H05-001](../issues/H05-001-safe-history-projection.md) | Safe deterministic history projection | Accepted design; QA pre-start disposition |
| 2 | [H05-002](../issues/H05-002-private-history-store.md) | Owner-private durable SQLite history | H05-001 |
| 3 | [H05-003](../issues/H05-003-history-application-services.md) | Same-operation capture and shared read queries | H05-001–002 |
| 4 | [H05-004](../issues/H05-004-opt-in-history-cli.md) | Explicit scan history CLI compatibility | H05-003 |
| 5 | [H05-005](../issues/H05-005-read-only-local-console.md) | Bounded read-only HTTP/static transport and foreground launch | H05-003 |
| 6 | [H05-006](../issues/H05-006-accessible-history-ui.md) | Accessible history/detail/help UI source | H05-005 contract |
| 7 | [H05-007](../issues/H05-007-embedded-console-packaging.md) | Reproducible embedded binary delivery | H05-005–006 |
| 8 | [H05-008](../issues/H05-008-local-history-release-evidence.md) | Integrated real workflow/release evidence | H05-001–007 |

The source-level UI smoke uses a throwaway embedding harness; production embedding/configuration belongs to H05-007. This prevents a circular dependency while requiring final packaged behavior. Temporary verification assets are never shipped placeholders. All children share the same pre-start gate and are assigned Todo, not In Progress.

## Contract traceability

- Product: [product design](../../../01-product/deprail-product-design-v1-ai.md), local-first web history/presentation, not hosted collaboration or a second scanner.
- Architecture: [technology baseline](../../../02-architecture/deprail-architecture-and-tech-stack-v1.md), Go/application/domain boundaries, SQLite, REST/OpenAPI and embedded React/TypeScript/Vite.
- Roadmap: [v0.5 local web/history row](../../../03-planning/deprail-roadmap-v1.md#release-roadmap).
- Release: [plan §14](../../../03-planning/deprail-development-plan-v0.5.0.md#14-detailed-contract-reconciliation-for-current-owner-review), including explicit #356 owner/target/exclusion dispositions.
- Execution: [FR-501–511](../requirements/FUNCTIONAL-REQUIREMENTS.md), failure/data, error, security, compatibility and test strategy, [ADR-0004](../../../adr/ADR-0004-local-scan-history.md), [history-v1 schema](../../../../schemas/history-v1/history.schema.json), [OpenAPI](../../../../schemas/openapi/v1/openapi.yaml).

## Out of Scope

No hosted/team service, PostgreSQL, accounts/RBAC, source upload, publishing, browser-triggered scan, remediation or source mutation, agent writes, new scanner, automatic installation, eviction/delete/export/repair/downgrade. No release tagging/merge permission is implied. Full Python/Java remediation remains deferred per release plan §14; representative multi-language history workflows are included.

## Inputs, Outputs, and Failure Behavior

Capture consumes typed same-operation graph/report/error/metadata, never rediscovery or parsed error strings. Store only allowlisted history-v1; the existing ScanReport/CLI JSON is unchanged. Missing report/context differs from trustworthy empty report/context. Explicit capture failures preserve scan output and original nonzero exit; additive exit6 applies only to otherwise-successful scan/save failure. Private storage atomically refuses quota/lock/disk-full/unsafe projections without altering old rows. Queries and UI cannot capture/mutate, hide corrupt/future storage as empty, or invent verified artifact provenance.

## Required Tests

Each child supplies deterministic offline consumer-visible unit/contract/security tests and actual targeted smoke as applicable. H05-008 requires reviewed-binary JS/Python/Java workflows, Linux/macOS/Windows packaging/permissions/recovery/API and browser/assistive-technology evidence. The general release checklist remains a separate release gate; no compile or mock substitutes for runtime/manual proof.

## Acceptance Criteria

- [ ] All eight product slices satisfy their exact observable acceptance/failure criteria and have owner-reviewed PR/CI/evidence links; no shipped scaffold or placeholder.
- [ ] The known local verification timeout has a linked resolution or explicit owner readiness disposition before runtime work starts.
- [ ] Ordinary scans remain independent of history/UI and retain machine stdout/report meaning; selected saves truthfully retain complete/partial/failed/cancelled state axes.
- [ ] The packaged console exposes only accepted read operations/local assets and satisfies actual hostile-input, auth, accessibility and platform/recovery gates.
- [ ] #356 predecessor gaps are executed or explicitly dispositioned for the release mode, never silently counted complete; deferred remediation is not delivered under this epic.
- [ ] Release identity, artifacts, SBOM/signature/provenance statuses and separate owner security/release decisions are recorded against the immutable reviewed binary before publication.

## Owner Review

The owner approved the exact design and residual-risk proposal in the linked conversation decision, separately recording technical and security/architecture acceptance. Each implementation PR still requires owner acceptance of every issue criterion and remaining risk. One person may provide both review records. No two R3 items run concurrently; at most two implementation issues/two review PRs active, one primary outcome per PR.

## Evidence Required

Per-child commands/exits/platform/commit, sanitized actual artifact or screenshot links, complete contract/security/manual evidence, owner review/merge, dependency-license/SBOM/provenance and candidate vulnerability review before introduction, final [release checklist](../../../RELEASE-CHECKLIST.md) dispositions and version-specific release evidence. Planning evidence remains separate from runtime proof.

## Rollback

Stop/withdraw the optional console/capture paths without altering default scanning. Preserve DB/WAL/SHM/backups and external artifacts; never auto-reset/delete/downgrade. Use a matching reviewed immutable binary/assets; a corrected release needs a new immutable tag.

## Final Acceptance

- [ ] Owner reviewed each child and the integrated acceptance criteria.
- [ ] Required CI and real workflow/browser/platform/recovery evidence is linked.
- [ ] Owner security/architecture and release decisions are separately recorded.
- [ ] All child/PR/evidence links and remaining risks are reconciled in tracking.
- [ ] Approved release mode has a complete checklist or explicit allowed preview dispositions; this epic's creation does not authorize publication.

## Bounded process capture prerequisite

[V05-QA-002 / #430](../issues/V05-QA-002-bounded-output-capture.md) restores the existing process capture limit after the #419 investigation found a real output-boundary bypass. This is a security correction, not a new history capability. Closed #419 and its Project Done status record accepted investigation only: the original copied-helper timeout and explicit owner runtime-start decision remain separately pending. Require reviewed capture correction and the existing readiness gates before H05 kickoff; do not infer readiness from issue closure or Sprint assignment.

## Darwin fixture readiness follow-up

Owner-merged capture fix PR #431 is included in the baseline for [V05-QA-003 / #432](../issues/V05-QA-003-darwin-helper-fixture.md). This follow-up changes only the darwin test helper to reference the already active test image; all production behavior, Linux/Windows copying, deadlines and assertions remain. Corrected-baseline targeted/full checks are recorded in its evidence; owner review/merge and the explicit runtime-start decision still precede H05 kickoff. Closed #419/#430 and Sprint assignment do not alone authorize feature code.

## H05-001 runtime-start decision

The earlier pending QA language records the publication and correction-review stages. [#432 owner acceptance](https://github.com/geoffrey-xiao/deprail/issues/432#issuecomment-5913275072), owner-merged [PR #433](https://github.com/geoffrey-xiao/deprail/pull/433), successful local verification and exact-head three-platform CI resolve the bounded QA correction gate; original failed runs and macOS root-cause uncertainty remain preserved.

Owner “可以 继续吧” separately authorizes only [H05-001 / #420](../issues/H05-001-safe-history-projection.md), with its checked DoR, Sprint 4, Project In Progress and synchronized reviewed-main `6c74cf9`. [Kickoff decision](https://github.com/geoffrey-xiao/deprail/issues/420#issuecomment-5913450360) records the pure projection/validation, privacy/boundary tests and real scan-input smoke scope. No storage, API, UI, dependency introduction, other H05 implementation, merge or release is authorized. Epic acceptance stays unchecked until its complete integrated evidence exists.
