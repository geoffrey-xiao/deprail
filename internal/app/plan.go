package app

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/geoffrey-xiao/deprail/internal/remediation"
	"github.com/geoffrey-xiao/deprail/internal/remediation/java"
	"github.com/geoffrey-xiao/deprail/internal/remediation/javascript"
	"github.com/geoffrey-xiao/deprail/internal/remediation/python"
)

var (
	ErrPlanUnsupported = errors.New("no planning adapter supports the finding component")
	ErrPlanIncomplete  = errors.New("planning requires a complete report")
)

type PlanOptions struct {
	CurrentRepositoryState string
	Adapters               []remediation.PlanningAdapter
}

func Plan(ctx context.Context, reportPath, findingKey, repositoryRoot string, options PlanOptions) (remediation.Plan, error) {
	if reportPath == "" {
		return remediation.Plan{}, &remediation.ReportError{Code: remediation.ReportRequired, Message: "an explicit normalized scan report is required"}
	}
	report, err := remediation.LoadReport(reportPath)
	if err != nil {
		return remediation.Plan{}, err
	}
	data, err := os.ReadFile(reportPath)
	if err != nil {
		return remediation.Plan{}, &remediation.ReportError{Code: remediation.ReportInvalid, Message: "normalized scan report is unreadable"}
	}
	currentState := options.CurrentRepositoryState
	if currentState == "" {
		currentState = report.RepositoryState
	}
	resolved, err := remediation.ResolveFinding(report, findingKey, currentState)
	if err != nil {
		return remediation.Plan{}, err
	}
	if resolved.Report.Status != remediation.ReportComplete {
		return remediation.Plan{}, ErrPlanIncomplete
	}
	root := repositoryRoot
	if root == "" {
		root = resolved.Report.RepositoryIdentity.Root
	}
	root, err = filepath.Abs(root)
	if err != nil {
		return remediation.Plan{}, fmt.Errorf("resolve repository root: %w", err)
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
	reportDigest := sha256.Sum256(data)
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
