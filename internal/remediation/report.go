package remediation

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
)

const (
	ReportSchemaVersion = "v0alpha1"
	ReportDocumentType  = "scan_report"
	ReportComplete      = "complete"
	ReportPartial       = "partial"
	ReportFailed        = "failed"
	maxReportBytes      = 16 << 20
)

type Report struct {
	SchemaVersion      string             `json:"schema_version"`
	DocumentType       string             `json:"document_type"`
	ReportID           string             `json:"report_id"`
	RepositoryIdentity RepositoryIdentity `json:"repository_identity"`
	SourceScanID       string             `json:"source_scan_id"`
	ArtifactDigests    []string           `json:"artifact_digests"`
	RepositoryState    string             `json:"repository_state"`
	Status             string             `json:"status"`
	Stale              bool               `json:"stale"`
	Findings           []ReportFinding    `json:"findings"`
}

type ReportFinding struct {
	StableKey      string            `json:"stable_key"`
	Workspace      WorkspaceIdentity `json:"workspace"`
	Component      Component         `json:"component"`
	Aliases        []string          `json:"aliases,omitempty"`
	CurrentVersion string            `json:"current_version"`
	DependencyPath []string          `json:"dependency_path,omitempty"`
	FixedVersions  []string          `json:"fixed_versions,omitempty"`
	Provenance     Provenance        `json:"provenance"`
}

type ResolvedFinding struct {
	Report  Report        `json:"report"`
	Finding ReportFinding `json:"finding"`
}

type ReportErrorCode string

const (
	ReportRequired   ReportErrorCode = "PLAN_REPORT_REQUIRED"
	ReportInvalid    ReportErrorCode = "PLAN_REPORT_INVALID"
	FindingNotFound  ReportErrorCode = "FINDING_NOT_FOUND"
	FindingAmbiguous ReportErrorCode = "FINDING_AMBIGUOUS"
	ReportInputStale ReportErrorCode = "PLAN_INPUT_STALE"
)

type ReportError struct {
	Code    ReportErrorCode
	Message string
}

func (e *ReportError) Error() string { return fmt.Sprintf("%s: %s", e.Code, e.Message) }

func LoadReport(path string) (Report, error) {
	if path == "" {
		return Report{}, &ReportError{Code: ReportRequired, Message: "an explicit normalized scan report is required"}
	}
	file, err := os.Open(path)
	if err != nil {
		return Report{}, &ReportError{Code: ReportInvalid, Message: "normalized scan report is unreadable"}
	}
	defer file.Close()

	decoder := json.NewDecoder(io.LimitReader(file, maxReportBytes))
	var report Report
	if err := decoder.Decode(&report); err != nil {
		return Report{}, &ReportError{Code: ReportInvalid, Message: "normalized scan report is malformed"}
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return Report{}, &ReportError{Code: ReportInvalid, Message: "normalized scan report contains trailing data"}
	}
	if err := ValidateReport(report); err != nil {
		return Report{}, err
	}
	return report, nil
}

func ResolveFinding(report Report, key string) (ResolvedFinding, error) {
	if err := ValidateReport(report); err != nil {
		return ResolvedFinding{}, err
	}
	if key == "" {
		return ResolvedFinding{}, &ReportError{Code: FindingNotFound, Message: "finding key is required"}
	}
	matches := make([]ReportFinding, 0, 1)
	for _, finding := range report.Findings {
		if finding.StableKey == key {
			matches = append(matches, finding)
		}
	}
	switch len(matches) {
	case 0:
		return ResolvedFinding{}, &ReportError{Code: FindingNotFound, Message: "finding key was not found in the supplied report"}
	case 1:
		return ResolvedFinding{Report: report, Finding: matches[0]}, nil
	default:
		return ResolvedFinding{}, &ReportError{Code: FindingAmbiguous, Message: "finding key resolves to multiple findings in the supplied report"}
	}
}

func LoadAndResolve(path, key string) (ResolvedFinding, error) {
	report, err := LoadReport(path)
	if err != nil {
		return ResolvedFinding{}, err
	}
	return ResolveFinding(report, key)
}

func ValidateReport(report Report) error {
	if report.SchemaVersion != ReportSchemaVersion || report.DocumentType != ReportDocumentType {
		return &ReportError{Code: ReportInvalid, Message: "normalized scan report schema identity is invalid"}
	}
	if report.ReportID == "" || report.SourceScanID == "" || report.RepositoryState == "" {
		return &ReportError{Code: ReportInvalid, Message: "normalized scan report provenance is incomplete"}
	}
	if report.RepositoryIdentity.Root == "" || report.RepositoryIdentity.Repository == "" {
		return &ReportError{Code: ReportInvalid, Message: "normalized scan report repository identity is incomplete"}
	}
	if report.Stale {
		return &ReportError{Code: ReportInputStale, Message: "normalized scan report is stale"}
	}
	if report.Status != ReportComplete && report.Status != ReportPartial {
		return &ReportError{Code: ReportInvalid, Message: "normalized scan report has no usable completeness state"}
	}
	if report.ArtifactDigests == nil || report.Findings == nil {
		return &ReportError{Code: ReportInvalid, Message: "normalized scan report collections must be arrays"}
	}
	if err := validateDigests(report.ArtifactDigests); err != nil {
		return &ReportError{Code: ReportInvalid, Message: err.Error()}
	}
	seen := make(map[string]struct{}, len(report.Findings))
	for _, finding := range report.Findings {
		if finding.StableKey == "" || finding.Workspace.ID == "" || finding.Component.PURL == "" || finding.CurrentVersion == "" {
			return &ReportError{Code: ReportInvalid, Message: "normalized finding identity is incomplete"}
		}
		if _, ok := seen[finding.StableKey]; ok {
			return &ReportError{Code: FindingAmbiguous, Message: "normalized report contains duplicate finding keys"}
		}
		seen[finding.StableKey] = struct{}{}
		if err := validateRelativePath(finding.Workspace.Path); err != nil {
			return &ReportError{Code: ReportInvalid, Message: "normalized finding workspace path is invalid"}
		}
		if err := validateDigests(finding.Provenance.ArtifactDigests); err != nil {
			return &ReportError{Code: ReportInvalid, Message: "normalized finding provenance is invalid"}
		}
	}
	return nil
}

func validateDigests(digests []string) error {
	seen := make(map[string]struct{}, len(digests))
	for _, digest := range digests {
		if len(digest) != 64 || digest != toLower(digest) || !isHex(digest) {
			return errors.New("artifact digests must be unique lowercase hexadecimal SHA-256 values")
		}
		if _, ok := seen[digest]; ok {
			return errors.New("artifact digests must be unique")
		}
		seen[digest] = struct{}{}
	}
	return nil
}

func toLower(value string) string {
	result := []byte(value)
	for i, char := range result {
		if char >= 'A' && char <= 'Z' {
			result[i] = char + ('a' - 'A')
		}
	}
	return string(result)
}

func isHex(value string) bool {
	for _, char := range value {
		if !(char >= '0' && char <= '9' || char >= 'a' && char <= 'f') {
			return false
		}
	}
	return true
}

func CanonicalizeReport(report *Report) {
	sort.Strings(report.ArtifactDigests)
	sort.Slice(report.Findings, func(i, j int) bool { return report.Findings[i].StableKey < report.Findings[j].StableKey })
	for i := range report.Findings {
		sort.Strings(report.Findings[i].Aliases)
		sort.Strings(report.Findings[i].FixedVersions)
	}
}
