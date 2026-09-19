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
	BaseScanID          string   `json:"base_scan_id"`
	HeadScanID          string   `json:"head_scan_id"`
	BaseArtifactDigests []string `json:"base_artifact_digests"`
	HeadArtifactDigests []string `json:"head_artifact_digests"`
	Changes             []Change `json:"changes"`
}

func findingIdentity(f Finding) string {
	return f.WorkspaceID + "\x00" + f.Component + "\x00" + f.TargetID
}

func aliasesEqual(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		if left[i] != right[i] {
			return false
		}
	}
	return true
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
	baseByIdentity := make(map[string]Finding, len(base.Findings))
	headByIdentity := make(map[string]Finding, len(head.Findings))
	for _, finding := range base.Findings {
		baseByIdentity[findingIdentity(finding)] = finding
	}
	for _, finding := range head.Findings {
		headByIdentity[findingIdentity(finding)] = finding
	}
	matchedBase, matchedHead := make(map[string]bool), make(map[string]bool)
	changes := make([]Change, 0, len(base.Findings)+len(head.Findings))
	identities := make([]string, 0, len(baseByIdentity)+len(headByIdentity))
	seen := make(map[string]struct{})
	for identity := range baseByIdentity {
		seen[identity] = struct{}{}
		identities = append(identities, identity)
	}
	for identity := range headByIdentity {
		if _, ok := seen[identity]; !ok {
			identities = append(identities, identity)
		}
	}
	sort.Strings(identities)
	for _, identity := range identities {
		b, inBase := baseByIdentity[identity]
		h, inHead := headByIdentity[identity]
		if inBase && inHead && b.StableKey != h.StableKey {
			bb, hh := b, h
			changes = append(changes, Change{Kind: Changed, Key: identity, Base: &bb, Head: &hh})
			matchedBase[b.StableKey] = true
			matchedHead[h.StableKey] = true
		}
	}
	keys := make([]string, 0, len(baseByKey)+len(headByKey))
	seen = make(map[string]struct{})
	for key := range baseByKey {
		if !matchedBase[key] {
			seen[key] = struct{}{}
			keys = append(keys, key)
		}
	}
	for key := range headByKey {
		if !matchedHead[key] {
			if _, ok := seen[key]; !ok {
				keys = append(keys, key)
			}
		}
	}
	sort.Strings(keys)
	for _, key := range keys {
		b, inBase := baseByKey[key]
		h, inHead := headByKey[key]
		switch {
		case inBase && inHead:
			bb, hh := b, h
			kind := Unchanged
			if bb.Component != hh.Component || bb.Version != hh.Version || bb.TargetID != hh.TargetID || !aliasesEqual(bb.VulnerabilityAliases, hh.VulnerabilityAliases) {
				kind = Changed
			}
			changes = append(changes, Change{Kind: kind, Key: key, Base: &bb, Head: &hh})
		case inBase:
			bb := b
			changes = append(changes, Change{Kind: Resolved, Key: key, Base: &bb})
		default:
			hh := h
			changes = append(changes, Change{Kind: Added, Key: key, Head: &hh})
		}
	}
	sort.Slice(changes, func(i, j int) bool { return changes[i].Key < changes[j].Key })
	return Comparison{BaseScanID: base.SourceScanID, HeadScanID: head.SourceScanID, BaseArtifactDigests: base.ArtifactDigests, HeadArtifactDigests: head.ArtifactDigests, Changes: changes}, nil
}
