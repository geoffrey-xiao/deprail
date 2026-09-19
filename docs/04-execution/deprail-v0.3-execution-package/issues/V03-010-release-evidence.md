# V03-010 Complete Cross-Platform Evidence and Release Gate

## Planning metadata

- Type: release
- Area: docs
- Priority: P1
- Risk: R2
- Epic: EPIC-004
- Dependencies: V03-006 through V03-009

## Goal
Prove v0.3 plan accuracy, read-only safety, schema compatibility, and cross-platform semantic equivalence.

## Acceptance criteria

- [ ] Linux, macOS, and Windows plan smoke passes.
- [ ] npm, Python, Java, and mixed fixtures have linked plan artifacts.
- [ ] Repository trees and lockfiles remain unchanged.
- [ ] Schema, golden, hostile-input, stale-input, and mutation tests pass.
- [ ] Release identity, checksums, SBOM, signing, provenance, and remaining risks are recorded.
- [ ] Owner and security/architecture reviewer record the v0.3 decision.
- [ ] v0.4 handoff contract is linked.

## Exclusions
No stable release approval without complete evidence and explicit owner decision.
