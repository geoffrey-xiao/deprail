# v0.3 Implementation Prompt

Implement only a reviewed v0.3 issue mapped to the execution package.

Before editing:

- read the issue, requirement row, schema, architecture boundary, and acceptance evidence;
- preserve v0.2 contracts;
- confirm the change cannot mutate the repository or execute package-manager scripts;
- identify tests and rollback evidence.

Implementation must use deterministic domain services, structured commands, bounded parsing, canonical path containment, explicit unknown/incomplete states, stable errors, schema validation, and provenance. Do not add `fix apply`, worktrees, package installation, network publishing, PR creation, web services, MCP writes, or unrelated scanner support.

Deliver code, tests, documentation, and evidence together. Report exact commands run and remaining risks.
