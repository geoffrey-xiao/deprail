# DepRail v0.3 Execution Package

**Status:** Planning baseline; implementation requires Definition of Ready approval
**Version:** v0.3.0
**Predecessor:** [v0.2 execution package](../deprail-v0.2-execution-package/README.md)
**Plan:** [v0.3 development plan](DEVELOPMENT-PLAN-v0.3.md)

## Purpose

This package is the durable execution context for v0.3 remediation planning. It preserves v0.2 contracts and adds only reviewable, deterministic remediation plans. It does not authorize repository mutation.

## Source order

1. [PRD-v0.3](PRD-v0.3.md)
2. [development plan](DEVELOPMENT-PLAN-v0.3.md)
3. [architecture](ARCHITECTURE-v0.3.md)
4. [requirements](requirements/)
5. [tracking](tracking/)
6. Epics and issues, after Definition of Ready

Repository workflow rules remain in [`AGENTS.md`](../../../AGENTS.md). Product, architecture, and roadmap baselines are [`docs/01-product/deprail-product-design-v1-ai.md`](../../01-product/deprail-product-design-v1-ai.md), [`docs/02-architecture/deprail-architecture-and-tech-stack-v1.md`](../../02-architecture/deprail-architecture-and-tech-stack-v1.md), and [`docs/03-planning/deprail-roadmap-v1.md`](../../03-planning/deprail-roadmap-v1.md).

## Release boundary

```text
v0.3: finding -> candidate analysis -> reviewable plan
v0.4: approved plan -> isolated apply -> verify -> rescan -> patch evidence
```

No v0.3 command may edit manifests or lockfiles, install packages, execute package-manager scripts, create worktrees, run verification, publish data, open PRs, or mutate a repository.

## Package status

No epics or implementation issues are created by this package. Issue #198 is the initial context-preparation item and must be completed before implementation planning begins.
