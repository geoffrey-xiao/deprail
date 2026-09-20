# V031-007 Harden Hostile-Label Rendering

## Planning metadata

- Type: security
- Area: cli
- Priority: P1
- Risk: R2
- Epic: EPIC-003
- Dependencies: V031-003, V031-005

## Goal
Ensure repository, scanner, path, package, and error labels remain data when rendered to terminals or redirected diagnostics.

## Acceptance criteria

- [ ] Newline, carriage return, tab, Unicode, long values, ANSI, and terminal controls are handled safely.
- [ ] TTY and redirected stderr use the same escaping and redaction policy.
- [ ] Credential-bearing URLs and sensitive environment values are not exposed.
- [ ] Escaping preserves enough safe context for diagnosis.
- [ ] Shell metacharacters remain data and are never executed.

## Evidence
Adversarial presenter tests, redaction cases, terminal capture review, and security review.

## Exclusions
No broad logging rewrite, scanner changes, network upload, or repository mutation.

## GitHub tracking
- Issue: [#250](https://github.com/geoffrey-xiao/deprail/issues/250)
- Parent epic: [#243](https://github.com/geoffrey-xiao/deprail/issues/243)
