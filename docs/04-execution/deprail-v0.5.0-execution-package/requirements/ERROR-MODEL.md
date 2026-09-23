# v0.5 Error Model

**Status:** Draft proposal; code strings and transport mappings require owner and independent architecture/security review.

## 1. Contract

Errors are structured with stable code, safe user-facing message, scope, and actionable context. HTTP status codes express transport-level outcome; they do not replace application error codes. A failed history/API operation never becomes an empty successful result, and an API error must not rewrite the underlying scan outcome.

This release must preserve existing scan/remediation error codes from the applicable v0.3/v0.4 contracts, including scanner and path errors such as `SCANNER_NOT_FOUND`, `SCANNER_TIMEOUT`, `SCANNER_OUTPUT_INVALID`, and `PATH_OUTSIDE_ROOT`. The table below lists candidate v0.5 categories only; no string is stable until crosswalked with the repository-wide [`ERROR-MODEL.md`](../../deprail-v0.3-execution-package/requirements/ERROR-MODEL.md) and approved.

## 2. Candidate history/API errors

| Candidate code | Meaning | Required behavior |
| --- | --- | --- |
| `HISTORY_UNAVAILABLE` | Local store cannot be opened/read or access is denied | Report store unavailable; do not return empty history or weaken permissions. |
| `HISTORY_WRITE_FAILED` | A history transaction could not be durably committed | Do not claim saved; retain scan outcome separately and report persistence failure. |
| `HISTORY_SCHEMA_UNSUPPORTED` | Store uses an unsupported newer/older schema | Refuse unsafe access/migration; preserve DB and offer version guidance. |
| `HISTORY_MIGRATION_FAILED` | Approved migration failed or was interrupted | Preserve recoverable prior data; disclose recovery state. |
| `HISTORY_CORRUPT` | DB/report row fails integrity or schema validation | Scope failure to store/record; do not silently repair/delete. |
| `HISTORY_RECORD_NOT_FOUND` | Requested stable scan identity is absent | Return not found, distinct from a valid record with no findings. |
| `HISTORY_ARTIFACT_MISSING` | Referenced raw artifact is unavailable | Preserve valid summary if allowed and disclose missing evidence. |
| `HISTORY_ARTIFACT_DIGEST_MISMATCH` | Retrieved artifact does not match recorded digest | Treat artifact as untrusted; do not claim provenance verified. |
| `API_REQUEST_INVALID` | Path, query, header, or content shape is invalid | Reject before expensive work; provide safe field/context detail. |
| `API_REQUEST_TOO_LARGE` | Request exceeds an approved body/page/response bound | Reject with documented bound; no partial success. |
| `API_VERSION_UNSUPPORTED` | Client requests an unsupported contract version | Fail explicitly and identify supported API version safely. |
| `API_METHOD_UNSUPPORTED` | Method or operation is not part of the approved surface | Reject; do not invoke side effects. |
| `API_ORIGIN_REJECTED` | Host/origin violates the approved local boundary | Reject without disclosing listener internals. |
| `API_TIMEOUT` | Bounded operation deadline elapsed | Report timeout; propagate cancellation and avoid success-shaped partial output. |
| `API_CANCELLED` | Operation was cancelled | Report cancellation distinctly from failed or complete scan. |
| `API_LISTENER_UNAVAILABLE` | Local service could not bind/serve on approved address | Do not fall back to wildcard/public binding; keep CLI usable. |

The final code set may merge/remove candidate categories only after the corresponding failure semantics remain observable. Do not add retry semantics or claim retryability absent an explicit contract.

## 3. Candidate HTTP mapping (not approved)

The OpenAPI design must assign status codes consistently. Candidate mapping only:

| Category | Candidate HTTP status |
| --- | ---: |
| Invalid path/query/content type | 400 |
| Request too large | 413 |
| Unsupported API version | 400 or 404, selected and documented in OpenAPI |
| Origin/Host rejected | 403 |
| Scan history item missing | 404 |
| Store unavailable/locked | 503 or 409, based on bounded conflict semantics |
| Integrity failure/corrupt record | 500 or 422, based on whether resource itself is invalid; must not leak raw data |
| Deadline exceeded | 504 |
| Cancelled | Request cancellation behavior documented; never fabricate a complete response |
| Unexpected internal failure | 500 with redacted safe message and stable request reference |

Final status/code mapping must be singular and testable; avoid multiple statuses for the same equivalent error absent a documented distinction. Do not expose DB driver messages, SQL, stack traces, filesystem absolute paths, credentials, or untrusted input verbatim.

## 4. Error envelope proposal

A candidate JSON envelope:

```json
{
  "error": {
    "code": "HISTORY_UNAVAILABLE",
    "message": "Local scan history is unavailable.",
    "requestId": "opaque-request-id"
  }
}
```

`requestId` is optional and remains undecided. Add fields only when they have a stable consumer contract. The success schema must never be overloaded as an error envelope, and `errors: []` is not a substitute for a failed HTTP request.

## 5. Acceptance checklist

- [ ] Crosswalk candidate codes with existing stable error contracts; retain established codes unchanged.
- [ ] Approve history/API candidate code strings, scope, safe message, and actionable guidance.
- [ ] Approve one OpenAPI status mapping per semantic category.
- [ ] Define behavior for mixed failures, e.g. valid scan metadata with a missing raw artifact.
- [ ] Define cancellation, timeout, listener startup and storage-write failure without false success.
- [ ] Ensure no raw SQL, DB messages, stack trace, token, credential URL, or full environment leaks.
- [ ] Validate OpenAPI examples and exercise each code/status through contract/integration tests.
- [ ] Record owner and independent reviewer approval; no runtime implementation before acceptance.
