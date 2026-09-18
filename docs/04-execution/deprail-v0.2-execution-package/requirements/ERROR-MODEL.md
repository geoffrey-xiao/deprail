# v0.2 Error Model

## Rules

Errors are observable contract data. Diagnostics MUST use stable codes, a safe message, scope, and actionable context where available. Wording may improve without changing the code.

Failures MUST NOT be converted into empty successful results.

## Required classes

| Code/class | Meaning | Result behavior |
| --- | --- | --- |
| `CONFIG_INVALID` | Invalid argument or configuration. | Exit `2`; no scan result. |
| `PATH_OUTSIDE_ROOT` | Traversal or symlink escape. | Security diagnostic; affected work is not trusted. |
| `SCANNER_NOT_FOUND` | Supported scanner unavailable. | Exit `3`; no false-safe result. |
| `SCANNER_TIMEOUT` | Scanner exceeded deadline. | Failed/partial according to workspace scope. |
| `SCANNER_OUTPUT_INVALID` | Output malformed or incompatible. | Failed/partial; retain safe raw evidence when possible. |
| `SCANNER_EXIT_NONZERO` | Unhandled scanner failure exit. | Failed/partial; vulnerability-found exit is interpreted by adapter contract. |
| `ARTIFACT_WRITE_FAILED` | Raw or report artifact could not be safely written. | Error; never claim complete persistence. |
| `DISCOVERY_INCOMPLETE` | Target scope cannot be fully determined. | Partial or failed; never safe. |

## Redaction

Messages and evidence MUST redact tokens, credential-bearing URLs, sensitive environment values, and unnecessary user information.
