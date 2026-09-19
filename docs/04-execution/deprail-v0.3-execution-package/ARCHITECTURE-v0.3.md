# DepRail v0.3 Architecture

## Architectural position

v0.3 extends the existing Go single-binary core with deterministic remediation planning. Domain code remains independent of Cobra, SQL, scanner-specific types, package-manager SDKs, and LLMs.

```text
cmd/deprail
  -> internal/app
    -> internal/domain
      -> internal/remediation
        -> planning ports
          -> read-only package-manager adapters
    -> internal/presenter
    -> internal/artifact / store
```

## Responsibilities

| Component | Responsibility |
| --- | --- |
| `internal/remediation` | Finding context, candidate model, stable plan identity, risks, plan generation |
| planning adapters | Read-only ownership, constraints, candidate compatibility, command templates, expected files |
| `internal/app` | `fix plan` orchestration, source validation, stale-input handling |
| `internal/presenter` | Terminal and JSON plan output |
| `schemas` | Versioned remediation-plan schema and examples |
| artifact/store | Plan persistence, source digests, safe atomic writes |

## Trust boundary

Manifests, lockfiles, scanner findings, package metadata, and registry responses are hostile inputs. Canonical root containment, symlink rejection, bounded parsing, credential redaction, no shell execution, and explicit network policy are mandatory.

## Planning flow

1. Load a finding and source scan identity.
2. Resolve workspace and dependency ownership.
3. Parse constraints and current resolution.
4. Enumerate candidates from available evidence.
5. Classify compatibility and migration risks.
6. Build a deterministic plan without writing or executing.
7. Render or persist schema-valid output.

## v0.3 versus v0.4

v0.3 produces commands and verification requirements as structured data only. v0.4 may consume an approved plan in an isolated worktree, execute package-manager changes, run tests/builds, rescan, and produce patch evidence.

## Compatibility

Existing scan, diff, policy, SARIF, exit-code, path, provenance, and output contracts remain unchanged. New schemas are additive and versioned. Any breaking change requires an ADR, converter, examples, and compatibility tests.

## Architecture review triggers

Require explicit review for network metadata, package-manager parsing changes, new write paths, schema-breaking changes, agent access, new scanner families, or any operation that could execute repository-provided scripts.
