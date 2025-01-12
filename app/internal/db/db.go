// package db はデータベースに対する操作を実装する
package db

import (
	"context"
	"errors"

	"github.com/supermarine1377/todoapp/app/common/apperrors"
	"github.com/supermarine1377/todoapp/app/internal/db/sqlite"
	"gorm.io/gorm"
)

// DB はデータベースを表す
type DB struct {
	g *gorm.DB
}

// Option はデータベースの設定を表す
type Option func(options *options) error

type options struct {
	kind *Kind
}

type Kind int
const (
	SQLite Kind = iota
)

var ErrInvalidDSN = errors.New("invalid DSN")

func WithDBKind(kind Kind) Option {
	return func(options *options) error {
		options.kind = &kind
		return nil
	}
}

// NewDB はDBを生成する
func NewDB(dsn string, opts ...Option) (*DB, error) {
	if dsn == "" {
		return nil, ErrInvalidDSN
	}

	var options options
	for _, opt := range opts {
		err := opt(&options)
		if err != nil {
			return nil, err
		}
	}
	if options.kind == nil {
		var defaultKind Kind = SQLite
		options.kind = &defaultKind
	}

	var dialector gorm.Dialector

	switch *options.kind {
	case SQLite:
		d, err := sqlite.New(dsn)
		if err != nil {
			return nil, err
		}
		dialector = d
	default:
		d, err := sqlite.New(dsn)
		if err != nil {
			return nil, err
		}
		dialector = d
	}

	g, err := gorm.Open(dialector)
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
