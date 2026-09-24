# ADR-0005: Solo-Owner Review and Approval Policy

- Status: Accepted
- Date: 2026-09-24
- Owners: `@geoffrey-xiao`
- Related issues: [#392 v0.5 readiness](https://github.com/geoffrey-xiao/deprail/issues/392)
- Supersedes: Independent-review and dual-approval gates in active project guidance and the v0.5.0 planning/execution documents

## Context

DepRail is maintained as a personal project by one owner. The current v0.5.0 plan and execution package require an independent architecture/security reviewer and separate owner approval before contract acceptance, implementation issue creation, runtime work, and release. The named account is not a repository collaborator, and waiting for another person has blocked the owner-directed work. The owner has explicitly decided that a second-person review or authorization is not a project requirement.

This decision changes governance, not the product scope or technical security boundary. The owner remains responsible for contract decisions, risk acceptance, implementation, and releases. Automated verification, evidence, safe defaults, and explicit failure behavior remain necessary because the application handles hostile repository data and local API/storage boundaries.

## Decision

1. The project owner is the sole required human reviewer and approver for plans, ADRs, schemas, compatibility decisions, security contracts, implementation PRs, Definition of Ready, and releases.
2. Independent reviews, collaborator reviews, and second-person approvals are optional and advisory. Their absence MUST NOT block work, merge, issue closure, or release when the owner's review, required CI, acceptance evidence, and release-specific gates are complete.
3. The owner records review and acceptance in the relevant ADR, issue, PR, checklist, or release evidence. Security/architecture assessment and owner acceptance remain distinguishable evidence items, but the same owner may perform and record both; a second identity is not required.
4. Issue and PR reviewer fields name `@geoffrey-xiao` by default. External reviewers may be invited voluntarily; an unanswered request or unavailable collaborator assignment is not a blocker.
5. High-risk changes still require the owner to inspect the change and record compatibility, security, failure, rollback, and residual-risk evidence as applicable. Required automated tests, CI checks, schema validation, platform evidence, and hostile-input checks are not waived by this policy.
6. This policy does not authorize an agent to merge or publish on the owner's behalf. Such actions still require the owner's explicit direction in the current conversation; that is owner control, not an external-review gate.
7. Existing historical decisions remain factual. Documents describing prior independent-review requirements are superseded for current and future work only when reconciled by this ADR; no missing external review is to be reported as completed.

## Alternatives Considered

- Keep independent review mandatory: rejected because the project has one maintainer and the requirement blocks owner-directed progress without adding an available independent reviewer.
- Remove all human review and evidence: rejected because owner review, reproducible verification, and risk records are necessary controls for security-sensitive behavior.
- Make independent review mandatory only for high-risk changes: rejected because it preserves the same unavailable-person blocker; high-risk assurance instead uses owner review plus specific automated and manual evidence.

## Consequences

### Positive

- The owner can complete planning, implementation, and release decisions without waiting for an unavailable second reviewer.
- Responsibility is explicit: one person owns contract acceptance and residual-risk decisions.
- Automated validation, security scenarios, compatibility evidence, and release evidence remain enforceable.

### Negative

- There is no required independent second set of eyes; the owner may miss design or security defects.
- The owner bears all review workload and must document accepted residual risk.

## Compatibility and Migration

This governance decision does not change DepRail's CLI, API, schema, storage, or release artifact contracts. It applies to active project guidance and v0.5.0 planning/execution documents after they reference this ADR. Historical plans, accepted ADRs, completed PRs, and retrospectives remain preserved; their recorded review history is not rewritten. Current open v0.5 planning records are updated to remove external-review blockers and to show the owner as the required reviewer.

Repository branch rules and release controls are not broadened or weakened by this ADR. The live main ruleset requires CI and a pull request but zero approving reviews; the release environment reviewer is the owner. Neither currently requires another person's approval.

## Validation

- Owner decision: this ADR is accepted directly by the project owner.
- Reconcile active v0.5 plan, execution package, issue contracts, repository guidance, and release checklist to this policy.
- Verify no active v0.5 Definition of Ready or release gate requires a second person; verify security, compatibility, CI, and evidence requirements remain.
- Record the policy decision on the linked v0.5 readiness issue and preserve the historical PR/review facts.
