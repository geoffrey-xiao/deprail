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


| Baseline and diff | V02-015 through V02-018; schema, comparison, classification, and CLI evidence | Normalization and compatibility review |
| Policy and exceptions | V02-019 through V02-021; policy fixtures, expiry cases, and exit matrix | Security and CI compatibility review |
| SARIF output | V02-022; SARIF schema validation and consumer smoke | Standards and output review |
| GitHub Action | V02-023 through V02-025; permissions, caching, and real pull-request evidence | Security and release review |
| Public preview | V02-026; artifact, cross-platform, smoke, rollback, and remaining-risk evidence | Owner and release review |
Every implementation issue must map to at least one row or explicitly record why it is documentation/process-only.
