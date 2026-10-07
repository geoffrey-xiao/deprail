# H05-008 Release Verification Readiness

Issue: [#427](https://github.com/geoffrey-xiao/deprail/issues/427). Procedure: [manual verification guide](../MANUAL-TEST-GUIDE.md). Assessment date:2026-10-07. Current synchronized implementation baseline: owner-merged main `c68c9f9`, including #455 timed-apply, #456 artifact fix and [#458 authoritative Release Combined repair](https://github.com/geoffrey-xiao/deprail/pull/458). Earlier baseline/failure evidence is preserved below. **Formal H05-008 runtime execution remains blocked by original prerequisite acceptance/dispositions**, not the now-delivered correction work. No release approval, owner no-go decision, tag or publication asserted.

Preparation: [PR #445](https://github.com/geoffrey-xiao/deprail/pull/445), Refs #427; [initial handoff](https://github.com/geoffrey-xiao/deprail/issues/427#issuecomment-6028034470). PR remains Draft and #427 Project Blocked. Corrective PRs are owner-merged; local/component proof is not blanket storage/security/browser-AT/supply-chain or formal release acceptance.

## Historical CI — preserved

Preparation CI [37551712270](https://github.com/geoffrey-xiao/deprail/actions/runs/37551712270) on `33f7944`: Ubuntu and macOS passed; Windows failed `go test ./...`, specifically `TestApplyEndToEndFailureBoundaries/cancellation` at `cmd/deprail/apply_e2e_test.go:215` with `unexpected end of JSON input`. Windows package build/native smoke did not run after that failure. Cause is not established by this log; no flaky-test classification, passing rerun, weaker assertion or runtime fix is claimed. This adds an unresolved verification blocker, distinct from the previously resolved #432 QA correction.

Preparation head `165e1ed` subsequently passed Ubuntu/macOS/Windows in [37552087392](https://github.com/geoffrey-xiao/deprail/actions/runs/37552087392), including native packaged smoke. UI head `d4abb39` passed three-OS CI/native smoke in [37568426603](https://github.com/geoffrey-xiao/deprail/actions/runs/37568426603) before owner-merged #451. Neither pass erases earlier failures or proves the final integrated candidate.

## Current blocker corrections and evidence

- Owner reported local testing OK after the row-only UI change. [Recorded bounded feedback](https://github.com/geoffrey-xiao/deprail/issues/427#issuecomment-6030612401): local UI smoke, not an identified release-binary/browser-AT/security/recovery matrix or publication authorization. #446/#448/#450 remain separate criterion-review records.
- #452 / #453 corrected the default optional artifact resolver (nil *artifact.Store in a nonnil interface). Four real-store/listener cases, make verify/integration and actual Darwin console passed; default detail retained five genuine findings with Unavailable, explicit root retained Verified, copied DB hash unchanged. Original diagnostic binary693eb8b9310f050ed26b39cd51d39dd420392c87dabb5ed198e51b848cb4605d. #453 merged into its temporary base; owner-merged #456 subsequently delivered the correction to main c554441. #452 is closed; no artifact-cutover blocker remains.
- Default-console exact head `bc4d732` [37569906460](https://github.com/geoffrey-xiao/deprail/actions/runs/37569906460): Ubuntu/macOS passed; Windows failed existing timed-apply timeout at apply_e2e_test.go:225 with unexpected end of JSON input; later Windows build/smoke skipped. No passing Windows claim or rerun.
- [#454 / PR #455](https://github.com/geoffrey-xiao/deprail/pull/455) isolates test-only timed-apply synchronization. Fixed500ms cancel and1s startup-inclusive deadline had no mutation-start barrier. Controlled preflight-cancellation smoke produced exit3/stdout0/PATH_OUTSIDE_ROOT and the exact decode error; original Windows stderr was absent, so original failed phase is not claimed. Both timed cases now observe a real mutation child ready marker before cancellation or an actual1s deadline. Original outcomes/exit/evidence/cleanup assertions preserved; persisted evidence validation, unchanged caller and worktree removal added. Focused actual-process/race, make verify and integration passed locally. Exact-head CI recorded on #455; no production JSON/error change.
- Cutover history: #453 stacked head 1644f30 passed three-OS CI/native CLI smoke and Ubuntu four-target cross-build in [37578685087](https://github.com/geoffrey-xiao/deprail/actions/runs/37578685087). Owner then merged #455 to main as 9181e76, but #453 to the still-temporary branch as 983c604; that merge did not deliver the artifact fix to main. #456 starts from synchronized main, cherry-picks only the artifact squash, and requires its own exact-head CI and owner review/merge. Its local real-process/HTTP scenarios, make verify and integration passed. No release verification from an isolated branch or a main missing the artifact fix; no autonomous merge.
- #456 exact head 0e3bc07 passed [37579457745](https://github.com/geoffrey-xiao/deprail/actions/runs/37579457745): all three OS tests/build/native CLI smoke and Ubuntu four-target cross-build. Owner then merged its bounded main-targeted correction as c554441. Earlier #453 failed/stacked runs remain history, not a current unmerged requirement.
- [#457 / #458](https://github.com/geoffrey-xiao/deprail/pull/458) repaired the owner-selected Release Combined: pinned frontend builds, one prepared SHA, complete private asset assembly before approval, and draft/upload/public transition. Exact head f83fc89 passed [37582051709](https://github.com/geoffrey-xiao/deprail/actions/runs/37582051709) on all three OS with native CLI smoke and Ubuntu cross-build; [local syntax/build/offline failure proof](H05-007-FU1-RELEASE-COMBINED-EVIDENCE.md) states genuine historical SPDX fixture/gh-recorder limits. Owner merged as c68c9f9; [completion record](https://github.com/geoffrey-xiao/deprail/issues/457#issuecomment-6034157699), Project Done, owner may close #457. Real Release Combined dispatch/SBOM/publication remains unexecuted.


## Reconciled prerequisite state

| Contract | Observed tracking state | Evidence and remaining gate |
| --- | --- | --- |
| QA #419 | Closed / Project Done | Historical investigation is preserved. [#432 explicit acceptance](https://github.com/geoffrey-xiao/deprail/issues/432#issuecomment-5913275072), merged #433 and the backlog's accepted kickoff addendum resolve the QA correction gate. Do not continue describing the original OutputCap correction or timeout fixture as unimplemented; original failures are not erased. |
| H05-001 #420 | Closed / Done | Prior accepted projection slice; integrated candidate behavior still requires H05-008 proof. |
| H05-002 #421 | Open / Review | [Store evidence](H05-002-STORE-EVIDENCE.md): permission/dependency/advisory review, runtime physical-growth evidence and manual restore rehearsal remain incomplete. Release-plan §15 already accepts the logical-versus-physical growth tradeoff under the private/no-delete boundaries; do not invent a missing planning decision or call it a hard disk quota. |
| H05-003 #422 | Closed / Done | Prior accepted service slice; integrated evidence still required. |
| H05-004 #423 | Open / Review | CLI capture acceptance has no completed issue disposition in the inspected live tracking. |
| H05-005 #424 / follow-ups #441, #452 | #424/#441 Open / Review; #452 Closed / Done | Transport/security criterion disposition remains outstanding. Artifact correction is delivered to main by owner-merged #456, with exact-head three-OS CI; no remaining cutover requirement. |
| H05-006 #425 / #446, #448, #450 | Open / Review | [UI evidence](H05-006-UI-EVIDENCE.md), merged #443/#447/#449/#451, real packaged headless desktop/320px/keyboard proof and owner local-test feedback exist. Actual supported browser/AT, physical zoom and supply-chain decisions remain incomplete; local success is not blanket release acceptance. |
| H05-007 #426 | Open / Review | Merged #444; exact-head CI [37549943755](https://github.com/geoffrey-xiao/deprail/actions/runs/37549943755) passed all three OS jobs at `5c78da1`. Automated review reported no major issues at that commit. [Packaging evidence](H05-007-PACKAGING-EVIDENCE.md) is component evidence, not release-binded native browser/platform or supply-chain approval. |
| Timed-apply QA #454 | Closed / Project Done | Owner-merged #455 to main 9181e76; [completion evidence](https://github.com/geoffrey-xiao/deprail/issues/454#issuecomment-6032036180), exact-head three-OS CI [37571076086](https://github.com/geoffrey-xiao/deprail/actions/runs/37571076086), focused/race/full/integration proof and original phase uncertainty retained. Delivered QA is not release acceptance. |
| Release Combined #457 | Open / Project Done; owner may close | Owner-merged #458 at c68c9f9; local actual-build/offline boundary proof and exact-head PR CI recorded. Approval environment has required owner reviewer geoffrey-xiao, prevent_self_review=false. Actual publication workflow/SBOM/release API execution remains unverified. |
| H05-008 #427 | Open / Blocked; preparation #445 Draft | Requires original H05-001–007/follow-up acceptance or dispositions, selected immutable candidate/version and final integrated runtime/security/recovery/platform/supply-chain evidence. Correction cutovers are complete; this record does not satisfy release runtime criteria. |

A merged PR and successful CI do not by themselves close unchecked acceptance criteria. No native APPROVED review is invented from the owner merging #444. User requests to continue authorize preparation, not automatic acceptance of missing platform or supply-chain evidence.

## Required owner dispositions and remaining evidence

1. Record acceptance or explicit bounded disposition for #421 and #423–#426. Keep distinct technical and security decisions where required.
   Main c68c9f9 now contains all three owner-merged correction PRs #455/#456/#458. No repeated merge is required; earlier #419/#432 gates remain resolved.
2. Resolve the disclosed Go dependency advisory, frontend `caniuse-lite` CC-BY-4.0 attribution/distribution implications and unverified publisher provenance. Existing registry integrity/audit results are not those decisions.
3. Supply actual authoritative release-binded candidate binaries/identity, or explicitly select a reviewed local candidate for prepublication verification. No v0.5 tag or suffix is chosen by this record.
4. Execute integrated JS/Python/Java scan/save/reopen/API/UI, security/read-only tree comparisons and offline recovery from disposable copies under the guide. Historical component tests do not prove this matrix.
5. Execute or explicitly disposition current stable Chrome/NVDA, Safari/VoiceOver and Firefox/Orca, and all four native binary targets. Cross-builds, Rosetta identity-only smoke and headless Chromium are not interchangeable with full native/browser/AT verification.
6. Record full artifact SBOM/signature/provenance status and separate owner technical, security/architecture and release decisions. Missing evidence remains a gate, not an implicitly approved preview gap.

### Requested preview gap proposal — not accepted dispositions

Owner selected preparation of a preview-with-gaps proposal, not blanket acceptance or publication. [PG-01–PG-12 proposal](H05-008-PREVIEW-GAP-PROPOSAL.md) names exact evidence, impact, mitigation, accountable owner and target gate for every gap. Native/browser-AT/recovery/physical-growth and provenance gaps are candidates for specific owner deferral; disclosed advisory, required attribution, original admission, core runtime and actual workflow/SBOM success are not automatically waivable. Proposed v0.5.0-preview.1 remains unselected/uncreated. All decision rows are pending; #427 stays Blocked.


## Release-checklist coverage

| General checklist section | Current disposition |
| --- | --- |
| Release identity | v0.5.0 preview-first plan exists; exact preview tag, release manifest and final artifact identity are not selected/verified here. |
| Plan and scope | Product/UI and corrective QA/artifact/publication repairs are merged to main c68c9f9. Original prerequisite acceptance/dispositions still prevent a completed release-scope claim. |
| Contracts and security | Existing approved contracts and component evidence retained; integrated runtime/privacy/recovery and owner risk decisions pending. |
| Automated/manual verification | Historical failures preserved; correction exact-head PR CI and local runtime/build proof passed. No final candidate/integrated release matrix or real Release Combined workflow pass inferred. |
| Artifacts/supply chain | H05-007 local hashes/budgets and host SBOM available; actual release artifacts and signature/provenance/full supply-chain dispositions unverified. |
| Publication/approval | Live release-approval required owner-reviewer rule verified, self-review allowed by ADR-0005; configuration unchanged. Actual dispatch/tag/publication is neither authorized nor executed. |
| Post-release | Not applicable before publication; no retrospective or release closure claimed. |

Technical readiness assessment: **blocked / not release-ready**. This is not an owner-issued release no-go. Preserve existing user DB/WAL/SHM/artifacts and immutable tags; do not downgrade, reset or overwrite to obtain a passing result.
