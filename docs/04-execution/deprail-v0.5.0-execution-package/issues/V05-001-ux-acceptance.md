# V05-001: Finalize History and Scan-Detail UX Acceptance
- GitHub Issue: [#394](https://github.com/geoffrey-xiao/deprail/issues/394)

- Epic: [EPIC-001 / #390](https://github.com/geoffrey-xiao/deprail/issues/390)
- Target: `v0.5.0`
- Status: Review; initial mockups merged via [PR #398](https://github.com/geoffrey-xiao/deprail/pull/398); follow-up UX proposals are in [PR #399](https://github.com/geoffrey-xiao/deprail/pull/399); owner acceptance and independent review remain pending.
- Type: docs
- Area: docs
- Priority: P0
- Risk: R2
- Owner: `@geoffrey-xiao`
- Reviewer: `@geoffrey-xiao` for planning-ticket oversight only, by explicit owner direction; not independent UX/security approval
- Dependencies: V05-000 / #387

## Value

Users need to locate retained scan operations and interpret operation outcome, report completeness, findings, workspace context, and provenance without mistaking incomplete or cancelled work for a clean completed scan.

## Scope

Complete `UX-DESIGN.md` with the accepted history/detail navigation and back/forward behavior; page-level wireframes or prototype for history, detail, empty, loading, unavailable, partial, failed, cancelled, missing-artifact, and stale-selection states; mapping from API responses/errors to UI; shadcn/ui-compatible component/dependency, code-ownership, icon/font, token, and embedded-build decisions; responsive, keyboard, screen-reader, contrast, reduced-motion, and hostile-data acceptance; and owner decision on whether About/local-data help is included.

## Out of scope

No React implementation, component installation, generated source, API/schema changes, scan-triggering controls, deletion/annotation, team settings, or remote asset/font/CDN dependency.

## Inputs, outputs, and failure behavior

Inputs: PRD, architecture, UX draft, FR-501–FR-511, API/error proposals, and privacy/security constraints. Outputs: reviewed-workflow-ready design evidence linked from `UX-DESIGN.md`. Ambiguous API or persistence behavior is recorded as a dependency/open question, not invented in the mockup. Design examples distinguish empty history from unavailable/corrupt data and preserve cancellation separately from report completeness.

## Required verification and evidence

Attach workflow assumptions, history/detail/empty/error wireframes or prototype, route/back-forward behavior, response-to-state map, component/dependency decision with license/build ownership, accessibility/responsive checklist, and representative untrusted repository values. Verify that all required status states are represented without color-only signaling.

- Link PNG previews and editable SVG sources for desktop/narrow history, scan detail, and the required loading, empty, and failure states from `UX-DESIGN.md`.

## Acceptance criteria

- A user can follow the history-to-detail workflow, identify unique history entries, and return using browser navigation.
- Every required UI state in `UX-DESIGN.md` is represented; empty history is distinct from loading and errors.
- Operation outcome and report completeness are visually and semantically distinct; partial/cancelled reports are never shown as clean completed runs.
- Keyboard, focus, screen-reader, contrast, narrow viewport, reduced-motion, and hostile-text acceptance is observable.
- The chosen UI primitives, dependencies, code ownership, and build direction are recorded without implying a dependency is already approved for runtime use.
- Independent review and runtime authorization remain outstanding until the v0.5 Definition of Ready passes.
