package app

import "fmt"

// EventType identifies an application lifecycle transition exposed to presenters.
type EventType string

const (
	EventInputValidated            EventType = "InputValidated"
	EventWorkspaceDiscoveryStarted EventType = "WorkspaceDiscoveryStarted"
	EventWorkspaceDiscovered       EventType = "WorkspaceDiscovered"
	EventScanPlanBuilt             EventType = "ScanPlanBuilt"
	EventWorkspaceScanStarted      EventType = "WorkspaceScanStarted"
	EventWorkspaceScanCompleted    EventType = "WorkspaceScanCompleted"
	EventArtifactStored            EventType = "ArtifactStored"
	EventNormalizationStarted      EventType = "NormalizationStarted"
	EventReportReady               EventType = "ReportReady"
	EventOperationCancelled        EventType = "OperationCancelled"
)

// Event is renderer-independent lifecycle data. Fields are populated only when
// relevant to the event type; renderers decide how to present the data.
type Event struct {
	Type                EventType
	WorkspaceID         string
	WorkspacePath       string
	WorkspaceCount      int
	CompletedWorkspaces int
	Status              string
	ErrorCode           string
}

// EventSink receives lifecycle events without coupling application code to a
// terminal implementation. Emit preserves event order for a single run.
type EventSink interface {
	Emit(Event) error
}

// EventValidator enforces lifecycle ordering and workspace identity invariants.
type EventValidator struct {
	phase      eventPhase
	workspaces map[string]workspacePhase
}

type eventPhase uint8

const (
	phaseInitial eventPhase = iota
	phaseDiscovery
	phasePlanned
	phaseScanning
	phaseNormalizing
	phaseReady
	phaseCancelled
)

type workspacePhase uint8

const (
	workspaceDiscovered workspacePhase = iota + 1
	workspaceScanning
)

// NewEventValidator returns a validator for one application run.
func NewEventValidator() *EventValidator {
	return &EventValidator{workspaces: make(map[string]workspacePhase)}
}

// Emit validates and accepts one event, rejecting invalid lifecycle transitions.
func (v *EventValidator) Emit(event Event) error {
	if v == nil {
		return fmt.Errorf("nil presentation event validator")
	}
	if err := event.validateEnvelope(); err != nil {
		return err
	}
	if v.phase == phaseReady || v.phase == phaseCancelled {
		return fmt.Errorf("event %q follows terminal event", event.Type)
	}

	switch event.Type {
	case EventInputValidated:
		if v.phase != phaseInitial {
			return fmt.Errorf("input validation must be the first event")
		}
		v.phase = phaseDiscovery
	case EventWorkspaceDiscoveryStarted:
		if v.phase != phaseDiscovery {
			return fmt.Errorf("workspace discovery started out of order")
		}
	case EventWorkspaceDiscovered:
		if v.phase != phaseDiscovery {
			return fmt.Errorf("workspace discovered out of order")
		}
		v.workspaces[event.WorkspaceID] = workspaceDiscovered
	case EventScanPlanBuilt:
		if v.phase != phaseDiscovery || len(v.workspaces) == 0 {
			return fmt.Errorf("scan plan built before workspace discovery")
		}
		v.phase = phasePlanned
	case EventWorkspaceScanStarted:
		if v.phase != phasePlanned && v.phase != phaseScanning {
			return fmt.Errorf("workspace scan started out of order")
		}
		if v.workspaces[event.WorkspaceID] != workspaceDiscovered {
			return fmt.Errorf("workspace %q was not discovered or is already scanning", event.WorkspaceID)
		}
		v.workspaces[event.WorkspaceID] = workspaceScanning
		v.phase = phaseScanning
	case EventWorkspaceScanCompleted:
		if v.phase != phaseScanning || v.workspaces[event.WorkspaceID] != workspaceScanning {
			return fmt.Errorf("workspace scan completed out of order")
		}
		v.workspaces[event.WorkspaceID] = workspaceDiscovered
	case EventArtifactStored:
		if v.phase != phaseScanning && v.phase != phasePlanned {
			return fmt.Errorf("artifact stored out of order")
		}
	case EventNormalizationStarted:
		if v.phase != phaseScanning || hasScanningWorkspace(v.workspaces) {
			return fmt.Errorf("normalization started before workspace scans completed")
		}
		v.phase = phaseNormalizing
	case EventReportReady:
		if v.phase != phaseNormalizing {
			return fmt.Errorf("report ready out of order")
		}
		v.phase = phaseReady
	case EventOperationCancelled:
		if v.phase == phaseInitial {
			return fmt.Errorf("operation cancelled before input validation")
		}
		v.phase = phaseCancelled
	default:
		return fmt.Errorf("unknown presentation event %q", event.Type)
	}
	return nil
}

func (e Event) validateEnvelope() error {
	if !knownEventType(e.Type) {
		return fmt.Errorf("unknown presentation event %q", e.Type)
	}
	if e.WorkspaceCount < 0 || e.CompletedWorkspaces < 0 {
		return fmt.Errorf("presentation event counts cannot be negative")
	}
	if e.CompletedWorkspaces > e.WorkspaceCount && e.WorkspaceCount > 0 {
		return fmt.Errorf("completed workspaces exceed workspace count")
	}
	if requiresWorkspaceIdentity(e.Type) && (e.WorkspaceID == "" || e.WorkspacePath == "") {
		return fmt.Errorf("event %q requires workspace identity", e.Type)
	}
	return nil
}

func requiresWorkspaceIdentity(eventType EventType) bool {
	switch eventType {
	case EventWorkspaceDiscovered, EventWorkspaceScanStarted, EventWorkspaceScanCompleted, EventArtifactStored:
		return true
	default:
		return false
	}
}

func hasScanningWorkspace(workspaces map[string]workspacePhase) bool {
	for _, phase := range workspaces {
		if phase == workspaceScanning {
			return true
		}
	}
	return false
}

func knownEventType(eventType EventType) bool {
	switch eventType {
	case EventInputValidated, EventWorkspaceDiscoveryStarted, EventWorkspaceDiscovered,
		EventScanPlanBuilt, EventWorkspaceScanStarted, EventWorkspaceScanCompleted,
		EventArtifactStored, EventNormalizationStarted, EventReportReady, EventOperationCancelled:
		return true
	default:
		return false
	}
}

// NopEventSink discards events when no presenter is requested.
type NopEventSink struct{}

func (NopEventSink) Emit(Event) error { return nil }
