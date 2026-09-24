# DepRail v0.5.0 Execution Package

| Attribute | Value |
| --- | --- |
| Release | v0.5.0; preview first |
| Status | Draft; external reviews are optional, and technical DoR/owner decisions remain open |
| Preparation issue | [#387](https://github.com/geoffrey-xiao/deprail/issues/387) |
| Release plan | [`../../03-planning/deprail-development-plan-v0.5.0.md`](../../03-planning/deprail-development-plan-v0.5.0.md) |
| Product baseline | [`../../01-product/deprail-product-design-v1-ai.md`](../../01-product/deprail-product-design-v1-ai.md) |
| Architecture baseline | [`../../02-architecture/deprail-architecture-and-tech-stack-v1.md`](../../02-architecture/deprail-architecture-and-tech-stack-v1.md) |
| Roadmap | [`../../03-planning/deprail-roadmap-v1.md`](../../03-planning/deprail-roadmap-v1.md) |
| Predecessor closeout | [#356 v0.4.0-preview.2](https://github.com/geoffrey-xiao/deprail/issues/356), still open |
| Owner | `@geoffrey-xiao` |
| Review authority | `@geoffrey-xiao` (sole required owner reviewer under [ADR-0005](../../adr/ADR-0005-solo-owner-review-policy.md)); external review is optional. |

This package translates the proposed v0.5 release outcome into reviewable UX, API, architecture, data, failure, security, compatibility, and test contracts. It is a draft subordinate to the product design, architecture, roadmap, and v0.5 release plan. Conflicts stop implementation and require a versioned decision; lower-level requirements cannot silently override those sources.

## Source order

1. Product design and architecture.
2. Whole-project roadmap.
3. [`v0.5.0 development plan`](../../03-planning/deprail-development-plan-v0.5.0.md).
4. This package and its requirements.
5. Owner-reviewed issue contracts.
6. Implementation, verification, and release evidence.

## Package documents

- [`PRD-v0.5.md`](PRD-v0.5.md) — outcome, scope, users, exclusions, proposed acceptance criteria.
- [`ARCHITECTURE-v0.5.md`](ARCHITECTURE-v0.5.md) — component boundaries and local system context.
- [ADR-0004](../../adr/ADR-0004-local-scan-history.md) — proposed SQLite/storage contract; exact technical decisions remain open. External-review gates are superseded by [ADR-0005](../../adr/ADR-0005-solo-owner-review-policy.md).
- [`UX-DESIGN.md`](UX-DESIGN.md) — workflows, information architecture, shadcn/ui-inspired direction, interaction and accessibility states.
- [`API-DESIGN.md`](API-DESIGN.md) — proposed local REST resource and operation contract; OpenAPI artifact remains a design-gate deliverable.
- [`requirements/FUNCTIONAL-REQUIREMENTS.md`](requirements/FUNCTIONAL-REQUIREMENTS.md) — observable behaviors and release gates.
- [`requirements/FAILURE-AND-DATA-CONTRACT.md`](requirements/FAILURE-AND-DATA-CONTRACT.md) — proposed persistence and failure semantics.
- [`requirements/ERROR-MODEL.md`](requirements/ERROR-MODEL.md) — candidate history/API error categories; codes/status mapping remain draft.
- [`requirements/SECURITY-REQUIREMENTS.md`](requirements/SECURITY-REQUIREMENTS.md) — local threat boundaries and security acceptance.
- [`requirements/COMPATIBILITY-MATRIX.md`](requirements/COMPATIBILITY-MATRIX.md) — compatibility commitments and unresolved decisions.
- [`requirements/TEST-STRATEGY.md`](requirements/TEST-STRATEGY.md) — contract, migration, UI, integration, and platform evidence.
- [`CONTRACT-CROSSWALK.md`](CONTRACT-CROSSWALK.md) — source-contract and FR-501–FR-511 traceability; open gates remain explicit.

- [`epics/README.md`](epics/README.md) and [`issues/README.md`](issues/README.md) — EPIC-001 / #390 and its planning-only subissue contracts.

## Release boundary

v0.4: approved remediation plan → isolated apply → verification → rescan → patch evidence.
v0.5 proposal: local scan history → local API → embedded browser console.
v0.6: team collaboration, PostgreSQL, identity, RBAC, and exceptions (not included here).

The v0.5 proposal is local-first and read-only from the browser. It does not add scan triggering, remediation, repository mutation, remote publishing, or team services. Exact endpoints, persistence schema, retention policy, UI routes, component dependencies, and serving command behavior are unapproved until the design gates below pass.

## Design-first gates

1. Reconcile user workflows and domain/resource vocabulary; disposition each unresolved #356 predecessor item.
2. The owner reviews and accepts [`UX-DESIGN.md`](UX-DESIGN.md), including the requested shadcn/ui style, accessibility states, and frontend dependency choice; external design review is optional.
3. The owner reviews and accepts [`API-DESIGN.md`](API-DESIGN.md), including OpenAPI 3.1, listener/origin/security, errors, bounds, and compatibility behavior.
4. The owner reviews and accepts the SQLite schema, lifecycle, migration, and recovery decision. Record the accepted decision as an ADR before storage implementation.
5. Complete security, failure/data, test, and compatibility contracts; record the owner's security/architecture assessment. No independent reviewer is required.
6. Derive ordered epics and implementation issues from these accepted contracts and update the v0.5 plan and Definition of Ready.
7. The owner records acceptance of the complete technical Definition of Ready; external approval is optional.

UX and API design may iterate in parallel only after their shared resource concepts are reconciled. No runtime implementation issue may be created or started before the technical gates and owner acceptance pass. The package does not mark technical gates complete without linked evidence.

## Solo-owner review policy

[ADR-0005](../../adr/ADR-0005-solo-owner-review-policy.md) supersedes all independent-review and dual-approval gates in this active package. The owner is the required reviewer and decision-maker. External review is optional; automated checks, security evidence, and explicit risk dispositions remain required.

The four technical v0.4 follow-ups remain open in [#356](https://github.com/geoffrey-xiao/deprail/issues/356). Its fifth, independent-review follow-up is preserved as history and superseded by ADR-0005; this does not complete or close the other four items.

## Planning-ticket history

The owner created the v0.5 planning epic and design/readiness subissues before the current solo-owner policy. [ADR-0005](../../adr/ADR-0005-solo-owner-review-policy.md) now makes external review optional for these and future contracts; technical design, evidence, and owner acceptance remain required.

## Definition of Ready

- [ ] Owner approves the v0.5 outcome, inclusions, exclusions, compatibility boundary, failure model, and preview acceptance gate.
- [ ] Every open #356 follow-up is explicitly assigned, carried forward, or deferred with owner and target; no predecessor gate is inferred complete.
- [ ] UX design is approved and browser interaction/accessibility criteria are observable.
- [ ] OpenAPI/API security design is reviewed; endpoint and data schemas are versioned and validated.
- [ ] SQLite schema, transaction, migration, retention/deletion, corruption, backup/recovery, and permissions decision is accepted.
- [ ] Security and failure contracts include local-origin, path, process, privacy, and sensitive-data boundaries.
- [ ] Compatibility matrix, cross-platform test strategy, packaging constraints, and rollback/recovery evidence are accepted.
- [ ] Owner reviews and records a decision on the exact current contract version; an independent reviewer is optional.
- [ ] Ordered issue contracts each have one primary outcome, dependencies, inputs/outputs, failures, acceptance evidence, owner reviewer, and rollback notes.
- [ ] Project/milestone/view metadata and evidence links are verified.

### Readiness evidence and blockers

The owner separately approved the then-current package in [#392](https://github.com/geoffrey-xiao/deprail/issues/392#issuecomment-5807906998), and the independent reviewer approved that earlier package in [PR #412](https://github.com/geoffrey-xiao/deprail/pull/412). These historical approvals predate the current [`OpenAPI 3.1 candidate`](../../../schemas/openapi/v1/openapi.yaml). PR #414 merged the later schema/security corrections; no independent APPROVED review is recorded, but [ADR-0005](../../adr/ADR-0005-solo-owner-review-policy.md) makes it optional. The owner must still accept the current technical decisions and link their evidence.

| Gate | Current evidence, not acceptance | Owner's next action / blocking decision |
| --- | --- | --- |
| Scope and predecessor, plan §12 items 1–2 | Draft plan [PR #388](https://github.com/geoffrey-xiao/deprail/pull/388), package [PR #389](https://github.com/geoffrey-xiao/deprail/pull/389), open [#387](https://github.com/geoffrey-xiao/deprail/issues/387). The owner assigned/targeted five #356 follow-ups in the [v0.4 retrospective](../../retrospectives/RETROSPECTIVE-v0.4.0-preview.2.md#follow-up-actions); ADR-0005 now supersedes the independent-review item, leaving four technical follow-ups open in #356. | Reconcile and accept the plan and each technical predecessor disposition in #387 without claiming the v0.4 gaps passed or silently importing them into v0.5. |
| Owner scope, compatibility, and preview gate | Proposed included/excluded work and release gate are in the [plan](../../03-planning/deprail-development-plan-v0.5.0.md#3-included-scope-proposed), [PRD](PRD-v0.5.md), and [crosswalk](CONTRACT-CROSSWALK.md#remaining-decisions-and-approval-gates). | Explicitly accept or revise the release boundary, history capture/CLI exit behavior, measured size limits or retention risk, and preview acceptance criteria; record any material change through the planning hierarchy. |
| UX and accessibility | [#394](https://github.com/geoffrey-xiao/deprail/issues/394) is closed; the [UX proposal](UX-DESIGN.md) and mockups were part of the full-package review approved in [PR #412](https://github.com/geoffrey-xiao/deprail/pull/412). | Component/dependency/tokens, supported browser/OS/assistive-technology matrix and actual-surface keyboard, screen-reader, contrast and responsive evidence remain to resolve; package review does not claim browser verification. |
| API/OpenAPI and local listener | [PR #413](https://github.com/geoffrey-xiao/deprail/pull/413) merged the candidate; the owner merged five schema/security corrections in [PR #414](https://github.com/geoffrey-xiao/deprail/pull/414). Its body accidentally triggered auto-close of [#396](https://github.com/geoffrey-xiao/deprail/issues/396); the issue is reopened and its Project status is Review. | No independent review of PR #414 is recorded; it is not a blocker under ADR-0005. Owner decisions remain on exact `history-v1` fields, safe diagnostics, SQL validation, row/API mapping, source-version compatibility, listener/token exposure, response bounds, and supported browser/platform evidence. |
| SQLite and data lifecycle | [#391](https://github.com/geoffrey-xiao/deprail/issues/391) closed after owner-approved [PR #407](https://github.com/geoffrey-xiao/deprail/pull/407); the owner-selected `history-v1` direction is in [merged PR #411](https://github.com/geoffrey-xiao/deprail/pull/411), and the then-current package was reviewed in [PR #412](https://github.com/geoffrey-xiao/deprail/pull/412). The ADR remains Proposed; no SQL/runtime schema is approved. | Exact projection fields, row validation, migration/lifecycle, retention-growth risk, opt-in capture and persistence-failure precedence remain explicit design decisions; review permissions, backups, recovery and artifact contracts before implementation. |
| Security, failure, test, compatibility and rollback | Draft [PR #408](https://github.com/geoffrey-xiao/deprail/pull/408) merged with three-OS CI; [#393](https://github.com/geoffrey-xiao/deprail/issues/393) remains in Review. The then-current threat/test contracts were included in PR #412's package review; FR/SEC/STORE scenarios remain planned evidence, not executed runtime checks. | Review the newer API local-origin/auth/error profile and failure/rollback changes; pin/select toolchains, SQLite driver, frontend/browser matrix, response limits and recovery owner/procedure. |
| Owner review and DoR sign-off | [#392](https://github.com/geoffrey-xiao/deprail/issues/392) records owner authorization to create planning tasks; PR #412's external approval and the owner comment both predate PR #414. PR #414 was merged by the owner; no independent APPROVED review is recorded, and none is required under ADR-0005. | Record the owner's decisions on the current technical contracts and complete the evidence-backed DoR; no second-person approval is required. |
| Ordered issue contracts and tracking | [#392](https://github.com/geoffrey-xiao/deprail/issues/392) is open in Project Review with [PR #409](https://github.com/geoffrey-xiao/deprail/pull/409); its [conditional delivery map](issues/V05-005-readiness-map.md#conditional-post-dor-delivery-map) is not an issue backlog. [EPIC-001 / #390](https://github.com/geoffrey-xiao/deprail/issues/390) has all six planning children; the [v0.5.0 milestone](https://github.com/geoffrey-xiao/deprail/milestone/11) and [Project view](https://github.com/orgs/geoffrey-xiao/projects/1/views/6) exist, with view filter `milestone:v0.5.0`. | After technical decisions and owner DoR acceptance, finalize one-outcome issue contracts with owner reviewer, risk, acceptance/evidence/rollback and verified labels/milestone/Project links. Do not create or activate implementation issues beforehand. |

**Current technical decisions still open:** The actual scan report is neither privacy-safe unchanged JSON nor v1alpha finding/error-schema-valid. [ADR-0004 records the validation finding and owner-selected separate projection](../../adr/ADR-0004-local-scan-history.md#owner-selected-direction-independent-history-projection-not-accepted). PR #414 merged the five schema corrections; its merge closed #396 and the issue was reopened after the accidental auto-close. External review is optional under ADR-0005. Exact `history-v1` fields, safe-diagnostic handling, SQL validation, row/API mapping, source-version compatibility, and the current security/bounds decisions remain owner DoR gates; no runtime work is authorized yet.

### Conservative technical recommendations (draft; not accepted)

These defaults consolidate the current proposals for owner review; they are not approved contracts or implementation authorization. Record an explicit owner decision and required evidence before changing any Definition-of-Ready checkbox.

| Decision | Recommended default | Acceptance/evidence still required |
| --- | --- | --- |
| Scope and scan capture | Keep v0.5 local, read-only, and limited to history/console. Save history only through an explicit `deprail scan --save-history`; existing scans remain unchanged by default. On a requested save failure, preserve the scan result, emit a typed stderr diagnostic, and use candidate exit code `6` only when scanning otherwise succeeds; keep the scan failure primary if both fail. | Approve the additive flag and exit-code mapping against the [v0.1 CLI contract](../deprail-v0.1-execution-package/requirements/CLI-CONTRACT.md); test all scan/persistence outcome combinations. |
| Persisted projection | Keep source `ScanReport` and CLI JSON unchanged. Use one separately versioned allowlisted `history-v1` projection; exclude absolute roots, raw error strings, secrets, and unclassified unknown fields. Obtain graph context from the same scan operation; if unavailable, store null workspace context rather than rediscovering or guessing. Reject unsafe/incomplete projections without changing the scan outcome or prior database. | Owner accepts exact field mapping, stable finding identity, typed diagnostics, row/projection checks, safe unknown-data policy, and evidence on real nonempty, error, partial, and failed reports. |
| Local API and response limits | Keep an ephemeral `127.0.0.1` listener, read-only `GET` routes, exact Host/Origin checks, no CORS, and process-scoped memory-only auth. If fragment bootstrap remains, scrub it before any request and record its brief browser/OS launch exposure as residual risk. Enforce the 1 MiB serialized UTF-8 response ceiling without truncation; return whole records only, page within the byte budget, and return a typed error if one record alone exceeds it. | The current 25-item worst-case sample exceeds 1 MiB. Measure representative payloads, Unicode/escaping, exact boundaries, request/load bounds, and supported platform behavior before accepting page/count/per-entry limits or the token bootstrap. |
| SQLite lifecycle and retention | Keep a canonical per-user data root outside scanned repositories, owner-only permissions, WAL with `synchronous=FULL`, one serialized writer, atomic entry/reference transactions, and validated SQLite online-backup snapshots before migrations. Do not auto-repair, downgrade, restore, evict, delete history, or delete shared artifacts. For the preview, prefer explicit opt-in plus documented manual cleanup over automatic deletion. | Select and test the driver/build strategy, finite busy timeout, migration/backup/restore procedure, and measured entry/store bounds—or explicitly accept unbounded growth. Any bound must refuse a new save visibly rather than silently prune prior records. |
| UX and platform support | Avoid new third-party runtime UI dependencies; use project-owned semantic HTML/CSS and local assets, with WCAG 2.2 AA as the target. Start with one explicit browser/assistive-technology pair per required OS: Chrome/NVDA on Windows, Safari/VoiceOver on macOS, Firefox/Orca on Linux. | Verify the actual surface: keyboard/focus, screen-reader status, contrast, 320 CSS px/200% zoom, reduced motion, hostile repository text, and browser Back/Forward on each declared combination. Versions and any unsupported combinations must be recorded. |
| #356 rollback/recovery | Name `@geoffrey-xiao` as the rollback owner; document immutable-tag recovery and manual database/artifact preservation. | Exercise the recovery procedure before stable release; preview gaps require explicit owner disposition and do not count as passed evidence. |
| #356 Python/Java coverage | Defer full remediation coverage outside v0.5 because this release adds read-only history, not remediation. Preserve language-neutral history behavior and exercise representative JavaScript, Python, and Java scan results where fixtures/workflows exist. | Record the explicit deferral in #356 and the release plan; do not imply v0.5 delivers remediation for those ecosystems. |
| #356 platform and provenance | Require representative local workflow smoke on Linux, macOS, and Windows before stable release. Record SBOM, signature, and provenance as supplied, unavailable, or deferred for every preview; never infer or fabricate them. | Verify checksums and available provenance against the immutable reviewed commit. Any stable-release gap requires an explicit owner disposition in the versioned release plan/evidence. |

**Status: awaiting owner decisions.** These recommendations do not close #387, #393, #396, or #356 and do not authorize implementation issues. Keep technical DoR items unchecked until the owner accepts or revises the choices and the evidence above exists; independent review is optional under [ADR-0005](../../adr/ADR-0005-solo-owner-review-policy.md).

## Runtime work remains gated by technical readiness

The owner authorized planning-task creation in [#392](https://github.com/geoffrey-xiao/deprail/issues/392#issuecomment-5807906998), and the independent reviewer approved the earlier package in PR #412; both are historical facts, not current contract acceptance. PR #414 was merged by the owner without an independent APPROVED review; that review is optional under ADR-0005. Do not begin runtime work until the remaining technical decisions, linked evidence, and owner acceptance of the complete DoR are recorded.
