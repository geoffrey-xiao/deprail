package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/geoffrey-xiao/deprail/internal/remediation"
	"github.com/geoffrey-xiao/deprail/internal/remediation/isolation"
)

type applyResult struct {
	SchemaVersion string   `json:"schema_version"`
	Outcome       string   `json:"outcome"`
	PlanID        string   `json:"plan_id"`
	PlanDigest    string   `json:"plan_digest"`
	SourceRoot    string   `json:"source_root"`
	SourceCommit  string   `json:"source_commit"`
	Diagnostics   []string `json:"diagnostics"`
}

func runFixApply(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("fix apply", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	planPath := flags.String("plan", "", "versioned remediation plan")
	approvalPath := flags.String("approval", "", "versioned approval record")
	root := flags.String("root", ".", "canonical repository root")
	format := flags.String("format", "terminal", "output format: terminal or json")
	dryRun := flags.Bool("dry-run", false, "validate without mutation")
	if err := flags.Parse(args); err != nil || *planPath == "" || *approvalPath == "" || (*format != "terminal" && *format != "json") || len(flags.Args()) != 0 {
		writeCLIError(stderr, "CONFIG_INVALID", "fix apply requires --plan and --approval and supports --root, --dry-run, and terminal or json output", "fix apply")
		return 2
	}
	plan, err := loadApplyJSON[remediation.Plan](*planPath)
	if err != nil {
		writeCLIError(stderr, "PLAN_INVALID", err.Error(), "plan")
		return 3
	}
	if err := plan.Validate(); err != nil {
		writeCLIError(stderr, "PLAN_INVALID", err.Error(), "plan")
		return 3
	}
	approval, err := loadApplyJSON[remediation.Approval](*approvalPath)
	if err != nil {
		writeCLIError(stderr, "APPROVAL_INVALID", err.Error(), "approval")
		return 3
	}
	canonical, err := isolation.CanonicalRepositoryRoot(context.Background(), *root)
	if err != nil {
		writeCLIError(stderr, "PATH_OUTSIDE_ROOT", err.Error(), "root")
		return 3
	}
	commit, err := currentCommit(canonical)
	if err != nil {
		writeCLIError(stderr, "SOURCE_INVALID", err.Error(), "root")
		return 3
	}
	if err := approval.ValidatePersisted(plan, canonical, commit, time.Now().UTC()); err != nil {
		var approvalErr *remediation.ApprovalError
		if errors.As(err, &approvalErr) {
			writeCLIError(stderr, string(approvalErr.Code), approvalErr.Message, "approval")
		} else {
			writeCLIError(stderr, "APPROVAL_INVALID", err.Error(), "approval")
		}
		return 3
	}
	result := applyResult{SchemaVersion: "v0alpha1", Outcome: "dry_run", PlanID: plan.PlanID, PlanDigest: remediation.PlanDigest(plan), SourceRoot: canonical, SourceCommit: commit, Diagnostics: []string{}}
	if !*dryRun {
		result.Outcome = "not_ready"
		result.Diagnostics = append(result.Diagnostics, "mutation, verification, rescan, cleanup, and evidence orchestration are not yet wired")
	}
	if *format == "json" {
		if err := json.NewEncoder(stdout).Encode(result); err != nil {
			writeCLIError(stderr, "OUTPUT_WRITE_FAILED", err.Error(), "stdout")
			return 3
		}
	} else {
		_, _ = fmt.Fprintf(stdout, "Outcome: %s\nPlan: %s\nSource: %s\n", result.Outcome, result.PlanID, result.SourceCommit)
		for _, diagnostic := range result.Diagnostics {
			_, _ = fmt.Fprintf(stdout, "Diagnostic: %s\n", diagnostic)
		}
	}
	if !*dryRun {
		return 3
	}
	return 0
}

func currentCommit(root string) (string, error) {
	git, err := exec.LookPath("git")
	if err != nil {
		return "", fmt.Errorf("git is unavailable: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	output, err := exec.CommandContext(ctx, git, "-C", root, "rev-parse", "HEAD").Output()
	if err != nil {
		return "", fmt.Errorf("resolve source commit: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

func loadApplyJSON[T any](path string) (T, error) {
	var value T
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return value, err
	}
	defer file.Close()
	decoder := json.NewDecoder(io.LimitReader(file, 16<<20))
	if err := decoder.Decode(&value); err != nil {
		return value, err
	}
	return value, nil
}
