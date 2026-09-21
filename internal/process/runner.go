package process

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type ErrorCode string

const (
	ErrInvalidRequest ErrorCode = "CONFIG_INVALID"
	ErrNotFound       ErrorCode = "SCANNER_NOT_FOUND"
	ErrTimeout        ErrorCode = "SCANNER_TIMEOUT"
	ErrCancelled      ErrorCode = "SCANNER_CANCELLED"
	ErrExitNonzero    ErrorCode = "SCANNER_EXIT_NONZERO"
	ErrOutputLimit    ErrorCode = "SCANNER_OUTPUT_LIMIT"
)

type Error struct {
	Code      ErrorCode
	Message   string
	Retryable bool
}

func (e *Error) Error() string { return fmt.Sprintf("%s: %s", e.Code, e.Message) }

func IsCode(err error, code ErrorCode) bool {
	var processErr *Error
	return errors.As(err, &processErr) && processErr.Code == code
}

type Request struct {
	Path      string
	Args      []string
	Dir       string
	Env       map[string]string
	Timeout   time.Duration
	OutputCap int64
}

type Result struct {
	Stdout   []byte
	Stderr   []byte
	ExitCode int
}

func (r Request) validate() error {
	if strings.TrimSpace(r.Path) == "" || r.Timeout <= 0 || r.OutputCap <= 0 {
		return &Error{Code: ErrInvalidRequest, Message: "path, positive timeout, and positive output cap are required"}
	}
	return nil
}

func Run(ctx context.Context, request Request) (Result, error) {
	if err := request.validate(); err != nil {
		return Result{}, err
	}
	ctx, cancel := context.WithTimeout(ctx, request.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, request.Path, request.Args...)
	cmd.Dir = request.Dir
	cmd.Env = approvedEnv(request.Env)
	configureProcessGroup(cmd)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &limitedBuffer{Buffer: &stdout, Limit: request.OutputCap}
	cmd.Stderr = &limitedBuffer{Buffer: &stderr, Limit: request.OutputCap}
	err := cmd.Run()
	if stdout.Len() >= int(request.OutputCap) || stderr.Len() >= int(request.OutputCap) {
		terminateProcessGroup(cmd)
		return Result{Stdout: stdout.Bytes(), Stderr: stderr.Bytes()}, &Error{Code: ErrOutputLimit, Message: "process output exceeded configured limit"}
	}
	if errors.Is(ctx.Err(), context.Canceled) {
		terminateProcessGroup(cmd)
		return Result{Stdout: stdout.Bytes(), Stderr: stderr.Bytes()}, &Error{Code: ErrCancelled, Message: "process was cancelled", Retryable: false}
	}
	if ctx.Err() != nil {
		terminateProcessGroup(cmd)
		return Result{Stdout: stdout.Bytes(), Stderr: stderr.Bytes()}, &Error{Code: ErrTimeout, Message: "process exceeded configured deadline", Retryable: true}
	}
	if err != nil {
		var execErr *exec.Error
		if errors.As(err, &execErr) || errors.Is(err, exec.ErrNotFound) || errors.Is(err, os.ErrNotExist) {
			return Result{}, &Error{Code: ErrNotFound, Message: "executable is not available"}
		}
		code := exitCode(err)
		return Result{Stdout: stdout.Bytes(), Stderr: stderr.Bytes(), ExitCode: code}, &Error{Code: ErrExitNonzero, Message: "process exited unsuccessfully"}
	}
	return Result{Stdout: stdout.Bytes(), Stderr: stderr.Bytes(), ExitCode: 0}, nil
}

func approvedEnv(extra map[string]string) []string {
	allowed := map[string]bool{"PATH": true, "HOME": true, "LANG": true, "LC_ALL": true, "SYSTEMROOT": true, "TEMP": true, "TMP": true, "USERPROFILE": true}
	values := make(map[string]string)
	for _, entry := range os.Environ() {
		key, value, ok := strings.Cut(entry, "=")
		if ok && allowed[strings.ToUpper(key)] {
			values[key] = value
		}
	}
	for key, value := range extra {
		if allowed[strings.ToUpper(key)] {
			values[key] = value
		}
	}
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	result := make([]string, 0, len(keys))
	for _, key := range keys {
		result = append(result, key+"="+values[key])
	}
	return result
}

type limitedBuffer struct {
	*bytes.Buffer
	Limit int64
}

func (b *limitedBuffer) Write(p []byte) (int, error) {
	remaining := b.Limit - int64(b.Len())
	if remaining <= 0 {
		return 0, io.ErrShortBuffer
	}
	if int64(len(p)) > remaining {
		p = p[:remaining]
		_, _ = b.Buffer.Write(p)
		return len(p), io.ErrShortBuffer
	}
	return b.Buffer.Write(p)
}
