# DepRail v0.3.1 Architecture

## Architectural position

v0.3.1 adds a presentation boundary to the existing Go application flow. Application services emit structured lifecycle events and domain results; replaceable presenters render terminal or machine output. Domain, discovery, scanner, artifact, and remediation semantics remain unchanged.

```text
cmd/deprail
  -> internal/app
    -> domain, discovery, scan, remediation
    -> presentation events
      -> internal/presenter
        -> terminal renderer or machine presenter
```

## Responsibilities

| Component | Responsibility |
| --- | --- |
| `cmd/deprail` | Translate flags, streams, and terminal options into application commands. |
| `internal/app` | Emit events at real work boundaries and preserve cancellation. |
| presentation contract | Define deterministic event fields and lifecycle semantics. |
| `internal/presenter` | Select output mode, render human summaries, preserve JSON purity, and encode labels. |
| terminal capability layer | Detect TTY, color, width, CI/non-TTY mode, and animation support. |
| domain and adapters | Remain independent of terminal libraries and presentation strings. |

## Event contract

The initial event vocabulary is:

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

Events carry structured data only. They do not contain ANSI sequences, spinner state, terminal control bytes, or renderer-specific layout instructions.

## Presentation flow

1. Validate input and resolve the canonical root.
2. Select output mode and terminal capabilities.
3. Execute application work unchanged.
4. Emit structured lifecycle events.
5. Render TTY progress to stderr when allowed.
6. Render the final human or machine result through the selected presenter.
7. Restore terminal state on cancellation or failure.

## Trust boundary

Repository paths, workspace IDs, package names, scanner output, and error details are hostile inputs. The renderer must escape terminal control sequences and unsafe line delimiters before output. Styling dependencies receive presentation data only and cannot access mutation, credentials, or unrestricted environment state.

## Dependency policy

Use a small internal renderer first. Lip Gloss may be selected for styling; Bubbles may be selected for bounded loader/progress elements after the event contract is stable. Bubble Tea is deferred until a genuinely interactive workflow is approved. PTerm is not selected without explicit review of global output, stderr routing, cancellation, concurrency, and JSON purity.

## Compatibility

Existing scan, discover, doctor, plan, JSON, stdout/stderr, exit-code, path, provenance, scanner, and mutation-boundary contracts remain unchanged. Any additive machine contract change requires schema review, examples, and compatibility tests. No v0.3.1 presentation dependency may be imported by domain packages.
