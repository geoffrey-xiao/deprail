# DepRail v0.5.0-preview.1 Publication Evidence

**Status:** Published preview; post-release criterion/security closeout remains pending.  
**Release:** [v0.5.0-preview.1](https://github.com/geoffrey-xiao/deprail/releases/tag/v0.5.0-preview.1), public prerelease, published `2026-10-07T10:54:16Z`.  
**Release issue:** [#427](https://github.com/geoffrey-xiao/deprail/issues/427); delivery epic [#418](https://github.com/geoffrey-xiao/deprail/issues/418).  
**Retrospective issue:** [#462](https://github.com/geoffrey-xiao/deprail/issues/462); [retrospective](../retrospectives/RETROSPECTIVE-v0.5.0-preview.1.md).  
**Owner / required reviewer:** `@geoffrey-xiao`; external review optional under [ADR-0005](../adr/ADR-0005-solo-owner-review-policy.md).  
**Controls:** [release checklist](../RELEASE-CHECKLIST.md), [release plan §§14–15](../03-planning/deprail-development-plan-v0.5.0.md#14-detailed-contract-reconciliation-for-current-owner-review), [manual guide](../04-execution/deprail-v0.5.0-execution-package/MANUAL-TEST-GUIDE.md).

This record distinguishes actual publication and observed post-publication smoke from unfinished integrated/native/browser/recovery and owner risk-disposition gates. It does not retrospectively mark missing evidence passed or assert stable readiness.

## Immutable identity and publication

- Source/tag target: **`48ada15e16a049af54f8b4f665f22892cc81d2b2`**, owner-merged dependency repair [#461](https://github.com/geoffrey-xiao/deprail/pull/461), atop owner-merged QA/artifact/publisher/evidence corrections.
- Actual [Release Combined run37610017998](https://github.com/geoffrey-xiao/deprail/actions/runs/37610017998): completed/success. `prepare`, `verify`, all four `build-artifacts` targets, `finalize` and `publish` succeeded. Prepare completed10:50:24Z; verify10:51:56Z; all binaries10:53:00Z; finalize10:53:22Z; publish started10:54:06Z and ended10:54:18Z.
- GitHub run approval history reports `geoffrey-xiao`, state `approved`, environment `release-approval`, empty comment. This is protected publication approval, **not** a recorded per-PG technical/security disposition.
- Successful inputs followed the selected preview: version0.5.0-preview.1, evidence_issue427, prerelease=true, confirm=RELEASE. Actual public tag/mode/source and binary identity agree.
- Prior [failed run37609702708](https://github.com/geoffrey-xiao/deprail/actions/runs/37609702708) at the same source is preserved: submitted version was `" 0.5.0-preview.1"` (leading ASCII32); strict prepare validation emitted Invalid release version/exit2. Verify/build/finalize/publish skipped. Corrected new dispatch succeeded; no validation weakening or failed-job retry substituted.
- [Comparison](https://github.com/geoffrey-xiao/deprail/compare/v0.4.0-preview.2...v0.5.0-preview.1). Release notes contain generated PR changelog and #427 link, not a dedicated installation/usage/preview-gap summary; explicit release-note closeout remains #427.

## Independently downloaded public assets

On Darwin23.6.0 arm64, downloaded all six official assets to private disposable `/tmp/deprail-v05-published-KZONJS` using `gh release download v0.5.0-preview.1 --repo geoffrey-xiao/deprail --dir /tmp/deprail-v05-published-KZONJS --pattern '*'`. Exit0. Independently computed SHA-256 and byte counts matched every release API digest/size. `shasum -a 256 -c SHA256SUMS` in that directory exited0: four binaries and SPDX **all OK**. The manifest itself separately matched its release API digest. Checksums prove byte consistency, not cryptographic publisher identity.

| Asset | Bytes | Independently observed SHA-256 |
| --- | ---: | --- |
| deprail-linux-amd64 | 13000967 | 123245b445b3d01281642db5f0ca93c32583dc095d9c0727b1553a5cdbb6cbd5 |
| deprail-darwin-amd64 | 13123904 | ce02fea7f746362778d3b30bd0ebd3af140ec4af755a0171bd2aeb5580c3895a |
| deprail-darwin-arm64 | 12445378 | 896853a0546f44e0aec3c5708042f5b26562523f84deb67724a9a31911cee295 |
| deprail-windows-amd64.exe | 13207040 | cc7a5fac973916084eb196f433bed7f2f27cb520e46c9bfb8171720d31c01674 |
| deprail-0.5.0-preview.1.spdx.json | 202901 | 249c3117fe01069022edb6308607889419261f973114c94fde7f430672b8f851 |
| SHA256SUMS | 452 | ebeebacaf356d1f99943bb867933285de1ade76a58608e4c4af1ea4205515922 |

Only the six intended uploaded resources are listed: four binaries, SPDX, checksum manifest. No raw scan reports/logs, signatures or provenance files are published as release assets. GitHub's automatically supplied source archives are distinct from this uploaded inventory. No independent GitHub attestation verification is claimed.

## Actual published-binary smoke

After checksum verification, chmod only the disposable native download and run:

| Command | Observed result |
| --- | --- |
| `/tmp/deprail-v05-published-KZONJS/deprail-darwin-arm64 --version` | Exit0: `deprail v0.5.0-preview.1 (tag=v0.5.0-preview.1 commit=48ada15e16a049af54f8b4f665f22892cc81d2b2)`. |
| `go version -m /tmp/deprail-v05-published-KZONJS/deprail-darwin-arm64` | Exit0: Go1.27.1, CGO_ENABLED=0, Darwin/arm64, trimpath, vcs.revision48ada15…, vcs.modified=false; selected x/text v0.39.0 with fixed checksum. Go module pseudoversion is build metadata, not a conflicting CLI tag. |
| `/tmp/deprail-v05-published-KZONJS/deprail-darwin-arm64 discover testdata/fixtures/mixed-repository --format json` | Exit0: overall complete; three complete workspaces frontend/npm, services/api/python, services/worker/maven; no diagnostics. |

The four fixture-file hashes match the prior [repair smoke record](../04-execution/deprail-v0.5.0-execution-package/tracking/H05-008-FU1-XTEXT-EVIDENCE.md#existing-verification-and-actual-consumer-smoke). This is bounded source-file consistency, not a full repository-tree comparison. Discovery is not a vulnerability scan or scan→save→reopen→API/UI validation. No findings fabricated, scanner installed, original history DB opened/migrated, user data changed or pipeline rerun for this retrospective.

## Supplied SBOM and supply-chain limits

Downloaded JSON parses as SPDX2.3; creation2026-10-07T10:53:18Z, Syft1.42.3 / Anchore, license-list3.28. It contains **119 package records / 333 relationships**, not119 unique shipped dependencies. PURL record categories: Golang82, GitHub20, npm8, Maven4, PyPI3, no-PURL2. The source/build inventory includes intentionally vulnerable test-fixture components such as lodash4.17.20, so presence is not proof they are linked into the CLI.

Observed runtime entries include React19.1.1/react-dom19.1.1/scheduler0.26.0. Five x/text records all say v0.39.0 (source/per-binary inventory); their declared license is NOASSERTION. Overall110 records have licenseDeclared NOASSERTION and92 licenseConcluded NOASSERTION. Separate source review found x/text's existing BSD terms, but SBOM supply does not establish complete license/notice approval, all toolchain-vendor visibility, final inventory classification or publisher provenance. These remain #427 PG-02/10/12 follow-ups, not an invented passing legal review.

[Dependency repair #460](https://github.com/geoffrey-xiao/deprail/issues/460) and [dated source/SDK applicability proof](../04-execution/deprail-v0.5.0-execution-package/tracking/H05-008-FU1-XTEXT-EVIDENCE.md) remain relevant: selected module is fixed; unchanged Go SDK vendor scope needs the specific unused-iterator assessment, not a broad patched-SDK/zero-vulnerability claim. Signatures/provenance are **not supplied as release assets**; no complete verification or explicit per-gap acceptance is inferred.

## Release-checklist disposition at post-publication assessment

| Checklist area | Observed evidence | Uncompleted / tracked gate |
| --- | --- | --- |
| Identity | Public preview tag, source48ada15, actual clean binary version/build metadata agree. | No stable/RC promotion; future tag immutable/new-version decision. |
| Plan/scope | Accepted v0.5 local history/read-only console; merged implementation and corrective PRs linked in retrospective. | Original #421/#423–427 and #460 remain open at inspection; every unchecked original criterion needs owner evidence/disposition. |
| Contracts/security | Existing approved contracts and bounded component/security evidence retained; no contract changes in retrospective. | Distinct original technical/security and exact preview-gap dispositions not established by protected publication approval. |
| Automated/manual checks | Actual authoritative run verify/build/finalize succeeded; prior final-head #461 CI37608443541 passed three OS/native smoke/four cross-builds. Native published identity/discovery observed above. | Full actual release-binary multi-language scan/history/API/UI/hostile state, platform/browser/AT/zoom/recovery/ENOSPC matrix not proven here. |
| Assets/supply chain | All six official downloaded digests/sizes agree; manifest verifies five listed artifacts; SPDX supplied and inspected. | Inventory/license classification/notice/publisher provenance remain incomplete; signature/provenance assets unavailable, no blanket accepted deferral. |
| Publication | Public prerelease, successful authoritative run and owner protected approval observed; complete assets prepared before publish. | Exact criterion/gap decisions and dedicated user-facing notes still need durable closeout. |
| Post-release | This version-specific evidence and retrospective cross-link #427/#462. Follow-ups have existing issue/PG references and owner/gates. | Owner retrospective acceptance and security/architecture review remain separate; no Master Checklist/issue closure/Project Done asserted. |

## Owner decisions, follow-ups and rollback

Owner chose preview publication in the conversation (“我觉得可以了 那我就去跑pipeline发布了”) and approved release-approval in the observed run; later reported release complete. Those are real publication decisions, not a recorded blanket PG-01–12 waiver, stable-ready claim or retrospective acceptance. [Retrospective follow-ups](../retrospectives/RETROSPECTIVE-v0.5.0-preview.1.md#follow-up-actions) assign every remaining action to existing issues and owner `@geoffrey-xiao`. Historical prepublication readiness/proposal is retained additively, not silently rewritten as a prepublication pass.

Rollback owner @geoffrey-xiao: preserve immutable v0.5.0-preview.1/tag/assets/evidence and original DB/WAL/SHM/artifacts. If unsuitable, mark preview accordingly and use a reviewed new version; never move/reuse tag, reset/downgrade user storage or overwrite the only recovery copy. Actual restore rehearsal remains #421/#427 PG-08.

- [x] Actual release/run/source/assets and independently exercised native identity/discovery recorded.
- [x] Retrospective and existing follow-up tracking linked.
- [ ] Owner reviewed all original and retrospective evidence/criteria and recorded bounded gap decisions.
- [ ] Owner separately recorded security/architecture/supply-chain closeout.
- [ ] Required release/retrospective closure accepted; no closure inferred from publication.

## Retrospective preparation verification

Checked **75 relative document file/heading links** across the six new/affected documents: zero missing targets. Actual published-asset/native smoke above is separate from document/PR checks. Disposable downloaded copies were removed after verification; official tag/assets and user history/source remain unchanged. No local test suite was rerun for documentation-only changes; final documentation PR's existing CI/review record is tracked in #462. Historical retrospectives and Master Checklist are untouched.
