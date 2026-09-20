# DepRail v0.3.1 Product Requirements

| Attribute | Value |
| --- | --- |
| Version | v0.3.1 |
| Status | Planning baseline |
| Outcome | Clear, safe, cross-platform CLI presentation without machine-contract regression |
| Predecessor | v0.3 remediation-planning release |
| Owner | Project owner |
| Reviewer | Architecture/security reviewer |

## Goal

Make the DepRail CLI explain ongoing work and final outcomes clearly for humans while preserving deterministic machine output, existing scan semantics, security boundaries, and cross-platform meaning.

## User jobs

- Understand which scan phase is active without guessing whether the process is idle.
- Read consistent command summaries, statuses, errors, tables, and next steps.
- Distinguish complete zero findings from partial and failed scans.
- Use interactive terminals, CI logs, redirected output, and narrow terminals safely.
- Consume JSON without progress, banners, ANSI escapes, or incidental text.
- Trust repository-derived labels not to control the terminal.

## Scope

- DepRail-owned structured presentation events and renderer boundary.
- Shared human presentation primitives across `doctor`, `discover`, `scan`, and `fix plan`.
- TTY capability detection, color policy, width handling, and non-TTY behavior.
- Progress driven only by real application lifecycle events.
- Explicit complete, partial, failed, and cancelled summaries.
- Preservation and verification of `--quiet`, `--verbose`, JSON, stdout/stderr, and exit contracts.
- Safe rendering of hostile repository, scanner, and error labels.
- Cross-platform terminal and redirected-output evidence.

## Out of scope

Scanner changes, new ecosystems, repository mutation, package installation, network services, telemetry, source upload, schema-breaking changes, autonomous approval, full-screen TUI, artificial progress, and v0.4 isolation/verification behavior.

## Product invariants

- Human presentation cannot alter domain results or scanner invocation semantics.
- JSON stdout contains machine data only.
- Progress and diagnostics use stderr.
- Interactive progress is TTY-only; non-TTY output is stable and line-oriented.
- Only complete zero-finding scans may say “No known vulnerabilities found.”
- Unknown or incomplete data never becomes a safe or successful result.
- Untrusted labels are escaped before terminal rendering.
- Local-first operation and no source upload remain unchanged.

## Acceptance

Representative evidence demonstrates useful TTY progress, consistent summaries, truthful complete/partial/failed states, preserved quiet/verbose behavior, JSON purity, safe hostile-label rendering, narrow and monochrome readability, cancellation cleanup, and equivalent Linux/macOS/Windows semantics without repository mutation.
