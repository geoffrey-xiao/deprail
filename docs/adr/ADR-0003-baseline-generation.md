# ADR-0003: Generate Trusted Baselines from Scan Results

**Status:** Proposed for v0.4.0-preview.1 scope review
**Date:** 2026-09-21
**Related issue:** [#329](https://github.com/geoffrey-xiao/deprail/issues/329)
**Release parent:** [#299](https://github.com/geoffrey-xiao/deprail/issues/299)

## Context

The repository has a validated `internal/baseline` contract and a `deprail diff` command, but no user-facing way to create a baseline from a real scan result. Requiring users to author baseline JSON manually is not an acceptable workflow. The v0.2 planning package intended baseline representation/storage and deterministic diff, while v0.4 currently depends on before/after scan identity for remediation evidence.

## Decision

Add a bounded `deprail baseline create` command in v0.4:

```text
deprail baseline create --scan scan.json --output baseline.json [--format json]
```

The command converts only a complete, valid scan result into the existing `v1alpha` baseline contract. It must preserve stable finding identity, artifact digests, deterministic ordering, and source scan identity. Partial, failed, malformed, stale, or unsupported scans fail explicitly and cannot become trusted baselines.

Baseline output uses the existing safe atomic, restrictive, non-overwriting persistence boundary. No automatic replacement, repository mutation, remote publishing, or policy-language expansion is included.

## Alternatives

- Manual JSON authoring: rejected; unsafe and not a real user workflow.
- Implicit baseline generation inside `diff`: rejected; hides provenance and makes persisted identity unclear.
- Defer to v0.5: rejected for this release because v0.4 release evidence and real before/after workflows require a supported baseline producer.

## Consequences

- v0.4 scope expands by one CLI capability and schema/fixture/test surface.
- Existing `diff` and `policy check` gain an end-to-end real-input workflow.
- The release gate gains baseline conversion, stale/incomplete failure, deterministic output, and cross-platform CLI evidence.
- Human review is required for the CLI and compatibility contract before implementation starts.
