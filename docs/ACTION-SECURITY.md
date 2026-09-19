# GitHub Action Security Boundaries

## Minimum permissions

The caller workflow should declare:

```yaml
permissions:
  contents: read
```

DepRail does not require write access, pull-request write access, issues access, or secrets. Fork pull requests run with the caller's read-only token and must not depend on secrets.

## Installation and versioning

Install `deprail` before invoking the composite Action and pin the executable to an immutable release or commit. The Action verifies that the installed binary reports the requested version and fails on mismatch.

## Cache and artifacts

The Action does not use a package or scanner cache. Reports are written under the runner temporary directory, not the checkout, and uploaded only as the named report artifact. Callers should set repository-approved retention and avoid uploading source files or credentials.

Cache keys, if introduced by a caller, must include repository identity, workflow/configuration inputs, DepRail version, and scanner database identity. Cache misses are safe; untrusted or stale cache data must never be treated as a complete scan.

## Failure behavior

Missing binaries, version mismatches, invalid paths, and invalid inputs fail explicitly. Scan exit status is preserved even when the report artifact upload runs after a partial or failed scan.

## Data handling

DepRail does not upload source or secrets by default. Avoid placing tokens, credential-bearing URLs, or sensitive environment values in command arguments or report paths. Diagnostics remain on stderr and machine reports remain on stdout/artifact files.
