package web

import (
	"testing"
	"testing/fstest"
)

func TestInventoryRejectsMissingOrMismatchedAssets(t *testing.T) {
	const index = `<html><meta name="deprail-api-version" content="v1"><script src="/assets/main-abcdef.js"></script></html>`
	cases := []struct {
		name  string
		index string
		asset bool
		extra bool
	}{
		{name: "matching", index: index, asset: true},
		{name: "missing", index: index},
		{name: "API version skew", index: `<html><meta name="deprail-api-version" content="v2"><script src="/assets/main-abcdef.js"></script></html>`, asset: true},
		{name: "stale reference", index: `<html><meta name="deprail-api-version" content="v1"><script src="/assets/other-abcdef.js"></script></html>`, asset: true},
		{name: "unreferenced asset", index: index, asset: true, extra: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			files := fstest.MapFS{"index.html": &fstest.MapFile{Data: []byte(tc.index)}}
			if tc.asset {
				files["assets/main-abcdef.js"] = &fstest.MapFile{Data: []byte("application")}
			}
			if tc.extra {
				files["assets/other-abcdef.js"] = &fstest.MapFile{Data: []byte("unlisted")}
			}
			paths, err := inventory(files, "v1")
			if tc.name == "matching" {
				if err != nil || len(paths) != 1 || paths[0] != "/assets/main-abcdef.js" {
					t.Fatalf("valid embedded build: paths=%q err=%v", paths, err)
				}
			} else if err == nil {
				t.Fatalf("%s accepted mismatched embedded build: %q", tc.name, paths)
			}
		})
	}
	if _, _, err := Assets("v2"); err == nil {
		t.Fatal("compiled v1 console accepted incompatible v2 API")
	}
}
