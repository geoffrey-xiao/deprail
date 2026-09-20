# DepRail v0.3.0-preview.1 Retrospective

**Tracking issue:** [#215](https://github.com/geoffrey-xiao/deprail/issues/215)
**Release:** [`v0.3.0-preview.1`](https://github.com/geoffrey-xiao/deprail/releases/tag/v0.3.0-preview.1)
**Source commit:** `72fc603bd8b9c852c2accadd9521ea22a9720991`
**Release workflow:** [Run 35488833192](https://github.com/geoffrey-xiao/deprail/actions/runs/35488833192)
**Status:** Preview published; retrospective pending owner and security/architecture review
**Release mode:** Preview; stable `v0.3.0` is not approved
**Owner:** `@geoffrey-xiao`
**Reviewer:** `[security/architecture reviewer to confirm]`

## Executive summary

The v0.3 cycle delivered the first read-only remediation-planning workflow. A user can select a finding from an explicit scan report and generate a deterministic, schema-versioned plan without modifying the repository, installing packages, invoking package-manager mutation, or publishing source data.

The preview release was published from a reviewed, CI-passing `main` commit. The protected release workflow built Linux amd64, macOS amd64, macOS arm64, and Windows amd64 artifacts; generated and published an SPDX SBOM; created and verified Cosign signatures; generated and verified build provenance; checked SHA-256 manifests; ran release smoke checks; and required manual approval before uploading assets.

The cycle also exposed process weaknesses. Release evidence was initially generated only in ignored `local_test/`, the first tracked evidence record lagged behind the final tag and workflow, release preparation required a separate metadata PR after the implementation work, and the release issue/project state needed explicit synchronization. Owner approval was recorded as **GO WITH APPROVED GAPS** for the preview, while security/architecture review remains a separate requirement for stable release.

## What we delivered

### Product and contract capabilities

- Versioned remediation-plan domain model and schema.
- Deterministic stable plan identity.
- Explicit finding resolution from a supplied scan report.
- Read-only JavaScript/npm planning.
- Read-only Python planning.
- Read-only Java/Maven/Gradle planning.
- PURL and workspace identity preservation from scan findings into plans.
- Candidate state and compatibility evidence preservation, including unknown metadata.
- Explicit affected manifests and lockfiles.
- Provenance with source scan and artifact digests.
- Safe external plan output.
- Rejection of repository-internal and traversal output paths.
- Rejection of overwriting an existing plan output.
- JSON and terminal presenters.
- Mixed-repository scan-to-plan integration evidence.

### Release and distribution capabilities

- Version identity reconciled to `0.3.0-preview.1` before tagging.
- Immutable tag `v0.3.0-preview.1`.
- Cross-platform artifact matrix:
  - Linux amd64.
  - macOS amd64.
  - macOS arm64.
  - Windows amd64.
- SHA-256 checksum manifest.
- SPDX SBOM.
- Keyless Cosign signatures and certificates.
- GitHub build-provenance attestations.
- Protected manual release approval before asset upload.
- Version-specific release evidence under `docs/release-evidence/`.
- Release evidence linked to issue #215 and final release workflow output.

## What went well

- The read-only boundary held in the real workflow. The mixed fixture tree and lockfiles remained unchanged after scan and plan operations.
- The scan-to-plan path preserved the ecosystem PURL `pkg:npm/lodash@4.17.20` and workspace `frontend`, preventing the earlier unsupported-planner failure mode.
- Safety failures were explicit and stable:
  - Existing output: `PLAN_WRITE_FAILED`.
  - Repository-internal output: `PLAN_OUTPUT_OUTSIDE_ROOT_REQUIRED`.
  - Traversal output: `PLAN_OUTPUT_OUTSIDE_ROOT_REQUIRED`.
- `make verify` passed from the reviewed release baseline.
- CI passed on Ubuntu, macOS, and Windows for the release preparation and evidence PRs.
- The release workflow failed closed behind the `release-approval` environment until manual approval was granted.
- All four platform artifacts were checksum-verified, signed, provenance-verified, and uploaded only after approval.
- SBOM generation and upload behavior stayed behind the controlled publication job rather than bypassing the approval boundary.
- The release asset set is complete and understandable: binaries, signatures, certificates, checksums, and SBOM.
- Owner approval explicitly used preview semantics: **GO WITH APPROVED GAPS**, without claiming stable-release readiness.
- The CLI UX refactor was kept as planning-only work and did not silently expand v0.3 implementation scope.

## Faults and mistakes

| Fault | Impact | Root cause | Correction | Prevention |
| --- | --- | --- | --- | --- |
| Initial release evidence lived only in ignored `local_test/` | Evidence was not durable or reviewable from the repository | The manual procedure generated artifacts before a canonical record was created | Added `docs/release-evidence/RELEASE-v0.3-preview.1-EVIDENCE.md` and an `AGENTS.md` storage rule | Create the versioned evidence record before running release smoke tests |
| Initial tracked evidence referenced an older candidate commit and development binary | Evidence lagged behind the actual tagged release | Evidence was captured before metadata reconciliation and tagging | Updated the record after publication with tag, commit, workflow, checksums, and asset links | Make evidence update a required post-tag release step |
| Release metadata remained on `0.2.0-preview.3` until a dedicated preparation PR | Tag/manifest identity could have diverged again | Version metadata was not reconciled before release preparation began | Updated `.release-please-manifest.json` to `0.3.0-preview.1` and verified before tagging | Add an automated pre-tag assertion comparing manifest, intended tag, and CLI identity |
| Release closeout required multiple PRs | Extra coordination and review overhead | Evidence, metadata, and final artifact evidence were separated late in the cycle | Used separate focused PRs and documented the sequence | Define the release PR sequence in the release plan before implementation freeze |
| Owner approval and security/architecture review were easy to conflate | A preview could appear fully accepted while the stable review remained open | Decision fields and PR comments were separate from the evidence record | Recorded owner decision explicitly and retained security review as a separate open field | Require separate owner and security-review sign-offs in the release template |
| Historical project issues and statuses were stale | Release progress was harder to read and required manual repair | Project automation did not synchronize status automatically | Explicitly moved issue #215 through Review to Done after merge | Add a release closeout checklist for issue, project, milestone, and evidence state |
| Release asset inspection initially relied on a screenshot from another release | Risk of confusing v0.2 and v0.3 asset inventories | The screenshot showed `v0.2.0-preview.4` names while reviewing v0.3 | Queried the actual GitHub release and verified v0.3 assets directly | Always inspect release assets by tag with `gh release view` |
| The release branch accumulated unrelated local untracked files | Increased risk of accidental inclusion or confusion during release work | Local build artifacts were present in the working tree | Excluded them from commits and documented the boundary | Add a clean-checkout preflight and keep build output outside the repository root |

## Security and supply-chain lessons

1. **Release identity is a security control.** Manifest version, tag, source commit, embedded CLI identity, evidence record, and release assets must agree.
2. **Public `.sig` and `.pem` files are expected.** They are verification material, not secrets. Private signing material must remain in the signing system.
3. **Checksums, signatures, SBOM, and provenance prove different properties.** The release record must link each control separately rather than treating one as a substitute for the others.
4. **Approval environments must be verified in practice.** The release workflow must visibly pause before asset upload and resume only after approval.
5. **Build artifacts and evidence artifacts have different audiences.** Binaries and verification metadata belong on the GitHub Release; raw scan JSON, command logs, and local transcripts belong in evidence storage or CI artifacts.
6. **Preview approval is not stable approval.** Explicitly accepted preview gaps must remain visible and must not be silently promoted to the stable release gate.
7. **The v0.4 mutation boundary is higher risk than v0.3 planning.** Any future write capability needs isolation, explicit approval, rollback, bounded processes, cancellation, and post-change verification before implementation.

## Process improvements

### Adopted during this cycle

- Version-specific release evidence records.
- Canonical evidence location under `docs/release-evidence/`.
- Explicit release metadata reconciliation before tagging.
- Clean reviewed `main` verification before tag creation.
- Protected manual artifact publication.
- Cross-platform artifact, checksum, SBOM, signature, and provenance verification.
- Explicit owner decision wording for preview releases.
- Separation of release asset verification from raw workflow evidence.
- Dedicated short-lived branches for release follow-up changes.

### Required before stable `v0.3.0`

- Complete and record security/architecture review.
- Confirm the v0.3 Master Checklist is reconciled rather than leaving release rows stale.
- Confirm all v0.3 issue and project statuses reflect the merged implementation and evidence.
- Add a pre-tag manifest/tag/embedded-identity assertion to the release workflow.
- Add runtime smoke checks for each downloaded platform artifact where feasible.
- Define stable policy for SBOM, signatures, provenance, rollback, and artifact retention.
- Preserve a fresh stable evidence record built from the stable tag; do not reuse preview checksums.

## v0.4 handoff recommendations

The next release should begin with contract work, not repository mutation code. Freeze these decisions before creating implementation issues:

- Isolated worktree versus temporary repository copy.
- Exact approval state required before mutation.
- Atomic write and rollback guarantees.
- Allowed file and path boundaries.
- Package-manager command representation and execution policy.
- Timeout, cancellation, output-limit, and process-tree behavior.
- Test/build/type-check command discovery and execution.
- Post-change rescan and before/after comparison.
- Patch evidence format and provenance.
- Behavior after partial failure.

The first v0.4 implementation slice should prove:

```text
approved plan -> isolated workspace -> controlled mutation -> diff -> rollback evidence
```

Package-manager execution, verification commands, and rescan automation should follow only after that isolation boundary is proven.

## Follow-up action register

| Action | Owner | Priority | Target | Acceptance evidence | Status |
| --- | --- | --- | --- | --- | --- |
| Record security/architecture review for v0.3 preview | `@geoffrey-xiao` / reviewer | P1 | v0.3 closeout | Review decision linked to release evidence | Open |
| Reconcile v0.3 Master Checklist and historical project statuses | `@geoffrey-xiao` | P1 | v0.3 closeout | Checklist and Project state match merged evidence | Open |
| Add pre-tag identity assertion | `@geoffrey-xiao` | P1 | v0.4 preparation | Release preflight rejects manifest/tag/CLI mismatch | Open |
| Create v0.4 development plan and handoff contract | `@geoffrey-xiao` | P0 | v0.4 planning | Reviewed plan with frozen safety boundaries | Open |
| Define isolated mutation and rollback contract | `@geoffrey-xiao` / reviewer | P0 | v0.4 Sprint 1 | Contract tests cover success, failure, cancellation, and rollback | Open |
| Add cross-platform downloaded-artifact runtime smoke | `@geoffrey-xiao` | P1 | v0.4 preparation | Linux/macOS/Windows artifact execution evidence | Open |

## Remaining risks

- Security/architecture review is not recorded in the v0.3 evidence record.
- The preview release is not evidence for stable `v0.3.0` without a fresh stable-tag workflow.
- The CLI currently depends on a separately installed OSV-Scanner; automatic installation remains out of scope.
- v0.4 write capability could introduce materially higher risk if isolation and rollback are underspecified.
- Historical v0.2/v0.3 GitHub issue statuses may require separate reconciliation and should not be treated as current release evidence.

## Retrospective acceptance

- [ ] Owner reviewed this retrospective.
- [ ] Security/architecture reviewer reviewed the security and process findings.
- [ ] Follow-up issues were created or explicitly accepted.
- [ ] Remaining stable-release risks are linked from the applicable release issue.
