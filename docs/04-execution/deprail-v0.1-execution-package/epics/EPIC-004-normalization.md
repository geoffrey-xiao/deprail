# EPIC-004 Normalization and Evidence

## Outcome

Stable components, aliases, evidence, fixed versions, keys, and ordering independent of source order.

## Scope

The epic owns the contracts and observable behavior implied by its linked issues. Cross-epic changes require an explicit dependency and reviewer.

## Non-goals

Do not add later-roadmap features, weaken completeness, or introduce undocumented platform-specific behavior.

## Issues

- `NORM-001`
- `NORM-002`
- `NORM-003`
- `NORM-004`
- `TEST-030`

## Key Decisions

Prefer deterministic, offline-testable domain behavior. Preserve provenance. Keep external tools behind versioned contracts and security boundaries.

## Acceptance

- [ ] Every linked issue meets its acceptance checklist.
- [ ] Required tests pass locally and in CI.
- [ ] Schemas, examples, and user documentation agree.
- [ ] Human review records risk and remaining limitations.
