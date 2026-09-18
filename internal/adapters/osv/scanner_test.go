package osv

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/geoffrey-xiao/deprail/internal/adapter"
)

func TestParseSortsRecordsAndPreservesEvidence(t *testing.T) {
	raw := adapter.RawResult{Stdout: []byte(`{"results":[{"packages":[{"package":{"name":"z","version":"1"},"vulnerabilities":[{"id":"OSV-2","aliases":["CVE-2","GHSA-2"],"severity":"HIGH","database_specific":{"fixed_version":"2"}}]},{"package":{"name":"a","version":"1"},"vulnerabilities":[{"id":"OSV-1"}]}]}]}`)}
	records, err := Parse(raw)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(records) != 2 || records[0].Component != "a" || records[1].Fixed != "2" || records[1].Aliases[0] != "CVE-2" {
		t.Fatalf("records = %#v", records)
	}
}

func TestParseRejectsMalformedJSON(t *testing.T) {
	if _, err := Parse(adapter.RawResult{Stdout: []byte("not json")}); !adapter.IsCode(err, adapter.ErrInvalidOutput) {
		t.Fatalf("error = %v", err)
	}
}

func TestScannerUsesV2SourceCommand(t *testing.T) {
	scanner := Scanner{Path: os.Args[0], Args: []string{"-test.run=TestScannerHelper", "mode=args"}, Timeout: time.Second, OutputCap: 1024}
	raw, err := scanner.Execute(context.Background(), adapter.Plan{Targets: []adapter.Target{{WorkspaceID: "root", RelativePath: ".", Ecosystem: "npm"}}})
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if !strings.HasPrefix(string(raw.Stdout), `{"results":[]}`) {
		t.Fatalf("stdout = %q", raw.Stdout)
	}
}

func TestScannerRejectsMultipleExecutionTargets(t *testing.T) {
	scanner := Scanner{Path: "missing", Timeout: time.Second, OutputCap: 1024}
	_, err := scanner.Execute(context.Background(), adapter.Plan{Targets: []adapter.Target{{WorkspaceID: "a", RelativePath: "a", Ecosystem: "npm"}, {WorkspaceID: "b", RelativePath: "b", Ecosystem: "npm"}}})
	if !adapter.IsCode(err, adapter.ErrInvalidPlan) {
		t.Fatalf("error = %v", err)
	}
}
