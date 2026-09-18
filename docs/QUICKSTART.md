# DepRail Quick Start

DepRail scans JavaScript, Python, and Java dependency workspaces through OSV-Scanner. Discovery is local and read-only. Scanning retains raw scanner output under `.deprail/artifacts/` and returns exit code `3` when the result is incomplete or the scanner fails.

## Install and verify

Build the binary with the repository toolchain:

```bash
go build -o deprail ./cmd/deprail
```

Check the environment before scanning:

```bash
./deprail doctor
./deprail doctor --format json
```

`doctor` never installs tools or accesses the network. Install OSV-Scanner separately, then rerun the command. Supported scanner compatibility currently requires an OSV-Scanner v1 version.

## Discover a repository

```bash
./deprail discover .
./deprail discover . --format json
```

Discovery detects npm/pnpm/Yarn, requirements/uv/Poetry, and Maven/Gradle workspaces. JSON paths are repository-relative and use `/` separators. A lockfile-free or malformed workspace is reported as incomplete; it is never treated as a safe empty result.

## Scan a repository

```bash
./deprail scan .
./deprail scan . --format json
./deprail scan . --format json --output scan.json
```

Machine output is written to stdout unless `--output` is provided. Diagnostics and guidance are written to stderr. Existing output files are not overwritten. Scan exit codes:

- `0`: complete execution
- `2`: invalid command or arguments
- `3`: scanner failure or incomplete result

## Troubleshooting

### `SCANNER_NOT_FOUND` or doctor reports unavailable

Install OSV-Scanner using your organization’s approved tool-management process. DepRail does not download or install scanners automatically.

### Exit code `3`

Inspect terminal diagnostics or the JSON `status` and `errors` fields. `partial` means some detected workspaces were omitted or failed. `failed` means no trustworthy scan result was produced. Do not interpret either status as safe.

### JSON output is empty or mixed with progress text

Use `--format json`. DepRail keeps machine data on stdout and diagnostics on stderr. Redirect them separately when scripting.

### Unsupported repository content

DepRail v0.1 supports dependency discovery for npm/pnpm/Yarn, requirements/uv/Poetry, and Maven/Gradle. Containers, SBOMs, licenses, secrets, IaC, Trivy, remediation, and baseline comparison are outside the v0.1 boundary.

## Verification

Repository contributors can run:

```bash
make verify
```

The CI matrix verifies Linux, macOS, and Windows builds and tests. Differences in path and process behavior must not change serialized domain meaning.
