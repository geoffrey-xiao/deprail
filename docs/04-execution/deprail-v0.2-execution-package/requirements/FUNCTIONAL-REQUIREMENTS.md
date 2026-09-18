# v0.2 Functional Requirements

## Normative requirements

- `FR-020-001`: Empty report collections MUST serialize as JSON arrays.
- `FR-020-002`: Commands MUST reject unexpected positional arguments and unsupported options.
- `FR-020-003`: Scan scope MUST be derived from the canonical requested root, not the caller working directory.
- `FR-020-004`: Every detected target MUST have an explicit completeness outcome.
- `FR-020-005`: Scanner failure, timeout, malformed output, missing tool, and incompatible version MUST NOT become an empty successful scan.
- `FR-020-006`: Raw scanner output MUST remain retained according to the artifact contract and include provenance.
- `FR-020-007`: Normalization MUST be order-independent for equivalent inputs.
- `FR-020-008`: Machine output MUST remain free of progress and diagnostic text.
- `FR-020-009`: Release builds MUST expose truthful version/tag/commit identity.
- `FR-020-010`: Discovery MUST remain read-only and offline.

## Deferred requirements

Baseline comparison, policy gates, SARIF, remediation, remote publishing, and additional scanner families require separate approved contracts.
