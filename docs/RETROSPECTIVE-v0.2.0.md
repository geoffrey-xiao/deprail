# DepRail v0.2.0 Release Retrospective

**Tracking issue:** [#194](https://github.com/geoffrey-xiao/deprail/issues/194)
**Release umbrella:** [#148](https://github.com/geoffrey-xiao/deprail/issues/148)
**Status:** Draft for owner and reviewer review
**Release mode:** Preview cycle; stable `v0.2.0` not yet published
**Owner:** `@geoffrey-xiao`
**Reviewer:** `@geoffreyxiaoai`

## Executive summary

The v0.2 cycle delivered the core dependency-security guardrail path: deterministic baseline comparison, workspace and alias identity, `diff`, typed policy evaluation, expiring exceptions, policy exit mapping, SARIF output, GitHub Action packaging, release artifacts, checksums, SBOM generation, keyless signing, provenance attestations, protected manual artifact publication, and user installation guidance.

The cycle also exposed release-process weaknesses. Preview releases were created before all identity and evidence gates were reconciled. One release workflow failed because the upload job lacked an explicit repository; one preview failed because the SBOM action attempted an early release upload with read-only contents permission; one preview carried a stale manifest identity; the manual environment initially had no required reviewer and therefore did not pause; and a merged branch was reused for a follow-up fix, causing avoidable branch/PR conflict.

The security-control implementation is now exercised successfully in `v0.2.0-preview.4`: SBOM generation, Cosign signing and verification, provenance generation and verification, checksums, and manual approval before asset publication all passed. Preview.4 still has a manifest/tag identity mismatch and is not stable-ready.

## What we delivered

### Product and contract capabilities

- Deterministic baseline model and tamper-safe local storage.
- Baseline comparison with changed/upgraded matching and provenance.
- Workspace and alias identity with version-independent PURL matching.
- `deprail diff` with versioned JSON and terminal output.
- Typed policy evaluation, severity/completeness gates, and exit mapping.
- Expiring exceptions and policy integration.
- SARIF output with valid levels, stable ordering, and provenance.
- Discovery and representative JavaScript, Python, Java, and mixed-repository scans.
- Raw artifact retention and deterministic normalization.

### Release and distribution capabilities

- GitHub Action packaging and security documentation.
- Cross-platform artifact matrix: Linux amd64, macOS amd64, macOS arm64, Windows amd64.
- SHA-256 checksum manifests and verification.
- SPDX SBOM generation.
- Keyless Cosign signatures and certificate publication.
- GitHub build-provenance attestations.
- Protected `release-approval` environment before release asset upload.
- Root `README.md` installation and usage guidance.
- Version-specific preview evidence records.

## What went well

- CI remained green across Linux, macOS, and Windows for the release-related PRs.
- The release upload failure was reproduced and corrected with an explicit `--repo` argument.
- Checksums were independently downloaded and verified for preview releases.
- The workflow failed closed when SBOM publication permissions were wrong; it did not silently publish an incomplete security bundle.
- The manual environment gate was eventually proven: preview.4 remained in `waiting` state until approval, then completed publication.
- Historical preview records were separated instead of mixing artifact checksums across tags.
- User-reported macOS smoke testing was retained as owner evidence rather than fabricated CI output.
- The release process explicitly distinguishes preview gaps from stable approval.

## Faults and mistakes

| Fault | Impact | Root cause | Correction | Prevention |
| --- | --- | --- | --- | --- |
| Release upload failed in the manual verification job | Initial `v0.2.0` workflow could not upload assets | `gh release upload` ran without checkout/repository context | Added explicit `--repo "${{ github.repository }}"` | Add a release-workflow smoke test that verifies upload context before publication |
| Preview.1 evidence and later preview records were mixed | Audit ambiguity and incorrect attribution risk | One evidence file was reused across immutable releases | Added version-specific evidence files | Require one evidence file per tag in the release checklist |
| Preview.3 SBOM step failed | Release stopped before approval | SBOM action attempted automatic release upload from a read-only build job | Disabled action-managed uploads; gated all release assets through the approval job | Review third-party action side effects and permissions before tagging |
| Preview.4 manifest remained `0.2.0-preview.3` | Artifact/tag/manifest identity mismatch; stable release blocked | Manifest was changed for the next preview but the tag was created before reconciliation | Recorded explicit no-go evidence; next corrected preview must align identity | Add a pre-tag assertion comparing manifest, tag plan, and intended tag |
| Manual approval environment initially had no required reviewer | Preview.2 publication proceeded automatically | Environment existed without protection rules | Added required reviewer configuration and verified preview.4 waiting state | Add environment protection verification to the release checklist and CI preflight |
| PR branch reused after merge | Follow-up SBOM fix conflicted with merged history | A new fix was pushed to an already-merged preparation branch | Closed superseded PR and created a clean branch from `origin/main` | One issue/PR branch per change; branch from current main before every fix |
| Issue #148 closed too early | Umbrella tracking lost later release work | First evidence PR used `Closes #148` | Reopened issue and converted references to `Refs #148` | Umbrella issues must use `Refs` until owner acceptance is complete |
| User checksum guidance failed for single downloads | Users saw missing-file failures | `shasum -c` checked all platform entries | Added `--ignore-missing` and Windows verification | Test documentation from a clean download directory on each platform |
| Windows guidance lacked verification and used PowerShell syntax for persistent PATH | Users could execute unverified binaries or get a broken persistent PATH | Installation examples were reviewed only from macOS/Linux perspective | Added PowerShell hash verification and `%USERPROFILE%\bin` guidance | Require platform-owner walkthroughs for installation docs |

## Security and supply-chain lessons

1. **Permissions must follow publication boundaries.** Build jobs need repository read, OIDC, and attestation permissions. Only the approved upload job needs `contents: write`.
2. **Third-party actions can perform hidden publication.** `anchore/sbom-action` uploaded release assets by default; those side effects must be disabled when publication is intentionally gated.
3. **An environment name is not an approval gate.** A required reviewer must be configured and verified. An unprotected environment runs immediately.
4. **Checksums, signatures, and provenance are different controls.** Each needs its own generated artifact, verification result, and evidence link.
5. **Release identity is a security control.** Manifest, tag, embedded CLI identity, source commit, and evidence record must agree before publication.
6. **Stable release evidence must be built from the stable tag.** Preview artifacts can inform readiness but must not be promoted or relabeled.

## Process improvements

### Adopted during this cycle

- Version-specific release evidence files.
- Explicit release owner, reviewer, rollback owner, and decision fields.
- Protected approval environment before asset upload.
- Least-privilege job permissions.
- CI generation and verification of SBOM, signatures, and provenance.
- Root-level installation documentation.
- Umbrella issue references instead of premature closing keywords.

### Required before stable `v0.2.0`

- Correct manifest/tag identity in a new immutable preview or stable build.
- Complete macOS Intel, macOS arm64, and Windows runtime smoke evidence.
- Stable artifact build from the final reviewed stable tag.
- Stable evidence record with fresh checksums and workflow URL.
- Explicit owner/security review and stable `GO` decision.
- Confirm environment reviewer configuration remains active.

## Follow-up action register

| Action | Owner | Priority | Target | Acceptance evidence | Status |
| --- | --- | --- | --- | --- | --- |
| Add pre-tag manifest/tag identity assertion | `@geoffrey-xiao` | P1 | v0.2.0 | CI or release preflight fails on mismatch | Open |
| Add release workflow side-effect review checklist | `@geoffreyxiaoai` | P1 | v0.2.0 | Security review records action permissions and uploads | Open |
| Complete platform runtime smoke evidence | `@geoffrey-xiao` | P1 | v0.2.0 | Attached macOS/Windows command output | Open |
| Publish stable artifacts from final stable tag | `@geoffrey-xiao` | P0 | v0.2.0 | Stable workflow, checksums, approval evidence | Blocked by identity/evidence |
| Add a release preflight command or job | `@geoffrey-xiao` | P1 | v0.3.0 | Preflight validates tag, manifest, permissions, environment | Open |
| Add clean-download documentation smoke | `@geoffreyxiaoai` | P2 | v0.3.0 | macOS/Linux/Windows install walkthrough | Open |
| Define stable SBOM/signing/provenance policy | `@geoffrey-xiao` | P1 | v0.2.0 | Policy decision linked from stable evidence | Open |

## Remaining risks

- Preview.4 is not stable-ready because its manifest value does not match its tag.
- The release workflow is triggered by a published release, so the release page exists before the asset-upload approval gate; a future stable flow may use draft releases.
- Formal stable acceptance has not been recorded.
- Platform runtime evidence is not yet attached to the preview.4 evidence record.

## Retrospective acceptance

- [ ] Owner reviewed this retrospective.
- [ ] Reviewer reviewed the security and process findings.
- [ ] Follow-up issues were created or explicitly accepted in this issue.
- [ ] Remaining stable-release risks are linked from issue #148.
