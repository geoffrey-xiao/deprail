package normalize

import (
	"reflect"
	"testing"
)

func TestNormalizeVulnerabilitiesIsPermutationInvariant(t *testing.T) {
	inputs := []VulnerabilityInput{
		{IDs: []string{"OSV-1", "CVE-1"}, Source: "osv"},
		{IDs: []string{"CVE-1", "GHSA-1"}, Source: "github"},
		{IDs: []string{"CVE-2"}, Source: "osv"},
	}
	permuted := []VulnerabilityInput{inputs[2], inputs[0], inputs[1]}
	if left, right := NormalizeVulnerabilities(inputs), NormalizeVulnerabilities(permuted); !reflect.DeepEqual(left, right) {
		t.Fatalf("permutation changed output: %#v != %#v", left, right)
	}
}

func TestNormalizeEvidenceIsIdempotent(t *testing.T) {
	high := 8.1
	input := []EvidenceInput{{Source: "osv", Severity: &SeverityObservation{Source: "osv", Severity: "high", Score: &high}, Fixed: &FixedVersion{Source: "osv", Version: "2.0"}}}
	first := NormalizeEvidence(input)
	second := NormalizeEvidence([]EvidenceInput{{Source: "osv", Severity: &first.Severities[0], Fixed: &first.FixedVersions[0]}})
	if !reflect.DeepEqual(first, second) {
		t.Fatalf("evidence normalization not idempotent: %#v != %#v", first, second)
	}
}

func FuzzStableFindingKey(f *testing.F) {
	f.Add("workspace", "pkg:npm/a@1", "1", "OSV-1")
	f.Fuzz(func(t *testing.T, workspace, purl, version, alias string) {
		input := FindingInput{WorkspaceID: workspace, ComponentPURL: purl, ComponentVersion: version, VulnerabilityAliases: []string{alias}}
		if StableFindingKey(input) != StableFindingKey(input) {
			t.Fatal("stable key changed on repeat")
		}
	})
}

func FuzzNormalizeComponentNeverPanics(f *testing.F) {
	f.Add("pkg", "1.0.0", "npm")
	f.Fuzz(func(t *testing.T, name, version, ecosystem string) {
		_, _ = NormalizeComponent(ComponentInput{Name: name, Version: version, Ecosystem: ecosystem})
	})
}
