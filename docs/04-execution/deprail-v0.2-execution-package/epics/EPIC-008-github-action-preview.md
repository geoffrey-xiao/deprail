# EPIC-008 GitHub Action and Public Preview

## Outcome

Validate the v0.2 guardrail on real pull requests and publish a reviewable public preview with explicit permissions and evidence.

## Scope

GitHub Action packaging, least-privilege permissions, caching, base/head pull-request flow, evidence artifacts, cross-platform smoke, and preview release acceptance.

## Out of Scope

Hosted DepRail services, source upload by default, automatic remediation, automatic merge, broad enterprise integrations, and signing/SBOM implementation beyond documented evidence status.

## Child issues

- V02-023 GitHub Action packaging.
- V02-024 Action permissions and caching.
- V02-025 Pull-request validation.
- V02-026 Public-preview release evidence.

## Dependencies

EPIC-006 diff, EPIC-007 policy/SARIF, and prerequisite release evidence.

## Acceptance

- A real pull request receives a deterministic new-risk result.
- Action permissions and cache boundaries are documented and least-privilege.
- Full evidence remains available without uploading repository source by default.
- Cross-platform and installation smoke evidence is recorded.
- The v0.2 public-preview decision records remaining risk and rollback.
