package verification

import (
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

func TestClassifyIncompleteScanWithholdsSafetyClaims(t *testing.T) {
	result, err := Classify([]remediation.ReportFinding{finding("missing"), finding("residual")}, []remediation.ReportFinding{finding("residual")}, false)
	if err != nil {
		t.Fatal(err)
	}
	if result[0] != (FindingTransition{"missing", Unknown}) || result[1] != (FindingTransition{"residual", Residual}) {
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

func TestSelectCommandsFiltersAndSortsSupportedChecks(t *testing.T) {
	commands, err := SelectCommands([]Command{
		{ID: "build", Kind: Build, Path: "/usr/bin/build", Enabled: true},
		{ID: "test", Kind: Test, Path: "/usr/bin/test", Enabled: true},
		{ID: "disabled", Kind: Test, Path: "/usr/bin/test", Enabled: false},
		{ID: "shell", Kind: CommandKind("shell"), Path: "/bin/sh", Enabled: true},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(commands) != 2 || commands[0].ID != "build" || commands[1].ID != "test" {
		t.Fatalf("commands = %#v", commands)
	}
}
