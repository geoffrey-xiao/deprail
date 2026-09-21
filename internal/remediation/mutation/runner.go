package mutation

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/geoffrey-xiao/deprail/internal/process"
	"github.com/geoffrey-xiao/deprail/internal/remediation/isolation"
)

type Request struct {
	Path        string
	Args        []string
	Workspace   isolation.Workspace
	Approved    bool
	Environment map[string]string
	Timeout     time.Duration
	OutputCap   int64
	DenyScripts bool
	DenyNetwork bool
}

func Run(ctx context.Context, request Request) (process.Result, error) {
	if err := validate(request); err != nil {
		return process.Result{}, err
	}
	args, err := securedArgs(request.Path, request.Args)
	if err != nil {
		return process.Result{}, err
	}
	return process.Run(ctx, process.Request{Path: request.Path, Args: args, Dir: request.Workspace.Path, Env: request.Environment, Timeout: request.Timeout, OutputCap: request.OutputCap})
}

func validate(request Request) error {
	if !request.Approved {
		return errors.New("mutation approval is required")
	}
	if !request.Workspace.Prepared() {
		return errors.New("verified isolated workspace is required")
	}
	if request.Path == "" || !filepath.IsAbs(request.Path) {
		return errors.New("mutation executable must be an absolute path")
	}
	resolved, err := filepath.EvalSymlinks(request.Path)
	if err != nil {
		return errors.New("mutation executable is unavailable")
	}
	info, err := os.Stat(resolved)
	if err != nil || !info.Mode().IsRegular() {
		return errors.New("mutation executable must be a regular file")
	}
	if !supportedExecutable(filepath.Base(resolved)) {
		return fmt.Errorf("unsupported mutation executable: %s", filepath.Base(resolved))
	}
	if request.Timeout <= 0 || request.OutputCap <= 0 {
		return errors.New("positive timeout and output cap are required")
	}

	if !request.DenyScripts || !request.DenyNetwork {
		return errors.New("scripts and network access must be denied")
	}
	return nil
}

func supportedExecutable(name string) bool {
	switch strings.ToLower(name) {
	case "npm", "npm.cmd", "pnpm", "pnpm.cmd", "yarn", "yarn.cmd", "uv", "uv.exe", "pip", "pip3", "pip.exe", "mvn", "mvn.cmd", "gradle", "gradle.bat":
		return true
	default:
		return false
	}
}

func securedArgs(path string, args []string) ([]string, error) {
	result := append([]string(nil), args...)
	base := strings.ToLower(filepath.Base(path))
	flags := []string{}
	switch base {
	case "npm", "npm.cmd", "pnpm", "pnpm.cmd", "yarn", "yarn.cmd":
		flags = []string{"--ignore-scripts", "--offline"}
	case "uv", "uv.exe":
		flags = []string{"--offline"}
	case "pip", "pip3", "pip.exe":
		flags = []string{"--no-index", "--disable-pip-version-check"}
	case "mvn", "mvn.cmd":
		flags = []string{"--offline"}
	case "gradle", "gradle.bat":
		flags = []string{"--offline"}
	default:
		return nil, fmt.Errorf("unsupported mutation executable: %s", base)
	}
	for _, flag := range flags {
		found := false
		for _, arg := range result {
			if arg == flag {
				found = true
				break
			}
		}
		if !found {
			result = append(result, flag)
		}
	}
	return result, nil
}
