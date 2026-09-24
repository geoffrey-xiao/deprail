# v0.5 Error Model

**Status:** Candidate history/API codes and OpenAPI mappings are documented in [`schemas/openapi/v1/openapi.yaml`](../../../../schemas/openapi/v1/openapi.yaml); owner technical acceptance remains pending. External review is optional under [ADR-0005](../../../adr/ADR-0005-solo-owner-review-policy.md).

## 1. Contract

Errors are structured with stable code, safe user-facing message, scope, and actionable context. HTTP status codes express transport-level outcome; they do not replace application error codes. A failed history/API operation never becomes an empty successful result, and an API error must not rewrite the underlying scan outcome.

The existing scan-error inventory is [`v0.1 ERROR-MODEL.md`](../../deprail-v0.1-execution-package/requirements/ERROR-MODEL.md); [`v0.2 ERROR-MODEL.md`](../../deprail-v0.2-execution-package/requirements/ERROR-MODEL.md) reinforces required result behavior and redaction. Preserve their established scan codes and meanings. [`v0.3.1 ERROR-MODEL.md`](../../deprail-v0.3.1-execution-package/requirements/ERROR-MODEL.md) requires cancelled operations to remain visibly cancelled without redefining report completeness. The separate [`v0.3 ERROR-MODEL.md`](../../deprail-v0.3-execution-package/requirements/ERROR-MODEL.md) covers remediation planning and is not a source of scan/API error codes; the v0.4 package has no `requirements/ERROR-MODEL.md`. The table below lists candidate v0.5 history/API categories only; no string is stable until crosswalked against the named applicable contracts and approved.

`operationOutcome=cancelled` is a history-entry state, not a candidate stable error code and not a `ScanReport.Status`. `API_CANCELLED` applies only to a cancelled read/API request; it does not describe or mutate a scan operation.

## 2. Candidate history/API errors and integrity states

Errors are stable application codes with safe messages. HTTP status describes the transport result and does not replace the application code. The OpenAPI candidate exposes only the read/API codes listed in its `ApiError` schema.

| Candidate code or state | Meaning | Required behavior |
| --- | --- | --- |
| `HISTORY_UNAVAILABLE` | Local store cannot be opened/read or access is denied | Report store unavailable; do not return empty history or weaken permissions. |
| `HISTORY_WRITE_FAILED` | A selected history projection cannot be safely validated or durably committed | Do not claim saved or create partial rows/references; retain the original scan result and report persistence failure separately. Not reachable through the read-only API. |
| `HISTORY_SCHEMA_UNSUPPORTED` | Store uses an unsupported schema | Refuse unsafe access/migration; preserve the database and offer version guidance. |
| `HISTORY_MIGRATION_FAILED` | Approved migration failed or was interrupted | Preserve recoverable prior data; disclose recovery state. |
| `HISTORY_CORRUPT` | Stored row/projection fails integrity, version, or consistency checks | Scope failure to the record/store; never silently repair, delete, or present empty success. |
| `HISTORY_ENTRY_NOT_FOUND` | Requested valid history-entry ID is absent | Return 404, distinct from a valid entry with no report/findings. |
| `HISTORY_ARTIFACT_MISSING` | Referenced raw artifact is absent | Preserve trustworthy metadata and disclose unavailable evidence. In API detail, represented as integrity state `missing`, not an error envelope. |
| `HISTORY_ARTIFACT_DIGEST_MISMATCH` | Retrieved artifact digest differs from the recorded digest | Do not trust bytes or claim verified provenance. In API detail, represented as integrity state `digest_mismatch`, not an error envelope. |
| `API_REQUEST_INVALID` | Identifier, query, cursor, or prohibited body is invalid | Reject before expensive work; use safe generic guidance and do not echo input. |
| `API_AUTH_UNAUTHORIZED` | Required process-scoped bearer credential is absent or invalid | Return 401 with a generic message and `WWW-Authenticate: Bearer`; do not distinguish secret values. |
| `API_RESPONSE_TOO_LARGE` | A fully serialized UTF-8 JSON success response exceeds the proposed 1 MiB byte bound | Return an explicit bounded 500 response before sending a partial body; never truncate into success. |
| `API_ROUTE_NOT_FOUND` | Path is not part of the candidate route table | Return 404; distinguish from unsupported method. |
| `API_VERSION_UNSUPPORTED` | Versioned API prefix is recognized but unsupported | Return 404 with a safe supported-version message. |
| `API_METHOD_UNSUPPORTED` | Method is not allowed for a known route | Return 405 with `Allow: GET`; invoke no side effect. |
| `API_ORIGIN_REJECTED` | Host, Origin, or Fetch Metadata violates the local boundary | Return 403 without disclosing listener details. |
| `API_BUSY` | The proposed bounded concurrent-request capacity is full | Return 503; do not queue unbounded work or retry implicitly. |
| `API_TIMEOUT` | Bounded read request deadline elapsed | Return 504; do not change history or scan outcome. |
| `API_CANCELLED` | Client disconnect cancels a read | Stop affected work; the disconnected client receives no fabricated response. |
| `API_INTERNAL_ERROR` | Unexpected internal API failure | Return a generic 500; keep driver details, SQL, paths, and stack traces out of the response. |
| `API_LISTENER_UNAVAILABLE` | Local service cannot bind/serve on the candidate loopback address | Report a safe startup diagnostic; never fall back to wildcard/public binding; keep CLI usable. |

All API routes are bodyless, so the candidate has no body-size error code or 413 response; a prohibited body fails as `API_REQUEST_INVALID`. Request-target and header-block size failures are rejected at the HTTP parser boundary with 414 or 431 and do not promise a stable API code or JSON envelope.

For the owner-selected `history-v1` direction, projection construction remains a pre-commit safety boundary: raw `ScanReport.Errors`, absolute roots, missing stable finding identity, and unreviewed required fields are not silently stored or discarded. An unsupported future projection schema is `HISTORY_SCHEMA_UNSUPPORTED`; a malformed row at a supported version is `HISTORY_CORRUPT`. Neither is empty success.

## 3. Candidate HTTP mapping

These mappings are concrete OpenAPI review proposals, not accepted stable behavior.

| HTTP status | Candidate code/condition |
| ---: | --- |
| 400 | `API_REQUEST_INVALID` |
| 401 | `API_AUTH_UNAUTHORIZED` |
| 403 | `API_ORIGIN_REJECTED` |
| 404 | `HISTORY_ENTRY_NOT_FOUND`; router-level `API_ROUTE_NOT_FOUND` or `API_VERSION_UNSUPPORTED` |
| 405 | `API_METHOD_UNSUPPORTED`, with `Allow: GET` |
| 414 | Request target exceeds 2,048 bytes; parser-level rejection before route execution, no stable API code/body guaranteed. |
| 431 | Request headers exceed 8,192 bytes; parser-level rejection before route execution, no stable API code/body guaranteed. |
| 500 | `API_RESPONSE_TOO_LARGE`, `HISTORY_CORRUPT`, or `API_INTERNAL_ERROR` |
| 503 | `HISTORY_UNAVAILABLE`, `HISTORY_SCHEMA_UNSUPPORTED`, `HISTORY_MIGRATION_FAILED`, or `API_BUSY` |
| 504 | `API_TIMEOUT` |

`HISTORY_ARTIFACT_MISSING` and `HISTORY_ARTIFACT_DIGEST_MISMATCH` are represented in a successful detail response as per-reference `integrity` states (`missing` and `digest_mismatch`), preserving trustworthy metadata without claiming evidence is available or verified. An unreadable reference is `unavailable`. Artifact bytes and filesystem paths are never returned.

`API_CANCELLED` and `API_LISTENER_UNAVAILABLE` do not produce HTTP responses: the former ends when the client disconnects; the latter is a startup diagnostic. `HISTORY_WRITE_FAILED` is an application/CLI persistence outcome, not reachable via the read-only API. A response over the size cap fails explicitly and is never truncated. Parser-level 414/431 rejections occur before API routing and may not use the JSON error envelope.

Do not expose DB driver messages, SQL, stack traces, filesystem absolute paths, credentials, or untrusted input verbatim. Router-level unknown-route and unsupported-version outcomes use the same safe JSON error envelope as declared operations.

## 4. Error envelope

The candidate JSON envelope contains only a stable code and safe message:

```json
{
  "error": {
    "code": "HISTORY_UNAVAILABLE",
    "message": "Local scan history is unavailable."
  }
}
```

Messages are static and do not echo input or expose secrets, raw paths, SQL, database-driver errors, or stack traces. The success schema is never overloaded as an error envelope, and `errors: []` is not a substitute for a failed HTTP request. A request ID is omitted because no stable consumer use case is established.

## 5. Acceptance checklist

- [x] Candidate codes are cross-walked to existing stable scan contracts without changing established meanings; see §6. Owner acceptance remains outstanding.
- [ ] Owner accepts history/API codes, scope, safe messages, and guidance.
- [ ] Owner accepts one HTTP mapping per reachable semantic category.
- [ ] Mixed artifact-integrity behavior is accepted; trustworthy detail remains distinct from missing/unverified evidence.
- [ ] Cancellation, timeout, listener startup, storage failure, and malformed/oversized inputs remain explicit without false success.
- [ ] No raw SQL, database message, stack trace, token, credential URL, or full environment is returned/logged.
- [x] OpenAPI candidate examples and status-specific error constraints were validated offline; runtime response-byte enforcement and contract/integration tests remain future work.
- [ ] Owner records acceptance of the error contract; no runtime implementation before the complete technical DoR passes. External review is optional under ADR-0005.

## 6. Candidate code-to-state and evidence crosswalk

These mappings connect consumer-visible behavior to the planned scenario IDs in [`TEST-STRATEGY.md`](TEST-STRATEGY.md#9-requirement-and-threat-evidence-matrix). They are contract proposals, not runtime test results.

| Code/state | Required visible behavior | Verification ID |
| --- | --- | --- |
| `HISTORY_UNAVAILABLE` | Failed list/detail state, distinct from valid empty history. | `FR-502` |
| `HISTORY_WRITE_FAILED` | Explicitly selected capture fails safely with no partial row or saved claim; scan outcome/report remain unchanged. | `STORE-01`, `SEC-07` |
| `HISTORY_SCHEMA_UNSUPPORTED`, `HISTORY_MIGRATION_FAILED`, `HISTORY_CORRUPT` | Explicit incompatibility/recovery; preserve data and never suggest reset or downgrade. | `FR-508`, `SEC-08` |
| `HISTORY_ENTRY_NOT_FOUND` | Stale-selection state, distinct from an entry with no report/findings. | `FR-503` |
| Artifact integrity `verified|missing|digest_mismatch|unavailable` | Preserve trustworthy detail; disclose unavailable/unverified evidence and never fabricate bytes or a clean result. | `FR-503`, `SEC-09` |
| `API_REQUEST_INVALID` | Reject malformed IDs, query values, cursors, or prohibited bodies without side effects or empty success. | `FR-507`, `SEC-03`, `SEC-04` |
| Transport-level 414/431 | Reject over-limit request targets/headers before routing; no stable API error body is promised. | `SEC-05` |
| `API_AUTH_UNAUTHORIZED` | Reject absent/invalid bearer without disclosing token details or returning local data. | `SEC-01`, `SEC-02` |
| `API_RESPONSE_TOO_LARGE`, `API_BUSY` | Explicit bounded failure; never truncate success or perform unbounded work. | `SEC-05` |
| `API_ROUTE_NOT_FOUND`, `API_VERSION_UNSUPPORTED`, `API_METHOD_UNSUPPORTED` | Distinguish unknown route, unsupported version, and unsupported method. | `FR-507`, `FR-509` |
| `API_ORIGIN_REJECTED` | Reject hostile Host/Origin/Fetch Metadata without listener disclosure. | `SEC-01`, `SEC-11` |
| `API_TIMEOUT`, `API_CANCELLED` | End only the affected request; do not rewrite saved history or scan outcome. | `FR-502`, `FR-503` |
| `API_LISTENER_UNAVAILABLE` | Safe startup failure with no public-bind fallback; existing CLI remains usable. | `FR-505`, `SEC-11` |

Existing v0.1/v0.2 scan codes, CLI exit meanings, report schemas, and artifact identities remain unchanged. New history/API codes and any CLI persistence exit behavior require a separate owner decision before they become stable.
