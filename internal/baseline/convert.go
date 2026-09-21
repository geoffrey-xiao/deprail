package baseline

import (
	"errors"
	"fmt"

	"github.com/geoffrey-xiao/deprail/internal/discovery"
	"github.com/geoffrey-xiao/deprail/internal/normalize"
)

const (
	ScanSchemaVersion = "v1alpha"
	ScanDocumentType  = "scan"
)

// ScanInput is the stable subset of a scan document needed to create a baseline.
// Callers must preserve the scan's repository state and scanner artifact digests.
type ScanInput struct {
	SchemaVersion   string
	DocumentType    string
	ScanID          string
	RepositoryState string
	Status          discovery.Completeness
	Findings        []ScanFinding
	Errors          []string
	ArtifactDigests []string
}

type ScanFinding struct {
	Component     string
	PURL          string
	Version       string
	Aliases       []string
	Severity      string
	TargetID      string
	WorkspaceID   string
	WorkspacePath string
}

type ConversionErrorCode string

const (
	ErrScanInvalid    ConversionErrorCode = "BASELINE_INPUT_INVALID"
	ErrScanIncomplete ConversionErrorCode = "BASELINE_INPUT_INCOMPLETE"
	ErrScanStale      ConversionErrorCode = "BASELINE_INPUT_STALE"
)

type ConversionError struct {
	Code    ConversionErrorCode
	Message string
}

func (e *ConversionError) Error() string { return fmt.Sprintf("%s: %s", e.Code, e.Message) }

func ConvertScan(input ScanInput) (Document, error) {
	if input.SchemaVersion != ScanSchemaVersion || input.DocumentType != ScanDocumentType || !validID(input.ScanID) {
		return Document{}, &ConversionError{Code: ErrScanInvalid, Message: "scan schema identity is invalid"}
	}
	if input.RepositoryState == "" {
		return Document{}, &ConversionError{Code: ErrScanStale, Message: "scan repository state is missing"}
	}
	if input.Status != discovery.Complete {
		return Document{}, &ConversionError{Code: ErrScanIncomplete, Message: "only complete scans can become trusted baselines"}
	}
	if len(input.Errors) != 0 {
		return Document{}, &ConversionError{Code: ErrScanIncomplete, Message: "scan contains errors"}
	}
	if input.Findings == nil || input.ArtifactDigests == nil {
		return Document{}, &ConversionError{Code: ErrScanInvalid, Message: "scan collections must be arrays"}
	}
	if err := validateDigests(input.ArtifactDigests); err != nil {
		return Document{}, &ConversionError{Code: ErrScanInvalid, Message: err.Error()}
	}

	document := Document{
		SchemaVersion:   SchemaVersion,
		DocumentType:    DocumentType,
		BaselineID:      "base-" + input.ScanID,
		SourceScanID:    input.ScanID,
		Status:          discovery.Complete,
		Findings:        make([]Finding, 0, len(input.Findings)),
		ArtifactDigests: append(make([]string, 0, len(input.ArtifactDigests)), input.ArtifactDigests...),
	}
	for _, inputFinding := range input.Findings {
		component := inputFinding.PURL
		if component == "" {
			component = inputFinding.Component
		}
		if component == "" || inputFinding.Version == "" {
			return Document{}, &ConversionError{Code: ErrScanInvalid, Message: "scan finding identity is incomplete"}
		}
		key := normalize.StableFindingKey(normalize.FindingInput{
			WorkspaceID:          inputFinding.WorkspaceID,
			ComponentPURL:        component,
			ComponentVersion:     inputFinding.Version,
			VulnerabilityAliases: inputFinding.Aliases,
		})
		document.Findings = append(document.Findings, Finding{
			StableKey:            key,
			WorkspaceID:          inputFinding.WorkspaceID,
			Component:            component,
			Version:              inputFinding.Version,
			TargetID:             inputFinding.TargetID,
			Severity:             inputFinding.Severity,
			VulnerabilityAliases: append([]string(nil), inputFinding.Aliases...),
		})
	}
	document.Canonicalize()
	if err := Validate(document); err != nil {
		return Document{}, &ConversionError{Code: ErrScanInvalid, Message: err.Error()}
	}
	return document, nil
}

func validateDigests(digests []string) error {
	seen := make(map[string]struct{}, len(digests))
	for _, digest := range digests {
		if len(digest) != 64 || digest != toLower(digest) || !isHex(digest) {
			return errors.New("scan artifact digests must be lowercase hexadecimal SHA-256 values")
		}
		if _, ok := seen[digest]; ok {
			return errors.New("scan artifact digests must be unique")
		}
		seen[digest] = struct{}{}
	}
	return nil
}

func toLower(value string) string {
	result := make([]byte, len(value))
	for i := range value {
		if value[i] >= 'A' && value[i] <= 'Z' {
			result[i] = value[i] + ('a' - 'A')
		} else {
			result[i] = value[i]
		}
	}
	return string(result)
}

func isHex(value string) bool {
	for _, char := range value {
		if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f')) {
			return false
		}
	}
	return true
}
