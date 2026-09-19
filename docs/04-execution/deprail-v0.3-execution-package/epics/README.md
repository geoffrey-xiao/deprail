# v0.3 Epics

These epics are a proposed decomposition of the v0.3 development plan. They are local execution contracts; GitHub epics and issues are created only after issue #198 completes Definition of Ready.

| Epic | GitHub issue | Outcome | Child issues |
| --- | --- | --- | --- |
| EPIC-001 | [#202](https://github.com/geoffrey-xiao/deprail/issues/202) | Deterministic remediation-plan model and schema | V03-001, V03-002 |
| EPIC-002 | [#203](https://github.com/geoffrey-xiao/deprail/issues/203) | Explicit report provenance and safe input/output boundaries | V03-003, V03-004 |
| EPIC-003 | [#204](https://github.com/geoffrey-xiao/deprail/issues/204) | Read-only package-manager planning adapters | V03-005, V03-006, V03-007 |
| EPIC-004 | [#205](https://github.com/geoffrey-xiao/deprail/issues/205) | CLI, presenters, persistence, and release evidence | V03-008, V03-009, V03-010 |

Dependencies flow from EPIC-001 to EPIC-004. No issue authorizes repository mutation or package-manager execution.
