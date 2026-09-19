package remediation

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"sort"
	"strings"
)

const SchemaVersion = "v0alpha1"

type CandidateState string

const (
	CandidateRecommended CandidateState = "recommended"
	CandidateViable      CandidateState = "viable"
	CandidateRejected    CandidateState = "rejected"
	CandidateUnavailable CandidateState = "unavailable"
	CandidateUnknown     CandidateState = "unknown"
)

type Plan struct {
	SchemaVersion      string             `json:"schema_version"`
	PlanID             string             `json:"plan_id"`
	CreatedFrom        CreatedFrom        `json:"created_from"`
	RepositoryIdentity RepositoryIdentity `json:"repository_identity"`
	WorkspaceIdentity  WorkspaceIdentity  `json:"workspace_identity"`
	FindingIdentity    FindingIdentity    `json:"finding_identity"`
	Component          Component          `json:"component"`
	CurrentState       CurrentState       `json:"current_state"`
	Candidates         []Candidate        `json:"candidates"`
	Recommendation     *Recommendation    `json:"recommendation"`
	AffectedFiles      []AffectedFile     `json:"affected_files"`
	Commands           []Command          `json:"commands"`
	Risks              []Risk             `json:"risks"`
	Assumptions        []Assumption       `json:"assumptions"`
	Verification       []Verification     `json:"verification"`
	Rollback           Rollback           `json:"rollback"`
	Provenance         Provenance         `json:"provenance"`
}

type CreatedFrom struct {
	ReportDigest    string `json:"report_digest"`
	SourceScanID    string `json:"source_scan_id"`
	Scanner         string `json:"scanner"`
	ScannerVersion  string `json:"scanner_version,omitempty"`
	RepositoryState string `json:"repository_state"`
}

type RepositoryIdentity struct {
	Root       string `json:"root"`
	Repository string `json:"repository"`
	Revision   string `json:"revision,omitempty"`
}

type WorkspaceIdentity struct {
	ID   string `json:"id"`
	Path string `json:"path"`
}

type FindingIdentity struct {
	StableKey      string   `json:"stable_key"`
	Aliases        []string `json:"aliases,omitempty"`
	DependencyPath []string `json:"dependency_path,omitempty"`
	CurrentVersion string   `json:"current_version"`
	FixedVersions  []string `json:"fixed_versions,omitempty"`
}

type Component struct {
	PURL    string `json:"purl"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

type CurrentState struct {
	Direct     bool   `json:"direct"`
	Constraint string `json:"constraint,omitempty"`
	Lockfile   string `json:"lockfile,omitempty"`
	Manifest   string `json:"manifest,omitempty"`
	Vulnerable bool   `json:"vulnerable"`
}

type Candidate struct {
	ID            string                `json:"id"`
	Version       string                `json:"version"`
	State         CandidateState        `json:"state"`
	Evidence      CompatibilityEvidence `json:"evidence"`
	Direct        bool                  `json:"direct"`
	MajorChange   bool                  `json:"major_change"`
	RejectionCode string                `json:"rejection_code,omitempty"`
	Reason        string                `json:"reason,omitempty"`
}

type CompatibilityEvidence struct {
	ConstraintSatisfied bool     `json:"constraint_satisfied"`
	LockfileResolution  bool     `json:"lockfile_resolution"`
	PeerCompatible      bool     `json:"peer_compatible"`
	RuntimeCompatible   bool     `json:"runtime_compatible"`
	EngineCompatible    bool     `json:"engine_compatible"`
	Unknown             []string `json:"unknown,omitempty"`
}

func (e CompatibilityEvidence) recommendable() bool {
	return len(e.Unknown) == 0 &&
		e.ConstraintSatisfied &&
		e.LockfileResolution &&
		e.PeerCompatible &&
		e.RuntimeCompatible &&
		e.EngineCompatible
}

type Recommendation struct {
	CandidateID string         `json:"candidate_id"`
	State       CandidateState `json:"state"`
	Reason      string         `json:"reason,omitempty"`
}

type AffectedFile struct {
	Path   string `json:"path"`
	Kind   string `json:"kind"`
	Effect string `json:"effect,omitempty"`
}

type Command struct {
	Executable       string   `json:"executable"`
	Arguments        []string `json:"arguments"`
	WorkingDirectory string   `json:"working_directory"`
}

type Risk struct {
	Code     string `json:"code"`
	Severity string `json:"severity"`
	Detected bool   `json:"detected"`
	Details  string `json:"details,omitempty"`
}

type Assumption struct {
	Code    string `json:"code"`
	Details string `json:"details,omitempty"`
}

type Verification struct {
	ID      string  `json:"id"`
	Command Command `json:"command"`
	Reason  string  `json:"reason,omitempty"`
}

type Rollback struct {
	Steps  []string `json:"steps,omitempty"`
	Limits string   `json:"limits,omitempty"`
}

type Provenance struct {
	ArtifactDigests []string `json:"artifact_digests"`
	Sources         []string `json:"sources"`
}

func (p *Plan) Canonicalize() {
	sort.Strings(p.FindingIdentity.Aliases)
	sort.Strings(p.FindingIdentity.FixedVersions)
	sort.Slice(p.Candidates, func(i, j int) bool { return p.Candidates[i].ID < p.Candidates[j].ID })
	for i := range p.Candidates {
		sort.Strings(p.Candidates[i].Evidence.Unknown)
	}
	sort.Slice(p.AffectedFiles, func(i, j int) bool { return p.AffectedFiles[i].Path < p.AffectedFiles[j].Path })
	sort.Slice(p.Commands, func(i, j int) bool { return commandKey(p.Commands[i]) < commandKey(p.Commands[j]) })
	sort.Slice(p.Risks, func(i, j int) bool { return p.Risks[i].Code < p.Risks[j].Code })
	sort.Slice(p.Assumptions, func(i, j int) bool { return p.Assumptions[i].Code < p.Assumptions[j].Code })
	sort.Slice(p.Verification, func(i, j int) bool { return p.Verification[i].ID < p.Verification[j].ID })
	sort.Strings(p.Provenance.ArtifactDigests)
	sort.Strings(p.Provenance.Sources)
}

func (p Plan) Validate() error {
	if p.SchemaVersion != SchemaVersion || p.PlanID == "" {
		return errors.New("plan schema version is unsupported or plan ID is missing")
	}
	if p.CreatedFrom.ReportDigest == "" || p.CreatedFrom.SourceScanID == "" || p.CreatedFrom.Scanner == "" || p.CreatedFrom.RepositoryState == "" {
		return errors.New("plan source report digest, scan ID, scanner, and repository state are required")
	}
	if p.RepositoryIdentity.Root == "" || p.WorkspaceIdentity.ID == "" || p.FindingIdentity.StableKey == "" || p.Component.PURL == "" {
		return errors.New("plan repository, workspace, finding, and component identity are required")
	}
	if p.PlanID != StablePlanID(p) {
		return errors.New("plan ID does not match canonical plan identity")
	}
	if p.WorkspaceIdentity.Path != "" {
		if err := validateRelativePath(p.WorkspaceIdentity.Path); err != nil {
			return fmt.Errorf("workspace path: %w", err)
		}
	}
	for _, file := range p.AffectedFiles {
		if file.Path == "" {
			return errors.New("affected file path is required")
		}
		if err := validateRelativePath(file.Path); err != nil {
			return fmt.Errorf("affected file %q: %w", file.Path, err)
		}
	}
	seen := make(map[string]struct{}, len(p.Candidates))
	for _, candidate := range p.Candidates {
		if candidate.ID == "" || candidate.Version == "" || !validCandidateState(candidate.State) {
			return fmt.Errorf("candidate identity or state is invalid: %q", candidate.ID)
		}
		if _, ok := seen[candidate.ID]; ok {
			return fmt.Errorf("candidate ID is duplicated: %q", candidate.ID)
		}
		seen[candidate.ID] = struct{}{}
		if candidate.State == CandidateRecommended && !candidate.Evidence.recommendable() {
			return fmt.Errorf("candidate %q cannot be recommended with incomplete or incompatible evidence", candidate.ID)
		}
	}
	if p.Recommendation != nil {
		candidate, ok := seenCandidate(p.Candidates, p.Recommendation.CandidateID)
		if !ok || candidate.State != CandidateRecommended || p.Recommendation.State != CandidateRecommended {
			return errors.New("recommendation must reference a recommended candidate")
		}
		if len(candidate.Evidence.Unknown) > 0 || !candidate.Evidence.recommendable() {
			return errors.New("recommendation cannot use incomplete or incompatible compatibility evidence")
		}
	}
	for _, command := range p.Commands {
		if err := validateCommand(command); err != nil {
			return err
		}
	}
	for _, verification := range p.Verification {
		if verification.ID == "" {
			return errors.New("verification ID is required")
		}
		if err := validateCommand(verification.Command); err != nil {
			return fmt.Errorf("verification %q: %w", verification.ID, err)
		}
	}
	return nil
}

func StablePlanID(plan Plan) string {
	plan.Canonicalize()
	identity := planIdentity{
		SchemaVersion:      plan.SchemaVersion,
		CreatedFrom:        plan.CreatedFrom,
		RepositoryIdentity: plan.RepositoryIdentity,
		WorkspaceIdentity:  plan.WorkspaceIdentity,
		FindingIdentity:    plan.FindingIdentity,
		Component:          plan.Component,
		CurrentState:       plan.CurrentState,
		Candidates:         identityCandidates(plan.Candidates),
		Recommendation:     identityRecommendation(plan.Recommendation),
		AffectedFiles:      identityAffectedFiles(plan.AffectedFiles),
		Commands:           plan.Commands,
		Risks:              identityRisks(plan.Risks),
		Assumptions:        identityAssumptions(plan.Assumptions),
		Verification:       identityVerifications(plan.Verification),
		Provenance:         plan.Provenance,
	}
	encoded, _ := json.Marshal(identity)
	digest := sha256.Sum256(encoded)
	return hex.EncodeToString(digest[:])
}

type planIdentity struct {
	SchemaVersion      string                  `json:"schema_version"`
	CreatedFrom        CreatedFrom             `json:"created_from"`
	RepositoryIdentity RepositoryIdentity      `json:"repository_identity"`
	WorkspaceIdentity  WorkspaceIdentity       `json:"workspace_identity"`
	FindingIdentity    FindingIdentity         `json:"finding_identity"`
	Component          Component               `json:"component"`
	CurrentState       CurrentState            `json:"current_state"`
	Candidates         []candidateIdentity     `json:"candidates"`
	Recommendation     *recommendationIdentity `json:"recommendation"`
	AffectedFiles      []affectedFileIdentity  `json:"affected_files"`
	Commands           []Command               `json:"commands"`
	Risks              []riskIdentity          `json:"risks"`
	Assumptions        []assumptionIdentity    `json:"assumptions"`
	Verification       []verificationIdentity  `json:"verification"`
	Provenance         Provenance              `json:"provenance"`
}

type recommendationIdentity struct {
	CandidateID string         `json:"candidate_id"`
	State       CandidateState `json:"state"`
}

type candidateIdentity struct {
	ID            string                `json:"id"`
	Version       string                `json:"version"`
	State         CandidateState        `json:"state"`
	Evidence      CompatibilityEvidence `json:"evidence"`
	Direct        bool                  `json:"direct"`
	MajorChange   bool                  `json:"major_change"`
	RejectionCode string                `json:"rejection_code,omitempty"`
}

type affectedFileIdentity struct {
	Path string `json:"path"`
	Kind string `json:"kind"`
}

type riskIdentity struct {
	Code     string `json:"code"`
	Severity string `json:"severity"`
	Detected bool   `json:"detected"`
}

func identityRecommendation(value *Recommendation) *recommendationIdentity {
	if value == nil {
		return nil
	}
	return &recommendationIdentity{CandidateID: value.CandidateID, State: value.State}
}

type assumptionIdentity struct {
	Code string `json:"code"`
}

type verificationIdentity struct {
	ID      string  `json:"id"`
	Command Command `json:"command"`
}

func identityCandidates(values []Candidate) []candidateIdentity {
	result := make([]candidateIdentity, 0, len(values))
	for _, value := range values {
		result = append(result, candidateIdentity{ID: value.ID, Version: value.Version, State: value.State, Evidence: value.Evidence, Direct: value.Direct, MajorChange: value.MajorChange, RejectionCode: value.RejectionCode})
	}
	return result
}

func identityAffectedFiles(values []AffectedFile) []affectedFileIdentity {
	result := make([]affectedFileIdentity, 0, len(values))
	for _, value := range values {
		result = append(result, affectedFileIdentity{Path: value.Path, Kind: value.Kind})
	}
	return result
}

func identityRisks(values []Risk) []riskIdentity {
	result := make([]riskIdentity, 0, len(values))
	for _, value := range values {
		result = append(result, riskIdentity{Code: value.Code, Severity: value.Severity, Detected: value.Detected})
	}
	return result
}

func identityAssumptions(values []Assumption) []assumptionIdentity {
	result := make([]assumptionIdentity, 0, len(values))
	for _, value := range values {
		result = append(result, assumptionIdentity{Code: value.Code})
	}
	return result
}

func identityVerifications(values []Verification) []verificationIdentity {
	result := make([]verificationIdentity, 0, len(values))
	for _, value := range values {
		result = append(result, verificationIdentity{ID: value.ID, Command: value.Command})
	}
	return result
}

func validCandidateState(state CandidateState) bool {
	switch state {
	case CandidateRecommended, CandidateViable, CandidateRejected, CandidateUnavailable, CandidateUnknown:
		return true
	default:
		return false
	}
}

func seenCandidate(candidates []Candidate, id string) (Candidate, bool) {
	for _, candidate := range candidates {
		if candidate.ID == id {
			return candidate, true
		}
	}
	return Candidate{}, false
}

func commandKey(command Command) string {
	encoded, _ := json.Marshal(command)
	return string(encoded)
}

func validateCommand(command Command) error {
	if strings.TrimSpace(command.Executable) == "" {
		return errors.New("command executable is required")
	}
	if command.Arguments == nil {
		return errors.New("command arguments must be an array")
	}
	if err := validateRelativePath(command.WorkingDirectory); err != nil {
		return fmt.Errorf("command working directory: %w", err)
	}
	return nil
}

func validateRelativePath(value string) error {
	if value == "" || strings.HasPrefix(value, "/") || strings.HasPrefix(value, "\\") || strings.Contains(value, ":") {
		return errors.New("must be a canonical repository-relative path")
	}
	if strings.Contains(value, "\\") {
		return errors.New("must use slash separators")
	}
	if value != "." && path.Clean(value) != value {
		return errors.New("must use canonical path spelling")
	}
	for _, part := range strings.Split(value, "/") {
		if part == ".." {
			return errors.New("must not traverse outside the repository")
		}
	}
	return nil
}
