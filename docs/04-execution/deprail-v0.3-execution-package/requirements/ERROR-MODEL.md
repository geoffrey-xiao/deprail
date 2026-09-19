# v0.3 Error Model

Errors remain structured with stable code, safe message, scope, and actionable context. Failures MUST NOT become empty successful plans or safe recommendations.

| Code/class | Meaning | Result behavior |
| --- | --- | --- |
| `CONFIG_INVALID` | Invalid command, option, or plan configuration | Exit `2`; no plan |
| `PLAN_REPORT_REQUIRED` | Explicit normalized scan report was not supplied | Input error; no plan |
| `PLAN_REPORT_INVALID` | Report is malformed, incompatible, or lacks required provenance | Input error; no plan |
| `FINDING_NOT_FOUND` | Finding/key absent from source report | Input error; no plan |
| `FINDING_AMBIGUOUS` | Input resolves to multiple findings | Input error; no plan |
| `PLAN_INPUT_STALE` | Source report or repository state changed | Reject plan; request fresh scan |
| `PLAN_INPUT_INVALID` | Malformed or incompatible source data | Failed/incomplete; preserve diagnostics |
| `REMEDIATION_UNSUPPORTED` | Ecosystem or manager is not supported | Explicit unsupported state; no false-safe recommendation |
| `CANDIDATE_METADATA_UNAVAILABLE` | Required package metadata unavailable | Candidate unavailable/unknown |
| `PLAN_ANALYSIS_INCOMPLETE` | Analysis could not establish all required facts | Incomplete plan; no recommendation where unsafe |
| `PLAN_OUTPUT_OUTSIDE_ROOT_REQUIRED` | Output path is inside target root or escapes allowed boundary | Reject before planning or writing |
| `PLAN_SCHEMA_INVALID` | Generated plan violates its schema | Error; do not write output |
| `PLAN_WRITE_FAILED` | Safe output/persistence failed | Error; do not claim saved plan |
| `PATH_OUTSIDE_ROOT` | Traversal or symlink escape | Security error; no plan |
| `PLAN_MUTATION_ATTEMPT` | A planning path attempted a write or execution | Security error; terminate operation |
