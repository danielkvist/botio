// Package server provides simple utilities to create a new http.Server.
package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
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

// WithBasicAuth receives an user, a password and a http.Handler
// and returns an Option that applies it with a basic Auth
// to the *http.Server.
func WithBasicAuth(user string, password string, h http.Handler) Option {
	return func(s *http.Server) {
		s.Handler = basicAuth(user, password, h)
	}
}

func basicAuth(username string, password string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, _ := r.BasicAuth()

		if username != user || pass != password {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	})
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
		return nil, fmt.Errorf("server.Handler equals to nil")
	case s.Addr == "":
		return nil, fmt.Errorf("server.Addr is empty")
	default:
		return s, nil
	}
}
