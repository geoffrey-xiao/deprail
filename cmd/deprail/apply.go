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
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/geoffrey-xiao/deprail/internal/remediation"
	"github.com/geoffrey-xiao/deprail/internal/remediation/isolation"
	"github.com/geoffrey-xiao/deprail/internal/remediation/mutation"
	"github.com/geoffrey-xiao/deprail/internal/remediation/verification"
)

type applyResult struct {
	SchemaVersion string   `json:"schema_version"`
	Outcome       string   `json:"outcome"`
	PlanID        string   `json:"plan_id"`
	PlanDigest    string   `json:"plan_digest"`
	SourceRoot    string   `json:"source_root"`
	SourceCommit  string   `json:"source_commit"`
	Workspace     string   `json:"workspace,omitempty"`
	Diagnostics   []string `json:"diagnostics"`
}

func runFixApply(args []string, stdout, stderr io.Writer) int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	return runFixApplyContext(ctx, args, stdout, stderr)
}

func runFixApplyContext(ctx context.Context, args []string, stdout, stderr io.Writer) int {
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
	canonical, err := isolation.CanonicalRepositoryRoot(ctx, *root)
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
		writeApprovalError(stderr, err)
		return 3
	}
	result := applyResult{SchemaVersion: "v0alpha1", Outcome: "dry_run", PlanID: plan.PlanID, PlanDigest: remediation.PlanDigest(plan), SourceRoot: canonical, SourceCommit: commit, Diagnostics: []string{}}
	if *dryRun {
		return writeApplyResult(result, *format, stdout, stderr)
	}

	store := remediation.ApprovalStore{Root: filepath.Join(canonical, ".deprail", "approvals", "consumed")}
	if err := store.Consume(approval.Token); err != nil {
		writeApprovalError(stderr, err)
		return 3
	}
	workspace, err := isolation.Create(ctx, canonical, commit)
	if err != nil {
		writeCLIError(stderr, "WORKTREE_CREATE_FAILED", err.Error(), "fix apply")
		return 3
	}
	result.Workspace = workspace.Path
	cleanup := func() error {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		return workspace.Remove(cleanupCtx)
	}
	defer cleanup()

	for _, command := range plan.Commands {
		workingDirectory := command.WorkingDirectory
		if workingDirectory == "" {
			workingDirectory = "."
		}
		if err := remediation.ValidateCommandWorkspace(plan, workingDirectory); err != nil {
			result.Outcome = "failed"
			result.Diagnostics = append(result.Diagnostics, "unauthorized mutation working directory: "+err.Error())
			_ = cleanup()
			return writeApplyResult(result, *format, stdout, stderr)
		}
		commandWorkspace, err := workspace.Subdirectory(workingDirectory)
		if err != nil {
			result.Outcome = "failed"
			result.Diagnostics = append(result.Diagnostics, "invalid mutation working directory: "+err.Error())
			_ = cleanup()
			return writeApplyResult(result, *format, stdout, stderr)
		}
		executable, err := exec.LookPath(command.Executable)
		if err != nil {
			result.Outcome = "failed"
			result.Diagnostics = append(result.Diagnostics, "mutation executable is unavailable: "+command.Executable)
			_ = cleanup()
			return writeApplyResult(result, *format, stdout, stderr)
		}
		_, err = mutation.Run(ctx, mutation.Request{
			Path: executable, Args: command.Arguments, Workspace: commandWorkspace, Approved: true,
			Timeout: 2 * time.Minute, OutputCap: 16 << 20, DenyScripts: true, DenyNetwork: true,
		})
		if err != nil {
			result.Outcome = "failed"
			result.Diagnostics = append(result.Diagnostics, err.Error())
			_ = cleanup()
			return writeApplyResult(result, *format, stdout, stderr)
		}
	}
	verificationCommands, err := applyVerificationCommands(plan)
	if err != nil {
		result.Outcome = "failed"
		result.Diagnostics = append(result.Diagnostics, "verification unavailable: "+err.Error())
		_ = cleanup()
		return writeApplyResult(result, *format, stdout, stderr)
	}
	_, err = verification.Run(ctx, workspace.Path, verificationCommands, 2*time.Minute, 16<<20)
	if err != nil {
		result.Outcome = "failed"
		result.Diagnostics = append(result.Diagnostics, "verification failed: "+err.Error())
		_ = cleanup()
		return writeApplyResult(result, *format, stdout, stderr)
	}
	result.Diagnostics = append(result.Diagnostics, fmt.Sprintf("verification complete: %d command(s)", len(verificationCommands)))
	result.Outcome = "applied"
	_ = cleanup()
	return writeApplyResult(result, *format, stdout, stderr)
}

func applyVerificationCommands(plan remediation.Plan) ([]verification.Command, error) {
	commands := make([]verification.Command, 0, len(plan.Verification))
	for _, item := range plan.Verification {
		executable, err := exec.LookPath(item.Command.Executable)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", item.ID, err)
		}
		commands = append(commands, verification.Command{
			ID:               item.ID,
			Kind:             verification.Test,
			Path:             executable,
			Args:             append([]string(nil), item.Command.Arguments...),
			WorkingDirectory: item.Command.WorkingDirectory,
			Reason:           item.Reason,
			Enabled:          true,
		})
	}
	return commands, nil
}

func writeApprovalError(stderr io.Writer, err error) {
	var approvalErr *remediation.ApprovalError
	if errors.As(err, &approvalErr) {
		writeCLIError(stderr, string(approvalErr.Code), approvalErr.Message, "approval")
		return
	}
	writeCLIError(stderr, "APPROVAL_INVALID", err.Error(), "approval")
}

func writeApplyResult(result applyResult, format string, stdout, stderr io.Writer) int {
	if format == "json" {
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
	if result.Outcome == "failed" {
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
