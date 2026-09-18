package osv

import (
	"context"
	"regexp"
	"strconv"
	"strings"

	"github.com/geoffrey-xiao/deprail/internal/adapter"
)

var versionPattern = regexp.MustCompile(`(?i)OSV-Scanner[^0-9]*v?([0-9]+)\.([0-9]+)\.([0-9]+)`)

// ParseVersion converts the stable portion of `osv-scanner --version` output
// into scanner-neutral metadata. Extra output is ignored.
func ParseVersion(output string) (adapter.Metadata, error) {
	match := versionPattern.FindStringSubmatch(output)
	if match == nil {
		return adapter.Metadata{}, &adapter.Error{Code: adapter.ErrInvalidOutput, Message: "OSV-Scanner version output is invalid"}
	}
	version := strings.Join(match[1:], ".")
	return adapter.Metadata{Name: "OSV-Scanner", Version: version, SupportedTargets: []string{"npm", "pnpm", "yarn", "pip", "uv", "poetry", "maven", "gradle"}}, nil
}

// Compatible accepts the v1 OSV-Scanner contract used by this adapter.
func Compatible(_ context.Context, metadata adapter.Metadata) error {
	if metadata.Name != "OSV-Scanner" {
		return &adapter.Error{Code: adapter.ErrUnsupportedTarget, Message: "scanner name is unsupported"}
	}
	parts := strings.Split(metadata.Version, ".")
	if len(parts) != 3 || parts[0] != "1" {
		return &adapter.Error{Code: adapter.ErrUnsupportedTarget, Message: "OSV-Scanner version is unsupported"}
	}
	for _, part := range parts {
		if _, err := strconv.Atoi(part); err != nil {
			return &adapter.Error{Code: adapter.ErrUnsupportedTarget, Message: "OSV-Scanner version is unsupported"}
		}
	}
	return nil
}
