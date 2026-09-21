package mutation

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/geoffrey-xiao/deprail/internal/process"
)

type Request struct {
	Path          string
	Args          []string
	WorkspaceRoot string
	Approved      bool
	Environment   map[string]string
	Timeout       time.Duration
	OutputCap     int64
	AllowScripts  bool
	AllowNetwork  bool
}

func Run(ctx context.Context, request Request) (process.Result, error) {
	if err := validate(request); err != nil {
		return process.Result{}, err
	}
	return process.Run(ctx, process.Request{Path: request.Path, Args: append([]string(nil), request.Args...), Dir: request.WorkspaceRoot, Env: request.Environment, Timeout: request.Timeout, OutputCap: request.OutputCap})
}

func validate(request Request) error {
	if !request.Approved {
		return errors.New("mutation approval is required")
	}
	if request.Path == "" || !filepath.IsAbs(request.Path) {
		return errors.New("mutation executable must be an absolute path")
	}
	if request.WorkspaceRoot == "" || !filepath.IsAbs(request.WorkspaceRoot) {
		return errors.New("mutation workspace must be an absolute path")
	}
	base := strings.ToLower(filepath.Base(request.Path))
	switch base {
	case "sh", "bash", "zsh", "fish", "cmd", "cmd.exe", "powershell", "powershell.exe", "pwsh", "pwsh.exe":
		return fmt.Errorf("shell executable is not allowed: %s", base)
	}
	if request.Timeout <= 0 || request.OutputCap <= 0 {
		return errors.New("positive timeout and output cap are required")
	}
	if request.AllowScripts || request.AllowNetwork {
		return errors.New("scripts and network access are denied by default")
	}
	return nil
}
