# V04-007: Generate Trusted Baselines from Scan Output

**Epic:** Release parent #299; extends the existing verification/rescan contract
**GitHub issue:** [#329](https://github.com/geoffrey-xiao/deprail/issues/329)
**ADR:** [ADR-0003](../../../adr/ADR-0003-baseline-generation.md)
**Status:** Proposed; scope-change review required

## Value

A user can produce the baseline document required by `deprail diff` and `deprail policy check` from an actual complete scan result. Manual baseline JSON authoring is not required.

## Scope

Implement the read-only command:

```text
deprail baseline create --scan scan.json --output baseline.json [--format json]
```

The command validates the existing versioned scan contract, converts it to the existing `v1alpha` baseline contract, preserves scan identity, artifact digests, and stable finding keys, sorts contract lists deterministically, and writes the result atomically with restrictive permissions. Existing output is not overwritten implicitly.

## Failure behavior

Malformed, unsupported, stale, partial, failed, or otherwise incomplete scan inputs fail with stable diagnostics and the documented non-success exit status. No invalid input may be persisted as a trusted baseline. Input and output path traversal, symlink escape, output-limit, permission, and interrupted-write failures remain explicit. The command performs no network access, installation, repository mutation, or automatic replacement.

## Acceptance

- A complete scan fixture converts to schema-valid baseline JSON.
- Repeated conversion of equivalent input produces byte-stable output and stable baseline identity.
- `deprail diff` and `deprail policy check` consume the generated baseline without manual edits.
- Invalid and incomplete scan fixtures fail and leave no output artifact.
- Existing output is preserved; writes are atomic and restrictive.
- Terminal and JSON diagnostics follow the CLI contract.
- Linux/macOS/Windows path and serialized-domain semantics are equivalent.

## Dependencies

- Existing baseline schema, validator, comparison, and policy contracts.
- Existing scan JSON contract and artifact-digest behavior.
- Existing process/path/output safety helpers.

## Exclusions

No policy-language expansion, new scanner family, remote publishing, automatic replacement, repository mutation, or release publication.

## Evidence

Contract tests, schema validation, deterministic-output fixture, failure fixtures, CLI smoke output, and cross-platform manual evidence linked from the issue and release evidence record.
