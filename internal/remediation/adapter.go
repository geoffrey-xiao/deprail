package remediation

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

// AdapterState describes whether an ecosystem adapter can produce trustworthy evidence.
type AdapterState string

const (
	AdapterSupported   AdapterState = "supported"
	AdapterUnsupported AdapterState = "unsupported"
	AdapterUnavailable AdapterState = "unavailable"
	AdapterUnknown     AdapterState = "unknown"
	AdapterRejected    AdapterState = "rejected"
)

// AdapterAssessment is a read-only capability result. It contains no executable action.
type AdapterAssessment struct {
	State  AdapterState `json:"state"`
	Reason string       `json:"reason,omitempty"`
}

// PlanningRequest contains normalized repository and finding evidence.
type PlanningRequest struct {
	Repository      RepositoryIdentity `json:"repository"`
	Workspace       WorkspaceIdentity  `json:"workspace"`
	Finding         FindingIdentity    `json:"finding"`
	Component       Component          `json:"component"`
	CurrentState    CurrentState       `json:"current_state"`
	RepositoryState string             `json:"repository_state"`
}

// PlanningEvidence is the normalized, read-only output consumed by plan generation.
type PlanningEvidence struct {
	State         AdapterState   `json:"state"`
	Reason        string         `json:"reason,omitempty"`
	Candidates    []Candidate    `json:"candidates"`
	AffectedFiles []AffectedFile `json:"affected_files"`
	Commands      []Command      `json:"commands"`
	Risks         []Risk         `json:"risks"`
	Assumptions   []Assumption   `json:"assumptions"`
	Verification  []Verification `json:"verification"`
}

// PlanningAdapter is the package-manager-independent read-only planning port.
// It intentionally has no execute, install, update, or script method.
type PlanningAdapter interface {
	Name() string
	Assess(context.Context, PlanningRequest) (AdapterAssessment, error)
	Plan(context.Context, PlanningRequest) (PlanningEvidence, error)
}

var (
	ErrInvalidPlanningRequest  = errors.New("invalid planning request")
	ErrInvalidPlanningEvidence = errors.New("invalid planning evidence")
)

func ValidatePlanningRequest(request PlanningRequest) error {
	if request.Repository.Root == "" || request.Repository.Repository == "" || request.RepositoryState == "" {
		return fmt.Errorf("repository identity and state are required: %w", ErrInvalidPlanningRequest)
	}
	if request.Workspace.ID == "" || validateRelativePath(request.Workspace.Path) != nil {
		return fmt.Errorf("workspace identity is invalid: %w", ErrInvalidPlanningRequest)
	}
	if request.Finding.StableKey == "" || request.Finding.CurrentVersion == "" {
		return fmt.Errorf("finding identity is incomplete: %w", ErrInvalidPlanningRequest)
	}
	if request.Component.PURL == "" || request.Component.Name == "" || request.Component.Version == "" {
		return fmt.Errorf("component identity is incomplete: %w", ErrInvalidPlanningRequest)
	}
	return nil
}

func ValidatePlanningEvidence(evidence PlanningEvidence) error {
	if !validAdapterState(evidence.State) {
		return fmt.Errorf("unsupported adapter state %q: %w", evidence.State, ErrInvalidPlanningEvidence)
	}
	if evidence.State != AdapterSupported && strings.TrimSpace(evidence.Reason) == "" {
		return fmt.Errorf("non-supported adapter evidence requires a reason: %w", ErrInvalidPlanningEvidence)
	}
	if evidence.Candidates == nil || evidence.AffectedFiles == nil || evidence.Commands == nil || evidence.Risks == nil || evidence.Assumptions == nil || evidence.Verification == nil {
		return fmt.Errorf("planning evidence collections must be arrays: %w", ErrInvalidPlanningEvidence)
	}
	seen := make(map[string]struct{}, len(evidence.Candidates))
	for _, candidate := range evidence.Candidates {
		if candidate.ID == "" || candidate.Version == "" || !validCandidateState(candidate.State) {
			return fmt.Errorf("candidate identity or state is invalid: %w", ErrInvalidPlanningEvidence)
		}
		if candidate.State == CandidateRecommended && !candidate.Evidence.recommendable() {
			return fmt.Errorf("recommended candidate has unsafe compatibility evidence: %w", ErrInvalidPlanningEvidence)
		}
		if _, exists := seen[candidate.ID]; exists {
			return fmt.Errorf("candidate IDs must be unique: %w", ErrInvalidPlanningEvidence)
		}
		seen[candidate.ID] = struct{}{}
	}
	for _, file := range evidence.AffectedFiles {
		if file.Path == "" || validateRelativePath(file.Path) != nil || file.Kind == "" {
			return fmt.Errorf("affected file is invalid: %w", ErrInvalidPlanningEvidence)
		}
	}
	for _, command := range evidence.Commands {
		if err := validateCommand(command); err != nil {
			return fmt.Errorf("planning command is invalid: %w", ErrInvalidPlanningEvidence)
		}
	}
	seenVerification := make(map[string]struct{}, len(evidence.Verification))
	for _, verification := range evidence.Verification {
		if verification.ID == "" {
			return fmt.Errorf("verification identity is required: %w", ErrInvalidPlanningEvidence)
		}
		if _, exists := seenVerification[verification.ID]; exists {
			return fmt.Errorf("verification IDs must be unique: %w", ErrInvalidPlanningEvidence)
		}
		seenVerification[verification.ID] = struct{}{}
		if err := validateCommand(verification.Command); err != nil {
			return fmt.Errorf("verification command is invalid: %w", ErrInvalidPlanningEvidence)
		}
	}
	return nil
}
func ValidateAdapterAssessment(assessment AdapterAssessment) error {
	if !validAdapterState(assessment.State) {
		return fmt.Errorf("unsupported adapter state %q: %w", assessment.State, ErrInvalidPlanningEvidence)
	}
	if assessment.State != AdapterSupported && strings.TrimSpace(assessment.Reason) == "" {
		return fmt.Errorf("non-supported adapter assessment requires a reason: %w", ErrInvalidPlanningEvidence)
	}
	return nil
}

func validAdapterState(state AdapterState) bool {
	switch state {
	case AdapterSupported, AdapterUnsupported, AdapterUnavailable, AdapterUnknown, AdapterRejected:
		return true
	default:
		return false
	}
}
