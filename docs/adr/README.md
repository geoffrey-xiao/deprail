# Architecture Decision Records

Record accepted or proposed architecture and contract decisions in this directory using the template at `docs/04-execution/deprail-v0.1-execution-package/templates/ADR-TEMPLATE.md`. The solo-owner review policy is [ADR-0005](ADR-0005-solo-owner-review-policy.md).

## Workflow

1. Open or use the linked GitHub Issue for the decision.
2. State the context, owner decision, alternatives, compatibility impact, security impact, optional reviewers, and evidence.
3. The owner decides before changing schemas, public interfaces, tool versions, storage, permissions, or security boundaries; external review is optional.
4. Create the ADR with the next numeric ID and link the GitHub Issue and related pull request.
5. Set the ADR status to `Accepted` after the owner records the decision.
6. Update affected contracts and implementation documentation in the same change.

ADR filenames use `ADR-NNNN-short-description.md`. Do not silently rewrite accepted decisions; supersede them with a new ADR that records migration and consequences.
