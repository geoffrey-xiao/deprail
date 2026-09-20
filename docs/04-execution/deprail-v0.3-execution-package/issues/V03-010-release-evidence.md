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
- [ ] The real mixed-repository scan-to-plan path is proven with the reviewed `deprail` binary; raw scan findings retain ecosystem PURLs and workspace identity.
- [ ] Repository trees and lockfiles remain unchanged.
- [ ] Schema, golden, hostile-input, stale-input, and mutation tests pass.
- [ ] Manual procedure and expected results are maintained in [`../MANUAL-TEST-GUIDE.md`](../MANUAL-TEST-GUIDE.md).
- [ ] Release identity, checksums, SBOM, signing, provenance, and remaining risks are recorded.
- [ ] Owner and security/architecture reviewer record the v0.3 decision.
- [ ] v0.4 handoff contract is linked.

## Required real-test evidence

The release gate MUST include a run of the actual built binary, not only package tests:

```bash
go build -o deprail ./cmd/deprail
./deprail doctor --format json > local_test/doctor-manual.json
./deprail scan testdata/fixtures/mixed-repository --format json --output local_test/mixed-repo-manual.json
FINDING_KEY="$(jq -r '.findings[0].TargetID' local_test/mixed-repo-manual.json)"
./deprail fix plan --report local_test/mixed-repo-manual.json --finding "$FINDING_KEY" --format json > local_test/fix-plan-fixed.json
```

Expected results:

- `doctor` exits `0`.
- The mixed scan exits `0`, reports `complete`, and contains findings.
- The selected finding contains an ecosystem PURL and workspace identity.
- `fix plan` exits `0` and emits a schema-valid JSON plan.
- The plan identifies the expected manifest and lockfile, candidate state, risks, verification requirements, and provenance.
- The repository tree and lockfiles are unchanged.

Observed local evidence for the current reviewed implementation:

- Scan: `complete`, 14 findings, no scan errors.
- First finding: `GHSA-29mw-wpgm-hmr9`, PURL `pkg:npm/lodash@4.17.20`, workspace `frontend`.
- Plan: exit `0`, stable plan ID generated, affected files `frontend/package.json` and `frontend/package-lock.json`.
- Automated checks: `go test ./...`, `go vet ./...`, and `git diff --check` passed locally; Linux, macOS, and Windows CI passed on PR #229.

Evidence artifacts are generated under ignored `local_test/`; release records MUST link the reviewed CI run and preserve the exact command output and platform/binary identity.
The current PRs are sequenced: PR #227 supplies scan `--output` parsing, and PR #229 supplies scan-to-plan PURL/workspace enrichment. The release-gate workflow must be rerun after both are present on the candidate release baseline; running PR #229 alone from `main` is expected to fail at scan argument parsing.

## Exclusions
No stable release approval without complete evidence and explicit owner decision.
## GitHub tracking

- Issue: [#215](https://github.com/geoffrey-xiao/deprail/issues/215)
- Parent epic: [#205](https://github.com/geoffrey-xiao/deprail/issues/205)
