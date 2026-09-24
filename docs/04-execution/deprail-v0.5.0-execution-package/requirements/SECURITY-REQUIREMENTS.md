# v0.5 Security Requirements

**Status:** Draft threat checklist; independent architecture/security review is mandatory.

## 1. Assets and trust boundaries

Assets include local scan history, raw scanner artifacts, provenance/digests, findings, repository paths/labels, SQLite schema and migrations, embedded UI assets, and the local API. Trust boundaries are the browser-to-local-HTTP connection, API-to-application translation, application-to-SQLite/artifact store, and process/filesystem boundary. Repository-derived reports, database rows, query parameters, headers, and browser requests are untrusted.

## 2. Threats and required controls

| ID | Threat | Required control and evidence reference |
| --- | --- | --- |
| SEC-01 | Remote website reaches loopback API (CSRF/DNS rebinding) | Strict Host/Origin checks, same-origin policy, no permissive CORS; hostile-origin browser/API scenario and sanitized request/response evidence. |
| SEC-02 | Other local process reads history | Explicit local trust model, canonical data root, restrictive owner-only permissions; POSIX mode/Windows ACL and unauthorized-user evidence. |
| SEC-03 | Path traversal via URL, encoded separator, symlink, or asset route | Canonical URL/path validation; embedded assets only; reject arbitrary filesystem parameters; hostile-path corpus on each OS. |
| SEC-04 | SQL injection or expensive query | Parameterized SQL, typed filters, allowlisted sort, bounded cursor/page/query work; SQL metacharacter and unknown-filter tests with DB unchanged. |
| SEC-05 | Oversized body/result or request flood | Approved finite request/response/page/concurrency/deadline bounds; boundary and concurrency evidence without truncation or resource exhaustion. |
| SEC-06 | XSS via repository labels/findings/diagnostics | Text rendering/escaping and safe URL policy; hostile Unicode/markup fixture in the actual browser, with no script execution. |
| SEC-07 | Credential/path leakage from source scan/project/error fields | Versioned allowlisted `history-v1` projection, safe typed diagnostics, and response allowlist; never persist or return absolute `ScanReport.repository_identity.root`, raw errors, credential URLs, full environments, source, SQL/stack traces, or unreviewed unknown fields. Exercise actual nonempty report and hostile path/diagnostic fixtures. |
| SEC-08 | Malicious/corrupt SQLite projection contents | Validate projection schema/version, repeated row/JSON fields and digest references at write/read boundaries; unsupported versions and corrupt/mismatched rows fail explicitly while preserving DB bytes without auto-repair/delete. |
| SEC-09 | Artifact substitution or loss | Verify digest/provenance; missing/mismatched artifact remains explicit and is never trusted, fabricated, or deleted. |
| SEC-10 | UI/API contract skew | Versioned contract and visible incompatibility failure; packaged asset/API mismatch scenario. |
| SEC-11 | Local server accidentally exposed to LAN/public | Loopback-only default and no unsafe fallback; inspect actual bound address on Linux/macOS/Windows. |
| SEC-12 | Static asset disclosure | Embedded allowlisted assets only, no filesystem fallback; traversal and host-file canary tests. |
| SEC-13 | Insecure dependency/build chain | Pin and review dependencies/toolchain, license/SBOM, and reproducible build evidence before package selection. |

Threat IDs link to the matching scenario and release-evidence contract in [`TEST-STRATEGY.md`](TEST-STRATEGY.md#9-requirement-and-threat-evidence-matrix). They are planned controls and evidence, not claims that security testing has passed.

## 3. API minimum controls (proposed)

- Allow only approved HTTP methods/paths; read-only candidate API has no state-changing routes.
- Validate Host, Origin, content type, identifier shape, pagination, and all enum values before storage access.
- No arbitrary path or SQL fragments; no shell/process execution from API input.
- Set finite request, response, query, concurrency, and lifetime limits; propagate cancellation.
- Return safe typed diagnostics; log only bounded metadata and stable IDs, not secrets or raw source.
- Avoid credentials in URLs and unrestricted browser storage. If a local bearer token is required, define generation, transport, storage, entropy, rotation, and redaction before implementation.

## 4. SQLite minimum controls (proposed)

- Canonical data root outside the scanned repository by default, restrictive permissions where supported, explicit user-facing location.
- Parameterized SQL and fixed schema/migration registry; no untrusted extension loading.
- Transactions for multi-record ingest; safe concurrent access and lock handling; disk-full and corruption paths.
- No silent schema downgrade, destructive repair, DB replacement, or artifact deletion.
- Define retention/deletion/backup and ensure raw artifact cleanup respects references and integrity.

## 5. Privacy

Default to local-only processing. No telemetry, source upload, remote history, external fonts/CDNs, or third-party analytics. Any later network behavior requires a separate explicit permission/scope and privacy review. UI diagnostics must redact credential-bearing URLs, user info, secrets, and sensitive environment values. Persist only information required for the approved user workflow.

## 6. Security acceptance before implementation

- [ ] Named independent reviewer examines API listener/origin and browser attack model.
- [ ] Data classification and exact persisted/exposed fields are reviewed.
- [x] Request/path/SQL/XSS/response-size threat cases are mapped to planned `SEC-*` verification IDs in [`TEST-STRATEGY.md`](TEST-STRATEGY.md#9-requirement-and-threat-evidence-matrix); actual security execution and review remain outstanding.
- [ ] Migration/backup/corruption/disk-full controls preserve data and permissions.
- [ ] Artifact digest and provenance verification behavior is explicit.
- [ ] Cross-platform bind, path, permission, and shutdown behavior is specified.
- [ ] No source upload, remote service, scanner install, or repository mutation enters scope implicitly.
- [ ] Human security review and residual-risk dispositions are recorded separately from owner approval.
