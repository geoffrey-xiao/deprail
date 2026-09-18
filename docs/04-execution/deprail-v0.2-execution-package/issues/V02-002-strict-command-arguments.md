# V02-002 Reject Unexpected Command Arguments

- Type: bug
- Area: cli
- Priority: P0
- Risk: R1
- Target version: 0.2.0
- Milestone: v0.2.0
- Sprint: Sprint 0
- Owner: TBD
- Reviewer: TBD
- Dependencies: None
- Status: Local planning; GitHub issue not created

## Definition of Ready

- [ ] Accepted arguments for every command are listed.
- [ ] Invalid argument exit behavior is identified.
- [ ] No-op or silently ignored paths are covered.
- [ ] Owner and reviewer are assigned.

## Goal

Prevent commands from silently ignoring unexpected positional arguments or options.

## Scope

Audit `doctor`, `discover`, and `scan` argument parsing. Reject unexpected input before discovery, scanner execution, writes, or network access.

## Out of Scope

Adding new CLI options, changing valid command syntax, or redesigning the command framework.

## Inputs, Outputs, and Failure Behavior

Invalid arguments produce a stable configuration error on stderr and exit code `2`. No partial command executes and no successful report or artifact is emitted.

## Required Tests

- Extra positional argument for `doctor`.
- Unsupported option for each command.
- Valid command invocation regression tests.
- Machine-output and stderr separation.

## Acceptance Criteria

- [ ] Unexpected arguments fail deterministically.
- [ ] Exit code is `2`.
- [ ] No command work occurs after argument failure.
- [ ] Valid v0.1 command forms remain accepted.

## Human Review

Review CLI compatibility and parser error ownership.

## Evidence Required

Command transcripts, exit codes, focused tests, and compatibility notes.

## Final Acceptance

- [ ] Owner reviewed every criterion during PR review.
- [ ] CI and evidence were reviewed.
- [ ] Remaining risk is recorded.
- [ ] PR review and merge evidence are linked.
