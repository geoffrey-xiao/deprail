package normalize

import "testing"

func TestNormalizeComponentBuildsCanonicalPURLs(t *testing.T) {
	tests := []struct{ name, version, ecosystem, want string }{
		{"@scope/pkg", "1.2.3", "npm", "pkg:npm/%40scope%2Fpkg@1.2.3"},
		{"Requests", "2.31.0", "pip", "pkg:pypi/requests@2.31.0"},
		{"com.example:service", "3.0", "maven", "pkg:maven/com.example/service@3.0"},
	}
	for _, test := range tests {
		got, err := NormalizeComponent(ComponentInput{Name: test.name, Version: test.version, Ecosystem: test.ecosystem})
		if err != nil || got.PURL != test.want {
			t.Errorf("%s: component = %#v, err = %v", test.ecosystem, got, err)
		}
	}
}

func TestNormalizeComponentRejectsIncompleteAndUnsupportedIdentity(t *testing.T) {
	for _, input := range []ComponentInput{{Name: "", Version: "1", Ecosystem: "npm"}, {Name: "com.example", Version: "1", Ecosystem: "maven"}, {Name: "x", Version: "1", Ecosystem: "unknown"}} {
		if _, err := NormalizeComponent(input); err == nil {
			t.Fatalf("input %#v unexpectedly normalized", input)
		}
	}
}
