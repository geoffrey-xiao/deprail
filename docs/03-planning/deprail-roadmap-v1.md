# DepRail Whole-Project Roadmap

| Attribute | Value |
| --- | --- |
| Version | 1.0 |
| Status | Current high-level product and architecture roadmap |
| Product baseline | `docs/01-product/deprail-product-design-v1-ai.md` |
| Architecture baseline | `docs/02-architecture/deprail-architecture-and-tech-stack-v1.md` |
| Release execution | One version-specific development plan per release |

## Purpose

This document defines the whole-project sequence from the scanning foundation to the v1.0 stable release. It is intentionally high level. It does not replace a release development plan, execution package, epic, or issue contract.

## Source hierarchy

```text
Product design + architecture
          ↓
This whole-project roadmap
          ↓
Release-specific development plan
          ↓
Release execution package
          ↓
Epics and issue contracts
          ↓
PRs, verification, and evidence
```

The product design and architecture documents define product intent and system boundaries. A release-specific plan selects the next bounded scope without silently changing those sources. Any contradiction requires an explicit decision and updates to the affected documents before issue creation.

## Release roadmap

| Release | Primary outcome | Main capabilities | Decision gate |
| --- | --- | --- | --- |
| v0.1 | Multi-language scanning foundation | Discovery, OSV-Scanner adapter, artifact retention, normalization, terminal/JSON output | Mixed-repository scans are stable and failures are explicit |
| v0.2 | CI guardrail and public preview | Baselines, diff, policy gates, exceptions, SARIF, GitHub Action, pull-request validation | New-risk gates work on real pull requests |
| v0.3 | Remediation planning | Reviewable dependency upgrade plans without repository mutation | Plans are accurate and risks are explicit |
| v0.4 | Remediation and verification | Isolated application, tests/builds, rescan, patch evidence | Failed changes remain contained and recoverable |
| v0.5 | Local web and history | SQLite history, local API, embedded React console | Repeated local use justifies persistent history |
| v0.6 | Team collaboration | PostgreSQL, identity, RBAC, exceptions, audit, notifications | Multi-user workflows are reliable and auditable |
| v0.7 | Safe agent access | Versioned MCP tools, scopes, approvals, cancellation, audit | Agents cannot exceed granted capabilities |
| v0.8 | Extension ecosystem | Trivy, open formats, process-isolated external adapters | Extensions remain secure and contract-tested |
| v0.9 | Release candidate | Contract freeze, migrations, performance, security, release rehearsals | Release-candidate gates pass |
| v1.0 | Stable product release | Stable contracts, support policy, signed artifacts, operations guidance | Explicit go/no-go decision |

## Cross-release invariants

Every release preserves these boundaries unless an approved compatibility decision changes them:

- Local-first operation; no default source upload.
- Deterministic domain logic independent of an LLM.
- External tools behind bounded, versioned adapter contracts.
- Canonical path containment and safe process invocation.
- Explicit complete, partial, and failed outcomes.
- Versioned machine output with stable keys and provenance.
- Reviewable changes; no autonomous mutation or merge.
- Product security behavior stays in application services, not prompts or Skills.
- Cross-platform semantic equivalence on Linux, macOS, and Windows.

## Scope gates

- Baseline comparison, policy, SARIF, and GitHub Action are v0.2 capabilities.
- Remediation planning starts in v0.3 and does not mutate repositories.
- Repository mutation and verification begin only in v0.4 with isolation and rollback.
- Web and team services are not prerequisites for the local core.
- MCP and agent write capabilities require explicit scopes, approval, and audit.
- New scanner families and broad security domains require separate compatibility and security decisions.

## v0.3.1 CLI UX refinement

The v0.3.1 preview line is the bounded follow-up to the v0.3 remediation-planning gate. It refines human-facing CLI presentation through structured application events, shared terminal rendering, TTY-aware progress, truthful complete/partial/failed summaries, safe hostile-label rendering, and cross-platform terminal evidence. It does not change scanner semantics, machine-readable schemas, repository mutation boundaries, or v0.4 verification scope.

The release-specific planning baseline is [`deprail-development-plan-v0.3.1.md`](deprail-development-plan-v0.3.1.md). Its execution package is [`deprail-v0.3.1-execution-package`](../04-execution/deprail-v0.3.1-execution-package/README.md). Both remain planning documents until owner and architecture/security review approve the Definition of Ready.

## Release planning rule

Before creating a new release's epics or issues, the owner must:

1. Read this roadmap and the current product/architecture baselines.
2. Create or update the release-specific development plan.
3. Reconcile the plan with product intent, architecture boundaries, prior release evidence, and deferred scope.
4. Freeze the release goal, included capabilities, exclusions, dependencies, risks, and acceptance gate.
5. Create or update the release execution package.
6. Only then create GitHub epics, issues, milestones, Project views, and parent-child links.

A release issue must map to a frozen release-plan item and at least one product, architecture, or requirements contract. No implementation starts from an unreviewed or contradictory plan.

## Change control

A material change to product stage, architecture boundary, release outcome, or deferred capability requires:

- A new or revised release-specific development plan.
- An explicit decision record explaining the change and displaced work.
- Synchronized product, architecture, requirements, execution, and tracking documents.
- Reconciliation of existing epics, issues, milestones, Project fields, and acceptance evidence.

Historical plan versions remain preserved. Do not silently rewrite a prior plan to hide a scope change.
