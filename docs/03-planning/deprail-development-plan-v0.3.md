# DepRail v0.3 Development Plan

| Attribute | Value |
| --- | --- |
| Release | v0.3 |
| Plan revision | 1.0 |
| Status | Draft for owner and architecture review |
| Whole-project roadmap | `docs/03-planning/deprail-roadmap-v1.md` |
| Product baseline | `docs/01-product/deprail-product-design-v1-ai.md` |
| Architecture baseline | `docs/02-architecture/deprail-architecture-and-tech-stack-v1.md` |
| Prior release plan | `docs/03-planning/deprail-development-plan-v0.2.md` |
| Delivery model | One project owner working with AI agents |
| Cadence | Two-week Sprints; one-week planning and contract reconciliation |

## 1. Purpose and reconciliation decision

The product design, architecture baseline, and whole-project roadmap agree that v0.3 is the remediation-planning stage. It must help a developer understand which dependency change is viable and produce a reviewable plan without changing the repository.

The v0.3 boundary is deliberately narrower than the complete remediation workflow:

```text
v0.3: finding -> candidate analysis -> reviewable plan
v0.4: approved plan -> isolated application -> verification -> rescan -> patch evidence
```

No v0.3 implementation may mutate a repository, install packages, invoke package-manager install scripts, create a worktree, open a pull request, publish results, or ask an LLM to decide security policy. Those capabilities remain deferred to v0.4 or later contracts.

This plan is a release-design artifact only. It does not create epics or GitHub issues. Issue creation begins after owner approval, v0.2 release-gate reconciliation, and execution-package updates.

## 2. Release outcome

A user can select a normalized finding and run a deterministic command that produces one or more candidate remediation plans. Each plan explains:

- the finding and affected component;
- direct or transitive dependency relationship;
- current and proposed versions;
- files that would change in a later application stage;
- package-manager commands that a later executor would run;
- expected lockfile or manifest effects;
- compatibility and major-version risks;
- required verification commands;
- rollback guidance for the later application stage;
- evidence, assumptions, and confidence;
- why a candidate was rejected or marked unavailable.

The plan is JSON-schema-valid, stable across repeated runs, safe to review or store, and contains no repository mutation.

## 3. User problems to solve

| User question | v0.3 answer |
| --- | --- |
| Which dependency should change? | Identify the direct owner or transitive path responsible for the finding. |
| Which versions are viable? | Enumerate non-vulnerable candidates compatible with declared constraints and known package metadata. |
| Is this a breaking change? | Classify major-version jumps and constraint changes as explicit migration risk. |
| What will change later? | List manifests, lockfiles, package-manager commands, and expected dependency effects without editing them. |
| What should I verify? | Produce deterministic tests, type checks, builds, rescans, and policy checks from repository evidence and configuration. |
| Why was a candidate rejected? | Preserve structured rejection reasons such as unavailable metadata, incompatible constraints, unsupported manager, or unresolved transitive owner. |
| Can an agent apply this now? | No. v0.3 returns a plan and requires a later explicit approval/apply capability. |

## 4. Scope

### 4.1 Finding selection and context

- Accept a finding identifier or a stable finding key from a scan/report artifact.
- Resolve the finding to a workspace, component, vulnerability aliases, affected version, fixed versions, dependency path, and provenance.
- Reject ambiguous or stale findings with stable errors; do not silently select a different finding.
- Support plans from retained local scan artifacts without requiring network access unless a future, explicitly enabled metadata provider is used.
- Preserve baseline and policy context when it affects prioritization or verification requirements.

### 4.2 Candidate analysis

- Prefer the minimum non-vulnerable version that satisfies declared constraints.
- Distinguish direct upgrades from transitive remediation through a direct dependency.
- Detect major-version, peer-dependency, runtime, engine, and lockfile compatibility risks when metadata is available.
- Represent multiple candidates rather than silently selecting the newest version.
- Record unavailable evidence as unknown; unknown must not become a safe or low-risk result.
- Keep candidate ordering deterministic and explain the ordering key.

### 4.3 Package-manager planning

Initial v0.3 planning adapters cover the ecosystems already supported by v0.2:

- npm/pnpm/Yarn;
- Python requirements/uv/Poetry;
- Maven/Gradle.

Each adapter is read-only in v0.3 and provides:

- dependency ownership and path interpretation;
- manifest and lockfile locations;
- version and constraint parsing;
- candidate compatibility checks;
- command templates for a future executor;
- expected files and lockfile effects;
- known verification command discovery.

Adapters must not execute package-manager mutation commands during planning.

### 4.4 Plan and explanation output

Provide:

```text
deprail fix plan <finding-or-key>
deprail fix plan <finding-or-key> --format json
deprail fix plan <finding-or-key> --output plan.json
```

Human output emphasizes the recommended candidate, risks, affected files, and verification steps. JSON is the normative machine contract. Diagnostics remain on stderr and machine data remains on stdout unless `--output` is used.

A plan contains a stable schema version, plan ID, source scan identity, finding identity, candidate list, selected recommendation, affected files, commands, risks, assumptions, verification steps, rollback description, and provenance.

### 4.5 Deterministic local storage

- Store plans as explicit versioned JSON files under the existing local data boundary when persistence is requested.
- Use content-addressed or stable plan identity where appropriate.
- Never overwrite an existing plan without an explicit output path and contract permission.
- Include source scan/report digest and repository state needed to detect stale plans.
- Do not store credentials, full environment variables, or source content unnecessarily.

## 5. Explicit exclusions

The following are not v0.3 deliverables:

- editing manifests or lockfiles;
- package-manager install, update, or script execution;
- repository mutation of any kind;
- isolated worktree creation;
- tests, builds, rescans, or verification execution by the remediation engine;
- automatic patch or pull-request creation;
- automatic approval, merge, or publication;
- web UI, hosted API, SQLite history, or team collaboration;
- MCP or agent write tools;
- Trivy or a new scanner family;
- a general-purpose policy language;
- autonomous vulnerability explanations that replace structured evidence;
- source upload or telemetry;
- a promise that every finding has a viable upgrade plan.

The plan may list future verification and application commands, but v0.3 does not execute them.

## 6. Architecture impact

### 6.1 New or expanded modules

| Module | v0.3 responsibility |
| --- | --- |
| `internal/remediation` | Finding context, candidate model, plan generation, risk classification, stable plan identity |
| `internal/remediation/npm` | Read-only npm/pnpm/Yarn constraint and ownership planning |
| `internal/remediation/python` | Read-only requirements/uv/Poetry planning |
| `internal/remediation/java` | Read-only Maven/Gradle planning |
| `internal/presenter` | Human and JSON remediation-plan output |
| `schemas` | Versioned remediation-plan schema and examples |
| `internal/artifact` | Plan/source evidence digests and safe output handling |
| `internal/app` | `fix plan` command orchestration and stale-input validation |

Names are illustrative; the execution package must preserve the inward dependency rule and avoid scanner-specific types in the domain.

### 6.2 Dependency direction

```text
CLI entry point
    -> application remediation service
        -> domain plan and candidate rules
            -> package-manager planning ports
                -> read-only filesystem/package metadata adapters
        -> presenters and plan storage
```

The remediation domain must not import Cobra, SQL, scanner-specific response types, or an LLM SDK. Package-manager adapters return normalized domain evidence through ports.

### 6.3 Trust boundaries

Repository manifests, lockfiles, package metadata, scanner findings, and registry responses are hostile inputs. Planning must:

- use canonical repository-root containment;
- reject traversal and symlink escape;
- parse rather than interpolate commands;
- preserve command arguments as arrays or structured command objects;
- cap metadata and parser input sizes;
- apply network access only through an explicit future permission;
- never execute returned scripts or shell fragments;
- redact credential-bearing URLs and environment values;
- distinguish missing metadata from safe compatibility.

## 7. Data and schema contracts

### 7.1 Remediation plan entity

The normative plan schema should include at least:

```text
schema_version
plan_id
created_from
repository_identity
workspace_identity
finding_identity
component
current_state
candidates
recommendation
affected_files
commands
risks
assumptions
verification
rollback
provenance
```

`created_from` includes the source scan/report digest, scanner identity, database/tool versions when available, and repository state. `plan_id` excludes timestamps and presentation wording so repeated equivalent inputs produce the same stable identity.

### 7.2 Candidate states

Candidates and recommendations must use explicit states, for example:

```text
recommended
viable
rejected
unavailable
unknown
```

A candidate cannot be `recommended` when required compatibility evidence is unknown. A plan may be emitted with no recommendation if analysis is incomplete, but it must explain why.

### 7.3 Compatibility and migration risk

The plan must preserve separate facts for:

- vulnerable versus non-vulnerable version range;
- declared constraint satisfaction;
- lockfile resolution availability;
- direct versus transitive ownership;
- major/minor/patch change;
- peer/runtime/engine compatibility;
- migration notes or missing migration evidence;
- expected verification commands.

Do not collapse these facts into one opaque confidence score.

### 7.4 Schema evolution

- Additive changes are allowed within the v0.3 schema line.
- Breaking changes require a new schema version, converter, examples, and compatibility tests.
- Golden examples must cover recommended, multiple-candidate, rejected, unavailable, stale, and incomplete plans.
- Descriptions, timestamps, evidence ordering, and presentation wording must not affect stable plan identity.

## 8. CLI and failure contract

Expected exit behavior:

| Condition | Result |
| --- | --- |
| Valid plan with a recommendation | Exit `0` |
| Valid plan with no viable candidate but complete analysis | Exit `0` with explicit no-recommendation state |
| Stale or missing finding input | Configuration/input error, stable code, non-zero exit |
| Unsupported ecosystem or package manager | Explicit unsupported/incomplete plan state; never a safe success |
| Malformed manifest or lockfile | Explicit planning failure or incomplete state with diagnostics |
| Missing package metadata | Unknown/unavailable evidence; no false-safe recommendation |
| Output write or digest failure | Error; do not claim a saved plan |
| Any attempted mutation path | Reject before execution |

The exact numeric mapping must be reconciled with the v0.2 CLI and error contracts before implementation. v0.3 must not silently repurpose existing exit codes.

## 9. Delivery stages

### Planning Sprint — contract reconciliation

- Confirm v0.2 stable gate status and carryover risks.
- Update the v0.3 execution package, requirements, traceability, and schema index.
- Freeze plan schema, candidate states, adapter ports, failure codes, and mutation boundary.
- Produce fixtures for each supported ecosystem and adversarial inputs.

### Sprint 1 — domain plan model

- Implement deterministic finding context and remediation-plan entities.
- Implement candidate state, stable plan identity, risk facts, and stale-input checks.
- Add schema, examples, golden serialization, and compatibility tests.

### Sprint 2 — package-manager planning adapters

- Implement read-only npm/pnpm/Yarn planning.
- Implement read-only Python planning.
- Implement read-only Maven/Gradle planning.
- Cover direct/transitive ownership, constraints, lockfiles, major-version risks, and unavailable metadata.

### Sprint 3 — CLI and evidence output

- Add `deprail fix plan` application orchestration.
- Add terminal and JSON presenters.
- Add safe output and optional local plan storage.
- Add mixed-repository and repeated-run determinism scenarios.

### Release hardening

- Run cross-platform smoke tests.
- Fuzz parsers and hostile plan inputs.
- Prove no repository mutation, package-manager execution, or network access occurs in default planning mode.
- Review schema compatibility and generated examples.
- Publish v0.3 preview only after owner and security review.

## 10. Test and evidence strategy

### Unit tests

- Candidate ordering and minimum-safe-version selection.
- Constraint and major-version classification.
- Direct/transitive ownership rules.
- Stable plan identity and order independence.
- Unknown versus unavailable versus rejected state handling.
- Stale scan/report detection.

### Contract tests

- Each package-manager adapter returns the same domain contract.
- Commands are structured arguments, never shell strings requiring interpolation.
- Adapter failures preserve explicit diagnostics and provenance.
- Plans validate against the versioned JSON Schema.

### Golden tests

Cover:

- direct patch upgrade;
- transitive upgrade through a direct owner;
- multiple viable candidates;
- major-version migration;
- peer/runtime incompatibility;
- no viable candidate;
- missing registry metadata;
- malformed manifest/lockfile;
- stale finding/report;
- mixed JavaScript/Python/Java repository.

### Security tests

- path traversal and symlink escape;
- shell metacharacters and command-argument boundaries;
- malicious manifest values;
- oversized manifests, lockfiles, and metadata;
- credential-bearing registry URLs;
- package-manager script fields;
- network-disabled default planning;
- interrupted and restrictive plan writes;
- attempted repository mutation;
- deterministic output under reordered input evidence.

### Cross-platform evidence

Linux, macOS, and Windows must produce semantically equivalent plan JSON for equivalent fixtures. Platform-specific command rendering may differ only where the contract explicitly permits it.

## 11. Security and compatibility risks

| Risk | Response | Release gate |
| --- | --- | --- |
| Planner recommends a version without enough compatibility evidence | Use explicit unknown state and withhold recommendation | Golden and adversarial candidate tests |
| Package-manager metadata causes command injection | Structured arguments, no shell, no execution in v0.3 | Process and hostile-input tests |
| A transitive fix requires a breaking direct upgrade | Classify migration risk and require review | Adapter contract evidence |
| Plans become stale after repository changes | Bind plan to source scan/repository identity and reject stale input | Stale-plan tests |
| Registry/network behavior makes plans nondeterministic | Offline-first metadata policy and provenance; no silent fallback | Repeated-run tests |
| AI treats a plan as approval to mutate | Plan schema and CLI state explicitly read-only; apply deferred to v0.4 | Interface and permission review |
| New package-manager support expands scope | Limit v0.3 to v0.2 ecosystems | Change-control review |
| Schema changes break agents or future web clients | Versioned schema, examples, converters, compatibility tests | Schema review |

## 12. Definition of Ready

Before creating v0.3 epics or issues:

- v0.2 release status and carryover risks are recorded;
- product, architecture, roadmap, and this plan agree;
- the v0.3 execution package exists;
- remediation-plan schema and examples are reviewed;
- candidate and failure contracts are explicit;
- package-manager adapter ownership is assigned;
- mutation, network, write, and approval boundaries are explicit;
- fixtures and security test scenarios are named;
- owner and reviewer are assigned;
- milestone, Project fields, and release tracking are prepared.

## 13. Definition of Done

v0.3 is complete when:

- `deprail fix plan` produces accurate, reviewable plans for supported ecosystems;
- no default plan operation mutates files, executes package managers, installs dependencies, or publishes data;
- candidate risks and unknown evidence are explicit;
- plan JSON is versioned, schema-valid, deterministic, and provenance-preserving;
- terminal output is concise and machine output is stable;
- stale, malformed, unsupported, incomplete, and adversarial inputs fail truthfully;
- Linux, macOS, and Windows semantics are equivalent;
- representative mixed-repository evidence is attached;
- documentation explains that v0.3 plans do not apply changes;
- owner and security reviewer approve the release evidence.

## 14. Decision gate

The v0.3 preview may proceed only when plan accuracy is demonstrated on representative repositories and the mutation boundary is proven. Stable v0.3 requires explicit evidence that no plan path silently edits the repository or executes package-manager operations.

The v0.4 handoff begins only after the v0.3 plan schema, candidate evidence, verification command model, and rollback model are stable enough for an isolated executor to consume them without reinterpretation.
