# v0.5 Approved Implementation Backlog

## Authorization and boundary — 2026-09-30

Owner `@geoffrey-xiao` accepted PR #417 head `b0649547d1e0f033f10f3f28920b74eab9b16e1b` with “approve，请继续”. [Conversation-derived decision record](https://github.com/geoffrey-xiao/deprail/pull/417#issuecomment-5907670878) separately records technical acceptance, security/architecture acceptance and verification disposition; it is not a fabricated native GitHub APPROVED review. Owner already merged #417 at `40e43a83eaac1b84dd038ee1ea4dccbd3e7234f6` before this handoff. New branch `docs/392-approved-implementation-backlog` was created from synchronized main; no commits were appended to the merged branch.

Authorization covers complete planning contracts and GitHub publication only. No runtime code, new-PR merge, dependency introduction, Sprint kickoff or release publication is authorized. Each delivery issue remains **Todo**; required runtime acceptance and owner review remain unchecked. The current handoff requires its own owner PR review.

## Published contracts

At initial publication, all ten rows were Project Todo with Sprint unassigned; owner/reviewer `@geoffrey-xiao`, milestone `v0.5.0` and Target Version `0.5.0` were verified. The owner subsequently requested Sprint assignment; live #418–#427 are Sprint 4. #419 is closed/Project Done for accepted investigation, not an automatic runtime-start disposition. Current #430 is a separate R3 bounded-capture correction in Sprint 4. H05 #420–#427 remain unstarted; their live Project fields are authoritative.

[CSV backlog](issue-backlog.csv) preserves the existing fourteen-column import convention and appends live GitHub number/URL and local contract path. Evidence cells name required future proof, not passed checks.

| Key / durable contract | GitHub | Dependencies / pre-start | Area / risk / priority / type |
|---|---|---|---|
| [V05-EPIC-002](../epics/EPIC-002-local-history-delivery.md) | [#418](https://github.com/geoffrey-xiao/deprail/issues/418) | Readiness #392; quality #419 | foundation / R3 / P0 / feature |
| [V05-QA-001](../issues/V05-QA-001-verification-timeout-disposition.md) | [#419](https://github.com/geoffrey-xiao/deprail/issues/419) | Accepted owner decision / merged #417 | test / R2 / P0 / bug |
| [H05-001 Safe projection](../issues/H05-001-safe-history-projection.md) | [#420](https://github.com/geoffrey-xiao/deprail/issues/420) | #419 | normalization / R3 / P0 / feature |
| [H05-002 Private store](../issues/H05-002-private-history-store.md) | [#421](https://github.com/geoffrey-xiao/deprail/issues/421) | #420; #419 global gate | foundation / R3 / P0 / feature |
| [H05-003 Application services](../issues/H05-003-history-application-services.md) | [#422](https://github.com/geoffrey-xiao/deprail/issues/422) | #420, #421; #419 global gate | foundation / R3 / P0 / feature |
| [H05-004 Opt-in CLI](../issues/H05-004-opt-in-history-cli.md) | [#423](https://github.com/geoffrey-xiao/deprail/issues/423) | #422; #419 global gate | cli / R3 / P0 / feature |
| [H05-005 Local transport](../issues/H05-005-read-only-local-console.md) | [#424](https://github.com/geoffrey-xiao/deprail/issues/424) | #422; #419 global gate | cli / R3 / P0 / feature |
| [H05-006 Accessible UI](../issues/H05-006-accessible-history-ui.md) | [#425](https://github.com/geoffrey-xiao/deprail/issues/425) | #424; #419 global gate | foundation / R2 / P0 / feature |
| [H05-007 Embedded packaging](../issues/H05-007-embedded-console-packaging.md) | [#426](https://github.com/geoffrey-xiao/deprail/issues/426) | #424, #425; #419 global gate | foundation / R2 / P0 / feature |
| [H05-008 Release evidence](../issues/H05-008-local-history-release-evidence.md) | [#427](https://github.com/geoffrey-xiao/deprail/issues/427) | #420–#426; #419 global gate | test / R3 / P0 / test |
| [V05-QA-002 Bounded capture](../issues/V05-QA-002-bounded-output-capture.md) | [#430](https://github.com/geoffrey-xiao/deprail/issues/430) | Accepted #419 investigation / merged #429; prerequisite to runtime-start security decision | foundation / R3 / P0 / bug |
| [V05-QA-003 Darwin fixture](../issues/V05-QA-003-darwin-helper-fixture.md) | [#432](https://github.com/geoffrey-xiao/deprail/issues/432) | Accepted #419 investigation / merged #429; capture correction #430 / merged #431 | test / R2 / P0 / bug |

This is a dependency DAG, not permission to start all issues. Work one issue at a time, at most two implementation issues/two review PRs active, never two R3 implementations concurrently. H05-006 can exercise the real H05-005 API through a throwaway embedding harness before H05-007 produces production assets; no shipped placeholder assets or cyclic release dependency. Store API close/reopen evidence belongs to H05-002; CLI capture is H05-004. HTTP owns wire envelopes/cursors; application services return typed domain results and metadata-only integrity outcomes. Unknown provenance remains null rather than causing blanket rejection.

## Pre-start and residual-risk gates

- #419 preserves the observed local `make verify` helper-process timeout as unresolved. Three-OS CI from the accepted proposal does not erase it. Resolve with evidence or obtain explicit owner pre-start disposition before runtime implementation; do not rerun merely to confirm, skip the assertion, blindly lengthen the deadline or assume a flake.
- #419 [controlled investigation evidence](V05-QA-001-INVESTIGATION.md) is owner-merged in PR #429; its original failed runs remain historical facts. Capture correction #430 is closed/Project Done after owner-merged PR #431 at `4f9fe3d`; it is not an isolated feature-branch prerequisite anymore.
- [#432 darwin fixture correction](../issues/V05-QA-003-darwin-helper-fixture.md) now implements the previously proposed active-image symlink for tests only. [Corrected-baseline evidence](V05-QA-003-EVIDENCE.md) records targeted and full verification success with unchanged deadlines/assertions. Owner review/merge and explicit runtime-start decision remain required; new passing checks complement rather than relabel the original failures.
- Candidate dependencies remain candidates. Complete direct/transitive license, vulnerability and provenance review before introduction; no automatic install or lifecycle scripts.
- Accepted local-browser fragment exposure and logical-versus-physical DB-growth risks remain explicitly bounded by the accepted contract. Actual browser/security/recovery/platform evidence belongs to delivery, not planning proof.
- #356's four technical carry-forward follow-ups remain separate; no silent closure or acceptance. External review remains optional under ADR-0005.
- Version intent is `0.5.0`, preview first. Baseline has no stable/RC tag; latest preview is `v0.4.0-preview.2`. `.release-please-manifest.json` was absent on synchronized main; no release configuration was invented. No tag/publication authorized.

## Exercised publication evidence

Contract admission smoke checked ten complete issue/epic templates, unchecked observable acceptance/final acceptance, local file/heading links and the acyclic nine-child dependency order. Final link admission is repeated after live-ID integration.

Actual GitHub commands used:

```text
gh issue list --repo geoffrey-xiao/deprail --state all --search 'in:title "<exact-key>"' --json number,title,state,url
gh issue create --repo geoffrey-xiao/deprail --title '<exact-key and outcome>' --body '<complete contract>' --milestone v0.5.0 --assignee geoffrey-xiao --label '<each required label>'
gh issue view <created-number> --repo geoffrey-xiao/deprail --json number,title,state,labels,milestone,assignees
gh api repos/geoffrey-xiao/deprail/issues/<child-number> --jq .id
gh api repos/geoffrey-xiao/deprail/issues/418/sub_issues -X POST -F sub_issue_id=<returned-database-id>
gh project list --owner geoffrey-xiao --format json
gh project field-list 1 --owner geoffrey-xiao --format json
gh project item-add 1 --owner geoffrey-xiao --url <created-issue-url> --format json
gh api graphql -f query='<updateProjectV2ItemFieldValue mutations with live IDs>'
gh api repos/geoffrey-xiao/deprail/issues/418/sub_issues?per_page=100
gh project item-list 1 --owner geoffrey-xiao --limit 1000 --format json
```

No exact-key duplicate existed. Read-back verified ten OPEN issues, matching labels/milestone/assignee, exactly nine native children and ten Project Todo items with Owner, Owner Reviewer, Target Version and dependency fields. Required `Owner Reviewer` text field was absent and created, then resolved/read back from the live project; no environment-specific IDs are prescribed here. Priority/Risk/Area mirror the verified labels. Evidence Link points to the actual owner decision. Existing `Release · v0.5.0` view filters `milestone:v0.5.0`; no Sprint assignment was invented.

Full local Markdown is the durable contract. GitHub bodies use immutable commit links for local sources and contain live parent/dependency numbers; subsequent contract updates must synchronize both. Runtime tests/builds were not re-run for this documentation-only publication, and no runtime behavior was changed. Planning CI and owner review of this new handoff remain distinct from future delivery acceptance.

## Current correction handoff

At capture review submission, V05-QA-002 / #430 was a native child of #418, not a replacement/reopening of #419; duplicate search found no existing correction. The owner's “请继续我们的任务” authorized that bounded process repair, not a helper-fixture/H05 change. The owner subsequently merged PR #431, and #430 is now closed/Project Done. Its historical failure records remain unchanged; current helper correction and runtime-start decision are tracked separately below.

## Darwin fixture follow-up

The owner's “请继续” following the #431 merge authorizes bounded V05-QA-003 / #432, not H05 kickoff. Exact-key/capability searches found only the closed diagnostic/capture issues and unrelated existing work, so #432 is a distinct correction, native child of #418, milestone `v0.5.0`, Sprint 4, R2/P0 with sole required owner reviewer `@geoffrey-xiao`. The branch starts from synchronized `4f9fe3d`, containing the owner-merged capture fix. Production verification/process/scanner/CLI code is unchanged; original assertions and one-second/20ms deadlines remain. H05 issues stay Todo until correction review/merge, runtime-start decision and issue-specific DoR pass.

## H05-001 authorized runtime kickoff

The preceding publication/correction stages remain historical. Owner-merged [#433](https://github.com/geoffrey-xiao/deprail/pull/433), successful exact-head three-platform CI, and [explicit #432 acceptance](https://github.com/geoffrey-xiao/deprail/issues/432#issuecomment-5913275072) resolve the QA correction gate; #432 is closed/Project Done with native Linked pull requests containing #433. Technical acceptance and test-fixture security/architecture assessment are distinct entries; residual macOS cold-start uncertainty and original failed runs remain preserved.

The owner's later “可以 继续吧” authorizes only [H05-001 / #420](../issues/H05-001-safe-history-projection.md), after its plan and issue-specific DoR, on synchronized reviewed-main `6c74cf9861981928c29f7a8cb09799e63f3fd6dd`. [Live kickoff record](https://github.com/geoffrey-xiao/deprail/issues/420#issuecomment-5913450360) records branch `feat/420-safe-history-projection`, milestone `v0.5.0`, Sprint 4, R3/P0, owner reviewer and Project In Progress. No second R3 issue, storage/API/UI integration, new dependency or release begins.

The nine Todo entries (#418 and #420–#427 before kickoff) had repeated artificial Evidence Link values cleared at owner request; the existing Linked pull requests field is unchanged. PRs must use native issue associations, not custom card URLs as substitutes. A closing association still follows explicit owner acceptance; do not claim unreviewed criteria passed merely to populate the field.
