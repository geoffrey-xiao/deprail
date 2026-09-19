package presenter

import (
	"encoding/json"

	"github.com/geoffrey-xiao/deprail/internal/baseline"
)

type SARIF struct {
	Version string `json:"version"`
	Schema  string `json:"$schema"`
	Runs    []Run  `json:"runs"`
}
type Run struct {
	Tool    Tool     `json:"tool"`
	Results []Result `json:"results"`
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

func SARIFFromBaseline(document baseline.Document) ([]byte, error) {
	rules := make([]Rule, 0, len(document.Findings))
	results := make([]Result, 0, len(document.Findings))
	for _, finding := range document.Findings {
		rules = append(rules, Rule{ID: finding.StableKey, Name: finding.Component})
		level := finding.Severity
		if level == "" || level == "unknown" {
			level = "warning"
		}
		results = append(results, Result{RuleID: finding.StableKey, Level: level, Message: Message{Text: finding.Component + " " + finding.Version}})
	}
	return json.MarshalIndent(SARIF{Version: "2.1.0", Schema: "https://json.schemastore.org/sarif-2.1.0.json", Runs: []Run{{Tool: Tool{Driver: Driver{Name: "deprail", Rules: rules}}, Results: results}}}, "", "  ")
}
