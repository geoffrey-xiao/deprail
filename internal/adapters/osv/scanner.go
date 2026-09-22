package osv

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/geoffrey-xiao/deprail/internal/adapter"
	"github.com/geoffrey-xiao/deprail/internal/process"
	"path/filepath"
	"sort"
	"time"
)

type Scanner struct {
	Path      string
	Dir       string
	Args      []string
	Timeout   time.Duration
	OutputCap int64
}

func (s Scanner) Metadata(ctx context.Context) (adapter.Metadata, error) {
	result, err := process.Run(ctx, process.Request{Path: s.Path, Args: []string{"--version"}, Dir: s.Dir, Timeout: s.Timeout, OutputCap: s.OutputCap})
	if err != nil {
		return adapter.Metadata{}, processToAdapterError(err)
	}
	return ParseVersion(string(result.Stdout))
}

func (s Scanner) Compatible(ctx context.Context, metadata adapter.Metadata) error {
	return Compatible(ctx, metadata)
}

func (s Scanner) Plan(ctx context.Context, targets []adapter.Target) (adapter.Plan, error) {
	plan := adapter.Plan{Targets: append([]adapter.Target(nil), targets...)}
	if err := adapter.ValidatePlan(plan); err != nil {
		return adapter.Plan{}, err
	}
	sort.Slice(plan.Targets, func(i, j int) bool { return plan.Targets[i].RelativePath < plan.Targets[j].RelativePath })
	return plan, nil
}

func (s Scanner) Execute(ctx context.Context, plan adapter.Plan) (adapter.RawResult, error) {
	if len(plan.Targets) != 1 {
		return adapter.RawResult{}, &adapter.Error{Code: adapter.ErrInvalidPlan, Message: "OSV execution requires one target per unit"}
	}
	target := plan.Targets[0]
	args := append([]string{}, s.Args...)
	args = append(args, "scan", "source", "--format", "json")
	if lockfile := targetLockfile(target); lockfile != "" {
		args = append(args, "--lockfile", lockfile)
	} else {
		args = append(args, target.RelativePath)
	}
	result, err := process.Run(ctx, process.Request{Path: s.Path, Args: args, Dir: s.Dir, Timeout: s.Timeout, OutputCap: s.OutputCap})
	if err != nil {
		raw := adapter.RawResult{Stdout: result.Stdout, Stderr: result.Stderr, ExitCode: result.ExitCode}
		if result.ExitCode == 1 {
			// OSV-Scanner v2 uses exit code 1 to report vulnerabilities found.
			return raw, nil
		}
		return raw, processToAdapterError(err)
	}
	return adapter.RawResult{Stdout: result.Stdout, Stderr: result.Stderr, ExitCode: result.ExitCode}, nil
}

func targetLockfile(target adapter.Target) string {
	for _, file := range target.PackageFiles {
		switch filepath.Base(filepath.FromSlash(file)) {
		case "package-lock.json", "npm-shrinkwrap.json", "pnpm-lock.yaml", "yarn.lock",
			"requirements.txt", "uv.lock", "poetry.lock", "pom.xml", "gradle.lockfile":
			return filepath.ToSlash(file)
		}
	}
	return ""
}

type rawResult struct {
	Results []rawWorkspace `json:"results"`
}
type rawWorkspace struct {
	Packages []rawPackage `json:"packages"`
}
type rawPackage struct {
	Package         rawPackageInfo     `json:"package"`
	Vulnerabilities []rawVulnerability `json:"vulnerabilities"`
}
type rawPackageInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
type rawVulnerability struct {
	ID               string          `json:"id"`
	Aliases          []string        `json:"aliases"`
	Severity         json.RawMessage `json:"severity"`
	DatabaseSpecific struct {
		FixedVersion string `json:"fixed_version"`
		Severity     string `json:"severity"`
	} `json:"database_specific"`
	Affected []rawAffected `json:"affected"`
}
type rawAffected struct {
	Ranges []rawRange `json:"ranges"`
}
type rawRange struct {
	Events []rawEvent `json:"events"`
}
type rawEvent struct {
	Fixed string `json:"fixed"`
}

func vulnerabilitySeverity(vulnerability rawVulnerability) string {
	if vulnerability.DatabaseSpecific.Severity != "" {
		return vulnerability.DatabaseSpecific.Severity
	}
	var text string
	if json.Unmarshal(vulnerability.Severity, &text) == nil {
		return text
	}
	var observations []struct {
		Score string `json:"score"`
	}
	if json.Unmarshal(vulnerability.Severity, &observations) == nil && len(observations) > 0 {
		return observations[0].Score
	}
	return ""
}

func vulnerabilityFixedVersion(vulnerability rawVulnerability) string {
	if vulnerability.DatabaseSpecific.FixedVersion != "" {
		return vulnerability.DatabaseSpecific.FixedVersion
	}
	versions := make([]string, 0)
	for _, affected := range vulnerability.Affected {
		for _, vulnerabilityRange := range affected.Ranges {
			for _, event := range vulnerabilityRange.Events {
				if event.Fixed != "" {
					versions = append(versions, event.Fixed)
				}
			}
		}
	}
	if len(versions) == 0 {
		return ""
	}
	sort.Strings(versions)
	return versions[0]
}

func (s Scanner) Normalize(_ context.Context, records []adapter.Record) ([]adapter.Finding, error) {
	findings := make([]adapter.Finding, 0, len(records))
	for _, record := range records {
		findings = append(findings, adapter.Finding{Component: record.Component, Version: record.Version, Aliases: append([]string(nil), record.Aliases...), Severity: record.Severity, Fixed: record.Fixed, TargetID: record.TargetID})
	}
	return findings, nil
}

func Parse(raw adapter.RawResult) ([]adapter.Record, error) {
	var document rawResult
	if err := json.Unmarshal(raw.Stdout, &document); err != nil {
		return nil, &adapter.Error{Code: adapter.ErrInvalidOutput, Message: "OSV-Scanner JSON output is invalid"}
	}
	records := make([]adapter.Record, 0)
	for _, workspace := range document.Results {
		for _, pkg := range workspace.Packages {
			for _, vuln := range pkg.Vulnerabilities {
				aliases := append([]string(nil), vuln.Aliases...)
				sort.Strings(aliases)
				records = append(records, adapter.Record{Component: pkg.Package.Name, Version: pkg.Package.Version, Aliases: aliases, Severity: vulnerabilitySeverity(vuln), Fixed: vulnerabilityFixedVersion(vuln), TargetID: vuln.ID})
			}
		}
	}
	sort.Slice(records, func(i, j int) bool {
		return records[i].TargetID+"\x00"+records[i].Component+"\x00"+records[i].Version < records[j].TargetID+"\x00"+records[j].Component+"\x00"+records[j].Version
	})
	return records, nil
}

func (s Scanner) Parse(_ context.Context, raw adapter.RawResult) ([]adapter.Record, error) {
	return Parse(raw)
}

func processToAdapterError(err error) error {
	var processErr *process.Error
	if !errors.As(err, &processErr) {
		return &adapter.Error{Code: adapter.ErrExecution, Message: "scanner execution failed"}
	}
	code := adapter.ErrExecution
	switch processErr.Code {
	case process.ErrNotFound:
		code = adapter.ErrScannerNotFound
	case process.ErrTimeout:
		code = adapter.ErrTimeout
	case process.ErrOutputLimit:
		code = adapter.ErrOutputLimit
	}
	return &adapter.Error{Code: code, Message: processErr.Message, Retryable: processErr.Retryable}
}

var _ adapter.ScannerAdapter = Scanner{}
