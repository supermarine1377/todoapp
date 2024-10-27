// package apiは、APIを実装する
package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/supermarine1377/todoapp/app/internal/api/handlers/healthz"
	"golang.org/x/sync/errgroup"
)

// Server は、APIのServerを表す
type Server struct {
	// Portは、APIを公開するポートを表す
	config Config
	e      *echo.Echo
}

// Configは、Serverの設定を抽象化する
type Config interface {
	Port() int
}

// NewServer は、Serverを作成する
func NewServer(config Config) *Server {
	e := echo.New()

	e.Logger.SetLevel(log.INFO)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	return &Server{
		config: config,
		e:      e,
	}
}

// Runは、Serverを起動する
// 以下の機能を備える
//
// 1. 定義されたHTTPリクエストを受け付ける
//
// 2. 引数のcontext.Contextを通じて処理の中断命令を検知したとき、Serverを終了する
//
// 3. 戻り値としては*echo.Echo.Start()のエラーを返却する
func (s *Server) Run(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	s.e.Add(http.MethodGet, "/healthz", healthz.Healthz)

	eg.Go(func() error {
		addr := fmt.Sprintf(":%d", s.config.Port())
		s.e.Logger.Info("Start sever")
		if err := s.e.Start(addr); err != http.ErrServerClosed {
			s.e.Logger.Error("failed to shutdown server: %w", err)
		}
		return nil
	})

	<-ctx.Done()
	s.e.Logger.Info("Shutting down server gracefully")
	if err := s.e.Shutdown(context.Background()); err != nil {
		s.e.Logger.Errorf("Failed to shutdown server: %w", err)
	}

	// Goメソッドで起動した別ゴールーチンの起動を待つ
	return eg.Wait()
}
