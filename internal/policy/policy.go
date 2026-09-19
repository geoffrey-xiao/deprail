package policy

import (
	"errors"
	"strings"

	"github.com/geoffrey-xiao/deprail/internal/baseline"
)

type Outcome string

const (
	Pass  Outcome = "pass"
	Warn  Outcome = "warn"
	Block Outcome = "block"
)

type Policy struct {
	BlockOnNew      bool   `json:"block_on_new"`
	MinimumSeverity string `json:"minimum_severity"`
}

type Decision struct {
	Outcome Outcome  `json:"outcome"`
	Reasons []string `json:"reasons"`
}

func Evaluate(policy Policy, comparison baseline.Comparison) (Decision, error) {
	minimum := strings.ToLower(policy.MinimumSeverity)
	if minimum != "" && minimum != "low" && minimum != "medium" && minimum != "high" && minimum != "critical" {
		return Decision{}, errors.New("invalid minimum severity")
	}
	decision := Decision{Outcome: Pass, Reasons: []string{}}
	for _, change := range comparison.Changes {
		if change.Kind != baseline.Added && change.Kind != baseline.Changed {
			continue
		}
		finding := change.Head
		if policy.BlockOnNew && change.Kind == baseline.Added {
			decision.Outcome = Block
			decision.Reasons = append(decision.Reasons, "new finding: "+change.Key)
			continue
		}
		if finding == nil {
			decision.Outcome = Block
			decision.Reasons = append(decision.Reasons, "unknown finding: "+change.Key)
			continue
		}
	}
	return decision, nil
}
