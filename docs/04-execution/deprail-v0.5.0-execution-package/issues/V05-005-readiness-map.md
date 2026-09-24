# V05-005: Assemble the Execution Package and Readiness Map
- GitHub Issue: [#392](https://github.com/geoffrey-xiao/deprail/issues/392)

- Epic: [EPIC-001 / #390](https://github.com/geoffrey-xiao/deprail/issues/390)
- Target: `v0.5.0`
- Status: Review in GitHub; draft readiness map in [PR #409](https://github.com/geoffrey-xiao/deprail/pull/409), independent review and v0.5 DoR pending.
- Type: decision
- Area: docs
- Priority: P0
- Risk: R3
- Owner: `@geoffrey-xiao`
- Reviewer: `@geoffrey-xiao` for planning-ticket oversight only, by explicit owner direction; not independent architecture/security approval
- Dependencies: V05-000 through V05-004

## Value

Implementation tickets must derive from accepted design and failure contracts, with explicit dependencies, evidence, compatibility, rollout, and rollback behavior; none should be inferred from a draft proposal.

## Scope

Integrate accepted UX, OpenAPI/API security, SQLite ADR/schema/lifecycle, failure/error, security/threat, compatibility, and test strategy into the v0.5 execution package. Reconcile source crosswalks and the v0.5 Definition of Ready; map FR-501–FR-511 and product/architecture/roadmap requirements to ordered implementation epics/issues; record one primary outcome, inputs/outputs, failures, required tests, observable acceptance, owner/reviewer, risk, evidence, rollback, milestone, labels, and GitHub parent hierarchy for each proposed implementation issue. Preserve #356's separately tracked follow-ups and link evidence.

## Out of scope

No runtime code, implementation branch, implementation issue creation/activation, release tag/publication, #356 closure, or claim of independent approval before the named reviewer records a decision.

## Inputs, outputs, and failure behavior

Inputs: V05-000 through V05-004, the reconciled planning hierarchy, all higher-level contracts, and current GitHub milestone/project state. Outputs: a complete execution package and a traceable draft issue map. Any missing design artifact, reviewer, owner decision, test evidence, or compatibility disposition leaves its DoR item unchecked and runtime implementation blocked.

## Required verification and evidence

Validate all local links and crosswalks; verify every proposed implementation issue maps to a frozen contract and has a single outcome and observable acceptance; verify milestone, required labels, subissue relationships, Project status, and view filters. No runtime verification is asserted.

## Acceptance criteria

- The v0.5 execution package is internally consistent with product design, architecture, roadmap, release plan, and #356.
- Every v0.5 DoR item is either supported by linked evidence or explicitly remains unchecked with an owner and next action.
- Any implementation epic/issue proposal remains unpublished/unstarted until the independent reviewer and owner separately approve the complete DoR.
- Project/milestone links and local/GitHub issue hierarchy are reconciled; no child is orphaned or attached to the wrong epic.

## Readiness snapshot

The [execution-package DoR map](../README.md#readiness-evidence-and-blockers) links every outstanding gate to evidence, owner, and next decision. Draft UX [#394](https://github.com/geoffrey-xiao/deprail/issues/394) and SQLite [#391](https://github.com/geoffrey-xiao/deprail/issues/391) issue closure do not constitute independent approval; API [#396](https://github.com/geoffrey-xiao/deprail/issues/396) and quality [#393](https://github.com/geoffrey-xiao/deprail/issues/393) remain in Review. [#387](https://github.com/geoffrey-xiao/deprail/issues/387) and this issue remain open. The five owner-dispositioned v0.4 follow-ups remain in [#356](https://github.com/geoffrey-xiao/deprail/issues/356), outside v0.5 scope.

## Conditional post-DoR delivery map

These are **candidate outcomes, not issue contracts, GitHub issues, an authorized epic, or a frozen sequence**. They become issue-ready only after approved UX/OpenAPI/SQLite/security/compatibility contracts and separate owner/independent reviewer DoR decisions. A changed design may change the decomposition. Each eventual issue must state exact owned files/interfaces, owner, named reviewer, approved versions/labels, one outcome, dependencies, rollback and linked acceptance evidence before creation; no candidate below grants a new CLI option, schema, endpoint, or dependency.

| Tentative order and one primary outcome | Contract and input gate | Observable output, failure, and evidence | Rollback / risk |
| --- | --- | --- | --- |
| A. Durable local history store | `FR-501`, `FR-508`, `STORE-01`, `SEC-08`, `SEC-09`; accepted [ADR-0004](../../../adr/ADR-0004-local-scan-history.md), numeric bounds, permissions/backup/recovery and driver selection. | Distinct occurrence IDs even for repeated scan IDs, atomic entry+digest references, explicit missing/corrupt/future-schema/lock/disk-full outcomes; fresh/interrupted/rollback and three-OS permission evidence. | Preserve committed DB, WAL and validated backup; no automatic repair or artifact deletion. R3. |
| B. Shared history capture and query service | `FR-501`–`FR-506`, `STORE-01`; A, approved capture trigger and CLI compatibility decision. | Selected operation saved once with truthful outcome/report states; read-only listing/detail retain provenance; failed save never claims persistence or overwrites the scan result; repeated/partial/cancelled/CLI-without-history scenarios. | Withdraw the selected capture integration without changing default CLI; retain existing DB and artifacts. R3. |
| C. Read-only local API transport | `FR-502`–`FR-507`, `FR-511`, `SEC-01`–`SEC-05`, `SEC-07`, `SEC-11`; B and validated OpenAPI/listener/security decisions. | Accepted schema and hostile-origin/path/SQL/size/timeout tests show bounded typed responses, no mutation or empty-success substitution; CLI and HTTP preserve shared domain meaning. | Stop/disable the listener; prior CLI and stored history remain usable. R3. |
| D. Accessible browser history and detail | `FR-502`–`FR-504`, `FR-510`, `FR-511`, `SEC-06`; C, accepted UX primitives, tokens and supported browser/assistive-technology matrix. | Actual-surface history/detail/help, keyboard/focus/status and hostile-text checks; complete/partial/failed/cancelled, empty/unavailable/incompatible and missing-evidence states remain distinct. | Remove the browser client without erasing local history or altering CLI meaning. R2. |
| E. Reproducible embedded delivery | `FR-505`, `FR-509`, `SEC-10`, `SEC-12`, `SEC-13`; C–D, selected pinned toolchains/dependencies and packaging targets. | Version-matched assets/API, visible mismatch and missing-asset failures, no filesystem fallback/CDN; clean builds and package identity/size/startup evidence on Linux/macOS/Windows. | Roll back to a matching reviewed binary/assets without downgrading or replacing the store. R2. |
| F. Release verification and evidence | `FR-501`–`FR-511`, `SEC-01`–`SEC-13`, `STORE-01`; A–E and [release checklist](../../../RELEASE-CHECKLIST.md). | Real representative workflows, browser and hostile-input evidence, release binary/tree comparison, version/checksum/SBOM/signature/provenance dispositions, separate reviewer and owner decisions; no unit-only release claim. | Document rollback owner, preserve immutable tags/data and publish a new corrected preview rather than moving a tag. R3. |

Proposed milestone for any future issue is [`v0.5.0`](https://github.com/geoffrey-xiao/deprail/milestone/11); each would be linked beneath an approved implementation epic in the [DepRail Project](https://github.com/orgs/geoffrey-xiao/projects/1/views/6). `area:storage`, `area:api`, and `area:web` are **not** in the current label set; the owner must select existing matching labels or approve a label-set change before publishing an issue. No reviewer is assigned here: `@geoffrey-xiao` owns the planning disposition, and the independent architecture/security reviewer remains to be named. The accepted contracts—not this provisional table—determine final issue boundaries and sequencing.

## Tracking verification (planning only)

- `gh api repos/geoffrey-xiao/deprail/issues/390/sub_issues` lists exactly the six planning children #387, #394, #396, #391, #393, and #392. No candidate delivery outcome above has a GitHub implementation issue.
- `gh project view 1 --owner geoffrey-xiao` resolves the live DepRail Project; its `Release · v0.5.0` view is number `6` with filter `milestone:v0.5.0`. #392 is open, assigned to the `v0.5.0` milestone, carries `area:docs`, `risk:R3`, `priority:P0`, `type:decision`, and is in Project `Review` after [PR #409](https://github.com/geoffrey-xiao/deprail/pull/409) opened; this is workflow state, not DoR approval.
- [#394](https://github.com/geoffrey-xiao/deprail/issues/394) and [#391](https://github.com/geoffrey-xiao/deprail/issues/391) are closed; [#396](https://github.com/geoffrey-xiao/deprail/issues/396) and [#393](https://github.com/geoffrey-xiao/deprail/issues/393) remain open in Project `Review`; [#387](https://github.com/geoffrey-xiao/deprail/issues/387) remains open in `In Progress`. Closure or Project state is not an independent contract acceptance record.
