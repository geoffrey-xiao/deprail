# H05-006 Frontend Dependency Evidence

Issue: [#425](https://github.com/geoffrey-xiao/deprail/issues/425).
Status: **dependency review incomplete; owner disposition required before approval**.

## Manifest and resolution

`web/package.json` pins the approved direct packages exactly: React/React DOM `19.1.1`, TypeScript `5.9.2`, Vite `7.3.6`, `@vitejs/plugin-react` `5.0.2`, `@types/react` `19.1.10`, `@types/react-dom` `19.1.7`, and `@types/node` `22.18.1`. Runtime dependencies are React and React DOM only. The build is `tsc --noEmit && vite build`; `web/.npmrc` defaults to `ignore-scripts=true`. Rollup is overridden to exact `4.63.0`: the initial lock's `4.64.0` stalled on even a minimal ReactDOM Vite build for 120 seconds; `4.52.4` built but `npm audit` flagged high-severity [GHSA-mw96-cpmx-2vgc](https://github.com/advisories/GHSA-mw96-cpmx-2vgc), so that vulnerable candidate was rejected. `4.63.0` built the full UI in 412 ms and is outside the advisory's affected `4.0.0–4.58.0` range.

The official Node `v22.23.3` Darwin arm64 archive SHA-256 `23b25245dcfb9af7262f8ff142e9e2e0af025368117329e7a7458a51e5922f53` matched Node's published `SHASUMS256.txt`; its bundled npm reports `10.9.9`. With that exact toolchain, `npm install --package-lock-only --ignore-scripts --no-fund --no-audit` and `npm ci --ignore-scripts --no-fund --no-audit` each exited 0. No dependency lifecycle scripts ran. The earlier host Node `26.7.0`/npm `11.19.0` resolution was superseded by this pinned verification.

Pinned npm `audit --audit-level=low` exited 0 with zero reported vulnerabilities after overriding Rollup. The 121 lock-resolved entries include platform variants; 71 packages install on this Darwin arm64 host. Registry advisory results were observed on 2026-10-06, not a guarantee of future safety. `package-lock.json` records registry tarball URLs, integrity hashes and reported licenses; the generated CycloneDX SBOM reflects the 71 installed packages only, not the whole cross-platform lock.

## Registry checks and provenance

Ran `npm view <package>@<version> license dist.integrity dist.attestations.provenance --json` against registry.npmjs.org for all eight direct dependencies. Registry license and integrity metadata matched the lock; React, React DOM, Vite, plugin-react and types metadata report MIT, TypeScript Apache-2.0. Registry provenance metadata was exposed for Vite and plugin-react (SLSA provenance predicate type); metadata was not returned for the other six packages by this query. Provenance attestations were not downloaded or cryptographically verified. Direct package metadata URLs are `https://registry.npmjs.org/<package>/<version>`; scoped type packages use their encoded npm registry paths. Lock integrity pins fetched tarball digests, but does not itself establish publisher provenance.

The lock's license strings are registry/package metadata, not an independent legal review. The SBOM contains hashes only for the packages actually installed on the host platform; the lock integrity values cover the complete resolved platform variants.
## License and provenance review risks

- The 121 locked entries report 112 MIT, two Apache-2.0, five ISC, one BSD-3-Clause, and one CC-BY-4.0 package (`caniuse-lite@1.0.30001814`); the root package has no declared license. `caniuse-lite` enters through `@vitejs/plugin-react` → `@babel/core` → `@babel/helper-compilation-targets` → `browserslist`. This is a build-time dependency, not a declared browser runtime dependency. Distribution/attribution implications require explicit owner disposition before release approval; the registry identifier alone does not settle them.
- Registry metadata did not establish provenance for six direct packages, and no attestations were verified. Transitive publisher provenance remains unverified; lock integrity is not a substitute. Owner security/supply-chain review must accept or mitigate this residual risk before release approval.
- No package lifecycle scripts ran. The pinned build and source-level real-API browser smoke are recorded separately; cross-platform package installation, provenance attestation verification, distribution attribution and packaged-binary review remain open for owner/release review.
