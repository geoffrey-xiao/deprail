# DepRail v0.3 Product Requirements

| Attribute | Value |
| --- | --- |
| Version | v0.3.0 |
| Status | Planning baseline |
| Outcome | Reviewable dependency remediation plans without repository mutation |
| Predecessor | v0.2 CI guardrail and preview |
| Owner | Project owner |
| Reviewer | Security/architecture reviewer |

## Goal

Help a developer understand which dependency change is viable, why it is risky, what files and commands a later change would involve, and how the result must be verified—without changing the repository.

## User jobs

- Select a finding from a trusted scan result.
- Identify direct or transitive dependency ownership.
- Compare viable non-vulnerable candidate versions.
- Detect constraint, peer, runtime, engine, lockfile, and major-version risks.
- Receive a deterministic plan suitable for human or future agent review.
- Understand why candidates are rejected or unavailable.

## Scope

- `deprail fix plan <finding-or-key>` with terminal and JSON output.
- Read-only npm/pnpm/Yarn, Python requirements/uv/Poetry, and Maven/Gradle planning adapters.
- Versioned remediation-plan schema with stable identity and provenance.
- Candidate states: `recommended`, `viable`, `rejected`, `unavailable`, `unknown`.
- Stale-input detection, explicit incomplete analysis, and safe plan output.
- Required future commands, files, risks, rollback, and verification steps represented but not executed.

## Out of scope

Repository mutation, package installation, package-manager scripts, worktrees, tests/builds/rescans, patches, pull requests, publishing, web/API history, MCP write tools, Trivy, new scanner families, and general policy-language work.

## Product invariants

- Deterministic core; no LLM is required for correctness.
- Unknown evidence never becomes a safe recommendation.
- No plan operation writes to the target repository.
- Structured commands use argument arrays, never shell interpolation.
- JSON is versioned, schema-valid, stable, and provenance-preserving.
- Local-first operation; no source upload by default.

## Acceptance

A representative mixed repository produces a schema-valid plan with direct/transitive ownership, candidates, risks, affected files, structured future commands, verification requirements, provenance, and no repository changes. Malformed, stale, unsupported, incomplete, and adversarial inputs fail truthfully.
