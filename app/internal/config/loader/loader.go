// package loader は、サーバーの設定を読み込む機能を提供する
package loader

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

// Config はサーバー外に設定された値を表す
type Config struct {
	// Port は、サーバーを公開するport
	Port int `env:"PORT" envDefault:"8080"`
	// DB はデータベースの設定
	DB DB
}

// DB はデータベースの設定を表す
type DB struct {
	// DSN はDBの接続情報を表す設定
	DSN string `env:"DATABASE_DSN"`
}

// Parse は、サーバーの設定をパースする
func Parse() (*Config, error) {
	var config Config
	if err := env.Parse(&config); err != nil {
		return nil, fmt.Errorf("failed to load server configuretions: %w", err)
	}
	return &config, nil
}
