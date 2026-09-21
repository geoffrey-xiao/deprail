# v0.4 Functional Requirements

| ID | Requirement | Acceptance evidence |
|---|---|---|
| FR-001 | Apply accepts only a valid versioned plan and explicit approval bound to source identity. | Invalid, expired, reused, or mismatched approval fails without mutation. |
| FR-002 | Dry-run performs validation and planning but no dependency or source mutation. | Source and workspace digests are unchanged. |
| FR-003 | Mutation occurs only inside the isolated workspace and only through a supported adapter. | Hostile path and arbitrary-command fixtures fail closed. |
| FR-004 | Every subprocess has direct argv, explicit cwd/env, deadline, cancellation, and bounded output. | Timeout, cancellation, non-zero, and output-limit evidence is emitted. |
| FR-005 | Verification runs only approved, discovered commands and reports incomplete verification explicitly. | Unsupported or missing verification does not become success. |
| FR-006 | Before/after scans preserve scanner identity and classify resolved, residual, introduced, and unknown findings deterministically. | Stable transition fixture and schema validation pass. |
| FR-007 | Failure, cancellation, or cleanup error retains evidence and never changes the caller repository. | Source-tree digest and cleanup evidence are recorded. |
| FR-008 | Machine output is versioned JSON on stdout; diagnostics and progress are on stderr. | JSON mode emits no banners or progress. |
| FR-009 | Evidence is redaction-safe and content-addressed; secrets and credential-bearing URLs are not persisted. | Redaction and artifact digest fixtures pass. |
| FR-010 | Supported behavior is equivalent on Linux, macOS, and Windows for serialized domain meaning. | Cross-platform matrix and manual evidence are complete. |
