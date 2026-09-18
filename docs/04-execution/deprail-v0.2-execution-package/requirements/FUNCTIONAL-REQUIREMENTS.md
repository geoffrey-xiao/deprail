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
- `FR-020-011`: Baseline comparison MUST reject incompatible or incomplete trusted inputs.
- `FR-020-012`: Diff MUST classify new, resolved, unchanged, and dependency-change results deterministically.
- `FR-020-013`: Policy MUST NOT pass incomplete scans, scanner failures, or expired exceptions silently.
- `FR-020-014`: Policy decisions MUST expose stable CI exit behavior.
- `FR-020-015`: SARIF output MUST validate and preserve finding identity and provenance.
- `FR-020-016`: The GitHub Action MUST use documented least-privilege permissions and preserve CLI failure semantics.

## Deferred requirements

Remediation, remote publishing, hosted history, additional scanner families, and SBOM/signing implementation require separate approved contracts.
