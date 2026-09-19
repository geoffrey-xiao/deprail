# v0.3 Epic and Issue Contract Crosswalk

This crosswalk verifies the local epic and issue decomposition against `DEVELOPMENT-PLAN-v0.3.md`, the functional requirements, traceability matrix, and Sprint 0 gate. It does not replace owner, architecture, or security review.

## Gate status

The plan's Definition of Ready is not yet complete. Issue #198 must record v0.2 carryover status, owner/reviewer decisions, schema/example review, adapter ownership, named fixtures/security scenarios, and tracking preparation before these contracts become implementation-ready or are imported into GitHub.

## Functional requirements

| Requirement | Covered by | Required evidence |
| --- | --- | --- |
| FR-030-001 explicit report and finding | V03-003, V03-008 | Argument, missing-input, and no-implicit-rescan tests |
| FR-030-002 repository/scan/artifact/workspace/state identity | V03-003 | Source-report fixtures and provenance assertions |
| FR-030-003 finding context and provenance | V03-001, V03-003, V03-005 | Domain and adapter contract tests |
| FR-030-004 read-only planning | V03-004, V03-008, V03-010 | Tree snapshots, process audit, mutation tests |
| FR-030-005 explicit candidate states | V03-001, V03-005 | Domain tests and state fixtures |
| FR-030-006 unknown withholds recommendation | V03-001, V03-006, V03-007 | Unknown/incomplete golden cases |
| FR-030-007 ownership, constraints, lockfiles, risks | V03-006, V03-007 | Ecosystem adapter fixtures and golden plans |
| FR-030-008 structured future commands | V03-001, V03-005, V03-006, V03-007 | Command-shape and working-directory assertions |
| FR-030-009 deterministic identity/serialization | V03-001, V03-002, V03-009 | Reordered-input property and golden tests |
| FR-030-010 source/report digest and stale rejection | V03-003, V03-008 | Stale report and repository-state tests |
| FR-030-011 shared ecosystem contract | V03-005, V03-006, V03-007 | Cross-adapter contract tests |
| FR-030-012 schema-valid JSON | V03-002, V03-009 | Schema validation and compatibility tests |
| FR-030-013 path/traversal/symlink safety | V03-004, V03-009 | Hostile path and safe-write tests |
| FR-030-014 explicit malformed/unsupported/incomplete handling | V03-003, V03-005, V03-006, V03-007, V03-008 | Stable diagnostics and no-recommendation fixtures |
| FR-030-015 human explanation output | V03-009 | Terminal presenter tests |
| FR-030-016 offline/local-first default | V03-004, V03-005, V03-010 | Network-disabled smoke and process audit |

## Development-plan delivery stages

| Plan stage | Epic/issue coverage | Exit evidence |
| --- | --- | --- |
| Planning Sprint | Issue #198 plus this crosswalk | Owner/reviewer decisions, frozen contracts, named fixtures |
| Sprint 1: domain model | EPIC-001; V03-001, V03-002 | Domain tests, schema, examples, golden and compatibility tests |
| Sprint 2: adapters | EPIC-003; V03-005 through V03-007 | npm/pnpm/Yarn, Python, Maven/Gradle read-only fixtures |
| Sprint 3: CLI/evidence | EPIC-002 and EPIC-004; V03-003, V03-004, V03-008, V03-009 | CLI, output, stale-input, safety, mixed-repository, determinism tests |
| Release hardening | V03-010 | Cross-platform smoke, fuzz/hostile inputs, network-disabled and mutation proof, review evidence |
| v0.4 handoff | V03-001, V03-002, V03-009, V03-010 | Executor-consumable plan schema, verification commands, rollback model |

## Traceability rows

| Traceability row | Issue coverage | Review owner |
| --- | --- | --- |
| Finding resolution and provenance | V03-001, V03-003 | Domain review |
| Direct/transitive ownership | V03-005, V03-006, V03-007 | Remediation review |
| Candidate compatibility and risks | V03-001, V03-006, V03-007 | Security/architecture review |
| No repository mutation | V03-004, V03-008, V03-010 | Security review |
| Stable plan identity | V03-001, V03-002, V03-009 | Compatibility review |
| Unknown/incomplete handling | V03-001, V03-003, V03-005 through V03-008 | Safety review |
| Schema-valid machine output | V03-002, V03-009 | Contract review |
| CLI behavior | V03-003, V03-008, V03-009 | CLI review |
| Cross-platform equivalence | V03-010 | Release review |
| v0.4 handoff | V03-001, V03-002, V03-009, V03-010 | Architecture review |

## Explicit plan boundaries

All issues remain limited to `finding -> candidate analysis -> reviewable plan`. They do not authorize manifest or lockfile edits, package installation, scripts, tests/builds/rescans, worktrees, patch/PR creation, publishing, web/team services, MCP writes, new scanners, policy-language work, source upload, telemetry, or autonomous approval.
