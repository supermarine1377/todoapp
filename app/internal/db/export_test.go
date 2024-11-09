package db

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/mattn/go-sqlite3"
	"github.com/supermarine1377/todoapp/app/internal/model/entity/task"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newMockDB() (*DB, sqlmock.Sqlmock, error) {
	conn, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	mock.ExpectQuery("select sqlite_version()").WillReturnRows(
		sqlmock.NewRows([]string{"sqlite_version()"}).AddRow("3.32.3"),
	)

	g, err := gorm.Open(sqlite.Dialector{Conn: conn}, &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}
	return &DB{g: g}, mock, nil
}

var ErrDummy = errors.New("dummy error")

type anyTime struct{}

func TestDB_InsertCtx(t *testing.T) {

	tests := []struct {
		name      string
		prepareDB func() (*DB, error)
		p         any
		wantErr   bool
	}{
		{
			name: "Failed to begin transaction",
			prepareDB: func() (*DB, error) {
				db, mock, err := newMockDB()
				if err != nil {
					return nil, err
				}
				mock.ExpectBegin().WillReturnError(ErrDummy)
				return db, nil
			},
			wantErr: true,
		},
		{
			name: "Failed to commit transaction",
			prepareDB: func() (*DB, error) {
				db, mock, err := newMockDB()
				if err != nil {
					return nil, err
				}
				mock.ExpectBegin().WillReturnError(nil)
				mock.ExpectCommit().WillReturnError(ErrDummy)
				return db, nil
			},
			p:       struct{}{},
			wantErr: true,
		},
		{
			name: "Failed to rollback transaction",
			prepareDB: func() (*DB, error) {
				db, mock, err := newMockDB()
				if err != nil {
					return nil, err
				}
				mock.ExpectBegin().WillReturnError(nil)
				mock.ExpectRollback().WillReturnError(ErrDummy)
				mock.ExpectCommit().WillReturnError(nil)
				return db, nil
			},
			p:       struct{}{},
			wantErr: true,
		},
		{
			name: "Insert successful",
			prepareDB: func() (*DB, error) {
				db, mock, err := newMockDB()
				if err != nil {
					return nil, err
				}
				mock.ExpectBegin().WillReturnError(nil)
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `tasks` (`title`,`created_at`,`updated_at`) VALUES (?,?,?)")).
					WithArgs("dummy", sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit().WillReturnError(nil)

				return db, nil
			},
			p:       &task.Task{Title: "dummy"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := tt.prepareDB()
			if err != nil {
				t.Fatal(err)
			}

			if err := db.InsertCtx(context.Background(), tt.p); (err != nil) != tt.wantErr {
				t.Errorf("DB.InsertCtx() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
