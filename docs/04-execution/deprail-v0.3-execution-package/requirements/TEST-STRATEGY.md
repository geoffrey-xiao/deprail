# v0.3 Test Strategy

## Unit

Candidate ordering, constraints, direct/transitive ownership, major-version risk, stable plan identity, canonicalization, stale-input checks, and unknown/unavailable state handling.

## Contract

Each package-manager adapter maps to the same domain plan contract; structured commands never require shell interpolation; plans validate against the schema; failures preserve stable diagnostics.

## Golden

Direct patch, transitive owner, multiple candidates, major migration, peer/runtime incompatibility, no candidate, missing metadata, malformed manifest/lockfile, stale input, and mixed repository.

## Security

Traversal, symlink escape, shell metacharacters, malicious manifests, oversized metadata, credential URLs, package scripts, network-disabled planning, restrictive writes, mutation attempts, and reordered evidence.

## Integration and smoke

Run `fix plan` against fixed npm, Python, Java, and mixed fixtures on Linux, macOS, and Windows. Assert repository tree and lockfiles are unchanged. Repeat with reordered source evidence and compare normalized JSON.

## Acceptance evidence

Retain command output, schema validation, golden diffs, mutation checks, platform CI links, and representative plan artifacts. Do not accept a test that only checks object construction or implementation details.
