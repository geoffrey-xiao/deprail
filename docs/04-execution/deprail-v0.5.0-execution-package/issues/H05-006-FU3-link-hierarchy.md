# H05-006-FU3: Refine Local Console Link Hierarchy

GitHub issue: [#450](https://github.com/geoffrey-xiao/deprail/issues/450).

## Planning metadata and Definition of Ready

- Parent [H05-006 / #425](https://github.com/geoffrey-xiao/deprail/issues/425); follows owner-merged [#449](https://github.com/geoffrey-xiao/deprail/pull/449), synchronized main `256cd54`.
- Type `bug`; area `cli`; risk `R1`; priority `P1`; milestone/target `v0.5.0`; Sprint4.
- Owner/reviewer `@geoffrey-xiao`; independent review optional.
- User requests further link polish after merging layout PR,2026-10-07.
- [x] Value, scope, observable evidence, failure boundary, dependencies and reviewer explicit below; no new capability or dependency.
- Contracts: [UX-DESIGN](../UX-DESIGN.md), [H05-006](H05-006-accessible-history-ui.md), [compatibility](../requirements/COMPATIBILITY-MATRIX.md). Presentation inside existing release scope, not new release authority.

## Goal and Scope

Reduce duplicate underlined link cues in a single history anchor; distinguish row detail affordance from metadata. Keep navigation tab styling stable on hover. Unify breadcrumb/footer link target, hover and focus affordances. Retain native anchors, exact routes, one anchor per row and intact keyboard/modified-click semantics. Keep ordinary inline text links underlined.

## Exclusions and Failure Behavior

No external hyperlinks, new routes/actions, disclosure behavior, API/auth/data/storage/scanner/dependency changes, nil-verifier/cancellation fixes or release acceptance. Errors/empty/unavailable, completeness and identifier data unchanged. No nested buttons inside anchors, no fake hrefs, no focus suppression or history mutation.

## Acceptance and Required Evidence

- [ ] Navigation hover/current states are consistent without competing underlines; brand remains a working home link.
- [ ] History row has clear title/arrow and one details affordance; hover/keyboard focus distinguish the entire row without duplicate link targets.
- [ ] Breadcrumb/footer text links have readable hover/focus and adequate targets; ordinary body links retain conventional underlines.
- [ ] Actual packaged desktop/320px list/detail/about, Tab/Enter/Back and native link destinations verified; sanitized screenshots, candidate identity and limitations recorded.

## Human Review and Final Acceptance

Owner reviews screenshot/keyboard evidence separately from prior UI/release/platform gates. Use Refs to the new issue; do not close prior issues simply because PRs merged. Owner should close prior #448 only after its criteria and required evidence are accepted.

[Actual link/keyboard evidence and screenshots](../tracking/H05-006-FU3-LINK-EVIDENCE.md) record the packaged candidate identity, hover/focus geometry, real navigation,320px layout and unexercised platform/modified-click limits.

- [ ] Owner criteria and remaining risks accepted.
- [ ] CI and exact candidate evidence linked.
- [ ] Owner-reviewed merge recorded before Project Done/closure.

## Rollback

Revert styling and rebuild matching embedded binary. No user-data migration/deletion.
