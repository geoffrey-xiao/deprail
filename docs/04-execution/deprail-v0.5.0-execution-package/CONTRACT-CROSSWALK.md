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
| v0.4 closeout #356 | Outstanding release follow-ups remain open and require disposition | Plan §1, package README | Each gap assigned/carried/deferred; no assumed completion |

## Unresolved decisions

1. Is scan history captured automatically during CLI scan or through an explicit local workflow?
2. Which scan report fields are stored/indexed, and how are raw content-addressed artifacts referenced/retained?
3. What is the data root, permission model, retention/deletion and export behavior?
4. What migration, transaction/concurrency, backup, disk-full, corruption and interrupted-write contract is accepted?
5. Which API resource shape, endpoint pagination, error codes, version policy and payload limits are required by approved UX?
6. What are listener startup/port/browser-open, Host/Origin/CORS/CSRF/auth, and shutdown behavior?
7. Which shadcn/ui-compatible implementation and pinned dependency/build strategy is accepted?
8. Which browsers, platform/browser combinations, bundle limits, and manual smoke criteria are supported?
9. Who is the independent architecture/security reviewer, and how are all #356 follow-ups dispositioned?

No implementation issue map is created until those decisions and Definition of Ready are accepted. Refer to the release plan for the design-first sequence.
