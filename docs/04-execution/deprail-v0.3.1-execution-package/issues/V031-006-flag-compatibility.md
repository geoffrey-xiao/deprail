# V031-006 Implement Quiet and Preserve Verbose Behavior

## Planning metadata

- Type: feature
- Area: cli
- Priority: P1
- Risk: R2
- Epic: EPIC-003
- Dependencies: V031-002, V031-003, V031-005

## Goal
Close the current `deprail scan --quiet` parser gap and preserve safe `--verbose` behavior across supported command output modes.

## Acceptance criteria

- [ ] `deprail scan --quiet` is accepted by the parser.
- [ ] Quiet mode suppresses successful progress and summaries while retaining required errors and diagnostics.
- [ ] Verbose mode preserves existing behavior and adds only safe diagnostics.
- [ ] Quiet and verbose work with JSON without contaminating stdout.
- [ ] `doctor`, `discover`, `scan`, `diff`, and `fix plan` behavior is documented or explicitly scoped.
- [ ] Parser, stream, and CLI smoke tests cover precedence and invalid combinations.

## Evidence
CLI contract update, parser tests, stream captures, and quiet/verbose smoke matrix.

## Exclusions
No unrelated flag redesign, scanner changes, or machine-schema break.
## GitHub tracking


- Issue: [#257](https://github.com/geoffrey-xiao/deprail/issues/257)
- Parent epic: [#243](https://github.com/geoffrey-xiao/deprail/issues/243)
## Definition of Ready

- [ ] Owner and named reviewer are assigned.
- [ ] v0.3.1 target, sprint, dependencies, and contract mapping are confirmed.
- [ ] Inputs, outputs, failure behavior, and required evidence are reviewed.
- [ ] Implementation remains blocked until release-level Definition of Ready approval.

## Contract mapping

- Requirements: `FUNCTIONAL-REQUIREMENTS.md` UX-006.
- CLI contract: `CLI-CONTRACT.md` quiet and verbose behavior.
- Evidence: parser, stream, precedence, and smoke cases.

## Inputs, outputs, and failure behavior

- Inputs: command arguments, output format, and stream mode.
- Outputs: documented quiet/verbose behavior with unchanged machine semantics.
- Failure: invalid combinations remain configuration errors; quiet parser support cannot convert scan failures into success.

## Final acceptance

- [ ] Owner reviewed every criterion.
- [ ] Architecture/security review completed.
- [ ] CI and required evidence links are recorded.
- [ ] Remaining risk and follow-up are documented.
