# DepRail v1alpha schemas

`deprail.schema.json` is the versioned JSON Schema for v1alpha project, discovery, scan, finding, provenance, and diagnostic documents.

Every document carries:

- `schema_version: "v1alpha"`
- `document_type`
- deterministic identifiers and declared ordering fields
- explicit `complete`, `partial`, or `failed` status where applicable
- provenance and raw-artifact digests for scanner evidence

Path schemas reject absolute paths. Repository-root containment, traversal, and symlink-escape checks remain runtime discovery responsibilities and are not delegated to regular-expression validation.

Unknown fields are allowed so adapters can preserve safe upstream data while the contract evolves additively. Stable keys must not include descriptions, timestamps, severity labels, or evidence ordering.

Examples live under `schemas/v1alpha/examples/`. `scan.json`, `partial-scan.json`, and `failed-scan.json` demonstrate complete, partial, and failed outcomes with array-valued collections. The `invalid-scan.json` fixture intentionally violates the schema and is used for negative validation checks.
