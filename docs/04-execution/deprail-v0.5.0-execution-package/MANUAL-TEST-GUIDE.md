# v0.5 Local History Manual Verification

Contract: [H05-008 / #427](issues/H05-008-local-history-release-evidence.md). Requirements: [test strategy](requirements/TEST-STRATEGY.md), [compatibility matrix](requirements/COMPATIBILITY-MATRIX.md), [general release checklist](../../RELEASE-CHECKLIST.md). This is a procedure, not executed evidence or release authority.

## Admission and identity

Before formal execution, record owner acceptance or explicit disposition for H05-001–007, including storage permissions/dependencies, frontend attribution/provenance and browser/AT gaps. The QA correction gate is resolved by [#432 acceptance](https://github.com/geoffrey-xiao/deprail/issues/432#issuecomment-5913275072); retain the original timeout investigation as history. A merged PR is not proof that its unchecked criteria passed.

Current admission requires both corrections on the same reviewed main candidate: #455 timed-apply is owner-merged as 9181e76; #453 artifact fix was merged into its temporary base, not main, so #452 now requires main-targeted #456 cutover. See [readiness assessment](tracking/H05-008-READINESS.md). Fixed startup timing is not mutation-ready evidence. Retain earlier failures; no rerun to erase them. Preparation #445 is Draft and formal #427 is Blocked pending artifact cutover and the original acceptance/evidence gates.

Use an owner-reviewed immutable main commit. Record `go version`, Node/npm versions, OS/architecture, scanner/version, binary `--version`, SHA-256 and build command. Select the actual intended preview tag separately; no preview suffix is prescribed here. Do not tag or publish during verification. Use the owner's authoritative release-binded artifacts when available; locally built candidates are not published binaries or provenance-attested artifacts.

Build inputs: Go `.go-version`, Node `22.23.3`, npm `10.9.9`, committed Go modules and frontend lock. From clean checkout:

```sh
make verify
make test-integration
npm --prefix web test
```

Record commands, outputs and exit codes separately. No passing rerun erases a failed run. On each supported target build the entrypoint with `go build -trimpath`, supplying `internal/buildinfo.Version`, `Tag` and `Commit` through linker flags. Cross-build targets: linux/amd64, darwin/amd64, darwin/arm64, windows/amd64. Launch each target natively or name the emulator and its limitations. Keep binary hashes, raw/gzip asset sizes, asset-only growth comparison, build time and listener-ready measurement; use H05-007 budgets, not a 1 MiB cap on SQLite/runtime growth.

## Isolated repositories and user storage

Never use the operator's only history database. Use a disposable OS user/profile or isolated VM with an empty user configuration root outside all scanned repositories. On macOS the path derives from the user's configuration directory, not `XDG_CONFIG_HOME`; verify the actual root before execution. Keep history DB, matching WAL/SHM, external artifact root and scan reports inside the disposable environment. Restrict permissions and avoid real tokens or credential-bearing URLs in fixtures.

Copy the tracked fixture sources to separate temporary repositories; exclude pre-existing `.deprail` artifact directories. Use `npm-basic`, `python-requirements`, `java-maven`, then `mixed-repository`. Do not run package installation scripts or synthesize vulnerabilities. Record fixture source commit and content digests. Store JSON/logs outside source fixtures. Discovery is offline; scans may use the scanner's advisory network, which must be recorded and explicitly allowed by the operator. Missing/offline database failures are errors, not complete zero-finding scans.

Create deterministic before/after manifests of relative path, file type, file content digest and symlink target. Include source manifests, lockfiles and package-manager state; exclude only explicitly declared DepRail raw-artifact output from scan comparisons. For console-only comparisons also compare the DB/WAL/SHM/artifact set, explaining any SQLite read-side lifecycle changes independently of repository writes.

## Representative capture and query

Run the actual candidate against each copied repository, with reports written outside it:

```text
deprail doctor --format json
deprail discover <fixture-root> --format json
deprail scan <fixture-root> --format json
deprail scan <fixture-root> --format json --save-history
```

Record exact command, scanner/network context, exit, stdout/stderr, raw artifact digests, report status and saved occurrence identity. Without `--save-history`, verify no history root/DB is initialized. With the flag, verify a validated entry can be reopened through the real console API. Compare established scan JSON meaning and exit behavior; do not demand identical occurrence timestamps across separate runs. Save the same fixture twice and verify distinct historyEntryIDs; do not confuse repeatable sourceScanID with occurrence identity.

Happy results require all detected targets complete. Exercise unsuccessful scans with an absent/incompatible scanner, malformed/limited output, timeout and cancellation in a separately identified controlled scanner harness. Such cases verify failure paths, not real vulnerability detection. Exercise a mixed successful/failed target case and preserve successful workspaces. For each available report compare the report status, operation outcome, workspace context, findings and provenance with the stored/API/UI representation. Absent reports/findings/workspaces remain unavailable, not empty or complete. Zero findings may be described as no known vulnerabilities only with complete report status. If a real fixture has no findings, record that limitation rather than manufacture findings.

## Packaged console and privacy

Start the actual candidate in an interactive terminal:

```text
deprail web --artifact-root <existing-trusted-artifact-root>
```

Press Enter to open the process bootstrap session. Noninteractive runs must explicitly use `--open`. Do not copy the fragment bearer into public logs, screenshots, request traces or shell history. Check the printed URL is bare loopback with an ephemeral port. The API is authenticated, read-only GET; no CORS, repository browsing, scan trigger, deletion or remediation actions.

Verify history/list/detail/about, cursor pages, Back/Forward, refresh, missing UUID, loading/read failures, empty history, complete/partial/failed/cancelled operations, absent report/workspaces/findings, missing/digest-mismatched artifacts. Compare the same saved entry through API and UI. Confirm labels/hostile markup render as inert text. Do not inject arbitrary production diagnostics into validated history; use legitimate typed records and separately identify negative-response injections.

Include the default command without any artifact-root flag. A saved digest-bearing detail must remain readable and show unavailable/unverified evidence, not a panic/connection failure or false verification. Then configure the explicit trusted root and confirm the same digest becomes verified only for matching bytes, with missing/mismatch states for absent/corrupt artifacts. Never infer an artifact root from repository/history data.

Verify fragment scrubbing before requests, memory-only session, no credential cookie/storage/referrer, safe reload/new-tab recovery and same-tab credential replacement. Check strict Host/Origin and forwarded-header rejection, absent/wrong bearer, unsupported methods/API version/query, hostile cursor/UUID/encoded paths, unlisted static assets and host-file canary refusal. The 1,048,576-byte API cap must fail or page whole records; never accept truncation. Missing/skewed compiled assets must fail without disk/network/stale fallback while existing CLI remains usable. Repeat repository manifests around console-only use; explain every difference.

## Browser and assistive-technology matrix

Use current stable Chrome/NVDA on Windows, Safari/VoiceOver on macOS, Firefox/Orca on Linux. Record exact browser, AT and OS versions and binary/source identity, not inferred installed versions. Each pairing must exercise keyboard-only controls, visible focus, heading focus and restoration, named landmarks, error/status announcements, cursor navigation, hostile text, reduced motion, 320 CSS px reflow and physical 200% zoom. Inspect rendered contrast and controls. Chromium automation/a11y snapshots are supplemental and do not prove VoiceOver/NVDA/Orca behavior. Unavailable combinations stay unchecked with an owner and explicit preview disposition or no-go.

## Storage and offline recovery

Use only disposable history copies. Verify owner-only POSIX permissions or protected same-user Windows ACL, bounded lock/admission failures, interrupted save atomicity, future schema/corrupt-row refusal, absent store without creation and reopen semantics. Do not fake a prior migration: no prior DepRail history schema is supported. Actual filesystem disk-full and forced SQLite page-limit tests are different evidence.

Recovery rehearsal:

1. Stop console and writers; preserve the original DB, matching WAL/SHM, artifacts and backups. Copy the whole stopped set to a separate location, hash it and retain originals unchanged.
2. Inspect copies with SQLite tooling in read-only mode: `PRAGMA user_version`, `PRAGMA quick_check`. Record exact tool/version and output; never run repair/reset/downgrade on the source.
3. Select a previously validated online backup, restore to a new separate disposable profile, verify supported schema, row validation and artifact digests, then launch the matching binary against the restored copy.
4. Compare saved entry identities/data before/after. Missing artifact bytes remain unavailable; no substitution or deletion. Preserve the only copies if backup validation fails. Stop restored processes and re-hash originals to prove preservation.

No validated backup available means the restore gate did not pass. No automatic restore, repair, eviction or deletion is permitted.

## Supply-chain, acceptance and publication boundary

Inventory each actual candidate artifact's name, size, source/build identity and SHA-256. Record SBOM, signature, certificate and provenance individually as supplied-and-verified, supplied-unverified, unavailable or explicitly owner-deferred. A checksum proves bytes, not publisher authenticity. The H05-006 SBOM covers one host's installed frontend packages, not every cross-platform lock entry or the entire binary. Review complete Go/frontend licenses, the disclosed x/text advisory, caniuse-lite attribution and unverified publisher attestations; a zero-result npm audit does not settle those risks.

Map every FR-501–511, SEC-01–13, STORE-01 and release-checklist row to a dated exact command/scenario, result, artifact and reviewer decision. Distinguish historical component evidence, current candidate runtime evidence and missing platform evidence. Owner technical acceptance, security/architecture assessment and release go/no-go are distinct records even for the same reviewer. Mark no-go whenever a critical scenario or prerequisite remains unaccepted. Publication, immutable tag creation and retrospective start only under separate owner authorization.
