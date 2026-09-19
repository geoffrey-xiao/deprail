# DepRail v0.2.0-preview.3 Release Evidence

**Status:** Preview release preparation; final evidence pending.
**Release mode:** Preview
**Target version:** `v0.2.0-preview.3`
**Release owner:** `@geoffrey-xiao`
**Reviewer:** `@geoffreyxiaoai`
**Rollback owner:** `@geoffrey-xiao`
**Decision:** `GO WITH APPROVED GAPS`
**Decision date:** `2026-09-19`

This record applies only to the immutable `v0.2.0-preview.3` tag. Preview 1 and preview 2 evidence remain in their own records.

## Release baseline

| Field | Evidence |
| --- | --- |
| Latest reviewed preview | `v0.2.0-preview.2` |
| Reviewed source commit | `[record after preview tag]` |
| Preview release tag | `v0.2.0-preview.3` |
| Manifest value | `0.2.0-preview.3` |
| Release workflow run | `[record after preview publication]` |
| Owner smoke report | Preview 2 macOS and Windows smoke reported passed by owner; preview 3 artifact output pending |

## Preview security controls

The release workflow generates and verifies the controls in the build job. Publication is blocked at the protected `release-approval` environment until an authorized reviewer presses the GitHub Actions approval button. The approval job then verifies checksums and uploads the already-reviewed assets.

| Control | Status | Evidence or disposition |
| --- | --- | --- |
| SBOM | Implemented in CI | SPDX JSON generated and attached after workflow verification |
| Signing | Implemented in CI | Keyless Cosign signatures and certificates generated and verified |
| Provenance | Implemented in CI | GitHub build attestations generated and verified |

## Artifact inventory and checksums

| Artifact | Platform | Size | SHA-256 | Smoke result | Evidence |
| --- | --- | ---: | --- | --- | --- |
| `deprail-linux-amd64` | Linux amd64 | `[pending]` | `[pending]` | `[pending]` | `[workflow URL]` |
| `deprail-darwin-amd64` | macOS amd64 | `[pending]` | `[pending]` | `[pending]` | `[owner output or URL]` |
| `deprail-darwin-arm64` | macOS arm64 | `[pending]` | `[pending]` | `[pending]` | `[owner output or URL]` |
| `deprail-windows-amd64.exe` | Windows amd64 | `[pending]` | `[pending]` | `[pending]` | `[owner output or URL]` |
| `SHA256SUMS` | All artifacts | — | `[pending]` | `[pending]` | `[release URL]` |

## Preview acceptance checklist

- [ ] Preview source commit and tag match the manifest and CLI identity.
- [ ] `make verify` passes from a clean checkout.
- [ ] All artifacts have checksums.
- [ ] Linux, macOS, and Windows artifact smoke output is attached.
- [ ] Representative JavaScript, Python, Java, and mixed scans are attached.
- [x] SBOM generation and verification are implemented in CI.
- [x] Signing and verification are implemented in CI.
- [x] Provenance generation and verification are implemented in CI.
- [x] Rollback owner and procedure are recorded.
- [ ] Release/security reviewer decision is recorded.
- [ ] Owner confirms preview release publication.

## Rollback

If artifact identity, checksums, smoke tests, or release evidence are incorrect: stop publication, do not move or reuse the tag, preserve failed evidence, create the next corrected preview/patch version, and record the cause and decision.
