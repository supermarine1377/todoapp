// package api_test は、APIレベルでのテストを実行する
package api_test

import (
	"context"
	"net/http"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supermarine1377/todoapp/app/internal/api"
	"golang.org/x/sync/errgroup"
)

type MockConfig struct{}

func (mc MockConfig) Port() int {
	return 8080
}

func TestServer_Run(t *testing.T) {
	tests := []struct {
		name       string
		prepareReq func() (*http.Request, error)
		statusCode int
	}{
		{
			name: "GET /healthz",
			prepareReq: func() (*http.Request, error) {
				return http.NewRequest(http.MethodGet, "http://localhost:8080/healthz", nil)
			},
			statusCode: http.StatusOK,
		},
		{
			name: "POST /task",
			prepareReq: func() (*http.Request, error) {
				return http.NewRequest(http.MethodPost, "http://localhost:8080/task", nil)
			},
			statusCode: http.StatusCreated,
		},
	}
	server := api.NewServer(MockConfig{})

	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)

	var wg sync.WaitGroup
	wg.Add(1)
	eg.Go(func() error {
		wg.Done()
		return server.Run(ctx)
	})
	// Wait until server's ready
	wg.Wait()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := tt.prepareReq()
			if err != nil {
				t.Fatal(err)
			}
			client := &http.Client{}
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
