package main

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/geoffrey-xiao/deprail/internal/adapter"
	"github.com/geoffrey-xiao/deprail/internal/app"
	"github.com/geoffrey-xiao/deprail/internal/discovery"
	"github.com/geoffrey-xiao/deprail/internal/remediation"
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
