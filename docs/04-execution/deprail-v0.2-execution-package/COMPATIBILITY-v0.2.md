# DepRail v0.2 Compatibility and Migration Record

## Policy

v0.2 is intended to be backward-compatible for the existing v0.1 command surface and machine-output meaning. Every contract change must be classified before implementation.

| Classification | Meaning | Required approval |
| --- | --- | --- |
| Unchanged | Existing behavior remains normative. | Issue review. |
| Clarified | Ambiguous behavior is made explicit without changing intended meaning. | PR review and evidence. |
| Extended | New optional behavior or metadata is added. | Compatibility review. |
| Breaking | Existing invocation or consumer behavior changes. | Explicit owner decision and migration record. |
| Deferred | Candidate behavior is not implemented in v0.2. | Backlog or decision link. |
| Accepted limitation | Known gap remains with recorded risk. | Owner acceptance during PR/release review. |

## Contract matrix

| Contract | v0.1 baseline | v0.2 position | Classification |
| --- | --- | --- | --- |
| Commands | `doctor`, `discover`, `scan` | Preserve names and core flows. | Unchanged |
| Invalid arguments | Some extra arguments were ignored. | Reject unexpected arguments with documented configuration error. | Clarified |
| Exit codes | `0` success, `2` config/argument, `3` scanner/incomplete. | Preserve meanings. | Unchanged |
| Empty collections | Some empty values could serialize as `null`. | Serialize contract collections as `[]`. | Clarified |
| Completeness | `complete`, `partial`, `failed`. | Preserve definitions and unsafe-result rules. | Unchanged |
| Paths | Repository-relative `/` paths with containment. | Preserve and test requested-root execution. | Unchanged |
| OSV-Scanner | Supported v2 range and v2 command/output. | Add real fixtures and exit-code coverage. | Extended verification |
| Version output | Build identity was not fully tag-derived. | Include tag/commit identity without false stable claims. | Extended |
| Schemas | Current v1alpha contracts. | No breaking schema change approved by this package. | Unchanged |
| Artifacts | Content-addressed raw results. | Preserve provenance and requested-root placement. | Unchanged |
| Configuration | `.deprail.yaml` and `.deprail/`. | No layout change approved. | Unchanged |
| Release modes | Preview procedure with explicit gaps. | Add preview/RC/stable mode checklist. | Extended process |

## Consumer migration

Consumers should:

- Treat empty collections as arrays.
- Treat unknown fields as forward-compatible where schema rules permit.
- Check completeness before describing a scan as safe or vulnerability-free.
- Use stable error codes rather than diagnostic wording.
- Avoid depending on descriptions, timestamps, severity labels, or evidence ordering for stable identity.
- Record the reported version and commit when comparing reports across releases.

## Breaking-change gate

A breaking change requires an ADR or compatibility decision containing affected consumers, old/new examples, migration instructions, schema or CLI evidence, rollback behavior, and owner approval during PR review.
