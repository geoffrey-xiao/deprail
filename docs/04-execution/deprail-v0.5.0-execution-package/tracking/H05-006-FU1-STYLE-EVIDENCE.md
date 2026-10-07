# H05-006-FU1 Console Control and Link Styling Evidence

Issue: [#446](https://github.com/geoffrey-xiao/deprail/issues/446); [contract](../issues/H05-006-FU1-ui-control-styles.md). Date: 2026-10-07. Owner/reviewer: `@geoffrey-xiao`; acceptance pending. Presentation-only follow-up to H05-006; no H05-008 release gate completed.

Review PR: [#447](https://github.com/geoffrey-xiao/deprail/pull/447), implementation commit `44e6fd8`. Linked with `Refs #446`; owner acceptance and CI outcome remain separate from the observed local smoke below.

## Change

- `web/src/style.css`: 44px minimum controls, no `.quiet` offset margin; aligned action groups with a separate responsive Refresh row; uniform navigation targets, hover and current states; inline breadcrumb styling; history hover and inset keyboard outline; wrapping footer and narrow title/actions.
- `web/src/main.tsx`: group existing paging controls, add a noninteractive `View scan details →` affordance inside each existing history anchor, group collection paging actions. No nested interactive elements or new action semantics.
- No API, authentication, database, capture, scanner, completeness, permissions, dependency or lock changes. The separately identified nil-artifact-verifier bug and Windows cancellation E2E failure remain unresolved and out of scope.

## Candidate identity and setup

Reviewed baseline: main `96a6f09`, branch `fix/h05-006-ui-control-styles`. Local candidate includes this issue's frontend changes. Go 1.27.1 on Darwin/arm64; Node 22.23.3/npm 10.9.9 from the existing pinned local toolchain. `make frontend` passed TypeScript checking and Vite build: 29 modules, 427ms Vite phase, CSS 9.56kB/gzip2.77kB, JS210.66kB/gzip65.66kB (rounded build output).

Command: `go build -trimpath -o local_test/0.5.0-preview.1/output/ui-style-446/bin/deprail ./cmd/deprail`.

Binary SHA-256: `94529403834d900e0a106c53b474c37195fe19daca4294dbfbaf053dba5181e4`. Development identity is not an attested release identity or preview tag.

Runtime: actual embedded binary and API, ephemeral loopback port, isolated HOME under `local_test/0.5.0-preview.1/output/ui-style-446/home`. A separate copy of the operator's fixture-only history supplied four real saved scans (npm twice, Python and Java). Source fixtures/history were not reset or repaired. Snapshot preparation's Apple SQLite CLI returned `unable to open database file (14)` after creating the destination; subsequent immutable standalone `PRAGMA quick_check` returned `ok`, and the actual Go reader reopened and rendered it. This is smoke-data preparation, not proof of the release backup/recovery gate. A private launcher handoff supplied the process bootstrap to managed Chromium; credential file and launcher were deleted and the owned console stopped with exit 0. The user's existing console was not stopped by the assistant.

## Observed browser checks

Browser: managed headless Chrome/150.0.7871.24, Darwin/arm64. Desktop1440×1100 and mobile320×900 CSS pixels. Screenshots have a 1.25 device scale; viewport values below are CSS pixels.

| Scenario | Observed result |
| --- | --- |
| Real list / desktop | Four real saved entries rendered with unchanged operation/report/count metadata. Previous, Next and Refresh each measured44px high and y=1091.375. Refresh is no longer displaced. Bootstrap fragment was scrubbed. |
| Link affordance | Every row remains one native anchor and includes visible underlined details text/arrow. About hover measured primary blue, underline, light blue background, 44px target. Current History/About states selected on their respective routes. |
| Keyboard detail | Tab from the auto-focused history heading reached the first row; solid3px outline with -3px inset remains inside the clipped panel. Enter navigated to npm detail with five findings; Back returned to list. |
| 320px list | document scrollWidth320, no horizontal overflow. Previous/Next both112px wide,62px high at the same y=1504.46875; Refresh288px wide/44px high on the next row. Long UUIDs and Python title wrapped. Header links remain44px high. |
| 320px detail/about | Both rendered without horizontal overflow; detail Refresh measured44px. About's active link was About. |
| Real not-found response | Navigation to a valid absent UUID produced `HISTORY_ENTRY_NOT_FOUND`; retry button44px, no horizontal overflow. No response mock was used. |
| Reload / auth recovery | Reload rendered the existing recovery page, not empty history; no browser JavaScript errors were reported. |
| Reduced motion | Browser emulation confirmed prefers-reduced-motion active at the narrow viewport. No new animation added. |

Screenshots:

- [Desktop history](H05-006-FU1-history-desktop.webp)
- [320px history](H05-006-FU1-history-mobile.webp)
- [Desktop detail](H05-006-FU1-detail-desktop.webp)

No permanent source-text/CSS-wiring tests were introduced: rendered geometry and actual keyboard/navigation smoke are the appropriate evidence. Frontend typecheck/build and binary compilation passed; CI result must be linked separately, never inferred from prior PRs.

## Limits and review

Four entries produce no next history cursor; enabled next-page transitions and simultaneous collection Next/First controls were not exercised in this fixture smoke. Their existing handlers/conditions are retained; shared wrapping layout is implemented, but this record does not claim a paginated-data runtime pass. No physical200% browser zoom, Safari/VoiceOver, Chrome/NVDA or Firefox/Orca validation was performed. Those remain the existing owner/platform release gates. Screenshot appearance is evidence for this bounded alignment request, not complete accessibility/security approval.

Owner reviews the visible desktop/mobile changes and residual limits. Use `Refs #446`; do not close the issue or mark it Done until owner review, required CI and merge evidence are complete. Rollback reverts CSS/markup and rebuilds the matching embedded binary, without data migration.
