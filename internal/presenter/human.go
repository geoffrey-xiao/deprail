package presenter

import (
	"fmt"
	"io"
)

func WriteHeader(w io.Writer, title string) error {
	_, err := fmt.Fprintf(w, "%s\n", title)
	return err
}

func WriteSection(w io.Writer, title string) error {
	_, err := fmt.Fprintf(w, "\n%s:\n", title)
	return err
}

func WriteStatus(w io.Writer, status string) error {
	_, err := fmt.Fprintf(w, "Status: %s\n", status)
	return err
}

func WriteSummary(w io.Writer, label string, value any) error {
	_, err := fmt.Fprintf(w, "%s: %v\n", label, value)
	return err
}

func WriteError(w io.Writer, code, message string) error {
	if code == "" {
		_, err := fmt.Fprintf(w, "Error: %s\n", message)
		return err
	}
	_, err := fmt.Fprintf(w, "Error [%s]: %s\n", code, message)
	return err
}
