package verification

import (
	"errors"
	"github.com/geoffrey-xiao/deprail/internal/remediation"
	"reflect"
	"sort"
)

type Transition string

const (
	Resolved   Transition = "resolved"
	Unchanged  Transition = "unchanged"
	Residual   Transition = "residual"
	Introduced Transition = "introduced"
	Unknown    Transition = "unknown"
)

type FindingTransition struct {
	StableKey string     `json:"stable_key"`
	State     Transition `json:"state"`
}

func Classify(before, after []remediation.ReportFinding, afterComplete bool) ([]FindingTransition, error) {
	beforeFindings, err := findingMap(before)
	if err != nil {
		return nil, err
	}
	afterFindings, err := findingMap(after)
	if err != nil {
		return nil, err
	}
	keys := make(map[string]struct{}, len(beforeFindings)+len(afterFindings))
	for key := range beforeFindings {
		keys[key] = struct{}{}
	}
	for key := range afterFindings {
		keys[key] = struct{}{}
	}
	result := make([]FindingTransition, 0, len(keys))
	for key := range keys {
		beforeFinding, existed := beforeFindings[key]
		afterFinding, remains := afterFindings[key]
		state := Unknown
		switch {
		case existed && remains && reflect.DeepEqual(beforeFinding, afterFinding):
			state = Unchanged
		case existed && remains:
			state = Residual
		case existed && afterComplete:
			state = Resolved
		case !existed:
			state = Introduced
		}
		result = append(result, FindingTransition{StableKey: key, State: state})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].StableKey < result[j].StableKey })
	return result, nil
}

func findingMap(findings []remediation.ReportFinding) (map[string]remediation.ReportFinding, error) {
	result := make(map[string]remediation.ReportFinding, len(findings))
	for _, finding := range findings {
		if finding.StableKey == "" {
			return nil, errors.New("finding stable key is required")
		}
		if _, exists := result[finding.StableKey]; exists {
			return nil, errors.New("duplicate finding stable key")
		}
		result[finding.StableKey] = finding
	}
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
