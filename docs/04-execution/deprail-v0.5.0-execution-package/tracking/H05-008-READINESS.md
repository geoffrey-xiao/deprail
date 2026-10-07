# H05-008 Release Verification Readiness

Issue: [#427](https://github.com/geoffrey-xiao/deprail/issues/427). Procedure: [manual verification guide](../MANUAL-TEST-GUIDE.md). Assessment date: 2026-10-07. Reviewed implementation baseline: merged main `96a6f09` containing [PR #444](https://github.com/geoffrey-xiao/deprail/pull/444). This record is preparatory: **formal H05-008 runtime execution is blocked by prerequisite acceptance**. No release approval, owner no-go decision, tag, publication or new runtime test result is asserted.

Preparation review: [PR #445](https://github.com/geoffrey-xiao/deprail/pull/445), linked with `Refs #427`; [issue handoff](https://github.com/geoffrey-xiao/deprail/issues/427#issuecomment-6028034470). The live issue Project is `Review` for preparation only, with a Blocked Reason naming the unmet formal execution gates. PR labels `area:test`, `risk:R3`, `priority:P0`, `type:test`, issue membership and Sprint 4 were verified after creation. This changes no acceptance decision.

Preparation CI [37551712270](https://github.com/geoffrey-xiao/deprail/actions/runs/37551712270) on `33f7944`: Ubuntu and macOS passed; Windows failed `go test ./...`, specifically `TestApplyEndToEndFailureBoundaries/cancellation` at `cmd/deprail/apply_e2e_test.go:215` with `unexpected end of JSON input`. Windows package build/native smoke did not run after that failure. Cause is not established by this log; no flaky-test classification, passing rerun, weaker assertion or runtime fix is claimed. This adds an unresolved verification blocker, distinct from the previously resolved #432 QA correction.

## Reconciled prerequisite state

| Contract | Observed tracking state | Evidence and remaining gate |
| --- | --- | --- |
| QA #419 | Closed / Project Done | Historical investigation is preserved. [#432 explicit acceptance](https://github.com/geoffrey-xiao/deprail/issues/432#issuecomment-5913275072), merged #433 and the backlog's accepted kickoff addendum resolve the QA correction gate. Do not continue describing the original OutputCap correction or timeout fixture as unimplemented; original failures are not erased. |
| H05-001 #420 | Closed / Done | Prior accepted projection slice; integrated candidate behavior still requires H05-008 proof. |
| H05-002 #421 | Open / Review | [Store evidence](H05-002-STORE-EVIDENCE.md): permission/dependency/advisory and physical-growth acceptance, plus manual restore rehearsal, are not recorded as complete. |
| H05-003 #422 | Closed / Done | Prior accepted service slice; integrated evidence still required. |
| H05-004 #423 | Open / Review | CLI capture acceptance has no completed issue disposition in the inspected live tracking. |
| H05-005 #424 | Open / Review | Transport acceptance has no completed issue disposition in the inspected live tracking. |
| H05-006 #425 | Open / Review | [Source UI evidence](H05-006-UI-EVIDENCE.md) and merged #443; actual supported browser/AT, physical zoom and supply-chain decisions remain unaccepted. |
| H05-007 #426 | Open / Review | Merged #444; exact-head CI [37549943755](https://github.com/geoffrey-xiao/deprail/actions/runs/37549943755) passed all three OS jobs at `5c78da1`. Automated review reported no major issues at that commit. [Packaging evidence](H05-007-PACKAGING-EVIDENCE.md) is component evidence, not release-binded native browser/platform or supply-chain approval. |
| H05-008 #427 | Open / Todo before readiness assessment | Requires completed H05-001–007 or explicit owner dispositions before formal execution. This guide/record does not satisfy its runtime acceptance criteria. |

A merged PR and successful CI do not by themselves close unchecked acceptance criteria. No native APPROVED review is invented from the owner merging #444. User requests to continue authorize preparation, not automatic acceptance of missing platform or supply-chain evidence.

## Required owner dispositions and remaining evidence

1. Record acceptance or explicit bounded disposition for #421 and #423–#426. Keep distinct technical and security decisions where required.
2. Resolve the disclosed Go dependency advisory, frontend `caniuse-lite` CC-BY-4.0 attribution/distribution implications and unverified publisher provenance. Existing registry integrity/audit results are not those decisions.
3. Supply actual authoritative release-binded candidate binaries/identity, or explicitly select a reviewed local candidate for prepublication verification. No v0.5 tag or suffix is chosen by this record.
4. Execute integrated JS/Python/Java scan/save/reopen/API/UI, security/read-only tree comparisons and offline recovery from disposable copies under the guide. Historical component tests do not prove this matrix.
5. Execute or explicitly disposition current stable Chrome/NVDA, Safari/VoiceOver and Firefox/Orca, and all four native binary targets. Cross-builds, Rosetta identity-only smoke and headless Chromium are not interchangeable with full native/browser/AT verification.
6. Record full artifact SBOM/signature/provenance status and separate owner technical, security/architecture and release decisions. Missing evidence remains a gate, not an implicitly approved preview gap.

## Release-checklist coverage

| General checklist section | Current disposition |
| --- | --- |
| Release identity | v0.5.0 preview-first plan exists; exact preview tag, release manifest and final artifact identity are not selected/verified here. |
| Plan and scope | Implementation merged through H05-007; outstanding issue acceptance/disposition above prevents a completed release-scope claim. |
| Contracts and security | Existing approved contracts and component evidence retained; integrated runtime/privacy/recovery and owner risk decisions pending. |
| Automated/manual verification | #444 three-OS component CI passed; no clean integrated release-candidate or representative workflow result claimed by this readiness record. |
| Artifacts/supply chain | H05-007 local hashes/budgets and host SBOM available; actual release artifacts and signature/provenance/full supply-chain dispositions unverified. |
| Publication/approval | Not authorized or executed. No protected publication gate or release workflow pass inferred from PR CI. |
| Post-release | Not applicable before publication; no retrospective or release closure claimed. |

Technical readiness assessment: **blocked / not release-ready**. This is not an owner-issued release no-go. Preserve existing user DB/WAL/SHM/artifacts and immutable tags; do not downgrade, reset or overwrite to obtain a passing result.
