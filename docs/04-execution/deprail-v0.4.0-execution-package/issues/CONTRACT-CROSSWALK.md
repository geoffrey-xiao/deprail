# v0.4 Epic and Issue Contract Crosswalk

This crosswalk maps the v0.4 context documents to the proposed local epics and issue contracts. It does not replace owner, architecture, security, or GitHub tracking approval.

## Gate status

The supplemental context package is under review in PR #312. Baseline generation is an additive scope change recorded in ADR-0003 and tracked by [#329](https://github.com/geoffrey-xiao/deprail/issues/329). The existing v0.4 Definition of Ready and ADR-0002 records remain authoritative. These proposed contracts are not implementation-ready until the owner accepts this crosswalk and the GitHub hierarchy is reconciled.

## Functional requirements

| Requirement | Covered by | Required evidence |
|---|---|---|
| FR-001 valid plan and approval binding | EPIC-002; V04-002 | stale, changed, expired, reused approval fixtures |
| FR-002 dry-run is non-mutating | EPIC-002; V04-002, V04-003 | source/workspace tree comparison |
| FR-003 isolated supported mutation | EPIC-001, EPIC-002; V04-001, V04-003 | path authorization and adapter audit |
| FR-004 bounded subprocess execution | EPIC-002; V04-003 | timeout, cancellation, output-limit, non-zero evidence |
| FR-005 bounded verification | EPIC-003; V04-004 | command discovery and incomplete verification fixtures |
| FR-006 deterministic before/after classification | EPIC-003; V04-004 | resolved/residual/introduced/unknown golden cases |
| FR-007 failure preserves source and evidence | EPIC-001, EPIC-004; V04-001, V04-005 | rollback, cleanup, and source digest evidence |
| FR-008 versioned stdout/stderr separation | EPIC-004; V04-005 | JSON contract and CLI smoke |
| FR-009 redaction and content-addressed evidence | EPIC-004; V04-005 | secret/redaction and artifact tests |
| FR-010 platform-equivalent meaning | EPIC-005; V04-006 | Linux/macOS/Windows matrix |
| FR-011 complete scan to baseline conversion | Release parent #299; V04-007 | schema validation, deterministic output, diff/policy consumption |
| FR-012 reject incomplete or invalid scans | Release parent #299; V04-007 | failure fixtures, stable diagnostics, no artifact |
| FR-013 safe baseline persistence | Release parent #299; V04-007 | atomicity, permissions, existing-output, source immutability |

## Sequence and dependencies

| Stage | Epic | Issues | Exit evidence |
|---|---|---|---|
| 1 | EPIC-001 isolated workspace | V04-001 | source immutability, hostile paths, cleanup/rollback |
| 2 | EPIC-002 approved mutation | V04-002, V04-003 | approval and bounded-process contract evidence |
| 3 | EPIC-003 verification/rescan plus baseline producer | V04-004, V04-007 | baseline conversion and finding-transition evidence |
| 4 | EPIC-004 evidence/failure | V04-005 | schema, redaction, artifact, terminal-outcome evidence |
| 5 | EPIC-005 platform acceptance | V04-006 | cross-platform manual and CI evidence |

Dependencies flow in sequence. EPIC-004 consumes outputs from EPIC-001 through EPIC-003. EPIC-005 validates the complete workflow.

## Explicit boundaries

No issue authorizes caller-worktree mutation, arbitrary shell commands, unbounded network/scripts, new ecosystems, automatic approval, commit/push/PR/merge, remote publication, web/team services, MCP writes, or source upload.
