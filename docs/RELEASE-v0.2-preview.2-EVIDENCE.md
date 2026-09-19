# DepRail v0.2.0-preview.2 Release Evidence

**Status:** Evidence pending
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
| Source commit | `[record after merge]` |
| Release tag | `v0.2.0-preview.2` |
| DepRail version | `v0.2.0-preview.2` |
| CLI identity output | `[attach artifact output]` |
| Release workflow run | `[attach workflow URL]` |
| Manifest value | `0.2.0-preview.2` |

## Artifact inventory and checksums

| Artifact | Platform | Size | SHA-256 | Smoke result | Evidence |
| --- | --- | ---: | --- | --- | --- |
| `deprail-linux-amd64` | Linux amd64 | `[pending]` | `[pending]` | `[pending]` | `[workflow URL]` |
| `deprail-darwin-amd64` | macOS amd64 | `[pending]` | `[pending]` | `[pending]` | `[manual output or URL]` |
| `deprail-darwin-arm64` | macOS arm64 | `[pending]` | `[pending]` | `[pending]` | `[manual output or URL]` |
| `deprail-windows-amd64.exe` | Windows amd64 | `[pending]` | `[pending]` | `[pending]` | `[manual output or URL]` |
| `SHA256SUMS` | All artifacts | — | `[pending]` | `[pending]` | `[release URL]` |

## Platform smoke

| Platform | Doctor | Discover | Scan failure behavior | Identity match | Evidence |
| --- | --- | --- | --- | --- | --- |
| Linux | `[pending]` | `[pending]` | `[pending]` | `[pending]` | `[URL]` |
| macOS | `[pending]` | `[pending]` | `[pending]` | `[pending]` | `[URL/output]` |
| Windows | `[pending]` | `[pending]` | `[pending]` | `[pending]` | `[URL/output]` |

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
