# DepRail v0.4.0-preview.2 Retrospective

**Tracking issue:** [#356](https://github.com/geoffrey-xiao/deprail/issues/356)  
**Release evidence:** [`RELEASE-v0.4.0-preview.2-EVIDENCE.md`](../release-evidence/RELEASE-v0.4.0-preview.2-EVIDENCE.md)  
**Release:** [v0.4.0-preview.2](https://github.com/geoffrey-xiao/deprail/releases/tag/v0.4.0-preview.2)  
**Workflow:** [Release Combined run 35719363746](https://github.com/geoffrey-xiao/deprail/actions/runs/35719363746)  
**Status:** Published preview; closeout pending independent architecture/security review and owner acceptance  
**Release mode:** Preview only; no stable-release claim  
**Owner:** `@geoffrey-xiao`  
**Reviewer:** Independent architecture/security reviewer to be named

## Executive summary

`v0.4.0-preview.2` delivered the first reviewed `deprail fix apply` orchestration path. The release adds plan-bound approvals, isolated application, direct-argv mutation boundaries, explicit verification commands, verification and rescan transitions, cleanup reporting, and evidence suitable for review. The preview was published from reviewed `main` SHA `0514f5f2f1450407e2784a0c020d3c55d7f9039c`.

The release workflow completed successfully across Linux amd64, macOS amd64, macOS arm64, and Windows amd64. The published release contains the four binaries, an SPDX SBOM, and `SHA256SUMS`; independently downloaded binaries passed checksum verification. Signatures and provenance were not present in the published assets and remain an explicitly accepted preview gap.

The owner decision is `go with approved gaps`. This is not stable-release approval.

## What we delivered

- Isolated approved remediation application with worktree and cleanup boundaries ([#332](https://github.com/geoffrey-xiao/deprail/issues/332), [#359](https://github.com/geoffrey-xiao/deprail/pull/359)).
- Approval binding to the plan and source identity, including stale or mismatched approval rejection ([#360](https://github.com/geoffrey-xiao/deprail/issues/360)).
- Explicit verification command input with bounded execution and safe path/executable validation ([#381](https://github.com/geoffrey-xiao/deprail/issues/381), [#382](https://github.com/geoffrey-xiao/deprail/pull/382)).
- Verification, rescan, transition, and cleanup evidence for successful application.
- OSV lockfile compatibility correction ([#378](https://github.com/geoffrey-xiao/deprail/pull/378)).
- Four-platform preview artifacts, SPDX SBOM, and independently verified checksums.
- Release evidence updated with the owner-approved gap disposition.

## What went well

- The explicit verification surface was narrowed to approved commands rather than inferred package-manager lifecycle scripts.
- Direct-argv execution, path validation, source binding, and cleanup behavior received focused review before publication.
- The real npm apply smoke completed with `outcome: applied`, verification/rescan completion, `resolved` transition, and successful cleanup.
- Python and Java scans completed with status `complete` on representative fixtures.
- Invalid Python and Java remediation keys correctly returned `FINDING_NOT_FOUND` with exit code `2`.
- Approval from a different source identity correctly returned `APPROVAL_SOURCE_MISMATCH` with exit code `3`.
- Release Combined completed all prepare, verify, publish, platform-build, and finalize jobs successfully.
- All four downloaded release binaries passed independent SHA-256 verification.

## Gaps and mistakes

| Gap or mistake | Impact | Root cause | Correction or disposition | Prevention |
|---|---|---|---|---|
| The release target was initially confused with an unrelated `v0.3.1-preview.2` publication | Evidence temporarily pointed at the wrong release | The release workflow was dispatched more than once during closeout | Reconciled evidence against the actual `v0.4.0-preview.2` tag and run | Verify exact tag, source SHA, and release URL before recording publication evidence |
| Python and Java fixtures contained no findings | Successful Python/Java remediation apply was not observable | Fixtures were suitable for scan coverage but not remediation mutation coverage | Accepted as a preview gap; do not fabricate findings | Add deterministic vulnerable fixtures with legitimate scanner evidence |
| Published assets had no signatures or provenance | Supply-chain authenticity evidence is incomplete | The authoritative release-binded path publishes checksums and SBOM but no observed signature/provenance assets | Accepted explicitly for this preview; not acceptable by default for stable release | Make signature and provenance requirements a release gate |
| Cross-platform representative smoke was incomplete | Platform artifact builds passed, but end-to-end behavior was not exercised on every target | Manual access and fixture coverage were limited | Accepted as a preview gap | Automate representative smoke against downloaded artifacts on each supported OS |
| Independent architecture/security review was not recorded before publication | Security sign-off remains incomplete | Owner approval proceeded with documented gaps | Keep release in preview status and obtain review before stable promotion | Require separate owner and security-review checklist decisions |

## Security and supply-chain lessons

1. A successful artifact build is not equivalent to signature or provenance verification.
2. Checksums establish content integrity only after the expected artifact is identified; they do not establish publisher identity.
3. Explicit verification commands must remain bounded, reviewable, and source-root constrained.
4. Approval records and source identity are security boundaries, not administrative metadata.
5. Preview approval with gaps must never be represented as stable-release readiness.
6. Ecosystem scan coverage does not prove ecosystem remediation coverage.

## Process improvements

### Adopted

- Version-specific release evidence and retrospective records.
- Explicit owner disposition using `go with approved gaps`.
- Independent downloaded-artifact checksum verification.
- Separate tracking of scan coverage versus remediation-apply coverage.
- Release workflow, artifact, and source-SHA URLs recorded together.

### Required before stable release

- Name and record an independent architecture/security reviewer.
- Add legitimate deterministic vulnerable Python and Java remediation fixtures.
- Run downloaded-binary smoke on Linux, macOS, and Windows.
- Supply and verify artifact signatures and provenance.
- Record rollback ownership and immutable-tag recovery procedures.

## Follow-up actions

| Action | Owner | Target | Acceptance evidence | Status |
|---|---|---|---|---|
| Add deterministic Python and Java remediation fixtures | `@geoffrey-xiao` | Before stable release | Successful plan, approval, apply, verify, rescan, and cleanup evidence | Open |
| Complete cross-platform downloaded-binary smoke | `@geoffrey-xiao` / CI owner | Before stable release | Linux/macOS/Windows command matrix and artifacts | Open |
| Add release artifact signatures and provenance | Release owner | Before stable release | Published signature/provenance assets and independent verification | Open |
| Name independent architecture/security reviewer | Project owner | Preview closeout | Review decision linked to #356 and this retrospective | Open |
| Record rollback owner and recovery procedure | Project owner | Preview closeout | Evidence checklist and owner recorded | Open |
| Complete retrospective owner acceptance | `@geoffrey-xiao` | Post-release closeout | All acceptance criteria reviewed | Open |

## Remaining risks and rollback

- The preview contains no observed signatures or provenance attestations.
- Python and Java remediation apply remain unverified.
- Cross-platform artifact execution remains less complete than platform build coverage.
- Architecture/security review remains pending.
- Rollback means preserving the immutable `v0.4.0-preview.2` tag and evidence, marking the preview unsuitable for stable promotion, and publishing a corrected new preview rather than moving the existing tag.

## Acceptance

- [ ] Owner reviewed every retrospective criterion.
- [ ] Architecture/security reviewer separately reviewed security and process findings.
- [x] Preview release evidence is linked.
- [x] Published assets, checksums, and SBOM status are recorded.
- [ ] Signatures and provenance are supplied or explicitly accepted as unavailable for this preview.
- [x] Release issue #356 is recorded.
- [ ] Release issue, evidence record, and retrospective are cross-linked.
