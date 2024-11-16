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

// ErrNoSuchFileOrDir は、DSNのファイルが見つからなったエラー
var ErrNoSuchFileOrDir = errors.New("no such file or directory")

// ErrIsADirは、DSNがディレクトリだったエラー
var ErrIsADir = errors.New("is a directory, not a file")

// ErrFileNotAccessble は、DSNのファイルが十分なパーミッションがなかったエラー
var ErrFileNotAccessble = errors.New("file is not accessible with read/write permissions")

// New は、SQLiteのgorm.Dialectorを返す
func New(config Config) (gorm.Dialector, error) {
	dsn := config.DSN()
	fi, err := os.Stat(dsn)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%s: %w", dsn, ErrNoSuchFileOrDir)
		}
	}
	if fi.IsDir() {
		return nil, fmt.Errorf("%s %w", dsn, ErrIsADir)
	}

	file, err := os.OpenFile(dsn, os.O_RDWR, 0)
	defer func() {
		_ = file.Close()
	}()

	if err != nil {
		return nil, fmt.Errorf("%s: %w: %w", dsn, ErrFileNotAccessble, err)
	}

	return sqlite.Open(dsn), nil
}
