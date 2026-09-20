package presenter

import (
	"strings"
	"testing"
)

func TestHumanPrimitivesProduceReadablePlainOutput(t *testing.T) {
	var output strings.Builder
	if err := WriteHeader(&output, "DepRail"); err != nil {
		t.Fatal(err)
	}
	if err := WriteStatus(&output, "complete"); err != nil {
		t.Fatal(err)
	}
	if err := WriteSummary(&output, "Findings", 0); err != nil {
		t.Fatal(err)
	}
	if err := WriteSection(&output, "Diagnostics"); err != nil {
		t.Fatal(err)
	}
	if err := WriteError(&output, "SCANNER_TIMEOUT", "scanner timed out"); err != nil {
		t.Fatal(err)
	}
	want := "DepRail\nStatus: complete\nFindings: 0\n\nDiagnostics:\nError [SCANNER_TIMEOUT]: scanner timed out\n"
	if output.String() != want {
		t.Fatalf("output = %q, want %q", output.String(), want)
	}
}
