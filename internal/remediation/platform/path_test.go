package platform

import "testing"

func TestRelativePathUsesStableSlashEncoding(t *testing.T) {
	path, err := RelativePath(`workspace\nested\package.json`)
	if err != nil {
		t.Fatal(err)
	}
	if path != "workspace/nested/package.json" {
		t.Fatalf("path = %q", path)
	}
}

func TestRelativePathRejectsEscapesAndAbsolutePaths(t *testing.T) {
	for _, path := range []string{"../outside", `/outside`, `C:\outside`} {
		if _, err := RelativePath(path); err == nil {
			t.Fatalf("accepted %q", path)
		}
	}
}
