# Toolchain and Verification

DepRail pins the Go toolchain in `.go-version`. Local verification uses the root `Makefile` and does not install tools or access the network implicitly.

```bash
make bootstrap
make generate
make test
make lint
make build
make verify
make test-integration
```

`make verify` checks the pinned Go version, runs generation, checks formatting, runs `go vet`, executes tests, and builds all packages. The GitHub Actions matrix runs the equivalent formatting, vet, test, and build checks on Linux, macOS, and Windows.

The repository currently has no integration-tagged tests and no generators. Their targets are still defined so later work has stable entry points. Tool and scanner upgrades require contract-fixture review.
