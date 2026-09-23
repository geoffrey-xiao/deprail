# DepRail v0.5.0 Development Plan

| Attribute | Value |
| --- | --- |
| Release | v0.5.0 |
| Plan revision | 1.0 |
| Status | Draft for owner and independent architecture/security review; implementation not authorized |
| Target mode | Preview first; stable follow-up only after complete evidence |
| Whole-project roadmap | [`deprail-roadmap-v1.md`](deprail-roadmap-v1.md) |
| Product baseline | [`../01-product/deprail-product-design-v1-ai.md`](../01-product/deprail-product-design-v1-ai.md) |
| Architecture baseline | [`../02-architecture/deprail-architecture-and-tech-stack-v1.md`](../02-architecture/deprail-architecture-and-tech-stack-v1.md) |
| Predecessor | v0.4.0-preview.2; release closeout issue [#356](https://github.com/geoffrey-xiao/deprail/issues/356) remains open |
| Preparation issue | [#387](https://github.com/geoffrey-xiao/deprail/issues/387) |
| Target milestone | [`v0.5.0`](https://github.com/geoffrey-xiao/deprail/milestone/11) |
| Project view | [Release · v0.5.0](https://github.com/orgs/geoffrey-xiao/projects/1/views/6) |
| Owner | `@geoffrey-xiao` |
| Independent reviewer | To be named; required before Definition of Ready approval |

This is a release planning baseline, not an implementation authorization. Proposed scope and design decisions remain subject to owner and independent architecture/security review. Do not create implementation issues or start runtime work until the gates in this plan are accepted and the execution package is complete.

## 1. Reconciliation and release decision

The planning hierarchy is normative in this order:

1. Product design: [`../01-product/deprail-product-design-v1-ai.md`](../01-product/deprail-product-design-v1-ai.md).
2. Architecture: [`../02-architecture/deprail-architecture-and-tech-stack-v1.md`](../02-architecture/deprail-architecture-and-tech-stack-v1.md).
3. Whole-project roadmap: [`deprail-roadmap-v1.md`](deprail-roadmap-v1.md).
4. This v0.5 release plan.
5. A v0.5 execution package, requirements, and reviewed issue contracts to be prepared after plan review.

The roadmap sets v0.5's outcome as a local web application and scan history. Product and architecture specify an embedded React/TypeScript/Vite console, SQLite local storage, and a local REST API, while keeping web as a collaboration/presentation layer over shared application services. v0.6 team services, identity, and PostgreSQL are later roadmap work. This plan selects that bounded v0.5 outcome; it does not authorize hosted services or an expanded collaboration product.

v0.4.0-preview.2 was published with an owner-approved `go with approved gaps`, but its issue #356 remains open. Its independent architecture/security review, rollback ownership/recovery procedure, Python/Java remediation coverage, cross-platform representative smoke, and signature/provenance dispositions are not recorded as complete in the release-evidence document. Before Definition of Ready approval, #387 must link an explicit disposition for each predecessor item: carry forward into v0.5 only when relevant, track separately, or accept/defer with an owner and target. None is presumed resolved by this plan.

| Capability | Product/architecture contract | Roadmap item | v0.5 proposal |
| --- | --- | --- | --- |
| Local scan history | Product design §Web Product Design; architecture §§5, 13, 16 | v0.5 local web and history | Include local persistence and browsing of completed scan history; define retained data, provenance, and failure semantics before implementation. |
| Local web console | Product design §§Web Product Design, Technical Architecture; architecture §§4, 13, 16 | v0.5 local web and history | Include an embedded React/TypeScript/Vite console using the existing application layer. User preference: shadcn/ui-inspired style, subject to a reviewed design-system and accessibility decision. |
| Local REST API | Architecture §§4, 5, 7.4, 13 | v0.5 local web and history | Include only the API surface required by the accepted local workflows; define OpenAPI contract and local security boundary before implementation. |
| Team collaboration | Product design §Web Product Design; architecture §§1, 7, 16 | v0.6 | Exclude PostgreSQL, identity, RBAC, remote publishing, and multi-user workflows. |
| Agent/MCP writes, new scanners, autonomous mutation | Product design §§Explicit Non-goals, Agent and Skill Design; architecture §§2, 7.7 | v0.7/v0.8 or later | Exclude; no API or UI affordance may silently authorize these capabilities. |

## 2. User outcome

A developer can use DepRail locally to review prior completed scans and inspect the current scan's project, workspace, findings, completeness, and provenance through a browser-based console. The console reads and updates only the local state permitted by reviewed v0.5 contracts, and it remains a presentation client of the shared application services rather than a second implementation of scan or policy behavior.

The release succeeds only when the stored history is trustworthy and understandable: failed or partial runs remain visibly distinct from complete runs, persisted scan data retains its source and artifact provenance, and local API/storage failures never appear as successful empty history or successful scans.

## 3. Included scope (proposed)

- Local persistence of scan history in SQLite, with an explicit schema, migration, transaction, corruption, and recovery contract.
- History workflows sufficient to list and inspect saved scans and navigate from a scan to its workspaces and findings, preserving existing report semantics and provenance.
- A local REST API backed by the existing application/domain services and described by a versioned OpenAPI 3.1 contract.
- An embedded React, TypeScript, and Vite console built as static assets and served by the local application. Use a shadcn/ui-inspired visual direction: accessible, composable primitives, consistent spacing/type/color tokens, and clear state surfaces. Final library/dependency and component decisions require design review.
- Explicit complete/partial/failed scan presentation and diagnostics for unavailable, malformed, corrupt, or incompatible stored data and API failures.
- Local-only data handling, safe path/storage boundaries, bounded requests, and cross-platform build/package/smoke coverage.
- Preview release evidence, user documentation, migration/compatibility notes, and a retrospective after publication.

Exact endpoints, resource models, persistence retention policy, UI routes, component inventory, and frontend dependency versions are undecided until the design gates below pass. They must not be inferred from this plan.

## 4. Explicit exclusions

- Hosted web services, remote publishing, PostgreSQL, identity, authentication for multi-user access, RBAC, team workflows, notifications, and shared/cloud history.
- MCP or agent write capabilities, autonomous approval, pull-request creation, commits, pushes, merges, or other repository mutation.
- New scanners, vulnerability intelligence, policy language, or changes to existing scan semantics unrelated to presenting and persisting results.
- Changes to existing CLI/JSON contracts unless a separately reviewed compatibility decision requires them.
- A general-purpose API/plugin platform, arbitrary third-party integrations, telemetry, source upload, or network access for default local workflows.
- A full-screen terminal UI or replacement of the CLI.
- Stable-release readiness claims before a separate complete stable release gate.

Any candidate capability outside the included scope requires an explicit owner decision and versioned scope/architecture change record before issues or implementation are added.

## 5. Architecture and design constraints

### Shared service boundary

The HTTP layer and CLI must call the same application services. UI components must not implement scanner, normalization, completeness, remediation, or policy decisions. Domain packages remain independent of React, HTTP transport, and SQL; transport adapters translate requests to application commands and domain results.

### Local API trust boundary

Before implementation, the API design must specify bind address and exposure, origin/CSRF behavior, request method and body limits, error mapping, cancellation/deadline behavior, and whether any mutating endpoint is needed. The default proposal is loopback-only with no remote publishing or source upload; do not expose a LAN/public listener or add an authentication bypass without an explicit security decision. Static asset serving and API routing must not permit path traversal or unintended file disclosure. CORS must be restrictive and unnecessary cross-origin access avoided. These are decision points for review, not implemented guarantees.

### SQLite and data lifecycle

The persistence design must define the stored records, schema versioning, migration direction, transaction boundaries, concurrent access, retention/deletion, disk exhaustion, database corruption, and consistency with raw content-addressed artifacts. Preserve scanner/tool/database provenance and raw artifact digests. A successful scan must not be lost or represented as persisted when its durable history write failed; the API/UI must distinguish the scan outcome from history-storage failure. No silent destructive migration or database replacement.

### UI system and accessibility

The user-selected direction is shadcn/ui style. The UX gate must decide the actual React component approach (including whether to use shadcn/ui components, compatible primitives, or a minimal project-owned layer), token ownership, iconography, typography, responsive breakpoints, and dependency/licensing/build implications. The interface must be keyboard navigable, have visible focus, meaningful labels/contrast, reduced-motion behavior, and status announcements; status cannot depend on color alone. This preference does not authorize a new design system or package dependency before review.

### Privacy and security

Keep source and history local by default. Treat database content, API inputs, repository-derived labels, findings, and paths as untrusted. Apply canonical path containment, parameterized SQL, request bounds, output escaping, restrictive permissions where supported, and secret redaction. Do not place credentials or full source contents into history or UI diagnostics unless the accepted data contract explicitly requires and protects them.

## 6. Design-first delivery sequence and gates

The following sequence is mandatory. UX and API design are completed and reviewed before any runtime implementation issue is created or started. Their work may proceed in parallel only after product scope and shared domain/resource assumptions have been reconciled.

| Order | Stage | Required deliverable | Gate / dependency |
| ---: | --- | --- | --- |
| 0 | Release reconciliation | Accepted scope, v0.4 follow-up dispositions, named owner and independent reviewer, compatibility and risk decisions | Complete #387 planning work; no implementation authorization. |
| 1 | Product workflows and UX research | Target users/jobs, primary navigation and routes, workflow/state diagrams, information architecture, wireframes/prototype, visual tokens, accessibility checklist, error/empty/loading/partial/failed states | Owner and design reviewer approve UX package; shadcn/ui direction and dependency strategy recorded. |
| 2 | API and data contract design | Versioned OpenAPI 3.1 proposal, resource and error models, endpoint/method matrix, pagination/filtering semantics if needed, request limits, bind/origin/security rules, cancellation semantics, compatibility/versioning policy | Independent architecture/security reviewer and owner approve contract; no endpoint or schema is assumed beforehand. |
| 3 | Persistence architecture decision | SQLite schema and migration plan, transactions/concurrency, artifact references, retention/deletion, corruption/disk-full behavior, backup/recovery, file permissions | ADR and data contract accepted; crosswalk to existing scan/artifact contracts. |
| 4 | Execution package and implementation decomposition | Requirements, threat model, test strategy, compatibility matrix, ordered epics/issues, acceptance evidence, rollback and release plan | UX, API, and persistence design gates all pass; Definition of Ready approved by owner and named reviewer. |
| 5 | Implementation | Ordered issue work derived from approved contracts (persistence/service first; API transport; UI integration; packaging/docs) | Only then create implementation issues and start code; one bounded issue at a time. |
| 6 | Release verification | Real local workflows, storage/API/UI failure evidence, cross-platform packaging/smoke, security review, artifacts, owner release decision | Version-specific evidence complete or gaps explicitly approved for preview. |

### UX design acceptance

Before coding, the UX package must demonstrate the agreed local workflows and, at minimum, specify:

- Landing/history view and scan selection/navigation.
- Scan detail with complete/partial/failed status, findings summary, workspaces, provenance, and diagnostics appropriate to accepted data contracts.
- Empty history, no findings, loading, cancellation, partial results, scan failure, storage/API unavailable, migration/corruption, and narrow viewport states.
- Keyboard interaction, focus management, screen-reader names/status, contrast, responsive behavior, and reduced motion.
- Reusable visual tokens and representative shadcn/ui-style components, plus the decision on actual library/dependencies.
- Explicitly excluded controls (team/account/publish/agent-write) so the design cannot imply unavailable capabilities.

### API design acceptance

Before API implementation, the reviewed OpenAPI contract and security review must define:

- Resources and operations required by accepted UX; read-only by default and no unspecified mutation endpoints.
- Stable response schemas and error codes that preserve complete/partial/failed distinctions and provenance.
- API versioning and additive/breaking change policy, serialization compatibility, and schema validation evidence.
- Local listener/origin/CORS/CSRF rules, path and file disclosure controls, request/body/response limits, and authentication decision for local browser access.
- Pagination or bounded history retrieval, ordering/determinism, filtering, timeouts/cancellation, and behavior on SQLite/API failures.
- No leakage of credentials, source content, or host-sensitive paths beyond approved fields.

## 7. Dependencies and sequencing

- **#387** owns preparation, contract reconciliation, the plan/execution package, named review, and Definition of Ready. It is currently `Todo` in the `v0.5.0` Project view.
- **#356** is still open for v0.4.0-preview.2 closeout. Its remaining items must receive explicit dispositions and owners; v0.5 cannot claim those gates complete.
- Product, architecture, and roadmap contracts remain prerequisites. This plan does not supersede them.
- UX flows/resource assumptions must be reconciled before the UX and API packages are finalized. After that shared prerequisite, UX and API design may proceed in parallel.
- Persistence schema/migration design depends on accepted API resource needs and existing scan/artifact contracts.
- Implementation issues depend on approved UX, API, persistence, execution package, and Definition of Ready.
- Release publication depends on implementation acceptance and version-specific evidence.

## 8. Compatibility and failure behavior

- Existing CLI commands, stable JSON schemas, exit codes, scan semantics, and artifact identity remain unchanged unless an explicit compatibility decision is approved.
- Local history is additive. Existing users without a database must retain current CLI behavior; local console/history availability must be diagnosable and must not make scanning depend on a running browser or persistent store unless the owner explicitly approves that contract.
- Schema migrations must be versioned, deterministic, tested against upgrade and failure paths, and non-destructive by default. Migration failure must preserve the prior usable data and report a typed, actionable error; exact recovery behavior must be designed before implementation.
- HTTP errors, timeouts, cancellation, malformed requests, oversized payloads, unavailable storage, full disk, corrupt DB, missing/corrupt embedded assets, and incompatible schema versions are explicit failures, never empty-success substitutes.
- Partial or failed scans remain distinguishable from complete scans through persistence, API, and UI. A failed history write cannot turn the scan into a false-success persistence claim.
- Repository paths use canonical containment; serialized relative paths use `/`. UI output escapes untrusted labels and does not leak secrets.

## 9. Proposed release acceptance gate

The gate is proposed and requires owner/reviewer approval. At minimum:

- Reviewed UX package, API/OpenAPI contract, persistence decision, threat model, and execution package are linked from implementation contracts before runtime work.
- Local web workflows run through shared application services and do not duplicate domain/security logic.
- A user can inspect deterministic local scan history and scan details, with correct provenance and complete/partial/failed behavior.
- SQLite migrations, concurrent access, retention/deletion, interrupted writes, disk-full and corruption behavior meet the reviewed contract without silent data loss.
- API requests/responses validate against the versioned contract; error, limits, origin/bind, cancellation, and path-safety behavior are exercised.
- UI matches reviewed flows; keyboard/screen-reader/focus/contrast/responsive/reduced-motion states are verified on the actual surface.
- Local-only operation is demonstrated; no unexpected source upload, public listener, credential exposure, or repository mutation occurs.
- Embedded assets and binary work on supported Linux, macOS, and Windows targets; local representative scan-history workflow and release artifacts are verified.
- Existing CLI/JSON compatibility is checked; unit, integration, contract, migration, security, and cross-platform evidence meet the execution package.
- Independent architecture/security review and owner decision are recorded separately; preview gaps, rollback owner/procedure, artifact checksums/SBOM/signature/provenance status, and retrospective follow-ups are explicit.

## 10. Risks and controls

| Risk | Control / required decision |
| --- | --- |
| Local API unexpectedly reachable or vulnerable to browser-origin abuse | Decide bind/origin/CSRF/auth behavior in API design; test actual packaged server and malicious-origin cases. |
| SQLite history stores sensitive or stale data indefinitely | Minimize stored fields, define retention/deletion and permissions, document local data location, and test deletion/recovery. |
| UI misrepresents partial or failed scans as safe/empty | Preserve domain outcome states end-to-end; test every state in API and browser. |
| UX-first design assumptions conflict with API/data model | Reconcile workflows and resource concepts first; design packages can iterate together; freeze both before coding. |
| shadcn/ui dependency or generated component code adds maintenance/licensing/build burden | Review package versions, generated code ownership, primitives, license, accessibility, and embedded build impact before adoption. |
| Embedded frontend increases binary/package size or harms startup | Measure built artifacts and startup/resource behavior against an agreed budget in the execution package. |
| Migrations or interrupted storage operations lose history | Versioned migrations, transaction/atomicity strategy, backup/recovery decision, old-schema fixtures, interruption and corruption tests. |
| v0.4 release gaps are mistaken for closed | Track #356 separately; explicit owner-approved carry-forward/deferral for each row; do not mark historical evidence complete by inference. |

## 11. Preview exclusions and displaced work

v0.5 consumes roadmap capacity planned for the local web/history milestone. Team collaboration and PostgreSQL remain displaced to v0.6; scoped MCP/agent work remains v0.7; Trivy/external adapters remain v0.8. No additional capability is displaced or added without a documented owner decision. If UX/API review shows that the local web outcome cannot meet the security or compatibility boundary, pause implementation and record a versioned scope decision rather than silently weakening the gate.

## 12. Definition of Ready

Runtime implementation remains blocked until all items below are linked with evidence:

- [ ] This plan reconciled with product, architecture, roadmap, and v0.4 predecessor evidence.
- [ ] Every #356 carry-forward item is explicitly accepted, assigned, or deferred by the owner; no pending item is implied complete.
- [ ] Owner confirms included/excluded capabilities, compatibility boundary, failure model, preview mode, and acceptance gate.
- [ ] Independent architecture/security reviewer is named and records review separately from owner acceptance.
- [ ] UX package meets the design acceptance section and records shadcn/ui/library decisions.
- [ ] API/OpenAPI and local API security contracts meet the API acceptance section.
- [ ] SQLite/data lifecycle ADR and migration/failure contracts are reviewed.
- [ ] Execution package, requirements, threat model, test/compatibility strategy, rollback/recovery plan, and evidence checklist exist.
- [ ] Every proposed implementation issue maps to an approved contract, has one primary outcome, explicit dependencies, observable acceptance, failure behavior, reviewer, and evidence.
- [ ] Milestone, Project view, issues, labels, and evidence links are synchronized and verified.

## 13. Release evidence and follow-up

Use [`../RELEASE-CHECKLIST.md`](../RELEASE-CHECKLIST.md) and a version-specific evidence record. Record reviewed commit and version identity; exact verification commands/results; API/OpenAPI/schema and migration evidence; Linux/macOS/Windows build and smoke results; browser accessibility and UX review evidence; artifact names/checksums/SBOM/signature/provenance status; owner and independent reviewer decisions; remaining risks and rollback procedure. After publication, create and link `docs/retrospectives/RETROSPECTIVE-v0.5.0-preview.N.md` or the applicable final release version and assign follow-up owners.
