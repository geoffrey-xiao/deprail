package evidence

import (
	"encoding/json"
	"testing"
)

func testRecord() Record {
	return Record{SchemaVersion: SchemaVersion, Outcome: Failed, PlanDigest: "plan", SourceCommit: "commit", SourceRoot: "/repo", WorkspaceID: "workspace", WorkspacePath: "/tmp/workspace", AuthorizedPaths: []string{"package.json"}, Operations: []Operation{{ID: "op", Path: "package.json", Effect: "update"}}, BeforeDigest: "before", FindingBeforeDigest: "finding-before", VerificationStatus: "failed", RescanStatus: "not_run", Commands: []CommandRecord{{ID: "scan", Tool: "tool", ExitCode: 1, Stdout: "token=abc Authorization: Bearer xyz", Stderr: "request failed for https://user:pass@example.invalid/x"}}, Diagnostics: []string{"secret=hidden"}, Cleanup: "discarded", ArtifactDigests: []string{}}
}

func TestCanonicalRedactsEvidence(t *testing.T) {
	record := testRecord().Canonical()
	if record.Commands[0].Stdout != "token=[REDACTED] Authorization=[REDACTED] [REDACTED]" {
		t.Fatalf("stdout = %q", record.Commands[0].Stdout)
	}
	if record.Commands[0].Stderr != "request failed for https://[REDACTED]@example.invalid/x" {
		t.Fatalf("stderr = %q", record.Commands[0].Stderr)
	}
	if record.Diagnostics[0] != "secret=[REDACTED]" {
		t.Fatalf("diagnostics = %q", record.Diagnostics[0])
	}
}
func TestMarshalAlwaysRedactsAndDigestIsStable(t *testing.T) {
	data, err := json.Marshal(testRecord())
	if err != nil {
		t.Fatal(err)
	}
	if string(data) == "" || contains(string(data), "abc") || contains(string(data), "pass") {
		t.Fatalf("secret persisted: %s", data)
	}
	left, err := Digest(testRecord())
	if err != nil {
		t.Fatal(err)
	}
	right, err := Digest(func() Record { r := testRecord(); r.Diagnostics = []string{"other", "secret=hidden"}; return r }())
	if err != nil {
		t.Fatal(err)
	}
	if left == right {
		t.Fatal("expected diagnostic content to affect digest")
	}
	invalid := testRecord()
	invalid.SchemaVersion = "v9"
	if _, err := Digest(invalid); err == nil {
		t.Fatal("expected schema rejection")
	}
}
func TestValidateRejectsUnknownOutcome(t *testing.T) {
	r := testRecord()
	r.Outcome = "success"
	if err := r.Validate(); err == nil {
		t.Fatal("expected outcome rejection")
	}
}
func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
