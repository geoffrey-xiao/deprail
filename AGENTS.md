# DepRail Agent Guide

## Purpose and source of truth

DepRail is a dependency security guardrail for multi-language repositories. The current checkout is documentation-first; implementation is expected to begin with Sprint 0. Treat the versioned contracts in `docs/` as normative until an approved code or ADR decision supersedes them.

Read in this order before changing behavior:

1. `docs/02-architecture/deprail-architecture-and-tech-stack-v1.md`
2. `docs/04-execution/deprail-v0.1-execution-package/PRD-v0.1.md`
3. `docs/04-execution/deprail-v0.1-execution-package/requirements/`
4. The linked issue in `docs/04-execution/deprail-v0.1-execution-package/issues/`
5. The applicable sprint and checklist in `docs/04-execution/deprail-v0.1-execution-package/tracking/`

The execution package is the active delivery contract. Do not mark checklist items complete without linked evidence.

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
