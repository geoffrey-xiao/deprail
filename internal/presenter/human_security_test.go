package presenter

import (
	"strings"
	"testing"
)

func TestSanitizeLabelEscapesTerminalControlsAndLineBreaks(t *testing.T) {
	input := "name\nnext\r\t\x1b[31m\x00 café"
	got := SanitizeLabel(input)
	want := `name\nnext\r\t\x1b[31m\x00 café`
	if got != want {
		t.Fatalf("sanitized label = %q, want %q", got, want)
	}
	if strings.ContainsAny(got, "\r\n\x1b") {
		t.Fatalf("sanitized label contains terminal controls: %q", got)
	}
}

func TestHumanPrimitiveSanitizesRenderedLabels(t *testing.T) {
	var output strings.Builder
	if err := WriteError(&output, "SCANNER\nTIMEOUT", "bad\x1b[2J\nvalue"); err != nil {
		t.Fatal(err)
	}
	if strings.ContainsAny(output.String(), "\r\n\x1b") {
		lines := strings.Split(output.String(), "\n")
		if len(lines) != 2 {
			t.Fatalf("unsafe multiline output: %q", output.String())
		}
	}
	if !strings.Contains(output.String(), `SCANNER\nTIMEOUT`) || !strings.Contains(output.String(), `\x1b[2J`) {
		t.Fatalf("escaped output = %q", output.String())
	}
}
