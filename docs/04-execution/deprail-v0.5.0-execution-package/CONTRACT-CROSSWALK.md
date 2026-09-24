# v0.5 Contract Crosswalk

**Status:** Draft traceability; final issue mapping follows design review.

| Source contract | v0.5 interpretation | Execution-package artifact | Approval/evidence gate |
| --- | --- | --- | --- |
| Product design §Web Product Design | Local console presents scan history/findings; not a second scanner | [`PRD-v0.5.md`](PRD-v0.5.md), [`UX-DESIGN.md`](UX-DESIGN.md) | Approved user workflows and visible state design |
| Product design §§Local first, Evidence first, Deterministic core | Keep reports/provenance local and domain meaning transport-independent | [`ARCHITECTURE-v0.5.md`](ARCHITECTURE-v0.5.md), [`SECURITY-REQUIREMENTS.md`](requirements/SECURITY-REQUIREMENTS.md) | Data classification, provenance roundtrip, no-upload evidence |
| Architecture §§4, 5, 7.3–7.5 | Shared app services; local REST; embedded React; SQLite | [`ARCHITECTURE-v0.5.md`](ARCHITECTURE-v0.5.md), [`API-DESIGN.md`](API-DESIGN.md), failure/data requirements | Reviewed contracts, OpenAPI and SQLite ADR before code |
| Architecture §§13, 16, 17 | Local mode, trust boundaries, release cross-platform evidence | Security, compatibility, test strategy | Browser-origin/path/migration and platform smoke |
| Roadmap v0.5 | Local web application and scan history | PRD and requirement IDs FR-501–FR-511 | Real local history workflow |
| v0.6 roadmap | Team service, identity, RBAC, exceptions, PostgreSQL later | PRD exclusions | No v0.6 scope in v0.5 endpoints/UI |
| v0.5 development plan §6 | UX and API design gates precede implementation | [`UX-DESIGN.md`](UX-DESIGN.md), [`API-DESIGN.md`](API-DESIGN.md), README DoR | Owner and independent reviewer approval before implementation issue creation |
| v0.5 development plan §8 | Failures and compatibility remain explicit | Failure/data and compatibility matrices | Failures do not become empty success; CLI regression |
| [`v0.1 scan error model`](../deprail-v0.1-execution-package/requirements/ERROR-MODEL.md) and [`v0.2 scan error model`](../deprail-v0.2-execution-package/requirements/ERROR-MODEL.md) | v0.1 enumerates established scan codes; v0.2 reinforces required result behavior and redaction | v0.5 [`requirements/ERROR-MODEL.md`](requirements/ERROR-MODEL.md) | Preserve scan codes/meanings; new history/API categories require separate approval and crosswalk |
| [`v0.3.1 cancellation error model`](../deprail-v0.3.1-execution-package/requirements/ERROR-MODEL.md) | Cancellation remains an explicit user-visible terminal operation state; this does not extend `ScanReport.Status` | v0.5 failure/data, API, UX contracts | Keep `operationOutcome=cancelled` distinct from report completeness |
| [`v0.3 remediation error model`](../deprail-v0.3-execution-package/requirements/ERROR-MODEL.md) | Remediation-plan errors only; not authoritative for scan/API error codes | v0.5 [`requirements/ERROR-MODEL.md`](requirements/ERROR-MODEL.md) | Do not infer v0.5 scan behavior from remediation-specific codes |
| v0.4 requirements | No `requirements/ERROR-MODEL.md` exists in the v0.4 package; its available requirements do not replace v0.1/v0.2 scan codes | v0.5 [`requirements/ERROR-MODEL.md`](requirements/ERROR-MODEL.md) | Name and review the actual contract instead of citing a nonexistent v0.4 error model |
| v0.4 closeout #356 | Outstanding release follow-ups remain open and require disposition | Plan §1, package README | Each gap assigned/carried/deferred; no assumed completion |

## FR-501–FR-511 execution traceability (draft)

The per-requirement verification IDs are defined in [`requirements/TEST-STRATEGY.md`](requirements/TEST-STRATEGY.md#9-requirement-and-threat-evidence-matrix). They specify the observable assertion, planned evidence artifact, and platform/browser matrix; they are not evidence that implementation tests have passed.

| Requirement | Design and failure contract | Required observable guard | Verification ID |
| --- | --- | --- | --- |
| FR-501 | PRD; [ADR-0004](../../adr/ADR-0004-local-scan-history.md); failure/data §§1, 4 | Repeated source `ScanReport.ScanID` values remain separate history entries; independently versioned `history-v1` retains source identity and only trustworthy, allowlisted provenance without altering CLI report bytes. | `FR-501` |
| FR-502 | API §§2–5; UX §§3–5; error model | Bounded, deterministic pages; empty results are distinct from loading, unavailable store, and API failure. | `FR-502` |
| FR-503 | API §§2–5; UX §§4–5; failure/data §§3–4 | Detail translates only safe projected fields; absent report, missing artifact, digest mismatch and invalid projection remain distinct, never fabricated evidence or empty success. | `FR-503` |
| FR-504 | Failure/data §§1–2; UX §5 | Operation outcome and optional report completeness remain independent through persistence, API, UI, and announcements. | `FR-504` |
| FR-505 | Architecture §§3, 6; v0.1 CLI contract; compatibility matrix | CLI command/output/exit behavior remains unchanged when the console or history store is absent/unavailable, unless a separate compatibility decision is approved. | `FR-505` |
| FR-506 | Architecture §§3, 5; API §1 | CLI and HTTP paths expose the same application/domain result; transport failures do not duplicate policy or scanner behavior. | `FR-506` |
| FR-507 | API §§3, 5–6; security requirements §§2–3 | Invalid identifiers/filters/routes/versions/methods, hostile origins, traversal, SQL metacharacters, and bounded-work failures are explicit and side-effect-free. | `FR-507` |
| FR-508 | [ADR-0004](../../adr/ADR-0004-local-scan-history.md); failure/data §§4–8; compatibility matrix | Fresh, supported, interrupted, corrupt, and future-schema outcomes preserve prior committed data under the proposed lifecycle; unresolved bounds/retention remain DoR gates. | `FR-508` |
| FR-509 | Architecture §6; API §7; UX §5–6 | Packaged assets match the API contract; missing assets or incompatible versions fail visibly without filesystem fallback. | `FR-509` |
| FR-510 | UX §§7–8 | Keyboard, assistive-technology, focus, status, contrast, responsive, and reduced-motion criteria pass on the actual packaged surface. | `FR-510` |
| FR-511 | PRD exclusions; API §3; UX §§3, 6 | Route and UI review confirms no hosted/team, remote publish, browser scan, remediation, repository mutation, or agent-write control. | `FR-511` |

## Remaining decisions and approval gates

ADR-0004's original schema-v1 and lifecycle proposal was owner-approved and merged in PR #407. Source-report validation exposed a conflict with its raw `report_json` assumption; the owner then selected the independent history-projection direction recorded in ADR-0004's additive section. The amended exact projection and all remaining lifecycle choices are still Proposed pending independent architecture/security review, numeric bounds, and retention-risk disposition.

**Owner-selected direction, exact contract not accepted:** [ADR-0004's additive finding and direction](../../adr/ADR-0004-local-scan-history.md#owner-selected-direction-independent-history-projection-not-accepted) show why existing `ScanReport` cannot be stored unchanged as a schema-valid privacy-safe v1alpha document. A separate allowlisted `history-v1` projection, distinct from SQLite schema v1, source report v1alpha, and `/api/v1`, leaves CLI JSON unchanged. The candidate storage layout, read projection, error mapping and future evidence matrix now reflect that decision; owner sign-off on exact fields, safe diagnostics, unknown-data policy and independent review are still needed before OpenAPI acceptance or implementation issues. A merged planning PR is not evidence of implementation readiness.

1. Approve or revise exact `history-v1` field mapping/validation and the candidate history-capture trigger and CLI compatibility behavior, including projection/persistence-failure and exit-code precedence.
2. Select evidence-backed per-entry, response, and store bounds, or explicitly accept unbounded growth; record disk-full refusal behavior.
3. Complete and validate the exact OpenAPI resource fields, pagination/cursor behavior, error mapping, version policy, and numeric request/response limits.
4. Review local listener bind/port/lifetime, Host/Origin/CORS/CSRF/authentication, and shutdown behavior.
5. Select the shadcn/ui-compatible component/dependency strategy and freeze supported browser/OS/accessibility combinations.
6. Name the independent architecture/security reviewer and record their approval separately from owner approval.
7. Disposition each #356 follow-up with owner, target, and evidence; no item is inferred complete.

No implementation issue map is authorized until these decisions and the complete Definition of Ready have linked evidence and separate reviewer/owner approval. Refer to the release plan for the design-first sequence.

## Conditional delivery trace

The [V05-005 readiness map](issues/V05-005-readiness-map.md#conditional-post-dor-delivery-map) groups the FR-501–FR-511 and SEC/STORE evidence into tentative store, shared-service, API, UI, packaging, and release-verification outcomes. It is not an implementation backlog: #396 and #393 remain in Review, ADR-0004 remains Proposed pending independent review, and the execution package's [DoR evidence map](README.md#readiness-evidence-and-blockers) keeps every unapproved gate explicit. Final issue ownership, labels, interfaces, dependencies and acceptance evidence follow separately approved contracts, not this provisional grouping.
