# v0.3.1 CLI Contract

## Commands

The existing command set is unchanged:

```text
deprail doctor [path]
deprail discover [path]
deprail scan [path]
deprail diff --base <ref> --head <ref>
deprail fix plan --report <scan-report> --finding <finding-key>
```

## Help behavior

- `deprail --help` and `deprail -h` print deterministic root usage and exit `0`.
- Every supported command and subcommand accepts `--help` and `-h`, prints usage to stdout, and exits `0`.
- Help output is plain, readable in TTY and non-TTY contexts, and contains no progress or ANSI control sequences.
- Help performs no discovery, scanner, package-manager, network, repository mutation, or artifact writes.

## Stream behavior

- Human output may use stdout according to the existing command contract.
- `--format json` writes machine data only to stdout.
- Diagnostics and progress use stderr.
- TTY progress is enabled only when stderr is interactive and the mode is not quiet.
- Redirected and CI output contains no animation or ANSI control sequences.
- The current scan parser rejects `--quiet`; v0.3.1 must implement the documented quiet mode rather than treat that rejection as compatibility.
- `--quiet` suppresses successful progress and summaries while retaining required errors.
- `--verbose` adds safe diagnostics without exposing secrets or changing result meaning.

## Human message lifecycle

Long-running commands use the following truthful sequence when the relevant boundary exists:

1. Start notice naming the operation and target, such as `Scanning <path>...` or `Discovering workspaces...`.
2. Real-work progress for discovered workspaces, scan targets, artifacts, or normalization; never invented percentages, durations, or activity.
3. Finding summary using only trustworthy counts and severities; unknown severity remains unknown.
4. One terminal outcome: `Complete`, `Partial`, `Failed`, or `Cancelled`.
5. A bounded next action when safe, such as a JSON output path, affected workspace, stable finding key, or read-only plan command.

For `--format json`, suppress start notices, progress, summaries, and guidance on both streams; retain only required diagnostics on stderr. JSON stdout contains only the versioned report.

| Outcome | Required meaning |
| --- | --- |
| `Complete` | Work finished; zero findings may say `No known vulnerabilities found.` |
| `Partial` | Some scope was omitted, failed, or incomplete; identify affected scope where known. |
| `Failed` | No trustworthy result or core execution failed; include stable error code and safe reason. |
| `Cancelled` | User/context cancellation stopped work; never claim report readiness or success. |

`--quiet` suppresses progress and successful summaries but retains required errors; an interactive one-shot start notice may remain for human terminal mode. `--verbose` adds safe diagnostics and evidence without changing result meaning.

Policy evaluation is separate from execution completeness. `policy check` must preserve and display its `pass`, `warn`, or `block` decision; `block` is not a failed execution, and `Complete` must not replace the policy decision.

## Outcome behavior

- `complete` with zero findings may render `No known vulnerabilities found.`
- `partial` must state that the scan is incomplete and identify affected scope safely.
- `failed` must state failure and expose the stable error code where applicable.
- `cancelled` must state cancellation and restore terminal state.

## Compatibility

The release does not add a full-screen mode, change exit-code meanings, infer missing inputs, mutate repositories, invoke package managers, or alter scanner arguments. Any exception requires an approved contract update before implementation.
