# v0.2 Data Contract

## Serialization invariants

- Empty findings, errors, diagnostics, workspaces, and artifact collections serialize as `[]` where the schema defines a collection.
- Stable JSON is deterministic for equivalent input.
- Repository-relative paths use `/`.
- Stable keys exclude timestamps, descriptions, severity labels, and evidence ordering.
- Unknown safe fields may be preserved but must not become a safe result or confidence claim.
- Scanner/database/tool provenance and raw artifact digests remain available.

## Completeness

`complete` means every detected target succeeded. `partial` means at least one target succeeded while another was omitted, incomplete, or failed. `failed` means no trustworthy result or a core execution failure.

A zero-finding report is only a no-known-vulnerabilities result when status is `complete`.

## Compatibility

No breaking schema change is approved by the v0.2 planning package. A schema change requires an updated schema, examples, compatibility record, consumer migration note, and reviewed acceptance evidence.
