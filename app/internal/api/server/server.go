// package server はサーバーを提供する
package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Server はサーバーを表す
type Server struct {
	e    *echo.Echo
	addr string
	port int
}

// Option はサーバーの設定を表す
type Option func(options *options) error

type options struct {
	port *int
}

// ErrInvalidPort はportが不正だった時のエラー
var ErrInvalidPort = errors.New("invalid port")

// WithPort はportを設定する
func WithPort(port int) Option {
	return func(options *options) error {
		if port <= 0 {
			return ErrInvalidPort
		}
		options.port = &port
		return nil
	}
}

const defaultPort = 8080

// New はServerを生成する
func New(addr string, opts ...Option) (*Server, error) {
	var options options
	for _, opt := range opts {
		err := opt(&options)
		if err != nil {
			return nil, err
		}
	}

	var port int
	if options.port == nil {
		port = defaultPort
	} else {
		port = *options.port
	}

	e := echo.New()

	return &Server{
		e:    e,
		addr: addr,
		port: port,
	}, nil
}

// Run はサーバーを起動する
func (s *Server) Run(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%d", s.addr, s.port)
	if err := s.e.Start(addr); err != http.ErrServerClosed {
		return err
	}
	return nil
}

// RegisterHandler はAPIハンドラーを登録する
func (s *Server) RegisterHandler(handler echo.HandlerFunc, path string, methods ...string) {
	for _, method := range methods {
		s.e.Add(method, path, handler)
	}
}

// Use はミドルウェアを登録する
func (s *Server) Use(middleware ...echo.MiddlewareFunc) {
	s.e.Use(middleware...)
}

// Shutdown はサーバーを終了する
func (s *Server) Shutdown(ctx context.Context) error {
	return s.e.Shutdown(ctx)
}
