# DepRail v0.4.0-preview.1 Retrospective

**Tracking issue:** [#334](https://github.com/geoffrey-xiao/deprail/issues/334)
**Release evidence:** [`RELEASE-v0.4.0-preview.1-EVIDENCE.md`](../release-evidence/RELEASE-v0.4.0-preview.1-EVIDENCE.md)
**Release issue:** [#325](https://github.com/geoffrey-xiao/deprail/issues/325)
**Release:** [v0.4.0-preview.1](https://github.com/geoffrey-xiao/deprail/releases/tag/v0.4.0-preview.1)
**Status:** Published preview; retrospective closeout pending owner acceptance and independent architecture/security review
**Release mode:** Preview only; no stable-release claim
**Owner:** `@geoffrey-xiao`
**Reviewer:** Independent architecture/security reviewer to be named

## Executive summary

`v0.4.0-preview.1` delivered the remediation safety primitives and trusted baseline generation needed for explainable dependency-security workflows. The release includes isolated workspace and rollback boundaries, approval binding, bounded mutation primitives, verification and rescan transition classification, redaction-safe evidence, cross-platform path semantics, and `deprail baseline create`.

The preview was published manually at reviewed commit `eaf1d4853d13dd8f7233ac0bcc53dea8d824b81a`. Release assets include macOS amd64/arm64, Linux amd64, Windows amd64 binaries, an SPDX SBOM, and `SHA256SUMS`. Downloaded release binaries passed checksum verification.

The release did not include `deprail fix apply`. That feature is explicitly deferred to `v0.4.0-preview.2` under [#332](https://github.com/geoffrey-xiao/deprail/issues/332).

## What we delivered

- Isolated workspace and rollback boundary ([#302](https://github.com/geoffrey-xiao/deprail/issues/302), [#319](https://github.com/geoffrey-xiao/deprail/pull/319)).
- Approval binding to remediation plan and source identity ([#307](https://github.com/geoffrey-xiao/deprail/issues/307), [#320](https://github.com/geoffrey-xiao/deprail/pull/320)).
- Bounded direct-argv mutation boundary ([#308](https://github.com/geoffrey-xiao/deprail/issues/308), [#321](https://github.com/geoffrey-xiao/deprail/pull/321)).
- Verification and rescan transition classification ([#309](https://github.com/geoffrey-xiao/deprail/issues/309), [#322](https://github.com/geoffrey-xiao/deprail/pull/322)).
- Versioned redaction-safe evidence ([#310](https://github.com/geoffrey-xiao/deprail/issues/310), [#323](https://github.com/geoffrey-xiao/deprail/pull/323)).
- Cross-platform relative-path semantics ([#311](https://github.com/geoffrey-xiao/deprail/issues/311), [#324](https://github.com/geoffrey-xiao/deprail/pull/324)).
- Trusted baseline generation from complete scans ([#329](https://github.com/geoffrey-xiao/deprail/issues/329), [#331](https://github.com/geoffrey-xiao/deprail/pull/331)).
- Manual scan, baseline, diff, policy-help, and fix-plan smoke coverage on macOS arm64.

## What went well

- The merged implementation PRs passed Linux, macOS, and Windows CI.
- `make verify` passed on reviewed `origin/main` commit `eaf1d48` on macOS arm64.
- The representative mixed repository scan completed with 14 findings and no scan errors.
- Baseline creation rejected unsafe or stale inputs and produced a complete baseline.
- Identical baseline comparison produced deterministic `unchanged` results.
- Release assets covered the intended four binary targets and included an SPDX SBOM.
- Downloaded release binaries passed `sha256sum -c` against the published `SHA256SUMS`.
- The `fix apply` scope gap was made explicit and tracked instead of being represented as implemented.

## Mistakes and corrections

| Mistake | Impact | Root cause | Correction | Prevention |
|---|---|---|---|---|
| Manual guide initially reused an existing baseline output path | `BASELINE_OUTPUT_EXISTS` interrupted reruns | The guide did not separate clean creation from overwrite-protection verification | Added cleanup of only the ignored manual baseline before the success path; retained the collision case as a negative test | Make manual procedures idempotent while preserving explicit no-overwrite checks |
| Identical baselines were interpreted as a meaningful security diff | Output contained only `unchanged` findings, which caused confusion | The smoke scenario tested determinism, not before/after remediation | Documented that real diff requires two scans from different repository states | Add a committed before/after remediation fixture for future diff evidence |
| Release workflow build job was disabled with `if: ${{ false }}` | CI did not generate signed/provenanced release artifacts | Workflow was retained as a non-runnable scaffold while manual publication proceeded | Recorded the gap in release evidence and retrospective; release remained preview-only | Never publish until the release workflow is executable or the manual supply-chain exception is explicitly approved |
| `fix apply` was absent from preview.1 despite appearing in broader v0.4 planning | Apply orchestration could not be manually verified | Preview sequencing was not explicit early enough | Deferred implementation to preview.2 and created backlog epic [#332](https://github.com/geoffrey-xiao/deprail/issues/332) | Reconcile each preview's included, excluded, and deferred capabilities before release work |

## Security and supply-chain lessons

1. A published checksum file is not equivalent to signatures or provenance attestations.
2. Manual publication must identify the exact reviewed commit and preserve immutable tag evidence.
3. Release workflow gates must be executable and observable; a disabled build job must be treated as a release gap.
4. The caller worktree and release evidence remain separate trust boundaries.
5. `fix apply` requires independent architecture/security review before any mutation capability ships.
6. Preview approval does not imply stable-release readiness.
7. SBOM, signing, provenance, checksums, and binary identity are separate evidence rows.

## Process improvements

### Adopted

- Version-specific release evidence under `docs/release-evidence/`.
- Explicit preview sequencing and deferred-scope tracking.
- Local manual output isolation under ignored `local_test/` paths.
- Published checksum download and verification.
- Separate GitHub epic and issue contract for `fix apply`.

### Required before preview.2

- Enable and review the release workflow build job.
- Run the workflow through SBOM, signing, provenance, and release-approval stages.
- Add a deterministic committed before/after remediation fixture for diff evidence.
- Implement and review `fix apply` under [#332](https://github.com/geoffrey-xiao/deprail/issues/332).
- Obtain named architecture/security review before closing release evidence.

## Follow-up actions

| Action | Owner | Target | Acceptance evidence | Status |
|---|---|---|---|---|
| Implement isolated `deprail fix apply` | `@geoffrey-xiao` / reviewer | `v0.4.0-preview.2` | #332 implementation PR, security tests, cross-platform smoke | Open |
| Enable and validate release CI | `@geoffrey-xiao` | Before preview.2 publication | Successful workflow with SBOM, signatures, provenance, and approval | Open |
| Add committed before/after diff fixture | `@geoffrey-xiao` | Before preview.2 evidence | Base/head scans produce a non-unchanged diff deterministically | Open |
| Name architecture/security reviewer | Project owner | Retrospective closeout | Review decision linked to #334 and release evidence | Open |
| Complete retrospective owner acceptance | `@geoffrey-xiao` | Post-release closeout | All acceptance checkboxes reviewed | Open |

## Remaining risks and rollback

- The release was manually published while the configured release workflow build job was disabled.
- Published assets have checksums and an SPDX SBOM, but no evidenced CI-generated signatures or provenance attestations.
- `fix apply` remains unimplemented and deferred to preview.2.
- Architecture/security review is not yet recorded.
- Rollback means preserving the immutable tag and evidence, marking the preview unsuitable for stable promotion, and publishing a corrected new preview tag rather than moving `v0.4.0-preview.1`.

## Acceptance

- [ ] Owner reviewed every retrospective criterion.
- [ ] Architecture/security reviewer separately reviewed security and process findings.
- [x] Preview release evidence is linked.
- [x] Published assets, checksums, and SBOM status are recorded.
- [ ] Signatures and provenance are supplied or explicitly accepted as unavailable for this preview.
- [x] Follow-up issue #332 and retrospective issue #334 are recorded.
- [ ] Release issue, evidence record, and retrospective are cross-linked.
