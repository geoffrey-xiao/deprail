# DepRail v0.5 Product Requirements

| Attribute | Value |
| --- | --- |
| Version | v0.5.0 |
| Status | Planning draft; owner and independent architecture/security review required |
| Roadmap outcome | Local web application and scan history |
| Release plan | [`../../03-planning/deprail-development-plan-v0.5.0.md`](../../03-planning/deprail-development-plan-v0.5.0.md) |
| Preparation issue | [#387](https://github.com/geoffrey-xiao/deprail/issues/387) |

## 1. Product outcome

A local DepRail user can review retained completed scan history in an embedded browser console, open a selected scan, and understand its findings, workspaces, completeness, and provenance. The browser is a presentation client over shared application services. Existing CLI scan behavior remains useful without opening the browser or depending on history storage unless an approved compatibility decision changes that boundary.

The history is local-only. A scan that is partial, failed, cancelled, or not durably recorded must not be presented as a complete, successful historical result. The v0.5 proposal does not trigger scans or mutate repositories from the browser.

## 2. Source reconciliation

- Product design defines local-first operation, scan history, an embedded React local console, SQLite local storage, and the web layer as a collaboration/presentation layer.
- Architecture defines a Go core, inward dependencies, REST/OpenAPI 3.1, React/TypeScript/Vite embedded assets, SQLite with migration and backup discipline, and the local mode as one binary with optional browser and private data directory.
- The roadmap assigns local web and history to v0.5; team services, identity/RBAC/PostgreSQL are v0.6.
- The v0.5 development plan selects only this local outcome and makes UX, API, persistence design and independent review prerequisites to implementation.
- The v0.4.0-preview.2 closeout issue #356 remains open. This PRD does not close or waive its review, rollback, Python/Java, platform-smoke, or artifact follow-ups.

## 3. User jobs

- Find a previous scan without reopening its repository or rerunning the scanner.
- Distinguish complete results from partial, failed, cancelled, or unavailable history.
- Inspect findings and workspace context while preserving scanner/tool/database provenance.
- Understand when history is unavailable, migration failed, or a stored artifact cannot be read.
- Continue using CLI scan and machine-readable output independently of the web console.

## 4. Proposed scope

- Persist an explicitly selected set of normalized scan reports and stable references to existing raw artifacts in local SQLite-backed history.
- List and inspect saved scans with bounded, deterministic pagination/ordering and stable scan identity.
- Provide an embedded local web console for history and scan details through the same application services used by CLI/core.
- Provide a minimal local REST API for only the console workflows accepted by UX; read-only in v0.5 unless a reviewed decision justifies a narrow mutation.
- Preserve scan completeness, diagnostics, workspace and finding relationships, and provenance through persistence, API, and UI.
- Version local storage/API contracts, document data location and lifecycle, and provide tested non-destructive migration/recovery semantics.
- Package frontend static assets with the Go application and verify representative local workflows on supported platforms.

Exact records, endpoints, retention, UI routes, launch UX, and dependencies are open design decisions. This PRD does not define a schema or authorize their implementation.

## 5. Explicit exclusions

- Hosted or team service; remote publishing or source upload; PostgreSQL; user identity, authentication for multiple users, RBAC, shared history, notifications, or audit workflows.
- Browser-triggered scanning, changing scan policy, scan deletion/annotation, exceptions, remediation, package-manager execution, repository mutation, patch creation, commits, pushes, pull requests, or merge.
- MCP or agent write access, autonomous approval, new scanners, new vulnerability databases, and scanner semantic changes.
- Web-specific reimplementation of scanner, normalizer, policy, completeness, or remediation logic.
- Changes to existing CLI and stable JSON contracts unless a separate reviewed compatibility decision authorizes them.
- Telemetry, analytics, an extension/plugin framework, full-screen TUI, or capability expansion beyond the roadmap outcome.

## 6. UX and design-system requirement

The console must follow a reviewed UX package before UI code. User preference is shadcn/ui-inspired: composed accessible primitives, clear typographic hierarchy, quiet surfaces, consistent tokens, responsive layouts, and deliberate empty/loading/error states. The design gate must select actual component source/library, package versions, code ownership, icon and font policy, and integration/build strategy; “shadcn style” alone does not approve adding dependencies or copying generated components.

The UI must preserve meaning without color, support keyboard navigation and visible focus, expose status changes to assistive technology, respect reduced motion, and render repository-controlled values as untrusted text. Refer to [`UX-DESIGN.md`](UX-DESIGN.md).

## 7. API and persistence gate

No endpoint or database schema is approved by this PRD. [`API-DESIGN.md`](API-DESIGN.md) proposes an intentionally small read-oriented surface for review. Before implementation, reviewers must approve a versioned OpenAPI contract, response/error mapping, bind and browser-origin controls, request limits, cancellation and compatibility semantics. A SQLite ADR must define schema, migrations, transaction/concurrency model, retention/deletion, artifact relationships, disk-full/corruption recovery, permissions, and data export/removal behavior.

## 8. Observable acceptance criteria (proposed)

- AC-501: Given saved scans, the console lists them in deterministic order with identity, time, status, and enough summary information to distinguish a complete scan from an incomplete one.
- AC-502: Selecting a saved scan shows its findings/workspaces and preserves accepted provenance; unavailable artifacts are reported explicitly rather than replaced by empty data.
- AC-503: A partial, failed, or cancelled scan remains visibly distinct at the list and detail surfaces; zero findings on an incomplete scan are never described as a clean scan.
- AC-504: An empty history has an intentional empty state distinct from API/storage failure.
- AC-505: Storage or migration failure is typed and actionable; no failed write is presented as a durable history record.
- AC-506: Oversized, malformed, unsupported, and stale API inputs fail without unintended filesystem access or data disclosure.
- AC-507: Local API is not exposed beyond the approved listener/origin boundary; browser-origin abuse and path traversal attempts fail safely.
- AC-508: The CLI continues to scan when the console is unused; CLI schemas and exit codes remain compatible absent an approved exception.
- AC-509: UI workflows satisfy reviewed keyboard, focus, contrast, screen-reader, responsive, and reduced-motion expectations on the actual surface.
- AC-510: Migrations from supported prior data versions either preserve usability or stop before destructive modification with a recoverable diagnostic.
- AC-511: The packaged local history workflow has evidence on Linux, macOS, and Windows, with semantically equivalent results.

The execution package and implementation issues must refine each criterion into deterministic test/manual evidence after design decisions are approved.

## 9. Failure behavior

Storage open/lock/read/write errors, disk exhaustion, corruption, incompatible schema, artifact missing/digest mismatch, API timeout/cancellation, malformed input, oversized input, listener failure, missing embedded assets, and client/server schema mismatch remain explicit errors. They must not be represented as empty history or complete zero-finding reports. See [`requirements/FAILURE-AND-DATA-CONTRACT.md`](requirements/FAILURE-AND-DATA-CONTRACT.md).

## 10. Compatibility and release

v0.5 adds optional local history and web delivery; existing CLI operation and serialized v1alpha scan report meaning remain stable. Any incompatible schema or CLI/API change needs a separate decision and migration evidence. Preview release approval requires cross-platform binary and browser smoke, local-only/security evidence, migration and failure coverage, actual artifact verification, and separate owner/security review. A passing unit suite alone is insufficient.

## 11. Required human decisions

- Approve exact scope, especially whether scan history is automatic or an explicit operation and whether any console operation writes state.
- Approve API version/path/resource and error contract, local listener/origin controls, and browser security model.
- Approve storage content, artifact retention, deletion/export, migration, backup/recovery, and permission model.
- Approve shadcn/ui implementation strategy, dependency policy, accessibility acceptance, and embedded build approach.
- Name an independent architecture/security reviewer and disposition the open #356 predecessor gaps.
