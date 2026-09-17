# v0.1 Data Contract

All documents include `schema_version`. Identifiers and ordering are deterministic.

## Project

Repository root identity, display name, configuration digest, workspaces, completeness, and diagnostics.

## Workspace

Stable workspace ID, relative path, ecosystem, package manager, manifests, lockfiles, scope, completeness, and warnings.

## Scan

Scan ID, timestamps, tool metadata, plan, status, findings, errors, and artifact digests. Time and run IDs are explicitly volatile.

## Component

Canonical PURL where possible, ecosystem, name, resolved version, workspace, directness, and dependency paths.

## Vulnerability

Canonical ID, sorted aliases, affected ranges, source records, severity observations, references, and fixed-version observations.

## Finding

Stable key, component, vulnerability, evidence, fixed versions, and state. A finding is unique for the normalized component version and vulnerability alias set within scope.

## Evidence

Adapter, tool version, source database, source ID, workspace, dependency path, raw artifact digest, and optional source location.

## Stability Rules

Stable keys exclude descriptions, timestamps, severity labels, and evidence ordering. Lists use declared sort keys. Unknown data is preserved where safe and never converted into confidence.
