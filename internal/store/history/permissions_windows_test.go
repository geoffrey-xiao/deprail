//go:build windows

package history

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestWindowsHistoryPathsHaveOwnerOnlyACLs(t *testing.T) {
	config, repository := testRoots(t)
	store := openTestWriter(t, config, repository, productionLimits)
	if err := verifyPrivateACL(filepath.Join(config, ".deprail"), true); err != nil {
		t.Fatalf("history directory ACL is not owner-only: %v", err)
	}
	if err := verifyPrivateACL(store.path, false); err != nil {
		t.Fatalf("database ACL is not owner-only: %v", err)
	}
}

func TestHistoryStoreRejectsWindowsReparseDataDirectory(t *testing.T) {
	config, repository := testRoots(t)
	outside := filepath.Join(config, "outside")
	if err := os.Mkdir(outside, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(config, ".deprail")); err != nil {
		t.Skipf("directory symlink/reparse creation unavailable: %v", err)
	}
	if _, err := openAt(context.Background(), config, repository, true, productionLimits); err == nil {
		t.Fatal("writer followed a reparse-point history directory")
	}
	if _, err := os.Lstat(filepath.Join(outside, "history.sqlite3")); !os.IsNotExist(err) {
		t.Fatalf("writer created a database through the reparse point: %v", err)
	}
}
