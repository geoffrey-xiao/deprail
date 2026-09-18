package process

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"
)

func TestRunUsesArgumentArrayAndReturnsOutput(t *testing.T) {
	result, err := Run(context.Background(), Request{Path: os.Args[0], Args: []string{"-test.run=TestProcessHelper", "--", "hello", "world"}, Timeout: time.Second, OutputCap: 1024})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if !strings.HasPrefix(string(result.Stdout), "hello world\n") || result.ExitCode != 0 {
		t.Fatalf("result = %#v", result)
	}
}

func TestRunRejectsMissingExecutableWithoutLeakingPath(t *testing.T) {
	_, err := Run(context.Background(), Request{Path: "/definitely/missing/deprail-scanner", Timeout: time.Second, OutputCap: 1024})
	if !IsCode(err, ErrNotFound) {
		t.Fatalf("error = %v, want %s", err, ErrNotFound)
	}
	if err.Error() == "" || contains(err.Error(), "definitely") {
		t.Fatalf("error leaked executable path: %v", err)
	}
}

func TestRunClassifiesTimeoutAndOutputLimit(t *testing.T) {
	_, timeoutErr := Run(context.Background(), Request{Path: os.Args[0], Args: []string{"-test.run=TestProcessHelper", "--", "sleep"}, Timeout: 20 * time.Millisecond, OutputCap: 1024})
	if !IsCode(timeoutErr, ErrTimeout) {
		t.Fatalf("timeout error = %v", timeoutErr)
	}
	_, limitErr := Run(context.Background(), Request{Path: os.Args[0], Args: []string{"-test.run=TestProcessHelper", "--", "large"}, Timeout: time.Second, OutputCap: 32})
	if !IsCode(limitErr, ErrOutputLimit) {
		t.Fatalf("limit error = %v", limitErr)
	}
}

func TestProcessHelper(t *testing.T) {
	mode := ""
	for _, arg := range os.Args[1:] {
		if arg == "hello" || arg == "sleep" || arg == "large" {
			mode = arg
		}
	}
	if mode == "" {
		return
	}
	switch mode {
	case "hello":
		os.Stdout.WriteString("hello world\n")
	case "sleep":
		time.Sleep(time.Second)
	case "large":
		for i := 0; i < 128; i++ {
			os.Stdout.WriteString("x")
		}
	}
}

func contains(s, needle string) bool {
	for i := 0; i+len(needle) <= len(s); i++ {
		if s[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}
