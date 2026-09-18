package app

import (
	"context"
	"os/exec"
	"runtime"
	"time"

	"github.com/geoffrey-xiao/deprail/internal/adapters/osv"
	"github.com/geoffrey-xiao/deprail/internal/buildinfo"
	"github.com/geoffrey-xiao/deprail/internal/process"
)

type DoctorReport struct {
	Version string        `json:"version"`
	Tag     string        `json:"tag"`
	Commit  string        `json:"commit"`
	OS      string        `json:"os"`
	Arch    string        `json:"arch"`
	Scanner DoctorScanner `json:"scanner"`
}
type DoctorScanner struct {
	Path       string `json:"path,omitempty"`
	Version    string `json:"version,omitempty"`
	Available  bool   `json:"available"`
	Compatible bool   `json:"compatible"`
	Error      string `json:"error,omitempty"`
	Help       string `json:"help,omitempty"`
}

func Doctor(ctx context.Context) DoctorReport {
	identity := buildinfo.Current()
	report := DoctorReport{
		Version: identity.Version,
		Tag:     identity.Tag,
		Commit:  identity.Commit,
		OS:      runtime.GOOS,
		Arch:    runtime.GOARCH,
		Scanner: DoctorScanner{Help: "Install OSV-Scanner and rerun deprail doctor"},
	}
	path, err := exec.LookPath("osv-scanner")
	if err != nil {
		report.Scanner.Error = "OSV-Scanner is not available"
		return report
	}
	report.Scanner.Path = path
	report.Scanner.Available = true
	result, err := process.Run(ctx, process.Request{Path: path, Args: []string{"--version"}, Timeout: 5 * time.Second, OutputCap: 4096})
	if err != nil {
		report.Scanner.Error = "OSV-Scanner version check failed"
		return report
	}
	metadata, err := osv.ParseVersion(string(result.Stdout))
	if err != nil {
		report.Scanner.Error = "OSV-Scanner version output is invalid"
		return report
	}
	report.Scanner.Version = metadata.Version
	if err := osv.Compatible(ctx, metadata); err != nil {
		report.Scanner.Error = "OSV-Scanner version is unsupported"
		return report
	}
	report.Scanner.Compatible = true
	report.Scanner.Help = ""
	return report
}
