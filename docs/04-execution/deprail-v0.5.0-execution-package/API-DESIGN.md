# v0.5 Local API Design Contract

**Status:** Draft proposal; no OpenAPI schema or runtime API is approved.
**Design inputs:** [`PRD-v0.5.md`](PRD-v0.5.md), [`UX-DESIGN.md`](UX-DESIGN.md), [`ARCHITECTURE-v0.5.md`](ARCHITECTURE-v0.5.md).

## 1. Goal and constraints

Define the smallest same-origin local API needed by the approved history and scan-detail UX. The API delegates to shared application services and returns versioned domain meaning. It does not reimplement scan logic, become a general extension platform, or imply team/cloud functionality. The proposed initial surface is read-only; browser-triggered scanning or mutation is excluded unless explicitly reconsidered through release change control.

An OpenAPI 3.1 document with JSON Schema-compatible examples and validation is a required design-stage deliverable before API implementation. This Markdown proposal is not a substitute for that artifact.

## 2. Proposed resource model (unapproved)

- **HistoryEntrySummary:** unique `historyEntryID` per stored operation, recorded time, `sourceScanID` from the existing report, `operationOutcome` (`completed|failed|cancelled`), optional `reportStatus` (`complete|partial|failed`), root display label under privacy contract, counts, and minimal provenance summary.
- **HistoryEntryDetail:** summary plus scan contract/version, optional report, diagnostics, workspaces, and stable links/identifiers for findings and artifacts.
- **WorkspaceSummary:** stable workspace ID, repository-relative path, ecosystem/manager metadata permitted by the data contract, and workspace completeness.
- **FindingSummary:** stable finding key, package identity, aliases, severity/source as available, affected workspace/path context, fixed-version data, and evidence references subject to existing report schema.
- **ApiError:** stable API error code, concise safe message, retryability only when meaningful, request/correlation ID if approved, and no stack trace/secrets/raw paths by default.

Existing report schemas remain the source of semantic meaning; the API must not invent new vulnerability or completeness states. Exact representation of an approved privacy-safe projection (bounded embedded detail versus paged child resources) remains open pending size, UX, and compatibility evidence.

**Unresolved report boundary:** [ADR-0004's validation finding](../../adr/ADR-0004-local-scan-history.md#post-proposal-validation-finding-source-report-is-not-the-proposed-stored-document) shows that the actual `ScanReport` contains an absolute repository root and does not validate as the proposed unchanged v1alpha scan document with findings/errors. `HistoryEntryDetail` and child collections must not embed raw `ScanReport` JSON, expose host paths, or claim a redacted copy is the original schema. The owner and independent reviewer must decide a privacy-safe, versioned history projection and its provenance/unknown-field behavior (or a separately reviewed CLI/schema compatibility change) before exact OpenAPI fields/examples can be frozen. A failure to form a trustworthy projection remains a typed failure, not a successful empty detail.

## 3. Proposed endpoint matrix (for design review)

| Method and path | Purpose | Proposed result | Write? |
| --- | --- | --- | --- |
| `GET /api/v1/health` | Local service readiness/schema version | Bounded readiness response; no environment dump | No |
| `GET /api/v1/scans` | List saved history | Ordered page of HistoryEntrySummary plus opaque continuation token | No |
| `GET /api/v1/scans/{historyEntryID}` | Load one history entry | HistoryEntryDetail or typed not-found/integrity error | No |
| `GET /api/v1/scans/{historyEntryID}/workspaces` | Load bounded workspace list if not embedded | Ordered page of WorkspaceSummary | No |
| `GET /api/v1/scans/{historyEntryID}/findings` | Load bounded finding list if not embedded | Ordered page of FindingSummary | No |

Only include separate child collection endpoints if UX/resource size needs justify them. Avoid parallel duplicate ways to fetch the same data. No `POST /scan`, delete, export/upload, policy mutation, remediation, approval, or exception endpoint is proposed.

## 4. Query and deterministic ordering proposal

- List order should be explicit and stable (candidate: recorded timestamp descending, unique `historyEntryID` as tie-breaker); final semantics require approval.
- Pagination must be bounded, opaque to clients, and stable under concurrent insertion. Prefer a cursor over unbounded offset paging if the chosen storage/index contract supports it.
- Maximum page size, filter allowlist, cursor version/expiry, and behavior when records are deleted between pages are open contract fields; never accept arbitrary SQL/order expressions.
- Query values and path identifiers are validated. Unknown filters/fields fail clearly rather than being silently ignored.

## 5. Response and error semantics

- Responses carry an explicit API/schema version according to the reviewed OpenAPI design; avoid redundant ad hoc version markers if URL versioning is sufficient.
- Existing report completeness values (`complete`, `partial`, `failed`) remain unchanged end-to-end. Cancellation is an execution outcome on the history entry, not a `ScanReport.Status` value or report-schema extension.
- Empty scan history is a successful empty page. Store unavailable, migration rejected, unsupported version, corrupt row, artifact digest mismatch, or timeout are explicit non-success responses.
- A missing history entry is distinct from an entry whose optional report has zero findings.
- Use consistent JSON error envelope and HTTP status mapping; candidate codes in [`requirements/ERROR-MODEL.md`](requirements/ERROR-MODEL.md) must be cross-walked with existing stable contracts before implementation. No new stable code is finalized here.
- User-visible errors are redacted and actionable; server diagnostics never include credentials, full environment, or untrusted raw content.

## 6. Local security model — decisions required

Proposed baseline for review: bind loopback only, serve UI and API from the same origin, reject unexpected Host/Origin values, no arbitrary CORS, no forwarded-header trust, no network fetches, no credentials in URLs, no mutation routes, and bounded bodies/timeouts/concurrency. A local API remains reachable by other local processes and can be targeted by hostile web origins; loopback alone is not sufficient protection.

The security review must decide and test:

- Bind address, port selection, port conflict, URL disclosure, process lifetime, and shutdown.
- Browser origin validation, DNS rebinding/Host handling, CORS, CSRF implications, and whether local authentication/token is needed.
- Strict method/path/content-type handling; request, response, page, and concurrency limits; cancellation/deadline propagation.
- Static asset/API route separation, encoded path normalization, traversal/symlink behavior, and no host filesystem exposure.
- Error, access, and debug log redaction and local data access/privacy expectations.

No LAN/public binding or authentication bypass is approved by this draft.

## 7. Compatibility and versioning

The design package must include an OpenAPI 3.1 file, schemas/examples, contract version policy, and consumer compatibility rules before code. Additive response evolution is not automatically safe if strict clients or generated types are used. Breaking changes need a new API contract/version or explicit migration. Existing CLI JSON and v1alpha report contracts are unaffected absent separate approval. Static UI asset version and API version mismatch must fail visibly and safely.

## 8. API design acceptance checklist

- [ ] Approved UX maps each endpoint to an observable screen/action.
- [ ] Resource representation and pagination are chosen from payload/usage evidence.
- [ ] OpenAPI 3.1 operations, schemas, examples, security, limits, and error responses are complete and validated.
- [ ] Data scope/provenance and sensitive-field redaction are reviewed.
- [ ] Bind/origin/CORS/CSRF/auth and path-serving decisions have independent security review.
- [ ] Cancellation, timeouts, malformed/oversized inputs, missing/corrupt records, and storage errors are mapped.
- [ ] API/CLI shared-service boundary and no-mutation scope are explicit.
- [ ] Compatibility and UI/API asset-version behavior are approved.
- [ ] Owner and independent reviewer accept the API contract before API implementation issues are created.

## 9. Recommended draft profile (unapproved)

The following recommendations make the candidate API reviewable. They do not approve an OpenAPI contract, numeric limits, listener behavior, or runtime implementation.

| Decision | Draft recommendation | Rationale and remaining evidence |
| --- | --- | --- |
| Read surface | Keep the listed endpoints `GET`-only; reject request bodies and unsupported methods. Do not add scan, delete, export, or mutation routes. | Matches the read-only v0.5 scope. |
| History representation | Return bounded summaries from the history list. Return one entry's operation/report metadata and diagnostics from detail; use separately paged workspace/finding collections rather than embedding unbounded child lists. | Keeps the list bounded and supports the UX states; validate fields and representative payload sizes in OpenAPI examples. |
| Artifact access | Do not expose raw artifact bytes or arbitrary artifact paths through the browser API. Return only approved evidence references and explicit missing/integrity outcomes. | Preserves the artifact-store boundary and prevents filesystem disclosure. |
| Ordering and pagination | Order by recorded timestamp descending, then `historyEntryID` descending; use an opaque keyset cursor bound to the API version and ordering. | Gives a deterministic tie-breaker and avoids offset drift under new history inserts. Cursor encoding, expiry, deletion behavior, and numeric page limits remain unapproved. |
| Versioning | Keep `/api/v1` as the candidate version boundary; do not repeat a response-level schema version unless embedding another independently versioned contract requires it. | Avoids redundant version fields while preserving an explicit transport version. |
| Local boundary | Retain loopback-only and same-origin serving as the candidate baseline; reject unexpected Host/Origin values, do not enable CORS or trust forwarded headers, and make no external network requests. | A local listener is still reachable by local processes and hostile browser origins; exact validation, port/lifecycle, DNS-rebinding, and authentication/token decisions require security review. |
| Bounds and cancellation | Require finite request, response, page, concurrency, and deadline bounds; reject an out-of-range page-size value as `API_REQUEST_INVALID`; report a server response-limit failure with the candidate `API_RESPONSE_TOO_LARGE` without truncation. | The actual numeric limits must follow payload and platform evidence; none is selected here. |
| Errors | Use the single candidate HTTP mapping in [`requirements/ERROR-MODEL.md`](requirements/ERROR-MODEL.md); preserve typed history/API codes and never substitute an empty success. | The mappings remain proposals until the OpenAPI source, mixed-integrity cases, and compatibility crosswalk are reviewed. |
| Error envelope | Prefer only `code` and a safe `message`; omit `requestId` unless a stable correlation use case is approved. | Avoids adding a field without a demonstrated consumer contract. |

Still unresolved before an OpenAPI artifact can be accepted: exact response fields and examples, numeric bounds, cursor format/lifetime, mixed artifact-integrity responses, local authentication, port selection and process lifecycle, and the crosswalk to existing report/error schemas. The independent-review and Definition of Ready gates remain unchanged.
