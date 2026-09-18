# v0.2 Test Strategy

## Required layers

- Unit tests for collection initialization, argument validation, version parsing, stable keys, and pure normalization.
- Contract tests for OSV-Scanner v2 command, output, exit codes, timeout, malformed output, and missing tool behavior.
- End-to-end tests for scanning from outside the requested root, mixed repositories, partial results, and artifact placement.
- Golden/schema tests for stable JSON arrays, deterministic ordering, and compatibility examples.
- Cross-platform smoke tests for Linux, macOS, Windows, amd64 artifacts, and arm64 smoke.
- Property/fuzz tests for ordering, parser hostility, and path containment where applicable.

## Required adversarial cases

- External symlink and traversal escape.
- Shell metacharacters in paths.
- Unicode and long paths.
- Malformed manifests and lockfiles.
- Huge scanner output.
- Timeout and cancellation.
- Non-zero scanner exits, including vulnerability-found exit.
- Credential-bearing URLs and sensitive environment values.
- Interrupted atomic writes.

## Evidence rule

A test is required when a plausible regression would otherwise escape review. Tests must assert consumer-visible behavior, boundaries, invariants, transitions, precedence, or real error behavior; they must not merely pin implementation details.
