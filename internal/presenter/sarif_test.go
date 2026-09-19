package presenter

import (
	"encoding/json"
	"testing"

	"github.com/geoffrey-xiao/deprail/internal/baseline"
	"github.com/geoffrey-xiao/deprail/internal/discovery"
)

func TestSARIFFromBaselineIsVersionedAndDeterministic(t *testing.T) {
	d := baseline.Document{SchemaVersion: baseline.SchemaVersion, DocumentType: baseline.DocumentType, BaselineID: "base", SourceScanID: "scan", Status: discovery.Complete, Findings: []baseline.Finding{{StableKey: "key", Component: "pkg:npm/a@1", Version: "1", Severity: "high"}}, ArtifactDigests: []string{}}
	first, err := SARIFFromBaseline(d)
	if err != nil {
		t.Fatal(err)
	}
	second, err := SARIFFromBaseline(d)
	if err != nil {
		t.Fatal(err)
	}
	if string(first) != string(second) {
		t.Fatal("SARIF output is nondeterministic")
	}
	var decoded SARIF
	if err := json.Unmarshal(first, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Version != "2.1.0" || len(decoded.Runs) != 1 || len(decoded.Runs[0].Results) != 1 {
		t.Fatalf("SARIF = %#v", decoded)
	}
}
