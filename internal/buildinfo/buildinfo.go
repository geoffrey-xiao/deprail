// Package buildinfo exposes truthful release metadata injected at build time.
package buildinfo

const (
	development = "development"
	unknown     = "unknown"
)

// These variables are intentionally linker-injectable. Empty values are treated
// as development/unknown rather than as a release claim.
var (
	Version = development
	Tag     = unknown
	Commit  = unknown
)

type Identity struct {
	Version string `json:"version"`
	Tag     string `json:"tag"`
	Commit  string `json:"commit"`
}

func Current() Identity {
	return Identity{
		Version: valueOr(Version, development),
		Tag:     valueOr(Tag, unknown),
		Commit:  valueOr(Commit, unknown),
	}
}

func valueOr(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
