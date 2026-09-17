# v0.1 Functional Requirements

## Discovery

- FR-DISC-001 detect supported manifests and lockfiles.
- FR-DISC-002 group files into stable workspaces.
- FR-DISC-003 apply ignore rules without escaping the root.
- FR-DISC-004 continue after one malformed project.
- FR-DISC-005 diagnose conflicts and missing authoritative lockfiles.
- FR-DISC-006 produce deterministic ProjectGraph output.

## Scanner

- FR-SCAN-001 check adapter availability and compatibility.
- FR-SCAN-002 create an explicit ScanPlan.
- FR-SCAN-003 execute with arguments, deadline, cancellation, and limits.
- FR-SCAN-004 retain tool metadata and raw artifacts.
- FR-SCAN-005 classify non-zero exit and malformed output.

## Normalization

- FR-NORM-001 canonicalize components with PURL.
- FR-NORM-002 merge vulnerability aliases deterministically.
- FR-NORM-003 preserve evidence, severity sources, and fixed-version provenance.
- FR-NORM-004 calculate stable keys and ordering.

## CLI and Output

- FR-CLI-001 implement discover, scan, and doctor.
- FR-CLI-002 separate data and diagnostics.
- FR-CLI-003 validate versioned JSON.
- FR-CLI-004 write files atomically and safely.

## Completeness and Errors

- FR-ERR-001 expose complete, partial, and failed.
- FR-ERR-002 never describe incomplete scans as safe.
- FR-ERR-003 use stable error codes and exit semantics.

## Release and Quality

- FR-REL-001 build and smoke-test supported platforms.
- FR-REL-002 publish checksums, an SBOM, and signing status.
- FR-REL-003 validate at least three real repositories manually.
