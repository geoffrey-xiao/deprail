# DepRail v0.2 Release Evidence Record

**Status:** Preview evidence collected; owner decision pending.
**Release mode:** Preview
**Target version:** `0.2.0-preview.1`
**Release owner:** `[not recorded]`
**Reviewer:** `[not recorded]`
**Decision:** `[not recorded]`

Missing or inconsistent evidence blocks the applicable release mode. Do not replace missing evidence with an unchecked claim.

## Source and identity

| Field | Evidence |
| --- | --- |
| Source commit | `28f4e3ea192d9d2c75f66a079e27ec54d29f28b1` |
| Release tag | `v0.2.0-preview.1` |
| DepRail version | `v0.2.0-preview.1` |
| CLI identity output | macOS Intel and Apple Silicon `doctor` smoke |
| Build command and Go version | Release workflow build; Go version recorded in workflow |
| Release workflow run | https://github.com/geoffrey-xiao/deprail/actions/runs/35433949073 |

## Pre-release gate

Record each gate from a clean detached checkout of the reviewed source commit. A missing result or non-zero exit code blocks the applicable release mode.

| Gate | Result | Exit code | Evidence |
| --- | --- | ---: | --- |
| Clean detached checkout at reviewed commit | Passed in release workflow | 0 | Release workflow |
| `go version` | Passed in release workflow | 0 | Release workflow |
| `make verify` | Passed | 0 | Release workflow |
| `git diff --check` | Passed locally | 0 | Local validation |

Attach command output or a CI run link for each row. The record is incomplete until every gate has an explicit result and exit code.

## Artifact inventory and checksums

Every artifact requires a recorded SHA-256 checksum and a smoke result from the same reviewed source commit.

| Artifact | Platform | Size | SHA-256 | Smoke result | Evidence |
| --- | --- | ---: | --- | --- | --- |
| `deprail-linux-amd64` | Linux amd64 | 3,760,288 | `ec3da498e07a2b1bac63e5d5682beb753478822bffde6d8a1e0552a0e226f72e` | Release smoke passed | Release workflow |
| `deprail-darwin-amd64` | macOS amd64 | 3,761,376 | `ed86ff1891db59ffce20ef98fab72d44386706f791f2fc065d3efe79647b143c` | Manual doctor passed | Release workflow / local |
| `deprail-darwin-arm64` | macOS arm64 | 3,514,258 | `54fb08db4e7a5bb3ecdcc8610290dd61bf7514368e10d0b46994b8b69fd17a82` | Manual smoke passed | Release workflow / owner report |
| `deprail-windows-amd64.exe` | Windows amd64 | 3,904,512 | `28d3f9b44717c9f423cc98a06821b11c8b1cba9a3b0aeb0d8626ac4555ed6625` | Release smoke passed | Release workflow |
| Checksum manifest | All artifacts | — | — | `shasum -a 256 -c SHA256SUMS`: all OK | Preview release |

## Platform and CLI smoke evidence

For every supported artifact, attach command output and the CI or manual-run link:

- `deprail doctor --format json`
- discovery against `testdata/fixtures/mixed-repository`
- scan behavior when OSV-Scanner is unavailable, including exit code `3`
- reported version, tag, and commit identity

| Platform | CI or run link | Doctor | Discover | Scan failure behavior | Identity match |
| --- | --- | --- | --- | --- | --- |
| Linux | https://github.com/geoffrey-xiao/deprail/actions/runs/35433949073 | Release smoke passed | Workflow | Workflow | Tag/commit injected |
| macOS | Local smoke plus release workflow | Passed | Passed | Passed | `darwin/amd64` and `darwin/arm64` verified |
| Windows | https://github.com/geoffrey-xiao/deprail/actions/runs/35433949073 | Release smoke passed | Workflow | Workflow | Tag/commit injected |

## Scanner and representative repositories

| Evidence | Result | Link or output |
| --- | --- | --- |
| OSV-Scanner version | `2.6.0` / `0.5.2` | Local macOS and release environment |
| JavaScript repository scan | Complete fixture workspace | PR #182 / local smoke |
| Python repository scan | Complete fixture workspace | PR #182 / local smoke |
| Java repository scan | Complete fixture workspace | PR #182 / local smoke |
| Mixed repository scan | Exit `0`, complete, 9 findings, no errors | PR #182 |

Every release-gating scan must complete successfully. Partial or failed scans block the relevant release mode.
## V02-025 pull-request validation evidence

| Evidence | Result | Link or output |
| --- | --- | --- |
| Controlled mixed-repository fixture | Complete discovery; 3 workspaces | `testdata/fixtures/mixed-repository` |
| Complete scan | Exit `0`; 9 findings; no scanner errors | Local transcript recorded during V02-025 validation |
| Cross-platform CI | Ubuntu, macOS, and Windows passed | PR #182 |
| Fixture pull-request change | Merged | [PR #182](https://github.com/geoffrey-xiao/deprail/pull/182) |

The fixture intentionally retains vulnerable dependencies so findings and policy/SARIF paths remain testable. A clean scan result is not claimed from this fixture; the recorded complete result means all detected workspaces and scanner executions completed without scanner errors.


## Supply-chain controls

| Control | Status | Evidence or gap rationale |
| --- | --- | --- |
| SBOM | Not implemented unless linked evidence is supplied | `[not recorded]` |
| Signing | Not implemented unless linked evidence is supplied | `[not recorded]` |
| Provenance | Not implemented unless linked evidence is supplied | `[not recorded]` |

A preview may carry explicitly approved gaps. An RC with missing SBOM, signing, or provenance evidence is not stable-ready. A stable release must have evidence for each control; otherwise the stable release is explicitly rejected and the missing control is recorded as a release blocker.

## Risks, rollback, and decision

- Known limitations: `[not recorded]`
- Remaining release risks: `[not recorded]`
- Last known-good tag: `[not recorded]`
- Rollback owner: `[not recorded]`
- Rollback procedure: stop publication, preserve evidence, create a new immutable corrective tag, and rerun the applicable mode checklist.
- Owner go/no-go decision: `[not recorded]`
- Decision date: `[not recorded]`

## Evidence review checklist

- [ ] Source commit and tag are recorded and match the build.
- [ ] Every required artifact has a checksum and smoke result.
- [ ] Platform CI links and manual scan evidence are attached.
- [ ] Scanner version and compatibility are recorded.
- [ ] SBOM, signing, and provenance status are explicit.
- [ ] Remaining risks and rollback owner are recorded.
- [ ] Release/security reviewer decision is recorded.
