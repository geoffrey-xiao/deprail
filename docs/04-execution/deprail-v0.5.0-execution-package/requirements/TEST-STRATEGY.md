# v0.5 Test Strategy

**Status:** Draft; exact commands/tool versions are assigned in the implementation execution package after design approval.

## 1. Principles

Test consumer-observable behavior across service, storage, API, browser, packaging, and platforms. Keep fixtures offline, deterministic, minimal, and license-safe. A passing unit suite alone does not satisfy release acceptance. Preserve CLI tests and add no UI/API framework-specific assertions where externally observable contracts suffice.

## 2. Domain and history contract tests

- Distinct runs with the same source `ScanReport.ScanID` have unique history-entry IDs and remain independently addressable; ordering is deterministic.
- Report completeness (`complete|partial|failed`) and operation outcome (`completed|failed|cancelled`) survive serialization/storage roundtrip as separate fields; cancelled operations never acquire a fabricated report status.
- Empty history differs from missing/corrupt/unavailable store.
- Finding keys, provenance, artifact digests, workspace identity, and diagnostics survive permitted data transforms.
- Duplicate ingestion/idempotency and conflict behavior once specified.
- Missing artifacts/digest mismatch are explicit and do not become empty evidence.

## 3. SQLite integration tests

- Fresh DB initialization, supported migration from each prior schema, latest-schema reopen.
- Migration interruption/failure, old DB preservation, backup/restore according to approved plan.
- Transaction rollback after mid-write failure and process-like interruption.
- Concurrent readers/writer, lock/busy timeout, cancellation, disk-full/quota simulation, permission denial, corrupt file/page, unsupported future schema.
- Retention/deletion and artifact reference accounting once defined.
- Atomic and restrictive filesystem behavior on Linux/macOS/Windows.

Use temporary directories/DBs; never require network, real user data, or host-global paths.

## 4. API contract and adversarial tests

- Validate operations and responses against reviewed OpenAPI 3.1 schema and examples.
- Verify bounded pages, stable cursors/order, invalid/unknown filters, unknown IDs, unsupported versions, malformed/oversized bodies, wrong content type/method.
- Assert HTTP status/error envelopes for store unavailable, corrupt data, artifact missing/mismatch, timeout, and read-request cancellation; separately verify persisted scan-operation cancellation does not mutate `ScanReport.Status`.
- Security cases: traversal and encoded separators, unusual Unicode, symlink/path boundary if any filesystem resolution exists, SQL metacharacters, hostile Host/Origin, DNS-rebinding patterns, forbidden CORS, credential-bearing headers, and oversized result/query workloads.
- Ensure APIs do not return raw credentials, full environments, unsafe absolute paths, or source file content.

## 5. Browser UX and accessibility verification

Exercise the actual packaged application/browser surface for history, scan detail, report completeness, completed/failed/cancelled operation outcomes, no-findings, empty history, loading, storage/API errors, incompatible schema, missing artifact, not found, repeated source scan IDs with distinct history entries, and narrow viewport. Verify shadcn-inspired tokens/components match approved design without asserting CSS implementation details. Test keyboard-only navigation, focus movement/restoration, accessible names/landmarks/table semantics/status announcements, contrast, reduced motion, and hostile repository labels. Capture browser/OS/version and evidence; no simulated component test substitutes for the actual-surface smoke.

## 6. Cross-platform and compatibility

Run unit/contract/integration and packaging checks on Linux, macOS, and Windows. Verify loopback-only listener, origin behavior, paths/permissions, SQLite locking/shutdown, embedded assets, and schema meaning on each supported OS. Compare semantic records and error categories; do not pin platform-specific path strings or incidental messages.

## 7. CLI and end-to-end checks

- Existing CLI commands and JSON/exit contracts pass regression coverage with console disabled and history unavailable according to accepted compatibility decisions.
- End-to-end representative local workflow: produce/use saved history entries from a fixed valid report fixture, including repeated source scan IDs and a cancelled operation with a returned report, open console, browse page/detail/findings, and confirm operation outcome, report completeness, and provenance remain distinct.
- Compare repository tree before/after browser/API workflows; browser must not mutate source repository or package files.
- Exercise no-data startup and storage failure without installing scanners or fetching network resources.

The exact method for creating history fixtures and whether CLI scan automatically records history is a design decision; do not invent a behavior in tests before it is approved.

## 8. Release manual evidence

Record reviewed binary/version/commit, OS/architecture/browser/version, exact commands, exit codes, listener address, history DB location/state, API results, screenshot/accessibility observations, tree comparison, artifact/checksum/SBOM/signature/provenance status, independent reviewer, owner decision, and remaining gaps. Follow `docs/RELEASE-CHECKLIST.md`; tests alone do not approve a release.
