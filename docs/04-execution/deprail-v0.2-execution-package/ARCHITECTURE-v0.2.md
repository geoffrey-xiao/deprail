# DepRail v0.2 Architecture Baseline

## Architectural position

v0.2 continues the Go single-binary architecture established in v0.1. Domain meaning remains inward-facing and independent of Cobra, SQL, scanner-specific types, and transport concerns.

```text
cmd/deprail/
  ↓
internal/app/
  ↓
internal/domain/
  ├── internal/discovery/
  ├── internal/scanplan/
  ├── internal/adapters/osv/
  ├── internal/process/
  ├── internal/artifact/
  ├── internal/normalize/
  ├── internal/policy/
  ├── internal/presenter/
  ├── internal/store/
  └── internal/remediation/ (later)
```

## Stable boundaries

- Entry points translate CLI inputs into application commands.
- Discovery detects targets and completeness without network access or writes.
- Scan plans describe work without embedding scanner-specific types in domain entities.
- Adapters invoke external tools through process ports and retain raw provenance.
- Normalization is deterministic and independently testable offline.
- Presenters separate machine data from diagnostics.
- Release tooling remains outside domain packages.

## v0.2 changes

### Output initialization

Report constructors must create empty slices/maps required by the serialized contract. Serialization code must not rely on nil collection behavior to communicate absence.

### Strict command parsing

Command parsers must reject unexpected positional arguments and unsupported options before application work begins. The parser owns configuration/argument errors; application services own scan failures.

### Version injection

Build metadata enters through a narrow version provider or linker-injected variables. Domain reports may receive release identity as declared run metadata, but stable finding keys must not depend on tag, commit, timestamp, severity label, or evidence order.

### Baselines and diff

Trusted scan results are compared through versioned baseline contracts. Base/head comparison and new/resolved/unchanged classification remain deterministic and independent of input ordering.

### Policy and SARIF

Typed policy evaluation owns completeness, severity, new-risk, and exception-expiry decisions. Policy does not become a general language in v0.2. SARIF rendering remains in presenters and preserves normalized identity and provenance.

### External adapter fixtures

The OSV adapter contract includes checked-in raw fixtures, supported-version metadata, command arguments, exit-code interpretation, output limits, timeout behavior, and malformed-output behavior. Fixtures must represent real supported output rather than approximated test structs.

### Requested-root enforcement

The application resolves and passes the canonical requested root to discovery, scanner execution, and artifact storage. The process current directory cannot implicitly redefine scan scope.

## Security invariants

- Never concatenate shell commands.
- Pass argument arrays directly.
- Resolve paths against the canonical root and reject traversal/symlink escape.
- Bound stdout/stderr and terminate timed-out processes according to platform contract.
- Preserve raw artifacts atomically with restrictive permissions.
- Redact credentials and sensitive environment values.
- Never upload repository source or scanner output automatically.

## Decisions still requiring ADRs

Create an ADR only if v0.2 selects one of these capabilities:

- Baseline or policy persistence.
- Remote publishing/history.
- SBOM generation or artifact signing.
- Additional scanner families.
- Schema-breaking output changes.
- Repository mutation or remediation.

## Verification expectations

Every architectural change requires focused contract tests, failure-path coverage, and a statement of whether serialized meaning is unchanged across platforms.
