# TEST-010 Test Mixed Monorepo Discovery Across Platforms

**Risk:** R2  
**Status:** Ready

## Goal

Validate JavaScript, Python, and Java boundaries, conflicts, paths, and completeness on all supported operating systems.

## Scope and Constraints

- Follow the v0.1 CLI, data, error, compatibility, and test contracts.
- Keep domain behavior deterministic and independently testable.
- Preserve evidence and return explicit completeness and stable errors.
- Do not add a shell, implicit network access, automatic installation, or unrelated roadmap work.
- Keep platform differences at transport boundaries; serialized domain meaning must remain stable.

## Inputs Outputs and Failure Behavior

Document accepted inputs, validation, output schema or interface, sort order, side effects, and every expected failure. A missing or malformed dependency must never become an empty successful result.

## Required Tests

- Happy path with the smallest representative fixture.
- Invalid and incomplete input.
- Deterministic repeat and input-order permutation where applicable.
- Linux, macOS, and Windows semantic behavior where paths or processes are involved.
- Security cases for hostile paths, arguments, output size, cancellation, or secret redaction as applicable.

## Acceptance Checklist

- [ ] The goal is demonstrable through the public or contracted interface.
- [ ] Unit and contract tests cover success and failure behavior.
- [ ] Output validates against the current schema when serialized.
- [ ] Diagnostics use stable codes and do not disclose secrets.
- [ ] Documentation and examples match implemented behavior.
- [ ] `make verify` passes with no unexplained generated changes.
- [ ] Required evidence is attached to the pull request.

## Human Review

Review contract boundaries, failure completeness, security assumptions, cross-platform behavior, test assertions, and any golden-file change. Run or inspect at least one manual representative case.

## Evidence

Record the commit, commands, test results, sample output or fixture, supported tool versions, remaining limitations, and reviewer decision.
