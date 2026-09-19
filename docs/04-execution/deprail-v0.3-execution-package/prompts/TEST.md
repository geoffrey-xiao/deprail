# v0.3 Test Prompt

Test the observable remediation-plan contract, not implementation wiring.

Required coverage:

- direct patch and transitive-owner plans;
- multiple candidates and deterministic recommendation;
- major, peer, runtime, engine, and lockfile risks;
- rejected, unavailable, unknown, stale, malformed, unsupported, and incomplete inputs;
- schema-valid JSON, stable ordering, stdout/stderr separation, and atomic output;
- traversal, symlink escape, shell metacharacters, malicious manifests, oversized metadata, credential URLs, and mutation attempts;
- unchanged repository tree and lockfiles after every plan;
- equivalent Linux/macOS/Windows semantics.

Retain command output, golden artifacts, schema results, and CI links. Never weaken an assertion to fit an implementation.
