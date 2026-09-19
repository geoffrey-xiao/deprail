# v0.3 CLI Contract

## Commands

```text
deprail fix plan <finding-or-key>
deprail fix plan <finding-or-key> --format json
deprail fix plan <finding-or-key> --output <new-file>
```

## Behavior

- Finding input is required and must resolve unambiguously.
- `--format json` writes only machine data to stdout.
- Diagnostics and guidance use stderr.
- `--output` writes atomically and does not overwrite unless explicitly permitted.
- Planning does not mutate the repository, invoke package managers, or execute scripts.
- Equivalent input produces equivalent plan JSON independent of finding/evidence order.

## Exit behavior

Existing v0.2 numeric meanings remain unchanged. The implementation must use the existing configuration/input error for missing or ambiguous findings and the existing incomplete/failure semantics for unusable analysis; any new code requires an error-model update before implementation.

## Human output

Show finding, current component, recommended candidate or no-recommendation reason, affected files, migration risks, future commands, verification requirements, provenance, and a clear read-only notice.
