package presenter

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/geoffrey-xiao/deprail/internal/baseline"
)

type SARIF struct {
	Version string `json:"version"`
	Schema  string `json:"$schema"`
	Runs    []Run  `json:"runs"`
}
type Run struct {
	Tool       Tool              `json:"tool"`
	Results    []Result          `json:"results"`
	Properties map[string]string `json:"properties,omitempty"`
}
type Tool struct {
	Driver Driver `json:"driver"`
}
type Driver struct {
	Name  string `json:"name"`
	Rules []Rule `json:"rules"`
}
type Rule struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type Result struct {
	RuleID  string  `json:"ruleId"`
	Level   string  `json:"level"`
	Message Message `json:"message"`
}
type Message struct {
	Text string `json:"text"`
}

func sarifLevel(severity string) string {
	switch severity {
	case "low":
		return "note"
	case "medium":
		return "warning"
	case "high", "critical":
		return "error"
	default:
		return "warning"
	}
}

func SARIFFromBaseline(document baseline.Document) ([]byte, error) {
	if err := baseline.Validate(document); err != nil {
		return nil, errors.New("cannot render invalid baseline: " + err.Error())
	}
	document.Canonicalize()
	rules := make([]Rule, 0, len(document.Findings))
	results := make([]Result, 0, len(document.Findings))
	for _, finding := range document.Findings {
		rules = append(rules, Rule{ID: finding.StableKey, Name: finding.Component})
		results = append(results, Result{RuleID: finding.StableKey, Level: sarifLevel(finding.Severity), Message: Message{Text: finding.Component + " " + finding.Version}})
	}
	properties := map[string]string{"baseline_id": document.BaselineID, "source_scan_id": document.SourceScanID}
	for i, digest := range document.ArtifactDigests {
		properties[fmt.Sprintf("artifact_digest_%d", i)] = digest
	}
	run := Run{Tool: Tool{Driver: Driver{Name: "deprail", Rules: rules}}, Results: results, Properties: properties}
	return json.MarshalIndent(SARIF{Version: "2.1.0", Schema: "https://json.schemastore.org/sarif-2.1.0.json", Runs: []Run{run}}, "", "  ")
}
