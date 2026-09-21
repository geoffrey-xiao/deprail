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
| Reviewed source commit | Pending release baseline capture |
| Release manifest | Pending |
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
- cross-platform relative-path semantics.

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

## Checklist reconciliation

| Area | Status | Evidence or gap |
|---|---|---|
| Release identity and immutable tag | Open | Baseline, manifest, and tag not captured |
| Plan and scope completion | Partial | Implementation slices merged; complete apply orchestration remains unresolved |
| Contract and security completion | Partial | Contracts exist; named release security/architecture review pending |
| Automated verification | Partial | Implementation PR CI passed; clean-checkout `make verify` pending |
| Manual verification | Open | Representative repository remediation smoke pending |
| Artifact and supply-chain verification | Open | Build artifacts, checksums, SBOM, signing, provenance pending |
| Publication and approval | Open | No release publication authorized |
| Post-release closure | Open | Retrospective cannot be completed before release decision |

## Required next evidence

- [ ] Capture release baseline: `origin/main`, tags, releases, manifest.
- [ ] Reconcile complete `deprail fix apply` scope and compatibility impact.
- [ ] Run `make verify` from a clean reviewed checkout.
- [ ] Run representative repository smoke coverage and record exit codes.
- [ ] Record platform, binary identity, artifacts, sizes, and SHA-256 checksums.
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
