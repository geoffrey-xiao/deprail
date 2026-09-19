# DepRail

DepRail is a read-only dependency security guardrail for JavaScript, Python, and Java repositories. It discovers dependency workspaces and scans them with OSV-Scanner.

## Install the preview release

Download the release for your operating system and CPU from:

<https://github.com/geoffrey-xiao/deprail/releases/tag/v0.2.0-preview.4>

| Platform | Artifact |
| --- | --- |
| macOS Intel | `deprail-darwin-amd64` |
| macOS Apple Silicon | `deprail-darwin-arm64` |
| Linux x86_64 | `deprail-linux-amd64` |
| Windows x86_64 | `deprail-windows-amd64.exe` |

### Verify the download

Download `SHA256SUMS` into the same directory as the binary. On macOS or Linux:

```bash
cd ~/Downloads
shasum -a 256 -c SHA256SUMS
```

The selected binary must report `OK`. Do not execute an artifact whose checksum fails.

The release also publishes an SPDX SBOM, Cosign signatures and certificates, and build provenance. Review those assets when your organization requires supply-chain verification.

### Install on macOS or Linux

This example installs the Apple Silicon artifact. Substitute `deprail-darwin-amd64` for an Intel Mac or `deprail-linux-amd64` for Linux:

```bash
mkdir -p "$HOME/bin"
install -m 755 "$HOME/Downloads/deprail-darwin-arm64" "$HOME/bin/deprail"

case ":$PATH:" in
  *":$HOME/bin:"*) ;;
  *) printf '\nexport PATH="$HOME/bin:$PATH"\n' >> "$HOME/.zshrc" ;;
esac

source "$HOME/.zshrc"
deprail --version
```

On Linux with Bash, put the same PATH line in `~/.bashrc` and run:

```bash
source "$HOME/.bashrc"
```

If macOS blocks the verified download because it has a quarantine attribute, remove that attribute only after checksum verification:

```bash
xattr -d com.apple.quarantine "$HOME/bin/deprail" 2>/dev/null || true
```

### Install on Windows

In PowerShell:

```powershell
New-Item -ItemType Directory -Force "$HOME\bin" | Out-Null
Copy-Item "$HOME\Downloads\deprail-windows-amd64.exe" "$HOME\bin\deprail.exe"
$env:Path = "$HOME\bin;$env:Path"
deprail.exe --version
```

For a persistent Windows PATH, add `$HOME\bin` through **System Settings → Environment Variables**, then open a new terminal.

## Check the environment

Install OSV-Scanner separately using your organization’s approved tool-management process. DepRail never installs tools or accesses the network during discovery.

```bash
deprail --version
deprail doctor
deprail doctor --format json
```

## Discover and scan a repository

```bash
deprail discover /path/to/repository
deprail discover /path/to/repository --format json

deprail scan /path/to/repository
deprail scan /path/to/repository --format json
```

Machine-readable data is written to stdout. Diagnostics and guidance are written to stderr. Use `--output scan.json` when a scan report should be written to a new file.

Exit codes:

- `0`: complete execution;
- `2`: invalid command or arguments;
- `3`: scanner failure or incomplete result.

Never interpret a partial or failed scan as safe. See [`docs/QUICKSTART.md`](docs/QUICKSTART.md) for troubleshooting and repository contributor instructions.
