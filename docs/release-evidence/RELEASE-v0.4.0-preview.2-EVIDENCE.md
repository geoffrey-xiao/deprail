# DepRail v0.4.0-preview.2 Release Readiness Evidence

**Status:** Readiness in progress; not approved for publication
**Release mode:** Preview (planned)
**Target version:** `v0.4.0-preview.2`
**Readiness issue:** [#356](https://github.com/geoffrey-xiao/deprail/issues/356)
**Parent scope:** [#299](https://github.com/geoffrey-xiao/deprail/issues/299)
**Fix apply epic:** [#332](https://github.com/geoffrey-xiao/deprail/issues/332)
**Release checklist:** [`docs/RELEASE-CHECKLIST.md`](../RELEASE-CHECKLIST.md)
**Owner:** `@geoffrey-xiao`
**Architecture/security reviewer:** Pending named assignment

Unchecked rows are open gates, not implied passes. This record must be updated from the reviewed immutable commit before any preview publication.

## Implementation status

| Scope | Issue | PR | Status | Evidence |
|---|---:|---:|---|---|
| Plan validation and approval binding | #337 | #348 | Merged; issue remains open in Project Review | Linux/macOS/Windows CI passed; owner/security review disposition pending |
| Isolated detached worktree lifecycle | #338 | #349 | Merged; issue closed | Linux/macOS/Windows CI passed; cleanup failure and cancellation coverage |
| Allowlisted adapter mutation | #336 | #350 | Merged; issue closed | Linux/macOS/Windows CI passed; executable and symlink validation |
| Bounded verification and rescan | #340 | #351 | Merged; issue closed | Linux/macOS/Windows CI passed; timeout/cancellation/output-limit coverage |
| Finding transition classification | #339 | #352 | Merged; issue closed | Linux/macOS/Windows CI passed; unchanged/residual distinction |
| Rollback, discard, cleanup evidence | #345 | #353 | Merged; issue closed | Linux/macOS/Windows CI passed; structured cleanup result |
| Apply evidence and deterministic output | #347 | #354 | Merged; issue closed | Linux/macOS/Windows CI passed; atomic redaction-safe store |
| Security and cross-platform apply coverage | #346 | #355 | Merged; issue closed | Linux/macOS/Windows CI passed; shell-metacharacter coverage |
| OSV-Scanner 2.6.0 lockfile compatibility and path boundaries | #356 / #379 | #378 | Merged as `89c1ab7`; issue closed | Linux/macOS/Windows CI passed; lockfile, traversal, symlink, and shared-workspace regression coverage |
| Explicit approved verification commands | #381 | #382 | Merged as `a38d3e8`; issue closed | Linux/macOS/Windows CI passed; shell/path/input-boundary coverage; Darwin explicit apply smoke passed |
| User-facing `deprail fix apply` orchestration | #332 | #359, #382 | Parent closed; verification follow-up merged | Isolated apply, evidence, and explicit verification coverage; complete release smoke still pending |

## Release identity

- [ ] Intended version and preview mode recorded in the release manifest.
- [ ] Latest stable, preview, and release-candidate tags checked.
- [ ] Reviewed `origin/main` commit recorded.
- [ ] New immutable `v0.4.0-preview.2` tag confirmed available.
- [ ] CLI identity matches version, tag, and source commit.
- [ ] Release owner and independent reviewer named.

## Plan and scope completion

- [ ] Preview.2 execution package and release outcome reconciled.
- [ ] #332 `deprail fix apply` implementation complete or owner-approved disposition recorded.
- [ ] All child issues have merged PR, CI, and acceptance evidence.
- [ ] No open release-blocking P0/P1 issue lacks explicit disposition.
- [ ] Included, excluded, and deferred capabilities documented.
- [ ] Compatibility impact, migration notes, remaining risks, and rollback scope recorded.

## Contract and security completion

- [ ] CLI and exit-code contracts for `fix apply` reviewed.
- [ ] Versioned plan, evidence, and output schemas validated.
- [ ] Determinism and stable ordering verified.
- [ ] Incomplete, stale, unsupported, cancelled, failed, and cleanup-failed states explicit.
- [x] Path containment and symlink boundaries covered by merged child PRs.
- [x] External-process, timeout, cancellation, and output-limit behavior covered by merged child PRs.
- [x] Package-script, shell, credential, and file-write boundaries have regression coverage.
- [ ] Independent architecture/security review recorded.

## Automated and manual verification

- [ ] `make verify` passes from a clean checkout of the reviewed release commit.
- [ ] `git diff --check` passes.
- [x] Linux CI passes for merged child PRs.
- [x] macOS CI passes for merged child PRs.
- [x] Windows CI passes for merged child PRs.
- [ ] Manual representative-repository `fix apply` smoke passes.
- [ ] Before/after scan produces resolved, unchanged, introduced, and residual transitions.
- [ ] Caller repository status and digest remain unchanged.
- [ ] Required JavaScript, Python, and Java fixtures are covered.

## Artifact and supply-chain verification

- [ ] Authoritative release-binded CI build job is enabled and reviewed.
- [ ] Artifacts built from reviewed immutable commit and tag.
- [ ] Supported platform matrix complete.
- [ ] Artifact names and sizes recorded.
- [ ] SHA-256 checksums independently verified.
- [ ] `doctor` identity matches release tag and commit.
- [ ] SBOM supplied or explicitly accepted as unavailable for preview mode.
- [ ] Signatures supplied or explicitly accepted as unavailable for preview mode.
- [ ] Provenance supplied or explicitly accepted as unavailable for preview mode.
- [ ] Public release assets exclude raw logs, scan reports, credentials, and internal evidence.

## Publication and approval

- [ ] Preview release notes explain outcome, installation, usage, scope, limitations, and feedback.
- [ ] GitHub Release marked as preview.
- [ ] Protected approval environment configured and verified.
- [ ] Release workflow completed successfully.
- [ ] Assets uploaded only after approval gate.
- [ ] Release, workflow, evidence, and comparison URLs recorded.
- [ ] Rollback owner and immutable-tag recovery procedure recorded.
- [ ] Owner records `go`, `go with approved gaps`, or `no-go`.

## Post-release closure

- [ ] Version-specific evidence complete.
- [ ] `RETROSPECTIVE-v0.4.0-preview.2.md` created.
- [ ] Retrospective linked from release issue and evidence.
- [ ] Follow-up actions have owners and acceptance evidence.
- [ ] Owner acceptance and security/architecture review recorded separately.
- [ ] Project status reflects actual release state.
- [ ] Release issue closed only after evidence and review are complete.

## Current blockers

1. Confirm the authoritative release-binded CI workflow and validate its preview.2 artifact controls before publication.
2. Name an independent architecture/security reviewer.
3. Run complete Python, Java, and cross-platform representative-repository release smoke.
4. Produce preview.2 artifact, SBOM, signature, provenance, checksum, and publication evidence.

## Manual command-surface observation

The merged `main` binary exposes `deprail fix apply`. On Darwin arm64, an explicit `node --check verify.js` command was approved and executed in the isolated worktree:

- `go test ./...` — exit `0`;
- `fix apply --help` — exit `0`;
- real apply — exit `0`, `outcome: applied`;
- verification and rescan complete;
- finding transition `resolved`;
- cleanup status `succeeded`.

Package-manager lifecycle scripts are not inferred or executed. The repository-local release workflows are not the authoritative publication path for this release; the complete release gate remains open for the release-binded CI and the blockers listed above.

## Decision record

- Owner decision: Pending
- Architecture/security review: Pending
- Rollback owner: Pending
- Publication URL: Pending
