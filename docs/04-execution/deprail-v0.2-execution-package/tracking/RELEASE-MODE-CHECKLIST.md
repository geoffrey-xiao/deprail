# v0.2 Release-Mode Evidence Checklist

This checklist is the release decision record for the v0.2 preview, release-candidate, and stable modes. A missing item blocks the mode; an unchecked claim is not evidence.

## Common release record

- [ ] Source commit is recorded and matches the build input.
- [ ] Release tag is recorded and immutable after publication.
- [ ] DepRail version, tag, and commit are present in `deprail doctor --format json`.
- [ ] Build commands and toolchain versions are recorded.
- [ ] CI run links cover Linux, macOS, and Windows.
- [ ] Artifact names, sizes, and SHA-256 checksums are recorded.
- [ ] OSV-Scanner version and compatibility result are recorded.
- [ ] Representative JavaScript, Python, Java, and mixed-repository scans are recorded.
- [ ] Scan status is `complete` for every release-gating scenario.
- [ ] Known risks, limitations, and rollback owner are recorded.

## Preview mode

### Entry criteria

- [ ] v0.2 scope and deferred capabilities are confirmed.
- [ ] Required CLI, schema, adapter, and output contracts are present on the tagged commit.
- [ ] Focused tests and three-platform CI pass.
- [ ] At least one mixed-repository scan produces schema-valid deterministic output.
- [ ] SBOM status is explicitly recorded as unavailable, planned, or supplied.
- [ ] Signing and provenance status is explicitly recorded as unavailable, planned, or supplied.

### Exit criteria

- [ ] Preview artifacts and checksums are published or attached to the release record.
- [ ] Manual scans cover representative repositories and retain command/output evidence.
- [ ] Known limitations and the signing/SBOM gap are visible in release notes.
- [ ] A reviewer records the decision to continue, revise, or stop.

## Release-candidate mode

### Entry criteria

- [ ] Preview findings and release defects are dispositioned.
- [ ] Contract and compatibility changes are frozen for the candidate.
- [ ] Linux, macOS, and Windows artifact smoke tests pass.
- [ ] amd64 and arm64 artifact evidence is present where supported.
- [ ] Checksums match downloaded artifacts and the recorded source commit.
- [ ] Representative repository scans are repeatable and complete.
- [ ] Rollback artifacts and immutable-tag procedure are rehearsed.

### Exit criteria

- [ ] No open release-blocking defect remains.
- [ ] Candidate artifact inventory and checksums are complete.
- [ ] SBOM status and signing/provenance status are explicit.
- [ ] Manual scan evidence and CI links are attached.
- [ ] The release owner records the go/no-go decision.

## Stable mode

### Entry criteria

- [ ] Candidate acceptance evidence is complete and reviewed.
- [ ] Support, compatibility, and migration notes are present.
- [ ] Release artifacts are reproducible from the recorded commit and tag.
- [ ] Checksums, platform smoke results, and representative scans are complete.
- [ ] SBOM and signing/provenance evidence is supplied, or the release is explicitly rejected as stable.
- [ ] Immutable stable tag and rollback procedure are confirmed.

### Exit criteria

- [ ] Stable artifacts are published under the approved immutable tag.
- [ ] Checksums and release evidence are publicly or internally available at the documented location.
- [ ] Installation smoke tests pass on supported platforms.
- [ ] Release notes identify known risks, scanner version, and support boundaries.
- [ ] The release owner records the final go/no-go decision.

## Rollback and immutable tags

- Never move or reuse a published release tag.
- [ ] Identify the last known-good tag and artifact checksums.
- [ ] Disable or withdraw the affected release channel without deleting evidence.
- [ ] Publish a corrective release with a new immutable tag.
- [ ] Record the failed commit, affected artifacts, user impact, and recovery steps.
- [ ] Re-run the applicable mode checklist before restoring the release channel.

## Evidence record

Attach the completed checklist with:

- source commit and tag;
- CI run URLs and platform results;
- artifact inventory and SHA-256 checksums;
- scanner version and representative scan commands;
- SBOM status;
- signing/provenance status;
- known risks, rollback record, and reviewer decision.

See [`EVIDENCE-GUIDE.md`](EVIDENCE-GUIDE.md) for the minimum evidence standard.
