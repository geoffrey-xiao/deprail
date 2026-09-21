# EPIC-006: Isolated Fix Apply Orchestration

**GitHub epic:** [#332](https://github.com/geoffrey-xiao/deprail/issues/332)

**Target:** `v0.4.0-preview.2`

**Status:** Todo; child issues created for preview.2 implementation

## Outcome

An explicitly approved, versioned remediation plan can be applied only in a temporary detached worktree, verified, rescanned, classified, and represented by durable redaction-safe evidence. The caller repository remains unchanged.

## Boundaries

Includes plan and approval validation, dry-run, canonical source binding, temporary worktree lifecycle, allowlisted adapter-owned mutations, bounded verification and rescan, transition classification, rollback/discard, cleanup, stable terminal/JSON output, and cross-platform evidence. Excludes caller-worktree mutation, automatic commit/push/PR/merge, autonomous approval, arbitrary shell commands, unbounded scripts/network access, new ecosystems, remote publishing, policy-language expansion, and release publication.

## Dependencies

EPIC-001 through EPIC-005, `deprail fix plan`, the versioned plan/evidence contracts, and baseline generation (#329/#331).

## Child issues

Implementation order is dependency-aware; security and cross-platform validation closes the epic.

1. [#337 — Validate apply plans and approval binding](https://github.com/geoffrey-xiao/deprail/issues/337)
2. [#338 — Implement isolated detached worktree lifecycle](https://github.com/geoffrey-xiao/deprail/issues/338)
3. [#336 — Implement allowlisted adapter mutations](https://github.com/geoffrey-xiao/deprail/issues/336)
4. [#340 — Add bounded verification and rescan orchestration](https://github.com/geoffrey-xiao/deprail/issues/340)
5. [#339 — Classify remediation finding transitions](https://github.com/geoffrey-xiao/deprail/issues/339)
6. [#345 — Implement rollback, discard, and cleanup evidence](https://github.com/geoffrey-xiao/deprail/issues/345)
7. [#346 — Implement apply evidence and deterministic output](https://github.com/geoffrey-xiao/deprail/issues/346)
8. [#347 — Add security and cross-platform apply coverage](https://github.com/geoffrey-xiao/deprail/issues/347)

## Acceptance

- Valid approved plans apply only in the temporary worktree; caller tree status and digest remain unchanged.
- Dry-run validates and reports without mutation.
- Stale, unapproved, malformed, unsupported, path-escaping, timed-out, cancelled, or failed operations return explicit stable failures and cleanup evidence.
- Verification and rescan classify resolved, unchanged, introduced, and residual findings deterministically.
- Rollback/discard is idempotent and tested after success and failure paths.
- Hostile-path, symlink, shell-metacharacter, output-limit, credential-redaction, process-tree, and interrupted-cleanup cases are covered.
- Linux/macOS/Windows smoke evidence and independent architecture/security review are attached before closure.
