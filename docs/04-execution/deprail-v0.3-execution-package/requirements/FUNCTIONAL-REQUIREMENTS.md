# v0.3 Functional Requirements

- `FR-030-001`: `deprail fix plan` MUST accept a finding identifier or stable finding key from a trusted scan/report.
- `FR-030-002`: Planning MUST resolve workspace, component, vulnerability aliases, current version, dependency path, and provenance.
- `FR-030-003`: Planning MUST be read-only and MUST NOT edit files, install packages, execute package-manager commands, or run scripts.
- `FR-030-004`: Candidates MUST distinguish recommended, viable, rejected, unavailable, and unknown states.
- `FR-030-005`: A recommendation MUST NOT be emitted when required compatibility evidence is unknown.
- `FR-030-006`: Plans MUST preserve direct/transitive ownership, constraints, lockfile effects, major-version risk, runtime/peer/engine risks, and assumptions.
- `FR-030-007`: Plans MUST include structured future commands and affected files without shell interpolation.
- `FR-030-008`: Plan identity and serialization MUST be deterministic for equivalent inputs.
- `FR-030-009`: Plans MUST bind to source scan/report digest and reject stale inputs explicitly.
- `FR-030-010`: npm/pnpm/Yarn, Python requirements/uv/Poetry, and Maven/Gradle adapters MUST share one normalized planning contract.
- `FR-030-011`: JSON output MUST validate against the versioned remediation-plan schema.
- `FR-030-012`: Missing metadata, malformed input, unsupported managers, and incomplete analysis MUST remain explicit and never become safe success.
- `FR-030-013`: Human output MUST identify recommendation, risks, affected files, and verification requirements.
- `FR-030-014`: Default planning MUST be offline/local-first unless a separately approved metadata capability is enabled.
