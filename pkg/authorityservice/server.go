package authorityservice

import (
	"context"
	"net/http"
)

// HTTPServer wraps the HTTP server lifecycle.
type HTTPServer struct {
	server  *http.Server
	tlsCert string
	tlsKey  string
}

// NewHTTPServer builds an HTTPServer from config and handler.
func NewHTTPServer(cfg *Config, handler http.Handler) *HTTPServer {
	return &HTTPServer{
		server: &http.Server{
			Addr:    cfg.ListenAddress,
			Handler: handler,
		},
		tlsCert: cfg.TLSCertFile,
		tlsKey:  cfg.TLSKeyFile,
	}
}

// Start runs the server (TLS if configured).
func (s *HTTPServer) Start() error {
	if s.tlsCert != "" && s.tlsKey != "" {
		return s.server.ListenAndServeTLS(s.tlsCert, s.tlsKey)
	}
	return s.server.ListenAndServe()
}

// Stop gracefully shuts down the server.
func (s *HTTPServer) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
