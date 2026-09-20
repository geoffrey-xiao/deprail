# v0.3 Manual Test Guide

This is the release-gate procedure for the v0.3 read-only remediation planning workflow. Run it from the repository root with the current built binary. Keep generated artifacts under the ignored `local_test/` directory.

## Scope and invariant

The test proves that a real scan finding can flow into a read-only remediation plan:

```text
repository -> discovery -> OSV-Scanner -> scan JSON -> fix plan JSON
```

The workflow must not modify manifests, lockfiles, source files, or package-manager state. v0.3 does not install packages, execute package-manager commands, run tests/builds/rescans, apply plans, or create patches.

## Prerequisites

- Go toolchain available.
- `jq` available.
- OSV-Scanner installed separately and available on `PATH`.
- Checked-in fixture `testdata/fixtures/mixed-repository` present.
- Run from the repository root.

Check the scanner:

```bash
osv-scanner --version
```

## Build the binary

```bash
go build -o deprail ./cmd/deprail
./deprail --version
```

Expected development-build shape:

```text
deprail development (tag=unknown commit=unknown)
```

Record the platform, Go version, binary version, and commit in the release evidence.

## Run doctor

```bash
./deprail doctor --format json > local_test/doctor-manual.json
printf 'doctor_exit=%s\n' "$?"
jq '{version, tag, commit}' local_test/doctor-manual.json
```

Expected exit code: `0`.

## Run the real mixed-repository scan

Use a new output filename because output publication is atomic and existing files are not overwritten:

```bash
./deprail scan testdata/fixtures/mixed-repository \
  --format json \
  --output local_test/mixed-repo-manual.json
printf 'scan_exit=%s\n' "$?"
```

Expected exit code: `0`.

Inspect scan completeness and finding count:

```bash
jq '{schema_version, document_type, scan_id, status, findings: (.findings | length), errors}' \
  local_test/mixed-repo-manual.json
```

Expected result for the current mixed fixture:

```json
{
  "schema_version": "v1alpha",
  "document_type": "scan",
  "status": "complete",
  "findings": 14,
  "errors": []
}
```

List findings using the current scan JSON fields:

```bash
jq -r '.findings[] | [.TargetID, .Component, .Version, .Fixed] | @tsv' \
  local_test/mixed-repo-manual.json
```

Select a finding:

```bash
FINDING_KEY="$(jq -r '.findings[0].TargetID' local_test/mixed-repo-manual.json)"
printf '%s\n' "$FINDING_KEY" | tee local_test/finding-key.txt
```

Do not use `.stable_key`, `.component.purl`, or `.current_version` against this raw scan document. Those fields belong to the normalized v0.3 planning-report contract.

## Verify enriched scan identity

The scan-to-plan integration must preserve planner identity:

```bash
jq -r '.findings[0] | {
  TargetID,
  Component,
  Version,
  purl,
  workspace_id,
  workspace_path,
  ecosystem
}' local_test/mixed-repo-manual.json
```

Expected shape for the first mixed-fixture finding:

```json
{
  "TargetID": "GHSA-29mw-wpgm-hmr9",
  "Component": "lodash",
  "Version": "4.17.20",
  "purl": "pkg:npm/lodash@4.17.20",
  "workspace_id": "frontend",
  "workspace_path": "frontend",
  "ecosystem": "npm"
}
```

The exact vulnerability key may change if the fixture changes; the PURL and workspace assertions must remain true for the selected npm finding.

## Generate the real remediation plan

```bash
./deprail fix plan \
  --report local_test/mixed-repo-manual.json \
  --finding "$FINDING_KEY" \
  --format json \
  > local_test/fix-plan-fixed.json
printf 'plan_exit=%s\n' "$?"
```

Expected exit code: `0`.

Inspect the generated plan:

```bash
jq '{
  plan_id,
  schema_version,
  component,
  workspace_identity,
  current_state,
  candidates,
  affected_files,
  risks,
  verification,
  provenance
}' local_test/fix-plan-fixed.json
```

Expected minimum evidence:

- Component PURL: `pkg:npm/lodash@4.17.20`.
- Workspace path: `frontend`.
- Manifest: `frontend/package.json`.
- Lockfile: `frontend/package-lock.json`.
- At least one candidate or an explicit no-recommendation explanation.
- Stable plan ID.
- Source scan ID and repository state.
- Verification requirements.
- No repository mutation.

The previous failure is a regression signal, not an acceptable release result:

```text
error[PLAN_UNSUPPORTED] no planning adapter supports the finding component (finding)
```

That error means the scan did not preserve the ecosystem PURL or planner identity.

## Safe output checks

External plan output must be outside the scanned repository root. `local_test/` is outside `testdata/fixtures/mixed-repository`, so it is suitable for this test:

```bash
./deprail fix plan \
  --report local_test/mixed-repo-manual.json \
  --finding "$FINDING_KEY" \
  --output local_test/fix-plan-external.json
printf 'external_plan_exit=%s\n' "$?"
jq '{plan_id, schema_version, component}' local_test/fix-plan-external.json
```

A second write to the same path must fail with `PLAN_WRITE_FAILED` and must not overwrite the existing plan.

Repository-internal and traversal outputs must fail with `PLAN_OUTPUT_OUTSIDE_ROOT_REQUIRED`:

```bash
./deprail fix plan \
  --report local_test/mixed-repo-manual.json \
  --finding "$FINDING_KEY" \
  --output testdata/fixtures/mixed-repository/plan.json

./deprail fix plan \
  --report local_test/mixed-repo-manual.json \
  --finding "$FINDING_KEY" \
  --output testdata/fixtures/mixed-repository/child/../plan.json
```

## Prove no mutation

Capture the repository state before and after the plan commands:

```bash
git status --short > local_test/tree-before.txt
git diff -- testdata/fixtures/mixed-repository > local_test/tree-diff-before.patch

# Run the scan and plan commands above.

git status --short > local_test/tree-after.txt
git diff -- testdata/fixtures/mixed-repository > local_test/tree-diff-after.patch
diff -u local_test/tree-diff-before.patch local_test/tree-diff-after.patch
```

Expected result: no diff caused by `fix plan`; manifests, lockfiles, and source files remain unchanged.

## Automated verification

```bash
go test ./...
go vet ./...
git diff --check
```

Record the exact command output, platform, binary identity, generated artifact paths, and any remaining limitation in the release-evidence record.

## Release evidence checklist

- [ ] Real binary built from the reviewed commit.
- [ ] `doctor` JSON captured.
- [ ] Complete mixed-repository scan captured.
- [ ] Scan contains vulnerabilities and enriched PURL/workspace identity.
- [ ] `fix plan` exits `0` and produces valid JSON.
- [ ] Plan contains candidate, affected-file, risk, verification, and provenance evidence.
- [ ] External output succeeds without overwriting.
- [ ] Internal/traversal output is rejected.
- [ ] Repository tree and lockfiles remain unchanged.
- [ ] `go test ./...` passes.
- [ ] `go vet ./...` passes.
- [ ] Linux, macOS, and Windows CI passes.
- [ ] Owner and security/architecture reviewer record the release decision.
