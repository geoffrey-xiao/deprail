// Package baseline defines versioned, deterministic local scan baselines.
package baseline

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/geoffrey-xiao/deprail/internal/discovery"
)

const (
	SchemaVersion = "v1alpha"
	DocumentType  = "baseline"
)

type Finding struct {
	StableKey string `json:"stable_key"`
	Component string `json:"component"`
	Version   string `json:"version"`
	TargetID  string `json:"target_id"`
}

type Document struct {
	SchemaVersion   string                 `json:"schema_version"`
	DocumentType    string                 `json:"document_type"`
	BaselineID      string                 `json:"baseline_id"`
	SourceScanID    string                 `json:"source_scan_id"`
	Status          discovery.Completeness `json:"status"`
	Findings        []Finding              `json:"findings"`
	ArtifactDigests []string               `json:"artifact_digests"`
}

func (d *Document) Canonicalize() {
	sort.Slice(d.Findings, func(i, j int) bool { return d.Findings[i].StableKey < d.Findings[j].StableKey })
	sort.Strings(d.ArtifactDigests)
}

func Validate(d Document) error {
	if d.SchemaVersion != SchemaVersion || d.DocumentType != DocumentType {
		return errors.New("baseline schema identity is invalid")
	}
	if strings.TrimSpace(d.BaselineID) == "" || strings.TrimSpace(d.SourceScanID) == "" {
		return errors.New("baseline identity is required")
	}
	if d.Status != discovery.Complete {
		return errors.New("only complete scans can become trusted baselines")
	}
	if d.Findings == nil || d.ArtifactDigests == nil {
		return errors.New("baseline collections must be arrays")
	}
	seen := make(map[string]struct{}, len(d.Findings))
	for _, finding := range d.Findings {
		if finding.StableKey == "" || finding.Component == "" || finding.Version == "" {
			return errors.New("baseline finding identity is incomplete")
		}
		if _, ok := seen[finding.StableKey]; ok {
			return errors.New("baseline finding keys must be unique")
		}
		seen[finding.StableKey] = struct{}{}
	}
	seenDigests := make(map[string]struct{}, len(d.ArtifactDigests))
	for _, digest := range d.ArtifactDigests {
		if len(digest) != sha256.Size*2 {
			return errors.New("baseline artifact digest is invalid")
		}
		if _, err := hex.DecodeString(digest); err != nil {
			return errors.New("baseline artifact digest is invalid")
		}
		if _, ok := seenDigests[digest]; ok {
			return errors.New("baseline artifact digests must be unique")
		}
		seenDigests[digest] = struct{}{}
	}
	return nil
}

type Store struct{ Root string }

func (s Store) Save(document Document) (string, error) {
	document.Canonicalize()
	if err := Validate(document); err != nil {
		return "", err
	}
	data, err := json.MarshalIndent(document, "", "  ")
	if err != nil {
		return "", fmt.Errorf("encode baseline: %w", err)
	}
	data = append(data, '\n')
	digest := sha256.Sum256(data)
	name := document.BaselineID + "-" + hex.EncodeToString(digest[:]) + ".json"
	if err := os.MkdirAll(s.Root, 0o700); err != nil {
		return "", fmt.Errorf("create baseline store: %w", err)
	}
	path := filepath.Join(s.Root, name)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return "", fmt.Errorf("write baseline: %w", err)
	}
	return path, nil
}

func (s Store) Load(path string) (Document, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Document{}, err
	}
	name := filepath.Base(path)
	if !strings.HasSuffix(name, ".json") {
		return Document{}, errors.New("baseline filename is invalid")
	}
	parts := strings.Split(strings.TrimSuffix(name, ".json"), "-")
	if len(parts) < 2 {
		return Document{}, errors.New("baseline filename is invalid")
	}
	expected := parts[len(parts)-1]
	digest := sha256.Sum256(data)
	if hex.EncodeToString(digest[:]) != expected {
		return Document{}, errors.New("baseline integrity check failed")
	}
	var document Document
	if err := json.Unmarshal(data, &document); err != nil {
		return Document{}, fmt.Errorf("decode baseline: %w", err)
	}
	if err := Validate(document); err != nil {
		return Document{}, err
	}
	document.Canonicalize()
	return document, nil
}
