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

- Issue: Pending creation
- Parent epic: Pending creation
