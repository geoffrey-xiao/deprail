# v0.4 Failure and Data Contract

## Outcome states

`dry_run`, `applied`, `verified`, `partial`, `failed`, `cancelled`, and `cleanup_failed` are mutually exclusive terminal outcomes. A zero-finding post-scan is meaningful only when the post-scan is complete.

## Required evidence

Every apply result records schema version, plan digest, source commit, canonical source root identity, authorized paths, workspace identity, operation records, command/tool identity, bounded stdout/stderr digests, before/after repository and finding digests, verification status, rescan status, rollback/cleanup status, and stable diagnostics.

## Error classes

Use stable codes for invalid or expired approval, source mismatch, path escape, scanner/tool missing, scanner timeout, scanner output invalid, process timeout, process cancellation, output limit, verification incomplete, mutation failed, rescan failed, artifact failure, and cleanup failure. Errors MUST retain successful evidence and MUST NOT be converted to empty success.

## Redaction

Tokens, credential-bearing URLs, authorization headers, secret environment values, and user data embedded in credentials are redacted before evidence persistence. Unknown data is preserved only when safe and never promoted to a safe result.
