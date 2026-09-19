# DepRail v0.2.0-preview.2 Release Evidence

**Status:** Artifact evidence collected; platform smoke evidence remains incomplete.
**Release mode:** Preview
**Target version:** `v0.2.0-preview.2`
**Release owner:** `@geoffrey-xiao`
**Reviewer:** `@geoffreyxiaoai`
**Rollback owner:** `@geoffrey-xiao`
**Decision:** `GO WITH APPROVED GAPS`
**Decision date:** `2026-09-19`

This record is immutable evidence for `v0.2.0-preview.2` only. Do not mix checksums or workflow results from preview.1.

## Source and identity

| Field | Evidence |
| --- | --- |
| Source commit | `05189eacdbda433eec9dc451cc3f5544bac32608` |
| Release tag | `v0.2.0-preview.2` |
| DepRail version | `v0.2.0-preview.2` |
| CLI identity output | Release artifact build embeds tag and commit; local CLI output pending |
| Release workflow run | https://github.com/geoffrey-xiao/deprail/actions/runs/35436690490 |
| Manifest value | `0.2.0-preview.2` |

## Artifact inventory and checksums

| Artifact | Platform | Size | SHA-256 | Smoke result | Evidence |
| --- | --- | ---: | --- | --- | --- |
| `deprail-linux-amd64` | Linux amd64 | 3,760,288 | `ec3da498e07a2b1bac63e5d5682beb753478822bffde6d8a1e0552a0e226f72e` | Doctor/discover smoke passed | Release workflow |
| `deprail-darwin-amd64` | macOS amd64 | 3,761,376 | `ed86ff1891db59ffce20ef98fab72d44386706f791f2fc065d3efe79647b143c` | Artifact checksum verified; runtime smoke pending | Release assets |
| `deprail-darwin-arm64` | macOS arm64 | 3,514,258 | `54fb08db4e7a5bb3ecdcc8610290dd61bf7514368e10d0b46994b8b69fd17a82` | Artifact checksum verified; runtime smoke pending | Release assets |
| `deprail-windows-amd64.exe` | Windows amd64 | 3,904,512 | `28d3f9b44717c9f423cc98a06821b11c8b1cba9a3b0aeb0d8626ac4555ed6625` | Artifact checksum verified; runtime smoke pending | Release assets |
| `SHA256SUMS` | All artifacts | — | `shasum -a 256 -c SHA256SUMS`: all OK | Checksum verification passed | Release workflow |

## Platform smoke

| Platform | Doctor | Discover | Scan failure behavior | Identity match | Evidence |
| --- | --- | --- | --- | --- | --- |
| Linux | Doctor/discover passed | Passed | Not run | Tag/commit injected | https://github.com/geoffrey-xiao/deprail/actions/runs/35436690490 |
| macOS | Pending | Pending | Pending | Build identity embedded | Release assets; runtime smoke needed |
| Windows | Pending | Pending | Pending | Build identity embedded | Release assets; runtime smoke needed |

## Scanner and representative repositories

| Evidence | Result | Link or output |
| --- | --- | --- |
| OSV-Scanner version | `[pending]` | `[record output]` |
| JavaScript repository scan | `[pending]` | `[evidence]` |
| Python repository scan | `[pending]` | `[evidence]` |
| Java repository scan | `[pending]` | `[evidence]` |
| Mixed repository scan | `[pending]` | `[evidence]` |

## Supply-chain controls

| Control | Status | Evidence or gap rationale |
| --- | --- | --- |
| SBOM | Not implemented for preview | Approved preview gap; stable release requires linked SBOM evidence |
| Signing | Not implemented for preview | Approved preview gap; no signed-artifact claim is made |
| Provenance | Formal attestation not published | Approved preview gap; workflow URL alone is not an attestation |

## Risks and rollback

- Known limitations: preview release; caller-installed OSV-Scanner; no remote baseline history; no automatic remediation or source mutation.
- Remaining release risks: SBOM, signing, and formal provenance remain unavailable for this preview; stable release is not approved.
- Rollback procedure: withdraw the preview reference, preserve artifacts and logs, restore the last known-good tag, create a new immutable preview tag, and rerun this checklist.

## Review checklist

- [ ] Source commit and tag match the build.
- [ ] All artifacts have checksums.
- [ ] Every supported platform has attached smoke evidence.
- [ ] Scanner and representative scan evidence are recorded.
- [x] SBOM, signing, and provenance gaps are explicit.
- [x] Release owner, reviewer, rollback owner, decision, and date are recorded.
