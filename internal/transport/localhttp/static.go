package localhttp

import (
	"errors"
	"io/fs"
	"net/http"
	"path"
	"regexp"
	"strings"
)

type staticAsset struct {
	contentType string
	body        []byte
}

type embeddedAssets struct {
	index []byte
	files map[string]staticAsset
}

var hashedAssetPath = regexp.MustCompile(`^/assets/[A-Za-z0-9_-]+-[A-Za-z0-9_-]{6,}\.(?:css|js|svg|png|webp|woff2)$`)

func loadAssets(source fs.FS, indexPath string, declared []AssetFile) (embeddedAssets, error) {
	if source == nil || !fs.ValidPath(indexPath) || indexPath == "." || len(declared) == 0 {
		return embeddedAssets{}, errors.New("embedded console assets are unavailable")
	}
	index, err := fs.ReadFile(source, indexPath)
	if err != nil || len(index) == 0 {
		return embeddedAssets{}, errors.New("embedded console entry document is invalid")
	}
	assets := embeddedAssets{index: index, files: make(map[string]staticAsset, len(declared))}
	total := len(index)
	for _, item := range declared {
		if !hashedAssetPath.MatchString(item.URLPath) || !fs.ValidPath(item.FSPath) || item.FSPath == "." || strings.Contains(item.FSPath, "\\") {
			return embeddedAssets{}, errors.New("embedded console asset declaration is invalid")
		}
		if _, duplicate := assets.files[item.URLPath]; duplicate {
			return embeddedAssets{}, errors.New("embedded console asset declaration is duplicated")
		}
		contentType, ok := staticContentType(path.Ext(item.URLPath))
		if !ok {
			return embeddedAssets{}, errors.New("embedded console asset type is not allowed")
		}
		body, err := fs.ReadFile(source, item.FSPath)
		if err != nil || len(body) == 0 {
			return embeddedAssets{}, errors.New("embedded console asset is invalid")
		}
		total += len(body)
		if total > responseLimit {
			return embeddedAssets{}, errors.New("embedded console assets exceed the configured limit")
		}
		assets.files[item.URLPath] = staticAsset{contentType: contentType, body: body}
	}
	return assets, nil
}

func staticContentType(extension string) (string, bool) {
	switch extension {
	case ".css":
		return "text/css; charset=utf-8", true
	case ".js":
		return "application/javascript; charset=utf-8", true
	case ".svg":
		return "image/svg+xml", true
	case ".png":
		return "image/png", true
	case ".webp":
		return "image/webp", true
	case ".woff2":
		return "font/woff2", true
	default:
		return "", false
	}
}

func isConsoleDocument(path string) bool {
	if path == "/console/" || path == "/console/history" || path == "/console/about" {
		return true
	}
	prefix := "/console/scans/"
	if !strings.HasPrefix(path, prefix) || strings.Contains(path[len(prefix):], "/") {
		return false
	}
	_, ok := parseHistoryID(path[len(prefix):])
	return ok
}

func (h *handler) serveStatic(w http.ResponseWriter, route route) {
	var body []byte
	contentType := "text/html; charset=utf-8"
	if route.kind == "document" {
		body = h.assets.index
	} else {
		asset, ok := h.assets.files[route.id]
		if !ok {
			h.writeError(w, http.StatusNotFound, "API_ROUTE_NOT_FOUND")
			return
		}
		body, contentType = asset.body, asset.contentType
	}
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Referrer-Policy", "no-referrer")
	w.Header().Set("Content-Security-Policy", consoleCSP)
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}
