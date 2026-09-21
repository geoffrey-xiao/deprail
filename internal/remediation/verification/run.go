package verification

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/geoffrey-xiao/deprail/internal/process"
)

type Result struct {
	Command Command
	Process process.Result
	Err     error
}

// Run executes selected verification commands in the isolated workspace. Results
// before the first failure are retained; no later command is started after a
// failure, timeout, cancellation, or output-limit error.
func Run(ctx context.Context, workspace string, commands []Command, timeout time.Duration, outputCap int64) ([]Result, error) {
	if workspace == "" || !filepath.IsAbs(workspace) {
		return nil, errors.New("verification workspace must be an absolute path")
	}
	root, err := filepath.EvalSymlinks(workspace)
	if err != nil {
		return nil, fmt.Errorf("resolve verification workspace: %w", err)
	}
	info, err := os.Stat(root)
	if err != nil || !info.IsDir() {
		return nil, errors.New("verification workspace must be a directory")
	}
	selected, err := SelectCommands(commands)
	if err != nil {
		return nil, err
	}
	results := make([]Result, 0, len(selected))
	for _, command := range selected {
		dir := filepath.Join(root, filepath.FromSlash(command.WorkingDirectory))
		resolvedDir, resolveErr := filepath.EvalSymlinks(dir)
		if resolveErr != nil {
			return results, fmt.Errorf("verification %q working directory: %w", command.ID, resolveErr)
		}
		relative, relErr := filepath.Rel(root, resolvedDir)
		if relErr != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
			return results, fmt.Errorf("verification %q working directory escapes workspace", command.ID)
		}
		result, runErr := process.Run(ctx, process.Request{Path: command.Path, Args: command.Args, Dir: resolvedDir, Timeout: timeout, OutputCap: outputCap})
		results = append(results, Result{Command: command, Process: result, Err: runErr})
		if runErr != nil {
			return results, fmt.Errorf("verification %q failed: %w", command.ID, runErr)
		}
	}
	return results, nil
}
