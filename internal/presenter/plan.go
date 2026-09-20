package presenter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/geoffrey-xiao/deprail/internal/remediation"
)

func WritePlanJSON(w io.Writer, plan remediation.Plan) error {
	canonical := plan
	canonical.Canonicalize()
	if err := canonical.Validate(); err != nil {
		return err
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
	if _, err := fmt.Fprintf(w, "Plan: %s\nFinding: %s\nRecommendation: %s\nRead-only: yes\n\n", plan.PlanID, plan.FindingIdentity.StableKey, state); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "Risks:"); err != nil {
		return err
	}
	for _, risk := range plan.Risks {
		if _, err := fmt.Fprintf(w, "- %s [%s] detected=%t %s\n", risk.Code, risk.Severity, risk.Detected, risk.Details); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintln(w, "Affected files:"); err != nil {
		return err
	}
	for _, file := range plan.AffectedFiles {
		if _, err := fmt.Fprintf(w, "- %s (%s) %s\n", file.Path, file.Kind, file.Effect); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintln(w, "Future commands:"); err != nil {
		return err
	}
	for _, command := range plan.Commands {
		if _, err := fmt.Fprintf(w, "- (cd %s && %s %s)\n", command.WorkingDirectory, command.Executable, strings.Join(command.Arguments, " ")); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintln(w, "Verification:"); err != nil {
		return err
	}
	for _, verification := range plan.Verification {
		if _, err := fmt.Fprintf(w, "- %s: %s %s\n", verification.ID, verification.Command.Executable, strings.Join(verification.Command.Arguments, " ")); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintf(w, "Provenance:\n- report: %s\n- scan: %s\n- repository state: %s\n", plan.CreatedFrom.ReportDigest, plan.CreatedFrom.SourceScanID, plan.CreatedFrom.RepositoryState); err != nil {
		return err
	}
	return nil
}

func WritePlanAtomic(path string, plan remediation.Plan) error {
	if path == "" {
		return fmt.Errorf("output path is required")
	}
	if err := rejectInsideRoot(path, plan.RepositoryIdentity.Root); err != nil {
		return err
	}
	var data bytes.Buffer
	if err := WritePlanJSON(&data, plan); err != nil {
		return err
	}
	return WriteAtomic(path, data.Bytes())
}

func rejectInsideRoot(path, root string) error {
	canonicalRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		return fmt.Errorf("resolve repository root: %w", err)
	}
	absolute, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	cleanRoot, err := filepath.Abs(canonicalRoot)
	if err != nil {
		return err
	}
	rel, err := filepath.Rel(cleanRoot, absolute)
	if err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return fmt.Errorf("output path is inside repository root")
	}
	parent, err := filepath.EvalSymlinks(filepath.Dir(absolute))
	if err == nil {
		if rel, e := filepath.Rel(cleanRoot, filepath.Join(parent, filepath.Base(absolute))); e == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
			return fmt.Errorf("output path is inside repository root")
		}
	}
	return nil
}
