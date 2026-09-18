package osv

import (
	"context"
	"os"
	"path/filepath"
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

func TestParseV2VulnerabilityFields(t *testing.T) {
	raw := adapter.RawResult{Stdout: []byte(`{"results":[{"packages":[{"package":{"name":"lodash","version":"4.17.20"},"vulnerabilities":[{"id":"GHSA-test","aliases":["CVE-test"],"database_specific":{"severity":"HIGH"},"severity":[{"score":"CVSS:3.1/AV:N"}],"affected":[{"ranges":[{"events":[{"introduced":"0"},{"fixed":"4.17.21"}]}]}]}]}]}]}`)}
	records, err := Parse(raw)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(records) != 1 || records[0].Severity != "HIGH" || records[0].Fixed != "4.17.21" {
		t.Fatalf("records = %#v", records)
	}
}

func TestScannerAcceptsVulnerabilityExitCode(t *testing.T) {
	scanner := Scanner{Path: os.Args[0], Args: []string{"-test.run=TestScannerHelper", "mode=vulnerable"}, Timeout: time.Second, OutputCap: 4096}
	raw, err := scanner.Execute(context.Background(), adapter.Plan{Targets: []adapter.Target{{WorkspaceID: "root", RelativePath: ".", Ecosystem: "npm"}}})
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	records, err := Parse(raw)
	if err != nil || len(records) != 1 {
		t.Fatalf("records = %#v, err = %v", records, err)
	}
}

func TestParseRejectsMalformedJSON(t *testing.T) {
	if _, err := Parse(adapter.RawResult{Stdout: []byte("not json")}); !adapter.IsCode(err, adapter.ErrInvalidOutput) {
		t.Fatalf("error = %v", err)
	}
}

func TestParseCheckedInV2Fixtures(t *testing.T) {
	empty, err := os.ReadFile(filepath.Join("..", "..", "..", "testdata", "osv", "v2-empty.json"))
	if err != nil {
		t.Fatal(err)
	}
	emptyRecords, err := Parse(adapter.RawResult{Stdout: empty})
	if err != nil || len(emptyRecords) != 0 {
		t.Fatalf("empty records = %#v, err = %v", emptyRecords, err)
	}

	findings, err := os.ReadFile(filepath.Join("..", "..", "..", "testdata", "osv", "v2-findings.json"))
	if err != nil {
		t.Fatal(err)
	}
	records, err := Parse(adapter.RawResult{Stdout: findings})
	if err != nil {
		t.Fatalf("findings parse: %v", err)
	}
	if len(records) < 4 {
		t.Fatalf("records = %d, want multiple v2 vulnerabilities", len(records))
	}
	if records[0].Component != "lodash" || records[0].Version != "4.17.20" {
		t.Fatalf("first record = %#v", records[0])
	}
	if records[0].Fixed == "" || records[0].Severity == "" || len(records[0].Aliases) == 0 {
		t.Fatalf("first record lost v2 evidence = %#v", records[0])
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
