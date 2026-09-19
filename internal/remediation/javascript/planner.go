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
	root, err := safeRepositoryPath(request.Repository.Root, filepath.ToSlash(filepath.Join(request.Workspace.Path, "package.json")))
	if err != nil {
		return remediation.AdapterAssessment{State: remediation.AdapterRejected, Reason: "workspace manifest escapes repository"}, nil
	}
	if _, err := os.Stat(root); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return remediation.AdapterAssessment{State: remediation.AdapterUnsupported, Reason: "package.json is unavailable"}, nil
		}
		return remediation.AdapterAssessment{State: remediation.AdapterRejected, Reason: "package.json is unreadable"}, nil
	}
	manager, _, err := locateLockfile(request)
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
		return emptyEvidence(remediation.AdapterRejected, err.Error()), nil
	}
	manifestPath, err := safeRepositoryPath(request.Repository.Root, filepath.ToSlash(filepath.Join(request.Workspace.Path, "package.json")))
	if err != nil {
		return emptyEvidence(remediation.AdapterRejected, "workspace manifest escapes repository"), nil
	}
	manifest, err := readManifest(manifestPath)
	if err != nil {
		return emptyEvidence(remediation.AdapterRejected, "package.json is malformed or unreadable"), nil
	}
	manager, lockfile, err := locateLockfile(request)
	if err != nil {
		return emptyEvidence(remediation.AdapterRejected, err.Error()), nil
	}
	if manager == "" {
		return emptyEvidence(remediation.AdapterUnavailable, "no supported JavaScript lockfile is available"), nil
	}
	lockPath, err := safeRepositoryPath(request.Repository.Root, lockfile)
	if err != nil {
		return emptyEvidence(remediation.AdapterRejected, "lockfile escapes repository"), nil
	}
	if err := validateLockfile(lockPath, manager); err != nil {
		return emptyEvidence(remediation.AdapterRejected, err.Error()), nil
	}

	constraint, direct := dependencyConstraint(manifest, request.Component.Name)
	owner := request.Component.Name
	if !direct {
		owner, direct = directOwner(manifest, request.Finding.DependencyPath)
		if !direct {
			return unknownEvidence("direct dependency owner is not identifiable", lockfile, request), nil
		}
		constraint, _ = dependencyConstraint(manifest, owner)
	}
	candidates := buildCandidates(request, constraint, direct)
	files := []remediation.AffectedFile{{Path: filepath.ToSlash(filepath.Join(request.Workspace.Path, "package.json")), Kind: "manifest", Effect: "dependency constraint"}, {Path: lockfile, Kind: "lockfile", Effect: "resolution"}}
	risks := candidateRisks(candidates)
	if len(candidates) == 0 {
		return remediation.PlanningEvidence{State: remediation.AdapterUnknown, Reason: "no non-vulnerable candidate version is available in retained evidence", Candidates: candidates, AffectedFiles: files, Commands: []remediation.Command{}, Risks: risks, Assumptions: []remediation.Assumption{}, Verification: []remediation.Verification{}}, nil
	}
	command := remediation.Command{Executable: manager, Arguments: commandArguments(manager, owner, candidates), WorkingDirectory: request.Workspace.Path}
	return remediation.PlanningEvidence{State: remediation.AdapterSupported, Candidates: candidates, AffectedFiles: files, Commands: []remediation.Command{command}, Risks: risks, Assumptions: []remediation.Assumption{}, Verification: []remediation.Verification{}}, nil
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

func directOwner(manifest packageManifest, dependencyPath []string) (string, bool) {
	for _, name := range dependencyPath {
		if _, ok := dependencyConstraint(manifest, name); ok {
			return name, true
		}
	}
	return "", false
}

func locateLockfile(request remediation.PlanningRequest) (string, string, error) {
	if request.CurrentState.Lockfile != "" {
		lockfile := filepath.ToSlash(filepath.Clean(request.CurrentState.Lockfile))
		if filepath.IsAbs(filepath.FromSlash(lockfile)) || strings.Contains(lockfile, "..") {
			return "", "", errors.New("retained lockfile path is unsafe")
		}
		manager := lockfileManager(lockfile)
		if manager == "" {
			return "", "", errors.New("retained lockfile type is unsupported")
		}
		if _, err := safeRepositoryPath(request.Repository.Root, lockfile); err != nil {
			return "", "", errors.New("retained lockfile escapes repository")
		}
		return manager, lockfile, nil
	}
	dir := filepath.Join(request.Repository.Root, filepath.FromSlash(request.Workspace.Path))
	candidates := []struct{ name, manager string }{{"package-lock.json", "npm"}, {"npm-shrinkwrap.json", "npm"}, {"pnpm-lock.yaml", "pnpm"}, {"yarn.lock", "yarn"}}
	for _, candidate := range candidates {
		path := filepath.Join(dir, candidate.name)
		if _, err := os.Stat(path); err == nil {
			return candidate.manager, filepath.ToSlash(filepath.Join(request.Workspace.Path, candidate.name)), nil
		}
	}
	return "", "", nil
}

func lockfileManager(path string) string {
	switch filepath.Base(filepath.FromSlash(path)) {
	case "package-lock.json", "npm-shrinkwrap.json":
		return "npm"
	case "pnpm-lock.yaml":
		return "pnpm"
	case "yarn.lock":
		return "yarn"
	default:
		return ""
	}
}

func safeRepositoryPath(root, relative string) (string, error) {
	canonicalRoot, err := remediation.CanonicalRepositoryRoot(root)
	if err != nil {
		return "", err
	}
	clean := filepath.Clean(filepath.Join(canonicalRoot, filepath.FromSlash(relative)))
	rel, err := filepath.Rel(canonicalRoot, clean)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", errors.New("path escapes repository")
	}
	resolved, err := filepath.EvalSymlinks(clean)
	if err != nil {
		return "", err
	}
	resolvedRel, err := filepath.Rel(canonicalRoot, resolved)
	if err != nil || resolvedRel == ".." || strings.HasPrefix(resolvedRel, ".."+string(filepath.Separator)) {
		return "", errors.New("resolved path escapes repository")
	}
	return resolved, nil
}

func validateLockfile(path, manager string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return errors.New("JavaScript lockfile is unreadable")
	}
	if len(data) == 0 {
		return errors.New("JavaScript lockfile is empty")
	}
	switch manager {
	case "npm":
		var document map[string]any
		if err := json.Unmarshal(data, &document); err != nil {
			return errors.New("npm lockfile is malformed")
		}
	case "pnpm":
		if !validPnpmLockfile(string(data)) {
			return errors.New("pnpm lockfile is malformed")
		}
	case "yarn":
		if !validYarnLockfile(string(data)) {
			return errors.New("Yarn lockfile is malformed")
		}
	}
	return nil
}

func validPnpmLockfile(data string) bool {
	for _, line := range strings.Split(data, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "lockfileVersion:") {
			return strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "lockfileVersion:")) != ""
		}
	}
	return false
}

func validYarnLockfile(data string) bool {
	if !strings.HasPrefix(strings.TrimSpace(data), "# yarn lockfile v") {
		return false
	}
	return strings.Contains(data, "\n  version ") || strings.Contains(data, "\nversion ")
}

func buildCandidates(request remediation.PlanningRequest, constraint string, direct bool) []remediation.Candidate {
	candidates := make([]remediation.Candidate, 0, len(request.Finding.FixedVersions))
	for _, raw := range request.Finding.FixedVersions {
		parsed, ok := parseVersion(raw)
		if !ok || !isNewer(parsed, request.Component.Version) {
			continue
		}
		candidate := remediation.Candidate{ID: raw, Version: raw, Direct: direct, Evidence: remediation.CompatibilityEvidence{Unknown: []string{"candidate_metadata_unavailable"}}}
		if constraint != "" && !satisfies(constraint, parsed) {
			candidate.State = remediation.CandidateRejected
			candidate.RejectionCode = "CONSTRAINT_CHANGE_REQUIRED"
			candidate.Reason = "candidate requires changing the manifest constraint"
			candidate.MajorChange = parseMajor(parsed) != parseMajorVersion(request.Component.Version)
		} else {
			candidate.State = remediation.CandidateViable
			candidate.Evidence.Unknown = append(candidate.Evidence.Unknown, "peer_compatibility_unknown", "runtime_compatibility_unknown", "engine_compatibility_unknown")
			candidate.MajorChange = parseMajor(parsed) != parseMajorVersion(request.Component.Version)
		}
		candidates = append(candidates, candidate)
	}
	sort.Slice(candidates, func(i, j int) bool { return compareVersions(candidates[i].Version, candidates[j].Version) < 0 })
	return candidates
}

func candidateRisks(candidates []remediation.Candidate) []remediation.Risk {
	risks := make([]remediation.Risk, 0)
	for _, candidate := range candidates {
		if candidate.MajorChange {
			risks = append(risks, remediation.Risk{Code: "MAJOR_VERSION_CHANGE", Severity: "high", Detected: true, Details: "candidate crosses a major version boundary"})
		}
		if candidate.State == remediation.CandidateRejected && candidate.RejectionCode == "CONSTRAINT_CHANGE_REQUIRED" {
			risks = append(risks, remediation.Risk{Code: "CONSTRAINT_CHANGE_REQUIRED", Severity: "high", Detected: true, Details: "candidate is outside the current manifest constraint"})
		}
	}
	return risks
}

func unknownEvidence(reason, lockfile string, request remediation.PlanningRequest) remediation.PlanningEvidence {
	return remediation.PlanningEvidence{State: remediation.AdapterUnknown, Reason: reason, Candidates: []remediation.Candidate{}, AffectedFiles: []remediation.AffectedFile{{Path: filepath.ToSlash(filepath.Join(request.Workspace.Path, "package.json")), Kind: "manifest"}, {Path: lockfile, Kind: "lockfile"}}, Commands: []remediation.Command{}, Risks: []remediation.Risk{}, Assumptions: []remediation.Assumption{}, Verification: []remediation.Verification{}}
}

func commandArguments(manager, name string, candidates []remediation.Candidate) []string {
	version := candidates[0].Version
	spec := name + "@" + version
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
