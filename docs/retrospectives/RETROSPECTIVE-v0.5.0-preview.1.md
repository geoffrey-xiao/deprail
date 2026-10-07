# DepRail v0.5.0-preview.1 Retrospective

**Tracking issue:** [#462](https://github.com/geoffrey-xiao/deprail/issues/462)  
**Release evidence:** [RELEASE-v0.5.0-preview.1-EVIDENCE.md](../release-evidence/RELEASE-v0.5.0-preview.1-EVIDENCE.md)  
**Release issue / delivery epic:** [#427](https://github.com/geoffrey-xiao/deprail/issues/427) / [#418](https://github.com/geoffrey-xiao/deprail/issues/418)  
**Release:** [v0.5.0-preview.1](https://github.com/geoffrey-xiao/deprail/releases/tag/v0.5.0-preview.1)  
**Workflow:** [Release Combined37610017998](https://github.com/geoffrey-xiao/deprail/actions/runs/37610017998)  
**Status:** Published preview; factual retrospective prepared for owner acceptance and separate security/architecture closeout.  
**Release mode:** Preview only; no stable-ready claim.  
**Owner / required reviewer:** `@geoffrey-xiao`; external review optional under [ADR-0005](../adr/ADR-0005-solo-owner-review-policy.md). Earlier retrospectives' independent-review requirements are historical, not current blockers.

## Executive summary

`v0.5.0-preview.1` delivers the local history/read-only console layer: explicit scan-history capture, independently versioned safe history projection, private SQLite storage, shared history queries, bounded authenticated loopback API and an embedded React console. Browser queries never start scans, remediate, mutate source, export/delete history or publish remotely. Ordinary scans remain independent of optional history/UI; operation outcome and report completeness stay distinct.

Published `2026-10-07T10:54:16Z` from owner-merged source **`48ada15e16a049af54f8b4f665f22892cc81d2b2`**. Actual authoritative prepare/verify/four artifact builds/finalize/publish succeeded after an initial strict-input failure. Owner `geoffrey-xiao` approved the protected release-approval gate. All four public binaries and supplied SPDX independently pass the published checksum manifest; native Darwin-arm64 published binary reports matching version/source and discovers three complete npm/Python/Maven workspaces.

Publication is proven; the full release-binary scan→save→reopen→API/UI, native/browser/AT/recovery matrix and original criterion-level gap decisions are **not** proven by that workflow. This preview does not inherit v0.4 risk acceptance or waive the current gaps. Owner publication approval and retrospective/security acceptance are separate records.

## What we delivered

| Capability | Implementation / issue | Boundary |
| --- | --- | --- |
| Safe history-v1 projection | [#434](https://github.com/geoffrey-xiao/deprail/pull/434) / #420 | Allowlisted same-operation graph/report/typed diagnostics; no raw host root/errors or invented provenance. |
| Private SQLite history | [#436](https://github.com/geoffrey-xiao/deprail/pull/436) / #421 | Restrictive permissions, atomic entry/reference commit, explicit corrupt/future/lock/quota failure, no automatic eviction/repair/downgrade. 256MiB is logical payload, not physical disk cap. |
| Shared capture/read services | [#438](https://github.com/geoffrey-xiao/deprail/pull/438) / #422 | Distinct operation occurrences, stable ordering/cursors and unavailable evidence semantics. |
| Explicit CLI history capture | [#439](https://github.com/geoffrey-xiao/deprail/pull/439) / #423 | `scan --save-history`; existing scan JSON/output preserved; exit6 only for save failure after otherwise-successful scan. |
| Read-only foreground local console | [#440](https://github.com/geoffrey-xiao/deprail/pull/440), [#442](https://github.com/geoffrey-xiao/deprail/pull/442) / #424/#441 | `web [--open]`; ephemeral loopback, bearer/origin/host guards, bounded whole-record responses, no query-triggered scan/write. |
| Embedded history/detail/about UI | [#443](https://github.com/geoffrey-xiao/deprail/pull/443), #447/#449/#451 / #425 | Explicit outcome/report/empty/unavailable/error states; semantic controls, keyboard/reflow evidence. Actual supported AT/physical zoom remains separate. |
| Four-target reproducible packaging | [#444](https://github.com/geoffrey-xiao/deprail/pull/444), [#458](https://github.com/geoffrey-xiao/deprail/pull/458) / #426/#457 | Pinned frontend assets, matching API/build identity, one frozen publication SHA, complete private artifacts before approval, draft/upload/public boundary. |
| Quality/artifact/dependency corrections | #431/#433/#455/#456/[#461](https://github.com/geoffrey-xiao/deprail/pull/461) | Bounded subprocess fixes, readiness-synchronized timed tests, safe optional verifier and minimum fixed x/text0.39.0. |

Excluded/deferred: hosted/team services, PostgreSQL/accounts/RBAC, browser-triggered scan or remediation, remote publishing/source upload, automatic scanner install, history deletion/export/cleanup/repair/downgrade, and added full Python/Java remediation. The latter stays outside v0.5 under [release plan §14](../03-planning/deprail-development-plan-v0.5.0.md#14-detailed-contract-reconciliation-for-current-owner-review); no scope reduction or new roadmap stage is introduced here.

## What went well

- Planning chose explicit data/privacy/failure contracts and solo-owner review before runtime delivery; separate history-v1 avoided serializing the raw CLI report into a local API.
- Real component/packaged browser evidence distinguished failed/cancelled operations from report completeness, unavailable data from empty success, and memory-only authentication from browser persistence. See [UI record](../04-execution/deprail-v0.5.0-execution-package/tracking/H05-006-UI-EVIDENCE.md).
- User-driven UI iteration narrowed control/layout/link affordances instead of adding browser mutation features; owner local testing feedback stayed bounded rather than being presented as AT certification.
- Timed-process failures, nil-interface behavior and wrong-base delivery were corrected without weakening consumer assertions or pretending isolated-branch CI proved main.
- Dependency remediation used a real failing-before/passing-after malformed-input iterator, minimum fixed version, checksum/graph checks and actual CLI/console. Other selected modules and Go1.27.1 stayed unchanged; the SDK-vendor copy was assessed separately, not hidden by a module upgrade.
- The repaired **actual** Release Combined completed one-source verification/build/SPDX/checksum assembly before owner approval and publication; all six downloaded API digests/sizes matched, and the manifest verified all five listed binary/SBOM artifacts.
- Actual published native identity is clean (`vcs.modified=false`), matching source48ada15, Go1.27.1, x/text0.39.0 and preview tag. Post-publication representative discovery succeeds without network/scanner installation.

## Gaps and mistakes

| Gap / mistake | Impact | Established root cause / limit | Correction / current disposition | Prevention |
| --- | --- | --- | --- | --- |
| Timed apply QA sometimes failed JSON decode on Windows | Native validation stopped before later smoke/build | Test startup/cancellation used fixed delays without observing actual mutation readiness; original failed phase was not recorded | #454/#455 synchronize on real child ready marker; original outcome/exit/cleanup assertions retained; corrected exact-head CI passed | Include real phase barrier and stderr/evidence in timed-child scenarios; don't label unexplained failures flaky or weaken assertions. |
| Optional artifact resolver held a typed nil in an interface | Default console could call an absent verifier instead of reporting unavailable | Nonnil interface containing nil *artifact.Store | #452/#456 deliver a truly absent verifier; real default/explicit-root artifact-state scenarios recorded | Exercise default optional-interface behavior against actual service/store/listener, not only explicit-root happy path. |
| Artifact fix first merged into a temporary base, not main | Passing isolated branch did not supply release prerequisites | #453 targeted the temporary branch | #456 clean main-targeted cutover after synchronized baseline, own CI/owner merge | Check PR base/delivery SHA and build release from main containing every prerequisite. |
| Original authoritative publishing path lacked coordinated frontend/source/artifact preparation | Component CI/build claims did not prove actual publication readiness | Missing pinned frontend build in release path, moving-main source exposure and publication before complete prepared asset set | #457/#458 freeze source, build embedded UI, assemble all four binaries/SBOM/checksums privately, then protected draft/upload/public; actual37610017998 now passed | Verify the owner-selected workflow end-to-end and inspect prepared assets at approval; don't substitute a different workflow/offline recorder. |
| Initial x/text assessment stopped too early or conflated module and SDK symbols | Either false whole-binary-safe claim or unnecessary toolchain upgrade was possible | Canonical module language graph differs from SDK-vendored norm; broad Form/private retained symbols do not prove the Iter trigger | #459/#460/#461 establish source/binary-bounded SDK applicability and upgrade selected module0.14→0.39; no patched-SDK/suppression claim | Separate versions, namespaces, actual trigger/callers, source and binary identity; preserve uncertainty and dated authoritative citations. |
| First publication dispatch had a leading space in version | prepare rejected input/exit2; all later jobs skipped | Logged input `" 0.5.0-preview.1"`, first codepoint32; entry mechanism not established | New correctly entered dispatch37610017998 succeeded; strict validation remained | Supply copy-safe exact input values and check field edges before dispatch; Re-run jobs retains invalid inputs. |
| Original criterion/PG dispositions not all durably recorded before publication | Protected approval proves publishing permission, not which verification gaps/risks were accepted | Prior preparation retained 27 original criteria and12 pending PG decisions; original #421/#423–427 still open at assessment | Record factual post-release evidence and retain missing technical/security dispositions; no retroactive pass/blanket waiver | Save exact criterion decisions and named gaps in the release record before protected publication approval. |
| Missing-scanner explicit capture reported HISTORY_WRITE_FAILED | Requested failed-operation history was not saved in the observed private-profile smoke | Same status/exit/diagnostic with old and fixed x/text; underlying cause not established | Keep prepublication observation under #423/#427 PG-05; no dependency-induced regression, successful capture or fix claimed | Exercise real absent-scanner and failed/cancelled history paths in the full released-binary matrix before accepting capture. |
| Native/browser/AT/recovery/physical-growth matrix remains incomplete | Artifact builds and component browser tests do not prove all supported user/security/recovery behavior | Existing evidence is scoped host/component/CI; required final matrix/rehearsals not supplied by publication | Existing issue/PG follow-ups below; no accepted deferral invented | Separate compile, CI native CLI, downloaded native workflow, browser/AT/zoom, logical quota and real ENOSPC evidence. |
| Supplied SBOM is not a completed license/provenance review | Inventory counts/fixture components can be mistaken for shipped dependency exposure; authenticity/compliance unproven | SPDX has119 records incl. source/per-target duplicates/test fixtures;110 licenseDeclared NOASSERTION; no signature/provenance assets | Record actual inventory/known runtime entries and unavailable/unverified statuses; existing PG-02/10/11/12 review remains | Classify runtime/build/test entries and SDK scope; review notices/licenses separately; checksum and SBOM supply are not attestations. |

## Security and supply-chain lessons

1. Private storage/read-only routes do not remove authentication, same-origin, size/deadline, path, permissions or corrupt-data boundaries. Transient OS/browser bootstrap exposure remains an accepted **design** tradeoff under the approved contracts, not an unrestricted runtime exemption.
2. Logical admission quotas do not cap SQLite/WAL/backups/artifacts or prove disk-full recovery. Never reset/repair/downgrade the only data copy to manufacture a pass.
3. Optional interfaces need genuine absent-capability semantics; unknown integrity must remain unavailable, never verified/clean.
4. A fixed selected module does not replace Go SDK vendor code. Conversely an old package version/private symbol is not by itself a reachable exploit; use precise trigger/source/binary evidence and owner security review, not generic suppression.
5. The published manifest independently verifies four binaries/SPDX; GitHub asset digests also match. Neither mechanism is a publisher signature or provenance verification.
6. SPDX contains React19.1.1/react-dom19.1.1/scheduler0.26.0 and selected x/text0.39.0. It also includes fixture lodash4.17.20 and duplicated records. Separate shipped exposure from test fixtures; don't equate119 records to119 unique runtime dependencies or complete license approval.
7. Actual protected owner publication approval is recorded. It does not replace per-criterion acceptance, named preview-gap dispositions, separate security assessment or stable go/no-go.

## Process improvements

### Applied in this release

- Reconciled contracts/ordered delivery, pinned script-free dependency installation and explicit capture/read-only boundaries.
- Real process-readiness synchronization, main-targeted cutover and one-source/prepared-asset publication.
- Source/SDK-aware vulnerability remediation and transparent negative evidence rather than broad risk acceptance.
- Version-specific retrospective/publication evidence with independent downloaded bytes and actual clean binary identity.

### Required follow-through

Use the table below as the accountable action set; these are not silently completed tasks or accepted waivers. Native/AT/recovery/signature/provenance deadlines come from the existing release plan; core failure-history and acceptance-record work should precede the next preview. No implementation or subsequent release is started by this document.

## Follow-up actions

| Action / tracking | Owner | Target gate | Observable acceptance | Status |
| --- | --- | --- | --- | --- |
| Original27 criterion and PG-01–12 decisions; [#427](https://github.com/geoffrey-xiao/deprail/issues/427), #421/#423–426 | @geoffrey-xiao | Post-release closeout; before next preview admission | Each criterion links actual proof or explicit bounded disposition; technical/security/publication decisions distinct, scope/owner/deadline recorded. | Open; publication alone did not check them. |
| Investigate missing-scanner save and full scan/save/reopen/API/UI failure matrix; [#423](https://github.com/geoffrey-xiao/deprail/issues/423)/#427 PG-05 | @geoffrey-xiao | Before next preview | Identify real cause; released-binary complete/partial/failed/cancelled/unavailable cases with truthful capture/status/error, source/DB/privacy comparisons; no fabricated findings. | Open; observed prepublication failure retained. |
| Downloaded native user/security workflow on four targets; [#426](https://github.com/geoffrey-xiao/deprail/issues/426)/#427 PG-06 | @geoffrey-xiao | Before stable v0.5; any earlier preview gap must be explicit per target | Actual binary identity/history/API/UI, permissions/launcher/shutdown and tree evidence on Linuxamd64/macOSamd64+arm64/Windowsamd64. | Open; cross-builds/checksums/native identity-only are insufficient. |
| Supported browser/AT/physical200% zoom; [#425](https://github.com/geoffrey-xiao/deprail/issues/425)/#427 PG-07 | @geoffrey-xiao | Before stable v0.5; named preview decisions separately | Chrome/NVDA, Safari/VoiceOver, Firefox/Orca versions and actual keyboard/focus/status/zoom scenarios. | Open; managed Chromium evidence not full AT proof. |
| Disposable backup/restore and physical-growth/ENOSPC; [#421](https://github.com/geoffrey-xiao/deprail/issues/421)/#427 PG-08/09 | @geoffrey-xiao | Before stable v0.5; preview disposition separately | Preserve stopped DB/WAL/SHM/artifacts, validate backup/restored copy and failure invariants; measure physical versus logical growth without data loss. | Open; page-limit tests are not real ENOSPC/restore rehearsal. |
| Distribution notices, SPDX classifications/licenses and publisher attestations; #425/#426/[#427](https://github.com/geoffrey-xiao/deprail/issues/427) PG-02/10/12 | @geoffrey-xiao | Publication evidence closeout; reassess before next preview/stable | Final format/runtime/build/test inventory, license/attribution obligations and per-provider verification/unavailable decisions; resolve NOASSERTION without guessing. | Open; actual SPDX supplied/inspected, complete review unproven. |
| Signatures/provenance; [#427](https://github.com/geoffrey-xiao/deprail/issues/427) PG-11 | @geoffrey-xiao | Before stable v0.5 or explicit versioned disposition | Supply/independently verify available signatures/provenance or record exact unavailable/owner-dispositioned status. | No such public assets; no blanket accepted waiver recorded. |
| Dedicated install/use/preview-limit release notes; [#427](https://github.com/geoffrey-xiao/deprail/issues/427) | @geoffrey-xiao | Post-release closeout; before next publication | Exact official installation/commands, included/excluded scope, named gaps and feedback/evidence path; preserve generated change links. | Open; current notes are changelog plus #427 link. |
| Dependency repair criterion/security closure; [#460](https://github.com/geoffrey-xiao/deprail/issues/460) | @geoffrey-xiao | Post-release closeout | Inspect fixed-pin/SDK assessment, exact-head CI and original criteria; close/Done only with owner review/evidence. | #461 merged and fixed module published; issue still open at inspection. |
| Retrospective owner acceptance/security closeout; [#462](https://github.com/geoffrey-xiao/deprail/issues/462) | @geoffrey-xiao | Post-release closeout | Review every retrospective criterion and follow-up; separate owner acceptance/security records and reviewed PR/CI/merge. | Pending; this document is not owner acceptance. |

## Remaining risks and rollback

- Full released-binary history/security/platform/browser/AT/recovery evidence and precise original gap decisions remain incomplete; missing-scanner capture observation is unresolved.
- SDK applicability is a specific documented unused-trigger assessment, not a whole-binary vulnerability clearance. Supplied SPDX is not complete license/attestation review; no signature/provenance assets were supplied.
- Preview publication does not prove stable readiness or justify automatically closing original issues, marking Master Checklist or launching v0.6.
- Rollback owner @geoffrey-xiao preserves immutable v0.5.0-preview.1/source/assets/evidence and all user DB/WAL/SHM/artifacts. Mark unsuitable preview if needed and publish a separately reviewed new version; never move/reuse tag, downgrade/reset storage or overwrite the only recovery copy. Current runbook/rehearsal remains #421/#427.

## Acceptance

- [ ] Owner reviewed every retrospective criterion and follow-up gate.
- [ ] Owner separately recorded security/architecture findings and precise remaining-risk decisions; external review optional.
- [x] Actual published release/source/workflow and protected owner publication approval are linked.
- [x] All six public asset digests/sizes independently matched; checksum manifest verified four binaries and supplied SPDX.
- [x] Actual native published-binary identity/discovery and SBOM observations/limits recorded in release evidence.
- [x] Release issue/evidence/retrospective/tracking issue cross-linked; every follow-up has issue, owner and acceptance gate.
- [ ] Original release/PG dispositions and retrospective accepted; no completion/closure inferred from publication or PR creation.
