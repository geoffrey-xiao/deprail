package normalize

import (
	"math"
	"sort"
)

type SeverityObservation struct {
	Source   string
	Severity string
	Score    *float64
}
type FixedVersion struct {
	Version string
	Source  string
}

type EvidenceInput struct {
	IDs      []string
	Source   string
	Severity *SeverityObservation
	Fixed    *FixedVersion
}

type Evidence struct {
	Severities    []SeverityObservation
	FixedVersions []FixedVersion
}

func NormalizeEvidence(inputs []EvidenceInput) Evidence {
	severities := make([]SeverityObservation, 0)
	fixed := make([]FixedVersion, 0)
	for _, input := range inputs {
		if input.Severity != nil && input.Severity.Source != "" && input.Severity.Severity != "" {
			severities = append(severities, *input.Severity)
		}
		if input.Fixed != nil && input.Fixed.Source != "" && input.Fixed.Version != "" {
			fixed = append(fixed, *input.Fixed)
		}
	}
	sort.Slice(severities, func(i, j int) bool {
		if severities[i].Source != severities[j].Source {
			return severities[i].Source < severities[j].Source
		}
		if severities[i].Severity != severities[j].Severity {
			return severities[i].Severity < severities[j].Severity
		}
		return compareScores(severities[i].Score, severities[j].Score) < 0
	})
	sort.Slice(fixed, func(i, j int) bool {
		if fixed[i].Source != fixed[j].Source {
			return fixed[i].Source < fixed[j].Source
		}
		return fixed[i].Version < fixed[j].Version
	})
	return Evidence{Severities: dedupeSeverity(severities), FixedVersions: dedupeFixed(fixed)}
}

func compareScores(left, right *float64) int {
	if left == nil {
		if right == nil {
			return 0
		}
		return -1
	}
	if right == nil {
		return 1
	}
	leftNaN, rightNaN := math.IsNaN(*left), math.IsNaN(*right)
	if leftNaN || rightNaN {
		if leftNaN == rightNaN {
			return 0
		}
		if leftNaN {
			return -1
		}
		return 1
	}
	if *left < *right {
		return -1
	}
	if *left > *right {
		return 1
	}
	return 0
}

func dedupeSeverity(values []SeverityObservation) []SeverityObservation {
	result := values[:0]
	for _, value := range values {
		duplicate := false
		for _, existing := range result {
			if existing.Source == value.Source && existing.Severity == value.Severity && sameScore(existing.Score, value.Score) {
				duplicate = true
				break
			}
		}
		if !duplicate {
			result = append(result, value)
		}
	}
	return result
}
func dedupeFixed(values []FixedVersion) []FixedVersion {
	result := values[:0]
	for _, value := range values {
		if len(result) == 0 || result[len(result)-1] != value {
			result = append(result, value)
		}
	}
	return result
}
func sameScore(a, b *float64) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}
