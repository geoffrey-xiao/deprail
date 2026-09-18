package buildinfo

import "testing"

func TestCurrentUsesTruthfulDevelopmentDefaults(t *testing.T) {
	oldVersion, oldTag, oldCommit := Version, Tag, Commit
	t.Cleanup(func() { Version, Tag, Commit = oldVersion, oldTag, oldCommit })
	Version, Tag, Commit = "", "", ""
	if got := Current(); got.Version != "development" || got.Tag != "unknown" || got.Commit != "unknown" {
		t.Fatalf("identity = %#v", got)
	}
}

func TestCurrentPreservesInjectedReleaseMetadata(t *testing.T) {
	oldVersion, oldTag, oldCommit := Version, Tag, Commit
	t.Cleanup(func() { Version, Tag, Commit = oldVersion, oldTag, oldCommit })
	Version, Tag, Commit = "v0.2.0", "v0.2.0", "abc123"
	if got := Current(); got.Version != Version || got.Tag != Tag || got.Commit != Commit {
		t.Fatalf("identity = %#v", got)
	}
}
