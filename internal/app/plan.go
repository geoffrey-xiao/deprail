package app

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/geoffrey-xiao/deprail/internal/remediation"
	"github.com/geoffrey-xiao/deprail/internal/remediation/java"
	"github.com/geoffrey-xiao/deprail/internal/remediation/javascript"
	"github.com/geoffrey-xiao/deprail/internal/remediation/python"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrPlanUnsupported = errors.New("no planning adapter supports the finding component")
	ErrPlanIncomplete  = errors.New("planning requires a complete report")
	ErrPlanRejected    = errors.New("planning input was rejected by the adapter")
	ErrPlanUnknown     = errors.New("planning result is unknown")
)

type PlanOptions struct {
	CurrentRepositoryState string
	Adapters               []remediation.PlanningAdapter
}

func Plan(ctx context.Context, reportPath, findingKey, repositoryRoot string, options PlanOptions) (remediation.Plan, error) {
	if reportPath == "" {
		return remediation.Plan{}, &remediation.ReportError{Code: remediation.ReportRequired, Message: "an explicit normalized scan report is required"}
	}
	report, err := loadPlanningReport(reportPath)
	if err != nil {
		return remediation.Plan{}, err
	}
	root := repositoryRoot
	if root == "" {
		root = report.RepositoryIdentity.Root
	}
	root, err = filepath.Abs(root)
	if err != nil {
		return remediation.Plan{}, fmt.Errorf("resolve repository root: %w", err)
	}
	currentState := options.CurrentRepositoryState
	if currentState == "" {
		currentState, err = CurrentRepositoryState(root)
		if err != nil {
			return remediation.Plan{}, &remediation.ReportError{Code: remediation.ReportInputStale, Message: "current repository state is unavailable"}
		}
	}
	resolved, err := remediation.ResolveFinding(report, findingKey, currentState)
	if err != nil {
		return remediation.Plan{}, err
	}
	if resolved.Report.Status != remediation.ReportComplete {
		return remediation.Plan{}, ErrPlanIncomplete
	}
	adapters := options.Adapters
	if len(adapters) == 0 {
		adapters = []remediation.PlanningAdapter{javascript.Adapter{}, python.Adapter{}, java.Adapter{}}
	}
	request := remediation.PlanningRequest{
		Repository:      resolved.Report.RepositoryIdentity,
		Workspace:       resolved.Finding.Workspace,
		Finding:         remediation.FindingIdentity{StableKey: resolved.Finding.StableKey, Aliases: resolved.Finding.Aliases, DependencyPath: resolved.Finding.DependencyPath, CurrentVersion: resolved.Finding.CurrentVersion, FixedVersions: resolved.Finding.FixedVersions},
		Component:       resolved.Finding.Component,
		RepositoryState: resolved.Report.RepositoryState,
	}
	request.Repository.Root = root
	var adapter remediation.PlanningAdapter
	for _, candidate := range adapters {
		if supportedPURL(candidate.Name(), request.Component.PURL) {
			adapter = candidate
			break
		}
	}
	if adapter == nil {
		return remediation.Plan{}, ErrPlanUnsupported
	}
	assessment, err := adapter.Assess(ctx, request)
	if err != nil {
		return remediation.Plan{}, err
	}
	if assessment.State == remediation.AdapterRejected {
		return remediation.Plan{}, fmt.Errorf("%w: %s", ErrPlanRejected, assessment.Reason)
	}
	if assessment.State == remediation.AdapterUnknown {
		return remediation.Plan{}, fmt.Errorf("%w: %s", ErrPlanUnknown, assessment.Reason)
	}
	if assessment.State == remediation.AdapterUnsupported || assessment.State == remediation.AdapterUnavailable {
		return remediation.Plan{}, fmt.Errorf("%w: %s", ErrPlanUnsupported, assessment.Reason)
	}
	evidence, err := adapter.Plan(ctx, request)
	if err != nil {
		return remediation.Plan{}, err
	}
	if err := remediation.ValidatePlanningEvidence(evidence); err != nil {
		return remediation.Plan{}, err
	}
	canonicalReport := report
	remediation.CanonicalizeReport(&canonicalReport)
	canonicalData, err := json.Marshal(canonicalReport)
	if err != nil {
		return remediation.Plan{}, fmt.Errorf("canonicalize normalized scan report: %w", err)
	}
	reportDigest := sha256.Sum256(canonicalData)
	plan := remediation.Plan{SchemaVersion: remediation.SchemaVersion, CreatedFrom: remediation.CreatedFrom{ReportDigest: hex.EncodeToString(reportDigest[:]), SourceScanID: resolved.Report.SourceScanID, Scanner: "osv-scanner", RepositoryState: resolved.Report.RepositoryState}, RepositoryIdentity: request.Repository, WorkspaceIdentity: request.Workspace, FindingIdentity: request.Finding, Component: request.Component, CurrentState: remediation.CurrentState{Direct: directCandidate(evidence), Manifest: firstFile(evidence, "manifest"), Lockfile: firstFile(evidence, "lockfile"), Vulnerable: true}, Candidates: evidence.Candidates, AffectedFiles: evidence.AffectedFiles, Commands: evidence.Commands, Risks: evidence.Risks, Assumptions: evidence.Assumptions, Verification: evidence.Verification, Rollback: remediation.Rollback{Steps: []string{"restore manifest and lockfile changes"}}, Provenance: remediation.Provenance{ArtifactDigests: resolved.Report.ArtifactDigests, Sources: []string{"normalized-scan-report", adapter.Name()}}}
	plan.Canonicalize()
	plan.PlanID = remediation.StablePlanID(plan)
	if err := plan.Validate(); err != nil {
		return remediation.Plan{}, err
	}
	return plan, nil
}

func supportedPURL(name, purl string) bool {
	return (name == "javascript" && strings.HasPrefix(purl, "pkg:npm/")) || (name == "python" && strings.HasPrefix(purl, "pkg:pypi/")) || (name == "java" && strings.HasPrefix(purl, "pkg:maven/"))
}
func firstFile(e remediation.PlanningEvidence, kind string) string {
	for _, file := range e.AffectedFiles {
		if file.Kind == kind {
			return file.Path
		}
	}
	return ""
}
func directCandidate(e remediation.PlanningEvidence) bool {
	for _, candidate := range e.Candidates {
		return candidate.Direct
	}
	return false
}
func loadPlanningReport(path string) (remediation.Report, error) {
	report, err := remediation.LoadReport(path)
	if err == nil {
		return report, nil
	}
	data, readErr := os.ReadFile(path)
	if readErr != nil {
		return remediation.Report{}, err
	}
	var scan ScanReport
	if json.Unmarshal(data, &scan) != nil || scan.SchemaVersion == "" || scan.DocumentType == "" {
		return remediation.Report{}, err
	}
	report = remediation.Report{
		SchemaVersion:      remediation.ReportSchemaVersion,
		DocumentType:       remediation.ReportDocumentType,
		ReportID:           scan.ScanID,
		SourceScanID:       scan.ScanID,
		RepositoryIdentity: scan.RepositoryIdentity,
		ArtifactDigests:    append([]string(nil), scan.ArtifactDigests...),
		RepositoryState:    scan.RepositoryState,
		Status:             string(scan.Status),
		Findings:           make([]remediation.ReportFinding, 0, len(scan.Findings)),
	}
	for _, finding := range scan.Findings {
		key := finding.TargetID
		if key == "" {
			key = finding.Component + "@" + finding.Version
		}
		workspace := remediation.WorkspaceIdentity{ID: finding.WorkspaceID, Path: finding.WorkspacePath}
		if workspace.ID == "" {
			workspace.ID = "root"
		}
		if workspace.Path == "" {
			workspace.Path = "."
		}
		purl := finding.PURL
		if purl == "" {
			purl = finding.Component
		}
		report.Findings = append(report.Findings, remediation.ReportFinding{
			StableKey: key,
			Workspace: workspace,
			Component: remediation.Component{PURL: purl, Name: finding.Component, Version: finding.Version},
			Aliases:   finding.Aliases, CurrentVersion: finding.Version,
			FixedVersions: nonEmptyFixed(finding.Fixed),
			Provenance:    remediation.Provenance{ArtifactDigests: scan.ArtifactDigests, Sources: []string{"scan-json"}},
		})
	}
	if err := remediation.ValidateReport(report); err != nil {
		return remediation.Report{}, err
	}
	return report, nil
}

func nonEmptyFixed(value string) []string {
	if value == "" {
		return []string{}
	}
	return []string{value}
}
