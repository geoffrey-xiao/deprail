# V031-010 Complete v0.3.1 Release Gate

## Planning metadata

- Type: docs
- Area: docs
- Priority: P1
- Risk: R2
- Epic: EPIC-004
- Dependencies: V031-009 and all v0.3.1 acceptance evidence

## Goal
Assemble the v0.3.1 evidence bundle and obtain owner and architecture/security approval without claiming unproven stable compatibility.

## Acceptance criteria

- [ ] General release checklist is complete.
- [ ] Development plan, execution package, traceability, and compatibility matrix are synchronized.
- [ ] All required manual and automated evidence links are recorded.
- [ ] No approved stable v0.3.0 tag is treated as existing unless verified; preview evidence is labeled carryover.
- [ ] Owner approval and architecture/security review are recorded.
- [ ] Remaining risks, deferred v0.4 work, and rollback notes are explicit.
- [ ] Retrospective handoff is prepared.

## Evidence
Release evidence document, checklist, CI links, manual captures, reviewer decisions, and retrospective link.

## Exclusions
No automatic issue closure, repository mutation, publication, or v0.4 execution.
## GitHub tracking

- Issue: [#253](https://github.com/geoffrey-xiao/deprail/issues/253)
- Parent epic: [#241](https://github.com/geoffrey-xiao/deprail/issues/241)
## Definition of Ready

- [ ] Owner and named release reviewer are assigned.
- [ ] v0.3.1 target, sprint, dependencies, and contract mapping are confirmed.
- [ ] Inputs, outputs, failure behavior, and required evidence are reviewed.
- [ ] Implementation remains blocked until release-level Definition of Ready approval.

## Contract mapping

- Requirements: `MASTER-CHECKLIST.md`, `TRACEABILITY.md`, and `RELEASE-CHECKLIST.md`.
- Release boundary: no stable claim without complete evidence and owner decision.
- Evidence: linked CI, manual captures, reviewer decisions, and retrospective handoff.

## Inputs, outputs, and failure behavior

- Inputs: all v0.3.1 implementation evidence and release metadata.
- Outputs: reviewable release evidence bundle and go/no-go record.
- Failure: missing, contradictory, or preview-only evidence blocks the applicable release claim.

## Final acceptance

- [ ] Owner reviewed every criterion.
- [ ] Architecture/security review completed.
- [ ] CI and required evidence links are recorded.
- [ ] Remaining risk and follow-up are documented.
