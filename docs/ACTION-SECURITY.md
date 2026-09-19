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

Cache keys must include all inputs that can change scan meaning:

- repository identity and immutable commit/tree SHA;
- dependency manifest and lockfile content digest;
- DepRail version/tag and commit;
- scanner/database identity;
- effective repository and DepRail configuration digest.

A cache miss is safe. A cache entry from another commit, dependency-input digest, tool/database identity, or configuration must not be restored as a trusted result. The Action itself does not use a package or scanner cache; these requirements apply to caller workflows.

Missing binaries, version mismatches, invalid paths, and invalid inputs fail explicitly. Scan exit status is preserved even when the report artifact upload runs after a partial or failed scan.

## Data handling

DepRail does not upload source or secrets by default. Avoid placing tokens, credential-bearing URLs, or sensitive environment values in command arguments or report paths. Diagnostics remain on stderr and machine reports remain on stdout/artifact files.
