package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/geoffrey-xiao/deprail/internal/remediation"
	"github.com/geoffrey-xiao/deprail/internal/remediation/isolation"
)

func runFixApprove(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("fix approve", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	planPath := flags.String("plan", "", "versioned remediation plan")
	root := flags.String("root", ".", "canonical repository root")
	expiresIn := flags.Duration("expires-in", time.Hour, "approval lifetime")
	output := flags.String("output", "", "write approval JSON to a new file")
	if err := flags.Parse(args); err != nil || *planPath == "" || *output == "" || *expiresIn <= 0 || len(flags.Args()) != 0 {
		writeCLIError(stderr, "CONFIG_INVALID", "fix approve requires --plan, --output, and a positive --expires-in", "fix approve")
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
	if plan.RepositoryIdentity.Root != canonical || plan.RepositoryIdentity.Revision != commit {
		writeCLIError(stderr, "APPROVAL_SOURCE_MISMATCH", "plan source identity does not match the current repository", "plan")
		return 3
	}
	approval, err := remediation.NewApproval(plan, canonical, time.Now().UTC().Add(*expiresIn))
	if err != nil {
		writeCLIError(stderr, "APPROVAL_INVALID", err.Error(), "approval")
		return 3
	}
	data, err := json.MarshalIndent(approval, "", "  ")
	if err != nil {
		writeCLIError(stderr, "APPROVAL_WRITE_FAILED", err.Error(), "output")
		return 3
	}
	if err := writeNewPrivateFile(*output, append(data, '\n')); err != nil {
		if errors.Is(err, os.ErrExist) {
			writeCLIError(stderr, "APPROVAL_OUTPUT_EXISTS", "approval output already exists", "output")
			return 3
		}
		writeCLIError(stderr, "APPROVAL_WRITE_FAILED", err.Error(), "output")
		return 3
	}
	_, _ = fmt.Fprintf(stdout, "Approval written: %s\n", *output)
	return 0
}

func writeNewPrivateFile(path string, data []byte) error {
	path = filepath.Clean(path)
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return err
	}
	if _, err := file.Write(data); err != nil {
		_ = file.Close()
		_ = os.Remove(path)
		return err
	}
	if err := file.Sync(); err != nil {
		_ = file.Close()
		_ = os.Remove(path)
		return err
	}
	if err := file.Close(); err != nil {
		_ = os.Remove(path)
		return err
	}
	return nil
}
