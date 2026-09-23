# DepRail v0.5 Architecture

**Status:** Draft; owner and independent architecture/security review required. No code authorization.

## 1. Context

v0.5 proposes a local scan-history store and embedded browser console served through a local REST API. It extends the Go single-binary core with SQLite storage and static React/TypeScript/Vite assets. Existing CLI and browser clients must use shared application services; the UI and HTTP transport do not own security or domain policy.

```text
Browser console (embedded static assets)
                 │ same-origin HTTP, reviewed local boundary
                 ▼
          Local HTTP transport ──── CLI transport
                 │                       │
                 └────────┬──────────────┘
                          ▼
                 Application services
                    │            │
                    ▼            ▼
              Domain reports  History port
                                   │
                                   ▼
                              SQLite store
                                   │
                                   ▼
                  Content-addressed raw artifacts
```

This is a context diagram, not an approved implementation design. Exact packages, API routes, process lifecycle, and data layout are design decisions.

## 2. Component responsibilities

| Component | Responsibility | Must not own |
| --- | --- | --- |
| CLI/web entry points | Translate user/HTTP input into application requests; format transport response | Domain policy or duplicated scan semantics |
| Application services | Coordinate scan/report/history use cases; validate access boundaries; map typed failures | HTTP framework or UI state |
| Domain | Deterministic report identity, report meaning, completeness, provenance | React, transport-specific types, SQL |
| Application/history service | Unique history-entry identity and operation lifecycle; define storage/retrieval use cases | Rendering or transport details |
| History port | Define operations and stable data contract for storing/retrieving accepted history entries | Rendering or transport details |
| SQLite adapter | Transactions, schema versions, query bounds, migration and recovery implementation | Scanner or policy decisions |
| Web console | Present API state and user navigation; escape untrusted text | Vulnerability matching, scan completeness calculation, secrets |
| Embedded asset server | Serve only packaged UI assets and API on approved origin | Arbitrary filesystem paths |
| Artifact store | Retain/read raw scanner evidence by verified digest | Scan history lifecycle decisions |

The Go domain remains free of Cobra, HTTP framework, SQLite/SQL, and scanner-specific types. Adapters point inward. This preserves the architecture dependency rule.

## 3. Proposed request and data flow

1. User starts the approved local-console entry point (command and browser-open behavior not yet selected).
2. The application chooses a canonical local data root and initializes/validates the history schema using the reviewed migration policy.
3. API requests are bounded, parsed, and translated into application queries; no query parameter becomes a path or SQL fragment.
4. History adapter returns a unique history entry, its source `ScanReport.ScanID`, an operation outcome, unchanged report completeness when a report exists, and stable artifact references; the application verifies applicable integrity/provenance contracts.
5. Transport returns a bounded response that keeps operation outcome (`completed|failed|cancelled`) distinct from report completeness (`complete|partial|failed`).
6. The browser presents cancellation or fatal operation failure as the top-level outcome; it never treats a report's retained completeness as proof the operation completed.
7. Shutdown stops accepting new requests, drains/cancels bounded work, closes storage cleanly, and reports lifecycle errors without exposing credentials.

Whether scans are automatically recorded, how existing reports enter history, and how the console is started are unresolved product/architecture decisions. Do not implement implicit recording or start a background server until approved.

## 4. Persistence proposal (unapproved)

SQLite is the roadmap/architecture choice for local indexed history; raw scanner artifacts remain in the existing content-addressed artifact store. SQLite should hold the smallest useful query/index metadata and a schema-validated normalized scan representation or stable reference to it, with digest and provenance linkage. The final choice requires an ADR and size/performance evidence.

The ADR must decide: canonical data directory; DB filename and permission model; WAL/journal/concurrency behavior; single-writer/multiple-reader strategy; transaction scope; report/blob size boundaries; artifact retention and dangling references; migration direction; backup/recovery; retention/deletion/export; disk-full/corruption behavior; process locking; and behavior when persistence is unavailable. No destructive migration, replacement, network sync, or implicit retention limit is approved by this draft.

## 5. Local HTTP trust boundary (unapproved)

Candidate baseline: loopback-only bind, UI and API served from one origin, no remote publishing, no cross-origin API access, bounded request/response bodies, restrictive Host/Origin checks, and no browser mutation routes in the initial proposal. These choices require independent review, platform verification, and threat analysis before implementation. Do not bind to wildcard/LAN/public interfaces, trust arbitrary forwarded headers, serve arbitrary file paths, or weaken browser-origin checks for convenience.

The API contract must define address/port selection, port conflict behavior, lifecycle, origin and DNS-rebinding defenses, CORS, CSRF relevance, authentication needs for local browser access, request limits, timeouts, cancellation, logging/redaction, and shutdown. A local listener is still a security boundary.

## 6. Frontend and asset boundary

React, TypeScript, and Vite follow the architecture baseline. The user asks for a shadcn/ui style; UX review must choose actual primitives/components and tokens. Build output is static and embedded into the binary using the repository's approved mechanism. The runtime must serve only compiled embedded assets, enforce SPA fallback only for known routes, and never join untrusted URL paths to the host filesystem. Asset version must be compatible with the API contract. Missing/invalid assets are an explicit startup or serving failure, not a blank success page.

Frontend toolchain pinning, lockfile policy, supply-chain review, code generation/ownership, bundle size and binary impact remain open decisions for the execution package.

## 7. Cross-cutting invariants

- Local-first; no default source upload, telemetry, or remote service.
- CLI remains functional with the console stopped and without history where the accepted contract permits it.
- Report completeness remains `complete`, `partial`, or `failed`; operation outcome and report status remain separate through storage, API, and UI. Cancellation does not extend the existing report schema.
- `historyEntryID` uniquely identifies a stored operation occurrence; `sourceScanID` and finding identifiers/provenance remain governed by current report contracts.
- Repository-relative paths serialize with `/`; canonical containment applies when resolving filesystem paths.
- Untrusted labels and findings are escaped in the UI; credentials and sensitive environment values are redacted.
- SQL values are parameterized; requests, queries, pages, and stored values are bounded.
- Migrations are versioned, deterministic, non-destructive by default, and recoverable.
- No browser or API route changes repository/package files or launches package-manager commands in this release.

## 8. Failure semantics

Storage unavailable, migration rejected, database corrupt/locked/full, artifact missing or digest mismatch, malformed/oversized request, timeout, request cancellation, cancelled scan operation, unsupported contract version, listener failure, and missing embedded UI assets are typed explicit failures. None becomes an empty successful scan/history response. A failure to save history must not rewrite report completeness or operation outcome; the user sees report completeness, operation outcome, and persistence outcome distinctly.

## 9. Decisions required before implementation

- Whether history capture is automatic during scan or explicit after scan; exact integration point and failed-write policy.
- Stored fields, report representation, artifact link integrity, data directory, retention/deletion and export.
- SQLite transaction/concurrency/backup/migration policy and disk/corruption recovery.
- REST endpoints, schemas, pagination, errors, versions, request limits, and status semantics.
- Local server startup, loopback/port/origin/auth/CSRF/CORS/shutdown behavior.
- UI navigation, shadcn/ui component/dependency approach, accessibility and embedded build.
- Platform matrix, bundle/startup budgets, supported browser versions, and release smoke.

Resolve material changes through an ADR/versioned plan decision, not an implementation shortcut.
