# V05-QA-002: Bounded Subprocess Capture Evidence

## Scope and identity

Primary issue [#430](https://github.com/geoffrey-xiao/deprail/issues/430), [local contract](../issues/V05-QA-002-bounded-output-capture.md), native child of #418, target `0.5.0`, milestone `v0.5.0`, Sprint 4, R3/P0. Owner and required reviewer `@geoffrey-xiao`; independent reviewer optional under ADR-0005.

Owner authorized the announced bounded capture repair with “请继续我们的任务”. Branch `fix/430-bounded-process-output` was created from synchronized owner-merged baseline `8febe78f3bbf04ca400c5e68d2e302797ab1a31c`. Platform: Darwin arm64; `go version go1.27.1 darwin/arm64`. No new version line, dependency, deadline, retry, error code, scanner/API contract, history/UI feature or merge authorization.

The earlier [#419 investigation](V05-QA-001-INVESTIGATION.md) remains unchanged. Its 155/32-byte finding was accepted as investigation, not as a passing bound or runtime-start disposition. Closed #419/Project Done does not itself resolve the original timeout readiness gate.

## Causal reproduction and smallest correction

A throwaway in-package diagnostic copied 155 deterministic bytes from `io.Pipe` into the actual `limitedBuffer` with limit 32, through `io.Copy`:

```text
Before: promoted_reader_from=true configured_cap=32 copied=155 captured=155 error=<nil>
After:  promoted_reader_from=false configured_cap=32 copied=32 captured=32 error=short buffer
```

This upgrades the earlier mechanism inference to an exercised cause: anonymous `*bytes.Buffer` embedding promotes `ReadFrom`, allowing a copy fast path to bypass the bounded `Write`. The correction uses a named `Buffer *bytes.Buffer` field and reads `b.Buffer.Len()`. Existing constructors and bounded Write remain; no forwarding method, new buffer, allocation/copy layer or alternative scanner path is introduced. Package documentation now states the per-stream bound and existing conservative at-cap error.

No Run error precedence, deadline/cancellation logic, environment allowlist, process grouping or public types changed. Reaching the cap still produces `SCANNER_OUTPUT_LIMIT`; this issue does not redefine exact-cap success semantics. Temporary type/fast-path diagnostics were removed, not kept as implementation-detail permanent tests.

## Permanent real-child regression: failing before, passing after

Added `TestRunBoundsCapturedOutput` uses the existing test executable as a real child, independently for stdout/stderr. The controlled helper exits without a testing-framework trailer. Cases retain exact deterministic prefixes for 31, 32, 33 bytes, a single 256KiB write and chunked 256KiB output, each with cap 32. The other stream remains empty; below-cap success and current at/over-cap failure classification are asserted.

```text
go test ./internal/process -run '^(TestCapturePipeDiagnostic|TestRunBoundsCapturedOutput)$' -v
```

- Before correction: exit **1**. Six real-child cases failed: stdout/stderr 33-byte, large-single-write and large-chunked cases. Actual retained sizes were **33** or **262,144**, expected **32**. Below/at cases were not failing; they guard unchanged behavior.
- After correction: exit **0**, all ten real-child boundary cases passed; pipe probe observed 32 retained bytes and `short buffer`.
- Permanent suite contains consumer-visible output/boundary assertions, not a promoted-interface/source-text test. Existing argument-array, missing-tool, timeout and output-limit tests remain unchanged.

## Actual production-runner smoke

A throwaway compiled Go caller invoked real `process.Run`, with independent stdout/stderr limit 32, real cwd, argument arrays and controlled child branches. Captured prefixes/errors were explicitly checked by the program; it exited **0**.

| Scenario | Retained stdout / stderr bytes | Exit field | Observed error | Elapsed ms |
|---|---|---:|---|---:|
| 155 bytes on stdout | 32 / 0 | 0 | SCANNER_OUTPUT_LIMIT | 8.177 |
| 155 bytes on stderr | 0 / 32 | 0 | SCANNER_OUTPUT_LIMIT | 5.901 |
| Small stdout/stderr and correct cwd | 6 / 5 | 0 | None | 5.210 |
| Explicit failure | 0 / 7 | 7 | SCANNER_EXIT_NONZERO | 5.127 |
| Ready marker then intentional sleep, 20ms deadline | 6 / 0 | 0 | SCANNER_TIMEOUT | 21.740 |
| Ready marker then sleep, cancel after 20ms | 6 / 0 | 0 | SCANNER_CANCELLED | 21.490 |
| Pre-cancelled context | 0 / 0 | 0 | SCANNER_CANCELLED | 0.024 |

The Exit field on classified timeout/cancellation/output-limit results retains existing behavior; it is not proof that those children exited successfully. No captured stream exceeded 32 bytes. Small stdout=`hello\n`, stderr=`done\n` and nonzero stderr=`failed\n` were exact. The child checked its cwd against the canonical controlled workspace.

## Actual CLI and race-instrumented CLI smoke

Built the actual `cmd/deprail` CLI from the corrected working tree, plus two offline controlled Go scanner binaries. Both expose OSV-Scanner version 2.0.0; the ordinary fixture emits `{"results":[]}`, while the noisy fixture attempts 32MiB of stdout in reusable 1KiB chunks. This is actual CLI/adapter/process execution with synthetic scanner fixtures, **not a live OSV/database or release workflow claim**.

Using isolated fixture workspace with `package.json` and `package-lock.json`, and selecting the controlled scanner via PATH:

```text
go build -o "$SMOKE/bin/runner" ./.deprail-430-smoke
"$SMOKE/bin/runner" "$SMOKE/workspace"
go build -o "$SMOKE/bin/deprail" ./cmd/deprail
go build -o "$SMOKE/bin/normal/osv-scanner" "$SMOKE/scanner.go"
go build -ldflags '-X main.captureMode=large' -o "$SMOKE/bin/large/osv-scanner" "$SMOKE/scanner.go"
PATH="$SMOKE/bin/large:$PATH" "$SMOKE/bin/deprail" scan "$SMOKE/workspace" --format json
PATH="$SMOKE/bin/normal:$PATH" "$SMOKE/bin/deprail" scan "$SMOKE/workspace" --format json
go build -race -o "$SMOKE/bin/deprail-race" ./cmd/deprail
PATH="$SMOKE/bin/large:$PATH" "$SMOKE/bin/deprail-race" scan "$SMOKE/workspace" --format json
PATH="$SMOKE/bin/normal:$PATH" "$SMOKE/bin/deprail-race" scan "$SMOKE/workspace" --format json
```

`$SMOKE` denotes the unique owned temporary directory; only its private prefix is replaced here. Runtime capture uses the unchanged application cap of 16MiB per stream. The fixture mode is selected at link time, not through an unapproved child environment variable; the environment allowlist is preserved.

| Executable | Fixture | Exit | JSON status | Findings | Errors / stderr |
|---|---|---:|---|---:|---|
| Actual CLI | Noisy | **3** | **failed** | 0 | `SCANNER_OUTPUT_LIMIT: process output exceeded configured limit`; stderr empty |
| Actual CLI | Ordinary | **0** | **complete** | 0 | Errors empty; stderr empty |
| Race-instrumented actual CLI | Noisy | **3** | **failed** | 0 | Same explicit output-limit failure; no race diagnostics |
| Race-instrumented actual CLI | Ordinary | **0** | **complete** | 0 | Errors empty; no race diagnostics |

Noisy zero findings are not reported as a safe/complete result. Existing scan JSON `errors` remain strings; no schema migration was inferred. Race smoke exercises the real CLI/process copy path; it is not a claim that a full `go test -race ./...` suite ran.

Local non-release CLI SHA-256: `419824e45ed7db5819542f0eed3007135a4aaae1ef7386e63683074f863f9697`. Source manifest SHA-256 `130a373c1c70746a1a260d0f11d68ac6ba583ce4c6bf3587481b70139a9ee50f`; lockfile SHA-256 `a44fa33038526c5a1e62febd1746578824aba11c1b99086c800261bc0ec3ccda` remained unchanged after normal/race/noisy runs. Expected `.deprail` artifact writes occurred only in the owned temporary workspace; no user's repository was scanned or mutated. Temporary source, binaries, fixtures and artifacts were removed after proof.

## Full local verification — known failure remains

Ran **one** `make verify` against the corrected working tree. Exit **2**:

- Generation, formatting check and vet passed.
- New/process tests passed (`internal/process`: 3.841s); CLI, OSV adapter, app and other reported packages passed.
- Unchanged `internal/remediation/verification.TestRunExecutesSelectedCommandsInWorkspace` failed at `run_test.go:54`: `verification "test" failed: SCANNER_TIMEOUT: process exceeded configured deadline` (1.00s).
- Build step was not reached by make after the test failure; the separate actual CLI/harness/race builds above exited 0.

This newly changed-baseline result complements the original failed run. It does not erase it, establish an OS root cause, justify a retry/deadline increase, authorize a darwin helper-fixture change or satisfy the owner runtime-start gate. The failure was not rerun to obtain a green result. Full local verification remains failing independently of the now-passing capture regression.

## Review and residual risk

Code/regression/smoke evidence is delivered for #430; technical review, separate owner security/architecture decision and owner-reviewed merge remain required. CI results must be linked separately for the submitted head. No local/full-suite pass or owner acceptance is claimed here.

This correction bounds retained stream bytes, not all child resources or exact allocator capacity. Existing fixed-size pipe-copy buffers/allocator slack remain. Process-tree termination timing and error precedence are unchanged; no redesign is implied. Windows/Linux actual execution is still represented only by the eventual CI jobs, not local runtime.

The original helper timeout requires its own bounded correction or explicit owner pre-start disposition. H05 #420–#427 stay unstarted until this capture correction is reviewed/merged and existing readiness gates are explicitly satisfied. This evidence creates no new release approval or automatic merge permission. Rollback would restore the vulnerability; preserve the failing regression and block dependent runtime starts rather than treating rollback as safe readiness.

## Owner security/architecture decision — 2026-10-01

The owner separately accepted the security/architecture assessment of merged [PR #431](https://github.com/geoffrey-xiao/deprail/pull/431) in the current conversation. The correction prevents `io.Copy`'s promoted `ReadFrom` fast path from bypassing the bounded writer, while preserving per-stream caps, existing output-limit classification, deadlines and cancellation. Real-child stdout/stderr tests cover below/at/over-cap and large single/chunked writes; runner and ordinary/noisy/race CLI smoke plus exact-head three-OS CI are linked above.

Accepted residual risk: retained byte counts are bounded; total child CPU/output and allocator-capacity overhead are not. This separate assessment does not resolve the copied-helper timing uncertainty or authorize H05 runtime by itself; #419 remains the global gate.

Durable owner record: [#430 comment](https://github.com/geoffrey-xiao/deprail/issues/430#issuecomment-5921486377).
