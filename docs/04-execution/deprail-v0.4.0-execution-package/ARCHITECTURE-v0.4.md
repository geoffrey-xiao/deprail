# DepRail v0.4 Architecture

**Status:** Planning draft; ADR-0002 is the accepted isolation decision

## Flow

`fix apply` validates a versioned plan, binds approval to the canonical source commit and authorized paths, creates a temporary detached worktree, applies only adapter-owned mutations, runs bounded verification, rescans, classifies transitions, persists evidence, and discards the worktree.

## Boundaries

- `cmd/deprail`: translates CLI input/output only.
- `internal/app`: owns application sequencing and mutation policy.
- `internal/domain`: owns deterministic plans, findings, transitions, and outcomes.
- `internal/remediation`: owns isolation, mutation, verification, and evidence ports.
- `internal/adapters`: owns package-manager and scanner integrations.
- `internal/process`: owns direct-argv bounded subprocess execution.
- `internal/artifact`: owns content-addressed raw evidence.

Domain packages MUST NOT import Cobra, SQL, shell commands, or scanner-specific types. Every subprocess receives an argument array, explicit cwd, approved environment, timeout, cancellation, and output limits.

## Isolation invariants

The caller root is canonicalized and symlink escapes are rejected. The default workspace is a temporary detached Git worktree created from the reviewed source commit. Writes are limited to authorized workspace paths and use atomic replacement. Cleanup is idempotent and records success or failure. The caller worktree is never mutated.

## Failure model

Approval mismatch, missing tools, non-zero exits, timeout, cancellation, malformed output, output limits, verification failure, rescan failure, and cleanup failure are explicit outcomes. No failure is represented as an empty successful result.
