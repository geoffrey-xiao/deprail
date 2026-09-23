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

## 3. Information architecture (proposal)

1. **Scan History** — primary route; paginated saved scans, status and summary filters only if justified by user research.
2. **Scan Detail** — selected scan metadata, completeness/diagnostics, finding summary, workspace list, findings, and provenance.
3. **About/Local data help** — data location and operational guidance only if the owner accepts the UX scope; no account/team settings.

Navigation must include a persistent product identity and a clear return-to-history path. Do not show controls for unsupported team, publish, exception-management, agent-write, or remediation capabilities.

## 4. Primary workflows

### Review history

1. Open the local console.
2. See a paginated ordered list, or an explicit empty-history state.
3. Identify a scan using its stable identity, timestamp, root label approved by the data contract, status, finding count, and concise source/tool summary.
4. Select a row/card to open details; preserve the selected scan in browser navigation/back behavior.

### Inspect scan detail

1. See scan outcome at top: Complete, Partial, Failed, or Cancelled; do not infer status from finding count.
2. See stable scan identity, completion/start time if available, repository/workspace summary, scanner/tool/database provenance, and diagnostics.
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
| Complete, zero findings | Clearly complete scan, zero findings; no ambiguity. |
| Complete, findings | Finding count/severity summary and details navigation. |
| Partial | Persistent warning label; explain omitted/incomplete scope and show diagnostics. |
| Failed/cancelled | Explicit outcome and cause/reference where safe; never “clean.” |
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

No design files, screenshots, or prototypes are asserted complete by this contract.
