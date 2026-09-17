# v0.1 Test Strategy

## Test Levels

Unit tests cover domain rules and deterministic transformations. Contract tests cover detectors, adapters, artifact storage, and presenters. Golden tests cover ProjectGraph and ScanReport. Integration tests cover processes and file output. End-to-end tests cover mixed repositories. Cross-platform smoke tests validate packaging.

## Fixture Principles

Fixtures are minimal, offline, license-safe, deterministic, and documented. Raw scanner fixtures include tool version and expected parse behavior. Golden updates require a reason and human review.

## Required Security Cases

External symlink, traversal, shell metacharacters, long path, Unicode, malformed manifest, huge stdout/stderr, timeout, cancellation, process-tree termination, credential-bearing URL, partial scanner failure, atomic-write interruption, and incompatible scanner version.

## Manual Release Verification

Install each platform artifact, run doctor, discover, terminal scan, JSON scan, redirected output, malformed-project scan, missing-scanner path, and one public mixed repository. Record commands, versions, and results.
