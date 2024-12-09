// package db はデータベースに対する操作を実装する
package db

import (
	"context"
	"errors"

	"github.com/supermarine1377/todoapp/app/common/apperrors"
	"github.com/supermarine1377/todoapp/app/internal/db/postgres"
	"github.com/supermarine1377/todoapp/app/internal/db/sqlite"
	"gorm.io/gorm"
)

// DB はデータベースを表す
type DB struct {
	g *gorm.DB
}

// DBConfig はデータベースの設定を抽象化する
type Config interface {
	DSN() string
	Type() string
}

// NewDB はDBを生成する
func NewDB(config Config) (*DB, error) {
	dbType := config.Type()
	var d gorm.Dialector

	switch dbType {
	case "sqlite":
		sqlite, err := sqlite.New(config)
		if err != nil {
			return nil, err
		}
		d = sqlite
	case "postgres":
		d = postgres.New(config)
	}
	g, err := gorm.Open(d)
	if err != nil {
		return nil, err
	}
	db, err := g.DB()
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &DB{g: g}, nil
}

// InsertCtx は、トランザクション内でデータを挿入する
// p は任意の型のポインタでなければならない
func (db *DB) InsertCtx(ctx context.Context, p any) error {
	if err := db.g.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Create(p).Error
	}); err != nil {
		return err
	}
	return nil
}

// SelectListCtx は、データの一覧を返す
// p は任意の型のポインタでなければならない
func (db *DB) SelectListCtx(ctx context.Context, p any, columns []string, offset, limit int) error {
	return db.g.WithContext(ctx).
		Select(columns).
		Offset(offset).
		Limit(limit).
		Find(p).
		Error
}

// SelectWithIDCtx は、与えたidのデータを返す
// p は任意の型のポインタでなければならない
func (db *DB) SelectWithIDCtx(ctx context.Context, p any, columns []string, id int) error {
	err := db.g.WithContext(ctx).
		Select(columns).
		First(p, id).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperrors.ErrNotFound
	}
	return err
}
