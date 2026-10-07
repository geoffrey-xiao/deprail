# H05-006-FU3 Link Hierarchy Evidence

Issue [#450](https://github.com/geoffrey-xiao/deprail/issues/450); [contract](../issues/H05-006-FU3-link-hierarchy.md). Date2026-10-07; baseline synchronized main256cd54 after owner-merged #449. Owner/reviewer @geoffrey-xiao; acceptance pending.

Review PR: [#451](https://github.com/geoffrey-xiao/deprail/pull/451), implementation commit `34b3ff9`, linked with `Refs #450`. Exact-head CI and owner acceptance are recorded separately; local screenshot proof does not imply either.

## Presentation change

Only `web/src/style.css` changed runtime behavior. Navigation retains tab/current styling without hover underlines. History titles are neutral by default with blue arrows; hovered/focused titles become blue/underlined. View details is a44px bordered visual affordance inside the same existing native anchor, not another button/link/tab stop. Breadcrumb/footer use44px padded link targets and hover/focus backgrounds/underlines. Ordinary inline anchor defaults and global3px focus remain. API, routes, markup, authentication, data, disclosure, scanner, dependencies and existing handlers unchanged; no external links introduced.

## Candidate and verification

Go1.27.1 Darwin/arm64; pinned Node22.23.3/npm10.9.9. `make frontend` passed TypeScript and Vite build:29 modules,427ms Vite phase; rounded CSS13.52kB/gzip3.42kB, JS213.18kB/gzip66.01kB. `go build -trimpath -o local_test/0.5.0-preview.1/output/ui-links-450/bin/deprail ./cmd/deprail` passed. Binary SHA256 `74381bcbd451a05fec809771ed7e89b60fb1cdba0be2b10b5622b05dbd5a2fc1`. Development identity, not a preview release claim.

Actual embedded console/API used a separate copy of the stopped standalone fixture-history snapshot in an isolated HOME; the user's current database/service was not modified/stopped. Four actual JS/Python/Java captures rendered. An explicit JS artifact root avoids the separately known nil-verifier bug; that bug is not fixed here. Private bootstrap handoff and launcher were removed, managed tab closed, owned console stopped exit0.

Managed headless Chrome150.0.7871.24, Darwin/arm64; desktop1440×1050 and narrow320×900 CSS pixels, screenshot scale1.25:

| Check | Observed result |
| --- | --- |
| Default history links | Title/details decoration none, details affordance44px; zero nested anchors/buttons/tabindex elements inside each row. Bootstrap fragment scrubbed. |
| Navigation hover | About light-blue background, no underline; existing active History indicator retained. |
| Row hover | Title blue rgb(29,78,216) and underlined; row background highlights whole target. |
| Keyboard | Tab from heading via Refresh reaches row; focus3px, inset-3px; highlighted row and underlined title. Enter navigates to actual npm detail. |
| Mouse details affordance | Clicking the decorative View details span navigates to npm detail through its existing enclosing anchor. |
| Return/home/Back | Breadcrumb destination /console/history,44px target and hover underline/background; click returns to list; brand returns to list; browser Back returns from details. |
| Footer |44px target; Tab from last row skips disabled paging and reaches About local history, with3px outline and underline. Click reaches About with correct active nav. |
|320px | document width/scrollWidth320; details/footer44px, readable wrapping; About no horizontal overflow. |
| Runtime errors | No browser JavaScript errors reported. |

Screenshots: [desktop](H05-006-FU3-links-desktop.webp), [320px](H05-006-FU3-links-mobile.webp).

No permanent CSS/source-wiring tests introduced. Actual native-anchor and rendered geometry smoke is the evidence; no fake API response was used. CI must reference its exact run/head separately. Modified-click/new-tab behavior was not exercised here; handler and href semantics are unchanged. Physical200% zoom, Safari/VoiceOver, NVDA/Orca and broad release/security acceptance remain separate gates. No cancellation or nil-verifier fix inferred.

Owner reviews hover/focus, visual affordance and native-link semantics before acceptance. Use Refs #450; no issue closure/Project Done or Master Checklist checkmarks before owner-reviewed criteria, required CI and merge evidence. Rollback reverts CSS and rebuilds matching embedded assets/binary without touching user history.
