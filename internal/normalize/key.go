package normalize

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"
)

type FindingInput struct {
	WorkspaceID          string
	ComponentPURL        string
	ComponentVersion     string
	VulnerabilityAliases []string
}

// StableFindingKey excludes mutable prose, timestamps, severity labels, and evidence order.
func StableFindingKey(input FindingInput) string {
	aliases := uniqueSorted(input.VulnerabilityAliases)
	canonical := strings.Join([]string{input.WorkspaceID, input.ComponentPURL, input.ComponentVersion, strings.Join(aliases, "\x1f")}, "\x00")
	sum := sha256.Sum256([]byte(canonical))
	return hex.EncodeToString(sum[:])
}

func SortFindingInputs(inputs []FindingInput) {
	sort.SliceStable(inputs, func(i, j int) bool {
		left, right := StableFindingKey(inputs[i]), StableFindingKey(inputs[j])
		if left != right {
			return left < right
		}
		return inputs[i].WorkspaceID < inputs[j].WorkspaceID
	})
}
