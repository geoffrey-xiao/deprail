# H05-008-FU1 x/text Remediation Evidence

Issue [#460](https://github.com/geoffrey-xiao/deprail/issues/460); contract [H05-008-FU1](../issues/H05-008-FU1-xtext-remediation.md). Date 2026-10-07. Base: owner-merged #459 / main `96f1078`. This is dependency-repair evidence, not integrated #427 acceptance or release authorization. Owner dependency/security review remains separate.

Review [PR #461](https://github.com/geoffrey-xiao/deprail/pull/461), Refs #460, implementation `7b1a8a5`. After creation, live gh verification confirmed area:foundation/risk:R3/priority:P0/type:bug on issue and PR, milestone v0.5.0, DepRail Project Sprint4/Review and owner reviewer @geoffrey-xiao. Track final-head existing three-OS CI in the PR checks/completion record rather than substituting prior-head checks; no owner approval/merge/Project Done/closure asserted.

## Decision and authoritative sources

Choose minimum fixed `golang.org/x/text v0.39.0`; retain Go `1.27.1` and all other production module pins. [GO-2026-5970 machine record](https://vuln.go.dev/ID/GO-2026-5970.json), alias CVE-2026-56852, describes infinite iteration on invalid UTF-8 in `norm.Iter`; fixed floor is 0.39.0. The current [Go module advisory index](https://vuln.go.dev/index/modules.json) returned four x/text records (GO-2020-0015, GO-2021-0113, GO-2022-1059, GO-2026-5970), with respective floors 0.3.3/0.3.7/0.3.8/0.39.0. This is dated authoritative-index review, not a scanner callgraph or universal zero-vulnerability guarantee.

[Fix 5ae8e578](https://github.com/golang/text/commit/5ae8e578e495731553eddba11b2d0e86c91a00ce) gives invalid-rune properties size 1 and includes input `f3 cc 80` with unchanged expected bytes. [Candidate go.mod](https://github.com/golang/text/blob/v0.39.0/go.mod) requires Go1.25.0, satisfied by pinned Go1.27.1. [Candidate license](https://github.com/golang/text/blob/v0.39.0/LICENSE) retains the existing BSD-3-Clause terms/Go Authors copyright; binary redistribution still requires the copyright/license/disclaimer, and this repair does not claim final-asset notice/publisher-provenance approval. Sources were fetched through Firecrawl CLI, not copied into release artifacts.

## Real before/after iterator smoke

A throwaway Go executable directly used the real selected module's `norm.Iter`, byte and string initialization, NFC and NFKC. It concatenated returned segments and asserted exact output; a segment-count guard `segments > len(input)+1` prevented endless outer iteration. Each executable was compiled before execution and also bounded externally to two seconds, so build/download time was not counted as a vulnerability timeout.

| Observed command/scenario | Result |
| --- | --- |
| `go build -trimpath -o /tmp/deprail-xtext-before /tmp/deprail-xtext-remediation-smoke.go`; `go version -m` | Exit0; Go1.27.1, x/text v0.14.0, old checksum `h1:ScX5w1eTa3QqT8oi6+ziP7dTV1S2+ALU0bI+0zXKWiQ=`. |
| Old executable, first invalid/NFC/bytes case | Printed input `f3cc80`; exited1 with `iterator did not advance` at the protective segment-count bound. **Not an external timeout**: the guard exposed the repeated nonprogress safely. The initial orchestration assertion incorrectly expected timeout; the actual exit1/nonprogress result is retained, not replaced by a passing rerun. |
| `go mod download -json golang.org/x/text@v0.39.0`; `go mod tidy`; `go mod verify` | Exit0; `all modules verified`; origin `https://go.googlesource.com/text`, tag v0.39.0, commit `b326f3d3c814ab79b3c516f4ac03c2314d8df65f`. SumDB `sum.golang.org`, proxy `https://proxy.golang.org,direct`. |
| Recompile `/tmp/deprail-xtext-after` from the same harness and run | Exit0, no timeout/stderr, all **12 exact cases** passed: invalid NFC/NFKC byte/string outputs `f3cc80`; valid decomposed `65cc81` outputs composed `c3a9`; empty inputs output empty bytes. |

Fixed checksums: module `h1:UbZz4pLOvn600D6Oh6GGEI6VAmndrEBLv8/6BEXzyus=`, go.mod `h1:3UwRclnC2g0TU9x8PZiyfOajCd1zaUNHF9cvqcQZ+ZM=`. Checksum verification establishes consistency, not publisher identity. No third-party-version/wiring test or copied upstream implementation is added permanently.

`go list -m all` compared with the isolated baseline `-modfile` graph changed exactly one selected version: removed x/text v0.14.0, added v0.39.0. All other selected module versions are unchanged.

## Distinct standard-library vendor assessment

The original [investigation](H05-002-STORE-EVIDENCE.md#additive-xtext-advisory-applicability-evidence--2026-10-07) remains historical: Go1.27.1 vendors x/text v0.37.0; its old invalid-size implementation and broad norm symbols were observed. Upgrading go.mod does not replace it. The official machine record lists **25** symbols, correcting the earlier prose's count22; broad `Form` method overlap does not establish the iterator trigger.

New authoritative evidence: Go member nicholashusin [explained on 2026-07-29](https://github.com/golang/go/issues/80142#issuecomment-5121323230) that build/tests still pass after deleting the vulnerable symbols and dependent code, so this triggering path is not reachable from the main Go repository despite vendoring and does not need backporting. This is the maintainer's applicability assessment, not an adopted blanket advisory suppression.

Local SDK source checks found idna calls only `norm.NFC.String`, `IsNormalString`, `QuickSpan`, and `Bytes`; no SDK caller of `norm.Iter`, `norm.NewReader` or `norm.NewWriter`. Iterator-specific `.nextMain` assignment occurs inside vendored iter.go's Init/InitString/Seek, not at idna callers. The actual dependency paths remain `schemas/history-v1 → jsonschema/v6 → x/text/language` for the selected module and `internal/transport/localhttp → net/http → vendor/x/net/idna → vendor/x/text/unicode/norm` for SDK normalization. Neither a module-only clean result nor private function-pointer retention is a whole-binary security clearance. No govulncheck result is claimed; the tool remains unavailable and was not installed.

Owner must review this precise source/binary-bounded applicability assessment with the repaired candidate. Future imports of module norm/SDK consumers or another toolchain require reassessment; no generic suppression or permanent exemption is introduced. Other PG-02–PG-12 release decisions remain pending.

## Existing verification and actual consumer smoke

Platform Darwin23.6.0 arm64; Go1.27.1, pinned Node22.23.3/npm10.9.9. `env PATH="/tmp/node-v22.23.3-darwin-arm64/bin:$PATH" make verify` and `make test-integration` each exited0. Generation, vet, full unit/contract/integration, frontend typecheck/production build and Go builds passed; explicit `go test ./schemas/history-v1 ./internal/store/history ./internal/transport/localhttp ./cmd/deprail` exited0. First unquoted PATH invocation exited127 before make because a host PATH entry contains spaces; quoting corrected the invocation, not repository code. Existing assertions/goldens/fixtures are unchanged.

`go build -trimpath -o /tmp/deprail-xtext-fixed ./cmd/deprail` exited0. SHA-256 **`bc1dfc803a70a31b87e3ff0a657bb725cffbae35925978adee5a8d28c50219ef`**; `go version -m` records Go1.27.1, x/text v0.39.0 and source `96f10788823e6baf426fce5efa3ffbb82b66aa5e`, `vcs.modified=true` (repair/evidence uncommitted at build). Actual `--version` exited0, development/tag/commit unknown. This is an identified repair-analysis binary, not a final immutable release candidate.

`go tool nm` still has **63 vendor norm symbols**, including Form.Bytes/String/IsNormalString, compInfo, doNormComposed and nextComposed, but no canonical-module norm symbols and no `(*Iter).Init/InitString/Next/Seek` entry points. `CGO_ENABLED=0 GOOS=<os> GOARCH=<arch> go list -deps ./cmd/deprail` exited0 for linux/amd64, darwin/amd64, darwin/arm64, windows/amd64: all retain vendor norm/idna and exclude canonical norm. Graph/symbol absence alone is not proof against compiler inlining; interpreted with exact SDK callers and the upstream member assessment, this supports the **specific unused iterator-trigger assessment**, not a claim that SDK source was patched or the whole binary has no vulnerabilities. These four graph checks are not four native executions.

Actual `discover testdata/fixtures/mixed-repository --format json` exited0 with three complete npm/Python/Maven workspaces, overall complete and no diagnostics. No vulnerability-scan success is inferred.

Actual packaged `web --open` used a fresh private canonical user profile and real embedded frontend/HTTP/SQLite. Only the OS launcher was replaced by a compiled private-file capture helper; native macOS opener behavior was not verified. The listener announced `http://127.0.0.1:56795/console/` without a token. Actual managed headless Chromium150.0.0.0, 1280×800, rendered **No saved scans**, refreshed real `GET /api/v1/scans` twice with200 and navigated to **Local history, read only**. Screenshot/observation were captured; browser runtime errors were empty. Bootstrap fragment was cleared; local/session storage counts0, cookie/referrer empty. SIGINT shutdown exited0 with only normal interactive guidance. The automation's initial network-URL metadata retained a bootstrap fragment; the isolated listener was stopped immediately after smoke, revoking that session, and no credential-bearing URL is copied into durable/public evidence. Diagnostic-tool behavior was reported.

History DB SHA-256 before/after browser read/shutdown: `2867ab7220e1f46a62f58d32cfc4c677385069af99287e9f5d56c9113b09a55b` (unchanged). The four representative source files' before/after hashes matched:

| File | SHA-256 |
| --- | --- |
| frontend/package-lock.json | ae2e8c395115a817fe85afd2bca20166f3ddfd4635fdf1ceceb5bb50a9a81f9c |
| frontend/package.json | f62a63c50dc467066a1cfeb27752670b5fd09a22365aac1d9cff6a549986659e |
| services/worker/pom.xml | 9bb7d52f45a95173cf83027059f3471e03c8befff56d9d3080713afea92cf8ad |
| services/api/requirements.txt | 1d277ef3981a3e49b02912a0f03fe1ab563539d7e4e1b5c1e6404a57b19d883f |

### Unchanged missing-scanner capture observation

For offline smoke PATH deliberately excluded OSV-Scanner; no scanner installed and no findings fabricated. Real `scan testdata/fixtures/mixed-repository --save-history --format json` exited3 with failed report/three SCANNER_NOT_FOUND errors and HISTORY_WRITE_FAILED. The first profile used macOS `/tmp` symlink spelling; using its canonical `/private/tmp` spelling still reported the save failure. An isolated old-v0.14.0 binary built with `go build -mod=mod -modfile=/tmp/deprail-xtext-baseline.mod -trimpath ... ./cmd/deprail`, without changing repaired production pins, produced the **same exit/status/save diagnostic** on the canonical profile. Therefore this observed failure is not introduced by the module upgrade; its full underlying cause is not established here. The console correctly showed the actual empty store. No successful capture or failed-scan history workflow is claimed; retain this observation under #427 PG-05 for the original integrated failure-history gate, rather than expanding this dependency fix or suppressing the diagnostic.

Limits: local macOS/managed Chromium repair smoke is not full native/browser/AT, representative real-scanner release workflow, restore rehearsal, final SBOM/license/provenance or release publication proof. Owner acceptance and final integrated #427 gates remain open.
