# DepRail v0.3.1 Preview.1 Release Evidence Record

**Status:** Draft release record; preview not published.
**Release mode:** Preview
**Target version:** `0.3.1-preview.1`
**Release owner:** `@geoffrey-xiao` — owner confirmation pending
**Reviewer:** Architecture/security reviewer — assignment pending
**Decision:** `GO / GO WITH APPROVED GAPS / NO-GO` — owner decision pending

This record is an evidence index, not a release approval. Unchecked items are open gaps. A preview may proceed only with explicit owner approval of documented gaps; this record must not be read as a stable-release claim.

Reusable gate: [`docs/RELEASE-CHECKLIST.md`](../RELEASE-CHECKLIST.md)

## Source and identity

| Field | Evidence |
| --- | --- |
| Reviewed source commit | `94ab2f0dcfb3b983f627b917da65ab5ee6480463` at record creation |
| Release tag | Not created; preview tag must be new and immutable |
| DepRail version | Current manifest remains `0.3.0-preview.1`; update only through the approved release process |
| Release workflow | Not run for `v0.3.1-preview.1` |
| Local verification platform | macOS arm64 / Darwin |
| Local Go version | `go1.27.1 darwin/arm64` |
| OSV-Scanner | Record exact installed version during release verification |
| Local artifact SHA-256 | Not generated for a release artifact |
| Tracking issue | [#253](https://github.com/geoffrey-xiao/deprail/issues/253) |
| Project | DepRail project, v0.3.1 release work |

## Scope and release boundary

v0.3.1 is a presentation-focused preview. Included scope:

- shared human presentation primitives;
- lifecycle events and truthful progress/outcome messages;
- actionable command help tips and examples;
- JSON stdout purity and stderr diagnostics separation;
- quiet/verbose behavior;
- hostile-label rendering safeguards;
- cross-platform presentation verification where evidence is available.

Deferred from this preview:

- semantic color integration into lifecycle messages ([#281](https://github.com/geoffrey-xiao/deprail/issues/281));
- full-screen terminal UI;
- repository mutation or remediation execution;
- scanner, ecosystem, schema, network, and publication changes outside the approved v0.3.1 package.

## Current verification evidence

The following smoke checks were run from the reviewed local checkout on macOS arm64:

| Gate | Result | Exit code | Evidence |
| --- | --- | ---: | --- |
| `go build -o /tmp/deprail-v031 ./cmd/deprail` | Passed | 0 | Local command run |
| `/tmp/deprail-v031 --help` | Passed | 0 | Local command run |
| All six command help pages | Passed; six `Tip:` and six `Example:` entries observed | 0 | Local command run |
| `discover testdata/fixtures/mixed-repository --format json` | Passed; no ANSI observed | 0 | Local command run |
| `doctor --format json` | Passed; no ANSI observed | 0 | Local command run |
| `scan testdata/fixtures/mixed-repository --format json` | Passed; no ANSI observed | 0 | Local command run |
| `go test ./...` | Passed on recent implementation PRs | 0 | CI and PR evidence |
| `git diff --check` | Passed on recent implementation PRs | 0 | CI and PR evidence |

The local commands above were smoke verification only. They do not replace the complete release checklist or artifact workflow.

## Implementation and CI evidence

| Area | Evidence | Status |
| --- | --- | --- |
| Lifecycle start notices | PR [#276](https://github.com/geoffrey-xiao/deprail/pull/276) | Merged |
| Actionable help tips | PR [#279](https://github.com/geoffrey-xiao/deprail/pull/279) | Merged |
| Safe color primitives | PR [#278](https://github.com/geoffrey-xiao/deprail/pull/278) | Merged; integration deferred |
| Command UX verification | Issue [#274](https://github.com/geoffrey-xiao/deprail/issues/274) | Owner-confirmed complete |
| Cross-platform CI | Required CI runs on implementation PRs | Passed where linked; aggregate release evidence pending |
| Cross-platform interactive TTY evidence | Issue [#254](https://github.com/geoffrey-xiao/deprail/issues/254) | Open gap unless separately captured |

## General release checklist reconciliation

| Checklist area | Status | Evidence or gap |
| --- | --- | --- |
| Release identity and immutable tag | Open | Version manifest and tag are not yet prepared |
| Plan and scope completion | Partial | v0.3.1 plan and implementation PRs exist; release disposition remains open |
| Contract and security completion | Partial | Runtime contracts are implemented; architecture/security review is not recorded |
| Automated verification | Partial | CI and local smoke evidence exist; clean release-checkout `make verify` record pending |
| Manual verification | Partial | Owner confirmed command UX verification; full release scenario bundle and platform captures pending |
| Artifact and supply-chain verification | Open | Release artifacts, checksums, SBOM, signing, and provenance not generated for this preview |
| Publication and approval | Open | No preview release workflow, notes, approval, or go/no-go decision yet |
| Post-release closure | Open | Retrospective and follow-up evidence are not yet prepared |

## Required release-gate work

Before publishing `v0.3.1-preview.1`:

- [ ] Confirm release mode, owner, and named architecture/security reviewer.
- [ ] Reconcile the v0.3.1 development plan, execution package, traceability, compatibility matrix, and master checklist.
- [ ] Update the version manifest through the approved release process.
- [ ] Run `make verify` from the reviewed clean checkout and record the result.
- [ ] Complete required Linux, macOS, and Windows artifact/smoke evidence.
- [ ] Record interactive, redirected, JSON, quiet, verbose, partial, failed, narrow, monochrome, hostile-label, and cancellation evidence or explicitly disposition each gap.
- [ ] Build artifacts from the reviewed immutable commit and record names, sizes, and SHA-256 checksums.
- [ ] Record SBOM, signing, and provenance status.
- [ ] Prepare release notes, rollback ownership, and immutable-tag recovery procedure.
- [ ] Obtain architecture/security review and owner go/no-go decision.

## Known risks and deferred work

- Interactive TTY evidence is not represented as complete in this record.
- The preview has not been published and has no release artifact identity yet.
- OSV-Scanner remains an external prerequisite; automatic installation is out of scope.
- Color integration is deferred to the backlog and is not required for this preview unless the owner changes scope.
- A preview approval must not be interpreted as stable compatibility approval.

## Decision and review

- Owner go/no-go decision: `[pending]`
- Decision date: `[pending]`
- Architecture/security review: `[pending]`
- Rollback owner: `[pending]`
- Rollback procedure: preserve the evidence and artifacts, withdraw the preview reference if necessary, and create a new immutable preview tag after resolving the documented gap.

## Evidence index

- Release checklist: [`docs/RELEASE-CHECKLIST.md`](../RELEASE-CHECKLIST.md)
- Development plan: [`docs/03-planning/deprail-development-plan-v0.3.1.md`](../03-planning/deprail-development-plan-v0.3.1.md)
- Master checklist: [`docs/04-execution/deprail-v0.3.1-execution-package/tracking/MASTER-CHECKLIST.md`](../04-execution/deprail-v0.3.1-execution-package/tracking/MASTER-CHECKLIST.md)
- Release gate issue: [#253](https://github.com/geoffrey-xiao/deprail/issues/253)
- Cross-platform evidence issue: [#254](https://github.com/geoffrey-xiao/deprail/issues/254)
- Local ignored captures: `local_test/` when collected; do not publish credentials or internal logs.
