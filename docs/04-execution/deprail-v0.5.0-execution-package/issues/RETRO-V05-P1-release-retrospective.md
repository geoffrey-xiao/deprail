# RETRO-V05-P1: Record v0.5.0-preview.1 Retrospective and Publication Evidence

GitHub issue: [#462](https://github.com/geoffrey-xiao/deprail/issues/462).

Documents: [retrospective](../../../retrospectives/RETROSPECTIVE-v0.5.0-preview.1.md) and [actual publication evidence](../../../release-evidence/RELEASE-v0.5.0-preview.1-EVIDENCE.md). Original release issue #427 retains unfinished criterion/security closeout.

## Planning metadata

- Type: docs; area: docs; priority: P1; risk: R1 (evidence accuracy; no runtime change).
- Target version/milestone: v0.5.0; published version v0.5.0-preview.1; Sprint 4.
- Owner/reviewer: @geoffrey-xiao; external review optional under ADR-0005.
- Parent/dependencies: #418/#427; owner-merged #461/main48ada15; actual published preview/run37610017998.
- Owner request: “ok 已经release 参考其他版本，开始写retrospective”. This authorizes version-specific documentation/evidence and review PR, not retrospective/security acceptance, issue closure, merge or another publication.
- Blocked reason: none for factual retrospective preparation; unfinished original release criteria stay explicitly open.

## Definition of Ready

- [x] Value: preserve what shipped, exact published identity, mistakes/corrections and actionable remaining risks without generalizing CI or preview publication into stable readiness.
- [x] Scope/exclusions, observable acceptance, failure boundaries and evidence scenarios are explicit.
- [x] Maps to product explainable verified changes, architecture secure local-first delivery, roadmap v0.5, release plan §13–15, #427 and RELEASE-CHECKLIST post-release closure.
- [x] Owner/reviewer/version/dependencies recorded; previous retrospectives inspected and ADR-0005 supersedes their old independent-review requirement.
- [x] Duplicate searches found no existing version-specific retrospective issue/local document; #427 owns integrated release evidence, not the separate version retrospective outcome.

## Goal and scope

Create `docs/retrospectives/RETROSPECTIVE-v0.5.0-preview.1.md` using existing version structure and `docs/release-evidence/RELEASE-v0.5.0-preview.1-EVIDENCE.md` as its factual publication record. Cross-link them with #427/current readiness and this issue. Record delivered history/projection/store/API/embedded UI scope, actual successful workflow/source/tag/assets, independent checksums and real native published-binary smoke, previous defects/corrections, supply-chain lessons, named existing follow-up issues/owners/gates, separate owner technical/security review, and immutable-tag/data-preserving rollback.

## Out of scope

No code/config/schema/storage/CLI/dependency/workflow changes, scan findings fabrication, full release-matrix claim, invented approval/waiver, new roadmap stage, automatic closure/merge/publishing or changes to accepted historical retrospectives/master checklist. No duplicate follow-up issues when #421/#423–427 already cover the unmet outcome. No displaced product work.

## Inputs, outputs and failure behavior

Inputs: actual public release, successful37610017998 and failed37609702708, immutable source48ada15, implementation/CI/local evidence and current tracking. Outputs: factual retrospective/publication evidence and cross-linked review PR. Unknown/unfinished criteria remain gaps; unavailable signatures/provenance are not verified; protected publication approval is not per-criterion security acceptance. Artifact identity/hash failure prevents execution or a passing claim; only disposable downloaded artifacts/profile are used, no user DB/source/tag mutation.

## Required verification

Independently download all six release assets; match GitHub byte counts/digests and published SHA256SUMS. Run actual Darwin-arm64 published binary --version and representative mixed-repository JSON discovery; inspect Go build provenance and supplied SPDX. Verify new/affected local document file/heading links and final documentation PR existing CI. No new permanent tests or rerunning the failed release pipeline to confirm known failure.

## Acceptance criteria

- [ ] Actual publication/version/source/run/owner approval/asset identity are reconciled, exact commands/results recorded and limitations explicit.
- [ ] Retrospective covers delivery, what went well, mistakes with impact/root cause/correction/prevention, security/supply-chain lessons, process improvements, remaining risks and rollback.
- [ ] Every follow-up has an existing issue/PG reference, accountable owner, target gate and observable acceptance; no untracked action or invented accepted deferral.
- [ ] Release evidence, retrospective, tracking issue and #427 are cross-linked; previous versions/master checklist unchanged.
- [ ] Actual published-binary smoke, relative document-link checks and final PR CI are recorded; owner acceptance/security review remain distinct and pending until performed.

## Human review and evidence

Owner reviews factual identity/scope, known failure interpretation, risk/disposition wording and follow-up gates. Repository documents remain the durable record; issue/PR retain CI and review links. Existing #460 is owner-merged with CI/evidence but still open: owner must check its acceptance/security evidence before closure; this retrospective must not silently close it or original #421/#423–427.

## Final acceptance

- [ ] Owner reviewed all retrospective criteria and accepted remaining risk.
- [ ] Owner separately recorded security/architecture assessment; external review optional.
- [ ] Verification/follow-up tracking reviewed and owner-reviewed merge linked before Done/closure.

## Rollback

Revert only this documentation change. Preserve published immutable tag/assets and user source/DB/WAL/SHM/artifacts; any product correction uses a reviewed new preview, never a moved tag or forced storage downgrade/reset.
