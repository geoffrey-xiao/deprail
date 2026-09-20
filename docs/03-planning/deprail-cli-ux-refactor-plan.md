# DepRail CLI UX Refactor Plan

**Status:** Future planning only; not authorized for current v0.3 implementation.

**Target stage:** Post-v0.3 CLI refinement, scheduled after the remediation-planning release gate.

**Purpose:** Improve the complete DepRail CLI experience—progress, summaries, errors, styling, interaction, and terminal ergonomics—without weakening machine-readable output, deterministic behavior, cross-platform semantics, or security boundaries.

## Motivation

The current CLI correctly separates machine output from diagnostics and reports explicit complete, partial, and failed states. Long-running scans can still appear idle while discovery, OSV-Scanner execution, artifact retention, and normalization are in progress. Human users also need consistent hierarchy, status language, errors, tables, colors, and empty-result summaries across commands.

The refactor should improve the whole terminal experience, not only add a spinner.

## Recommended architecture

Use a two-layer design:

### Layer 1: DepRail-owned presentation contract

Application services emit structured events and results. They do not write terminal strings directly.

Potential events:

```text
InputValidated
WorkspaceDiscoveryStarted
WorkspaceDiscovered
ScanPlanBuilt
WorkspaceScanStarted
WorkspaceScanCompleted
ArtifactStored
NormalizationStarted
ReportReady
OperationCancelled
```

The presentation layer owns:

- Output mode.
- TTY capabilities.
- Color and animation policy.
- Progress lifecycle.
- Human summaries.
- Error and incomplete-result presentation.
- Safe rendering of untrusted labels.

This keeps domain and application behavior independent of terminal libraries and allows tests to assert structured events without comparing ANSI output.

### Layer 2: replaceable terminal renderers

Start with a small internal renderer and explicit capability detection. The renderer may use:

- [Lip Gloss](https://github.com/charmbracelet/lipgloss) selectively for styling, layout, borders, and colors.
- Selected [Bubbles](https://github.com/charmbracelet/bubbles) components for a spinner or progress element when needed.

Do not make the scan pipeline depend directly on a full-screen TUI framework.

[Bubble Tea](https://github.com/charmbracelet/bubbletea) should remain a later option for genuinely interactive workflows such as remediation review, approval, plan inspection, or multi-step terminal screens. It is not required for ordinary `scan`, `discover`, or `doctor` output.

[PTerm](https://github.com/pterm/pterm) is an alternative for a faster high-level implementation, but it should not be adopted until its global printer behavior, TTY detection, stderr routing, JSON purity, cancellation, and concurrent-output behavior are reviewed against DepRail’s contracts.

## Scope

### Whole-CLI visual system

Define consistent presentation primitives for:

- Command headers.
- Sections and subsections.
- Success, warning, error, and informational states.
- Tables and key/value summaries.
- Finding counts and severity summaries.
- Workspace status.
- Plan candidates and risks.
- Guidance and next steps.
- Stable error-code display.
- Color and monochrome output.
- Terminal width and narrow-window behavior.

Visual styling must remain useful without color and must not encode meaning only through color.

### Interactive progress

Provide concise progress messages on stderr for interactive TTY runs:

```text
Discovering workspaces...
Scanning 3 workspaces...
Scanning frontend... done
Scanning services/api... done
Scanning services/worker... done
Normalizing findings...
```

An optional spinner or loader may be used while a phase is active, but it must be a presentation detail rather than domain or scanner behavior. Progress must come from real application events; never invent percentages or claim completion before the underlying phase completes.

Progress phases should correspond to real boundaries:

1. Validate input and resolve the repository root.
2. Discover workspaces and authoritative files.
3. Build scanner targets.
4. Execute scanner targets.
5. Retain raw artifacts.
6. Normalize findings.
7. Render the final report.

### Empty-result messaging

After a `complete` scan with zero findings:

```text
No known vulnerabilities found.
```

Never use that message for `partial` or `failed` scans:

```text
Scan incomplete: 1 workspace failed.
```

```text
Scan failed: OSV-Scanner was not found.
```

### Existing flag compatibility

The v0.1 CLI contract already defines `--quiet` and `--verbose` for `deprail scan`. The refactor must preserve and reconcile those contracts; they are not optional future flags.

Recommended semantics:

- `--quiet`: suppress progress and successful human summaries while retaining errors and required diagnostics.
- `--verbose`: preserve existing behavior and expose additional safe diagnostics.
- Both flags must remain compatible with `--format json` and never contaminate machine stdout.
- Any semantic change requires a CLI contract update and compatibility evidence.

### Machine-output and TTY contract

The refactor must preserve:

- `--format json` writes JSON data only to stdout.
- Diagnostics and progress use stderr.
- Redirected stdout remains valid JSON with no banners or ANSI escapes.
- Interactive progress is enabled only when stderr is a TTY.
- Non-TTY progress is disabled by default; a future structured-progress option may explicitly opt in to machine-consumable progress.
- Non-TTY diagnostics are stable, line-oriented, and free of carriage-return animation.
- JSON output remains deterministic and unchanged except for approved schema changes.
- No progress or summary contains secrets, credential-bearing URLs, or raw sensitive environment values.

### Hostile-label rendering

Repository-derived values are untrusted input. Workspace IDs, paths, package names, error details, and scanner labels must be safely encoded before terminal rendering.

The future renderer must:

- Escape newlines, carriage returns, tabs where unsafe, and terminal control sequences.
- Prevent raw ANSI sequences from repository-controlled labels.
- Preserve enough readable context for users to identify the workspace.
- Apply the same protection to TTY and redirected stderr output.

## Non-goals

This plan does not authorize:

- Scanner behavior changes.
- New scanners or package ecosystems.
- Network services or telemetry.
- Repository mutation or remediation execution.
- Changes to JSON schemas without compatibility review.
- Hidden progress output in machine-readable streams.
- Artificial delays or activity indicators disconnected from real work.
- Replacing all command parsing or application architecture without a measured need.

## Proposed implementation slices

1. **Presentation contract**
   - Define structured events, output modes, severity, lifecycle, TTY capabilities, and compatibility rules.
2. **Terminal capability layer**
   - Detect TTY, color, width, animation support, CI/non-TTY mode, and cancellation safely.
3. **Static human renderer**
   - Implement deterministic line-oriented summaries and errors before animation.
4. **Safe styling system**
   - Add selective Lip Gloss styling with monochrome and narrow-terminal fallbacks.
5. **Progress renderer**
   - Add TTY-only progress and optional spinner/Bubbles components based on real events.
6. **Application event wiring**
   - Emit discovery, scanner, artifact, normalization, and render events from real boundaries.
7. **Empty and incomplete summaries**
   - Distinguish complete zero findings from partial and failed scans.
8. **Flag compatibility**
   - Reconcile and test existing `--quiet` and `--verbose` semantics.
9. **Hostile-label hardening**
   - Escape control characters and add adversarial presenter tests.
10. **Cross-platform behavior**
    - Verify macOS, Linux, and Windows output, cancellation, terminal cleanup, and redirected logs.
11. **Interactive TUI evaluation**
    - Evaluate Bubble Tea only after a real interactive plan/review use case is defined.
12. **Regression and release evidence**
    - Prove JSON stdout purity, stderr routing, no-animation redirection, deterministic summaries, safe labels, and no secret leakage.

## Acceptance criteria for a future issue

- Interactive scan runs provide useful phase/workspace feedback without false progress.
- The whole CLI uses consistent visual hierarchy, statuses, errors, tables, and summaries.
- Complete zero-finding scans explicitly say no known vulnerabilities were found.
- Partial and failed scans never present a safe-looking empty result.
- Existing `--quiet` and `--verbose` contracts are preserved and documented.
- JSON stdout remains parseable and free of progress text or ANSI escapes.
- TTY progress is disabled by default for non-TTY and CI execution.
- Hostile workspace labels cannot inject control characters or terminal escapes.
- Narrow terminals and monochrome output remain readable.
- macOS, Linux, and Windows smoke scenarios pass.
- Repository contents and scanner invocation semantics remain unchanged.
- Manual evidence includes terminal, redirected, JSON, quiet, verbose, empty, partial, failed, hostile-label, and cancellation scenarios.

## Delivery gate

Do not implement this refactor during v0.3 remediation planning. Create a separate implementation issue and release-scoped PR after the v0.3 release gate, with explicit compatibility review for CLI behavior, existing flags, output contracts, security, and dependency selection.
