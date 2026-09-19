package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestDiscoverJSONWritesDataToStdoutOnly(t *testing.T) {
	fixture := fixturePath(t, "npm-basic")
	var stdout, stderr bytes.Buffer
	code := run([]string{"discover", fixture, "--format", "json"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit code = %d, stderr = %q", code, stderr.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
	var graph struct {
		DocumentType string `json:"document_type"`
		Completeness string `json:"completeness"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &graph); err != nil {
		t.Fatalf("stdout is not JSON: %v", err)
	}
	if graph.DocumentType != "project" || graph.Completeness != "complete" {
		t.Fatalf("graph = %#v", graph)
	}
}

func TestDiscoverIncompleteProjectReturnsCodeThreeAndJSON(t *testing.T) {
	fixture := fixturePath(t, "npm-basic")
	var stdout, stderr bytes.Buffer
	code := run([]string{"discover", fixture, "--format", "json"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("fixture should remain complete, exit code = %d, stderr = %q", code, stderr.String())
	}

	root := t.TempDir()
	writeTestFile(t, filepath.Join(root, "package.json"), `{"name":"incomplete"}`)
	stdout.Reset()
	stderr.Reset()
	code = run([]string{"discover", root, "--format", "json"}, &stdout, &stderr)
	if code != 3 {
		t.Fatalf("exit code = %d, want 3; stderr = %q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"completeness": "partial"`) {
		t.Fatalf("stdout = %q, want partial graph", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty for graph result", stderr.String())
	}
}
func TestDoctorJSONIncludesBuildIdentity(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"doctor", "--format", "json"}, &stdout, &stderr)
	if code != 0 && code != 3 {
		t.Fatalf("exit code = %d, want successful or scanner failure", code)
	}
	var report struct {
		Version string `json:"version"`
		Tag     string `json:"tag"`
		Commit  string `json:"commit"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("stdout is not JSON: %v", err)
	}
	if report.Version == "" || report.Tag == "" || report.Commit == "" {
		t.Fatalf("build identity = %#v", report)
	}
}

func TestReleaseBuildReportsInjectedIdentity(t *testing.T) {
	repositoryRoot := filepath.Join(filepath.Dir(fixturePath(t, "npm-basic")), "..", "..")
	binary := filepath.Join(t.TempDir(), "deprail")
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}
	build := exec.Command("go", "build",
		"-ldflags",
		"-X github.com/geoffrey-xiao/deprail/internal/buildinfo.Version=v0.2.0 -X github.com/geoffrey-xiao/deprail/internal/buildinfo.Tag=v0.2.0-rc.1 -X github.com/geoffrey-xiao/deprail/internal/buildinfo.Commit=abc123",
		"-o", binary, "./cmd/deprail",
	)
	build.Dir = repositoryRoot
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build release binary: %v\n%s", err, output)
	}
	command := exec.Command(binary, "doctor", "--format", "json")
	command.Env = []string{"PATH=" + t.TempDir()}
	output, err := command.Output()
	if exitErr, ok := err.(*exec.ExitError); !ok || exitErr.ExitCode() != 3 {
		t.Fatalf("doctor error = %v, want scanner failure exit code 3", err)
	}
	var report struct {
		Version string `json:"version"`
		Tag     string `json:"tag"`
		Commit  string `json:"commit"`
	}
	if err := json.Unmarshal(output, &report); err != nil {
		t.Fatalf("doctor output is not JSON: %v", err)
	}
	if report.Version != "v0.2.0" || report.Tag != "v0.2.0-rc.1" || report.Commit != "abc123" {
		t.Fatalf("release identity = %#v", report)
	}
}

func TestDiscoverRejectsMultiplePaths(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"discover", ".", "other"}, &stdout, &stderr)
	if code != 2 || stdout.Len() != 0 || !strings.Contains(stderr.String(), "CONFIG_INVALID") {
		t.Fatalf("code=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
}

func TestRejectsUnexpectedCommandArguments(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "doctor positional", args: []string{"doctor", "unexpected"}},
		{name: "discover option", args: []string{"discover", "--unsupported"}},
		{name: "scan option", args: []string{"scan", "--unsupported"}},
		{name: "doctor option", args: []string{"doctor", "--unsupported"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			if code := run(test.args, &stdout, &stderr); code != 2 {
				t.Fatalf("exit code = %d, want 2; stdout=%q stderr=%q", code, stdout.String(), stderr.String())
			}
			if stdout.Len() != 0 || !strings.Contains(stderr.String(), "CONFIG_INVALID") {
				t.Fatalf("stdout=%q stderr=%q", stdout.String(), stderr.String())
			}
		})
	}
}

func TestScanFromOutsideRootUsesTargetAndContainsArtifacts(t *testing.T) {
	raw := []byte(`{"results":[{"packages":[{"package":{"name":"target-only","version":"1.0.0"},"vulnerabilities":[{"id":"OSV-TARGET","aliases":["CVE-TARGET"],"database_specific":{"severity":"HIGH"},"affected":[{"ranges":[{"events":[{"introduced":"0"},{"fixed":"1.0.1"}]}]}]}]}]}]}`)
	sum := sha256.Sum256(raw)
	expectedDigest := hex.EncodeToString(sum[:])
	repositoryRoot := filepath.Join(filepath.Dir(fixturePath(t, "npm-basic")), "..", "..")
	cliPath := filepath.Join(t.TempDir(), "deprail")
	if runtime.GOOS == "windows" {
		cliPath += ".exe"
	}
	build := exec.Command("go", "build", "-o", cliPath, "./cmd/deprail")
	build.Dir = repositoryRoot
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build deprail: %v\n%s", err, output)
	}
	for _, absolute := range []bool{true, false} {
		t.Run(map[bool]string{true: "absolute target", false: "relative target"}[absolute], func(t *testing.T) {
			target := t.TempDir()
			writeTestFile(t, filepath.Join(target, "package.json"), `{"name":"target"}`)
			writeTestFile(t, filepath.Join(target, "package-lock.json"), `{"name":"target","lockfileVersion":3,"packages":{}}`)
			writeTestFile(t, filepath.Join(target, "target.marker"), "")
			caller := t.TempDir()
			writeTestFile(t, filepath.Join(caller, "caller.marker"), "")
			bin := t.TempDir()
			scannerName := "osv-scanner"
			if runtime.GOOS == "windows" {
				scannerName += ".exe"
			}
			executable, err := os.ReadFile(mustExecutable(t))
			if err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(bin, scannerName), executable, 0o755); err != nil {
				t.Fatal(err)
			}
			targetArg := target
			if !absolute {
				targetArg, err = filepath.Rel(caller, target)
				if err != nil {
					t.Fatal(err)
				}
			}
			home := t.TempDir()
			command := exec.Command(cliPath, "scan", targetArg, "--format", "json")
			command.Dir = caller
			command.Env = append(os.Environ(),
				"PATH="+bin+string(os.PathListSeparator)+os.Getenv("PATH"),
				"HOME="+home,
			)
			stdout, stderr, err := commandOutput(command)
			if err != nil {
				t.Fatalf("scan: %v stdout=%q stderr=%q", err, stdout, stderr)
			}
			var report struct {
				Status   string `json:"status"`
				Findings []struct {
					TargetID string `json:"TargetID"`
				} `json:"findings"`
				ArtifactDigests []string `json:"artifact_digests"`
			}
			if err := json.Unmarshal(stdout, &report); err != nil {
				t.Fatalf("stdout is not JSON: %v", err)
			}
			if report.Status != "complete" || len(report.Findings) != 1 || report.Findings[0].TargetID != "OSV-TARGET" {
				t.Fatalf("report = %#v", report)
			}
			if len(report.ArtifactDigests) != 1 || report.ArtifactDigests[0] != expectedDigest {
				t.Fatalf("artifact digests = %#v, want %q", report.ArtifactDigests, expectedDigest)
			}
			artifactPath := filepath.Join(target, ".deprail", "artifacts", expectedDigest[:2], expectedDigest)
			stored, err := os.ReadFile(artifactPath)
			if err != nil || string(stored) != string(raw) {
				t.Fatalf("stored artifact = %q, err=%v", stored, err)
			}
			if _, err := os.Stat(filepath.Join(caller, ".deprail")); !os.IsNotExist(err) {
				t.Fatalf("caller received artifact tree: %v", err)
			}
			scannerCWD, err := os.ReadFile(filepath.Join(home, "scanner-cwd"))
			if err != nil {
				t.Fatal(err)
			}
			targetInfo, err := os.Stat(target)
			if err != nil {
				t.Fatal(err)
			}
			scannerInfo, err := os.Stat(strings.TrimSpace(string(scannerCWD)))
			if err != nil || !os.SameFile(targetInfo, scannerInfo) {
				t.Fatalf("scanner cwd=%q target=%q err=%v", scannerCWD, target, err)
			}
		})
	}
}

func commandOutput(command *exec.Cmd) ([]byte, []byte, error) {
	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	err := command.Run()
	return stdout.Bytes(), stderr.Bytes(), err
}

func TestMain(m *testing.M) {
	if strings.HasPrefix(filepath.Base(os.Args[0]), "osv-scanner") {
		runScannerHelper()
		return
	}
	os.Exit(m.Run())
}

func runScannerHelper() {
	cwd, err := os.Getwd()
	if err != nil {
		os.Exit(9)
	}
	if home := os.Getenv("HOME"); home != "" {
		_ = os.WriteFile(filepath.Join(home, "scanner-cwd"), []byte(cwd), 0o600)
	}
	for _, arg := range os.Args[1:] {
		if arg == "--version" {
			_, _ = os.Stdout.Write([]byte("osv-scanner 1.0.0"))
			os.Exit(0)
		}
	}
	if len(os.Args) == 0 {
		os.Exit(9)
	}
	targetArg := os.Args[len(os.Args)-1]
	target := targetArg
	if !filepath.IsAbs(target) {
		target = filepath.Join(cwd, target)
	}
	targetID := "OSV-CALLER"
	if _, err := os.Stat(filepath.Join(filepath.Clean(target), "target.marker")); err == nil {
		targetID = "OSV-TARGET"
	}
	output := `{"results":[{"packages":[{"package":{"name":"target-only","version":"1.0.0"},"vulnerabilities":[{"id":"` + targetID + `","aliases":["CVE-TARGET"],"database_specific":{"severity":"HIGH"},"affected":[{"ranges":[{"events":[{"introduced":"0"},{"fixed":"1.0.1"}]}]}]}]}]}]}`
	_, _ = os.Stdout.Write([]byte(output))
	os.Exit(0)
}

func mustExecutable(t *testing.T) string {
	t.Helper()
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	return executable
}

func fixturePath(t *testing.T, fixture string) string {
	t.Helper()
	_, source, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("locate command test source")
	}
	return filepath.Join(filepath.Dir(source), "..", "..", "testdata", "fixtures", fixture)
}

func writeTestFile(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}
