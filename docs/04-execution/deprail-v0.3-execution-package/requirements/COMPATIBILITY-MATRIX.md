# v0.3 Compatibility Matrix

| Surface | Required behavior | Compatibility rule |
| --- | --- | --- |
| Go/CLI | Existing v0.2 commands and exit meanings remain valid | No silent changes |
| JSON | Existing scan/diff/policy/SARIF schemas remain valid | New remediation schema is versioned |
| OS | Linux, macOS, Windows plan semantics are equivalent | Path/command presentation differences must be explicit |
| Ecosystems | npm/pnpm/Yarn; requirements/uv/Poetry; Maven/Gradle | Unsupported managers fail explicitly |
| Scanner evidence | Existing OSV provenance and finding keys remain consumable | No reinterpretation of stable identities |
| Configuration | Existing `.deprail.yaml` remains valid | New planning settings are additive and schema-validated |
| Storage | Existing artifacts remain readable | Plan storage is additive and digest-bound |
| Agent/web future clients | Plans are structured, versioned, and read-only | No write capability implied |

Any breaking change requires an ADR, migration note, compatibility tests, and owner/security approval.
