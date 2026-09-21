package verification

import (
	"path/filepath"
	"testing"

	"github.com/geoffrey-xiao/deprail/internal/remediation"
)

func finding(key string) remediation.ReportFinding { return remediation.ReportFinding{StableKey: key} }

func TestClassifyCompleteScanDeterministically(t *testing.T) {
	result, err := Classify([]remediation.ReportFinding{finding("residual"), finding("resolved")}, []remediation.ReportFinding{finding("introduced"), finding("residual")}, true)
	if err != nil {
		t.Fatal(err)
	}
	want := []FindingTransition{{"introduced", Introduced}, {"residual", Residual}, {"resolved", Resolved}}
	for i := range want {
		if result[i] != want[i] {
			t.Fatalf("result[%d] = %#v, want %#v", i, result[i], want[i])
		}
	}
}

func TestClassifyIncompleteScanWithholdsMissingSafetyClaims(t *testing.T) {
	result, err := Classify([]remediation.ReportFinding{finding("missing"), finding("residual")}, []remediation.ReportFinding{finding("introduced"), finding("residual")}, false)
	if err != nil {
		t.Fatal(err)
	}
	if result[0] != (FindingTransition{"introduced", Introduced}) || result[1] != (FindingTransition{"missing", Unknown}) || result[2] != (FindingTransition{"residual", Residual}) {
		t.Fatalf("result = %#v", result)
	}
}

func TestClassifyRejectsDuplicateOrMissingKeys(t *testing.T) {
	if _, err := Classify([]remediation.ReportFinding{finding("same"), finding("same")}, nil, true); err == nil {
		t.Fatal("expected duplicate rejection")
	}
	if _, err := Classify([]remediation.ReportFinding{{}}, nil, true); err == nil {
		t.Fatal("expected missing key rejection")
	}
}

func TestSelectCommandsFiltersShellsAndSortsChecks(t *testing.T) {
	buildPath, err := filepath.Abs("build")
	if err != nil {
		t.Fatal(err)
	}
	testPath, err := filepath.Abs("test")
	if err != nil {
		t.Fatal(err)
	}
	commands, err := SelectCommands([]Command{
		{ID: "build", Kind: Build, Path: buildPath, WorkingDirectory: ".", Enabled: true},
		{ID: "test", Kind: Test, Path: testPath, WorkingDirectory: ".", Enabled: true},
		{ID: "disabled", Kind: Test, Path: testPath, WorkingDirectory: ".", Enabled: false},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(commands) != 2 || commands[0].ID != "build" || commands[1].ID != "test" {
		t.Fatalf("commands = %#v", commands)
	}
	if _, err := SelectCommands([]Command{{ID: "shell", Kind: Test, Path: "/bin/sh", WorkingDirectory: ".", Args: []string{"-c", "echo unsafe"}, Enabled: true}}); err == nil {
		t.Fatal("expected shell rejection")
	}
}
