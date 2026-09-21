package evidence

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/url"
	"regexp"
	"sort"
	"strings"
)

const SchemaVersion = "v0alpha1"

type Outcome string

const (
	DryRun        Outcome = "dry_run"
	Applied       Outcome = "applied"
	Verified      Outcome = "verified"
	Partial       Outcome = "partial"
	Failed        Outcome = "failed"
	Cancelled     Outcome = "cancelled"
	CleanupFailed Outcome = "cleanup_failed"
)

type Operation struct {
	ID     string `json:"id"`
	Path   string `json:"path"`
	Effect string `json:"effect"`
}
type CommandRecord struct {
	ID           string `json:"id"`
	Tool         string `json:"tool"`
	ExitCode     int    `json:"exit_code"`
	Stdout       string `json:"stdout,omitempty"`
	Stderr       string `json:"stderr,omitempty"`
	StdoutDigest string `json:"stdout_digest"`
	StderrDigest string `json:"stderr_digest"`
}
type Record struct {
	SchemaVersion       string          `json:"schema_version"`
	Outcome             Outcome         `json:"outcome"`
	PlanDigest          string          `json:"plan_digest"`
	SourceCommit        string          `json:"source_commit"`
	SourceRoot          string          `json:"source_root"`
	WorkspaceID         string          `json:"workspace_id"`
	AuthorizedPaths     []string        `json:"authorized_paths"`
	Operations          []Operation     `json:"operations"`
	BeforeDigest        string          `json:"before_digest"`
	AfterDigest         string          `json:"after_digest,omitempty"`
	FindingBeforeDigest string          `json:"finding_before_digest"`
	FindingAfterDigest  string          `json:"finding_after_digest,omitempty"`
	DiffDigest          string          `json:"diff_digest,omitempty"`
	ArtifactDigests     []string        `json:"artifact_digests"`
	VerificationStatus  string          `json:"verification_status"`
	RescanStatus        string          `json:"rescan_status"`
	Commands            []CommandRecord `json:"commands"`
	Diagnostics         []string        `json:"diagnostics"`
	Cleanup             string          `json:"cleanup"`
}

func (r Record) Validate() error {
	if r.SchemaVersion != SchemaVersion || !validOutcome(r.Outcome) || r.PlanDigest == "" || r.SourceCommit == "" || r.SourceRoot == "" || r.WorkspaceID == "" || r.BeforeDigest == "" || r.VerificationStatus == "" || r.RescanStatus == "" || r.Cleanup == "" {
		return errors.New("evidence identity, outcome, and statuses are required")
	}
	if r.AuthorizedPaths == nil || r.Operations == nil || r.Commands == nil || r.Diagnostics == nil || r.ArtifactDigests == nil {
		return errors.New("evidence collections must be present")
	}
	return nil
}
func validOutcome(o Outcome) bool {
	switch o {
	case DryRun, Applied, Verified, Partial, Failed, Cancelled, CleanupFailed:
		return true
	}
	return false
}
func (r Record) Canonical() Record {
	r.Diagnostics = redactAll(r.Diagnostics)
	sort.Strings(r.AuthorizedPaths)
	sort.Strings(r.ArtifactDigests)
	sort.Strings(r.Diagnostics)
	sort.Slice(r.Operations, func(i, j int) bool { return r.Operations[i].ID < r.Operations[j].ID })
	sort.Slice(r.Commands, func(i, j int) bool { return r.Commands[i].ID < r.Commands[j].ID })
	for i := range r.Commands {
		r.Commands[i].Stdout = Redact(r.Commands[i].Stdout)
		r.Commands[i].Stderr = Redact(r.Commands[i].Stderr)
	}
	return r
}
func (r Record) MarshalJSON() ([]byte, error) {
	type plain Record
	return json.Marshal(plain(r.Canonical()))
}
func Digest(r Record) (string, error) {
	r = r.Canonical()
	if err := r.Validate(); err != nil {
		return "", err
	}
	data, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	d := sha256.Sum256(data)
	return hex.EncodeToString(d[:]), nil
}

var sensitiveAssignment = regexp.MustCompile(`(?i)(token|password|secret|authorization|api[_-]?key|aws_secret_access_key)\s*[:=]\s*([^\s&]+)`)
var authHeader = regexp.MustCompile(`(?i)(authorization\s*:\s*bearer\s+)[^\s]+`)
var userInfo = regexp.MustCompile(`(?i)(https?://)[^\s/@]+(?::[^\s/@]*)?@`)

func Redact(value string) string {
	if value == "" {
		return value
	}
	value = authHeader.ReplaceAllString(value, "$1[REDACTED]")
	value = sensitiveAssignment.ReplaceAllString(value, "$1=[REDACTED]")
	value = userInfo.ReplaceAllString(value, "$1[REDACTED]@")
	if parsed, err := url.Parse(value); err == nil && parsed.User != nil {
		parsed.User = nil
		value = parsed.String()
	}
	return value
}
func redactAll(v []string) []string {
	r := make([]string, len(v))
	for i, s := range v {
		r[i] = Redact(s)
	}
	return r
}

var _ = strings.Builder{}
