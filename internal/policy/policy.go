package policy

import (
	"errors"
	"fmt"
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
type Input struct {
	Comparison baseline.Comparison
	Complete   bool
	Errors     []string
}
type Decision struct {
	Outcome Outcome  `json:"outcome"`
	Reasons []string `json:"reasons"`
}

func severityRank(value string) int {
	switch strings.ToLower(value) {
	case "low":
		return 1
	case "medium":
		return 2
	case "high":
		return 3
	case "critical":
		return 4
	default:
		return 0
	}
}

func Evaluate(policy Policy, input Input) (Decision, error) {
	minimum := strings.ToLower(policy.MinimumSeverity)
	if minimum != "" && severityRank(minimum) == 0 {
		return Decision{}, errors.New("invalid minimum severity")
	}
	if !input.Complete || len(input.Errors) > 0 {
		return Decision{Outcome: Block, Reasons: []string{"comparison is incomplete or contains scanner errors"}}, nil
	}
	decision := Decision{Outcome: Pass, Reasons: []string{}}
	for _, change := range input.Comparison.Changes {
		if change.Kind != baseline.Added && change.Kind != baseline.Changed {
			continue
		}
		if change.Head == nil {
			decision.Outcome = Block
			decision.Reasons = append(decision.Reasons, "unknown finding: "+change.Key)
			continue
		}
		if policy.BlockOnNew && change.Kind == baseline.Added {
			decision.Outcome = Block
			decision.Reasons = append(decision.Reasons, "new finding: "+change.Key)
			continue
		}
		if minimum != "" {
			level := severityRank(change.Head.Severity)
			if level == 0 {
				decision.Outcome = Block
				decision.Reasons = append(decision.Reasons, fmt.Sprintf("unknown severity: %s", change.Key))
				continue
			}
			if level >= severityRank(minimum) {
				decision.Outcome = Block
				decision.Reasons = append(decision.Reasons, fmt.Sprintf("severity %s meets threshold: %s", change.Head.Severity, change.Key))
			}
		}
	}
	return decision, nil
}
