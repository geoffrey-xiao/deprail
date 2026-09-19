// Package baseline defines versioned, deterministic local scan baselines.
package baseline

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	StableKey            string   `json:"stable_key"`
	WorkspaceID          string   `json:"workspace_id,omitempty"`
	Component            string   `json:"component"`
	Version              string   `json:"version"`
	TargetID             string   `json:"target_id"`
	VulnerabilityAliases []string `json:"vulnerability_aliases,omitempty"`
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
	for i := range d.Findings {
		sort.Strings(d.Findings[i].VulnerabilityAliases)
	}
	sort.Slice(d.Findings, func(i, j int) bool { return d.Findings[i].StableKey < d.Findings[j].StableKey })
	sort.Strings(d.ArtifactDigests)
}

func validID(value string) bool {
	if value == "" {
		return false
	}
	for _, r := range value {
		if !(r == '-' || r == '_' || r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9') {
			return false
		}
	}
	return true
}

func Validate(d Document) error {
	if d.SchemaVersion != SchemaVersion || d.DocumentType != DocumentType {
		return errors.New("baseline schema identity is invalid")
	}
	if !validID(d.BaselineID) || !validID(d.SourceScanID) {
		return errors.New("baseline identity is invalid")
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
		aliases := make(map[string]struct{}, len(finding.VulnerabilityAliases))
		for _, alias := range finding.VulnerabilityAliases {
			if !validID(alias) {
				return errors.New("baseline vulnerability alias is invalid")
			}
			if _, ok := aliases[alias]; ok {
				return errors.New("baseline vulnerability aliases must be unique")
			}
			aliases[alias] = struct{}{}
		}
		if _, ok := seen[finding.StableKey]; ok {
			return errors.New("baseline finding keys must be unique")
		}
		seen[finding.StableKey] = struct{}{}
	}
	seenDigests := make(map[string]struct{}, len(d.ArtifactDigests))
	for _, digest := range d.ArtifactDigests {
		if len(digest) != sha256.Size*2 || strings.ToLower(digest) != digest {
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
	root, err := filepath.Abs(s.Root)
	if err != nil {
		return "", fmt.Errorf("resolve baseline store: %w", err)
	}
	path := filepath.Join(root, name)
	temp, err := os.CreateTemp(root, ".baseline-*")
	if err != nil {
		return "", fmt.Errorf("create baseline temporary file: %w", err)
	}
	tempName := temp.Name()
	defer os.Remove(tempName)
	if err := temp.Chmod(0o600); err != nil {
		temp.Close()
		return "", fmt.Errorf("restrict baseline temporary file: %w", err)
	}
	if _, err := temp.Write(data); err != nil {
		temp.Close()
		return "", fmt.Errorf("write baseline: %w", err)
	}
	if err := temp.Sync(); err != nil {
		temp.Close()
		return "", fmt.Errorf("sync baseline: %w", err)
	}
	if err := temp.Close(); err != nil {
		return "", fmt.Errorf("close baseline: %w", err)
	}
	if err := os.Link(tempName, path); err != nil {
		return "", fmt.Errorf("publish baseline without overwrite: %w", err)
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
	decoder := json.NewDecoder(strings.NewReader(string(data)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&document); err != nil {
		return Document{}, fmt.Errorf("decode baseline: %w", err)
	}
	var extra struct{}
	if err := decoder.Decode(&extra); err != io.EOF {
		return Document{}, errors.New("baseline contains trailing data")
	}
	if err := Validate(document); err != nil {
		return Document{}, err
	}
	document.Canonicalize()
	return document, nil
}
