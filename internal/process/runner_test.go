package process

import (
	"context"
	"os"
	"strconv"
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

func TestRunBoundsCapturedOutput(t *testing.T) {
	for _, stream := range []string{"stdout", "stderr"} {
		for _, test := range []struct {
			name    string
			count   int
			chunked bool
		}{
			{"below", 31, false},
			{"at", 32, false},
			{"over", 33, false},
			{"large-single-write", 256 * 1024, false},
			{"large-chunked", 256 * 1024, true},
		} {
			t.Run(stream+"/"+test.name, func(t *testing.T) {
				mode := "bounded"
				if test.chunked {
					mode = "bounded-chunks"
				}
				result, err := Run(context.Background(), Request{
					Path: os.Args[0], Args: []string{"-test.run=^TestProcessHelper$", "--", mode, stream, strconv.Itoa(test.count)},
					Timeout: time.Second, OutputCap: 32,
				})
				if test.count >= 32 {
					if !IsCode(err, ErrOutputLimit) {
						t.Fatalf("error = %v, want %s", err, ErrOutputLimit)
					}
				} else if err != nil || result.ExitCode != 0 {
					t.Fatalf("exit=%d error=%v", result.ExitCode, err)
				}
				output, other := result.Stdout, result.Stderr
				if stream == "stderr" {
					output, other = result.Stderr, result.Stdout
				}
				if len(output) != min(test.count, 32) {
					t.Fatalf("captured %d bytes, want %d", len(output), min(test.count, 32))
				}
				if string(output) != strings.Repeat("x", min(test.count, 32)) || len(other) != 0 {
					t.Fatalf("unexpected captured prefix or other-stream bytes: %q, %d", output, len(other))
				}
			})
		}
	}
}

func TestProcessHelper(t *testing.T) {
	mode := ""
	for _, arg := range os.Args[1:] {
		if arg == "hello" || arg == "sleep" || arg == "large" || arg == "bounded" || arg == "bounded-chunks" {
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
	if mode == "bounded" || mode == "bounded-chunks" {
		count, err := strconv.Atoi(os.Args[len(os.Args)-1])
		if err != nil || count < 0 {
			os.Exit(2)
		}
		output := os.Stdout
		if os.Args[len(os.Args)-2] == "stderr" {
			output = os.Stderr
		}
		if mode == "bounded" {
			_, _ = output.WriteString(strings.Repeat("x", count))
		} else {
			chunk := strings.Repeat("x", 1024)
			for remaining := count; remaining > 0; {
				size := min(remaining, len(chunk))
				_, _ = output.WriteString(chunk[:size])
				remaining -= size
			}
		}
		os.Exit(0)
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
