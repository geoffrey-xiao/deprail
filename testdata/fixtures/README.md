# DepRail fixture catalog

These fixtures are minimal, synthetic, offline, deterministic, and license-safe. They are discovery inputs, not applications to build or install.

## Catalog

| Fixture | Intended detection | Expected authoritative inputs |
| --- | --- | --- |
| `npm-basic` | npm | `package.json`, `package-lock.json` |
| `python-requirements` | Python requirements | `requirements.txt` |
| `python-uv` | Python uv | `pyproject.toml`, `uv.lock` |
| `java-maven` | Maven | `pom.xml` |
| `mixed-repository` | mixed monorepo | frontend npm, API requirements, worker Maven |

Do not run package installation or network-dependent commands against fixtures. Keep manifests small and pinned. When changing a fixture, explain the detection behavior and update the contract test. Fixture updates require human review when they change golden or expected discovery behavior.
