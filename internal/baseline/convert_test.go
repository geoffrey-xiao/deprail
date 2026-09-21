package baseline

import (
	"reflect"
	"testing"

	"github.com/geoffrey-xiao/deprail/internal/discovery"
)

func TestConvertScanIsDeterministicAndPreservesIdentity(t *testing.T) {
	first := ScanInput{
		SchemaVersion:   ScanSchemaVersion,
		DocumentType:    ScanDocumentType,
		ScanID:          "scan-001",
		RepositoryState: "tree-001",
		Status:          discovery.Complete,
		Findings: []ScanFinding{
			{PURL: "pkg:npm/b@2.0.0", Version: "2.0.0", WorkspaceID: "api", TargetID: "OSV-B", Aliases: []string{"CVE-B", "OSV-B"}},
			{PURL: "pkg:npm/a@1.0.0", Version: "1.0.0", WorkspaceID: "root", TargetID: "OSV-A", Aliases: []string{"OSV-A", "CVE-A"}},
		},
		Errors:          []string{},
		ArtifactDigests: []string{"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
	}
	second := first
	second.Findings = []ScanFinding{
		{PURL: "pkg:npm/a@1.0.0", Version: "1.0.0", WorkspaceID: "root", TargetID: "OSV-A", Aliases: []string{"CVE-A", "OSV-A"}},
		{PURL: "pkg:npm/b@2.0.0", Version: "2.0.0", WorkspaceID: "api", TargetID: "OSV-B", Aliases: []string{"OSV-B", "CVE-B"}},
	}

	left, err := ConvertScan(first)
	if err != nil {
		t.Fatal(err)
	}
	right, err := ConvertScan(second)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(left, right) {
		t.Fatalf("equivalent scans produced different baselines:\nleft=%#v\nright=%#v", left, right)
	}
	if left.BaselineID != "base-scan-001" || left.SourceScanID != "scan-001" {
		t.Fatalf("identity = %#v", left)
	}
	if err := Validate(left); err != nil {
		t.Fatal(err)
	}
}

func TestConvertScanRejectsIncompleteInvalidAndStaleInputs(t *testing.T) {
	base := ScanInput{SchemaVersion: ScanSchemaVersion, DocumentType: ScanDocumentType, ScanID: "scan-001", RepositoryState: "tree-001", Status: discovery.Complete, Findings: []ScanFinding{}, Errors: []string{}, ArtifactDigests: []string{}}
	cases := []struct {
		name string
		edit func(*ScanInput)
		code ConversionErrorCode
	}{
		{name: "partial", edit: func(input *ScanInput) { input.Status = discovery.Partial }, code: ErrScanIncomplete},
		{name: "failed", edit: func(input *ScanInput) { input.Status = discovery.Failed }, code: ErrScanIncomplete},
		{name: "errors", edit: func(input *ScanInput) { input.Errors = []string{"scanner failed"} }, code: ErrScanIncomplete},
		{name: "stale", edit: func(input *ScanInput) { input.RepositoryState = "" }, code: ErrScanStale},
		{name: "schema", edit: func(input *ScanInput) { input.DocumentType = "baseline" }, code: ErrScanInvalid},
		{name: "duplicate digest", edit: func(input *ScanInput) {
			input.ArtifactDigests = []string{"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}
		}, code: ErrScanInvalid},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			input := base
			test.edit(&input)
			_, err := ConvertScan(input)
			if err == nil {
				t.Fatal("invalid scan was accepted")
			}
			conversionErr, ok := err.(*ConversionError)
			if !ok || conversionErr.Code != test.code {
				t.Fatalf("error = %T %v, want %s", err, err, test.code)
			}
		})
	}
}
