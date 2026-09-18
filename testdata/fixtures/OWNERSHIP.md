# Fixture ownership and updates

Each fixture has one intended detection purpose documented in `README.md`. The contributor changing a fixture owns:

- explaining the input and expected detection;
- keeping it offline, minimal, deterministic, and license-safe;
- updating contract or golden tests when behavior changes;
- recording tool/version provenance for raw scanner fixtures;
- attaching verification output to the pull request.

Do not refresh fixture or golden output broadly without explaining each semantic change and obtaining human review.
