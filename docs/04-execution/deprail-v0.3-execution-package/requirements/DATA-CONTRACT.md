# v0.3 Data Contract

## RemediationPlan

Required fields:

```text
schema_version
plan_id
created_from
repository_identity
workspace_identity
finding_identity
component
current_state
candidates
recommendation
affected_files
commands
risks
assumptions
verification
rollback
provenance
```

`created_from` includes source report/scan digest, scanner/tool identity, and repository state. `plan_id` excludes timestamps, prose, evidence order, and presentation formatting.

## Candidate states

```text
recommended | viable | rejected | unavailable | unknown
```

Unknown or unavailable evidence cannot be represented as safe compatibility. A valid plan may contain no recommendation when analysis is incomplete.

## Invariants

- Stable JSON is deterministic and schema-valid.
- Repository-relative paths use `/`.
- Commands are structured objects with executable, arguments, and canonical repository-relative `working_directory`; `.` represents repository root.
- `working_directory` MUST pass root containment and symlink checks and MUST NOT be absolute or traverse outside the repository.
- Findings, candidates, risks, and evidence are canonicalized before serialization.
- Provenance and raw source digests remain available.
- Plan output never embeds secrets or unnecessary source content.

Breaking changes require a new schema version, converter, examples, compatibility tests, and owner/security review.
