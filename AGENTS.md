# DepRail Agent Guide

## Purpose and source of truth

DepRail is a dependency security guardrail for multi-language repositories. The current checkout is documentation-first; implementation is expected to begin with Sprint 0. Treat the versioned contracts in `docs/` as normative until an approved code or ADR decision supersedes them.

Before changing behavior or creating release work, read the planning hierarchy in this order:

1. `docs/01-product/deprail-product-design-v1-ai.md`
2. `docs/02-architecture/deprail-architecture-and-tech-stack-v1.md`
3. `docs/03-planning/deprail-roadmap-v1.md`
4. The applicable release plan, such as `docs/03-planning/deprail-development-plan-v0.2.md`
5. The applicable execution package under `docs/04-execution/`
6. The linked epic, issue contract, requirements, sprint, and checklist

Product and architecture define intent and boundaries. The roadmap defines whole-project sequencing. The release plan selects the current release scope. The execution package defines implementation contracts. Lower-level documents MUST NOT silently contradict higher-level documents.

The execution package is the active delivery contract only after it has been reconciled with the product, architecture, roadmap, and release plan. Do not mark checklist items complete without linked evidence.

## Mandatory planning and scope gate

No new release, epic, issue, implementation branch, or PR may begin until the scope is reconciled across:

- Product design.
- Architecture and technology baseline.
- Whole-project roadmap.
- Release-specific development plan.
- Execution package requirements and checklists.

The release plan MUST state the user outcome, included capabilities, exclusions, dependencies, compatibility impact, failure behavior, acceptance gate, risks, and displaced work. Every epic and issue MUST map to a release-plan item and at least one product, architecture, or requirements contract. Candidate capabilities MUST NOT enter implementation without an explicit inclusion decision.

If plans or context change:

1. Stop issue creation and implementation for the affected scope.
2. Create a new versioned release plan or an explicit ADR/decision record.
3. Reconcile affected product, architecture, requirements, execution, tracking, and GitHub records.
4. Record what changed, why, displaced work, compatibility impact, and remaining risk.
5. Resume only after the updated Definition of Ready is satisfied.

Historical plans MUST remain preserved. Do not silently rewrite a prior plan to hide a scope change.

## Sprint and version kickoff

Before starting a new Sprint or version stage such as v0.1, first establish the GitHub tracking set:

1. Confirm the stage, Sprint, scope, dependencies, acceptance evidence, risk, and named reviewer from the reconciled release plan and execution package.
2. Verify that the corresponding GitHub milestone, Project fields, labels, and views exist; create or update them from `tracking/GITHUB-PROJECT-SETUP.md`.
3. Check GitHub for existing issues before creating anything. Do not create duplicates.
4. Map each local `EPIC-*` to its GitHub tracking item and each implementation-ready file under `docs/.../issues/` to one GitHub Issue. Import `tracking/issue-backlog.csv` and copy the authoritative issue body when an issue is absent.
5. Record the GitHub Issue or Project URL/number in the local evidence or tracking record, and preserve the local Markdown as the durable contract.
6. Do not start implementation until the issue is assigned to the correct milestone/Sprint and satisfies Definition of Ready: value, scope, dependencies, contracts, failure behavior, acceptance tests, risk, reviewer, and evidence are explicit.

At stage kickoff, reconcile the GitHub state with the local roadmap, release plan, epics, issue backlog, Sprint checklist, and Master Checklist. GitHub is the workflow and review source of truth; repository documents remain the durable contract.

## Project identity and scope

- Product name: `DepRail`.
- CLI and command names: lowercase `deprail`; entry point `cmd/deprail`.
- Recommended MCP binary: `deprail-mcp`.
- Configuration: `.deprail.yaml`.
- Local data: `.deprail/`.
- Repository: `git@github.com:geoffrey-xiao/deprail.git`.
- Public positioning: dependency security guardrail for multi-language repositories.

The v0.1 goal is a stable, explainable, machine-readable scan for JavaScript, Python, and Java repositories. Required commands are:

```text
deprail doctor
deprail discover [path]
deprail scan [path]
```

V0.1 includes discovery, OSV-Scanner integration, raw artifact retention, deterministic normalization, completeness states, terminal/JSON output, and Linux/macOS/Windows smoke coverage. It explicitly excludes baseline diff, policy gates, SARIF, remediation or file mutation, web/history, Trivy, container/SBOM/license/secret/IaC scanning, remote publishing, and automatic scanner installation.

## Architecture rules

Use a Go single-binary core with inward dependencies:

```text
cmd/deprail/                 CLI entry point
cmd/deprail-mcp/             MCP entry point (later)
internal/app/                application services
internal/domain/             deterministic entities and rules
internal/discovery/          walkers, detectors, completeness
internal/scanplan/           project-to-adapter plans
internal/adapters/osv/       OSV integration
internal/process/            bounded subprocess execution
internal/artifact/           content-addressed raw results
internal/normalize/          PURLs, aliases, severity, stable keys
internal/policy/             later policy decisions
internal/remediation/        later planning and verification
internal/presenter/          terminal and machine formats
internal/store/              later persistence
schemas/                     versioned JSON Schema and examples
web/                         later React/Vite application
testdata/                    fixed repository and scanner fixtures
docs/                        architecture, ADRs, and contributor guidance
```

Domain packages must not import Cobra, SQL, or scanner-specific types. Entry points translate transport inputs into application commands. Adapters implement ports. Keep platform-specific behavior at transport/adapter boundaries; serialized domain meaning must remain equivalent across platforms. Do not add a server, web UI, database, or plugin system before the documented roadmap stage.

## Non-negotiable invariants

### Determinism and contracts

- Stable JSON is versioned and validates against the current schema.
- Stable keys exclude descriptions, timestamps, severity labels, and evidence ordering.
- Sort declared lists by their contract keys; results must be order-independent.
- Repository-relative paths use `/` in serialized output.
- Preserve scanner/database/tool provenance and raw artifact digests.
- Unknown data is preserved when safe; it must not become confidence or a safe result.
- Schema, CLI, error, adapter, or storage contract changes require explicit human approval and compatibility evidence.

### Completeness and failure

- Every detected target has an explicit completeness state.
- `complete` means all detected targets succeeded; `partial` means some succeeded and some were omitted, incomplete, or failed; `failed` means no trustworthy result or a core execution failure.
- Zero findings may be described as no known vulnerabilities only when status is `complete`.
- Missing tools, incompatible versions, non-zero exits, timeouts, malformed output, output limits, and artifact failures are errors—not empty successful scans.
- Partial success retains successful workspace results and emits stable workspace-scoped diagnostics.
- Use stable error codes from `requirements/ERROR-MODEL.md`, including `SCANNER_NOT_FOUND`, `SCANNER_TIMEOUT`, `SCANNER_OUTPUT_INVALID`, and `PATH_OUTSIDE_ROOT`.

### Security

Treat repository files, manifests, paths, symlinks, scanner output, and network responses as hostile inputs.

- Never concatenate shell commands or invoke a shell for scanner/package-manager operations; pass argument arrays directly.
- Set explicit working directories, approved environment variables, deadlines, cancellation, and bounded stdout/stderr capture.
- Resolve paths against a canonical repository root; reject traversal and symlink escape.
- Discovery is read-only and must not access the network.
- Do not install tools or run package install scripts automatically.
- Write atomically with restrictive permissions; never overwrite unless the contract permits it.
- Redact tokens, credential-bearing URLs, user information in credentials, and sensitive environment values.
- Do not upload source or log secrets.
- Read-only agent capabilities are the default. Any future mutation, publishing, or PR creation requires an independent permission scope, plan/dry-run output, exact commands/files, and verification evidence.

## CLI and I/O conventions

- Machine data goes to stdout; diagnostics, progress, and guidance go to stderr.
- `--format json` emits no banners or progress text.
- Color is enabled only for a TTY.
- CLI options and exit codes follow `docs/04-execution/deprail-v0.1-execution-package/requirements/CLI-CONTRACT.md`.
- Exit code `0` is successful execution, `2` is configuration/argument error, and `3` means scanner failure or incomplete result. Codes `1`, `4`, and `5` retain their documented meanings/reserved status.
- File output is atomic and safe.
- Keep human terminal output concise; verbose mode may expose additional evidence but never secrets.

## Development workflow

Work one issue at a time. Before implementing any task—especially a feature—first state the plan in the response before making code or configuration changes. The plan must name the goal, linked contracts/issues, intended files and symbols, out-of-scope changes, risk level, and verification commands or scenarios. Wait for no extra approval unless the task is materially ambiguous; the required plan is the normal first step, not a substitute for implementation.

After the plan, implement only the stated scope. If investigation changes the plan, stop before the next edit and state the revised plan and why. Keep scope bounded to the issue; do not silently add retries, telemetry, unrelated abstractions, or roadmap features.

Every implementation item should deliver code, tests, documentation, and evidence together. Use the issue format in `templates/ISSUE-TEMPLATE.md` and the PR format in `templates/PULL-REQUEST-TEMPLATE.md`. Record architectural decisions using `templates/ADR-TEMPLATE.md`.

### Branch and merge policy

- Keep `main` releasable and do not develop directly on it.
- Before creating or checking out an issue branch, synchronize the local baseline:
  ```bash
  git fetch origin main
  git switch main
  git pull --ff-only origin main
  git switch -c <issue-branch> main
  ```
  If switching or fast-forwarding fails, preserve local work and report the blocker; never reset or overwrite changes to force synchronization.
- Create one short-lived branch per issue or bounded task from the synchronized local `main`.
- Push the branch and open a pull request linked to the issue; do not push implementation commits directly to `main`.
- Keep one primary outcome per pull request. Include scope, risk, contract impact, verification results, evidence, and rollback notes.
- Merge only after required CI and human review pass. Use a squash merge tied to the issue ID, then delete the branch.
- Emergency security changes may use an expedited path, but still require a linked issue, review, verification, and a follow-up record.

### Commit and pull request traceability

- Every commit on an issue branch should reference exactly one issue when practical.
- Use this commit format:

  ```text
  <type>(<issue-key>): <imperative summary> (#<github-issue-number>)
  ```

  Examples: `feat(s0-001): establish repository baseline (#26)`, `fix(disc-002): reject symlink escape (#7)`, and `docs(s0-003): add pull request templates (#28)`.
- Squash commit titles must contain the GitHub issue reference, such as `(#27)`.
- Every pull request must link its issue with `Closes #N` or `Refs #N`.
- Every pull request must carry matching `area`, `risk`, `priority`, and `type` labels from the project label set.
- Before requesting review, agents must check that the pull request has the required labels and issue link.
- A missing issue reference or required label is a process defect; fix it before review or merge rather than deferring it.

### Project status synchronization

GitHub does not currently move DepRail Project items automatically when a pull request opens. Agents MUST synchronize the project item explicitly:

- When implementation starts, add the issue to the `DepRail` project if absent and set Project Status to `In Progress`.
- Immediately after opening the pull request, set the linked issue's Project Status to `Review`.
- Before requesting review, verify the issue link, labels, project membership, and Project Status with `gh`.
- During PR review, the owner reviews the acceptance criteria and required evidence. After the reviewed PR merges and required CI/evidence pass, the agent or maintainer may set Project Status to `Done`; no separate post-merge owner-acceptance step is required.
- If project-write permission is unavailable, report the exact missing permission and leave the issue status unchanged; never claim synchronization occurred.

Resolve the project, item, Status field, and option IDs from the live project rather than hard-coding environment-specific IDs:

```bash
gh project list --owner geoffrey-xiao --format json
gh project item-list <project-number> --owner geoffrey-xiao --format json
gh project field-list <project-number> --owner geoffrey-xiao --format json
gh project item-edit --project-id <project-id> --id <item-id> --field-id <status-field-id> --single-select-option-id <review-option-id>
```

The final `item-edit` command is required after PR creation because the current repository has no automatic PR-to-Project status mutation.
The complete creation and verification sequence is documented in [`docs/04-execution/deprail-v0.2-execution-package/tracking/WORKFLOW-GUIDE.md`](docs/04-execution/deprail-v0.2-execution-package/tracking/WORKFLOW-GUIDE.md). Agents MUST verify labels after both issue and PR creation; command flags are not evidence that labels were applied.

### Release version baseline

Before starting development for a new product-version line such as `0.2.0`, agents MUST establish the current release baseline:

- Fetch the current branch and all tags: `git fetch origin main --tags`.
- Verify the synchronized `origin/main` commit and inspect the latest stable tag and any preview or release-candidate tags:
  ```bash
  git log --oneline -1 origin/main
  git tag --list 'v*' --sort=-v:refname
  gh release list --repo geoffrey-xiao/deprail --limit 20
  git show origin/main:.release-please-manifest.json
  ```
- Treat the latest stable `vMAJOR.MINOR.PATCH` tag as the release baseline. Preview tags such as `v0.1.0-preview.1` and release candidates such as `v0.1.0-rc.1` do not replace the stable baseline.
- Record the next intended version in the issue or release plan before implementation. For example, `0.2.0` follows the completed `0.1.x` line; it must not be inferred from an arbitrary branch or unpublished local tag.
- Create release tags only from a reviewed, CI-passing `main` commit. Tags are immutable: never move, overwrite, or reuse a version tag.
- A bug fix after `v0.1.0` increments the patch version (`v0.1.1`); a new backward-compatible product line increments the minor version (`v0.2.0`).
- If no prior stable tag exists, document the initial version decision explicitly before creating the first preview or stable tag.

### Issue lifecycle and closure

- Before starting a new issue, check the previous issue and pull request. Remind the project owner to close the previous issue if its acceptance, verification, review, and evidence are complete.
- Do not silently abandon or leave completed issues open. If the prior issue is incomplete, state the missing acceptance item and keep it open or mark it blocked with an owner and reason.
- When an issue is finished, prepare a completion report that identifies each acceptance item, required command, CI result, human review result, evidence link, and remaining risk. Record owner review from the PR when applicable.
- The project owner must inspect and check every acceptance item during PR review. An agent may move the Project item to `Done` after the reviewed PR merges and required verification/evidence are complete; do not claim completion before those conditions.
- The project owner may close issues manually after the reviewed PR merges and the required evidence is attached. A merged PR without review, CI, or required evidence does not prove completion.
- Link a merged pull request with a closing keyword where appropriate, and update the local checklist only with linked evidence and the completed review record.

Respect the execution controls:

- At most two implementation issues and two review PRs are active.
- Never parallelize two R3 items.
- One issue has one primary outcome.
- Keep generated or golden-file changes explained and human-reviewed.
- Do not alter `tracking/MASTER-CHECKLIST.md` until acceptance has linked evidence.
- Human review is mandatory for schemas, permissions, compatibility, external-process behavior, network/file writes, migrations, credentials, publishing, and broad golden updates.

The repository's intended verification entry points are:

```text
make bootstrap
make generate
make test
make lint
make build
make verify
make test-integration
```

These commands are part of the Sprint 0 baseline and may not exist in the documentation-only checkout yet. Do not claim they pass until the repository actually contains the toolchain and they have been run.

## Testing expectations

Use minimal, offline, license-safe, deterministic fixtures. Match the test level to the behavior:

- Unit tests: domain invariants and pure deterministic transforms.
- Contract tests: detectors, adapters, artifact store, and presenters.
- Golden tests: stable `ProjectGraph` and `ScanReport` serialization.
- Integration tests: controlled child processes, SQLite/file output when introduced.
- End-to-end tests: mixed JavaScript/Python/Java repositories.
- Property/fuzz tests: order independence, stable keys, parser and hostile-path behavior.
- Cross-platform smoke tests: Linux, macOS, and Windows semantics and packaging.

Cover happy, invalid, incomplete, partial, failed, boundary, repeated/deterministic, and adversarial cases as applicable. Required security cases include external symlinks, traversal, shell metacharacters, long paths, Unicode, malformed manifests, huge output, timeout, cancellation, process-tree termination, credential-bearing URLs, partial scanner failure, interrupted atomic writes, and incompatible scanner versions. Never weaken an assertion to fit an implementation or refresh broad goldens without explaining each semantic change.

## Current implementation order

Start with the Sprint 0 baseline:

1. Repository/module structure, license, contribution guidance, and ADR baseline (`S0-001`).
2. Pinned toolchain and fast GitHub Actions verification (`S0-002`).
3. Issue, PR, and ADR templates (`S0-003`).
4. v1alpha schemas and fixtures (`SCHEMA-001`, `FIXTURE-001`).
5. Discovery vertical slice: detector contract, secure walker, npm detector, `discover`, and golden tests.
6. Remaining language discovery, OSV adapter/process/artifact contracts, normalization, reporting, and doctor.

Do not begin web, team services, autonomous writes, or agent mutation features before the prior release gate is proven in real use.

## Completion and reporting

A change is complete only when its observable behavior, failure behavior, tests, documentation, contracts, and evidence are complete; local and CI verification agree; and human review records the decision and remaining risk. Reports must name changed behavior and files, commands actually run and their results, contract decisions, security/compatibility impact, remaining limitations, and any required human review. Never claim an unrun test, command, or manual validation passed.
