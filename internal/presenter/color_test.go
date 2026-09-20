package presenter

import (
	"strings"
	"testing"
)

func TestWriteColoredHeaderStylesOnlyFixedMessage(t *testing.T) {
	var output strings.Builder
	if err := WriteColoredHeader(&output, "Scanning evil\nvalue", ColorCyan, true); err != nil {
		t.Fatal(err)
	}
	if got := output.String(); got != "\x1b[36mScanning evil\\nvalue\x1b[0m\n" {
		t.Fatalf("colored output = %q", got)
	}
}

func TestWriteColoredHeaderFallsBackToPlainOutput(t *testing.T) {
	var output strings.Builder
	if err := WriteColoredHeader(&output, "Complete", ColorGreen, false); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(output.String(), "\x1b") || output.String() != "Complete\n" {
		t.Fatalf("plain output = %q", output.String())
	}
}
