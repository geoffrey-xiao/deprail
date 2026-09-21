# V04-003: Execute Bounded Supported Mutation

**Epic:** EPIC-002
**Status:** Proposed; implementation blocked

## Scope

Define package-manager adapter ports and direct-argv execution with explicit cwd, approved environment, deadlines, output caps, cancellation, process-tree termination, and script/network policy.

## Acceptance

Only supported ecosystems and approved operations execute. Missing tools, incompatible versions, non-zero exits, timeout, cancellation, output limits, scripts, network, and shell metacharacters produce explicit failures without caller mutation.
