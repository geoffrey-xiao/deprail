# EPIC-002: Approval and Controlled Mutation

**Status:** Proposed

## Outcome

Only an explicit, unexpired approval for an unchanged plan and source identity can authorize supported dependency mutations.

## Boundaries

Includes plan/approval binding, dry-run, adapter selection, direct-argv package-manager execution, environment/network/script policy, deadlines, cancellation, and output limits. Excludes arbitrary commands, new ecosystems, and autonomous approval.

## Acceptance

Approval mismatch, expiry, reuse, path changes, unsupported tools, non-zero exits, timeouts, and output limits fail explicitly without caller-worktree mutation.
