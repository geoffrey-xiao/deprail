package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/geoffrey-xiao/deprail/internal/app"
	"github.com/geoffrey-xiao/deprail/internal/artifact"
	"github.com/geoffrey-xiao/deprail/internal/store/history"
)

func TestWebDetailArtifactIntegrity(t *testing.T) {
	for _, test := range []struct {
		name      string
		withRoot  bool
		integrity string
	}{
		{name: "no root", integrity: "unavailable"},
		{name: "matching artifact", withRoot: true, integrity: "verified"},
		{name: "missing artifact", withRoot: true, integrity: "missing"},
		{name: "mismatched artifact", withRoot: true, integrity: "digest_mismatch"},
	} {
		t.Run(test.name, func(t *testing.T) {
			setTestUserConfigRoot(t, canonicalTestPath(t, t.TempDir()))
			repository := canonicalTestPath(t, t.TempDir())
			artifactRoot := canonicalTestPath(t, t.TempDir())
			artifactStore := artifact.Store{Root: artifactRoot, MaxBytes: 1024}
			raw, err := artifactStore.Put([]byte("offline scanner evidence"))
			if err != nil {
				t.Fatal(err)
			}
			artifactPath := filepath.Join(artifactRoot, raw.Digest[:2], raw.Digest)
			switch test.integrity {
			case "missing":
				if err := os.Remove(artifactPath); err != nil {
					t.Fatal(err)
				}
			case "digest_mismatch":
				if err := os.WriteFile(artifactPath, []byte("corrupt evidence"), 0o600); err != nil {
					t.Fatal(err)
				}
			}

			fixture, err := os.ReadFile(filepath.Join("..", "..", "schemas", "history-v1", "examples", "nonempty.json"))
			if err != nil {
				t.Fatal(err)
			}
			var entry history.Entry
			if err := json.Unmarshal(fixture, &entry.Projection); err != nil {
				t.Fatal(err)
			}
			entry.Projection.Report.ArtifactDigests = []string{raw.Digest}
			entry.ArtifactRefs = []string{raw.Digest}
			writer, err := history.OpenWriter(context.Background(), history.Options{RepositoryRoot: repository})
			if err != nil {
				t.Fatal(err)
			}
			appendErr := writer.Append(context.Background(), entry)
			closeErr := writer.Close()
			if appendErr != nil || closeErr != nil {
				t.Fatalf("save fixture: append=%v close=%v", appendErr, closeErr)
			}

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			client := &http.Client{Timeout: 3 * time.Second}
			launch := func(target string) error {
				bootstrap, err := url.Parse(target)
				if err != nil {
					t.Fatal(err)
				}
				read := func(path string, destination any) {
					t.Helper()
					request, err := http.NewRequestWithContext(ctx, http.MethodGet, bootstrap.Scheme+"://"+bootstrap.Host+path, nil)
					if err != nil {
						t.Fatal(err)
					}
					request.Header.Set("Authorization", "Bearer "+strings.TrimPrefix(bootstrap.Fragment, "token="))
					response, err := client.Do(request)
					if err != nil {
						t.Fatalf("read saved scan: %v", err)
					}
					body, readErr := io.ReadAll(response.Body)
					_ = response.Body.Close()
					if readErr != nil || response.StatusCode != http.StatusOK {
						t.Fatalf("saved scan status=%d read=%v body=%s", response.StatusCode, readErr, body)
					}
					if err := json.Unmarshal(body, destination); err != nil {
						t.Fatal(err)
					}
				}
				path := "/api/v1/scans/" + entry.Projection.HistoryEntryID
				var detail app.HistoryDetail
				read(path, &detail)
				if detail.Summary.HistoryEntryID != entry.Projection.HistoryEntryID || detail.Report == nil || detail.Report.RepositoryState != entry.Projection.Report.RepositoryState || detail.Summary.FindingCount == nil || *detail.Summary.FindingCount != 1 {
					t.Fatalf("artifact resolution lost saved report metadata: %#v", detail)
				}
				if len(detail.ArtifactReferences) != 1 || detail.ArtifactReferences[0].Digest != raw.Digest || detail.ArtifactReferences[0].Integrity != test.integrity {
					t.Fatalf("artifact references=%#v, want digest=%s integrity=%s", detail.ArtifactReferences, raw.Digest, test.integrity)
				}
				var findings struct {
					CollectionState string            `json:"collectionState"`
					Items           []app.FindingView `json:"items"`
				}
				read(path+"/findings", &findings)
				if findings.CollectionState != "available" || len(findings.Items) != 1 || findings.Items[0].ComponentName != "lodash" || findings.Items[0].StableFindingKey != entry.Projection.Report.Findings[0].StableFindingKey {
					t.Fatalf("saved findings became unavailable: %#v", findings)
				}
				cancel()
				return nil
			}
			args := []string{"--open"}
			if test.withRoot {
				args = append(args, "--artifact-root", artifactRoot)
			}
			var stdout, stderr bytes.Buffer
			if code := runWebWithContext(ctx, args, &stdout, &stderr, nil, false, webTestAssets(), launch); code != 0 {
				t.Fatalf("web exit=%d stderr=%q", code, stderr.String())
			}
		})
	}
}
