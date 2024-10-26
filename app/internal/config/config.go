// package config は、サーバーの設定を提供する
// 設定を操作することはできない
// 設定を読み込む機能はload packageで実装する
package config

import config_loader "github.com/supermarine1377/todoapp/app/internal/config/internal/loader"

// Config は、サーバーの設定を表す
type Config struct {
	config config_loader.Config
}

func New() (*Config, error) {
	config, err := config_loader.Parse()
	if err != nil {
		return nil, err
	}
	return &Config{
		config: *config,
	}, nil
}

func (c *Config) Port() int {
	return c.config.Port
}
