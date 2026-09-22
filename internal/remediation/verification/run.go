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

var supportedTools = map[string]struct{}{
	"gradle": {}, "gradle.bat": {}, "mvn": {}, "mvn.cmd": {},
	"node": {}, "node.exe": {}, "npm": {}, "npm.cmd": {}, "npm.exe": {}, "pip": {}, "pip3": {}, "pip.exe": {},
	"pnpm": {}, "pnpm.cmd": {}, "python": {}, "python3": {}, "python.exe": {},
	"uv": {}, "uv.exe": {}, "yarn": {}, "yarn.cmd": {},
}

var ErrNoCommands = errors.New("no verification commands are available")

func trustedTool(path string) bool {
	_, ok := supportedTools[strings.ToLower(filepath.Base(path))]
	return ok
}

// ValidateCommandSpec validates a plan-level command before approval.
func ValidateCommandSpec(path string, args []string) error {
	base := strings.ToLower(filepath.Base(path))
	if _, ok := supportedTools[base]; !ok {
		return fmt.Errorf("unsupported verification executable %q", path)
	}
	if base == "node" || base == "node.exe" {
		if len(args) != 2 || args[0] != "--check" || filepath.IsAbs(args[1]) {
			return errors.New("node verification must be exactly: node --check <relative-file>")
		}
		clean := filepath.Clean(args[1])
		if clean != args[1] || clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
			return errors.New("node verification input must remain repository-relative")
		}
	}
	return nil
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
	if len(selected) == 0 {
		return nil, ErrNoCommands
	}
	results := make([]Result, 0, len(selected))
	for _, command := range selected {
		if !trustedTool(command.Path) {
			return results, fmt.Errorf("verification %q uses an untrusted executable", command.ID)
		}
		dir := filepath.Join(root, filepath.FromSlash(command.WorkingDirectory))
		resolvedDir, resolveErr := filepath.EvalSymlinks(dir)
		if resolveErr != nil {
			return results, fmt.Errorf("verification %q working directory: %w", command.ID, resolveErr)
		}
		relative, relErr := filepath.Rel(root, resolvedDir)
		if relErr != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
			return results, fmt.Errorf("verification %q working directory escapes workspace", command.ID)
		}
		environment := map[string]string{
			"HOME":        root,
			"TEMP":        root,
			"TMP":         root,
			"USERPROFILE": root,
		}
		result, runErr := process.Run(ctx, process.Request{Path: command.Path, Args: command.Args, Dir: resolvedDir, Env: environment, Timeout: timeout, OutputCap: outputCap})
		results = append(results, Result{Command: command, Process: result, Err: runErr})
		if runErr != nil {
			return results, fmt.Errorf("verification %q failed: %w", command.ID, runErr)
		}
	}
	return results, nil
}
