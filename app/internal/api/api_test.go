// package api_test は、APIレベルでのテストを実行する
package api_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/cenkalti/backoff/v4"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/supermarine1377/todoapp/app/internal/api"
	"github.com/supermarine1377/todoapp/app/internal/model/entity/task"
	"golang.org/x/sync/errgroup"
)

const testDSN = "test.db"

func prepareDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", testDSN)
	if err != nil {
		return nil, err
	}
	f, err := os.Open("../../../_migration/sqlite/schema.sql")
	if err != nil {
		return nil, err
	}
	if err := execSQL(db, f); err != nil {
		return nil, err
	}
	f, err = os.Open("_test/test_data.sql")
	if err != nil {
		return nil, err
	}
	if err := execSQL(db, f); err != nil {
		return nil, err
	}
	return db, nil
}

func execSQL(db *sql.DB, f *os.File) error {
	defer f.Close()
	var buff bytes.Buffer
	_, err := io.Copy(&buff, f)
	if err != nil {
		return err
	}
	_, err = db.Exec(buff.String())
	if err != nil {
		return err
	}
	return nil
}

type MockConfig struct{}

func (mc MockConfig) Port() int {
	return 8080
}

func (mc MockConfig) DSN() string {
	return testDSN
}

func jsonUnmarshal[T any](res io.Reader) (any, error) {
	var buff bytes.Buffer
	_, err := io.Copy(&buff, res)
	if err != nil {
		return nil, err
	}
	var t T
	if err := json.Unmarshal(buff.Bytes(), &t); err != nil {
		return nil, err
	}
	return t, nil
}

var tasksInDB = task.Tasks{
	{
		ID:        1,
		Title:     "Write report",
		CreatedAt: 1,
		UpdatedAt: 1,
	},
	{
		ID:        2,
		Title:     "Review code",
		CreatedAt: 2,
		UpdatedAt: 2,
	},
	{
		ID:        3,
		Title:     "Plan meeting",
		CreatedAt: 2,
		UpdatedAt: 2,
	},
	{
		ID:        4,
		Title:     "Update documentation",
		CreatedAt: 2,
		UpdatedAt: 2,
	},
	{
		ID:        5,
		Title:     "Fix bug #42",
		CreatedAt: 2,
		UpdatedAt: 3,
	},
}

var newTask = task.Task{
	ID:    6,
	Title: "hoge",
}

var tasksAfterInsert = append(tasksInDB, &newTask)

func TestServer_Run(t *testing.T) {
	t.Cleanup(func() {
		_ = os.Remove(testDSN)
	})
	db, err := prepareDB()
	defer func() {
		_ = db.Close()
	}()
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name        string
		prepareReq  func() (*http.Request, error)
		statusCode  int
		testResFunc func(t *testing.T, res io.Reader) error
	}{
		{
			name: "GET tasks",
			prepareReq: func() (*http.Request, error) {
				return http.NewRequest(http.MethodGet, "http://localhost:8080/tasks", nil)
			},
			statusCode: http.StatusOK,
			testResFunc: func(t *testing.T, res io.Reader) error {
				tasks, err := jsonUnmarshal[task.Tasks](res)
				if err != nil {
					return err
				}
				assert.Equal(t, tasksInDB, tasks)
				return nil
			},
		},
		{
			name: "POST task (Invalid request)",
			prepareReq: func() (*http.Request, error) {
				return http.NewRequest(http.MethodPost, "http://localhost:8080/tasks", nil)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "POST task (appropiate request)",
			prepareReq: func() (*http.Request, error) {
				b := strings.NewReader(`{"title": "hoge"}`)
				req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/tasks", b)
				if err != nil {
					return nil, err
				}
				req.Header.Add("Content-Type", "application/json")
				return req, nil
			},
			statusCode: http.StatusCreated,
		},
		{
			name: "GET tasks after task creation",
			prepareReq: func() (*http.Request, error) {
				return http.NewRequest(http.MethodGet, "http://localhost:8080/tasks", nil)
			},
			statusCode: http.StatusOK,
			testResFunc: func(t *testing.T, res io.Reader) error {
				tasks, err := jsonUnmarshal[task.Tasks](res)
				if err != nil {
					return err
				}
				if diff := cmp.Diff(
					tasks,
					tasksAfterInsert,
					cmpopts.IgnoreFields(task.Task{}, "CreatedAt", "UpdatedAt")); diff != "" {
					t.Error(diff)
				}
				return nil
			},
		},
		{
			name: "GET task id =6",
			prepareReq: func() (*http.Request, error) {
				return http.NewRequest(http.MethodGet, "http://localhost:8080/tasks/6", nil)
			},
			statusCode: http.StatusOK,
			testResFunc: func(t *testing.T, res io.Reader) error {
				ta, err := jsonUnmarshal[task.Task](res)
				if err != nil {
					return err
				}
				if diff := cmp.Diff(
					ta,
					newTask,
					cmpopts.IgnoreFields(task.Task{}, "CreatedAt", "UpdatedAt")); diff != "" {
					t.Error(diff)
				}
				return nil
			},
		},
	}
	server := api.NewServer(MockConfig{})

	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return server.Run(ctx)
	})

	client := &http.Client{}
	if err := backoff.Retry(func() error {
		req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/healthz", nil)
		if err != nil {
			t.Fatal(err)
		}
		res, err := client.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			return errors.New("unexpected status code")
		}
		return nil
	}, backoff.NewExponentialBackOff()); err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := tt.prepareReq()
			if err != nil {
				t.Fatal(err)
			}
			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			assert.Equal(t, tt.statusCode, res.StatusCode)
			if tt.testResFunc != nil {
				err := tt.testResFunc(t, res.Body)
				if err != nil {
					t.Fatal(err)
				}
			}
		})
	}

	cancel()

	// Server.Run()の戻り値を検証
	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}

}
