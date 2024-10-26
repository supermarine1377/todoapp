package api_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/supermarine1377/todoapp/app/internal/api"
	"golang.org/x/sync/errgroup"
)

type MockConfig struct{}

func (mc MockConfig) Port() int {
	return 8080
}

func TestServer_Run(t *testing.T) {
	t.Run("Test Sever.Run()", func(t *testing.T) {
		server := api.NewServer(MockConfig{})

		ctx, cancel := context.WithCancel(context.Background())
		eg, ctx := errgroup.WithContext(ctx)
		eg.Go(func() error {
			return server.Run(ctx)
		})

		resp, err := http.Get("http://localhost:8080/healthz")
		if err != nil {
			t.Errorf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Error("got unexpected http status code")
		}

		cancel()

		// Server.Run()の戻り値を検証
		if err := eg.Wait(); err != nil {
			t.Fatal(err)
		}
	})
}
