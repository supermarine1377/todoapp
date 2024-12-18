// package apiは、APIを実装する
package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/supermarine1377/todoapp/app/common/logger"
	"github.com/supermarine1377/todoapp/app/internal/api/handlers/healthz"
	"github.com/supermarine1377/todoapp/app/internal/api/handlers/task"
	"github.com/supermarine1377/todoapp/app/internal/api/my_middleware"
	"github.com/supermarine1377/todoapp/app/internal/db"
	"github.com/supermarine1377/todoapp/app/internal/repository"
	"golang.org/x/sync/errgroup"
)

// Server は、APIのServerを表す
type Server struct {
	// Portは、APIを公開するポートを表す
	config Config
	e      *echo.Echo
	db     *db.DB
}

// Configは、Serverの設定を抽象化する
type Config interface {
	Port() int
	DSN() string
}

// NewServer は、Serverを作成する
func NewServer(config Config) (*Server, error) {
	db, err := db.NewDB(config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	e := echo.New()

	logger := slog.New(logger.NewHandler())
	slog.SetDefault(logger)

	e.Use(my_middleware.Recover(logger))
	e.Use(my_middleware.RequestID())
	e.Use(my_middleware.Log(logger))

	return &Server{
		config: config,
		e:      e,
		db:     db,
	}, nil
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

	{
		s.e.Add(http.MethodGet, "/healthz", healthz.Healthz)
	}
	{
		tr := repository.NewTaskRepository(s.db)
		th := task.NewTaskHandler(tr)
		s.e.Add(http.MethodPost, "/tasks", th.Create)
		s.e.Add(http.MethodGet, "/tasks", th.List)
		s.e.Add(http.MethodGet, "/tasks/:id", th.Get)
	}

	eg.Go(func() error {
		addr := fmt.Sprintf(":%d", s.config.Port())
		slog.Info("start server")
		if err := s.e.Start(addr); err != http.ErrServerClosed {
			slog.Error("failed to start server", "err", err)
			s.e.Logger.Errorf("failed tto start server: %w", err)
		}
		return nil
	})

	<-ctx.Done()
	s.e.Logger.Info("Shutting down server gracefully")
	if err := s.e.Shutdown(context.Background()); err != nil {
		slog.Error("failed to shutdown server", "err", err)
	}

	// Goメソッドで起動した別ゴールーチンの起動を待つ
	return eg.Wait()
}
