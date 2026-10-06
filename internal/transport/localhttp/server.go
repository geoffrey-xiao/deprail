package localhttp

import (
	"context"
	"crypto/rand"
	"errors"
	"io"
	"io/fs"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/geoffrey-xiao/deprail/internal/app"
)

const (
	requestTargetLimit = 2048
	headerLimit        = 8192
	requestTimeout     = 10 * time.Second
	shutdownTimeout    = 10 * time.Second
	responseLimit      = 1 << 20
	requestSlots       = 8
	consoleCSP         = "default-src 'none'; script-src 'self'; style-src 'self'; img-src 'self' data:; font-src 'self'; connect-src 'self'; base-uri 'none'; form-action 'none'; frame-ancestors 'none'; object-src 'none'"
)

type AssetFile struct {
	URLPath string
	FSPath  string
}

type Config struct {
	History    *app.HistoryService
	AssetFS    fs.FS
	IndexPath  string
	AssetFiles []AssetFile
}

type Server struct {
	listener net.Listener
	http     *http.Server
	handler  *handler
	cancel   context.CancelFunc
}

func NewServer(config Config) (*Server, error) {
	assets, err := loadAssets(config.AssetFS, config.IndexPath, config.AssetFiles)
	if err != nil {
		return nil, err
	}
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}
	host := listener.Addr().String()
	if !strings.HasPrefix(host, "127.0.0.1:") {
		_ = listener.Close()
		return nil, errors.New("listener address is not loopback")
	}
	var token [32]byte
	if _, err := rand.Read(token[:]); err != nil {
		_ = listener.Close()
		return nil, errors.New("unable to create local API credential")
	}
	base, cancel := context.WithCancel(context.Background())
	h := &handler{
		host: host, origin: "http://" + host, token: token[:], history: config.History,
		assets: assets, slots: make(chan struct{}, requestSlots),
	}
	httpServer := &http.Server{
		Handler: h, MaxHeaderBytes: headerLimit, ReadHeaderTimeout: requestTimeout,
		WriteTimeout: requestTimeout, IdleTimeout: requestTimeout, ErrorLog: log.New(io.Discard, "", 0),
		BaseContext: func(net.Listener) context.Context { return base },
	}
	return &Server{listener: listener, http: httpServer, handler: h, cancel: cancel}, nil
}

func (s *Server) ConsoleURL() string {
	return "http://" + s.handler.host + "/console/"
}

func (s *Server) BootstrapURL() string {
	return "http://" + s.handler.host + "/console/#token=" + s.handler.bootstrapToken()
}

func (s *Server) Serve() error {
	err := s.http.Serve(s.listener)
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

func (s *Server) Shutdown() error {
	s.handler.revoke()
	s.cancel()
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	err := s.http.Shutdown(ctx)
	if err != nil {
		_ = s.http.Close()
	}
	_ = s.listener.Close()
	return err
}

func (s *Server) Close() error {
	s.handler.revoke()
	s.cancel()
	err := s.http.Close()
	_ = s.listener.Close()
	return err
}
