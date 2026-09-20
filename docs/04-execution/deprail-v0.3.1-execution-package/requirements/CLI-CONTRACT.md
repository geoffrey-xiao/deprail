# v0.3.1 CLI Contract

## Commands

The existing command set is unchanged:

```text
deprail doctor [path]
deprail discover [path]
deprail scan [path]
deprail fix plan --report <scan-report> --finding <finding-key>
```

## Stream behavior

- Human output may use stdout according to the existing command contract.
- `--format json` writes machine data only to stdout.
- Diagnostics and progress use stderr.
- TTY progress is enabled only when stderr is interactive and the mode is not quiet.
- Redirected and CI output contains no animation or ANSI control sequences.
- `--quiet` suppresses successful progress and summaries while retaining required errors.
- `--verbose` adds safe diagnostics without exposing secrets or changing result meaning.

## Outcome behavior

- `complete` with zero findings may render `No known vulnerabilities found.`
- `partial` must state that the scan is incomplete and identify affected scope safely.
- `failed` must state failure and expose the stable error code where applicable.
- `cancelled` must state cancellation and restore terminal state.

## Compatibility

The release does not add a full-screen mode, change exit-code meanings, infer missing inputs, mutate repositories, invoke package managers, or alter scanner arguments. Any exception requires an approved contract update before implementation.
