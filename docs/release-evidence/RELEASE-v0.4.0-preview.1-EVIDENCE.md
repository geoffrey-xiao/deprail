# DepRail v0.4.0 Preview.1 Release Evidence Record

**Status:** Release preparation; not approved for publication
**Release mode:** Preview
**Target version:** `v0.4.0-preview.1`
**Release issue:** [#325](https://github.com/geoffrey-xiao/deprail/issues/325)
**Parent scope:** [#299](https://github.com/geoffrey-xiao/deprail/issues/299)
**Decision:** Pending owner go/no-go

This record is a release-gate index. Unchecked rows are open gaps, not implied passes.

Reusable gate: [`docs/RELEASE-CHECKLIST.md`](../RELEASE-CHECKLIST.md)

## Source and identity

| Field | Evidence |
|---|---|
| Reviewed source commit | [`eaf1d48`](https://github.com/geoffrey-xiao/deprail/commit/eaf1d48) (`origin/main` after documentation PR #333 merge) |
| Latest stable tag | None; latest preview is `v0.3.1-preview.1` |
| Existing preview tags | `v0.3.1-preview.1`, `v0.3.0-preview.1`, `v0.2.0-preview.4`, `v0.2.0-preview.3`, `v0.2.0-preview.2`, `v0.2.0-preview.1`, `v0.1.0-preview.1` |
| Release manifest | Not required for this manually run release; manual version/tag/release identity will be recorded here |
| Release tag | Not created; publication not authorized |
| DepRail version | Pending release workflow |
| Release owner | `@geoffrey-xiao` |
| Architecture/security reviewer | Pending named review |
| Release workflow | Not runnable as configured; `.github/workflows/release.yml` currently has `build.if: ${{ false }}` |

## Scope

Included implementation contracts and primitives:

- isolated workspace and rollback boundary;
- approval binding to plan and source identity;
- bounded direct-argv mutation boundary;
- verification command selection and finding transitions;
- versioned redaction-safe evidence;
- cross-platform relative-path semantics;
- trusted baseline generation from complete scan results under [#329](https://github.com/geoffrey-xiao/deprail/issues/329), merged by [#331](https://github.com/geoffrey-xiao/deprail/pull/331).

The baseline-generation manual workflow is recorded in `local_test/v0.4.0-preview.1/guide/MANUAL-TEST-GUIDE.md`; generated artifacts remain local-only.

The current merged slices do not include a user-facing `deprail fix apply` orchestration flow. That feature is explicitly deferred to `v0.4.0-preview.2` under backlog epic [#332](https://github.com/geoffrey-xiao/deprail/issues/332); preview.1 must not claim it is implemented or manually verified.

## Implementation evidence

| Area | Issue | PR | Status |
|---|---:|---:|---|
| Isolation and rollback | #302 | #319 | Merged; CI passed on Linux/macOS/Windows |
| Approval binding | #307 | #320 | Merged; CI passed on Linux/macOS/Windows |
| Bounded mutation | #308 | #321 | Merged; CI passed on Linux/macOS/Windows |
| Verification and rescan | #309 | #322 | Merged; CI passed on Linux/macOS/Windows |
| Evidence contract | #310 | #323 | Merged; CI passed on Linux/macOS/Windows |
| Cross-platform semantics | #311 | #324 | Merged; CI passed on Linux/macOS/Windows |
| Trusted baseline generation | #329 | #331 | Merged; CI passed on Linux/macOS/Windows; reviewed binary smoke passed on macOS arm64 |

## Checklist reconciliation

| Area | Status | Evidence or gap |
|---|---|---|
| Release identity and immutable tag | Partial | Candidate source captured at `eaf1d48`; release tag and final CLI identity remain pending |
| Plan and scope completion | Partial | Preview.1 baseline/diff scope is implemented; complete apply orchestration is intentionally deferred to preview.2 under #332 |
| Contract and security completion | Partial | Contracts exist; named release security/architecture review pending |
| Automated verification | Complete for merged preview.1 scope | `make verify` passed on macOS arm64; PR #331 and #333 CI passed on Linux/macOS/Windows |
| Manual verification | Partial | Scan, baseline, diff, policy-help, and fix-plan workflow passed on macOS arm64; fix apply is deferred to preview.2 |
| Artifact and supply-chain verification | Partial | Local candidate artifacts built for macOS arm64, Linux amd64, and Windows amd64 with verified SHA-256 checksums; the release CI build job is disabled, so CI-generated SBOM/signatures/provenance are not available |
| Publication and approval | Blocked | Manual publication would bypass the configured release workflow; owner decision alone does not provide CI-generated supply-chain evidence |
| Post-release closure | Open | Retrospective cannot be completed before release decision |

## Required next evidence

- [x] Capture release baseline: `origin/main`, tags, releases, and manifest status. Manifest is not required for the owner-run manual release process.
- [x] Run `make verify` from reviewed `origin/main` `eaf1d48` on macOS arm64; `go generate`, `go vet`, `go test`, and `go build` completed successfully.
- [x] Run representative repository scan and baseline/diff smoke; macOS arm64, source `7903612`, binary `/tmp/deprail-v0.4.0-main`, scan complete with 14 findings, baseline create exit `0`, identical diff exit `0`.
- [x] Reconcile complete `deprail fix apply` scope for preview.1: explicitly deferred to preview.2 under [#332](https://github.com/geoffrey-xiao/deprail/issues/332); no implementation claim is made for preview.1.
- [x] Build candidate artifacts from `eaf1d48`: `deprail-darwin-arm64` (5,737,698 bytes), `deprail-linux-amd64` (6,165,199 bytes), and `deprail-windows-amd64.exe` (6,174,720 bytes). Checksums are in local-only `local_test/v0.4.0-preview.1/artifacts/SHA256SUMS` and verified with `shasum -a 256 -c`.
- [x] Record reviewed candidate identity: macOS artifact reports `deprail 0.4.0-preview.1 (tag=unreleased-candidate commit=eaf1d48...)`; final immutable tag remains pending.
- [ ] Prepare preview notes and rollback procedure.
- [ ] Enable and review the release workflow build job before publication; current `build.if: ${{ false }}` prevents the `manual-verify` job from running.
- [ ] Record owner go/no-go decision.
- [ ] Create and link retrospective after publication.

## Review and decision

- Owner decision: `[pending]`
- Decision date: `[pending]`
- Architecture/security review: `[pending]`
- Rollback owner: `[pending]`
- Rollback procedure: preserve evidence, do not move immutable tags, and publish a new preview only after resolving approved gaps.
