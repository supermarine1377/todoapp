// package db はデータベースに対する操作を実装する
package db

import (
	"context"
	"errors"

	"github.com/supermarine1377/todoapp/app/common/apperrors"
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
