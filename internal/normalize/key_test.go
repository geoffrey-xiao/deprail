package normalize

import "testing"

func TestStableFindingKeyIgnoresAliasOrderAndMutableFields(t *testing.T) {
	base := FindingInput{WorkspaceID: "service", ComponentPURL: "pkg:npm/a@1", ComponentVersion: "1", VulnerabilityAliases: []string{"CVE-1", "OSV-1"}}
	permuted := FindingInput{WorkspaceID: "service", ComponentPURL: "pkg:npm/a@1", ComponentVersion: "1", VulnerabilityAliases: []string{"OSV-1", "CVE-1"}}
	if StableFindingKey(base) != StableFindingKey(permuted) {
		t.Fatal("alias order changed stable key")
	}
	if len(StableFindingKey(base)) != 64 {
		t.Fatalf("key = %q", StableFindingKey(base))
	}
}

func TestStableFindingKeySeparatesScopeAndIdentity(t *testing.T) {
	base := FindingInput{WorkspaceID: "service-a", ComponentPURL: "pkg:npm/a@1", ComponentVersion: "1", VulnerabilityAliases: []string{"OSV-1"}}
	changed := []FindingInput{
		{WorkspaceID: "service-b", ComponentPURL: base.ComponentPURL, ComponentVersion: base.ComponentVersion, VulnerabilityAliases: base.VulnerabilityAliases},
		{WorkspaceID: base.WorkspaceID, ComponentPURL: "pkg:npm/a@2", ComponentVersion: "2", VulnerabilityAliases: base.VulnerabilityAliases},
		{WorkspaceID: base.WorkspaceID, ComponentPURL: base.ComponentPURL, ComponentVersion: base.ComponentVersion, VulnerabilityAliases: []string{"OSV-2"}},
	}
	for _, input := range changed {
		if StableFindingKey(base) == StableFindingKey(input) {
			t.Fatalf("changed identity collided: %#v", input)
		}
	}
}

func TestStableFindingKeySeparatesVulnerabilityIDsWithoutAliases(t *testing.T) {
	base := FindingInput{WorkspaceID: "service", ComponentPURL: "pkg:npm/a@1", ComponentVersion: "1", VulnerabilityID: "OSV-1"}
	other := base
	other.VulnerabilityID = "OSV-2"
	if StableFindingKey(base) == StableFindingKey(other) {
		t.Fatal("vulnerability IDs with empty aliases collided")
	}
}

func TestSortFindingInputsIsDeterministic(t *testing.T) {
	inputs := []FindingInput{{WorkspaceID: "z", ComponentPURL: "pkg:npm/z@1", ComponentVersion: "1", VulnerabilityAliases: []string{"OSV-2"}}, {WorkspaceID: "a", ComponentPURL: "pkg:npm/a@1", ComponentVersion: "1", VulnerabilityAliases: []string{"OSV-1"}}}
	SortFindingInputs(inputs)
	if StableFindingKey(inputs[0]) > StableFindingKey(inputs[1]) {
		t.Fatalf("inputs = %#v", inputs)
	}
}
