# v0.5 Functional Requirements

**Status:** Draft; acceptance requires owner and independent architecture/security review.

## FR-501 — Preserve history and report identities

Each stored history entry has a unique, stable `historyEntryID`. Preserve the source report's `ScanReport.ScanID` as `sourceScanID`; it may repeat across runs and is never used as the history primary key. Preserve existing report schema, provenance, and report completeness without rewriting it to manufacture an operation outcome.

## FR-502 — Browse bounded scan history

The local console can request and display an explicitly bounded, deterministically ordered page of saved scans. An empty page is distinguishable from loading, API failure, unavailable store, and unsupported storage schema.

## FR-503 — Inspect a saved scan

The user can select a saved history entry by its unique entry ID and inspect the approved detail set, including operation outcome, report completeness when present, workspace scope, finding summaries, provenance, and diagnostics. A missing raw artifact/digest mismatch is visibly reported and never silently replaced with fabricated or empty evidence.

## FR-504 — Preserve incomplete states

Report completeness (`complete|partial|failed`) and history operation outcome (`completed|failed|cancelled`) remain distinct in persistence, API responses, lists, details, and accessible announcements. A cancelled operation must not appear as a completed clean scan even when its returned report retains `complete`; zero findings on an incomplete report must not be represented as clean.

## FR-505 — Optional local console

Existing CLI operations do not require a browser, web server, or successful history initialization unless a separately approved contract defines otherwise. Console startup failure cannot alter existing CLI output semantics.

## FR-506 — Shared application behavior

HTTP and CLI use the same application/domain semantics. API-specific transport errors map from application/storage errors; no security or scan policy is duplicated in the React client.

## FR-507 — Safe local API

Only reviewed read operations are exposed. Inputs, page sizes, paths, and outputs are bounded and validated. Unsupported methods/fields and invalid identifiers fail explicitly without filesystem escape, SQL injection, or side-effect execution.

## FR-508 — Versioned storage lifecycle

Persistent records identify the schema version; migration behavior is explicit, deterministic, and preserves prior data on failure. The accepted history retention/deletion behavior is documented and observable.

## FR-509 — Embedded UI consistency

The shipped console is built from reviewed frontend sources, served from embedded static assets, and matched to the API contract it consumes. Missing/incompatible assets and API versions are explicit errors.

## FR-510 — Accessible interactions

Users can navigate the accepted history and detail workflows with keyboard and assistive technology. Focus, names, status announcements, contrast, narrow layouts, and reduced-motion behavior meet the reviewed UX checklist.

## FR-511 — No unsupported controls

The v0.5 UI/API does not expose hosted accounts, multi-user collaboration, remote publishing, browser-triggered scans, remediation writes, repository mutation, or agent write scopes unless an approved versioned scope decision supersedes this package.

## Requirement-to-design traceability

| Requirement | Design source | Evidence expected |
| --- | --- | --- |
| FR-501–504 | [`PRD-v0.5.md`](../PRD-v0.5.md), [`API-DESIGN.md`](../API-DESIGN.md), [`FAILURE-AND-DATA-CONTRACT.md`](FAILURE-AND-DATA-CONTRACT.md) | Storage/API contract tests and browser status-state smoke |
| FR-505–507 | [`ARCHITECTURE-v0.5.md`](../ARCHITECTURE-v0.5.md), [`API-DESIGN.md`](../API-DESIGN.md), [`SECURITY-REQUIREMENTS.md`](SECURITY-REQUIREMENTS.md) | CLI regression, API integration and hostile-request evidence |
| FR-508 | [`FAILURE-AND-DATA-CONTRACT.md`](FAILURE-AND-DATA-CONTRACT.md) | Migration/recovery and interrupted-write tests |
| FR-509–510 | [`UX-DESIGN.md`](../UX-DESIGN.md), [`ARCHITECTURE-v0.5.md`](../ARCHITECTURE-v0.5.md) | Actual packaged browser and accessibility checks |
| FR-511 | [`PRD-v0.5.md`](../PRD-v0.5.md) | API route review and UI workflow inspection |

Detailed FR-501–FR-511 design traceability is in [`CONTRACT-CROSSWALK.md`](../CONTRACT-CROSSWALK.md); planned scenario IDs and evidence requirements are in [`TEST-STRATEGY.md`](TEST-STRATEGY.md#9-requirement-and-threat-evidence-matrix). These are future acceptance requirements, not runtime pass claims.

Implementation issues must link each requirement to a frozen contract and specific observable acceptance evidence. No requirement alone authorizes a new capability.
