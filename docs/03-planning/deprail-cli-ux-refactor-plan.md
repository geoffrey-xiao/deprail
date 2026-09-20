# DepRail CLI UX Refactor Plan

**Status:** Future planning only; not authorized for current v0.3 implementation.

**Target stage:** Post-v0.3 CLI refinement, scheduled after the remediation-planning release gate.

**Purpose:** Make interactive CLI execution easier to understand without weakening machine-readable output, deterministic behavior, cross-platform semantics, or security boundaries.

## Motivation

The current CLI correctly separates machine output from diagnostics and reports explicit complete, partial, and failed states. Long-running scans can still appear idle while discovery, OSV-Scanner execution, artifact retention, and normalization are in progress. A future refactor should provide useful progress feedback for human terminal users while preserving the existing automation contract.

## Scope

### Interactive progress

Provide concise progress messages on stderr for interactive terminal runs:

```text
Discovering workspaces...
Scanning 3 workspaces...
Scanning frontend... done
Scanning services/api... done
Scanning services/worker... done
Normalizing findings...
```

An optional spinner or loader may be used while a phase is active, but it must be a presentation detail rather than a domain or scanner behavior.

Progress phases should be based on real application events:

1. Validate input and resolve the repository root.
2. Discover workspaces and authoritative files.
3. Plan scanner targets.
4. Execute scanner targets.
5. Retain raw artifacts.
6. Normalize findings.
7. Render the final report.

Do not invent progress percentages or claim completion before the underlying phase completes.

### Empty-result messaging

After a `complete` scan with zero findings, show an explicit human message:

```text
No known vulnerabilities found.
```

Never use that message for `partial` or `failed` scans. Those states must remain visibly incomplete or failed:

```text
Scan incomplete: 1 workspace failed.
```

```text
Scan failed: OSV-Scanner was not found.
```

### Machine-output contract

The refactor must preserve these rules:

- `--format json` writes JSON data only to stdout.
- Diagnostics and progress use stderr.
- Redirected stdout remains valid JSON with no banners or ANSI escapes.
- Progress is disabled when stderr is not a TTY unless an explicit future flag enables structured progress.
- JSON output remains deterministic and unchanged except for approved schema changes.
- No progress message contains secrets, credential-bearing URLs, or raw sensitive environment values.

### Quiet and non-interactive behavior

Define a consistent future interface for:

- `--quiet`: suppress human progress while retaining errors.
- `--verbose`: expose additional diagnostics without exposing secrets.
- CI/non-TTY execution: no spinner, no carriage-return animation, stable line-oriented diagnostics.
- Cancellation: stop progress cleanly and report the stable cancellation/failure code.

Do not add these flags until their compatibility and exit behavior are specified.

### Terminal presentation

Refactor terminal rendering around a small presentation layer that owns:

- TTY detection.
- Color and animation policy.
- Progress lifecycle.
- Human summary formatting.
- Error and incomplete-result presentation.

The application and domain layers must emit structured phase events or results; they must not write terminal strings directly.

## Non-goals

This plan does not authorize:

- Scanner behavior changes.
- New scanners or package ecosystems.
- Network services or telemetry.
- Repository mutation or remediation execution.
- Changes to JSON schemas without compatibility review.
- Hidden progress output in machine-readable streams.
- Artificial delays or activity indicators disconnected from real work.

## Proposed implementation slices

1. **CLI presentation contract**
   - Define progress events, severity, lifecycle, TTY policy, and quiet behavior.
2. **Progress renderer**
   - Implement line-oriented non-TTY output first, then optional TTY animation.
3. **Application event wiring**
   - Emit discovery, scanner, artifact, normalization, and render phases from real boundaries.
4. **Empty and incomplete summaries**
   - Distinguish complete zero findings from partial and failed scans.
5. **Cross-platform terminal behavior**
   - Verify macOS, Linux, and Windows output, cancellation, and terminal cleanup.
6. **Regression and release evidence**
   - Prove JSON stdout purity, stderr diagnostics, no-animation redirection, deterministic summaries, and no secret leakage.

## Acceptance criteria for a future issue

- Interactive scan runs provide useful phase/workspace feedback without false progress.
- Complete zero-finding scans explicitly say no known vulnerabilities were found.
- Partial and failed scans never present a safe-looking empty result.
- JSON stdout remains parseable and free of progress text or ANSI escapes.
- Non-TTY output is stable, line-oriented, and animation-free.
- `--quiet` and `--verbose` behavior is documented and tested if introduced.
- macOS, Linux, and Windows smoke scenarios pass.
- Repository contents and scanner invocation semantics remain unchanged.
- Manual evidence includes terminal, redirected, JSON, empty, partial, failed, and cancellation scenarios.

## Delivery gate

Do not implement this refactor during v0.3 remediation planning. Create a separate issue and release-scoped PR after the v0.3 release gate, with explicit compatibility review for CLI behavior and output contracts.
