# H05-006-FU2 Reference-Aligned Local Console Layout

Issue: [#448](https://github.com/geoffrey-xiao/deprail/issues/448); [contract](../issues/H05-006-FU2-reference-layout.md). Date: 2026-10-07. Owner/reviewer `@geoffrey-xiao`; acceptance pending. Baseline synchronized main `28e832a` includes owner-merged #447. This is bounded presentation work, not H05-008 release acceptance or disposition of the previously recorded Windows cancellation/nil-verifier defects.

## Reference reconciliation

Compared the versioned history-desktop, history-mobile and scan-detail-desktop proposals linked by [UX-DESIGN](../UX-DESIGN.md). Their synthetic examples do not override the accepted API or semantic paginated-list contract.

| Reference intent | Implemented presentation | Deliberately not copied |
| --- | --- | --- |
| Quiet local-only/read-only context | Persistent local-only badge; history-stays-local notice; existing read-only actions only | Synthetic illustrative-data labels; hosted/scan/mutation controls |
| Compact comparable history | Desktop CSS-grid rows in a native list/anchor; separate operation, report, findings/workspace and recorded columns; panel-level Refresh | Expandable/grouped repository table, filters, total history count or started timestamp unavailable from API |
| Narrow cards with labeled statuses | Two-column property cards below 1024px; complete UUID/source identity wrapping; labels visible at narrow width | Truncated fictional entry IDs or omitted provenance |
| Findings-first desktop detail | Two explicit outcome/report cards; findings/workspaces/recorded metric cards; main findings/diagnostics with supporting workspace/record/provenance/artifact column | Invented severity summary, view-all controls or workspace-detail routes |
| Less repetitive technical text | Native keyboard-operated finding disclosure retains exact PURL, aliases and stable key; vulnerability/workspace/fixed version remain visible | Dropped evidence, guessed versions or remote icons/fonts |

Files: `web/src/main.tsx` and `web/src/style.css`. Only presentation and safe existing refresh-action placement changed; empty history now has a usable Refresh action matching its guidance. API, capture, data, auth, scanner behavior and dependencies unchanged.

## Exact local candidate

Go1.27.1 Darwin/arm64; Node22.23.3/npm10.9.9 from the existing pinned local toolchain. Commands:

```text
make frontend
go build -trimpath -o local_test/0.5.0-preview.1/output/ui-layout-448/bin/deprail ./cmd/deprail
```

Both passed; frontend includes TypeScript checking. Vite29 modules,427ms phase; rounded output CSS12.78kB/gzip3.32kB, JS213.18kB/gzip66.01kB, HTML0.56kB/gzip0.31kB. Binary SHA-256: `42f4f4663cf7de8dd992275f52a37d5f7d4347d31f85ac42abe3fe173b87a712`. This is a development build, not an attested preview release.

Actual embedded binary/API ran in an isolated HOME with a copy of the previous standalone smoke snapshot, not the user's active DB. Initial four real JS/Python/Java capture entries remained separate, with five npm findings in detail. The existing JS artifact directory was supplied explicitly; no nil-verifier workaround is represented as a fix. A private launcher handoff bootstrapped managed Chromium; it and the credential file were removed after browser closure. Both owned foreground consoles stopped with exit0; no user console was killed.

## Observed actual-surface checks

Browser: managed headless Chrome/150.0.7871.24 on Darwin/arm64; device scale1.25. Viewports measured in CSS pixels.

| Scenario | Observation |
| --- | --- |
| Desktop1440×1050 | Four native row links, aligned columns, independent Completed/Complete badges, exact UUID/source/count/time fields. Actual row height110.484375px; document scrollWidth1440. Bootstrap fragment scrubbed. |
| Keyboard selection | Tab from auto-focused heading reaches Refresh then first row; row focus3px with -3px inset. Enter navigates to actual npm detail; browser Back returns to list. |
| Detail hierarchy | Main headings Findings then Diagnostics; supporting headings Workspaces, Record information, Provenance, Artifact integrity. Five actual findings; two status cards remain independent. |
| Disclosure | Focus summary then Enter opens first native details; PURL, aliases and complete stable key visible. First key `b1e7ca3ec48c21d2ad1c43492f3f46ad4a6d639afc9cfd1d3e6d3e106a7e4d99`. No evidence field removed. |
| Narrow320×900 | History labels static/visible; refresh44px; Previous/Next both62px at the same y; document width/scrollWidth320. Detail becomes one288px column, disclosure stays accessible, no page overflow. Reduced-motion preference emulated. |
| Medium768×1050 / mobile390×1000 | Medium history uses two341px property columns; document width/scrollWidth768.390px screenshot shows readable cards and navigation. |
| Real offline failure + refresh | Actual scan of a new isolated npm fixture with PATH=/usr/bin:/bin, --save-history, returned3. Refresh added the saved occurrence. UI correctly showed Completed operation / Failed report, zero recorded findings without claiming clean, and SCANNER_NOT_FOUND diagnostic. The persistent failed-report warning and non-clean empty-collection message remained. No mock API response or vulnerability generated. |
| Empty profile | Separate fresh HOME, no DB. Empty UI rendered No saved scans; Refresh stayed usable and empty; no history paths created; no320px overflow. |
| Error/about/recovery | Actual absent UUID showed HISTORY_ENTRY_NOT_FOUND; About active state and no narrow overflow; reload showed existing auth-recovery page. Browser JS errors empty. |

Screenshots of the initial real four-entry snapshot:

- [Desktop history](H05-006-FU2-history-desktop.webp)
- [390px history](H05-006-FU2-history-mobile.webp)
- [320px history](H05-006-FU2-history-320.webp)
- [Desktop detail](H05-006-FU2-detail-desktop.webp)
- [320px detail with evidence expanded](H05-006-FU2-detail-320.webp)

No source-text/wiring tests added for layout; actual packaged geometry, native disclosure, keyboard navigation and real API state checks are the verification. CI outcome is separate and must reference its exact head/run.

## Remaining limits and owner review

No physical200% browser zoom, Safari/VoiceOver, Chrome/NVDA or Firefox/Orca proof. This small data set has no next cursor; multi-page transitions and simultaneous collection paging actions are not claimed exercised. No report-start time/total-entry count invented. Reference drawings are interpreted under current contracts, not pixel-identical implementation mandates. No new backup/recovery/security/release gate claimed. Prior #446 acceptance remains separate even though #447 merged.

Owner reviews actual screenshots, data-field preservation, keyboard disclosure and residual platform/CI limits. PR uses `Refs #448`; no final acceptance boxes or Master Checklist changes until linked owner-reviewed evidence. Revert CSS/markup and rebuild matching assets/binary to roll back; no migration or data deletion.
