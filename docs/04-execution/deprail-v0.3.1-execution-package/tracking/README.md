# v0.3.1 Tracking Baseline

This directory records release readiness after the v0.3.1 epic and issue contracts were created.

## Current status

- Release: v0.3.1
- Stage: Planning and contract reconciliation
- Implementation authorized: No
- Context issue: #238
- Development plan: `docs/03-planning/deprail-development-plan-v0.3.1.md`
- Execution package: `../README.md`

## Required gates

1. Owner reviews the development plan and this execution package.
2. Architecture/security review confirms event, renderer, dependency, and hostile-input boundaries.
3. Acceptance and evidence scenarios are confirmed.
4. Definition of Ready is recorded.
5. Runtime implementation begins only after the Definition of Ready is approved.

The local epics and issue contracts are mapped to GitHub milestone `v0.3.1` and the DepRail release project view. No implementation work is authorized by backlog creation alone.

## GitHub hierarchy and backlog workflow

Native sub-issue relationships are established:

- EPIC-001 `#244`: `#255`, `#246`
- EPIC-002 `#242`: `#245`, `#256`
- EPIC-003 `#243`: `#252`, `#257`, `#250`, `#247`
- EPIC-004 `#241`: `#254`, `#253`

Project #1 now has a `Backlogs` option in the existing `Status` field (`d9e5dfde`, gray, “Queued backlog work”). Release view #4 (`Release · v0.3.1`) exposes that field and remains filtered to `milestone:v0.3.1`. Use `Backlogs` for approved future work that is intentionally not yet scheduled for implementation.
