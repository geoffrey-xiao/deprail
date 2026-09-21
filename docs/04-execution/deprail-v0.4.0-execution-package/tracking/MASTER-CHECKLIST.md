# v0.4.0 Master Checklist

**Release:** `v0.4.0-preview.1`
**Parent:** [#299](https://github.com/geoffrey-xiao/deprail/issues/299)
**Readiness:** [#301](https://github.com/geoffrey-xiao/deprail/issues/301)
**Isolation ADR:** [#303](https://github.com/geoffrey-xiao/deprail/issues/303)

## Planning and review

- [x] v0.4 development plan exists.
- [x] Product, architecture, roadmap, and v0.3 evidence are reconciled.
- [x] Execution package exists.
- [x] Definition of Ready is recorded and owner-approved.
- [x] Architecture/security reviewer approval is recorded for ADR-0002.
- [x] Mutation boundary is accepted separately from implementation approval.

## Isolation and rollback gate

- [x] ADR-0002 is accepted.
- [x] ISO-001 through ISO-011 are reviewed.
- [ ] Hostile path and symlink fixtures exist.
- [ ] Source-tree immutability fixture exists.
- [ ] Interrupted atomic-write fixture exists.
- [ ] Timeout and cancellation cleanup fixtures exist.
- [ ] Repeated cleanup is proven idempotent.
- [ ] Linux/macOS/Windows behavior is recorded.

## Implementation gate

- [ ] #302 has an implementation branch and contract tests.
- [ ] No package-manager mutation is included in the isolation slice.
- [ ] No caller worktree mutation is possible.
- [ ] No commit, push, pull request, merge, or publication behavior is included.
- [ ] Architecture/security review is recorded before merge.

## Evidence gate

- [ ] Workspace identity and source commit are captured.
- [ ] Plan digest and approval binding are captured.
- [ ] Before/after trees and digests are captured.
- [ ] Cleanup and rollback outcomes are explicit.
- [ ] Credentials and unrestricted environment values are absent.
- [ ] Partial and failed outcomes are distinguishable from success.
