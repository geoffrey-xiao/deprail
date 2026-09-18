# GitHub Project Setup

## Recommendation GitHub Issues and Projects

Use Issues as the source of truth, Projects for flow, Milestones for S0-S4, and pull requests for implementation evidence.

## Milestones

Create Sprint 0, Sprint 1, Sprint 2, Sprint 3, Sprint 4, and v0.1 Preview.

## Labels

`area:foundation`, `area:discovery`, `area:adapter`, `area:normalization`, `area:cli`, `area:test`, `area:docs`, `risk:R0` through `risk:R3`, `priority:P0` through `priority:P2`, `status:blocked`, and `type:bug|feature|test|docs|decision`.
`status:blocked` is a synchronized search label for items whose Project Status is `Blocked`; the Project Status field is authoritative. Do not add lifecycle labels for `In Progress`, `Review`, or `Done`.

## Project Fields

Status (`Todo`, `In Progress`, `Review`, `Blocked`, `Done`), Sprint, Milestone, Priority, Risk, Area, Owner, Reviewer, Target Version, Dependencies, Blocked Reason, and Evidence Link.

## Views

The single long-lived `DepRail` project uses release-specific views over the shared issue history. The current view is `Release · v0.1 Preview`, a table filtered to milestone `v0.1 Preview`. When a future release has committed work, create a matching view such as `Release · v0.2`; do not create a separate project. Keep the existing view layout and status flow unless tracking needs justify a later change.

## Automation

New issues enter `Todo`. Definition of Ready is recorded in the issue body rather than represented by a separate status. Creating a branch and starting implementation moves the item to `In Progress`. An opened PR moves the item to `Review`; because the current repository has no automatic PR-to-Project mutation, the agent or maintainer MUST set this explicitly with `gh project item-edit` and verify it with `gh project item-list`. An external dependency or unresolved decision moves it to `Blocked`; `Blocked Reason`, an owner, and a next-check date are required. Review rework returns to `In Progress`. Merge does not move to `Done` until verification, evidence, and owner acceptance are complete.

### Required status synchronization

After opening a PR, resolve the live `DepRail` project and field IDs, locate the linked issue's project item, and set the Status field to the `Review` option. Do not copy project IDs from another repository or assume PR creation changed the item. If the token lacks project-write permission, report the blocker and leave the status unchanged rather than claiming success.

## Issue Implementation Sequence

For each issue, follow this sequence without reusing a previously completed issue branch:

1. Confirm the GitHub issue, local contract, dependencies, reviewer, and Definition of Ready.
2. Start from the current `main` and create one dedicated branch for that issue.
3. Implement only that issue's primary outcome, including its tests and evidence.
4. Run the required verification commands and record the results.
5. Commit with the issue key and number, push the branch, and open one labeled pull request linked to the issue.
6. Leave the issue and local checklist open until CI, human review, linked evidence, and owner acceptance are complete.
7. Merge only through the reviewed pull request; then update the issue and checklist with the evidence.

## Import Method

Import `issue-backlog.csv`, verify IDs and dependencies, add issue bodies from `issues/`, and then assign milestones and reviewers.
