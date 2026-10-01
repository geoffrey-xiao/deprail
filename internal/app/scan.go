package app

import (
	"context"
	"errors"
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
	Scanner       adapter.ScannerAdapter
	Artifacts     artifact.Store
	Events        EventSink
	Capture       CaptureSink
	CaptureCancel context.Context
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
	if options.Capture == nil {
		return scan(ctx, root, options, nil)
	}
	id, idErr := randomHistoryEntryID()
	state := &scanCaptureState{
		recordedAtUS: time.Now().UTC().UnixMicro(),
		unclassified: idErr != nil,
	}
	report, operationErr := scan(ctx, root, options, state)
	var reportRef *ScanReport
	if state.reportAvailable {
		reportRef = &report
	}
	input := CaptureInput{
		HistoryEntryID: id, RecordedAtUS: state.recordedAtUS, Graph: state.graph,
		Report: reportRef, OperationError: operationErr,
		Diagnostics: append([]normalize.HistoryDiagnosticInput{}, state.diagnostics...),
		Provenance:  state.provenance, UnclassifiedSource: state.unclassified,
	}
	captureCtx, cancel := context.WithTimeout(context.Background(), captureTimeout)
	defer cancel()
	if options.CaptureCancel != nil {
		stop := context.AfterFunc(options.CaptureCancel, cancel)
		defer stop()
	}
	result, captureErr := options.Capture.Capture(captureCtx, input)
	if captureErr == nil && (!result.Saved || result.HistoryEntryID != id) {
		captureErr = errors.New("capture sink did not confirm the occurrence")
	}
	if captureErr != nil {
		return report, &ScanCaptureError{OperationErr: operationErr, PersistenceErr: stableHistoryError(captureErr)}
	}
	return report, operationErr
}

func scan(ctx context.Context, root string, options ScanOptions, state *scanCaptureState) (ScanReport, error) {
	if options.Events == nil {
		options.Events = NopEventSink{}
	}
	if err := emitScanEvent(state, options.Events, Event{Type: EventInputValidated}); err != nil {
		return ScanReport{}, err
	}
	if err := emitScanEvent(state, options.Events, Event{Type: EventWorkspaceDiscoveryStarted}); err != nil {
		return ScanReport{}, err
	}
	graph, err := Discover(ctx, root, DiscoverOptions{})
	if err != nil {
		state.addBoundaryError(err, "discover")
		return ScanReport{Status: discovery.Failed, Findings: []adapter.Finding{}, Errors: []string{err.Error()}, ArtifactDigests: []string{}}, err
	}
	if state != nil {
		state.graph = &graph
		state.addGraph(graph)
	}
	for _, workspace := range graph.Workspaces {
		if err := emitScanEvent(state, options.Events, Event{Type: EventWorkspaceDiscovered, WorkspaceID: workspace.WorkspaceID, WorkspacePath: workspace.RelativePath}); err != nil {
			return ScanReport{}, err
		}
	}
	scanRoot, err := filepath.Abs(root)
	if err != nil {
		state.markUnclassified()
		return ScanReport{Status: discovery.Failed, Findings: []adapter.Finding{}, Errors: []string{"resolve scanner root"}, ArtifactDigests: []string{}}, err
	}
	scanRoot, err = filepath.EvalSymlinks(scanRoot)
	if err != nil {
		state.markUnclassified()
		return ScanReport{Status: discovery.Failed, Findings: []adapter.Finding{}, Errors: []string{"resolve scanner root"}, ArtifactDigests: []string{}}, err
	}
	repositoryState, stateErr := CurrentRepositoryState(scanRoot)
	if stateErr != nil {
		state.markUnclassified()
		return ScanReport{Status: discovery.Failed, Findings: []adapter.Finding{}, Errors: []string{"resolve repository state"}, ArtifactDigests: []string{}}, stateErr
	}
	scanID := "scan-" + repositoryState[:16]
	report := ScanReport{
		SchemaVersion: remediation.ReportSchemaVersion, DocumentType: remediation.ReportDocumentType,
		ScanID: scanID, RepositoryIdentity: remediation.RepositoryIdentity{Root: scanRoot, Repository: filepath.Base(scanRoot)},
		RepositoryState: repositoryState, Status: graph.Completeness, Findings: []adapter.Finding{},
		Errors: []string{}, ArtifactDigests: []string{},
	}
	if state != nil {
		state.reportAvailable = true
	}
	if options.Scanner == nil {
		options.Scanner = osv.Scanner{Path: "osv-scanner", Dir: scanRoot, Timeout: 2 * time.Minute, OutputCap: 16 << 20}
	}
	if state != nil {
		if metadata, metadataErr := options.Scanner.Metadata(ctx); metadataErr == nil {
			if metadata.Name != "" {
				state.provenance.ScannerName = new(metadata.Name)
			}
			if metadata.Version != "" {
				state.provenance.ScannerVersion = new(metadata.Version)
			}
		}
	}
	if options.Artifacts.Root == "" {
		options.Artifacts = artifact.Store{Root: filepath.Join(scanRoot, ".deprail", "artifacts"), MaxBytes: 16 << 20}
	}
	units, err := scanplan.Build(graph)
	if err != nil {
		state.addBoundaryError(err, "plan")
		return ScanReport{Status: discovery.Failed, Findings: []adapter.Finding{}, Errors: []string{err.Error()}, ArtifactDigests: []string{}}, err
	}
	if err := emitScanEvent(state, options.Events, Event{Type: EventScanPlanBuilt, WorkspaceCount: len(units)}); err != nil {
		return ScanReport{}, err
	}
	normalizationStarted := false
	for _, unit := range units {
		if err := emitScanEvent(state, options.Events, Event{Type: EventWorkspaceScanStarted, WorkspaceID: unit.Target.WorkspaceID, WorkspacePath: unit.Target.RelativePath}); err != nil {
			return ScanReport{}, err
		}
		completeWorkspace := func(status string) error {
			return emitScanEvent(state, options.Events, Event{Type: EventWorkspaceScanCompleted, WorkspaceID: unit.Target.WorkspaceID, WorkspacePath: unit.Target.RelativePath, Status: status})
		}
		plan, planErr := options.Scanner.Plan(ctx, []adapter.Target{unit.Target})
		if planErr != nil {
			state.addWorkspaceError(planErr, "plan", unit.Target.WorkspaceID)
			if errors.Is(planErr, context.Canceled) {
				_ = emitScanEvent(state, options.Events, Event{Type: EventOperationCancelled, ErrorCode: "CANCELLED"})
				return report, planErr
			}
			report.Errors = append(report.Errors, planErr.Error())
			report.Status = discovery.Partial
			if err := completeWorkspace("failed"); err != nil {
				return report, err
			}
			continue
		}
		raw, execErr := options.Scanner.Execute(ctx, plan)
		if len(raw.Stdout) > 0 {
			stored, storeErr := options.Artifacts.Put(raw.Stdout)
			if storeErr != nil {
				state.addWorkspaceError(storeErr, "artifact", unit.Target.WorkspaceID)
				report.Errors = append(report.Errors, storeErr.Error())
				report.Status = discovery.Partial
			} else {
				report.ArtifactDigests = append(report.ArtifactDigests, stored.Digest)
				if err := emitScanEvent(state, options.Events, Event{Type: EventArtifactStored, WorkspaceID: unit.Target.WorkspaceID, WorkspacePath: unit.Target.RelativePath, Status: "stored"}); err != nil {
					return ScanReport{}, err
				}
			}
		}
		if execErr != nil {
			state.addWorkspaceError(execErr, "execute", unit.Target.WorkspaceID)
			if errors.Is(execErr, context.Canceled) {
				_ = emitScanEvent(state, options.Events, Event{Type: EventOperationCancelled, ErrorCode: "CANCELLED"})
				return report, execErr
			}
			report.Errors = append(report.Errors, execErr.Error())
			report.Status = discovery.Partial
			if err := completeWorkspace("failed"); err != nil {
				return report, err
			}
			continue
		}
		records, parseErr := options.Scanner.Parse(ctx, raw)
		if parseErr != nil {
			state.addWorkspaceError(parseErr, "parse", unit.Target.WorkspaceID)
			if errors.Is(parseErr, context.Canceled) {
				_ = emitScanEvent(state, options.Events, Event{Type: EventOperationCancelled, ErrorCode: "CANCELLED"})
				return report, parseErr
			}
			report.Errors = append(report.Errors, parseErr.Error())
			report.Status = discovery.Partial
			if err := completeWorkspace("failed"); err != nil {
				return report, err
			}
			continue
		}
		if !normalizationStarted {
			if err := emitScanEvent(state, options.Events, Event{Type: EventNormalizationStarted}); err != nil {
				return report, err
			}
			normalizationStarted = true
		}
		findings, normErr := options.Scanner.Normalize(ctx, records)
		if normErr != nil {
			state.addWorkspaceError(normErr, "normalize", unit.Target.WorkspaceID)
			if errors.Is(normErr, context.Canceled) {
				_ = emitScanEvent(state, options.Events, Event{Type: EventOperationCancelled, ErrorCode: "CANCELLED"})
				return report, normErr
			}
			report.Errors = append(report.Errors, normErr.Error())
			report.Status = discovery.Partial
			if err := completeWorkspace("failed"); err != nil {
				return report, err
			}
			continue
		}
		validFindings := make([]adapter.Finding, 0, len(findings))
		for i := range findings {
			component, componentErr := normalize.NormalizeComponent(normalize.ComponentInput{
				Name: findings[i].Component, Version: findings[i].Version,
				Ecosystem: purlEcosystem(unit.Target.Ecosystem), WorkspaceID: unit.Target.WorkspaceID,
			})
			if componentErr != nil {
				state.addWorkspaceError(componentErr, "component", unit.Target.WorkspaceID)
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
		if err := emitScanEvent(state, options.Events, Event{Type: EventWorkspaceScanCompleted, WorkspaceID: unit.Target.WorkspaceID, WorkspacePath: unit.Target.RelativePath}); err != nil {
			return ScanReport{}, err
		}
	}
	if !normalizationStarted {
		if err := emitScanEvent(state, options.Events, Event{Type: EventNormalizationStarted}); err != nil {
			return ScanReport{}, err
		}
	}
	if len(report.Findings) == 0 && len(report.Errors) > 0 {
		report.Status = discovery.Failed
	}
	if err := emitScanEvent(state, options.Events, Event{Type: EventReportReady, Status: string(report.Status), WorkspaceCount: len(units), CompletedWorkspaces: len(units)}); err != nil {
		return ScanReport{}, err
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
