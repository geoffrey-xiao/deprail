package presenter

import (
	"fmt"
	"io"
	"regexp"
	"strings"
	"unicode"
)

var credentialURLPattern = regexp.MustCompile(`(?i)(https?://)[^/@\s]+:[^/@\s]+@`)

func WriteHeader(w io.Writer, title string) error {
	_, err := fmt.Fprintf(w, "%s\n", SanitizeLabel(title))
	return err
}

func WriteSection(w io.Writer, title string) error {
	_, err := fmt.Fprintf(w, "\n%s:\n", SanitizeLabel(title))
	return err
}

func WriteStatus(w io.Writer, status string) error {
	_, err := fmt.Fprintf(w, "Status: %s\n", SanitizeLabel(status))
	return err
}

func WriteSummary(w io.Writer, label string, value any) error {
	_, err := fmt.Fprintf(w, "%s: %s\n", SanitizeLabel(label), SanitizeLabel(fmt.Sprint(value)))
	return err
}

func WriteError(w io.Writer, code, message string) error {
	if code == "" {
		_, err := fmt.Fprintf(w, "Error: %s\n", SanitizeLabel(message))
		return err
	}
	_, err := fmt.Fprintf(w, "Error [%s]: %s\n", SanitizeLabel(code), SanitizeLabel(message))
	return err
}

func SanitizeLabel(value string) string {
	value = credentialURLPattern.ReplaceAllString(value, `${1}[REDACTED]@`)
	var builder strings.Builder
	for _, character := range value {
		switch {
		case character == '\n':
			builder.WriteString(`\n`)
		case character == '\r':
			builder.WriteString(`\r`)
		case character == '\t':
			builder.WriteString(`\t`)
		case character == '\x1b':
			builder.WriteString(`\x1b`)
		case unicode.IsControl(character):
			fmt.Fprintf(&builder, `\x%02x`, character)
		default:
			builder.WriteRune(character)
		}
	}
	return builder.String()
}
