# DepRail v0.1 PRD Multi-language Dependency Scanning Prototype

| Attribute | Value |
| --- | --- |
| Version | v0.1 |
| Status | Implementation Ready |
| Owner | Project Owner |
| Target duration | Sprint 0 plus S1-S4, approximately nine weeks |
| Core users | Multi-language developers, open-source maintainers, AI coding-tool users |

## 1. Product Goal

v0.1 tests one proposition: a user can run one command in a repository containing JavaScript, Python, and Java and receive a stable, explainable, machine-readable dependency vulnerability report.

It must correctly detect project boundaries and authoritative files, invoke OSV-Scanner through a controlled process, normalize raw output, distinguish complete/partial/failed status, produce terminal and JSON reports, and remain semantically deterministic for identical inputs.

## 2. Target Users and Jobs

### Persona A Multi-language Developer

Scans a frontend, Python service, and Java service together without learning three audit tools.

### Persona B Open-source Maintainer

Reproduces a report and sees the affected component, dependency path, fixed version, and evidence source.

### Persona C AI Coding-tool User

Consumes stable JSON rather than parsing colored or changing terminal text.

## 3. User Stories

- Run `deprail discover .` to confirm scan scope.
- Run `deprail scan .` to get a short action-oriented report.
- Consume versioned JSON without diagnostic noise.
- Retain scanner and vulnerability-database provenance for every finding.
- Detect partial failure through visible status and exit code.
- Replay normalization offline from saved raw results.

## 4. v0.1 Scope

Must include discovery, scan, and doctor commands; npm/pnpm/Yarn, requirements/uv/Poetry, and Maven/Gradle discovery; OSV adapter; core schemas; PURL-first identity; alias merging; terminal and JSON; completeness; and Linux/macOS/Windows smoke tests.

Excluded: baseline diff, policy gates, SARIF, remediation, file mutation, Git worktrees, web/history, Trivy, container/SBOM/license/secret/IaC scanning, a vulnerability database, remote publishing, and automatic scanner installation.

## 5. Core Flows

```bash
deprail doctor
deprail discover .
deprail discover . --format json
deprail scan .
deprail scan . --format json --output scan.json
deprail scan . --verbose
```

Doctor checks the DepRail version, operating system, configuration, OSV-Scanner presence, and compatibility without printing secrets. Discovery performs no network access or writes. Scan discovers workspaces, creates a plan, runs OSV-Scanner, saves raw artifacts, normalizes, and emits a report.

## 6. Product Behavior

Discovery emits repository-relative slash-separated JSON paths, refuses symlink escape, continues past one malformed manifest, diagnoses conflicting lockfiles, marks lockfile-free projects incomplete, and orders results by path, ecosystem, and package manager.

Scanner execution uses argument arrays, cancellation, timeouts, bounded output, recorded metadata, and stable errors. Non-zero exit, timeout, bad JSON, and incompatible versions are not successful empty scans. Raw output is digest-addressed for replay.

Normalization prefers PURL, merges only the same component/version and alias set, preserves source severities, records fixed-version provenance, and is order-independent. Terminal output is concise; verbose output preserves evidence. JSON validates against its schema. Machine output and diagnostics remain separated.

## 7. Completeness Model

| Status | Definition | User meaning |
| --- | --- | --- |
| complete | Every detected target was scanned and parsed successfully | Report covers detected scope |
| partial | At least one target succeeded and another was omitted, failed, or incomplete | The repository must not be considered safe |
| failed | No trustworthy result or a core execution failure | The report cannot support a security decision |

Zero findings may be described as no known vulnerabilities only for `complete` status.

## 8. Non-functional Requirements

Support Linux, macOS, Windows and release smoke tests on amd64/arm64. Target discovery P95 under two seconds and a warm full scan under two minutes for common repositories. Domain output is deterministic except declared run metadata. Do not run install scripts, concatenate shells, escape repository paths, upload source, or log credentials. Adapter execution and normalization are independently testable offline.

## 9. Success Metrics

- At least 85% scan success across 15 public repositories, rising to 90% before v0.2.
- JavaScript, Python, Java, and mixed-monorepo fixtures pass.
- Repeated normalization is stable.
- Scanner failure never creates a false-safe report.
- A new user completes the first scan in five minutes excluding an initial database download.

## 10. End-to-end Acceptance

### AC-E2E-001 Mixed Repository

Given npm, uv, and Maven workspaces, scanning produces three workspaces, a schema-valid complete report, and component, vulnerability, and evidence for every finding.

### AC-E2E-002 One Workspace Fails

A malformed lockfile does not stop other workspaces; status is partial, a stable workspace-scoped error is present, and no safe message is shown.

### AC-E2E-003 Scanner Missing

A missing OSV-Scanner returns an execution-failure exit code and installation guidance on stderr, without an empty success report.

### AC-E2E-004 Deterministic Replay

Permuting finding and evidence input order for a fixed raw fixture produces identical domain JSON except allowed run metadata.

### AC-E2E-005 Hostile Path

External symlinks and shell metacharacters never escape the repository or execute additional commands; a security diagnostic is emitted.

## 11. Release Gate

All P0 requirements trace to tests, the Master Checklist is complete, three-platform smoke tests pass, three real repositories are manually verified, no Critical or High release defect is open, and release artifacts include checksums, an SBOM, and either signatures or an explicit preview-stage signing gap.
