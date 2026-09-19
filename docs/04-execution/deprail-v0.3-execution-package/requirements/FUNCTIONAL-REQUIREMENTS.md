# v0.3 Functional Requirements

- `FR-030-001`: `deprail fix plan` MUST require an explicit normalized scan-report input and a finding key; it MUST NOT infer a latest report or rescan implicitly.
- `FR-030-002`: The report MUST provide repository identity, source scan ID, artifact digests, workspace, finding key, and repository state.
- `FR-030-003`: Planning MUST resolve workspace, component, vulnerability aliases, current version, dependency path, and provenance.
- `FR-030-004`: Planning MUST be read-only and MUST NOT edit files, install packages, execute package-manager commands, or run scripts.
- `FR-030-005`: Candidates MUST distinguish recommended, viable, rejected, unavailable, and unknown states.
- `FR-030-006`: A recommendation MUST NOT be emitted when required compatibility evidence is unknown.
- `FR-030-007`: Plans MUST preserve direct/transitive ownership, constraints, lockfile effects, major-version risk, runtime/peer/engine risks, and assumptions.
- `FR-030-008`: Plans MUST include structured future commands with executable, arguments, and canonical repository-relative working directory.
- `FR-030-009`: Plan identity and serialization MUST be deterministic for equivalent inputs.
- `FR-030-010`: Plans MUST bind to source scan/report digest and reject stale inputs explicitly.
- `FR-030-011`: npm/pnpm/Yarn, Python requirements/uv/Poetry, and Maven/Gradle adapters MUST share one normalized planning contract.
- `FR-030-012`: JSON output MUST validate against the versioned remediation-plan schema.
- `FR-030-013`: Output paths inside the target repository, traversal paths, and symlink escapes MUST be rejected.
- `FR-030-014`: Missing metadata, malformed input, unsupported managers, and incomplete analysis MUST remain explicit and never become safe success.
- `FR-030-015`: Human output MUST identify recommendation, risks, affected files, and verification requirements.
- `FR-030-016`: Default planning MUST be offline/local-first unless a separately approved metadata capability is enabled.
