# DepRail v0.3.1 Development Plan

| Attribute | Value |
| --- | --- |
| Release | v0.3.1 |
| Plan revision | 1.0 |
| Status | Draft for owner and architecture/security review |
| Target mode | Preview first; stable follow-up only after evidence |
| Whole-project roadmap | `docs/03-planning/deprail-roadmap-v1.md` |
| Product baseline | `docs/01-product/deprail-product-design-v1-ai.md` |
| Architecture baseline | `docs/02-architecture/deprail-architecture-and-tech-stack-v1.md` |
| Predecessor release | `v0.3.0-preview.1` |
| Prior release plan | `docs/04-execution/deprail-v0.3-execution-package/DEVELOPMENT-PLAN-v0.3.md` |
| UX planning contract | `docs/03-planning/deprail-cli-ux-refactor-plan.md` |
| Context issue | [#238](https://github.com/geoffrey-xiao/deprail/issues/238) |
| Target milestone | `v0.3.1` |
| Owner | `@geoffrey-xiao` |
| Reviewer | `@geoffreyxiaoai` / architecture-security reviewer to confirm |

## 1. Purpose and reconciliation decision

The roadmap defines v0.3.1 as the appropriate bounded follow-up for CLI presentation refinement after the v0.3 remediation-planning gate. The product design makes the CLI the primary developer interface, while the architecture requires application services and presenters to remain separate and machine-readable contracts to remain stable.

This release improves human-facing terminal clarity without changing scan semantics, scanner behavior, remediation planning meaning, repository safety boundaries, or cross-platform domain behavior.

The existing UX plan is the normative scope source for this release:

- `docs/03-planning/deprail-cli-ux-refactor-plan.md`
- Issue [#230](https://github.com/geoffrey-xiao/deprail/issues/230)
- Merged planning PR [#231](https://github.com/geoffrey-xiao/deprail/pull/231)

This document converts that future-only plan into a release-scoped planning baseline. It does not authorize implementation until the v0.3.1 execution package and Definition of Ready are reviewed.

The release boundary remains:

```text
v0.3.0: finding -> candidate analysis -> reviewable plan
v0.3.1: clear, safe, cross-platform terminal presentation
v0.4: approved plan -> isolated application -> verification -> rescan -> patch evidence
```

### Release sequencing and displaced work

No v0.4 capability is displaced or removed. v0.3.1 is a bounded preview refinement scheduled between the v0.3 planning gate and v0.4 mutation/verification work; it consumes only presentation-focused planning and implementation capacity. The v0.4 start gate remains unchanged and still requires its own execution package, isolation, rollback, verification, rescan, and patch-evidence contracts.

The repository has no approved stable `v0.3.0` tag at this planning point. The `v0.3.0-preview.1` release evidence and retrospective are therefore implementation carryover context, not a stable compatibility baseline. v0.3.1 compatibility must be checked against the latest stable repository contracts and explicitly recorded release evidence; if no stable v0.x tag exists, the release record must state that fact rather than treating the preview as stable.

## 2. Release outcome

A developer using `deprail` receives clear, consistent, actionable terminal output during discovery, scanning, reporting, and plan inspection while automation receives unchanged, deterministic machine output.

The release must make the following true:

- Human output has a consistent visual hierarchy and status vocabulary.
- Long-running scans show progress only when the presentation channel supports it.
- Complete zero-finding scans are clearly distinguished from partial and failed scans.
- JSON stdout remains pure and deterministic.
- Diagnostics and progress remain on stderr.
- Existing `--quiet` and `--verbose` contracts remain compatible.
- Repository-derived labels cannot inject terminal control sequences.
- Non-TTY and CI output remains stable, line-oriented, and free of animation.
- Linux, macOS, and Windows preserve equivalent serialized meaning and usable terminal behavior.

## 3. User problems to solve

| User problem | v0.3.1 answer |
| --- | --- |
| A scan appears idle while work is running | Show concise progress from real application lifecycle events on interactive TTY stderr. |
| Output differs between commands | Use shared presentation primitives for headers, sections, statuses, tables, summaries, and errors. |
| Zero findings can be confused with an incomplete scan | Render explicit complete, partial, and failed outcomes. |
| CI logs contain animation or mixed machine output | Disable interactive progress for non-TTY and keep JSON stdout pure. |
| Existing quiet/verbose behavior is unclear | Preserve, document, and contract-test both flags. |
| Repository labels can contain control characters | Encode untrusted labels before terminal rendering. |
| Color-only output is hard to read | Provide monochrome and narrow-terminal fallbacks. |
| A full-screen TUI would add unnecessary complexity | Start with a static renderer and defer Bubble Tea to a real interactive workflow. |

## 4. Scope

### 4.1 Presentation contract

Define a DepRail-owned presentation contract between application services and renderers. Application services emit structured events and results rather than terminal strings.

Initial event vocabulary:

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

Each event must have deterministic fields, explicit lifecycle semantics, and no presentation-specific ANSI or terminal escape data.

The presentation contract owns:

- output mode;
- human versus JSON behavior;
- TTY capabilities;
- color and animation policy;
- progress lifecycle;
- human summaries;
- error and incomplete-result presentation;
- safe rendering of untrusted labels;
- cancellation and terminal cleanup behavior.

### 4.2 Whole-CLI visual system

Define reusable presentation primitives for:

- command headers;
- sections and subsections;
- success, warning, error, and informational states;
- tables and key/value summaries;
- finding counts and severity summaries;
- workspace status;
- plan candidates and risks;
- guidance and next steps;
- stable error-code display;
- color and monochrome output;
- narrow terminal widths.

Visual meaning must not depend on color alone. Human output must remain useful when color is disabled or unavailable.

- `deprail doctor`;
- `deprail discover`;
- `deprail scan`;
- `deprail diff`;
- `deprail fix plan`.

The CLI also provides deterministic root and subcommand help through `deprail --help` and `--help`/`-h` on every supported command. Help is plain, successful, and side-effect free.

Human message lifecycle contract:

```text
start notice -> real-work progress -> trustworthy finding summary
             -> one truthful outcome -> bounded next action
```

The contract applies consistently to discover, scan, doctor, diff, fix plan, and policy check. Messages must distinguish complete, partial, failed, and cancelled work without implying remediation or success that did not occur.

### 4.3 Interactive progress

Provide concise progress messages on stderr only for interactive TTY runs. Progress must correspond to real boundaries:

1. Validate input and resolve the repository root.
2. Discover workspaces and authoritative files.
3. Build scanner targets.
4. Execute scanner targets.
5. Retain raw artifacts.
6. Normalize findings.
7. Render the final report.

Example shape:

```text
Discovering workspaces...
Scanning 3 workspaces...
Scanning frontend... done
Scanning services/api... done
Normalizing findings...
```

A spinner or loader is optional and must remain a presentation detail. The implementation MUST NOT invent percentages, durations, completion states, or work that the underlying application has not reported.

### 4.4 Outcome summaries

Human output must distinguish:

```text
No known vulnerabilities found.
```

only for a `complete` scan with zero findings.

For incomplete and failed scans, use explicit messages such as:

```text
Scan incomplete: 1 workspace failed.
Scan failed: OSV-Scanner was not found.
```

A partial or failed scan must never look like a successful empty result.

### 4.5 Existing flag compatibility

The current implementation does not yet accept `--quiet` for `deprail scan`; this is a known CLI contract gap, not behavior to preserve. v0.3.1 must either implement the documented quiet behavior or explicitly defer and remove it from the release acceptance criteria before implementation approval. The preferred scope is to implement it and add parser, stream, and smoke coverage.

Required semantics:

- `--quiet` suppresses progress and successful human summaries while retaining errors and required diagnostics.
- `--verbose` preserves existing behavior and exposes additional safe diagnostics.
- Both remain compatible with `--format json`.
- Neither flag contaminates machine stdout.
- Any semantic change requires a CLI contract update and compatibility evidence.

### 4.6 Machine-output and TTY contract

The release MUST preserve these invariants:

- `--format json` writes JSON data only to stdout.
- Diagnostics and progress use stderr.
- Redirected stdout remains valid JSON with no banners or ANSI escapes.
- Interactive progress is enabled only when stderr is a TTY.
- Non-TTY progress is disabled by default.
- Non-TTY diagnostics are stable, line-oriented, and free of carriage-return animation.
- JSON output remains deterministic and unchanged except for explicitly reviewed additive changes.
- Human presentation never leaks secrets, credential-bearing URLs, or sensitive environment values.

A future structured-progress option may be considered separately; it is not part of this release unless the contract is explicitly added and reviewed.

### 4.7 Hostile-label rendering

Repository-derived values are hostile inputs. The renderer must protect:

- workspace IDs;
- repository-relative paths;
- package names;
- error details;
- scanner labels;
- plan and candidate descriptions.

The implementation must:

- escape newlines, carriage returns, and tabs where unsafe;
- prevent raw ANSI and terminal control sequences;
- preserve enough readable context to identify the workspace or component;
- apply the same protection to TTY and redirected stderr output;
- avoid logging secrets or credential-bearing URLs.

### 4.8 Dependency strategy

Start with a small DepRail-owned internal renderer and explicit capability detection.

Selective styling may use Lip Gloss. Selected Bubbles components may be used for a spinner or progress element after the static renderer and event contract are stable.

Bubble Tea is deferred unless a real interactive workflow such as remediation review, approval, or plan inspection is approved. PTerm is not selected by this plan; adopting it would require a separate review of global printer behavior, TTY detection, stderr routing, JSON purity, cancellation, and concurrent output.

## 5. Architecture impact

The implementation must preserve the inward dependency rule:

```text
CLI entry point
    -> application services
        -> structured application events and domain results
            -> presentation contract
                -> terminal renderer or machine presenter
```

### 5.1 Expected ownership

| Area | v0.3.1 responsibility |
| --- | --- |
| `cmd/deprail` | Translate CLI flags and streams into application presentation options. |
| `internal/app` | Emit structured lifecycle events at real work boundaries. |
| `internal/presenter` | Own human summaries, JSON purity, output modes, statuses, and safe rendering. |
| `internal/domain` | Remain independent of terminal libraries and ANSI output. |
| `internal/process` | Preserve cancellation, deadlines, bounded output, and process safety. |
| Terminal capability layer | Detect TTY, color, width, CI/non-TTY mode, and animation support. |

The exact file and package names are implementation decisions for the execution package. Domain and scanner packages must not import terminal styling or TUI dependencies.

### 5.2 Dependency direction

Application services may emit plain event structs through an interface or channel. Renderers consume those events and write to the approved stream. JSON presenters remain normative for machine output.

No renderer may alter domain results, scanner invocation arguments, exit status meaning, or artifact persistence behavior.

### 5.3 Trust boundaries

Terminal output is a presentation boundary, not a trust boundary. Repository values, scanner output, and error details remain untrusted. Rendering must encode them safely before output.

Terminal styling dependencies must not gain access to credentials, unrestricted environment values, repository mutation capabilities, or network permissions.

## 6. Explicit exclusions

The following are not v0.3.1 deliverables:

- scanner behavior changes;
- new scanners or package ecosystems;
- repository mutation or remediation execution;
- package-manager installation, update, or script execution;
- network services, telemetry, or source upload;
- changes to JSON schemas without compatibility review;
- hidden progress output in machine-readable streams;
- artificial delays or activity indicators disconnected from real work;
- replacing all command parsing or application architecture without a measured need;
- full-screen TUI for ordinary scan, discover, doctor, or plan commands;
- autonomous approval, merge, patch creation, or publication;
- v0.4 isolation, rollback, verification execution, or rescan behavior.

## 7. Delivery stages

### Planning Sprint — contract reconciliation

- Confirm v0.3 preview evidence, owner decision, and remaining review gaps.
- Reconcile product, architecture, roadmap, v0.3.1 plan, UX plan, and CLI contracts.
- Freeze presentation events, output modes, TTY capabilities, status vocabulary, and compatibility rules.
- Define the acceptance matrix and manual evidence scenarios.
- Create the v0.3.1 execution package and implementation issue map.

### Sprint 1 — presentation contract and static renderer

- Define event structs and lifecycle semantics.
- Add terminal capability detection.
- Implement deterministic line-oriented human summaries and errors.
- Preserve existing JSON and stderr/stdout behavior.
- Add empty, partial, failed, and cancellation presentation.

### Sprint 2 — styling and progress

- Add selective styling with monochrome and narrow-terminal fallbacks.
- Wire progress to real application events.
- Add TTY-only loader/progress behavior.
- Confirm non-TTY output has no animation or carriage-return updates.
- Add hostile-label rendering and adversarial tests.

### Sprint 3 — cross-command integration and hardening

- Integrate presentation behavior across `doctor`, `discover`, `scan`, `diff`, and `fix plan`.
- Implement and test `--quiet`; preserve and test `--verbose`.
- Test redirected output, CI behavior, cancellation, terminal cleanup, and concurrent output.
- Review dependency footprint and renderer ownership.
- Produce manual and cross-platform evidence.

### Release hardening

- Run the general release checklist.
- Verify JSON stdout purity and deterministic output.
- Verify TTY and non-TTY behavior on Linux, macOS, and Windows.
- Verify empty, partial, failed, hostile-label, narrow-terminal, monochrome, quiet, verbose, and cancellation scenarios.
- Review compatibility and security impact.
- Publish a v0.3.1 preview only after owner and security/architecture review.

## 8. Test and evidence strategy

### Unit tests

- Event lifecycle ordering and terminal-state invariants.
- TTY capability detection and output-mode selection.
- Status and summary selection for complete, partial, failed, and cancelled operations.
- Quiet and verbose precedence.
- Terminal width and monochrome layout decisions.
- Stable escaping of hostile labels and control sequences.

### Presenter contract tests

- JSON stdout contains only valid JSON.
- Progress and diagnostics remain on stderr.
- Human rendering does not alter domain or machine results.
- Repeated equivalent inputs produce deterministic human summaries where the contract requires it.
- Non-TTY output is line-oriented and contains no animation control sequences.
- Error codes and incomplete states remain visible and truthful.

### Golden and snapshot tests

Cover:

- complete scan with findings;
- complete scan with zero findings;
- partial scan with workspace failures;
- failed scan with missing scanner;
- `doctor`, `discover`, `scan`, and `fix plan` summaries;
- quiet output;
- verbose output;
- monochrome output;
- narrow terminal output;
- hostile workspace IDs, paths, package names, and scanner labels;
- cancellation and interrupted progress.

Snapshot tests must assert semantic presentation contracts, not incidental ANSI sequences or library internals.

### Security tests

- ANSI and terminal escape injection through repository-derived labels.
- Newline, carriage-return, tab, Unicode, and long-label handling.
- Credential-bearing URL and sensitive-environment redaction.
- JSON contamination through progress or diagnostics.
- Concurrent output interleaving.
- Cancellation and terminal cleanup.
- Shell metacharacters remain data and are never executed.

### Cross-platform evidence

Linux, macOS, and Windows must preserve equivalent serialized output and status meaning. Platform-specific terminal capabilities may change styling or line endings only where explicitly documented.

Manual evidence must include:

- interactive TTY output;
- redirected stderr/stdout;
- JSON output;
- quiet and verbose modes;
- complete zero findings;
- partial scan;
- failed scan;
- hostile labels;
- narrow terminal;
- monochrome output;
- cancellation;
- scanner failure and timeout behavior.

## 9. Dependencies and risks

### Dependencies

- The v0.3.0-preview.1 release evidence and retrospective remain carryover context; no approved stable v0.3.0 tag exists at this planning point.
- Issue #230 and PR #231 provide the approved UX scope.
- Existing CLI, JSON, stdout/stderr, and flag contracts remain normative.
- Existing application lifecycle boundaries must be sufficient to emit real progress events.
- Any new terminal dependency requires dependency, license, compatibility, and security review.

### Risks

| Risk | Response | Release gate |
| --- | --- | --- |
| Progress contaminates JSON or stdout | Route all progress through presentation and stderr policy | JSON purity tests and redirected-output evidence |
| False progress claims completion | Emit only from real application lifecycle events | Event contract and integration tests |
| `--quiet` or `--verbose` behavior changes silently | Define quiet implementation and preserve verbose semantics with compatibility evidence | CLI contract review and smoke evidence |
| `--quiet` is currently rejected by the scan parser | Implement quiet mode explicitly rather than treating rejection as compatibility | Parser, stream, and smoke coverage |
| Whole-CLI coverage omits `diff` | Include `diff` in shared presenter integration and evidence | Cross-command smoke coverage |
| Hostile labels inject terminal control sequences | Escape and test all repository-derived labels | Adversarial presenter tests |
| Styling makes output unreadable in CI or narrow terminals | Disable animation and provide monochrome/width fallbacks | Cross-platform and narrow-terminal evidence |
| Terminal library leaks global output or cancellation behavior | Keep renderer behind DepRail-owned interfaces and review dependencies | Architecture/security review |
| CLI presentation changes scanner or domain semantics | Keep events at application boundaries and presenters read-only | Regression tests and unchanged scanner invocation evidence |
| Scope expands into interactive remediation | Defer Bubble Tea and mutation workflows to a separate contract | Definition of Done and change-control review |
| Concurrent progress becomes nondeterministic | Serialize renderer writes and test event ordering | Concurrent-output tests |

## 10. Definition of Ready

Before creating implementation issues or editing runtime code:

- the v0.3.0 preview release and carryover risks are recorded;
- product, architecture, roadmap, UX plan, and this plan agree;
- the v0.3.1 execution package exists;
- presentation events and renderer ownership are explicit;
- CLI, JSON, stdout/stderr, TTY, non-TTY, quiet, verbose, and exit contracts are reconciled;
- hostile-label and terminal-control requirements are named;
- fixtures and manual evidence scenarios are named;
- dependency and license review scope is defined;
- owner and reviewer are assigned;
- tracking context is prepared; no implementation issues are created by this planning step.

## 11. Definition of Done

v0.3.1 is complete when:

- the whole CLI uses consistent human presentation primitives;
- interactive progress reflects real application phases without false claims;
- complete, partial, failed, and cancelled outcomes are distinct;
- `--quiet` is implemented with the documented suppression semantics and covered by parser, stream, and smoke evidence;
- JSON stdout remains parseable, deterministic, and free of progress or ANSI escapes;
- non-TTY and CI output is stable and line-oriented;
- hostile labels cannot inject terminal control sequences;
- narrow and monochrome output remains readable;
- Linux, macOS, and Windows evidence is complete;
- repository contents and scanner invocation semantics remain unchanged;
- release artifacts and evidence pass the general release checklist;
- owner and security/architecture reviewer approve the release evidence;
- the v0.3.1 retrospective records follow-up actions and remaining risks.

## 12. Decision gate

The v0.3.1 preview may proceed only when the presentation contract, compatibility behavior, and security boundaries are reviewed and representative terminal evidence is complete.

A stable v0.3.1 release requires fresh stable-tag evidence, cross-platform smoke results, explicit owner approval, and security/architecture review. The v0.3.1 work must not weaken the v0.3 planning boundary or begin v0.4 repository mutation capabilities.

The v0.4 handoff remains separate. Any interactive remediation review, approval, mutation, verification, rollback, or rescan workflow requires its own v0.4 contract and release gate.
