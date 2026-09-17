# ADR-0001 Go Core and Dependency Boundaries

- Status: Accepted baseline
- Date: 2026-09-18
- Owners: Project owner
- Related issues: S0-001

## Context

DepRail must provide a portable local CLI for multi-language dependency scanning while keeping discovery, scanner execution, normalization, and reporting deterministic and testable. CLI, storage, and scanner implementations will evolve independently. Repository contents and scanner output are untrusted inputs.

## Decision

Use a Go single-binary core with module path `github.com/geoffrey-xiao/deprail`. Keep dependencies directed inward:

- `cmd/` contains entry-point wiring only.
- `internal/app/` coordinates application services.
- `internal/domain/` owns deterministic entities and rules.
- `internal/discovery/`, `internal/scanplan/`, and `internal/adapters/` implement ports and integrations.
- `internal/process/` owns bounded, cancellable subprocess execution.
- `internal/artifact/` owns content-addressed raw results.
- `internal/normalize/` owns canonical identities, aliases, evidence, and stable keys.
- `internal/presenter/` owns terminal and machine output.
- `internal/policy/`, `internal/remediation/`, and `internal/store/` remain isolated until their roadmap stages.

Domain code must not import Cobra, SQL, or scanner-specific types. External processes use argument arrays, explicit working directories, deadlines, cancellation, bounded output, and approved environment variables. Public serialized contracts live under `schemas/` and use explicit versions.

## Alternatives considered

- A script-first implementation: rejected because cross-platform process control, deterministic contracts, and a distributable single binary are core requirements.
- A framework-centered architecture: rejected because domain behavior must remain independent of CLI, storage, and transport choices.
- In-process scanner plugins: rejected because untrusted extensions must not share the process trust boundary.

## Consequences

### Positive

- One portable local binary with a small operational surface.
- Domain behavior can be tested without scanners, databases, or CLI wiring.
- Scanner failures, raw evidence, and schema evolution have explicit boundaries.
- Future CLI, web, API, and MCP entry points can reuse application services.

### Negative

- Initial interfaces and package boundaries require discipline before feature work.
- Adapters need contract fixtures and controlled process integration tests.
- Future storage and web layers must preserve the inward dependency direction.

## Compatibility and migration

The v0.1 implementation begins with discovery, the OSV adapter, artifact storage, normalization, and terminal/JSON output. Later policy, remediation, storage, web, and MCP work must add adapters or application services without moving scanner or transport types into the domain package.

## Validation

- `go test ./...` must pass for the baseline module and all committed Go packages.
- S0-001 review must verify module boundaries, failure semantics, security assumptions, and cross-platform portability.
- Subsequent schema and CLI contracts must link back to this decision when they alter these boundaries.
