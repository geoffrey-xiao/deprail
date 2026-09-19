package python

import (
	"context"
	"errors"
	"os"
	"path/filepath"
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
		return unknown("transitive fixed version cannot be mapped to a direct owner version", request, lockfile), nil
	}
	_ = owner
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
			if !within(root, resolved) {
				return "", errors.New("Python manifest escapes repository root")
			}
			return resolved, nil
		}
	}
	return "", errors.New("Python manifest is unavailable")
}
func within(root, path string) bool {
	rel, err := filepath.Rel(root, path)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator))
}
func locateManager(r remediation.PlanningRequest) (string, string, error) {
	root, err := remediation.CanonicalRepositoryRoot(r.Repository.Root)
	if err != nil {
		return "", "", err
	}
	dir := filepath.Join(root, filepath.FromSlash(r.Workspace.Path))
	if r.CurrentState.Lockfile != "" {
		p := filepath.Join(root, filepath.FromSlash(r.CurrentState.Lockfile))
		if !within(root, p) {
			return "", "", errors.New("retained Python lockfile escapes repository root")
		}
		if _, e := os.Stat(p); e != nil {
			return "", "", errors.New("retained Python lockfile is unavailable")
		}
		manager := managerFor(filepath.Base(p))
		if manager == "" {
			return "", "", errors.New("retained Python lockfile type is unsupported")
		}
		if err := validateLockfile(p, manager); err != nil {
			return "", "", err
		}
		return manager, filepath.ToSlash(r.CurrentState.Lockfile), nil
	}
	for _, v := range []struct{ file, manager string }{{"uv.lock", "uv"}, {"poetry.lock", "poetry"}, {"requirements.txt", "pip"}} {
		p := filepath.Join(dir, v.file)
		if _, e := os.Stat(p); e == nil {
			if err := validateLockfile(p, v.manager); err != nil {
				return "", "", err
			}
			return v.manager, filepath.ToSlash(filepath.Join(r.Workspace.Path, v.file)), nil
		}
	}
	return "", "", nil
}
func validateLockfile(path, manager string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return errors.New("Python lockfile is unreadable")
	}
	text := strings.TrimSpace(string(data))
	if text == "" {
		return errors.New("Python lockfile is empty")
	}
	switch manager {
	case "uv":
		if !strings.Contains(text, "[[package]]") && !strings.Contains(text, "version =") {
			return errors.New("Python uv lockfile is malformed")
		}
	case "poetry":
		if !strings.Contains(text, "[[package]]") && !strings.Contains(text, "[metadata]") {
			return errors.New("Python Poetry lockfile is malformed")
		}
	case "pip":
		for _, line := range strings.Split(text, "\n") {
			line = strings.TrimSpace(line)
			if line != "" && !strings.HasPrefix(line, "#") {
				return nil
			}
		}
		return errors.New("Python requirements lockfile is malformed")
	}
	return nil
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

func findConstraint(data, name string) (string, bool) {
	for _, line := range strings.Split(data, "\n") {
		clean := strings.TrimSpace(strings.SplitN(line, "#", 2)[0])
		lower, target := strings.ToLower(clean), strings.ToLower(name)
		for start := 0; ; {
			i := strings.Index(lower[start:], target)
			if i < 0 {
				break
			}
			i += start
			beforeOK := i == 0 || !isNameChar(lower[i-1])
			end := i + len(target)
			afterOK := end == len(lower) || !isNameChar(lower[end])
			if beforeOK && afterOK {
				expr := strings.TrimSpace(clean[end:])
				expr = strings.Trim(expr, " \t\"'[](),")
				if expr != "" && strings.ContainsAny(expr, "=<>~^") {
					return expr, true
				}
			}
			start = end
		}
	}
	return "", false
}
func isNameChar(b byte) bool {
	return b == '_' || b == '-' || b == '.' || b >= 'a' && b <= 'z' || b >= '0' && b <= '9'
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
	parts := strings.Split(c, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		op := "=="
		for _, candidate := range []string{"~=", "==", ">=", "<=", "^", ">", "<"} {
			if strings.HasPrefix(part, candidate) {
				op = candidate
				part = strings.TrimSpace(strings.TrimPrefix(part, candidate))
				break
			}
		}
		b, ok := pv(part)
		if !ok {
			return false
		}
		ok = false
		switch op {
		case ">=":
			ok = newer(v, fmtv(b)) || v == b
		case "<=":
			ok = !newer(v, fmtv(b))
		case ">":
			ok = newer(v, fmtv(b))
		case "<":
			ok = newer(b, fmtv(v))
		case "^":
			ok = v.a == b.a && (newer(v, fmtv(b)) || v == b)
		case "~=":
			ok = v.a == b.a && v.b == b.b && (newer(v, fmtv(b)) || v == b)
		default:
			ok = v == b
		}
		if !ok {
			return false
		}
	}
	return true
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
	sort.Slice(out, func(i, j int) bool {
		a, _ := pv(out[i].Version)
		b, _ := pv(out[j].Version)
		if a != b {
			return newer(b, fmtv(a))
		}
		return out[i].Version < out[j].Version
	})
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
