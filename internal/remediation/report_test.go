package remediation

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

const reportDigest = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

func testReport() Report {
	return Report{
		SchemaVersion:      ReportSchemaVersion,
		DocumentType:       ReportDocumentType,
		ReportID:           "scan-1",
		RepositoryIdentity: RepositoryIdentity{Root: ".", Repository: "repo", Revision: "rev-1"},
		SourceScanID:       "scan-1",
		ArtifactDigests:    []string{reportDigest},
		RepositoryState:    "tree-1",
		Status:             ReportComplete,
		Findings: []ReportFinding{{
			StableKey:      "finding-1",
			Workspace:      WorkspaceIdentity{ID: "service", Path: "."},
			Component:      Component{PURL: "pkg:npm/example@1.0.0", Name: "example", Version: "1.0.0"},
			Aliases:        []string{"OSV-1"},
			CurrentVersion: "1.0.0",
			DependencyPath: []string{"service", "example"},
			Provenance:     Provenance{ArtifactDigests: []string{reportDigest}, Sources: []string{"scan.json"}},
		}},
	}
}

func writeReport(t *testing.T, report Report) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "report.json")
	data, err := marshalJSON(report)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func marshalJSON(value any) ([]byte, error) {
	return json.Marshal(value)
}
func TestResolveFindingRequiresExactUniqueKey(t *testing.T) {
	report := testReport()
	resolved, err := ResolveFinding(report, "finding-1", "tree-1")
	if err != nil {
		t.Fatal(err)
	}
	if resolved.Report.SourceScanID != report.SourceScanID || resolved.Finding.StableKey != "finding-1" {
		t.Fatalf("resolved = %#v", resolved)
	}

	if _, err := ResolveFinding(report, "missing", "tree-1"); !hasReportCode(err, FindingNotFound) {
		t.Fatalf("missing finding error = %v", err)
	}

	report.Findings = append(report.Findings, report.Findings[0])
	if _, err := ResolveFinding(report, "finding-1", "tree-1"); !hasReportCode(err, FindingAmbiguous) {
		t.Fatalf("duplicate finding error = %v", err)
	}
}

func TestLoadAndResolveRejectsImplicitOrInvalidReports(t *testing.T) {
	if _, err := LoadReport(""); !hasReportCode(err, ReportRequired) {
		t.Fatalf("missing report error = %v", err)
	}

	path := filepath.Join(t.TempDir(), "malformed.json")
	if err := os.WriteFile(path, []byte("{} trailing"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadReport(path); !hasReportCode(err, ReportInvalid) {
		t.Fatalf("malformed report error = %v", err)
	}

	report := testReport()
	if _, err := ResolveFinding(report, "finding-1", "changed-tree"); !hasReportCode(err, ReportInputStale) {
		t.Fatalf("stale report error = %v", err)
	}
}

func TestLoadAndResolvePreservesProvenance(t *testing.T) {
	path := writeReport(t, testReport())
	resolved, err := LoadAndResolve(path, "finding-1", "tree-1")
	if err != nil {
		t.Fatal(err)
	}
	if resolved.Report.RepositoryIdentity.Repository != "repo" || resolved.Report.SourceScanID != "scan-1" || len(resolved.Report.ArtifactDigests) != 1 {
		t.Fatalf("provenance = %#v", resolved.Report)
	}
}
func TestLoadReportRejectsOversizedInput(t *testing.T) {
	path := filepath.Join(t.TempDir(), "oversized.json")
	data, err := marshalJSON(testReport())
	if err != nil {
		t.Fatal(err)
	}
	data = append(data, bytes.Repeat([]byte(" "), maxReportBytes-len(data)+1)...)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadReport(path); !hasReportCode(err, ReportInvalid) {
		t.Fatalf("oversized report error = %v", err)
	}
}

func TestValidateReportRejectsMalformedProvenanceAndPath(t *testing.T) {
	report := testReport()
	report.ArtifactDigests = []string{"not-a-digest"}
	if err := ValidateReport(report); !hasReportCode(err, ReportInvalid) {
		t.Fatalf("digest error = %v", err)
	}

	report = testReport()
	report.Findings[0].Workspace.Path = "../outside"
	if err := ValidateReport(report); !hasReportCode(err, ReportInvalid) {
		t.Fatalf("path error = %v", err)
	}
}

func hasReportCode(err error, want ReportErrorCode) bool {
	var reportErr *ReportError
	return errors.As(err, &reportErr) && reportErr.Code == want
}
