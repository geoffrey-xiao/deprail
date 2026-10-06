# H05-005 Transport Evidence

Issue: [#424](https://github.com/geoffrey-xiao/deprail/issues/424)  
Branch: `feat/h05-005-local-console`  
Evidence status: implementation verification only; owner review, acceptance, PR review, and merge remain pending.

## Delivered behavior

- `deprail web` is wired as a foreground command with `--open` and optional `--artifact-root`; the default command does not launch a browser. Non-TTY use requires `--open`, and launcher errors are static and do not expose credential URLs or local paths.
- `internal/transport/localhttp` binds only `127.0.0.1` on an ephemeral port, serves exact console/static allowlists and five read-only API operations, and fails before binding when assets are absent or invalid.
- API requests enforce exact Host, optional same-origin Origin and Fetch Metadata, a process-scoped 256-bit bearer, bodyless GETs, canonical routes/cursors, target/header/response bounds, an eight-request admission limit, a ten-second deadline, and bounded shutdown. Shutdown revokes the credential before draining.
- History queries use the existing read-only application/store services. Missing storage is a successful empty page without directory creation; typed storage failures remain errors. Optional artifact verification uses a canonical, readable directory and returns integrity metadata only.
- Workspace continuation cursors contain a parent UUID, ordinal, and SHA-256 workspace identity; neither the cursor nor its application continuation value carries a workspace name.
- Test-only assets are built from in-memory `fstest.MapFS` fixtures. No fallback assets, embedded production UI, or filesystem serving were added.

## Automated verification

- `make verify` — passed after the final implementation changes (`go generate`, `go vet`, all package tests, and `go build`).
- `make test-integration` — passed after the final implementation changes.
- Focused transport/CLI tests — passed, including all five API operations against a real listener and a real SQLite history store, empty/missing history behavior, typed failure mapping, read-only database digest comparison, actual artifact digest verification with metadata-only responses and symlinked-path rejection, cursor continuation/tampering/parent/resource rejection, static allowlists and headers, Host/Origin/Fetch Metadata/auth/method/body/query rejection, target/header/response limits, deadline handling, eight-slot saturation, token revocation, active-query cancellation during shutdown, non-TTY launch failure, and TTY reopen.
- A separately built Darwin `deprail` binary printed `web --help` and was launched with `web --open`; the expected fail-closed result was exit `3` with `API_LISTENER_UNAVAILABLE` because H05-007 production assets are not yet connected.
- Windows amd64 and Linux amd64 console packages — cross-compiled successfully with `GOOS=windows` and `GOOS=linux`; `-exec=/usr/bin/true` means these were compile checks only, not runtime tests on those operating systems.

## Browser smoke

On the Darwin workstation, a temporary test-only server used the production `localhttp.Server` and in-memory shell assets. Chromium observed:

1. Bootstrap fragment consumed and removed before the authenticated health fetch; API status became `ready`.
2. Reload without a fragment showed the safe “reopen from the active terminal” recovery state.
3. Reopening the same tab with the active-session bootstrap URL handled the fragment change, authenticated successfully, and removed the fragment again.
4. `location.hash` was empty after bootstrap; local and session storage remained empty.

A screenshot of the test-only shell showing `API ready` was captured by the browser tool at `/var/folders/gz/dx70pnsn0pdb5y7m5fdhtnqr0000gn/T/omp-sshots-159b6fd7dd4fa267.webp`. The image is a transient runner artifact, not a production UI screenshot. The browser harness used a deterministic test-only credential; the CLI lifecycle tests exercised the real random process credential against the actual local HTTP listener.

## Remaining scope and risk

- H05-007 still owns reviewed production asset embedding; until that package is connected, the production `run()` path has no assets and `deprail web` safely exits with `API_LISTENER_UNAVAILABLE`. The test assets are not a fallback.
- H05-006 owns the production UI and must preserve the tested fragment-removal, recovery, and same-tab fragment-change behavior. The browser smoke validates transport and browser mechanics, not that future UI bundle.
- Linux and Windows runtime/browser behavior, CI, and owner security/contract acceptance are not established by these local checks. Do not mark the H05-005 acceptance boxes complete until owner review and required PR/CI evidence are recorded.
