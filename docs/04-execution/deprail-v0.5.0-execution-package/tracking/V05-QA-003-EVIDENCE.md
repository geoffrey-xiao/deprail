# V05-QA-003: Darwin Verification Helper Correction Evidence

## Authorization, baseline and scope

Primary [issue #432](https://github.com/geoffrey-xiao/deprail/issues/432), [local contract](../issues/V05-QA-003-darwin-helper-fixture.md), native child of #418. Target `0.5.0`, milestone `v0.5.0`, Sprint 4, R2/P0, owner and required reviewer `@geoffrey-xiao`; external review optional under ADR-0005.

The owner's “请继续” authorizes the announced test-only fixture correction, not H05 kickoff or automatic merge. Branch `fix/432-darwin-verification-helper` starts from synchronized main `4f9fe3d98aa5cf920cfa3eb443caede67f36589c`, containing owner-merged capture fix PR #431. Issue #430 is now closed/Project Done. This is not verification of an isolated branch missing the capture prerequisite.

Only `internal/remediation/verification/run_test.go:testTool` changes: on darwin, use `os.Executable()` and create the allowlisted `npm` symlink to the running test image. Errors fail fixture creation, with no fallback or skip. Linux/Windows keep the existing byte copy, including `npm.exe`; no original test body, helper mode, assertion, one-second happy-path deadline or 20ms timeout changes. No production verification/process/scanner/CLI, permission, schema, storage, dependency or error contract changes.

The [earlier investigation](V05-QA-001-INVESTIGATION.md) measured cold-image first output over three seconds and new symlinks at 4–7ms; the [capture correction evidence](V05-QA-002-EVIDENCE.md) independently recorded the same copied-helper failure during full local verification. Those failed runs remain failed historical facts. Their exact OS security/scheduling cause remains unknown; this correction removes the avoidable image-copy fixture sensitivity, not claims to identify or bypass an OS security subsystem.

## Actual fixture smoke

A temporary in-package smoke used the **real corrected `testTool` and real verification.Run**, not a mock runner or an independent replacement for the fixture. Inputs were isolated temporary workspace, allowlisted helper path, explicit argument array, timeout and 1,024-byte output cap. Original helper modes supplied success, exit7 and intentional sleep; a temporary child compared its actual cwd with the canonical expected workspace and printed only a boolean, not a host path.

```text
gofmt -w internal/remediation/verification/run_test.go internal/remediation/verification/fixture_smoke_test.go
go test ./internal/remediation/verification -run '^TestActualFixtureSmoke$' -v
```

| Scenario | Deadline ms | Elapsed ms | Exit field | Actual stdout/stderr/error |
|---|---:|---:|---:|---|
| Original `ok` helper | 1000 | 7.495 | 0 | `ok\nPASS\n`, empty stderr, no error |
| Original `fail` helper | 1000 | 5.422 | 7 | Empty stdout, `failed\n`, SCANNER_EXIT_NONZERO |
| Original intentional `sleep` | 20 | 21.687 | 0 | SCANNER_TIMEOUT |
| Pre-cancelled helper | 1000 | 0.175 | 0 | Empty stdout/stderr, SCANNER_CANCELLED |
| Canonical workspace child | 1000 | 5.843 | 0 | `{"workspace_match":true}\nPASS\n`, empty stderr, no error |

Timeout/cancellation Exit fields reflect the existing classified result shape, not successful child execution. The stdout `PASS` trailer is from the original Go test helper; it was neither suppressed nor changed in permanent code.

The first temporary smoke invocation exited **1** at its workspace scenario: its extra child flag was incorrectly outside the Go test `--` separator. The preceding four scenarios above had already produced the recorded results. Corrected only the throwaway harness's argument list and reran the failed workspace scenario, which exited **0** with `workspace_match=true`. This harness error is disclosed, not counted as a successful first run or as a production/fixture timeout. No permanent assertion or deadline was changed to make it pass. The temporary smoke file was removed before full verification; no implementation-specific symlink/source-text test was added.

## Targeted actual behavior suites

```text
go test ./internal/remediation/verification ./internal/process -count=1
```

Exit **0**; verification **0.238s**, process **0.607s**. This uncached run executes the existing workspace command, stop-after-failure/partial results, 20ms timeout, cancellation, untrusted/empty command and workspace escape cases, plus #430's bounded stdout/stderr regression. No known original failed check was rerun merely to confirm it before correcting the fixture.

## Full corrected-baseline verification

After removing the temporary smoke, ran **one**:

```text
make verify
```

Exit **0**, **26.62s**:

- `go generate ./...` passed.
- Formatting check and `go vet ./...` passed.
- `go test ./...` passed, including `cmd/deprail` **23.405s**, application **7.411s**, verification **2.633s** and the process suite (cached in this full run, already executed uncached above).
- `go build -ldflags "-X github.com/geoffrey-xiao/deprail/internal/buildinfo.Version=development -X github.com/geoffrey-xiao/deprail/internal/buildinfo.Tag=unknown -X github.com/geoffrey-xiao/deprail/internal/buildinfo.Commit=unknown" ./...` passed.

This is a new result on corrected source, not a relabeling of the original `make verify` exit2 or a success-by-retry substitution. No full-verification retry was needed. No new deadline, retry, prewarming, test skip or production override was used.

Platform/toolchain: Darwin arm64, pinned Go `1.27.1` (enforced by the successful make bootstrap prerequisite). Linux/Windows native execution is pending the submitted PR's three-OS CI and must be linked separately; cross-compilation is not reported as runtime proof. The locally changed surface is a private test fixture; there is no changed CLI/UI surface or release binary to visually validate.

## Review and runtime gate

Implementation and actual helper/targeted/full evidence are delivered for #432. Owner technical acceptance and separate security/architecture assessment, required CI and owner-reviewed merge remain pending. Production boundaries remain unchanged; the already merged #431 supplies the capture correction.

The owner can record the timeout readiness gate's resolution after reviewing this correction and its exact local/CI outcomes. That decision is distinct from starting any H05 issue, which still requires its own DoR, dependencies and reviewed-main synchronization. No history/UI feature, release approval, automatic merge or owner disposition is inferred from the local pass. Original uncertainty and failed runs remain linked permanently.

Rollback: a reviewed revert restores the cold-copy fixture sensitivity without altering production/user data; preserve the evidence and reopen the relevant readiness gate rather than conceal the failure.

## Additive owner-reviewed completion

Three-platform [CI run 36725991235](https://github.com/geoffrey-xiao/deprail/actions/runs/36725991235) passed on tested head `b6a96ee34b1053b1c860fa9bf22dad17d19c741b`: Ubuntu 26s, macOS 1m5s, Windows 1m40s. Owner-merged [PR #433](https://github.com/geoffrey-xiao/deprail/pull/433) produced reviewed-main `6c74cf9861981928c29f7a8cb09799e63f3fd6dd`.

[Owner acceptance](https://github.com/geoffrey-xiao/deprail/issues/432#issuecomment-5913275072) records the explicit “确认验收并完成关联” decision: all acceptance criteria confirmed, technical acceptance and separate security/architecture assessment recorded, remaining macOS cold-start cause uncertainty accepted. No native APPROVED review is invented. API readback verified #432 closed, Project Done, and native Linked pull requests contains #433. The accidentally introduced custom card-link display and repeated links on nine Todo items were removed; no browser inspection was performed.

Earlier pending language describes review submission. This completion resolves QA, not a blanket runtime or release authorization. The later [H05-001 kickoff](https://github.com/geoffrey-xiao/deprail/issues/420#issuecomment-5913450360) separately records owner “可以 继续吧”, checked DoR and a branch from synchronized reviewed-main. Historical failed checks and unproven OS root cause remain unchanged.
