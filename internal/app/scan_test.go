package app

import (
	"context"
	"encoding/json"
	"github.com/geoffrey-xiao/deprail/internal/adapter"
	"github.com/geoffrey-xiao/deprail/internal/artifact"
	"path/filepath"
	"strings"
	"testing"
)

type emptyScanner struct {
	fakeScanner
}

func (emptyScanner) Parse(context.Context, adapter.RawResult) ([]adapter.Record, error) {
	return nil, nil
}

func (emptyScanner) Normalize(context.Context, []adapter.Record) ([]adapter.Finding, error) {
	return nil, nil
}

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

func TestScanEmptyCollectionsSerializeAsArrays(t *testing.T) {
	report, err := Scan(context.Background(), fixturePath(t, "npm-basic"), ScanOptions{
		Scanner:   emptyScanner{},
		Artifacts: artifact.Store{Root: t.TempDir(), MaxBytes: 1024},
	})
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	data, err := json.Marshal(report)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	for _, field := range []string{`"findings":[]`, `"errors":[]`} {
		if !strings.Contains(string(data), field) {
			t.Fatalf("serialized report %s missing %s", data, field)
		}
	}
	if !strings.Contains(string(data), `"artifact_digests":[`) {
		t.Fatalf("serialized report %s missing artifact array", data)
	}
}

func TestScanFailureSerializesEmptyCollections(t *testing.T) {
	report, err := Scan(context.Background(), filepath.Join(t.TempDir(), "missing"), ScanOptions{})
	if err == nil {
		t.Fatal("expected scan failure")
	}
	data, marshalErr := json.Marshal(report)
	if marshalErr != nil {
		t.Fatalf("marshal: %v", marshalErr)
	}
	for _, field := range []string{`"findings":[]`, `"artifact_digests":[]`} {
		if !strings.Contains(string(data), field) {
			t.Fatalf("serialized report %s missing %s", data, field)
		}
	}
	if !strings.Contains(string(data), `"errors":[`) {
		t.Fatalf("serialized report %s missing errors array", data)
	}
}
