# v0.3 Master Checklist

## Context and readiness

- [ ] v0.2 release status and carryover risks are recorded.
- [ ] Product, architecture, roadmap, and v0.3 plan are reconciled.
- [ ] Owner and architecture/security reviewer are recorded.
- [ ] Project view `Release v0.3` and milestone `v0.3.0` are active.
- [ ] Issue #198 context preparation is accepted.

## Contracts

- [ ] Remediation-plan schema and examples are reviewed.
- [ ] Candidate states and stable plan identity are frozen.
- [ ] CLI, stdout/stderr, output, and failure behavior are frozen.
- [ ] Package-manager adapter ports are frozen.
- [ ] Mutation, network, package-script, and approval boundaries are explicit.
- [ ] v0.4 handoff contract is documented.

## Implementation evidence

- [ ] npm/pnpm/Yarn planning is contract-tested.
- [ ] Python planning is contract-tested.
- [ ] Maven/Gradle planning is contract-tested.
- [ ] Deterministic JSON and schema validation pass.
- [ ] Stale, incomplete, malformed, unsupported, and unknown cases are covered.
- [ ] Security and hostile-input tests pass.
- [ ] Linux, macOS, and Windows smoke evidence is attached.
- [ ] Repository mutation is proven absent.

## Release review

- [ ] Documentation explains planning-only behavior.
- [ ] Remaining risks and unsupported cases are explicit.
- [ ] Owner decision is recorded.
- [ ] Security/architecture review is recorded.
- [ ] Release evidence links all acceptance rows.
