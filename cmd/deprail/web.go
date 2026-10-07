package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/geoffrey-xiao/deprail/internal/app"
	"github.com/geoffrey-xiao/deprail/internal/artifact"
	domain "github.com/geoffrey-xiao/deprail/internal/domain/history"
	"github.com/geoffrey-xiao/deprail/internal/presenter"
	"github.com/geoffrey-xiao/deprail/internal/store/history"
	"github.com/geoffrey-xiao/deprail/internal/transport/localhttp"
)

func runWeb(args []string, stdout, stderr io.Writer, assets localhttp.Config) int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	return runWebWithContext(ctx, args, stdout, stderr, os.Stdin, presenter.IsTerminal(os.Stdin), assets, launchBrowser)
}

func runWebWithContext(
	ctx context.Context,
	args []string,
	stdout, stderr io.Writer,
	stdin io.Reader,
	interactive bool,
	assets localhttp.Config,
	launch func(string) error,
) int {
	flags := flag.NewFlagSet("web", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	openBrowser := flags.Bool("open", false, "open the local console")
	artifactRoot := flags.String("artifact-root", "", "existing trusted artifact directory")
	if err := flags.Parse(args); err != nil || len(flags.Args()) != 0 {
		writeCLIError(stderr, "CONFIG_INVALID", "web accepts only --open and --artifact-root", "arguments")
		return 2
	}
	artifactRootProvided := false
	flags.Visit(func(value *flag.Flag) {
		if value.Name == "artifact-root" {
			artifactRootProvided = true
		}
	})
	if artifactRootProvided && *artifactRoot == "" {
		writeCLIError(stderr, "CONFIG_INVALID", "artifact root must name an existing directory", "artifact-root")
		return 2
	}
	if !interactive && !*openBrowser {
		writeCLIError(stderr, "CONFIG_INVALID", "non-interactive web requires --open", "web")
		return 2
	}

	var verifier app.ArtifactVerifier
	if *artifactRoot != "" {
		canonical, err := canonicalArtifactRoot(*artifactRoot)
		if err != nil {
			writeCLIError(stderr, "CONFIG_INVALID", "artifact root must be an existing directory", "artifact-root")
			return 2
		}
		store := artifact.Store{Root: canonical, MaxBytes: 16 << 20}
		verifier = &store
	}

	readStore, openErr := history.OpenReadOnly(ctx, history.Options{})
	if readStore != nil {
		defer func() { _ = readStore.Close() }()
	}
	var queries *app.HistoryService
	if openErr != nil {
		queries = &app.HistoryService{Reader: failedHistoryReader{err: openErr}}
	} else {
		queries = &app.HistoryService{Reader: readStore}
	}
	queries.Artifacts = verifier
	assets.History = queries
	server, err := localhttp.NewServer(assets)
	if err != nil {
		writeCLIError(stderr, "API_LISTENER_UNAVAILABLE", "local console assets or listener are unavailable", "web")
		return 3
	}

	serveResult := make(chan error, 1)
	go func() { serveResult <- server.Serve() }()
	if _, err := fmt.Fprintln(stdout, server.ConsoleURL()); err != nil {
		_ = server.Shutdown()
		writeCLIError(stderr, "OUTPUT_WRITE_FAILED", "could not write the local console address", "stdout")
		return 3
	}
	if *openBrowser {
		if err := launch(server.BootstrapURL()); err != nil {
			writeCLIError(stderr, "API_LISTENER_UNAVAILABLE", "the browser could not be opened", "web")
			if !interactive {
				_ = server.Shutdown()
				return 3
			}
		}
	}
	if interactive {
		_, _ = fmt.Fprintln(stderr, "Press Enter to open or reopen this console; Ctrl-C stops it.")
	}
	enter := make(chan struct{}, 1)
	if interactive {
		go watchEnter(stdin, enter)
	}
	for {
		select {
		case <-ctx.Done():
			_ = server.Shutdown()
			return 0
		case serveErr := <-serveResult:
			_ = server.Close()
			if serveErr == nil && ctx.Err() != nil {
				return 0
			}
			writeCLIError(stderr, "API_LISTENER_UNAVAILABLE", "the local console stopped unexpectedly", "web")
			return 3
		case <-enter:
			if err := launch(server.BootstrapURL()); err != nil {
				writeCLIError(stderr, "API_LISTENER_UNAVAILABLE", "the browser could not be opened", "web")
			}
		}
	}
}

func canonicalArtifactRoot(root string) (string, error) {
	absolute, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	canonical, err := filepath.EvalSymlinks(absolute)
	if err != nil {
		return "", err
	}
	info, err := os.Stat(canonical)
	if err != nil || !info.IsDir() {
		return "", errors.New("artifact root is not an existing directory")
	}
	directory, err := os.Open(canonical)
	if err != nil {
		return "", errors.New("artifact root is not readable")
	}
	if err := directory.Close(); err != nil {
		return "", errors.New("artifact root is not readable")
	}
	return filepath.Clean(canonical), nil
}

type failedHistoryReader struct{ err error }

func (r failedHistoryReader) Get(context.Context, string) (domain.Entry, error) {
	return domain.Entry{}, r.err
}

func (r failedHistoryReader) List(context.Context, domain.Page) (domain.PageResult, error) {
	return domain.PageResult{}, r.err
}

func watchEnter(input io.Reader, enter chan<- struct{}) {
	reader := bufio.NewReaderSize(input, 256)
	for {
		for {
			value, err := reader.ReadByte()
			if err != nil {
				return
			}
			if value == '\n' {
				select {
				case enter <- struct{}{}:
				default:
				}
				break
			}
		}
	}
}

func launchBrowser(target string) error {
	name, args := browserCommand(target)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	command := exec.CommandContext(ctx, name, args...)
	command.Stdout = io.Discard
	command.Stderr = io.Discard
	return command.Run()
}

func browserCommand(target string) (string, []string) {
	switch runtime.GOOS {
	case "darwin":
		return "open", []string{target}
	case "windows":
		return "rundll32", []string{"url.dll,FileProtocolHandler", target}
	default:
		return "xdg-open", []string{target}
	}
}
