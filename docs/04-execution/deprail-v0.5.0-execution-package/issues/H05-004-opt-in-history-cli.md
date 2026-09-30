# H05-004 Opt-in History CLI
- GitHub Issue: [#423](https://github.com/geoffrey-xiao/deprail/issues/423).
- GitHub parent: [#418](https://github.com/geoffrey-xiao/deprail/issues/418). Published dependencies: #422; global pre-start gate #419.

## Planning metadata

- Type: `feature`
- Area: `cli`
- Priority: `P0`
- Risk: `R3`
- Target version: `v0.5.0`
- Milestone: [`v0.5.0`](https://github.com/geoffrey-xiao/deprail/milestone/11)
- Sprint: Unassigned
- Owner: `@geoffrey-xiao`
- Reviewer: `@geoffrey-xiao`
- Dependencies: [H05-003](H05-003-history-application-services.md); accepted planning/design contracts and parent-owned timeout follow-up disposition.
- Blocked reason: None unless Project Status is Blocked; runtime remains gated on acceptance below.

## Definition of Ready

- [x] Value and user impact are stated.
- [x] Scope and explicit exclusions are stated.
- [x] Inputs, outputs, and failure behavior are defined.
- [x] Required tests or smoke scenarios are named.
- [x] Acceptance criteria are observable.
- [x] Owner reviewer is assigned; external reviewer is optional.
- [x] Dependencies and target version are recorded.

## Goal

Expose history capture only through explicit `deprail scan --save-history`, preserving all existing default scan output and exit behavior.

## Scope

Own `cmd/deprail` scan flag parsing/wiring and CLI usage documentation. Add boolean `--save-history`; absent flag must not initialize, create or write history. With flag, connect one selected scan operation to H05-003 capture boundary, including terminal completed/failed/cancelled outcomes where typed safe context exists. Preserve exact existing scan stdout bytes and JSON schema in all paths; history diagnostics only go to stderr using static typed safe message/code, never raw path, token, scanner error or secret. Existing nonzero scan/cancellation exit code remains primary if save also fails. If scan exits 0 and requested persistence fails, return additive exit code `6`; successful requested save returns existing success code. Do not repurpose exits 1–5. Never print “saved” until commit success. Update CLI help to describe opt-in, private data location, quota refusal and residual local retention/growth boundaries.

Trace: [FR-505](../requirements/FUNCTIONAL-REQUIREMENTS.md), [STORE-01](../requirements/FAILURE-AND-DATA-CONTRACT.md), [CLI compatibility](../requirements/COMPATIBILITY-MATRIX.md), [error model](../requirements/ERROR-MODEL.md), [ADR-0004 selections](../../../adr/ADR-0004-local-scan-history.md#detailed-engineering-selections-for-owner-review), [product baseline](../../../01-product/deprail-product-design-v1-ai.md), [architecture baseline](../../../02-architecture/deprail-architecture-and-tech-stack-v1.md), [release plan §14](../../../03-planning/deprail-development-plan-v0.5.0.md).

## Out of Scope

No default/automatic capture, CLI history listing or deletion, browser/HTTP launch, scan output or signature changes, changed existing exit codes, raw error forwarding, or retrying persistence.

## Inputs, Outputs, and Failure Behavior

Input is existing scan args plus boolean opt-in. Output is unchanged scan stdout/report, existing scanner stderr semantics plus safe history diagnostic only on requested failure, and exit precedence: scan nonzero wins; otherwise save failure gives 6; otherwise existing success. No flag makes no history-store call or filesystem side effect. A failed save never claims successful persistence. Invalid flag use follows existing CLI argument-error convention without starting scan.

## Required Tests

- CLI compatibility: no flag proves no store open/create/write, byte-for-byte stdout equality and existing exit values.
- Opt-in successful scan writes exactly one occurrence; stdout bytes remain identical; successful message only follows commit.
- Successful scan + busy timeout, permission refusal, quota, disk full, invalid projection returns 6, unchanged stdout and safe stderr with no secret/path leakage.
- Scanner nonzero + save failure preserves scanner exit as primary and records distinct safe history diagnostic; cancellation similarly retains established primary code.
- Repeated invocation same source ID creates distinct occurrence; no automatic capture from other commands.
- Smoke actual CLI against isolated HOME/config root: default scan leaves it absent; opt-in scan creates private store, reopening shows one row; injected failure observes exit 6 and unchanged stdout.

## Acceptance Criteria

- [ ] `--save-history` is the sole explicit trigger; without it there is no history initialization or write.
- [ ] Existing stdout/report JSON is byte-for-byte unchanged with and without requested history.
- [ ] Exit code 6 occurs only for otherwise-successful scan with requested save failure; existing nonzero scan/cancellation remains primary.
- [ ] History failure stderr is typed, static and privacy-safe; no false persistence success message.
- [ ] Help explains opt-in and local retention/quota behavior without implying hard disk quota.
- [ ] Consumer-visible precedence, failure/security, and actual CLI state smoke results are supplied after implementation.

## Owner Review

`@geoffrey-xiao` reviews compatibility and each criterion; records technical acceptance separately from security/architecture acceptance. External review optional. CLI runtime behavior and CI are future evidence, not planning claims.

## Evidence Required

- Verification commands or scenarios: CLI argument and exit-precedence tests; isolated HOME smoke for default absence, opt-in persistence, injected failure; stdout hash/byte comparison.
- Expected artifacts, logs, screenshots, or links: captured exit/stdout/stderr, database row count and permission evidence, privacy scan, CI results, owner decisions.

## Final Acceptance

- [ ] Owner reviewed every acceptance criterion during PR review.
- [ ] Required verification and CI results were reviewed.
- [ ] Owner review and remaining risk are recorded.
- [ ] Evidence links are attached.
- [ ] Owner review and merge evidence are linked.

## Rollback

Withdraw the opt-in flag wiring without altering ordinary scan behavior; preserve already committed database/artifacts. Owner `@geoffrey-xiao` records rollback and any recovery action.
