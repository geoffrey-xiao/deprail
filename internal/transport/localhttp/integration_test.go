package localhttp

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/geoffrey-xiao/deprail/internal/app"
	"github.com/geoffrey-xiao/deprail/internal/artifact"
	domain "github.com/geoffrey-xiao/deprail/internal/domain/history"
	"github.com/geoffrey-xiao/deprail/internal/normalize"
	"github.com/geoffrey-xiao/deprail/internal/store/history"
)

const (
	testEntryOne = "00000000-0000-4000-8000-000000000101"
	testEntryTwo = "00000000-0000-4000-8000-000000000102"
)

func TestHTTPHistoryRoutesPageAndPreserveStoredMeaning(t *testing.T) {
	service, databasePath := openTransportHistory(t)
	before := fileDigest(t, databasePath)
	server := startTransportServer(t, service)
	token := tokenFromBootstrap(t, server)
	client := &http.Client{Timeout: 3 * time.Second}

	status, headers, body := requestJSON(t, client, server, http.MethodGet, "/api/v1/health", token, true)
	if status != http.StatusOK || headers.Get("Cache-Control") != "no-store" || headers.Get("Access-Control-Allow-Origin") != "" {
		t.Fatalf("health status=%d headers=%v body=%s", status, headers, body)
	}
	var health healthResponse
	if err := json.Unmarshal(body, &health); err != nil || health != (healthResponse{APIVersion: "v1", Status: "ready"}) {
		t.Fatalf("health response=%s err=%v", body, err)
	}

	status, _, body = requestJSON(t, client, server, http.MethodGet, "/api/v1/scans?pageSize=1", token, true)
	if status != http.StatusOK {
		t.Fatalf("history page status=%d body=%s", status, body)
	}
	var first struct {
		Items      []app.HistorySummary `json:"items"`
		NextCursor *string              `json:"nextCursor"`
	}
	if err := json.Unmarshal(body, &first); err != nil || len(first.Items) != 1 || first.Items[0].HistoryEntryID != testEntryTwo || first.NextCursor == nil {
		t.Fatalf("first history page=%s err=%v", body, err)
	}
	query := url.Values{"pageSize": {"1"}, "cursor": {*first.NextCursor}}
	status, _, body = requestJSON(t, client, server, http.MethodGet, "/api/v1/scans?"+query.Encode(), token, true)
	var second struct {
		Items      []app.HistorySummary `json:"items"`
		NextCursor *string              `json:"nextCursor"`
	}
	if err := json.Unmarshal(body, &second); err != nil || status != http.StatusOK || len(second.Items) != 1 || second.Items[0].HistoryEntryID != testEntryOne || second.NextCursor != nil {
		t.Fatalf("second history page status=%d response=%s err=%v", status, body, err)
	}

	status, _, body = requestJSON(t, client, server, http.MethodGet, "/api/v1/scans/"+testEntryOne, token, true)
	var detail app.HistoryDetail
	if err := json.Unmarshal(body, &detail); err != nil || status != http.StatusOK || detail.Report == nil || detail.Report.SourceSchemaVersion != "v1alpha" || len(detail.ArtifactReferences) != 1 || detail.ArtifactReferences[0].Integrity != "unavailable" {
		t.Fatalf("history detail status=%d response=%s err=%v", status, body, err)
	}
	if strings.Contains(string(body), databasePath) || strings.Contains(string(body), `"path"`) {
		t.Fatalf("detail disclosed a host path: %s", body)
	}

	status, _, body = requestJSON(t, client, server, http.MethodGet, "/api/v1/scans/"+testEntryOne+"/workspaces?pageSize=1", token, true)
	var firstWorkspace struct {
		CollectionState string              `json:"collectionState"`
		Items           []app.WorkspaceView `json:"items"`
		NextCursor      *string             `json:"nextCursor"`
	}
	if err := json.Unmarshal(body, &firstWorkspace); err != nil || status != http.StatusOK || firstWorkspace.CollectionState != "available" || len(firstWorkspace.Items) != 1 || firstWorkspace.NextCursor == nil {
		t.Fatalf("first workspace page status=%d response=%s err=%v", status, body, err)
	}
	query = url.Values{"pageSize": {"1"}, "cursor": {*firstWorkspace.NextCursor}}
	status, _, body = requestJSON(t, client, server, http.MethodGet, "/api/v1/scans/"+testEntryTwo+"/workspaces?"+query.Encode(), token, true)
	var invalidCursor apiError
	if err := json.Unmarshal(body, &invalidCursor); err != nil || status != http.StatusBadRequest || invalidCursor.Error.Code != "API_REQUEST_INVALID" {
		t.Fatalf("cross-parent workspace cursor status=%d response=%s err=%v", status, body, err)
	}
	query = url.Values{"pageSize": {"1"}, "cursor": {*firstWorkspace.NextCursor}}
	status, _, body = requestJSON(t, client, server, http.MethodGet, "/api/v1/scans/"+testEntryOne+"/workspaces?"+query.Encode(), token, true)
	var secondWorkspace struct {
		CollectionState string              `json:"collectionState"`
		Items           []app.WorkspaceView `json:"items"`
		NextCursor      *string             `json:"nextCursor"`
	}
	if err := json.Unmarshal(body, &secondWorkspace); err != nil || status != http.StatusOK || secondWorkspace.CollectionState != "available" || len(secondWorkspace.Items) != 1 || secondWorkspace.NextCursor != nil || firstWorkspace.Items[0].WorkspaceID == secondWorkspace.Items[0].WorkspaceID {
		t.Fatalf("second workspace page status=%d response=%s err=%v", status, body, err)
	}

	status, _, body = requestJSON(t, client, server, http.MethodGet, "/api/v1/scans/"+testEntryOne+"/findings?pageSize=1", token, true)
	var firstFinding struct {
		CollectionState string            `json:"collectionState"`
		Items           []app.FindingView `json:"items"`
		NextCursor      *string           `json:"nextCursor"`
	}
	if err := json.Unmarshal(body, &firstFinding); err != nil || status != http.StatusOK || firstFinding.CollectionState != "available" || len(firstFinding.Items) != 1 || firstFinding.NextCursor == nil {
		t.Fatalf("first finding page status=%d response=%s err=%v", status, body, err)
	}

	query = url.Values{"pageSize": {"1"}, "cursor": {*firstFinding.NextCursor}}
	status, _, body = requestJSON(t, client, server, http.MethodGet, "/api/v1/scans/"+testEntryOne+"/workspaces?"+query.Encode(), token, true)
	if err := json.Unmarshal(body, &invalidCursor); err != nil || status != http.StatusBadRequest || invalidCursor.Error.Code != "API_REQUEST_INVALID" {
		t.Fatalf("cross-resource cursor status=%d response=%s err=%v", status, body, err)
	}
	query = url.Values{"pageSize": {"1"}, "cursor": {*firstFinding.NextCursor}}
	status, _, body = requestJSON(t, client, server, http.MethodGet, "/api/v1/scans/"+testEntryOne+"/findings?"+query.Encode(), token, true)
	var secondFinding struct {
		CollectionState string            `json:"collectionState"`
		Items           []app.FindingView `json:"items"`
		NextCursor      *string           `json:"nextCursor"`
	}
	if err := json.Unmarshal(body, &secondFinding); err != nil || status != http.StatusOK || secondFinding.CollectionState != "available" || len(secondFinding.Items) != 1 || secondFinding.NextCursor != nil || firstFinding.Items[0].StableFindingKey == secondFinding.Items[0].StableFindingKey {
		t.Fatalf("second finding page status=%d response=%s err=%v", status, body, err)
	}

	status, _, body = requestJSON(t, client, server, http.MethodGet, "/api/v1/scans/"+testEntryTwo+"/workspaces", token, true)
	var unavailable struct {
		CollectionState string `json:"collectionState"`
		Reason          string `json:"reason"`
		Items           any    `json:"items"`
		NextCursor      any    `json:"nextCursor"`
	}
	if err := json.Unmarshal(body, &unavailable); err != nil || status != http.StatusOK || unavailable.CollectionState != "unavailable" || unavailable.Reason != "source_graph_unavailable" || unavailable.Items != nil || unavailable.NextCursor != nil {
		t.Fatalf("unavailable workspace page status=%d response=%s err=%v", status, body, err)
	}
	status, _, body = requestJSON(t, client, server, http.MethodGet, "/api/v1/scans/"+testEntryTwo+"/findings", token, true)
	var unavailableFindings struct {
		CollectionState string  `json:"collectionState"`
		Reason          string  `json:"reason"`
		ReportStatus    *string `json:"reportStatus"`
		Items           any     `json:"items"`
		NextCursor      any     `json:"nextCursor"`
	}
	if err := json.Unmarshal(body, &unavailableFindings); err != nil || status != http.StatusOK || unavailableFindings.CollectionState != "unavailable" || unavailableFindings.Reason != "source_report_unavailable" || unavailableFindings.ReportStatus != nil || unavailableFindings.Items != nil || unavailableFindings.NextCursor != nil {
		t.Fatalf("unavailable finding page status=%d response=%s err=%v", status, body, err)
	}
	if after := fileDigest(t, databasePath); after != before {
		t.Fatalf("read-only API changed history database: before=%x after=%x", before, after)
	}
}

func TestHTTPArtifactReferencesAreMetadataOnlyAndRejectSymlinkedPaths(t *testing.T) {
	service, databasePath := openTransportHistory(t)
	before := fileDigest(t, databasePath)
	root := t.TempDir()
	artifacts := artifact.Store{Root: root, MaxBytes: 1024}
	stored, err := artifacts.Put([]byte("raw scanner artifact"))
	if err != nil {
		t.Fatalf("write test artifact: %v", err)
	}
	service.Artifacts = &artifacts
	server := startTransportServer(t, service)
	client := &http.Client{Timeout: 3 * time.Second}
	token := tokenFromBootstrap(t, server)
	detailPath := "/api/v1/scans/" + testEntryOne
	status, _, body := requestJSON(t, client, server, http.MethodGet, detailPath, token, true)
	var detail app.HistoryDetail
	if err := json.Unmarshal(body, &detail); err != nil || status != http.StatusOK || len(detail.ArtifactReferences) != 1 || detail.ArtifactReferences[0].Digest != stored.Digest || detail.ArtifactReferences[0].Integrity != "verified" {
		t.Fatalf("verified artifact detail status=%d response=%s err=%v", status, body, err)
	}
	if strings.Contains(string(body), "raw scanner artifact") {
		t.Fatalf("artifact response disclosed bytes: %s", body)
	}

	prefix := filepath.Join(root, stored.Digest[:2])
	if err := os.Remove(filepath.Join(prefix, stored.Digest)); err != nil {
		t.Fatalf("remove test artifact: %v", err)
	}
	if err := os.Remove(prefix); err != nil {
		t.Fatalf("remove artifact prefix directory: %v", err)
	}
	outside := t.TempDir()
	if err := os.Symlink(outside, prefix); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	status, _, body = requestJSON(t, client, server, http.MethodGet, detailPath, token, true)
	if err := json.Unmarshal(body, &detail); err != nil || status != http.StatusOK || len(detail.ArtifactReferences) != 1 || detail.ArtifactReferences[0].Integrity != "digest_mismatch" {
		t.Fatalf("symlinked artifact detail status=%d response=%s err=%v", status, body, err)
	}
	if after := fileDigest(t, databasePath); after != before {
		t.Fatalf("artifact verification changed history database: before=%x after=%x", before, after)
	}
}

func TestMissingHistoryStoreIsEmptyAndRemainsUninitialized(t *testing.T) {
	configRoot, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"HOME", "XDG_CONFIG_HOME", "APPDATA", "USERPROFILE"} {
		t.Setenv(key, configRoot)
	}
	userConfig, err := os.UserConfigDir()
	if err != nil {
		t.Fatal(err)
	}
	dataDirectory := filepath.Join(userConfig, ".deprail")
	if _, err := os.Lstat(dataDirectory); !os.IsNotExist(err) {
		t.Fatalf("test history directory exists before read-only open: %v", err)
	}
	reader, err := history.OpenReadOnly(context.Background(), history.Options{})
	if err != nil {
		t.Fatalf("open missing history read-only: %v", err)
	}
	t.Cleanup(func() { _ = reader.Close() })
	server := startTransportServer(t, &app.HistoryService{Reader: reader})
	status, _, body := requestJSON(t, &http.Client{Timeout: 3 * time.Second}, server, http.MethodGet, "/api/v1/scans", tokenFromBootstrap(t, server), true)
	var page struct {
		Items      []app.HistorySummary `json:"items"`
		NextCursor *string              `json:"nextCursor"`
	}
	if err := json.Unmarshal(body, &page); err != nil || status != http.StatusOK || len(page.Items) != 0 || page.Items == nil || page.NextCursor != nil {
		t.Fatalf("missing-store page status=%d response=%s err=%v", status, body, err)
	}

	status, _, body = requestJSON(t, &http.Client{Timeout: 3 * time.Second}, server, http.MethodGet, "/api/v1/scans/"+testEntryOne, tokenFromBootstrap(t, server), true)
	var missingEntry apiError
	if err := json.Unmarshal(body, &missingEntry); err != nil || status != http.StatusNotFound || missingEntry.Error.Code != "HISTORY_ENTRY_NOT_FOUND" {
		t.Fatalf("missing entry status=%d response=%s err=%v", status, body, err)
	}
	if _, err := os.Lstat(dataDirectory); !os.IsNotExist(err) {
		t.Fatalf("read-only history access initialized the store: %v", err)
	}
}

func TestHTTPStoreFailuresRemainTypedInsteadOfEmptySuccess(t *testing.T) {
	cases := []struct {
		storeCode history.ErrorCode
		status    int
		apiCode   string
	}{
		{storeCode: history.ErrCorrupt, status: http.StatusInternalServerError, apiCode: "HISTORY_CORRUPT"},
		{storeCode: history.ErrSchemaUnsupported, status: http.StatusServiceUnavailable, apiCode: "HISTORY_SCHEMA_UNSUPPORTED"},
		{storeCode: history.ErrMigrationFailed, status: http.StatusServiceUnavailable, apiCode: "HISTORY_MIGRATION_FAILED"},
		{storeCode: history.ErrUnavailable, status: http.StatusServiceUnavailable, apiCode: "HISTORY_UNAVAILABLE"},
	}
	for _, testCase := range cases {
		t.Run(testCase.apiCode, func(t *testing.T) {
			server := startTransportServer(t, &app.HistoryService{Reader: failingTransportHistoryReader{err: &history.Error{Code: testCase.storeCode}}})
			status, _, body := requestJSON(t, &http.Client{Timeout: 3 * time.Second}, server, http.MethodGet, "/api/v1/scans", tokenFromBootstrap(t, server), true)
			var response apiError
			if err := json.Unmarshal(body, &response); err != nil || status != testCase.status || response.Error.Code != testCase.apiCode {
				t.Fatalf("typed history failure status=%d response=%s err=%v", status, body, err)
			}
		})
	}
}
func TestHTTPStaticRoutesUseOnlyEmbeddedAllowlist(t *testing.T) {
	server := startTransportServer(t, nil)
	client := &http.Client{Timeout: 3 * time.Second}
	for _, route := range []string{"/console/", "/console/history", "/console/about", "/console/scans/" + testEntryOne} {
		response, err := client.Get(strings.TrimSuffix(server.ConsoleURL(), "/console/") + route)
		if err != nil {
			t.Fatalf("GET %s: %v", route, err)
		}
		body, readErr := io.ReadAll(response.Body)
		_ = response.Body.Close()
		if readErr != nil || response.StatusCode != http.StatusOK || string(body) != "console smoke document" {
			t.Fatalf("static route %s status=%d body=%q err=%v", route, response.StatusCode, body, readErr)
		}
		if response.Header.Get("Content-Security-Policy") != consoleCSP || response.Header.Get("Cache-Control") != "no-store" || response.Header.Get("Referrer-Policy") != "no-referrer" {
			t.Fatalf("static route %s security headers=%v", route, response.Header)
		}
	}
	response, err := client.Get(strings.TrimSuffix(server.ConsoleURL(), "/console/") + "/assets/console-abcdef.js")
	if err != nil {
		t.Fatal(err)
	}
	asset, readErr := io.ReadAll(response.Body)
	_ = response.Body.Close()
	if readErr != nil || response.StatusCode != http.StatusOK || string(asset) != "console fixture asset" || response.Header.Get("Content-Type") != "application/javascript; charset=utf-8" {
		t.Fatalf("asset response status=%d type=%q body=%q err=%v", response.StatusCode, response.Header.Get("Content-Type"), asset, readErr)
	}
	response, err = client.Get(strings.TrimSuffix(server.ConsoleURL(), "/console/") + "/assets/unlisted-abcdef.js")
	if err != nil {
		t.Fatal(err)
	}
	body, _ := io.ReadAll(response.Body)
	_ = response.Body.Close()
	if response.StatusCode != http.StatusNotFound || strings.Contains(string(body), "unlisted-abcdef") {
		t.Fatalf("unlisted static path status=%d body=%q", response.StatusCode, body)
	}
}

func TestStaticNavigationAcceptsInitialFetchMetadataNone(t *testing.T) {
	server := startTransportServer(t, nil)
	request := httptest.NewRequest(http.MethodGet, "/console/", nil)
	request.Host = server.handler.host
	request.Header.Set("Sec-Fetch-Site", "none")
	response := httptest.NewRecorder()
	server.handler.ServeHTTP(response, request)
	if response.Code != http.StatusOK || response.Body.String() != "console smoke document" {
		t.Fatalf("initial navigation status=%d body=%q", response.Code, response.Body.String())
	}
	if response.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Fatalf("static navigation enabled CORS: %v", response.Header())
	}

	nullOrigin := httptest.NewRequest(http.MethodGet, "/console/", nil)
	nullOrigin.Host = server.handler.host
	nullOrigin.Header.Set("Origin", "null")
	rejected := httptest.NewRecorder()
	server.handler.ServeHTTP(rejected, nullOrigin)
	if rejected.Code != http.StatusForbidden {
		t.Fatalf("null static origin status=%d body=%q", rejected.Code, rejected.Body.String())
	}
}

func openTransportHistory(t *testing.T) (*app.HistoryService, string) {
	t.Helper()
	configRoot, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"HOME", "XDG_CONFIG_HOME", "APPDATA", "USERPROFILE"} {
		t.Setenv(key, configRoot)
	}
	repositoryRoot, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	writer, err := history.OpenWriter(context.Background(), history.Options{RepositoryRoot: repositoryRoot})
	if err != nil {
		t.Fatalf("open test history writer: %v", err)
	}
	artifactDigest := sha256.Sum256([]byte("raw scanner artifact"))
	digest := hex.EncodeToString(artifactDigest[:])
	workspaces := []domain.Workspace{
		{WorkspaceID: "workspace-b", Path: "packages/b", Ecosystem: "npm", PackageManager: "npm", DiscoveryCompleteness: "complete"},
		{WorkspaceID: "workspace-a", Path: "packages/a", Ecosystem: "npm", PackageManager: "npm", DiscoveryCompleteness: "complete"},
	}
	findings := []normalize.HistoryFindingInput{
		{Component: "example-a", PURL: "pkg:npm/example-a@1.0.0", Version: "1.0.0", WorkspaceID: "workspace-a", WorkspacePath: "packages/a", Ecosystem: "npm", TargetID: "GHSA-aaaa-bbbb-cccc", Aliases: []string{"CVE-2026-0001"}, Severity: "high", FixedVersion: "1.0.1"},
		{Component: "example-b", PURL: "pkg:npm/example-b@2.0.0", Version: "2.0.0", WorkspaceID: "workspace-a", WorkspacePath: "packages/a", Ecosystem: "npm", TargetID: "GHSA-dddd-eeee-ffff", Aliases: []string{"CVE-2026-0002"}, Severity: "moderate", FixedVersion: "2.0.1"},
	}
	for _, fixture := range []struct {
		id          string
		recordedAt  int64
		withContext bool
	}{
		{id: testEntryOne, recordedAt: 1000, withContext: true},
		{id: testEntryTwo, recordedAt: 2000},
	} {
		input := normalize.HistoryInput{
			HistoryEntryID: fixture.id, RecordedAtUS: fixture.recordedAt, OperationOutcome: "completed",
			Diagnostics: []normalize.HistoryDiagnosticInput{},
		}
		refs := []string{}
		if fixture.withContext {
			label := "example repository"
			input.RepositoryLabel = &label
			input.Workspaces = workspaces
			input.Report = &normalize.HistoryReportInput{
				SchemaVersion: "v1alpha", ScanID: "source-scan-1", Status: "complete",
				RepositoryState: strings.Repeat("a", 64), Findings: findings, ArtifactDigests: []string{digest},
			}
			refs = []string{digest}
		}
		projection, err := normalize.ProjectHistory(input)
		if err != nil {
			t.Fatalf("project history fixture %s: %v", fixture.id, err)
		}
		if err := writer.Append(context.Background(), domain.Entry{Projection: projection, ArtifactRefs: refs}); err != nil {
			t.Fatalf("append history fixture %s: %v", fixture.id, err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close test history writer: %v", err)
	}
	reader, err := history.OpenReadOnly(context.Background(), history.Options{})
	if err != nil {
		t.Fatalf("open test history reader: %v", err)
	}
	t.Cleanup(func() { _ = reader.Close() })
	userConfig, err := os.UserConfigDir()
	if err != nil {
		t.Fatal(err)
	}
	return &app.HistoryService{Reader: reader}, filepath.Join(userConfig, ".deprail", "history.sqlite3")
}

func startTransportServer(t *testing.T, service *app.HistoryService) *Server {
	t.Helper()
	assets := fstest.MapFS{
		"index.html":               &fstest.MapFile{Data: []byte("console smoke document")},
		"assets/console-abcdef.js": &fstest.MapFile{Data: []byte("console fixture asset")},
	}
	server, err := NewServer(Config{
		History: service, AssetFS: assets, IndexPath: "index.html",
		AssetFiles: []AssetFile{{URLPath: "/assets/console-abcdef.js", FSPath: "assets/console-abcdef.js"}},
	})
	if err != nil {
		t.Fatalf("create local HTTP server: %v", err)
	}
	serveResult := make(chan error, 1)
	go func() { serveResult <- server.Serve() }()
	t.Cleanup(func() {
		_ = server.Close()
		<-serveResult
	})
	return server
}

func tokenFromBootstrap(t *testing.T, server *Server) string {
	t.Helper()
	parsed, err := url.Parse(server.BootstrapURL())
	if err != nil || !strings.HasPrefix(parsed.Fragment, "token=") {
		t.Fatalf("invalid bootstrap URL: %v", err)
	}
	return strings.TrimPrefix(parsed.Fragment, "token=")
}

func requestJSON(t *testing.T, client *http.Client, server *Server, method, requestPath, token string, sameOrigin bool) (int, http.Header, []byte) {
	t.Helper()
	origin := strings.TrimSuffix(server.ConsoleURL(), "/console/")
	request, err := http.NewRequest(method, origin+requestPath, nil)
	if err != nil {
		t.Fatalf("create request %s: %v", requestPath, err)
	}
	request.Header.Set("Authorization", "Bearer "+token)
	if sameOrigin {
		request.Header.Set("Origin", origin)
		request.Header.Set("Sec-Fetch-Site", "same-origin")
	}
	response, err := client.Do(request)
	if err != nil {
		t.Fatalf("request %s: %v", requestPath, err)
	}
	body, readErr := io.ReadAll(response.Body)
	_ = response.Body.Close()
	if readErr != nil {
		t.Fatalf("read response %s: %v", requestPath, readErr)
	}
	return response.StatusCode, response.Header, body
}

func fileDigest(t *testing.T, path string) [sha256.Size]byte {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read history DB: %v", err)
	}
	return sha256.Sum256(content)
}

type failingTransportHistoryReader struct {
	err error
}

func (r failingTransportHistoryReader) List(context.Context, domain.Page) (domain.PageResult, error) {
	return domain.PageResult{}, r.err
}

func (r failingTransportHistoryReader) Get(context.Context, string) (domain.Entry, error) {
	return domain.Entry{}, r.err
}
