package app

import (
	"context"
	"testing"

	"github.com/geoffrey-xiao/deprail/internal/adapter"
	"github.com/geoffrey-xiao/deprail/internal/artifact"
)

type fakeScanner struct{}

func (fakeScanner) Metadata(context.Context) (adapter.Metadata, error) {
	return adapter.Metadata{Name: "fake", Version: "1"}, nil
}
func (fakeScanner) Compatible(context.Context, adapter.Metadata) error { return nil }
func (fakeScanner) Plan(_ context.Context, targets []adapter.Target) (adapter.Plan, error) {
	return adapter.Plan{Targets: targets}, nil
}
func (fakeScanner) Execute(context.Context, adapter.Plan) (adapter.RawResult, error) {
	return adapter.RawResult{Stdout: []byte(`{"results":[]`)}, nil
}
func (fakeScanner) Parse(context.Context, adapter.RawResult) ([]adapter.Record, error) {
	return []adapter.Record{{Component: "demo", Version: "1", TargetID: "OSV-1"}}, nil
}
func (fakeScanner) Normalize(context.Context, []adapter.Record) ([]adapter.Finding, error) {
	return []adapter.Finding{{Component: "demo", Version: "1", TargetID: "OSV-1"}}, nil
}

func TestScanRetainsArtifactAndFindings(t *testing.T) {
	report, err := Scan(context.Background(), fixturePath(t, "npm-basic"), ScanOptions{Scanner: fakeScanner{}, Artifacts: artifact.Store{Root: t.TempDir(), MaxBytes: 1024}})
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	if len(report.Findings) != 1 || len(report.ArtifactDigests) != 1 {
		t.Fatalf("report = %#v", report)
	}
}
