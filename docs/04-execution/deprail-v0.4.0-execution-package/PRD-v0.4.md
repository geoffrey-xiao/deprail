# DepRail v0.4 PRD

**Status:** Planning draft; owner and architecture/security review required
**Release:** `v0.4.0-preview.1`
**Roadmap outcome:** Remediation application and verification

## User outcome

A developer can create a trusted baseline from a complete scan, review an approved remediation plan, apply it in an isolated workspace, run bounded verification, rescan the result, and receive patch evidence without changing the caller's repository or hiding failure.

## Included

- `deprail baseline create` conversion from a complete scan result to the existing validated baseline contract.
- Explicit rejection of incomplete, failed, malformed, stale, or unsupported scan inputs.
- `deprail fix apply` plan validation and dry-run.
- Explicit approval bound to plan and source identity.
- Isolated Git worktree lifecycle.
- Allowlisted dependency mutation through supported package-manager adapters.
- Bounded tests, builds, and type checks.
- Before/after rescan and finding transition classification.
- Atomic rollback/discard and durable redaction-safe evidence.

## Excluded

No caller-worktree mutation, automatic commit/push/PR/merge, autonomous approval, arbitrary shell commands, unbounded scripts/network, new ecosystems, web/team services, MCP write tools, automatic baseline replacement, remote baseline publishing, policy-language expansion, or stable-release claim.

## Success gate

Failed changes remain contained and recoverable; the source tree is unchanged; evidence identifies every command, path, tool, outcome, and rollback state; serialized meaning is equivalent across supported platforms.
