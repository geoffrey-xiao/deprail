# H05-008-FU1: Remediate x/text GO-2026-5970 and Verify SDK Applicability

GitHub issue: [#460](https://github.com/geoffrey-xiao/deprail/issues/460).

Observed technical evidence: [before/after, SDK applicability and actual consumer smoke](../tracking/H05-008-FU1-XTEXT-EVIDENCE.md). Owner acceptance boxes below remain unchecked.

## Planning metadata

- Parent: H05-008 / #427; private-store dependency evidence #421; delivery epic #418.
- Type: bug; area: foundation; priority: P0; risk: R3 (dependency/security).
- Target version/milestone: v0.5.0; Sprint 4.
- Owner/reviewer: @geoffrey-xiao; independent review optional under ADR-0005.
- Dependencies: owner-merged #459, synchronized main `96f1078`; no active review PRs at kickoff.
- Owner scope request: “所以要怎么做 你能帮我做吗” after the advisory investigation. This authorizes a technical repair and review PR, not residual-risk acceptance, autonomous merge, release dispatch, tagging or publication.
- Blocked reason: none for this bounded repair; integrated release acceptance remains #427.

## Definition of Ready

- [x] Value: remove the selected module version affected by the malformed-UTF-8 normalization iterator vulnerability and establish the exact SDK/binary applicability rather than treating a dependency warning as a proven exploit.
- [x] Scope, exclusions, compatibility/failure behavior and rollback are explicit below.
- [x] Reconciled with product verified dependency changes, architecture secure Go build chain, roadmap v0.5 local history/console, release plan §§14–15, SEC-13 and H05-008 PG-01; no new capability or displaced work.
- [x] Required observable smoke, existing verification, owner/reviewer, version and evidence requirements are explicit.
- [x] Public Go advisory index lists four x/text advisories with maximum fixed floor `0.39.0`; candidate requires Go `1.25.0`, compatible with pinned `1.27.1`, and retains the existing BSD-3-Clause license class. No new package scripts or automatic tools installed.

## Goal and scope

Upgrade only the selected `golang.org/x/text` module from `v0.14.0` to the minimum fixed `v0.39.0` in `go.mod`/`go.sum`. Record checksum verification, license/source/fix review, module/import graph and actual binary identity. Preserve other dependency pins and Go `1.27.1` unless concrete current-SDK evidence proves that toolchain remediation is necessary; any such expansion requires a revised plan before editing.

Distinguish canonical module normalization from Go SDK-vendored normalization. Go upstream's [member assessment](https://github.com/golang/go/issues/80142#issuecomment-5121323230) says the triggering iterator path is unused by the standard library. Verify actual SDK callers and binary iterator entry points. Existing `Form` methods/private retained norm symbols alone neither prove the iterator reachable nor prove a patched vendor copy.

## Out of scope

No vulnerability suppression, generic risk acceptance, Go version-line change, unrelated upgrades, API/CLI/schema/storage/permissions changes, migrations, frontend changes, scanner installation, telemetry/retries, license-packaging redesign, signing/provenance, full platform/browser/AT/recovery acceptance, dispatch/tag/publish/merge. Existing release gaps stay open.

## Inputs, outputs and failure behavior

Inputs: authoritative GO-2026-5970/CVE-2026-56852 record, upstream fix/test, current selected graph and SDK, and unchanged fixed project/history fixtures. Outputs: fixed selected module plus reproducible before/after malformed-input and actual binary/consumer evidence. A timed-out old iterator reproducer is not a clean result; compile first and separately bound execution. Fixed iteration must finish with exact preserved invalid bytes, including byte/string initialization. Dependency/checksum/build/regression failures block the candidate and do not alter user data or authorize publication. Existing CLI errors, completeness, exit codes and local-console authentication remain unchanged.

## Required verification and acceptance

- [ ] Bounded upstream input `f3 cc 80` reproduces failure with the old module; the fixed module completes NFC/NFKC byte/string iteration with exact expected bytes and valid/empty controls.
- [ ] Selected module is `v0.39.0`, checksum verification succeeds, other production pins remain unchanged, and reviewed Go minimum/license/source/known-advisory evidence is recorded.
- [ ] Current SDK vendor/import/callers and actual CLI symbols establish whether the triggering iterator entry points are reachable; no claim that the older SDK vendor source is patched.
- [ ] Existing `make verify` and `make test-integration` pass on the pinned local toolchain; specific schema/history consumer tests and actual mixed-repository JSON discovery/packaged local-console smoke pass without changing contracts.
- [ ] Exact PR-head three-OS CI passes; local-only or browser/OS limitations remain explicit rather than being generalized into release approval.
- [ ] Original PG-01/store evidence records the technical correction and review link additively; #427 and other release/owner-acceptance gaps remain open.

## Owner review and evidence

Owner reviews selected dependency compatibility/security, SDK applicability reasoning, exact binary/source/checksums, local/CI evidence and remaining release gaps. Preserve the previous investigation as history and link this issue from the original evidence. Attach actual commands, exits, before/after result, source/advisory citations, binary identity, module hashes, platform/browser identity and reviewed CI link to the dedicated remediation evidence record and PR. No permanent test that merely pins a dependency/version or copies a third-party implementation; throwaway vulnerability smoke plus existing consumer regression checks are appropriate.

## Final acceptance

- [ ] Owner inspected every criterion and separately recorded dependency/security review.
- [ ] Required local and exact-head CI evidence accepted; remaining risk recorded.
- [ ] Owner-reviewed merge linked before Project Done/closure. No release go/no-go inferred.

## Rollback

Revert only dependency pins/evidence with an explicit security-regression disposition. Preserve DB/WAL/SHM, artifact and source bytes; no migration/reset/overwrite. Do not publish the vulnerable baseline as an approved repaired candidate.

## Authoritative references

- [GO-2026-5970](https://pkg.go.dev/vuln/GO-2026-5970), [machine record](https://vuln.go.dev/ID/GO-2026-5970.json), [current module advisory index](https://vuln.go.dev/index/modules.json).
- [Upstream fix and malformed-input regression](https://github.com/golang/text/commit/5ae8e578e495731553eddba11b2d0e86c91a00ce).
- [Pinned candidate Go minimum](https://github.com/golang/text/blob/v0.39.0/go.mod), [BSD license](https://github.com/golang/text/blob/v0.39.0/LICENSE).
