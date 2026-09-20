package presenter

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/geoffrey-xiao/deprail/internal/remediation"
)

func WritePlanJSON(w io.Writer, plan remediation.Plan) error {
	return json.NewEncoder(w).Encode(plan)
}

func WritePlanTerminal(w io.Writer, plan remediation.Plan) error {
	state := "none"
	if plan.Recommendation != nil {
		state = string(plan.Recommendation.State)
	}
	if _, err := fmt.Fprintf(w, "Plan: %s\nFinding: %s\nRecommendation: %s\nCandidates: %d\nCommands: %d\n", plan.PlanID, plan.FindingIdentity.StableKey, state, len(plan.Candidates), len(plan.Commands)); err != nil {
		return err
	}
	for _, candidate := range plan.Candidates {
		if _, err := fmt.Fprintf(w, "Candidate: %s %s\n", candidate.State, candidate.Version); err != nil {
			return err
		}
	}
	return nil
}
