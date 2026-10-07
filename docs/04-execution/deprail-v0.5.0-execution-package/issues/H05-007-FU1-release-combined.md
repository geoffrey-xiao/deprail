# H05-007-FU1: Make Release Combined Build and Publish One Complete v0.5 Source

GitHub issue: [#457](https://github.com/geoffrey-xiao/deprail/issues/457).

Review PR: [#458](https://github.com/geoffrey-xiao/deprail/pull/458), Refs #457; implementation dea8f1c. Exact-head CI and owner publishing/security review remain separate from local/offline evidence.

## Planning metadata

- Parent: H05-007 / #426; release evidence #427; delivery epic #418.
- Type: bug; area: foundation; priority: P0; risk: R3 (publishing).
- Target version/milestone: v0.5.0; Sprint 4.
- Owner/reviewer: @geoffrey-xiao; independent review optional under ADR-0005.
- Dependencies: owner-merged #455/#456, main c554441; existing pinned frontend/embedded asset contract from #444.
- Owner scope decision: conversation approval “可以 继续” after the four-defect inspection on #427. This authorizes the repair and review PR, not dispatch, tagging, publishing or autonomous merge.
- Blocked reason: none for this bounded repair; formal release acceptance remains separate.

## Definition of Ready

- [x] Value: the owner's authoritative Release Combined must build the embedded console and cannot publicly release before preparation succeeds.
- [x] Scope, exclusions, source/build inputs and failure boundaries are explicit below.
- [x] Runtime smoke and observable acceptance are explicit; owner/reviewer and dependencies recorded.
- [x] Maps to release plan §13, product open-source release controls, architecture §11, H05-007 and the release checklist's identity/artifact/publication gates.

## Goal and scope

Repair only `.github/workflows/release-combined.yml`: install Node22.23.3/npm10.9.9 and generate embedded assets in verification and artifact jobs; freeze prepare's main SHA and use it for every subsequent checkout, binary identity and tag; finish all four artifacts, checksums and existing SPDX SBOM before the protected approval/publication job. Preserve current dispatch inputs, four artifact names, version/tag conventions and release-approval environment.

Use the existing `make frontend` and CI toolchain setup patterns. Only publication needs contents:write. Create a draft release with the complete prepared assets, then make it public only after successful upload. No clobber or moving/reusing an existing immutable tag. A publication-time failure may leave a new tag or draft; retain it and require owner recovery/new-version disposition, never automatic deletion or retry.

## Out of scope

No signatures/provenance implementation, native matrix expansion, dependencies, schema/CLI/storage changes, new release version selection, H05-008 acceptance, dispatch or public release. Native/browser/AT and supply-chain gaps remain explicit under release-plan §14; prior v0.4 approval is not v0.5 approval. No displaced product work.

## Inputs, outputs and failure behavior

Inputs remain version, evidence_issue, prerelease, confirm=RELEASE. Prepare records one main commit before any verification. Outputs: four binaries carrying the selected version/tag/source SHA, SHA256SUMS and the existing version-specific SPDX file, assembled privately as a workflow artifact before approval. Failed source/frontend/test/platform/SBOM/checksum work prevents publication. Failed asset upload leaves the release nonpublic. All checkouts and tag identity use the prepared SHA even if main advances. Existing tags fail closed; no overwrite.

## Required verification and acceptance

- [ ] Fresh archived source plus pinned frontend build yields real embedded binaries; all four targets compile; native host identity and representative discovery work.
- [ ] Every source-dependent job and tag/binary identity uses one prepared SHA; independent main movement cannot change publication source.
- [ ] Failed verification/build/assembly prevents the publication job; asset upload failure cannot make the release public.
- [ ] Checksums verify all intended prepared binary/SBOM assets; no logs/credentials/private artifacts uploaded as release assets.
- [ ] Workflow syntax checked; actual shell build/publication paths exercised offline, with limitations distinguished from live GitHub execution.
- [ ] New PR exact-head existing three-OS CI passes; no release dispatch performed during repair.

## Human review and evidence

Owner reviews publishing permissions, job dependencies, immutable-source identity, draft/upload/public transition, failure recovery and unchanged exclusions. Technical and security/publishing acceptance remain distinct from release go/no-go. Record clean-source failure proof, exact commands/results, source/binary hashes, offline publication failure scenarios, syntax check, CI and reviewer decisions in the linked tracking evidence. Do not mark Master Checklist or release acceptance from this repair.

[Actual syntax/build/offline failure-boundary evidence](../tracking/H05-007-FU1-RELEASE-COMBINED-EVIDENCE.md) records the four clean builds, identity, canonical-path caveat and explicit limits of offline publication/SBOM fixtures.

## Final acceptance

- [ ] Owner criteria/security/residual risk reviewed.
- [ ] Required verification and exact-head CI accepted.
- [ ] Owner-reviewed main merge recorded before Project Done/closure.

## Rollback

Revert this workflow change without dispatching it. Preserve existing immutable tags/releases and all user data. A failed real publication requires an explicit owner recovery decision; never force-push/move a tag or clobber assets.
