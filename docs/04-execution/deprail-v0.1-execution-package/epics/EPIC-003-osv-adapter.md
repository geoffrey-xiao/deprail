# EPIC-003 OSV Scanner Adapter

## Outcome

A bounded, version-aware, replayable OSV integration that never converts failure into empty success.

## Scope

The epic owns the contracts and observable behavior implied by its linked issues. Cross-epic changes require an explicit dependency and reviewer.

## Non-goals

Do not add later-roadmap features, weaken completeness, or introduce undocumented platform-specific behavior.

## Issues

- `ADAPTER-001`
- `PROC-001`
- `OSV-001`
- `OSV-002`
- `ART-001`
- `TEST-020`

## Key Decisions

Prefer deterministic, offline-testable domain behavior. Preserve provenance. Keep external tools behind versioned contracts and security boundaries.

## Acceptance

- [ ] Every linked issue meets its acceptance checklist.
- [ ] Required tests pass locally and in CI.
- [ ] Schemas, examples, and user documentation agree.
- [ ] Human review records risk and remaining limitations.
