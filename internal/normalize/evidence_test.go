package normalize

import "testing"

func TestNormalizeEvidencePreservesConflictingSources(t *testing.T) {
	high := 8.1
	low := 3.2
	got := NormalizeEvidence([]EvidenceInput{
		{Source: "osv", Severity: &SeverityObservation{Source: "osv", Severity: "high", Score: &high}, Fixed: &FixedVersion{Source: "osv", Version: "2.0"}},
		{Source: "github", Severity: &SeverityObservation{Source: "github", Severity: "moderate", Score: &low}, Fixed: &FixedVersion{Source: "github", Version: "1.9"}},
		{Source: "osv", Severity: &SeverityObservation{Source: "osv", Severity: "high", Score: &high}, Fixed: &FixedVersion{Source: "osv", Version: "2.0"}},
	})
	if len(got.Severities) != 2 || got.Severities[0].Source != "github" || len(got.FixedVersions) != 2 {
		t.Fatalf("evidence = %#v", got)
	}
}

func TestNormalizeVulnerabilitiesCarriesEvidence(t *testing.T) {
	got := NormalizeVulnerabilities([]VulnerabilityInput{{IDs: []string{"OSV-1"}, Source: "osv", Severity: &SeverityObservation{Source: "osv", Severity: "high"}}})
	if len(got) != 1 || len(got[0].Severities) != 1 {
		t.Fatalf("vulnerability = %#v", got)
	}
}
