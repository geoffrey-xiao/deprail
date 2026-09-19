# EPIC-003 Read-Only Planning Adapters

## Outcome
Supported v0.2 ecosystems produce normalized remediation candidates without executing package-manager operations.

## Scope

- Shared planning adapter port.
- npm/pnpm/Yarn ownership and constraint analysis.
- Python requirements/uv/Poetry analysis.
- Maven/Gradle analysis.
- Direct/transitive ownership, candidate compatibility, lockfile, peer/runtime/engine, and major-version risks.

## Issues

- V03-005: Define read-only package-manager planning adapter contract.
- V03-006: Implement JavaScript planning adapter.
- V03-007: Implement Python and Java planning adapters.

## Acceptance

Adapters produce the same domain plan contract, return explicit unknown/unavailable/rejected states, preserve provenance, use structured future commands with working directories, and never execute package managers or scripts.
## GitHub tracking

- Parent issue: [#204](https://github.com/geoffrey-xiao/deprail/issues/204)
- Child issues: [#210](https://github.com/geoffrey-xiao/deprail/issues/210), [#211](https://github.com/geoffrey-xiao/deprail/issues/211), [#212](https://github.com/geoffrey-xiao/deprail/issues/212)
