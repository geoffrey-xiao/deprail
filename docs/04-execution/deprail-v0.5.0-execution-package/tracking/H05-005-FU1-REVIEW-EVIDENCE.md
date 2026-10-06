# H05-005-FU1 Transport Review Evidence

- Issue: [#441](https://github.com/geoffrey-xiao/deprail/issues/441); source review: [PR #440](https://github.com/geoffrey-xiao/deprail/pull/440).
- Status: local implementation evidence only. Owner acceptance, follow-up PR review/merge, and release approval remain pending.

## Contract corrections

- The local `http.Server` now sets `ReadTimeout: 10s` alongside the existing header, write and idle timeouts. A real TCP client withholding a declared request byte after a rejected GET is disconnected at a shortened test read deadline; a subsequent authenticated health request still succeeds. The test also checks the production deadline is the contract value before shortening it for the scenario.
- A history cursor now checks both its UUID and UTC-microsecond tuple against a validated stored occurrence. An invented UUID and an altered timestamp return `400 API_REQUEST_INVALID` rather than silently skipping entries. Store corruption/unavailability remains typed.
- A findings cursor must match a stable key in the selected validated parent report before advancing; an invented key returns `400 API_REQUEST_INVALID`. Existing real-listener pagination verifies legitimate continuations without skips and rejects wrong parent/resource use.
- Unauthorized API requests include `WWW-Authenticate: Bearer` with the typed JSON 401. Both absent and invalid credentials were tested through the listener.

## Verification

- Focused live-listener test command: `go test -count=1 -timeout=90s -run 'TestHTTPUnreadRejectedBodyConnectionIsBounded|TestHTTPRejectsInventedHistoryAndFindingCursorTuples|TestHTTPRejectsCrossOriginMalformedAndUnsupportedRequests|TestHTTPHistoryRoutesPageAndPreserveStoredMeaning' ./internal/transport/localhttp` — passed.
- Focused application query command: `go test -count=1 -timeout=90s -run 'TestWorkspaceHashedOrdinalCursorContinuesAndRejectsTampering|TestCaptureAndQuery' ./internal/app` — passed.
- `make verify` — passed (`go generate`, `go vet`, unit/contract tests, and `go build`).
- `make test-integration` — passed on the Darwin workstation.
- `GOOS=windows GOARCH=amd64 go test -run '^$' -exec=/usr/bin/true ./cmd/deprail ./internal/transport/localhttp ./internal/app` and the same command with `GOOS=linux` — passed as compile-only checks, not runtime tests.
- Exact-head CI and owner security/contract review remain pending. No production UI or asset bundle is claimed.

## Remaining risk

The 10-second read deadline bounds connection lifetime, not the number of unauthenticated connections simultaneously open during that interval. The eight-slot admission bound remains an API work bound. Same-user local processes remain within the documented local trust model; production assets/UI still belong to H05-007/H05-006.
