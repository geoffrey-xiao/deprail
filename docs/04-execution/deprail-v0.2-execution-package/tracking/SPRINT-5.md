# v0.2 Sprint 5 — Baselines and Diff

## Goal

Produce deterministic base/head comparison and change classification.

## Planned items

- V02-015 baseline representation and storage.
- V02-016 base/head scan comparison.
- V02-017 new/resolved/unchanged classification.
- V02-018 deterministic diff output and `deprail diff`.

## Dependencies

Prerequisite hardening and stable finding keys must be complete. V02-016 depends on V02-015; V02-017 depends on V02-016; V02-018 depends on V02-017.

## Exit criteria

- Compatible baselines round-trip and validate.
- Real fixtures produce deterministic classifications.
- `deprail diff` has documented JSON, terminal, and failure behavior.
