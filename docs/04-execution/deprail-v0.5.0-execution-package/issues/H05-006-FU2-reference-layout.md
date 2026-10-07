# H05-006-FU2: Align Local Console Information Layout with UX References

GitHub issue: [#448](https://github.com/geoffrey-xiao/deprail/issues/448).

## Planning metadata

- Parent: [H05-006 / #425](https://github.com/geoffrey-xiao/deprail/issues/425); follows owner-merged [#447](https://github.com/geoffrey-xiao/deprail/pull/447).
- Type `bug`; area `cli`; priority `P1`; risk `R1`; milestone/target `v0.5.0`; Sprint 4.
- Owner/reviewer: `@geoffrey-xiao`; independent review optional.
- Baseline: synchronized main `28e832a`. User requested comparison with the existing UI design and further presentation improvement, 2026-10-07.
- Contracts: [UX design](../UX-DESIGN.md), [H05-006](H05-006-accessible-history-ui.md), [compatibility](../requirements/COMPATIBILITY-MATRIX.md). This stays inside accepted local-history product/architecture/release scope; reference drawings do not authorize new API/data capabilities.

## Definition of Ready

- [x] User value: denser comparable history, findings-first detail and clear local/read-only context.
- [x] Included/excluded scope, failure semantics, observable evidence, owner and version explicit below.
- [x] Existing merged UI/assets available; no new dependency or unresolved release-gate acceptance needed for this R1 presentation-only follow-up.

## Scope

Desktop history uses a compact column-aligned semantic list, not grouped/expandable repository tables. Separate operation outcome from report completeness with labeled cells; keep UUID, source identity, workspace count and recorded time truthful and available. Narrow history becomes readable cards. Add local-only/read-only context and a panel-level refresh action. Detail prioritizes findings, with workspace/provenance/record information in a supporting column and secondary finding identifiers disclosed through native details. Status warnings, diagnostics and artifact-integrity limitations remain visible; exact stored fields are never dropped.

## Exclusions

No search/filter/grouping, totals/start-time/severity statistics unsupported by API, new routes, API/storage/auth/scanner changes, dependencies/fonts/CDN, raw-artifact access, scan/delete/export actions, nil-verifier or Windows cancellation fixes, release approval/publication. No implication that prior UI/CI gaps are accepted by a merge.

## Failure Behavior

Keep empty/unavailable/error states distinct, report/operation independent, missing artifact warnings persistent, and hostile content inert. Disclosure is local presentation only; no additional request or mutation. Native anchors, keyboard/focus, Back/Forward and memory-only bootstrap remain.

## Verification and Acceptance

- [ ] Real history displays separate labeled operation/report states, exact identity/count/time fields, compact desktop columns and narrow cards.
- [ ] Detail findings precede auxiliary metadata; all original finding/provenance/artifact data is visible or accessible through keyboard-operated disclosure, without invented summary data.
- [ ] Refresh/pagination/navigation operate on the actual candidate; 320px and medium widths have no page overflow, focus and labels remain readable.
- [ ] Real packaged list/detail/about/not-found/auth-recovery exercised; sanitized desktop/mobile screenshots and exact candidate identity recorded.
- [ ] Owner accepts remaining browser/AT/CI limits separately; release acceptance and publication not inferred.

## Evidence and Human Review

[Reference comparison, actual candidate results and screenshots](../tracking/H05-006-FU2-LAYOUT-EVIDENCE.md) record desktop/medium/mobile layout, keyboard disclosure and real offline-failed/empty/error/recovery cases. Owner reviews every criterion and residual risks; use `Refs`, not automatic closure until accepted.

## Final Acceptance

- [ ] Owner criteria and remaining risks reviewed.
- [ ] Required CI/evidence reviewed and linked.
- [ ] Owner-reviewed merge recorded before Project Done/closure.

## Rollback

Revert presentation markup/styles, rebuild matching assets/binary. Preserve user history and existing contracts.
