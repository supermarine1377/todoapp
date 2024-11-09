// package db はデータベースに対する操作を実装する
package db

import (
	"context"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB はデータベースを表す
type DB struct {
	g *gorm.DB
}

// Config はデータベースを抽象化する
type Config interface {
	DSN() string
}

// NewDB はDBを生成する
func NewDB(config Config) (*DB, error) {
	g, err := gorm.Open(sqlite.Open(config.DSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &DB{g: g}, nil
}

// InsertCtx は、トランザクション内でデータを挿入する
func (db *DB) InsertCtx(ctx context.Context, p any) error {
	if err := db.g.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Create(p).Error
	}); err != nil {
		return err
	}
	return nil
}
