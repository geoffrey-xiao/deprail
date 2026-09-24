# DepRail General Release Checklist

Use this checklist for every preview, release candidate, and stable release. Copy the checklist into the applicable release evidence record or link the completed evidence rows from that record. An unchecked item is a gap, not an implied pass.

This checklist is a release gate. It does not replace the version-specific development plan, execution-package master checklist, release evidence record, or retrospective.

## Release identity

- [ ] Intended version and release mode are recorded: preview, release candidate, or stable.
- [ ] Latest stable, preview, and release-candidate tags are checked.
- [ ] Reviewed `origin/main` commit is recorded.
- [ ] Release manifest matches the intended version.
- [ ] The release tag is new and will be immutable.
- [ ] DepRail CLI identity will match the version, tag, and source commit.
- [ ] Release owner and owner-review record are identified; external reviewer is optional.

## Plan and scope completion

- [ ] The applicable development plan is approved and its release outcome is clear.
- [ ] Every release-scoped implementation issue is complete or has an explicit owner-approved disposition.
- [ ] Every implementation item has a merged PR, CI result, and acceptance evidence.
- [ ] No release-blocking P0/P1 issue remains open without an explicit disposition.
- [ ] Included capabilities are implemented and documented.
- [ ] Excluded and deferred capabilities are documented.
- [ ] Compatibility impact and migration notes are recorded.
- [ ] Remaining risks and rollback scope are explicit.

## Contract and security completion

- [ ] CLI and exit-code contracts are reviewed.
- [ ] Machine-readable schemas and examples are versioned and validated.
- [ ] Stable identity and deterministic ordering are verified.
- [ ] Failure, incomplete, stale, unsupported, and unknown states are explicit.
- [ ] Path containment and symlink boundaries are tested.
- [ ] External-process, timeout, cancellation, and output-limit behavior is tested where applicable.
- [ ] Network, package-script, credential, and file-write boundaries are documented.
- [ ] Owner security/architecture review and risk disposition are recorded for schema, compatibility, process, file-write, credential, publishing, and release changes.

## Automated and manual verification

- [ ] Unit tests pass.
- [ ] Contract, integration, golden, and schema tests pass where applicable.
- [ ] Security and hostile-input tests pass.
- [ ] Determinism and repeated-run checks pass.
- [ ] `make verify` passes from a clean checkout of the reviewed commit.
- [ ] `git diff --check` passes.
- [ ] Linux CI passes.
- [ ] macOS CI passes.
- [ ] Windows CI passes.
- [ ] Manual representative-repository smoke passes.
- [ ] Required ecosystem fixtures are covered.
- [ ] Repository trees, manifests, lockfiles, and package-manager state remain unchanged when the release is planning-only.

## Artifact and supply-chain verification

- [ ] Artifacts are built from the reviewed immutable commit and tag.
- [ ] The supported platform matrix is complete or every gap is explicitly accepted for the release mode.
- [ ] Artifact names and sizes are recorded.
- [ ] SHA-256 checksums are generated and independently verified.
- [ ] `doctor` or equivalent identity output matches the release tag and commit.
- [ ] SBOM status is recorded as supplied, unavailable, or explicitly deferred.
- [ ] Signature status is recorded as supplied, unavailable, or explicitly deferred.
- [ ] Provenance status is recorded as supplied, unavailable, or explicitly deferred.
- [ ] Release assets contain only intended public resources.
- [ ] Raw logs, scan reports, credentials, and internal evidence are kept outside public release assets.

## Publication and approval

- [ ] Release notes explain the user outcome, installation, usage, scope, limitations, and feedback path.
- [ ] The GitHub Release is marked preview, release candidate, or stable correctly.
- [ ] The protected approval environment is configured and verified.
- [ ] The release workflow completed successfully.
- [ ] Assets were uploaded only after the required approval gate.
- [ ] Release, workflow, evidence, and comparison URLs are recorded.
- [ ] Rollback owner and immutable-tag recovery procedure are recorded.
- [ ] Owner records the release decision: go, go with approved gaps, or no-go.

## Post-release closure

- [ ] The version-specific release evidence record is complete.
- [ ] The version-specific retrospective is created under `docs/retrospectives/`.
- [ ] The retrospective is linked from the release issue and release evidence record.
- [ ] Retrospective follow-up actions have owners and acceptance evidence.
- [ ] Owner acceptance and the owner's security/architecture review are recorded separately; one person may complete both.
- [ ] The GitHub Project status reflects the actual release state.
- [ ] The release issue is closed only after required evidence and owner review are complete.

## Evidence record

The completed release record must link, at minimum:

- source commit, manifest value, and immutable tag;
- development plan and master checklist status;
- issue, PR, CI, and owner-review links;
- platform matrix and smoke output;
- artifact names, sizes, and SHA-256 checksums;
- scanner/toolchain versions and representative scans;
- SBOM, signing, and provenance status;
- release notes and public release URL;
- known risks, rollback procedure, and final owner decision;
- retrospective and follow-up actions.

## Release-mode guidance

### Preview

A preview may proceed with explicitly documented and owner-approved gaps. It must not be described as stable-ready.

### Release candidate

A release candidate requires disposition of known preview defects, complete artifact identity/checksum evidence, and no open release-blocking P0/P1 defect without explicit owner disposition.

### Stable

A stable release requires complete checklist evidence, owner-reviewed artifacts, platform smoke results, representative scans, owner security/architecture assessment, and an explicit owner go/no-go decision.
