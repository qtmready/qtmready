package server

import (
	"context"
	"crypto/tls"
	"log/slog"
	"net"
	"net/http"
	"time"
)

type (
	Server struct {
		config *Config        // server configuration
		mux    *http.ServeMux // request multiplexer
		self   *http.Server   // server instance
	}

	Option func(*Server)
)

func (s *Server) add(path string, handler http.Handler) {
	if s.mux == nil {
		s.mux = http.NewServeMux()
	}

	s.mux.Handle(path, handler)
}

func (s *Server) Start(ctx context.Context) error {
	if s.config == nil {
		slog.Warn("nomad: no configuration provider, using default configuration")

		s.config = &DefaultConfig
	}

	if s.mux == nil {
		s.mux = http.NewServeMux()
	}

	s.self = &http.Server{
		Addr:                         s.config.Address(),
		Handler:                      s.mux,
		DisableGeneralOptionsHandler: false,
		ReadHeaderTimeout:            time.Second * 30,
		WriteTimeout:                 time.Second * 30,
		IdleTimeout:                  time.Second * 60,
		MaxHeaderBytes:               http.DefaultMaxHeaderBytes,
		TLSNextProto:                 map[string]func(*http.Server, *tls.Conn, http.Handler){},
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
		ConnContext: func(ctx context.Context, c net.Conn) context.Context {
			return ctx
		},
	}

	slog.Info("nomad: starting", "port", s.config.Port, "ssl", s.config.EnableSSL)

	return s.self.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.self.Close()
}

func WithConfig(c *Config) Option {
	return func(s *Server) {
		s.config = c
	}
}

func WithHandler(path string, handler http.Handler) Option {
	return func(s *Server) {
		if s.mux == nil {
			s.mux = http.NewServeMux()
		}

		s.mux.Handle(path, handler)
	}
}

func New(opts ...Option) *Server {
	s := &Server{}

	for _, opt := range opts {
		opt(s)
	}

	return s
}
