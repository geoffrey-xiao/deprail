package localhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/geoffrey-xiao/deprail/internal/app"
	domain "github.com/geoffrey-xiao/deprail/internal/domain/history"
)

func TestHTTPRejectsCrossOriginMalformedAndUnsupportedRequests(t *testing.T) {
	server := startTransportServer(t, nil)
	client := &http.Client{Timeout: 3 * time.Second}
	token := tokenFromBootstrap(t, server)
	origin := strings.TrimSuffix(server.ConsoleURL(), "/console/")

	assertHTTPStatus := func(method, path, auth, requestOrigin, fetchSite string, want int) *http.Response {
		t.Helper()
		request, err := http.NewRequest(method, origin+path, nil)
		if err != nil {
			t.Fatal(err)
		}
		if auth != "" {
			request.Header.Set("Authorization", auth)
		}
		if requestOrigin != "" {
			request.Header.Set("Origin", requestOrigin)
		}
		if fetchSite != "" {
			request.Header.Set("Sec-Fetch-Site", fetchSite)
		}
		response, err := client.Do(request)
		if err != nil {
			t.Fatalf("request %s %s: %v", method, path, err)
		}
		body, readErr := io.ReadAll(response.Body)
		_ = response.Body.Close()
		if readErr != nil || response.StatusCode != want {
			t.Fatalf("request %s %s status=%d want=%d body=%q err=%v", method, path, response.StatusCode, want, body, readErr)
		}
		if response.Header.Get("Access-Control-Allow-Origin") != "" {
			t.Fatalf("request %s %s enabled CORS: %v", method, path, response.Header)
		}
		return response
	}

	assertHTTPStatus(http.MethodGet, "/api/v1/health", "", "", "", http.StatusUnauthorized)
	assertHTTPStatus(http.MethodGet, "/api/v1/health", "Bearer wrong-token", "", "", http.StatusUnauthorized)
	assertHTTPStatus(http.MethodGet, "/api/v1/health", "Bearer "+token, "https://attacker.invalid", "same-origin", http.StatusForbidden)
	assertHTTPStatus(http.MethodGet, "/api/v1/health", "Bearer "+token, "null", "", http.StatusForbidden)
	assertHTTPStatus(http.MethodGet, "/api/v1/health", "Bearer "+token, "", "cross-site", http.StatusForbidden)
	assertHTTPStatus(http.MethodGet, "/api/v1/health?unexpected=1", "Bearer "+token, "", "", http.StatusBadRequest)
	assertHTTPStatus(http.MethodGet, "/api/v1/scans?pageSize=01", "Bearer "+token, "", "", http.StatusBadRequest)
	assertHTTPStatus(http.MethodPost, "/api/v1/unknown", "Bearer "+token, "", "", http.StatusNotFound)
	assertHTTPStatus(http.MethodGet, "/api/v1/scans?pageSize=1&pageSize=2", "Bearer "+token, "", "", http.StatusBadRequest)
	assertHTTPStatus(http.MethodGet, "/api/v1/scans/not-a-uuid", "Bearer "+token, "", "", http.StatusBadRequest)
	assertHTTPStatus(http.MethodGet, "/api/v1/scans/"+testEntryOne+"/findings?cursor=not-a-cursor", "Bearer "+token, "", "", http.StatusBadRequest)
	assertHTTPStatus(http.MethodGet, "/api/v2/health", "Bearer "+token, "", "", http.StatusNotFound)
	assertHTTPStatus(http.MethodGet, "/api/v1/unknown", "Bearer "+token, "", "", http.StatusNotFound)

	methodResponse := assertHTTPStatus(http.MethodPost, "/api/v1/health", "Bearer "+token, "", "", http.StatusMethodNotAllowed)
	assertHTTPStatus(http.MethodHead, "/api/v1/health", "Bearer "+token, "", "", http.StatusMethodNotAllowed)
	if methodResponse.Header.Get("Allow") != http.MethodGet {
		t.Fatalf("unsupported method Allow=%q", methodResponse.Header.Get("Allow"))
	}

	request, err := http.NewRequest(http.MethodGet, origin+"/api/v1/health", strings.NewReader("ignored"))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Authorization", "Bearer "+token)
	response, err := client.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	body, _ := io.ReadAll(response.Body)
	_ = response.Body.Close()
	if response.StatusCode != http.StatusBadRequest {
		t.Fatalf("GET body status=%d body=%q", response.StatusCode, body)
	}

	request, err = http.NewRequest(http.MethodGet, origin+"/api/v1/health", nil)
	if err != nil {
		t.Fatal(err)
	}
	request.Host = "localhost" + strings.TrimPrefix(origin, "http://127.0.0.1")
	request.Header.Set("Authorization", "Bearer "+token)
	response, err = client.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	body, _ = io.ReadAll(response.Body)
	_ = response.Body.Close()
	if response.StatusCode != http.StatusForbidden {
		t.Fatalf("host mismatch status=%d body=%q", response.StatusCode, body)
	}
}

func TestHTTPBoundsTargetsHeadersResponsesAndExpiredContexts(t *testing.T) {
	server := startTransportServer(t, nil)
	target := httptest.NewRequest(http.MethodGet, "/api/v1/health?x="+strings.Repeat("a", requestTargetLimit), nil)
	target.Host = server.handler.host
	targetRecorder := httptest.NewRecorder()
	server.handler.ServeHTTP(targetRecorder, target)
	if targetRecorder.Code != http.StatusRequestURITooLong {
		t.Fatalf("oversized target status=%d body=%q", targetRecorder.Code, targetRecorder.Body.String())
	}

	header := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	header.Host = server.handler.host
	header.Header.Set("X-Large", strings.Repeat("x", headerLimit))
	headerRecorder := httptest.NewRecorder()
	server.handler.ServeHTTP(headerRecorder, header)
	if headerRecorder.Code != http.StatusRequestHeaderFieldsTooLarge {
		t.Fatalf("oversized header status=%d body=%q", headerRecorder.Code, headerRecorder.Body.String())
	}

	type responsePayload struct {
		Large string `json:"large"`
	}
	remaining := responseLimit - len(`{"large":"`) - len(`"}`)
	payload := responsePayload{Large: strings.Repeat("é", remaining/2)}
	if remaining%2 != 0 {
		payload.Large += "a"
	}
	encoded, err := json.Marshal(payload)
	if err != nil || len(encoded) != responseLimit {
		t.Fatalf("exact-boundary payload bytes=%d want=%d err=%v", len(encoded), responseLimit, err)
	}
	exactRecorder := httptest.NewRecorder()
	server.handler.writeJSON(exactRecorder, http.StatusOK, payload)
	if exactRecorder.Code != http.StatusOK || exactRecorder.Body.Len() != responseLimit {
		t.Fatalf("exact-boundary response status=%d bytes=%d", exactRecorder.Code, exactRecorder.Body.Len())
	}
	payload.Large += "a"
	overRecorder := httptest.NewRecorder()
	server.handler.writeJSON(overRecorder, http.StatusOK, payload)
	if overRecorder.Code != http.StatusInternalServerError || overRecorder.Body.Len() > responseLimit || !strings.Contains(overRecorder.Body.String(), "API_RESPONSE_TOO_LARGE") {
		t.Fatalf("one-byte-over response status=%d bytes=%d body=%q", overRecorder.Code, overRecorder.Body.Len(), overRecorder.Body.String())
	}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer cancel()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil).WithContext(ctx)
	request.Host = server.handler.host
	request.Header.Set("Authorization", "Bearer "+tokenFromBootstrap(t, server))
	recorder := httptest.NewRecorder()
	server.handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusGatewayTimeout || !strings.Contains(recorder.Body.String(), "API_TIMEOUT") {
		t.Fatalf("expired request status=%d body=%q", recorder.Code, recorder.Body.String())
	}
}

type blockedHistoryReader struct {
	started chan struct{}
	release chan struct{}
}

func (r *blockedHistoryReader) List(ctx context.Context, page domain.Page) (domain.PageResult, error) {
	r.started <- struct{}{}
	select {
	case <-r.release:
		return domain.PageResult{Entries: []domain.Summary{}}, nil
	case <-ctx.Done():
		return domain.PageResult{}, ctx.Err()
	}
}

func (r *blockedHistoryReader) Get(context.Context, string) (domain.Entry, error) {
	return domain.Entry{}, errors.New("unexpected detail request")
}

func TestHTTPAdmissionIsBoundedAtEightRequests(t *testing.T) {
	reader := &blockedHistoryReader{started: make(chan struct{}, requestSlots+1), release: make(chan struct{})}
	service := &app.HistoryService{Reader: reader}
	server := startTransportServer(t, service)
	token := tokenFromBootstrap(t, server)
	var wait sync.WaitGroup
	for range requestSlots {
		wait.Add(1)
		go func() {
			defer wait.Done()
			request := httptest.NewRequest(http.MethodGet, "/api/v1/scans", nil)
			request.Host = server.handler.host
			request.Header.Set("Authorization", "Bearer "+token)
			recorder := httptest.NewRecorder()
			server.handler.ServeHTTP(recorder, request)
			if recorder.Code != http.StatusOK {
				t.Errorf("admitted request status=%d body=%q", recorder.Code, recorder.Body.String())
			}
		}()
	}
	for range requestSlots {
		select {
		case <-reader.started:
		case <-time.After(2 * time.Second):
			close(reader.release)
			wait.Wait()
			t.Fatal("request did not reach the history reader")
		}
	}
	request := httptest.NewRequest(http.MethodGet, "/api/v1/scans", nil)
	request.Host = server.handler.host
	request.Header.Set("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()
	server.handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusServiceUnavailable || !strings.Contains(recorder.Body.String(), "API_BUSY") {
		close(reader.release)
		wait.Wait()
		t.Fatalf("ninth request status=%d body=%q", recorder.Code, recorder.Body.String())
	}
	close(reader.release)
	wait.Wait()
}

func TestHTTPRejectsNonemptyGETBodyBeforeCallingHistory(t *testing.T) {
	server := startTransportServer(t, nil)
	request := httptest.NewRequest(http.MethodGet, "/api/v1/health", bytes.NewBufferString("untrusted"))
	request.Host = server.handler.host
	request.Header.Set("Authorization", "Bearer "+tokenFromBootstrap(t, server))
	recorder := httptest.NewRecorder()
	server.handler.ServeHTTP(recorder, request)
	var response apiError
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil || recorder.Code != http.StatusBadRequest || response.Error.Code != "API_REQUEST_INVALID" {
		t.Fatalf("GET body response status=%d body=%q err=%v", recorder.Code, recorder.Body.String(), err)
	}
}

func TestHTTPShutdownRevokesBearerCredential(t *testing.T) {
	server := startTransportServer(t, nil)
	request := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	request.Host = server.handler.host
	request.Header.Set("Authorization", "Bearer "+tokenFromBootstrap(t, server))
	if !server.handler.authenticated(request) {
		t.Fatal("bootstrap credential did not authenticate before shutdown")
	}
	if err := server.Shutdown(); err != nil {
		t.Fatalf("shutdown local server: %v", err)
	}
	if server.handler.authenticated(request) {
		t.Fatal("shutdown retained the process bearer credential")
	}
}

func TestHTTPShutdownCancelsInFlightQueryAndClosesPromptly(t *testing.T) {
	reader := &blockedHistoryReader{started: make(chan struct{}, 1), release: make(chan struct{})}
	server := startTransportServer(t, &app.HistoryService{Reader: reader})
	request, err := http.NewRequest(http.MethodGet, strings.TrimSuffix(server.ConsoleURL(), "/console/")+"/api/v1/scans", nil)
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Authorization", "Bearer "+tokenFromBootstrap(t, server))
	requestDone := make(chan error, 1)
	go func() {
		response, err := (&http.Client{Timeout: 3 * time.Second}).Do(request)
		if response != nil {
			_ = response.Body.Close()
		}
		requestDone <- err
	}()
	select {
	case <-reader.started:
	case <-time.After(2 * time.Second):
		t.Fatal("query did not reach the history reader")
	}
	shutdownDone := make(chan error, 1)
	go func() { shutdownDone <- server.Shutdown() }()
	select {
	case err := <-shutdownDone:
		if err != nil {
			t.Fatalf("shutdown local server: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("shutdown did not cancel and drain the in-flight query")
	}
	select {
	case <-requestDone:
	case <-time.After(2 * time.Second):
		t.Fatal("in-flight client request remained blocked after shutdown")
	}
}
