# v0.5 Security Requirements

**Status:** Draft threat checklist; owner security assessment and evidence remain pending. Independent review is optional under [ADR-0005](../../../adr/ADR-0005-solo-owner-review-policy.md).

## 1. Assets and trust boundaries

Assets include local scan history, raw scanner artifacts, provenance/digests, findings, repository paths/labels, SQLite schema and migrations, embedded UI assets, and the local API. Trust boundaries are the browser-to-local-HTTP connection, API-to-application translation, application-to-SQLite/artifact store, and process/filesystem boundary. Repository-derived reports, database rows, query parameters, headers, and browser requests are untrusted.

## 2. Threats and required controls

| ID | Threat | Required control and evidence reference |
| --- | --- | --- |
| SEC-01 | Remote website reaches loopback API (CSRF/DNS rebinding) | Candidate requires exact Host and same-origin checks, rejects hostile Origin/Fetch Metadata, requires bearer auth on every API route, and emits no permissive CORS; exercise hostile origins with valid/missing/invalid credentials and retain sanitized evidence. |
| SEC-02 | Other local process reads history | Explicit local trust model and canonical data root with restrictive owner-only permissions. The process-scoped bearer token limits browser API access but does not isolate against a hostile process running as the same user; validate POSIX mode/Windows ACL and unauthorized-user behavior. |
| SEC-03 | Path traversal via URL, encoded separator, symlink, or asset route | Canonical URL/path validation; embedded assets only; reject arbitrary filesystem parameters; hostile-path corpus on each OS. |
| SEC-04 | SQL injection or expensive query | Candidate exposes no filters or arbitrary sort; use parameterized fixed queries, bounded cursor/page work, and deterministic order. Exercise SQL metacharacters and unknown query parameters with DB unchanged. |
| SEC-05 | Oversized body/result or request flood | Candidate request-target/header/page/concurrency/deadline bounds are explicit. The response cap counts fully serialized UTF-8 JSON bytes, not schema code points; schema-valid values may exceed it and must fail before sending a partial body. Test maximum Unicode/escaping cases, exact byte boundaries, concurrency, and slow requests without truncation or resource exhaustion; exact limits remain unapproved. |
| SEC-06 | XSS via repository labels/findings/diagnostics | Text rendering/escaping and safe URL policy; hostile Unicode/markup fixture in the actual browser, with no script execution. |
| SEC-07 | Credential/path leakage from source scan/project/error fields | Versioned allowlisted `history-v1` projection, safe typed diagnostics, and response allowlist; never persist or return absolute `ScanReport.repository_identity.root`, raw errors, credential URLs, full environments, source, SQL/stack traces, or unreviewed unknown fields. Exercise actual nonempty report and hostile path/diagnostic fixtures. |
| SEC-08 | Malicious/corrupt SQLite projection contents | Validate projection schema/version, repeated row/JSON fields and digest references at write/read boundaries; unsupported versions and corrupt/mismatched rows fail explicitly while preserving DB bytes without auto-repair/delete. |
| SEC-09 | Artifact substitution or loss | Verify digest/provenance; missing/mismatched artifact remains explicit and is never trusted, fabricated, or deleted. |
| SEC-10 | UI/API contract skew | Versioned contract and visible incompatibility failure; packaged asset/API mismatch scenario. |
| SEC-11 | Local server accidentally exposed to LAN/public | Loopback-only default and no unsafe fallback; inspect actual bound address on Linux/macOS/Windows. |
| SEC-12 | Static asset disclosure | Embedded allowlisted assets only, no filesystem fallback; traversal and host-file canary tests. |
| SEC-13 | Insecure dependency/build chain | Pin and review dependencies/toolchain, license/SBOM, and reproducible build evidence before package selection. |

Threat IDs link to the matching scenario and release-evidence contract in [`TEST-STRATEGY.md`](TEST-STRATEGY.md#9-requirement-and-threat-evidence-matrix). They are planned controls and evidence, not claims that security testing has passed.

## 3. API minimum controls (candidate for owner review)

- Bind only to `127.0.0.1` on an OS-selected ephemeral port; never fall back to wildcard/LAN/public binding. Run only within an explicit foreground local-console command, not as a scan side effect or background service; stop accepting on shutdown and drain for at most the proposed 10 seconds.
- Expose only the five candidate `GET /api/v1` routes. Reject all request bodies, unsupported methods, unknown query parameters, malformed IDs/cursors, and page sizes outside 1–50 before storage access. No state-changing route is specified.
- Require exact `Host` for the selected loopback address/port. If `Origin` is present, require exact same-origin and reject `null`/mismatch; if `Sec-Fetch-Site` is present, require `same-origin`. Emit no CORS allow headers and trust no forwarded headers.
- Require an opaque cryptographically random process-scoped 256-bit bearer token on every API route, including health. Accept it only in the `Authorization` header; revoke it when the process ends; keep it only in browser memory. Candidate bootstrap passes it in the initial URL fragment, which the client must clear with `history.replaceState` before any request. A reload/new tab loses the token and must show safe reopen guidance from the active CLI session, with no cookie/storage fallback. The token is never sent as a request target, referrer, or persisted browser value. Fragment exposure in browser/OS launch state is a material unresolved owner security decision.
- Candidate finite bounds are: request target 2,048 bytes (transport 414 over limit); headers 8,192 bytes (transport 431 over limit); eight concurrent requests; 10-second request deadline and shutdown drain; 1 MiB response; page size 1–50 (default 25); cursor at most 512 characters. These are design proposals, not accepted performance limits. Parser-boundary failures may not use the API error envelope; route-level failures are typed. Never truncate success or queue unbounded work.
- No arbitrary paths, SQL, shell/process execution, artifact bytes, or filesystem paths; use parameterized bounded reads and propagate disconnect cancellation. No external network access.
- Redact authorization headers/tokens, sensitive request headers, raw paths, repository data, SQL/driver details, stack traces, and full environment from logs, diagnostics, and responses.

## 4. SQLite minimum controls (proposed)

- Canonical data root outside the scanned repository by default, restrictive permissions where supported, explicit user-facing location.
- Parameterized SQL and fixed schema/migration registry; no untrusted extension loading.
- Transactions for multi-record ingest; safe concurrent access and lock handling; disk-full and corruption paths.
- No silent schema downgrade, destructive repair, DB replacement, or artifact deletion.
- Define retention/deletion/backup and ensure raw artifact cleanup respects references and integrity.

## 5. Privacy

Default to local-only processing. No telemetry, source upload, remote history, external fonts/CDNs, or third-party analytics. Any later network behavior requires a separate explicit permission/scope and privacy review. UI diagnostics must redact credential-bearing URLs, user info, secrets, and sensitive environment values. Persist only information required for the approved user workflow.

## 6. Security acceptance before implementation

- [ ] Owner examines API listener/origin and browser attack model and records residual risk; independent review is optional.
- [ ] Data classification and exact persisted/exposed fields are reviewed.
- [x] Request/path/SQL/XSS/response-size threat cases are mapped to planned `SEC-*` verification IDs in [`TEST-STRATEGY.md`](TEST-STRATEGY.md#9-requirement-and-threat-evidence-matrix); actual security execution and review remain outstanding.
- [ ] Migration/backup/corruption/disk-full controls preserve data and permissions.
- [ ] Artifact digest and provenance verification behavior is explicit.
- [ ] Cross-platform bind, path, permission, and shutdown behavior is specified.
- [ ] No source upload, remote service, scanner install, or repository mutation enters scope implicitly.
- [ ] Owner records security assessment and residual-risk dispositions separately from technical owner acceptance; the same owner may complete both.

## 7. Proposed console and transport details

The selected candidate is `deprail web`, a foreground local-console command. It prints only the non-secret loopback URL (`http://127.0.0.1:<port>/console/`), never the bearer fragment/token or credential-bearing URL. It does not launch a browser by default. Optional `--open` passes the exact fragment URL to the OS launcher by argv, never a shell. In a stdin TTY, Enter reopens that active process URL without printing the credential; a launcher failure warns and permits an explicit Enter retry. Non-TTY stdin without `--open` is a usage error (exit `2` before startup); a failed non-TTY initial launch stops the listener and exits `3`, not an unusable success. Listener failure exits `3`, normal signal shutdown exits `0`, and invalid/inaccessible supplied artifact-root is configuration exit `2`. Token is 256-bit, process-scoped, revoked at exit, removed before rendering/request and held only in browser memory. Reload recovery uses Enter in the active TTY or restart with `--open` for noninteractive invocation. Browser/OS fragment exposure remains an explicit residual owner risk decision.

An optional CLI-only `--artifact-root <path>` accepts an existing trusted artifact-store directory; it does not walk/scan a repository. Canonicalize the directory at startup and resolve only validated digest-addressed artifacts beneath it, refusing symlink/reparse escape. Invalid/inaccessible supplied roots cause safe startup failure. Without the option, artifact references are `unavailable`. Never persist/infer the root, expose it in API/UI/errors/logs, accept a path from API, or return bytes.

Static `/console/` and explicitly embedded asset paths are separate from API routes, require exact loopback Host but not API bearer auth, and never fall back to filesystem, directory listing, arbitrary path or SPA catch-all. Unknown paths return bounded non-reflecting 404. Entry document CSP is exactly `default-src 'none'; script-src 'self'; style-src 'self'; img-src 'self' data:; font-src 'self'; connect-src 'self'; base-uri 'none'; form-action 'none'; frame-ancestors 'none'; object-src 'none'`; use `Cache-Control: no-store` and `Referrer-Policy: no-referrer`. Embedded assets only; repository strings are data, never executable; no external assets/network. Apply Fetch Metadata same-origin enforcement to API requests only, since legitimate initial top-level navigation can carry `Sec-Fetch-Site: none`.

Exact Host is `127.0.0.1:<selected-port>`; reject aliases and forwarded headers. API `Origin` exact-match and `Sec-Fetch-Site: same-origin` checks apply when present; reject null/mismatch. No CORS allow headers. All five API routes, including health, require bearer auth.

Request processing applies parser limits first (request target 2,048 bytes → transport 414; aggregate headers 8,192 bytes → transport 431, neither promises API JSON), then exact Host, framing/body rejection, method, route/version, query/ID validation, origin/fetch policy, authentication, and bounded admission. Bodies are prohibited; invalid body/framing is 400, unsupported method 405, hostile origin 403, invalid credentials 401. Admission permits at most eight in-flight requests and has no waiting queue; saturation is 503 `API_BUSY`. The 10-second deadline covers work and serialization; timeout is 504 `API_TIMEOUT` where possible. Disconnect cancels storage work and sends no fabricated response.

Cursor is canonical unpadded base64url, at most 512 bytes, with bytes defined in API-DESIGN §12: version/resource bytes, optional parent UUID, then tuple. History uses signed int64 UTC microseconds + UUID; findings use raw 32-byte stable key; workspaces use zero-based uint32 ordinal + SHA-256 of workspaceID UTF-8, resolved against immutable validated ordering and hash-checked before continuing. Integers are big-endian; UUID bytes follow canonical textual order; exact decoded sizes are 26/54/50 bytes for history/workspaces/findings. It contains no workspace/repository strings or secret and is not authorization. Malformed, noncanonical, wrong resource/version/parent, invalid ordinal or hash mismatch is 400 `API_REQUEST_INVALID`. Whole-record paging advances only through emitted records; an empty first page is terminal only when no records exist. No skip on short byte-budget pages.

The route error vocabulary is exactly OpenAPI's existing categories: 400 `API_REQUEST_INVALID`; 401 `API_AUTH_UNAUTHORIZED`; 403 `API_ORIGIN_REJECTED`; 404 `HISTORY_ENTRY_NOT_FOUND`, `API_ROUTE_NOT_FOUND`, `API_VERSION_UNSUPPORTED`; 405 `API_METHOD_UNSUPPORTED`; 500 `API_RESPONSE_TOO_LARGE`, `HISTORY_CORRUPT`, `API_INTERNAL_ERROR`; 503 `API_BUSY`, `HISTORY_UNAVAILABLE`, `HISTORY_SCHEMA_UNSUPPORTED`, `HISTORY_MIGRATION_FAILED`; 504 `API_TIMEOUT`. Parser 414/431 are transport exceptions. Errors are safe and never expose SQL, driver detail, paths, stack, raw errors, token, or echoed input.

These engineering details supersede earlier draft wording that left launch/cursor mechanics unresolved, but remain proposed pending owner acceptance and implementation evidence.

