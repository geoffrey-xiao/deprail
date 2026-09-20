# v0.3.1 Functional Requirements

| ID | Requirement | Priority |
| --- | --- | --- |
| UX-001 | Human output uses shared headers, sections, statuses, summaries, and error primitives. | Must |
| UX-002 | Application lifecycle events are structured and renderer-independent. | Must |
| UX-003 | Interactive progress is emitted only for real work and only on TTY stderr. | Must |
| UX-004 | Non-TTY output is stable, line-oriented, and non-animated. | Must |
| UX-005 | Complete zero findings, partial results, failures, and cancellation are visibly distinct. | Must |
| UX-006 | `--quiet` and `--verbose` preserve existing documented behavior. | Must |
| UX-007 | JSON stdout contains no progress, banners, diagnostics, or ANSI escapes. | Must |
| UX-008 | Repository-derived labels cannot inject terminal control sequences. | Must |
| UX-009 | Narrow and monochrome output remains readable. | Must |
| UX-010 | Linux, macOS, and Windows preserve equivalent result meaning. | Must |
| UX-011 | Renderer changes do not mutate repositories or change scanner invocation. | Must |
| UX-012 | Terminal state is restored after cancellation and failure. | Must |

## Non-functional requirements

- No artificial delays or invented progress percentages.
- No secrets, credential-bearing URLs, or sensitive environment values in user output.
- Deterministic machine output and stable error codes remain mandatory.
- Terminal dependencies remain behind presenter boundaries.
