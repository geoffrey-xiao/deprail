package remediation

import "testing"

func testPlan() Plan {
	plan := Plan{
		SchemaVersion:      SchemaVersion,
		CreatedFrom:        CreatedFrom{ReportDigest: "report-digest", SourceScanID: "scan-1", Scanner: "osv", RepositoryState: "tree-digest"},
		RepositoryIdentity: RepositoryIdentity{Root: ".", Repository: "repo", Revision: "rev-1"},
		WorkspaceIdentity:  WorkspaceIdentity{ID: "service", Path: "services/api"},
		FindingIdentity:    FindingIdentity{StableKey: "finding-1", Aliases: []string{"OSV-2", "CVE-1"}, DependencyPath: []string{"app", "library", "vulnerable"}, CurrentVersion: "1.0.0"},
		Component:          Component{PURL: "pkg:npm/example@1.0.0", Name: "example", Version: "1.0.0"},
		CurrentState:       CurrentState{Direct: true, Constraint: "^1.0.0", Manifest: "package.json", Vulnerable: true},
		Candidates:         []Candidate{{ID: "candidate-2", Version: "1.2.0", State: CandidateViable, Evidence: CompatibilityEvidence{ConstraintSatisfied: true, LockfileResolution: true, PeerCompatible: true, RuntimeCompatible: true, EngineCompatible: true}}, {ID: "candidate-1", Version: "1.1.0", State: CandidateRecommended, Evidence: CompatibilityEvidence{ConstraintSatisfied: true, LockfileResolution: true, PeerCompatible: true, RuntimeCompatible: true, EngineCompatible: true}}},
		Recommendation:     &Recommendation{CandidateID: "candidate-1", State: CandidateRecommended},
		AffectedFiles:      []AffectedFile{{Path: "package-lock.json", Kind: "lockfile"}, {Path: "package.json", Kind: "manifest"}},
		Commands:           []Command{{Executable: "npm", Arguments: []string{"install", "example@1.1.0"}, WorkingDirectory: "."}},
		Risks:              []Risk{{Code: "peer_dependency", Severity: "medium", Detected: false, Details: "presentation text"}},
		Assumptions:        []Assumption{{Code: "metadata_local", Details: "presentation text"}},
		Verification:       []Verification{{ID: "tests", Command: Command{Executable: "npm", Arguments: []string{"test"}, WorkingDirectory: "services/api"}, Reason: "presentation text"}},
		Rollback:           Rollback{Steps: []string{"restore package files"}},
		Provenance:         Provenance{ArtifactDigests: []string{"b", "a"}, Sources: []string{"scanner", "manifest"}},
	}
	plan.PlanID = StablePlanID(plan)
	return plan
}

func TestStablePlanIDIgnoresOrderingAndProse(t *testing.T) {
	left := testPlan()
	right := testPlan()
	right.Candidates[0], right.Candidates[1] = right.Candidates[1], right.Candidates[0]
	right.AffectedFiles[0], right.AffectedFiles[1] = right.AffectedFiles[1], right.AffectedFiles[0]
	right.Provenance.ArtifactDigests[0], right.Provenance.ArtifactDigests[1] = right.Provenance.ArtifactDigests[1], right.Provenance.ArtifactDigests[0]
	right.Risks[0].Details = "different prose"
	right.Assumptions[0].Details = "different prose"
	right.Verification[0].Reason = "different prose"
	right.Rollback.Steps[0] = "different prose"
	right.Recommendation.Reason = "different prose"
	if StablePlanID(left) != StablePlanID(right) {
		t.Fatal("equivalent plans produced different stable IDs")
	}
}

func TestStablePlanIDPreservesDependencyPathOrder(t *testing.T) {
	plan := testPlan()
	reversed := testPlan()
	reversed.FindingIdentity.DependencyPath[0], reversed.FindingIdentity.DependencyPath[2] = reversed.FindingIdentity.DependencyPath[2], reversed.FindingIdentity.DependencyPath[0]
	if StablePlanID(plan) == StablePlanID(reversed) {
		t.Fatal("reversed dependency paths produced the same stable ID")
	}
}
func TestValidateRejectsFailedCompatibilityRecommendation(t *testing.T) {
	plan := testPlan()
	plan.Candidates[0].Evidence.PeerCompatible = false
	plan.PlanID = StablePlanID(plan)
	if err := plan.Validate(); err == nil {
		t.Fatal("expected failed compatibility evidence to reject recommendation")
	}
}

func TestValidateRejectsUnsupportedSchemaAndMissingScanner(t *testing.T) {
	plan := testPlan()
	plan.SchemaVersion = "invalid"
	plan.PlanID = StablePlanID(plan)
	if err := plan.Validate(); err == nil {
		t.Fatal("expected unsupported schema to be rejected")
	}

	plan = testPlan()
	plan.CreatedFrom.Scanner = ""
	plan.PlanID = StablePlanID(plan)
	if err := plan.Validate(); err == nil {
		t.Fatal("expected missing scanner provenance to be rejected")
	}
}

func TestValidateRejectsNoncanonicalPaths(t *testing.T) {
	for _, path := range []string{"services/./api", "services//api", "services/api/"} {
		plan := testPlan()
		plan.WorkspaceIdentity.Path = path
		plan.PlanID = StablePlanID(plan)
		if err := plan.Validate(); err == nil {
			t.Fatalf("path %q was accepted", path)
		}
	}
}

func TestValidateRejectsUnknownRecommendation(t *testing.T) {
	plan := testPlan()
	for i := range plan.Candidates {
		if plan.Candidates[i].ID == "candidate-1" {
			plan.Candidates[i].Evidence.Unknown = []string{"peer metadata"}
		}
	}
	if err := plan.Validate(); err == nil {
		t.Fatal("expected unknown evidence to reject recommendation")
	}
}

func TestValidateRequiresContainedRelativeCommandDirectory(t *testing.T) {
	for _, directory := range []string{"/tmp", "../outside", `..\\outside`, "C:/repo"} {
		plan := testPlan()
		plan.Commands[0].WorkingDirectory = directory
		if err := plan.Validate(); err == nil {
			t.Fatalf("working directory %q was accepted", directory)
		}
	}
}

func TestValidateCommandWorkspaceCoversAuthorizedPaths(t *testing.T) {
	plan := testPlan()
	plan.AffectedFiles = []AffectedFile{
		{Path: "services/api/package.json", Kind: "manifest"},
		{Path: "services/api/package-lock.json", Kind: "lockfile"},
	}
	for _, directory := range []string{".", "services/api"} {
		if err := ValidateCommandWorkspace(plan, directory); err != nil {
			t.Fatalf("directory %q rejected: %v", directory, err)
		}
	}
	if err := ValidateCommandWorkspace(plan, "services/admin"); err == nil {
		t.Fatal("expected workspace outside authorized paths to be rejected")
	}
	plan.AffectedFiles = nil
	if err := ValidateCommandWorkspace(plan, "."); err == nil {
		t.Fatal("expected command without authorized paths to be rejected")
	}
}

func TestValidateAcceptsExplicitNoRecommendation(t *testing.T) {
	plan := testPlan()
	plan.Recommendation = nil
	plan.Candidates[1].State = CandidateUnknown
	plan.Candidates[1].Evidence.Unknown = []string{"compatibility metadata"}
	plan.PlanID = StablePlanID(plan)
	if err := plan.Validate(); err != nil {
		t.Fatalf("no-recommendation plan rejected: %v", err)
	}
}
