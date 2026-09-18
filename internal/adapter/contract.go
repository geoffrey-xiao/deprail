// Package adapter defines scanner-independent contracts for external vulnerability scanners.
package adapter

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

// ErrorCode identifies stable adapter failure classes.
type ErrorCode string

const (
	ErrInvalidPlan       ErrorCode = "CONFIG_INVALID"
	ErrUnsupportedTarget ErrorCode = "SCANNER_VERSION_UNSUPPORTED"
	ErrExecution         ErrorCode = "SCANNER_EXIT_NONZERO"
	ErrTimeout           ErrorCode = "SCANNER_TIMEOUT"
	ErrOutputLimit       ErrorCode = "SCANNER_OUTPUT_LIMIT"
	ErrInvalidOutput     ErrorCode = "SCANNER_OUTPUT_INVALID"
)

// Error is a safe, stable adapter error. Message must not contain secrets or raw scanner output.
type Error struct {
	Code      ErrorCode
	Message   string
	Scope     string
	Retryable bool
}

func (e *Error) Error() string {
	if e.Scope == "" {
		return fmt.Sprintf("%s: %s", e.Code, e.Message)
	}
	return fmt.Sprintf("%s (%s): %s", e.Code, e.Scope, e.Message)
}

// Metadata describes a scanner and the target kinds it accepts.
type Metadata struct {
	Name             string
	Version          string
	SupportedTargets []string
}

// Target is a repository-relative scan target. Paths use slash separators.
type Target struct {
	WorkspaceID  string
	RelativePath string
	Ecosystem    string
	PackageFiles []string
}

// Plan is an immutable, deterministic set of scanner targets.
type Plan struct {
	Targets []Target
}

// RawResult contains bounded scanner output and execution metadata.
type RawResult struct {
	Stdout    []byte
	Stderr    []byte
	ExitCode  int
	Truncated bool
}

// Record is scanner-independent parsed evidence. Adapters map tool output into this shape.
type Record struct {
	Component string
	Version   string
	Aliases   []string
	Severity  string
	Fixed     string
	TargetID  string
}

// Finding is the normalized adapter boundary; downstream normalization owns identity and ordering.
type Finding struct {
	Component string
	Version   string
	Aliases   []string
	Severity  string
	Fixed     string
	TargetID  string
}

// ScannerAdapter separates metadata, planning, execution, parsing, and normalization.
type ScannerAdapter interface {
	Metadata(ctx context.Context) (Metadata, error)
	Compatible(ctx context.Context, metadata Metadata) error
	Plan(ctx context.Context, targets []Target) (Plan, error)
	Execute(ctx context.Context, plan Plan) (RawResult, error)
	Parse(ctx context.Context, raw RawResult) ([]Record, error)
	Normalize(ctx context.Context, records []Record) ([]Finding, error)
}

// ValidatePlan checks adapter-independent plan invariants before execution.
func ValidatePlan(plan Plan) error {
	if len(plan.Targets) == 0 {
		return &Error{Code: ErrInvalidPlan, Message: "scan plan has no targets"}
	}
	seen := make(map[string]struct{}, len(plan.Targets))
	for _, target := range plan.Targets {
		if target.WorkspaceID == "" || target.RelativePath == "" || target.Ecosystem == "" {
			return &Error{Code: ErrInvalidPlan, Message: "scan target identity is required"}
		}
		if target.RelativePath[0] == '/' || target.RelativePath == ".." || strings.HasPrefix(target.RelativePath, "../") || strings.Contains(target.RelativePath, "/../") {
			return &Error{Code: ErrInvalidPlan, Message: "scan target path must be workspace-relative", Scope: target.RelativePath}
		}
		if _, ok := seen[target.WorkspaceID]; ok {
			return &Error{Code: ErrInvalidPlan, Message: "scan target IDs must be unique", Scope: target.WorkspaceID}
		}
		seen[target.WorkspaceID] = struct{}{}
	}
	return nil
}

// IsCode reports whether err or one of its causes has the requested stable code.
func IsCode(err error, code ErrorCode) bool {
	var adapterErr *Error
	return errors.As(err, &adapterErr) && adapterErr.Code == code
}
