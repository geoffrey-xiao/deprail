# v0.1 CLI Contract

## Global Rules

- Machine output is stable and versioned.
- Data goes to stdout; diagnostics go to stderr.
- Color is enabled only for a TTY.
- Paths in JSON are repository-relative and slash-separated.
- `--format json` never includes banners or progress text.
- Cancellation and timeout produce explicit errors.

## Commands

### `deprail discover [path]`

Options: `--format terminal|json`, `--config`, `--no-ignore`, and `--verbose`. It performs local read-only discovery and returns project graph, completeness, warnings, and diagnostics.

### `deprail scan [path]`

Options: `--format terminal|json`, `--output`, `--adapter`, `--timeout`, `--verbose`, and `--quiet`. It discovers, plans, executes, stores raw output, normalizes, and reports.

### `deprail doctor`

Reports DepRail, operating system, configuration, adapter availability, tool version, compatibility, and safe remediation guidance.

## Exit Codes

| Code | Meaning |
| --- | --- |
| 0 | Execution succeeded and policy passed |
| 1 | Execution succeeded but policy blocked |
| 2 | Configuration or argument error |
| 3 | Scanner failed or result is incomplete |
| 4 | Remediation or verification failed; reserved in v0.1 |
| 5 | Approval or permission required; reserved in v0.1 |

## Output Safety

Redact tokens, user info in credential URLs, and sensitive environment values. File output is atomic and does not overwrite unless the command contract explicitly allows it.
