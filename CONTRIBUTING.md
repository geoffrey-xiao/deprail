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


## Branch and pull request policy

`main` is the releasable integration branch. Do not develop directly on it. Create a short-lived branch from the current `main` for each issue or bounded task:

```text
chore/s0-002-toolchain-ci
feat/discovery-walker
fix/path-containment
```

Push the branch and open a pull request linked to the issue. Keep one primary outcome per pull request. The pull request must include scope, risk, contract impact, verification results, evidence, and rollback notes. Required CI and human review must pass before a squash merge to `main`; delete the branch after merge. Emergency security changes may use an expedited path but still require an issue, review, verification, and follow-up record.

## Change workflow

1. Link the issue and identify scope, exclusions, risk, contracts, and acceptance evidence.
2. Present a short implementation plan before editing.
3. Add or update focused tests, fixtures, and documentation with the behavior change.
4. Run the narrow checks first, then `make verify` when the Sprint 0 toolchain is available.
5. Attach command results, sample output, remaining risk, and required human review to the pull request.

Use the repository templates under `docs/04-execution/deprail-v0.1-execution-package/templates/` for issues, pull requests, and ADRs. Schema, external-process, file-write, permission, credential, migration, and publishing changes require human review.

## Code of conduct and licensing

Contributions are made under the Apache-2.0 license in `LICENSE`. Be specific, respectful, and security-conscious when reporting defects or proposing changes.
