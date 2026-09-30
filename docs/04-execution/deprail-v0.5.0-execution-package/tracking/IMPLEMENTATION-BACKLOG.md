# v0.5 Approved Implementation Backlog

## Authorization and boundary — 2026-09-30

Owner `@geoffrey-xiao` accepted PR #417 head `b0649547d1e0f033f10f3f28920b74eab9b16e1b` with “approve，请继续”. [Conversation-derived decision record](https://github.com/geoffrey-xiao/deprail/pull/417#issuecomment-5907670878) separately records technical acceptance, security/architecture acceptance and verification disposition; it is not a fabricated native GitHub APPROVED review. Owner already merged #417 at `40e43a83eaac1b84dd038ee1ea4dccbd3e7234f6` before this handoff. New branch `docs/392-approved-implementation-backlog` was created from synchronized main; no commits were appended to the merged branch.

Authorization covers complete planning contracts and GitHub publication only. No runtime code, new-PR merge, dependency introduction, Sprint kickoff or release publication is authorized. Each delivery issue remains **Todo**; required runtime acceptance and owner review remain unchecked. The current handoff requires its own owner PR review.

## Published contracts

All rows are assigned to `@geoffrey-xiao`, owner reviewer `@geoffrey-xiao`, milestone `v0.5.0`, Target Version `0.5.0`, Project Todo, Sprint unassigned. Exact `area`, `risk`, `priority` and `type` labels were read back after creation. Every child is a native GitHub subissue of #418, not a child of readiness epic #390. #418 depends on the accepted readiness baseline.

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

This is a dependency DAG, not permission to start all issues. Work one issue at a time, at most two implementation issues/two review PRs active, never two R3 implementations concurrently. H05-006 can exercise the real H05-005 API through a throwaway embedding harness before H05-007 produces production assets; no shipped placeholder assets or cyclic release dependency. Store API close/reopen evidence belongs to H05-002; CLI capture is H05-004. HTTP owns wire envelopes/cursors; application services return typed domain results and metadata-only integrity outcomes. Unknown provenance remains null rather than causing blanket rejection.

## Pre-start and residual-risk gates

- #419 preserves the observed local `make verify` helper-process timeout as unresolved. Three-OS CI from the accepted proposal does not erase it. Resolve with evidence or obtain explicit owner pre-start disposition before runtime implementation; do not rerun merely to confirm, skip the assertion, blindly lengthen the deadline or assume a flake.
- #419 [controlled investigation evidence](V05-QA-001-INVESTIGATION.md) is available for owner review: cold-image delay is observed, exact historical/OS cause remains uncertain, no correction applied. A separate real smoke captured 155 stdout bytes for a 32-byte output cap; that unresolved security boundary must be included in the owner decision, not mistaken for passing output-limit evidence. Runtime pre-start disposition remains pending.
- Issue-specific Definition of Ready and owner review, dependency availability and synchronization from reviewed main are required before a runtime branch. This handoff does not grant permission to implement #419's fix either; that issue must begin with its own bounded plan.
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
