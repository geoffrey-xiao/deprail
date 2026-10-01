# H05-002 Private History Store — Implementation Evidence

**GitHub issue:** [#421](https://github.com/geoffrey-xiao/deprail/issues/421)

**Status:** Draft [PR #436](https://github.com/geoffrey-xiao/deprail/pull/436) open; Project Status `Review`; the first Windows CI run failed on owner verification for a newly created history directory. The safe ownership correction is in the PR follow-up; exact-head CI rerun and owner acceptance remain pending.
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
| `GOOS=windows GOARCH=amd64 go test -c -o /tmp/deprail-history-windows.test.exe ./internal/store/history`; PR #436 CI run [36802315026](https://github.com/geoffrey-xiao/deprail/actions/runs/36802315026) | Cross-compilation passed. Ubuntu/macOS runtime jobs passed; Windows runtime failed because newly created path ownership did not match the token user. Follow-up only assigns current-user ownership to paths created by this process; existing mismatched ownership remains a refusal. New exact-head CI run pending. |
| `go mod tidy` | Passed; SQLite is direct and existing `golang.org/x/sys` is pinned at the compatible `v0.47.0`. |
| SBOM/checksum/source-diff validation | `jq` validated CycloneDX 1.6 structure, 12 dependency components and allowed license IDs; `go mod verify` reported all modules verified; `git diff --check` passed. |
| `make verify` | Passed: `go generate ./...`, `go vet ./...`, `go test ./...`, and `go build ./...` all completed successfully on the local Darwin/arm64 host. |

Local `make verify` passed. Exact-head multi-OS rerun and owner review are still required before issue acceptance. Local POSIX mode tests ran on Darwin; Windows ACL runtime evidence remains pending.

## Limits and unresolved acceptance evidence

- SQLite `max_page_count` is the deterministic full-storage failure injection; it is not an actual filesystem `ENOSPC` test. The child-process test exits after inserting row and refs but before commit, then reopens and verifies neither uncommitted row nor ref survived.
- No older DepRail schema exists to migrate. Current v1 opens without migration; this branch tests creation, unique online snapshots, snapshot integrity, no-overwrite naming and cancellation preservation. It does not invent a fake migration or claim migration-failure/restore rehearsal evidence. Manual restore remains a future release gate; no restore action is implemented.
- The store commits digest references only. Missing/digest-mismatch artifact integrity states and artifact retrieval remain the dependent H05-003 responsibility; this package never reads, deletes or garbage-collects raw artifacts.
- The 256 MiB bound covers summed projection JSON only. SQLite pages/indexes, WAL/SHM, backups and external artifacts can use more disk; no hard filesystem quota is claimed.
- Owner review must separately accept schema/permission/dependency impact, the disclosed `x/text` advisory and logical-versus-physical growth risk. Issue #421 remains open at Project Status `Review` until evidence, owner review and CI are recorded.
