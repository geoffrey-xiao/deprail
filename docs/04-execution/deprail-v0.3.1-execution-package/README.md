# DepRail v0.3.1 Execution Package

**Status:** Planning baseline; implementation requires Definition of Ready approval
**Version:** v0.3.1
**Predecessor:** [v0.3 execution package](../deprail-v0.3-execution-package/README.md)
**Plan:** [v0.3.1 development plan](../../03-planning/deprail-development-plan-v0.3.1.md)
**UX source plan:** [CLI UX refactor plan](../../03-planning/deprail-cli-ux-refactor-plan.md)
**Context:** [Issue #238](https://github.com/geoffrey-xiao/deprail/issues/238)

## Purpose

This package is the durable execution context for the v0.3.1 compatibility-preserving CLI UX refinement. It defines presentation contracts, terminal behavior, evidence, and release gates without changing scan semantics or authorizing repository mutation.

## Source order

1. [PRD-v0.3.1](PRD-v0.3.1.md)
2. [development plan](../../03-planning/deprail-development-plan-v0.3.1.md)
3. [architecture](ARCHITECTURE-v0.3.1.md)
4. [requirements](requirements/)
5. [tracking](tracking/)
6. Implementation issues only after Definition of Ready approval

Product, architecture, and roadmap baselines remain [`docs/01-product/deprail-product-design-v1-ai.md`](../../01-product/deprail-product-design-v1-ai.md), [`docs/02-architecture/deprail-architecture-and-tech-stack-v1.md`](../../02-architecture/deprail-architecture-and-tech-stack-v1.md), and [`docs/03-planning/deprail-roadmap-v1.md`](../../03-planning/deprail-roadmap-v1.md).

## Release boundary

```text
v0.3: finding -> candidate analysis -> reviewable plan
v0.3.1: clear, safe, cross-platform terminal presentation
v0.4: approved plan -> isolated apply -> verify -> rescan -> patch evidence
```

No v0.3.1 command may edit manifests or lockfiles, install packages, execute package-manager scripts, create worktrees, run remediation verification, publish data, open pull requests, or mutate a repository.

## Package status

This package intentionally contains no epics or implementation issue files. The release remains in planning until the owner and architecture/security reviewer approve the contract, evidence matrix, dependencies, and Definition of Ready.
