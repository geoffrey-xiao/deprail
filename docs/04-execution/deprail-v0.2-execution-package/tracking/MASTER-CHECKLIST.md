# v0.2 Master Checklist

This checklist is closed only by linked evidence. It is a local planning contract; GitHub tracking is intentionally not created in this documentation pass.

## Stage kickoff

- [ ] v0.1 release decision and remaining risks are recorded.
- [ ] Latest stable, preview, and RC tags are checked.
- [ ] `origin/main` and release manifest baseline are checked.
- [ ] V02-000 accepted: v0.2.0 context, target version, baseline, and Sprint 0 scope are reviewed.
- [ ] v0.2 target version decision is recorded.
- [ ] Owners and reviewers are assigned before implementation.

## Sprint 0 — Contract and release baseline

- [ ] V02-001 accepted: empty collections serialize as arrays.
- [ ] V02-002 accepted: unexpected command arguments fail correctly.
- [ ] V02-004 accepted: release-mode evidence checklist exists.
- [ ] V02-005 accepted: real OSV-Scanner v2 fixtures exist.
- [ ] V02-006 accepted: scanner exit-code matrix is tested.
- [ ] V02-007 accepted: requested-root scanner execution is verified.
- [ ] V02-009 accepted: schema/examples reflect array collections.
- [ ] V02-013 accepted: release evidence process is defined.
- [ ] V02-014 accepted: candidate capability boundary is decided.

## Sprint 1 — Product reliability

- [ ] V02-003 accepted: build version identity is truthful.
- [ ] V02-008 accepted: outside-root scan regression is permanent.
- [ ] V02-010 accepted: deterministic ordering/artifacts are verified.
- [ ] V02-011 accepted: release identity matches tag, CLI, artifacts, and commit.
- [ ] V02-012 accepted: preview/RC/stable smoke procedure is executable.

## Sprint 5 — Baselines and diff

- [ ] V02-015 accepted: baseline representation and storage are versioned and trusted.
- [ ] V02-016 accepted: base/head comparison is deterministic.
- [ ] V02-017 accepted: new/resolved/unchanged classification is explicit.
- [ ] V02-018 accepted: `deprail diff` output and exit behavior are stable.

## Sprint 6 — Policy and SARIF

- [ ] V02-019 accepted: typed policy evaluation is deterministic.
- [ ] V02-020 accepted: expiring exceptions are bounded and auditable.
- [ ] V02-021 accepted: policy exit behavior is stable.
- [ ] V02-022 accepted: SARIF output validates and preserves provenance.

## Sprint 7 — GitHub Action and public preview

- [ ] V02-023 accepted: GitHub Action packaging is reproducible.
- [ ] V02-024 accepted: permissions and caching are least-privilege and documented.
- [ ] V02-025 accepted: real pull-request validation works.
- [ ] V02-026 accepted: public-preview release evidence is complete.

## Product acceptance

- [ ] Existing v0.1 commands remain compatible.
- [ ] Empty complete reports are schema-valid and use arrays.
- [ ] Invalid arguments never silently execute another command.
- [ ] Scanner failures and incomplete results never become safe results.
- [ ] Requested-root execution cannot escape or be replaced by caller CWD.
- [ ] Equivalent inputs produce equivalent semantic JSON.

## Platform and release

- [ ] Linux, macOS, and Windows PR verification passes.
- [ ] Required amd64 artifacts are built and smoke-tested.
- [ ] arm64 smoke evidence is recorded.
- [ ] Checksums are generated and reviewed.
- [ ] Manual representative-repository scans are recorded.
- [ ] SBOM status is explicit.
- [ ] Signing/provenance status is explicit.
- [ ] Release tag, artifact version, CLI version, and commit identity agree.

## Closure

- [ ] Every implementation item has linked PR, CI, and acceptance evidence.
- [ ] Owner acceptance is recorded during PR review.
- [ ] Remaining risks and deferred candidates are recorded.
- [ ] Final release decision is explicit.
