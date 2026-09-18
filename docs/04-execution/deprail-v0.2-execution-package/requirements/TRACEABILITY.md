# v0.2 Traceability Matrix

| Requirement | Planned evidence | Owner/reviewer decision point |
| --- | --- | --- |
| Empty collections serialize as arrays | JSON contract test and schema example | PR review |
| Unexpected arguments fail | CLI contract test and exit-code assertion | PR review |
| OSV v2 behavior is reproducible | Real raw fixtures and exit-code matrix | Adapter review |
| Requested-root scanning is safe | Outside-root end-to-end test and artifact-path evidence | Security review |
| Version identity is truthful | Release build and doctor JSON smoke output | Release review |
| Deterministic normalization remains stable | Golden/property tests | Normalization review |
| Cross-platform meaning remains equivalent | Linux/macOS/Windows CI and artifact smoke | Release review |
| Release gaps are explicit | Release evidence record and checklist links | Owner review |

Every implementation issue must map to at least one row or explicitly record why it is documentation/process-only.
