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
	"github.com/geoffrey-xiao/deprail/internal/normalize"
	"github.com/geoffrey-xiao/deprail/internal/remediation"
	"github.com/geoffrey-xiao/deprail/internal/scanplan"
)

type ScanOptions struct {
	Scanner   adapter.ScannerAdapter
	Artifacts artifact.Store
}
type ScanReport struct {
	SchemaVersion      string                         `json:"schema_version"`
	DocumentType       string                         `json:"document_type"`
	ScanID             string                         `json:"scan_id"`
	RepositoryIdentity remediation.RepositoryIdentity `json:"repository_identity"`
	RepositoryState    string                         `json:"repository_state"`
	Status             discovery.Completeness         `json:"status"`
	Findings           []adapter.Finding              `json:"findings"`
	Errors             []string                       `json:"errors"`
	ArtifactDigests    []string                       `json:"artifact_digests"`
}

func Scan(ctx context.Context, root string, options ScanOptions) (ScanReport, error) {
	graph, err := Discover(ctx, root, DiscoverOptions{})
	if err != nil {
		return ScanReport{Status: discovery.Failed, Findings: []adapter.Finding{}, Errors: []string{err.Error()}, ArtifactDigests: []string{}}, err
	}
	scanRoot, err := filepath.Abs(root)
	if err != nil {
		return ScanReport{Status: discovery.Failed, Findings: []adapter.Finding{}, Errors: []string{"resolve scanner root"}, ArtifactDigests: []string{}}, err
	}
	scanRoot, err = filepath.EvalSymlinks(scanRoot)
	if err != nil {
		return ScanReport{Status: discovery.Failed, Findings: []adapter.Finding{}, Errors: []string{"resolve scanner root"}, ArtifactDigests: []string{}}, err
	}
	repositoryState, stateErr := CurrentRepositoryState(scanRoot)
	if stateErr != nil {
		return ScanReport{Status: discovery.Failed, Findings: []adapter.Finding{}, Errors: []string{"resolve repository state"}, ArtifactDigests: []string{}}, stateErr
	}
	scanID := "scan-" + repositoryState[:16]
	report := ScanReport{SchemaVersion: remediation.ReportSchemaVersion, DocumentType: remediation.ReportDocumentType, ScanID: scanID, RepositoryIdentity: remediation.RepositoryIdentity{Root: scanRoot, Repository: filepath.Base(scanRoot)}, RepositoryState: repositoryState, Status: graph.Completeness, Findings: []adapter.Finding{}, Errors: []string{}, ArtifactDigests: []string{}}
	if options.Scanner == nil {
		options.Scanner = osv.Scanner{Path: "osv-scanner", Dir: scanRoot, Timeout: 2 * time.Minute, OutputCap: 16 << 20}
	}
	if options.Artifacts.Root == "" {
		options.Artifacts = artifact.Store{Root: filepath.Join(scanRoot, ".deprail", "artifacts"), MaxBytes: 16 << 20}
	}
	units, err := scanplan.Build(graph)
	if err != nil {
		return ScanReport{Status: discovery.Failed, Findings: []adapter.Finding{}, Errors: []string{err.Error()}, ArtifactDigests: []string{}}, err
	}
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
		validFindings := make([]adapter.Finding, 0, len(findings))
		for i := range findings {
			component, componentErr := normalize.NormalizeComponent(normalize.ComponentInput{
				Name: findings[i].Component, Version: findings[i].Version,
				Ecosystem: purlEcosystem(unit.Target.Ecosystem), WorkspaceID: unit.Target.WorkspaceID,
			})
			if componentErr != nil {
				report.Errors = append(report.Errors, componentErr.Error())
				report.Status = discovery.Partial
				continue
			}
			findings[i].PURL = component.PURL
			findings[i].WorkspaceID = unit.Target.WorkspaceID
			findings[i].WorkspacePath = unit.Target.RelativePath
			findings[i].Ecosystem = unit.Target.Ecosystem
			validFindings = append(validFindings, findings[i])
		}
		report.Findings = append(report.Findings, validFindings...)
	}
	if len(report.Findings) == 0 && len(report.Errors) > 0 {
		report.Status = discovery.Failed
	}
	return report, nil
}

func ScanTerminal(report ScanReport) string {
	return fmt.Sprintf("Status: %s\nFindings: %d\nErrors: %d\n", report.Status, len(report.Findings), len(report.Errors))
}

func purlEcosystem(ecosystem string) string {
	switch ecosystem {
	case "javascript":
		return "npm"
	case "python":
		return "pip"
	case "java":
		return "maven"
	default:
		return ecosystem
	}
}
