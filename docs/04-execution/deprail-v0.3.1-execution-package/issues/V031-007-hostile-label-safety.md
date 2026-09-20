# V031-007 Harden Hostile-Label Rendering

## Planning metadata

- Type: feature
- Area: cli
- Priority: P1
- Epic: EPIC-003
- Dependencies: V031-003

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
## Definition of Ready

- [ ] Owner and named security reviewer are assigned.
- [ ] v0.3.1 target, sprint, dependencies, and contract mapping are confirmed.
- [ ] Inputs, outputs, failure behavior, and required evidence are reviewed.
- [ ] Implementation remains blocked until release-level Definition of Ready approval.

## Contract mapping

- Requirements: `FUNCTIONAL-REQUIREMENTS.md` UX-008.
- Security: hostile-label and redaction requirements in `ERROR-MODEL.md`.
- Evidence: adversarial terminal, redaction, and shell-metacharacter cases.

## Inputs, outputs, and failure behavior

- Inputs: repository-derived labels, paths, scanner output, and error details.
- Outputs: escaped, redacted, terminal-safe labels.
- Failure: unsafe or unrepresentable values are safely escaped or rejected; they never become terminal control data.

## Final acceptance

- [ ] Owner reviewed every criterion.
- [ ] Security review completed.
- [ ] CI and required evidence links are recorded.
- [ ] Remaining risk and follow-up are documented.
