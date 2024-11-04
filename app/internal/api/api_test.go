// package api_test は、APIレベルでのテストを実行する
package api_test

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/cenkalti/backoff/v4"
	"github.com/stretchr/testify/assert"
	"github.com/supermarine1377/todoapp/app/internal/api"
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
	_, _ = io.Copy(&buff, f)
	_, err := db.Exec(buff.String())
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

func TestServer_Run(t *testing.T) {
	defer func() {
		_ = os.Remove(testDSN)
	}()
	db, err := prepareDB()
	defer func() {
		_ = db.Close()
	}()
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name       string
		prepareReq func() (*http.Request, error)
		statusCode int
	}{
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
		})
	}

	cancel()

	// Server.Run()の戻り値を検証
	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}

}
