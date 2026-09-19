package remediation

import (
	"context"
	"errors"
	"testing"
)

type offlinePlanningAdapter struct{}

func (offlinePlanningAdapter) Name() string { return "offline-fixture" }

func (offlinePlanningAdapter) Assess(context.Context, PlanningRequest) (AdapterAssessment, error) {
	return AdapterAssessment{State: AdapterSupported}, nil
}

func (offlinePlanningAdapter) Plan(context.Context, PlanningRequest) (PlanningEvidence, error) {
	return PlanningEvidence{
		State: AdapterSupported,
		Candidates: []Candidate{{
			ID: "example-1.2.0", Version: "1.2.0", State: CandidateViable,
			Evidence: CompatibilityEvidence{ConstraintSatisfied: true, LockfileResolution: true, PeerCompatible: true, RuntimeCompatible: true, EngineCompatible: true},
		}},
		AffectedFiles: []AffectedFile{{Path: "package.json", Kind: "manifest"}},
		Commands:      []Command{{Executable: "npm", Arguments: []string{"install", "example@1.2.0"}, WorkingDirectory: "services/api"}},
		Risks:         []Risk{}, Assumptions: []Assumption{}, Verification: []Verification{{
			ID: "test", Command: Command{Executable: "npm", Arguments: []string{"test"}, WorkingDirectory: "services/api"},
		}},
	}, nil
}

func validPlanningRequest() PlanningRequest {
	return PlanningRequest{
		Repository:      RepositoryIdentity{Root: "/repo", Repository: "example", Revision: "rev-1"},
		Workspace:       WorkspaceIdentity{ID: "api", Path: "services/api"},
		Finding:         FindingIdentity{StableKey: "finding-1", CurrentVersion: "1.0.0"},
		Component:       Component{PURL: "pkg:npm/example@1.0.0", Name: "example", Version: "1.0.0"},
		RepositoryState: "tree-1",
	}
}

func TestPlanningAdapterContractIsPackageManagerIndependent(t *testing.T) {
	var adapter PlanningAdapter = offlinePlanningAdapter{}
	request := validPlanningRequest()
	if err := ValidatePlanningRequest(request); err != nil {
		t.Fatal(err)
	}
	assessment, err := adapter.Assess(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	if err := ValidateAdapterAssessment(assessment); err != nil {
		t.Fatal(err)
	}
	evidence, err := adapter.Plan(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	if err := ValidatePlanningEvidence(evidence); err != nil {
		t.Fatal(err)
	}
	if evidence.Commands[0].Executable != "npm" || evidence.Commands[0].WorkingDirectory != request.Workspace.Path {
		t.Fatalf("commands = %#v", evidence.Commands)
	}
}

func TestPlanningAdapterStatesRequireExplicitReasons(t *testing.T) {
	for _, state := range []AdapterState{AdapterUnsupported, AdapterUnavailable, AdapterUnknown, AdapterRejected} {
		evidence := PlanningEvidence{State: state, Candidates: []Candidate{}, AffectedFiles: []AffectedFile{}, Commands: []Command{}, Risks: []Risk{}, Assumptions: []Assumption{}, Verification: []Verification{}}
		if err := ValidatePlanningEvidence(evidence); !errors.Is(err, ErrInvalidPlanningEvidence) {
			t.Fatalf("state %q error = %v", state, err)
		}
		evidence.Reason = "fixture reason"
		if err := ValidatePlanningEvidence(evidence); err != nil {
			t.Fatalf("state %q: %v", state, err)
		}
	}
}

func TestPlanningAdapterValidationRejectsUnsafeOrAmbiguousEvidence(t *testing.T) {
	request := validPlanningRequest()
	request.Workspace.Path = "../outside"
	if err := ValidatePlanningRequest(request); !errors.Is(err, ErrInvalidPlanningRequest) {
		t.Fatalf("unsafe request error = %v", err)
	}

	evidence := PlanningEvidence{State: AdapterSupported, Candidates: []Candidate{{ID: "same", Version: "1", State: CandidateUnknown}, {ID: "same", Version: "2", State: CandidateRejected}}, AffectedFiles: []AffectedFile{}, Commands: []Command{}, Risks: []Risk{}, Assumptions: []Assumption{}, Verification: []Verification{}}
	if err := ValidatePlanningEvidence(evidence); !errors.Is(err, ErrInvalidPlanningEvidence) {
		t.Fatalf("duplicate candidate error = %v", err)
	}
}
