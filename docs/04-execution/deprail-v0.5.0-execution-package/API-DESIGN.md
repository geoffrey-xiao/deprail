# v0.5 Local API Design Contract

**Status:** Draft proposal; no OpenAPI schema or runtime API is approved.
**Design inputs:** [`PRD-v0.5.md`](PRD-v0.5.md), [`UX-DESIGN.md`](UX-DESIGN.md), [`ARCHITECTURE-v0.5.md`](ARCHITECTURE-v0.5.md).

## 1. Goal and constraints

Define the smallest same-origin local API needed by the approved history and scan-detail UX. The API delegates to shared application services and returns versioned domain meaning. It does not reimplement scan logic, become a general extension platform, or imply team/cloud functionality. The proposed initial surface is read-only; browser-triggered scanning or mutation is excluded unless explicitly reconsidered through release change control.

An OpenAPI 3.1 document with JSON Schema-compatible examples and validation is a required design-stage deliverable before API implementation. This Markdown proposal is not a substitute for that artifact.

## 2. Proposed resource model (unapproved)

- **ScanSummary:** stable scan ID, recorded time, root display label under privacy contract, outcome (`complete|partial|failed|cancelled`), finding/workspace counts, and minimal provenance summary.
- **ScanDetail:** summary plus scan contract/version, diagnostics, workspaces, and stable links/identifiers for findings and artifacts.
- **WorkspaceSummary:** stable workspace ID, repository-relative path, ecosystem/manager metadata permitted by the data contract, and workspace completeness.
- **FindingSummary:** stable finding key, package identity, aliases, severity/source as available, affected workspace/path context, fixed-version data, and evidence references subject to existing report schema.
- **ApiError:** stable API error code, concise safe message, retryability only when meaningful, request/correlation ID if approved, and no stack trace/secrets/raw paths by default.

Existing report schemas govern nested finding/report meaning. Do not invent new vulnerability semantics. Exact response representation (full embedded report versus paged child resources) is an open decision that must be resolved against size, UX, and compatibility evidence.

## 3. Proposed endpoint matrix (for design review)

| Method and path | Purpose | Proposed result | Write? |
| --- | --- | --- | --- |
| `GET /api/v1/health` | Local service readiness/schema version | Bounded readiness response; no environment dump | No |
| `GET /api/v1/scans` | List saved history | Ordered page of ScanSummary plus opaque continuation token | No |
| `GET /api/v1/scans/{scanID}` | Load one scan detail | ScanDetail or typed not-found/integrity error | No |
| `GET /api/v1/scans/{scanID}/workspaces` | Load bounded workspace list if not embedded | Ordered page of WorkspaceSummary | No |
| `GET /api/v1/scans/{scanID}/findings` | Load bounded finding list if not embedded | Ordered page of FindingSummary | No |

Only include separate child collection endpoints if UX/resource size needs justify them. Avoid parallel duplicate ways to fetch the same data. No `POST /scan`, delete, export/upload, policy mutation, remediation, approval, or exception endpoint is proposed.

## 4. Query and deterministic ordering proposal

- List order should be explicit and stable (candidate: recorded timestamp descending, stable scan ID as tie-breaker); final semantics require approval.
- Pagination must be bounded, opaque to clients, and stable under concurrent insertion. Prefer a cursor over unbounded offset paging if the chosen storage/index contract supports it.
- Maximum page size, filter allowlist, cursor version/expiry, and behavior when records are deleted between pages are open contract fields; never accept arbitrary SQL/order expressions.
- Query values and path identifiers are validated. Unknown filters/fields fail clearly rather than being silently ignored.

## 5. Response and error semantics

- Responses carry an explicit API/schema version according to the reviewed OpenAPI design; avoid redundant ad hoc version markers if URL versioning is sufficient.
- Existing complete/partial/failed/cancelled scan states remain unchanged end-to-end.
- Empty scan history is a successful empty page. Store unavailable, migration rejected, unsupported version, corrupt row, artifact digest mismatch, or timeout are explicit non-success responses.
- A missing scan resource is distinct from a scan whose report has zero findings.
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
