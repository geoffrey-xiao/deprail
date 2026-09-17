# GitHub Project Setup

## Recommendation GitHub Issues and Projects

Use Issues as the source of truth, Projects for flow, Milestones for S0-S4, and pull requests for implementation evidence.

## Milestones

Create Sprint 0, Sprint 1, Sprint 2, Sprint 3, Sprint 4, and v0.1 Preview.

## Labels

`area:foundation`, `area:discovery`, `area:adapter`, `area:normalization`, `area:cli`, `area:test`, `area:docs`, `risk:R0` through `risk:R3`, `priority:P0` through `priority:P2`, `status:blocked`, and `type:bug|feature|test|docs|decision`.

## Project Fields

Status, Sprint, Milestone, Priority, Risk, Area, Owner, Reviewer, Target Version, Dependencies, Blocked Reason, and Evidence Link.

## Views

Board by status; Sprint table; high-risk review; blocked items; missing evidence; release gate.

## Automation

New issue enters Backlog. Assigned Sprint enters Ready when Definition of Ready is checked. An opened PR moves the item to Review. Merge does not move to Done until evidence and acceptance are complete.

## Import Method

Import `issue-backlog.csv`, verify IDs and dependencies, add issue bodies from `issues/`, and then assign milestones and reviewers.
