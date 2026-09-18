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

Board by status; Sprint table; high-risk review; blocked items; missing evidence; release gate.

## Automation

New issues enter `Todo`. Definition of Ready is recorded in the issue body rather than represented by a separate status. Creating a branch and starting implementation moves the item to `In Progress`. An opened PR moves the item to `Review`. An external dependency or unresolved decision moves it to `Blocked`; `Blocked Reason`, an owner, and a next-check date are required. Review rework returns to `In Progress`. Merge does not move to `Done` until verification, evidence, and owner acceptance are complete.

## Import Method

Import `issue-backlog.csv`, verify IDs and dependencies, add issue bodies from `issues/`, and then assign milestones and reviewers.
