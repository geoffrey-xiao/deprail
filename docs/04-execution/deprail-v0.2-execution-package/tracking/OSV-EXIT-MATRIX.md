# OSV-Scanner v2 Exit Matrix

Supported contract: OSV-Scanner v2, currently verified with 2.6.0.

| Condition | Process result | Adapter result | Safe success? |
| --- | --- | --- | --- |
| No findings, valid JSON | exit 0 | parse empty findings | Yes |
| Findings, valid JSON | exit 1 | accept output, parse findings | Yes, with findings |
| Findings exit with malformed JSON | exit 1 | `SCANNER_OUTPUT_INVALID` | No |
| Unexpected non-zero exit | exit other than 1 | `SCANNER_EXIT_NONZERO` | No |
| Missing executable | process not found | `SCANNER_NOT_FOUND` | No |
| Timeout | deadline exceeded | `SCANNER_TIMEOUT` | No |
| Cancellation | context canceled | execution failure | No |
| Output limit exceeded | capture limit | `SCANNER_OUTPUT_LIMIT` | No |
| Unsupported scanner version | valid version output | `SCANNER_VERSION_UNSUPPORTED` | No |

The adapter accepts exit 1 only as the scanner's vulnerability-found signal. The captured stdout must still parse as valid OSV JSON; otherwise the scan fails explicitly. No retry or automatic installation is performed.

Evidence is implemented in `internal/adapters/osv/failure_matrix_test.go`, with controlled child-process modes for zero findings, valid findings, malformed output, unexpected exit, timeout, cancellation, missing executable, and output limits. Version compatibility is covered by `internal/adapters/osv/metadata_test.go`.
