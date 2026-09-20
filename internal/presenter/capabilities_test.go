package presenter

import (
	"bytes"
	"io"
	"testing"
)

func TestSelectCapabilitiesUsesStderrTTYOnlyForInteractiveProgress(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	isTerminal := func(writer io.Writer) bool { return writer == stderr }
	caps := SelectCapabilities(CapabilityOptions{
		Mode:       OutputTerminal,
		Stdout:     stdout,
		Stderr:     stderr,
		Color:      true,
		CanCancel:  true,
		CanRestore: true,
		IsTerminal: isTerminal,
	})
	if !caps.Interactive || !caps.Color || !caps.Animation || !caps.RestoreOnExit {
		t.Fatalf("interactive capabilities = %#v", caps)
	}
}

func TestSelectCapabilitiesDisablesInteractiveFeaturesForJSONAndCI(t *testing.T) {
	for _, options := range []CapabilityOptions{
		{Mode: OutputJSON, Color: true, CI: false, IsTerminal: func(io.Writer) bool { return true }},
		{Mode: OutputTerminal, Color: true, CI: true, IsTerminal: func(io.Writer) bool { return true }},
	} {
		caps := SelectCapabilities(options)
		if caps.Interactive || caps.Color || caps.Animation {
			t.Fatalf("unsafe interactive capabilities = %#v", caps)
		}
	}
}

func TestSelectCapabilitiesUsesExplicitWidthAndSafeFallback(t *testing.T) {
	if got := SelectCapabilities(CapabilityOptions{Width: 120}).Width; got != 120 {
		t.Fatalf("width = %d, want 120", got)
	}
	if got := SelectCapabilities(CapabilityOptions{}).Width; got != 80 {
		t.Fatalf("fallback width = %d, want 80", got)
	}
}
