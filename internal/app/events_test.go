package app

import "testing"

func TestEventValidatorAcceptsCompleteLifecycleAndRepeatedWorkspace(t *testing.T) {
	validator := NewEventValidator()
	events := []Event{
		{Type: EventInputValidated},
		{Type: EventWorkspaceDiscoveryStarted},
		{Type: EventWorkspaceDiscovered, WorkspaceID: "workspace-1", WorkspacePath: "apps/one"},
		{Type: EventScanPlanBuilt},
		{Type: EventWorkspaceScanStarted, WorkspaceID: "workspace-1", WorkspacePath: "apps/one"},
		{Type: EventWorkspaceScanCompleted, WorkspaceID: "workspace-1", WorkspacePath: "apps/one"},
		{Type: EventWorkspaceScanStarted, WorkspaceID: "workspace-1", WorkspacePath: "apps/one"},
		{Type: EventWorkspaceScanCompleted, WorkspaceID: "workspace-1", WorkspacePath: "apps/one"},
		{Type: EventNormalizationStarted},
		{Type: EventReportReady, Status: "complete"},
	}
	for _, event := range events {
		if err := validator.Emit(event); err != nil {
			t.Fatalf("event %q: %v", event.Type, err)
		}
	}
}

func TestEventValidatorRejectsInvalidOrderingAndIdentity(t *testing.T) {
	cases := []struct {
		name   string
		events []Event
	}{
		{
			name: "completion before start",
			events: []Event{
				{Type: EventInputValidated},
				{Type: EventWorkspaceDiscoveryStarted},
				{Type: EventWorkspaceDiscovered, WorkspaceID: "workspace-1", WorkspacePath: "apps/one"},
				{Type: EventScanPlanBuilt},
				{Type: EventWorkspaceScanCompleted, WorkspaceID: "workspace-1", WorkspacePath: "apps/one"},
			},
		},
		{
			name: "empty workspace identity",
			events: []Event{
				{Type: EventInputValidated},
				{Type: EventWorkspaceDiscoveryStarted},
				{Type: EventWorkspaceDiscovered},
			},
		},
		{
			name: "report after cancellation",
			events: []Event{
				{Type: EventInputValidated},
				{Type: EventOperationCancelled},
				{Type: EventReportReady},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			validator := NewEventValidator()
			for index, event := range tc.events {
				err := validator.Emit(event)
				if index == len(tc.events)-1 {
					if err == nil {
						t.Fatalf("final event %q was accepted", event.Type)
					}
					return
				}
				if err != nil {
					t.Fatalf("event %q: %v", event.Type, err)
				}
			}
		})
	}
}

func TestEventValidatorRejectsEventsAfterTerminalState(t *testing.T) {
	validator := NewEventValidator()
	for _, event := range []Event{
		{Type: EventInputValidated},
		{Type: EventOperationCancelled},
	} {
		if err := validator.Emit(event); err != nil {
			t.Fatal(err)
		}
	}
	if err := validator.Emit(Event{Type: EventReportReady}); err == nil {
		t.Fatal("report-ready event accepted after cancellation")
	}
}
