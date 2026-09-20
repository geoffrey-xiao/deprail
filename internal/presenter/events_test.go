package presenter

import "testing"

func TestEventContractRecognizesLifecycleVocabulary(t *testing.T) {
	events := []EventType{
		EventInputValidated,
		EventWorkspaceDiscoveryStarted,
		EventWorkspaceDiscovered,
		EventScanPlanBuilt,
		EventWorkspaceScanStarted,
		EventWorkspaceScanCompleted,
		EventArtifactStored,
		EventNormalizationStarted,
		EventReportReady,
		EventOperationCancelled,
	}
	for _, eventType := range events {
		if err := (Event{Type: eventType}).Validate(); err != nil {
			t.Fatalf("event %q: %v", eventType, err)
		}
	}
}

func TestEventValidateRejectsInvalidEnvelope(t *testing.T) {
	cases := []Event{
		{Type: EventType("Unknown")},
		{Type: EventReportReady, WorkspaceCount: -1},
		{Type: EventReportReady, WorkspaceCount: 1, CompletedWorkspaces: 2},
	}
	for _, event := range cases {
		if err := event.Validate(); err == nil {
			t.Fatalf("event %#v was accepted", event)
		}
	}
}

func TestNopEventSinkAcceptsValidatedAndUnvalidatedEvents(t *testing.T) {
	var sink NopEventSink
	if err := sink.Emit(Event{Type: EventReportReady}); err != nil {
		t.Fatalf("validated event: %v", err)
	}
	if err := sink.Emit(Event{Type: EventType("future")}); err != nil {
		t.Fatalf("sink should not interpret events: %v", err)
	}
}
