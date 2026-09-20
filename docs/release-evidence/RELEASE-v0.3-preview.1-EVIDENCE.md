# DepRail v0.3 Release Evidence Record

**Status:** Preview published; final owner decision pending.
**Release mode:** Preview
**Target version:** `0.3.0-preview.1`
**Release owner:** `[owner to confirm]`
**Reviewer:** `[security/architecture reviewer to confirm]`
**Decision:** `[GO / GO WITH APPROVED GAPS / NO-GO — owner to confirm]`

Release workflow and artifact publication completed successfully. The owner decision and security/architecture review remain open.

Missing or inconsistent evidence blocks the applicable release mode. Do not replace missing evidence with an unchecked claim.

## Source and identity

| Field | Evidence |
| --- | --- |
| Source commit | `72fc603bd8b9c852c2accadd9521ea22a9720991` |
| Release tag | [`v0.3.0-preview.1`](https://github.com/geoffrey-xiao/deprail/releases/tag/v0.3.0-preview.1) |
| DepRail version | `v0.3.0-preview.1` |
| Release workflow | [Run 35488833192](https://github.com/geoffrey-xiao/deprail/actions/runs/35488833192) |
| Host platform | Release artifacts built on GitHub Actions Linux runner |
| Local verification platform | macOS arm64 / Darwin |
| Go version | `go1.27.1 darwin/arm64` locally; release workflow used `.go-version` |
| OSV-Scanner | `2.6.0`; OSV-Scalibr `0.5.2` locally |
| Local binary SHA-256 | `0e3c5f8260d8a2a2ba546d270b4c69fc813d971efa321c93f772bfd225c980af` |
| Tracking issue | [#215](https://github.com/geoffrey-xiao/deprail/issues/215) |
| Local artifact directory | `local_test/`; ignored and retained for archival |

## Pre-release gate

The following commands ran from the repository root using the binary built from the source commit above. The exact captured output is under `local_test/`.

| Gate | Result | Exit code | Evidence |
| --- | --- | ---: | --- |
| `go build -o deprail ./cmd/deprail` | Passed | 0 | `local_test/preview-identity.txt` |
| `./deprail --version` | Passed locally; release identity verified by workflow artifacts | 0 | `local_test/preview-identity.txt`; release workflow |
| `./deprail doctor --format json` | Passed | 0 | `local_test/doctor-preview.json`; release smoke |
| Mixed repository scan | Passed; complete, 14 findings, no errors | 0 | `local_test/mixed-repo-preview.json` |
| `fix plan` from real scan finding | Passed; schema `v0alpha1` | 0 | `local_test/fix-plan-preview.json` |
| `make verify` from reviewed `origin/main` | Passed | 0 | Release preparation verification |
| Release workflow | Passed | 0 | [Run 35488833192](https://github.com/geoffrey-xiao/deprail/actions/runs/35488833192) |

The release workflow ran from the tagged reviewed commit and completed both build and protected publication jobs.

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

The release matrix was built, signed, attested, checksum-verified, and published by the protected release workflow.

| Artifact | Platform | SHA-256 | Smoke result | Evidence |
| --- | --- | --- | --- | --- |
| `deprail-linux-amd64` | Linux amd64 | `d42748daa9fc1ad06d70ce95bab619e7a610c11ad9259ae1334a346a78898295` | Release smoke passed | [`SHA256SUMS`](https://github.com/geoffrey-xiao/deprail/releases/download/v0.3.0-preview.1/SHA256SUMS) |
| `deprail-darwin-amd64` | macOS amd64 | `1e01b61e61e53472c904b06d3b7982e815e1e2a0df1914a96a04a053b78b001a` | Built and checksum-verified | [`SHA256SUMS`](https://github.com/geoffrey-xiao/deprail/releases/download/v0.3.0-preview.1/SHA256SUMS) |
| `deprail-darwin-arm64` | macOS arm64 | `d4c5c580036e51c82017b6b7474c8840a59c9a5f023b8385191563173761baaa` | Built and checksum-verified; local workflow also passed | [`SHA256SUMS`](https://github.com/geoffrey-xiao/deprail/releases/download/v0.3.0-preview.1/SHA256SUMS) |
| `deprail-windows-amd64.exe` | Windows amd64 | `0c132805c4b4857ddd77a01b3ae2cba204fce2530969dba9ec7991cd0c853c7a` | Built and checksum-verified | [`SHA256SUMS`](https://github.com/geoffrey-xiao/deprail/releases/download/v0.3.0-preview.1/SHA256SUMS) |
| Checksum manifest | All artifacts | — | `sha256sum -c SHA256SUMS`: all OK | [`SHA256SUMS`](https://github.com/geoffrey-xiao/deprail/releases/download/v0.3.0-preview.1/SHA256SUMS) |

## Platform and CI evidence

| Platform | Evidence | Status |
| --- | --- | --- |
| Linux amd64 | Release workflow build, smoke, checksum, signature, and provenance verification | Passed |
| macOS amd64 | Release workflow build, checksum, signature, and provenance verification | Passed |
| macOS arm64 | Release workflow build, checksum, signature, and provenance verification; local smoke passed | Passed |
| Windows amd64 | Release workflow build, checksum, signature, and provenance verification | Passed |

The release workflow completed successfully: [Run 35488833192](https://github.com/geoffrey-xiao/deprail/actions/runs/35488833192).

## Supply-chain controls

| Control | Status | Evidence |
| --- | --- | --- |
| SBOM | Generated and published | [`deprail-v0.3.0-preview.1.spdx.json`](https://github.com/geoffrey-xiao/deprail/releases/download/v0.3.0-preview.1/deprail-v0.3.0-preview.1.spdx.json) |
| Signing | Generated and verified for release artifacts | `.sig` and `.pem` assets in the [GitHub Release](https://github.com/geoffrey-xiao/deprail/releases/tag/v0.3.0-preview.1) |
| Provenance | Generated and verified by the release workflow | [Run 35488833192](https://github.com/geoffrey-xiao/deprail/actions/runs/35488833192) |

A preview may carry explicitly approved gaps. These gaps do not approve stable `v0.3.0`.

## Risks, rollback, and decision

- Known limitations: preview-only, OSV-Scanner must be installed separately, no package-manager mutation, and no automatic remediation.
- Remaining release risks: owner decision and security/architecture review remain open; stable release approval is not implied.
- Last known-good tag: `v0.3.0-preview.1`
- Rollback owner: `[owner to confirm]`
- Rollback procedure: withdraw the preview reference, preserve artifacts and logs, restore the last known-good tag, create a new immutable preview tag, and rerun the applicable checklist.
- Owner go/no-go decision: `[owner to confirm]`
- Decision date: `[owner to confirm]`

## Evidence review checklist

- [x] Source commit and preview tag are recorded and match the build.
- [x] All four release artifacts have checksums and workflow verification.
- [x] Representative mixed-repository scan and plan evidence is attached locally.
- [x] Read-only and unsafe-output behavior is recorded.
- [x] Clean-checkout, release-workflow, SBOM, signing, and provenance evidence is attached.
- [ ] Owner decision is recorded.
- [ ] Security/architecture review is recorded.
