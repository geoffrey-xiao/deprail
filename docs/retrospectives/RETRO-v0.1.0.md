# DepRail v0.1.0 Development Retrospective

**Status:** Development complete for the preview baseline; release gate still open
**Preview tag:** [`v0.1.0-preview.1`](https://github.com/geoffrey-xiao/deprail/tree/v0.1.0-preview.1)
**Baseline commit:** [`f6fb5fd`](https://github.com/geoffrey-xiao/deprail/commit/f6fb5fd05ca1ebaad6bf24a9346b1cb4184e455c)
**Period:** Sprint 0 through Sprint 4 follow-up
**Audience:** Project owner, maintainers, and implementation agents

## 1. Executive summary

The v0.1 development period established a working Go CLI for dependency discovery and OSV-Scanner-backed scanning. The repository now has the documented command surface, secure discovery boundaries, bounded external-process execution, raw artifact retention, deterministic normalization primitives, machine-readable output, cross-platform CI, release workflows, and a preview tag.

The period also exposed important process and contract gaps. The largest were an OSV-Scanner v1 assumption after the project moved to v2, incorrect handling of OSV-Scanner's vulnerability-found exit code, scanner execution initially using the caller's working directory instead of the requested scan root, inconsistent GitHub Project status updates, and release automation that required an explicit pause before the release process was ready.

The work is not yet a final production release. The v0.1 release contract still requires linked acceptance evidence, three real-repository checks, reviewed release artifacts and checksums, and explicit handling of the documented SBOM and signing gaps. The Master Checklist remains the authority for final acceptance.

## 2. Intended outcome and delivered scope

The PRD defined the v0.1 proposition as one command-line workflow for JavaScript, Python, and Java dependency repositories that produces a stable, explainable, machine-readable report. The delivered implementation covers:

- `deprail doctor`, `deprail discover [path]`, and `deprail scan [path]`.
- npm/pnpm/Yarn, requirements/uv/Poetry, and Maven/Gradle discovery boundaries.
- Secure repository walking with symlink containment and incomplete-scope diagnostics.
- Argument-array process execution with timeouts, cancellation, bounded output, and process-group handling.
- OSV-Scanner v2 compatibility (`>=2.0.0 <3.0.0`), including the v2 `scan source` command and v2 vulnerability JSON shape.
- Explicit complete, partial, and failed scan states.
- Content-addressed raw scan artifacts.
- Deterministic alias, severity, fixed-version, and finding normalization work.
- Terminal and JSON presenters with atomic file output.
- Linux, macOS, and Windows CI verification.
- Release-please configuration, tag-triggered artifact workflow, checksums, and protected manual approval gates.
- Project and agent workflow guidance for synchronized branches, issue linkage, labels, and Project status.

Primary evidence includes the v0.1 PRD, release procedure, Master Checklist, merged implementation history, CI runs, and the preview tag:

- [v0.1 PRD](../04-execution/deprail-v0.1-execution-package/PRD-v0.1.md)
- [v0.1 release procedure](../RELEASE-v0.1-PREVIEW.md)
- [v0.1 Master Checklist](../04-execution/deprail-v0.1-execution-package/tracking/MASTER-CHECKLIST.md)
- [OSV-Scanner v2 PR #85](https://github.com/geoffrey-xiao/deprail/pull/85)
- [Project workflow PR #86](https://github.com/geoffrey-xiao/deprail/pull/86)
- [Release baseline PR #88](https://github.com/geoffrey-xiao/deprail/pull/88)
- [Preview tag](https://github.com/geoffrey-xiao/deprail/tree/v0.1.0-preview.1)

## 3. What went well

### 3.1 Contracts were written before implementation

The repository established architecture, requirements, issue contracts, schemas, error codes, and a Master Checklist before the main implementation sequence. This gave the implementation a reviewable boundary and made later failures identifiable as contract mismatches instead of ambiguous behavior.

### 3.2 Security boundaries were treated as product behavior

The implementation consistently treated paths, symlinks, manifests, scanner output, subprocesses, and credentials as hostile inputs. Argument arrays, bounded output, explicit working directories, path containment, and read-only discovery are now part of the design rather than afterthoughts.

### 3.3 Cross-platform verification was present early

The CI matrix runs formatting, vet, tests, and builds on Ubuntu, macOS, and Windows. This caught process and path assumptions earlier than a single-platform workflow would have. Merged v0.1 PRs received all-platform CI evidence before merge.

### 3.4 The v2 mismatch was corrected end-to-end

The initial v1-only adapter was not papered over. The correction updated version compatibility, v2 command invocation, v2 JSON parsing, vulnerability-found exit semantics, scan-root handling, fixture dependencies, documentation, and regression coverage. The real Lodash fixture produced a complete report with findings and fixed versions after the correction.

### 3.5 Release controls became safer during the period

Release artifact publication is gated by a protected approval environment. Automatic release-please checks were paused when the team needed more control, while manual dispatch remains available. The first preview tag is immutable and points to the reviewed `main` commit.

### 3.6 The process learned from operational failures

The project added durable guidance for synchronized branch baselines, issue-to-project status transitions, release version baselines, immutable tags, and stable-versus-prerelease version progression. These rules now live in `AGENTS.md` and tracking documentation instead of only in conversation history.

## 4. What went wrong and what we learned

| Area | Observed failure | Root cause | Immediate response | Remaining improvement |
| --- | --- | --- | --- | --- |
| Scanner compatibility | DepRail rejected installed OSV-Scanner 2.6.0 as unsupported. | The adapter hard-coded an OSV-Scanner v1 contract while the installed tool and desired baseline had moved to v2. | Added v2 compatibility, command invocation, parser support, and tests in PR #85. | Pin and exercise the external scanner contract in a real integration fixture before implementation merge. |
| Scanner exit semantics | A vulnerable package caused `SCANNER_EXIT_NONZERO` and a failed DepRail scan. | OSV-Scanner v2 uses exit code `1` to mean vulnerabilities were found; DepRail treated every non-zero exit as execution failure. | Accepted the documented vulnerability-found result and parsed the returned JSON. | Maintain an explicit external-tool exit-code matrix as part of adapter contract tests. |
| Scanner JSON parsing | The first v2 parser expected a scalar severity and legacy fixed-version location. | The v2 JSON shape was assumed from existing tests instead of verified against real output. | Added v2 severity-array, database severity, affected-range, and fixed-event parsing. | Keep a checked-in real v2 raw-output fixture and validate it against the supported scanner version. |
| Scan root | A scan initially invoked the scanner from the process working directory, causing it to inspect the wrong repository. | The serialized graph intentionally used `.` and that value was reused as the subprocess directory. | Canonicalized the requested scan root for scanner execution and artifact storage. | Add an end-to-end test that invokes `deprail scan` from outside the target repository. |
| Output representation | Successful empty scans serialize `findings` and `errors` as `null` rather than `[]`. | Nil Go slices are encoded directly without report initialization. | Documented that `null` means no entries during manual testing. | Normalize empty collections to arrays before final release; this is clearer and safer for machine consumers. |
| CLI argument handling | `./deprail doctor discover .` ran doctor and silently ignored the extra command. | The doctor flag parser does not reject unexpected positional arguments. | Manual guidance now says to invoke one command at a time. | Reject unexpected doctor arguments with `CONFIG_INVALID` and add a regression test. |
| Fixture maintenance | The `left-pad` fixture could not be installed reliably. | The fixture depended on an obsolete package and did not model the intended vulnerability test. | Replaced it with Lodash 4.17.20 and committed the lockfile. | Add dependency-install and fixture-refresh checks that are offline or use a controlled fixture policy. |
| Local generated files | `node_modules` and `.deprail` appeared as untracked local files. | Ignore rules were incomplete when the fixture was exercised manually. | Added `node_modules/` and `.deprail/` to `.gitignore`. | Add all generated local paths to the repository baseline before fixture work begins. |
| Release automation | Release-please was initially configured to run on every push to `main`. | Automation was enabled before the preview release operating model was settled. | Paused the push trigger; manual dispatch remains available. | Define and review release-mode transitions before enabling automation. |
| Project workflow | Opening a PR did not move the linked Project item to `Review`. | GitHub Project mutation automation was absent; the assumption existed only in documentation. | Manually synchronized items and documented the required CLI procedure. | Decide whether to add approved project automation or retain the explicit procedure with verification. |
| Branch hygiene | Work was started from a branch that was not always freshly synchronized with `main`. | Baseline synchronization was a process expectation, not an enforced first step. | Added `fetch`, `switch main`, `pull --ff-only`, and branch-creation instructions. | Make the baseline check part of the issue-start checklist and PR template. |
| Versioning | The relationship between tags, previews, stable releases, and the CLI version was initially unclear. | Release identity, manifest state, and hard-coded doctor version were not unified. | Created `v0.1.0-preview.1` and documented immutable tag progression. | Inject the tag-derived version into the binary and make release evidence use one version source. |

## 5. Process assessment

### Keep

- Issue-first implementation with explicit scope, risk, acceptance, and reviewer focus.
- One primary outcome per PR.
- Small, deterministic fixtures and focused contract tests.
- CI on every PR and on `main` after merge.
- Human approval for publishing, credentials, permissions, schemas, and external-process behavior.
- Immutable tags and release artifacts built from reviewed `main`.
- Explicit failure states instead of converting failures to empty success.

### Stop

- Do not assume an external tool's version, command shape, output shape, or exit semantics from an older fixture.
- Do not use a broad non-zero exit rule when the tool has a documented findings exit code.
- Do not create a branch from a stale local `main`.
- Do not assume GitHub Projects automatically synchronize with PR lifecycle events.
- Do not treat a merged PR, preview tag, or passing CI run as owner acceptance of the entire release gate.

### Start

- Start every external adapter with a real versioned command/output/exit-code fixture.
- Start every release line with a baseline check covering tags, releases, manifest, and `origin/main`.
- Start every PR with explicit issue linkage, project membership, labels, and status verification.
- Start every CLI contract test with invalid-argument and machine-output-shape cases.
- Start release preparation with an evidence checklist, not only a build command.

## 6. Prioritized improvement plan

| Priority | Improvement | Owner | Exit condition |
| --- | --- | --- | --- |
| P0 | Normalize empty report collections to `[]` and reject unexpected `doctor` arguments. | CLI/Reporting maintainer | Contract tests pass; JSON examples show arrays; invalid doctor arguments return code `2`. |
| P0 | Add a real OSV-Scanner v2 raw-output fixture and an exit-code compatibility matrix. | Adapter maintainer | Fixture covers no findings, findings, malformed output, and scanner failures. |
| P0 | Complete the v0.1 release gate evidence. | Project owner / release maintainer | Master Checklist links CI, artifacts, checksums, manual repository scans, and explicit remaining gaps. |
| P1 | Inject the Git tag and commit into the binary version output. | CLI maintainer | `doctor --format json` reports the same version identity as the release tag. |
| P1 | Add an end-to-end scan-from-outside-root test. | Scan maintainer | Scanner and artifacts stay within the requested repository root. |
| P1 | Decide whether Project status transitions should remain procedural or gain approved automation. | Project owner | Tracking issue records the decision, permissions, rollback, and verification path. |
| P1 | Add a release-mode checklist for preview, RC, and stable releases. | Release maintainer | A release cannot proceed without tag baseline, version decision, CI, artifact, and approval evidence. |
| P2 | Measure the PRD success metrics across representative public repositories. | Project owner / QA | At least 15 repositories are recorded with success, partial, failure, timing, and cause data. |
| P2 | Add dependency-fixture maintenance guidance and controlled refresh commands. | Test/fixture maintainer | Fixture updates are reproducible and do not require committing install directories. |

## 7. Release readiness decision

**Current decision: preview-ready, not final-release-ready.**

The `v0.1.0-preview.1` tag is suitable for controlled preview validation because the implementation is merged, `main` CI passed, and the release workflow has verification and approval gates. It must not be described as a final v0.1.0 release until the following evidence is attached:

- All applicable Master Checklist items have linked acceptance evidence.
- Linux, macOS, and Windows smoke results are recorded for the release artifact set.
- amd64 artifacts and arm64 smoke results are recorded.
- Checksums are reviewed and attached.
- Three real repositories are manually scanned and their results recorded.
- SBOM and signing status remain explicit, with owner-approved preview limitations if still unavailable.
- The final stable version decision is recorded separately from the preview tag.

## 8. Next-release working agreement

Before starting a `0.2.0` line:

1. Confirm the final v0.1 acceptance decision and remaining risks.
2. Confirm the latest stable tag; preview and RC tags do not replace it.
3. Update the release issue, manifest, and target version deliberately.
4. Synchronize local `main` from `origin/main` before creating branches.
5. Preserve the v0.1 contracts unless a reviewed compatibility decision changes them.
6. Require a real external-tool fixture for every adapter contract change.
7. Keep release PR, Project status, issue evidence, and tag identity synchronized.

## 9. Evidence index

- [v0.1 PRD](../04-execution/deprail-v0.1-execution-package/PRD-v0.1.md)
- [Release procedure](../RELEASE-v0.1-PREVIEW.md)
- [Master Checklist](../04-execution/deprail-v0.1-execution-package/tracking/MASTER-CHECKLIST.md)
- [GitHub Project setup](../04-execution/deprail-v0.1-execution-package/tracking/GITHUB-PROJECT-SETUP.md)
- [OSV-003 issue](https://github.com/geoffrey-xiao/deprail/issues/84)
- [OSV-003 implementation PR](https://github.com/geoffrey-xiao/deprail/pull/85)
- [Tracking workflow PR](https://github.com/geoffrey-xiao/deprail/pull/86)
- [Release baseline PR](https://github.com/geoffrey-xiao/deprail/pull/88)
- [Preview tag](https://github.com/geoffrey-xiao/deprail/tree/v0.1.0-preview.1)
