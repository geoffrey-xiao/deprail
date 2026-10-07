# H05-007-FU1 Release Combined Repair Evidence

Issue [#457](https://github.com/geoffrey-xiao/deprail/issues/457); [contract](../issues/H05-007-FU1-release-combined.md); release gate [#427](https://github.com/geoffrey-xiao/deprail/issues/427). Date 2026-10-07. Owner/reviewer @geoffrey-xiao. Source baseline c5544416de23fe2f908d4b51d05485f288d017c3 contains owner-merged #455/#456. Owner authorized this four-defect repair, not dispatch/publication or autonomous merge. Owner publishing/security review remains required.

## Reproduced before repair

A fresh archive of main directly running go build failed exit1: web/assets.go:11:12 pattern dist: no matching files found. Default local Node26.7.0/npm11.19.0 make generate failed exit2: expected Node22.23.3. This is local evidence, not a claim about the GitHub runner's installed Node. Inspection found independent floating-main checkouts and a tag/public release created before platform builds/assembly. Prior v0.4 Release Combined success is not v0.5 build evidence.

## Correction and failure contract

- Verification and each artifact job use existing pinned Node22.23.3/npm10.9.9 setup and locked ignore-scripts frontend build. make verify builds the shared frontend prerequisite once.
- Prepare outputs one source_sha. Every later checkout, binary Commit and annotated tag use it, not a later main or a newly resolved tag.
- Graph: prepare → verify → build-artifacts (four targets) → finalize (existing SPDX, checksums, private assembled artifact) → protected publish. Default success dependencies block publishing after preparation failure; GitHub orchestration has not been dispatched in this repair.
- Only publish has contents:write. No unused OIDC/attestation or preapproval release-asset permissions. Existing SBOM action explicitly disables automatic release upload.
- Finalize hashes and verifies the four named binaries and SPDX file; only those and SHA256SUMS enter the prepared release artifact. Publish verifies the same hashes before creating any tag.
- Protected publish creates a new immutable source-bound tag, then a draft release with the six explicit prepared assets; public transition runs only after successful creation/upload. No --clobber, forced tag mutation, automated rollback or retry. A publication failure can leave a tag/draft, requiring explicit owner recovery/new-version disposition.

## Syntax and actual build verification

- go run github.com/rhysd/actionlint/cmd/actionlint@v1.7.7 .github/workflows/release-combined.yml: exit0, no diagnostics.
- Executed actual YAML prepare shell with diagnostic-only version 0.0.0-preview.457; emitted source_sha=c5544416de23fe2f908d4b51d05485f288d017c3. This is not a selected release version.
- Fresh local clone make verify with pinned Node/npm: exit0 on canonical /private/tmp checkout, Go1.27.1, Darwin/arm64. Initial execution under macOS /tmp alias failed existing app fixture/history path checks; that failure is preserved. Canonical-path execution passed without changing product/tests or weakening assertions. [INFERENCE] Alias-sensitive fixture source paths caused the initial failures; no alias-portability fix is included here.
- Executed actual YAML Build platform artifact shell in four separate clean clones after each locked make frontend. All four compiled; go version -m recorded matching source revision and vcs.modified=false. Earlier sequential smoke produced dirty metadata from preceding untracked dist outputs; not used as clean-job evidence. Separate clean clones match the actual isolated job model.
- Actual Darwin arm64 and Darwin amd64/Rosetta binaries --version exited0 and printed v0.0.0-preview.457 / matching tag / exact c554441 source SHA. Actual host discover testdata/fixtures/mixed-repository --format json exited0 with complete npm/Python/Maven graph. Linux/Windows binaries were cross-built and inspected, not executed natively.

Diagnostic-only artifact SHA-256 (host build, not published/release-binded assets):

| Target | SHA-256 |
| --- | --- |
| linux-amd64 | a802b0296ec8efb7cb8cbb33f7ab0ec5139896cd0dd02872afe8d5cf44a9f2a6 |
| darwin-amd64 | 6850d543c7989e1113c51111d02d7f2842ddb713f7eacd6fcf99dd770cc9eabb |
| darwin-arm64 | b567ef76e85cff5d895c47d3eab74edf5b60e7560b2ee7ced8e35026ee5ede6a |
| windows-amd64 | 7520e8c2ba8a294430c92157a8cedeb09d38885ba94b657db9e4d887d4b3e7d0 |

## Actual shell publication boundaries — offline only

Extracted the actual YAML checksum/publication scripts. Used real local Git clones and bare remotes, real SHA-256 via macOS shasum-compatible adapter, and an explicit offline gh transport recorder. No real GitHub release API or repository tag push was invoked. Checksum input used a genuine existing v0.3 SPDX file solely as a fixture; it is not a matching v0.5 SBOM and not evidence that the Anchore action generated one. Live SBOM/action/protected-environment/publication execution remains unverified.

| Scenario | Observed result |
| --- | --- |
| Checksums | All four binary inputs and genuine SPDX fixture verified; five OK results. |
| Wrong source SHA | exit1, no gh call or publication. |
| Missing Windows artifact | exit1 at checksum verification; no tag, gh call or public transition. |
| Asset-create/upload transport failure | exit42, source-bound local tag retained, no public-transition call. |
| Existing tag after failed creation/upload | exit2; immutable tag not moved/reused; no public transition. |
| Main advances after prepare | Local main moved to 6fa07a517fe2a370792f8772907cfe57876b6d5d; frozen checkout/tag remained c554441. |
| Successful offline transport | exit0; draft creation included all six explicit assets and prepared SHA, followed by public-transition call; local tag resolves to prepared SHA. Not a real published release. |

Throwaway checkouts/scripts/transport adapters are removed after evidence capture. No permanent mock/wiring tests, user DB/service changes or release dispatch. Existing PR CI is separate and recorded with the review PR.

## Remaining release gates and rollback

This repairs the owner-selected authoritative workflow only. Native four-target representative workflows, browser/AT, recovery, dependency advisory/attribution/publisher provenance, signature/provenance disposition and owner release decision remain under #427. No release approval, completed release checklist or accepted preview gap is inferred. Owner reviews job permissions and source/artifact/publication transition before merge. Roll back the workflow without running it; preserve immutable tags, any draft/public releases and user data.
