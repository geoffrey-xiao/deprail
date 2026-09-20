# DepRail v0.3.1 Manual Test Guide

## Purpose

Run this guide against the release candidate to collect human, stream, safety, and cross-platform evidence for the v0.3.1 release gate.

Do not mark a scenario passed without saving the command, exit code, stdout, stderr, platform, and binary identity.

## Candidate setup

Run from the repository root:

```bash
go build -o deprail ./cmd/deprail
./deprail --version
shasum -a 256 ./deprail
```

Record:

- commit and tag
- host OS and architecture
- Go version
- OSV-Scanner version
- binary SHA-256

Use a copy of `testdata/fixtures/mixed-repository` for mutation checks. Set `ROOT` to that copy and `BIN` to the built binary.

## Evidence capture convention

For each scenario:

```bash
set +e
"$BIN" <arguments> >local_test/<name>.stdout 2>local_test/<name>.stderr
status=$?
printf '%s\n' "$status" >local_test/<name>.exit
set -e
```

Record the platform, terminal/redirected mode, expected result, observed result, and artifact paths in the release evidence report.

## 1. Help and version

Expected: exit `0`, deterministic plain stdout, empty stderr, no scanner or repository access.

```bash
"$BIN" --help
"$BIN" -h
"$BIN" discover --help
"$BIN" scan --help
"$BIN" doctor --help
"$BIN" diff --help
"$BIN" fix plan --help
"$BIN" policy check --help
"$BIN" --version
```

Confirm every command help page includes a concise `Tip:` and copy/paste `Example:`; confirm `scan --help` documents `--format`, `--output`, `--quiet`, and `--verbose`; confirm `fix plan --help` documents `--output`.

Confirm the recommended first-run sequence is discover, scan, then optional diff/policy/fix-plan workflows.

Negative help cases:

```bash
"$BIN" scna --help
"$BIN" fix typo --help
"$BIN" policy typo --help
```

## 2. Human terminal output

Run in an interactive terminal:

```bash
"$BIN" doctor
"$BIN" discover testdata/fixtures/mixed-repository
"$BIN" scan testdata/fixtures/mixed-repository
"$BIN" scan testdata/fixtures/mixed-repository --quiet
"$BIN" scan testdata/fixtures/mixed-repository --verbose
```

Confirm:

- start notices identify long-running work (`Discovering workspaces...`, `Scanning...`, `Comparing baselines...`, or `Planning remediation...`) where applicable;
- start notices appear on interactive stderr, not machine stdout;
- output is readable when color is unavailable or disabled;
- when color is enabled, only fixed message prefixes are styled;
- repository paths, package names, URLs, and error details remain sanitized data;
- progress corresponds to real lifecycle boundaries;
- no invented percentages, durations, or completion claims;
- complete, partial, failed, and cancelled outcomes are distinct.

## 3. JSON stream purity

```bash
"$BIN" discover testdata/fixtures/mixed-repository --format json >local_test/discover.json 2>local_test/discover.stderr
"$BIN" doctor --format json >local_test/doctor.json 2>local_test/doctor.stderr
"$BIN" scan testdata/fixtures/mixed-repository --format json >local_test/scan.json 2>local_test/scan.stderr
```

Expected:

- each stdout file parses as JSON;
- stdout contains no progress, banners, diagnostics, ANSI escapes, or carriage returns;
- progress and diagnostics, if any, are confined to stderr;
- result meaning and exit codes remain unchanged.

## 4. Redirected non-TTY output

```bash
"$BIN" scan testdata/fixtures/mixed-repository >local_test/redirected.stdout 2>local_test/redirected.stderr
"$BIN" scan testdata/fixtures/mixed-repository --format json >local_test/redirected.json 2>local_test/redirected-json.stderr
```

Check captured files for ANSI escape bytes (`0x1b`), carriage returns, spinner frames, and animation. Expected: stable line-oriented output with no ANSI or carriage-return updates. Redirected output must remain plain even when the invoking environment supports color.

## 5. Quiet and verbose behavior

```bash
"$BIN" scan testdata/fixtures/mixed-repository --quiet >local_test/quiet.stdout 2>local_test/quiet.stderr
"$BIN" scan testdata/fixtures/mixed-repository --verbose >local_test/verbose.stdout 2>local_test/verbose.stderr
"$BIN" scan testdata/fixtures/mixed-repository --quiet --format json >local_test/quiet.json 2>local_test/quiet-json.stderr
```

Expected:

- quiet retains the single interactive `Scanning...` start notice but suppresses progress and successful summaries;
- required errors remain visible;
- verbose adds only safe diagnostics;
- JSON remains valid and uncontaminated in both modes.

## 6. Outcome states

Exercise and capture:

- complete scan with zero findings: only this state may say `No known vulnerabilities found.`;
- complete scan with findings;
- partial scan with one workspace failing;
- failed scan with no trustworthy findings;
- cancelled scan.

Confirm partial output identifies affected scope, failed output exposes stable error codes, and cancelled output does not claim report completion.

## 7. Scanner failure matrix

Using controlled scanner fixtures, exercise:

- scanner missing;
- non-zero scanner exit;
- timeout;
- malformed scanner output;
- output-limit failure;
- partial workspace failure.

Expected: non-zero documented exit where applicable, stable error code, truthful status, preserved successful workspace results for partial scans, and no empty-success fallback.

## 8. Hostile labels and redaction

Exercise repository-derived values containing:

- newline, carriage return, tab, Unicode, and long values;
- ANSI and other terminal control bytes;
- shell metacharacters;
- credential-bearing URLs such as `https://user:token@example.test/path`.

Expected: labels remain data, output stays line-oriented, terminal controls are escaped, credentials are redacted, and no command is executed.

## 9. Cancellation and terminal cleanup

Cancel an interactive scan during discovery or scanner execution using `Ctrl-C`:

```bash
"$BIN" scan testdata/fixtures/mixed-repository
```

Expected:

- cancellation is reported truthfully;
- `OperationCancelled` is the terminal lifecycle event;
- no `ReportReady` claim follows cancellation;
- terminal state is restored;
- stderr does not leave a spinner, partial escape sequence, or corrupted prompt.

Repeat with redirected output and JSON mode.

## 10. Read-only and artifact checks

Before and after representative scans, compare the fixture tree:

```bash
git -C "$ROOT" status --short
find "$ROOT" -type f -print | sort >local_test/tree-before.txt
"$BIN" scan "$ROOT" --format json >local_test/scan-readonly.json 2>local_test/scan-readonly.stderr
find "$ROOT" -type f -print | sort >local_test/tree-after.txt
diff -u local_test/tree-before.txt local_test/tree-after.txt
```

Expected: no repository mutation. Raw artifacts, if produced, remain under the approved `.deprail/` location with expected permissions.

## 11. Cross-platform matrix

Run the applicable sections on Linux, macOS, and Windows. Record:

- platform and architecture;
- terminal versus redirected mode;
- command and exit code;
- stdout/stderr captures;
- whether differences are only styling and not serialized meaning.

The current CI matrix proves build/test execution on all three platforms. Interactive TTY captures must be collected separately where required by the release decision.

## Evidence index

Create one evidence file per candidate under `local_test/` and link it from issue #253. Do not claim scenarios that were not run. Record known limitations, especially missing interactive-TTY captures.
