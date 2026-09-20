package presenter

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

// Validate checks the event envelope without interpreting terminal output.
func (e Event) Validate() error {
	if !knownEventType(e.Type) {
		return fmt.Errorf("unknown presentation event %q", e.Type)
	}
	if e.WorkspaceCount < 0 || e.CompletedWorkspaces < 0 {
		return fmt.Errorf("presentation event counts cannot be negative")
	}
	if e.CompletedWorkspaces > e.WorkspaceCount && e.WorkspaceCount > 0 {
		return fmt.Errorf("completed workspaces exceed workspace count")
	}
	return nil
}

// EventSink receives lifecycle events without coupling application code to a
// terminal implementation. Emit must preserve event order for a single run.
type EventSink interface {
	Emit(Event) error
}

// NopEventSink discards events when no presenter is requested.
type NopEventSink struct{}

func (NopEventSink) Emit(Event) error { return nil }

func knownEventType(eventType EventType) bool {
	switch eventType {
	case EventInputValidated,
		EventWorkspaceDiscoveryStarted,
		EventWorkspaceDiscovered,
		EventScanPlanBuilt,
		EventWorkspaceScanStarted,
		EventWorkspaceScanCompleted,
		EventArtifactStored,
		EventNormalizationStarted,
		EventReportReady,
		EventOperationCancelled:
		return true
	default:
		return false
	}
}
