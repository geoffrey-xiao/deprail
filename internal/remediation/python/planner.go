package python

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/geoffrey-xiao/deprail/internal/remediation"
)

type Adapter struct{}

func (Adapter) Name() string { return "python" }
func (Adapter) Assess(_ context.Context, request remediation.PlanningRequest) (remediation.AdapterAssessment, error) {
	if err := remediation.ValidatePlanningRequest(request); err != nil {
		return remediation.AdapterAssessment{State: remediation.AdapterRejected, Reason: err.Error()}, nil
	}
	if _, err := manifestPath(request); err != nil {
		return remediation.AdapterAssessment{State: remediation.AdapterRejected, Reason: err.Error()}, nil
	}
	if _, _, err := locateManager(request); err != nil {
		return remediation.AdapterAssessment{State: remediation.AdapterRejected, Reason: err.Error()}, nil
	}
	manager, _, _ := locateManager(request)
	if manager == "" {
		return remediation.AdapterAssessment{State: remediation.AdapterUnavailable, Reason: "no supported Python manifest or lockfile is available"}, nil
	}
	return remediation.AdapterAssessment{State: remediation.AdapterSupported}, nil
}
func (Adapter) Plan(ctx context.Context, request remediation.PlanningRequest) (remediation.PlanningEvidence, error) {
	if err := ctx.Err(); err != nil {
		return remediation.PlanningEvidence{}, err
	}
	if err := remediation.ValidatePlanningRequest(request); err != nil {
		return empty(remediation.AdapterRejected, err.Error()), nil
	}
	manifest, err := manifestPath(request)
	if err != nil {
		return empty(remediation.AdapterRejected, err.Error()), nil
	}
	manager, lockfile, err := locateManager(request)
	if err != nil {
		return empty(remediation.AdapterRejected, err.Error()), nil
	}
	if manager == "" {
		return empty(remediation.AdapterUnavailable, "no supported Python manifest or lockfile is available"), nil
	}
	data, err := os.ReadFile(manifest)
	if err != nil {
		return empty(remediation.AdapterRejected, "Python manifest is unreadable"), nil
	}
	constraint, direct := findConstraint(string(data), request.Component.Name)
	owner := request.Component.Name
	if !direct {
		owner, direct = ownerFromPath(string(data), request.Finding.DependencyPath)
		if !direct {
			return unknown("direct dependency owner is not identifiable", request, lockfile), nil
		}
		constraint, _ = findConstraint(string(data), owner)
	}
	planned := candidates(request, constraint, direct)
	files := []remediation.AffectedFile{{Path: filepath.ToSlash(filepath.Join(request.Workspace.Path, filepath.Base(manifest))), Kind: "manifest"}, {Path: lockfile, Kind: "lockfile"}}
	risks := risks(planned)
	if len(planned) == 0 {
		return remediation.PlanningEvidence{State: remediation.AdapterUnknown, Reason: "no retained Python candidate is available", Candidates: planned, AffectedFiles: files, Commands: []remediation.Command{}, Risks: risks, Assumptions: []remediation.Assumption{}, Verification: []remediation.Verification{}}, nil
	}
	viable := firstViable(planned)
	if viable == nil {
		return remediation.PlanningEvidence{State: remediation.AdapterUnknown, Reason: "all retained Python candidates require constraint changes", Candidates: planned, AffectedFiles: files, Commands: []remediation.Command{}, Risks: risks, Assumptions: []remediation.Assumption{}, Verification: []remediation.Verification{}}, nil
	}
	return remediation.PlanningEvidence{State: remediation.AdapterSupported, Candidates: planned, AffectedFiles: files, Commands: []remediation.Command{{Executable: manager, Arguments: command(manager, owner, viable.Version), WorkingDirectory: request.Workspace.Path}}, Risks: risks, Assumptions: []remediation.Assumption{}, Verification: []remediation.Verification{}}, nil
}
func manifestPath(r remediation.PlanningRequest) (string, error) {
	root, err := remediation.CanonicalRepositoryRoot(r.Repository.Root)
	if err != nil {
		return "", err
	}
	for _, name := range []string{"pyproject.toml", "requirements.txt", "setup.cfg"} {
		p := filepath.Join(root, filepath.FromSlash(r.Workspace.Path), name)
		if _, e := os.Stat(p); e == nil {
			resolved, e := filepath.EvalSymlinks(p)
			if e != nil {
				return "", e
			}
			return resolved, nil
		}
	}
	return "", errors.New("Python manifest is unavailable")
}
func locateManager(r remediation.PlanningRequest) (string, string, error) {
	root, err := remediation.CanonicalRepositoryRoot(r.Repository.Root)
	if err != nil {
		return "", "", err
	}
	dir := filepath.Join(root, filepath.FromSlash(r.Workspace.Path))
	if r.CurrentState.Lockfile != "" {
		p := filepath.Join(root, filepath.FromSlash(r.CurrentState.Lockfile))
		if _, e := os.Stat(p); e != nil {
			return "", "", errors.New("retained Python lockfile is unavailable")
		}
		return managerFor(filepath.Base(p)), filepath.ToSlash(r.CurrentState.Lockfile), nil
	}
	for _, v := range []struct{ file, manager string }{{"uv.lock", "uv"}, {"poetry.lock", "poetry"}, {"requirements.txt", "pip"}} {
		if _, e := os.Stat(filepath.Join(dir, v.file)); e == nil {
			return v.manager, filepath.ToSlash(filepath.Join(r.Workspace.Path, v.file)), nil
		}
	}
	return "", "", nil
}
func managerFor(name string) string {
	switch name {
	case "uv.lock":
		return "uv"
	case "poetry.lock":
		return "poetry"
	case "requirements.txt":
		return "pip"
	}
	return ""
}

var depPattern = regexp.MustCompile(`(?m)(?:^|["'\s])([A-Za-z0-9_.-]+)\s*(==|~=|>=|<=|>|<|\^)?\s*([0-9]+(?:\.[0-9]+){0,2})`)

func findConstraint(data, name string) (string, bool) {
	for _, m := range depPattern.FindAllStringSubmatch(data, -1) {
		if strings.EqualFold(m[1], name) {
			return m[2] + m[3], true
		}
	}
	return "", false
}
func ownerFromPath(data string, path []string) (string, bool) {
	for _, name := range path {
		if _, ok := findConstraint(data, name); ok {
			return name, true
		}
	}
	return "", false
}

type ver struct{ a, b, c int }

func pv(s string) (ver, bool) {
	p := strings.Split(strings.TrimSpace(s), ".")
	n := [3]int{}
	if len(p) < 1 || len(p) > 3 {
		return ver{}, false
	}
	for i := range p {
		v, e := strconv.Atoi(p[i])
		if e != nil {
			return ver{}, false
		}
		n[i] = v
	}
	return ver{n[0], n[1], n[2]}, true
}
func newer(a ver, s string) bool {
	b, ok := pv(s)
	return ok && (a.a > b.a || a.a == b.a && (a.b > b.b || a.b == b.b && a.c > b.c))
}
func satisfies(c string, v ver) bool {
	c = strings.TrimSpace(c)
	if c == "" {
		return true
	}
	op := "=="
	for _, x := range []string{"~=", "==", ">=", "<=", "^", ">", "<"} {
		if strings.HasPrefix(c, x) {
			op = x
			c = strings.TrimSpace(strings.TrimPrefix(c, x))
			break
		}
	}
	b, ok := pv(c)
	if !ok {
		return false
	}
	switch op {
	case ">=":
		return newer(v, fmtv(b)) || v == b
	case "<=":
		return !newer(v, fmtv(b))
	case ">":
		return newer(v, fmtv(b))
	case "<":
		return newer(b, fmtv(v))
	case "^":
		return v.a == b.a && (newer(v, fmtv(b)) || v == b)
	case "~=":
		return v.a == b.a && v.b == b.b && (newer(v, fmtv(b)) || v == b)
	default:
		return v == b
	}
}
func fmtv(v ver) string { return strconv.Itoa(v.a) + "." + strconv.Itoa(v.b) + "." + strconv.Itoa(v.c) }
func candidates(r remediation.PlanningRequest, c string, direct bool) []remediation.Candidate {
	out := []remediation.Candidate{}
	for _, raw := range r.Finding.FixedVersions {
		v, ok := pv(raw)
		if !ok || !newer(v, r.Component.Version) {
			continue
		}
		x := remediation.Candidate{ID: raw, Version: raw, Direct: direct, MajorChange: v.a != parseMajor(r.Component.Version), Evidence: remediation.CompatibilityEvidence{Unknown: []string{"lockfile_resolution_unknown", "runtime_compatibility_unknown"}}}
		if c != "" && !satisfies(c, v) {
			x.State = remediation.CandidateRejected
			x.RejectionCode = "CONSTRAINT_CHANGE_REQUIRED"
			x.Reason = "candidate requires changing the Python constraint"
		} else {
			x.State = remediation.CandidateViable
		}
		out = append(out, x)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Version < out[j].Version })
	return out
}
func firstViable(values []remediation.Candidate) *remediation.Candidate {
	for i := range values {
		if values[i].State == remediation.CandidateViable {
			return &values[i]
		}
	}
	return nil
}
func parseMajor(s string) int { v, _ := pv(s); return v.a }
func risks(c []remediation.Candidate) []remediation.Risk {
	r := []remediation.Risk{}
	for _, x := range c {
		if x.MajorChange {
			r = append(r, remediation.Risk{Code: "MAJOR_VERSION_CHANGE", Severity: "high", Detected: true})
		}
		if x.State == remediation.CandidateRejected {
			r = append(r, remediation.Risk{Code: "CONSTRAINT_CHANGE_REQUIRED", Severity: "high", Detected: true})
		}
	}
	return r
}
func command(m, n, v string) []string {
	switch m {
	case "uv":
		return []string{"lock", "--upgrade-package", n + ">=" + v}
	case "poetry":
		return []string{"update", n}
	default:
		return []string{"install", n + "==" + v}
	}
}
func empty(s remediation.AdapterState, r string) remediation.PlanningEvidence {
	return remediation.PlanningEvidence{State: s, Reason: r, Candidates: []remediation.Candidate{}, AffectedFiles: []remediation.AffectedFile{}, Commands: []remediation.Command{}, Risks: []remediation.Risk{}, Assumptions: []remediation.Assumption{}, Verification: []remediation.Verification{}}
}
func unknown(reason string, r remediation.PlanningRequest, lock string) remediation.PlanningEvidence {
	return remediation.PlanningEvidence{State: remediation.AdapterUnknown, Reason: reason, Candidates: []remediation.Candidate{}, AffectedFiles: []remediation.AffectedFile{{Path: lock, Kind: "lockfile"}}, Commands: []remediation.Command{}, Risks: []remediation.Risk{}, Assumptions: []remediation.Assumption{}, Verification: []remediation.Verification{}}
}
