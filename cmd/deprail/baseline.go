package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/geoffrey-xiao/deprail/internal/adapter"
	"github.com/geoffrey-xiao/deprail/internal/app"
	"github.com/geoffrey-xiao/deprail/internal/baseline"
)

const maxScanInputBytes = 16 << 20

func runBaseline(args []string, stdout, stderr io.Writer) int {
	if len(args) < 1 || args[0] != "create" {
		writeCLIError(stderr, "CONFIG_INVALID", "baseline requires the create subcommand", "baseline")
		return 2
	}
	return runBaselineCreate(args[1:], stdout, stderr)
}

func runBaselineCreate(args []string, stdout, stderr io.Writer) int {
	flags := newFlagSet("baseline create")
	scanPath := flags.String("scan", "", "complete scan JSON input")
	outputPath := flags.String("output", "", "baseline JSON output path")
	format := flags.String("format", "terminal", "output format: terminal or json")
	if err := flags.Parse(args); err != nil || *scanPath == "" || *outputPath == "" || (*format != "terminal" && *format != "json") || len(flags.Args()) != 0 {
		writeCLIError(stderr, "CONFIG_INVALID", "baseline create requires --scan and --output and supports terminal or json output", "baseline create")
		return 2
	}

	input, err := loadScanInput(*scanPath)
	if err != nil {
		writeCLIError(stderr, scanInputErrorCode(err), err.Error(), "scan")
		return 3
	}
	document, err := baseline.ConvertScan(input)
	if err != nil {
		var conversionErr *baseline.ConversionError
		if errors.As(err, &conversionErr) {
			writeCLIError(stderr, string(conversionErr.Code), conversionErr.Message, "scan")
		} else {
			writeCLIError(stderr, "BASELINE_INPUT_INVALID", err.Error(), "scan")
		}
		return 3
	}
	if err := baseline.SaveAs(*outputPath, document); err != nil {
		if errors.Is(err, fs.ErrExist) {
			writeCLIError(stderr, "BASELINE_OUTPUT_EXISTS", "output file already exists", "output")
		} else {
			writeCLIError(stderr, "BASELINE_WRITE_FAILED", err.Error(), "output")
		}
		return 3
	}
	if *format == "json" {
		data, err := baseline.Marshal(document)
		if err != nil {
			writeCLIError(stderr, "OUTPUT_WRITE_FAILED", err.Error(), "stdout")
			return 3
		}
		_, err = stdout.Write(data)
	} else {
		_, err = fmt.Fprintf(stdout, "Baseline: %s\nFindings: %d\n", *outputPath, len(document.Findings))
	}
	if err != nil {
		writeCLIError(stderr, "OUTPUT_WRITE_FAILED", err.Error(), "stdout")
		return 3
	}
	return 0
}

func newFlagSet(name string) *flag.FlagSet {
	flags := flag.NewFlagSet(name, flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	return flags
}

func loadScanInput(path string) (baseline.ScanInput, error) {
	file, err := os.Open(path)
	if err != nil {
		return baseline.ScanInput{}, fmt.Errorf("scan input is unreadable")
	}
	defer file.Close()
	data, err := io.ReadAll(io.LimitReader(file, maxScanInputBytes+1))
	if err != nil {
		return baseline.ScanInput{}, fmt.Errorf("scan input is unreadable")
	}
	if len(data) > maxScanInputBytes {
		return baseline.ScanInput{}, fmt.Errorf("scan input exceeds the input limit")
	}
	if err := requireScanArrays(data); err != nil {
		return baseline.ScanInput{}, err
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	var report app.ScanReport
	if err := decoder.Decode(&report); err != nil {
		return baseline.ScanInput{}, fmt.Errorf("scan input is malformed")
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return baseline.ScanInput{}, fmt.Errorf("scan input contains trailing data")
	}
	var wire struct {
		Findings []struct {
			TargetID       string `json:"target_id"`
			LegacyTargetID string `json:"TargetID"`
		} `json:"findings"`
	}
	if err := json.Unmarshal(data, &wire); err != nil {
		return baseline.ScanInput{}, fmt.Errorf("scan input is malformed")
	}
	findings := make([]baseline.ScanFinding, 0, len(report.Findings))
	for index, finding := range report.Findings {
		converted := scanFinding(finding)
		if converted.TargetID == "" && index < len(wire.Findings) {
			converted.TargetID = wire.Findings[index].TargetID
			if converted.TargetID == "" {
				converted.TargetID = wire.Findings[index].LegacyTargetID
			}
		}
		findings = append(findings, converted)
	}
	return baseline.ScanInput{
		SchemaVersion:   report.SchemaVersion,
		DocumentType:    report.DocumentType,
		ScanID:          report.ScanID,
		RepositoryState: report.RepositoryState,
		Status:          report.Status,
		Findings:        findings,
		Errors:          cloneStrings(report.Errors),
		ArtifactDigests: cloneStrings(report.ArtifactDigests),
	}, nil
}

func scanFinding(finding adapter.Finding) baseline.ScanFinding {
	return baseline.ScanFinding{
		Component:     finding.Component,
		PURL:          finding.PURL,
		Version:       finding.Version,
		Aliases:       append([]string(nil), finding.Aliases...),
		Severity:      finding.Severity,
		TargetID:      finding.TargetID,
		WorkspaceID:   finding.WorkspaceID,
		WorkspacePath: finding.WorkspacePath,
	}
}

func requireScanArrays(data []byte) error {
	var object map[string]json.RawMessage
	if err := json.Unmarshal(data, &object); err != nil {
		return fmt.Errorf("scan input is malformed")
	}
	for _, name := range []string{"findings", "errors", "artifact_digests"} {
		value, ok := object[name]
		if !ok || bytes.Equal(bytes.TrimSpace(value), []byte("null")) || len(value) == 0 || value[0] != '[' {
			return fmt.Errorf("scan %s must be an array", name)
		}
	}
	return nil
}

func cloneStrings(values []string) []string {
	if values == nil {
		return nil
	}
	return append(make([]string, 0, len(values)), values...)
}

func scanInputErrorCode(err error) string {
	if errors.Is(err, os.ErrNotExist) {
		return "BASELINE_INPUT_INVALID"
	}
	return "BASELINE_INPUT_INVALID"
}
