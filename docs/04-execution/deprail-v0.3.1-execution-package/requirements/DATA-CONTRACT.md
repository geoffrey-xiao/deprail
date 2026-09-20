# v0.3.1 Data Contract

## Presentation events

Events are internal application data, not a serialized public schema. Each event must have:

- a stable event name;
- an explicit lifecycle position;
- deterministic identifiers and counts;
- no terminal escape sequences;
- no secrets or unrestricted environment values;
- enough information for human progress without requiring renderer inference.

Initial events are `InputValidated`, `WorkspaceDiscoveryStarted`, `WorkspaceDiscovered`, `ScanPlanBuilt`, `WorkspaceScanStarted`, `WorkspaceScanCompleted`, `ArtifactStored`, `NormalizationStarted`, `ReportReady`, and `OperationCancelled`.

## Machine output

Existing versioned JSON schemas, stable keys, provenance, completeness states, diagnostics, and ordering remain authoritative. v0.3.1 does not require a schema version change.

If implementation discovers a necessary additive schema change, it must include examples, compatibility tests, a decision record, and explicit review before merge.

## Human output

Human presentation is not a source of truth for domain data. It may summarize only values present in structured results and events. Labels are escaped before rendering.
