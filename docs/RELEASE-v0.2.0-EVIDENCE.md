# DepRail v0.2.0 Stable Release Evidence

**Status:** Stable release preparation; final approval pending.
**Release mode:** Stable
**Target version:** `v0.2.0`
**Release owner:** `@geoffrey-xiao`
**Reviewer:** `@geoffreyxiaoai`
**Rollback owner:** `@geoffrey-xiao`
**Decision:** `[GO / NO-GO — owner to confirm]`
**Decision date:** `[owner to confirm]`

This record applies only to the immutable `v0.2.0` tag. Preview evidence remains in the version-specific preview records.

## Release baseline

| Field | Evidence |
| --- | --- |
| Latest approved preview | `v0.2.0-preview.2` |
| Reviewed source commit | `[record after stable tag]` |
| Stable release tag | `v0.2.0` |
| Manifest value | `0.2.0` |
| Release workflow run | `[record after stable publication]` |
| Release owner smoke report | Preview 2 macOS and Windows smoke reported passed by owner; stable artifact output pending |

## Required stable controls

| Control | Status | Evidence or disposition |
| --- | --- | --- |
| SBOM | `[required or approved exception]` | Stable policy requires linked evidence or explicit owner-approved exception |
| Signing | `[required or approved exception]` | Stable policy requires linked evidence or explicit owner-approved exception |
| Provenance | `[required or approved exception]` | Stable policy requires linked evidence or explicit owner-approved exception |

## Artifact inventory and checksums

| Artifact | Platform | Size | SHA-256 | Smoke result | Evidence |
| --- | --- | ---: | --- | --- | --- |
| `deprail-linux-amd64` | Linux amd64 | `[pending]` | `[pending]` | `[pending]` | `[workflow URL]` |
| `deprail-darwin-amd64` | macOS amd64 | `[pending]` | `[pending]` | `[pending]` | `[owner output or URL]` |
| `deprail-darwin-arm64` | macOS arm64 | `[pending]` | `[pending]` | `[pending]` | `[owner output or URL]` |
| `deprail-windows-amd64.exe` | Windows amd64 | `[pending]` | `[pending]` | `[pending]` | `[owner output or URL]` |
| `SHA256SUMS` | All artifacts | — | `[pending]` | `[pending]` | `[release URL]` |

## Stable acceptance checklist

- [ ] Stable source commit and tag match the manifest and CLI identity.
- [ ] `make verify` passes from a clean checkout.
- [ ] All artifacts have checksums.
- [ ] Linux, macOS, and Windows artifact smoke output is attached.
- [ ] Representative JavaScript, Python, Java, and mixed scans are attached.
- [ ] SBOM status is implemented or explicitly approved as an exception.
- [ ] Signing status is implemented or explicitly approved as an exception.
- [ ] Provenance status is implemented or explicitly approved as an exception.
- [ ] Rollback owner and procedure are recorded.
- [ ] Release/security reviewer decision is recorded.
- [ ] Owner gives explicit stable GO decision.

## Rollback

If artifact identity, checksums, smoke tests, or release evidence are incorrect: stop publication, do not move or reuse the tag, preserve failed evidence, create the next corrected preview/patch version, and record the cause and decision.
