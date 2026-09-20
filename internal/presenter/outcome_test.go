package presenter

import (
	"strings"
	"testing"
)

func TestWriteOutcomeSummaryDistinguishesStatuses(t *testing.T) {
	cases := []struct {
		name     string
		status   string
		findings int
		want     string
	}{
		{"complete zero", "complete", 0, "No known vulnerabilities found."},
		{"partial", "partial", 0, "Scan incomplete: partial results."},
		{"failed", "failed", 0, "Scan failed."},
		{"cancelled", "cancelled", 0, "Scan cancelled."},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var output strings.Builder
			if err := WriteOutcomeSummary(&output, tc.status, tc.findings, nil); err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(output.String(), tc.want) {
				t.Fatalf("output = %q, want %q", output.String(), tc.want)
			}
		})
	}
}

func TestWriteOutcomeSummaryIncludesScopeAndStableErrorCode(t *testing.T) {
	var output strings.Builder
	if err := WriteOutcomeSummary(&output, "partial", 0, []string{"SCANNER_TIMEOUT: scanner timed out"}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "Affected scope: one or more workspaces") {
		t.Fatalf("output lacks affected scope: %q", output.String())
	}
	if !strings.Contains(output.String(), "Error [SCANNER_TIMEOUT]:") {
		t.Fatalf("output lacks stable error code: %q", output.String())
	}
}

func TestWriteOutcomeSummaryDoesNotClaimNoVulnerabilitiesForFindings(t *testing.T) {
	var output strings.Builder
	if err := WriteOutcomeSummary(&output, "complete", 2, nil); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(output.String(), "No known vulnerabilities found.") || !strings.Contains(output.String(), "Findings: 2") {
		t.Fatalf("output = %q", output.String())
	}
}
