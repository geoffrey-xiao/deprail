package java

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/geoffrey-xiao/deprail/internal/remediation"
)

type Adapter struct{}

func (Adapter) Name() string { return "java" }
func (Adapter) Assess(_ context.Context, request remediation.PlanningRequest) (remediation.AdapterAssessment, error) {
	if err := remediation.ValidatePlanningRequest(request); err != nil {
		return remediation.AdapterAssessment{State: remediation.AdapterRejected, Reason: err.Error()}, nil
	}
	manager, _, err := locate(request)
	if err != nil {
		return remediation.AdapterAssessment{State: remediation.AdapterRejected, Reason: err.Error()}, nil
	}
	if manager == "" {
		return remediation.AdapterAssessment{State: remediation.AdapterUnavailable, Reason: "no supported Java build manifest is available"}, nil
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
	manager, manifest, err := locate(request)
	if err != nil {
		return empty(remediation.AdapterRejected, err.Error()), nil
	}
	if manager == "" {
		return empty(remediation.AdapterUnavailable, "no supported Java build manifest is available"), nil
	}
	data, err := os.ReadFile(manifest)
	if err != nil {
		return empty(remediation.AdapterRejected, "Java build manifest is unreadable"), nil
	}
	name := request.Component.Name
	constraint, direct := findDependency(string(data), manager, name)
	owner := name
	if !direct {
		owner, direct = ownerFromPath(string(data), manager, request.Finding.DependencyPath)
		if !direct {
			return unknown("direct dependency owner is not identifiable", request, filepath.ToSlash(filepath.Join(request.Workspace.Path, filepath.Base(manifest)))), nil
		}
		return unknown("transitive fixed version cannot be mapped to a direct owner version", request, filepath.ToSlash(filepath.Join(request.Workspace.Path, filepath.Base(manifest)))), nil
	}
	_ = owner
	candidates := build(request, constraint, direct)
	files := []remediation.AffectedFile{{Path: filepath.ToSlash(filepath.Join(request.Workspace.Path, filepath.Base(manifest))), Kind: "manifest"}}
	risks := risk(candidates)
	if len(candidates) == 0 {
		return remediation.PlanningEvidence{State: remediation.AdapterUnknown, Reason: "no retained Java candidate is available", Candidates: candidates, AffectedFiles: files, Commands: []remediation.Command{}, Risks: risks, Assumptions: []remediation.Assumption{}, Verification: []remediation.Verification{}}, nil
	}
	viable := firstViable(candidates)
	if viable == nil {
		return remediation.PlanningEvidence{State: remediation.AdapterUnknown, Reason: "all retained Java candidates require constraint changes", Candidates: candidates, AffectedFiles: files, Commands: []remediation.Command{}, Risks: risks, Assumptions: []remediation.Assumption{}, Verification: []remediation.Verification{}}, nil
	}
	if manager == "gradle" {
		return remediation.PlanningEvidence{State: remediation.AdapterUnknown, Reason: "no verified Gradle upgrade command is available", Candidates: candidates, AffectedFiles: files, Commands: []remediation.Command{}, Risks: risks, Assumptions: []remediation.Assumption{}, Verification: []remediation.Verification{}}, nil
	}
	return remediation.PlanningEvidence{State: remediation.AdapterSupported, Candidates: candidates, AffectedFiles: files, Commands: []remediation.Command{{Executable: manager, Arguments: command(manager, owner, viable.Version), WorkingDirectory: request.Workspace.Path}}, Risks: risks, Assumptions: []remediation.Assumption{}, Verification: []remediation.Verification{}}, nil
}

type pom struct {
	XMLName      xml.Name `xml:"project"`
	Dependencies []struct {
		Group    string `xml:"groupId"`
		Artifact string `xml:"artifactId"`
		Version  string `xml:"version"`
	} `xml:"dependencies>dependency"`
}

func locate(r remediation.PlanningRequest) (string, string, error) {
	root, err := remediation.CanonicalRepositoryRoot(r.Repository.Root)
	if err != nil {
		return "", "", err
	}
	dir := filepath.Join(root, filepath.FromSlash(r.Workspace.Path))
	for _, v := range []struct{ file, manager string }{{"pom.xml", "mvn"}, {"build.gradle", "gradle"}, {"build.gradle.kts", "gradle"}} {
		p := filepath.Join(dir, v.file)
		if _, e := os.Stat(p); e == nil {
			resolved, e := filepath.EvalSymlinks(p)
			if e != nil {
				return "", "", e
			}
			if !within(root, resolved) {
				return "", "", errors.New("Java build manifest escapes repository root")
			}
			return v.manager, resolved, nil
		}
	}
	return "", "", nil
}
func within(root, path string) bool {
	rel, err := filepath.Rel(root, path)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator))
}

var gradleDep = regexp.MustCompile(`([A-Za-z0-9_.-]+):([A-Za-z0-9_.-]+):([0-9][A-Za-z0-9_.-]*)`)

func findDependency(data, manager, name string) (string, bool) {
	if manager == "mvn" {
		var p pom
		if xml.Unmarshal([]byte(data), &p) != nil {
			return "", false
		}
		for _, d := range p.Dependencies {
			if d.Artifact == name || d.Group+":"+d.Artifact == name {
				return d.Version, true
			}
		}
		return "", false
	}
	for _, m := range gradleDep.FindAllStringSubmatch(data, -1) {
		if m[2] == name || m[1]+":"+m[2] == name {
			return m[3], true
		}
	}
	return "", false
}
func ownerFromPath(data, manager string, path []string) (string, bool) {
	for _, n := range path {
		if _, ok := findDependency(data, manager, n); ok {
			return n, true
		}
	}
	return "", false
}

type ver struct {
	a, b, c int
}

func pv(s string) (ver, bool) {
	s = strings.TrimSpace(s)
	p := strings.Split(s, ".")
	if len(p) < 2 {
		return ver{}, false
	}
	n := [3]int{}
	for i := 0; i < len(p) && i < 3; i++ {
		value := 0
		digits := 0
		for _, r := range p[i] {
			if r < '0' || r > '9' {
				break
			}
			value = value*10 + int(r-'0')
			digits++
		}
		if digits == 0 {
			return ver{}, false
		}
		n[i] = value
	}
	return ver{n[0], n[1], n[2]}, true
}

func newer(a ver, s string) bool {
	b, ok := pv(s)
	return ok && (a.a > b.a || a.a == b.a && (a.b > b.b || a.b == b.b && a.c > b.c))
}
func build(r remediation.PlanningRequest, c string, d bool) []remediation.Candidate {
	out := []remediation.Candidate{}
	for _, raw := range r.Finding.FixedVersions {
		v, ok := pv(raw)
		if !ok || !newer(v, r.Component.Version) {
			continue
		}
		x := remediation.Candidate{ID: raw, Version: raw, Direct: d, MajorChange: v.a != major(r.Component.Version), Evidence: remediation.CompatibilityEvidence{Unknown: []string{"lockfile_resolution_unknown", "runtime_compatibility_unknown"}}}
		if c != "" && raw == c {
			x.State = remediation.CandidateViable
		} else if c != "" {
			x.State = remediation.CandidateRejected
			x.RejectionCode = "CONSTRAINT_CHANGE_REQUIRED"
			x.Reason = "candidate requires changing the Java constraint"
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

func fmtv(v ver) string {
	return fmt.Sprintf("%d.%d.%d", v.a, v.b, v.c)
}
func firstViable(values []remediation.Candidate) *remediation.Candidate {
	for i := range values {
		if values[i].State == remediation.CandidateViable {
			return &values[i]
		}
	}
	return nil
}
func major(s string) int { v, _ := pv(s); return v.a }
func risk(c []remediation.Candidate) []remediation.Risk {
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
	if m == "gradle" {
		return []string{"dependencies", "--configuration", "runtimeClasspath", n + ":" + v}
	}
	return []string{"org.codehaus.mojo:versions-maven-plugin:use-dep-version", "-Dincludes=" + n, "-DdepVersion=" + v}
}
func empty(s remediation.AdapterState, r string) remediation.PlanningEvidence {
	return remediation.PlanningEvidence{State: s, Reason: r, Candidates: []remediation.Candidate{}, AffectedFiles: []remediation.AffectedFile{}, Commands: []remediation.Command{}, Risks: []remediation.Risk{}, Assumptions: []remediation.Assumption{}, Verification: []remediation.Verification{}}
}
func unknown(reason string, r remediation.PlanningRequest, file string) remediation.PlanningEvidence {
	return remediation.PlanningEvidence{State: remediation.AdapterUnknown, Reason: reason, Candidates: []remediation.Candidate{}, AffectedFiles: []remediation.AffectedFile{{Path: file, Kind: "manifest"}}, Commands: []remediation.Command{}, Risks: []remediation.Risk{}, Assumptions: []remediation.Assumption{}, Verification: []remediation.Verification{}}
}
