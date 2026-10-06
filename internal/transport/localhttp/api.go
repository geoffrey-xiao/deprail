package localhttp

import (
	"context"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/geoffrey-xiao/deprail/internal/app"
	domain "github.com/geoffrey-xiao/deprail/internal/domain/history"
)

type handler struct {
	host    string
	origin  string
	token   []byte
	tokenMu sync.RWMutex
	history *app.HistoryService
	assets  embeddedAssets
	slots   chan struct{}
}

type route struct {
	kind string
	id   string
}

type apiError struct {
	Error apiErrorDetail `json:"error"`
}

type apiErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type healthResponse struct {
	APIVersion string `json:"apiVersion"`
	Status     string `json:"status"`
}

type historyPageResponse struct {
	Items      []app.HistorySummary `json:"items"`
	NextCursor *string              `json:"nextCursor"`
}

type workspacePageResponse struct {
	CollectionState string  `json:"collectionState"`
	Reason          string  `json:"reason,omitempty"`
	Items           any     `json:"items"`
	NextCursor      *string `json:"nextCursor"`
}

type findingPageResponse struct {
	CollectionState string  `json:"collectionState"`
	Reason          string  `json:"reason,omitempty"`
	ReportStatus    *string `json:"reportStatus"`
	Items           any     `json:"items"`
	NextCursor      *string `json:"nextCursor"`
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()
	if len(r.RequestURI) > requestTargetLimit {
		writeTransportFailure(w, http.StatusRequestURITooLong)
		return
	}
	if requestHeaderBytes(r) > headerLimit {
		writeTransportFailure(w, http.StatusRequestHeaderFieldsTooLarge)
		return
	}
	if r.Host != h.host {
		h.writeError(w, http.StatusForbidden, "API_ORIGIN_REJECTED")
		return
	}
	if !canonicalRequestPath(r) {
		h.writeError(w, http.StatusBadRequest, "API_REQUEST_INVALID")
		return
	}
	if r.ContentLength != 0 || len(r.TransferEncoding) != 0 {
		h.writeError(w, http.StatusBadRequest, "API_REQUEST_INVALID")
		return
	}
	route, found := resolveRoute(r.URL.Path, h.assets)
	if !found {
		h.writeError(w, http.StatusNotFound, "API_ROUTE_NOT_FOUND")
		return
	}
	if route.kind == "unknownAPI" {
		h.writeError(w, http.StatusNotFound, "API_ROUTE_NOT_FOUND")
		return
	}
	if route.kind == "unsupportedVersion" {
		h.writeError(w, http.StatusNotFound, "API_VERSION_UNSUPPORTED")
		return
	}
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		h.writeError(w, http.StatusMethodNotAllowed, "API_METHOD_UNSUPPORTED")
		return
	}
	if route.kind == "asset" || route.kind == "document" {
		if hasQuery(r) {
			h.writeError(w, http.StatusBadRequest, "API_REQUEST_INVALID")
			return
		}
		if !h.originAllowed(r) {
			h.writeError(w, http.StatusForbidden, "API_ORIGIN_REJECTED")
			return
		}
		h.serveStatic(w, route)
		return
	}

	request, err := parseAPIRequest(r, route)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "API_REQUEST_INVALID")
		return
	}
	if !h.originAllowed(r) || !sameOriginFetch(r) {
		h.writeError(w, http.StatusForbidden, "API_ORIGIN_REJECTED")
		return
	}
	if !h.authenticated(r) {
		h.writeError(w, http.StatusUnauthorized, "API_AUTH_UNAUTHORIZED")
		return
	}
	select {
	case h.slots <- struct{}{}:
		defer func() { <-h.slots }()
	default:
		h.writeError(w, http.StatusServiceUnavailable, "API_BUSY")
		return
	}
	response, err := h.execute(ctx, request)
	if errors.Is(ctx.Err(), context.DeadlineExceeded) || errors.Is(err, context.DeadlineExceeded) {
		h.writeError(w, http.StatusGatewayTimeout, "API_TIMEOUT")
		return
	}
	if ctx.Err() != nil {
		return
	}
	if err != nil {
		status, code := mapHistoryError(err)
		h.writeError(w, status, code)
		return
	}
	h.writeJSON(w, http.StatusOK, response)
}

func requestHeaderBytes(r *http.Request) int {
	total := len("Host") + 2 + len(r.Host) + 2
	for name, values := range r.Header {
		for _, value := range values {
			total += len(name) + 2 + len(value) + 2
		}
	}
	return total
}

func canonicalRequestPath(r *http.Request) bool {
	if r.URL.RawPath != "" || !utf8.ValidString(r.URL.Path) || strings.ContainsAny(r.URL.Path, "\\\x00\r\n") {
		return false
	}
	rawPath, _, _ := strings.Cut(r.RequestURI, "?")
	if rawPath != r.URL.Path || strings.Contains(rawPath, "%") {
		return false
	}
	for _, part := range strings.Split(rawPath, "/") {
		if part == "." || part == ".." {
			return false
		}
	}
	return true
}

func hasQuery(r *http.Request) bool {
	return r.URL.RawQuery != "" || r.URL.ForceQuery
}

func resolveRoute(path string, assets embeddedAssets) (route, bool) {
	if isConsoleDocument(path) {
		return route{kind: "document"}, true
	}
	if _, ok := assets.files[path]; ok {
		return route{kind: "asset", id: path}, true
	}
	switch path {
	case "/api/v1/health":
		return route{kind: "health"}, true
	case "/api/v1/scans":
		return route{kind: "list"}, true
	case "/api/v1/scans/":
		return route{}, false
	}
	const prefix = "/api/v1/scans/"
	if strings.HasPrefix(path, prefix) {
		remaining := strings.TrimPrefix(path, prefix)
		if id, suffix, found := strings.Cut(remaining, "/"); found {
			switch suffix {
			case "workspaces":
				if _, ok := parseHistoryID(id); !ok {
					return route{kind: "invalidID"}, true
				}
				return route{kind: "workspaces", id: id}, true
			case "findings":
				if _, ok := parseHistoryID(id); !ok {
					return route{kind: "invalidID"}, true
				}
				return route{kind: "findings", id: id}, true
			default:
				return route{}, false
			}
		}
		if _, ok := parseHistoryID(remaining); ok {
			return route{kind: "detail", id: remaining}, true
		}
		return route{kind: "invalidID"}, true
	}
	if strings.HasPrefix(path, "/api/") {
		parts := strings.Split(strings.TrimPrefix(path, "/api/"), "/")
		if len(parts) > 0 && strings.HasPrefix(parts[0], "v") && parts[0] != "v1" {
			return route{kind: "unsupportedVersion"}, true
		}
		return route{kind: "unknownAPI"}, true
	}
	return route{}, false
}

type apiRequest struct {
	route     route
	pageSize  int
	history   *domain.Cursor
	workspace *app.WorkspaceCursor
	finding   *app.FindingCursor
}

func parseAPIRequest(r *http.Request, route route) (apiRequest, error) {
	request := apiRequest{route: route}
	if route.kind == "invalidID" {
		return request, errors.New("invalid history ID")
	}
	query, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return request, err
	}
	if route.kind == "health" || route.kind == "detail" {
		if len(query) != 0 || r.URL.ForceQuery {
			return request, errors.New("unexpected query")
		}
		return request, nil
	}
	if route.kind != "list" && route.kind != "workspaces" && route.kind != "findings" {
		return request, errors.New("invalid route")
	}
	for key, values := range query {
		if key != "pageSize" && key != "cursor" || len(values) != 1 {
			return request, errors.New("invalid query")
		}
	}
	request.pageSize = 25
	if values, ok := query["pageSize"]; ok {
		value := values[0]
		n, err := strconv.Atoi(value)
		if err != nil || n < 1 || n > 50 || strconv.Itoa(n) != value {
			return request, errors.New("invalid page size")
		}
		request.pageSize = n
	}
	if values, ok := query["cursor"]; ok {
		value := values[0]
		switch route.kind {
		case "list":
			request.history, err = decodeHistoryCursor(value)
		case "workspaces":
			request.workspace, err = decodeWorkspaceCursor(value, route.id)
		case "findings":
			request.finding, err = decodeFindingCursor(value, route.id)
		}
		if err != nil {
			return request, err
		}
	}
	return request, nil
}

func (h *handler) originAllowed(r *http.Request) bool {
	values := r.Header.Values("Origin")
	if len(values) == 0 {
		return true
	}
	return len(values) == 1 && values[0] == h.origin
}

func sameOriginFetch(r *http.Request) bool {
	values := r.Header.Values("Sec-Fetch-Site")
	return len(values) == 0 || len(values) == 1 && values[0] == "same-origin"
}

func (h *handler) authenticated(r *http.Request) bool {
	values := r.Header.Values("Authorization")
	if len(values) != 1 || !strings.HasPrefix(values[0], "Bearer ") {
		return false
	}
	encoded := []byte(strings.TrimPrefix(values[0], "Bearer "))
	var candidate [32]byte
	n, err := base64.RawURLEncoding.Strict().Decode(candidate[:], encoded)
	if err != nil || n != len(candidate) {
		return false
	}
	h.tokenMu.RLock()
	defer h.tokenMu.RUnlock()
	return len(h.token) == len(candidate) && subtle.ConstantTimeCompare(candidate[:], h.token) == 1
}

func (h *handler) execute(ctx context.Context, request apiRequest) (any, error) {
	if request.route.kind == "health" {
		return healthResponse{APIVersion: "v1", Status: "ready"}, nil
	}
	if h.history == nil {
		return nil, &app.HistoryError{Code: app.HistoryUnavailable}
	}
	switch request.route.kind {
	case "list":
		page, err := h.history.List(ctx, app.HistoryPageRequest{Limit: request.pageSize, Before: request.history, MaxBytes: responseLimit})
		if err != nil {
			return nil, err
		}
		var next *string
		if page.HasMore {
			if page.Next == nil {
				return nil, &app.HistoryError{Code: app.HistoryCorrupt}
			}
			encoded, err := encodeHistoryCursor(*page.Next)
			if err != nil {
				return nil, &app.HistoryError{Code: app.HistoryCorrupt}
			}
			next = &encoded
		}
		return historyPageResponse{Items: page.Items, NextCursor: next}, nil
	case "detail":
		return h.history.Detail(ctx, request.route.id)
	case "workspaces":
		page, err := h.history.Workspaces(ctx, app.WorkspacePageRequest{
			HistoryEntryID: request.route.id, Limit: request.pageSize,
			Before: request.workspace, MaxBytes: responseLimit,
		})
		if err != nil {
			return nil, err
		}
		if !page.Available {
			return workspacePageResponse{CollectionState: "unavailable", Reason: page.Reason, Items: nil}, nil
		}
		var next *string
		if page.HasMore {
			encoded, err := encodeWorkspaceCursor(request.route.id, page.Next)
			if err != nil {
				return nil, &app.HistoryError{Code: app.HistoryCorrupt}
			}
			next = &encoded
		}
		return workspacePageResponse{CollectionState: "available", Items: page.Items, NextCursor: next}, nil
	case "findings":
		page, err := h.history.Findings(ctx, app.FindingPageRequest{
			HistoryEntryID: request.route.id, Limit: request.pageSize,
			Before: request.finding, MaxBytes: responseLimit,
		})
		if err != nil {
			return nil, err
		}
		if !page.Available {
			return findingPageResponse{CollectionState: "unavailable", Reason: page.Reason, ReportStatus: nil, Items: nil}, nil
		}
		var next *string
		if page.HasMore {
			encoded, err := encodeFindingCursor(request.route.id, page.Next)
			if err != nil {
				return nil, &app.HistoryError{Code: app.HistoryCorrupt}
			}
			next = &encoded
		}
		return findingPageResponse{CollectionState: "available", ReportStatus: page.ReportStatus, Items: page.Items, NextCursor: next}, nil
	default:
		return nil, apiStatusError{status: http.StatusNotFound, code: "API_ROUTE_NOT_FOUND"}
	}
}

type apiStatusError struct {
	status int
	code   string
}

func (e apiStatusError) Error() string { return e.code }

func mapHistoryError(err error) (int, string) {
	if errors.Is(err, context.DeadlineExceeded) {
		return http.StatusGatewayTimeout, "API_TIMEOUT"
	}
	var statusErr apiStatusError
	if errors.As(err, &statusErr) {
		return statusErr.status, statusErr.code
	}
	var historyErr *app.HistoryError
	if errors.As(err, &historyErr) {
		switch historyErr.Code {
		case app.HistoryRequestInvalid:
			return http.StatusBadRequest, "API_REQUEST_INVALID"
		case app.HistoryEntryNotFound:
			return http.StatusNotFound, "HISTORY_ENTRY_NOT_FOUND"
		case app.HistoryUnavailable:
			return http.StatusServiceUnavailable, "HISTORY_UNAVAILABLE"
		case app.HistorySchemaUnsupported:
			return http.StatusServiceUnavailable, "HISTORY_SCHEMA_UNSUPPORTED"
		case app.HistoryMigrationFailed:
			return http.StatusServiceUnavailable, "HISTORY_MIGRATION_FAILED"
		case app.HistoryCorrupt:
			return http.StatusInternalServerError, "HISTORY_CORRUPT"
		case app.HistoryResponseTooLarge:
			return http.StatusInternalServerError, "API_RESPONSE_TOO_LARGE"
		}
	}
	return http.StatusInternalServerError, "API_INTERNAL_ERROR"
}

func (h *handler) writeError(w http.ResponseWriter, status int, code string) {
	message := errorMessage(code)
	if message == "" {
		status, code, message = http.StatusInternalServerError, "API_INTERNAL_ERROR", errorMessage("API_INTERNAL_ERROR")
	}
	h.writeJSON(w, status, apiError{Error: apiErrorDetail{Code: code, Message: message}})
}

func errorMessage(code string) string {
	switch code {
	case "API_AUTH_UNAUTHORIZED":
		return "Authentication is required."
	case "API_BUSY":
		return "The local API is temporarily busy."
	case "API_INTERNAL_ERROR":
		return "The local API request could not be completed."
	case "API_METHOD_UNSUPPORTED":
		return "The method is not allowed for this route."
	case "API_ORIGIN_REJECTED":
		return "The request origin is not allowed."
	case "API_REQUEST_INVALID":
		return "The request is invalid."
	case "API_RESPONSE_TOO_LARGE":
		return "The response exceeds the configured limit."
	case "API_ROUTE_NOT_FOUND":
		return "The requested API route was not found."
	case "API_TIMEOUT":
		return "The local API request timed out."
	case "API_VERSION_UNSUPPORTED":
		return "The requested API version is not supported."
	case "HISTORY_CORRUPT":
		return "The stored history entry is invalid."
	case "HISTORY_ENTRY_NOT_FOUND":
		return "The requested history entry was not found."
	case "HISTORY_MIGRATION_FAILED":
		return "Local history could not be opened safely."
	case "HISTORY_SCHEMA_UNSUPPORTED":
		return "The local history version is not supported."
	case "HISTORY_UNAVAILABLE":
		return "Local scan history is unavailable."
	default:
		return ""
	}
}

func (h *handler) writeJSON(w http.ResponseWriter, status int, value any) {
	body, err := json.Marshal(value)
	if err != nil || !utf8.Valid(body) {
		body = []byte(`{"error":{"code":"API_INTERNAL_ERROR","message":"The local API request could not be completed."}}`)
		status = http.StatusInternalServerError
	} else if len(body) > responseLimit {
		body = []byte(`{"error":{"code":"API_RESPONSE_TOO_LARGE","message":"The response exceeds the configured limit."}}`)
		status = http.StatusInternalServerError
	}
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	_, _ = w.Write(body)
}

func writeTransportFailure(w http.ResponseWriter, status int) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	if status == http.StatusRequestURITooLong {
		_, _ = w.Write([]byte("request target too long\n"))
	} else {
		_, _ = w.Write([]byte("request headers too large\n"))
	}
}

func (h *handler) revoke() {
	h.tokenMu.Lock()
	for i := range h.token {
		h.token[i] = 0
	}
	h.tokenMu.Unlock()
}

func (h *handler) bootstrapToken() string {
	h.tokenMu.RLock()
	defer h.tokenMu.RUnlock()
	return base64.RawURLEncoding.EncodeToString(h.token)
}
