# DepRail v0.2.0 Product Requirements Document

| Attribute | Value |
| --- | --- |
| Version | v0.2.0 |
| Status | Planning baseline |
| Owner | Project Owner |
| Predecessor | v0.1 preview line |
| Target | A reviewed, contract-stable hardening release |
| Core users | Multi-language developers, maintainers, AI coding-tool users |

## 1. Product goal

Make the v0.1 scanning workflow dependable enough for repeated developer and CI use: predictable machine output, strict CLI behavior, explainable external-tool failures, reproducible adapter contracts, and release artifacts whose identity is unambiguous.

v0.2 is not authorized to expand into a server, web UI, remote publishing service, automatic remediation system, or broad policy platform. Those remain future product decisions unless this PRD is amended.

## 2. User jobs

- Run a scan repeatedly and receive the same semantic result for the same inputs.
- Distinguish an empty complete result from a failed or incomplete scan without interpreting implementation details.
- Understand exactly which scanner version, command, exit code, and raw artifact produced a finding.
- Run DepRail from outside the target repository without scanning the wrong working directory.
- Identify the exact CLI build from its tag and commit metadata.
- Validate a release artifact on supported operating systems before trusting it.

## 3. v0.2 scope

### Required hardening

- Normalize empty findings, errors, and artifact collections to JSON arrays.
- Reject unexpected command arguments with the documented configuration error and exit code.
- Add real OSV-Scanner v2 raw-output fixtures and an explicit exit-code matrix.
- Ensure scanner execution, path resolution, and artifact placement use the requested scan root.
- Add scan-from-outside-root end-to-end coverage.
- Inject Git tag and commit identity into CLI version output.
- Preserve deterministic ordering and stable error codes.
- Document and exercise preview, RC, and stable release modes.
- Improve release evidence capture for platforms, artifacts, checksums, manual scans, SBOM, and signing gaps.
- Keep issue/PR metadata and evidence workflow consistent with the reviewed-merge closure policy.

### Candidate capabilities requiring separate approval

- Baseline comparison.
- Policy gates.
- Additional scanner ecosystems.
- SARIF output.
- Remote publishing or history.
- SBOM generation or signing implementation.

Each candidate requires an issue contract, compatibility analysis, security review, and explicit inclusion decision before implementation.

## 4. Out of scope

Unless this PRD is amended, v0.2 excludes web/history, team services, remote publishing, automatic dependency remediation or file mutation, automatic scanner installation, Trivy, container scanning, license/secret/IaC scanning, a hosted vulnerability database, and autonomous repository mutation.

## 5. Stable user flows

```bash
deprail doctor
deprail discover .
deprail discover . --format json
deprail scan .
deprail scan . --format json --output scan.json
deprail scan /path/to/repository --format json
```

The v0.2 line must preserve v0.1 command names, documented exit codes, machine-output separation, path normalization, completeness semantics, and scanner failure behavior unless the compatibility record explicitly approves a change.

## 6. Product behavior

- Empty successful collections serialize as `[]`, never `null`.
- Unexpected arguments fail before work begins and do not silently run a different command.
- `complete`, `partial`, and `failed` retain their v0.1 meanings.
- A zero-finding report is safe to describe as no known vulnerabilities only when status is `complete`.
- Scanner findings, scanner failures, malformed output, timeouts, and missing tools remain distinguishable.
- Repository-relative paths use `/` and cannot escape the canonical root.
- Version output identifies the release tag and source commit when available; development builds identify their non-release state.
- Stable JSON remains schema-valid, deterministic, and free of diagnostic noise.

## 7. Non-functional requirements

- Linux, macOS, and Windows pull-request verification remains required.
- Release artifacts cover the documented amd64 targets and arm64 smoke requirements.
- Discovery remains offline and read-only.
- External processes use argument arrays, explicit working directories, bounded output, deadlines, cancellation, and approved environment values.
- No credentials, source uploads, or sensitive environment values enter reports or logs.
- Adapter fixtures and normalization tests remain runnable offline and deterministically.

## 8. Success metrics

- At least 90% scan success across the selected representative repository set, with every failure categorized.
- No known false-safe result caused by scanner failure, malformed output, timeout, or incomplete discovery.
- Repeated scans of fixed fixtures produce semantically identical JSON apart from declared run metadata.
- A new user can complete doctor, discovery, and first scan using the Quick Start without undocumented steps.
- Release evidence identifies the exact artifact, commit, scanner compatibility, checksum, platform smoke result, and remaining supply-chain gaps.

## 9. End-to-end acceptance

### AC-020-001 Stable empty output

A complete scan with no findings serializes `findings`, `errors`, and artifact collections as empty arrays and validates against the current schema.

### AC-020-002 Strict command arguments

Unexpected arguments to `doctor` and other commands fail with the documented configuration error and do not execute a partial command.

### AC-020-003 OSV-Scanner v2 contract

Checked-in v2 fixtures cover no findings, findings, malformed JSON, incompatible output, non-zero vulnerability-found exit, scanner failure, timeout, and missing tool behavior.

### AC-020-004 Requested-root execution

Running `deprail scan /target/repository` from another working directory scans only `/target/repository` and writes artifacts within the requested root according to the artifact contract.

### AC-020-005 Version identity

A release artifact reports the same version identity as its Git tag and includes the source commit; development builds do not falsely identify as stable releases.

### AC-020-006 Release evidence

A candidate release has passing cross-platform CI, reviewed checksums, platform smoke evidence, manual repository evidence, and explicit SBOM/signing status.

## 10. Release gate

The v0.2 release requires the v0.2 Master Checklist to be complete, all P0 requirements linked to evidence, reviewed pull requests, cross-platform verification, artifact and checksum evidence, representative repository scans, and an explicit owner decision on remaining risk.
