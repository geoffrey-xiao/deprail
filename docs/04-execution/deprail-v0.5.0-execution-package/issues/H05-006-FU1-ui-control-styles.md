# H05-006-FU1: Align Local Console Controls and Link Affordances

GitHub issue: [#446](https://github.com/geoffrey-xiao/deprail/issues/446).

## Planning metadata

- Parent: [H05-006 / #425](https://github.com/geoffrey-xiao/deprail/issues/425).
- Type: `bug`; Area: `cli` (existing local-console label convention); Priority: `P1`; Risk: `R1`.
- Target version/milestone: `v0.5.0`; Sprint: 4.
- Owner/reviewer: `@geoffrey-xiao`; independent review optional.
- Dependencies: merged H05-006 UI and H05-007 embedding. No dependency on completed H05-008 release acceptance for this bounded styling follow-up.
- Contracts: [H05-006](H05-006-accessible-history-ui.md), [UX design](../UX-DESIGN.md), [compatibility matrix](../requirements/COMPATIBILITY-MATRIX.md). Preserves product/architecture boundaries and v0.5 release-plan local-history scope.
- User inclusion decision: requested button alignment and link styling on the actual local history screenshot, 2026-10-07.

## Definition of Ready

- [x] User screenshot shows Refresh vertically displaced by the generic `.quiet` margin, inconsistent navigation links and weak history-row link affordance.
- [x] Scope, exclusions, observable criteria and evidence below are explicit.
- [x] Owner/reviewer, version, dependencies and low-risk presentation-only boundary are recorded.

## Goal and Scope

Align pagination/refresh buttons; use consistent semantic navigation links with active, hover and keyboard-focus states; make existing whole-row history links visibly discoverable without adding nested controls. Keep collection paging controls in aligned wrapping groups, including narrow viewports. Reuse existing React/CSS primitives; no dependencies or icons/fonts fetched at runtime.

## Out of Scope

No API/auth/storage/CLI semantics, sorting, data/schema, labels for scan completeness, capture, scanner behavior, empty-verifier panic fix, remediation cancellation fix, release publication or acceptance of unresolved browser/AT/supply-chain risks. No redesign of existing status/report meaning.

## Inputs, Outputs and Failure Behavior

Input: existing typed history/detail/about/error/recovery UI. Output: CSS/layout/link-affordance improvements only. Preserve disabled/loading behavior, safe navigation, native anchors, text-only hostile data and visible failure states. User data is unchanged.

## Required Verification

Pinned frontend build/typecheck and actual embedded candidate binary against real saved fixture histories in an isolated copy. Inspect desktop and 320px layout, control geometry, keyboard focus, row navigation, Back/Forward, about links, error/recovery and reduced motion. Capture sanitized screenshots. Automation does not prove physical 200% zoom or supported screen-reader combinations.

## Acceptance Criteria

- [ ] Pagination and Refresh controls share height/alignment; collection paging groups have consistent spacing without offset margins.
- [ ] Navigation links have uniform targets and unmistakable current/hover/focus states; history links include a visible details affordance without nested interactive elements.
- [ ] At 320px controls wrap without horizontal overflow, long labels/IDs wrap, and focus remains visible despite history-panel clipping.
- [ ] Actual packaged list/detail/about and keyboard navigation work against real history data; no API/auth/data meaning changes.
- [ ] Evidence records exact candidate identity, browser/viewport, screenshots and limitations; owner review remains required.

## Human Review and Evidence

Owner reviews screenshots and rendered behavior, separately from unresolved H05-006/H05-008 release decisions. [Actual candidate evidence and screenshots](../tracking/H05-006-FU1-STYLE-EVIDENCE.md) record desktop/320px geometry, keyboard navigation, real detail/error/recovery and the remaining platform/paging limits. Link the PR with `Refs`, not automatic issue closure.

## Final Acceptance

- [ ] Owner reviewed all criteria and remaining risks.
- [ ] CI and runtime evidence accepted.
- [ ] Reviewed merge linked; issue closure/project Done permitted only afterward.

## Rollback

Revert presentation changes and rebuild the matching embedded frontend/binary. No history migration or user-data cleanup.
