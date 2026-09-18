# v0.1 Compatibility Matrix

## Project Discovery

| Ecosystem | Supported inputs | v0.1 level |
| --- | --- | --- |
| JavaScript | package.json, package-lock, npm-shrinkwrap | Full |
| JavaScript | pnpm workspace and lockfile | Discovery |
| JavaScript | Yarn workspaces and common lockfiles | Discovery |
| Python | requirements files, uv lock/project | Full discovery |
| Python | Poetry project and lockfile | Discovery |
| Java | Maven pom.xml | Full discovery |
| Java | Gradle settings and build files | Discovery |

## Platforms

Linux, macOS, and Windows are required. amd64 release artifacts are required; arm64 receives at least release smoke coverage. Semantic JSON output must be equivalent across platforms.

## OSV-Scanner

The pinned supported range is OSV-Scanner `>=2.0.0 <3.0.0`. The v2 CLI `scan source` command and JSON output are required. Missing, incompatible, or malformed output produces a stable diagnostic. Automatic installation is not supported.

## Explicitly Unsupported

Containers, SBOM input, licenses, secrets, IaC, package installation, remote repositories, source upload, and automatic remediation are outside v0.1.
