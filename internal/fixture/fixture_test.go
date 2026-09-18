package fixture_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestCatalogContainsRepresentativeOfflineFixtures(t *testing.T) {
	_, source, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("locate fixture test source")
	}
	root := filepath.Join(filepath.Dir(source), "..", "..", "testdata", "fixtures")
	expected := map[string][]string{
		"npm-basic":           {"package.json", "package-lock.json"},
		"python-requirements": {"requirements.txt"},
		"python-uv":           {"pyproject.toml", "uv.lock"},
		"java-maven":          {"pom.xml"},
		"mixed-repository":    {"frontend/package.json", "services/api/requirements.txt", "services/worker/pom.xml"},
	}
	for fixture, files := range expected {
		for _, relative := range files {
			path := filepath.Join(root, fixture, filepath.FromSlash(relative))
			info, err := os.Stat(path)
			if err != nil {
				t.Errorf("fixture %s missing %s: %v", fixture, relative, err)
				continue
			}
			if info.IsDir() {
				t.Errorf("fixture %s entry %s is a directory", fixture, relative)
			}
		}
	}
}
