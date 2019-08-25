// Package server provides simple utilities to create a new http.Server.
package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/danielkvist/botio/jwt"
)

// Option defines an option for a new *http.Server.
type Option func(s *http.Server)

// WithHandler receives an http.Handler and returns an Option
// that applies it to the *http.Server.
func WithHandler(h http.Handler) Option {
	return func(s *http.Server) {
		s.Handler = h
	}
}

// WithJWTAuth receives a signing key and returns an Option that applies JWT
// tokens authentication to the *http.Server usin ghe *http.Server's Handler.
func WithJWTAuth(key string) Option {
	return func(s *http.Server) {
		s.Handler = jwt.Middleware(key, s.Handler)
	}
}

// WithListenAddr receives an address and returns a Option
// that applies it to the *http.Server's Addr field.
func WithListenAddr(addr string) Option {
	return func(s *http.Server) {
		s.Addr = addr
	}
}

// WithReadTimeout receives a timeout and returns an Option
// that applies it to the *http.Server's ReadTiemout field.
func WithReadTimeout(t time.Duration) Option {
	return func(s *http.Server) {
		s.ReadTimeout = t
	}
}

// WithWriteTimeout receives a timeout and returns an Option
// that assigns it to the *http.Server's WriteTiemout field.
func WithWriteTimeout(t time.Duration) Option {
	return func(s *http.Server) {
		s.WriteTimeout = t
	}
}

// WithIdleTimeout receives a timeout and returns an Option
// that assigns it to the *http.Server's IdleTiemout field.
func WithIdleTimeout(t time.Duration) Option {
	return func(s *http.Server) {
		s.IdleTimeout = t
	}
}

// WithGracefulShutdown receives a done and a quit channel and returns
// an Option that initializes a goroutine listening the done channel
// to try to shutdown the *http.Server gracefully.
func WithGracefulShutdown(done <-chan struct{}, quit chan<- struct{}) Option {
	return func(s *http.Server) {
		go gracefulShutdown(s, done, quit)
	}
}

func gracefulShutdown(s *http.Server, done <-chan struct{}, quit chan<- struct{}) {
	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	s.SetKeepAlivesEnabled(false)
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("while trying to shutdown server gracefully: %v", err)
	}

	quit <- struct{}{}
}

// WithTLS returns an Option that changes the TLS configuration of the *http.Server.
func WithTLS() Option {
	tlsConf := &tls.Config{
		PreferServerCipherSuites: true,
		MinVersion:               tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}

	return func(s *http.Server) {
		s.TLSConfig = tlsConf
	}
}

// New can receive zero or multiple Option and returns a new *http.Server
// or an error if before return:
//
// * There is no server.Handler specified.
// * There is no server.Addr defined.
func New(options ...Option) (*http.Server, error) {
	s := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	for _, opt := range options {
		opt(s)
	}

	switch {
	case s.Handler == nil:
		return nil, fmt.Errorf("server.Handler cannot be nil")
	case s.Addr == "":
		return nil, fmt.Errorf("server.Addr cannot be an empty string")
	default:
		return s, nil
	}
}
