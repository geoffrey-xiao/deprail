# DepRail v0.2 Development Plan

| Attribute | Value |
| --- | --- |
| Release | v0.2 |
| Plan revision | 1.0 |
| Status | Current release-specific development plan |
| Whole-project roadmap | `deprail-roadmap-v1.md` |
| Product baseline | `docs/01-product/deprail-product-design-v1-ai.md` |
| Architecture baseline | `docs/02-architecture/deprail-architecture-and-tech-stack-v1.md` |
| Prior planning snapshot | `deprail-development-plan-v1.md` |
| Delivery model | One project owner working with AI agents |
| Cadence | Two-week Sprints; one-week Sprint 0 |

## 1. Purpose and reconciliation decision

This revision keeps the product and architecture roadmap aligned with the execution package. The product and architecture documents define v0.2 as the CI guardrail stage: baseline diff, policy gates, SARIF, and GitHub Action preview.

The existing v0.2 hardening issues remain valid prerequisites. They are not, by themselves, the complete v0.2 product outcome. Existing V02 issue identifiers are retained for traceability; they are not renumbered.

## 2. Canonical stage boundaries

| Stage | Outcome | Required evidence |
| --- | --- | --- |
| v0.1 foundation closure | Discovery, OSV adapter, normalization, terminal/JSON output, safe failures, deterministic contracts | Mixed-repository scan, adapter fixtures, schema-valid output, cross-platform verification |
| v0.2 CI guardrail | Baselines, diff, policy, SARIF, GitHub Action, real pull-request validation | New-risk gate works on real pull requests |
| v0.3 remediation planning | Reviewable plans without repository mutation | Accurate plans and explicit risks |
| v0.4 remediation and verification | Isolated application, verification, rescan, patch delivery | Failed changes remain contained |
| v0.5 local web/history | SQLite history, local API, embedded React | Repeated local use justifies history |
| v0.6 team collaboration | PostgreSQL, identity, RBAC, exceptions, audit | Multi-user workflow is reliable |
| v0.7 agent access | MCP tools, scopes, approvals, audit | Agents cannot exceed granted capabilities |
| v0.8 extensions | Trivy, open formats, external adapter protocol | Extensions remain isolated and contract-tested |
| v0.9 release candidate | Contract freeze, performance, security, migration, release hardening | RC gates pass |
| v1.0 stable | Stable contracts and production release | Go/no-go decision |

## 3. v0.2 prerequisite hardening

The current execution package's EPIC-001 through EPIC-004 work is retained as the prerequisite hardening gate:

- EPIC-001: empty collections, strict arguments, truthful version identity, release-mode checklist.
- EPIC-002: real OSV-Scanner v2 fixtures, exit-code behavior, requested-root execution.
- EPIC-003: outside-root regression, report examples, deterministic artifacts and ordering.
- EPIC-004: release identity, release-mode smoke, artifact and supply-chain evidence.

These tasks protect the v0.2 guardrail from false-safe results, unstable output, and unreproducible scanner behavior. They must not be marked as the complete v0.2 product outcome.

## 4. v0.2 CI guardrail work

### S5 — Baselines and diff

Add a new epic and implementation issues for:

- Baseline scan/result representation.
- Baseline storage and compatibility.
- Base-versus-head comparison.
- New, resolved, and unchanged finding classification.
- Dependency additions, upgrades, removals, and affected-workspace changes.
- Deterministic diff JSON and terminal output.
- `deprail diff` command and exit behavior.

### S6 — Policy, exceptions, and SARIF

Add a new epic and implementation issues for:

- Typed policy evaluation through the application layer.
- New/high/critical finding gates.
- Required-complete-scan behavior.
- Expiring exceptions with rationale, scope, approver, and review condition.
- Stable policy exit behavior.
- `deprail policy check` command.
- SARIF output with schema and consumer-validation coverage.

Policy must remain deterministic and must not be implemented as prompt logic or a general policy language prematurely.

### S7 — GitHub Action and public preview

Add a new epic and implementation issues for:

- GitHub Action packaging.
- Least-privilege permissions and caching documentation.
- Pull-request base/head scan flow.
- Artifact and evidence publication behavior.
- Real pull-request validation.
- Cross-platform release and installation smoke.
- v0.2 public-preview acceptance and rollback evidence.

## 5. Required cross-document updates

Before creating the new CI-guardrail implementation issues, update the v0.2 execution package so these documents agree with this plan:

- `PRD-v0.2.md`
- `ARCHITECTURE-v0.2.md`
- `requirements/FUNCTIONAL-REQUIREMENTS.md`
- `requirements/CLI-CONTRACT.md`
- `requirements/DATA-CONTRACT.md`
- `requirements/ERROR-MODEL.md`
- `requirements/TRACEABILITY.md`
- `tracking/MASTER-CHECKLIST.md`
- v0.2 epic and issue indexes

The v0.2 PRD must no longer classify baseline comparison, policy gates, or SARIF as unapproved candidate capabilities. Candidate capability decisions remain separate for capabilities not in the product/architecture v0.2 baseline.

## 6. Execution controls

- One primary issue per PR.
- Maximum two implementation issues and two review PRs active.
- No more than one high-risk task in implementation.
- Existing V02 identifiers remain stable.
- Every issue maps to a product or architecture requirement.
- Every PR includes tests, documentation, evidence, compatibility impact, security impact, and rollback.
- No implementation begins for a new S5-S7 capability until its reconciled contract is present in the execution package.

## 7. Decision gate

The v0.2 release is not complete until:

- Prerequisite hardening is complete.
- Baselines and diff work on real branch comparisons.
- Policy gates classify new risk and incomplete scans correctly.
- SARIF output validates for supported consumers.
- The GitHub Action works on real pull requests with documented permissions and caching.
- Cross-platform verification and release evidence are complete.
- The owner records the remaining-risk and public-preview decision.

## 8. Historical plan

`deprail-development-plan-v1.md` remains preserved as the prior planning snapshot. It must not be silently rewritten; future changes belong in a new version or an explicit decision record.
