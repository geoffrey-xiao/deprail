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

## Unresolved decisions

1. What uniquely identifies a history entry across repeated runs with the same `ScanReport.ScanID`, and how are ingestion retries made idempotent without deduplicating distinct runs?
2. How is history captured, and which optional report fields and raw content-addressed artifacts are stored/referenced/retained?
3. What is the data root, permission model, retention/deletion and export behavior?
4. What migration, transaction/concurrency, backup, disk-full, corruption and interrupted-write contract is accepted?
5. Which API resource shape, endpoint pagination, error codes, version policy and payload limits are required by approved UX?
6. What are listener startup/port/browser-open, Host/Origin/CORS/CSRF/auth, and shutdown behavior?
7. Which shadcn/ui-compatible implementation and pinned dependency/build strategy is accepted?
8. Which browsers, platform/browser combinations, bundle limits, and manual smoke criteria are supported?
9. Who is the independent architecture/security reviewer, and how are all #356 follow-ups dispositioned?

No implementation issue map is created until those decisions and Definition of Ready are accepted. Refer to the release plan for the design-first sequence.
