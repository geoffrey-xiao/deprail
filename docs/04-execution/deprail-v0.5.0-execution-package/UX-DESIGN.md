# v0.5 Local Console UX Design Contract

**Status:** Draft design brief; not a final visual design or implementation authorization.
**Product/design baseline:** [`PRD-v0.5.md`](PRD-v0.5.md), [`ARCHITECTURE-v0.5.md`](ARCHITECTURE-v0.5.md), and [`../../03-planning/deprail-development-plan-v0.5.0.md`](../../03-planning/deprail-development-plan-v0.5.0.md).

## 1. Design goal

Make saved local scan results easy to find and trustworthy to interpret. The console should emphasize scan status, time, repository/workspace context, and findings while retaining source/tool provenance and actionable diagnostics. It is not a replacement for CLI workflows or a hosted security dashboard.

## 2. Visual direction

User preference: shadcn/ui style. Proposed characteristics, pending design review:

- Quiet neutral surfaces with semantic accent colors, subtle borders, restrained shadows, clear typography, consistent spacing, and dense but readable tables/cards.
- Composable, accessible controls rather than bespoke interactions; visual consistency across history, scan detail, dialogs, and banners.
- Use semantic tokens for background, foreground, border, muted, primary, destructive, warning, success, and focus-ring states. Exact palette, font, icon set, radius, and spacing scale remain to be selected.
- Status communicates through text/icon/shape as well as color. Critical state remains legible in monochrome and high-contrast modes.
- Mobile/narrow desktop layout preserves identity and completeness first; secondary provenance is progressively disclosed without being dropped.

This brief does not select package versions, copy generated components, or approve a dependency. The design gate must compare using shadcn/ui source components, compatible accessible primitives, or a small project-owned layer against React/Vite embedding, license, maintenance, build, and accessibility needs.

### Draft visual proposal (unapproved)

PNG previews and editable SVG sources:

- Scan history, desktop: [`history-desktop.png`](design/ux/history-desktop.png) ([SVG](design/ux/history-desktop.svg)).
- Scan history, narrow layout: [`history-mobile.png`](design/ux/history-mobile.png) ([SVG](design/ux/history-mobile.svg)).
- Scan detail, desktop: [`scan-detail-desktop.png`](design/ux/scan-detail-desktop.png) ([SVG](design/ux/scan-detail-desktop.svg)).
- Loading, empty, and recovery states: [`state-patterns.png`](design/ux/state-patterns.png) ([SVG](design/ux/state-patterns.svg)).
- Scan detail, narrow layout: [`scan-detail-mobile.png`](design/ux/scan-detail-mobile.png) ([SVG](design/ux/scan-detail-mobile.svg)).
- Scan detail outcome/report variants: [`scan-detail-statuses.png`](design/ux/scan-detail-statuses.png) ([SVG](design/ux/scan-detail-statuses.svg)).

These are proposal artifacts, not accepted schemas or implementation authorization. All names, IDs, counts, timestamps, and finding examples are synthetic. Routes, exact response fields, ordering, page size, pagination behavior, finding-detail depth, and component dependencies remain open for the API/UX review.

- Proposed client navigation: `/history` and `/history/{historyEntryID}`; both routes are directly addressable. This does not define an API route.
- From a history row/card, open its detail route as an in-app history entry; browser Back returns to the originating history list, and Forward restores detail if the entry is still available.
- When `/history/{historyEntryID}` is opened directly or in a new tab, browser Back retains normal browser behavior and may leave the console. The visible `Back to scan history` link is always available as the in-app fallback. A missing or stale entry shows the not-found state with that link; do not invent an in-app history entry for direct navigation.
- Restoring the list's exact page, cursor, or filter selection remains deferred until the API pagination and URL-state contracts are approved; these are client routes, not API routes.
- A project-level expandable table was considered but deferred for v0.5. It could help users browse repeated scans per repository, but FR-502 and the proposed `GET /api/v1/scans` page scan entries; no stable project-group identity or bounded group/child-pagination contract is approved. Grouping risks splitting a project’s scans across pages and adds nested-table keyboard/screen-reader complexity. Revisit only with user evidence and an explicit data/API contract update.
- The detail proposal places operation outcome before report completeness and keeps workspaces, findings, provenance, and diagnostics distinct. The existing desktop and new narrow examples show a cancelled operation with a complete returned report; the status board separately illustrates completed/partial, completed/failed, failed/no-report, and completed/complete/zero-findings. Only the last case may present zero findings as a complete result.
- The state board distinguishes loading, empty history, store unavailable, incompatible/corrupt storage, missing or mismatched raw artifacts, and stale selection. Retry is shown only for a safe read; no destructive recovery control is proposed.
- The mockups use system UI fonts and local vector shapes only; no remote fonts, icons, or images are required. The shadcn/ui direction remains a visual reference, not a dependency or component-source decision.
- Repository-controlled strings must render as inert text. Review cases include a label such as `<script>alert(1)</script>` and a long path; text wraps or truncates without changing meaning or causing horizontal page overflow.

Exact color tokens and dependency/build decisions require owner and independent UX/security review. Contrast, keyboard, screen-reader, responsive, and reduced-motion acceptance remain to be verified on the eventual implementation surface.


## 3. Information architecture (proposal)

1. **Scan History** — primary route; paginated saved scans, status and summary filters only if justified by user research.
2. **Scan Detail** — selected scan metadata, completeness/diagnostics, finding summary, workspace list, findings, and provenance.
3. **About/Local data help — recommendation, pending owner decision:** include a static, read-only view explaining local-only operation and pointing to approved data-location guidance. Do not expose absolute host paths or offer reset, delete, or export controls.

Navigation includes a persistent product identity and a clear return-to-history path. Do not show controls for unsupported team, publish, exception-management, agent-write, or remediation capabilities. The owner must accept or exclude About/Local data help before UX acceptance; if excluded, retain a concise local-only note on history. Exact data location remains deferred until its storage contract is approved.

## 4. Primary workflows

### Review history

1. Open the local console.
2. See a paginated ordered list, or an explicit empty-history state.
3. Identify an entry using its unique history-entry identity, timestamp, root label approved by the data contract, execution outcome, report completeness when present, finding count, and concise source/tool summary. Repeated runs may share a source scan ID but remain separate entries.
4. Select a row/card to open details by `historyEntryID`; preserve the selected entry in browser navigation/back behavior.

### Inspect scan detail

1. See execution outcome at top: Completed, Failed, or Cancelled. If a report exists, show its completeness separately as Complete, Partial, or Failed; a cancelled execution must never appear as a completed clean scan because its report retained an earlier completeness value.
2. See unique history-entry identity, source `ScanReport.ScanID`, completion/start time if available, repository/workspace summary, scanner/tool/database provenance, and diagnostics.
3. Review finding summary and findings with stable identifiers and key fields; open details without losing current context.
4. If a raw artifact is unavailable or unverifiable, show that explicitly and do not fabricate missing evidence.

### Recover from unavailable data

- Distinguish empty history from unavailable API/store, loading, corrupt data, incompatible version, and artifact retrieval failure.
- Show actionable guidance and safe retry where it does not repeat a mutation (v0.5 is proposed read-only).
- Do not present stale cached data as current unless visibly labeled and the API/data contract permits it.

## 5. Required states and copy intent

| State | Required representation |
| --- | --- |
| Initial loading | Stable skeleton/progress label; no invented percentage. |
| Empty history | Clear no-saved-scans message plus safe next step; distinct from error. |
| Completed operation, complete report, zero findings | Clearly show completed execution and zero findings; no ambiguity. |
| Completed operation, complete report, findings | Finding count/severity summary and details navigation. |
| Partial report, any operation outcome | Persistent warning label; explain omitted/incomplete scope and show diagnostics while retaining the separate execution outcome. |
| Failed operation with no report | Show Failed outcome, report unavailable, findings unavailable, and diagnostics; never show zero findings. |
| Failed report after completed operation | Show Completed outcome and Failed report separately; findings are unavailable and diagnostics remain visible; never call this clean. |
| Cancelled execution | Show Cancelled as the operation outcome, separately from any returned report completeness; never treat it as a completed scan. |
| API/storage unavailable | Error state distinct from empty; retry/action only when safe. |
| Migration/storage corruption | Explain history availability and recovery path without suggesting destructive reset. |
| Missing artifact | Preserve report metadata; disclose evidence unavailable/integrity failure. |
| Not found/stale selection | User-visible not-found state; return to history. |

Text labels and aria status announcements convey state. Color alone is never sufficient.

### API response-to-state map (proposal)

The endpoint names and error codes below are candidates from `API-DESIGN.md` and `requirements/ERROR-MODEL.md`. HTTP status mapping and wire schemas remain unapproved; this map defines the intended user-visible distinction, not a frozen transport contract.

| Response or condition | Required UI state and behavior |
| --- | --- |
| Initial `GET /api/v1/scans` request | Loading skeleton with an announced loading label; no invented progress percentage. |
| Successful history collection with no entries | Empty history only after a successful response confirms the collection is empty; never use this state for a request or storage failure. |
| Successful history page with entries | Render each unique history entry. Show operation outcome separately from optional report completeness. For `completed` + `failed`, label the report `Failed` and do not show zero findings as a clean result. |
| Initial detail or child-collection request | Show loading state in the affected detail/collection region; do not replace the overall entry or another loaded section with an empty state. |
| Successful `GET /api/v1/scans/{historyEntryID}` with a report | Render the entry identity, operation outcome, report completeness, provenance, workspaces, findings, and diagnostics. Keep operation outcome and report completeness independent. |
| Successful detail with no report | Render trustworthy operation metadata and diagnostics, label report/findings as unavailable, and do not synthesize an empty report or a zero-finding result. |
| Successful workspace collection with entries | Render the returned workspace page with its parent entry context and workspace completeness. |
| Successful workspace collection confirmed empty | Show that no workspaces were recorded only after successful completion of the collection; a failed or not-yet-loaded collection is not empty. |
| Successful findings collection with entries | Render the returned finding page with stable identifiers and workspace/provenance context; retain the parent report state. |
| Successful findings collection confirmed empty, complete report | Show no findings only when the API confirms the collection is empty/end-of-collection and the report is `complete`. |
| Successful findings collection confirmed empty, partial/failed/no report | Keep the partial, failed, or unavailable report state and diagnostics; never describe this as a clean zero-finding result. |
| Workspace/finding child request fails | Show an error scoped to that section and retain other trustworthy detail; do not substitute an empty collection. |
| `HISTORY_UNAVAILABLE` | Store-unavailable state; explain that history could not be read and offer only safe read retry. |
| `HISTORY_SCHEMA_UNSUPPORTED`, `HISTORY_MIGRATION_FAILED`, or `HISTORY_CORRUPT` | Compatibility/recovery state; preserve existing data and do not suggest reset or deletion. |
| `HISTORY_ENTRY_NOT_FOUND` | Stale-selection state with a visible return-to-history link; do not substitute another entry. |
| `HISTORY_ARTIFACT_MISSING` or `HISTORY_ARTIFACT_DIGEST_MISMATCH` | Keep trustworthy metadata if permitted, disclose unavailable/unverified evidence, and never substitute an empty artifact or report. |
| `API_REQUEST_INVALID`, `API_REQUEST_TOO_LARGE`, `API_VERSION_UNSUPPORTED`, `API_METHOD_UNSUPPORTED`, or `API_ORIGIN_REJECTED` | Explicit safe request/compatibility error; never render empty history or child data. |
| `API_TIMEOUT` or `API_CANCELLED` on a read | End the affected loading state with a request-level message. Do not rewrite saved operation outcome or report completeness. |

The screen presents `operationOutcome` (`completed|failed|cancelled`) and, when a report exists, `reportStatus` (`complete|partial|failed`) as independent values. Candidate codes do not replace the approved error contract, and an error response is never treated as a successful empty page.

## 6. Component and implementation direction (proposal; not approved)

The visual direction is shadcn/ui-inspired, not a commitment to install shadcn/ui or copy generated components. The conservative design baseline is a small project-owned layer built from semantic HTML and CSS; no new runtime dependency is selected or authorized here.

| Decision area | Proposed direction | Approval/build boundary |
| --- | --- | --- |
| Primitives | Native links, buttons, headings, lists, semantic tables, status text, alerts, skeletons, and pagination; use links for history-to-detail navigation. Avoid dialogs for primary detail navigation. | Verify keyboard and assistive-technology behavior on the actual browser surface. |
| Component ownership | Hand-authored components and styles belong to DepRail and are reviewed in the repository. | If shadcn/ui source or another library is later preferred, record copied-source ownership, package/version, license, maintenance, and bundle impact before adoption. |
| Dependencies | No third-party component, icon, chart, or font dependency is selected by this UX proposal. | Any runtime package remains subject to owner and independent review; this document does not approve installation. |
| Tokens | Project-owned CSS custom properties for background, foreground, border, muted, primary, destructive, warning, success, and focus states. | Exact palette, spacing, type scale, and contrast values remain design-review decisions. |
| Icons and fonts | System font stack and local/inline SVG icons; no remote fonts, images, CDN, or analytics. | Keep assets local and review any future dependency or license change. |
| Embedded build | React, TypeScript, and Vite static output embedded in the Go application, consistent with the architecture baseline. | `web/` is the source ownership boundary; package manager, pinned versions, embed path/tool, and bundle/startup budget remain execution-package decisions. |

The proposal preserves the requested shadcn/ui visual characteristics without freezing a library or implementation. Owner and independent UX/security review must approve the actual design-system, dependency, and build decisions before runtime work.

## 7. Accessibility and responsive acceptance

- All navigation, filters, pagination, and finding actions are keyboard reachable with a logical focus order; visible focus remains clear in every state.
- The history-to-detail flow works with keyboard only: focus a row/card link, open detail, use the visible return link, then browser Back and Forward. Focus returns to a predictable location.
- Headings and landmarks express page structure; controls have programmatic labels; tables identify column headers; links have distinguishable names.
- Loading, errors, cancellation, and status changes are announced without repeatedly interrupting assistive technology. Outcome and report completeness are both conveyed in text.
- Proposed contrast target is WCAG 2.2 AA: at least 4.5:1 for normal text, 3:1 for large text, and 3:1 for meaningful UI/focus indicators. Verify all interactive and status states; this target remains subject to design review.
- At 320 CSS px width and 200% zoom, essential status, identity, findings information, and navigation remain available without page-level horizontal scrolling; wide tables use an intentional responsive pattern. The 390 px narrow-detail mockup is illustrative and does not verify this acceptance.
- Reduced-motion preferences are respected; movement is nonessential and disabled or simplified when `prefers-reduced-motion` is enabled.
- Repository-controlled strings remain inert text. Review markup-like labels such as `<script>alert(1)</script>`, long paths, Unicode, and direction-control characters; they must not become executable markup, navigation targets, or misleading visual labels.
- Select and document the supported browser/OS and screen-reader combination in the compatibility/test plan. Record the environment and observed gaps; a static mockup is not accessibility verification.

Color alone never communicates status. PNG/SVG inspection can verify layout and represented states only; keyboard, screen-reader, contrast, reduced-motion, and hostile-text acceptance require the actual browser surface before implementation is accepted.

## 8. UX evidence package required before implementation

The design package links draft history (desktop/narrow), detail (desktop/narrow), outcome/report-state variants, empty/loading/recovery wireframes, and editable sources above. Before UX acceptance, the owner and design reviewer must review and record:

- User/workflow assumptions and scope decisions, including the About/Local data help decision.
- History/detail navigation, direct-route behavior, and browser back/forward behavior, including stale-selection recovery.
- The response-to-state map from API success, candidate errors, and report/operation values to visible UI states.
- The component/dependency, code ownership, icon/font, token, and embedded-build proposal, with every unapproved choice clearly marked.
- Responsive breakpoints, keyboard flow, focus behavior, contrast target, screen-reader names/status, reduced-motion behavior, and hostile-data cases.
- Owner and independent reviewer acceptance separately; neither is implied by a merged design PR.

The SVG/PNG examples are synthetic proposal artifacts. They do not approve an API, dependency, storage behavior, visual system, or runtime implementation, and they do not replace browser accessibility evidence.
