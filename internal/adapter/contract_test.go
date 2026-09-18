package adapter

import (
	"context"
	"testing"
)

type contractAdapter struct{}

func (contractAdapter) Metadata(context.Context) (Metadata, error) {
	return Metadata{Name: "fixture", Version: "1", SupportedTargets: []string{"npm"}}, nil
}
func (contractAdapter) Compatible(context.Context, Metadata) error { return nil }
func (contractAdapter) Plan(_ context.Context, targets []Target) (Plan, error) {
	plan := Plan{Targets: append([]Target(nil), targets...)}
	return plan, ValidatePlan(plan)
}
func (contractAdapter) Execute(context.Context, Plan) (RawResult, error) {
	return RawResult{ExitCode: 0}, nil
}
func (contractAdapter) Parse(context.Context, RawResult) ([]Record, error)     { return nil, nil }
func (contractAdapter) Normalize(context.Context, []Record) ([]Finding, error) { return nil, nil }

func TestScannerAdapterContractCanBeImplementedWithoutToolTypes(t *testing.T) {
	var scanner ScannerAdapter = contractAdapter{}
	metadata, err := scanner.Metadata(context.Background())
	if err != nil || metadata.Name != "fixture" {
		t.Fatalf("metadata = %#v, err = %v", metadata, err)
	}
	if _, err := scanner.Plan(context.Background(), []Target{{WorkspaceID: "root", RelativePath: ".", Ecosystem: "npm"}}); err != nil {
		t.Fatalf("plan: %v", err)
	}
}

func TestValidatePlanRejectsMissingIdentityAndDuplicateTargets(t *testing.T) {
	for name, plan := range map[string]Plan{
		"missing identity": {Targets: []Target{{WorkspaceID: "", RelativePath: "app", Ecosystem: "npm"}}},
		"duplicate": {Targets: []Target{
			{WorkspaceID: "app", RelativePath: "app", Ecosystem: "npm"},
			{WorkspaceID: "app", RelativePath: "app", Ecosystem: "npm"},
		}},
	} {
		t.Run(name, func(t *testing.T) {
			if err := ValidatePlan(plan); !IsCode(err, ErrInvalidPlan) {
				t.Fatalf("error = %v, want %s", err, ErrInvalidPlan)
			}
		})
	}
}

func TestValidatePlanRejectsEscapingPaths(t *testing.T) {
	for _, path := range []string{"/tmp/outside", "../outside", "app/../../outside"} {
		t.Run(path, func(t *testing.T) {
			err := ValidatePlan(Plan{Targets: []Target{{WorkspaceID: "app", RelativePath: path, Ecosystem: "npm"}}})
			if !IsCode(err, ErrInvalidPlan) {
				t.Fatalf("error = %v, want invalid plan", err)
			}
		})
	}
}
