# DepRail v0.2 Release Procedure

**Status:** Planning procedure; not a release declaration
**Scope:** v0.2.0 preview, RC, and stable release modes

## 1. Release identity

Before preparing a release, record:

- latest stable tag;
- latest preview and RC tags;
- current `origin/main` commit;
- `.release-please-manifest.json` value;
- intended release mode and version;
- release owner and reviewer.

Never reuse an existing tag. Keep preview, RC, and stable identities distinct:

```text
v0.2.0-preview.1 → v0.2.0-rc.1 → v0.2.0
```

Bug fixes use the next patch version. Do not create arbitrary fix tags.

## 2. Release modes

### Preview

Used for controlled validation. It may carry explicitly approved SBOM/signing gaps and does not imply production acceptance.

### Release candidate

Used when the release gate is believed complete. All known P0/P1 release defects require resolution or explicit owner disposition.

### Stable

Requires complete Master Checklist evidence, reviewed artifacts, platform smoke results, representative repository scans, and an explicit owner release decision.

## 3. Pre-release gate

Run the following from a clean checkout of the reviewed commit:

```bash
git fetch origin main --tags
git checkout --detach <reviewed-commit>
go version
make verify
git diff --check
```

Record the commit, Go version, commands, exit codes, and CI run URLs in the release evidence record. Any failure blocks the release mode.

The release mode determines the required decision:

| Mode | Required decision |
| --- | --- |
| Preview | Continue only with explicit owner acceptance of documented gaps. |
| RC | Proceed only after preview defects are dispositioned and artifact identity/checksums are verified. |
| Stable | Proceed only with complete evidence and explicit owner go/no-go approval. |

- [ ] v0.2 Master Checklist is current.
- [ ] Target version is consistent across manifest, tag plan, CLI identity, and release record.
- [ ] Required PRs are merged with owner review recorded.
- [ ] `make verify` passes from a clean checkout.
- [ ] Pull-request CI passes on Linux, macOS, and Windows.
- [ ] Required scanner version and compatibility evidence are recorded.
- [ ] Representative repository scan evidence is attached.
- [ ] Remaining risk is explicit.

## 4. Artifact verification

Build only from the reviewed release commit. The required matrix remains:

```text
linux/amd64
darwin/amd64
darwin/arm64
windows/amd64
```

For each artifact:

- run `doctor --format json`;
- run discovery against the fixed mixed-repository fixture;
- run scanner smoke only where the supported scanner is available;
- confirm reported version/tag/commit identity;
- generate and review SHA-256 checksums;
- record output and environment without secrets.

## 5. Supply-chain status

- SBOM generation must be explicitly marked implemented or unavailable.
- Signing and provenance must be explicitly marked implemented or unavailable.
- A preview may carry approved gaps; a stable release must satisfy the current release policy or record an owner-approved exception.
- Never claim a control based solely on workflow configuration if repository settings or generated evidence do not prove it.

## 6. Publication sequence

1. Confirm release baseline and mode.
2. Confirm Master Checklist evidence.
3. Build artifacts from the reviewed commit.
4. Run verification and smoke checks.
5. Review checksums and release notes.
6. Create a new immutable tag.
7. Publish the GitHub Release only when the release mode permits publication.
8. Run protected approval gates before artifact upload.
9. Attach artifacts, checksums, evidence, and remaining risks.
10. Record the final owner decision.
11. Create and review the version-specific retrospective.
12. Link follow-up actions and remaining risks from the release issue.

The retrospective is required after preview, release-candidate, and stable publication. It does not replace release evidence or owner/security approval. Use the reusable [`docs/RELEASE-CHECKLIST.md`](RELEASE-CHECKLIST.md) for the complete pre-publication and post-publication gate.

## 7. Rollback

If artifact identity, checksums, smoke tests, or release evidence are incorrect:

- stop publication;
- do not move or reuse the tag;
- preserve the failed evidence;
- create the next corrected preview/RC/patch version;
- record the cause and decision in the release record.

## 8. Release evidence record

Use a version-specific evidence record, such as [`docs/RELEASE-v0.2-preview.2-EVIDENCE.md`](RELEASE-v0.2-preview.2-EVIDENCE.md). Preserve superseded preview records unchanged and never mix artifact checksums across tags.

The release record must link:

- source commit and tag;
- CI matrix;
- artifact names and SHA-256 values;
- platform smoke output;
- scanner version and adapter evidence;
- representative repository scans;
- SBOM/signing status;
- reviewer and owner decision;
- remaining risks and rollback notes.
## 9. Post-release retrospective

Create `docs/retrospectives/RETROSPECTIVE-v<version>.md` after every published release. The record MUST include:

- delivered capabilities and release outcomes;
- what went well;
- mistakes, impact, root cause, correction, and prevention;
- security and supply-chain lessons;
- process improvements;
- remaining risks;
- follow-up actions with owners and acceptance evidence;
- separate owner and security/architecture review status.

Link the retrospective from the release issue and version-specific release evidence. Follow-up actions MUST be tracked as issues or explicitly accepted before the retrospective is marked complete. Preserve historical retrospectives unchanged; corrections require an additive follow-up record.
