package evidence

import "testing"

func testRecord() Record {
	return Record{SchemaVersion: SchemaVersion, Outcome: Failed, PlanDigest: "plan", SourceCommit: "commit", SourceRoot: "/repo", BeforeDigest: "before", Commands: []CommandRecord{{ID: "scan", Tool: "tool", ExitCode: 1, Stdout: "token=abc", Stderr: "https://user:pass@example.invalid/x"}}, Diagnostics: []string{"secret=hidden"}, Cleanup: "discarded"}
}

func TestCanonicalRedactsEvidence(t *testing.T) {
	record := testRecord().Canonical()
	if record.Commands[0].Stdout != "token=[REDACTED]" {
		t.Fatalf("stdout = %q", record.Commands[0].Stdout)
	}
	if record.Commands[0].Stderr != "https://user@example.invalid/x" {
		t.Fatalf("stderr = %q", record.Commands[0].Stderr)
	}
	if record.Diagnostics[0] != "secret=[REDACTED]" {
		t.Fatalf("diagnostics = %q", record.Diagnostics[0])
	}
}

func TestDigestRequiresValidVersionedRecord(t *testing.T) {
	left, err := Digest(testRecord())
	if err != nil {
		t.Fatal(err)
	}
	right, err := Digest(testRecord())
	if err != nil {
		t.Fatal(err)
	}
	if left != right {
		t.Fatalf("digest changed: %q != %q", left, right)
	}
	invalid := testRecord()
	invalid.SchemaVersion = "v9"
	if _, err := Digest(invalid); err == nil {
		t.Fatal("expected schema rejection")
	}
}
