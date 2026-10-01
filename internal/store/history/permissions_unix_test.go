//go:build !windows

package history

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestPOSIXHistoryModesAreOwnerPrivate(t *testing.T) {
	config, repository := testRoots(t)
	openTestWriter(t, config, repository, productionLimits)
	for path, want := range map[string]os.FileMode{
		filepath.Join(config, ".deprail"):                    0o700,
		filepath.Join(config, ".deprail", "history.sqlite3"): 0o600,
	} {
		info, err := os.Stat(path)
		if err != nil {
			t.Fatal(err)
		}
		if info.Mode().Perm() != want {
			t.Errorf("%s permissions = %04o, want %04o", filepath.Base(path), info.Mode().Perm(), want)
		}
	}
}

func TestPOSIXReadOnlyRefusesBroadDatabaseMode(t *testing.T) {
	config, repository := testRoots(t)
	store := openTestWriter(t, config, repository, productionLimits)
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(config, ".deprail", "history.sqlite3")
	if err := os.Chmod(path, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := openAt(context.Background(), config, repository, false, productionLimits); !IsCode(err, ErrUnavailable) {
		t.Fatalf("read-only broad-mode error = %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o644 {
		t.Fatalf("read-only open modified database permissions to %04o", info.Mode().Perm())
	}
}

func TestHistoryStoreRejectsSymlinkDataDirectory(t *testing.T) {
	config, repository := testRoots(t)
	outside := filepath.Join(config, "outside")
	if err := os.Mkdir(outside, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(config, ".deprail")); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	if _, err := openAt(context.Background(), config, repository, true, productionLimits); err == nil {
		t.Fatal("writer followed a symlinked history directory")
	}
	if _, err := os.Lstat(filepath.Join(outside, "history.sqlite3")); !os.IsNotExist(err) {
		t.Fatalf("writer created a database through the symlink: %v", err)
	}
}

func TestHistoryStoreRejectsSymlinkDatabase(t *testing.T) {
	config, repository := testRoots(t)
	directory := filepath.Join(config, ".deprail")
	if err := ensurePrivateDirectory(directory); err != nil {
		t.Fatal(err)
	}
	outside := filepath.Join(config, "outside.db")
	if err := os.WriteFile(outside, []byte("unchanged"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(directory, "history.sqlite3")); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	if _, err := openAt(context.Background(), config, repository, true, productionLimits); err == nil {
		t.Fatal("writer followed a symlinked database")
	}
	data, err := os.ReadFile(outside)
	if err != nil || string(data) != "unchanged" {
		t.Fatalf("symlink target changed: %q, %v", data, err)
	}
}
