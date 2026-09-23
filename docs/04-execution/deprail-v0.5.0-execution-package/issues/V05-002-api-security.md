# V05-002: Freeze the Local API and OpenAPI Security Contract
- GitHub Issue: [#396](https://github.com/geoffrey-xiao/deprail/issues/396)

- Epic: [EPIC-001 / #390](https://github.com/geoffrey-xiao/deprail/issues/390)
- Target: `v0.5.0`
- Status: Review; draft API profile in merged [PR #404](https://github.com/geoffrey-xiao/deprail/pull/404); error-mapping corrections in follow-up [PR #405](https://github.com/geoffrey-xiao/deprail/pull/405); complete OpenAPI schema, numeric bounds, listener/auth decisions, independent security review, and Definition of Ready remain pending.
- Type: decision
- Area: docs
- Priority: P0
- Risk: R3
- Owner: `@geoffrey-xiao`
- Reviewer: `@geoffrey-xiao` for planning-ticket oversight only, by explicit owner direction; not independent architecture/security approval
- Dependencies: V05-001; V05-000 / #387

## Value

A precise local API contract prevents transport code and the browser from inventing scan semantics, leaking local data, or accepting hostile requests.

## Scope

Resolve the resource and response representation required by accepted UX; endpoint set and read-only boundary; deterministic bounded pagination/cursors; OpenAPI 3.1 operations, schemas, examples, error responses, and version policy; HTTP status/error mapping; listener bind/port/startup/shutdown; Host/Origin, DNS-rebinding, CORS, CSRF, and authentication/token decisions; method/content-type/path/query validation; request/response/page/concurrency/time bounds; cancellation; redaction/logging; and CLI/shared-application compatibility. Add and validate the OpenAPI artifact and crosswalk it with existing error/report schemas.

## Out of scope

No HTTP handler, listener, API client, browser fetch implementation, generated runtime code, mutation route, browser-triggered scan, export/upload, or implementation issue.

## Inputs, outputs, and failure behavior

Inputs: accepted UX resource/state map, `API-DESIGN.md`, architecture, security/error requirements, and current versioned report contracts. Outputs: approved-ready OpenAPI 3.1 source/examples, decision record for unresolved local listener controls, and an endpoint/error/failure traceability table. Malformed, unknown, oversized, stale, unauthorized-origin, unavailable-store, corrupt-record, timeout, and cancellation outcomes remain explicit; they never become empty successful data. If a decision conflicts with existing report/CLI meaning, preserve the current meaning and record the unresolved compatibility decision.

## Required verification and evidence

Validate OpenAPI syntax and examples against the selected validator; crosswalk status/error codes and schema versions; map each operation to an accepted UI workflow and security case; document bounds and safe diagnostics. No runtime test is implied by this design issue.

## Acceptance criteria

- Every endpoint maps to a UX need and has request/response schemas, bounds, error mappings, and examples in OpenAPI 3.1.
- Bind address, port conflict, origin/Host, DNS-rebinding, CORS/CSRF/auth, lifecycle, shutdown, and redaction decisions have explicit evidence and failure behavior.
- Only approved read operations are specified; no arbitrary SQL, filesystem path, or process operation is exposed.
- Existing report completeness, operation outcome, CLI output, and error meanings remain compatible unless separately approved.
- Independent security review and runtime authorization remain outstanding until the v0.5 Definition of Ready passes.
