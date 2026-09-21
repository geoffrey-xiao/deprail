package main

import "io"

func writeRootHelp(w io.Writer) error {
	_, err := io.WriteString(w, `Usage: deprail <command> [options]

Start here:
  deprail doctor
  deprail discover .
  deprail scan .

Commands:
  discover       discover dependency workspaces
  scan           scan dependencies for known vulnerabilities
  baseline create create a baseline from a complete scan
  diff           compare two scan baselines
  fix plan       produce a read-only remediation plan
  policy check   evaluate a baseline policy

Output:
  Human output is plain terminal text.
  JSON result data is written to stdout; incidental diagnostics use stderr.
  Use --help on any command for examples and command-specific flags.

Options:
  --help         show this help
  --version      show build information
`)
	return err
}

func writeCommandHelp(w io.Writer, command string) error {
	var text string
	switch command {
	case "discover":
		text = "Usage: deprail discover [path] [--format terminal|json] [--verbose]\n\nDiscover dependency workspaces and completeness.\n\nTip: start with `deprail discover .`; use `--format json` for automation.\nExample: deprail discover . --format json\n"
	case "scan":
		text = "Usage: deprail scan [path] [--format terminal|json] [--output path] [--quiet] [--verbose]\n\nScan dependencies for known vulnerabilities.\n\nTip: use JSON for automation and --verbose for safe diagnostics. Quiet mode keeps required errors.\nExample: deprail scan . --format json --output scan.json\n"
	case "diff":
		text = "Usage: deprail diff --base path --head path [--format terminal|json]\n\nCompare two scan baselines.\n\nTip: use JSON when another tool will consume the comparison.\nExample: deprail diff --base base.json --head head.json --format json\n"
	case "baseline create":
		text = "Usage: deprail baseline create --scan path --output path [--format terminal|json]\n\nCreate a trusted baseline from a complete scan result.\n\nTip: only complete, error-free scans are accepted; output is never overwritten.\nExample: deprail baseline create --scan scan.json --output baseline.json --format json\n"
	case "doctor":
		text = "Usage: deprail doctor [--format terminal|json]\n\nInspect local tool availability without scanning the repository.\n\nTip: run doctor first when scanner availability is uncertain.\nExample: deprail doctor --format json\n"
	case "fix plan":
		text = "Usage: deprail fix plan --report path --finding key [--format terminal|json] [--output path]\n\nProduce a read-only remediation plan; no files are mutated.\n\nTip: obtain the finding key from a completed scan report.\nExample: deprail fix plan --report scan.json --finding FINDING_KEY --format json --output plan.json\n"
	case "policy check":
		text = "Usage: deprail policy check --baseline path [--format json|sarif]\n\nEvaluate a baseline policy and preserve its pass, warn, or block decision.\n\nTip: use SARIF when integrating with a compatible code-scanning workflow.\nExample: deprail policy check --baseline baseline.json --format sarif\n"
	default:
		return writeRootHelp(w)
	}
	_, err := io.WriteString(w, text)
	return err
}

func hasHelp(args []string) bool {
	for _, arg := range args {
		if arg == "--help" || arg == "-h" {
			return true
		}
	}
	return false
}
