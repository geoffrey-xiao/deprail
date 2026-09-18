# v0.2 Sprint 7 — GitHub Action and Public Preview

## Goal

Validate the guardrail on real pull requests and prepare the v0.2 public preview.

## Planned items

- V02-023 GitHub Action packaging.
- V02-024 Action permissions and caching.
- V02-025 pull-request validation.
- V02-026 public-preview release evidence.

## Dependencies

EPIC-006 and EPIC-007 must provide stable diff, policy, and SARIF contracts. V02-024 depends on V02-023; V02-025 depends on V02-023 and V02-024; V02-026 depends on V02-023 through V02-025.

## Exit criteria

- A real pull request receives a truthful policy result and evidence artifact.
- Permissions, caching, and privacy boundaries are documented.
- Cross-platform and release evidence supports an explicit public-preview decision.
