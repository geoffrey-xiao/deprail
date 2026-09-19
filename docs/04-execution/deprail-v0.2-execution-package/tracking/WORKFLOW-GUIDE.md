# DepRail Issue and Pull Request Workflow Guide

This guide prevents metadata validation failures and keeps issue, PR, and Project records connected.

## 1. Prepare the issue body

Use the applicable template under `.github/ISSUE_TEMPLATE/`. Before creating the issue, complete:

- Planning metadata: type, area, priority, risk, target version, milestone, sprint, owner, reviewer, dependencies.
- Definition of Ready.
- Goal, scope, exclusions, inputs, outputs, and failure behavior.
- Required tests or smoke scenarios.
- Observable acceptance criteria.
- Evidence requirements.

## 2. Create the issue

```bash
gh issue create \
  --repo geoffrey-xiao/deprail \
  --title "[V02-001] Short imperative title" \
  --label "area:normalization" \
  --label "risk:R1" \
  --label "priority:P0" \
  --label "type:bug" \
  --milestone "v0.2.0" \
  --body-file issue.md
```

Exactly one label from each category is required:

```text
area:foundation|discovery|adapter|normalization|cli|test|docs
risk:R0|R1|R2|R3
priority:P0|P1|P2
type:bug|feature|test|docs|decision
```

Always verify the result; do not assume the creation command attached every label:

```bash
gh issue view <issue-number> \
  --repo geoffrey-xiao/deprail \
  --json number,labels,milestone,state
```

Repair missing labels before implementation:

```bash
gh issue edit <issue-number> \
  --repo geoffrey-xiao/deprail \
  --add-label "area:normalization" \
  --add-label "risk:R1" \
  --add-label "priority:P0" \
  --add-label "type:bug"
```

## 3. Add and configure the Project item

```bash
gh project item-add 1 \
  --owner geoffrey-xiao \
  --url https://github.com/geoffrey-xiao/deprail/issues/<issue-number>
```

Resolve the live item ID and set the Project fields. Do not reuse IDs from another repository or stale output:

```bash
gh project item-list 1 \
  --owner geoffrey-xiao \
  --format json
```

At minimum set:

```text
Status: Todo
Target Version: 0.2.0
Milestone: v0.2.0
Priority, Risk, Area, Sprint, Owner, Reviewer, Dependencies
```

Verify the item after editing:

```bash
gh project item-list 1 \
  --owner geoffrey-xiao \
  --format json
```

## 4. Start implementation and create the branch

Before editing, set the issue's Project Status to `In Progress`; the GitHub issue itself remains open:

```bash
gh project item-edit \
  --project-id <project-id> \
  --id <item-id> \
  --field-id <status-field-id> \
  --single-select-option-id <in-progress-option-id>
gh project item-list <project-number> \
  --owner geoffrey-xiao \
  --format json
```

Verify the item reports `In Progress` before creating the branch. After the PR opens, the required next transition is `Review`; the issue is not closed until the reviewed PR merges and evidence is complete.


## 5. Run the OMP pre-PR review

Before invoking the reviewer, run the applicable focused verification and record its command and output. The reviewer must read the linked issue, applicable requirements, changed files, tests, and current verification evidence.

The reviewer checks:

- scope, dependencies, acceptance criteria, and contract compatibility;
- failure, incomplete, boundary, security, and deterministic behavior;
- path containment, subprocess, file-write, credential, and permission risks;
- focused tests, cross-platform evidence, documentation, issue link, labels, and rollback.

Present the findings to the owner without editing the branch. Use these classifications:

```text
BLOCKER  Blocks PR creation unless the owner explicitly changes the plan.
MAJOR    Requires an owner decision before PR creation.
MINOR    Owner decides whether to fix now or record for later.
NOTE     Context only; no action required.
```

Stop after reporting the findings. Do not automatically fix, suppress, accept, or reject any finding. The owner chooses the next action. If the owner requests fixes, apply only those requested changes, rerun focused verification, and rerun the OMP review. If the owner accepts remaining risk, record that decision and its rationale before proceeding. Human review remains required after the PR opens.

## 6. Prepare and create the PR

Use the PR template. The PR body must contain exactly one unique linked issue number. For an implementation PR, use a closing keyword:

```markdown
Closes #<issue-number>
```

Use `Fixes` or `Resolves` when appropriate. Use `Refs` only when the PR does not complete the issue.

Create the PR with labels matching the issue:

```bash
gh pr create \
  --repo geoffrey-xiao/deprail \
  --base main \
  --head <issue-branch> \
  --title "fix(v02-001): short summary (#<issue-number>)" \
  --label "area:normalization" \
  --label "risk:R1" \
  --label "priority:P0" \
  --label "type:bug" \
  --body-file pr.md
```

Verify linkage and labels:

```bash
gh pr view <pr-number> \
  --repo geoffrey-xiao/deprail \
  --json number,labels,closingIssuesReferences,state
```

If labels are missing, repair them explicitly:

```bash
gh pr edit <pr-number> \
  --repo geoffrey-xiao/deprail \
  --add-label "area:normalization" \
  --add-label "risk:R1" \
  --add-label "priority:P0" \
  --add-label "type:bug"
```

## 7. Synchronize Project review status

After the PR opens, locate the linked issue’s Project item and set Status to `Review` using the live Project and field IDs. Verify the result with `gh project item-list`.

Do not claim synchronization if the token lacks Project write permission.

## 8. Run checks before requesting review

```bash
gh pr checks <pr-number> --repo geoffrey-xiao/deprail
```

The metadata check requires:

- valid labels on the PR;
- matching labels on the linked issue;
- required issue and PR sections;
- one unique linked issue;
- required Project synchronization checklist in the PR body.

## 9. Merge and close

During PR review, the owner reviews acceptance criteria, evidence, compatibility, security impact, and remaining risk. After the reviewed PR merges and required CI/evidence pass:

1. Attach final evidence to the issue.
2. Set the Project item to `Done`.
3. `Closes #N` automatically closes the issue when the PR merges.

A `Refs #N` relationship does not provide the same automatic closure behavior.
