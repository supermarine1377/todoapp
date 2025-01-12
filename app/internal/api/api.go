// package apiは、APIを実装する
package api

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/supermarine1377/todoapp/app/common/logger"
	"github.com/supermarine1377/todoapp/app/internal/api/server"
	"github.com/supermarine1377/todoapp/app/internal/db"
	"golang.org/x/sync/errgroup"
)

// Configは、APIの設定を抽象化する
type Config interface {
	Port() int
	DSN() string
}

// HTTPServe は、提供された設定とコンテキストでHTTPサーバーを起動します。
// ヘルスチェックやタスク管理エンドポイントのためのハンドラを登録します。
// サーバーは別のGoroutineで実行され、contextが完了するのを待って
// サーバーをgracefully shutdownします。
//
// パラメータ:
//   - ctx: サーバーのライフサイクルを制御するためのコンテキスト。
//   - config: ポートやDSNを含むサーバーの設定。
//
// 戻り値:
//   - error: サーバーの起動またはシャットダウンに失敗した場合のエラー。
func HTTPServe(ctx context.Context, config Config) error {
	s, err := server.New(
		"localhost",
		server.WithPort(config.Port()),
	)
	if err != nil {
		return err
	}

	db, err := db.NewDB(config.DSN())
	if err != nil {
		return err
	}

	Routes(s, db)

	logger := slog.New(logger.NewHandler())
	slog.SetDefault(logger)

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		slog.Info("start server")
		if err := s.Run(ctx); err != nil {
			slog.Error("failed to start server", "err", err)
			return fmt.Errorf("failed to run server: %w", err)
		}
		return nil
	})

	<-ctx.Done()
	slog.Info("shutting down server gracefully")
	if err := s.Shutdown(context.Background()); err != nil {
		slog.Error("failed to shutdown server", "err", err)
	}

	// Goメソッドで起動した別ゴールーチンの起動を待つ
	return eg.Wait()
}
