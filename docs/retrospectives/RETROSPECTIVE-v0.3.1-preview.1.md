# DepRail v0.3.1-preview.1 Retrospective

**Tracking issue:** [#239](https://github.com/geoffrey-xiao/deprail/issues/239)
**Release evidence:** [`RELEASE-v0.3.1-preview.1-EVIDENCE.md`](../release-evidence/RELEASE-v0.3.1-preview.1-EVIDENCE.md)
**Release gate:** [#253](https://github.com/geoffrey-xiao/deprail/issues/253)
**Status:** Preview published; retrospective closeout pending final owner/security review and supply-chain evidence
**Release mode:** Preview only; no stable-release claim
**Owner:** `@geoffrey-xiao`
**Reviewer:** `@geoffreyxiaoai` / architecture-security reviewer to confirm

## Executive summary

The v0.3.1 cycle delivered a presentation-focused CLI refinement without changing scanner semantics, domain meaning, repository mutation boundaries, or machine-readable result contracts. The implementation added structured lifecycle events, terminal capability detection, shared human output primitives, truthful outcome summaries, interactive progress, quiet/verbose behavior, hostile-label protection, lifecycle start notices, and actionable help tips.

The release work also exposed weaknesses in version and publication automation. Release Please repeatedly calculated `0.4.0-preview.1` for the intended `0.3.1-preview.1` preview and was replaced with an explicit release workflow. The combined workflow then required fixes for shell regex quoting, pinned Go setup, Git tag identity, and job decomposition. These corrections are recorded as process lessons, not hidden as successful first-pass behavior.

This retrospective records the published preview and remains open only for final evidence and review closeout.

## What we delivered

- Structured lifecycle presentation events with real application boundaries.
- Terminal capability selection and non-TTY safeguards.
- Shared human headers, sections, statuses, summaries, and safe labels.
- Explicit complete, partial, failed, and cancelled outcome semantics.
- Interactive progress on approved TTY streams only.
- Quiet and verbose scan behavior without contaminating JSON stdout.
- Hostile-label encoding and credential-bearing URL redaction boundaries.
- Lifecycle start notices for discovery, scan, baseline comparison, and remediation planning.
- Actionable tips and examples on root and command help.
- JSON stdout purity and stderr diagnostics separation.
- Initial safe color primitives; full color integration deferred to backlog [#281].

## What went well

- Application services remained separate from terminal-specific rendering.
- JSON output stayed free of incidental banners, progress, and ANSI sequences in smoke verification.
- Help output became directly actionable: all six command pages include a tip and example.
- CI passed across Linux, macOS, and Windows for the implementation PRs.
- Hostile labels remained presentation data rather than executable terminal input.
- The combined release workflow completed source verification, cross-platform artifact builds, checksums, and SBOM publication for `v0.3.1-preview.1`.
- Release failures were caught before publication; no incorrect `0.4.0-preview.1` release was published.
- The preview was published only after explicit owner approval and remains clearly separated from stable-release approval.

## Mistakes and corrections

| Mistake | Impact | Root cause | Correction | Prevention |
| --- | --- | --- | --- | --- |
| Release Please calculated `0.4.0-preview.1` instead of the approved `0.3.1-preview.1` | Release preparation was blocked and several incorrect release PRs had to be closed | Conventional-commit bump rules did not match the preview line and action-level override support was misunderstood | Replaced Release Please with an explicit workflow using owner-selected release inputs | Keep release version calculation in one tested workflow with exact preview/RC rules |
| Release workflow used an unsupported `release-as` action input | Workflow dispatch rejected the input | Action version and manifest-mode behavior were not verified before configuration | Removed the unsupported input and abandoned the Release Please path | Verify third-party action inputs against the pinned action version before merging |
| Shell `sed` expression expanded `$#` | Version resolver failed before verification | Regex was embedded in a double-quoted shell expression | Escaped the end anchor and added a focused fix PR | Avoid shell interpolation in regexes; use a tested script or single-quoted expressions |
| Release runner used Go 1.27.0 while the repository requires Go 1.27.1 | Release verification failed before tag creation | Release workflow omitted `actions/setup-go` | Install Go from `.go-version` before verification | Every release workflow must install the pinned toolchain explicitly |
| Annotated tag creation lacked Git identity | Combined release failed before tag creation | CI runner does not have a configured committer identity | Configure the GitHub Actions bot identity | Treat tag creation as a separately tested publication boundary |
| Combined workflow initially hid all verification behind `make verify` | Diagnosis was slower when verification failed | Release pipeline favored compactness over failure visibility | Split preparation, verification, publication, matrix builds, and finalization into jobs | Keep release jobs independently observable and retryable |
| Manual guide changes were stored under ignored `local_test/` | Durable review and release evidence were difficult to maintain | Local evidence path was not tracked by default | Added the canonical release evidence record under `docs/release-evidence/` | Create the versioned evidence record before manual capture |

## Security and supply-chain lessons

1. Release identity, tag identity, embedded CLI identity, source commit, evidence record, and assets must be reconciled before publication.
2. A release workflow must fail before creating a tag when verification, version validation, or approval prerequisites fail.
3. GitHub environment approval must remain before artifact publication, not after public assets are uploaded.
4. Checksums, SBOM, signatures, and provenance are separate controls and must each have evidence.
5. Release workflows require least-privilege permissions and explicit tag/release boundaries.
6. Shell parsing of tags and versions is security-sensitive; malformed or hostile tag data must not alter commands.
7. Preview approval does not imply stable compatibility approval.
8. Deferred v0.4 mutation work remains outside this release and requires a separate isolation, rollback, and verification contract.

## Process improvements

### Adopted

- Version-specific release evidence under `docs/release-evidence/`.
- Explicit release mode and owner confirmation.
- Pinned Go setup from `.go-version`.
- Separate visible release verification steps.
- Matrix platform artifact builds.
- Protected release-approval environment.
- Automatic prerelease counter resolution from existing tags.
- Explicit deferral of nonessential color integration.

### Required before final preview publication

- Complete and review the combined release workflow, including signing and provenance attestation.
- Run the combined workflow successfully through artifact upload.
- Record checksums, SBOM, signatures, provenance, and platform smoke results.
- Record owner go/no-go and architecture/security review separately.
- Link the final release tag and workflow run from the evidence record.

## Follow-up actions

| Action | Owner | Target | Acceptance evidence | Status |
| --- | --- | --- | --- | --- |
| Complete combined release workflow signing and provenance | `@geoffrey-xiao` | v0.3.1 closeout | Successful signed and attested artifact workflow | Open |
| Complete v0.3.1 preview evidence record | `@geoffrey-xiao` | v0.3.1 closeout | Tag, artifact, checksum, SBOM, and workflow links | Open |
| Record architecture/security review | `@geoffrey-xiao` / reviewer | v0.3.1 closeout | Review decision linked to #253 | Open |
| Integrate semantic lifecycle colors | `@geoffrey-xiao` | Backlog / v0.4 planning | Issue [#281](https://github.com/geoffrey-xiao/deprail/issues/281) acceptance evidence | Deferred |
| Create v0.4 mutation and rollback contract | `@geoffrey-xiao` / reviewer | v0.4 planning | Reviewed contract and safety tests | Open |

## Remaining risks and rollback

- The preview release has not yet completed the combined artifact workflow.
- Signing and provenance are not yet integrated into the new combined pipeline.
- Interactive TTY evidence remains a release-gate decision point.
- OSV-Scanner remains an external prerequisite; automatic installation is out of scope.
- Rollback requires preserving the evidence and artifacts, withdrawing the preview reference if necessary, and creating a new immutable preview tag after correction.

## Acceptance

- [ ] Owner reviewed every retrospective criterion.
- [ ] Security/architecture reviewer reviewed security and process findings.
- [ ] Preview release evidence is complete.
- [ ] Artifact checksums, SBOM, signatures, and provenance are linked.
- [ ] Follow-up issues and owners are recorded.
- [ ] Release issue, evidence record, and retrospective are cross-linked.
