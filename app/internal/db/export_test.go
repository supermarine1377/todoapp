package db

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	"github.com/supermarine1377/todoapp/app/common/apperrors"
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

func TestDB_ListCtx(t *testing.T) {
	type args struct {
		columns []string
		p       any
		offset  int
		limit   int
	}
	tests := []struct {
		name      string
		args      args
		prepareDB func() (*DB, error)
		wantErr   bool
	}{
		{
			name: "Select successful",
			args: args{
				columns: []string{"title", "created_at", "updated_at"},
				p:       &task.Tasks{},
				offset:  5,
				limit:   10,
			},
			prepareDB: func() (*DB, error) {
				db, mock, err := newMockDB()
				if err != nil {
					return nil, err
				}
				mock.ExpectQuery("SELECT `title`,`created_at`,`updated_at` FROM `tasks` LIMIT 10 OFFSET 5").
					WillReturnRows(sqlmock.NewRows([]string{"title", "created_at", "updated_at"}).
						AddRow("hoge", 1, 1))
				return db, nil
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := tt.prepareDB()
			if err != nil {
				t.Fatal(err)
			}

			if err := db.SelectListCtx(
				context.Background(), tt.args.p, tt.args.columns, tt.args.offset, tt.args.limit,
			); (err != nil) != tt.wantErr {
				t.Errorf("DB.ListCtx() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDB_SelectWithIDCtx(t *testing.T) {
	type args struct {
		p       any
		columns []string
		id      int
	}
	tests := []struct {
		name      string
		args      args
		prepareDB func() (*DB, error)
		wantErr   bool
		error     error
	}{
		{
			name: "Select successful",
			args: args{
				p:       &task.Task{},
				columns: []string{"id", "title", "created_at", "updated_at"},
				id:      1,
			},
			prepareDB: func() (*DB, error) {
				db, mock, err := newMockDB()
				if err != nil {
					return nil, err
				}
				sql := "SELECT `id`,`title`,`created_at`,`updated_at` FROM `tasks` WHERE `tasks`.`id` = ? ORDER BY `tasks`.`id` LIMIT 1"
				mock.ExpectQuery(regexp.QuoteMeta(sql)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "created_at", "updated_at"}).
						AddRow(1, "hoge", 0, 0),
					)
				return db, nil
			},
			wantErr: false,
		},
		{
			name: "Select successful but got no rows",
			args: args{
				p:       &task.Task{},
				columns: []string{"id", "title", "created_at", "updated_at"},
				id:      1,
			},
			prepareDB: func() (*DB, error) {
				db, mock, err := newMockDB()
				if err != nil {
					return nil, err
				}
				sql := "SELECT `id`,`title`,`created_at`,`updated_at` FROM `tasks` WHERE `tasks`.`id` = ? ORDER BY `tasks`.`id` LIMIT 1"
				mock.ExpectQuery(regexp.QuoteMeta(sql)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "created_at", "updated_at"}))
				return db, nil
			},
			wantErr: true,
			error:   apperrors.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := tt.prepareDB()
			if err != nil {
				t.Fatal(err)
			}
			err = db.SelectWithIDCtx(context.Background(), tt.args.p, tt.args.columns, tt.args.id)
			if !tt.wantErr {
				require.NoError(t, err)
			}
			if tt.wantErr {
				require.Error(t, err, tt.error)
			}
		})
	}
}
