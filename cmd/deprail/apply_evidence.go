package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/geoffrey-xiao/deprail/internal/process"
	"github.com/geoffrey-xiao/deprail/internal/remediation"
	"github.com/geoffrey-xiao/deprail/internal/remediation/evidence"
)

type applyEvidenceState struct {
	BeforeDigest        string
	AfterDigest         string
	FindingBeforeDigest string
	FindingAfterDigest  string
	VerificationStatus  string
	RescanStatus        string
	Commands            []evidence.CommandRecord
}

func newApplyEvidenceState(plan remediation.Plan) (applyEvidenceState, error) {
	findingBeforeDigest, err := digestApplyValue([]remediation.ReportFinding{planFinding(plan)})
	if err != nil {
		return applyEvidenceState{}, fmt.Errorf("digest finding before state: %w", err)
	}
	return applyEvidenceState{
		FindingBeforeDigest: findingBeforeDigest,
		VerificationStatus:  "not_run",
		RescanStatus:        "not_run",
		Commands:            []evidence.CommandRecord{},
	}, nil
}

func applyEvidenceCommand(id, tool string, result process.Result) evidence.CommandRecord {
	stdout := string(result.Stdout)
	stderr := string(result.Stderr)
	return evidence.CommandRecord{
		ID:           id,
		Tool:         filepath.Base(tool),
		ExitCode:     result.ExitCode,
		Stdout:       stdout,
		Stderr:       stderr,
		StdoutDigest: digestApplyText(stdout),
		StderrDigest: digestApplyText(stderr),
	}
}

func digestApplyText(value string) string {
	redacted := evidence.Redact(value)
	sum := sha256.Sum256([]byte(redacted))
	return hex.EncodeToString(sum[:])
}

func digestApplyValue(value any) (string, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}

func applyEvidenceOperations(plan remediation.Plan) []evidence.Operation {
	operations := make([]evidence.Operation, 0, len(plan.AffectedFiles)+len(plan.Commands))
	for _, file := range plan.AffectedFiles {
		effect := file.Effect
		if effect == "" {
			effect = "update"
		}
		operations = append(operations, evidence.Operation{ID: "file:" + file.Path, Path: file.Path, Effect: effect})
	}
	for index, command := range plan.Commands {
		workingDirectory := command.WorkingDirectory
		if workingDirectory == "" {
			workingDirectory = "."
		}
		operations = append(operations, evidence.Operation{ID: fmt.Sprintf("command:%d", index), Path: workingDirectory, Effect: "execute"})
	}
	return operations
}

func buildApplyEvidence(plan remediation.Plan, result applyResult, state applyEvidenceState, workspacePath, cleanup string) (evidence.Record, error) {
	record := evidence.Record{
		SchemaVersion:       evidence.SchemaVersion,
		Outcome:             evidence.Outcome(result.Outcome),
		PlanDigest:          result.PlanDigest,
		SourceCommit:        result.SourceCommit,
		SourceRoot:          result.SourceRoot,
		WorkspaceID:         plan.WorkspaceIdentity.ID,
		WorkspacePath:       workspacePath,
		AuthorizedPaths:     make([]string, 0, len(plan.AffectedFiles)),
		Operations:          applyEvidenceOperations(plan),
		BeforeDigest:        state.BeforeDigest,
		AfterDigest:         state.AfterDigest,
		FindingBeforeDigest: state.FindingBeforeDigest,
		FindingAfterDigest:  state.FindingAfterDigest,
		ArtifactDigests:     append([]string{}, plan.Provenance.ArtifactDigests...),
		VerificationStatus:  state.VerificationStatus,
		RescanStatus:        state.RescanStatus,
		Commands:            append([]evidence.CommandRecord{}, state.Commands...),
		Diagnostics:         append([]string{}, result.Diagnostics...),
		Cleanup:             cleanup,
	}
	for _, file := range plan.AffectedFiles {
		record.AuthorizedPaths = append(record.AuthorizedPaths, file.Path)
	}
	record.ArtifactDigests = append(record.ArtifactDigests, result.RescanArtifactDigests...)
	if err := record.Validate(); err != nil {
		return evidence.Record{}, err
	}
	return record, nil
}

func persistApplyEvidence(root string, plan remediation.Plan, result applyResult, state applyEvidenceState, workspacePath, cleanup string) (string, string, error) {
	record, err := buildApplyEvidence(plan, result, state, workspacePath, cleanup)
	if err != nil {
		return "", "", err
	}
	return (evidence.Store{Root: root}).Save(record)
}
