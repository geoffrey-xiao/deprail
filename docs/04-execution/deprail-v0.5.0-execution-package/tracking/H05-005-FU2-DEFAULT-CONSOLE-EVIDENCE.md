# H05-005-FU2 Default Console Artifact Evidence

Issue [#452](https://github.com/geoffrey-xiao/deprail/issues/452); [contract](../issues/H05-005-FU2-default-artifact-verifier.md). Baseline synchronized owner-merged main d6362f3; owner/reviewer @geoffrey-xiao. Date2026-10-07. Owner acceptance and release approval are not inferred.

## Problem and correction

Owner local testing reported failed detail requests using default web --open. Inspection found a nil *artifact.Store assigned to HistoryService.Artifacts, producing a nonnil interface and calling Verify through a nil receiver. Existing H05-005/FR-503 and failure/data contracts require unavailable integrity without a resolver, not a crash or falsely verified evidence.

cmd/deprail/web.go now declares its optional resolver using the existing app.ArtifactVerifier port. Without a configured root the interface remains nil; an explicitly validated root still installs the same artifact.Store. No artifact.Store nil special-case, reflection, inferred root, extra filesystem access, retry, schema/API/storage/permission or dependency change.

## Regression and repository verification

- go test ./cmd/deprail -run '^TestWebDetailArtifactIntegrity$' -count=1 -v: passed all four cases (no root unavailable, matching verified, missing missing, corrupt digest_mismatch).
- Regression uses actual validated digest-bearing history storage and the CLI's real authenticated HTTP listener. Detail retains entry identity/report/finding count/digest; subsequent findings read returns the real fixture's lodash finding. No mock resolver or source-text/wiring assertions.
- env PATH="/tmp/node-v22.23.3-darwin-arm64/bin:$PATH" make verify: exit0, pinned frontend build, generation, vet, complete Go test suite and build passed.
- Same pinned environment make test-integration: exit0; full integration-tagged suite passed.
- Existing reported failure is prior evidence; no check was rerun merely to confirm the owner's observation.

### Stacked CI correction

Original head bc4d732 failed Windows timed-apply timeout in [37569906460](https://github.com/geoffrey-xiao/deprail/actions/runs/37569906460); preserve that failure. PR #453 is rebased onto unmerged #455 head 5857187ca196ccbd1165c8493ba7916d05689e25, whose three-OS tests/builds/native CLI smoke passed in [37571076086](https://github.com/geoffrey-xiao/deprail/actions/runs/37571076086). Temporary PR base is test/v05-qa-004-timed-apply, keeping artifact-verifier changes isolated in the review diff. Combined-head verification is recorded on #453, not inferred from #455. Owner must review/merge #455 first, then retarget #453 to main and review its final diff/checks. This stack is not an owner-reviewed main or a formal release candidate; no merge or release approval is inferred.

### Main cutover correction

Owner merged #455 into main as 9181e769ed4e5d3faee752a5481d79e3327e7e5a, then merged #453 as 983c604dc49f9df389609233030578bdbe4dd820 into its still-temporary test/v05-qa-004-timed-apply base. The artifact correction therefore did not enter main. Combined #453 head 1644f30 passed [37578685087](https://github.com/geoffrey-xiao/deprail/actions/runs/37578685087) on Ubuntu/macOS/Windows, including native CLI smoke and Ubuntu four-target cross-build; local focused real-process/HTTP scenarios, make verify and integration passed. That evidence is retained, not main-cutover proof.

Fresh branch fix/h05-005-artifact-main-cutover starts from synchronized main 9181e76 and cherry-picks only the six-file #453 squash commit as c2e5be0. No timed-apply duplicate, new behavior or commit to an already-merged branch. A new main-targeted PR under #452 requires its own exact-head verification and owner review/merge. #452 remains open until main delivery and acceptance are evidenced; no release approval is inferred.

Main-based branch local verification: go test ./cmd/deprail -run '^(TestWebDetailArtifactIntegrity|TestApplyEndToEndFailureBoundaries)$' -count=1 -v passed all four artifact states and five real-process failure scenarios. Cancellation observed mutation readiness, exit3/cancelled; timeout exit3/partial; both retained valid durable evidence, succeeded cleanup, unchanged caller and removed worktree. Pinned Node22.23.3/npm10.9.9 make verify and make test-integration exited0. These are newly exercised main-based branch results, not final main delivery or a published binary.

## Actual packaged console smoke

Native development binary: local_test/0.5.0-preview.1/output/default-console-452/bin/deprail; go build -trimpath; SHA256693eb8b9310f050ed26b39cd51d39dd420392c87dabb5ed198e51b848cb4605d. Darwin/arm64, Go1.27.1; Node22.23.3/npm10.9.9; Chrome150.0.7871.24 managed headless,1440x1050 CSS viewport. Development identity only, not release-candidate approval.

Used a separate HOME copied from the stopped four-entry real JS/Python/Java fixture-history snapshot. User's active history/service was not modified or stopped. Launchers privately captured bootstrap URLs; browser scrubbed token fragment; private launcher/bootstrap files removed; tabs closed; both owned processes stopped exit0.

1. HOME=<isolated profile> PATH=<private launcher>:$PATH <binary> web --open, with no artifact-root flag. Actual list→npm-basic detail displayed five genuine captured findings, original digest a284b8766be38be9f8e9d36f2d779bd02accaea3154ae0d41f88a739e3ad2192, Unavailable and the unverified-evidence warning. Browser Back→detail repeated successfully. No JS errors or horizontal overflow. [Default-root screenshot](H05-005-FU2-no-root-detail.webp).
2. Same binary/profile web --open --artifact-root local_test/0.5.0-preview.1/output/run-20261007-103425/repos/npm-basic/.deprail/artifacts. Same record/digest and five findings; Verified and no unverified warning; no JS errors. [Configured-root screenshot](H05-005-FU2-with-root-detail.webp).
3. SHA256 of the source stopped DB and smoke DB after both console sessions matched exactly: fa942844dcd8484d29fd437779d3a0149395cbd47eb081a49c4e7de1e17c4607. This proves copied database bytes unchanged, not a broader release tree/permissions/platform recovery claim.

## Review, risk and rollback

Owner reviews absent versus explicit resolver behavior and every issue criterion. Cross-platform exact-head CI is recorded with the review PR; Darwin browser proof is not native Windows/Linux/AT proof. Windows cancellation cause, integrated H05-008, dependency advisory/attribution/provenance and other prerequisite decisions remain separate. No Master Checklist or acceptance checkmarks changed. Revert transport construction/regression and rebuild matching assets to roll back; retain history/raw artifacts. No merge/tag/publication authorized.
