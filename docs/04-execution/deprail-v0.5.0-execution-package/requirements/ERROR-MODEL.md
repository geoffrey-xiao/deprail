# v0.5 Error Model

**Status:** Draft proposal; code strings and transport mappings require owner and independent architecture/security review.

## 1. Contract

Errors are structured with stable code, safe user-facing message, scope, and actionable context. HTTP status codes express transport-level outcome; they do not replace application error codes. A failed history/API operation never becomes an empty successful result, and an API error must not rewrite the underlying scan outcome.

The existing scan-error inventory is [`v0.1 ERROR-MODEL.md`](../../deprail-v0.1-execution-package/requirements/ERROR-MODEL.md); [`v0.2 ERROR-MODEL.md`](../../deprail-v0.2-execution-package/requirements/ERROR-MODEL.md) reinforces required result behavior and redaction. Preserve their established scan codes and meanings. [`v0.3.1 ERROR-MODEL.md`](../../deprail-v0.3.1-execution-package/requirements/ERROR-MODEL.md) requires cancelled operations to remain visibly cancelled without redefining report completeness. The separate [`v0.3 ERROR-MODEL.md`](../../deprail-v0.3-execution-package/requirements/ERROR-MODEL.md) covers remediation planning and is not a source of scan/API error codes; the v0.4 package has no `requirements/ERROR-MODEL.md`. The table below lists candidate v0.5 history/API categories only; no string is stable until crosswalked against the named applicable contracts and approved.

`operationOutcome=cancelled` is a history-entry state, not a candidate stable error code and not a `ScanReport.Status`. `API_CANCELLED` applies only to a cancelled read/API request; it does not describe or mutate a scan operation.

## 2. Candidate history/API errors

| Candidate code | Meaning | Required behavior |
| --- | --- | --- |
| `HISTORY_UNAVAILABLE` | Local store cannot be opened/read or access is denied | Report store unavailable; do not return empty history or weaken permissions. |
| `HISTORY_WRITE_FAILED` | A history transaction could not be durably committed | Do not claim saved; retain scan outcome separately and report persistence failure. |
| `HISTORY_SCHEMA_UNSUPPORTED` | Store uses an unsupported newer/older schema | Refuse unsafe access/migration; preserve DB and offer version guidance. |
| `HISTORY_MIGRATION_FAILED` | Approved migration failed or was interrupted | Preserve recoverable prior data; disclose recovery state. |
| `HISTORY_CORRUPT` | DB/report row fails integrity or schema validation | Scope failure to store/record; do not silently repair/delete. |
| `HISTORY_ENTRY_NOT_FOUND` | Requested unique history-entry ID is absent | Return not found, distinct from a valid entry whose report has no findings. |
| `HISTORY_ARTIFACT_MISSING` | Referenced raw artifact is unavailable | Preserve valid summary if allowed and disclose missing evidence. |
| `HISTORY_ARTIFACT_DIGEST_MISMATCH` | Retrieved artifact does not match recorded digest | Treat artifact as untrusted; do not claim provenance verified. |
| `API_REQUEST_INVALID` | Query, header, content shape, path identifier, filter, or page-size value is invalid | Reject before expensive work; provide safe field/context detail. |
| `API_REQUEST_TOO_LARGE` | A request body exceeds an approved body limit | Reject with 413 and no partial success; this candidate applies only to a body-bearing operation if one is approved. |
| `API_RESPONSE_TOO_LARGE` | A valid request would require a response above the server's approved response bound | Return an explicit bounded failure; never truncate into an apparently successful response. Candidate code requires review. |
| `API_ROUTE_NOT_FOUND` | Path is not part of the approved API route table | Return not found; do not classify an unknown route as an unsupported method. |
| `API_VERSION_UNSUPPORTED` | Client requests an unsupported contract version | Fail explicitly and identify supported API version safely. |
| `API_METHOD_UNSUPPORTED` | Method is not allowed for a known route | Reject without invoking side effects; distinguish from an unknown route. |
| `API_ORIGIN_REJECTED` | Host/origin violates the approved local boundary | Reject without disclosing listener internals. |
| `API_TIMEOUT` | Bounded read/API request deadline elapsed | Report timeout; do not change persisted history or scan operation outcome. |
| `API_CANCELLED` | A read/API request was cancelled by the caller or deadline | Terminate the request without changing scan operation outcome or persisted report state. |
| `API_LISTENER_UNAVAILABLE` | Local service could not bind/serve on approved address | Do not fall back to wildcard/public binding; keep CLI usable. |

The final code set may merge/remove candidate categories only after the corresponding failure semantics remain observable. Do not add retry semantics or claim retryability absent an explicit contract.

## 3. Candidate HTTP mapping (not approved)

The table recommends one candidate HTTP status per API error category to make the draft deterministic. These mappings remain unapproved until the OpenAPI contract is reviewed.

| Category | Recommended candidate HTTP status |
| --- | ---: |
| Invalid path identifier, query, content type, filter, or page-size value (`API_REQUEST_INVALID`) | 400 |
| Request body exceeds its bound (`API_REQUEST_TOO_LARGE`) | 413; candidate only for an approved body-bearing operation |
| A valid request would exceed the response-size bound (`API_RESPONSE_TOO_LARGE`) | 500; do not truncate the response |
| Unknown API route (`API_ROUTE_NOT_FOUND`) | 404 |
| Unsupported API version (`API_VERSION_UNSUPPORTED`) | 404 |
| Unsupported method on a known route (`API_METHOD_UNSUPPORTED`) | 405 with an `Allow` header |
| Host or origin rejected (`API_ORIGIN_REJECTED`) | 403 |
| History entry missing (`HISTORY_ENTRY_NOT_FOUND`) | 404 |
| Store unavailable/locked, unsupported schema, or migration failure | 503 |
| Corrupt store/report or required artifact missing/digest mismatch | 500 |
| Server-side request deadline (`API_TIMEOUT`) | 504 |
| Unexpected internal failure | 500 |

`API_CANCELLED` describes a cancelled read request; when the client has disconnected, stop work and return no fabricated response. `API_LISTENER_UNAVAILABLE` is a startup failure, not an HTTP response. `HISTORY_WRITE_FAILED` is not reachable through the proposed read-only API and remains a distinct application/CLI persistence outcome. The current GET-only surface has no request-body route, so `API_REQUEST_TOO_LARGE` is unreachable unless a body-bearing operation is separately approved. Page-size bounds are invalid request parameters; output-limit failures use the distinct candidate `API_RESPONSE_TOO_LARGE`.

Do not expose DB driver messages, SQL, stack traces, filesystem absolute paths, credentials, or untrusted input verbatim. The recommended statuses and their mapping to existing stable error contracts require owner and independent architecture/security review before acceptance.

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

- [x] Crosswalk candidate codes to existing stable error contracts and retain established codes unchanged; candidate-code/state/evidence mapping is in §6. Owner/reviewer approval remains outstanding.
- [ ] Approve history/API candidate code strings, scope, safe message, and actionable guidance.
- [ ] Approve one OpenAPI status mapping per semantic category.
- [ ] Define behavior for mixed failures, e.g. valid scan metadata with a missing raw artifact.
- [ ] Define cancellation, timeout, listener startup and storage-write failure without false success.
- [ ] Ensure no raw SQL, DB messages, stack trace, token, credential URL, or full environment leaks.
- [ ] Validate OpenAPI examples and exercise each code/status through contract/integration tests.
- [ ] Record owner and independent reviewer approval; no runtime implementation before acceptance.

## 6. Candidate code-to-state and evidence crosswalk

Each candidate maps to a consumer-visible outcome and a planned scenario in [`TEST-STRATEGY.md`](TEST-STRATEGY.md#9-requirement-and-threat-evidence-matrix). These are future verification requirements, not completed tests or finalized wire semantics.

| Candidate code | Required visible behavior | Verification ID |
| --- | --- | --- |
| `HISTORY_UNAVAILABLE` | Failed list/detail state, distinct from valid empty history; safe retry only when appropriate. | `FR-502` |
| `HISTORY_WRITE_FAILED` | For separately approved CLI capture only: safe stderr persistence diagnostic, no saved claim, scan report/outcome preserved; no read-only API response. | `STORE-01` |
| `HISTORY_SCHEMA_UNSUPPORTED`, `HISTORY_MIGRATION_FAILED`, `HISTORY_CORRUPT` | Explicit incompatibility/recovery state; preserve existing DB/backup and never suggest automatic reset or downgrade. | `FR-508`, `SEC-08` |
| `HISTORY_ENTRY_NOT_FOUND` | Stale-selection/not-found state, distinct from an entry with no report or findings. | `FR-503` |
| `HISTORY_ARTIFACT_MISSING`, `HISTORY_ARTIFACT_DIGEST_MISMATCH` | Preserve only trustworthy metadata; disclose unavailable/unverified evidence, never fabricate or trust bytes. | `FR-503`, `SEC-09` |
| `API_REQUEST_INVALID` | Safe client error before expensive work; no empty-success or side effect. | `FR-507`, `SEC-03`, `SEC-04` |
| `API_REQUEST_TOO_LARGE`, `API_RESPONSE_TOO_LARGE` | Explicit bounded failure, never partial/truncated success. Request-body case applies only if a body-bearing route is separately approved. | `SEC-05` |
| `API_ROUTE_NOT_FOUND`, `API_VERSION_UNSUPPORTED`, `API_METHOD_UNSUPPORTED` | Distinguish unknown route, unsupported contract version, and unsupported method using the candidate status mapping in §3; never render empty history. | `FR-507`, `FR-509` |
| `API_ORIGIN_REJECTED` | Reject the hostile origin without disclosing listener internals. | `SEC-01` |
| `API_TIMEOUT`, `API_CANCELLED` | End only the affected read/request; no fabricated response and no rewrite of stored history or scan outcome. | `FR-502`, `FR-503` |
| `API_LISTENER_UNAVAILABLE` | Safe startup failure before the browser surface exists; existing CLI remains usable and no public-bind fallback occurs. | `FR-505`, `SEC-11` |

Existing v0.1/v0.2 scan codes, CLI exit meanings, report schemas, and artifact identities remain unchanged. Candidate history/API codes and any additive CLI exit code require separate owner and independent-review approval before they become stable.
