---
name: issue-create
description: Reconcile DepRail planning contracts and create ordered local epics, implementation-ready sub-issues, and synchronized GitHub tracking items at release or stage kickoff.
---

# DepRail Issue Creation and Release Planning

Use this skill when starting a release, sprint, epic, or coordinated implementation sequence that must be represented in both repository documents and GitHub Issues/Projects.

## Non-negotiable rule

Do not create implementation issues from a single request, PRD, or roadmap item. Reconcile the complete planning hierarchy first. If the documents disagree, stop issue creation and record the conflict as an ADR or versioned release-plan decision before continuing.

## Required planning read order

Read these sources in order and retain the release/stage decisions they establish:

1. `docs/01-product/deprail-product-design-v1-ai.md`
2. `docs/02-architecture/deprail-architecture-and-tech-stack-v1.md`
3. `docs/03-planning/deprail-roadmap-v1.md`
4. The applicable release plan under `docs/03-planning/`, such as `deprail-development-plan-v0.4.md`
5. The applicable execution package under `docs/04-execution/`
6. The linked product requirements, epic contracts, issue contracts, sprint plan, and checklist
7. `docs/RELEASE-CHECKLIST.md` when the work is release-related
8. `docs/04-execution/deprail-v0.1-execution-package/tracking/GITHUB-PROJECT-SETUP.md` and the applicable workflow guide

Lower-level documents cannot silently override product, architecture, roadmap, or release-plan decisions.

## Inputs to collect

Before creating anything, identify:

- Release or stage name and intended version.
- User outcome and included capabilities.
- Explicit exclusions and displaced work.
- Required dependencies and sequencing.
- Compatibility and migration impact.
- Failure, incomplete, cancellation, and rollback behavior.
- Acceptance evidence and required manual smoke scenarios.
- Risk level and security/architecture review requirement.
- Named owner and independent reviewer.
- Target milestone and Project fields.
- Existing local epics/issues and existing GitHub issues.

## Scope reconciliation

Produce a private planning table before mutation:

| Capability | Product contract | Architecture boundary | Roadmap/release item | Execution/requirements contract | Decision |
|---|---|---|---|---|---|
| ... | link/path | link/path | link/path | link/path | include/exclude/defer |

Every included capability must map to at least one authoritative contract and exactly one implementation-ready issue. Do not add inferred capabilities, retries, telemetry, abstractions, or roadmap features without an explicit inclusion decision.

If scope changes after reconciliation:

1. Stop issue creation.
2. Create a new versioned release plan or ADR.
3. Record changed scope, reason, displaced work, compatibility impact, and remaining risk.
4. Reconcile local documents, tracking, and GitHub state.
5. Resume only when Definition of Ready is satisfied.

## Ordered epic and sub-issue design

Create epics as GitHub Issues. An epic is not a project board placeholder.

Each epic must contain:

- User outcome.
- Included capabilities.
- Explicit exclusions.
- Dependencies.
- Ordered implementation stages.
- Child issue list.
- Release/compatibility impact.
- Failure behavior.
- Acceptance gate.
- Required evidence.
- Owner and independent reviewer.

Each child issue must have one primary outcome and be implementation-ready. Keep child issues small enough for one branch and one PR. A child issue must not own two independent primary outcomes merely because they share a subsystem.

Use dependency-aware ordering:

1. Contracts and decision records.
2. Domain/schema changes.
3. Security boundaries and path/process primitives.
4. Adapters or external integrations.
5. Application orchestration.
6. Presenters and CLI/API wiring.
7. Contract/integration tests.
8. Cross-platform and hostile-input coverage.
9. Release evidence and documentation.

Represent dependencies explicitly as `#N` links in local and GitHub issue bodies. Do not rely on issue number order to imply execution order.

Recommended child table:

| Order | Issue key | Primary outcome | Dependencies | Evidence gate |
|---:|---|---|---|---|
| 1 | `Vxx-001` | Contract/decision | None | Reviewed contract |
| 2 | `Vxx-002` | Security boundary | #... | Boundary tests |
| 3 | `Vxx-003` | Core implementation | #... | Unit/contract tests |
| 4 | `Vxx-004` | Integration wiring | #... | CLI/integration smoke |
| 5 | `Vxx-005` | Cross-platform acceptance | #... | Linux/macOS/Windows evidence |

## Local documentation creation

For a new release or epic, create or update only the applicable durable contracts:

- Versioned release plan under `docs/03-planning/`.
- Epic contract under the active execution package `epics/`.
- Issue contracts under the active execution package `issues/`.
- Requirements/ADR when the change affects schema, CLI, permissions, compatibility, external processes, writes, credentials, publishing, or other architectural boundaries.
- Release evidence/checklist references when release gates are affected.

Use existing templates and naming conventions. Preserve historical plans. Never rewrite an accepted historical plan to hide scope changes. Do not mark checklists complete without linked evidence.

## Issue body requirements

Use `docs/04-execution/deprail-v0.1-execution-package/templates/ISSUE-TEMPLATE.md` or the current release template. Every issue body must include:

```markdown
# <issue key> <title>

## Planning metadata
- Type:
- Area:
- Priority:
- Risk:
- Target version:
- Milestone:
- Sprint:
- Owner:
- Reviewer:
- Dependencies:
- Blocked reason:

## Definition of Ready
- [ ] Value and user impact are stated.
- [ ] Scope and explicit exclusions are stated.
- [ ] Inputs, outputs, and failure behavior are defined.
- [ ] Required tests or smoke scenarios are named.
- [ ] Acceptance criteria are observable.
- [ ] Owner and reviewer are assigned.
- [ ] Dependencies and target version are recorded.

## Goal

## Scope

## Out of Scope

## Inputs, Outputs, and Failure Behavior

## Required Tests

## Acceptance Criteria

## Human Review

## Evidence Required

## Final Acceptance
```

Acceptance criteria must describe consumer-observable behavior, boundaries, transitions, and real failure modes. Avoid asserting implementation details, source text, incidental defaults, or mock behavior.

## Duplicate and state checks

Before creating any GitHub issue:

1. Search existing issues by exact title, issue key, capability, and linked epic.
2. Search local `docs/**/issues/`, `tracking/`, and the release plan.
3. Reuse or update the existing contract when it covers the same outcome.
4. Create a new issue only for genuinely missing scope or an explicitly versioned follow-up.
5. Record duplicate/disposition decisions in the epic or release tracking record.

Never create duplicate epics or create a new issue just because an existing issue is closed; create a follow-up only when the new acceptance criteria are materially different.

## GitHub creation workflow

Use the repository and current release milestone. Resolve live IDs; never hard-code project IDs or status option IDs from memory.

1. Verify the GitHub milestone exists; create/update it from the tracking setup guide if authorized.
2. Verify the shared `DepRail` Project exists.
3. Create the epic issue first.
4. Create child issues in dependency order, each linking its parent epic and dependencies.
5. Apply required labels: matching `area:*`, `risk:*`, `priority:*`, and `type:*`.
6. Add every issue to the shared project.
7. Set new issues to Project `Todo`.
8. Verify issue state, milestone, labels, project membership, and status with `gh`.
9. Record GitHub issue URLs/numbers in local epic/issue contracts or tracking records.

Typical commands:

```bash
gh issue list --repo geoffrey-xiao/deprail --state all --search '<title or key>' --json number,title,state,url
gh issue create --repo geoffrey-xiao/deprail --title '<title>' --body-file /tmp/issue.md --label area:cli --label risk:R3 --label priority:P0 --label type:feature --milestone v0.4.0
gh project list --owner geoffrey-xiao --format json
gh project item-add <project-number> --owner geoffrey-xiao --url <issue-url> --format json
gh project item-list <project-number> --owner geoffrey-xiao --format json
```

Resolve the live Project Status field and option IDs before calling `gh project item-edit`. The initial state is normally `Todo`; implementation start changes it to `In Progress`; an opened PR changes it to `Review`; reviewed merged work with passing CI and evidence changes it to `Done`.

If project-write permission is unavailable, report the exact failure and do not claim synchronization occurred.

## Release-start output

At the end of planning, provide a tracking map:

| Order | Local contract | GitHub issue | Project status | Milestone | Dependencies | Owner/reviewer |
|---:|---|---|---|---|---|---|
| 0 | Epic | #... | Todo | v... | None | ... |
| 1 | Issue | #... | Todo | v... | #... | ... |

Also record:

- Included capabilities.
- Explicit exclusions.
- Displaced/deferred work.
- Definition-of-Ready gaps.
- Required human review.
- Required release evidence.
- Remaining risks and rollback plan.

## Implementation and closure handoff

Do not begin implementation merely because issues exist. Each issue must satisfy Definition of Ready and be assigned to the active release/milestone.

When implementation begins:

- Synchronize `origin/main` before creating the issue branch.
- Set Project status to `In Progress`.
- Work one issue at a time.
- Do not parallelize two high-risk release items without an explicit contract.

When a PR opens:

- Link exactly one primary issue with `Closes #N` only when all acceptance criteria are complete; otherwise use `Refs #N`.
- Verify labels, issue linkage, project membership, and set Project status to `Review`.
- Record commands actually run and evidence paths.

After reviewed merge:

- Attach acceptance and CI evidence.
- Record owner and independent reviewer decisions separately.
- Set Project status to `Done` only after required evidence passes.
- Close the issue only when all acceptance criteria and release evidence are complete.
