package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/geoffrey-xiao/deprail/internal/transport/localhttp"
	"github.com/geoffrey-xiao/deprail/web"
)

func TestWebHelpAndNoninteractiveRequirements(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := runWithConsoleAssets([]string{"web", "--help"}, &stdout, &stderr, localhttp.Config{}); code != 0 {
		t.Fatalf("web help exit=%d stderr=%q", code, stderr.String())
	}
	for _, expected := range []string{"deprail web", "foreground-only", "--open", "--artifact-root"} {
		if !strings.Contains(stdout.String(), expected) {
			t.Fatalf("web help omitted %q: %s", expected, stdout.String())
		}
	}

	stdout.Reset()
	stderr.Reset()
	if code := runWebWithContext(context.Background(), nil, &stdout, &stderr, nil, false, localhttp.Config{}, nil); code != 2 {
		t.Fatalf("noninteractive web without --open exit=%d", code)
	}
	if stdout.Len() != 0 || !strings.Contains(stderr.String(), "CONFIG_INVALID") || strings.Contains(stderr.String(), "http://") {
		t.Fatalf("noninteractive rejection stdout=%q stderr=%q", stdout.String(), stderr.String())
	}

	stdout.Reset()
	stderr.Reset()
	if code := runWithConsoleAssets([]string{"web", "--unknown=secret"}, &stdout, &stderr, localhttp.Config{}); code != 2 {
		t.Fatalf("unknown web flag exit=%d", code)
	}
	if stdout.Len() != 0 || strings.Contains(stderr.String(), "secret") {
		t.Fatalf("invalid flag was reflected: stdout=%q stderr=%q", stdout.String(), stderr.String())
	}
}

func TestWebStartsReadOnlyLoopbackAndPrintsOnlyBareConsoleURL(t *testing.T) {
	configRoot := canonicalTestPath(t, t.TempDir())
	setTestUserConfigRoot(t, configRoot)
	userConfig, err := os.UserConfigDir()
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var stdout, stderr bytes.Buffer
	opened := false
	launch := func(target string) error {
		opened = true
		parsed, err := url.Parse(target)
		if err != nil || parsed.Scheme != "http" || parsed.Hostname() != "127.0.0.1" || parsed.Path != "/console/" || !strings.HasPrefix(parsed.Fragment, "token=") {
			t.Fatalf("launcher target violates bootstrap contract: %q err=%v", target, err)
		}
		token := strings.TrimPrefix(parsed.Fragment, "token=")
		request, err := http.NewRequest(http.MethodGet, parsed.Scheme+"://"+parsed.Host+"/api/v1/health", nil)
		if err != nil {
			t.Fatal(err)
		}
		request.Header.Set("Authorization", "Bearer "+token)
		response, err := (&http.Client{Timeout: 2 * time.Second}).Do(request)
		if err != nil {
			t.Fatalf("request actual local listener: %v", err)
		}
		body, readErr := io.ReadAll(response.Body)
		_ = response.Body.Close()
		var health struct {
			APIVersion string `json:"apiVersion"`
			Status     string `json:"status"`
		}
		if readErr != nil || response.StatusCode != http.StatusOK || json.Unmarshal(body, &health) != nil || health.APIVersion != "v1" || health.Status != "ready" {
			t.Fatalf("actual local API status=%d body=%q err=%v", response.StatusCode, body, readErr)
		}
		listRequest, err := http.NewRequest(http.MethodGet, parsed.Scheme+"://"+parsed.Host+"/api/v1/scans", nil)
		if err != nil {
			t.Fatal(err)
		}
		listRequest.Header.Set("Authorization", "Bearer "+token)
		listResponse, err := (&http.Client{Timeout: 2 * time.Second}).Do(listRequest)
		if err != nil {
			t.Fatalf("request actual local history API: %v", err)
		}
		listBody, listErr := io.ReadAll(listResponse.Body)
		_ = listResponse.Body.Close()
		var emptyPage struct {
			Items      []json.RawMessage `json:"items"`
			NextCursor *string           `json:"nextCursor"`
		}
		if listErr != nil || listResponse.StatusCode != http.StatusOK || json.Unmarshal(listBody, &emptyPage) != nil || emptyPage.Items == nil || len(emptyPage.Items) != 0 || emptyPage.NextCursor != nil {
			t.Fatalf("missing-store history response status=%d body=%q err=%v", listResponse.StatusCode, listBody, listErr)
		}
		document, err := (&http.Client{Timeout: 2 * time.Second}).Get(parsed.Scheme + "://" + parsed.Host + "/console/")
		if err != nil {
			t.Fatalf("request packaged console: %v", err)
		}
		html, readErr := io.ReadAll(document.Body)
		_ = document.Body.Close()
		if readErr != nil || document.StatusCode != http.StatusOK || !bytes.Contains(html, []byte(`<meta name="deprail-api-version" content="v1"`)) {
			t.Fatalf("packaged console status=%d err=%v", document.StatusCode, readErr)
		}
		cancel()
		return nil
	}
	files, paths, err := web.Assets(localhttp.APIVersion)
	if err != nil {
		t.Fatalf("compiled console assets: %v", err)
	}
	assets := localhttp.Config{AssetFS: files, IndexPath: "index.html"}
	for _, path := range paths {
		assets.AssetFiles = append(assets.AssetFiles, localhttp.AssetFile{URLPath: path, FSPath: path[1:]})
	}
	code := runWebWithContext(ctx, []string{"--open"}, &stdout, &stderr, nil, false, assets, launch)
	if code != 0 || !opened {
		t.Fatalf("web exit=%d opened=%t stdout=%q stderr=%q", code, opened, stdout.String(), stderr.String())
	}
	consoleURL := strings.TrimSpace(stdout.String())
	parsed, err := url.Parse(consoleURL)
	if err != nil || parsed.Scheme != "http" || parsed.Hostname() != "127.0.0.1" || parsed.Path != "/console/" || parsed.Fragment != "" || strings.Contains(stdout.String(), "token=") {
		t.Fatalf("CLI output is not only the bare console URL: %q err=%v", stdout.String(), err)
	}
	if stderr.Len() != 0 {
		t.Fatalf("normal noninteractive startup wrote diagnostics: %q", stderr.String())
	}
	if _, err := os.Lstat(filepath.Join(userConfig, ".deprail")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("web initialized history storage: %v", err)
	}
}

func TestWebFailsClosedWithoutAssetsOrWithInvalidArtifactRoot(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := runWebWithContext(context.Background(), []string{"--open"}, &stdout, &stderr, nil, false, localhttp.Config{}, func(string) error {
		t.Fatal("browser launch ran without validated console assets")
		return nil
	}); code != 3 {
		t.Fatalf("web without embedded assets exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	if stdout.Len() != 0 || !strings.Contains(stderr.String(), "API_LISTENER_UNAVAILABLE") {
		t.Fatalf("asset failure stdout=%q stderr=%q", stdout.String(), stderr.String())
	}

	stdout.Reset()
	stderr.Reset()
	if code := runWithConsoleAssets([]string{"web", "--open"}, &stdout, &stderr, localhttp.Config{}); code != 3 {
		t.Fatalf("injected missing assets exit=%d", code)
	}
	if stdout.Len() != 0 || !strings.Contains(stderr.String(), "API_LISTENER_UNAVAILABLE") {
		t.Fatalf("injected asset failure stdout=%q stderr=%q", stdout.String(), stderr.String())
	}

	stdout.Reset()
	stderr.Reset()
	missing := filepath.Join(t.TempDir(), "missing-artifacts")
	if code := runWebWithContext(context.Background(), []string{"--open", "--artifact-root", missing}, &stdout, &stderr, nil, false, webTestAssets(), func(string) error {
		t.Fatal("browser launch ran with an invalid artifact root")
		return nil
	}); code != 2 {
		t.Fatalf("invalid artifact root exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	if stdout.Len() != 0 || strings.Contains(stderr.String(), missing) || !strings.Contains(stderr.String(), "CONFIG_INVALID") {
		t.Fatalf("artifact root failure disclosed path: stdout=%q stderr=%q", stdout.String(), stderr.String())
	}
}

func TestWebNonTTYLauncherFailureStopsListenerWithoutLeakingURL(t *testing.T) {
	var stdout, stderr bytes.Buffer
	bootstrap := ""
	launch := func(target string) error {
		bootstrap = target
		return errors.New("private launcher detail")
	}
	code := runWebWithContext(context.Background(), []string{"--open"}, &stdout, &stderr, nil, false, webTestAssets(), launch)
	if code != 3 || bootstrap == "" {
		t.Fatalf("failed noninteractive launch exit=%d bootstrap-set=%t", code, bootstrap != "")
	}
	if strings.Contains(stdout.String(), "token=") || strings.Contains(stderr.String(), "token=") || strings.Contains(stderr.String(), "private launcher detail") || !strings.Contains(stderr.String(), "API_LISTENER_UNAVAILABLE") {
		t.Fatalf("launcher failure leaked details: stdout=%q stderr=%q", stdout.String(), stderr.String())
	}
	consoleURL := strings.TrimSpace(stdout.String())
	response, err := (&http.Client{Timeout: time.Second}).Get(consoleURL)
	if err == nil {
		_ = response.Body.Close()
		t.Fatalf("noninteractive launcher failure left listener active at %q", consoleURL)
	}
	if !strings.Contains(bootstrap, "token=") {
		t.Fatal("browser launch did not receive the process bootstrap credential")
	}
}

func TestWebInteractiveEnterRetriesCurrentSession(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var stdout, stderr bytes.Buffer
	launches := 0
	var firstTarget string
	launch := func(target string) error {
		launches++
		if launches == 1 {
			firstTarget = target
			return errors.New("simulated opener failure")
		}
		if target != firstTarget {
			t.Fatalf("reopen changed the active session URL: first=%q next=%q", firstTarget, target)
		}
		parsed, err := url.Parse(target)
		if err != nil || !strings.HasPrefix(parsed.Fragment, "token=") {
			t.Fatalf("retry omitted bootstrap credential: target=%q err=%v", target, err)
		}
		request, err := http.NewRequest(http.MethodGet, parsed.Scheme+"://"+parsed.Host+"/api/v1/health", nil)
		if err != nil {
			t.Fatal(err)
		}
		request.Header.Set("Authorization", "Bearer "+strings.TrimPrefix(parsed.Fragment, "token="))
		response, err := (&http.Client{Timeout: 2 * time.Second}).Do(request)
		if err != nil {
			t.Fatalf("retry local API request: %v", err)
		}
		_ = response.Body.Close()
		if response.StatusCode != http.StatusOK {
			t.Fatalf("retry API status=%d", response.StatusCode)
		}
		cancel()
		return nil
	}
	code := runWebWithContext(ctx, []string{"--open"}, &stdout, &stderr, strings.NewReader("\n"), true, webTestAssets(), launch)
	if code != 0 || launches != 2 || strings.Contains(stdout.String(), "token=") || !strings.Contains(stderr.String(), "Press Enter") {
		t.Fatalf("interactive retry exit=%d launches=%d stdout=%q stderr=%q", code, launches, stdout.String(), stderr.String())
	}
}

func webTestAssets() localhttp.Config {
	return localhttp.Config{
		AssetFS: fstest.MapFS{
			"index.html":               &fstest.MapFile{Data: []byte(`<html><body><main id="status">Starting local console.</main><script src="/assets/console-abcdef.js"></script></body></html>`)},
			"assets/console-abcdef.js": &fstest.MapFile{Data: []byte(`(()=>{const status=document.getElementById("status");const consume=()=>{const token=location.hash.startsWith("#token=")?location.hash.slice(7):"";history.replaceState(null,"",location.pathname+location.search);if(!token){status.textContent="Reopen this console from the active terminal.";return}fetch("/api/v1/health",{headers:{Authorization:"Bearer "+token}}).then(response=>response.ok?response.json():Promise.reject()).then(value=>{status.textContent=value.status==="ready"?"API ready":"Local API unavailable"}).catch(()=>{status.textContent="Local API unavailable"})};window.addEventListener("hashchange",consume);consume()})()`)},
		},
		IndexPath: "index.html",
		AssetFiles: []localhttp.AssetFile{
			{URLPath: "/assets/console-abcdef.js", FSPath: "assets/console-abcdef.js"},
		},
	}
}
