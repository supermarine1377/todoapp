// package appは、必要な依存性を注入してtodoappのサーバーを起動する
package app

import (
	"context"

	"github.com/supermarine1377/todoapp/app/internal/api"
	"github.com/supermarine1377/todoapp/app/internal/config"
	"github.com/supermarine1377/todoapp/app/internal/config/loader"
)

// Run はサーバーを起動する
func Run(ctx context.Context) error {
	config, err := config.New(loader.Parse)
	if err != nil {
		return err
	}
	server := api.NewServer(config)

	if err := server.Run(ctx); err != nil {
		return err
	}

	return nil
}
