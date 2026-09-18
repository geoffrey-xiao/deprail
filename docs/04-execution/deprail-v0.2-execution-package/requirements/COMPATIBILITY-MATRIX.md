# v0.2 Compatibility Matrix

| Surface | Supported baseline | v0.2 requirement | Verification |
| --- | --- | --- | --- |
| OS | Linux, macOS, Windows | Preserve semantic output and exit meanings. | Matrix CI and smoke artifacts. |
| Architecture | amd64; arm64 smoke | Record artifact and smoke coverage explicitly. | Release evidence. |
| Go | Repository `.go-version` | Use pinned toolchain in CI and release builds. | `make verify`. |
| OSV-Scanner | v2 supported range | Exercise real v2 output and exit matrix. | Adapter contract fixtures. |
| Node/npm fixtures | npm lockfile fixtures | Keep fixtures reproducible; do not commit install directories. | Offline/controlled fixture checks. |
| Python fixtures | requirements/uv/Poetry boundaries | Preserve discovery semantics. | Detector and mixed-repository tests. |
| Java fixtures | Maven/Gradle boundaries | Preserve discovery semantics. | Detector and mixed-repository tests. |
| JSON consumers | Current v1alpha contracts | Empty collections are arrays; no breaking schema change approved. | Schema validation and golden tests. |

Unsupported tools or incompatible versions MUST produce explicit diagnostics rather than silent fallback.
