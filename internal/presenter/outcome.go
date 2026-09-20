package presenter

import (
	"fmt"
	"io"
)

// WriteOutcomeSummary renders status-specific human meaning without changing results.
func WriteOutcomeSummary(w io.Writer, status string, findings int, errors []string) error {
	switch status {
	case "complete":
		if findings == 0 {
			if err := WriteHeader(w, "No known vulnerabilities found."); err != nil {
				return err
			}
		} else if err := WriteSummary(w, "Findings", findings); err != nil {
			return err
		}
	case "partial":
		if err := WriteHeader(w, "Scan incomplete: partial results."); err != nil {
			return err
		}
		if err := WriteSummary(w, "Findings", findings); err != nil {
			return err
		}
	case "failed":
		if err := WriteHeader(w, "Scan failed."); err != nil {
			return err
		}
	case "cancelled":
		if err := WriteHeader(w, "Scan cancelled."); err != nil {
			return err
		}
	default:
		if err := WriteStatus(w, status); err != nil {
			return err
		}
		if err := WriteSummary(w, "Findings", findings); err != nil {
			return err
		}
	}
	for _, diagnostic := range errors {
		if err := WriteError(w, diagnosticCode(diagnostic), diagnostic); err != nil {
			return err
		}
	}
	return nil
}

func diagnosticCode(diagnostic string) string {
	if len(diagnostic) > 0 && diagnostic[0] == '[' {
		if end := len(diagnostic) - 1; end > 1 {
			for index := 1; index < len(diagnostic); index++ {
				if diagnostic[index] == ']' {
					return diagnostic[1:index]
				}
			}
		}
	}
	return ""
}

func formatOutcomeError(code, message string) string {
	if code == "" {
		return message
	}
	return fmt.Sprintf("[%s] %s", code, message)
}
