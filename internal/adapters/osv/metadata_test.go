package osv

import (
	"context"
	"testing"

	"github.com/geoffrey-xiao/deprail/internal/adapter"
)

func TestParseVersionAndCompatibility(t *testing.T) {
	metadata, err := ParseVersion("OSV-Scanner 1.9.0\n")
	if err != nil || metadata.Version != "1.9.0" {
		t.Fatalf("metadata = %#v, err = %v", metadata, err)
	}
	if err := Compatible(context.Background(), metadata); err != nil {
		t.Fatalf("compatible: %v", err)
	}
	if err := Compatible(context.Background(), adapter.Metadata{Name: "OSV-Scanner", Version: "2.0.0"}); !adapter.IsCode(err, adapter.ErrUnsupportedTarget) {
		t.Fatalf("error = %v", err)
	}
}

func TestParseVersionRejectsMalformedOutput(t *testing.T) {
	if _, err := ParseVersion("osv scanner unavailable"); !adapter.IsCode(err, adapter.ErrInvalidOutput) {
		t.Fatalf("error = %v", err)
	}
}
