package evidence

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/url"
	"regexp"
)

const SchemaVersion = "v0alpha1"

type Outcome string

const (
	Applied       Outcome = "applied"
	Verified      Outcome = "verified"
	Partial       Outcome = "partial"
	Failed        Outcome = "failed"
	Cancelled     Outcome = "cancelled"
	CleanupFailed Outcome = "cleanup_failed"
)

type CommandRecord struct {
	ID       string `json:"id"`
	Tool     string `json:"tool"`
	ExitCode int    `json:"exit_code"`
	Stdout   string `json:"stdout,omitempty"`
	Stderr   string `json:"stderr,omitempty"`
}

type Record struct {
	SchemaVersion string          `json:"schema_version"`
	Outcome       Outcome         `json:"outcome"`
	PlanDigest    string          `json:"plan_digest"`
	SourceCommit  string          `json:"source_commit"`
	SourceRoot    string          `json:"source_root"`
	BeforeDigest  string          `json:"before_digest"`
	AfterDigest   string          `json:"after_digest,omitempty"`
	Commands      []CommandRecord `json:"commands"`
	Diagnostics   []string        `json:"diagnostics"`
	Cleanup       string          `json:"cleanup"`
}

func (r Record) Validate() error {
	if r.SchemaVersion != SchemaVersion || r.Outcome == "" || r.PlanDigest == "" || r.SourceCommit == "" || r.SourceRoot == "" || r.BeforeDigest == "" || r.Cleanup == "" {
		return errors.New("evidence identity and outcome are required")
	}
	if r.Commands == nil || r.Diagnostics == nil {
		return errors.New("evidence collections must be present")
	}
	return nil
}

func (r Record) Canonical() Record {
	r.Diagnostics = redactAll(r.Diagnostics)
	for i := range r.Commands {
		r.Commands[i].Stdout = Redact(r.Commands[i].Stdout)
		r.Commands[i].Stderr = Redact(r.Commands[i].Stderr)
	}
	return r
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
	digest := sha256.Sum256(data)
	return hex.EncodeToString(digest[:]), nil
}

var tokenPattern = regexp.MustCompile(`(?i)(token|password|secret|authorization)=([^\s&]+)`)

func Redact(value string) string {
	if value == "" {
		return value
	}
	value = tokenPattern.ReplaceAllString(value, "$1=[REDACTED]")
	if parsed, err := url.Parse(value); err == nil && parsed.User != nil {
		parsed.User = url.User(parsed.User.Username())
		value = parsed.String()
	}
	return value
}

func redactAll(values []string) []string {
	result := make([]string, len(values))
	for i, value := range values {
		result[i] = Redact(value)
	}
	return result
}
