# V05-005: Assemble the Execution Package and Readiness Map
- GitHub Issue: [#392](https://github.com/geoffrey-xiao/deprail/issues/392)

- Epic: [EPIC-001 / #390](https://github.com/geoffrey-xiao/deprail/issues/390)
- Target: `v0.5.0`
- Status: Open in Project Review; [PR #417](https://github.com/geoffrey-xiao/deprail/pull/417) reconciles the exact engineering proposal, schema/examples, boundary evidence and eight-slice unpublished map. The owner must accept the current contract and separately disposition security/architecture risk before implementation issues are published. Earlier PR #409/#412/#414 approvals remain historical; external review is optional under ADR-0005.
- Type: decision
- Area: docs
- Priority: P0
- Risk: R3
- Owner: `@geoffrey-xiao`
- Reviewer: `@geoffrey-xiao` (owner); independent review is optional under [ADR-0005](../../../adr/ADR-0005-solo-owner-review-policy.md).
- Dependencies: V05-000 through V05-004

## Value

Implementation tickets must derive from accepted design and failure contracts, with explicit dependencies, evidence, compatibility, rollout, and rollback behavior; none should be inferred from a draft proposal.

## Scope

Integrate owner-accepted UX, OpenAPI/API security, SQLite ADR/schema/lifecycle, failure/error, security/threat, compatibility, and test strategy into the v0.5 execution package. Reconcile source crosswalks and the v0.5 Definition of Ready; map FR-501–FR-511 and product/architecture/roadmap requirements to ordered implementation epics/issues; record one primary outcome, inputs/outputs, failures, required tests, observable acceptance, owner/reviewer, risk, evidence, rollback, milestone, labels, and GitHub parent hierarchy for each proposed implementation issue. Preserve the #356 history, link its superseded external-review disposition, and keep its four technical follow-ups separately tracked.

## Out of scope

No runtime code, implementation branch, issue publication or activation, release tag/publication, or #356 closure. The owner authorized task creation in #392 before PR #414; do not publish implementation issues until current API design decisions and the technical Definition of Ready are accepted by the owner.

## Inputs, outputs, and failure behavior

Inputs: V05-000 through V05-004, the reconciled planning hierarchy, all higher-level contracts, and current GitHub milestone/project state. Outputs: a complete execution package and a traceable draft issue map. Any missing design artifact, owner decision, test evidence, or compatibility disposition leaves its DoR item unchecked and runtime implementation blocked; an optional external reviewer is not a readiness dependency.

## Required verification and evidence

Validate all local links and crosswalks; verify every proposed implementation issue maps to a frozen contract and has a single outcome and observable acceptance; verify milestone, required labels, subissue relationships, Project status, and view filters. No runtime verification is asserted.

## Acceptance criteria

- The v0.5 execution package is internally consistent with product design, architecture, roadmap, release plan, and #356.
- Every v0.5 DoR item is either supported by linked evidence or explicitly remains unchecked with an owner and next action.
- Any implementation epic/issue proposal remains unpublished/unstarted until the owner accepts the complete technical DoR; independent review is optional under ADR-0005.
- Project/milestone links and local/GitHub issue hierarchy are reconciled; no child is orphaned or attached to the wrong epic.

## Readiness snapshot

The [execution-package DoR map](../README.md#readiness-evidence-and-blockers) links current evidence and blockers. UX [#394](https://github.com/geoffrey-xiao/deprail/issues/394) and SQLite [#391](https://github.com/geoffrey-xiao/deprail/issues/391) are closed; API [#396](https://github.com/geoffrey-xiao/deprail/issues/396) is open in Project Review; quality [#393](https://github.com/geoffrey-xiao/deprail/issues/393) remains in Review. PR #412's independent approval applied to the earlier package only. No independent approval of PR #414 is recorded, and none is required under ADR-0005. The former fifth #356 independent-review follow-up is superseded and owner-dispositioned; four technical follow-ups remain open. #387 and #392 remain open; current technical decisions and owner DoR acceptance are still pending.

## Conditional post-DoR delivery map

The original A–F table below is preserved as the earlier conditional grouping. The current unpublished planning decomposition in §Current review decomposition refines that grouping for owner review; it is not a GitHub implementation backlog or runtime authorization. The owner must accept the current contracts and planning DoR before implementation issue publication, then each published issue must satisfy its own DoR before code starts. Proposed names/files below identify future ownership boundaries, not existing implemented APIs or packages.

| Tentative order and one primary outcome | Contract and input gate | Observable output, failure, and evidence | Rollback / risk |
| --- | --- | --- | --- |
| A. Durable local history store | `FR-501`, `FR-508`, `STORE-01`, `SEC-08`, `SEC-09`; accepted [ADR-0004](../../../adr/ADR-0004-local-scan-history.md), numeric bounds, permissions/backup/recovery and driver selection. | Distinct occurrence IDs even for repeated scan IDs, atomic entry+digest references, explicit missing/corrupt/future-schema/lock/disk-full outcomes; fresh/interrupted/rollback and three-OS permission evidence. | Preserve committed DB, WAL and validated backup; no automatic repair or artifact deletion. R3. |
| B. Shared history capture and query service | `FR-501`–`FR-506`, `STORE-01`; A, approved capture trigger and CLI compatibility decision. | Selected operation saved once with truthful outcome/report states; read-only listing/detail retain provenance; failed save never claims persistence or overwrites the scan result; repeated/partial/cancelled/CLI-without-history scenarios. | Withdraw the selected capture integration without changing default CLI; retain existing DB and artifacts. R3. |
| C. Read-only local API transport | `FR-502`–`FR-507`, `FR-511`, `SEC-01`–`SEC-05`, `SEC-07`, `SEC-11`; B and validated OpenAPI/listener/security decisions. | Accepted schema and hostile-origin/path/SQL/size/timeout tests show bounded typed responses, no mutation or empty-success substitution; CLI and HTTP preserve shared domain meaning. | Stop/disable the listener; prior CLI and stored history remain usable. R3. |
| D. Accessible browser history and detail | `FR-502`–`FR-504`, `FR-510`, `FR-511`, `SEC-06`; C, accepted UX primitives, tokens and supported browser/assistive-technology matrix. | Actual-surface history/detail/help, keyboard/focus/status and hostile-text checks; complete/partial/failed/cancelled, empty/unavailable/incompatible and missing-evidence states remain distinct. | Remove the browser client without erasing local history or altering CLI meaning. R2. |
| E. Reproducible embedded delivery | `FR-505`, `FR-509`, `SEC-10`, `SEC-12`, `SEC-13`; C–D, selected pinned toolchains/dependencies and packaging targets. | Version-matched assets/API, visible mismatch and missing-asset failures, no filesystem fallback/CDN; clean builds and package identity/size/startup evidence on Linux/macOS/Windows. | Roll back to a matching reviewed binary/assets without downgrading or replacing the store. R2. |
| F. Release verification and evidence | `FR-501`–`FR-511`, `SEC-01`–`SEC-13`, `STORE-01`; A–E and [release checklist](../../../RELEASE-CHECKLIST.md). | Real representative workflows, browser and hostile-input evidence, release binary/tree comparison, version/checksum/SBOM/signature/provenance dispositions, distinct owner security assessment and release decision; no unit-only release claim. | Document rollback owner, preserve immutable tags/data and publish a new corrected preview rather than moving a tag. R3. |

Proposed milestone for any future issue is [`v0.5.0`](https://github.com/geoffrey-xiao/deprail/milestone/11); each would be linked beneath an approved implementation epic in the [DepRail Project](https://github.com/orgs/geoffrey-xiao/projects/1/views/6). `area:storage`, `area:api`, and `area:web` are not in the current label set. The owner authorized task creation in #392 before PR #414; ADR-0005 now makes independent review optional. Under the current plan, do not publish implementation issues until the remaining technical decisions, owner acceptance, and complete evidence-backed DoR are reconciled. The accepted contracts determine final issue boundaries and sequencing.

## Tracking verification (planning only)

- `gh api repos/geoffrey-xiao/deprail/issues/390/sub_issues` lists exactly the six planning children #387, #394, #396, #391, #393, and #392. No candidate delivery outcome above has a GitHub implementation issue.
- `gh project view 1 --owner geoffrey-xiao` resolves the live DepRail Project; its `Release · v0.5.0` view is number `6`, filtered to `milestone:v0.5.0`. #392 is open in Project Review after PR #409 merged, with the v0.5.0 milestone and required labels; its status is workflow state, not DoR approval. The earlier owner authorization predates PR #414; ADR-0005 now removes independent-review authorization as a gate.
- [#394](https://github.com/geoffrey-xiao/deprail/issues/394) and [#391](https://github.com/geoffrey-xiao/deprail/issues/391) are closed; [#396](https://github.com/geoffrey-xiao/deprail/issues/396) is open in Project Review after correction; [#393](https://github.com/geoffrey-xiao/deprail/issues/393) remains open in Review; [#387](https://github.com/geoffrey-xiao/deprail/issues/387) remains open in In Progress. Closure or Project state is not owner technical acceptance.

## Current review decomposition

All rows have owner and required reviewer `@geoffrey-xiao`, target/milestone `v0.5.0`, priority `P0`, no assigned Sprint until kickoff, and no GitHub implementation issue yet. Proposed interfaces are design names, not implemented symbols. Use existing label vocabulary; no new `area:storage/api/web` labels are required. The rows are an unpublished planning map for the owner's contract/DoR review, not issue contracts ready to activate.

| Order/key | One primary outcome and owned boundary | Contracts/dependencies | Observable acceptance and failure evidence | Labels and rollback |
| --- | --- | --- | --- | --- |
| 1 / H05-001 | Safe deterministic history projection; proposed `internal/domain/history` entities plus `internal/normalize` history transform; history schema/examples. | Failure/data, error model, `FR-501`, `FR-503–504`, `SEC-07`; accepted design prerequisite. | Real nonempty/partial/failed/cancelled source values map to strict allowlists and stable keys; null vs empty, unknown source/diagnostic, duplicate identity, unsafe path and entry boundaries refuse rather than fabricate safe data. | `area:normalization`, `risk:R3`, `type:feature`; withdraw projection/capture support without changing CLI report bytes. |
| 2 / H05-002 | Owner-private durable SQLite history store; proposed `internal/store/history` adapter. | H05-001; ADR-0004, `FR-508`, `STORE-01`, `SEC-02/08`. | Unique occurrences, atomic rows+refs, read-only absent-store behavior, quota and lock/disk-full/cancellation rollback, future/corrupt-schema refusal, private modes/ACL, consistent backup and recovery on declared platforms. | `area:foundation`, `risk:R3`, `type:feature`; preserve DB/WAL/SHM/backups, no reset/downgrade/artifact deletion. |
| 3 / H05-003 | Same-operation capture/query application service; proposed `internal/app` history service and optional scan capture boundary. | H05-001–002; `FR-501–506`, `STORE-01`, artifact integrity contract. | Capture same graph/report/typed diagnostics once; identical-ID retry vs conflicting write; safe early failure/cancellation and bounded terminal capture; read projections expose only trustworthy metadata and explicit artifact integrity. | `area:foundation`, `risk:R3`, `type:feature`; remove optional integration, retain original scan contract and stored data. |
| 4 / H05-004 | Opt-in CLI history capture; `cmd/deprail` scan wiring and documentation. | H05-003; CLI compatibility, `FR-505`, `STORE-01`. | No flag means no DB initialization/write; selected success/failure/cancellation retains original JSON/stdout and primary nonzero scan exit; exit `6` only otherwise-successful scan/save failure with safe stderr and no false saved claim. | `area:cli`, `risk:R3`, `type:feature`; withdraw opt-in flag integration without altering existing scan behavior. |
| 5 / H05-005 | Bounded read-only HTTP/static transport and foreground console command; proposed local transport adapter plus `cmd/deprail` web wiring. | H05-003; exact API/launch/security contract, `FR-502/507/511`, `SEC-01–05/11–12`; build uses fixed contract smoke assets until H05-007 provides reviewed real UI. | Five authenticated GET operations and known embedded static routes; no filesystem fallback; hostile Host/Origin/query/path/body, auth, cursor, byte cap, deadline/admission/shutdown and artifact-root boundary scenarios; no scan/mutation route. | `area:cli`, `risk:R3`, `type:feature`; stop listener, preserve independent CLI and history. |
| 6 / H05-006 | Accessible history/detail/help React UI source; `web/` source and lockfile. | H05-005 contract; UX, compatibility, `FR-502–504/510–511`, `SEC-06`; actual packaged verification depends on H05-007. | Live served surface renders distinct operation/report/empty/unavailable/error states; keyboard/focus/status, hostile text, narrow/zoom, contrast/reduced motion, auth-recovery and Back/Forward checks. Do not claim unit-component tests prove accessibility. | `area:foundation`, `risk:R2`, `type:feature`; withdraw UI assets, keep store/CLI intact. |
| 7 / H05-007 | Reproducible embedded binary delivery; frontend build/embedding and platform packaging configuration. | H05-005–006; compatibility/dependency review, `FR-509`, `SEC-10/12/13`. | Exact lockfile build and license/SBOM/vulnerability disposition; matching assets/API, missing/skew failures, known-route/no-disk fallback, measured asset/startup budgets and four declared binary targets. | `area:foundation`, `risk:R2`, `type:feature`; use matching reviewed binary/assets, never downgrade store automatically. |
| 8 / H05-008 | Integrated release acceptance and evidence, not new features; test/verification procedures and version-specific evidence. | H05-001–007; all FR/SEC/STORE rows and general release checklist. | Representative JS/Python/Java history workflows; Linux/macOS/Windows actual binary and browser/AT results; recovery rehearsal, repository-tree comparison, checksums/SBOM/signature/provenance disposition and separate owner security/release decisions. | `area:test`, `risk:R3`, `type:test`; explicit no-go or approved preview gap; immutable tags preserved, corrected release gets a new tag. |

The references to temporary contract smoke assets in H05-005 define a throwaway verification aid, not a shipped placeholder UI or substitute acceptance. H05-006 and H05-007 together deliver the complete production console; no intermediate scaffold is release-ready.

At publication, convert each row into the repository issue template with exact owned symbols/files, inputs/outputs, scope/exclusions, failure behavior, commands/scenarios, acceptance/evidence, reviewer, rollback, dependencies and live GitHub links. Search existing issues first; publish nothing until owner contract/DoR acceptance. Keep at most two implementation issues/two review PRs active and never parallelize two R3 implementation items.
