# v0.2 CLI Contract

## Commands

```text
deprail doctor
deprail discover [path]
deprail scan [path]
deprail diff --base <input> --head <input>
deprail policy check --baseline <input> [path]
```

The v0.2 commands preserve v0.1 behavior and add deterministic baseline comparison and policy evaluation. Primary options remain versioned contracts.

## Argument behavior

- Missing optional path defaults to the documented repository context.
- Unexpected positional arguments fail with the configuration/argument error.
- Unsupported options fail before discovery, scanner execution, or writes.
- `--format json` sends only machine data to stdout.
- Diagnostics, progress, and guidance use stderr.
- `--output` writes atomically and does not overwrite unless explicitly permitted by the contract.

## Exit codes

| Code | Meaning |
| --- | --- |
| `0` | Successful execution; scan or policy status permits success. |
| `2` | Invalid command, argument, option, policy, or configuration. |
| `3` | Scanner failure, incomplete scope, or unusable scan result. |
| `4` | Policy violation or required guardrail block. |
| `1`, `5` | Reserved/documented meanings remain unchanged from v0.1. |

## Version identity

Version output MUST distinguish release builds from development builds and SHOULD include:

- semantic version or development marker;
- Git tag when available;
- source commit when available.

Version identity is metadata, not part of stable finding keys.
