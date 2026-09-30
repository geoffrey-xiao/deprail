# V05-QA-001: Helper Lifecycle Investigation

## Decision and authorization boundary

Investigation for [#419](https://github.com/geoffrey-xiao/deprail/issues/419), authorized by the owner's `go` instruction after [PR #428](https://github.com/geoffrey-xiao/deprail/pull/428) merged. **No correction or pre-start risk disposition is approved in this record.** Runtime issues remain unstarted. Owner review must distinguish technical acceptance, security/architecture disposition and the readiness decision.

Observed cold-start delay is sufficient to exceed the existing one-second test deadline. This is evidence of a startup-fixture sensitivity, not proof of the exact cause of the earlier failed run or of a particular OS security subsystem. The separate output-capture finding below prevents claiming that all bounded-process invariants passed.

## Baseline and original failure

- Source baseline: `325aa23627f8b08ad6b0e6535b2a79a9e968a941`, owner-merged #428; branch `investigate/419-verification-helper-timeout` from synchronized main.
- Platform/toolchain: Darwin arm64, `go version go1.27.1 darwin/arm64`.
- Original [planning smoke record](../requirements/TEST-STRATEGY.md#11-planning-contract-smoke-evidence-pr-417): `make verify` exit **2**, `TestRunExecutesSelectedCommandsInWorkspace` failed with `SCANNER_TIMEOUT: process exceeded configured deadline`. Generation/vet passed; separate `make build` exit 0. This result remains failed; it was not rerun to confirm or overwritten with a successful sample.
- `run_test.go:14–29` copies the currently running test executable to a new `npm`/`npm.exe` file for each fixture. Line 52 gives `verification.Run` one second. The `ok` helper only writes `ok\n`; it has no intentional sleep.
- `process.Run` starts its deadline before `exec.CommandContext(...).Run()`, so executable startup consumes the same budget as command execution. No production/test deadline, assertion, retry or scanner semantics changed.

## Exact diagnostic commands and method

```text
git fetch origin main
git switch main
git pull --ff-only origin main
git switch -c investigate/419-verification-helper-timeout main
go version
uname -sm
go test -c -o /tmp/deprail-419-verification.test ./internal/remediation/verification
go build -o /tmp/deprail-419-diagnostic/npm ./.deprail-419-diagnostic
/tmp/deprail-419-diagnostic/npm /tmp/deprail-419-verification.test
/usr/bin/codesign -dv --verbose=2 /tmp/deprail-419-verification.test
/usr/bin/log show --last 10m --style json --predicate 'process == "syspolicyd" AND eventMessage CONTAINS "deprail-419"'
go build -o /tmp/deprail-419-diagnostic/npm ./.deprail-419-diagnostic
/tmp/deprail-419-diagnostic/npm --bounds
```

All compile/run commands exited 0. Compilation used the unchanged package and did not execute the reported failing test. Only `-test.run=^TestVerificationHelper$ -- ok/fail/sleep` helper modes were selected for direct observations. Temporary harness/source/binaries were removed after collecting evidence; no permanent tests were added.

Direct observer: `subprocess.Popen([helper, '-test.run=^TestVerificationHelper$', '--', mode], cwd=controlled_workspace, env=approved_environment, stdout=PIPE, stderr=PIPE, start_new_session=True)`. A monotonic timer records Popen return, selector-readiness on stdout, and communicate/reap completion. For `ok`, stdout was exactly `ok\nPASS\n`, stderr empty and exit 0, so the readiness event corresponds to real output, not an empty EOF. The observer waits up to five seconds to measure natural startup; this is not a changed production/test timeout. No network, shell command, package install, raw environment dump or privileged log access was used.

The environment retained only the production allowlist keys and set HOME/TEMP/TMP/USERPROFILE to the controlled workspace. Each fresh copy was byte-identical to the compiled helper, mode 0700. Each warm observation immediately reused that copy. Six separate symlink names pointed to the same already-executed compiled helper. For the cancellation observation, SIGKILL targeted the isolated process group after 20ms without output, then the child was reaped.

Compiled helper: **4,756,434 bytes**, SHA-256 `c68bc2a59b6092d11f229cfd0cec058d9df5b2492b2fdf3cd76250b8b963ae57`. Mach-O arm64, ad-hoc/linker-signed; no TeamIdentifier. No signature was modified.

## Direct startup observations

Times are milliseconds, measured on this host; six samples are not a statistical reliability guarantee.

| Scenario | Popen return | First stdout | Completion | Exit |
|---|---:|---:|---:|---:|
| Original compiled helper, first launch | 3.527 | 3041.486 | 3041.880 | 0 |
| Fresh copy 0 | 6.194 | 203.056 | 203.394 | 0 |
| Fresh copy 1 | 2.695 | 173.334 | 173.668 | 0 |
| Fresh copy 2 | 3.454 | 174.665 | 175.226 | 0 |
| Fresh copy 3 | 3.440 | 3041.953 | 3042.433 | 0 |
| Fresh copy 4 | 3.193 | 203.367 | 203.681 | 0 |
| Fresh copy 5 | 3.498 | 185.371 | 185.683 | 0 |

The corresponding same-copy warm completions were **4.544, 6.199, 5.269, 10.098, 4.891, 5.270ms**. Six first launches through newly created symlinks completed in **6.397, 4.573, 4.064, 4.338, 6.455, 5.223ms**, each with `ok\nPASS\n`, empty stderr and exit 0.

An explicit `fail` helper exited **7** with `failed\n` on stderr in **4.539ms**. A deliberately sleeping helper was killed at the 20ms observation boundary, exited **-9**, and was reaped **0.434ms** after the kill request; total **23.981ms**. This direct observer result is distinct from production runner classification.

## Actual unchanged verification/process boundary

The throwaway Go harness called the real `verification.Run` with the same command kind, working directory `.`, one-second deadline, 1,024-byte cap and copied helper bytes as the existing test. Eight fresh/warm pairs, each a separate controlled scenario rather than an automatic retry:

| Pair | Fresh completion (ms) | Same-copy warm completion (ms) |
|---|---:|---:|
| 0 | 199.832 | 2.929 |
| 1 | 289.257 | 4.467 |
| 2 | 194.433 | 2.830 |
| 3 | 195.241 | 3.723 |
| 4 | 187.285 | 2.677 |
| 5 | 199.995 | 2.986 |
| 6 | 248.929 | 4.837 |
| 7 | 197.351 | 2.946 |

All sixteen returned exit 0 and `ok\nPASS\n`. This shows that these scenarios fit the deadline; it does **not** replace the original full-suite failure or demonstrate that all cold launches do.

The harness also used a real auxiliary child that compares its actual cwd against the canonical requested workspace, prints `ready workspace_match=true`, and optionally sleeps or emits 128 bytes:

| Actual process.Run scenario | Elapsed (ms) | Observed result |
|---|---:|---|
| Workspace/startup | 2.981 | Exit 0; `ready workspace_match=true` |
| Intentional one-second sleep, 20ms deadline | 20.979 | Ready marker proves child ran in workspace; `SCANNER_TIMEOUT` |
| Pre-cancelled context | 0.039 | No stdout; `SCANNER_CANCELLED` |
| OutputCap=32, ready marker plus 128 bytes | 3.586 | `SCANNER_OUTPUT_LIMIT`; capture bound did not hold, see below |

This is smoke evidence of real subprocess behavior, not mock forwarding or source-text testing.

## Causal conclusion and limits

Observed long delay lies **after Popen returns and before the helper's first stdout**. The `ok` branch contains no sleep, and reuse/symlink observations avoid that delay in this sample. [INFERENCE] Repeated executable copying introduces cold-image startup sensitivity unrelated to the intended workspace-command assertion. A one-second budget can therefore kill a correct helper before that assertion becomes observable.

The earlier failed run contains no lifecycle timing, so its exact failure phase cannot be reconstructed. Go/runtime initialization, filesystem/security checks and scheduling are not separately isolated by this observer. The read-only syspolicyd query exited 0 but returned **zero matching entries**; there is no evidence here that Gatekeeper/AMFI/syspolicyd specifically caused the delay. No logs or security controls were disabled, cleared or bypassed.

## Separate output-capture finding — unresolved

A dedicated second smoke called the unchanged `process.Run` with `OutputCap=32`; actual JSON output was:

```json
{"captured_stdout_bytes":155,"configured_cap":32,"error":"SCANNER_OUTPUT_LIMIT: process output exceeded configured limit"}
```

Classification alone is not proof of bounded memory/capture. This finding is separate from the original helper timeout and is **not fixed or dispositioned** by this investigation. `limitedBuffer` embeds `*bytes.Buffer` and defines a bounded `Write`; [INFERENCE] a promoted `ReadFrom` fast path can bypass that method when the exec pipe is copied. That mechanism still requires its own focused diagnosis/regression and bounded correction review; no such production change is authorized here. The observation and risk are tracked in #419 so they cannot be lost during a readiness decision. Owner security/architecture review must account for this known boundary failure; do not describe the process boundary as fully verified or quietly declare the runtime gate safe.

## Proposed bounded correction and next decision

Recommended timeout-fixture proposal, **not applied**: on darwin only, have `testTool` create the allowlisted `npm` symlink to the existing running test binary rather than copy executable bytes. Keep Linux/Windows fixture behavior, one-second happy-path deadline, 20ms timeout, cancellation/output-limit/workspace assertions and production runner unchanged. This is supported by the six symlink observations but still requires owner authorization and tests against the actual fixture; it is not proven as a permanent fix.

After the owner authorizes a correction, the exact verification gate is targeted verification/process behavior checks and `make verify` on the synchronized corrected baseline, then three-OS CI and owner review. New results complement rather than erase the original failed run. Investigate/review the separate capture-bound finding in an explicitly scoped issue before any assertion that the external-process security boundary is acceptable. If choosing a residual-risk disposition instead, the owner must explicitly name scope, rationale, remaining uncertainty/security finding and the next required gate; this record does not grant that disposition.

Current acceptance: diagnostic evidence delivered; original failure preserved; correction, owner security/architecture decision and runtime pre-start disposition **remain pending**. No H05 issue starts, issue closure, release approval or merge authorization follows from these observations.

## Additive update — 2026-10-01

The two follow-up corrections have since been owner-merged: bounded capture in [PR #431](https://github.com/geoffrey-xiao/deprail/pull/431), with real-child cap regressions and three-OS CI run [36720770813](https://github.com/geoffrey-xiao/deprail/actions/runs/36720770813); and the Darwin helper fixture in [PR #433](https://github.com/geoffrey-xiao/deprail/pull/433), with actual fixture probes, corrected-baseline `make verify`, and three-OS CI run [36725991235](https://github.com/geoffrey-xiao/deprail/actions/runs/36725991235). Owner acceptance of the fixture and its residual cold-start uncertainty is recorded in [#432](https://github.com/geoffrey-xiao/deprail/issues/432#issuecomment-5913275072). The corrected integrated branch later passed local `make verify` and exact-head three-OS CI in owner-merged [PR #434](https://github.com/geoffrey-xiao/deprail/pull/434), run [36740438976](https://github.com/geoffrey-xiao/deprail/actions/runs/36740438976).

The owner accepted the #430 capture correction's security/architecture assessment in the current conversation; it is recorded in [V05-QA-002-EVIDENCE](V05-QA-002-EVIDENCE.md#owner-securityarchitecture-decision---2026-10-01). The combined #419 readiness disposition remains pending owner review/merge of the closeout PR. This updates the historical status without relabeling the original failure. Keep #421 blocked until the disposition is accepted.
