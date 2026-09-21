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
| Reviewed source commit | [`7903612`](https://github.com/geoffrey-xiao/deprail/commit/7903612) (`origin/main` after PR #331 merge) |
| Latest stable tag | None; latest preview is `v0.3.1-preview.1` |
| Existing preview tags | `v0.3.1-preview.1`, `v0.3.0-preview.1`, `v0.2.0-preview.4`, `v0.2.0-preview.3`, `v0.2.0-preview.2`, `v0.2.0-preview.1`, `v0.1.0-preview.1` |
| Release manifest | Not required for this manually run release; manual version/tag/release identity will be recorded here |
| Release tag | Not created; publication not authorized |
| DepRail version | Pending release workflow |
| Release owner | `@geoffrey-xiao` |
| Architecture/security reviewer | Pending named review |
| Release workflow | Pending |

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

The current merged slices do **not** yet prove a complete user-facing `deprail fix apply` orchestration flow. That scope gap must be reconciled and owner-approved before publication.

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
| Release identity and immutable tag | Partial | Baseline captured at `7903612`; manual version/tag identity and tag are still pending |
| Plan and scope completion | Partial | Implementation slices merged; complete apply orchestration remains unresolved |
| Contract and security completion | Partial | Contracts exist; named release security/architecture review pending |
| Automated verification | Complete for merged baseline slice | PR #331 CI passed on Linux/macOS/Windows; reviewed `origin/main` binary smoke passed locally |
| Manual verification | Open | Scan/baseline/diff workflow passed on macOS arm64; representative remediation smoke pending |
| Artifact and supply-chain verification | Open | Build artifacts, checksums, SBOM, signing, provenance pending |
| Publication and approval | Open | No release publication authorized |
| Post-release closure | Open | Retrospective cannot be completed before release decision |

## Required next evidence

- [x] Capture release baseline: `origin/main`, tags, releases, and manifest status. Manifest is not required for the owner-run manual release process.
- [x] Run `make verify` from a clean reviewed checkout. Result: passed on macOS arm64 from baseline branch; `go generate`, `go vet`, `go test`, and `go build` completed successfully.
- [x] Run representative repository scan and baseline/diff smoke; macOS arm64, source `7903612`, binary `/tmp/deprail-v0.4.0-main`, scan complete with 14 findings, baseline create exit `0`, identical diff exit `0`.
- [ ] Reconcile complete `deprail fix apply` scope and compatibility impact.
- [ ] Record platform, binary identity, artifacts, sizes, and SHA-256 checksums.
- Manual evidence artifacts: `local_test/v0.4.0-preview.1/output/scan.json`, `baseline-create.json`, `diff.json`; baseline artifact is under the ignored fixture `.deprail/manual` directory.
- [ ] Record SBOM, signing, and provenance status.
- [ ] Obtain named architecture/security review.
- [ ] Prepare preview notes and rollback procedure.
- [ ] Record owner go/no-go decision.
- [ ] Create and link retrospective after publication.

## Review and decision

- Owner decision: `[pending]`
- Decision date: `[pending]`
- Architecture/security review: `[pending]`
- Rollback owner: `[pending]`
- Rollback procedure: preserve evidence, do not move immutable tags, and publish a new preview only after resolving approved gaps.
