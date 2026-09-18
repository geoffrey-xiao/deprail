package osv

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/geoffrey-xiao/deprail/internal/adapter"
)

func TestScannerFailureMatrix(t *testing.T) {
	tests := []struct {
		name    string
		mode    string
		want    adapter.ErrorCode
		timeout time.Duration
		cap     int64
	}{
		{"missing", "missing", adapter.ErrScannerNotFound, time.Second, 1024},
		{"nonzero", "exit", adapter.ErrExecution, time.Second, 1024},
		{"timeout", "sleep", adapter.ErrTimeout, 20 * time.Millisecond, 1024},
		{"oversized", "large", adapter.ErrOutputLimit, time.Second, 32},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			scanner := Scanner{Path: os.Args[0], Args: []string{"-test.run=TestScannerHelper", "mode=" + test.mode}, Timeout: test.timeout, OutputCap: test.cap}
			if test.mode == "missing" {
				scanner.Path = "/definitely/missing/osv-scanner"
			}
			_, err := scanner.Execute(context.Background(), adapter.Plan{Targets: []adapter.Target{{WorkspaceID: "root", RelativePath: ".", Ecosystem: "npm"}}})
			if !adapter.IsCode(err, test.want) {
				t.Fatalf("error = %v, want %s", err, test.want)
			}
		})
	}
}

func TestScannerCancellationIsFailure(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := (Scanner{Path: os.Args[0], Timeout: time.Second, OutputCap: 1024}).Execute(ctx, adapter.Plan{Targets: []adapter.Target{{WorkspaceID: "root", RelativePath: ".", Ecosystem: "npm"}}})
	if err == nil {
		t.Fatal("expected cancelled execution failure")
	}
}

func TestScannerExitMatrix(t *testing.T) {
	tests := []struct {
		name         string
		mode         string
		wantErr      adapter.ErrorCode
		wantParseErr adapter.ErrorCode
		wantRecord   bool
	}{
		{name: "zero findings exit zero", mode: "empty", wantRecord: false},
		{name: "valid findings exit one", mode: "vulnerable", wantRecord: true},
		{name: "invalid output exit one", mode: "invalid-vulnerable", wantParseErr: adapter.ErrInvalidOutput},
		{name: "unexpected exit", mode: "exit", wantErr: adapter.ErrExecution},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			scanner := Scanner{Path: os.Args[0], Args: []string{"-test.run=TestScannerHelper", "mode=" + test.mode}, Timeout: time.Second, OutputCap: 4096}
			raw, err := scanner.Execute(context.Background(), adapter.Plan{Targets: []adapter.Target{{WorkspaceID: "root", RelativePath: ".", Ecosystem: "npm"}}})
			if test.wantErr != "" {
				if !adapter.IsCode(err, test.wantErr) {
					t.Fatalf("execute error = %v, want %s", err, test.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("execute: %v", err)
			}
			records, parseErr := scanner.Parse(context.Background(), raw)
			if test.wantParseErr != "" {
				if !adapter.IsCode(parseErr, test.wantParseErr) {
					t.Fatalf("parse error = %v, want %s", parseErr, test.wantParseErr)
				}
				return
			}
			if test.wantRecord {
				if parseErr != nil || len(records) == 0 {
					t.Fatalf("records = %#v, parse error = %v", records, parseErr)
				}
			} else if parseErr != nil || len(records) != 0 {
				t.Fatalf("empty records = %#v, parse error = %v", records, parseErr)
			}
		})
	}
}

func TestScannerHelper(t *testing.T) {
	mode := ""
	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "mode=") {
			mode = strings.TrimPrefix(arg, "mode=")
		}
	}
	if mode == "" {
		return
	}
	switch mode {
	case "exit":
		os.Exit(7)
	case "sleep":
		time.Sleep(time.Second)
	case "large":
		for range 128 {
			os.Stdout.WriteString("x")
		}
	case "empty":
		_, _ = os.Stdout.WriteString(`{"results":[]}`)
		os.Exit(0)
	case "cwd":
		cwd, err := os.Getwd()
		if err != nil {
			os.Exit(9)
		}
		_, _ = fmt.Fprintln(os.Stderr, cwd)
	case "args":
		want := []string{"scan", "source", "--format", "json", "."}
		for i := range len(os.Args) - len(want) + 1 {
			match := true
			for j := range want {
				if os.Args[i+j] != want[j] {
					match = false
					break
				}
			}
			if match {
				_, _ = os.Stdout.WriteString(`{"results":[]}`)
				return
			}
		}
		os.Exit(8)
	case "vulnerable":
		_, _ = os.Stdout.WriteString(`{"results":[{"packages":[{"package":{"name":"lodash","version":"4.17.20"},"vulnerabilities":[{"id":"GHSA-test","aliases":["CVE-test"],"database_specific":{"severity":"HIGH"},"severity":[{"score":"CVSS:3.1/AV:N"}],"affected":[{"ranges":[{"events":[{"introduced":"0"},{"fixed":"4.17.21"}]}]}]}]}]}]}`)
		os.Exit(1)
	case "invalid-vulnerable":
		_, _ = os.Stdout.WriteString("not json")
		os.Exit(1)
	}
}
