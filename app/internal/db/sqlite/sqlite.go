// package sqlite は、データベースとしてSQLiteを使う場合の機能を提供する
package sqlite

import (
	"errors"
	"fmt"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Config interface {
	DSN() string
}

// ErrSQLiteFileNotFound はSQLiteのファイルが見つからなかったエラー
var ErrSQLiteFileNotFound = errors.New("error: SQLite file not found")

// ErrPathIsDirectory は指定のパスがファイルでなくディレクトリを指しているときのエラー
var ErrPathIsDirectory = errors.New("error: path points to a directory, not a file")

// ErrFileLacksPermissions はDSNファイルの読み込み/書き込み権限がユーザーにないことを表す
var ErrFileLacksPermissions = errors.New("error: file lacks read/write permissions for user")

// New は、SQLiteのgorm.Dialectorを返す
func New(config Config) (gorm.Dialector, error) {
	dsn := config.DSN()
	fi, err := os.Stat(dsn)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%s: %w", dsn, ErrSQLiteFileNotFound)
		}
	}
	if fi.IsDir() {
		return nil, fmt.Errorf("%s %w", dsn, ErrPathIsDirectory)
	}

	file, err := os.OpenFile(dsn, os.O_RDWR, 0)
	defer func() {
		_ = file.Close()
	}()

	if err != nil {
		return nil, fmt.Errorf("%s: %w: %w", dsn, ErrFileLacksPermissions, err)
	}

	return sqlite.Open(dsn), nil
}
