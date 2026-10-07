# H05-005-FU2: Preserve Unavailable Artifact Integrity Without a Configured Root

GitHub issue: [#452](https://github.com/geoffrey-xiao/deprail/issues/452).

Review PR: [#453](https://github.com/geoffrey-xiao/deprail/pull/453); implementation commit bd79d15. Refs #452; owner acceptance and exact-head CI remain separate from local proof.

Main-cutover follow-up: #453 was owner-merged into its temporary #455 branch, not main. Fresh branch fix/h05-005-artifact-main-cutover carries the same artifact correction onto reviewed main 9181e76; [cutover evidence](../tracking/H05-005-FU2-DEFAULT-CONSOLE-EVIDENCE.md#main-cutover-correction) preserves the merge history. Issue remains open pending the new main-targeted PR, exact-head CI and owner-reviewed main delivery.

Main-targeted review: [#456](https://github.com/geoffrey-xiao/deprail/pull/456), Refs #452; migrated artifact commit c2e5be0 plus cutover evidence cc01b30. Final-head CI and owner review/merge remain separate.

## Planning metadata

- Parent: H05-005 / #424; release verification #427; delivery epic #418.
- Type: bug; area: cli; priority: P0; risk: R2.
- Target version/milestone: v0.5.0; Sprint 4.
- Owner/reviewer: @geoffrey-xiao; independent review optional under ADR-0005.
- Dependencies: owner-merged H05-005 transport and H05-006/007 assets; synchronized main d6362f3.
- Contracts: release plan included local history/detail/error scope; H05-005 Inputs/Outputs; FR-503; FAILURE-AND-DATA-CONTRACT projection mapping requires artifact integrity unavailable without a resolver.
- Blocked reason: none for this existing-contract correction; formal release acceptance remains separately gated.

## Definition of Ready

- [x] Default console detail failure was reported in real local testing and recorded in #427; do not rerun merely to confirm the report.
- [x] Scope, failure semantics, test/scenario and rollback specified below.
- [x] No API/schema/storage/permission or release-scope change; owner/reviewer and dependency baseline explicit.

## Goal

An existing saved scan remains readable through default deprail web --open without --artifact-root; digest references are unavailable rather than causing a panic/connection failure or being falsely verified.

## Scope and exclusions

Correct optional verifier construction in cmd/deprail/web.go using the existing application port. Add a real-listener CLI regression using actual history storage and a digest-bearing entry. Verify configured root still distinguishes verified/missing/mismatch and absent root retains metadata/findings. No nil handling added to artifact.Store, no reflection, inferred repository roots, automatic artifact access, retry, migration, API changes or cancellation correction.

## Inputs, outputs and failure behavior

Inputs: validated saved history with artifact digests; optional existing trusted artifact root. Outputs: existing detail JSON/console, original metadata and digest identities; unavailable integrity when resolver absent. Explicit roots continue using current digest verification. Invalid root still fails closed before listener startup. History/artifacts remain read-only.

## Required tests and acceptance

- [ ] Real CLI listener returns HTTP200 detail with unchanged entry/report/digest and integrity unavailable when no root was supplied; subsequent read remains usable.
- [ ] Explicit trusted root yields verified for matching content and explicit missing/mismatch states for failed integrity reads.
- [ ] Actual compiled console displays unavailable warning without artifact root and verified evidence with explicit root; no JavaScript/API crash or history mutation.
- [ ] Focused Go regression, make verify, make test-integration and exact-head cross-platform CI/evidence recorded.

## Human review and evidence

Owner reviews every criterion and optional artifact-access boundary. Record binary/source identity, exact commands/results, sanitized actual-browser evidence, repository/history preservation and CI links. User's default-path failure is prior evidence; regression and post-fix runtime proof do not replace owner release decisions.

[Default/configured-root runtime and regression evidence](../tracking/H05-005-FU2-DEFAULT-CONSOLE-EVIDENCE.md) includes native binary identity, genuine saved findings, unavailable/verified screenshots and unchanged database hash.

## Final acceptance

- [ ] Owner technical/security effect and residual risk reviewed.
- [ ] Required CI/runtime evidence accepted.
- [ ] Owner-reviewed merge recorded before issue closure/Project Done.

## Rollback

Revert transport construction/test change and rebuild matching assets; preserve DB/WAL/SHM, raw artifacts and immutable release tags. No automatic downgrade or destructive reset.
