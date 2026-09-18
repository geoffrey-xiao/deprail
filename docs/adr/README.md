# Architecture Decision Records

Record accepted or proposed architecture and contract decisions in this directory using the template at `docs/04-execution/deprail-v0.1-execution-package/templates/ADR-TEMPLATE.md`.

## Workflow

1. Open a GitHub Issue using the architecture decision template.
2. State the context, decision required, alternatives, compatibility impact, security impact, reviewers, and evidence.
3. Discuss the decision before changing schemas, public interfaces, tool versions, storage, permissions, or security boundaries.
4. Create the ADR with the next numeric ID and link the GitHub Issue and related pull request.
5. Set the ADR status to `Accepted` only after human review.
6. Update affected contracts and implementation documentation in the same change.

ADR filenames use `ADR-NNNN-short-description.md`. Do not silently rewrite accepted decisions; supersede them with a new ADR that records migration and consequences.
