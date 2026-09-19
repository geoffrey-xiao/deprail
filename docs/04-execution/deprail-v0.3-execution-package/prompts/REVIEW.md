# v0.3 Review Prompt

Review the change against the linked v0.3 requirement and evidence row.

Check:

- no repository mutation, package installation, script execution, or unapproved network access;
- finding and dependency ownership correctness;
- candidate states and unknown evidence handling;
- direct/transitive, constraint, lockfile, peer/runtime, and major-version risks;
- stable plan identity and order independence;
- schema, CLI, stdout/stderr, exit, path, and provenance contracts;
- hostile input, redaction, bounds, and cross-platform behavior;
- v0.4 handoff clarity.

Reject false-safe recommendations, silent fallbacks, shell interpolation, implementation-only tests, and undocumented contract changes.
