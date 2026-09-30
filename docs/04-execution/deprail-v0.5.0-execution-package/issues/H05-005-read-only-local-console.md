# H05-005: Implement Read-Only Local Console Transport

## Planning metadata

- Epic: [EPIC-002 — Local History Delivery](../epics/EPIC-002-local-history-delivery.md)
- Type: `feature`
- Area: `cli`
- Priority: `P0`
- Risk: `R3`
- Target version: `v0.5.0`
- Milestone: [`v0.5.0`](https://github.com/geoffrey-xiao/deprail/milestone/11)
- Sprint: Unassigned
- Owner: `@geoffrey-xiao`
- Reviewer: `@geoffrey-xiao` (owner); independent review is not a gate under ADR-0005.
- Dependencies: H05-003 same-operation capture/query service; accepted API, launch and security contracts and prerequisite dispositions recorded in [EPIC-002](../epics/EPIC-002-local-history-delivery.md). H05-007 may provide reviewed production assets; H05-005 smoke may use temporary contract assets only.
- Blocked reason: None (Todo, not started); global timeout disposition and issue-specific pre-start gates apply.

## Definition of Ready

- [x] Value and user impact are stated by accepted proposal [PR #417](https://github.com/geoffrey-xiao/deprail/pull/417) and owner decision [#417](https://github.com/geoffrey-xiao/deprail/pull/417#issuecomment-5907670878).
- [x] Scope and explicit exclusions are specified in [API design §11–13](../API-DESIGN.md#11-proposed-console-launch-and-static-service-contract).
- [x] Inputs, outputs, and failure behavior are specified by API design and [OpenAPI](../../../../schemas/openapi/v1/openapi.yaml).
- [x] Required behavioral and security scenarios are identified below and in [test strategy](../requirements/TEST-STRATEGY.md).
- [x] Observable contract acceptance is specified; runtime acceptance remains unchecked below.
- [x] Owner reviewer is assigned; independent review is not required.
- [x] Dependencies and target version are recorded; prerequisite disposition remains a pre-start gate.

## Goal

Provide an explicitly invoked, foreground-only, read-only local console transport over the shared H05-003 application service. Users can inspect captured history without changing scan behavior or creating/migrating storage.

## Scope

Own `cmd/deprail` foreground `web` command and `internal/transport/localhttp` HTTP/static adapter. Serve exactly the five authenticated GET operations in [OpenAPI](../../../../schemas/openapi/v1/openapi.yaml) and static `/console/`, `/console/history`, `/console/scans/{historyEntryID}` (canonical UUIDv4), `/console/about`, plus declared hashed embedded assets. Use app queries only. Bind `127.0.0.1:0`; exact Host, Origin/API Fetch Metadata, process-scoped 256-bit bearer, memory-only fragment bootstrap, exact CSP and no-store/no-referrer follow API §§11–12. Print only bare origin/path. Non-TTY requires --open; failed non-TTY initial launch exits3 and stops listener; TTY Enter recovery and safe warnings follow the accepted contract.

Transport owns API envelope serialization and exact binary/base64url cursor codec over the app's continuation values. Enforce all query/parser/auth/method/body/byte/deadline/admission/shutdown bounds and parent/resource/version/canonical cursor validation. Consume an allowlisted embedded asset FS from H05-007; temporary transport-test assets are never production assets or a shipped fallback.

Traceability: [product baseline](../../../01-product/deprail-product-design-v1-ai.md), [architecture baseline](../../../02-architecture/deprail-architecture-and-tech-stack-v1.md), [FR-502–507/511](../requirements/FUNCTIONAL-REQUIREMENTS.md), API design, security SEC-01–05/07/11–12, STORE-01 and [release plan §14](../../../03-planning/deprail-development-plan-v0.5.0.md#14-detailed-contract-reconciliation-for-current-owner-review).

## Out of Scope

No scan, mutation, delete, export, upload, arbitrary filesystem serving, filesystem fallback, background service, CORS, token persistence, frontend production UI, or history initialization/migration. Temporary assets are throwaway smoke inputs and MUST NOT ship. Do not alter existing scan signature, CLI JSON, output, or exit semantics.

## Inputs, Outputs, and Failure Behavior

Inputs: H05-003 query service, validated history storage, accepted route/security contract, optional existing trusted artifact root, and H05-007 production assets when available. Outputs: bounded typed responses and allowlisted static content; one safe startup origin line without credentials. Missing store is read-only empty; corrupt/incompatible store is typed failure, never empty-success. Invalid root fails safely before startup; artifact read failures become `unavailable`. Invalid host/origin/auth/input, unsupported routes/methods, request bounds, saturation, timeout, oversize response and listener/launcher failures follow exact API/CLI status semantics. Stop accepting work on signal, cancel requests, drain at most ten seconds, close storage and revoke token. No request can initiate a write or scan.

## Required Tests

- Exercise every one of five GETs against real application/storage behavior; absent, corrupt and incompatible stores remain distinct from empty history.
- Runtime test wrong Host/Origin, absent/malformed bearer auth, disallowed method/body/query/path and prove no storage mutation or scan; verify fragment is removed before the first request and token is absent from logs/stdout.
- Test UTF-8 byte boundary and one-byte-over, whole-record short pages, no skipped row, canonical cursor rejection, altered resource/version/tuple, and child cursor used for another parent.
- Exercise request-target/header caps, eight-slot admission, deadline cancellation, bounded shutdown and listener failure; confirm CLI remains usable.
- Test exact static route and asset allowlist, CSP, no filesystem fallback, symlink/reparse escape rejection for artifact root, and recovery/reopen behavior in actual browser.

## Acceptance Criteria

- [ ] Only the five documented authenticated GETs and exact static allowlist are served; no mutation route, scan action, filesystem fallback, or storage initialization exists.
- [ ] Host/origin/auth/CSP/token lifecycle, foreground launch and non-TTY/recovery behavior match the contract; credentials never print or persist.
- [ ] Complete serialized responses never exceed 1,048,576 UTF-8 bytes; pagination is whole-record and cursor/parent validation prevents skips or cross-resource reuse.
- [ ] Parser, admission, deadline, shutdown, trusted-root, typed error and failure cases are bounded and observable as specified.
- [ ] Existing CLI output/exit behavior and history bytes remain unchanged; all behavioral scenarios above pass against the running transport.

## Owner Review

Owner reviews each observable criterion, security boundary and residual fragment exposure against runtime evidence. This contract's planning DoR does not assert runtime implementation, technical/security acceptance or release approval. No independent reviewer gate.

## Evidence Required

- Verification commands or scenarios: future targeted transport behavioral tests; launched CLI smoke with actual local API and browser; malicious Host/Origin/token/cursor probes; paging boundary and shutdown/recovery scenarios; existing CLI compatibility scenario.
- Expected artifacts, logs, screenshots, or links: sanitized request/response assertions, browser recovery evidence, bounded byte measurements and platform launch results; redact tokens, roots and repository data.

## Final Acceptance

- [ ] Owner reviewed every acceptance criterion during PR review.
- [ ] Required verification and CI results were reviewed.
- [ ] Owner review and remaining risk are recorded.
- [ ] Evidence links are attached.
- [ ] Owner review and merge evidence are linked.

## Rollback

Disable/remove the console command and listener integration; preserve database, WAL, artifacts and existing CLI. No automatic store downgrade, reset or artifact deletion.
