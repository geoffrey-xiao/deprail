# H05-006: Implement Accessible Local History UI
- GitHub Issue: [#425](https://github.com/geoffrey-xiao/deprail/issues/425).
- GitHub parent: [#418](https://github.com/geoffrey-xiao/deprail/issues/418). Published dependencies: #424; global pre-start gate #419.

## Planning metadata

- Epic: [EPIC-002 — Local History Delivery](../epics/EPIC-002-local-history-delivery.md)
- Type: `feature`
- Area: `foundation`
- Priority: `P0`
- Risk: `R2`
- Target version: `v0.5.0`
- Milestone: [`v0.5.0`](https://github.com/geoffrey-xiao/deprail/milestone/11)
- Sprint: Unassigned
- Owner: `@geoffrey-xiao`
- Reviewer: `@geoffrey-xiao` (owner); independent review is not a gate under ADR-0005.
- Dependencies: H05-005 route, response, auth and recovery contract; accepted UX and security contract. Packaged acceptance belongs to H05-007; no H05-007 dependency is needed for source-level smoke.
- Blocked reason: None (Todo, not started); global timeout disposition and issue-specific pre-start gates apply.

## Definition of Ready

- [x] Value/user impact supported by accepted proposal [PR #417](https://github.com/geoffrey-xiao/deprail/pull/417) and owner decision [#417](https://github.com/geoffrey-xiao/deprail/pull/417#issuecomment-5907670878).
- [x] Scope/exclusions specified by [UX design](../UX-DESIGN.md) and [API design](../API-DESIGN.md).
- [x] Inputs, state distinctions and auth recovery specified by UX/API contracts.
- [x] Actual-surface smoke scenarios are named below.
- [x] Observable source behavior specified; runtime acceptance remains unchecked.
- [x] Owner reviewer assigned; independent review is not required.
- [x] H05-005 dependency and target version recorded; implementation remains gated by prerequisite disposition.

## Goal

Deliver the complete accessible browser UI source for browsing local scan history and operation details through the read-only local API, without implying scan or mutation capabilities.

## Scope

Own the full UI source and lockfile in `web/` (including entry, routing, state, styles, tests as appropriate and exact dependency lock). Consume H05-005 typed API states; maintain operation outcome separately from report completeness, and distinguish empty, unavailable, incompatible, error, partial and missing-evidence states. Use only local tokens; remove bootstrap fragment before first request/render, retain bearer only in memory, and provide safe reopen/recovery guidance without a stored credential fallback. Provide history/detail/about/help navigation, Back/Forward, keyboard operation, visible focus, meaningful status announcements, zoom/reflow/narrow layouts, contrast and reduced-motion behavior; render all untrusted findings, labels and diagnostics as text, never executable markup. Follow [UX primitives](../UX-DESIGN.md), product [PRD](../PRD-v0.5.md) FR-502–FR-504, FR-510–FR-511, architecture [ARCHITECTURE-v0.5](../ARCHITECTURE-v0.5.md), and [SEC-06](../requirements/SECURITY-REQUIREMENTS.md). Source smoke may use a throwaway embed harness and actual accepted API; production packaging/embedded-binary gate is H05-007.

Before introducing candidate packages/lockfile, complete the required direct/transitive license, vulnerability and provenance review using the exact selected compatibility versions. Use locked, reviewed build steps with dependency lifecycle scripts disabled by default; no third-party component/font/icon/CDN runtime dependency. H05-007 reuses these review records for the full packaged binary rather than postponing pre-introduction review.

## Out of Scope

No backend, storage, CLI, API contract change, scan/mutation controls, CDN/network dependency, production embedding/packaging, shipped test assets, mock-only behavioral claims or H05-007 circular dependency.

## Inputs, Outputs, and Failure Behavior

Inputs: versioned API responses/errors, auth fragment at initial bootstrap, local UI state and embedded assets. Output: navigable history/detail/help surface, with accessible and safe distinction of completed/failed/cancelled operation, complete/partial/failed report, empty/unavailable/error, artifact-integrity states and auth loss. API timeout, incompatibility, corruption, malformed data or lost token must never appear as empty or successful history; recovery asks user to reopen from active local console process. Hostile text is inert. No token in URL after bootstrap, persistent storage, referrer or logs.

## Required Tests

- Run source UI in a throwaway embed harness against a real contract-compatible API; exercise history/detail/help, actual keyboard focus, Back/Forward and auth recovery; do not use a mock echo as proof.
- Exercise representative operation/report combinations and distinguish no report, empty report, unavailable collection, API failure and artifact-integrity states.
- Render hostile HTML/script-shaped finding and diagnostic text and verify it remains inert; verify token fragment removed before first request and credentials absent from persistent storage.
- Inspect actual surface with keyboard-only, focus visibility, narrow viewport/zoom, reduced motion and screen-reader status/navigation scenarios; record browser/AT results.

## Acceptance Criteria

- [ ] Complete history, detail and help routes consume typed API meaning and visibly distinguish all required operation/report/data/error states.
- [ ] Keyboard navigation, focus, status announcements, zoom/reflow and responsive navigation work on actual rendered surface; hostile text remains inert.
- [ ] Token bootstrap/recovery meets memory-only security contract and never renders or logs the credential URL.
- [ ] Source-level real-API smoke passes; no synthetic component/mock assertion is represented as accessibility or runtime proof. Packaged binary acceptance remains H05-007/H05-008.

## Owner Review

Owner reviews actual UI behavior and accessibility/security evidence. Planning DoR is not actual runtime or packaged acceptance. Independent reviewer is not a gate.

## Evidence Required

- Verification commands or scenarios: future locked source build; throwaway embed/browser run against real API; keyboard/AT, narrow/zoom, hostile-text, state-transition, Back/Forward and token-recovery scenarios.
- Expected artifacts, logs, screenshots, or links: browser/AT matrix, actual-surface screenshots and sanitized state/API evidence; no credential or repository-sensitive content.

## Final Acceptance

- [ ] Owner reviewed every acceptance criterion during PR review.
- [ ] Required verification and CI results were reviewed.
- [ ] Owner review and remaining risk are recorded.
- [ ] Evidence links are attached.
- [ ] Owner review and merge evidence are linked.

## Rollback

Withdraw UI source/assets without modifying stored history or changing existing CLI behavior; H05-007 may roll back to a matching reviewed binary/assets.
