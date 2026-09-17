# DepRail v0.1 Development Execution Package

This folder turns the v0.1 product and architecture decisions into ready-to-run requirements, epics, issues, tracking checklists, AI prompts, and review templates.

## Use Order

1. Read `PRD-v0.1.md`.
2. Review every file in `requirements/`.
3. Create milestones and fields from `tracking/GITHUB-PROJECT-SETUP.md`.
4. Start with `tracking/SPRINT-0.md` and the three `S0-*` issues.
5. Use `prompts/IMPLEMENT.md`, `TEST.md`, and `REVIEW.md` for bounded AI work.
6. Update `tracking/MASTER-CHECKLIST.md` only when linked evidence exists.

## Structure

- `requirements/` defines observable contracts.
- `epics/` groups outcomes and dependency boundaries.
- `issues/` contains implementation-ready work items.
- `tracking/` contains Sprint and release control.
- `prompts/` constrains AI implementation, testing, and review.
- `templates/` standardizes ADRs, Issues, and Pull Requests.

## Recommended Management Tools

Use GitHub Issues for atomic work, Milestones for Sprint 0 and S1-S4, GitHub Projects for flow and risk fields, and pull requests for review evidence. Keep the Markdown checklists in the repository as the durable execution contract.

## Execution Rules

- One issue has one primary outcome.
- At most two implementation issues and two review PRs are active.
- R3 work is never parallelized with another R3 item.
- Every issue delivers code, tests, documentation, and evidence together.
- No failure becomes an empty success.
- Contract or schema changes require explicit human approval.
- Generated and golden-file changes are explained and reviewed.

## Definition of Complete

An item is complete only when acceptance checkboxes are backed by commands, test results, screenshots or sample output where appropriate, and a human review decision. A merged implementation without its required evidence remains incomplete.
