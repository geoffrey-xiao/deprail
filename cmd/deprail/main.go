package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/geoffrey-xiao/deprail/internal/app"
	"github.com/geoffrey-xiao/deprail/internal/baseline"
	"github.com/geoffrey-xiao/deprail/internal/buildinfo"
	"github.com/geoffrey-xiao/deprail/internal/discovery"
	"github.com/geoffrey-xiao/deprail/internal/policy"
	"github.com/geoffrey-xiao/deprail/internal/presenter"
	"github.com/geoffrey-xiao/deprail/internal/remediation"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 1 && (args[0] == "--help" || args[0] == "-h") {
		if err := writeRootHelp(stdout); err != nil {
			return 3
		}
		return 0
	}
	if len(args) == 1 && (args[0] == "--version" || args[0] == "version") {
		identity := buildinfo.Current()
		_, _ = fmt.Fprintf(stdout, "deprail %s (tag=%s commit=%s)\n", identity.Version, identity.Tag, identity.Commit)
		return 0
	}
	if len(args) == 0 {
		writeCLIError(stderr, "CONFIG_INVALID", "a command is required", "command")
		return 2
	}
	if hasHelp(args[1:]) {
		command := args[0]
		switch command {
		case "discover", "scan", "doctor", "diff":
		case "baseline":
			if len(args) < 2 || args[1] != "create" {
				writeCLIError(stderr, "CONFIG_INVALID", "baseline requires the create subcommand", "baseline")
				return 2
			}
			command = "baseline create"
		case "fix":
			if len(args) < 2 || (args[1] != "plan" && args[1] != "approve" && args[1] != "apply") {
				writeCLIError(stderr, "CONFIG_INVALID", "fix requires the plan, approve, or apply subcommand", "fix")
				return 2
			}
			command = "fix " + args[1]
		case "policy":
			if len(args) < 2 || args[1] != "check" {
				writeCLIError(stderr, "CONFIG_INVALID", "policy requires the check subcommand", "policy")
				return 2
			}
			command = "policy check"
		default:
			writeCLIError(stderr, "CONFIG_INVALID", fmt.Sprintf("unknown command %q", command), "command")
			return 2
		}
		if err := writeCommandHelp(stdout, command); err != nil {
			return 3
		}
		return 0
	}
	switch args[0] {
	case "discover":
		return runDiscover(args[1:], stdout, stderr)
	case "scan":
		return runScan(args[1:], stdout, stderr)
	case "fix":
		if len(args) < 2 {
			writeCLIError(stderr, "CONFIG_INVALID", "fix requires the plan, approve, or apply subcommand", "fix")
			return 2
		}
		switch args[1] {
		case "plan":
			return runFixPlan(args[2:], stdout, stderr)
		case "approve":
			return runFixApprove(args[2:], stdout, stderr)
		case "apply":
			return runFixApply(args[2:], stdout, stderr)
		default:
			writeCLIError(stderr, "CONFIG_INVALID", "fix requires the plan, approve, or apply subcommand", "fix")
			return 2
		}
	case "diff":
		return runDiff(args[1:], stdout, stderr)
	case "baseline":
		return runBaseline(args[1:], stdout, stderr)
	case "policy":
		if len(args) < 2 || args[1] != "check" {
			writeCLIError(stderr, "CONFIG_INVALID", "policy requires the check subcommand", "policy")
			return 2
		}
		return runPolicyCheck(args[2:], stdout, stderr)
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
	writeStartNotice(*format, stderr, "Discovering workspaces...")
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

func normalizeScanArgs(args []string) ([]string, error) {
	flags := make([]string, 0, len(args))
	positionals := make([]string, 0, 1)
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--" {
			positionals = append(positionals, args[i:]...)
			break
		}
		switch arg {
		case "--format", "--output":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("%s requires a value", arg)
			}
			flags = append(flags, arg, args[i+1])
			i++
		case "--verbose", "--quiet":
			flags = append(flags, arg)
		default:
			if strings.HasPrefix(arg, "--format=") || strings.HasPrefix(arg, "--output=") {
				flags = append(flags, arg)
				continue
			}
			positionals = append(positionals, arg)
		}
	}
	return append(flags, positionals...), nil
}

func writeCLIError(w io.Writer, code, message, scope string) {
	_ = presenter.WriteError(w, code, scope+": "+message)
}
func writeStartNotice(format string, w io.Writer, message string) {
	if format == "terminal" && presenter.IsTerminal(w) {
		_ = presenter.WriteHeader(w, message)
	}
}

func runScan(args []string, stdout, stderr io.Writer) int {
	normalized, err := normalizeScanArgs(args)
	if err != nil {
		writeCLIError(stderr, "CONFIG_INVALID", err.Error(), "arguments")
		return 2
	}
	flags := flag.NewFlagSet("scan", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	format := flags.String("format", "terminal", "output format: terminal or json")
	output := flags.String("output", "", "write output atomically to a file")
	verbose := flags.Bool("verbose", false, "include diagnostics in terminal output")
	quiet := flags.Bool("quiet", false, "suppress successful terminal summaries and progress")
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
	scanOptions := app.ScanOptions{}
	capabilities := presenter.SelectCapabilities(presenter.CapabilityOptions{
		Mode:   presenter.OutputMode(*format),
		Stdout: stdout,
		Stderr: stderr,
	})
	writeStartNotice(*format, stderr, "Scanning...")
	if capabilities.Interactive && !*quiet {
		scanOptions.Events = presenter.ProgressSink{Writer: stderr}
	}
	report, scanErr := app.Scan(context.Background(), root, scanOptions)
	var rendered bytes.Buffer
	var renderErr error
	if *format == "json" {
		renderErr = presenter.WriteScanJSON(&rendered, report)
	} else if !*quiet || report.Status != discovery.Complete || len(report.Errors) > 0 {
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

func runFixPlan(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("fix plan", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	report := flags.String("report", "", "normalized scan report")
	finding := flags.String("finding", "", "finding stable key")
	root := flags.String("root", "", "repository root")
	repositoryState := flags.String("repository-state", "", "current repository state")
	verificationPath := flags.String("verification", "", "explicit verification command JSON")
	output := flags.String("output", "", "write JSON plan atomically to an external file")
	format := flags.String("format", "terminal", "output format: terminal or json")
	if err := flags.Parse(args); err != nil || *report == "" || *finding == "" || (*format != "terminal" && *format != "json") || len(flags.Args()) != 0 {
		writeCLIError(stderr, "CONFIG_INVALID", "fix plan requires --report and --finding and supports --verification, terminal or json output", "fix plan")
		return 2
	}
	var explicitVerification []remediation.Verification
	if *verificationPath != "" {
		loaded, err := loadVerificationCommands(*verificationPath)
		if err != nil {
			writeCLIError(stderr, "CONFIG_INVALID", err.Error(), "verification")
			return 2
		}
		explicitVerification = loaded
	}
	if *output != "" {
		if err := app.ValidatePlanOutput(*report, *output, *root); err != nil && errors.Is(err, remediation.ErrUnsafePath) {
			writeCLIError(stderr, "PLAN_OUTPUT_OUTSIDE_ROOT_REQUIRED", err.Error(), "output")
			return 2
		}
	}
	writeStartNotice(*format, stderr, "Planning remediation...")
	plan, err := app.Plan(context.Background(), *report, *finding, *root, app.PlanOptions{CurrentRepositoryState: *repositoryState, Verification: explicitVerification})
	if err != nil {
		var reportErr *remediation.ReportError
		if errors.As(err, &reportErr) {
			writeCLIError(stderr, string(reportErr.Code), reportErr.Message, "fix plan")
			if reportErr.Code == remediation.FindingNotFound || reportErr.Code == remediation.FindingAmbiguous {
				return 2
			}
		} else if errors.Is(err, app.ErrPlanIncomplete) {
			writeCLIError(stderr, "PLAN_INCOMPLETE", err.Error(), "report")
		} else if errors.Is(err, app.ErrPlanUnsupported) {
			writeCLIError(stderr, "PLAN_UNSUPPORTED", err.Error(), "finding")
		} else if errors.Is(err, app.ErrPlanRejected) {
			writeCLIError(stderr, "PLAN_REJECTED", err.Error(), "repository")
		} else if errors.Is(err, app.ErrPlanUnknown) {
			writeCLIError(stderr, "PLAN_UNKNOWN", err.Error(), "repository")
		} else {
			writeCLIError(stderr, "PLAN_FAILED", err.Error(), "fix plan")
		}
		return 3
	}
	if *output != "" {
		if err := presenter.WritePlanAtomic(*output, plan); err != nil {
			switch {
			case errors.Is(err, presenter.ErrPlanOutputOutsideRoot):
				writeCLIError(stderr, "PLAN_OUTPUT_OUTSIDE_ROOT_REQUIRED", err.Error(), "output")
			case errors.Is(err, presenter.ErrPlanSchemaInvalid):
				writeCLIError(stderr, "PLAN_SCHEMA_INVALID", err.Error(), "plan")
			default:
				writeCLIError(stderr, "PLAN_WRITE_FAILED", err.Error(), "output")
			}
			return 3
		}
		return 0
	}
	var outputErr error
	if *format == "json" {
		outputErr = presenter.WritePlanJSON(stdout, plan)
	} else {
		outputErr = presenter.WritePlanTerminal(stdout, plan)
	}

	if outputErr != nil {
		writeCLIError(stderr, "OUTPUT_WRITE_FAILED", outputErr.Error(), "stdout")
		return 3
	}
	return 0
}
func loadVerificationCommands(path string) ([]remediation.Verification, error) {
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var input struct {
		Verification []remediation.Verification `json:"verification"`
	}
	decoder := json.NewDecoder(io.LimitReader(file, 16<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		return nil, err
	}
	if len(input.Verification) == 0 {
		return nil, errors.New("verification input must contain at least one command")
	}
	return input.Verification, nil
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
			_ = presenter.WriteHeader(stdout, "DepRail: "+report.Version)
			_ = presenter.WriteSummary(stdout, "Tag", report.Tag)
			_ = presenter.WriteSummary(stdout, "Commit", report.Commit)
			_ = presenter.WriteSummary(stdout, "OS", report.OS+"/"+report.Arch)
			_ = presenter.WriteSummary(stdout, "OSV-Scanner", report.Scanner.Version)
		} else {
			_ = presenter.WriteHeader(stdout, "DepRail: "+report.Version)
			_ = presenter.WriteSummary(stdout, "Tag", report.Tag)
			_ = presenter.WriteSummary(stdout, "Commit", report.Commit)
			_ = presenter.WriteSummary(stdout, "OS", report.OS+"/"+report.Arch)
			_ = presenter.WriteSummary(stdout, "OSV-Scanner", report.Scanner.Error)
		}
		if report.Scanner.Error != "" {
			_ = presenter.WriteSummary(stderr, "Help", report.Scanner.Help)
		}
	}
	if !report.Scanner.Available || !report.Scanner.Compatible {
		return 3
	}
	return 0
}
func runDiff(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("diff", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	basePath := flags.String("base", "", "base baseline path")
	headPath := flags.String("head", "", "head baseline path")
	format := flags.String("format", "terminal", "output format: terminal or json")
	if err := flags.Parse(args); err != nil || *basePath == "" || *headPath == "" || (*format != "terminal" && *format != "json") || len(flags.Args()) != 0 {
		writeCLIError(stderr, "CONFIG_INVALID", "diff requires --base and --head and supports terminal or json output", "diff")
		return 2
	}
	writeStartNotice(*format, stderr, "Comparing baselines...")
	load := func(input string) (baseline.Document, error) {
		info, err := os.Stat(input)
		if err != nil {
			return baseline.Document{}, fmt.Errorf("input %q is not a readable baseline path: %w", input, err)
		}
		if info.IsDir() {
			return baseline.Document{}, fmt.Errorf("input %q is a directory, not a baseline", input)
		}
		return baseline.LoadFile(input)
	}
	base, err := load(*basePath)
	if err != nil {
		writeCLIError(stderr, "BASELINE_INVALID", err.Error(), "base")
		return 3
	}
	head, err := load(*headPath)
	if err != nil {
		writeCLIError(stderr, "BASELINE_INVALID", err.Error(), "head")
		return 3
	}
	comparison, err := baseline.Compare(base, head)
	if err != nil {
		writeCLIError(stderr, "BASELINE_INVALID", err.Error(), "comparison")
		return 3
	}
	if *format == "json" {
		document := struct {
			SchemaVersion string `json:"schema_version"`
			DocumentType  string `json:"document_type"`
			baseline.Comparison
		}{SchemaVersion: baseline.SchemaVersion, DocumentType: "diff", Comparison: comparison}
		if err := json.NewEncoder(stdout).Encode(document); err != nil {
			writeCLIError(stderr, "OUTPUT_WRITE_FAILED", err.Error(), "output")
			return 3
		}
		return 0
	}
	if err := presenter.WriteHeader(stdout, fmt.Sprintf("Diff %s -> %s", comparison.BaseScanID, comparison.HeadScanID)); err != nil {
		return 3
	}
	for _, change := range comparison.Changes {
		if err := presenter.WriteSummary(stdout, string(change.Kind), change.Key); err != nil {
			return 3
		}
	}
	return 0
}
func runPolicyCheck(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("policy check", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	baselinePath := flags.String("baseline", "", "baseline path")
	format := flags.String("format", "json", "output format: json or sarif")
	if err := flags.Parse(args); err != nil || *baselinePath == "" || (*format != "json" && *format != "sarif") || len(flags.Args()) != 0 {
		writeCLIError(stderr, "CONFIG_INVALID", "policy check requires --baseline and --format json or sarif", "policy check")
		return 2
	}
	document, err := baseline.LoadFile(*baselinePath)
	if err != nil {
		writeCLIError(stderr, "SCAN_UNUSABLE", err.Error(), "baseline")
		return 3
	}
	if *format == "sarif" {
		data, err := presenter.SARIFFromBaseline(document)
		if err != nil {
			writeCLIError(stderr, "SCAN_UNUSABLE", err.Error(), "baseline")
			return 3
		}
		if _, err := stdout.Write(append(data, '\n')); err != nil {
			writeCLIError(stderr, "OUTPUT_WRITE_FAILED", err.Error(), "output")
			return 3
		}
		return 0
	}
	decision, err := policy.Evaluate(policy.Policy{}, policy.Input{Complete: true})
	if err != nil {
		writeCLIError(stderr, "CONFIG_INVALID", err.Error(), "policy")
		return 2
	}
	if err := json.NewEncoder(stdout).Encode(decision); err != nil {
		writeCLIError(stderr, "OUTPUT_WRITE_FAILED", err.Error(), "output")
		return 3
	}
	return policy.ExitCode(decision)
}
