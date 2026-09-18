package app

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/geoffrey-xiao/deprail/internal/adapter"
	"github.com/geoffrey-xiao/deprail/internal/adapters/osv"
	"github.com/geoffrey-xiao/deprail/internal/artifact"
	"github.com/geoffrey-xiao/deprail/internal/discovery"
	"github.com/geoffrey-xiao/deprail/internal/scanplan"
)

type ScanOptions struct {
	Scanner   adapter.ScannerAdapter
	Artifacts artifact.Store
}
type ScanReport struct {
	Status          discovery.Completeness `json:"status"`
	Findings        []adapter.Finding      `json:"findings"`
	Errors          []string               `json:"errors"`
	ArtifactDigests []string               `json:"artifact_digests"`
}

func Scan(ctx context.Context, root string, options ScanOptions) (ScanReport, error) {
	graph, err := Discover(ctx, root, DiscoverOptions{})
	if err != nil {
		return ScanReport{Status: discovery.Failed, Errors: []string{err.Error()}}, err
	}
	scanRoot, err := filepath.Abs(root)
	if err != nil {
		return ScanReport{Status: discovery.Failed, Errors: []string{"resolve scanner root"}}, err
	}
	scanRoot, err = filepath.EvalSymlinks(scanRoot)
	if err != nil {
		return ScanReport{Status: discovery.Failed, Errors: []string{"resolve scanner root"}}, err
	}
	if options.Scanner == nil {
		options.Scanner = osv.Scanner{Path: "osv-scanner", Dir: scanRoot, Timeout: 2 * time.Minute, OutputCap: 16 << 20}
	}
	if options.Artifacts.Root == "" {
		options.Artifacts = artifact.Store{Root: filepath.Join(scanRoot, ".deprail", "artifacts"), MaxBytes: 16 << 20}
	}
	units, err := scanplan.Build(graph)
	if err != nil {
		return ScanReport{Status: discovery.Failed, Errors: []string{err.Error()}}, err
	}
	report := ScanReport{Status: graph.Completeness}
	for _, unit := range units {
		plan, planErr := options.Scanner.Plan(ctx, []adapter.Target{unit.Target})
		if planErr != nil {
			report.Errors = append(report.Errors, planErr.Error())
			report.Status = discovery.Partial
			continue
		}
		raw, execErr := options.Scanner.Execute(ctx, plan)
		if len(raw.Stdout) > 0 {
			stored, storeErr := options.Artifacts.Put(raw.Stdout)
			if storeErr != nil {
				report.Errors = append(report.Errors, storeErr.Error())
				report.Status = discovery.Partial
			} else {
				report.ArtifactDigests = append(report.ArtifactDigests, stored.Digest)
			}
		}
		if execErr != nil {
			report.Errors = append(report.Errors, execErr.Error())
			report.Status = discovery.Partial
			continue
		}
		records, parseErr := options.Scanner.Parse(ctx, raw)
		if parseErr != nil {
			report.Errors = append(report.Errors, parseErr.Error())
			report.Status = discovery.Partial
			continue
		}
		findings, normErr := options.Scanner.Normalize(ctx, records)
		if normErr != nil {
			report.Errors = append(report.Errors, normErr.Error())
			report.Status = discovery.Partial
			continue
		}
		report.Findings = append(report.Findings, findings...)
	}
	if len(report.Findings) == 0 && len(report.Errors) > 0 {
		report.Status = discovery.Failed
	}
	return report, nil
}

func ScanTerminal(report ScanReport) string {
	return fmt.Sprintf("Status: %s\nFindings: %d\nErrors: %d\n", report.Status, len(report.Findings), len(report.Errors))
}
