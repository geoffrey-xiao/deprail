package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/geoffrey-xiao/deprail/internal/adapter"
	"github.com/geoffrey-xiao/deprail/internal/app"
	"github.com/geoffrey-xiao/deprail/internal/discovery"
	"github.com/geoffrey-xiao/deprail/internal/process"
	"github.com/geoffrey-xiao/deprail/internal/remediation"
	"github.com/geoffrey-xiao/deprail/internal/remediation/evidence"
	"github.com/geoffrey-xiao/deprail/internal/remediation/isolation"
	"github.com/geoffrey-xiao/deprail/internal/remediation/verification"
)

func TestApplyVerificationCommandsResolvesSupportedTools(t *testing.T) {
	dir := t.TempDir()
	name := "npm"
	if runtime.GOOS == "windows" {
		name = "npm.exe"
	}
	tool := filepath.Join(dir, name)
	data, err := os.ReadFile(os.Args[0])
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(tool, data, 0o700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir)

	commands, err := applyVerificationCommands(remediation.Plan{Verification: []remediation.Verification{{
		ID:      "tests",
		Command: remediation.Command{Executable: "npm", Arguments: []string{"test"}, WorkingDirectory: "services/api"},
		Reason:  "run repository tests",
	}}})
	if err != nil {
		t.Fatal(err)
	}
	if len(commands) != 1 {
		t.Fatalf("commands = %#v", commands)
	}
	command := commands[0]
	if command.ID != "tests" || command.Kind != verification.Test || command.Path != tool || !command.Enabled || command.WorkingDirectory != "services/api" {
		t.Fatalf("command = %#v", command)
	}
	if len(command.Args) != 1 || command.Args[0] != "test" {
		t.Fatalf("args = %#v", command.Args)
	}
}

func TestApplyVerificationCommandsRejectsMissingTools(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	_, err := applyVerificationCommands(remediation.Plan{Verification: []remediation.Verification{{
		ID:      "tests",
		Command: remediation.Command{Executable: "npm", Arguments: []string{"test"}, WorkingDirectory: "."},
	}}})
	if err == nil {
		t.Fatal("expected missing verification tool error")
	}
}

func TestRescanApplyWorkspaceRejectsMissingScanner(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	report, err := rescanApplyWorkspace(t.Context(), t.TempDir())
	if err == nil {
		t.Fatal("expected missing scanner error")
	}
	if report.Status != discovery.Failed || len(report.Errors) != 1 {
		t.Fatalf("report = %#v", report)
	}
}

func TestApplyTransitionClassificationWithholdsIncompleteResolution(t *testing.T) {
	plan := remediation.Plan{FindingIdentity: remediation.FindingIdentity{StableKey: "finding"}}
	complete, err := verification.Classify([]remediation.ReportFinding{planFinding(plan)}, nil, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(complete) != 1 || complete[0].State != verification.Resolved {
		t.Fatalf("complete transitions = %#v", complete)
	}
	incomplete, err := verification.Classify([]remediation.ReportFinding{planFinding(plan)}, nil, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(incomplete) != 1 || incomplete[0].State != verification.Unknown {
		t.Fatalf("incomplete transitions = %#v", incomplete)
	}
}

func TestRescanFindingsForPlanUsesPlanningKeyAndFiltersUnrelated(t *testing.T) {
	plan := remediation.Plan{FindingIdentity: remediation.FindingIdentity{StableKey: "OSV-1"}}
	report := app.ScanReport{
		Findings: []adapter.Finding{
			{TargetID: "OSV-2", Component: "other", Version: "1", PURL: "pkg:npm/other@1", WorkspaceID: "root", WorkspacePath: "."},
			{TargetID: "OSV-1", Component: "demo", Version: "1", PURL: "pkg:npm/demo@1", WorkspaceID: "root", WorkspacePath: "."},
		},
	}
	findings := rescanFindingsForPlan(report, plan)
	if len(findings) != 1 || findings[0].StableKey != "OSV-1" {
		t.Fatalf("findings = %#v", findings)
	}
	transitions, err := verification.Classify([]remediation.ReportFinding{planFinding(plan)}, findings, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(transitions) != 1 || transitions[0].State != verification.Residual {
		t.Fatalf("transitions = %#v", transitions)
	}
}

func TestRecordCleanupOutcomeExposesPartialFailure(t *testing.T) {
	result := applyResult{Outcome: "failed", Diagnostics: []string{}}
	recordCleanupOutcome(&result, "failed", isolation.CleanupResult{
		GitRemoved:        true,
		FilesystemRemoved: false,
		Retryable:         true,
		Err:               errors.New("git worktree remove failed"),
	})
	if result.Outcome != "cleanup_failed" || result.CleanupStatus != "partial" || !result.CleanupGitRemoved || result.CleanupFilesystemRemoved || !result.CleanupRetryable {
		t.Fatalf("result = %#v", result)
	}
	if len(result.Diagnostics) != 1 || result.Diagnostics[0] != "cleanup failed after failed: git worktree remove failed" {
		t.Fatalf("diagnostics = %#v", result.Diagnostics)
	}
}

func TestApplyFailureOutcomeDistinguishesCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if got := applyFailureOutcome(ctx, errors.New("cancelled")); got != "cancelled" {
		t.Fatalf("cancelled context outcome = %q", got)
	}
	if got := applyFailureOutcome(context.Background(), &process.Error{Code: process.ErrCancelled}); got != "cancelled" {
		t.Fatalf("cancelled process outcome = %q", got)
	}
	if got := applyFailureOutcome(context.Background(), errors.New("failed")); got != "failed" {
		t.Fatalf("failed outcome = %q", got)
	}
	if got := applyOperationOutcome(context.Background(), errors.New("failed"), true); got != "partial" {
		t.Fatalf("partial failure outcome = %q", got)
	}
	if got := applyOperationOutcome(context.Background(), &process.Error{Code: process.ErrCancelled}, true); got != "cancelled" {
		t.Fatalf("partial cancellation outcome = %q", got)
	}
}

func TestWriteApplyResultReturnsFailureExitForTerminalFailures(t *testing.T) {
	for _, outcome := range []string{"partial", "failed", "cancelled", "cleanup_failed"} {
		var stdout, stderr bytes.Buffer
		if code := writeApplyResult(applyResult{Outcome: outcome, Diagnostics: []string{}}, "json", &stdout, &stderr); code != 3 {
			t.Fatalf("outcome %q exit code = %d", outcome, code)
		}
	}
	var stdout, stderr bytes.Buffer
	if code := writeApplyResult(applyResult{Outcome: "applied", Diagnostics: []string{}}, "json", &stdout, &stderr); code != 0 {
		t.Fatalf("applied exit code = %d", code)
	}
}

func TestPersistApplyEvidenceRedactsAndRetainsIdentity(t *testing.T) {
	plan := remediation.Plan{
		WorkspaceIdentity: remediation.WorkspaceIdentity{ID: "workspace"},
		FindingIdentity:   remediation.FindingIdentity{StableKey: "finding"},
		AffectedFiles:     []remediation.AffectedFile{{Path: "package.json", Kind: "manifest", Effect: "update"}},
		Commands:          []remediation.Command{{Executable: "npm", WorkingDirectory: "."}},
	}
	state, err := newApplyEvidenceState(plan)
	if err != nil {
		t.Fatal(err)
	}
	state.BeforeDigest = "before"
	state.AfterDigest = "after"
	state.FindingAfterDigest = "finding-after"
	state.VerificationStatus = "complete"
	state.RescanStatus = "complete"
	state.Commands = append(state.Commands, applyEvidenceCommand("mutation:0", "/usr/bin/npm", process.Result{
		Stdout: []byte("token=secret"),
		Stderr: []byte("Authorization: Bearer credential"),
	}))
	result := applyResult{
		Outcome:               "partial",
		PlanDigest:            "plan",
		SourceCommit:          "commit",
		SourceRoot:            "/repo",
		RescanArtifactDigests: []string{},
		Diagnostics:           []string{"secret=hidden"},
	}
	path, digest, err := persistApplyEvidence(filepath.Join(t.TempDir(), "evidence"), plan, result, state, "/tmp/worktree", "partial")
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(digest) != 64 || len(data) == 0 {
		t.Fatalf("evidence path/digest/data = %q/%q/%q", path, digest, data)
	}
	if len(data) == 0 || bytes.Contains(data, []byte("token=secret")) || bytes.Contains(data, []byte("Bearer credential")) || bytes.Contains(data, []byte("secret=hidden")) {
		t.Fatalf("sensitive evidence persisted: %s", data)
	}
	var record evidence.Record
	if err := json.Unmarshal(data, &record); err != nil {
		t.Fatal(err)
	}
	if err := record.Validate(); err != nil {
		t.Fatal(err)
	}
	if record.WorkspaceID != "workspace" || record.WorkspacePath != "/tmp/worktree" || record.Outcome != evidence.Partial || record.Cleanup != "partial" {
		t.Fatalf("record identity/outcome = %#v", record)
	}
}
