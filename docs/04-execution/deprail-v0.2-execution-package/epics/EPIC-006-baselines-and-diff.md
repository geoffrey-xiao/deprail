# EPIC-006 Baselines and Diff

## Outcome

Turn scan history into a deterministic comparison that identifies dependency and vulnerability risk introduced by a change.

## Scope

Baseline representation and storage, base/head comparison, new/resolved/unchanged classification, deterministic diff output, and the `deprail diff` command.

## Out of Scope

Policy decisions, exceptions, SARIF rendering, repository mutation, remote publishing, and hosted history.

## Child issues

- V02-015 Baseline representation and storage.
- V02-016 Base/head scan comparison.
- V02-017 New/resolved/unchanged classification.
- V02-018 Deterministic diff output and `deprail diff`.

## Dependencies

Prerequisite hardening EPIC-001 through EPIC-004; especially stable scan JSON, stable keys, and artifact provenance.

## Acceptance

- Equivalent baselines compare deterministically.
- Added, resolved, unchanged, and affected dependency changes are distinguishable.
- Diff output is schema-valid and machine-output safe.
- `deprail diff` works without a server or network publishing service.
- Real fixture evidence is linked for the v0.2 gate.
