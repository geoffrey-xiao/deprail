# DepRail v0.3 Release Evidence Record

**Status:** Preview evidence collected; owner decision pending.
**Release mode:** Preview
**Target version:** `0.3.0-preview.1`
**Release owner:** `[owner to confirm]`
**Reviewer:** `[security/architecture reviewer to confirm]`
**Decision:** `[GO / GO WITH APPROVED GAPS / NO-GO — owner to confirm]`

Missing or inconsistent evidence blocks the applicable release mode. Do not replace missing evidence with an unchecked claim.

## Source and identity

| Field | Evidence |
| --- | --- |
| Source commit | `2ef17f3b6932fd41001465870416059af0d40662` |
| Release tag | `Not created; candidate preview evidence only` |
| DepRail version | Development build: `tag=unknown commit=unknown` |
| Host platform | macOS arm64 / Darwin |
| Go version | `go1.27.1 darwin/arm64` |
| OSV-Scanner | `2.6.0`; OSV-Scalibr `0.5.2` |
| Binary SHA-256 | `0e3c5f8260d8a2a2ba546d270b4c69fc813d971efa321c93f772bfd225c980af` |
| Tracking issue | [#215](https://github.com/geoffrey-xiao/deprail/issues/215) |
| Local artifact directory | `local_test/`; ignored and retained for attachment or archival |

## Pre-release gate

The following commands ran from the repository root using the binary built from the source commit above. The exact captured output is under `local_test/`.

| Gate | Result | Exit code | Evidence |
| --- | --- | ---: | --- |
| `go build -o deprail ./cmd/deprail` | Passed | 0 | `local_test/preview-identity.txt` |
| `./deprail --version` | Passed; development identity reported | 0 | `local_test/preview-identity.txt` |
| `./deprail doctor --format json` | Passed | 0 | `local_test/doctor-preview.json` |
| Mixed repository scan | Passed; complete, 14 findings, no errors | 0 | `local_test/mixed-repo-preview.json` |
| `fix plan` from real scan finding | Passed; schema `v0alpha1` | 0 | `local_test/fix-plan-preview.json` |
| `go test ./...` | Passed | 0 | `local_test/go-test.txt` |
| `go vet ./...` | Passed | 0 | `local_test/go-vet.txt` |
| `git diff --check` | Passed | 0 | `local_test/git-diff-check.txt` |

A clean detached checkout and `make verify` have not yet been recorded for the preview candidate.

## Representative repository workflow

Workflow: repository discovery/scan -> OSV-Scanner -> scan JSON -> read-only remediation plan.

| Evidence | Result |
| --- | --- |
| Fixture | `testdata/fixtures/mixed-repository` |
| Scan | Exit `0`; status `complete`; 14 findings; errors `[]` |
| Selected finding | `GHSA-29mw-wpgm-hmr9` / `lodash@4.17.20` |
| Ecosystem identity | PURL `pkg:npm/lodash@4.17.20`; ecosystem `npm` |
| Workspace identity | ID/path `frontend` |
| Plan | Exit `0`; plan ID `8fa59ebf6c03a4d557d29528705b41a48c3da6cbae8d0108acc93a64ea06bfa3` |
| Affected files | `frontend/package.json`; `frontend/package-lock.json` |
| Candidate | `lodash@4.17.21`, viable with unknown compatibility metadata preserved |
| Provenance | Three artifact digests; sources `javascript`, `normalized-scan-report` |

Detailed artifacts: `local_test/scan-summary.json`, `local_test/finding-identity.json`, and `local_test/plan-summary.json`.

## Output safety and read-only behavior

| Scenario | Result | Evidence |
| --- | --- | --- |
| External plan output | Passed, exit `0` | `local_test/fix-plan-external.json` |
| Repeated output to existing path | Rejected, exit `3`, `PLAN_WRITE_FAILED` | `local_test/repeat-stderr.txt` |
| Repository-internal output | Rejected, exit `2`, `PLAN_OUTPUT_OUTSIDE_ROOT_REQUIRED` | `local_test/internal-stderr.txt` |
| Traversal output | Rejected, exit `2`, `PLAN_OUTPUT_OUTSIDE_ROOT_REQUIRED` | `local_test/traversal-stderr.txt` |
| Fixture tree comparison | Unchanged | `local_test/tree-diff-comparison.txt` |

## Artifact inventory and checksums

Only the local macOS arm64 development binary has been built and smoke-tested in this evidence run. The release matrix is incomplete.

| Artifact | Platform | SHA-256 | Smoke result | Evidence |
| --- | --- | --- | --- | --- |
| `deprail` | macOS arm64 | `0e3c5f8260d8a2a2ba546d270b4c69fc813d971efa321c93f772bfd225c980af` | Doctor, scan, and plan passed | `local_test/deprail.sha256`; artifacts above |
| `deprail-linux-amd64` | Linux amd64 | Not built | Missing | Requires candidate release build |
| `deprail-darwin-amd64` | macOS amd64 | Not built | Missing | Requires candidate release build |
| `deprail-windows-amd64.exe` | Windows amd64 | Not built | Missing | Requires candidate release build |
| Checksum manifest | All artifacts | Not generated | Missing | Requires candidate release build |

## Platform and CI evidence

| Platform | Evidence | Status |
| --- | --- | --- |
| macOS arm64 | Local real-binary run from candidate commit | Passed |
| Linux | No candidate binary smoke attached | Missing |
| Windows | No candidate binary smoke attached | Missing |

PR #231 CI passed on Linux, macOS, and Windows, but it verified documentation changes and is not binary artifact smoke evidence for this preview candidate.

## Supply-chain controls

| Control | Status | Gap rationale |
| --- | --- | --- |
| SBOM | Not generated | Preview evidence does not include a release artifact SBOM |
| Signing | Not implemented/verified | No signed-artifact claim is made |
| Provenance | Source commit and local checksum recorded; attestation not published | Release workflow evidence remains required |

A preview may carry explicitly approved gaps. These gaps do not approve stable `v0.3.0`.

## Risks, rollback, and decision

- Known limitations: preview-only, OSV-Scanner must be installed separately, no package-manager mutation, no automatic remediation, and only macOS arm64 binary smoke is attached.
- Remaining release risks: Linux/Windows artifact evidence, full release matrix, clean-checkout gate, SBOM/signing/provenance, owner decision, and security/architecture review.
- Last known-good tag: `[owner to confirm]`
- Rollback owner: `[owner to confirm]`
- Rollback procedure: withdraw the preview reference, preserve artifacts and logs, restore the last known-good tag, create a new immutable preview tag, and rerun the applicable checklist.
- Owner go/no-go decision: `[owner to confirm]`
- Decision date: `[owner to confirm]`

## Evidence review checklist

- [ ] Source commit and preview tag are recorded and match the build.
- [x] Local macOS arm64 binary has a checksum and smoke result.
- [ ] Linux, macOS amd64, and Windows artifacts have checksums and smoke results.
- [x] Representative mixed-repository scan and plan evidence is attached locally.
- [x] Read-only and unsafe-output behavior is recorded.
- [ ] Clean-checkout and release-workflow evidence is attached.
- [ ] SBOM, signing, and provenance status are dispositioned by the owner.
- [ ] Remaining risks, rollback owner, and release decision are recorded.
- [ ] Security/architecture review is recorded.
