# DepRail Project Documentation Package

This archive consolidates the current product, architecture, delivery-plan, and v0.1 execution documentation for DepRail. The public project name, CLI commands, configuration paths, and code entry points consistently use DepRail.

## Recommended Reading Order

1. `00-overview/PROJECT-NAMING.md` - project and technical naming rules.
2. `01-product/deprail-product-design-v1.docx` - product position, users, scope, roadmap, and commercial boundaries.
3. `01-product/deprail-product-design-v1-ai.md` - AI-readable product design source.
4. `02-architecture/deprail-architecture-and-tech-stack-v1.md` - system architecture, technology choices, security boundaries, and version roadmap.
5. `03-planning/deprail-development-plan-v1.md` - Sprint plan, quality gates, and human-AI delivery model.
6. `04-execution/deprail-v0.1-execution-package/README.md` - the working entry point for v0.1.

## Package Structure

```text
deprail-documentation-v1-en/
|-- README.md
|-- 00-overview/
|   `-- PROJECT-NAMING.md
|-- 01-product/
|   |-- deprail-product-design-v1.docx
|   `-- deprail-product-design-v1-ai.md
|-- 02-architecture/
|   `-- deprail-architecture-and-tech-stack-v1.md
|-- 03-planning/
|   `-- deprail-development-plan-v1.md
`-- 04-execution/
    `-- deprail-v0.1-execution-package/
        |-- README.md
        |-- PRD-v0.1.md
        |-- requirements/
        |-- epics/
        |-- issues/
        |-- tracking/
        |-- prompts/
        `-- templates/
```

## Document Responsibilities

| Category | Question answered | Primary readers |
| --- | --- | --- |
| Product design | Why build it, for whom, and what is in or out of scope? | Owner, contributors, prospective users |
| Architecture | How is it implemented and where are the security boundaries? | Developers, architecture reviewers, AI agents |
| Development plan | How does the product evolve and how is each stage accepted? | Owner, maintainers |
| v0.1 execution package | What should be built now, how is work split, and how is it verified? | Developers, AI agents, reviewers |

## Current Execution Entry Point

Start implementation with:

- `04-execution/deprail-v0.1-execution-package/tracking/SPRINT-0.md`
- `04-execution/deprail-v0.1-execution-package/tracking/MASTER-CHECKLIST.md`
- `04-execution/deprail-v0.1-execution-package/issues/S0-001-repository-baseline.md`
- `04-execution/deprail-v0.1-execution-package/issues/S0-002-toolchain-and-ci.md`
- `04-execution/deprail-v0.1-execution-package/issues/S0-003-project-templates.md`

`S0-002-toolchain-and-ci.md` includes the Sprint 0 CI and GitHub Actions bootstrap requirements.
