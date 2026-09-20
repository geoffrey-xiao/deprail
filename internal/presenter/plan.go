package presenter

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/geoffrey-xiao/deprail/internal/remediation"
)

var (
	ErrPlanOutputOutsideRoot = errors.New("plan output is outside the allowed boundary")
	ErrPlanSchemaInvalid     = errors.New("generated plan schema is invalid")
	ErrPlanWriteFailed       = errors.New("plan persistence failed")
)

func WritePlanJSON(w io.Writer, plan remediation.Plan) error {
	canonical := plan
	canonical.Canonicalize()
	if err := canonical.Validate(); err != nil {
		return fmt.Errorf("%w: %v", ErrPlanSchemaInvalid, err)
	}
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(canonical)
}

func WritePlanTerminal(w io.Writer, plan remediation.Plan) error {
	state := "no recommendation"
	if plan.Recommendation != nil {
		state = string(plan.Recommendation.State) + " " + plan.Recommendation.CandidateID
	}
	if err := WriteHeader(w, "Plan: "+plan.PlanID); err != nil {
		return err
	}
	if err := WriteSummary(w, "Finding", plan.FindingIdentity.StableKey); err != nil {
		return err
	}
	if err := WriteSummary(w, "Recommendation", state); err != nil {
		return err
	}
	if err := WriteSummary(w, "Read-only", "yes"); err != nil {
		return err
	}
	if err := WriteSection(w, "Risks"); err != nil {
		return err
	}
	for _, risk := range plan.Risks {
		if _, err := fmt.Fprintf(w, "- %s [%s] detected=%t %s\n", risk.Code, risk.Severity, risk.Detected, risk.Details); err != nil {
			return err
		}
	}
	if err := WriteSection(w, "Affected files"); err != nil {
		return err
	}
	for _, file := range plan.AffectedFiles {
		if _, err := fmt.Fprintf(w, "- %s (%s) %s\n", file.Path, file.Kind, file.Effect); err != nil {
			return err
		}
	}
	if err := WriteSection(w, "Future commands"); err != nil {
		return err
	}
	for _, command := range plan.Commands {
		if _, err := fmt.Fprintf(w, "- (cd %s && %s %s)\n", command.WorkingDirectory, command.Executable, strings.Join(command.Arguments, " ")); err != nil {
			return err
		}
	}
	if err := WriteSection(w, "Verification"); err != nil {
		return err
	}
	for _, verification := range plan.Verification {
		if _, err := fmt.Fprintf(w, "- %s: %s %s\n", verification.ID, verification.Command.Executable, strings.Join(verification.Command.Arguments, " ")); err != nil {
			return err
		}
	}
	if err := WriteSection(w, "Provenance"); err != nil {
		return err
	}
	if err := WriteSummary(w, "report", plan.CreatedFrom.ReportDigest); err != nil {
		return err
	}
	if err := WriteSummary(w, "scan", plan.CreatedFrom.SourceScanID); err != nil {
		return err
	}
	return WriteSummary(w, "repository state", plan.CreatedFrom.RepositoryState)
}

func WritePlanAtomic(path string, plan remediation.Plan) error {
	var data bytes.Buffer
	if err := WritePlanJSON(&data, plan); err != nil {
		return err
	}
	if _, err := remediation.WriteExternalOutput(plan.RepositoryIdentity.Root, path, data.Bytes()); err != nil {
		if errors.Is(err, remediation.ErrUnsafePath) {
			return fmt.Errorf("%w: %v", ErrPlanOutputOutsideRoot, err)
		}
		return fmt.Errorf("%w: %v", ErrPlanWriteFailed, err)
	}
	return nil
}
