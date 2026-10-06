# H05-005-FU1: Close Local-Console Transport Review Findings

- GitHub issue: [#441](https://github.com/geoffrey-xiao/deprail/issues/441); parent H05-005 [#424](https://github.com/geoffrey-xiao/deprail/issues/424).
- Owner and owner reviewer: `@geoffrey-xiao`; target release `v0.5.0`, Sprint 4; `area:cli`, `risk:R3`, `priority:P0`, `type:bug`.
- Source: three contract defects identified in [PR #440 review](https://github.com/geoffrey-xiao/deprail/pull/440): rejected unread bodies can hold connections; invented history/finding cursor tuples can skip records; API 401 lacks the declared bearer challenge.
- Existing contracts: [API design §§11–13](../API-DESIGN.md), [OpenAPI v1](../../../../schemas/openapi/v1/openapi.yaml), [SEC-01/05](../requirements/SECURITY-REQUIREMENTS.md), and [H05-005](H05-005-read-only-local-console.md). This fixes accepted behavior, with no new release capability or displaced work.

## Goal and scope

Bound unread rejected-body connections with the existing 10-second request limit. Confirm history cursor UUID/timestamp against a stored occurrence and finding cursor key against the selected immutable report projection; invented tuples return `400 API_REQUEST_INVALID`, while corrupt/unavailable storage remains a typed error. Emit `WWW-Authenticate: Bearer` for unauthorized API responses. Preserve valid whole-record continuation and exact parent/resource binding.

Production UI/assets (H05-006/H05-007), new routes, scans/writes, schema and storage migration, retries and token persistence remain out of scope. The previously merged H05-005 branch is not reused.

## Acceptance and evidence

- [ ] An unread rejected-body connection closes within its read deadline; the listener still serves a healthy request.
- [ ] Fabricated/mismatched history and finding cursor tuples fail `400`; legitimate next-page cursors continue without skips.
- [ ] Missing or invalid API bearer returns `401` and the exact declared challenge; valid authorization remains unaffected.
- [ ] Focused live-listener regressions, `make verify`, `make test-integration`, and exact-head CI are linked.
- [ ] Owner reviews each criterion, security/contract effects and remaining risk on the follow-up PR.

Local command results and tested scenarios: [H05-005-FU1 transport review evidence](../tracking/H05-005-FU1-REVIEW-EVIDENCE.md). Acceptance remains unchecked until exact-head CI and owner review.

No owner acceptance or release approval is implied by local verification. Rollback reverts only this follow-up and leaves history and artifacts unchanged.
