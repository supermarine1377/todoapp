// package config は、サーバーの設定を提供する
// 設定を操作することはできない
// 設定を読み込む機能はload packageで実装する
package config

import config_loader "github.com/supermarine1377/todoapp/app/internal/config/loader"

// Config は、サーバーの設定を表す
type Config struct {
	config config_loader.Config
}

// ConfigLoaderFunc はサーバーの設定を読み込む関数
type ConfigLoaderFunc func() (*config_loader.Config, error)

// New はサーバーの設定を返却する
func New(clf ConfigLoaderFunc) (*Config, error) {
	config, err := clf()
	if err != nil {
		return nil, err
	}
	return &Config{
		config: *config,
	}, nil
}

// Port はサーバーが稼働するポートを返す
func (c *Config) Port() int {
	return c.config.Port
}

// DSN はデータベースのDSNを返す
func (c *Config) DSN() string {
	return c.config.DB.DSN
}
