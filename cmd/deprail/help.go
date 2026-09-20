package main

import "io"

func writeRootHelp(w io.Writer) error {
	_, err := io.WriteString(w, `Usage: deprail <command> [options]

Commands:
  discover       discover dependency workspaces
  scan           scan dependencies for known vulnerabilities
  doctor         inspect local tool availability
  diff           compare two scan baselines
  fix plan       produce a read-only remediation plan
  policy check   evaluate a baseline policy

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
		text = "Usage: deprail discover [path] [--format terminal|json] [--verbose]\n\nDiscover dependency workspaces.\n"
	case "scan":
		text = "Usage: deprail scan [path] [--format terminal|json] [--output path] [--quiet] [--verbose]\n\nScan dependencies for known vulnerabilities.\n"
	case "doctor":
		text = "Usage: deprail doctor [--format terminal|json]\n\nInspect local tool availability.\n"
	case "diff":
		text = "Usage: deprail diff --base path --head path [--format terminal|json]\n\nCompare two scan baselines.\n"
	case "fix plan":
		text = "Usage: deprail fix plan --report path --finding key [--format terminal|json]\n\nProduce a read-only remediation plan.\n"
	case "policy check":
		text = "Usage: deprail policy check --baseline path [--format json|sarif]\n\nEvaluate a baseline policy.\n"
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
