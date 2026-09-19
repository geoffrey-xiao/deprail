# v0.3 CLI Contract

## Commands

```text
deprail fix plan --report <scan-report> --finding <finding-key>
deprail fix plan --report <scan-report> --finding <finding-key> --format json
deprail fix plan --report <scan-report> --finding <finding-key> --output <path-outside-repository>
```

The report input is mandatory. v0.3 does not infer a latest report, rescan implicitly, or search an unindexed artifact directory. The report must carry or reference repository identity, source scan ID, artifact digests, workspace, finding key, and repository state.

## Behavior

- `--report` must resolve to a readable, schema-valid normalized scan report.
- `--finding` must resolve to exactly one finding in that report.
- `--format json` writes only machine data to stdout.
- Diagnostics and guidance use stderr.
- `--output` is optional, writes atomically, and MUST resolve outside the target repository root.
- Output paths inside the repository, traversal paths, and symlink escapes fail before planning.
- Planning does not mutate the repository, invoke package managers, or execute scripts.
- Equivalent input produces equivalent plan JSON independent of finding/evidence order.

## Exit behavior

Existing v0.2 numeric meanings remain unchanged. The implementation must use the existing configuration/input error for missing or ambiguous findings and the existing incomplete/failure semantics for unusable analysis; any new code requires an error-model update before implementation.

## Human output

Show finding, current component, recommended candidate or no-recommendation reason, affected files, migration risks, future commands, verification requirements, provenance, and a clear read-only notice.
