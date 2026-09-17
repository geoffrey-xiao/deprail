# DepRail Project Naming Standard

## Official Names

- Product: `DepRail`
- Recommended GitHub repository: `deprail`
- CLI executable: `deprail`
- Go command entry point: `cmd/deprail`
- Recommended MCP executable: `deprail-mcp`
- Repository configuration: `.deprail.yaml`
- Local data directory: `.deprail/`

## Command Examples

```bash
deprail doctor
deprail discover .
deprail scan .
deprail diff --base origin/main --head HEAD
deprail policy check --baseline .deprail/baseline.json --format sarif
deprail fix plan GHSA-xxxx-xxxx-xxxx
deprail fix apply --plan .deprail/plans/plan-001.json --verify
deprail serve
```

## Naming Rules

Documentation, code, issues, CI configuration, examples, and release artifacts must use the full name `DepRail`. Commands, filenames, directory names, and package names use lowercase `deprail`. Historical codenames and commands must not appear in current project materials.

Public positioning statement:

> DepRail is a dependency security guardrail for multi-language repositories.
