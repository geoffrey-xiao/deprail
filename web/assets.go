package web

import (
	"bytes"
	"embed"
	"errors"
	"io/fs"
	"regexp"
)

//go:embed dist
var distribution embed.FS

var assetName = regexp.MustCompile(`^[A-Za-z0-9_-]+-[A-Za-z0-9_-]{6,}\.(?:css|js|svg|png|webp|woff2)$`)
var assetReference = regexp.MustCompile(`(?:src|href)="(/assets/[^"]+)"`)

// Assets returns only the resources compiled into this binary. A missing, stale,
// or unsupported build cannot be served from the working directory.
func Assets(apiVersion string) (fs.FS, []string, error) {
	files, err := fs.Sub(distribution, "dist")
	if err != nil {
		return nil, nil, err
	}
	paths, err := inventory(files, apiVersion)
	if err != nil {
		return nil, nil, err
	}
	return files, paths, nil
}

func inventory(files fs.FS, apiVersion string) ([]string, error) {
	index, err := fs.ReadFile(files, "index.html")
	if err != nil || !bytes.Contains(index, []byte(`<meta name="deprail-api-version" content="`+apiVersion+`"`)) {
		return nil, errors.New("embedded console API version is unavailable or incompatible")
	}
	root, err := fs.ReadDir(files, ".")
	if err != nil || len(root) != 2 || root[0].Name() != "assets" || !root[0].IsDir() || root[1].Name() != "index.html" || root[1].IsDir() {
		return nil, errors.New("embedded console contains unexpected files")
	}
	entries, err := fs.ReadDir(files, "assets")
	if err != nil || len(entries) == 0 {
		return nil, errors.New("embedded console has no declared assets")
	}
	paths := make([]string, 0, len(entries))
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !assetName.MatchString(name) {
			return nil, errors.New("embedded console contains an undeclared asset")
		}
		paths = append(paths, "/assets/"+name)
	}
	references := assetReference.FindAllSubmatch(index, -1)
	if len(references) != len(paths) {
		return nil, errors.New("embedded console asset inventory does not match its entry document")
	}
	for _, match := range references {
		if _, err := fs.Stat(files, string(match[1][1:])); err != nil {
			return nil, errors.New("embedded console references a missing asset")
		}
	}
	for _, path := range paths {
		found := false
		for _, match := range references {
			if string(match[1]) == path {
				if found {
					return nil, errors.New("embedded console contains a duplicate asset reference")
				}
				found = true
			}
		}
		if !found {
			return nil, errors.New("embedded console contains an unreferenced asset")
		}
	}
	return paths, nil
}
