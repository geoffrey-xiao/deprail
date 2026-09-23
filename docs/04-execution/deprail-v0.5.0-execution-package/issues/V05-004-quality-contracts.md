# V05-004: Complete Failure, Security, Compatibility, and Test Contracts
- GitHub Issue: [#393](https://github.com/geoffrey-xiao/deprail/issues/393)

- Epic: [EPIC-001 / #390](https://github.com/geoffrey-xiao/deprail/issues/390)
- Target: `v0.5.0`
- Status: Todo; planning/design deliverable
- Type: test
- Area: docs
- Priority: P0
- Risk: R3
- Owner: `@geoffrey-xiao`
- Reviewer: `@geoffrey-xiao` for planning-ticket oversight only, by explicit owner direction; not independent architecture/security approval
- Dependencies: V05-001, V05-002, V05-003; V05-000 / #387

## Value

The release needs observable evidence for failure preservation, security boundaries, compatibility, and packaged cross-platform behavior—not only a happy-path UI or unit suite.

## Scope

Finalize the failure/data, error, security, compatibility, and test-strategy contracts. Map every FR-501–FR-511 and relevant API/storage/UI failure to a consumer-visible result, deterministic test or manual scenario, evidence artifact, and platform/browser matrix. Cover operation outcome versus report completeness, missing/corrupt/unavailable history, migration interruption, transaction rollback, locking/disk-full/permissions, missing/digest-mismatched artifacts, hostile Host/Origin/path/SQL/XSS input, oversized requests, cancellation/timeouts, embedded asset/API version mismatch, CLI-without-web behavior, and Linux/macOS/Windows semantics. Set toolchain/browser/SQLite-driver decision inputs without selecting unreviewed dependencies.

## Out of scope

No test harness, fixtures requiring source mutation, runtime test code, scanner changes, frontend dependency installation, release publication, or relaxation of an existing CLI/schema contract.

## Inputs, outputs, and failure behavior

Inputs: accepted UX, API/OpenAPI and SQLite contracts plus existing versioned report/error behavior. Outputs: updated requirement/test crosswalk, exact future verification scenarios, threat-to-control evidence map, and compatibility/rollback/release-evidence matrix. Every storage/API/UI failure remains distinguishable from empty history and complete zero findings; unresolved contract conflicts remain explicit blockers.

## Required verification and evidence

Check traceability from each FR to design, failure mode, observable assertion/manual scenario, and evidence. Include offline deterministic fixture constraints, actual-browser requirements, CLI regression coverage, and cross-platform artifact smoke. This issue defines verification; it does not claim those future runtime checks have passed.

## Acceptance criteria

- All specified failures and adversarial cases map to explicit behavior, protection, and evidence.
- Compatibility preserves CLI/report semantics unless a separate reviewed change is accepted.
- Linux/macOS/Windows and supported-browser evidence requirements are explicit and comparable.
- Required verification distinguishes genuine behavior from implementation-specific or mock-only assertions.
- Independent security review and runtime authorization remain outstanding until the v0.5 Definition of Ready passes.
