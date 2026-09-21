# v0.4 Test Strategy

## Contract tests

Validate plan/approval binding, stable JSON schemas, path authorization, error codes, redaction, finding transitions, and artifact digests.

## Integration tests

Use controlled temporary Git repositories and child processes to cover detached worktrees, source immutability, hostile paths, atomic writes, non-zero exits, timeout, cancellation, process-tree termination, output limits, cleanup retry/idempotence, and malformed scanner output.

## End-to-end tests

Exercise representative JavaScript, Python, and Java repositories through dry-run, approved apply, verification, rescan, partial failure, and rollback flows. Tests remain offline and use fixed fixtures.

## Cross-platform matrix

Run required smoke scenarios on Linux, macOS, and Windows. Compare serialized domain meaning, not platform-specific temporary paths or process wording.

## Manual release evidence

Use the v0.4 manual guide to record binary identity, exact commands, exit codes, source-tree comparison, evidence artifact paths, reviewer, and remaining risk. Tests alone do not satisfy release approval.
