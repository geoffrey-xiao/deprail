# v0.1 Error Model

## Structure

```json
{"code":"SCANNER_NOT_FOUND","message":"OSV-Scanner is not available","scope":"repository","retryable":false,"details":{},"help":"Run deprail doctor"}
```

Errors contain a stable code, concise message, affected scope, retryability, safe structured details, and optional guidance. Internal stack traces appear only in debug output.

## Stable Error Codes

`CONFIG_INVALID`, `PATH_OUTSIDE_ROOT`, `MANIFEST_INVALID`, `LOCKFILE_CONFLICT`, `DISCOVERY_INCOMPLETE`, `SCANNER_NOT_FOUND`, `SCANNER_VERSION_UNSUPPORTED`, `SCANNER_TIMEOUT`, `SCANNER_EXIT_NONZERO`, `SCANNER_OUTPUT_LIMIT`, `SCANNER_OUTPUT_INVALID`, `ARTIFACT_WRITE_FAILED`, `NORMALIZATION_FAILED`, and `OUTPUT_WRITE_FAILED`.

## Rules

Errors never contain credentials. Partial success retains successful workspace results. Status and exit code reflect the most severe incomplete condition. Human messages may improve without changing stable codes.
