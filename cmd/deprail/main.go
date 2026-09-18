package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/geoffrey-xiao/deprail/internal/app"
	"github.com/geoffrey-xiao/deprail/internal/discovery"
	"github.com/geoffrey-xiao/deprail/internal/presenter"
	"io"
	"os"
	"path/filepath"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		writeCLIError(stderr, "CONFIG_INVALID", "a command is required", "command")
		return 2
	}
	switch args[0] {
	case "discover":
		return runDiscover(args[1:], stdout, stderr)
	case "scan":
		return runScan(args[1:], stdout, stderr)
	case "doctor":
		return runDoctor(args[1:], stdout, stderr)
	default:
		writeCLIError(stderr, "CONFIG_INVALID", fmt.Sprintf("unknown command %q", args[0]), "command")
		return 2
	}
}

func runDiscover(args []string, stdout, stderr io.Writer) int {
	normalized, err := normalizeDiscoverArgs(args)
	if err != nil {
		writeCLIError(stderr, "CONFIG_INVALID", err.Error(), "arguments")
		return 2
	}
	flags := flag.NewFlagSet("discover", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	format := flags.String("format", "terminal", "output format: terminal or json")
	config := flags.String("config", "", "configuration file")
	noIgnore := flags.Bool("no-ignore", false, "do not ignore default directories")
	verbose := flags.Bool("verbose", false, "include diagnostics in terminal output")
	if err := flags.Parse(normalized); err != nil {
		writeCLIError(stderr, "CONFIG_INVALID", err.Error(), "arguments")
		return 2
	}
	if *format != "terminal" && *format != "json" {
		writeCLIError(stderr, "CONFIG_INVALID", "format must be terminal or json", "format")
		return 2
	}
	if *config != "" {
		writeCLIError(stderr, "CONFIG_INVALID", "configuration files are not supported by discover yet", *config)
		return 2
	}
	paths := flags.Args()
	if len(paths) > 1 {
		writeCLIError(stderr, "CONFIG_INVALID", "discover accepts at most one path", "path")
		return 2
	}
	root := "."
	if len(paths) == 1 {
		root = paths[0]
	}
	if err := app.ValidateDiscoverRoot(root); err != nil {
		writeCLIError(stderr, "CONFIG_INVALID", "path must be an existing directory", filepath.Clean(root))
		return 2
	}
	graph, err := app.Discover(context.Background(), root, app.DiscoverOptions{NoIgnore: *noIgnore})
	if err != nil {
		writeCLIError(stderr, "CONFIG_INVALID", "discovery failed", filepath.Clean(root))
		return 2
	}
	var outputErr error
	if *format == "json" {
		outputErr = presenter.WriteDiscoveryJSON(stdout, graph)
	} else {
		outputErr = presenter.WriteDiscoveryTerminal(stdout, graph, *verbose)
	}
	if outputErr != nil {
		writeCLIError(stderr, "OUTPUT_WRITE_FAILED", "could not write discovery output", "stdout")
		return 3
	}
	if graph.Completeness != discovery.Complete {
		return 3
	}
	return 0
}

func normalizeDiscoverArgs(args []string) ([]string, error) {
	flags := make([]string, 0, len(args))
	positionals := make([]string, 0, 1)
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--format", "--config":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("%s requires a value", arg)
			}
			flags = append(flags, arg, args[i+1])
			i++
		case "--no-ignore", "--verbose":
			flags = append(flags, arg)
		default:
			if len(arg) >= 9 && (arg[:9] == "--format=" || arg[:9] == "--config=") {
				flags = append(flags, arg)
				continue
			}
			positionals = append(positionals, arg)
		}
	}
	return append(flags, positionals...), nil
}

func writeCLIError(w io.Writer, code, message, scope string) {
	_, _ = fmt.Fprintf(w, "error[%s] %s (%s)\n", code, message, scope)
}

func runScan(args []string, stdout, stderr io.Writer) int {
	normalized, err := normalizeDiscoverArgs(args)
	if err != nil {
		writeCLIError(stderr, "CONFIG_INVALID", err.Error(), "arguments")
		return 2
	}
	flags := flag.NewFlagSet("scan", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	format := flags.String("format", "terminal", "output format: terminal or json")
	output := flags.String("output", "", "write output atomically to a file")
	verbose := flags.Bool("verbose", false, "include diagnostics in terminal output")
	if err := flags.Parse(normalized); err != nil {
		writeCLIError(stderr, "CONFIG_INVALID", err.Error(), "arguments")
		return 2
	}
	if *format != "terminal" && *format != "json" {
		writeCLIError(stderr, "CONFIG_INVALID", "format must be terminal or json", "format")
		return 2
	}
	root := "."
	if len(flags.Args()) > 1 {
		writeCLIError(stderr, "CONFIG_INVALID", "scan accepts at most one path", "arguments")
		return 2
	}
	if len(flags.Args()) == 1 {
		root = flags.Args()[0]
	}
	if err := app.ValidateDiscoverRoot(root); err != nil {
		writeCLIError(stderr, "PATH_OUTSIDE_ROOT", err.Error(), root)
		return 2
	}
	report, scanErr := app.Scan(context.Background(), root, app.ScanOptions{})
	var rendered bytes.Buffer
	var renderErr error
	if *format == "json" {
		renderErr = presenter.WriteScanJSON(&rendered, report)
	} else {
		renderErr = presenter.WriteScanTerminal(&rendered, report, *verbose)
	}
	err = renderErr
	if err == nil && *output != "" {
		err = presenter.WriteAtomic(*output, rendered.Bytes())
	} else if err == nil {
		_, err = stdout.Write(rendered.Bytes())
	}
	if err != nil {
		writeCLIError(stderr, "OUTPUT_WRITE_FAILED", err.Error(), "output")
		return 3
	}
	if scanErr != nil || report.Status != discovery.Complete {
		return 3
	}
	return 0
}

func runDoctor(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("doctor", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	format := flags.String("format", "terminal", "output format: terminal or json")
	if err := flags.Parse(args); err != nil {
		writeCLIError(stderr, "CONFIG_INVALID", err.Error(), "arguments")
		return 2
	}
	if len(flags.Args()) > 0 {
		writeCLIError(stderr, "CONFIG_INVALID", "doctor does not accept positional arguments", "arguments")
		return 2
	}
	if *format != "terminal" && *format != "json" {
		writeCLIError(stderr, "CONFIG_INVALID", "format must be terminal or json", "format")
		return 2
	}
	report := app.Doctor(context.Background())
	if *format == "json" {
		_ = json.NewEncoder(stdout).Encode(report)
	} else {
		if report.Scanner.Compatible {
			_, _ = fmt.Fprintf(stdout, "DepRail: %s\nTag: %s\nCommit: %s\nOS: %s/%s\nOSV-Scanner: %s\n", report.Version, report.Tag, report.Commit, report.OS, report.Arch, report.Scanner.Version)
		} else {
			_, _ = fmt.Fprintf(stdout, "DepRail: %s\nTag: %s\nCommit: %s\nOS: %s/%s\nOSV-Scanner: %s\n", report.Version, report.Tag, report.Commit, report.OS, report.Arch, report.Scanner.Error)
		}
		if report.Scanner.Error != "" {
			_, _ = fmt.Fprintf(stderr, "Help: %s\n", report.Scanner.Help)
		}
	}
	if !report.Scanner.Available || !report.Scanner.Compatible {
		return 3
	}
	return 0
}
