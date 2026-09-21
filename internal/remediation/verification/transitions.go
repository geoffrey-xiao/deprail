package verification

import (
	"errors"
	"sort"

	"github.com/geoffrey-xiao/deprail/internal/remediation"
)

type Transition string

const (
	Resolved   Transition = "resolved"
	Residual   Transition = "residual"
	Introduced Transition = "introduced"
	Unknown    Transition = "unknown"
)

type FindingTransition struct {
	StableKey string     `json:"stable_key"`
	State     Transition `json:"state"`
}

func Classify(before, after []remediation.ReportFinding, afterComplete bool) ([]FindingTransition, error) {
	beforeKeys, err := findingKeys(before)
	if err != nil {
		return nil, err
	}
	afterKeys, err := findingKeys(after)
	if err != nil {
		return nil, err
	}
	keys := make(map[string]struct{}, len(beforeKeys)+len(afterKeys))
	for key := range beforeKeys {
		keys[key] = struct{}{}
	}
	for key := range afterKeys {
		keys[key] = struct{}{}
	}
	result := make([]FindingTransition, 0, len(keys))
	for key := range keys {
		_, existed := beforeKeys[key]
		_, remains := afterKeys[key]
		state := Unknown
		switch {
		case existed && remains:
			state = Residual
		case existed && afterComplete:
			state = Resolved
		case !existed && afterComplete:
			state = Introduced
		}
		result = append(result, FindingTransition{StableKey: key, State: state})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].StableKey < result[j].StableKey })
	return result, nil
}

func findingKeys(findings []remediation.ReportFinding) (map[string]struct{}, error) {
	keys := make(map[string]struct{}, len(findings))
	for _, finding := range findings {
		if finding.StableKey == "" {
			return nil, errors.New("finding stable key is required")
		}
		if _, exists := keys[finding.StableKey]; exists {
			return nil, errors.New("duplicate finding stable key")
		}
		keys[finding.StableKey] = struct{}{}
	}
	return keys, nil
}
