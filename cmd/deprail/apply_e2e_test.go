package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/geoffrey-xiao/deprail/internal/app"
	"github.com/geoffrey-xiao/deprail/internal/remediation"
	"github.com/geoffrey-xiao/deprail/internal/remediation/evidence"
)

type applyE2ECase struct {
	name          string
	fixture       string
	workspacePath string
	executable    string
	mutationFile  string
	verifyArgs    []string
	marker        bool
	mutationArgs  []string
}

func applyE2EExecutable(name string) string {
	if runtime.GOOS != "windows" {
		return name
	}
	switch name {
	case "npm", "mvn":
		return name + ".cmd"
	case "pip":
		return name + ".exe"
	default:
		return name
	}
}

func TestApplyEndToEndRepresentativeFixtures(t *testing.T) {
	cases := []applyE2ECase{
		{name: "javascript", fixture: "npm-basic", executable: "npm", mutationFile: "package.json", verifyArgs: []string{"test"}, mutationArgs: []string{"install"}},
		{name: "python", fixture: "python-requirements", executable: "pip", mutationFile: "requirements.txt", verifyArgs: []string{"check"}, mutationArgs: []string{"install"}},
		{name: "nested javascript", fixture: "npm-basic", workspacePath: "services/web", executable: "npm", mutationFile: "services/web/package.json", verifyArgs: []string{"test"}, mutationArgs: []string{"install"}},
		{name: "java", fixture: "java-maven", executable: "mvn", mutationFile: "pom.xml", verifyArgs: []string{"test"}, mutationArgs: []string{"install"}},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			repo, commit := prepareApplyE2ERepository(t, test)
			tools := prepareApplyE2ETools(t, test.executable)
			t.Setenv("PATH", tools+string(os.PathListSeparator)+os.Getenv("PATH"))
			planPath, approvalPath := writeApplyE2EInputs(t, repo, commit, test)
			beforeState, err := app.CurrentRepositoryState(repo)
			if err != nil {
				t.Fatal(err)
			}
			beforeFile, err := os.ReadFile(filepath.Join(repo, test.mutationFile))
			if err != nil {
				t.Fatal(err)
			}

			result, code, stderr := runApplyE2E(t, context.Background(), repo, planPath, approvalPath)
			if code != 0 || result.Outcome != "applied" {
				t.Fatalf("apply code=%d outcome=%q stderr=%q result=%#v", code, result.Outcome, stderr, result)
			}
			if result.EvidencePath == "" || result.EvidenceDigest == "" || result.CleanupStatus != "succeeded" {
				t.Fatalf("apply evidence/cleanup = %#v", result)
			}
			if len(result.Transitions) != 1 || result.Transitions[0].State != "residual" {
				t.Fatalf("transitions = %#v", result.Transitions)
			}
			if _, err := os.Stat(result.EvidencePath); err != nil {
				t.Fatalf("evidence path: %v", err)
			}
			data, err := os.ReadFile(result.EvidencePath)
			if err != nil {
				t.Fatal(err)
			}
			var record evidence.Record
			if err := json.Unmarshal(data, &record); err != nil {
				t.Fatal(err)
			}
			if err := record.Validate(); err != nil {
				t.Fatal(err)
			}
			if record.Outcome != evidence.Applied || record.WorkspaceID == "" || record.WorkspacePath == "" || record.BeforeDigest == "" || record.AfterDigest == "" {
				t.Fatalf("evidence record = %#v", record)
			}
			afterState, err := app.CurrentRepositoryState(repo)
			if err != nil {
				t.Fatal(err)
			}
			if afterState != beforeState {
				t.Fatalf("caller state changed: before=%s after=%s", beforeState, afterState)
			}
			afterFile, err := os.ReadFile(filepath.Join(repo, test.mutationFile))
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(beforeFile, afterFile) {
				t.Fatalf("caller file changed: before=%q after=%q", beforeFile, afterFile)
			}

			var reuseOut, reuseErr bytes.Buffer
			reuseCode := runFixApplyContext(context.Background(), []string{"--plan", planPath, "--approval", approvalPath, "--root", repo, "--format", "json"}, &reuseOut, &reuseErr)
			reuseStderr := reuseErr.String()
			if reuseCode != 3 || !strings.Contains(reuseStderr, "APPROVAL_USED") {
				t.Fatalf("approval reuse code=%d stderr=%q", reuseCode, reuseStderr)
			}
		})
	}
}

func TestApplyEndToEndDryRunPreservesApproval(t *testing.T) {
	test := applyE2ECase{name: "dry-run", fixture: "npm-basic", executable: "npm", mutationFile: "package.json", verifyArgs: []string{"test"}, mutationArgs: []string{"install"}}
	repo, commit := prepareApplyE2ERepository(t, test)
	tools := prepareApplyE2ETools(t, test.executable)
	t.Setenv("PATH", tools+string(os.PathListSeparator)+os.Getenv("PATH"))
	planPath, approvalPath := writeApplyE2EInputs(t, repo, commit, test)
	beforeState, err := app.CurrentRepositoryState(repo)
	if err != nil {
		t.Fatal(err)
	}
	beforeFile, err := os.ReadFile(filepath.Join(repo, test.mutationFile))
	if err != nil {
		t.Fatal(err)
	}

	result, code, stderr, err := invokeApplyE2E(context.Background(), repo, []string{"--plan", planPath, "--approval", approvalPath, "--root", repo, "--dry-run", "--format", "json"})
	if err != nil {
		t.Fatal(err)
	}
	if code != 0 || result.Outcome != "dry_run" || result.EvidencePath != "" {
		t.Fatalf("dry-run code=%d outcome=%q evidence=%q stderr=%q", code, result.Outcome, result.EvidencePath, stderr)
	}
	if _, err := os.Stat(filepath.Join(repo, ".deprail", "evidence")); !os.IsNotExist(err) {
		t.Fatalf("dry-run wrote evidence: %v", err)
	}
	afterState, err := app.CurrentRepositoryState(repo)
	if err != nil {
		t.Fatal(err)
	}
	if afterState != beforeState {
		t.Fatalf("dry-run changed caller state: before=%s after=%s", beforeState, afterState)
	}
	afterFile, err := os.ReadFile(filepath.Join(repo, test.mutationFile))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(beforeFile, afterFile) {
		t.Fatalf("dry-run changed caller file: before=%q after=%q", beforeFile, afterFile)
	}

	applied, applyCode, applyStderr := runApplyE2E(t, context.Background(), repo, planPath, approvalPath)
	if applyCode != 0 || applied.Outcome != "applied" {
		t.Fatalf("approval was consumed by dry-run: code=%d outcome=%q stderr=%q", applyCode, applied.Outcome, applyStderr)
	}
}

func TestApplyEndToEndFailureBoundaries(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		marker      bool
		context     func() (context.Context, context.CancelFunc)
		wantOutcome string
	}{
		{name: "hostile argument is direct argv", args: []string{"install", "$(touch deprail-pwned)"}, wantOutcome: "applied"},
		{name: "malformed rescan", args: []string{"install"}, marker: true, wantOutcome: "partial"},
		{name: "cancellation", args: []string{"install", "--deprail-e2e-sleep"}, context: func() (context.Context, context.CancelFunc) { return context.WithCancel(context.Background()) }, wantOutcome: "cancelled"},
		{name: "timeout", args: []string{"install", "--deprail-e2e-sleep"}, context: func() (context.Context, context.CancelFunc) {
			return context.WithTimeout(context.Background(), time.Second)
		}, wantOutcome: "partial"},
		{name: "mutation failure", args: []string{"install", "--deprail-e2e-fail"}, wantOutcome: "partial"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			caseData := applyE2ECase{name: test.name, fixture: "npm-basic", executable: "npm", mutationFile: "package.json", verifyArgs: []string{"test"}, mutationArgs: test.args, marker: test.marker}
			repo, commit := prepareApplyE2ERepository(t, caseData)
			tools := prepareApplyE2ETools(t, caseData.executable)
			t.Setenv("PATH", tools+string(os.PathListSeparator)+os.Getenv("PATH"))
			planPath, approvalPath := writeApplyE2EInputs(t, repo, commit, caseData)
			ctx := context.Background()
			var cancel context.CancelFunc
			if test.context != nil {
				ctx, cancel = test.context()
				defer cancel()
			}
			if test.name == "cancellation" {
				resultCh := make(chan struct {
					result applyResult
					code   int
					stderr string
					err    error
				})
				go func() {
					result, code, stderr, err := invokeApplyE2E(ctx, repo, []string{"--plan", planPath, "--approval", approvalPath, "--root", repo, "--format", "json"})
					resultCh <- struct {
						result applyResult
						code   int
						stderr string
						err    error
					}{result, code, stderr, err}
				}()
				time.Sleep(500 * time.Millisecond)
				cancel()
				outcome := <-resultCh
				if outcome.err != nil {
					t.Fatal(outcome.err)
				}
				if outcome.code != 3 || outcome.result.Outcome != test.wantOutcome {
					t.Fatalf("apply code=%d outcome=%q stderr=%q", outcome.code, outcome.result.Outcome, outcome.stderr)
				}
				if outcome.result.EvidencePath == "" || outcome.result.CleanupStatus != "succeeded" {
					t.Fatalf("failure evidence/cleanup = %#v", outcome.result)
				}
				return
			}
			result, code, stderr := runApplyE2E(t, ctx, repo, planPath, approvalPath)
			wantCode := 3
			if test.wantOutcome == "applied" {
				wantCode = 0
			}
			if code != wantCode || result.Outcome != test.wantOutcome {
				t.Fatalf("apply code=%d outcome=%q stderr=%q result=%#v", code, result.Outcome, stderr, result)
			}
			if result.EvidencePath == "" || result.CleanupStatus != "succeeded" {
				t.Fatalf("failure evidence/cleanup = %#v", result)
			}
			if test.name == "hostile argument is direct argv" {
				if _, err := os.Stat(filepath.Join(repo, "deprail-pwned")); !os.IsNotExist(err) {
					t.Fatalf("shell metacharacter escaped argv: %v", err)
				}
			}
		})
	}
}

func prepareApplyE2ERepository(t *testing.T, test applyE2ECase) (string, string) {
	t.Helper()
	repo := t.TempDir()
	destination := repo
	if test.workspacePath != "" && test.workspacePath != "." {
		destination = filepath.Join(repo, test.workspacePath)
		if err := os.MkdirAll(destination, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	copyDirectory(t, fixturePath(t, test.fixture), destination)
	if test.marker {
		writeTestFile(t, filepath.Join(repo, "malformed-scanner.marker"), "malformed\n")
	}
	gitApplyE2E(t, repo, "init", "-q")
	gitApplyE2E(t, repo, "add", ".")
	gitApplyE2E(t, repo, "-c", "user.name=DepRail Test", "-c", "user.email=test@example.invalid", "commit", "-qm", "initial")
	return repo, strings.TrimSpace(string(gitApplyE2E(t, repo, "rev-parse", "HEAD")))
}

func prepareApplyE2ETools(t *testing.T, executable string) string {
	t.Helper()
	bin := t.TempDir()
	data, err := os.ReadFile(mustExecutable(t))
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{executable, "osv-scanner"} {
		targetName := name
		if runtime.GOOS == "windows" {
			targetName += ".exe"
		}
		path := filepath.Join(bin, targetName)
		if err := os.WriteFile(path, data, 0o755); err != nil {
			t.Fatal(err)
		}
		if runtime.GOOS == "windows" && (name == "npm" || name == "mvn") {
			wrapper := fmt.Sprintf("@echo off\r\n\"%%~dp0%s\" %%*\r\n", targetName)
			if err := os.WriteFile(filepath.Join(bin, name+".cmd"), []byte(wrapper), 0o644); err != nil {
				t.Fatal(err)
			}
		}
	}
	return bin
}

func writeApplyE2EInputs(t *testing.T, repo, commit string, test applyE2ECase) (string, string) {
	t.Helper()
	workspacePath := test.workspacePath
	if workspacePath == "" {
		workspacePath = "."
	}
	plan := remediation.Plan{
		SchemaVersion:      remediation.SchemaVersion,
		CreatedFrom:        remediation.CreatedFrom{ReportDigest: "report-digest", SourceScanID: "scan-e2e", Scanner: "osv-scanner", RepositoryState: commit},
		RepositoryIdentity: remediation.RepositoryIdentity{Root: repo, Repository: filepath.Base(repo), Revision: commit},
		WorkspaceIdentity:  remediation.WorkspaceIdentity{ID: "root", Path: workspacePath},
		FindingIdentity:    remediation.FindingIdentity{StableKey: "OSV-CALLER", CurrentVersion: "1.0.0", FixedVersions: []string{"1.0.1"}},
		Component:          remediation.Component{PURL: "pkg:npm/target-only@1.0.0", Name: "target-only", Version: "1.0.0"},
		CurrentState:       remediation.CurrentState{Vulnerable: true},
		AffectedFiles:      []remediation.AffectedFile{{Path: test.mutationFile, Kind: "manifest", Effect: "update"}},
		Commands:           []remediation.Command{{Executable: applyE2EExecutable(test.executable), Arguments: test.mutationArgs, WorkingDirectory: workspacePath}},
		Verification:       []remediation.Verification{{ID: "tests", Command: remediation.Command{Executable: applyE2EExecutable(test.executable), Arguments: test.verifyArgs, WorkingDirectory: workspacePath}}},
		Provenance:         remediation.Provenance{ArtifactDigests: []string{}, Sources: []string{"e2e-test"}},
	}
	plan.Canonicalize()
	plan.PlanID = remediation.StablePlanID(plan)
	approval, err := remediation.NewApproval(plan, repo, time.Now().UTC().Add(time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	planPath := filepath.Join(root, "plan.json")
	approvalPath := filepath.Join(root, "approval.json")
	writeJSONApplyE2E(t, planPath, plan)
	writeJSONApplyE2E(t, approvalPath, approval)
	return planPath, approvalPath
}

func runApplyE2E(t *testing.T, ctx context.Context, repo, planPath, approvalPath string) (applyResult, int, string) {
	t.Helper()
	result, code, stderr, err := invokeApplyE2E(ctx, repo, []string{"--plan", planPath, "--approval", approvalPath, "--root", repo, "--format", "json"})
	if err != nil {
		t.Fatal(err)
	}
	return result, code, stderr
}

func invokeApplyE2E(ctx context.Context, repo string, args []string) (applyResult, int, string, error) {
	var stdout, stderr bytes.Buffer
	code := runFixApplyContext(ctx, args, &stdout, &stderr)
	var result applyResult
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		return applyResult{}, code, stderr.String(), err
	}
	return result, code, stderr.String(), nil
}

func copyDirectory(t *testing.T, source, destination string) {
	t.Helper()
	err := filepath.WalkDir(source, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		if relative == "." {
			return nil
		}
		target := filepath.Join(destination, relative)
		if entry.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, 0o644)
	})
	if err != nil {
		t.Fatal(err)
	}
}

func gitApplyE2E(t *testing.T, repo string, args ...string) []byte {
	t.Helper()
	command := exec.Command("git", append([]string{"-C", repo}, args...)...)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, output)
	}
	return output
}

func writeJSONApplyE2E(t *testing.T, path string, value any) {
	t.Helper()
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, append(data, '\n'), 0o600); err != nil {
		t.Fatal(err)
	}
}
