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

These are proposal artifacts, not accepted schemas or implementation authorization. All names, IDs, counts, timestamps, and finding examples are synthetic. Routes, exact response fields, ordering, page size, pagination behavior, finding-detail depth, and component dependencies remain open for the API/UX review.

- Proposed client navigation: `/history` and `/history/{historyEntryID}`; the history entry remains addressable, and both the visible back link and browser Back return to the list. This does not define an API route.
- Desktop history remains a flat scan-entry table with repository context in the first column; narrow layouts use stacked scan cards. Each row/card is one link to `/history/{historyEntryID}`. The repository label and a visible focus treatment expose the detail navigation; the visible back link and browser Back return to history.
- A project-level expandable table was considered but deferred for v0.5. It could help users browse repeated scans per repository, but FR-502 and the proposed `GET /api/v1/scans` page scan entries; no stable project-group identity or bounded group/child-pagination contract is approved. Grouping risks splitting a project’s scans across pages and adds nested-table keyboard/screen-reader complexity. Revisit only with user evidence and an explicit data/API contract update.
- The detail proposal places operation outcome before report completeness and keeps workspaces, findings, provenance, and diagnostics distinct. The cancelled/complete example is intentionally adversarial: `complete` describes the returned report, not successful completion of the operation.
- The state board distinguishes loading, empty history, store unavailable, incompatible/corrupt storage, missing or mismatched raw artifacts, and stale selection. Retry is shown only for a safe read; no destructive recovery control is proposed.
- The mockups use system UI fonts and local vector shapes only; no remote fonts, icons, or images are required. The shadcn/ui direction remains a visual reference, not a dependency or component-source decision.
- Repository-controlled strings must render as inert text. Review cases include a label such as `<script>alert(1)</script>` and a long path; text wraps or truncates without changing meaning or causing horizontal page overflow.

Exact color tokens and dependency/build decisions require owner and independent UX/security review. Contrast, keyboard, screen-reader, responsive, and reduced-motion acceptance remain to be verified on the eventual implementation surface.


## 3. Information architecture (proposal)

1. **Scan History** — primary route; paginated saved scans, status and summary filters only if justified by user research.
2. **Scan Detail** — selected scan metadata, completeness/diagnostics, finding summary, workspace list, findings, and provenance.
3. **About/Local data help** — data location and operational guidance only if the owner accepts the UX scope; no account/team settings.

Navigation must include a persistent product identity and a clear return-to-history path. Do not show controls for unsupported team, publish, exception-management, agent-write, or remediation capabilities.

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
| Failed execution/report | Explicitly distinguish execution failure from report completeness; never “clean.” |
| Cancelled execution | Show Cancelled as the operation outcome, separately from any returned report completeness; never treat it as a completed scan. |
| API/storage unavailable | Error state distinct from empty; retry/action only when safe. |
| Migration/storage corruption | Explain history availability and recovery path without suggesting destructive reset. |
| Missing artifact | Preserve report metadata; disclose evidence unavailable/integrity failure. |
| Not found/stale selection | User-visible not-found state; return to history. |

Text labels and aria status announcements convey state. Color alone is never sufficient.

## 6. Component mapping (illustrative)

Final components follow the design-system decision. Candidate primitives include: page header, breadcrumbs/back link, status badge, summary card, table/list, pagination, filter control, tabs only if they improve scan detail comprehension, alert/banner, empty state, skeleton, tooltip, and a finding detail panel. Prefer semantic HTML/table patterns for tabular data. Dialogs are not the default for primary scan detail navigation because they impair deep linking and keyboard/screen-reader orientation.

Use shadcn/ui-style composition without adopting a component just to match a screenshot. Avoid heavy charting for v0.5 unless reviewed user research establishes a specific need; trend analytics are not part of the current proposal.

## 7. Accessibility and responsive acceptance

- All navigation, filters, pagination, and finding actions are reachable by keyboard and have a logical focus order.
- Visible focus indicator persists in all themes and states; focus is restored predictably after navigation/close.
- Headings and landmarks express the page structure; controls have programmatic labels; tables identify column headers.
- Dynamic loading/error/status updates are announced without repeatedly interrupting assistive technology.
- Text, icons, charts, and focus states meet the agreed contrast threshold; status is not color-only.
- At 320 CSS px width, essential status, identity, finding information, and navigation remain available without page-level horizontal scrolling; wide tables use an intentional responsive pattern.
- Reduced-motion preferences are respected; animation is nonessential.
- Test screen-reader names/status in a supported browser/OS combination and document scope and remaining gaps.

## 8. UX evidence package required before implementation

The design reviewer must attach or link:

- User/workflow assumptions and scope decisions.
- Low- or high-fidelity page wireframes/prototype for history, detail, empty and error states; no code scaffold is needed.
- Route/navigation and browser back/forward behavior.
- Interaction/state inventory mapped to API response and failure cases.
- Token/component proposal with shadcn/ui adoption/dependency and code ownership decision.
- Responsive breakpoints and keyboard/accessibility review.
- Representative untrusted-data rendering examples.
- Owner and independent reviewer acceptance, separately recorded.

These linked PNG/SVG mockups are first-pass design drafts; they do not complete the UX evidence package until the response/state map, component and dependency decisions, accessibility review, owner acceptance, and independent review are recorded.
