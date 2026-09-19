# DepRail v0alpha1 remediation-plan schema

`remediation-plan.schema.json` is the versioned JSON Schema for read-only v0.3 remediation plans.

The schema is additive to the existing `schemas/v1alpha` contracts. It requires:

- `schema_version: "v0alpha1"` and a 64-character lowercase hexadecimal `plan_id`;
- source report digest, scan ID, scanner identity, and repository state;
- repository, workspace, finding, and component identity;
- explicit candidate states: `recommended`, `viable`, `rejected`, `unavailable`, and `unknown`;
- structured commands with executable, argument arrays, and repository-relative `working_directory`;
- affected files, risks, assumptions, verification, rollback, and provenance.

Schema validation enforces required fields, version identity, structured commands, and basic relative-path shape. Runtime canonical-root and symlink containment remains an application responsibility.

Examples cover recommended, multiple-candidate, no-recommendation, unknown, rejected, stale, and incomplete plans. Plan identity rules are implemented and tested by `internal/remediation`.

Within the v0alpha1 schema line, additive fields are compatible. Breaking changes require a new schema version, converter, examples, compatibility tests, and owner/security review.
