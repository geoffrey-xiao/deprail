# DepRail v0.2 Release Evidence Record

**Status:** Evidence collection pending; this record is not release approval.
**Release mode:** Preview / RC / Stable (select one before completing)
**Target version:** `0.2.0`
**Release owner:** `[not recorded]`
**Reviewer:** `[not recorded]`
**Decision:** `[not recorded]`

Missing or inconsistent evidence blocks the applicable release mode. Do not replace missing evidence with an unchecked claim.

## Source and identity

| Field | Evidence |
| --- | --- |
| Source commit | `[not recorded]` |
| Release tag | `[not recorded]` |
| DepRail version | `[not recorded]` |
| CLI identity output | Link to `deprail doctor --format json` output |
| Build command and Go version | `[not recorded]` |
| Release workflow run | Link to GitHub Actions run |

## Artifact inventory and checksums

Every artifact requires a recorded SHA-256 checksum and a smoke result from the same reviewed source commit.

| Artifact | Platform | Size | SHA-256 | Smoke result | Evidence |
| --- | --- | ---: | --- | --- | --- |
| `deprail-linux-amd64` | Linux amd64 | `[not recorded]` | `[not recorded]` | `[not recorded]` | `[not recorded]` |
| `deprail-darwin-amd64` | macOS amd64 | `[not recorded]` | `[not recorded]` | `[not recorded]` | `[not recorded]` |
| `deprail-darwin-arm64` | macOS arm64 | `[not recorded]` | `[not recorded]` | `[not recorded]` | `[not recorded]` |
| `deprail-windows-amd64.exe` | Windows amd64 | `[not recorded]` | `[not recorded]` | `[not recorded]` | `[not recorded]` |

Checksum manifest: `[not recorded]`

## Platform and CLI smoke evidence

For every supported artifact, attach command output and the CI or manual-run link:

- `deprail doctor --format json`
- discovery against `testdata/fixtures/mixed-repository`
- scan behavior when OSV-Scanner is unavailable, including exit code `3`
- reported version, tag, and commit identity

| Platform | CI or run link | Doctor | Discover | Scan failure behavior | Identity match |
| --- | --- | --- | --- | --- | --- |
| Linux | `[not recorded]` | `[not recorded]` | `[not recorded]` | `[not recorded]` | `[not recorded]` |
| macOS | `[not recorded]` | `[not recorded]` | `[not recorded]` | `[not recorded]` | `[not recorded]` |
| Windows | `[not recorded]` | `[not recorded]` | `[not recorded]` | `[not recorded]` | `[not recorded]` |

## Scanner and representative repositories

| Evidence | Result | Link or output |
| --- | --- | --- |
| OSV-Scanner version | `[not recorded]` | `[not recorded]` |
| JavaScript repository scan | `[not recorded]` | `[not recorded]` |
| Python repository scan | `[not recorded]` | `[not recorded]` |
| Java repository scan | `[not recorded]` | `[not recorded]` |
| Mixed repository scan | `[not recorded]` | `[not recorded]` |

Every release-gating scan must complete successfully. Partial or failed scans block the relevant release mode.

## Supply-chain controls

| Control | Status | Evidence or gap rationale |
| --- | --- | --- |
| SBOM | Not implemented unless linked evidence is supplied | `[not recorded]` |
| Signing | Not implemented unless linked evidence is supplied | `[not recorded]` |
| Provenance | Not implemented unless linked evidence is supplied | `[not recorded]` |

A preview may carry explicitly approved gaps. A stable release must satisfy the release policy or record an owner-approved rejection/exception.

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
