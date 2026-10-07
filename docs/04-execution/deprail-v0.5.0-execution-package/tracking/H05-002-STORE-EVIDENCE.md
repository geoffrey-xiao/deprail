# H05-002 Private History Store — Implementation Evidence

**GitHub issue:** [#421](https://github.com/geoffrey-xiao/deprail/issues/421)

**Current status — 2026-10-07:** [PR #436](https://github.com/geoffrey-xiao/deprail/pull/436) was owner-merged 2026-10-01; original issue #421 remains open in Review with criteria/security dispositions pending. Exact-head CI [36806173168](https://github.com/geoffrey-xiao/deprail/actions/runs/36806173168) passed Ubuntu/macOS/Windows; initial Windows failures remain recorded below. [Individual STORE-01–07 worksheet](H05-008-READINESS.md#h05-002--421--private-store-criteria) maps evidence and missing candidate/platform/recovery decisions. Dated advisory findings below distinguish selected-module and Go-SDK vendor copies; no security acceptance or release readiness follows from merge.
**Branch:** `feat/h05-002-private-history-store` from reviewed main `da59b10`.

## Delivered behavior

- `internal/store/history` provides per-user `OpenReadOnly`, explicit `OpenWriter`, idempotent/conflict-safe `Append`, validated `Get`, and bounded keyset `List`. The DB path is `os.UserConfigDir()/.deprail/history.sqlite3`; reads of an absent store create no paths, and writers require a repository boundary.
- SQLite `user_version=1` uses the exact H05-002 row/index shape, WAL, per-connection foreign keys, writer `synchronous=FULL`, `_txlock=immediate`, a 2,000 ms busy timeout, serialized local writers, atomic entry/reference inserts and in-transaction entry/byte admission. No retries, eviction, artifact deletion, automatic repair, restore or downgrade.
- Inputs and persisted rows validate against the embedded `history-v1` JSON Schema and `normalize.ValidateProjection`. Reads check canonical JSON, UUIDv4, projection/row summary and source-report metadata, sorted artifact-reference equality and supported schema version before returning data. Errors expose stable codes without paths or SQL details.
- POSIX paths enforce current-user ownership and directory `0700` / DB and backup `0600`; Windows paths use a protected DACL containing only the current user. Path components reject symlinks/reparse points, and repository containment fails closed. Backups use SQLite's online backup API, unique sibling names, restrictive permissions and schema/row/SQLite integrity validation before returning a completed path.
- The history schema package now offers a production JSON Schema validator so storage checks the committed schema on both write and read; the existing examples and schema tests use the same path-pattern matcher.

## Dependency, license and vulnerability review

Pinned driver: [`modernc.org/sqlite v1.59.0`](https://pkg.go.dev/modernc.org/sqlite@v1.59.0), Go module checksum `h1:X1es1GpqBlS/5T+vbM4HLUdaa8OtQx468DF2vrx+38A=`, `go.mod` checksum `h1:+paeT2A3iPRHkQDwG7oA6Tk0zQd5woMEI8q7orfry8k=`. The tag resolves to upstream GitLab commit `c96a4e6cb22254bf70026502a781a54a053c2cf0`; its declared Go baseline is 1.25.0, below this repository's pinned Go 1.27.1. The driver is CGo-free and declares BSD-3-Clause; the embedded SQLite engine is public domain. `go.mod`/`go.sum` pin the complete module graph.

The archived [CycloneDX 1.6 SBOM](H05-002-dependencies.cdx.json) lists the reviewed Go dependency components and license identifiers: BSD-3-Clause, MIT and Apache-2.0. The isolated driver closure scan with OSV-Scanner 2.6.0 reported no issues. The final lockfile-scope command `osv-scanner scan source --lockfile go.mod` exited 0 and reported 0 affected packages. The SBOM also records `GO-2026-5970` / `CVE-2026-56852` for the pre-existing `golang.org/x/text v0.14.0`, reachable through the pre-existing `jsonschema/v6` dependency (`go mod why -m golang.org/x/text`). The scanner's affected-package count is zero; the advisory remains disclosed rather than being treated as a safe-version guarantee. The selected SQLite driver closure itself is not affected. This issue does not upgrade the unrelated `x/text` dependency; owner security review must disposition the recorded residual risk.

Dependency review proves candidate provenance and a clean SQLite closure; it does not replace owner acceptance of this R3 storage/dependency change.

## Behavioral verification

| Command/scenario | Observed result |
| --- | --- |
| `go test ./internal/store/history ./schemas/history-v1` | Passed. Covers empty read without creation; current-v1 initialization/reopen; row/projection and future-version refusal; deterministic keyset order; identical/conflicting retries; repeated source scan IDs; atomic row+ref rollback; child-process exit before commit; cancellation; external SQLite lock bounded near 2 seconds; SQLite `max_page_count` failure; exact/over byte and entry quotas; POSIX modes and refusal of broad read-only DB permissions; symlink/repository-boundary rejection; WAL/foreign-key/FULL settings; unique online backups and validation. |
| `go run ./internal/store/history/smoke` (throwaway harness, then removed) | Printed `{"missing_read_only_created_no_paths":true,"entry_reopened":true,"list_returned_one":true,"artifact_reference_round_trip":true}` on the local Darwin/arm64 host. It called the public store API, opened a missing store read-only, explicitly wrote one validated occurrence/ref, closed/reopened, then fetched and listed it. |
| `GOOS=windows GOARCH=amd64 go test -c -o /tmp/deprail-history-windows.test.exe ./internal/store/history`; exact-head PR #436 CI [36806173168](https://github.com/geoffrey-xiao/deprail/actions/runs/36806173168) | Cross-compilation passed. The three OS jobs passed on the latest PR head. Earlier Windows failures exposed newly created ownership and the two same-user ACEs (current-object full access and inherit-only child access); the final run passed the strict protected-DACL checks. |
| `go mod tidy` | Passed; SQLite is direct and existing `golang.org/x/sys` is pinned at the compatible `v0.47.0`. |
| SBOM/checksum/source-diff validation | `jq` validated CycloneDX 1.6 structure, 12 dependency components and allowed license IDs; `go mod verify` reported all modules verified; `git diff --check` passed. |
| `make verify` | Passed: `go generate ./...`, `go vet ./...`, `go test ./...`, and `go build ./...` all completed successfully on the local Darwin/arm64 host. |

Local `make verify` passed, and latest exact-head CI run [36806173168](https://github.com/geoffrey-xiao/deprail/actions/runs/36806173168) passed Ubuntu/macOS/Windows. Owner technical/security review and acceptance remain required. POSIX mode tests ran locally on Darwin; Windows ACL runtime tests passed in CI.

## Limits and unresolved acceptance evidence

- SQLite `max_page_count` is the deterministic full-storage failure injection; it is not an actual filesystem `ENOSPC` test. The child-process test exits after inserting row and refs but before commit, then reopens and verifies neither uncommitted row nor ref survived.
- No older DepRail schema exists to migrate. Current v1 opens without migration; this branch tests creation, unique online snapshots, snapshot integrity, no-overwrite naming and cancellation preservation. It does not invent a fake migration or claim migration-failure/restore rehearsal evidence. Manual restore remains a future release gate; no restore action is implemented.
- The store commits digest references only. Missing/digest-mismatch artifact integrity states and artifact retrieval remain the dependent H05-003 responsibility; this package never reads, deletes or garbage-collects raw artifacts.
- The 256 MiB bound covers summed projection JSON only. SQLite pages/indexes, WAL/SHM, backups and external artifacts can use more disk; no hard filesystem quota is claimed.
- Owner review must separately accept schema/permission/dependency impact, the disclosed `x/text` advisory and logical-versus-physical growth risk. Issue #421 remains open at Project Status `Review` until evidence, owner review and CI are recorded.

## Additive x/text advisory applicability evidence — 2026-10-07

The warning is confirmed as a real, current Go vulnerability record, not merely an unverified SBOM label. The Go vulnerability database record [GO-2026-5970](https://pkg.go.dev/vuln/GO-2026-5970) (alias [CVE-2026-56852](https://www.cve.org/CVERecord?id=CVE-2026-56852)) was published 2026-07-14 and modified 2026-08-10. It describes an infinite loop on invalid UTF-8 in `golang.org/x/text/unicode/norm`; affected package version range is before v0.39.0, and it lists 22 affected `norm` symbols including `Iter.Init`, `Iter.InitString`, `Iter.Next`, `Iter.Seek` and multiple `Form` methods. The record links [Go issue 80142](https://go.dev/issue/80142), [fix CL 794100](https://go.dev/cl/794100), and the machine-readable [vuln database record](https://vuln.go.dev/ID/GO-2026-5970.json). Thus pinned v0.14.0 falls within the affected range; a corrective dependency recommendation is `golang.org/x/text` v0.39.0 or later. No dependency change or remediation has been made here.

Local evidence was collected on 2026-10-07 at source checkout `a01df71368708ce8081c6051116ad436695c5c45`:

| Command | Observed result |
| --- | --- |
| `go list -m -json golang.org/x/text` | Selected module is `golang.org/x/text v0.14.0`, indirect, checksum `h1:ScX5w1eTa3QqT8oi6+ziP7dTV1S2+ALU0bI+0zXKWiQ=`, module Go version 1.18. |
| `go mod why -m golang.org/x/text` | Dependency path is `github.com/geoffrey-xiao/deprail/schemas/history-v1` → `github.com/santhosh-tekuri/jsonschema/v6` → `golang.org/x/text/language`. |
| `go list -deps -f '{{if eq .ImportPath "golang.org/x/text/unicode/norm"}}{{.ImportPath}}{{end}}' ./...` | Printed no exact unprefixed package path. This tests the selected module copy, not `vendor/golang.org/x/text/unicode/norm` in the Go standard library. |
| `go list -deps -f '{{if eq .ImportPath "golang.org/x/text/language"}}{{.ImportPath}} {{join .Imports " "}}{{end}}' ./...` | `golang.org/x/text/language` is loaded, and its listed imports are `errors`, `fmt`, `golang.org/x/text/internal/language`, `golang.org/x/text/internal/language/compact`, `sort`, `strconv`, `strings`; no `unicode/norm` edge appears. |
| `command -v govulncheck` | No path/output; `govulncheck` is not installed or available on `PATH`. No scanner callgraph result is available. |

Applicability assessment has **two distinct copies**. The selected module v0.14.0 is version-affected but supplies `x/text/language`, not the unprefixed affected `unicode/norm` package, in the examined application graph. That does not prove the entire CLI excludes normalization code: the compiled Go standard library imports a separate vendor copy, documented below. No whole-binary “not affected,” no demonstrated exploit and no scanner callgraph pass is claimed. The archived SBOM warning remains disclosed.

Recommendation for PG-01: record the module and toolchain findings separately and keep security disposition **pending before preview publication**. A scoped module update to `golang.org/x/text >= v0.39.0` can remove the selected-module version finding after compatibility verification, but **does not replace Go's standard-library vendor copy**. The toolchain copy requires patch/reachability assessment or a verified fixed toolchain under explicit owner approval; do not offer a go.mod-only update as complete remediation. No dependency/toolchain upgrade or suppression is performed here.

Limits: the public record was retrieved from pkg.go.dev through Firecrawl CLI on 2026-10-07; its publication/modification dates do not establish the archived SBOM scan date. Module checksums are not publisher-authenticity/license approval. `govulncheck` is unavailable and was not installed automatically; this investigation instead records real build symbols, four-target package graphs and patch-source evidence. A future scanner result must still distinguish the module and Go-vendor scopes rather than treating scanner absence or zero module findings as safe. Owner security acceptance remains pending.

### Actual binary and standard-library vendor follow-through

Platform Darwin 23.6.0 arm64, installed Go 1.27.1. Temporary `go build -trimpath -o /tmp/deprail-h05-008-applicability ./cmd/deprail` exited 0. Actual `--version` exited 0 (`development`, tag/commit `unknown`); `discover testdata/fixtures/mixed-repository --format json` exited 0 with three complete npm/Python/Maven workspaces and no diagnostics. This is discovery/runtime proof, **not** a vulnerability scan or final release workflow.

Binary SHA-256 `f36c34b2648e2c1b25c26d30ca8420c21ebf030b149797b3bf3f66e0a8d753bc`. `go version -m` records source a01df71368708ce8081c6051116ad436695c5c45 and `vcs.modified=true` because evidence documents were being updated; it is an analysis binary, not an owner-selected immutable release candidate. Its module list still includes x/text v0.14.0.

| Analysis | Observed result and interpretation |
| --- | --- |
| `CGO_ENABLED=0 GOOS=<os> GOARCH=<arch> go list -deps ./cmd/deprail` for linux/amd64, darwin/amd64, darwin/arm64, windows/amd64 | All four exclude exact unprefixed `golang.org/x/text/unicode/norm` but include `vendor/golang.org/x/text/unicode/norm` and `vendor/golang.org/x/net/idna`. These are import graphs, not four native runtime passes. |
| `go mod why vendor/golang.org/x/text/unicode/norm` | `internal/transport/localhttp` → `net/http` → `vendor/golang.org/x/net/idna` → `vendor/golang.org/x/text/unicode/norm`. This is a package-dependency path, not proof an attacker can invoke a vulnerable function. |
| `go tool nm /tmp/deprail-h05-008-applicability` | 63 symbols in the `vendor/golang.org/x/text/unicode/norm` namespace, including `Form.Bytes`, `Form.String`, `compInfo` and `doNormComposed`; no unprefixed norm symbol namespace. An initial substring count included vendor symbols and was corrected by namespace inspection. |
| Installed `$GOROOT/src/vendor/modules.txt` | Go 1.27.1 declares vendored x/text **v0.37.0**, below the module advisory's v0.39.0 floor. `go version -m`'s module list alone does not expose this standard-library copy. |
| Local vendor `forminfo.go:249–267`, `iter.go:427–441` versus [submitted fix CL 794100](https://go.dev/cl/794100) / [commit 5ae8e578](https://go.googlesource.com/text/+/5ae8e578e495731553eddba11b2d0e86c91a00ce) | Upstream fixes invalid-character properties to size 1 rather than assuming nonzero character size. Installed `compInfo` retains `size: uint8(sz)` without that invalid-size rewrite; `doNormComposed` advances by `int(i.info.size)`. This is evidence against assuming the SDK copy is fixed from its Go version label alone; no complete patch-equivalence or malicious-input callgraph proof is claimed. |

[Upstream issue #80142](https://go.dev/issue/80142) and submitted CL were fetched through Firecrawl on 2026-10-07. The fix merged 2026-06-26; the issue also records a request about Go release-branch vendoring. No third-party suppression in that discussion is adopted as DepRail evidence or policy.

**Proposed disposition, not owner acceptance:** module copy's affected package is unused, but standard-library normalization code is linked and its fixed status/runtime attacker-input reachability is not established. PG-01 remains a preview security gate. Before publication, either verify a fixed approved toolchain plus selected-module remediation, or supply a precise call/patch assessment and an explicit source/binary-bounded owner security decision. Do not mark the entire binary nonaffected, upgrade the toolchain implicitly, suppress the advisory, or publish on import-prefix absence.

## Additive repaired-candidate assessment — H05-008-FU1

Technical follow-up [#460](https://github.com/geoffrey-xiao/deprail/issues/460) / [repair evidence](H05-008-FU1-XTEXT-EVIDENCE.md) upgrades the selected module to minimum fixed v0.39.0; bounded real old-v0.14.0 iteration fails to advance, while all12 fixed malformed/valid/empty NFC/NFKC byte/string cases pass. Module checksums, existing full/integration and actual binary/discovery/packaged-console smoke are recorded. The old investigation above remains historical, including its uncertainty and 22-symbol prose count; the official machine record actually lists25.

New [Go member assessment](https://github.com/golang/go/issues/80142#issuecomment-5121323230), exact SDK idna callers and repaired-binary symbol inventory support the specific **unused vulnerable iterator trigger** in the standard-library vendor scope. Go1.27.1/vendored x/text0.37.0 is retained: Form/private norm symbols remain linked, but there are no Iter.Init/InitString/Next/Seek entry points; no patched-vendor, blanket nonaffected, scanner-clearance or permanent-suppression claim. Go.mod-only removal of the version warning is now paired with this precise applicability assessment, not presented alone.

The repair is a review candidate, not yet owner-merged/security-accepted or published. #421/#427 criteria and other release gaps remain open. Actual no-scanner history-save smoke reported HISTORY_WRITE_FAILED with both old and fixed modules; the dedicated evidence preserves that independent integrated-workflow limit under PG-05.
