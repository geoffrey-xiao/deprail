# v0.3 Traceability Matrix

| Requirement | Planned evidence | Review |
| --- | --- | --- |
| Finding resolution and provenance | Source-report fixtures and plan examples | Domain review |
| Direct/transitive ownership | npm/Python/Java adapter contract tests | Remediation review |
| Candidate compatibility and risks | Constraint/major/peer/runtime golden cases | Security/architecture review |
| No repository mutation | Before/after tree, process audit, mutation-attempt tests | Security review |
| Stable plan identity | Reordered-input property tests and golden JSON | Compatibility review |
| Unknown/incomplete handling | Error fixtures and no-recommendation cases | Safety review |
| Schema-valid machine output | JSON Schema validation and compatibility tests | Contract review |
| CLI behavior | Argument, stdout/stderr, output-file, exit tests | CLI review |
| Cross-platform equivalence | Linux/macOS/Windows fixture smoke | Release review |
| v0.4 handoff | Plan consumed by a fixture executor design without reinterpretation | Architecture review |

Every future implementation issue MUST map to one or more rows or explain why it is documentation/process-only.
