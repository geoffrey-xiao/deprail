# OSV-Scanner v2 raw fixtures

These fixtures were captured from OSV-Scanner 2.6.0 (`osv-scalibr` 0.5.2) using:

```text
osv-scanner scan source --format json testdata/fixtures/python-uv
osv-scanner scan source --format json testdata/fixtures/npm-basic
```

- `v2-empty.json` is the complete empty-result response for the `python-uv` fixture.
- `v2-findings.json` is the vulnerability response for the `npm-basic` fixture (`lodash` 4.17.20), including aliases, severity, affected ranges, and fixed events.
- The scanner's absolute source path was replaced with the repository-relative fixture path so the checked-in sample remains deterministic across checkouts. No vulnerability content was changed.
- The output is retained as raw scanner JSON. Parser tests assert empty results, finding identity, aliases, severity, and fixed-version extraction.

The findings fixture is sourced from the public OSV database response emitted by the scanner and contains no repository secrets or private source content.
