package presenter

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/geoffrey-xiao/deprail/internal/app"
	"github.com/geoffrey-xiao/deprail/internal/discovery"
)

func TestWriteScanJSONAndTerminal(t *testing.T) {
	report := app.ScanReport{Status: discovery.Partial, Errors: []string{"scanner failed"}}
	var jsonOutput bytes.Buffer
	if err := WriteScanJSON(&jsonOutput, report); err != nil || !strings.Contains(jsonOutput.String(), `"status": "partial"`) {
		t.Fatalf("json = %q, err = %v", jsonOutput.String(), err)
	}
	var terminal bytes.Buffer
	if err := WriteScanTerminal(&terminal, report, true); err != nil || !strings.Contains(terminal.String(), "Error: scanner failed") {
		t.Fatalf("terminal = %q, err = %v", terminal.String(), err)
	}
}

func TestWriteAtomicDoesNotOverwrite(t *testing.T) {
	path := filepath.Join(t.TempDir(), "scan.json")
	if err := WriteAtomic(path, []byte("first")); err != nil {
		t.Fatal(err)
	}
	if err := WriteAtomic(path, []byte("second")); err == nil {
		t.Fatal("expected overwrite rejection")
	}
	data, err := os.ReadFile(path)
	if err != nil || string(data) != "first" {
		t.Fatalf("data = %q, err = %v", data, err)
	}
}
