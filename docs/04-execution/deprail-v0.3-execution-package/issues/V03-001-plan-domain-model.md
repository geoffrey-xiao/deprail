# V03-001 Define Remediation-Plan Domain Model

## Planning metadata

- Type: feature
- Area: normalization
- Priority: P1
- Risk: R2
- Epic: EPIC-001
- Dependencies: Issue #198, v0.3 data contract

## Goal
Define deterministic domain entities for remediation plans, candidates, risks, structured future commands, verification, rollback, and provenance.

## Acceptance criteria

- [ ] Domain types cover every required plan field.
- [ ] Candidate states are explicit: recommended, viable, rejected, unavailable, unknown.
- [ ] Stable plan identity excludes timestamps, prose, evidence order, and presentation.
- [ ] Commands include executable, arguments, and contained repository-relative working directory.
- [ ] Unknown evidence cannot become a safe recommendation.
- [ ] Domain package imports no Cobra, SQL, scanner-specific, package-manager SDK, or LLM types.
- [ ] Unit/property tests prove canonicalization and order independence.

## Evidence
Domain tests, serialized examples, and a review of dependency direction.

## Exclusions
No CLI, filesystem mutation, package-manager execution, network metadata, or schema migration.
## GitHub tracking

- Issue: [#206](https://github.com/geoffrey-xiao/deprail/issues/206)
- Parent epic: [#202](https://github.com/geoffrey-xiao/deprail/issues/202)
