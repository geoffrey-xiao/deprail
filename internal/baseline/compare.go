package baseline

import (
	"errors"
	"sort"
)

type ChangeKind string

const (
	Added     ChangeKind = "added"
	Resolved  ChangeKind = "resolved"
	Unchanged ChangeKind = "unchanged"
	Changed   ChangeKind = "changed"
)

type Change struct {
	Kind ChangeKind `json:"kind"`
	Key  string     `json:"key"`
	Base *Finding   `json:"base,omitempty"`
	Head *Finding   `json:"head,omitempty"`
}

type Comparison struct {
	BaseScanID string   `json:"base_scan_id"`
	HeadScanID string   `json:"head_scan_id"`
	Changes    []Change `json:"changes"`
}

func Compare(base, head Document) (Comparison, error) {
	if err := Validate(base); err != nil {
		return Comparison{}, errors.New("base baseline is invalid: " + err.Error())
	}
	if err := Validate(head); err != nil {
		return Comparison{}, errors.New("head baseline is invalid: " + err.Error())
	}
	base.Canonicalize()
	head.Canonicalize()
	baseByKey := make(map[string]Finding, len(base.Findings))
	headByKey := make(map[string]Finding, len(head.Findings))
	for _, finding := range base.Findings {
		baseByKey[finding.StableKey] = finding
	}
	for _, finding := range head.Findings {
		headByKey[finding.StableKey] = finding
	}
	keys := make([]string, 0, len(baseByKey)+len(headByKey))
	seen := make(map[string]struct{}, len(baseByKey)+len(headByKey))
	for key := range baseByKey {
		seen[key] = struct{}{}
		keys = append(keys, key)
	}
	for key := range headByKey {
		if _, ok := seen[key]; !ok {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)
	changes := make([]Change, 0, len(keys))
	for _, key := range keys {
		baseFinding, inBase := baseByKey[key]
		headFinding, inHead := headByKey[key]
		switch {
		case inBase && inHead:
			b, h := baseFinding, headFinding
			kind := Unchanged
			if b.Component != h.Component || b.Version != h.Version || b.TargetID != h.TargetID {
				kind = Changed
			}
			changes = append(changes, Change{Kind: kind, Key: key, Base: &b, Head: &h})
		case inBase:
			b := baseFinding
			changes = append(changes, Change{Kind: Resolved, Key: key, Base: &b})
		default:
			h := headFinding
			changes = append(changes, Change{Kind: Added, Key: key, Head: &h})
		}
	}
	return Comparison{BaseScanID: base.SourceScanID, HeadScanID: head.SourceScanID, Changes: changes}, nil
}
