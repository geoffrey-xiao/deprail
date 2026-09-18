# v0.2 CLI Contract

## Commands

```text
deprail doctor
deprail discover [path]
deprail scan [path]
```

The command names and primary options remain v0.1-compatible.

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
| `0` | Successful execution; scan status is complete when applicable. |
| `2` | Invalid command, argument, option, or configuration. |
| `3` | Scanner failure, incomplete scope, or unusable scan result. |
| `1`, `4`, `5` | Reserved/documented meanings remain unchanged from v0.1. |

## Version identity

Version output MUST distinguish release builds from development builds and SHOULD include:

- semantic version or development marker;
- Git tag when available;
- source commit when available.

Version identity is metadata, not part of stable finding keys.
