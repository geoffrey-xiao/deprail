# V05-QA-004 Timed Apply Evidence

Issue [#454](https://github.com/geoffrey-xiao/deprail/issues/454); [contract](../issues/V05-QA-004-timed-apply-startup.md); parent release verification #427. Synchronized main d6362f3,2026-10-07. Owner/reviewer @geoffrey-xiao; owner acceptance and release decisions pending.

## Observations and cause boundary

- Historical preparation run [37551712270](https://github.com/geoffrey-xiao/deprail/actions/runs/37551712270) at33f7944 and UI run [37564839801](https://github.com/geoffrey-xiao/deprail/actions/runs/37564839801) at7667fa3 failed Windows cancellation with unexpected end of JSON input at apply_e2e_test.go:215. Subsequent passing runs do not erase these observations.
- Default-console run [37569906460](https://github.com/geoffrey-xiao/deprail/actions/runs/37569906460) atbc4d732 passed Ubuntu/macOS but failed Windows timeout at apply_e2e_test.go:225 with the same decode error. Windows later package build/native smoke skipped.
- Existing test cancelled500ms after starting apply, or created a1s context deadline before repository/worktree startup. No controlled mutation readiness was checked. Production preflight errors write stderr and exit3 before a result exists; generic JSON decoding then obscures the actual phase/error.
- Controlled throwaway preflight-cancellation diagnostic using real fixture repository/plan/approval and runFixApplyContext produced: exit3, stdoutBytes0, PATH_OUTSIDE_ROOT=true, worktreeError=false, jsonDecode=unexpected end of JSON input. Diagnostic file was removed. This proves the early-cancellation mechanism, not the exact phase of the original Windows failures: those logs did not capture stderr. No old failed scenario was rerun merely for confirmation, classified as flaky, or suppressed.

## Correction

Only cmd/deprail/apply_e2e_test.go and the existing cmd/deprail/main_test.go controlled mutation helper change runtime test behavior. Production code, JSON/error contracts, process timeouts, permission/environment/approval/worktree behavior remain unchanged.

Controlled helper emits a private ready marker before its existing sleep. Both timed cases wait for actual mutation readiness, early operation completion, filesystem failure or30s startup bound. Early completion/failure remains fatal with exit/stderr context; startup failure cancels/drains the operation. The500ms guess is gone. The1s standard-library deadline now starts only after readiness; a test-only Background-based context bridge propagates the genuine deadline cause through the already-running apply call and derived subprocess contexts. Cancellation remains context.Canceled, deadline remains context.DeadlineExceeded.

The first bridge implementation yielded cancelled for the timeout case; the unchanged partial assertion caught it. Corrected the bridge's context propagation, not the expected outcome. No values exist on its Background source; exposing the internal cancel context would bypass its deadline Err semantics, so the bridge's Value returns nil. Owner should review this small test-only adapter and its derived-context behavior.

Original exit3, cancelled/partial outcome, evidence-path and successful-cleanup assertions remain. Strengthened proof additionally reads/validates persisted evidence with the correct outcome, compares caller repository state, and checks isolated worktree removal. Decode failures now include exit/stderr instead of a bare EOF-like JSON error; they are not accepted or converted to results.

## Verification exercised

Darwin/arm64, Go1.27.1, pinned Node22.23.3/npm10.9.9:

- go test ./cmd/deprail -run '^TestApplyEndToEndFailureBoundaries/(cancellation|timeout)$' -count=1 -v: passed both actual native child/worktree cases. Observed cancellation exit3/outcome cancelled and timeout exit3/outcome partial; readiness observed, validated evidence retained, cleanup succeeded, caller unchanged, worktree removed.
- Same focused command with -race: both passed; no race report.
- Pinned make verify: exit0; frontend, generation, vet, full Go tests/build passed.
- Pinned make test-integration: exit0; full integration-tagged suite passed.
- The changed surface is the native E2E harness: actual apply application/CLI execution and real controlled subprocesses/isolated worktrees were exercised, not a mock result or source assertion. Package-manager fixtures are deliberately offline; these do not claim real npm/pip/Maven remediation or full release browser/AT evidence.

## Review and release boundary

Exact-head three-OS CI/native package smoke recorded with the review PR. Previous failures and the original Windows phase uncertainty remain explicit. After owner-reviewed QA merge, dependent default-console PR #453 must synchronize onto that baseline and obtain its own passing exact-head CI; an isolated QA branch pass is not a release/main verification claim.

Preparation #445 was returned to Draft; #427 Project Blocked with actual unmet formal acceptance/evidence gates, keeping at most two active Review PRs. No issue closure, Master Checklist checkmarks, owner security/risk disposition, tag or publication. Revert test/helper change to roll back without altering production binaries or user data.
