# v0.2 Evidence and Review Guide

## Purpose

Evidence proves observable behavior and release readiness. A merged PR is not evidence by itself; the PR must identify the commands, scenarios, artifacts, and results reviewed.

## Issue evidence

Every implementation issue must define:

- required command, test, or smoke scenario;
- expected observable result;
- failure or incomplete-result behavior;
- artifact, log, schema, screenshot, or link to attach;
- reviewer and remaining-risk record.

## PR review evidence

During PR review, the owner checks:

- every acceptance criterion;
- scope and exclusions;
- contract impact;
- required tests and failure cases;
- CI results;
- security and compatibility impact;
- rollback behavior;
- remaining risk.

This review is the owner-acceptance point. After the reviewed PR merges and the required evidence is attached, the issue may move directly to `Done`.

## Evidence by change type

| Change | Minimum evidence |
| --- | --- |
| CLI | Focused invocation, expected output, invalid-input behavior, exit code. |
| JSON/schema | Schema validation, empty/boundary examples, deterministic comparison. |
| Adapter | Real raw fixture, command/working directory, exit matrix, malformed/timeout/missing-tool cases. |
| Security boundary | Hostile path/symlink/metacharacter scenario and safe failure result. |
| Release | Commit, tag, artifact names, checksums, platform smoke output, scanner version, reviewer decision. |
| Documentation/process | Link verification, formatting check, and review against the governing contract. |

## Release evidence

A candidate release record must include:

- source commit and tag;
- CI matrix links;
- artifact names and hashes;
- supported scanner version;
- Linux/macOS/Windows smoke results;
- amd64/arm64 evidence;
- representative repository scans;
- SBOM status;
- signing/provenance status;
- known risks and owner decision.

## Insufficient evidence

The following are not sufficient alone:

- “tests passed” without command and result;
- a merged PR without acceptance review;
- a zero-finding scan with incomplete status;
- a screenshot without the command or fixture;
- a generated artifact without its source commit and checksum;
- a release tag without platform smoke evidence.
