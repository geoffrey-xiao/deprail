# EPIC-001 Engineering Foundation and Contracts

## Outcome

A reproducible repository with pinned tools, versioned contracts, CI, fixtures, and contribution controls.

## Scope

The epic owns the contracts and observable behavior implied by its linked issues. Cross-epic changes require an explicit dependency and reviewer.

## Non-goals

Do not add later-roadmap features, weaken completeness, or introduce undocumented platform-specific behavior.

## Issues

- `S0-001`
- `S0-002`
- `S0-003`
- `SCHEMA-001`
- `FIXTURE-001`

## Key Decisions

Prefer deterministic, offline-testable domain behavior. Preserve provenance. Keep external tools behind versioned contracts and security boundaries.

## Acceptance

- [ ] Every linked issue meets its acceptance checklist.
- [ ] Required tests pass locally and in CI.
- [ ] Schemas, examples, and user documentation agree.
- [ ] Human review records risk and remaining limitations.
