package javascript

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/geoffrey-xiao/deprail/internal/remediation"
)

type Adapter struct{}

func (Adapter) Name() string { return "javascript" }

func (Adapter) Assess(_ context.Context, request remediation.PlanningRequest) (remediation.AdapterAssessment, error) {
	if err := remediation.ValidatePlanningRequest(request); err != nil {
		return remediation.AdapterAssessment{State: remediation.AdapterRejected, Reason: err.Error()}, nil
	}
	manifest := filepath.Join(request.Repository.Root, filepath.FromSlash(request.Workspace.Path), "package.json")
	if _, err := os.Stat(manifest); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return remediation.AdapterAssessment{State: remediation.AdapterUnsupported, Reason: "package.json is unavailable"}, nil
		}
		return remediation.AdapterAssessment{State: remediation.AdapterRejected, Reason: "package.json is unreadable"}, nil
	}
	manager, _, err := locateLockfile(request.Repository.Root, request.Workspace.Path)
	if err != nil {
		return remediation.AdapterAssessment{State: remediation.AdapterRejected, Reason: err.Error()}, nil
	}
	if manager == "" {
		return remediation.AdapterAssessment{State: remediation.AdapterUnavailable, Reason: "no supported JavaScript lockfile is available"}, nil
	}
	return remediation.AdapterAssessment{State: remediation.AdapterSupported}, nil
}

func (Adapter) Plan(ctx context.Context, request remediation.PlanningRequest) (remediation.PlanningEvidence, error) {
	if err := ctx.Err(); err != nil {
		return remediation.PlanningEvidence{}, err
	}
	if err := remediation.ValidatePlanningRequest(request); err != nil {
		return remediation.PlanningEvidence{State: remediation.AdapterRejected, Reason: err.Error(), Candidates: []remediation.Candidate{}, AffectedFiles: []remediation.AffectedFile{}, Commands: []remediation.Command{}, Risks: []remediation.Risk{}, Assumptions: []remediation.Assumption{}, Verification: []remediation.Verification{}}, nil
	}
	workspaceRoot := filepath.Join(request.Repository.Root, filepath.FromSlash(request.Workspace.Path))
	manifestPath := filepath.Join(workspaceRoot, "package.json")
	manifest, err := readManifest(manifestPath)
	if err != nil {
		return emptyEvidence(remediation.AdapterRejected, "package.json is malformed or unreadable"), nil
	}
	manager, lockfile, err := locateLockfile(request.Repository.Root, request.Workspace.Path)
	if err != nil {
		return emptyEvidence(remediation.AdapterRejected, err.Error()), nil
	}
	if manager == "" {
		return emptyEvidence(remediation.AdapterUnavailable, "no supported JavaScript lockfile is available"), nil
	}
	if err := validateLockfile(filepath.Join(request.Repository.Root, filepath.FromSlash(lockfile)), manager); err != nil {
		return emptyEvidence(remediation.AdapterRejected, err.Error()), nil
	}
	constraint, direct := dependencyConstraint(manifest, request.Component.Name)
	if constraint == "" {
		direct = false
		constraint = request.CurrentState.Constraint
	}
	candidates := make([]remediation.Candidate, 0, len(request.Finding.FixedVersions))
	for _, version := range request.Finding.FixedVersions {
		parsed, ok := parseVersion(version)
		if !ok || (constraint != "" && !satisfies(constraint, parsed)) || !isNewer(parsed, request.Component.Version) {
			continue
		}
		candidates = append(candidates, remediation.Candidate{
			ID: version, Version: version, State: remediation.CandidateViable, Direct: direct,
			MajorChange: parseMajor(parsed) != parseMajorVersion(request.Component.Version),
			Evidence: remediation.CompatibilityEvidence{
				ConstraintSatisfied: true, LockfileResolution: true, PeerCompatible: true,
				RuntimeCompatible: true, EngineCompatible: true,
			},
		})
	}
	sort.Slice(candidates, func(i, j int) bool { return compareVersions(candidates[i].Version, candidates[j].Version) < 0 })
	files := []remediation.AffectedFile{{Path: filepath.ToSlash(filepath.Join(request.Workspace.Path, "package.json")), Kind: "manifest", Effect: "dependency constraint"}, {Path: lockfile, Kind: "lockfile", Effect: "resolution"}}
	commands := make([]remediation.Command, 0, 1)
	commands = append(commands, remediation.Command{Executable: manager, Arguments: commandArguments(manager, request.Component.Name, candidates), WorkingDirectory: request.Workspace.Path})
	risks := make([]remediation.Risk, 0)
	for _, candidate := range candidates {
		if candidate.MajorChange {
			risks = append(risks, remediation.Risk{Code: "MAJOR_VERSION_CHANGE", Severity: "high", Detected: true, Details: "candidate crosses a major version boundary"})
			break
		}
	}
	if len(candidates) == 0 {
		return remediation.PlanningEvidence{State: remediation.AdapterUnknown, Reason: "no non-vulnerable candidate version is available in retained evidence", Candidates: candidates, AffectedFiles: files, Commands: commands, Risks: risks, Assumptions: []remediation.Assumption{}, Verification: []remediation.Verification{}}, nil
	}
	return remediation.PlanningEvidence{State: remediation.AdapterSupported, Candidates: candidates, AffectedFiles: files, Commands: commands, Risks: risks, Assumptions: []remediation.Assumption{}, Verification: []remediation.Verification{}}, nil
}

type packageManifest struct {
	Dependencies         map[string]string `json:"dependencies"`
	DevDependencies      map[string]string `json:"devDependencies"`
	PeerDependencies     map[string]string `json:"peerDependencies"`
	OptionalDependencies map[string]string `json:"optionalDependencies"`
}

func readManifest(path string) (packageManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return packageManifest{}, err
	}
	var manifest packageManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return packageManifest{}, err
	}
	return manifest, nil
}

func dependencyConstraint(manifest packageManifest, name string) (string, bool) {
	for _, deps := range []map[string]string{manifest.Dependencies, manifest.DevDependencies, manifest.PeerDependencies, manifest.OptionalDependencies} {
		if value, ok := deps[name]; ok {
			return value, true
		}
	}
	return "", false
}

func locateLockfile(root, workspace string) (string, string, error) {
	dir := filepath.Join(root, filepath.FromSlash(workspace))
	candidates := []struct{ name, manager string }{{"package-lock.json", "npm"}, {"npm-shrinkwrap.json", "npm"}, {"pnpm-lock.yaml", "pnpm"}, {"yarn.lock", "yarn"}}
	found := make([]struct{ name, manager string }, 0, 1)
	for _, candidate := range candidates {
		if _, err := os.Stat(filepath.Join(dir, candidate.name)); err == nil {
			found = append(found, candidate)
		}
	}
	if len(found) > 1 {
		return "", "", errors.New("conflicting JavaScript lockfiles are present")
	}
	if len(found) == 0 {
		return "", "", nil
	}
	return found[0].manager, filepath.ToSlash(filepath.Join(workspace, found[0].name)), nil
}

func validateLockfile(path, manager string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return errors.New("JavaScript lockfile is unreadable")
	}
	if len(data) == 0 {
		return errors.New("JavaScript lockfile is empty")
	}
	if manager == "npm" {
		var document map[string]any
		if err := json.Unmarshal(data, &document); err != nil {
			return errors.New("npm lockfile is malformed")
		}
	}
	return nil
}

func commandArguments(manager, name string, candidates []remediation.Candidate) []string {
	version := ""
	if len(candidates) > 0 {
		version = candidates[0].Version
	}
	spec := name
	if version != "" {
		spec += "@" + version
	}
	switch manager {
	case "pnpm":
		return []string{"update", spec}
	case "yarn":
		return []string{"upgrade", spec}
	default:
		return []string{"install", spec}
	}
}

func emptyEvidence(state remediation.AdapterState, reason string) remediation.PlanningEvidence {
	return remediation.PlanningEvidence{State: state, Reason: reason, Candidates: []remediation.Candidate{}, AffectedFiles: []remediation.AffectedFile{}, Commands: []remediation.Command{}, Risks: []remediation.Risk{}, Assumptions: []remediation.Assumption{}, Verification: []remediation.Verification{}}
}

type version struct{ major, minor, patch int }

func parseVersion(value string) (version, bool) {
	parts := strings.Split(strings.TrimPrefix(strings.TrimSpace(value), "v"), ".")
	if len(parts) < 2 || len(parts) > 3 {
		return version{}, false
	}
	nums := [3]int{}
	for i := range parts {
		n, err := strconv.Atoi(strings.TrimLeft(parts[i], "0"))
		if err != nil {
			if parts[i] == "0" {
				n = 0
			} else {
				return version{}, false
			}
		}
		nums[i] = n
	}
	return version{nums[0], nums[1], nums[2]}, true
}
func parseMajorVersion(value string) int { v, _ := parseVersion(value); return v.major }
func parseMajor(v version) int           { return v.major }
func compareVersions(a, b string) int {
	av, aok := parseVersion(a)
	bv, bok := parseVersion(b)
	if !aok || !bok {
		return strings.Compare(a, b)
	}
	if av.major != bv.major {
		if av.major < bv.major {
			return -1
		}
		return 1
	}
	if av.minor != bv.minor {
		if av.minor < bv.minor {
			return -1
		}
		return 1
	}
	if av.patch < bv.patch {
		return -1
	}
	if av.patch > bv.patch {
		return 1
	}
	return 0
}
func isNewer(a version, current string) bool {
	cv, ok := parseVersion(current)
	return ok && (a.major > cv.major || a.major == cv.major && (a.minor > cv.minor || a.minor == cv.minor && a.patch > cv.patch))
}
func satisfies(constraint string, v version) bool {
	constraint = strings.TrimSpace(constraint)
	if constraint == "" || constraint == "*" {
		return true
	}
	if strings.HasPrefix(constraint, "^") {
		base, ok := parseVersion(strings.TrimPrefix(constraint, "^"))
		return ok && v.major == base.major && compareVersions(fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch), fmt.Sprintf("%d.%d.%d", base.major, base.minor, base.patch)) >= 0
	}
	if strings.HasPrefix(constraint, "~") {
		base, ok := parseVersion(strings.TrimPrefix(constraint, "~"))
		return ok && v.major == base.major && v.minor == base.minor && compareVersions(fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch), fmt.Sprintf("%d.%d.%d", base.major, base.minor, base.patch)) >= 0
	}
	base, ok := parseVersion(constraint)
	return ok && v == base
}
