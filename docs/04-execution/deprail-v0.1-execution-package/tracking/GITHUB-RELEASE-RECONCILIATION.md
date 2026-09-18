# GitHub Issue and Local Contract Reconciliation

Use this procedure at the end of each release, before the release issue is accepted or closed. It reconciles the durable issue contracts in `docs/` with the active GitHub workflow record without turning local Markdown into a second status database.

## Sources of truth

- Local issue documents define the durable goal, scope, dependencies, contracts, acceptance criteria, risk, review requirements, and evidence expectations.
- GitHub Issues define active workflow state, assignees, reviewers, pull requests, CI, comments, acceptance discussion, and closure.
- The GitHub Project defines status, release views, grouping, and evidence links.

When sources disagree, resolve the difference explicitly. Do not silently overwrite either source.

## When to run

Run this procedure during the release gate, after implementation and review work is complete but before the release issue is accepted. For DepRail v0.1, attach the result to `REL-001` evidence.

## Reconciliation steps

### 1. Enumerate the release scope

Identify the release milestone and Target Version. Export or list every issue in that scope, including completed, open, blocked, and re-planned work.

Confirm that each release issue has:

- one local issue key;
- one GitHub Issue;
- the expected milestone and Target Version;
- the expected area, risk, priority, and type labels;
- a Project item.

### 2. Compare contract content

For each issue, compare the local document and GitHub Issue body:

- title and issue key;
- goal;
- scope and out of scope;
- inputs, outputs, and failure behavior;
- dependencies;
- acceptance criteria;
- required tests;
- risk and human review requirements;
- evidence expectations.

The GitHub Issue should contain a local-source comment such as:

```markdown
<!-- Local source: docs/04-execution/deprail-v0.1-execution-package/issues/<issue-file>.md -->
```

Update both sources when the contract changed. Preserve the local document as the durable contract and record material decisions in the GitHub Issue discussion.

### 3. Verify implementation evidence

For completed issues, verify:

- linked pull request;
- required CI checks;
- commands or scenarios actually run;
- human review result;
- remaining risk or limitation;
- Evidence Link in the Project.

A merged pull request is not sufficient evidence by itself.

### 4. Verify workflow state

Check that Project Status matches reality:

- `Todo`: not started;
- `In Progress`: implementation active;
- `Review`: pull request or acceptance review active;
- `Blocked`: external dependency or decision prevents progress;
- `Done`: verification, evidence, owner acceptance, and closure are complete.

Blocked items must have an owner, Blocked Reason, and next-check date. The `status:blocked` label is a search aid; Project Status remains authoritative.

### 5. Verify release views and fields

Confirm that the release view uses the intended milestone or Target Version filter and contains every release issue. Check that:

- completed issues remain visible for release accounting;
- issues without release assignment are not silently included;
- Sprint and Milestone values are consistent;
- risk, priority, area, owner, reviewer, and evidence fields are populated where required.

### 6. Record reconciliation evidence

Add a release evidence record containing:

- release name and date;
- issue range or list reconciled;
- local documents compared;
- GitHub Project and release-view URLs;
- discrepancies found;
- decisions made;
- commands or queries used;
- reviewer and owner decision;
- remaining risks.

Do not mark the release checklist complete until this evidence is linked.

## Safe verification examples

These commands inspect state without mutating the repository or project:

```text
gh issue list --state all --milestone "v0.1 Preview"
gh pr list --state all --search "milestone:v0.1 Preview"
gh project view 1 --owner geoffrey-xiao
gh project item-list 1 --owner geoffrey-xiao
```

Use a token with the minimum required GitHub scope. Never print token values, environment dumps, or credential-bearing URLs.

## Failure handling

If an issue has no local contract, no GitHub Issue, missing release metadata, missing evidence, or conflicting acceptance state:

1. Leave the release item open.
2. Record the discrepancy and owner in the release evidence.
3. Create or update the missing contract or GitHub record.
4. Re-run reconciliation.
5. Do not mark the issue or release complete until the discrepancy is resolved or explicitly re-planned.
