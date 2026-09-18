# Contributing to DepRail

DepRail is developed from the contracts in `docs/`. Before changing behavior, read `AGENTS.md`, the v0.1 PRD, the applicable requirements, and the linked issue.

## Development rules

- Work on one issue with one primary outcome.
- Keep changes within the issue scope and preserve the documented security boundaries.
- Treat repository files and scanner output as hostile input.
- Preserve deterministic serialized output and explicit `complete`, `partial`, and `failed` states.
- Never turn missing tools, malformed output, timeouts, or non-zero exits into empty success.
- Do not change schemas, public interfaces, tool versions, or security boundaries without explicit review.
- Keep machine data on stdout and diagnostics on stderr.

## Sprint and version kickoff

Before starting a new Sprint or version stage such as v0.1:

1. Confirm scope, dependencies, acceptance evidence, risk, and reviewer from the execution package.
2. Verify the matching GitHub milestone, Project fields, labels, and views against `tracking/GITHUB-PROJECT-SETUP.md`.
3. Check for existing GitHub Issues before creating new ones; never create duplicates.
4. Map each local epic to GitHub tracking and each file in `docs/.../issues/` to one GitHub Issue. Import `tracking/issue-backlog.csv` and use the local issue body when an issue is missing.
5. Link GitHub Issue/Project numbers or URLs back to the local tracking/evidence record.
6. Do not implement until the issue is assigned to the correct Sprint/milestone and meets Definition of Ready.

GitHub Issues are the workflow and review source of truth. The repository’s epics, issue files, requirements, and checklists remain the durable execution contract.


## Branch and pull request policy

`main` is the releasable integration branch. Do not develop directly on it. Create a short-lived branch from the current `main` for each issue or bounded task:

```text
<type>/<issue-key>-<short-slug>
```

Use lowercase, hyphenated, specific names with one of these types: `feat`, `fix`, `chore`, `docs`, `test`, or `refactor`. Examples:

```text
chore/s0-002-toolchain-ci
feat/schema-001-v1alpha
fix/disc-002-symlink-escape
docs/s0-003-template-workflow
```

Include the local issue key whenever an issue exists. Push the branch and open a pull request linked to the issue; do not push implementation commits directly to `main`.

Push the branch and open a pull request linked to the issue. Keep one primary outcome per pull request. The pull request must include scope, risk, contract impact, verification results, evidence, and rollback notes. Required CI and human review must pass before a squash merge to `main`; delete the branch after merge. Emergency security changes may use an expedited path but still require an issue, review, verification, and follow-up record.

## Commit and pull request traceability

- Every commit on an issue branch should reference exactly one issue when practical.
- Use:

  ```text
  <type>(<issue-key>): <imperative summary> (#<github-issue-number>)
  ```

  Examples: `feat(s0-001): establish repository baseline (#26)`, `fix(disc-002): reject symlink escape (#7)`, and `docs(s0-003): add pull request templates (#28)`.
- Squash commit titles must contain the GitHub issue reference, such as `(#27)`.
- Every pull request must link its issue with `Closes #N` or `Refs #N`.
- Every pull request must carry matching `area`, `risk`, `priority`, and `type` labels.
- Agents must check the issue link and required labels before requesting review.
- A missing issue reference or required label is a process defect and must be fixed before review or merge.

## Issue lifecycle and closure

Before starting a new issue, check the previous issue and pull request. If the previous work meets its acceptance criteria, verification, review, and evidence requirements, remind the project owner to close it. If it is incomplete, keep it open or mark it blocked with an owner and reason.

When an issue is finished, prepare a completion report for the owner covering every acceptance item, required command, CI result, human review result, evidence link, and remaining risk. The owner must inspect and check every acceptance item and explicitly confirm completion before closure. An agent may explain evidence and draft the closure comment, but may close the issue only after the owner authorizes it. A merged pull request alone does not prove completion, and local checklists require linked evidence and owner confirmation.

## Defect triage

Create a separate GitHub bug issue when a defect is outside the current issue, changes a contract/schema/security boundary, affects another component, needs follow-up, blocks work, or cannot be fixed without expanding the current pull request. Keep it in the current issue only when it is a direct acceptance failure and the fix preserves one primary outcome.

Record reproduction, expected and observed behavior, impact, risk, priority, dependencies, regression-test requirements, reviewer, and evidence. Apply matching `type:bug`, `area:*`, `risk:*`, and `priority:*` labels, and link the bug from related issues and pull requests. Critical security, false-safe, data-loss, path-escape, command-injection, and credential-leak defects must be tracked immediately.

Do not hide important defects in commits, TODOs, or review comments. Keep the bug open until regression evidence, verification, human review, and owner acceptance are complete.

## Change workflow

1. Link the issue and identify scope, exclusions, risk, contracts, and acceptance evidence.
2. Present a short implementation plan before editing.
3. Add or update focused tests, fixtures, and documentation with the behavior change.
4. Run the narrow checks first, then `make verify` when the Sprint 0 toolchain is available.
5. Attach command results, sample output, remaining risk, and required human review to the pull request.

Use the repository templates under `docs/04-execution/deprail-v0.1-execution-package/templates/` for issues, pull requests, and ADRs. Schema, external-process, file-write, permission, credential, migration, and publishing changes require human review.

## Code of conduct and licensing

Contributions are made under the Apache-2.0 license in `LICENSE`. Be specific, respectful, and security-conscious when reporting defects or proposing changes.
