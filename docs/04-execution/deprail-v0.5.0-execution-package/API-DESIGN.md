# v0.5 Local API Design Contract

**Status:** Candidate OpenAPI 3.1 artifact and local security profile are available for owner and independent architecture/security review; no API or runtime implementation is approved.
**Design inputs:** [`PRD-v0.5.md`](PRD-v0.5.md), [`UX-DESIGN.md`](UX-DESIGN.md), [`ARCHITECTURE-v0.5.md`](ARCHITECTURE-v0.5.md).

## 1. Goal and constraints

Define the smallest same-origin local API needed by the candidate history and scan-detail UX. The API delegates to shared application services and returns versioned domain meaning. It does not reimplement scan logic, become a general extension platform, or imply team/cloud functionality. The proposed initial surface is read-only; browser-triggered scanning or mutation is excluded unless explicitly reconsidered through release change control.

The candidate OpenAPI 3.1 contract is [`schemas/openapi/v1/openapi.yaml`](../../../schemas/openapi/v1/openapi.yaml). It is a required design-stage artifact before API implementation, not an approved runtime contract.

## 2. Proposed resource model (unapproved)

- **HistoryEntrySummary:** unique `historyEntryID` per stored operation, recorded time, `sourceScanID` from the existing report, `operationOutcome` (`completed|failed|cancelled`), optional `reportStatus` (`complete|partial|failed`), root display label under privacy contract, counts, and minimal provenance summary.
- **HistoryEntryDetail:** summary plus allowlisted projected report metadata (if any), safe diagnostics, and explicit artifact-integrity states; never raw report or project JSON.
- **WorkspaceSummary:** stable workspace ID, repository-relative path and safe ecosystem/manager fields only when captured from the same operation's validated graph. Any exposed `discoveryCompleteness` is discovery-only, not per-workspace scan success; absent graph context is unavailable, not a confirmed empty workspace collection.
- **FindingSummary:** key derived through the existing stable-finding identity rule, package identity, aliases, available severity and fixed version, and validated repository-relative workspace context. Raw artifact digests are report-level references; do not attach one to a finding or invent a dependency path/evidence source without a verified source mapping.
- **ApiError:** stable API error code and concise safe message; no request ID, stack trace, secrets, raw paths, or echoed input.

Existing CLI/report semantics remain the source of scan meaning; the API must not invent vulnerability, completeness, or workspace-scan states. The OpenAPI candidate uses separately paged workspace/finding collections; response fields remain a strict projection of the candidate `history-v1` model, not a copy of its stored JSON.

**Owner-selected report boundary, not yet accepted:** [ADR-0004's direction](../../adr/ADR-0004-local-scan-history.md#owner-selected-direction-independent-history-projection-not-accepted) keeps existing `ScanReport`/CLI serialization unchanged and proposes a separate, allowlisted `history-v1` storage projection. `HistoryEntryDetail` and child collections translate only approved projected fields; they must not embed the stored JSON wholesale, raw `ScanReport`, raw errors, absolute host paths, or fabricated tool/evidence provenance. [`FAILURE-AND-DATA-CONTRACT.md`](requirements/FAILURE-AND-DATA-CONTRACT.md#candidate-history-v1-projection-object) defines candidate source-to-projection fields, including stable keys from `normalize.StableFindingKey`, nullable report/collections and safe diagnostics. The owner and independent reviewer must approve that exact field/unknown-data/redaction mapping before OpenAPI schemas/examples are frozen. A missing or unsafe projection is a typed failure, not an empty successful detail.

## 3. Proposed endpoint matrix (for design review)

| Method and path | Purpose | Proposed result | Write? |
| --- | --- | --- | --- |
| `GET /api/v1/health` | Local API process readiness and API version | Bounded readiness response; no environment dump | No |
| `GET /api/v1/scans` | List saved history | Ordered page of HistoryEntrySummary plus opaque continuation token | No |
| `GET /api/v1/scans/{historyEntryID}` | Load one history entry | Projected detail with artifact-integrity states, or typed not-found error | No |
| `GET /api/v1/scans/{historyEntryID}/workspaces` | Load the bounded workspace collection | Ordered page of WorkspaceSummary | No |
| `GET /api/v1/scans/{historyEntryID}/findings` | Load the bounded finding collection | Ordered page of FindingSummary | No |

The candidate OpenAPI file includes the five operations above. The separate child collections are paged to keep finding/workspace responses bounded. No `POST /scan`, delete, export/upload, policy mutation, remediation, approval, or exception endpoint is specified.

## 4. Query and deterministic ordering

The versioned candidate is [`schemas/openapi/v1/openapi.yaml`](../../../schemas/openapi/v1/openapi.yaml). It defines `pageSize` (default 25, range 1–50) and an opaque, URL-safe keyset `cursor` (maximum 512 characters); there are no filters, offsets, arbitrary sort expressions, or client-provided SQL.

- History is ordered by `recordedAt DESC, historyEntryID DESC`.
- Workspaces are ordered by `workspaceID ASC`; findings are ordered by `stableFindingKey ASC`.
- A cursor is bound to API version, resource, parent history-entry ID where applicable, and the last ordering tuple. It is not an authorization token and contains no repository data.
- No cursor expiry is proposed because v0.5 has no history deletion or automatic eviction. A contract-version change invalidates cursors from the older API version.
- Keyset pagination avoids offset shifts. It does not promise a database snapshot across requests; clients restart pagination to include newly recorded entries.
- Unknown query parameters, malformed cursors, unsupported cursor versions, and page sizes outside the contract fail with `API_REQUEST_INVALID`.

Every successful response is limited to 1 MiB of fully serialized UTF-8 JSON bytes, including framing and required escaping. The server must measure the complete encoded response before sending; it never streams or truncates a partial success. A response above the bound returns `500 API_RESPONSE_TOO_LARGE`, including when all fields satisfy their schemas, because JSON Schema `maxLength` counts code points rather than encoded bytes. The bounded error envelope must itself fit. These candidate limits remain subject to payload, load, and independent-review evidence.

## 5. Response and error semantics

- `/api/v1` is the transport-version boundary. Responses do not repeat an ad hoc API `schemaVersion`; the source report version and `history-v1` projection version remain distinct.
- Existing report completeness (`complete|partial|failed`) and operation outcome (`completed|failed|cancelled`) remain independent. A missing report is `null`, not a zero-finding report.
- A successful empty history page differs from store/API failure. Unavailable, incompatible, corrupt, timed-out, unauthorized, invalid, or response-size-limited outcomes use the typed error envelope in [`requirements/ERROR-MODEL.md`](requirements/ERROR-MODEL.md). Parser-level request-target/header rejections are the explicit 414/431 exceptions described below.
- A valid detail record remains a `200` response when an associated raw artifact is missing, unreadable, or digest-mismatched. Each artifact reference carries an explicit integrity state; the API never serves artifact bytes or filesystem paths. This preserves trustworthy report metadata without claiming the evidence is available or verified.
- The API error envelope contains only a stable code and safe static message. It does not echo input, expose SQL/driver errors, or add a request ID without an approved use case.
- All API routes are bodyless; a prohibited body is `400 API_REQUEST_INVALID`, and there is no `API_REQUEST_TOO_LARGE` body error. Request targets above 2,048 bytes are rejected with transport-level 414; headers above 8,192 bytes with transport-level 431. These parser-boundary responses do not promise an API JSON error envelope. A disconnected client cancels work without a fabricated response; listener startup errors are CLI diagnostics, not HTTP responses.
- Existing v0.1/v0.2 scan error codes and CLI output/exit meanings remain unchanged.

## 6. Local security profile — candidate for independent review

The following decisions are explicit proposals for the OpenAPI/security review. None authorizes a listener implementation.

- Bind only `127.0.0.1:0`; let the OS allocate an ephemeral port. Never retry on wildcard, LAN, or public interfaces. The browser URL uses exactly `http://127.0.0.1:<selected-port>`; `localhost`, forwarded headers, and alternate host aliases are not trusted.
- Run the listener only within an explicit local-console command, not as a background service and never as a side effect of `deprail scan`. The command remains alive until shutdown, stops accepting requests on termination, drains bounded requests for at most 10 seconds, then closes storage. The command name is a separate CLI-contract decision.
- Require an opaque 256-bit process-scoped bearer token on every API route, including health. Generate it from a cryptographically secure random source; revoke it on process shutdown; keep it in browser memory only; never store it in local/session storage or cookies. Missing/invalid credentials return `401 API_AUTH_UNAUTHORIZED`.
- Candidate bootstrap transfers the token only in an initial URL fragment; fragments are not sent in HTTP requests. The client must read it and immediately remove it with `history.replaceState` before making a request. The token is then sent only in the `Authorization` header. It must never appear in a path, query, referrer, log, diagnostic, or persisted browser storage. A reload or new tab loses the in-memory token; show a safe recovery state instructing the user to reopen the console from the active local CLI session, with no cookie/storage fallback. Fragment bootstrap remains a material security-review decision because it can be visible briefly in browser/OS launch state.
- Require the exact `Host` for the selected loopback port. If `Origin` is present, require an exact same-origin match; reject `null` and mismatched origins. If `Sec-Fetch-Site` is present, require `same-origin`. Do not emit CORS allow headers or trust forwarded headers. The bearer token remains required when `Origin` is absent.
- Accept `GET` only. Reject `HEAD`, `OPTIONS`, other methods, unknown query parameters, and all request bodies; return `405` for an unsupported method and `400 API_REQUEST_INVALID` for a prohibited body or invalid input. No filesystem path, SQL fragment, process operation, or artifact bytes are accepted from or returned to the client.
- Candidate finite bounds: request target 2,048 bytes (over-limit transport 414); headers 8,192 bytes (over-limit transport 431); eight concurrent requests; 10-second request deadline; 10-second shutdown drain; 1 MiB maximum response; page size 1–50 (default 25); cursor at most 512 characters. Saturation returns `503 API_BUSY`; a deadline returns `504 API_TIMEOUT`; an oversized response returns `500 API_RESPONSE_TOO_LARGE` without truncation. Parser-level 414/431 responses may not carry the API JSON envelope.
- Return `Cache-Control: no-store` for API responses. Never log authorization headers, tokens, sensitive request headers, raw paths, repository data, or database/stack details. Make no external network requests.
- If binding fails, report `API_LISTENER_UNAVAILABLE` safely and leave existing CLI commands usable.

The profile addresses DNS rebinding, hostile browser origins, local API authorization, request/resource bounds, and lifecycle behavior as one reviewable proposal. Exact controls, token bootstrap, browser behavior, and residual risks require independent architecture/security approval.

## 7. Compatibility and versioning

The candidate source uses OpenAPI `3.1.0`, contract `info.version: 1.0.0`, and the `/api/v1` URL boundary. These versions are distinct from SQLite `user_version`, the `history-v1` projection, and a source `ScanReport` version. The OpenAPI file and examples must be validated before implementation; additive response evolution is not assumed safe for strict clients.

Breaking API changes require a new API version or an explicit reviewed migration. Existing CLI JSON, v1alpha report meaning, scan error codes, and exit-code behavior remain unchanged absent a separate compatibility decision. The embedded UI and API must ship as a matching version; mismatch fails visibly rather than serving empty or stale success.

## 8. API design acceptance checklist

The OpenAPI structure, local references, response examples, and status-specific error constraints have been checked for this review candidate. No schema-derived maximum response size is claimed: `maxLength` counts Unicode code points, not serialized bytes. The 1 MiB UTF-8 response cap is a runtime byte check and remains unverified. These checks do not constitute owner or independent-review acceptance.

- [ ] UX acceptance maps every operation to a reviewed screen/action.
- [ ] Resource representation and pagination are accepted against payload and usage evidence.
- [ ] OpenAPI 3.1 operations, schemas, examples, error mappings, and candidate bounds are independently reviewed and accepted.
- [ ] Exact `history-v1` projection fields, provenance, diagnostics, and sensitive-field redaction are accepted.
- [ ] Bind/origin/CORS/CSRF/authentication, token transfer, process lifecycle, and path-serving decisions have independent security review.
- [ ] Cancellation, timeouts, malformed/oversized inputs, missing/corrupt records, artifact-integrity states, and storage errors are mapped.
- [ ] API/CLI shared-service boundary and no-mutation scope are explicit and compatible.
- [ ] API/UI version behavior and the supported-client compatibility policy are accepted.
- [ ] Owner and independent reviewer separately accept the API contract before API implementation issues are created.

## 9. Candidate profile and remaining review decisions

The OpenAPI file is a concrete review candidate, not an accepted public contract. It defines the five read-only operations in §3 and rejects unbounded, duplicate, or mutation-oriented alternatives.

| Decision | Candidate | Evidence or remaining risk |
| --- | --- | --- |
| Resources | `GET /health`, `/scans`, `/scans/{historyEntryID}`, and paged `/workspaces` and `/findings` child collections. | Each operation maps to history/detail workflows in `UX-DESIGN.md`; no scan, delete, export, publish, remediation, or mutation route exists. |
| History projection | Strict response allowlists; null report and unavailable collections remain distinct from confirmed empty collections. Finding identity uses `normalize.StableFindingKey`; `TargetID` remains only the vulnerability ID. | Exact persisted `history-v1` fields, typed diagnostics, unknown-data policy, and same-operation graph capture remain separate review gates. The current `app.ScanReport` has no tool/database metadata field, so none is invented. |
| Artifact integrity | Detail returns digest plus `verified`, `missing`, `digest_mismatch`, or `unavailable`; a safe detail remains `200` when metadata is trustworthy. | The API exposes neither artifact bytes nor paths. Owner/reviewer must accept this mixed-integrity behavior against `SEC-09` and UX. |
| Pagination | Deterministic keyset; page size 25 by default, maximum 50; 512-character versioned cursor; no expiry while v0.5 has no deletion/eviction. | No snapshot-isolation guarantee across requests. Cursor and insertion behavior require review; consumers restart to include new entries. |
| Response bounds | Candidate cap: 1 MiB per fully serialized UTF-8 JSON response; no truncation. | Earlier schema-size figures used character-count maxima and are withdrawn as byte-budget evidence. Schema-valid payloads may exceed the wire cap and receive `API_RESPONSE_TOO_LARGE`. Before implementation, measure serialized bytes with Unicode/escaping and exact-boundary cases; record runtime, memory, and platform evidence. |
| Listener and auth | `127.0.0.1:0`, same-origin UI/API, process-scoped 256-bit bearer token, no CORS, strict Host/Origin/Fetch-Metadata checks. | Fragment-only bootstrap is not sent over HTTP but is visible briefly in the browser/OS launch state; independent reviewer must accept or replace this mechanism. |
| Runtime bounds | 2,048-byte request target, 8,192-byte headers, eight concurrent requests, 10-second request deadline and shutdown drain. | These remain proposed limits; actual platform/load evidence and owner/reviewer acceptance remain required before implementation. |
| Errors | Typed, safe API envelope; 400/401/403/404/405/500/503/504 mappings, plus parser-level 414/431 rejection for request-target/header limits without a JSON-envelope guarantee. No request bodies. Artifact-integrity states are part of a successful detail response. | OpenAPI response schemas constrain each declared status to its allowed code set; runtime status/code behavior remains unverified. |
| Versioning | `/api/v1`, OpenAPI 3.1, `info.version: 1.0.0`; no redundant response API-version field. | UI/API mismatch must fail visibly. Any breaking change requires a separately reviewed version change. |

The remaining acceptance gates are independent review of this exact artifact, owner acceptance of the candidate decisions, exact storage-projection approval, supported browser and platform decisions, and the complete v0.5 Definition of Ready. Until those gates are linked, this contract remains proposed and no implementation issue is ready.

## 10. Candidate validation evidence

- `openapi-spec-validator 0.9.0` accepted the parsed OpenAPI 3.1 document with local references.
- `openapi_schema_validator.OAS31Validator` validated all 62 operation-response examples and 15 component-response examples, including the stricter status/code schemas and unavailable collection states.
- Earlier schema-size figures used character-count maxima and are withdrawn; they did not establish UTF-8 wire-byte maxima. The candidate limit applies to the complete serialized response, including JSON escaping. Schema-valid values may exceed 1 MiB and must produce a bounded `500 API_RESPONSE_TOO_LARGE` response without partial output. Runtime Unicode and exact-boundary tests remain required.
- A schema-valid default-size page of 25 findings with maximum-length four-byte Unicode strings and 32 aliases per item serialized to 1,544,660 UTF-8 bytes (4,615,860 bytes with ASCII escaping). It exceeds the proposed 1 MiB cap and demonstrates why code-point limits do not guarantee a wire-byte bound.
- These are offline contract checks only. They do not claim runtime, real-payload, browser, platform, performance, or security-test evidence.
